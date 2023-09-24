package trace

import (
	"errors"
	"math"
	"strconv"

	"analyzer/debug"
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
*   pre (*traceElementPre): The pre element of the channel operation
 */
type traceElementChannel struct {
	routine int
	tpre    int
	tpost   int
	vpre    vectorClock
	vpost   vectorClock
	id      int
	opC     opChannel
	cl      bool
	oId     int
	qSize   int
	pos     string
	sel     *traceElementSelect
	partner *traceElementChannel
	pre     *traceElementPre
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
func addTraceElementChannel(routine int, numberOfRoutines int, tpre string,
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
		vpre:    newVectorClock(numberOfRoutines),
		vpost:   newVectorClock(numberOfRoutines),
		id:      id_int,
		opC:     opC_int,
		cl:      cl_bool,
		oId:     oId_int,
		qSize:   qSize_int,
		pos:     pos,
	}

	elem_pre := traceElementPre{
		elem:     &elem,
		elemType: Chan,
	}

	elem.pre = &elem_pre

	res1 := addElementToTrace(&elem_pre)
	res2 := addElementToTrace(&elem)
	return errors.Join(res1, res2)
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
func (ch *traceElementChannel) getVpre() *vectorClock {
	return &ch.vpre
}

/*
 * Get the vector clock at the end of the event
 * Returns:
 *   vectorClock: The vector clock at the end of the event
 */
func (ch *traceElementChannel) getVpost() *vectorClock {
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
 * Function to find communication partner for channel operations
 * This includes send/receive pairs with or without select, as well as
 * close operations for send/receive/select, which where
 * executed because of a close operation on the channel.
 * If a partner is found, the partner field of the element is set.
 */
func (ch *traceElementChannel) findPartner() {
	if ch.tpost == 0 { // if tpost is 0, the operation was not finished
		debug.Log("Channel operation "+ch.toString()+" was not finished", 3)
		return
	}

	if ch.opC == close { // close operation has no partner
		debug.Log("Add close operation "+ch.toString()+" to closeOperations", 3)
		closeOperations = append(closeOperations, ch)
	}

	// check if partner is already in channelOperations
	i := 0
	for _, partner := range channelOperations {
		// check for send receive
		if ch.id == partner.id && ch.opC != partner.opC && ch.oId == partner.oId {
			debug.Log(
				"Found partner for channel operation"+partner.toString()+" <-> "+ch.toString(),
				3,
			)
			ch.partner = partner
			partner.partner = ch
			break
		}

		// check new close
		if ch.opC == close {
			debug.Log("Check for new close operation "+ch.toString(), 3)
			if ch.id == partner.id && partner.cl {
				partner.partner = ch
				break
			}
		}
		i++
	}

	// remove the partner element, if an partner was found
	if ch.partner != nil {
		debug.Log("Partner found. Remove "+ch.partner.toString()+" from channelOperations", 3)
		channelOperations = append(channelOperations[:i], channelOperations[i+1:]...)
	}

	// check if partner is already in closeOperations
	for _, partner := range closeOperations {
		debug.Log("Check for close partner "+partner.toString(), 3)
		if ch.id == partner.id {
			if ch.opC == close && partner.cl {
				debug.Log("Found close partner "+partner.toString()+" for "+ch.toString(), 3)
				partner.partner = ch
			} else if ch.cl && partner.opC == close {
				debug.Log("Found close partner "+partner.toString()+" for "+ch.toString(), 3)
				ch.partner = partner
			}
		}
	}

	// if partner is not found, add to channelOperations
	if ch.partner == nil && ch.opC != close {
		debug.Log("No partner found. Add "+ch.toString()+" to channelOperations", 3)
		channelOperations = append(channelOperations, ch)
	}
}

/*
 * Update and calculate the vector clock of the element
 * Args:
 *   vc (vectorClock): The current vector clocks
 * TODO: implement
 */
func (ch *traceElementChannel) calculateVectorClock(vc *[]vectorClock) {

}

/*
 * Function to check if there are channel operations without partner
 */
func checkChannelOperations() {
	for _, elem := range channelOperations {
		debug.Log("Channel operation "+elem.toString()+" has no partner", 1)
	}
}
