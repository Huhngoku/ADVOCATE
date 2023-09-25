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
 *   numberOfRoutines (int): The number of routines in the trace
 *   vc (*[]VectorClock): The vector clocks
 * Returns:
 *   (VectorClock): The vector clock for the read
 */
func Write(routine int, id int, numberOfRoutines int, vc *[]VectorClock) VectorClock {
	newLw(id, numberOfRoutines)
	lw[id] = (*vc)[routine]
	(*vc)[routine].Inc(routine)
	return (*vc)[routine]
}

/*
 * Calculate the new vector clock for a read operation and update cv
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the atomic variable
 *   numberOfRoutines (int): The number of routines in the trace
 *   vc (*[]VectorClock): The vector clocks
 * Returns:
 *   (VectorClock): The vector clock for the read
 */
func Read(routine int, id int, numberOfRoutines int, vc *[]VectorClock) VectorClock {
	newLw(id, numberOfRoutines)
	(*vc)[routine] = (*vc)[routine].Sync(lw[id])
	(*vc)[routine].Inc(routine)
	return (*vc)[routine]
}

/*
 * Calculate the new vector clock for a swap operation and update cv. A swap
 * operation is a read and a write.
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the atomic variable
 *   numberOfRoutines (int): The number of routines in the trace
 *   cv (*[]VectorClock): The vector clocks
 * Returns:
 *   (VectorClock): The vector clock for the read
 */
func Swap(routine int, id int, numberOfRoutines int, cv *[]VectorClock) VectorClock {
	_ = Read(routine, id, numberOfRoutines, cv)
	return Write(routine, id, numberOfRoutines, cv)
}
