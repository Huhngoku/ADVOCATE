package trace

import (
	"sort"
)

var traces = make(map[int][]traceElement)

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
