package vectorClock

// vector clocks for the successful do
var oSuc map[int]VectorClock = make(map[int]VectorClock)

/*
 * Create a new oSuc if needed
 * Args:
 *   index (int): The id of the atomic variable
 *   nRout (int): The number of routines in the trace
 */
func newOSuc(index int, nRout int) {
	if _, ok := oSuc[index]; !ok {
		oSuc[index] = NewVectorClock(nRout)
	}
}

/*
 * Update and calculate the vector clocks given a successful do operation
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the atomic variable
 *   vc (map[int]VectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func DoSuc(routine int, id int, vc map[int]VectorClock) VectorClock {
	newOSuc(id, vc[id].size)
	oSuc[id] = vc[routine]
	vc[routine] = vc[routine].Inc(routine)
	return vc[routine].Copy()
}

/*
 * Update and calculate the vector clocks given a unsuccessful do operation
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the atomic variable
 *   vc (map[int]VectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func DoFail(routine int, id int, vc map[int]VectorClock) VectorClock {
	newOSuc(id, vc[id].size)
	vc[routine] = vc[routine].Sync(oSuc[id])
	vc[routine] = vc[routine].Inc(routine)
	return vc[routine].Copy()
}
