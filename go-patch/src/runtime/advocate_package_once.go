package runtime

type advocateOnceElement struct {
	tpre  uint64 // global timer at the beginning of the execution
	tpost uint64 // global timer at the end of the execution
	id    uint64 // id of the once
	suc   bool   // true if the do on the once was called for the first time
	file  string // file where the operation was called
	line  int    // line where the operation was called
}

func (elem advocateOnceElement) isAdvocateTraceElement() {}

/*
 * Return a string representation of the element
 * Return:
 * 	string representation of the element "O,'tpre','tpost','id','suc','file':'line"
 *    'tpre' (number): global timer at the beginning of the execution
 *    'tpost' (number): global timer at the end of the execution
 *    'id' (number): id of the once
 *    'suc' (t/f): true if the do on the once was called for the first time, false otherwise
 *    'file' (string): file where the operation was called
 *    'line' (string): line where the operation was called
 */
func (elem advocateOnceElement) toString() string {
	res := "O,"
	res += uint64ToString(elem.tpre) + ","
	res += uint64ToString(elem.tpost) + ","
	res += uint64ToString(elem.id) + ","
	if elem.suc {
		res += "t"
	} else {
		res += "f"
	}
	res += "," + elem.file + ":" + intToString(elem.line)
	return res
}

/*
 * Get the operation
 */
func (elem advocateOnceElement) getOperation() Operation {
	return OperationOnce
}

/*
 * Get the file
 */
func (elem advocateOnceElement) getFile() string {
	return elem.file
}

/*
 * Get the line
 */
func (elem advocateOnceElement) getLine() int {
	return elem.line
}

/*
 * AdvocateOncePre adds a once to the trace
 * Args:
 * 	id: id of the once
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateOncePre(id uint64) int {
	_, file, line, _ := Caller(2)
	timer := GetNextTimeStep()
	elem := advocateOnceElement{id: id, tpre: timer, file: file, line: line}
	return insertIntoTrace(elem)
}

/*
 * Add the end counter to an operation of the trace
 * Args:
 * 	index: index of the operation in the trace
 * 	suc: true if the do on the once was called for the first time, false otherwise
 */
func AdvocateOncePost(index int, suc bool) {
	if index == -1 {
		return
	}
	timer := GetNextTimeStep()
	elem := currentGoRoutine().getElement(index).(advocateOnceElement)
	elem.tpost = timer
	elem.suc = suc

	currentGoRoutine().updateElement(index, elem)
}
