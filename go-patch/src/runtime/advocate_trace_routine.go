package runtime

// type to save in the trace for routines
type advocateSpawnElement struct {
	id    uint64 // id of the routine
	timer uint64 // global timer
	file  string // file where the routine was created
	line  int    // line where the routine was created
}

func (elem advocateSpawnElement) isAdvocateTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "G,'id'"
 *    'id' (number): id of the routine
 */
func (elem advocateSpawnElement) toString() string {
	return "G," + uint64ToString(elem.timer) + "," + uint64ToString(elem.id) + "," + elem.file + ":" + intToString(elem.line)
}

/*
 * Get the operation
 */
func (elem advocateSpawnElement) getOperation() Operation {
	return OperationSpawn
}

/*
 * Get the file
 */
func (elem advocateSpawnElement) getFile() string {
	return elem.file
}

/*
 * Get the line
 */
func (elem advocateSpawnElement) getLine() int {
	return elem.line
}

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
	callerRoutine.addToTrace(advocateSpawnElement{id: newID, timer: timer,
		file: file, line: int(line)})
}

// type to save in the trace for routines
type advocateTraceSpawnedElement struct {
	id    uint64 // id of the routine
	timer uint64 // global timer
	file  string // file where the routine was created
	line  int    // line where the routine was created
}

func (elem advocateTraceSpawnedElement) isAdvocateTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "G,'id'"
 *    'id' (number): id of the routine
 */
func (elem advocateTraceSpawnedElement) toString() string {
	return "g," + uint64ToString(elem.timer) + "," + uint64ToString(elem.id) + "," + elem.file + ":" + intToString(elem.line)
}
