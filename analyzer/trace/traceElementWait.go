package trace

import (
	"errors"
	"math"
	"strconv"
)

// enum for opW
type opW int

const (
	ChangeOp opW = iota
	WaitOp
)

/*
* traceElementWait is a trace element for a wait group statement
* Fields:
*   tpre (int): The timestamp at the start of the event
*   tpost (int): The timestamp at the end of the event
*   vpre (vectorClock): The vector clock at the start of the event
*   vpost (vectorClock): The vector clock at the end of the event
*   id (int): The id of the wait group
*   opW (opW): The operation on the wait group
*   delta (int): The delta of the wait group
*   val (int): The value of the wait group
*   pos (string): The position of the wait group in the code
*   pre (*traceElementPre): The pre element of the wait group
 */
type traceElementWait struct {
	routine int
	tpre    int
	tpost   int
	vpre    vectorClock
	vpost   vectorClock
	id      int
	opW     opW
	delta   int
	val     int
	pos     string
	pre     *traceElementPre
}

/*
 * Create a new wait group trace element
 * Args:
 *   routine (int): The routine id
 *   numberOfRoutines (int): The number of routines in the trace
 *   tpre (string): The timestamp at the start of the event
 *   tpost (string): The timestamp at the end of the event
 *   id (string): The id of the wait group
 *   opW (string): The operation on the wait group
 *   delta (string): The delta of the wait group
 *   val (string): The value of the wait group
 *   pos (string): The position of the wait group in the code
 */
func addTraceElementWait(routine int, numberOfRoutines int, tpre string,
	tpost string, id string, opW string, delta string, val string,
	pos string) error {
	tpre_int, err := strconv.Atoi(tpre)
	if err != nil {
		return errors.New("tpre is not an integer")
	}

	tpost_int, err := strconv.Atoi(tpost)
	if err != nil {
		return errors.New("tpost is not an integer")
	}

	id_int, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id is not an integer")
	}

	opW_op := ChangeOp
	if opW == "W" {
		opW_op = WaitOp
	} else if opW != "A" {
		return errors.New("op is not a valid operation")
	}

	delta_int, err := strconv.Atoi(delta)
	if err != nil {
		return errors.New("delta is not an integer")
	}

	val_int, err := strconv.Atoi(val)
	if err != nil {
		return errors.New("val is not an integer")
	}

	elem := traceElementWait{
		routine: routine,
		tpre:    tpre_int,
		tpost:   tpost_int,
		vpre:    newVectorClock(numberOfRoutines),
		vpost:   newVectorClock(numberOfRoutines),
		id:      id_int,
		opW:     opW_op,
		delta:   delta_int,
		val:     val_int,
		pos:     pos}

	// create the pre event
	elem_pre := traceElementPre{
		elem:     &elem,
		elemType: Wait,
	}

	err1 := addElementToTrace(&elem_pre)
	err2 := addElementToTrace(&elem)

	return errors.Join(err1, err2)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (wa *traceElementWait) getRoutine() int {
	return wa.routine
}

/*
 * Get the timestamp at the start of the event
 * Returns:
 *   int: The timestamp at the start of the event
 */
func (wa *traceElementWait) getTpre() int {
	return wa.tpre
}

/*
 * Get the timestamp at the end of the event
 * Returns:
 *   int: The timestamp at the end of the event
 */
func (wa *traceElementWait) getTpost() int {
	return wa.tpost
}

/*
 * Get the vector clock at the begin of the event
 * Returns:
 *   vectorClock: The vector clock at the begin of the event
 */
func (wa *traceElementWait) getVpre() *vectorClock {
	return &wa.vpre
}

/*
 * Get the vector clock at the end of the event
 * Returns:
 *   vectorClock: The vector clock at the end of the event
 */
func (wa *traceElementWait) getVpost() *vectorClock {
	return &wa.vpost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (wa *traceElementWait) getTsort() int {
	if wa.tpost == 0 {
		// add at the end of the trace
		return math.MaxInt
	}
	return wa.tpost
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (wa *traceElementWait) toString() string {
	return "W" + strconv.Itoa(wa.id) + "," + strconv.Itoa(wa.tpre) + "," +
		strconv.Itoa(wa.tpost) + "," + "," +
		strconv.Itoa(wa.delta) + "," + strconv.Itoa(wa.val) + "," + wa.pos
}

/*
 * Update and calculate the vector clock of the element
 * Args:
 *   vc (vectorClock): The current vector clocks
 * TODO: implement
 */
func (wa *traceElementWait) calculateVectorClock(vc *[]vectorClock) {}
