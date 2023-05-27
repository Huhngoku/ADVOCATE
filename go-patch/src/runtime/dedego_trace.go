// DEDEGO-FILE-START

package runtime

type operation int // enum for operation

const (
	opMutLock operation = iota
	opMutRLock
	opMutTryLock
	opMutRTryLock
	opMutUnlock
	opMutRUnlock

	opWgAdd
	opWgWait

	opChanSend
	opChanRecv
	opChanClose
)

type prePost int // enum for pre/post
const (
	pre prePost = iota
	post
	none
)

type dedegoTraceElement interface {
	isDedegoTraceElement()
	toString() string
}

/*
 * Return a string representation of the trace
 * Return:
 * 	string representation of the trace
 */
func TraceToString() string {
	res := ""
	for i, elem := range currentGoRoutine().Trace {
		if i != 0 {
			res += ";"
		}
		res += elem.toString()
	}

	return res
}

/*
 * Add an operation to the trace
 * Args:
 *  elem: element to add to the trace
 * Return:
 * 	index of the element in the trace
 */
func insertIntoTrace(elem dedegoTraceElement) int {
	return currentGoRoutine().addToTrace(elem)
}

/*
 * Print the trace of the current routines
 */
func PrintTrace() {
	routineId := GetRoutineId()
	println("Routine", routineId, ":", TraceToString())
}

// ============================= Mutex =============================

// type to save in the trace for mutexe
type dedegoTraceMutexElement struct {
	id           uint32    // id of the mutex
	counterStart int32     // counter of the operation
	counterEnd   int32     // counter of the operation
	op           operation // operation
	rw           bool      // true if it is a rwmutex
	suc          bool      // success of the operation, only for tryLock
}

func (elem dedegoTraceMutexElement) isDedegoTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element
 */
func (elem dedegoTraceMutexElement) toString() string {
	res := "M," + uint32ToString(elem.id) + "," + int32ToString(elem.counterStart) + ","
	res += int32ToString(elem.counterEnd) + ","

	if elem.rw {
		res += "RW,"
	} else {
		res += "-,"
	}

	switch elem.op {
	case opMutLock:
		res += "L"
	case opMutRLock:
		res += "LR"
	case opMutTryLock:
		res += "T"
	case opMutRTryLock:
		res += "TR"
	case opMutUnlock:
		res += "U"
	case opMutRUnlock:
		res += "UR"
	}

	if elem.op == opMutTryLock || elem.op == opMutRTryLock {
		if elem.suc {
			res += ",succ"
		} else {
			res += ",fail"
		}
	}
	return res
}

/*
 * Add a mutex lock to the trace
 * Args:
 * 	id: id of the mutex
 *  rw: true if it is a rwmutex
 *  r: true if it is a rlock operation
 * Return:
 * 	index of the operation in the trace
 */
func DedegoLock(id uint32, rw bool, r bool) int {
	c := updateCounter()
	op := opMutLock
	if r {
		op = opMutRLock
	}
	elem := dedegoTraceMutexElement{id: id, counterStart: c, op: op, rw: rw, suc: true}
	return insertIntoTrace(elem)
}

/*
 * Add a mutex trylock to the trace
 * Args:
 * 	id: id of the mutex
 * 	rw: true if it is a rwmutex
 * 	r: true if it is a rlock operation
 * Return:
 * 	index of the operation in the trace
 */
func DedegoTryLock(id uint32, rw bool, r bool) int {
	c := updateCounter()
	op := opMutTryLock
	if r {
		op = opMutRTryLock
	}
	elem := dedegoTraceMutexElement{id: id, counterStart: c, op: op, rw: rw}
	return insertIntoTrace(elem)
}

/*
 * Add a mutex unlock to the trace
 * Args:
 * 	id: id of the mutex
 * 	rw: true if it is a runlock
 * 	r: true if it is a rlock operation
 * Return:
 * 	index of the operation in the trace
 */
func DedegoUnlock(id uint32, rw bool, r bool) int {
	c := updateCounter()
	op := opMutUnlock
	if r {
		op = opMutRUnlock
	}
	elem := dedegoTraceMutexElement{id: id, counterStart: c, counterEnd: c, op: op, rw: rw, suc: true}
	return insertIntoTrace(elem)
}

// ============================= WaitGroup ===========================

type dedegoTraceWaitGroupElement struct {
	id           uint32    // id of the waitgroup
	counterStart int32     // counter of the operation
	counterEnd   int32     // counter of the operation
	op           operation // operation
	delta        int       // delta of the waitgroup
	val          int32     // value of the waitgroup after the operation
}

func (elem dedegoTraceWaitGroupElement) isDedegoTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element
 */
func (elem dedegoTraceWaitGroupElement) toString() string {
	res := "W," + uint32ToString(elem.id) + "," + int32ToString(elem.counterStart) + ","
	res += int32ToString(elem.counterEnd) + ","

	switch elem.op {
	case opWgAdd:
		res += "A,"
	case opWgWait:
		res += "W,"
	}
	res += intToString(elem.delta) + "," + int32ToString(elem.val)
	return res
}

/*
 * Add a waitgroup add or done to the trace
 * Args:
 * 	id: id of the waitgroup
 *  delta: delta of the waitgroup
 * 	val: value of the waitgroup after the operation
 * Return:
 * 	index of the operation in the trace
 */
func DedegoAdd(id uint32, delta int, val int32) int {
	c := updateCounter()
	elem := dedegoTraceWaitGroupElement{id: id, counterStart: c, op: opWgWait, delta: delta, val: val}
	return insertIntoTrace(elem)

}

/*
 * Add a waitgroup wait to the trace
 * Args:
 * 	id: id of the waitgroup
 * Return:
 * 	index of the operation in the trace
 */
func DedegoWait(id uint32) int {
	c := updateCounter()
	elem := dedegoTraceWaitGroupElement{id: id, counterStart: c, counterEnd: c, op: opWgWait}
	return insertIntoTrace(elem)
}

// ============================= Channel =============================
// TODO: add channels into chan code

type dedegoTraceChannelElement struct {
	id           uint32    // id of the channel
	counterStart int32     // counter of the start of the operation
	counterEnd   int32     // counter of the end of the operation
	op           operation // operation
	partnerId    uint32    // id of the channel with wich the communication took place
}

func (elem dedegoTraceChannelElement) isDedegoTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element
 */
func (elem dedegoTraceChannelElement) toString() string {
	res := "C" + uint32ToString(elem.id) + "," + int32ToString(elem.counterStart) + ","
	res += int32ToString(elem.counterEnd) + ","

	switch elem.op {
	case opChanSend:
		res += ",S"
	case opChanRecv:
		res += ",R"
	case opChanClose:
		res += ",C"
	}
	res += "," + uint32ToString(elem.partnerId)
	return res
}

/*
 * Add a channel send to the trace
 * Args:
 * 	id: id of the channel
 * Return:
 * 	index of the operation in the trace
 */
func DedegoSend(id uint32) int {
	c := updateCounter()
	elem := dedegoTraceChannelElement{id: id, op: opChanSend, counterStart: c}
	return insertIntoTrace(elem)
}

/*
 * Add a channel recv to the trace
 * Args:
 * 	id: id of the channel
 * 	isPre: true if pre event, post otherwise
 * 	partnerId: id of the channel with wich the communication took place
 * Return:
 * 	index of the operation in the trace
 */
func DedegoRecv(id uint32) int {
	c := updateCounter()
	elem := dedegoTraceChannelElement{id: id, op: opChanRecv, counterStart: c}
	return insertIntoTrace(elem)
}

/*
 * Add a channel close to the trace
 * Args:
 * 	id: id of the channel
 * Return:
 * 	index of the operation in the trace
 */
func DedegoClose(id uint32) int {
	c := updateCounter()
	elem := dedegoTraceChannelElement{id: id, op: opChanClose, counterStart: c}
	return insertIntoTrace(elem)
}

// ============================= Finish ================================
/*
 * Add the end counter to an operation of the trace. For try use DedegoFinishTry.
 * Args:
 * 	index: index of the operation in the trace
 */
func DedegoFinish(index int) {
	c := updateCounter()
	// only needed to fix tests
	if currentGoRoutine() == nil {
		return
	}

	switch elem := currentGoRoutine().Trace[index].(type) {
	case dedegoTraceMutexElement:
		elem.counterEnd = int32(c)
		currentGoRoutine().Trace[index] = elem
	case dedegoTraceWaitGroupElement:
		elem.counterEnd = int32(c)
		currentGoRoutine().Trace[index] = elem
	case dedegoTraceChannelElement:
		elem.counterEnd = int32(c)
		currentGoRoutine().Trace[index] = elem
	}
}

/*
 * Add the end counter to an try operation of the trace
 * Args:
 * 	index: index of the operation in the trace
 * 	suc: true if the try was successful, false otherwise
 */
func DedegoFinishTry(index int, suc bool) {
	c := updateCounter()
	switch elem := currentGoRoutine().Trace[index].(type) {
	case dedegoTraceMutexElement:
		elem.counterEnd = int32(c)
		elem.suc = suc
		currentGoRoutine().Trace[index] = elem
	default:
		panic("DedegoFinishTry called on non mutex")
	}
}

// DEDUGO-FILE-END
