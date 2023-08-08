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
 *   tpost (int): The timestamp of the event
 *   id (int): The id of the atomic variable
 */
type traceElementAtomic struct {
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
func AddTraceElementAtomic(routine int, tpost string, id string,
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

	elem := traceElementAtomic{tpost_int, id_int, operation_int}

	return addElementToTrace(routine, elem)
}

func (elem traceElementAtomic) getSimpleString() string {
	return "A" + strconv.Itoa(elem.id) + "," + strconv.Itoa(elem.tpost) + "," +
		strconv.Itoa(int(elem.operation))
}
