package trace

import (
	"errors"
	"strconv"
)

// enum for operation
type opAtomic int

const (
	LoadOp opAtomic = iota
	StoreOp
	AddOp
	SwapOp
	CompSwapOp
)

/*
 * Struct to save an atomic event in the trace
 * Fields:
 *   routine (int): The routine id
 *   tpost (int): The timestamp of the event
 *   id (int): The id of the atomic variable
 *   operation (int, enum): The operation on the atomic variable
 */
type traceElementAtomic struct {
	routine   int
	tpost     int
	id        int
	operation opAtomic
}

/*
 * Create a new atomic trace element
 * Args:
 *   routine (int): The routine id
 *   tpost (string): The timestamp of the event
 *   id (string): The id of the atomic variable
 */
func addTraceElementAtomic(routine int, tpost string, id string,
	operation string) error {
	tpost_int, err := strconv.Atoi(tpost)
	if err != nil {
		return errors.New("tpost is not an integer")
	}

	id_int, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("tpost is not an integer")
	}

	var operation_int opAtomic = 0
	switch operation {
	case "L":
		operation_int = LoadOp
	case "S":
		operation_int = StoreOp
	case "A":
		operation_int = AddOp
	case "W":
		operation_int = SwapOp
	case "C":
		operation_int = CompSwapOp
	default:
		return errors.New("operation is not a valid operation")
	}

	elem := traceElementAtomic{routine, tpost_int, id_int, operation_int}

	return addElementToTrace(&elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (elem *traceElementAtomic) getRoutine() int {
	return elem.routine
}

/*
 * Get the tpost of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpost of the element
 */
func (elem *traceElementAtomic) getTpre() int {
	return elem.tpost
}

/*
 * Get the tpost of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpost of the element
 */
func (elem *traceElementAtomic) getTpost() int {
	return elem.tpost
}

/*
 * Get the simple string representation of the element.
 * Returns:
 *   string: The simple string representation of the element
 */
func (elem *traceElementAtomic) toString() string {
	return "A" + strconv.Itoa(elem.id) + "," + strconv.Itoa(elem.tpost) + "," +
		strconv.Itoa(int(elem.operation))
}
