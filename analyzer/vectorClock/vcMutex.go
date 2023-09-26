package vectorClock

// vector clocks for last release times
var relW map[int]VectorClock = make(map[int]VectorClock)
var relR map[int]VectorClock = make(map[int]VectorClock)

/*
 * Create a new relW and relR if needed
 * Args:
 *   index (int): The id of the atomic variable
 *   nRout (int): The number of routines in the trace
 */
func newRel(index int, nRout int) {
	if _, ok := relW[index]; !ok {
		relW[index] = NewVectorClock(nRout)
		relR[index] = NewVectorClock(nRout)
	}
}

/*
 * Update and calculate the vector clocks given a lock operation
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the mutex
 *   vc (map[int]VectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func Lock(routine int, id int, vc map[int]VectorClock) VectorClock {
	newRel(id, vc[id].size)
	vc[routine].Sync(relW[id])
	vc[routine].Sync(relR[id])
	return vc[routine].Inc(routine).Copy()
}

/*
 * Update and calculate the vector clocks given a unlock operation
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the mutex
 *   vc (map[int]VectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func Unlock(routine int, id int, vc map[int]VectorClock) VectorClock {
	newRel(id, vc[id].size)
	relW[id] = vc[routine]
	relR[id] = vc[routine]
	return vc[routine].Inc(routine).Copy()
}

/*
 * Update and calculate the vector clocks given a rlock operation
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the mutex
 *   vc (map[int]VectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func RLock(routine int, id int, vc map[int]VectorClock) VectorClock {
	newRel(id, vc[id].size)
	vc[routine].Sync(relW[id])
	return vc[routine].Inc(routine).Copy()
}

/*
 * Update and calculate the vector clocks given a runlock operation
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the mutex
 *   vc (map[int]VectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func RUnlock(routine int, id int, vc map[int]VectorClock) VectorClock {
	newRel(id, vc[id].size)
	relR[id].Sync(vc[routine])
	return vc[routine].Inc(routine).Copy()
}
