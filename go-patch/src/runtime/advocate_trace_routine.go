package runtime

/*
 * AdvocateSpawnCaller adds a routine spawn to the trace
 * Args:
 * 	callerRoutine: routine that created the new routine
 * 	newID: id of the new routine
 * 	file: file where the routine was created
 * 	line: line where the routine was created
 */
func AdvocateSpawnCaller(callerRoutine *AdvocateRoutine, newID uint64, file string, line int32) {
	timer := GetAdvocateCounter()

	elem := "G," + uint64ToString(timer) + "," + uint64ToString(newID) + "," + file + ":" + int32ToString(line)

	callerRoutine.addToTrace(elem)
}
