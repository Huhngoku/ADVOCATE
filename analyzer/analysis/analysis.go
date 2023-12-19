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

var lastAddValue = make(map[int]map[int]int)      // map[id]map[routine]value
var lastAddPos = make(map[int]map[int]string)     // map[id]map[routine]pos
var lastAddVc = make(map[int]map[int]VectorClock) // map[id]map[routine]vc

var lastDoneValue = make(map[int]map[int]int)      // map[id]map[routine]value
var lastDonePos = make(map[int]map[int]string)     // map[id]map[routine]pos
var lastDoneVc = make(map[int]map[int]VectorClock) // map[id]map[routine]vc

var doneCount = make(map[int]int)

/*
 * Update the last value of a wait group for an add operation
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the wait group
 *   delta (int): The delta of the operation
 *   pos (string): The position of the operation
 *   vc (VectorClock): The vector clock of the operation
 */
func checkForDoneBeforeAddAdd(routine int, id int, delta int, pos string, vc VectorClock) {
	updateDoneBeforeAdd(lastAddValue, lastAddPos, lastAddVc, routine, id, delta, pos, vc)
}

/*
 * Update the last value of a wait group for a done operation and check if a done before add can happen
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the wait group
 *   delta (int): The delta of the operation (made positive)
 *   pos (string): The position of the operation
 *   vc (VectorClock): The vector clock of the operation
 */
func checkForDoneBeforeAddDone(routine int, id int, delta int, pos string, vc VectorClock) {
	updateDoneBeforeAdd(lastDoneValue, lastDonePos, lastDoneVc, routine, id, delta, pos, vc)
	doneCount[id] += delta
	checkForDoneBeforeAdd(routine, id, vc)
}

/*
 * Check if a wg counter can become negative.
 * This can happen if #(delta add pre) < #(delta done pre) + #(delta done concurrent)
 * Args:
 */

func checkForDoneBeforeAdd(routine int, id int, vc VectorClock) {
	numberDonePre := lastDoneValue[id][routine] - 1
	numberTotal := doneCount[id]
	// max value in lastAddValue[id]
	numberAddPre := 0
	for _, value := range lastAddValue[id] {
		if value > numberAddPre && GetHappensBefore(lastAddVc[id][routine], vc) == Before {
			numberAddPre = value
		}
	}

	if numberAddPre < numberTotal {
		found := "Possible done before add:\n"
		found += "\tdone: " + lastDonePos[id][routine] + "\n"
		found += "\tadd : " // TODO: add position of add
		logging.Result(found, logging.CRITICAL)
	}

	println("checkForDoneBeforeWait: ", numberAddPre, numberDonePre, numberTotal)
}

/*
 * Update the last value of a wait group for a add or done operation for counting
 * the number of
 * Args:
 * lastValue (map[int]map[int]int): The last value map of a wait group
 * lastPos (map[int]map[int]string): The last position map of a wait group
 * lastVc (map[int]map[int]VectorClock): The last vector clock map of a wait group
 * routine (int): The routine id
 * id (int): The id of the wait group
 * delta (int): The delta of the operation
 * pos (string): The position of the operation
 * vc (VectorClock): The vector clock of the operation
 */
func updateDoneBeforeAdd(lastValue map[int]map[int]int,
	lastPos map[int]map[int]string, lastVc map[int]map[int]VectorClock,
	routine int, id int, delta int, pos string, vc VectorClock) {
	// create map if not exists
	if _, ok := lastValue[id]; !ok {
		lastValue[id] = make(map[int]int)
		lastPos[id] = make(map[int]string)
		lastVc[id] = make(map[int]VectorClock)
	}

	// get max value of map[id]
	max := 0
	for _, value := range lastValue[id] {
		if value > max && GetHappensBefore(lastVc[id][routine], vc) == Before {
			max = value
		}
	}
	lastPos[id][routine] = pos
	lastVc[id][routine] = vc.Copy()
	lastValue[id][routine] = max + delta
}

// /*
//  * Update the last list of concurrent done operations for a wait group
//  * Args:
//  *   id (int): The id of the wait group
//  *   vc (VectorClock): The vector clock of the operation
//  * Returns:
//  *   int: The number of concurrent done operations
//  */
// func updateDoneBeforeWaitCon(id int, vc VectorClock) int {

// for i, elem := range contDoneGroupes[id] {
// 	if GetHappensBefore(elem.VectorClock, vc) == Concurrent {
// 		contDoneGroupes[id][i].int++
// 		return contDoneGroupes[id][i].int
// 	}
// }
// contDoneGroupes[id] = append(contDoneGroupes[id], struct {
// 	VectorClock
// 	int
// }{vc.Copy(), 0})
// return 0
// }
