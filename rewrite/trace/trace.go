package trace

import (
	"errors"
	"sort"
)

var traces = make(map[int][]TraceElement)

/*
* Add an element to the trace
* Args:
*   routine (int): The routine id
*   element (TraceElement): The element to add
* Returns:
*   error: An error if the routine does not exist
 */
func addElementToTrace(element TraceElement) error {
	routine := element.GetRoutine()
	traces[routine] = append(traces[routine], element)
	return nil
}

func AddEmptyRoutine(routine int) {
	traces[routine] = make([]TraceElement, 0)
}

/*
 * Sort the trace by tsort
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

func Sort() {
	for routine, trace := range traces {
		if len(traces[routine]) <= 1 {
			continue
		}
		traces[routine] = sortTrace(trace)
	}
}

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
	for routine := 0; routine < len(traces); routine++ {
		for j := 0; j < len(traces[routine]); j++ {
			if traces[routine][j].GetPos() == pos {
				return &traces[routine][j], nil
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
	for routine, localTrace := range traces {
		for _, elem := range localTrace {
			if elem.GetTSort() >= startTime && !contains(excludedRoutines, routine) {
				elem.SetTsortWithoutNotExecuted(elem.GetTSort() + steps)
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
