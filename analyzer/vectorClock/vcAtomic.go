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
 */
func Write(routine int, id int, vc map[int]VectorClock) {
	newLw(id, vc[id].size)
	lw[id] = vc[routine].Copy()
	vc[routine] = vc[routine].Inc(routine)
}

/*
 * Calculate the new vector clock for a read operation and update cv
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the atomic variable
 *   numberOfRoutines (int): The number of routines in the trace
 *   vc (map[int]VectorClock): The vector clocks
 */
func Read(routine int, id int, vc map[int]VectorClock) {
	newLw(id, vc[id].size)
	vc[routine] = vc[routine].Sync(lw[id])
	vc[routine] = vc[routine].Inc(routine)
}

/*
 * Calculate the new vector clock for a swap operation and update cv. A swap
 * operation is a read and a write.
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the atomic variable
 *   numberOfRoutines (int): The number of routines in the trace
 *   cv (map[int]VectorClock): The vector clocks
 */
func Swap(routine int, id int, cv map[int]VectorClock) {
	Read(routine, id, cv)
	Write(routine, id, cv)
}
