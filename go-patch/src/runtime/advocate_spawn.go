// ADVOCATE-FILE_START
package runtime

/*
 * Wait and record for a spawn operation.
 * This function must be inserted at the beginning of each spawned routine.
 */
func AdvocateSpawnWait() {
	_, file, line, _ := Caller(1)
	_, _ = WaitForReplayPath(AdvocateReplaySpawn, file, line)

	elem := advocateTraceSpawnedElement{
		id:    GetRoutineId(),
		timer: GetAdvocateCounter(),
		file:  file,
		line:  line,
	}
	insertIntoTrace(elem)
}

// ADVOCATE-FILE-END
