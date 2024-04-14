package trace

import (
	"analyzer/analysis"
	"analyzer/clock"
	"analyzer/logging"
	"errors"
	"sort"
	"strconv"
)

var (
	traces map[int][]TraceElement = make(map[int][]TraceElement)

	// current happens before vector clocks
	currentVCHb = make(map[int]clock.VectorClock)

	// current must happens before vector clocks
	currentVCWmhb = make(map[int]clock.VectorClock)

	// channel without partner
	channelWithoutPartner = make(map[int]map[int]*TraceElementChannel) // id -> opId -> element

	currentIndex     = make(map[int]int)
	numberOfRoutines = 0
	fifo             bool
	result           string

	analysisCases map[string]bool
)

/*
* Add an element to the trace
* Args:
*   routine (int): The routine id
*   element (TraceElement): The element to add
* Returns:
*   error: An error if the routine does not exist
 */
func AddElementToTrace(element TraceElement) error {
	routine := element.GetRoutine()
	traces[routine] = append(traces[routine], element)
	return nil
}

/*
* Add an empty routine to the trace
* Args:
*   routine (int): The routine id
 */
func AddEmptyRoutine(routine int) {
	traces[routine] = make([]TraceElement, 0)
}

/*
 * Sort the trace by tSort
 */
type sortByTSort []TraceElement

func (a sortByTSort) Len() int      { return len(a) }
func (a sortByTSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortByTSort) Less(i, j int) bool {
	return a[i].GetTSort() < a[j].GetTSort()
}

/*
 * Sort a trace by tpost
 * Args:
 *   trace ([]traceElement): The trace to sort
 * Returns:
 *   ([]traceElement): The sorted trace
 */
func sortTrace(trace []TraceElement) []TraceElement {
	sort.Sort(sortByTSort(trace))
	return trace
}

/*
 * Sort all traces by tpost
 */
func Sort() {
	for routine, trace := range traces {
		traces[routine] = sortTrace(trace)
	}
}

/*
 * Get the traces
 * Returns:
 *   map[int][]traceElement: The traces
 */
func GetTraces() *map[int][]TraceElement {
	return &traces
}

/*
 * Get the trace of the given routine
 * Args:
 *   id (int): The id of the routine
 * Returns:
 *   []traceElement: The trace of the routine
 */
func GetTraceFromId(id int) []TraceElement {
	return traces[id]
}

/*
 * Given the file and line info, return the routine and index of the element
 * in trace.
 * Args:
 *   tID (string): The tID of the element
 * Returns:
 *   error: An error if the element does not exist
 *   int: The routine of the element
 *   int: The index of the element in the trace of the routine
 */
func GetTraceElementFromTID(tID string) (*TraceElement, error) {
	if tID == "" {
		return nil, errors.New("tID is empty")
	}

	for routine, trace := range traces {
		for index, elem := range trace {
			if elem.GetTID() == tID {
				return &traces[routine][index], nil
			}
		}
	}
	return nil, errors.New("Element " + tID + " does not exist")
}

/*
 * Shorten the trace by removing all elements after the given time
 * Args:
 *   time (int): The time to shorten the trace to
 *   incl (bool): True if an element with the same time should stay included in the trace
 */
func ShortenTrace(time int, incl bool) {
	for routine, trace := range traces {
		for index, elem := range trace {
			if incl && elem.GetTSort() > time {
				traces[routine] = traces[routine][:index]
				break
			}
			if !incl && elem.GetTSort() >= time {
				traces[routine] = traces[routine][:index]
				break
			}
		}
	}
}

/*
 * Remove the element with the given tID from the trace
 * Args:
 *   tID (string): The tID of the element to remove
 */
func RemoveElementFromTrace(tID string) {
	for routine, trace := range traces {
		for index, elem := range trace {
			if elem.GetTID() == tID {
				traces[routine] = append(traces[routine][:index], traces[routine][index+1:]...)
				break
			}
		}
	}
}

/*
 * Shorten the trace of the given routine by removing all elements after and equal the given time
 * Args:
 *   routine (int): The routine to shorten
 *   time (int): The time to shorten the trace to
 */
func ShortenRoutine(routine int, time int) {
	for index, elem := range traces[routine] {
		if elem.GetTSort() >= time {
			traces[routine] = traces[routine][:index]
			break
		}
	}
}

func ShortenRoutineIndex(routine int, index int, incl bool) {
	if incl {
		traces[routine] = traces[routine][:index+1]
	} else {
		traces[routine] = traces[routine][:index]
	}
}

/*
 * Switch the timer of two elements
 * Args:
 *   element1 (traceElement): The first element
 *   element2 (traceElement): The second element
 */
func SwitchTimer(element1 *TraceElement, element2 *TraceElement) {
	routine1 := (*element1).GetRoutine()
	routine2 := (*element2).GetRoutine()
	tSort1 := (*element1).GetTSort()
	for index, elem := range traces[routine1] {
		if elem.GetTSort() == (*element1).GetTSort() {
			traces[routine1][index].SetTSort((*element2).GetTSort())
		}
	}
	for index, elem := range traces[routine2] {
		if elem.GetTSort() == (*element2).GetTSort() {
			traces[routine2][index].SetTSort(tSort1)
			break
		}
	}

}

/*
 * Move the time of elements back by steps, excluding the routines in
 * excludedRoutines
 * Args:
 *   startTime (int): The time to start moving back from
 *   steps (int): The number of steps to move back
 *   excludedRoutines ([]int): The routines to exclude
 */
// func MoveTimeBack(startTime int, steps int, excludedRoutines []int) {
// 	println("Move Time Back")
// 	println("Start Time: ", startTime)
// 	println("Steps: ", steps)
// 	for routine, localTrace := range traces {
// 		for _, elem := range localTrace {
// 			if elem.GetTSort() >= startTime && !contains(excludedRoutines, routine) {
// 				elem.SetTSortWithoutNotExecuted(elem.GetTSort() + steps)
// 			}
// 		}
// 	}
// 	Sort()
// }

func contains(slice []int, elem int) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}

/*
 * Set the number of routines
 * Args:
 *   n (int): The number of routines
 */
func SetNumberOfRoutines(n int) {
	numberOfRoutines = n
}

/*
* Calculate vector clocks
* MARK: run analysis
* Args:
*   assume_fifo (bool): True to assume fifo ordering in buffered channels
*   ignoreCriticalSections (bool): True to ignore critical sections when updating
*   	vector clocks
*   analysisCasesMap (map[string]bool): The analysis cases to run
 */
func RunAnalysis(assumeFifo bool, ignoreCriticalSections bool, analysisCasesMap map[string]bool) string {

	logging.Debug("Analyze the trace...", logging.INFO)

	fifo = assumeFifo

	analysisCases = analysisCasesMap
	analysis.InitAnalysis(analysisCases)

	for i := 1; i <= numberOfRoutines; i++ {
		currentVCHb[i] = clock.NewVectorClock(numberOfRoutines)
		currentVCWmhb[i] = clock.NewVectorClock(numberOfRoutines)
	}

	currentVCHb[1] = currentVCHb[1].Inc(1)
	currentVCWmhb[1] = currentVCWmhb[1].Inc(1)

	for elem := getNextElement(); elem != nil; elem = getNextElement() {

		switch e := elem.(type) {
		case *TraceElementAtomic:
			logging.Debug("Update vector clock for atomic operation "+e.ToString()+
				" for routine "+strconv.Itoa(e.GetRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *TraceElementChannel:
			logging.Debug("Update vector clock for channel operation "+e.ToString()+
				" for routine "+strconv.Itoa(e.GetRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *TraceElementMutex:
			if ignoreCriticalSections {
				logging.Debug("Ignore critical section "+e.ToString()+
					" for routine "+strconv.Itoa(e.GetRoutine()), logging.DEBUG)
				e.updateVectorClockAlt()
			} else {
				logging.Debug("Update vector clock for mutex operation "+e.ToString()+
					" for routine "+strconv.Itoa(e.GetRoutine()), logging.DEBUG)
				e.updateVectorClock()
			}
		case *TraceElementFork:
			logging.Debug("Update vector clock for routine operation "+e.ToString()+
				" for routine "+strconv.Itoa(e.GetRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *TraceElementSelect:
			logging.Debug("Update vector clock for select operation "+e.ToString()+
				" for routine "+strconv.Itoa(e.GetRoutine()), logging.DEBUG)
			cases := e.GetCases()
			ids := make([]int, 0)
			opTypes := make([]int, 0)
			for _, c := range cases {
				switch c.opC {
				case Send:
					ids = append(ids, c.GetID())
					opTypes = append(opTypes, 0)
				case Recv:
					ids = append(ids, c.GetID())
					opTypes = append(opTypes, 1)
				}
			}
			analysis.CheckForLeakSelectRun(ids, opTypes, currentVCHb[e.routine].Copy(), e.tID)
			e.updateVectorClock()
		case *TraceElementWait:
			logging.Debug("Update vector clock for go operation "+e.ToString()+
				" for routine "+strconv.Itoa(e.GetRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *TraceElementCond:
			logging.Debug("Update vector clock for cond operation "+e.ToString()+
				" for routine "+strconv.Itoa(e.GetRoutine()), logging.DEBUG)
			e.updateVectorClock()
		}

		// check for leak
		if analysisCases["leak"] && elem.getTpost() == 0 {
			switch e := elem.(type) {
			case *TraceElementChannel:
				switch e.opC {
				case Send:
					analysis.CheckForLeakChannelStuck(elem.GetID(), currentVCHb[e.routine],
						elem.GetTID(), 0, e.qSize != 0)
				case Recv:
					analysis.CheckForLeakChannelStuck(elem.GetID(), currentVCHb[e.routine],
						elem.GetTID(), 1, e.qSize != 0)
				}
			case *TraceElementMutex:
				analysis.CheckForLeakMutex(elem.GetID(), elem.GetTID())
			case *TraceElementWait:
				analysis.CheckForLeakWait(elem.GetTID())
			case *TraceElementSelect:
				cases := e.GetCases()
				ids := make([]int, 0)
				opTypes := make([]int, 0)
				for _, c := range cases {
					switch c.opC {
					case Send:
						ids = append(ids, c.GetID())
						opTypes = append(opTypes, 0)
					case Recv:
						ids = append(ids, c.GetID())
						opTypes = append(opTypes, 1)
					}
				}
				analysis.CheckForLeakSelectStuck(ids, currentVCHb[e.routine], e.tID, opTypes, e.tPre)
			case *TraceElementCond:
				analysis.CheckForLeakCond(elem.GetTID())
			}
		}

	}

	if analysisCases["selectWithoutPartner"] {
		analysis.CheckForSelectCaseWithoutPartner()
	}

	if analysisCases["leak"] {
		analysis.CheckForLeak()
	}

	if analysisCases["doneBeforeAdd"] {
		analysis.CheckForDoneBeforeAdd()
	}

	if analysisCases["cyclicDeadlock"] {
		analysis.CheckForCyclicDeadlock()
	}

	logging.Debug("Analysis completed", logging.INFO)
	return result
}

func getNextElement() TraceElement {
	// find the local trace, where the element on which currentIndex points to
	// has the smallest tpost
	var minTSort = -1
	var minRoutine = -1
	for routine, trace := range traces {
		// no more elements in the routine trace
		if currentIndex[routine] == -1 {
			continue
		}
		// ignore non executed operations
		if trace[currentIndex[routine]].GetTSort() == 0 {
			continue
		}
		if minTSort == -1 || trace[currentIndex[routine]].GetTSort() < minTSort {
			minTSort = trace[currentIndex[routine]].GetTSort()
			minRoutine = routine
		}
	}

	// all elements have been processed
	if minRoutine == -1 {
		return nil
	}

	// return the element and increase the index
	element := traces[minRoutine][currentIndex[minRoutine]]
	increaseIndex(minRoutine)
	return element
}

func increaseIndex(routine int) {
	currentIndex[routine]++
	if currentIndex[routine] >= len(traces[routine]) {
		currentIndex[routine] = -1
	}
}

// MARK: Shift

/*
 * Shift all elements with time greater or equal to startTSort by shift
 * Only shift forward
 * Args:
 *   startTPre (int): The time to start shifting
 *   shift (int): The shift
 */
func ShiftTrace(startTPre int, shift int) bool {
	if shift <= 0 {
		return false
	}

	for routine, trace := range traces {
		for index, elem := range trace {
			if elem.GetTPre() >= startTPre {
				traces[routine][index].SetTSortWithoutNotExecuted(elem.GetTSort() + shift)
			}
		}
	}

	return true
}

/*
 * Shift all elements that are concurrent or HB-later than the element such
 * that they are after the element without changeing the order of these elements
 * Args:
 *   element (traceElement): The element
 */
func ShiftConcurrentOrAfterToAfter(element *TraceElement) {
	elemsToShift := make([]TraceElement, 0)
	minTime := -1

	println("ELEM: ", (*element).GetTID(), (*element).GetVC().ToString())

	for _, trace := range traces {
		for _, elem := range trace {
			if elem.GetTID() == (*element).GetTID() {
				continue
			}

			if !(clock.GetHappensBefore(elem.GetVC(), (*element).GetVC()) == clock.Before) {
				println("SHIFT: ", elem.GetTID(), elem.GetVC().ToString())
				elemsToShift = append(elemsToShift, elem)
				if minTime == -1 || elem.GetTPre() < minTime {
					minTime = elem.GetTPre()
				}
			}
		}
	}

	distance := (*element).GetTPre() - minTime + 1

	for _, elem := range elemsToShift {
		tSort := elem.GetTPre()
		elem.SetTSort(tSort + distance)
	}
}

/*
 * Shift the element to be after all elements, that are concurrent to it
 * Args:
 *   element (traceElement): The element
 */
func ShiftConcurrentToBefore(element *TraceElement) {
	lastConcurrentTime := (*element).GetTPre()
	elementsToShift := make([]TraceElement, 0)

	for _, trace := range traces {
		for _, elem := range trace {
			hb := clock.GetHappensBefore(elem.GetVC(), (*element).GetVC())
			if elem.GetTID() == (*element).GetTID() || hb == clock.After {
				elementsToShift = append(elementsToShift, elem)
			} else if hb == clock.Concurrent && elem.GetTPre() > lastConcurrentTime {
				lastConcurrentTime = elem.GetTPre()
			}
		}
	}

	distance := lastConcurrentTime - (*element).GetTPre()

	for _, elem := range elementsToShift {
		tSort := elem.GetTPre()
		elem.SetTSort(tSort + distance)
	}
}

/*
 * Remove all elements that are concurrent to the element
 * Args:
 *   element (traceElement): The element
 */
func RemoveConcurrent(element *TraceElement) {
	for routine, trace := range traces {
		result := make([]TraceElement, 0)
		for _, elem := range trace {
			if elem.GetTID() == (*element).GetTID() {
				result = append(result, elem)
				continue
			}

			if clock.GetHappensBefore((*element).GetVC(), elem.GetVC()) != clock.Concurrent {
				result = append(result, elem)
			}
		}
		traces[routine] = result
	}
}

/*
 * For each routine, get the earliest element that is concurrent to the element
 * Args:
 *   element (traceElement): The element
 * Returns:
 *   map[int]traceElement: The earliest concurrent element for each routine
 */
func GetConcurrentEarliest(element *TraceElement) map[int]*TraceElement {
	concurrent := make(map[int]*TraceElement)
	for routine, trace := range traces {
		for _, elem := range trace {
			if elem.GetTID() == (*element).GetTID() {
				continue
			}

			if clock.GetHappensBefore((*element).GetVC(), elem.GetVC()) == clock.Concurrent {
				concurrent[routine] = &elem
			}
		}
	}
	return concurrent
}

/*
 * Shift all elements with time greater or equal to startTSort by shift
 * Only shift back
 * Args:
 *   routine (int): The routine to shift
 *   startTSort (int): The time to start shifting
 *   shift (int): The shift
 * Returns:
 *   bool: True if the shift was successful, false otherwise (shift <= 0)
 * TODO: is this allowed or will it create problems?
 */
func ShiftRoutine(routine int, startTSort int, shift int) bool {
	if shift <= 0 {
		return false
	}

	for index, elem := range traces[routine] {
		if elem.GetTPre() >= startTSort {
			traces[routine][index].SetTSortWithoutNotExecuted(elem.GetTSort() + shift)
		}
	}

	return true
}

/*
 * Get the partial trace of all element between startTime and endTime incluseve.
 * Args:
 *  startTime (int): The start time
 *  endTime (int): The end time
 * Returns:
 *  map[int][]TraceElement: The partial trace
 */
func GetPartialTrace(startTime int, endTime int) map[int][]*TraceElement {
	result := make(map[int][]*TraceElement)
	println("\n\n")
	for routine, trace := range traces {
		for index, elem := range trace {
			if _, ok := result[routine]; !ok {
				result[routine] = make([]*TraceElement, 0)
			}
			time := elem.GetTSort()
			if time >= startTime && time <= endTime {
				result[routine] = append(result[routine], &traces[routine][index])
			}
		}
	}

	return result
}
