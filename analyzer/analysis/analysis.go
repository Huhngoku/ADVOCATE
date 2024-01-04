package analysis

import (
	"analyzer/logging"
	"strconv"
)

// vc of close on channel
var closeVC = make(map[int]VectorClock)
var closePos = make(map[int]string)

// last send and receive on channel
var lastSend = make(map[int]VectorClock)
var lastRecv = make(map[int]VectorClock)

// last receive for each routine and each channel
var lastRecvRoutine = make(map[int]map[int]VectorClock)
var lastRecvRoutinePos = make(map[int]map[int]string)

// most recent send, used for detection of send on closed
var hasSend = make(map[int]bool)
var mostRecentSend = make(map[int]VectorClock)
var mostRecentSendPosition = make(map[int]string)

// most recent send, used for detection of received on closed
var hasReceived = make(map[int]bool)
var mostRecentReceive = make(map[int]VectorClock)
var mostRecentReceivePosition = make(map[int]string)

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

var addVcs = make(map[int]map[int][]VectorClock) // id -> routine -> []vc
var addPos = make(map[int]map[int][]string)      // id -> routine -> []pos

var doneVcs = make(map[int]map[int][]VectorClock) // id -> routine -> []vc
var donePos = make(map[int]map[int][]string)      // id > routine -> []pos

func checkForDoneBeforeAdd(routine int, id int, delta int, pos string, vc VectorClock) {
	if delta > 0 {
		checkForDoneBeforeAddAdd(routine, id, pos, vc, delta)
	} else if delta < 0 {
		checkForDoneBeforeAddDone(routine, id, pos, vc)
	} else {
		// checkForImpossibleWait(routine, id, pos, vc)
	}
}

func checkForDoneBeforeAddAdd(routine int, id int, pos string, vc VectorClock, delta int) {
	// if necessary, create maps and lists
	if _, ok := addVcs[id]; !ok {
		addVcs[id] = make(map[int][]VectorClock)
		addPos[id] = make(map[int][]string)
	}
	if _, ok := addVcs[id][routine]; !ok {
		addVcs[id][routine] = make([]VectorClock, 0)
		addPos[id][routine] = make([]string, 0)
	}

	// add the vector clock and position to the list
	for i := 0; i < delta; i++ {
		addVcs[id][routine] = append(addVcs[id][routine], vc.Copy())
		addPos[id][routine] = append(addPos[id][routine], pos)
	}

	// for now, test new vector clock against all done vector clocks
	// TODO: make this more efficient
	// for r, vcs := range doneVcs[id] {
	// 	for i, vcDone := range vcs {
	// 		happensBefore := GetHappensBefore(vcDone, vc)
	// 		if happensBefore == Concurrent {
	// 			found := "Found concurrent Add and Done on same waitgroup:\n"
	// 			found += "\tdone: " + donePos[id][r][i] + "\n"
	// 			found += "\tadd: " + addPos[id][routine][len(addPos[id][routine])-1]
	// 			logging.Result(found, logging.CRITICAL)
	// 		}
	// 	}
	// }
}

func checkForDoneBeforeAddDone(routine int, id int, pos string, vc VectorClock) {
	// if necessary, create maps and lists
	if _, ok := doneVcs[id]; !ok {
		doneVcs[id] = make(map[int][]VectorClock)
		donePos[id] = make(map[int][]string)
	}
	if _, ok := doneVcs[id][routine]; !ok {
		doneVcs[id][routine] = make([]VectorClock, 0)
		donePos[id][routine] = make([]string, 0)
	}

	// add the vector clock and position to the list
	doneVcs[id][routine] = append(doneVcs[id][routine], vc.Copy())
	donePos[id][routine] = append(donePos[id][routine], pos)

	// for now, test new vector clock against all add vector clocks
	// TODO: make this more efficient
	// for r, vcs := range addVcs[id] {
	// 	for i, vcAdd := range vcs {
	// 		happensBefore := GetHappensBefore(vcAdd, vc)
	// 		if happensBefore == Concurrent {
	// 			found := "Found concurrent Add and Done on same waitgroup:\n"
	// 			found += "\tdone: " + donePos[id][routine][len(donePos[id][routine])-1] + "\n"
	// 			found += "\tadd: " + addPos[id][r][i]
	// 			logging.Result(found, logging.CRITICAL)
	// 		}
	// 	}
	// }
}

/*
 * Check if a wait group counter could become negative
 * For each done operation, count the number of add operations a that happen before
 * the done operation, the number of done operations d that happen before the done operation
 * and the number of done operations d' that happen concurrent to the done operation.
 * If a < d + d', then the counter could become negative.
 * In this case, print a warning.
 */
func CheckForDoneBeforeAdd() {
	for id, vcs := range doneVcs { // for all waitgroups id
		for routine, vcs := range vcs { // for all routines
			for op, vcDone := range vcs { // for all done operations
				// count the number of add operations a that happen before or concurrent to the done operation
				countAdd := 0
				addPosList := []string{}
				for routineAdd, vcs := range addVcs[id] { // for all routines
					for opAdd, vcAdd := range vcs { // for all add operations
						happensBefore := GetHappensBefore(vcAdd, vcDone)
						if happensBefore == Before {
							countAdd++
						} else if happensBefore == Concurrent {
							addPosList = append(addPosList, addPos[id][routineAdd][opAdd])
						}
					}
				}
				// count the number of done operations d that happen before the done operation
				countDone := 0
				donePosList := []string{}
				for routine2, vcs := range doneVcs[id] { // for all routines
					for op2, vcDone2 := range vcs { // for all done operations
						if routine2 == routine && op2 == op {
							continue
						}
						happensBefore := GetHappensBefore(vcDone2, vcDone)
						if happensBefore == Before {
							countDone++
						} else if happensBefore == Concurrent {
							countDone++
							donePosList = append(donePosList, donePos[id][routine2][op2])
						}
					}
				}

				if countAdd < countDone {
					found := "Possible negative waitgroup counter:\n"
					found += "\tdone: " + donePos[id][routine][op] + "\n"
					found += "\tdone/add: "
					for i, pos := range donePosList {
						if i != 0 {
							found += ";"
						}
						found += pos
					}
					found += ";"
					for i, pos := range addPosList {
						if i != 0 {
							found += ";"
						}
						found += pos
					}
					logging.Result(found, logging.CRITICAL)
				}
			}
		}
	}
}
