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
		for routine, mrs := range mostRecentSend {
			logging.Debug("Check for possible send on closed channel "+
				strconv.Itoa(id)+" with "+
				mrs[id].Vc.ToString()+" and "+closeData[id].Vc.ToString(),
				logging.DEBUG)

			happensBefore := clock.GetHappensBefore(closeData[id].Vc, mrs[id].Vc)
			if mrs[id].TID != "" && happensBefore == clock.Concurrent {

				file1, line1, tPre1, err := infoFromTID(mrs[id].TID) // send
				if err != nil {
					logging.Debug(err.Error(), logging.ERROR)
					return
				}

				file2, line2, tPre2, err := infoFromTID(pos) // close
				if err != nil {
					logging.Debug(err.Error(), logging.ERROR)
					return
				}

				arg1 := logging.TraceElementResult{ // send
					RoutineID: routine,
					ObjID:     id,
					TPre:      tPre1,
					ObjType:   "CS",
					File:      file1,
					Line:      line1,
				}

				arg2 := logging.TraceElementResult{ // close
					RoutineID: closeData[id].Routine,
					ObjID:     id,
					TPre:      tPre2,
					ObjType:   "CC",
					File:      file2,
					Line:      line2,
				}

				logging.Result(logging.CRITICAL, logging.PSendOnClosed,
					"send", []logging.ResultElem{arg1}, "close", []logging.ResultElem{arg2})
			}
		}
	}
	// check if there is an earlier receive, that could happen concurrently to close
	if analysisCases["receiveOnClosed"] && hasReceived[id] {
		for routine, mrr := range mostRecentReceive {
			logging.Debug("Check for possible receive on closed channel "+
				strconv.Itoa(id)+" with "+
				mrr[id].Vc.ToString()+" and "+closeData[id].Vc.ToString(),
				logging.DEBUG)

			happensBefore := clock.GetHappensBefore(closeData[id].Vc, mrr[id].Vc)
			if mrr[id].TID != "" && (happensBefore == clock.Concurrent || happensBefore == clock.Before) {

				file1, line1, tPre1, err := infoFromTID(mrr[id].TID) // recv
				if err != nil {
					logging.Debug(err.Error(), logging.ERROR)
					return
				}

				file2, line2, tPre2, err := infoFromTID(pos) // close
				if err != nil {
					logging.Debug(err.Error(), logging.ERROR)
					return
				}

				arg1 := logging.TraceElementResult{ // recv
					RoutineID: routine,
					ObjID:     id,
					TPre:      tPre1,
					ObjType:   "CR",
					File:      file1,
					Line:      line1,
				}

				arg2 := logging.TraceElementResult{ // close
					RoutineID: closeData[id].Routine,
					ObjID:     id,
					TPre:      tPre2,
					ObjType:   "CC",
					File:      file2,
					Line:      line2,
				}

				logging.Result(logging.WARNING, logging.PRecvOnClosed,
					"recv", []logging.ResultElem{arg1}, "close", []logging.ResultElem{arg2})
			}
		}
	}

}

func foundSendOnClosedChannel(routineID int, id int, posSend string) {
	if _, ok := closeData[id]; !ok {
		return
	}

	posClose := closeData[id].TID
	if posClose == "" || posSend == "" || posClose == "\n" || posSend == "\n" {
		return
	}

	file1, line1, tPre1, err := infoFromTID(posSend)
	if err != nil {
		logging.Debug(err.Error(), logging.ERROR)
		return
	}

	file2, line2, tPre2, err := infoFromTID(posClose)
	if err != nil {
		logging.Debug(err.Error(), logging.ERROR)
		return
	}

	arg1 := logging.TraceElementResult{ // send
		RoutineID: routineID,
		ObjID:     id,
		TPre:      tPre1,
		ObjType:   "CS",
		File:      file1,
		Line:      line1,
	}

	arg2 := logging.TraceElementResult{ // close
		RoutineID: closeData[id].Routine,
		ObjID:     id,
		TPre:      tPre2,
		ObjType:   "CC",
		File:      file2,
		Line:      line2,
	}

	logging.Result(logging.CRITICAL, logging.ASendOnClosed,
		"send", []logging.ResultElem{arg1}, "close", []logging.ResultElem{arg2})

}

func foundReceiveOnClosedChannel(routineID int, id int, posRecv string) {
	if _, ok := closeData[id]; !ok {
		return
	}

	posClose := closeData[id].TID
	if posClose == "" || posRecv == "" || posClose == "\n" || posRecv == "\n" {
		return
	}

	file1, line1, tPre1, err := infoFromTID(posRecv)
	if err != nil {
		logging.Debug(err.Error(), logging.ERROR)
		return
	}

	file2, line2, tPre2, err := infoFromTID(posClose)
	if err != nil {
		logging.Debug(err.Error(), logging.ERROR)
		return
	}

	arg1 := logging.TraceElementResult{ // recv
		RoutineID: routineID,
		ObjID:     id,
		TPre:      tPre1,
		ObjType:   "CR",
		File:      file1,
		Line:      line1,
	}

	arg2 := logging.TraceElementResult{ // close
		RoutineID: closeData[id].Routine,
		ObjID:     id,
		TPre:      tPre2,
		ObjType:   "CC",
		File:      file2,
		Line:      line2,
	}

	logging.Result(logging.WARNING, logging.ARecvOnClosed,
		"recv", []logging.ResultElem{arg1}, "close", []logging.ResultElem{arg2})
}

/*
 * Check for a close on a closed channel.
 * Must be called, before the current close operation is added to closePos
 * Args:
 *  routineID (int): the id of the routine
 * 	id (int): the id of the channel
 * 	pos (string): the position of the close in the program
 */
func checkForClosedOnClosed(routineID int, id int, pos string) {
	if oldClose, ok := closeData[id]; ok {
		if oldClose.TID == "" || oldClose.TID == "\n" || pos == "" || pos == "\n" {
			return
		}

		file1, line1, tPre1, err := infoFromTID(oldClose.TID)
		if err != nil {
			logging.Debug(err.Error(), logging.ERROR)
			return
		}

		file2, line2, tPre2, err := infoFromTID(oldClose.TID)
		if err != nil {
			logging.Debug(err.Error(), logging.ERROR)
			return
		}

		arg1 := logging.TraceElementResult{
			RoutineID: routineID,
			ObjID:     id,
			TPre:      tPre1,
			ObjType:   "CC",
			File:      file1,
			Line:      line1,
		}

		arg2 := logging.TraceElementResult{
			RoutineID: oldClose.Routine,
			ObjID:     id,
			TPre:      tPre2,
			ObjType:   "CC",
			File:      file2,
			Line:      line2,
		}

		logging.Result(logging.CRITICAL, logging.ACloseOnClosed,
			"close", []logging.ResultElem{arg1}, "close", []logging.ResultElem{arg2})
	}
}
