package vectorClock

// vector clocks for last write times
var lw map[int]VectorClock = make(map[int]VectorClock)

/*
 * Create a new lw if needed
 * Args:
 *   index (int): The id of the atomic variable
 *   nRout (int): The number of routines in the trace
 */
func newLw(index int, nRout int) {
	if _, ok := lw[index]; !ok {
		lw[index] = NewVectorClock(nRout)
	}
}

/*
 * Calculate the new vector clock for a write operation and update cv
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the atomic variable
 *   vc (*map[int]VectorClock): The vector clocks
 * Returns:
 *   (VectorClock): The vector clock for the read
 */
func Write(routine int, id int, vc map[int]VectorClock) VectorClock {
	newLw(id, vc[id].size)
	lw[id] = vc[routine]
	return vc[routine].Inc(routine).Copy()
}

/*
 * Calculate the new vector clock for a read operation and update cv
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the atomic variable
 *   numberOfRoutines (int): The number of routines in the trace
 *   vc (map[int]VectorClock): The vector clocks
 * Returns:
 *   (VectorClock): The vector clock for the read
 */
func Read(routine int, id int, vc map[int]VectorClock) VectorClock {
	newLw(id, vc[id].size)
	vc[routine].Sync(lw[id])
	return vc[routine].Inc(routine).Copy()
}

/*
 * Calculate the new vector clock for a swap operation and update cv. A swap
 * operation is a read and a write.
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the atomic variable
 *   numberOfRoutines (int): The number of routines in the trace
 *   cv (map[int]VectorClock): The vector clocks
 * Returns:
 *   (VectorClock): The vector clock for the read
 */
func Swap(routine int, id int, cv map[int]VectorClock) VectorClock {
	_ = Read(routine, id, cv)
	return Write(routine, id, cv).Copy()
}
