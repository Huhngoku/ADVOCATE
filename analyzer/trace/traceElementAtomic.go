package trace

import (
	"analyzer/logging"
	vc "analyzer/vectorClock"
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
 *   vpost (vectorClock): The vector clock at the end of the event
 *   id (int): The id of the atomic variable
 *   operation (int, enum): The operation on the atomic variable
 */
type traceElementAtomic struct {
	routine int
	tpost   int
	vpost   vc.VectorClock
	id      int
	opA     opAtomic
}

/*
 * Create a new atomic trace element
 * Args:
 *   routine (int): The routine id
 *   numberOfRoutines (int): The number of routines in the trace
 *   tpost (string): The timestamp of the event
 *   id (string): The id of the atomic variable
 *   operation (string): The operation on the atomic variable
 */
func AddTraceElementAtomic(routine int, numberOfRoutines int, tpost string,
	id string, operation string) error {
	tpost_int, err := strconv.Atoi(tpost)
	if err != nil {
		return errors.New("tpost is not an integer")
	}

	id_int, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("tpost is not an integer")
	}

	var opA_int opAtomic = 0
	switch operation {
	case "L":
		opA_int = LoadOp
	case "S":
		opA_int = StoreOp
	case "A":
		opA_int = AddOp
	case "W":
		opA_int = SwapOp
	case "C":
		opA_int = CompSwapOp
	default:
		return errors.New("operation is not a valid operation")
	}

	elem := traceElementAtomic{
		routine: routine,
		tpost:   tpost_int,
		vpost:   vc.NewVectorClock(numberOfRoutines),
		id:      id_int,
		opA:     opA_int,
	}

	return addElementToTrace(&elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (at *traceElementAtomic) getRoutine() int {
	return at.routine
}

/*
 * Get the tpost of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpost of the element
 */
func (at *traceElementAtomic) getTpre() int {
	return at.tpost
}

/*
 * Get the tpost of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpost of the element
 */
func (at *traceElementAtomic) getTpost() int {
	return at.tpost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (at *traceElementAtomic) getTsort() int {
	return at.tpost
}

/*
 * Get the vector clock at the begin of the event. It is equal to the vector clock
 * at the end of the event.
 * Returns:
 *   vectorClock: The vector clock at the begin of the event
 */
// func (at *traceElementAtomic) getVpre() *vc.VectorClock {
// 	return &at.vpost
// }

/*
 * Get the vector clock at the end of the event
 * Returns:
 *   vectorClock: The vector clock at the end of the event
 */
func (at *traceElementAtomic) getVpost() *vc.VectorClock {
	return &at.vpost
}

/*
 * Get the simple string representation of the element.
 * Returns:
 *   string: The simple string representation of the element
 */
func (at *traceElementAtomic) toString() string {
	return "A," + strconv.Itoa(at.tpost) + "," + strconv.Itoa(at.id) + "," +
		strconv.Itoa(int(at.opA))
}

/*
 * Update and calculate the vector clock of the element
 */
func (at *traceElementAtomic) updateVectorClock() {
	switch at.opA {
	case LoadOp:
		at.vpost = vc.Read(at.routine, at.id, currentVectorClocks)
	case StoreOp, AddOp:
		at.vpost = vc.Write(at.routine, at.id, currentVectorClocks)
	case SwapOp, CompSwapOp:
		at.vpost = vc.Swap(at.routine, at.id, currentVectorClocks)
	default:
		err := "Unknown operation: " + at.toString()
		logging.Debug(err, logging.ERROR)
	}
}
