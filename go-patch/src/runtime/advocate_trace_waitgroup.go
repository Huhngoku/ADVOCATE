package runtime

type advocateWaitGroupElement struct {
	id    uint64    // id of the waitgroup
	op    Operation // operation
	delta int       // delta of the waitgroup
	val   int32     // value of the waitgroup after the operation
	file  string    // file where the operation was called
	line  int       // line where the operation was called
	tPre  uint64    // global timer
	tPost uint64    // global timer
}

func (elem advocateWaitGroupElement) isAdvocateTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "W,'tPre','tPost','id','op','delta','val','file':'line'"
 *    'tPre' (number): global before the operation
 *    'tPost' (number): global after the operation
 *    'id' (number): id of the mutex
 *	  'op' (A/W): A if it is an add or Done, W if it is a wait
 *	  'delta' (number): delta of the waitgroup, positive for add, negative for done, 0 for wait
 *	  'val' (number): value of the waitgroup after the operation
 *    'file' (string): file where the operation was called
 *    'line' (number): line where the operation was called
 */
func (elem advocateWaitGroupElement) toString() string {
	res := "W,"
	res += uint64ToString(elem.tPre) + "," + uint64ToString(elem.tPost) + ","
	res += uint64ToString(elem.id) + ","
	switch elem.op {
	case OperationWaitgroupAddDone:
		res += "A,"
	case OperationWaitgroupWait:
		res += "W,"
	}

	res += intToString(elem.delta) + "," + int32ToString(elem.val)
	res += "," + elem.file + ":" + intToString(elem.line)
	return res
}

/*
 * Get the operation
 */
func (elem advocateWaitGroupElement) getOperation() Operation {
	return elem.op
}

/*
 * Get the file
 */
func (elem advocateWaitGroupElement) getFile() string {
	return elem.file
}

/*
 * Get the line
 */
func (elem advocateWaitGroupElement) getLine() int {
	return elem.line
}

/*
 * AdvocateWaitGroupAdd adds a waitgroup add or done to the trace
 * Args:
 * 	id: id of the waitgroup
 *  delta: delta of the waitgroup
 * 	val: value of the waitgroup after the operation
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateWaitGroupAdd(id uint64, delta int, val int32) int {
	var file string
	var line int
	if delta > 0 {
		_, file, line, _ = Caller(2)
	} else {
		_, file, line, _ = Caller(3)
	}
	timer := GetAdvocateCounter()
	elem := advocateWaitGroupElement{id: id, op: OperationWaitgroupAddDone,
		delta: delta, val: val, file: file, line: line, tPre: timer, tPost: timer}
	return insertIntoTrace(elem)

}

/*
 * AdvocateWaitGroupWaitPre adds a waitgroup wait to the trace
 * Args:
 * 	id: id of the waitgroup
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateWaitGroupWaitPre(id uint64) int {
	_, file, line, _ := Caller(2)
	timer := GetAdvocateCounter()
	elem := advocateWaitGroupElement{id: id, op: OperationWaitgroupWait,
		file: file, line: line, tPre: timer}
	return insertIntoTrace(elem)
}
