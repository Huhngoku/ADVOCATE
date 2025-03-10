package trace

import (
	"analyzer/analysis"
	"analyzer/clock"
	"errors"
	"math"
	"strconv"
)

type OpCond int

const (
	WaitCondOp OpCond = iota
	SignalOp
	BroadcastOp
)

/*
 * TraceElementCond is a trace element for a condition variable
 * MARK: Struct
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
	opC     OpCond
	pos     string
	tID     string
	vc      clock.VectorClock
}

/*
 * Create a new condition variable trace element
 * MARK: New
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
	var op OpCond
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

// MARK: Getter

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
 * Get the operation of the element
 * Returns:
 *   (OpCond): The operation of the element
 */
func (co *TraceElementCond) GetOpCond() OpCond {
	return co.opC
}

/*
 * Get the vector clock of the element
 * Returns:
 *   VectorClock: The vector clock of the element
 */
func (co *TraceElementCond) GetVC() clock.VectorClock {
	return co.vc
}

/*
 * Get all to element concurrent wait, broadcast and signal operations on the same condition variable
 * Args:
 *   element (traceElement): The element
 *   filter ([]string): The types of the elements to return
 * Returns:
 *   []*traceElement: The concurrent elements
 */
func GetConcurrentWaitgroups(element *TraceElement) map[string][]*TraceElement {
	res := make(map[string][]*TraceElement)
	res["broadcast"] = make([]*TraceElement, 0)
	res["signal"] = make([]*TraceElement, 0)
	res["wait"] = make([]*TraceElement, 0)
	for _, trace := range traces {
		for _, elem := range trace {
			switch elem.(type) {
			case *TraceElementCond:
			default:
				continue
			}

			if elem.GetTID() == (*element).GetTID() {
				continue
			}

			e := elem.(*TraceElementCond)

			if e.opC == WaitCondOp {
				continue
			}

			if clock.GetHappensBefore((*element).GetVC(), e.GetVC()) == clock.Concurrent {
				e := elem.(*TraceElementCond)
				if e.opC == SignalOp {
					res["signal"] = append(res["signal"], &elem)
				} else if e.opC == BroadcastOp {
					res["broadcast"] = append(res["broadcast"], &elem)
				} else if e.opC == WaitCondOp {
					res["wait"] = append(res["wait"], &elem)
				}
			}
		}
	}
	return res
}

// MARK: Setter

/*
 * Set the tPre and tPost of the element
 * Args:
 *   time (int): The tPre and tPost of the element
 */
func (co *TraceElementCond) SetT(time int) {
	co.tPre = time
	co.tPost = time
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
func (co *TraceElementCond) SetTWithoutNotExecuted(tSort int) {
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
 * MARK: ToString
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
 * MARK: VectorClock
 */
func (co *TraceElementCond) updateVectorClock() {
	switch co.opC {
	case WaitCondOp:
		analysis.CondWait(co.id, co.routine, currentVCHb, co.tPost == 0)
	case SignalOp:
		analysis.CondSignal(co.id, co.routine, currentVCHb)
	case BroadcastOp:
		analysis.CondBroadcast(co.id, co.routine, currentVCHb)
	}

	co.vc = currentVCHb[co.routine].Copy()
}

// MARK: Copy

/*
 * Copy the element
 * Returns:
 *   (TraceElement): The copy of the element
 */
func (co *TraceElementCond) Copy() TraceElement {
	return &TraceElementCond{
		routine: co.routine,
		tPre:    co.tPre,
		tPost:   co.tPost,
		id:      co.id,
		opC:     co.opC,
		pos:     co.pos,
		tID:     co.tID,
		vc:      co.vc.Copy(),
	}
}
