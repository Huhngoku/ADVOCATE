package trace

import (
	"sort"

	"analyzer/debug"
)

var trace []traceElement = make([]traceElement, 0)

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
	debug.Log("Find partners finished", 2)
}

/*
* Calculate vector clocks
 */
func CalculateVectorClocks() {
	debug.Log("Calculate vector clocks...", 2)

	debug.Log("Vector clock calculation finished", 2)
}

/*
 * Sort the trace by tpre
 */
type sortByTPost []traceElement

func (a sortByTPost) Len() int           { return len(a) }
func (a sortByTPost) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByTPost) Less(i, j int) bool { return a[i].getTpost() < a[j].getTpost() }
func Sort() {
	debug.Log("Sort Trace...", 2)
	sort.Sort(sortByTPost(trace))
	debug.Log("Trace sorted", 2)
}
