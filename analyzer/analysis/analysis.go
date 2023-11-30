package analysis

import (
	"analyzer/logging"
	"strconv"
)

// vc of close on channel
var closeVC map[int]VectorClock = make(map[int]VectorClock)
var closePos map[int]string = make(map[int]string)

// last send and receive on channel
var lastSend map[int]VectorClock = make(map[int]VectorClock)
var lastRecv map[int]VectorClock = make(map[int]VectorClock)

// last receive for each routine and each channel
var lastRecvRoutine map[int]map[int]VectorClock = make(map[int]map[int]VectorClock)
var lastRecvRoutinePos map[int]map[int]string = make(map[int]map[int]string)

// most recent send, used for detection of send on closed
var hasSend map[int]bool = make(map[int]bool)
var mostRecentSend map[int]VectorClock = make(map[int]VectorClock)
var mostRecentSendPosition map[int]string = make(map[int]string)

// most recent send, used for detection of received on closed
var hasReceived map[int]bool = make(map[int]bool)
var mostRecentReceive map[int]VectorClock = make(map[int]VectorClock)
var mostRecentReceivePosition map[int]string = make(map[int]string)

/*
Check if a send or receive on a closed channel is possible
It it is possible, print a warning or error
Args:

	id (int): the id of the channel
	pos (string): the position of the close in the program
*/
func checkForPotentialCommunicationOnClosedChannel(id int, pos string) {
	// check if there is an earlier send, that could happen concurrently to close
	if hasSend[id] {
		logging.Debug("Check for possible send on closed channel "+
			strconv.Itoa(id)+" with "+
			mostRecentSend[id].ToString()+" and "+closeVC[id].ToString(),
			logging.DEBUG)
		happensBefore := GetHappensBefore(closeVC[id], mostRecentSend[id])
		if happensBefore == Concurrent {
			found := "Possible send on closed channel:\n"
			found += "\tclose: " + pos + "\n"
			found += "\tsend : " + mostRecentSendPosition[id]
			logging.Result(found, logging.CRITICAL)
		}
	}
	// check if there is an earlier receive, that could happen concurrently to close
	if hasReceived[id] {
		logging.Debug("Check for possible receive on closed channel "+
			strconv.Itoa(id)+" with "+
			mostRecentReceive[id].ToString()+" and "+closeVC[id].ToString(),
			logging.DEBUG)
		happensBefore := GetHappensBefore(closeVC[id], mostRecentReceive[id])
		if happensBefore == Concurrent || happensBefore == Before {
			found := "Possible receive on closed channel:\n"
			found += "\tclose: " + pos + "\n"
			found += "\trecv : " + mostRecentReceivePosition[id]
			logging.Result(found, logging.WARNING)
		}
	}

}

func foundReceiveOnClosedChannel(posClose string, posRecv string) {
	found := "Found receive on closed channel:\n"
	found += "\tclose: " + posClose + "\n"
	found += "\trecv : " + posRecv
	logging.Result(found, logging.WARNING)
}

func checkForConcurrentRecv(routine int, id int, pos string, vc map[int]VectorClock) {
	if _, ok := lastRecvRoutine[routine]; !ok {
		lastRecvRoutine[routine] = make(map[int]VectorClock)
		lastRecvRoutinePos[routine] = make(map[int]string)
	}

	lastRecvRoutine[routine][id] = vc[routine].Copy()
	lastRecvRoutinePos[routine][id] = pos

	for r, elem := range lastRecvRoutine {
		if r == routine {
			continue
		}

		if elem[id].clock == nil {
			continue
		}

		happensBefore := GetHappensBefore(elem[id], vc[routine])
		if happensBefore == Concurrent {
			found := "Found concurrent Recv on same channel:\n"
			found += "\trecv: " + pos + "\n"
			found += "\trecv : " + lastRecvRoutinePos[r][id]
			logging.Result(found, logging.CRITICAL)
		}
	}
}

/*
 * Check for a close on a closed channel.
 * Must be called, before the current close operation is added to closePos
 * Args:
 * 	id (int): the id of the channel
 * 	pos (string): the position of the close in the program
 */
func checkForClosedOnClosed(id int, pos string) {
	if posOld, ok := closePos[id]; ok {
		found := "Found close on closed channel:\n"
		found += "\tclose: " + pos + "\n"
		found += "\tclose: " + posOld
		logging.Result(found, logging.CRITICAL)
	}
}
