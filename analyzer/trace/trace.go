package trace

import (
	"errors"
	"fmt"
)

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

func CheckTraceChannel() bool {
	res := true
	for i, routine := range trace {
		for j, element := range routine {
			switch elem := element.(type) {
			case traceElementChannel:
				if elem.partner == nil {
					fmt.Println(i, j, "Error")
					res = false
				} else {
					fmt.Println(i, j, "Ok")
				}
			}
		}
	}
	return res
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
	if _, ok := trace[routine]; !ok {
		return errors.New("routine does not exist")
	}
	trace[routine] = append(trace[routine], element)
	return nil
}
