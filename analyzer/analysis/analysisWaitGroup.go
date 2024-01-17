package analysis

import (
	"analyzer/logging"
	"sort"
)

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
					createDoneBeforeAddMessage(id, routine, op, addPosList, donePosList)
				}
			}
		}
	}
}

func createDoneBeforeAddMessage(id int, routine int, op int, addPosList []string, donePosList []string) {
	uniquePos := make(map[string]bool)
	sort.Strings(addPosList)
	sort.Strings(donePosList)

	found := "Possible negative waitgroup counter:\n"
	found += "\tdone: " + donePos[id][routine][op] + "\n"
	found += "\tdone/add: "
	for i, pos := range donePosList {
		if uniquePos[pos] {
			continue
		}
		if i != 0 {
			found += ";"
		}
		found += pos
		uniquePos[pos] = true
	}
	found += ";"
	for i, pos := range addPosList {
		if uniquePos[pos] {
			continue
		}
		if i != 0 {
			found += ";"
		}
		found += pos
		uniquePos[pos] = true
	}
	logging.Result(found, logging.CRITICAL)
}
