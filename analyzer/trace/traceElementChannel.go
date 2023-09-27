package trace

import (
	"errors"
	"math"
	"strconv"

	"analyzer/debug"
	vc "analyzer/vectorClock"
)

// enum for opC
type opChannel int

const (
	send opChannel = iota
	recv
	close
)

/*
* traceElementChannel is a trace element for a channel
* Fields:
*   routine (int): The routine id
*   tpre (int): The timestamp at the start of the event
*   tpost (int): The timestamp at the end of the event
*   vpre (vectorClock): The vector clock at the start of the event
*   vpost (vectorClock): The vector clock at the end of the event
*   id (int): The id of the channel
*   opC (int, enum): The operation on the channel
*   exec (int, enum): The execution status of the operation
*   oId (int): The id of the other communication
*   qSize (int): The size of the channel queue
*   qCountPre (int): The number of elements in the queue before the operation
*   qCountPost (int): The number of elements in the queue after the operation
*   pos (string): The position of the channel operation in the code
*   sel (*traceElementSelect): The select operation, if the channel operation is part of a select, otherwise nil
*   partner (*traceElementChannel): The partner of the channel operation
 */
type traceElementChannel struct {
	routine int
	tpre    int
	tpost   int
	vpre    vc.VectorClock
	vpost   vc.VectorClock
	id      int
	opC     opChannel
	cl      bool
	oId     int
	qSize   int
	pos     string
	sel     *traceElementSelect
	partner *traceElementChannel
}

/*
* Create a new channel trace element
* Args:
*   routine (int): The routine id
*   numberOfRoutines (int): The number of routines in the trace
*   tpre (string): The timestamp at the start of the event
*   tpost (string): The timestamp at the end of the event
*   id (string): The id of the channel
*   opC (string): The operation on the channel
*   cl (string): Whether the channel was finished because it was closed
*   oId (string): The id of the other communication
*   qSize (string): The size of the channel queue
*   pos (string): The position of the channel operation in the code
 */
func AddTraceElementChannel(routine int, numberOfRoutines int, tpre string,
	tpost string, id string, opC string, cl string, oId string, qSize string,
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

	var opC_int opChannel = 0
	switch opC {
	case "S":
		opC_int = send
	case "R":
		opC_int = recv
	case "C":
		opC_int = close
	default:
		return errors.New("opC is not a valid value")
	}

	cl_bool, err := strconv.ParseBool(cl)
	if err != nil {
		return errors.New("suc is not a boolean")
	}

	oId_int, err := strconv.Atoi(oId)
	if err != nil {
		return errors.New("oId is not an integer")
	}

	qSize_int, err := strconv.Atoi(qSize)
	if err != nil {
		return errors.New("qSize is not an integer")
	}

	elem := traceElementChannel{
		routine: routine,
		tpre:    tpre_int,
		tpost:   tpost_int,
		// vpre:    vc.NewVectorClock(numberOfRoutines),
		vpost: vc.NewVectorClock(numberOfRoutines),
		id:    id_int,
		opC:   opC_int,
		cl:    cl_bool,
		oId:   oId_int,
		qSize: qSize_int,
		pos:   pos,
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
	if ch.partner == nil {
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
 * Get the vector clock at the end of the event
 * Returns:
 *   vectorClock: The vector clock at the end of the event
 */
func (ch *traceElementChannel) getVpost() *vc.VectorClock {
	return &ch.vpost
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (ch *traceElementChannel) toString() string {
	return ch.toStringSep(",")
}

func (ch *traceElementChannel) toStringSep(sep string) string {
	return "C," + strconv.Itoa(ch.tpre) + sep + strconv.Itoa(ch.tpost) + sep +
		strconv.Itoa(ch.id) + sep + strconv.Itoa(int(ch.opC)) + sep +
		strconv.Itoa(ch.oId) + sep + strconv.Itoa(ch.qSize) + sep + ch.pos
}

// list to store operations where partner has not jet been found
var channelOperations = make([]*traceElementChannel, 0)

// list to store close operations, to find operations, that were executed
// because of a close on a channel
var closeOperations = make([]*traceElementChannel, 0)

/*
 * Update and calculate the vector clock of the element
 * TODO: implement
 */
func (ch *traceElementChannel) updateVectorClock() {
	if ch.qSize == 0 { // unbuffered channel
		switch ch.opC {
		case send, recv:
			partnerRoutine := ch.findUnbufferedPartner()
			if partnerRoutine != -1 {
				vc.Unbuffered(ch.routine, partnerRoutine, currentVectorClocks)
				// advance index of receive routine, send routine is already advanced
				increaseIndex(partnerRoutine)
			}
		}
	}
}

func (ch *traceElementChannel) findUnbufferedPartner() int {
	for routine, trace := range traces {
		if currentIndex[routine] == -1 {
			continue
		}
		elem := trace[currentIndex[routine]]
		switch e := elem.(type) {
		case *traceElementChannel:
			if ch.routine != e.getRoutine() && e.id == ch.id && e.oId == ch.oId {
				return routine
			}
		default:
			continue
		}
	}
	debug.Log("Could not find unbuffered partner for "+ch.toString(), 1)
	return -1
}
