package trace

import (
	"analyzer/analysis"
	"analyzer/clock"
	"analyzer/logging"
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
* TraceElementWait is a trace element for a wait group statement
* Fields:
*   tpre (int): The timestamp at the start of the event
*   tpost (int): The timestamp at the end of the event
*   id (int): The id of the wait group
*   opW (opW): The operation on the wait group
*   delta (int): The delta of the wait group
*   val (int): The value of the wait group
*   pos (string): The position of the wait group in the code
*   tID (string): The id of the trace element, contains the position and the tpre
 */
type TraceElementWait struct {
	routine int
	tPre    int
	tPost   int
	id      int
	opW     opW
	delta   int
	val     int
	pos     string
	tID     string
	vc      clock.VectorClock
}

/*
 * Create a new wait group trace element
 * Args:
 *   routine (int): The routine id
 *   tpre (string): The timestamp at the start of the event
 *   tpost (string): The timestamp at the end of the event
 *   id (string): The id of the wait group
 *   opW (string): The operation on the wait group
 *   delta (string): The delta of the wait group
 *   val (string): The value of the wait group
 *   pos (string): The position of the wait group in the code
 */
func AddTraceElementWait(routine int, tpre string,
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

	elem := TraceElementWait{
		routine: routine,
		tPre:    tpre_int,
		tPost:   tpost_int,
		id:      id_int,
		opW:     opW_op,
		delta:   delta_int,
		val:     val_int,
		pos:     pos,
		tID:     pos + "@" + tpre,
	}

	return AddElementToTrace(&elem)
}

/*
 * Get the id of the element
 * Returns:
 *   int: The id of the element
 */
func (wa *TraceElementWait) GetID() int {
	return wa.id
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (wa *TraceElementWait) GetRoutine() int {
	return wa.routine
}

/*
 * Get the timestamp at the start of the event
 * Returns:
 *   int: The timestamp at the start of the event
 */
func (wa *TraceElementWait) GetTPre() int {
	return wa.tPre
}

/*
 * Set the tpre of the element.
 * Args:
 *   tPre (int): The tpre of the element
 */
func (wa *TraceElementWait) SetTPre(tPre int) {
	wa.tPre = tPre
	if wa.tPost != 0 && wa.tPost < tPre {
		wa.tPost = tPre
	}
}

/*
 * Get the timestamp at the end of the event
 * Returns:
 *   int: The timestamp at the end of the event
 */
func (wa *TraceElementWait) getTpost() int {
	return wa.tPost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (wa *TraceElementWait) GetTSort() int {
	if wa.tPost == 0 {
		// add at the end of the trace
		return math.MaxInt
	}
	return wa.tPost
}

/*
 * Get the position of the operation.
 * Returns:
 *   string: The position of the element
 */
func (wa *TraceElementWait) GetPos() string {
	return wa.pos
}

/*
 * Get the tID of the element.
 * Returns:
 *   string: The tID of the element
 */
func (wa *TraceElementWait) GetTID() string {
	return wa.tID
}

/*
 * Set the timer, that is used for the sorting of the trace
 * Args:
 *   tSort (int): The timer of the element
 */
func (wa *TraceElementWait) SetTSort(tSort int) {
	wa.SetTPre(tSort)
	wa.tPost = tSort
}

/*
 * Set the timer, that is used for the sorting of the trace, only if the original
 * value was not 0
 * Args:
 *   tSort (int): The timer of the element
 */
func (wa *TraceElementWait) SetTSortWithoutNotExecuted(tSort int) {
	wa.SetTPre(tSort)
	if wa.tPost != 0 {
		wa.tPost = tSort
	}
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (wa *TraceElementWait) ToString() string {
	res := "W,"
	res += strconv.Itoa(wa.tPre) + "," + strconv.Itoa(wa.tPost) + ","
	res += strconv.Itoa(wa.id) + ","
	switch wa.opW {
	case ChangeOp:
		res += "A,"
	case WaitOp:
		res += "W,"
	}

	res += strconv.Itoa(wa.delta) + "," + strconv.Itoa(wa.val)
	res += "," + wa.pos
	return res
}

/*
 * Update and calculate the vector clock of the element
 */
func (wa *TraceElementWait) updateVectorClock() {
	switch wa.opW {
	case ChangeOp:
		analysis.Change(wa.routine, wa.id, wa.delta, wa.tID, currentVCHb)
	case WaitOp:
		analysis.Wait(wa.routine, wa.id, wa.tID, currentVCHb)
	default:
		err := "Unknown operation on wait group: " + wa.ToString()
		logging.Debug(err, logging.ERROR)
	}

	wa.vc = currentVCHb[wa.routine].Copy()
}

func (wa *TraceElementWait) GetDelta() int {
	return wa.delta
}

/*
 * Get the vector clock of the element
 * Returns:
 *   VectorClock: The vector clock of the element
 */
func (wa *TraceElementWait) GetVC() clock.VectorClock {
	return wa.vc
}
