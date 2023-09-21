// DEDEGO-FILE-START

package runtime

import (
	at "runtime/internal/atomic"
)

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

type dedegoAtomicMapElem struct {
	addr      uint64
	operation int
}

var dedegoDisabled bool = false
var dedegoAtomicMap map[uint64]dedegoAtomicMapElem = make(map[uint64]dedegoAtomicMapElem)
var dedegoAtomicMapLock mutex

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
 * Return the trace of the routine with id 'id'
 * Args:
 * 	id: id of the routine
 * Return:
 * 	string representation of the trace of the routine
 * 	bool: true if the routine exists, false otherwise
 */
func TraceToStringById(id uint64) (string, bool) {
	lock(DedegoRoutinesLock)
	defer unlock(DedegoRoutinesLock)
	if trace, ok := DedegoRoutines[id]; ok {
		return traceToString(trace), true
	}
	return "", false
}

/*
 * Get the trace of the routine with id 'id'.
 * To minimized the needed ram the trace is sent to the channel 'c' in chunks
 * of 1000 elements.
 * Args:
 * 	id: id of the routine
 * 	c: channel to send the trace to
 *  atomic: it true, the atomic trace is returned
 */
func TraceToStringByIdChannel(id int, c chan<- string) {
	// lock(DedegoRoutinesLock)
	// defer unlock(DedegoRoutinesLock)

	if trace, ok := DedegoRoutines[uint64(id)]; ok {
		res := ""
		for i, elem := range *trace {
			if i != 0 {
				res += ";"
			}
			res += elem.toString()

			if i%1000 == 0 {
				c <- res
				res = ""
			}
		}
		c <- res
	}

}

// }

/*
 * Return the trace of all traces
 * Return:
 * 	string representation of the trace of all routines
 */
func AllTracesToString() string {
	// write warning if projectPath is empty
	res := ""
	lock(DedegoRoutinesLock)
	defer unlock(DedegoRoutinesLock)

	for i := 1; i <= len(DedegoRoutines); i++ {
		res += ""
		trace := DedegoRoutines[uint64(i)]
		if trace == nil {
			panic("Trace is nil")
		}
		res += traceToString(trace) + "\n"

	}
	return res
}

/*
* Print the trace of all routines
 */
func PrintAllTraces() {
	print(AllTracesToString())
}

/*
 * Return the number of routines in the trace
 * Return:
 *	number of routines in the trace
 */
func GetNumberOfRoutines() int {
	lock(DedegoRoutinesLock)
	defer unlock(DedegoRoutinesLock)
	return len(DedegoRoutines)
}

/* Enable the collection of the trace */
func InitAtomics(size int) {
	// link runtime with atomic via channel to receive information about
	// atomic events
	c := make(chan at.AtomicElem, size)
	at.DedegoAtomicLink(c)
	go func() {
		for atomic := range c {
			lock(&dedegoAtomicMapLock)
			dedegoAtomicMap[atomic.Index] = dedegoAtomicMapElem{
				addr:      atomic.Addr,
				operation: atomic.Operation,
			}
			unlock(&dedegoAtomicMapLock)
		}
	}()

	dedegoDisabled = false
}

/* Disable the collection of the trace */
func DisableTrace() {
	at.DedegoAtomicUnlink()
	dedegoDisabled = true
}

// ============================= Routine ===========================

// type to save in the trace for routines
type dedegoTraceSpawnElement struct {
	id    uint64 // id of the routine
	timer uint64 // global timer
}

func (elem dedegoTraceSpawnElement) isDedegoTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "G,'id'"
 *    'id' (number): id of the routine
 */
func (elem dedegoTraceSpawnElement) toString() string {
	return "G," + uint64ToString(elem.timer) + "," + uint64ToString(elem.id)
}

/*
 * Add a routine spawn to the trace
 * Args:
 * 	id: id of the routine
 */
func DedegoSpawn(callerRoutine *DedegoRoutine, newId uint64) {
	timer := GetDedegoCounter()
	callerRoutine.addToTrace(dedegoTraceSpawnElement{id: newId, timer: timer})
}

// ============================= Mutex =============================

// type to save in the trace for mutexe
type dedegoTraceMutexElement struct {
	id    uint64    // id of the mutex
	op    operation // operation
	rw    bool      // true if it is a rwmutex
	suc   bool      // success of the operation, only for tryLock
	file  string    // file where the operation was called
	line  int       // line where the operation was called
	tPre  uint64    // global timer at begin of operation
	tPost uint64    // global timer at end of operation
}

func (elem dedegoTraceMutexElement) isDedegoTraceElement() {}

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
func (elem dedegoTraceMutexElement) toString() string {
	res := "M,"
	res += uint64ToString(elem.tPre) + "," + uint64ToString(elem.tPost) + ","
	res += uint64ToString(elem.id) + ","

	if elem.rw {
		res += "R,"
	} else {
		res += "-,"
	}

	switch elem.op {
	case opMutLock:
		res += "L"
	case opMutRLock:
		res += "R"
	case opMutTryLock:
		res += "T"
	case opMutRTryLock:
		res += "Y"
	case opMutUnlock:
		res += "U"
	case opMutRUnlock:
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
 * Add a mutex lock to the trace
 * Args:
 * 	id: id of the mutex
 *  rw: true if it is a rwmutex
 *  r: true if it is a rlock operation
 * Return:
 * 	index of the operation in the trace
 */
func DedegoMutexLockPre(id uint64, rw bool, r bool) int {
	op := opMutLock
	if r {
		op = opMutRLock
	}
	_, file, line, _ := Caller(2)
	timer := GetDedegoCounter()
	elem := dedegoTraceMutexElement{id: id, op: op, rw: rw, suc: true,
		file: file, line: line, tPre: timer}
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
func DedegoMutexLockTry(id uint64, rw bool, r bool) int {
	op := opMutTryLock
	if r {
		op = opMutRTryLock
	}
	_, file, line, _ := Caller(2)
	timer := GetDedegoCounter()
	elem := dedegoTraceMutexElement{id: id, op: op, rw: rw, file: file,
		line: line, tPre: timer}
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
func DedegoUnlockPre(id uint64, rw bool, r bool) int {
	op := opMutUnlock
	if r {
		op = opMutRUnlock
	}
	_, file, line, _ := Caller(2)
	timer := GetDedegoCounter()
	elem := dedegoTraceMutexElement{id: id, op: op, rw: rw, suc: true,
		file: file, line: line, tPre: timer, tPost: timer}
	return insertIntoTrace(elem)
}

/*
 * Add the end counter to an operation of the trace. For try use DedegoPostTry.
 *   Also used for wait group
 * Args:
 * 	index: index of the operation in the trace
 * 	c: number of the send
 */
func DedegoPost(index int) {
	// internal elements are not in the trace
	if index == -1 {
		return
	}

	// only needed to fix tests
	if currentGoRoutine() == nil {
		return
	}

	timer := GetDedegoCounter()

	switch elem := currentGoRoutine().getElement(index).(type) {
	case dedegoTraceMutexElement:
		elem.tPost = timer
		currentGoRoutine().updateElement(index, elem)
	case dedegoTraceWaitGroupElement:
		elem.tPost = timer
		currentGoRoutine().updateElement(index, elem)

	default:
		panic("DedegoPost called on non mutex, waitgroup or channel")
	}
}

/*
 * Add the end counter to an try operation of the trace
 * Args:
 * 	index: index of the operation in the trace
 * 	suc: true if the try was successful, false otherwise
 */
func DedegoPostTry(index int, suc bool) {
	// internal elements are not in the trace
	if index == -1 {
		return
	}

	switch elem := currentGoRoutine().getElement(index).(type) {
	case dedegoTraceMutexElement:
		elem.suc = suc
		elem.tPost = GetDedegoCounter()
		currentGoRoutine().updateElement(index, elem)
	default:
		panic("DedegoPostTry called on non mutex")
	}
}

// ============================= WaitGroup ===========================

type dedegoTraceWaitGroupElement struct {
	id    uint64    // id of the waitgroup
	op    operation // operation
	delta int       // delta of the waitgroup
	val   int32     // value of the waitgroup after the operation
	file  string    // file where the operation was called
	line  int       // line where the operation was called
	tPre  uint64    // global timer
	tPost uint64    // global timer
}

func (elem dedegoTraceWaitGroupElement) isDedegoTraceElement() {}

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
func (elem dedegoTraceWaitGroupElement) toString() string {
	res := "W,"
	res += uint64ToString(elem.tPre) + "," + uint64ToString(elem.tPost) + ","
	res += uint64ToString(elem.id) + ","
	switch elem.op {
	case opWgAdd:
		res += "A,"
	case opWgWait:
		res += "W,"
	}

	res += intToString(elem.delta) + "," + int32ToString(elem.val)
	res += "," + elem.file + ":" + intToString(elem.line)
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
func DedegoWaitGroupAdd(id uint64, delta int, val int32) int {
	var file string
	var line int
	if delta > 0 {
		_, file, line, _ = Caller(2)
	} else {
		_, file, line, _ = Caller(3)
	}
	timer := GetDedegoCounter()
	elem := dedegoTraceWaitGroupElement{id: id, op: opWgWait, delta: delta,
		val: val, file: file, line: line, tPre: timer, tPost: timer}
	return insertIntoTrace(elem)

}

/*
 * Add a waitgroup wait to the trace
 * Args:
 * 	id: id of the waitgroup
 * Return:
 * 	index of the operation in the trace
 */
func DedegoWaitGroupWaitPre(id uint64) int {
	_, file, line, _ := Caller(2)
	timer := GetDedegoCounter()
	elem := dedegoTraceWaitGroupElement{id: id, op: opWgWait, file: file,
		line: line, tPre: timer}
	return insertIntoTrace(elem)
}

// ============================= Channel =============================

type dedegoTraceChannelElement struct {
	id     uint64    // id of the channel
	op     operation // operation
	qSize  uint32    // size of the channel, 0 for unbuffered
	opId   uint64    // id of the operation
	file   string    // file where the operation was called
	line   int       // line where the operation was called
	tPre   uint64    // global timer before the operation
	tPost  uint64    // global timer after the operation
	closed bool      // true if the channel operation was finished, because the channel was closed at another routine
}

func (elem dedegoTraceChannelElement) isDedegoTraceElement() {}

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
func (elem dedegoTraceChannelElement) toString() string {
	return elem.toStringSep(",", true)
}

/*
* Get a string representation of the element given a separator
* Args:
* 	sep: separator to use
* 	showPos: true if the position of the operation should be shown
* Return:
* 	string representation of the element
 */
func (elem dedegoTraceChannelElement) toStringSep(sep string, showPos bool) string {
	res := "C" + sep
	res += uint64ToString(elem.tPre) + sep + uint64ToString(elem.tPost) + sep
	res += uint64ToString(elem.id) + sep

	switch elem.op {
	case opChanSend:
		res += "S"
	case opChanRecv:
		res += "R"
	case opChanClose:
		res += "C"
	}

	if elem.closed {
		res += sep + "t"
	} else {
		res += sep + "f"
	}

	res += sep + uint64ToString(elem.opId)
	res += sep + uint32ToString(elem.qSize)
	if showPos {
		res += sep + elem.file + ":" + intToString(elem.line)
	}
	return res
}

/*
 * Add a channel send to the trace. If the channel send was created by an atomic
 * operation, add this to the trace as well
 * Args:
 * 	id: id of the channel
 * 	opId: id of the operation
 * 	qSize: size of the channel, 0 for unbuffered
 * Return:
 * 	index of the operation in the trace, return -1 if it is a atomic operation
 */
var dedegoCounterAtomic uint64

func DedegoChanSendPre(id uint64, opId uint64, qSize uint) int {
	_, file, line, _ := Caller(3)
	// internal channels to record atomic operations
	if isSuffix(file, "dedego_atomic.go") {
		dedegoCounterAtomic += 1
		DedegoAtomic(dedegoCounterAtomic)
		// they are not recorded in the trace
		return -1
	}
	timer := GetDedegoCounter()
	elem := dedegoTraceChannelElement{id: id, op: opChanSend, opId: opId,
		file: file, line: line, tPre: timer, qSize: uint32(qSize)}
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
 * Add a channel recv to the trace
 * Args:
 * 	id: id of the channel
 * 	opId: id of the operation
 * 	qSize: size of the channel
 * Return:
 * 	index of the operation in the trace
 */
func DedegoChanRecvPre(id uint64, opId uint64, qSize uint) int {
	_, file, line, _ := Caller(3)
	// do not record channel operation of internal channel to record atomic operations
	if isSuffix(file, "dedego_trace.go") {
		return -1
	}

	timer := GetDedegoCounter()
	elem := dedegoTraceChannelElement{id: id, op: opChanRecv, opId: opId,
		file: file, line: line, tPre: timer, qSize: uint32(qSize)}
	return insertIntoTrace(elem)
}

/*
 * Add a channel close to the trace
 * Args:
 * 	id: id of the channel
 * Return:
 * 	index of the operation in the trace
 */
func DedegoChanClose(id uint64, qSize uint) int {
	_, file, line, _ := Caller(2)
	timer := GetDedegoCounter()
	elem := dedegoTraceChannelElement{id: id, op: opChanClose, file: file,
		line: line, tPre: timer, tPost: timer, qSize: uint32(qSize)}
	return insertIntoTrace(elem)
}

/*
 * Set the operation as successfully finished
 * Args:
 * 	index: index of the operation in the trace
 */
func DedegoChanPost(index int) {
	if index == -1 {
		return
	}
	elem := currentGoRoutine().getElement(index).(dedegoTraceChannelElement)
	elem.tPost = GetDedegoCounter()
	currentGoRoutine().updateElement(index, elem)
}

func DedegoChanPostCausedByClose(index int) {
	if index == -1 {
		return
	}
	elem := currentGoRoutine().getElement(index).(dedegoTraceChannelElement)
	elem.closed = true
	currentGoRoutine().updateElement(index, elem)
}

// ============================= Select ==============================

type dedegoTraceSelectElement struct {
	tPre    uint64                      // global timer before the operation
	tPost   uint64                      // global timer after the operation
	id      uint64                      // id of the select
	cases   []dedegoTraceChannelElement // cases of the select
	chosen  int                         // index of the chosen case in cases (0 indexed, -1 for default)
	nsend   int                         // number of send cases
	defa    bool                        // set true if a default case exists
	defaSel bool                        // set true if a default case was chosen
	file    string                      // file where the operation was called
	line    int                         // line where the operation was called
}

func (elem dedegoTraceSelectElement) isDedegoTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "S,'tPre','tPost','id','cases','opId','file':'line'"
 *    'tPre' (number): global timer before the operation
 *    'tPost' (number): global timer after the operation
 *    'id' (number): id of the mutex
 *	  'cases' (string): cases of the select, d for default
 *    'chosen' (number): index of the chosen case in cases (0 indexed, -1 for default)
 *	  'opId' (number): id of the operation on the channel
 *    'file' (string): file where the operation was called
 *    'line' (number): line where the operation was called
 */
func (elem dedegoTraceSelectElement) toString() string {
	res := "S,"
	res += uint64ToString(elem.tPre) + "," + uint64ToString(elem.tPost) + ","
	res += uint64ToString(elem.id) + ","

	for i, ca := range elem.cases {
		if i != 0 {
			res += "~"
		}
		res += ca.toStringSep(".", false)
	}

	if elem.defa {
		if len(elem.cases) != 0 {
			res += "~"
		}
		if elem.defaSel {
			res += "D"
		} else {
			res += "d"
		}
	}

	res += "," + elem.file + ":" + intToString(elem.line)
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
func DedegoSelectPre(cases *[]scase, nsends int, block bool) int {
	timer := GetDedegoCounter()
	if cases == nil {
		return -1
	}

	id := GetDedegoObjectId()
	caseElements := make([]dedegoTraceChannelElement, len(*cases))
	_, file, line, _ := Caller(2)

	for i, ca := range *cases {
		caseElements[i] = dedegoTraceChannelElement{id: ca.c.id,
			qSize: uint32(ca.c.dataqsiz), tPre: timer}
	}

	elem := dedegoTraceSelectElement{id: id, cases: caseElements, nsend: nsends,
		defa: !block, file: file, line: line, tPre: timer}
	return insertIntoTrace(elem)
}

/*
 * Post event for select in case of an non-default case
 * Args:
 * 	index: index of the operation in the trace
 * 	c: channel of the chosen case
 * 	chosenIndex: index of the chosen case in the select
 * 	lockOrder: order of the locks
 * 	rClosed: true if the channel was closed at another routine
 */
func DedegoSelectPost(index int, c *hchan, chosenIndex int,
	lockOrder []uint16, rClosed bool) {
	if index == -1 || c == nil {
		return
	}

	timer := GetDedegoCounter()
	elem := currentGoRoutine().getElement(index).(dedegoTraceSelectElement)

	elem.chosen = chosenIndex
	elem.tPost = timer

	for i, op := range lockOrder {
		opChan := opChanRecv
		if op < uint16(elem.nsend) {
			opChan = opChanSend
		}
		elem.cases[i].op = opChan
	}

	if chosenIndex == -1 { // default case
		elem.defaSel = true
	} else {
		elem.cases[chosenIndex].tPost = timer
		elem.cases[chosenIndex].closed = rClosed
		send := false
		if elem.cases[chosenIndex].op == opChanSend {
			send = true
		} else if elem.cases[chosenIndex].op == opChanRecv {
			send = false
		} else {
			panic("DedegoSelectPostNonDefault called on non sending channel op")
		}
		// set oId
		if send {
			c.numberSend++
			elem.cases[chosenIndex].opId = c.numberSend
		} else {
			c.numberRecv++
			elem.cases[chosenIndex].opId = c.numberRecv
		}

	}

	currentGoRoutine().updateElement(index, elem)
}

/*
* Add a new select element to the trace if the select has exactly one
* non-default case and a default case
* Args:
* 	c: channel of the non-default case
* 	send: true if the non-default case is a send, false otherwise
* Return:
* 	index of the operation in the trace
 */
func DedegoSelectPreOneNonDef(c *hchan, send bool) int {
	if c == nil {
		return -1
	}

	id := GetDedegoObjectId()
	timer := GetDedegoCounter()

	opChan := opChanRecv
	if send {
		opChan = opChanSend
	}

	caseElements := make([]dedegoTraceChannelElement, 1)
	caseElements[0] = dedegoTraceChannelElement{id: c.id,
		qSize: uint32(c.dataqsiz), tPre: timer, op: opChan}

	nSend := 0
	if send {
		nSend = 1
	}

	_, file, line, _ := Caller(2)
	elem := dedegoTraceSelectElement{id: id, cases: caseElements, nsend: nSend,
		defa: true, file: file, line: line, tPre: timer}
	return insertIntoTrace(elem)
}

/*
 * Add the selected case for a select with one non-default and one default case
 * Args:
 * 	index: index of the operation in the trace
 * 	res: 0 for the non-default case, -1 for the default case
 */
func DedegoSelectPostOneNonDef(index int, res bool, oId uint64) {
	if index == -1 {
		return
	}

	timer := GetDedegoCounter()
	elem := currentGoRoutine().getElement(index).(dedegoTraceSelectElement)

	if res {
		elem.chosen = 0
		elem.cases[0].tPost = timer
		elem.cases[0].opId = oId
	} else {
		elem.chosen = -1
		elem.defaSel = true
	}
	elem.tPost = timer

	currentGoRoutine().updateElement(index, elem)
}

// ============================= Atomic ================================
type dedegoTraceAtomicElement struct {
	timer     uint64 // global timer
	index     uint64 // index of the atomic event in dedegoAtomicMap
	operation int    // type of operation
}

func (elem dedegoTraceAtomicElement) isDedegoTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "A,'addr'"
 *    'addr' (number): address of the atomic variable
 */
// enum for atomic operation, must be the same as in dedego_atomic.go
const (
	LoadOp = iota
	StoreOp
	AddOp
	SwapOp
	CompSwapOp
)

func (elem dedegoTraceAtomicElement) toString() string {
	lock(&dedegoAtomicMapLock)
	mapElement := dedegoAtomicMap[elem.index]
	res := "A," + uint64ToString(elem.timer) + "," +
		uint64ToString(mapElement.addr) + ","
	switch mapElement.operation {
	case LoadOp:
		res += "L"
	case StoreOp:
		res += "S"
	case AddOp:
		res += "A"
	case SwapOp:
		res += "W"
	case CompSwapOp:
		res += "C"
	default:
		res += "U"
	}
	unlock(&dedegoAtomicMapLock)
	return res
}

func DedegoAtomic(index uint64) {
	timer := GetDedegoCounter()
	elem := dedegoTraceAtomicElement{index: index, timer: timer}
	insertIntoTrace(elem)
}

// ======================= Once ============================
type dedegoOnceElement struct {
	tpre  uint64 // global timer at the beginning of the execution
	tpost uint64 // global timer at the end of the execution
	id    uint64 // id of the once
	suc   bool   // true if the do on the once was called for the first time
	file  string // file where the operation was called
	line  int    // line where the operation was called
}

func (elem dedegoOnceElement) isDedegoTraceElement() {}

/*
 * Return a string representation of the element
 * Return:
 * 	string representation of the element "O,'tpre','tpost','id','suc','file':'line"
 *    'tpre' (number): global timer at the beginning of the execution
 *    'tpost' (number): global timer at the end of the execution
 *    'id' (number): id of the once
 *    'suc' (t/f): true if the do on the once was called for the first time, false otherwise
 *    'file' (string): file where the operation was called
 *    'line' (string): line where the operation was called
 */
func (elem dedegoOnceElement) toString() string {
	res := "O,"
	res += uint64ToString(elem.tpre) + ","
	res += uint64ToString(elem.tpost) + ","
	res += uint64ToString(elem.id) + ","
	if elem.suc {
		res += "t"
	} else {
		res += "f"
	}
	res += "," + elem.file + ":" + intToString(elem.line)
	return res
}

func DedegoOncePre(id uint64) int {
	_, file, line, _ := Caller(2)
	timer := GetDedegoCounter()
	elem := dedegoOnceElement{id: id, tpre: timer, file: file, line: line}
	return insertIntoTrace(elem)
}

func DedegoOncePost(index int, suc bool) {
	if index == -1 {
		return
	}
	timer := GetDedegoCounter()
	elem := currentGoRoutine().getElement(index).(dedegoOnceElement)
	elem.tpost = timer
	elem.suc = suc

	currentGoRoutine().updateElement(index, elem)
}

// DEDEGO-FILE-END
