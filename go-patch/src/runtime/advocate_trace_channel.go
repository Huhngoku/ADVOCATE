package runtime

type advocateChannelElement struct {
	id     uint64    // id of the channel
	op     Operation // operation
	qSize  uint32    // size of the channel, 0 for unbuffered
	opID   uint64    // id of the operation
	file   string    // file where the operation was called
	line   int       // line where the operation was called
	tPre   uint64    // global timer before the operation
	tPost  uint64    // global timer after the operation
	closed bool      // true if the channel operation was finished, because the channel was closed at another routine
}

func (elem advocateChannelElement) isAdvocateTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "C,'tPre','tPost','id','op','pId','file':'line'"
 *    'tPre' (number): global timer before the operation
 *    'tPost' (number): global timer after the operation
 *    'id' (number): id of the channel
 *	  'op' (S/R/C): S if it is a send, R if it is a receive, C if it is a close
 *	  'pId' (number): id of the channel with witch the communication took place
 *    'file' (string): file where the operation was called
 *    'line' (number): line where the operation was called
 */
func (elem advocateChannelElement) toString() string {
	return elem.toStringSep(",", true)
}

/*
 * Get the operation
 */
func (elem advocateChannelElement) getOperation() Operation {
	return elem.op
}

/*
 * Get the file
 */
func (elem advocateChannelElement) getFile() string {
	return elem.file
}

/*
 * Get the line
 */
func (elem advocateChannelElement) getLine() int {
	return elem.line
}

/*
* Get a string representation of the element given a separator
* Args:
* 	sep: separator to use
* 	showPos: true if the position of the operation should be shown
* Return:
* 	string representation of the element
 */
func (elem advocateChannelElement) toStringSep(sep string, showPos bool) string {
	res := "C" + sep
	res += uint64ToString(elem.tPre) + sep + uint64ToString(elem.tPost) + sep
	res += uint64ToString(elem.id) + sep

	switch elem.op {
	case OperationChannelSend:
		res += "S"
	case OperationChannelRecv:
		res += "R"
	case OperationChannelClose:
		res += "C"
	default:
		panic("Unknown channel operation" + intToString(int(elem.op)))
	}

	if elem.closed {
		res += sep + "t"
	} else {
		res += sep + "f"
	}

	res += sep + uint64ToString(elem.opID)
	res += sep + uint32ToString(elem.qSize)
	if showPos {
		res += sep + elem.file + ":" + intToString(elem.line)
	}
	return res
}

var advocateCounterAtomic uint64

/*
 * AdvocateChanSendPre adds a channel send to the trace.
 * If the channel send was created by an atomic
 * operation, add this to the trace as well
 * Args:
 * 	id: id of the channel
 * 	opId: id of the operation
 * 	qSize: size of the channel, 0 for unbuffered
 * Return:
 * 	index of the operation in the trace, return -1 if it is a atomic operation
 */
func AdvocateChanSendPre(id uint64, opID uint64, qSize uint) int {
	_, file, line, _ := Caller(3)
	// internal channels to record atomic operations
	if isSuffix(file, "advocate_atomic.go") {
		advocateCounterAtomic++
		lock(&advocateAtomicMapRoutineLock)
		advocateAtomicMapRoutine[advocateCounterAtomic] = GetRoutineID()
		unlock(&advocateAtomicMapRoutineLock)
		AdvocateAtomic(advocateCounterAtomic)

		// they are not recorded in the trace
		return -1
	}
	timer := GetAdvocateCounter()
	elem := advocateChannelElement{id: id, op: OperationChannelSend,
		opID: opID, file: file, line: line, tPre: timer, qSize: uint32(qSize)}
	return insertIntoTrace(elem)
}

/*
 * Helper function to check if a string ends with a suffix
 * Args:
 * 	s: string to check
 * 	suffix: suffix to check
 * Return:
 * 	true if s ends with suffix, false otherwise
 */
func isSuffix(s, suffix string) bool {
	if len(suffix) > len(s) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}

/*
 * AdvocateChanRecvPre adds a channel recv to the trace
 * Args:
 * 	id: id of the channel
 * 	opId: id of the operation
 * 	qSize: size of the channel
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateChanRecvPre(id uint64, opID uint64, qSize uint) int {
	_, file, line, _ := Caller(3)
	// do not record channel operation of internal channel to record atomic operations
	if isSuffix(file, "advocate_trace.go") {
		return -1
	}

	timer := GetAdvocateCounter()
	elem := advocateChannelElement{id: id, op: OperationChannelRecv,
		opID: opID, file: file, line: line, tPre: timer, qSize: uint32(qSize)}
	return insertIntoTrace(elem)
}

/*
 * AdvocateChanClose adds a channel close to the trace
 * Args:
 * 	id: id of the channel
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateChanClose(id uint64, qSize uint) int {
	_, file, line, _ := Caller(2)
	timer := GetAdvocateCounter()
	elem := advocateChannelElement{id: id, op: OperationChannelClose,
		file: file, line: line, tPre: timer, tPost: timer, qSize: uint32(qSize)}
	return insertIntoTrace(elem)
}

/*
 * AdvocateChanPost sets the operation as successfully finished
 * Args:
 * 	index: index of the operation in the trace
 */
func AdvocateChanPost(index int) {
	if index == -1 {
		return
	}

	elem := currentGoRoutine().getElement(index).(advocateChannelElement)
	elem.tPost = GetAdvocateCounter()
	currentGoRoutine().updateElement(index, elem)
}

/*
 * AdvocateChanPostCausedByClose sets the operation as successfully finished
 * Args:
 * 	index: index of the operation in the trace
 */
func AdvocateChanPostCausedByClose(index int) {
	if index == -1 {
		return
	}
	elem := currentGoRoutine().getElement(index).(advocateChannelElement)
	elem.closed = true
	currentGoRoutine().updateElement(index, elem)
}
