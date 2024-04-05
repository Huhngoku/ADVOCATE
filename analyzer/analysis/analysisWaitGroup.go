package analysis

import (
	"analyzer/logging"
	"sort"
)

func checkForDoneBeforeAddChange(routine int, id int, delta int, pos string, vc VectorClock) {
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
	if _, ok := addWait[id]; !ok {
		addWait[id] = make(map[int][]VectorClockTID)
	}
	if _, ok := addWait[id][routine]; !ok {
		addWait[id][routine] = make([]VectorClockTID, 0)
	}

	// add the vector clock and position to the list
	for i := 0; i < delta; i++ {
		addWait[id][routine] = append(addWait[id][routine], VectorClockTID{vc.Copy(), pos})
	}
}

func checkForDoneBeforeAddDone(routine int, id int, pos string, vc VectorClock) {
	// if necessary, create maps and lists
	if _, ok := doneWait[id]; !ok {
		doneWait[id] = make(map[int][]VectorClockTID)

	}
	if _, ok := doneWait[id][routine]; !ok {
		doneWait[id][routine] = make([]VectorClockTID, 0)
	}

	// add the vector clock and position to the list
	doneWait[id][routine] = append(doneWait[id][routine], VectorClockTID{vc.Copy(), pos})

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
	for id, vcs := range doneWait { // for all waitgroups id
		for routine, vc := range vcs { // for all routines
			for op, vcDone := range vc { // for all done operations
				// count the number of add operations a that happen before or concurrent to the done operation
				countAdd := 0
				addPosList := []string{}
				for routineAdd, vcs := range addWait[id] { // for all routines
					for opAdd, vcAdd := range vcs { // for all add operations
						happensBefore := GetHappensBefore(vcAdd.vc, vcDone.vc)
						if happensBefore == Before {
							countAdd++
						} else if happensBefore == Concurrent {
							addPosList = append(addPosList, addWait[id][routineAdd][opAdd].tID)
						}
					}
				}
				// count the number of done operations d that happen before the done operation
				countDone := 1 // self
				donePosListConc := []string{}
				for routine2, vcs := range doneWait[id] { // for all routines
					for op2, vcDone2 := range vcs { // for all done operations
						if routine2 == routine && op2 == op {
							continue
						}
						happensBefore := GetHappensBefore(vcDone2.vc, vcDone.vc)
						if happensBefore == Before {
							countDone++
						} else if happensBefore == Concurrent {
							countDone++
							donePosListConc = append(donePosListConc, doneWait[id][routine2][op2].tID)
						}
					}
				}

				if countAdd < countDone {
					createDoneBeforeAddMessage(id, routine, op, addPosList, donePosListConc)
				}
			}
		}
	}
}

func createDoneBeforeAddMessage(id int, routine int, op int, addPosList []string, donePosListConc []string) {
	sort.Strings(addPosList)
	sort.Strings(donePosListConc)

	// add/done contains the add and done, that are concurrent to the done operation
	// add and done is separated by a #

	found := "Possible negative waitgroup counter:\n"
	found += "\tdone: " + doneWait[id][routine][op].tID + "\n"
	found += "\tadd/done: "

	for _, pos := range donePosListConc {
		found += pos + ";"
	}

	for _, pos := range addPosList {
		found += pos + ";"
	}

	logging.Result(found, logging.CRITICAL)
}
