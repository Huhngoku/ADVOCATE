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
func CurrentTraceToString() string {
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
 * Return a string representation of the trace
 * Args:
 * 	trace: trace to convert to string
 * Return:
 * 	string representation of the trace
 */
func traceToString(trace *[]dedegoTraceElement) string {
	res := ""
	for i, elem := range *trace {
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
	println("Routine", routineId, ":", CurrentTraceToString())
}

/*
 * Return the trace of all traces
 * Args:
 * 	all: if true, return the trace of all routines, else return only the non empty traces
 * Return:
 * 	string representation of the trace of all routines
 */
func AllTracesToString(all bool) string {
	res := ""
	lock(DedegoRoutinesLock)
	for id, trace := range DedegoRoutines {
		traceString := traceToString(trace)
		if all || traceString != "" {
			res += uint64ToString(id) + ":" + traceString + "\n"
		}
	}
	unlock(DedegoRoutinesLock)
	return res
}

/*
 * Print the trace of all routines
 * Args:
 * 	all: if true, print the trace of all routines, else print only the non empty traces
 */
func PrintAllTraces(all bool) {
	print(AllTracesToString(all))
}

// ============================= Routine ===========================

// type to save in the trace for routines
type dedegoTraceSpawnElement struct {
	id uint64 // id of the routine
}

func (elem dedegoTraceSpawnElement) isDedegoTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "S,'id'"
 *    'id' (number): id of the routine
 */
func (elem dedegoTraceSpawnElement) toString() string {
	return "S," + uint64ToString(elem.id)
}

/*
 * Add a routine spawn to the trace
 * Args:
 * 	id: id of the routine
 */
func DedegoSpawn(id uint64) {
	insertIntoTrace(dedegoTraceSpawnElement{id: id})
}

// ============================= Mutex =============================

// type to save in the trace for mutexe
type dedegoTraceMutexElement struct {
	id   uint64    // id of the mutex
	exec bool      // set true if the operation was successfully finished
	op   operation // operation
	rw   bool      // true if it is a rwmutex
	suc  bool      // success of the operation, only for tryLock
}

func (elem dedegoTraceMutexElement) isDedegoTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "M,'id','rw','op','exec','suc'"
 *    'id' (number): id of the mutex
 *    'rw' (R/-): R if it is a rwmutex, otherwise -
 *	  'op' (L/LR/T/TR/U/UR): L if it is a lock, LR if it is a rlock, T if it is a trylock, TR if it is a rtrylock, U if it is an unlock, UR if it is an runlock
 *	  'exec' (e/o): e if the operation was successfully finished, o otherwise
 *	  'suc' (s/f): s if the trylock was successful, f otherwise
 */
func (elem dedegoTraceMutexElement) toString() string {
	res := "M," + uint64ToString(elem.id) + ","

	if elem.rw {
		res += "R,"
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

	if elem.exec {
		res += ",e"
	} else {
		res += ",o"
	}

	if elem.op == opMutTryLock || elem.op == opMutRTryLock {
		if elem.suc {
			res += ",s"
		} else {
			res += ",f"
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
func DedegoLock(id uint64, rw bool, r bool) int {
	op := opMutLock
	if r {
		op = opMutRLock
	}
	elem := dedegoTraceMutexElement{id: id, op: op, rw: rw, suc: true}
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
func DedegoTryLock(id uint64, rw bool, r bool) int {
	op := opMutTryLock
	if r {
		op = opMutRTryLock
	}
	elem := dedegoTraceMutexElement{id: id, op: op, rw: rw}
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
func DedegoUnlock(id uint64, rw bool, r bool) int {
	op := opMutUnlock
	if r {
		op = opMutRUnlock
	}
	elem := dedegoTraceMutexElement{id: id, op: op, rw: rw, suc: true}
	return insertIntoTrace(elem)
}

// ============================= WaitGroup ===========================

type dedegoTraceWaitGroupElement struct {
	id    uint64    // id of the waitgroup
	exec  bool      // set true if the operation was successfully finished
	op    operation // operation
	delta int       // delta of the waitgroup
	val   int32     // value of the waitgroup after the operation
}

func (elem dedegoTraceWaitGroupElement) isDedegoTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "W,'id','op','exec','delta','val'"
 *    'id' (number): id of the mutex
 *	  'op' (A/W): A if it is an add or Done, W if it is a wait
 *	  'exec' (e/o): e if the operation was successfully finished, o otherwise
 *	  'delta' (number): delta of the waitgroup, positive for add, negative for done, 0 for wait
 *	  'val' (number): value of the waitgroup after the operation
 */
func (elem dedegoTraceWaitGroupElement) toString() string {
	res := "W," + uint64ToString(elem.id) + ","
	switch elem.op {
	case opWgAdd:
		res += "A,"
	case opWgWait:
		res += "W,"
	}

	if elem.exec {
		res += "e,"
	} else {
		res += "o,"
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
func DedegoAdd(id uint64, delta int, val int32) int {
	elem := dedegoTraceWaitGroupElement{id: id, op: opWgWait, delta: delta, val: val}
	return insertIntoTrace(elem)

}

/*
 * Add a waitgroup wait to the trace
 * Args:
 * 	id: id of the waitgroup
 * Return:
 * 	index of the operation in the trace
 */
func DedegoWait(id uint64) int {
	elem := dedegoTraceWaitGroupElement{id: id, op: opWgWait}
	return insertIntoTrace(elem)
}

// ============================= Channel =============================

type dedegoTraceChannelElement struct {
	id   uint64    // id of the channel
	exec bool      // set true if the operation was successfully finished
	op   operation // operation
	opId uint64    // id of the operation
}

func (elem dedegoTraceChannelElement) isDedegoTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "C,'id','op','exec','pId'"
 *    'id' (number): id of the mutex
 *	  'op' (S/R/C): S if it is a send, R if it is a receive, C if it is a close
 *	  'exec' (e/o): e if the operation was successfully finished, o otherwise
 *	  'pId' (number): id of the channel with wich the communication took place
 */
func (elem dedegoTraceChannelElement) toString() string {
	res := "C," + uint64ToString(elem.id) + ","

	switch elem.op {
	case opChanSend:
		res += "S"
	case opChanRecv:
		res += "R"
	case opChanClose:
		res += "C"
	}

	if elem.exec {
		res += ",e"
	} else {
		res += ",o"
	}

	res += "," + uint64ToString(elem.opId)
	return res
}

/*
 * Add a channel send to the trace
 * Args:
 * 	id: id of the channel
 * 	opId: id of the operation
 * Return:
 * 	index of the operation in the trace
 */
func DedegoSend(id uint64, opId uint64) int {
	elem := dedegoTraceChannelElement{id: id, op: opChanSend, opId: opId}
	return insertIntoTrace(elem)
}

/*
 * Add a channel recv to the trace
 * Args:
 * 	id: id of the channel
 * 	opId: id of the operation
 * Return:
 * 	index of the operation in the trace
 */
func DedegoRecv(id uint64, opId uint64) int {
	elem := dedegoTraceChannelElement{id: id, op: opChanRecv, opId: opId}
	return insertIntoTrace(elem)
}

/*
 * Add a channel close to the trace
 * Args:
 * 	id: id of the channel
 * Return:
 * 	index of the operation in the trace
 */
func DedegoClose(id uint64) int {
	elem := dedegoTraceChannelElement{id: id, op: opChanClose}
	return insertIntoTrace(elem)
}

// ============================= Select ==============================

type dedegoTraceSelectElement struct {
	id     uint64   // id of the select
	cases  []string // cases of the select
	send   []bool   // true if the case is a send, false if it is a receive
	nsend  int      // number of send cases
	chosen int      // index of the chosen case
	exec   bool     // set true if the operation was successfully finished
	defa   bool     // set true if a default case exists
}

func (elem dedegoTraceSelectElement) isDedegoTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "S,'id','cases','chosen','exec'"
 *    'id' (number): id of the mutex
 *	  'cases' (string): cases of the select, id and r/s, separated by '.', d for default
 *    'chosen' (number): index of the chosen case in cases (0 indexed, -1 for default)
 *	  'exec' (e/o): e if the operation was successfully finished, o otherwise
 */
func (elem dedegoTraceSelectElement) toString() string {
	res := "S," + uint64ToString(elem.id) + ","
	for i, ca := range elem.cases {
		if i != 0 {
			res += "."
		}
		res += ca
		if elem.send[i] {
			res += "s"
		} else {
			res += "r"
		}
	}
	if elem.defa {
		if len(elem.cases) != 0 {
			res += "."
		}
		res += "d"
	}

	res += "," + intToString(elem.chosen)
	if elem.exec {
		res += ",e"
	} else {
		res += ",o"
	}
	return res
}

/*
 * Add a select to the trace
 * Args:
 * 	cases: cases of the select
 * 	nsends: number of send cases
 * 	block: true if the select is blocking (has no default), false otherwise
 * Return:
 * 	index of the operation in the trace
 */
func DedegoSelect(cases *[]scase, nsends int, block bool) int {
	id := GetDedegoObjectId()
	casesStr := make([]string, len(*cases))
	for i, ca := range *cases {
		casesStr[i] = uint64ToString(ca.c.id)
	}

	elem := dedegoTraceSelectElement{id: id, cases: casesStr, nsend: nsends, defa: !block}
	return insertIntoTrace(elem)
}

// ============================= Finish ================================
/*
 * Add the end counter to an operation of the trace. For try use DedegoFinishTry.
 * Args:
 * 	index: index of the operation in the trace
 */
func DedegoFinish(index int) {
	// only needed to fix tests
	if currentGoRoutine() == nil {
		return
	}

	switch elem := currentGoRoutine().Trace[index].(type) {
	case dedegoTraceMutexElement:
		elem.exec = true
		currentGoRoutine().Trace[index] = elem
	case dedegoTraceWaitGroupElement:
		elem.exec = true
		currentGoRoutine().Trace[index] = elem
	case dedegoTraceChannelElement:
		elem.exec = true
		currentGoRoutine().Trace[index] = elem
	default:
		panic("DedegoFinish called on non mutex, waitgroup or channel")
	}
}

/*
 * Add the end counter to an try operation of the trace
 * Args:
 * 	index: index of the operation in the trace
 * 	suc: true if the try was successful, false otherwise
 */
func DedegoFinishTry(index int, suc bool) {
	switch elem := currentGoRoutine().Trace[index].(type) {
	case dedegoTraceMutexElement:
		elem.exec = true
		elem.suc = suc
		currentGoRoutine().Trace[index] = elem
	default:
		panic("DedegoFinishTry called on non mutex")
	}
}

/*
 * Add the chosen case to the select
 * Args:
 * 	index: index of the operation in the trace
 * 	chosen: index of the chosen case
 */
func DedegoFinishSelect1(index int, chosen int) {
	elem := currentGoRoutine().Trace[index].(dedegoTraceSelectElement)

	elem.exec = true
	elem.chosen = chosen
	currentGoRoutine().Trace[index] = elem
}

/*
 * Add the lock order to the select
 * Args:
 * 	index: index of the operation in the trace
 * 	lockOrder: lock order of the select
 */
func DedegoFinishSelect2(index int, lockOrder []uint16) {
	elem := currentGoRoutine().Trace[index].(dedegoTraceSelectElement)
	send := make([]bool, len(lockOrder))
	for i, lo := range lockOrder {
		send[i] = (lo < uint16(elem.nsend))
	}
	elem.send = send
	currentGoRoutine().Trace[index] = elem
}

// DEDUGO-FILE-END
