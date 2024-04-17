package trace

import (
	"errors"
	"strconv"

	"analyzer/analysis"
	"analyzer/clock"
	"analyzer/logging"
)

// enum for opC
type OpChannel int

const (
	Send OpChannel = iota
	Recv
	Close
)

var waitingReceive = make([]*TraceElementChannel, 0)
var maxOpID = make(map[int]int)

/*
* TraceElementChannel is a trace element for a channel
* MARK: Struct
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
*   sel (*traceElementSelect): The select operation, if the channel operation
*       is part of a select, otherwise nil
*   partner (*TraceElementChannel): The partner of the channel operation
*   tID (string): The id of the trace element, contains the position and the tpre
 */
type TraceElementChannel struct {
	routine int
	tPre    int
	tPost   int
	id      int
	opC     OpChannel
	cl      bool
	oID     int
	qSize   int
	pos     string
	sel     *TraceElementSelect
	partner *TraceElementChannel
	tID     string
	vc      clock.VectorClock
}

/*
* Create a new channel trace element
* MARK: New
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
*   tID (string): The id of the trace element, contains the position and the tpre
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

	var opCInt OpChannel
	switch opC {
	case "S":
		opCInt = Send
	case "R":
		opCInt = Recv
	case "C":
		opCInt = Close
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

	tIDStr := pos + "@" + strconv.Itoa(tPreInt)

	elem := TraceElementChannel{
		routine: routine,
		tPre:    tPreInt,
		tPost:   tPostInt,
		id:      idInt,
		opC:     opCInt,
		cl:      clBool,
		oID:     oIDInt,
		qSize:   qSizeInt,
		pos:     pos,
		tID:     tIDStr,
	}

	// check if partner was already processed, otherwise add to channelWithoutPartner
	if tPostInt != 0 {
		if _, ok := channelWithoutPartner[idInt][oIDInt]; ok {
			elem.partner = channelWithoutPartner[idInt][oIDInt]
			channelWithoutPartner[idInt][oIDInt].partner = &elem
			delete(channelWithoutPartner[idInt], oIDInt)
		} else {
			if _, ok := channelWithoutPartner[idInt]; !ok {
				channelWithoutPartner[idInt] = make(map[int]*TraceElementChannel)
			}

			channelWithoutPartner[idInt][oIDInt] = &elem
		}
	}

	return AddElementToTrace(&elem)
}

// MARK: Getter

/*
* Get the partner of the channel operation
* Returns:
*   *TraceElementChannel: The partner of the channel operation
 */
func (ch *TraceElementChannel) GetPartner() *TraceElementChannel {
	return ch.partner
}

/*
 * Get the id of the element
 * Returns:
 *   int: The id of the element
 */
func (ch *TraceElementChannel) GetID() int {
	return ch.id
}

/*
	* Get the routine of the element
 * Returns:
 *   int: The routine of the element
*/
func (ch *TraceElementChannel) GetRoutine() int {
	return ch.routine
}

/*
 * Get the tpre of the element
 * Returns:
 *   int: The tpre of the element
 */
func (ch *TraceElementChannel) GetTPre() int {
	return ch.tPre
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   float32: The time of the element
 */
func (ch *TraceElementChannel) GetTSort() int {
	if ch.tPost == 0 {
		// if operation was not executed, return tPre. When updating vc, check that tPre is not 0
		return ch.tPre
	}
	return ch.tPost
}

/*
 * Get the position of the operation.
 * Returns:
 *   string: The position of the element
 */
func (ch *TraceElementChannel) GetPos() string {
	return ch.pos
}

/*
 * Get the tID of the element.
 * Returns:
 *   string: The tID of the element
 */
func (ch *TraceElementChannel) GetTID() string {
	return ch.tID
}

/*
 * Get the oID of the element
 * Returns:
 *   int: The oID of the element
 */
func (ch *TraceElementChannel) GetOID() int {
	return ch.oID
}

/*
 * Check if the channel operation is buffered
 * Returns:
 *   bool: Whether the channel operation is buffered
 */
func (ch *TraceElementChannel) IsBuffered() bool {
	return ch.qSize != 0
}

/*
 * Get the type of the operation
 * Returns:
 *   OpChannel: The type of the operation
 */
func (ch *TraceElementChannel) Operation() OpChannel {
	return ch.opC
}

/*
 * Get the vector clock of the element
 * Returns:
 *   VectorClock: The vector clock of the element
 */
func (ch *TraceElementChannel) GetVC() clock.VectorClock {
	return ch.vc
}

/*
 * Get the tpost of the element
 * Returns:
 *   int: The tpost of the element
 */
func (ch *TraceElementChannel) getTpost() int {
	return ch.tPost
}

// MARK: Setter

/*
* Set the tpre of the element.
* Args:
 *   tPre (int): The tpre of the element
*/
func (ch *TraceElementChannel) SetTPre(tPre int) {
	ch.tPre = tPre
	if ch.tPost != 0 && ch.tPost < tPre {
		ch.tPost = tPre
	}
}

/*
 * Set the partner of the channel operation
 * Args:
 *   partner (*TraceElementChannel): The partner of the channel operation
 */
func (ch *TraceElementChannel) SetPartner(partner *TraceElementChannel) {
	ch.partner = partner
}

/*
 * Set the tpost of the element.
 * Args:
 *   tPost (int): The tpost of the element
 */
func (ch *TraceElementChannel) SetTPost(tPost int) {
	ch.tPost = tPost
}

/*
 * Set the timer, that is used for the sorting of the trace
 * Args:
 *   tSort (int): The timer of the element
 */
func (ch *TraceElementChannel) SetTSort(tpost int) {
	ch.SetTPre(tpost)
	ch.tPost = tpost
}

/*
 * Set the timer, that is used for the sorting of the trace, only if the original
 * value was not 0
 * Args:
 *   tSort (int): The timer of the element
 */
func (ch *TraceElementChannel) SetTSortWithoutNotExecuted(tSort int) {
	ch.SetTPre(tSort)
	if ch.tPost != 0 {
		ch.tPost = tSort
	}
}

/*
 * Set the oID of the element
 * Args:
 *   oID (int): The oID of the element
 */
func (ch *TraceElementChannel) SetOID(oID int) {
	ch.oID = oID
}

// MARK: ToString

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (ch *TraceElementChannel) ToString() string {
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
func (ch *TraceElementChannel) toStringSep(sep string, pos bool) string {
	res := "C" + sep
	res += strconv.Itoa(ch.tPre) + sep + strconv.Itoa(ch.tPost) + sep
	res += strconv.Itoa(ch.id) + sep

	switch ch.opC {
	case Send:
		res += "S"
	case Recv:
		res += "R"
	case Close:
		res += "C"
	default:
		panic("Unknown channel operation" + strconv.Itoa(int(ch.opC)))
	}

	res += sep + "f"

	res += sep + strconv.Itoa(ch.oID)
	res += sep + strconv.Itoa(ch.qSize)
	if pos {
		res += sep + ch.pos
	}
	return res
}

/*
 * Update and calculate the vector clock of the element
 * MARK: Vector Clock
 */
func (ch *TraceElementChannel) updateVectorClock() {
	// hold back receive operations, until the send operation is processed
	for _, elem := range waitingReceive {
		if elem.oID <= maxOpID[ch.id] {
			waitingReceive = waitingReceive[1:]
			elem.updateVectorClock()
		}
	}
	if ch.IsBuffered() && ch.tPost != 0 {
		if ch.opC == Send {
			maxOpID[ch.id] = ch.oID
		} else if ch.opC == Recv {
			logging.Debug("Holding back", logging.INFO)
			if ch.oID > maxOpID[ch.id] && !ch.cl {
				waitingReceive = append(waitingReceive, ch)
				return
			}
		}
	}

	if !ch.IsBuffered() { // unbuffered channel
		switch ch.opC {
		case Send:
			partner := ch.findPartner()
			if partner != -1 {
				logging.Debug("Update vector clock of channel operation: "+
					traces[partner][currentIndex[partner]].ToString(),
					logging.DEBUG)
				pos := traces[partner][currentIndex[partner]].(*TraceElementChannel).tID
				analysis.Unbuffered(ch.routine, partner, ch.id, ch.tID,
					pos, currentVCHb, ch.tPost)
				// advance index of receive routine, send routine is already advanced
				increaseIndex(partner)
			} else {
				if ch.cl { // recv on closed channel
					logging.Debug("Update vector clock of channel operation: "+
						ch.ToString(), logging.DEBUG)
					analysis.RecvC(ch.routine, ch.id, ch.tID,
						currentVCHb, ch.tPost, false)
				} else {
					logging.Debug("Could not find partner for "+ch.tID, logging.INFO)
				}
			}

		case Recv: // should not occur, but better save than sorry
			partner := ch.findPartner()
			if partner != -1 {
				logging.Debug("Update vector clock of channel operation: "+
					traces[partner][currentIndex[partner]].ToString(), logging.DEBUG)
				tID := traces[partner][currentIndex[partner]].(*TraceElementChannel).tID
				analysis.Unbuffered(partner, ch.routine, ch.id, tID,
					ch.tID, currentVCHb, ch.tPost)
				// advance index of receive routine, send routine is already advanced
				increaseIndex(partner)
			} else {
				if ch.cl { // recv on closed channel
					logging.Debug("Update vector clock of channel operation: "+
						ch.ToString(), logging.DEBUG)
					analysis.RecvC(ch.routine, ch.id, ch.tID,
						currentVCHb, ch.tPost, false)
				} else {
					logging.Debug("Could not find partner for "+ch.tID, logging.INFO)
				}
			}
		case Close:
			analysis.Close(ch.routine, ch.id, ch.tID, currentVCHb, ch.tPost)
		default:
			err := "Unknown operation: " + ch.ToString()
			logging.Debug(err, logging.ERROR)
		}
	} else { // buffered channel
		switch ch.opC {
		case Send:
			logging.Debug("Update vector clock of channel operation: "+
				ch.ToString(), logging.DEBUG)
			analysis.Send(ch.routine, ch.id, ch.oID, ch.qSize, ch.tID,
				currentVCHb, fifo, ch.tPost)
		case Recv:
			if ch.cl { // recv on closed channel
				logging.Debug("Update vector clock of channel operation: "+
					ch.ToString(), logging.DEBUG)
				analysis.RecvC(ch.routine, ch.id, ch.tID, currentVCHb, ch.tPost, true)
			} else {
				logging.Debug("Update vector clock of channel operation: "+
					ch.ToString(), logging.DEBUG)
				analysis.Recv(ch.routine, ch.id, ch.oID, ch.qSize, ch.tID,
					currentVCHb, fifo, ch.tPost)
			}
		case Close:
			logging.Debug("Update vector clock of channel operation: "+
				ch.ToString(), logging.DEBUG)
			analysis.Close(ch.routine, ch.id, ch.tID, currentVCHb, ch.tPost)
		default:
			err := "Unknown operation: " + ch.ToString()
			logging.Debug(err, logging.ERROR)
		}
	}

	ch.vc = currentVCHb[ch.routine].Copy()
	if ch.partner != nil {
		ch.partner.vc = currentVCHb[ch.partner.routine].Copy()
	}
}

/*
 * Find the partner of the channel operation
 * MARK: Partner
 * Returns:
 *   int: The routine id of the partner
 */
func (ch *TraceElementChannel) findPartner() int {
	// return -1 if closed by channel
	if ch.cl {
		return -1
	}

	for routine, trace := range traces {
		if currentIndex[routine] == -1 {
			continue
		}
		// if routine == ch.routine {
		// 	continue
		// }
		elem := trace[currentIndex[routine]]
		switch e := elem.(type) {
		case *TraceElementChannel:
			if e.id == ch.id && e.oID == ch.oID {
				return routine
			}
		case *TraceElementSelect:
			if e.chosenCase.tPost != 0 &&
				e.chosenCase.oID == ch.id &&
				e.chosenCase.oID == ch.oID {
				return routine
			}
		default:
			continue
		}
	}
	return -1
}
