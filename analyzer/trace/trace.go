package trace

import (
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
	for _, elem := range trace {
		switch e := elem.(type) {
		case *traceElementChannel:
			e.findPartner()
		case *traceElementSelect:
			for _, c := range e.cases {
				c.findPartner()
			}
		}
	}
}

/*
 * Sort the trace by tpre
 */
type sortByTPost []traceElement

func (a sortByTPost) Len() int           { return len(a) }
func (a sortByTPost) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByTPost) Less(i, j int) bool { return a[i].getTpost() < a[j].getTpost() }
func Sort()                              { sort.Sort(sortByTPost(trace)) }

/*
 * Check if all channel operations have a partner
 * Returns:
 *   bool: True if all channel operations have a partner, false otherwise
 * TODO:
 *   remove
 */
func CheckTraceChannel() bool {
	res := true
	for i, element := range trace {
		switch elem := element.(type) {
		case *traceElementChannel:
			if elem.opC == 2 { // close
				continue
			}
			if elem.partner == nil {
				fmt.Println(i, elem.toString(), "Error")
				res = false
			} else {
				fmt.Println(i, elem.toString(), "Ok")
			}
		}

	}
	return res
}
