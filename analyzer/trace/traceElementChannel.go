package trace

import (
	"errors"
	"math"
	"strconv"

	"analyzer/analysis"
	"analyzer/logging"
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
*   numberOfRoutines (int): The number of routines in the trace
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

/*
 * Update and calculate the vector clock of the element
 */
func (ch *traceElementChannel) updateVectorClock() {
	// hold back receive operations, until the send operation is processed
	for _, elem := range waitingReceive {
		if elem.oID <= maxOpID[ch.id] {
			waitingReceive = waitingReceive[1:]
			elem.updateVectorClock()
		}
	}
	if ch.qSize != 0 {
		if ch.opC == send {
			maxOpID[ch.id] = ch.oID
		} else if ch.opC == recv {
			logging.Debug("Holding back", logging.INFO)
			if ch.oID > maxOpID[ch.id] && !ch.cl {
				waitingReceive = append(waitingReceive, ch)
				return
			}
		}
	}

	if ch.qSize == 0 { // unbuffered channel
		switch ch.opC {
		case send:
			partner := ch.findUnbufferedPartner()
			if partner != -1 {
				logging.Debug("Update vector clock of channel operation: "+
					traces[partner][currentIndex[partner]].toString(),
					logging.DEBUG)
				pos := traces[partner][currentIndex[partner]].(*traceElementChannel).pos
				analysis.Unbuffered(ch.routine, partner, ch.id, ch.pos,
					pos, currentVectorClocks)
				// advance index of receive routine, send routine is already advanced
				increaseIndex(partner)
			} else {
				if ch.cl { // recv on closed channel
					logging.Debug("Update vector clock of channel operation: "+
						ch.toString(), logging.DEBUG)
					analysis.RecvC(ch.routine, ch.id, ch.pos,
						currentVectorClocks)
				} else {
					logging.Debug("Could not find partner for "+ch.pos, logging.INFO)
				}
			}
		case recv: // should not occur, but better save than sorry
			partner := ch.findUnbufferedPartner()
			if partner != -1 {
				logging.Debug("Update vector clock of channel operation: "+
					traces[partner][currentIndex[partner]].toString(), logging.DEBUG)
				pos := traces[partner][currentIndex[partner]].(*traceElementChannel).pos
				analysis.Unbuffered(partner, ch.routine, ch.id, pos,
					ch.pos, currentVectorClocks)
				// advance index of receive routine, send routine is already advanced
				increaseIndex(partner)
			} else {
				if ch.cl { // recv on closed channel
					logging.Debug("Update vector clock of channel operation: "+
						ch.toString(), logging.DEBUG)
					analysis.RecvC(ch.routine, ch.id, ch.pos,
						currentVectorClocks)
				} else {
					logging.Debug("Could not find partner for "+ch.pos, logging.INFO)
				}
			}
		case close:
			analysis.Close(ch.routine, ch.id, ch.pos, currentVectorClocks)
		default:
			err := "Unknown operation: " + ch.toString()
			logging.Debug(err, logging.ERROR)
		}
	} else { // buffered channel
		switch ch.opC {
		case send:
			logging.Debug("Update vector clock of channel operation: "+
				ch.toString(), logging.DEBUG)
			analysis.Send(ch.routine, ch.id, ch.oID, ch.qSize, ch.pos,
				currentVectorClocks, fifo)
		case recv:
			if ch.cl { // recv on closed channel
				logging.Debug("Update vector clock of channel operation: "+
					ch.toString(), logging.DEBUG)
				analysis.RecvC(ch.routine, ch.id, ch.pos, currentVectorClocks)
			} else {
				logging.Debug("Update vector clock of channel operation: "+
					ch.toString(), logging.DEBUG)
				analysis.Recv(ch.routine, ch.id, ch.oID, ch.qSize, ch.pos,
					currentVectorClocks, fifo)
			}
		case close:
			logging.Debug("Update vector clock of channel operation: "+
				ch.toString(), logging.DEBUG)
			analysis.Close(ch.routine, ch.id, ch.pos, currentVectorClocks)
		default:
			err := "Unknown operation: " + ch.toString()
			logging.Debug(err, logging.ERROR)
		}
	}
}

func (ch *traceElementChannel) findUnbufferedPartner() int {
	// return -1 if closed by channel
	if ch.cl {
		return -1
	}

	for routine, trace := range traces {
		if currentIndex[routine] == -1 {
			continue
		}
		if routine == ch.routine {
			continue
		}
		elem := trace[currentIndex[routine]]
		switch e := elem.(type) {
		case *traceElementChannel:
			if e.id == ch.id && e.oID == ch.oID {
				return routine
			}
		case *traceElementSelect:
			if e.chosenCase.oID == ch.id &&
				e.chosenCase.oID == ch.oID {
				return routine
			}
		default:
			continue
		}
	}
	return -1
}
