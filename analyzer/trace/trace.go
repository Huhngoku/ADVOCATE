package trace

import (
	"sort"

	"analyzer/debug"
)

var trace []traceElement = make([]traceElement, 0)
var currentVectorClocks []vectorClock = make([]vectorClock, 0)

/*
* Get the trace
* Returns:
*   []traceElement: The trace
 */
func GetTrace() []traceElement {
	return trace
}

/*
* Add an element to the trace
* Args:
*   routine (int): The routine id
*   element (TraceElement): The element to add
* Returns:
*   error: An error if the routine does not exist
 */
func addElementToTrace(element traceElement) error {
	trace = append(trace, element)
	return nil
}

/*
 * Function to start the search for all partner elements
 */
func FindPartner() {
	debug.Log("Find partners...", 2)
	for _, elem := range trace {
		switch e := elem.(type) {
		case *traceElementChannel:
			debug.Log("Find partner for channel operation "+e.toString(), 3)
			e.findPartner()
		case *traceElementSelect:
			debug.Log("Find partner for select operation "+e.toString(), 3)
			for _, c := range e.cases {
				c.findPartner()
			}
		case *traceElementMutex:
			debug.Log("Find partner for mutex operation "+e.toString(), 3)
			e.findPartner()
		}
	}

	// check if there are operations without partner
	checkChannelOperations()
	debug.Log("Partners found", 2)
}

/*
* Calculate vector clocks
 */
func CalculateVectorClocks(numberOfRoutines int) {
	debug.Log("Calculate vector clocks...", 2)

	// create current vector clock
	currentVectorClocks = make([]vectorClock, numberOfRoutines)
	for i := 0; i < numberOfRoutines; i++ {
		currentVectorClocks[i] = newVectorClock(numberOfRoutines)
	}

	for _, elem := range trace {
		switch e := elem.(type) {
		case *traceElementAtomic:
			debug.Log("Calculate vector clock for atomic operation "+e.toString(), 3)
			e.calculateVectorClock(&currentVectorClocks)
		case *traceElementChannel:
			debug.Log("Calculate vector clock for channel operation "+e.toString(), 3)
			e.calculateVectorClock(&currentVectorClocks)
		case *traceElementMutex:
			debug.Log("Calculate vector clock for mutex operation "+e.toString(), 3)
			e.calculateVectorClock(&currentVectorClocks)
		case *traceElementRoutine:
			debug.Log("Calculate vector clock for routine operation "+e.toString(), 3)
			e.calculateVectorClock(&currentVectorClocks)
		case *traceElementSelect:
			debug.Log("Calculate vector clock for select operation "+e.toString(), 3)
			e.calculateVectorClock(&currentVectorClocks)
		case *traceElementWait:
			debug.Log("Calculate vector clock for go operation "+e.toString(), 3)
			e.calculateVectorClock(&currentVectorClocks)
		case *traceElementPre:
			debug.Log("Calculate vector clock for pre operation "+e.toString(), 3)
			e.calculateVectorClock(&currentVectorClocks)
		}
	}

	debug.Log("Vector clock calculation completed", 2)
}

/*
 * Sort the trace by tpre
 */
type sortByTPost []traceElement

func (a sortByTPost) Len() int      { return len(a) }
func (a sortByTPost) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortByTPost) Less(i, j int) bool {
	return a[i].getTsort() < a[j].getTsort()
}
func Sort() {
	debug.Log("Sort Trace...", 2)
	sort.Sort(sortByTPost(trace))
	debug.Log("Trace sorted", 2)
}
