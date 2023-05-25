// DEDEGO_FILE_START

package runtime

/*
 * Get the id of the current routine
 * Return:
 * 	id of the current routine
 */
func GetRoutineId() uint64 {
	return getg().goid
}

// DEDEGO-FILE-END
