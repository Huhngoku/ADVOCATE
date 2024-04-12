package analysis

import "analyzer/clock"

// vector clock for each wait group
var wg map[int]clock.VectorClock = make(map[int]clock.VectorClock)

/*
 * Create a new wg if needed
 * Args:
 *   index (int): The id of the wait group
 *   nRout (int): The number of routines in the trace
 */
func newWg(index int, nRout int) {
	if _, ok := wg[index]; !ok {
		wg[index] = clock.NewVectorClock(nRout)
	}
}

/*
 * Calculate the new vector clock for a add or done operation and update cv
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the wait group
 *   delta (int): The delta of the wait group
 *   tID (string): The id of the trace element, contains the position and the tpre
 *   vc (map[int]VectorClock): The vector clocks
 */
func Change(routine int, id int, delta int, tID string, vc map[int]clock.VectorClock) {
	newWg(id, vc[id].GetSize())
	wg[id] = wg[id].Sync(vc[routine])
	vc[routine] = vc[routine].Inc(routine)

	if analysisCases["doneBeforeAdd"] {
		checkForDoneBeforeAddChange(routine, id, delta, tID, vc[routine])
	}
}

/*
 * Calculate the new vector clock for a wait operation and update cv
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the wait group
 *   vc (*map[int]VectorClock): The vector clocks
 *   notLeak (bool): If the wait group is not leaked (tpost = 0)
 */
func Wait(routine int, id int, tID string, vc map[int]clock.VectorClock, notLeak bool) {
	newWg(id, vc[id].GetSize())
	if notLeak {
		vc[routine] = vc[routine].Sync(wg[id])
		vc[routine] = vc[routine].Inc(routine)
	}
}
