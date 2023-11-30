package trace

import (
	"errors"
	"math"
	"strconv"
)

// enum for opC
type opChannel int

const (
	send opChannel = iota
	recv
	close
)

var waitingReceive = make([]*traceElementChannel, 0)
var maxOpID = make(map[int]int)

/*
* traceElementChannel is a trace element for a channel
* Fields:
*   routine (int): The routine id
*   tpre (int): The timestamp at the start of the event
*   tpost (int): The timestamp at the end of the event
*   id (int): The id of the channel
*   opC (int, enum): The operation on the channel
*   exec (int, enum): The execution status of the operation
*   oId (int): The id of the other communication
*   qSize (int): The size of the channel queue
*   qCount (int): The number of elements in the queue after the operation
*   pos (string): The position of the channel operation in the code
*   sel (*traceElementSelect): The select operation, if the channel operation is part of a select, otherwise nil
*   partner (*traceElementChannel): The partner of the channel operation
 */
type traceElementChannel struct {
	routine int
	tpre    int
	tpost   int
	id      int
	opC     opChannel
	cl      bool
	oID     int
	qSize   int
	pos     string
	sel     *traceElementSelect
}

/*
* Create a new channel trace element
* Args:
*   routine (int): The routine id
*   tPre (string): The timestamp at the start of the event
*   tPost (string): The timestamp at the end of the event
*   id (string): The id of the channel
*   opC (string): The operation on the channel
*   cl (string): Whether the channel was finished because it was closed
*   oId (string): The id of the other communication
*   qSize (string): The size of the channel queue
*   pos (string): The position of the channel operation in the code
 */
func AddTraceElementChannel(routine int, tPre string,
	tPost string, id string, opC string, cl string, oID string, qSize string,
	pos string) error {
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

	var opCInt opChannel
	switch opC {
	case "S":
		opCInt = send
	case "R":
		opCInt = recv
	case "C":
		opCInt = close
	default:
		return errors.New("opC is not a valid value")
	}

	clBool, err := strconv.ParseBool(cl)
	if err != nil {
		return errors.New("suc is not a boolean")
	}

	oIDInt, err := strconv.Atoi(oID)
	if err != nil {
		return errors.New("oId is not an integer")
	}

	qSizeInt, err := strconv.Atoi(qSize)
	if err != nil {
		return errors.New("qSize is not an integer")
	}

	elem := traceElementChannel{
		routine: routine,
		tpre:    tPreInt,
		tpost:   tPostInt,
		id:      idInt,
		opC:     opCInt,
		cl:      clBool,
		oID:     oIDInt,
		qSize:   qSizeInt,
		pos:     pos,
	}

	return addElementToTrace(&elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (ch *traceElementChannel) getRoutine() int {
	return ch.routine
}

/*
 * Get the tpre of the element
 * Returns:
 *   int: The tpre of the element
 */
func (ch *traceElementChannel) getTpre() int {
	return ch.tpre
}

/*
 * Get the tpost of the element
 * Returns:
 *   int: The tpost of the element
 */
func (ch *traceElementChannel) getTpost() int {
	return ch.tpost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   float32: The time of the element
 */
func (ch *traceElementChannel) getTsort() int {
	if ch.tpost == 0 {
		// add to the end of the trace
		return math.MaxInt
	}
	return ch.tpost
}

/*
 * Get the vector clock at the begin of the event
 * Returns:
 *   vectorClock: The vector clock at the begin of the event
 */
// func (ch *traceElementChannel) getVpre() *vc.VectorClock {
// 	return &ch.vpre
// }

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (ch *traceElementChannel) toString() string {
	return ch.toStringSep(",", true)
}

/*
 * Get the simple string representation of the element
 * Args:
 *   sep (string): The separator between the values
 *   pos (bool): Whether the position should be included
 * Returns:
 *   string: The simple string representation of the element
 */
func (ch *traceElementChannel) toStringSep(sep string, pos bool) string {
	res := "C," + strconv.Itoa(ch.tpre) + sep + strconv.Itoa(ch.tpost) + sep +
		strconv.Itoa(ch.id) + sep + strconv.Itoa(int(ch.opC)) + sep +
		strconv.Itoa(ch.oID) + sep + strconv.Itoa(ch.qSize)
	if pos {
		res += sep + ch.pos
	}
	return res
}
