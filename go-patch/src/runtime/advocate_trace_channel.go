package runtime

var advocateCounterAtomic uint64

// MARK: Pre

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
		AdvocateAtomicPre(advocateCounterAtomic)

		// they are not recorded in the trace
		return -1
	}
	timer := GetNextTimeStep()
	elem := "C," + uint64ToString(timer) + ",0," + uint64ToString(id) + ",S,f," +
		uint64ToString(opID) + "," + uint32ToString(uint32(qSize)) + "," +
		file + ":" + intToString(line)

	return insertIntoTrace(elem, false)
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

	timer := GetNextTimeStep()
	elem := "C," + uint64ToString(timer) + ",0," + uint64ToString(id) + ",R,f," +
		uint64ToString(opID) + "," + uint32ToString(uint32(qSize)) + "," +
		file + ":" + intToString(line)
	return insertIntoTrace(elem, false)
}

// MARK: Close

/*
 * AdvocateChanClose adds a channel close to the trace
 * Args:
 * 	id: id of the channel
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateChanClose(id uint64, qSize uint) int {
	_, file, line, _ := Caller(2)
	timer := uint64ToString(GetNextTimeStep())
	elem := "C," + timer + "," + timer + "," + uint64ToString(id) + ",C,f,0," +
		uint32ToString(uint32(qSize)) + "," + file + ":" + intToString(line)

	return insertIntoTrace(elem, false)
}

// MARK: Post

/*
 * AdvocateChanPost sets the operation as successfully finished
 * Args:
 * 	index: index of the operation in the trace
 */
func AdvocateChanPost(index int) {
	if index == -1 {
		return
	}

	elem := currentGoRoutine().getElement(index)

	split := splitStringAtCommas(elem, []int{2, 3})
	split[1] = uint64ToString(GetNextTimeStep())
	elem = mergeString(split)

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

	elem := currentGoRoutine().getElement(index)
	split := splitStringAtCommas(elem, []int{2, 3, 5, 6})
	split[1] = uint64ToString(GetNextTimeStep())
	split[3] = "t"
	elem = mergeString(split)

	currentGoRoutine().updateElement(index, elem)
}
