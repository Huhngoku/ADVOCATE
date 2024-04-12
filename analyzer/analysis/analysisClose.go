package analysis

import (
	"analyzer/clock"
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
func checkForCommunicationOnClosedChannel(id int, pos string) {
	// check if there is an earlier send, that could happen concurrently to close
	if analysisCases["sendOnClosed"] && hasSend[id] {
		logging.Debug("Check for possible send on closed channel "+
			strconv.Itoa(id)+" with "+
			mostRecentSend[id].vc.ToString()+" and "+closeData[id].vc.ToString(),
			logging.DEBUG)
		happensBefore := clock.GetHappensBefore(closeData[id].vc, mostRecentSend[id].vc)
		if happensBefore == clock.Concurrent {
			found := "Possible send on closed channel:\n"
			found += "\tclose: " + pos + "\n"
			found += "\tsend : " + mostRecentSend[id].tID
			logging.Result(found, logging.CRITICAL)
		}
	}
	// check if there is an earlier receive, that could happen concurrently to close
	if analysisCases["receiveOnClosed"] && hasReceived[id] {
		logging.Debug("Check for possible receive on closed channel "+
			strconv.Itoa(id)+" with "+
			mostRecentReceive[id].vc.ToString()+" and "+closeData[id].vc.ToString(),
			logging.DEBUG)
		happensBefore := clock.GetHappensBefore(closeData[id].vc, mostRecentReceive[id].vc)
		if happensBefore == clock.Concurrent || happensBefore == clock.Before {
			found := "Possible receive on closed channel:\n"
			found += "\tclose: " + pos + "\n"
			found += "\trecv : " + mostRecentReceive[id].tID
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
	if posOld, ok := closeData[id]; ok {
		found := "Found close on closed channel:\n"
		found += "\tclose: " + pos + "\n"
		found += "\tclose: " + posOld.tID
		logging.Result(found, logging.CRITICAL)
	}
}
