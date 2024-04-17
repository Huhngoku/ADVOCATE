package runtime

// type to save in the trace for mutexe
type advocateMutexElement struct {
	id    uint64    // id of the mutex
	op    Operation // operation
	rw    bool      // true if it is a rwmutex
	suc   bool      // success of the operation, only for tryLock
	file  string    // file where the operation was called
	line  int       // line where the operation was called
	tPre  uint64    // global timer at begin of operation
	tPost uint64    // global timer at end of operation
}

func (elem advocateMutexElement) isAdvocateTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "M,'tre','tPost','id','rw','op','suc','file':'line'"
 *    't' (number): global timer
 *    'id' (number): id of the mutex
 *    'rw' (R/-): R if it is a rwmutex, otherwise -
 *	  'op' (L/R/T/Y/U/N): L if it is a lock, R if it is a rlock, T if it is a trylock, Y if it is a rtrylock, U if it is an unlock, N if it is an runlock
 *	  'suc' (t/f): s if the trylock was successful, f otherwise
 *    'file' (string): file where the operation was called
 *    'line' (number): line where the operation was called
 */
func (elem advocateMutexElement) toString() string {
	res := "M,"
	res += uint64ToString(elem.tPre) + "," + uint64ToString(elem.tPost) + ","
	res += uint64ToString(elem.id) + ","

	if elem.rw {
		res += "R,"
	} else {
		res += "-,"
	}

	switch elem.op {
	case OperationMutexLock, OperationRWMutexLock:
		res += "L"
	case OperationRWMutexRLock:
		res += "R"
	case OperationMutexTryLock, OperationRWMutexTryLock:
		res += "T"
	case OperationRWMutexTryRLock:
		res += "Y"
	case OperationMutexUnlock, OperationRWMutexUnlock:
		res += "U"
	case OperationRWMutexRUnlock:
		res += "N"
	}

	if elem.suc {
		res += ",t"
	} else {
		res += ",f"
	}
	res += "," + elem.file + ":" + intToString(elem.line)
	return res
}

/*
 * Get the operation
 */
func (elem advocateMutexElement) getOperation() Operation {
	return elem.op
}

/*
 * Get the file
 */
func (elem advocateMutexElement) getFile() string {
	return elem.file
}

/*
 * Get the line
 */
func (elem advocateMutexElement) getLine() int {
	return elem.line
}

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
	var op Operation
	if !rw { // Mutex
		if !r { // Lock
			op = OperationMutexLock
		} else { // rLock, invalid case
			panic("Tried to RLock a non-RW Mutex")
		}
	} else { // RWMutex
		if !r { // Lock
			op = OperationRWMutexLock
		} else { // RLock
			op = OperationRWMutexRLock
		}
	}

	_, file, line, _ := Caller(2)
	timer := GetAdvocateCounter()
	elem := advocateMutexElement{id: id, op: op, rw: rw, suc: true,
		file: file, line: line, tPre: timer}
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
	var op Operation
	if !rw { // Mutex
		if !r { // Lock
			op = OperationMutexTryLock
		} else { // rLock, invalid case
			panic("Tried to TryRLock a non-RW Mutex")
		}
	} else { // RWMutex
		if !r { // Lock
			op = OperationRWMutexTryLock
		} else { // RLock
			op = OperationRWMutexTryRLock
		}
	}

	_, file, line, _ := Caller(2)
	timer := GetAdvocateCounter()
	elem := advocateMutexElement{id: id, op: op, rw: rw, file: file,
		line: line, tPre: timer}
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
	var op Operation
	if !rw { // Mutex
		if !r { // Lock
			op = OperationMutexUnlock
		} else { // rLock, invalid case
			panic("Tried to RUnlock a non-RW Mutex")
		}
	} else { // RWMutex
		if !r { // Lock
			op = OperationRWMutexUnlock
		} else { // RLock
			op = OperationRWMutexRUnlock
		}
	}
	_, file, line, _ := Caller(2)
	timer := GetAdvocateCounter()
	elem := advocateMutexElement{id: id, op: op, rw: rw, suc: true,
		file: file, line: line, tPre: timer, tPost: timer}
	return insertIntoTrace(elem)
}

/*
 * AdvocatePost adds the end counter to an operation of the trace.
 * For try use AdvocatePostTry.
 * Also used for wait group
 * Args:
 * 	index: index of the operation in the trace
 * 	c: number of the send
 */
func AdvocatePost(index int) {
	// internal elements are not in the trace
	if index == -1 {
		return
	}

	// only needed to fix tests
	if currentGoRoutine() == nil {
		return
	}

	timer := GetAdvocateCounter()

	switch elem := currentGoRoutine().getElement(index).(type) {
	case advocateMutexElement:
		elem.tPost = timer
		currentGoRoutine().updateElement(index, elem)
	case advocateWaitGroupElement:
		elem.tPost = timer
		currentGoRoutine().updateElement(index, elem)

	default:
		panic("AdvocatePost called on non mutex, waitgroup or channel")
	}
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

	switch elem := currentGoRoutine().getElement(index).(type) {
	case advocateMutexElement:
		elem.suc = suc
		elem.tPost = GetAdvocateCounter()
		currentGoRoutine().updateElement(index, elem)
	default:
		panic("AdvocatePostTry called on non mutex")
	}
}
