package analysis

var currentWaits = make(map[int][]int) // -> id -> routine

/*
 * Update and calculate the vector clocks given a wait operation
 * Args:
 *   id (int): The id of the condition variable
 *   routine (int): The routine id
 *   vc (map[int]VectorClock): The current vector clocks
 */
func CondWait(id int, routine int, vc map[int]VectorClock) {
	currentWaits[id] = append(currentWaits[id], routine)
	vc[routine].Inc(routine)
}

/*
 * Update and calculate the vector clocks given a signal operation
 * Args:
 *   id (int): The id of the condition variable
 *   routine (int): The routine id
 *   vc (map[int]VectorClock): The current vector clocks
 */
func CondSignal(id int, routine int, vc map[int]VectorClock) {
	if len(currentWaits[id]) != 0 {
		waitRoutine := currentWaits[id][0]
		currentWaits[id] = currentWaits[id][1:]
		vc[waitRoutine].Sync(vc[routine])
	}
	vc[routine].Inc(routine)
}

/*
 * Update and calculate the vector clocks given a broadcast operation
 * Args:
 *   id (int): The id of the condition variable
 *   routine (int): The routine id
 *   vc (map[int]VectorClock): The current vector clocks
 */
func CondBroadcast(id int, routine int, vc map[int]VectorClock) {
	for _, waitRoutine := range currentWaits[id] {
		vc[waitRoutine].Sync(vc[routine])
	}
	currentWaits[id] = []int{}
	vc[routine].Inc(routine)
}
