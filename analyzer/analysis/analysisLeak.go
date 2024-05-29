package analysis

import (
	"analyzer/clock"
	"analyzer/logging"
)

/*
 * Run for channel operation without a post event. Check if the operation has
 * a possible communication partner in mostRecentSend, mostRecentReceive or closeData.
 * If so, add an error or warning to the result.
 * If not, add to leakingChannels, for later check.
 * MARK: Channel Stuck
 * Args:
 *   routineID (int): The routine id
 *   objID (int): The channel id
 *   vc (VectorClock): The vector clock of the operation
 *   tID (string): The trace id
 *   opType (int): An identifier for the type of the operation (send = 0, recv = 1)
 *   buffered (bool): If the channel is buffered
 */
func CheckForLeakChannelStuck(routineID int, objID int, vc clock.VectorClock, tID string, opType int,
	buffered bool) {
	logging.Debug("Checking channel for for leak channel", logging.INFO)

	if objID == -1 {
		objType := "C"
		if opType == 0 {
			objType += "S"
		} else {
			objType += "R"
		}

		file, line, tPre, err := infoFromTID(tID)
		if err != nil {
			logging.Debug("Error in infoFromTID", logging.ERROR)
			return
		}

		arg1 := logging.TraceElementResult{
			RoutineID: routineID, ObjID: objID, TPre: tPre, ObjType: objType, File: file, Line: line}

		logging.Result(logging.CRITICAL, logging.LNil,
			"Channel", []logging.ResultElem{arg1}, "", []logging.ResultElem{})

		return
	}

	// if !buffered {
	foundPartner := false

	if opType == 0 { // send
		for partnerRout, mrr := range mostRecentReceive {
			if _, ok := mrr[objID]; ok {
				if clock.GetHappensBefore(mrr[objID].Vc, vc) == clock.Concurrent {

					var bugType logging.ResultType = logging.LUnbufferedWith
					if buffered {
						bugType = logging.LBufferedWith
					}

					file1, line1, tPre1, err := infoFromTID(tID)
					if err != nil {
						logging.Debug("Error in infoFromTID", logging.ERROR)
						return
					}
					file2, line2, tPre2, err := infoFromTID(mrr[objID].TID)
					if err != nil {
						logging.Debug("Error in infoFromTID", logging.ERROR)
						return
					}

					arg1 := logging.TraceElementResult{
						RoutineID: routineID, ObjID: objID, TPre: tPre1, ObjType: "CS", File: file1, Line: line1}
					arg2 := logging.TraceElementResult{
						RoutineID: partnerRout, ObjID: objID, TPre: tPre2, ObjType: "CR", File: file2, Line: line2}

					logging.Result(logging.CRITICAL, bugType,
						"channel", []logging.ResultElem{arg1}, "partner", []logging.ResultElem{arg2})

					foundPartner = true
				}
			}
		}
	} else if opType == 1 { // recv
		for partnerRout, mrs := range mostRecentSend {
			if _, ok := mrs[objID]; ok {
				if clock.GetHappensBefore(mrs[objID].Vc, vc) == clock.Concurrent {

					var bugType logging.ResultType = logging.LUnbufferedWith
					if buffered {
						bugType = logging.LBufferedWith
					}

					file1, line1, tPre1, err1 := infoFromTID(tID)
					if err1 != nil {
						logging.Debug("Error in infoFromTID", logging.ERROR)
						return
					}
					file2, line2, tPre2, err2 := infoFromTID(mrs[objID].TID)
					if err2 != nil {
						logging.Debug("Error in infoFromTID", logging.ERROR)
						return
					}

					arg1 := logging.TraceElementResult{
						RoutineID: routineID, ObjID: objID, TPre: tPre1, ObjType: "CR", File: file1, Line: line1}
					arg2 := logging.TraceElementResult{
						RoutineID: partnerRout, ObjID: objID, TPre: tPre2, ObjType: "CS", File: file2, Line: line2}

					logging.Result(logging.CRITICAL, bugType,
						"channel", []logging.ResultElem{arg1}, "partner", []logging.ResultElem{arg2})

					foundPartner = true
				}
			}
		}

	}

	if !foundPartner {
		leakingChannels[objID] = append(leakingChannels[objID], VectorClockTID2{routineID, objID, vc, tID, opType, -1, buffered, false})
	}
}

/*
 * Run for channel operation with a post event. Check if the operation would be
 * possible communication partner for a stuck operation in leakingChannels.
 * If so, add an error or warning to the result and remove the stuck operation.
 * MARK: ChannelRun
 * Args:
 *   routineID (int): The routine id
 *   objID (int): The channel id
 *   vc (VectorClock): The vector clock of the operation
 *   tID (string): The trace id
 *   opType (int): An identifier for the type of the operation (send = 0, recv = 1, close = 2)
 *   buffered (bool): If the channel is buffered
 */
func CheckForLeakChannelRun(routineID int, objID int, vcTID VectorClockTID, opType int, buffered bool) bool {
	logging.Debug("Checking channel for for leak channels", logging.INFO)
	res := false
	if opType == 0 || opType == 2 { // send or close
		for i, vcTID2 := range leakingChannels[objID] {
			if vcTID2.val != 1 {
				continue
			}

			if clock.GetHappensBefore(vcTID2.vc, vcTID.Vc) == clock.Concurrent {
				var bugType logging.ResultType = logging.LUnbufferedWith
				if buffered {
					bugType = logging.LBufferedWith
				}

				file1, line1, tPre1, err1 := infoFromTID(vcTID2.tID) // leaking
				if err1 != nil {
					logging.Debug("Error in infoFromTID", logging.ERROR)
					return res
				}
				file2, line2, tPre2, err2 := infoFromTID(vcTID.TID) // partner
				if err2 != nil {
					logging.Debug("Error in infoFromTID", logging.ERROR)
					return res
				}

				objType := "C"
				if opType == 0 {
					objType += "S"
				} else {
					objType += "C"
				}

				arg1 := logging.TraceElementResult{
					RoutineID: routineID, ObjID: objID, TPre: tPre1, ObjType: "CR", File: file1, Line: line1}
				arg2 := logging.TraceElementResult{
					RoutineID: vcTID2.routine, ObjID: objID, TPre: tPre2, ObjType: objType, File: file2, Line: line2}

				logging.Result(logging.CRITICAL, bugType,
					"channel", []logging.ResultElem{arg1}, "partner", []logging.ResultElem{arg2})

				res = true

				// remove the stuck operation from the list. If it is a select, remove all operations with the same val
				if vcTID2.val == -1 {
					leakingChannels[objID] = append(leakingChannels[objID][:i], leakingChannels[objID][i+1:]...)
				} else {
					for j, vcTID3 := range leakingChannels[objID] {
						if vcTID3.val == vcTID2.val {
							leakingChannels[objID] = append(leakingChannels[objID][:j], leakingChannels[objID][j+1:]...)
						}
					}
				}
			}
		}
	} else if opType == 1 { // recv
		for i, vcTID2 := range leakingChannels[objID] {
			objType := "C"
			if vcTID2.val == 0 {
				objType += "S"
			} else if vcTID2.val == 2 {
				objType += "C"
			} else {
				continue
			}

			if clock.GetHappensBefore(vcTID2.vc, vcTID.Vc) == clock.Concurrent {

				var bugType logging.ResultType = logging.LUnbufferedWith
				if buffered {
					bugType = logging.LBufferedWith
				}

				file1, line1, tPre1, err1 := infoFromTID(vcTID2.tID) // leaking
				if err1 != nil {
					logging.Debug("Error in infoFromTID", logging.ERROR)
					return res
				}
				file2, line2, tPre2, err2 := infoFromTID(vcTID.TID) // partner
				if err2 != nil {
					logging.Debug("Error in infoFromTID", logging.ERROR)
					return res
				}

				arg1 := logging.TraceElementResult{
					RoutineID: routineID, ObjID: objID, TPre: tPre1, ObjType: objType, File: file1, Line: line1}
				arg2 := logging.TraceElementResult{
					RoutineID: vcTID2.routine, ObjID: objID, TPre: tPre2, ObjType: "CR", File: file2, Line: line2}

				logging.Result(logging.CRITICAL, bugType,
					"channel", []logging.ResultElem{arg1}, "partner", []logging.ResultElem{arg2})

				res = true

				// remove the stuck operation from the list. If it is a select, remove all operations with the same val
				if vcTID2.val == -1 {
					leakingChannels[objID] = append(leakingChannels[objID][:i], leakingChannels[objID][i+1:]...)
				} else {
					for j, vcTID3 := range leakingChannels[objID] {
						if vcTID3.val == vcTID2.val {
							leakingChannels[objID] = append(leakingChannels[objID][:j], leakingChannels[objID][j+1:]...)
						}
					}
				}
			}
		}
	}
	return res
}

/*
 * After all operations have been analyzed, check if there are still leaking
 * operations without a possible partner.
 */
func CheckForLeak() {
	// channel
	for _, vcTIDs := range leakingChannels {
		buffered := false
		for _, vcTID := range vcTIDs {
			if vcTID.tID == "" {
				continue
			}

			found := false
			var partner VectorClockTID
			for _, c := range selectCases {
				if c.id != vcTID.id {
					continue
				}

				if (c.send && vcTID.typeVal == 0) || (!c.send && vcTID.typeVal == 1) {
					continue
				}

				hb := clock.GetHappensBefore(c.vcTID.Vc, vcTID.vc)
				if hb == clock.Concurrent {
					found = true
					if c.buffered {
						buffered = true
					}
					partner = c.vcTID
					break
				}

				if c.buffered {
					if (c.send && hb == clock.Before) || (!c.send && hb == clock.After) {
						found = true
						buffered = true
						partner = c.vcTID
						break
					}
				}
			}

			foundStr := ""
			if found {
				if vcTID.sel {
					if buffered {
						foundStr = "Leak on select with possible buffered partner:\n"
					} else {
						foundStr = "Leak on select with possible unbuffered partner:\n"
					}
					foundStr += "\tselect: " + vcTID.tID + "\n"
				} else {
					if buffered { // BUG: get unbuffered but should be bufferd
						foundStr = "Leak on buffered channel with possible partner:\n"
					} else {
						foundStr = "Leak on unbuffered channel with possible partner:\n"
					}
					foundStr += "\tchannel: " + vcTID.tID + "\n"
				}
				foundStr += "\tpartner: " + partner.TID
				logging.Result(foundStr, logging.CRITICAL)
			} else {
				if vcTID.sel {
					foundStr = "Leak on select without possible partner:\n"
					foundStr += "\tselect: " + vcTID.tID + "\n"
				} else {
					if buffered {
						foundStr = "Leak on buffered channel without possible partner:\n"
					} else {
						foundStr = "Leak on unbuffered channel without possible partner:\n"
					}
					foundStr += "\tchannel: " + vcTID.tID + "\n"
				}
				foundStr += "\tpartner: -"
				logging.Result(foundStr, logging.CRITICAL)
			}
		}
	}
}

/*
 * Run for select operation without a post event. Check if the operation has
 * a possible communication partner in mostRecentSend, mostRecentReceive or closeData.
 * If so, add an error or warning to the result.
 * If not, add all elements to leakingChannels, for later check.
 * MARK: SelectStuck
 * Args:
 *   routineID (int): The routine id
 *   ids (int): The channel ids
 *   buffered ([]bool): If the channels are buffered
 *   vc (VectorClock): The vector clock of the operation
 *   tID (string): The trace id
 *   opTypes ([]int): An identifier for the type of the operations (send = 0, recv = 1)
 *   idSel (int): The id of the select operation
 *   tPre (int): The tpre of the select operations. Used to connect the operations of the
 *     same select statement in leakingChannels.
 */
func CheckForLeakSelectStuck(routineID int, ids []int, buffered []bool, vc clock.VectorClock, tID string, opTypes []int, tPre int) {
	foundPartner := false

	if len(ids) == 0 {
		found := "Leak on select with only nil channels:\n"
		found += "\tselect: " + tID + "\n"
		found += "\tpartner: -"
		logging.Result(found, logging.CRITICAL)
		return
	}

	for i, id := range ids {
		if opTypes[i] == 0 { // send
			for _, mrr := range mostRecentReceive {
				if _, ok := mrr[id]; ok {
					if clock.GetHappensBefore(vc, mrr[id].Vc) == clock.Concurrent {
						found := ""
						if buffered[i] {
							found = "Leak on select with possible buffered partner:\n"
						} else {
							found = "Leak on select with possible unbuffered partner:\n"
						}
						found += "\tselect: " + tID + "\n"
						found += "\tpartner: " + mrr[id].TID + "\n"
						logging.Result(found, logging.CRITICAL)
						foundPartner = true
					}
				}
			}
		} else if opTypes[i] == 1 { // recv
			for _, mrs := range mostRecentSend {
				if _, ok := mrs[id]; ok {
					if clock.GetHappensBefore(vc, mrs[id].Vc) == clock.Concurrent {
						found := ""
						if buffered[i] {
							found = "Leak on select with possible buffered partner:\n"
						} else {
							found = "Leak on select with possible unbuffered partner:\n"
						}
						found += "\tselect: " + tID + "\n"
						found += "\tpartner: " + mrs[id].TID
						logging.Result(found, logging.CRITICAL)
						foundPartner = true
					}
				}
			}
			if _, ok := closeData[id]; ok {
				found := ""
				if buffered[i] {
					found = "Leak on select with possible buffered partner:\n"
				} else {
					found = "Leak on select with possible unbuffered partner:\n"
				}
				found += "\tselect: " + tID + "\n"
				found += "\tpartner: " + closeData[id].TID
				logging.Result(found, logging.CRITICAL)
				foundPartner = true
			}
		}
	}

	if !foundPartner {
		for i, id := range ids {
			// add all select operations to leaking Channels,
			leakingChannels[id] = append(leakingChannels[id], VectorClockTID2{routineID, id, vc, tID, opTypes[i], tPre, buffered[i], true})
		}
	}
}

/*
 * Run for mutex operation without a post event. Show an error in the results
 * MARK: Mutex
 * Args:
 *   routineID (int): The routine id
 *   id (int): The mutex id
 *   tID (string): The trace id
 *   op (int): The operation on the mutex
 */
func CheckForLeakMutex(routineID int, id int, tID string, op int) {
	file1, line1, tPre1, err := infoFromTID(tID)
	if err != nil {
		logging.Debug("Error in infoFromTID", logging.ERROR)
		return
	}

	file2, line2, tPre2, err := infoFromTID(mostRecentAcquireTotal[id].TID)
	if err != nil {
		logging.Debug("Error in infoFromTID", logging.ERROR)
		return
	}

	objType1 := "M"
	if op == 0 { // lock
		objType1 += "L"
	} else if op == 1 { // rlock
		objType1 += "R"
	} else { // only lock and rlock can lead to leak
		return
	}

	objType2 := "M"
	if mostRecentAcquireTotal[id].Val == 0 { // lock
		objType2 += "L"
	} else if mostRecentAcquireTotal[id].Val == 1 { // rlock
		objType2 += "R"
	} else if mostRecentAcquireTotal[id].Val == 2 { // TryLock
		objType2 += "T"
	} else if mostRecentAcquireTotal[id].Val == 3 { // TryRLock
		objType2 += "Y"
	} else { // only lock and rlock can lead to leak
		return
	}

	arg1 := logging.TraceElementResult{
		RoutineID: routineID, ObjID: id, TPre: tPre1, ObjType: objType1, File: file1, Line: line1}

	arg2 := logging.TraceElementResult{
		RoutineID: mostRecentAcquireTotal[id].Routine, ObjID: id, TPre: tPre2, ObjType: objType2, File: file2, Line: line2}

	logging.Result(logging.CRITICAL, logging.LMutex,
		"mutex", []logging.ResultElem{arg1}, "last", []logging.ResultElem{arg2})

}

/*
 * Add the most recent acquire operation for a mutex
 * Args:
 *   routine (int): The routine id
 *   id (int): The mutex id
 *   tID (string): The trace id
 *   vc (VectorClock): The vector clock of the operation
 *   op (int): The operation on the mutex
 */
func addMostRecentAcquireTotal(routine int, id int, tID string, vc clock.VectorClock, op int) {
	mostRecentAcquireTotal[id] = VectorClockTID3{Routine: routine, Vc: vc, TID: tID, Val: op}
}

/*
 * Run for wait group operation without a post event. Show an error in the results
 * MARK: WaitGroup
 * Args:
 *   routine (int): The routine id
 *   id (int): The wait group id
 *   tID (string): The trace id
 */
func CheckForLeakWait(routine int, id int, tID string) {
	file, line, tPre, err := infoFromTID(tID)
	if err != nil {
		logging.Debug("Error in infoFromTID", logging.ERROR)
		return
	}

	arg := logging.TraceElementResult{
		RoutineID: routine, ObjID: id, TPre: tPre, ObjType: "WW", File: file, Line: line}

	logging.Result(logging.CRITICAL, logging.LWaitGroup,
		"wait", []logging.ResultElem{arg}, "", []logging.ResultElem{})
}

/*
 * Run for conditional varable operation without a post event. Show an error in the results
 * MARK: Cond
 * Args:
 *   routine (int): The routine id
 *   id (int): The conditional variable id
 *   tID (string): The trace id
 */
func CheckForLeakCond(routine int, id int, tID string) {
	file, line, tPre, err := infoFromTID(tID)
	if err != nil {
		logging.Debug("Error in infoFromTID", logging.ERROR)
		return
	}

	arg := logging.TraceElementResult{
		RoutineID: routine, ObjID: id, TPre: tPre, ObjType: "NW", File: file, Line: line}

	logging.Result(logging.CRITICAL, logging.LCond,
		"cond", []logging.ResultElem{arg}, "", []logging.ResultElem{})
}
