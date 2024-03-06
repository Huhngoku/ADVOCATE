package analysis

// vector clock for each wait group
var wg map[int]VectorClock = make(map[int]VectorClock)

/*
 * Create a new wg if needed
 * Args:
 *   index (int): The id of the wait group
 *   nRout (int): The number of routines in the trace
 */
func newWg(index int, nRout int) {
	if _, ok := wg[index]; !ok {
		wg[index] = NewVectorClock(nRout)
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
func Change(routine int, id int, delta int, tID string, vc map[int]VectorClock) {
	newWg(id, vc[id].size)
	wg[id] = wg[id].Sync(vc[routine])
	vc[routine] = vc[routine].Inc(routine)

	if analysisCases["doneBeforeAdd"] {
		checkForDoneBeforeAdd(routine, id, delta, tID, vc[routine])
	}
}

/*
 * Calculate the new vector clock for a wait operation and update cv
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the wait group
 *   vc (*map[int]VectorClock): The vector clocks
 */
func Wait(routine int, id int, vc map[int]VectorClock) {
	newWg(id, vc[id].size)
	vc[routine] = vc[routine].Sync(wg[id])
	vc[routine] = vc[routine].Inc(routine)
}
