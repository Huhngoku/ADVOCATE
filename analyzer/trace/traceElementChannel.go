package trace

import (
	"errors"
	"strconv"
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
*   tpre (int): The timestamp at the start of the event
*   tpost (int): The timestamp at the end of the event
*   id (int): The id of the channel
*   opC (int, enum): The operation on the channel
*   exec (int, enum): The execution status of the operation
*   oId (int): The id of the other communication
*   qSize (int): The size of the channel queue
*   qCountPre (int): The number of elements in the queue before the operation
*   qCountPost (int): The number of elements in the queue after the operation
*   pos (string): The position of the channel operation in the code
 */
type traceElementChannel struct {
	tpre       int
	tpost      int
	id         int
	opC        opChannel
	exec       bool
	oId        int
	qSize      int
	qCountPre  int
	qCountPost int
	pos        string
}

/*
* Create a new channel trace element
* Args:
*   routine (int): The routine id
*   tpre (string): The timestamp at the start of the event
*   tpost (string): The timestamp at the end of the event
*   id (string): The id of the channel
*   opC (string): The operation on the channel
*   exec (string): The execution status of the operation
*   oId (string): The id of the other communication
*   qSize (string): The size of the channel queue
*   qCountPre (string): The number of elements in the queue before the operation
*   qCountPost (string): The number of elements in the queue after the operation
*   pos (string): The position of the channel operation in the code
 */
func AddTraceElementChannel(routine int, tpre string, tpost string, id string,
	opC string, exec string, oId string, qSize string, qCountPre string,
	qCountPost string, pos string) error {
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

	exec_bool := false
	if exec == "e" {
		exec_bool = true
	}

	oId_int, err := strconv.Atoi(oId)
	if err != nil {
		return errors.New("oId is not an integer")
	}

	qSize_int, err := strconv.Atoi(qSize)
	if err != nil {
		return errors.New("qSize is not an integer")
	}

	qCountPre_int, err := strconv.Atoi(qCountPre)
	if err != nil {
		return errors.New("qCountPre is not an integer")
	}

	qCountPost_int, err := strconv.Atoi(qCountPost)
	if err != nil {
		return errors.New("qCountPost is not an integer")
	}

	elem := traceElementChannel{tpre_int, tpost_int, id_int, opC_int, exec_bool,
		oId_int, qSize_int, qCountPre_int, qCountPost_int, pos}

	return addElementToTrace(routine, elem)
}

func (elem traceElementChannel) getSimpleString() string {
	return "A" + strconv.Itoa(elem.tpre) + "," + strconv.Itoa(elem.tpost) + "," +
		strconv.Itoa(elem.id) + "," + strconv.Itoa(int(elem.opC)) + "," +
		strconv.Itoa(elem.oId)
}
