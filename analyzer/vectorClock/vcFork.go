package vectorClock

/*
 * Update the vector clocks given a fork operation
 * Args:
 *   oldRout (int): The id of the old routine
 *   newRout (int): The id of the new routine
 *   vc (VectorClock): The current vector clocks
 */
func Fork(oldRout int, newRout int, vc map[int]VectorClock) {
	vc[newRout] = vc[oldRout].Copy()
	vc[newRout] = vc[newRout].Inc(newRout)
	vc[oldRout] = vc[oldRout].Inc(oldRout)
}
