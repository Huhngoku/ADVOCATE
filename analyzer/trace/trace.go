package trace

import (
	"analyzer/debug"
	"fmt"
	"sort"
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
func addElementToTrace(routine int, element traceElement) error {
	trace = append(trace, element)
	return nil
}

/*
 * Function to start the search for all partner elements
 * TODO: only channel is implemented, missing select, mutex, ...
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

/*
 * Check if all channel and mutex operations have a partner (if they should have one)
 * Returns:
 *   bool: True if all channel operations have a partner, false otherwise
 * TODO:
 *   remove
 */
func CheckTrace() bool {
	res := true
	for _, element := range trace {
		switch elem := element.(type) {
		case *traceElementChannel:
			if elem.opC == 2 { // close
				continue
			}
			if elem.partner == nil {
				fmt.Println(elem.toString(), "Error")
				res = false
			} else {
				fmt.Println(elem.toString(), "Ok")
			}
		case *traceElementMutex:
			if elem.partner == nil {
				fmt.Println(elem.toString(), "Error")
				res = false
			} else {
				fmt.Println(elem.toString(), "Ok")
			}
		}

	}
	return res
}
