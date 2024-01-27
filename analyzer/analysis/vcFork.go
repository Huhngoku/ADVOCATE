package analysis

/*
 * Update the vector clocks given a fork operation
 * Args:
 *   oldRout (int): The id of the old routine
 *   newRout (int): The id of the new routine
 *   vcHb (map[int]VectorClock): The current hb vector clocks
 *   vcMhb (map[int]VectorClock): The current mhb vector clocks
 */
func Fork(oldRout int, newRout int, vcHb map[int]VectorClock, vcMhb map[int]VectorClock) {
	vcHb[newRout] = vcHb[oldRout].Copy()
	vcHb[newRout] = vcHb[newRout].Inc(newRout)
	vcHb[oldRout] = vcHb[oldRout].Inc(oldRout)

	vcMhb[newRout] = vcMhb[oldRout].Copy()
	vcMhb[newRout] = vcMhb[newRout].Inc(newRout)
	vcMhb[oldRout] = vcMhb[oldRout].Inc(oldRout)
}
