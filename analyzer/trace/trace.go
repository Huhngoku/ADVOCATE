package trace

import "errors"

var trace map[int][]traceElement = make(map[int][]traceElement)

/*
* Create a new routine in the trace
* Args:
*   routine (int): The routine id
 */
func NewRoutine(routine int) error {
	if _, ok := trace[routine]; ok {
		return errors.New("routine already exists")
	}
	trace[routine] = make([]traceElement, 0)
	return nil
}

/*
* Add an element to the trace
* Args:
*   routine (int): The routine id
*   element (TraceElement): The element to add
 */
func addElementToTrace(routine int, element traceElement) error {
	if _, ok := trace[routine]; !ok {
		return errors.New("routine does not exist")
	}
	trace[routine] = append(trace[routine], element)
	return nil
}
