package trace

import (
	"analyzer/analysis"
	"analyzer/logging"
	"sort"
	"strconv"
)

var traces map[int][]traceElement = make(map[int][]traceElement)
var currentVectorClocks map[int]analysis.VectorClock = make(map[int]analysis.VectorClock)
var currentIndex map[int]int = make(map[int]int)
var numberOfRoutines int = 0
var fifo bool
var result string

/*
* Add an element to the trace
* Args:
*   routine (int): The routine id
*   element (TraceElement): The element to add
* Returns:
*   error: An error if the routine does not exist
 */
func addElementToTrace(element traceElement) error {
	routine := element.getRoutine()
	traces[routine] = append(traces[routine], element)
	return nil
}

/*
 * Sort the trace by tsort
 */
type sortByTSort []traceElement

func (a sortByTSort) Len() int      { return len(a) }
func (a sortByTSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortByTSort) Less(i, j int) bool {
	return a[i].getTsort() < a[j].getTsort()
}

/*
 * Sort a trace by tpost
 * Args:
 *   trace ([]traceElement): The trace to sort
 * Returns:
 *   ([]traceElement): The sorted trace
 */
func sortTrace(trace []traceElement) []traceElement {
	sort.Sort(sortByTSort(trace))
	return trace
}

func Sort() {
	for routine, trace := range traces {
		traces[routine] = sortTrace(trace)
	}
}

/*
<<<<<<< Updated upstream
=======
 * Get the traces
 * Returns:
 *   map[int][]traceElement: The traces
 */
func GetTraces() *map[int][]TraceElement {
	return &traces
}

/*
 * Given the file and line info, return the routine and index of the element
 * in trace.
 * Args:
 *   pos (string): The position of the element
 * Returns:
 *   error: An error if the element does not exist
 *   int: The routine of the element
 *   int: The index of the element in the trace of the routine
 */
func GetTraceElementFromPos(pos string) (*TraceElement, error) {
	for routine, trace := range traces {
		for index, elem := range trace {
			if elem.GetPos() == pos {
				return &traces[routine][index], nil
			}
		}
	}
	return nil, errors.New("Element " + pos + " does not exist")
}

/*
 * Move the time of elements back by steps, excluding the routines in
 * excludedRoutines
 * Args:
 *   startTime (int): The time to start moving back from
 *   steps (int): The number of steps to move back
 *   excludedRoutines ([]int): The routines to exclude
 */
func MoveTimeBack(startTime int, steps int, excludedRoutines []int) {
	println("Move Time Back")
	println("Start Time: ", startTime)
	println("Steps: ", steps)
	for routine, localTrace := range traces {
		for _, elem := range localTrace {
			if elem.GetTSort() >= startTime && !contains(excludedRoutines, routine) {
				elem.SetTSortWithoutNotExecuted(elem.GetTSort() + steps)
			}
		}
	}
	Sort()
}

func contains(slice []int, elem int) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}

/*
>>>>>>> Stashed changes
 * Set the number of routines
 * Args:
 *   n (int): The number of routines
 */
func SetNumberOfRoutines(n int) {
	numberOfRoutines = n
}

/*
* Calculate vector clocks
* Args:
*   assume_fifo (bool): True to assume fifo ordering in buffered channels
 */
func RunAnalysis(assume_fifo bool) string {
	logging.Debug("Analyze the trace...", logging.INFO)

	fifo = assume_fifo

	for i := 1; i <= numberOfRoutines; i++ {
		currentVectorClocks[i] = analysis.NewVectorClock(numberOfRoutines)
	}

	currentVectorClocks[1] = currentVectorClocks[1].Inc(1)

	for elem := getNextElement(); elem != nil; elem = getNextElement() {
		// ignore non executed operations
		if elem.getTpost() == 0 {
			logging.Debug("Skip vector clock calculation for "+elem.toString(), logging.DEBUG)
			continue
		}

		switch e := elem.(type) {
		case *traceElementAtomic:
			logging.Debug("Update vector clock for atomic operation "+e.toString()+
				" for routine "+strconv.Itoa(e.getRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementChannel:
			logging.Debug("Update vector clock for channel operation "+e.toString()+
				" for routine "+strconv.Itoa(e.getRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementMutex:
			logging.Debug("Update vector clock for mutex operation "+e.toString()+
				" for routine "+strconv.Itoa(e.getRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementFork:
			logging.Debug("Update vector clock for routine operation "+e.toString()+
				" for routine "+strconv.Itoa(e.getRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementSelect:
			logging.Debug("Update vector clock for select operation "+e.toString()+
				" for routine "+strconv.Itoa(e.getRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementWait:
			logging.Debug("Update vector clock for go operation "+e.toString()+
				" for routine "+strconv.Itoa(e.getRoutine()), logging.DEBUG)
			e.updateVectorClock()
		}

	}

	logging.Debug("Analysis completed", logging.INFO)
	return result
}

func getNextElement() traceElement {
	// find the local trace, where the element on which currentIndex points to
	// has the smallest tpost
	var minTSort int = -1
	var minRoutine int = -1
	for routine, trace := range traces {
		// no more elements in the routine trace
		if currentIndex[routine] == -1 {
			continue
		}
		// ignore non executed operations
		if trace[currentIndex[routine]].getTsort() == 0 {
			continue
		}
		if minTSort == -1 || trace[currentIndex[routine]].getTsort() < minTSort {
			minTSort = trace[currentIndex[routine]].getTsort()
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
