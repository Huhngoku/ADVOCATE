package trace

import (
	"analyzer/debug"
	vc "analyzer/vectorClock"
	"sort"
)

var trace []traceElement = make([]traceElement, 0)
var currentVectorClocks map[int]vc.VectorClock = make(map[int]vc.VectorClock)
var numberOfRoutines int = 0

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
 * Set the number of routines
 * Args:
 *   n (int): The number of routines
 */
func SetNumberOfRoutines(n int) {
	numberOfRoutines = n
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
func CalculateVectorClocks() {
	debug.Log("Calculate vector clocks...", 2)

	for i := 1; i <= numberOfRoutines; i++ {
		currentVectorClocks[i] = vc.NewVectorClock(numberOfRoutines)
	}

	for _, elem := range trace {
		switch e := elem.(type) {
		case *traceElementAtomic:
			debug.Log("Update vector clock for atomic operation "+e.toString(), 3)
			e.updateVectorClock()
		case *traceElementChannel:
			debug.Log("Update vector clock for channel operation "+e.toString(), 3)
			e.updateVectorClock()
		case *traceElementMutex:
			debug.Log("Update vector clock for mutex operation "+e.toString(), 3)
			e.updateVectorClock()
		case *traceElementRoutine:
			debug.Log("Update vector clock for routine operation "+e.toString(), 3)
			e.updateVectorClock()
		case *traceElementSelect:
			debug.Log("Update vector clock for select operation "+e.toString(), 3)
			e.updateVectorClock()
		case *traceElementWait:
			debug.Log("Update vector clock for go operation "+e.toString(), 3)
			e.updateVectorClock()
		}
	}

	debug.Log("Vector clock calculation completed", 2)
}

// TODO: change to interlace not sort
/*
 * Sort the trace by tpost
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
