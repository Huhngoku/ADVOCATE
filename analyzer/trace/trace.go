package trace

import (
	"analyzer/logging"
	vc "analyzer/vectorClock"
)

var traces map[int][]traceElement = make(map[int][]traceElement)
var currentVectorClocks map[int]vc.VectorClock = make(map[int]vc.VectorClock)
var currentIndex map[int]int = make(map[int]int)
var numberOfRoutines int = 0
var fifo bool

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
func CalculateVectorClocks(assume_fifo bool) {
	logging.Log("Calculate vector clocks...", logging.INFO)

	fifo = assume_fifo

	for i := 1; i <= numberOfRoutines; i++ {
		currentVectorClocks[i] = vc.NewVectorClock(numberOfRoutines)
	}

	for elem := getNextElement(); elem != nil; elem = getNextElement() {
		// ignore non executed operations
		if elem.getTpost() == 0 {
			logging.Log("Skip vector clock calculation for "+elem.toString(), logging.DEBUG)
			continue
		}

		switch e := elem.(type) {
		case *traceElementAtomic:
			logging.Log("Update vector clock for atomic operation "+e.toString(), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementChannel:
			logging.Log("Update vector clock for channel operation "+e.toString(), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementMutex:
			logging.Log("Update vector clock for mutex operation "+e.toString(), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementRoutine:
			logging.Log("Update vector clock for routine operation "+e.toString(), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementSelect:
			logging.Log("Update vector clock for select operation "+e.toString(), logging.DEBUG)
			e.updateVectorClock()
		case *traceElementWait:
			logging.Log("Update vector clock for go operation "+e.toString(), logging.DEBUG)
			e.updateVectorClock()
		}

		for i := 1; i <= numberOfRoutines; i++ {
			logging.Log(currentVectorClocks[i].ToString(), logging.DEBUG)
		}

	}

	logging.Log("Vector clock calculation completed", logging.INFO)
}

func getNextElement() traceElement {
	// find the local trace, where the element on which currentIndex points to
	// has the smallest tpost
	var minTpost int = -1
	var minRoutine int = -1
	for routine, trace := range traces {
		// no more elements in the routine trace
		if currentIndex[routine] == -1 {
			continue
		}
		// ignore non executed operations
		if trace[currentIndex[routine]].getTpost() == 0 {
			continue
		}
		if minTpost == -1 || trace[currentIndex[routine]].getTpost() < minTpost {
			minTpost = trace[currentIndex[routine]].getTpost()
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
