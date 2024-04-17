package analysis

import "analyzer/clock"

/*
 * Update the vector clocks given a fork operation
 * Args:
 *   oldRout (int): The id of the old routine
 *   newRout (int): The id of the new routine
 *   vcHb (map[int]VectorClock): The current hb vector clocks
 *   vcMhb (map[int]VectorClock): The current mhb vector clocks
 */
func Fork(oldRout int, newRout int, vcHb map[int]clock.VectorClock, vcMhb map[int]clock.VectorClock) {
	vcHb[oldRout] = vcHb[oldRout].Inc(oldRout)
	vcHb[newRout] = vcHb[oldRout].Copy()
	vcHb[newRout] = vcHb[newRout].Inc(newRout)

	vcMhb[oldRout] = vcMhb[oldRout].Inc(oldRout)
	vcMhb[newRout] = vcMhb[oldRout].Copy()
	vcMhb[newRout] = vcMhb[newRout].Inc(newRout)
}
