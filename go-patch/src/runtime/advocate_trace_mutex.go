package runtime

/*
 * AdvocateMutexLockPre adds a mutex lock to the trace
 * Args:
 * 	id: id of the mutex
 *  rw: true if it is a rwmutex
 *  r: true if it is a rlock operation
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateMutexLockPre(id uint64, rw bool, r bool) int {
	var op string
	var rwStr string
	if !rw { // Mutex
		rwStr = "-"
		if !r { // Lock
			op = "L"
		} else { // rLock, invalid case
			panic("Tried to RLock a non-RW Mutex")
		}
	} else { // RWMutex
		rwStr = "R"
		if !r { // Lock
			op = "L"
		} else { // RLock
			op = "R"
		}
	}

	_, file, line, _ := Caller(2)
	timer := GetAdvocateCounter()

	elem := "M," + uint64ToString(timer) + ",0," + uint64ToString(id) + "," +
		rwStr + "," + op + "," + file + ":" + uint64ToString(uint64(line))

	return insertIntoTrace(elem)
}

/*
 * AdvocateMutexLockTry adds a mutex trylock to the trace
 * Args:
 * 	id: id of the mutex
 * 	rw: true if it is a rwmutex
 * 	r: true if it is a rlock operation
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateMutexLockTry(id uint64, rw bool, r bool) int {
	var op string
	var rwStr string
	if !rw { // Mutex
		rwStr = "-"
		if !r { // Lock
			op = "T"
		} else { // rLock, invalid case
			panic("Tried to TryRLock a non-RW Mutex")
		}
	} else { // RWMutex
		rwStr = "R"
		if !r { // Lock
			op = "T"
		} else { // RLock
			op = "Y"
		}
	}

	_, file, line, _ := Caller(2)
	timer := GetAdvocateCounter()

	elem := "M," + uint64ToString(timer) + ",0," + uint64ToString(id) + "," +
		rwStr + "," + op + "," + file + ":" + uint64ToString(uint64(line))

	return insertIntoTrace(elem)
}

/*
 * AdvocateUnlockPre adds a mutex unlock to the trace
 * Args:
 * 	id: id of the mutex
 * 	rw: true if it is a runlock
 * 	r: true if it is a rlock operation
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateUnlockPre(id uint64, rw bool, r bool) int {
	var op string
	var rwStr string
	if !rw { // Mutex
		rwStr = "-"
		if !r { // Lock
			op = "U"
		} else { // rLock, invalid case
			panic("Tried to RUnlock a non-RW Mutex")
		}
	} else { // RWMutex
		rwStr = "R"
		if !r { // Lock
			op = "U"
		} else { // RLock
			op = "N"
		}
	}
	_, file, line, _ := Caller(2)
	timer := GetAdvocateCounter()

	elem := "M," + uint64ToString(timer) + ",0," + uint64ToString(id) + "," +
		rwStr + "," + op + "," + file + ":" + uint64ToString(uint64(line))

	return insertIntoTrace(elem)
}

/*
 * AdvocateMutexPost adds the end counter to an operation of the trace.
 * For try use AdvocatePostTry.
 * Also used for wait group
 * Args:
 * 	index: index of the operation in the trace
 * 	c: number of the send
 */
func AdvocateMutexPost(index int) {
	// internal elements are not in the trace
	if index == -1 {
		return
	}

	// only needed to fix tests

	if currentGoRoutine() == nil {
		return
	}

	timer := GetAdvocateCounter()

	elem := currentGoRoutine().getElement(index)
	split := splitStringAtCommas(elem, []int{2, 3})
	split[1] = uint64ToString(timer)
	elem = mergeString(split)

	currentGoRoutine().updateElement(index, elem)
}

/*
 * AdvocatePostTry adds the end counter to an try operation of the trace
 * Args:
 * 	index: index of the operation in the trace
 * 	suc: true if the try was successful, false otherwise
 */
func AdvocatePostTry(index int, suc bool) {
	// internal elements are not in the trace
	if index == -1 {
		return
	}

	timer := GetAdvocateCounter()

	elem := currentGoRoutine().getElement(index)
	split := splitStringAtCommas(elem, []int{2, 3, 6, 7})

	split[1] = uint64ToString(timer)
	split[3] = boolToString(suc)

	elem = mergeString(split)

	currentGoRoutine().updateElement(index, elem)

}
