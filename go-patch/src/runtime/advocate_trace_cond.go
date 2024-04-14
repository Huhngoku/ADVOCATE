package runtime

type advocateCondElement struct {
	tpre  uint64 // global timer at the beginning of the execution
	tpost uint64 // global timer at the end of the execution
	id    uint64 // id of the cond
	op    Operation
	file  string // file where the operation was called
	line  int    // line where the operation was called
}

func (elem advocateCondElement) isAdvocateTraceElement() {}

/*
 * Return a string representation of the element
 * Return:
 * 	string representation of the element "C,'tpre','tpost','id','op','file':'line"
 *    'tpre' (number): global timer at the beginning of the execution
 *    'tpost' (number): global timer at the end of the execution
 *    'id' (number): id of the cond
 *    'op' (W/S/B): W if it is a wait, S if it is a signal, B if it is a broadcast
 *    'file' (string): file where the operation was called
 *    'line' (string): line where the operation was called
 */
func (elem advocateCondElement) toString() string {
	res := "N,"
	res += uint64ToString(elem.tpre) + ","
	res += uint64ToString(elem.tpost) + ","
	res += uint64ToString(elem.id) + ","
	switch elem.op {
	case OperationCondWait:
		res += "W"
	case OperationCondSignal:
		res += "S"
	case OperationCondBroadcast:
		res += "B"
	}
	res += "," + elem.file + ":" + intToString(elem.line)
	return res
}

/*
 * Get the operation
 */
func (elem advocateCondElement) getOperation() Operation {
	return elem.op
}

/*
 * Get the file
 */
func (elem advocateCondElement) getFile() string {
	return elem.file
}

/*
 * Get the line
 */
func (elem advocateCondElement) getLine() int {
	return elem.line
}

/*
 * AdvocateCondPre adds a cond wait to the trace
 * Args:
 * 	id: id of the cond
 * 	op: 0 for wait, 1 for signal, 2 for broadcast
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateCondPre(id uint64, op int) int {
	_, file, line, _ := Caller(2)
	timer := GetAdvocateCounter()
	var opC Operation
	switch op {
	case 0:
		opC = OperationCondWait
	case 1:
		opC = OperationCondSignal
	case 2:
		opC = OperationCondBroadcast
	default:
		panic("Unknown cond operation")
	}
	elem := advocateCondElement{id: id, tpre: timer, file: file, line: line, op: opC}
	return insertIntoTrace(elem)
}

/*
 * AdvocateCondPost adds the end counter to an operation of the trace
 * Args:
 * 	index: index of the operation in the trace
 */
func AdvocateCondPost(index int) {
	if index == -1 {
		return
	}
	timer := GetAdvocateCounter()
	elem := currentGoRoutine().getElement(index).(advocateCondElement)
	elem.tpost = timer

	currentGoRoutine().updateElement(index, elem)
}
