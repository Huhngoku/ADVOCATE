package trace

import (
	"analyzer/logging"
	vc "analyzer/vectorClock"
	"sort"
	"strconv"
)

var traces map[int][]traceElement = make(map[int][]traceElement)
var currentVectorClocks map[int]vc.VectorClock = make(map[int]vc.VectorClock)
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
type sortByTPost []traceElement

func (a sortByTPost) Len() int      { return len(a) }
func (a sortByTPost) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortByTPost) Less(i, j int) bool {
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
	sort.Sort(sortByTPost(trace))
	return trace
}

func Sort() {
	for routine, trace := range traces {
		traces[routine] = sortTrace(trace)
	}
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
* Args:
*   assume_fifo (bool): True to assume fifo ordering in buffered channels
 */
func RunAnalysis(assume_fifo bool) string {
	logging.Log("Calculate vector clocks...", logging.INFO)

	fifo = assume_fifo

	for i := 1; i <= numberOfRoutines; i++ {
		currentVectorClocks[i] = vc.NewVectorClock(numberOfRoutines)
	}

	currentVectorClocks[1] = currentVectorClocks[1].Inc(1)

	for elem := getNextElement(); elem != nil; elem = getNextElement() {
		// ignore non executed operations
		if elem.getTpost() == 0 {
			logging.Log("Skip vector clock calculation for "+elem.toString(), logging.DEBUG)
			continue
		}

		switch e := elem.(type) {
		case *traceElementAtomic:
			logging.Log("Update vector clock for atomic operation "+e.toString()+
				" for routine "+strconv.Itoa(e.getRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementChannel:
			logging.Log("Update vector clock for channel operation "+e.toString()+
				" for routine "+strconv.Itoa(e.getRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementMutex:
			logging.Log("Update vector clock for mutex operation "+e.toString()+
				" for routine "+strconv.Itoa(e.getRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementRoutine:
			logging.Log("Update vector clock for routine operation "+e.toString()+
				" for routine "+strconv.Itoa(e.getRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementSelect:
			logging.Log("Update vector clock for select operation "+e.toString()+
				" for routine "+strconv.Itoa(e.getRoutine()), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementWait:
			logging.Log("Update vector clock for go operation "+e.toString()+
				" for routine "+strconv.Itoa(e.getRoutine()), logging.DEBUG)
			e.updateVectorClock()
		}

		for i := 1; i <= numberOfRoutines; i++ {
			logging.Log(currentVectorClocks[i].ToString(), logging.DEBUG)
		}

	}

	logging.Log("Vector clock calculation completed", logging.INFO)
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
