package vectorClock

var relW map[int]VectorClock = make(map[int]VectorClock)
var relR map[int]VectorClock = make(map[int]VectorClock)

func newRel(index int, nRout int) {
	if _, ok := relW[index]; !ok {
		relW[index] = NewVectorClock(nRout)
		relR[index] = NewVectorClock(nRout)
	}
}

/*
 * Update and calculate the vector clocks given a lock operation
 * Args:
 *   vc (vectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func Lock(routine int, id int, nRout int, vc *[]VectorClock) VectorClock {
	newRel(id, nRout)
	(*vc)[routine] = (*vc)[routine].Sync(relW[id])
	(*vc)[routine] = (*vc)[routine].Sync(relR[id])
	(*vc)[routine].Inc(routine)
	return (*vc)[routine]
}

/*
 * Update and calculate the vector clocks given a unlock operation
 * Args:
 *   vc (vectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func Unlock(routine int, id int, nRout int, vc *[]VectorClock) VectorClock {
	newRel(id, nRout)
	relW[id] = (*vc)[routine]
	relR[id] = (*vc)[routine]
	(*vc)[routine].Inc(routine)
	return (*vc)[routine]
}

/*
 * Update and calculate the vector clocks given a rlock operation
 * Args:
 *   vc (vectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func RLock(routine int, id int, nRout int, vc *[]VectorClock) VectorClock {
	newRel(id, nRout)
	(*vc)[routine] = (*vc)[routine].Sync(relW[id])
	(*vc)[routine].Inc(routine)
	return (*vc)[routine]
}

/*
 * Update and calculate the vector clocks given a runlock operation
 * Args:
 *   vc (vectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func RUnlock(routine int, id int, nRout int, vc *[]VectorClock) VectorClock {
	newRel(id, nRout)
	relR[id] = (*vc)[routine].Sync(relR[id])
	(*vc)[routine].Inc(routine)
	return (*vc)[routine]
}
