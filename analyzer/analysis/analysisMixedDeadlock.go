package analysis

import (
	"analyzer/logging"
	"strconv"
)

/*
 * Add a lock to the lockSet of a routine. Also save the vector clock of the acquire
 * Args:
 *   routine (int): The routine id
 *   lock (int): The id of the mutex
 *   pos (string): The position of the mutex operation
 *   vc (VectorClock): The current vector clock
 */
func lockSetAddLock(routine int, lock int, pos string, vc VectorClock) {
	if _, ok := lockSet[routine]; !ok {
		lockSet[routine] = make(map[int]string)
	}

	if posOld, ok := lockSet[routine][lock]; ok {
		errorMsg := "Lock " + strconv.Itoa(lock) +
			" already in lockSet for routine " + strconv.Itoa(routine)
		logging.Debug(errorMsg, logging.ERROR)

		// this is a double locking
		found := "Double locking:\n"
		found += "\tlock1: " + posOld + "\n"
		found += "\tlock2: " + pos
		logging.Result(found, logging.CRITICAL)
	}

	lockSet[routine][lock] = pos
	mostRecentAcquire[routine][lock] = vc
}

/*
 * Remove a lock from the lockSet of a routine
 * Args:
 *   routine (int): The routine id
 *   lock (int): The id of the mutex
 */
func lockSetRemoveLock(routine int, lock int) {
	if _, ok := lockSet[routine][lock]; !ok {
		errorMsg := "Lock " + strconv.Itoa(lock) +
			" not in lockSet for routine " + strconv.Itoa(routine)
		logging.Debug(errorMsg, logging.ERROR)
		return
	}
	delete(lockSet[routine], lock)
}

func checkForMixedDeadlock(routineSend int, routineRevc int) {
	for mSend, posSend := range lockSet[routineSend] {
		if _, ok := mostRecentAcquire[routineSend][mSend]; !ok { // Acq(tS, y) not defined
			continue
		}
		for mRecv, posRecv := range lockSet[routineRevc] {
			if _, ok := mostRecentAcquire[routineRevc][mRecv]; !ok { // Acq(tR, y) not defined
				continue
			}

			// no mixed deadlock possible if Acq(tS, y) <wmhb Acq(tR, y)
			weakHappensBefore := GetHappensBefore(mostRecentAcquire[routineSend][mSend], mostRecentAcquire[routineRevc][mRecv])
			if weakHappensBefore != Concurrent {
				continue
			}

			// found potential mixed deadlock
			found := "Potential mixed deadlock:\n"
			found += "\tlock1: " + posSend + "\n"
			found += "\tlock2: " + posRecv
		}
	}
}
