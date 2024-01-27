package analysis

import (
	"analyzer/logging"
	"strconv"
)

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
