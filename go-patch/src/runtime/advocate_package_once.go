package runtime

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

	elem := "O," + uint64ToString(timer) + ",0," + uint64ToString(id) + ",f," +
		file + ":" + intToString(line)

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
	elem := currentGoRoutine().getElement(index)

	split := splitStringAtCommas(elem, []int{2, 3, 4, 5})
	split[1] = uint64ToString(timer)
	if suc {
		split[3] = "t"
	}
	elem = mergeString(split)

	currentGoRoutine().updateElement(index, elem)
}
