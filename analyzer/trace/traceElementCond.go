package trace

import (
	"analyzer/analysis"
	"analyzer/clock"
	"errors"
	"math"
	"strconv"
)

type opCond int

const (
	WaitCondOp opCond = iota
	SignalOp
	BroadcastOp
)

/*
 * TraceElementCond is a trace element for a condition variable
 * Fields:
 *   routine (int): The routine id
 *   tpre (int): The timestamp at the start of the event
 *   tpost (int): The timestamp at the end of the event
 *   id (int): The id of the condition variable
 *   opC (opCond): The operation on the condition variable
 *   pos (string): The position of the condition variable operation in the code
 *   tID (string): The id of the trace element, contains the position and the tpre
 */
type TraceElementCond struct {
	routine int
	tPre    int
	tPost   int
	id      int
	opC     opCond
	pos     string
	tID     string
	vc      clock.VectorClock
}

/*
 * Create a new condition variable trace element
 * Args:
 *   routine (int): The routine id
 *   tPre (string): The timestamp at the start of the event
 *   tPost (string): The timestamp at the end of the event
 *   id (string): The id of the condition variable
 *   opC (string): The operation on the condition variable
 *   pos (string): The position of the condition variable operation in the code
 *   tID (string): The id of the trace element, contains the position and the tpre
 */
func AddTraceElementCond(routine int, tPre string, tPost string, id string, opN string, pos string) error {
	tPreInt, err := strconv.Atoi(tPre)
	if err != nil {
		return errors.New("tpre is not an integer")
	}
	tPostInt, err := strconv.Atoi(tPost)
	if err != nil {
		return errors.New("tpost is not an integer")
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id is not an integer")
	}
	var op opCond
	switch opN {
	case "W":
		op = WaitCondOp
	case "S":
		op = SignalOp
	case "B":
		op = BroadcastOp
	default:
		return errors.New("op is not a valid operation")
	}

	tIDStr := pos + "@" + strconv.Itoa(tPreInt)

	elem := TraceElementCond{
		routine: routine,
		tPre:    tPreInt,
		tPost:   tPostInt,
		id:      idInt,
		opC:     op,
		pos:     pos,
		tID:     tIDStr,
	}

	return AddElementToTrace(&elem)
}

/*
 * Get the id of the element
 * Returns:
 *   int: The id of the element
 */
func (co *TraceElementCond) GetID() int {
	return co.id
}

/*
 * Get the routine of the element
 * Returns:
 *   (int): The routine id
 */
func (co *TraceElementCond) GetRoutine() int {
	return co.routine
}

/*
 * Get the tpre of the element.
 * Returns:
 *   int: The tpre of the element
 */
func (co *TraceElementCond) GetTPre() int {
	return co.tPre
}

/*
 * Set the tpre of the element.
 * Args:
 *   tPre (int): The tpre of the element
 */
func (co *TraceElementCond) SetTPre(tPre int) {
	co.tPre = tPre
	if co.tPost != 0 && co.tPost < tPre {
		co.tPost = tPre
	}
}

/*
 * Get the tpost of the element.
 * Returns:
 *   int: The tpost of the element
 */
func (co *TraceElementCond) getTpost() int {
	return co.tPost
}

/*
 * Get the timer, that is used for sorting the trace
 * Returns:
 *   (int): The timer of the element
 * TODO: check if tPre is correct
 */
func (co *TraceElementCond) GetTSort() int {
	t := co.tPre
	if co.opC == WaitCondOp {
		t = co.tPost
	}
	if t == 0 {
		// add at the end of the trace
		return math.MaxInt
	}
	return t
}

/*
 * Get the position of the operation
 * Returns:
 *   (string): The position of the operation
 */
func (co *TraceElementCond) GetPos() string {
	return co.pos
}

/*
 * Get the tID of the element.
 * Returns:
 *   string: The tID of the element
 */
func (co *TraceElementCond) GetTID() string {
	return co.tID
}

/*
 * Set the timer that is used for sorting the trace
 * Args:
 *   tSort (int): The timer of the element
 * TODO: check if tPre is correct
 */
func (co *TraceElementCond) SetTSort(tSort int) {
	co.SetTPre(tSort)
	if co.opC == WaitCondOp {
		co.tPost = tSort
	}
}

/*
 * Set the timer, that is used for the sorting of the trace, only if the original
 * value was not 0
 * Args:
 *   tSort (int): The timer of the element
 */
func (co *TraceElementCond) SetTSortWithoutNotExecuted(tSort int) {
	co.SetTPre(tSort)
	if co.opC == WaitCondOp {
		if co.tPost != 0 {
			co.tPost = tSort
		}
		return
	}
	if co.tPre != 0 {
		co.tPre = tSort
	}
	return
}

/*
 * Get the string representation of the element
 * Returns:
 *   (string): The string representation of the element
 */
func (co *TraceElementCond) ToString() string {
	res := "N,"
	res += strconv.Itoa(co.tPre) + "," + strconv.Itoa(co.tPost) + ","
	res += strconv.Itoa(co.id) + ","
	switch co.opC {
	case WaitCondOp:
		res += "W"
	case SignalOp:
		res += "S"
	case BroadcastOp:
		res += "B"
	}
	res += "," + co.pos
	return res
}

/*
 * Update the vector clock of the trace and element
 */
func (co *TraceElementCond) updateVectorClock() {
	switch co.opC {
	case WaitCondOp:
		analysis.CondWait(co.id, co.routine, currentVCHb)
	case SignalOp:
		analysis.CondSignal(co.id, co.routine, currentVCHb)
	case BroadcastOp:
		analysis.CondBroadcast(co.id, co.routine, currentVCHb)
	}

	co.vc = currentVCHb[co.routine].Copy()
}

/*
 * Get the vector clock of the element
 * Returns:
 *   VectorClock: The vector clock of the element
 */
func (co *TraceElementCond) GetVC() clock.VectorClock {
	return co.vc
}
