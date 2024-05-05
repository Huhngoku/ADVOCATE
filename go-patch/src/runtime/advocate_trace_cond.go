package runtime

/*
 * AdvocateCondPre adds a cond wait to the trace
 * MARK: Pre
 * Args:
 * 	id: id of the cond
 * 	op: 0 for wait, 1 for signal, 2 for broadcast
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateCondPre(id uint64, op int) int {
	_, file, line, _ := Caller(2)
	timer := GetNextTimeStep()
	var opC string
	switch op {
	case 0:
		opC = "W"
	case 1:
		opC = "S"
	case 2:
		opC = "B"
	default:
		panic("Unknown cond operation")
	}

	elem := "N," + uint64ToString(timer) + ",0," + uint64ToString(id) +
		"," + opC + "," + file + ":" + uint64ToString(uint64(line))
	return insertIntoTrace(elem, false)
}

/*
 * AdvocateCondPost adds the end counter to an operation of the trace
 * MARK: Post
 * Args:
 * 	index: index of the operation in the trace
 */
func AdvocateCondPost(index int) {
	if index == -1 {
		return
	}
	timer := GetNextTimeStep()
	elem := currentGoRoutine().getElement(index)

	split := splitStringAtCommas(elem, []int{2, 3})
	split[1] = uint64ToString(timer)

	elem = mergeString(split)

	currentGoRoutine().updateElement(index, elem)
}
