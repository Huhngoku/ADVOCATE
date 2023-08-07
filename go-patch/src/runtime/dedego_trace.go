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
	getFile() string
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
 * Get the file where the element was called
 * Return:
 * 	empty string
 */
func (elem dedegoTraceSpawnElement) getFile() string {
	return ""
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
	id        uint64    // id of the mutex
	exec      bool      // set true if the operation was successfully finished
	op        operation // operation
	rw        bool      // true if it is a rwmutex
	suc       bool      // success of the operation, only for tryLock
	file      string    // file where the operation was called
	line      int       // line where the operation was called
	timerPre  uint64    // global timer at begin of operation
	timerPost uint64    // global timer at end of operation
}

func (elem dedegoTraceMutexElement) isDedegoTraceElement() {}

/*
 * Get the file where the element was called
 * Return:
 * 	file where the element was called
 */
func (elem dedegoTraceMutexElement) getFile() string {
	return elem.file
}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "M,'tre','tpost','id','rw','op','exec','suc','file':'line'"
 *    't' (number): global timer
 *    'id' (number): id of the mutex
 *    'rw' (R/-): R if it is a rwmutex, otherwise -
 *	  'op' (L/LR/T/TR/U/UR): L if it is a lock, LR if it is a rlock, T if it is a trylock, TR if it is a rtrylock, U if it is an unlock, UR if it is an runlock
 *	  'exec' (e/o): e if the operation was successfully finished, o otherwise
 *	  'suc' (s/f): s if the trylock was successful, f otherwise
 *    'file' (string): file where the operation was called
 *    'line' (number): line where the operation was called
 */
func (elem dedegoTraceMutexElement) toString() string {
	res := "M,"
	res += uint64ToString(elem.timerPre) + "," + uint64ToString(elem.timerPost) + ","
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

	if elem.suc {
		res += ",s"
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
		file: file, line: line, timerPre: timer}
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
		line: line, timerPre: timer}
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
		file: file, line: line, timerPre: timer, timerPost: timer}
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

	lock(currentGoRoutine().lock)
	defer unlock(currentGoRoutine().lock)

	switch elem := currentGoRoutine().Trace[index].(type) {
	case dedegoTraceMutexElement:
		elem.exec = true
		elem.timerPost = timer
		currentGoRoutine().Trace[index] = elem
	case dedegoTraceWaitGroupElement:
		elem.exec = true
		elem.timerPost = timer
		currentGoRoutine().Trace[index] = elem

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

	lock(currentGoRoutine().lock)
	defer unlock(currentGoRoutine().lock)

	switch elem := currentGoRoutine().Trace[index].(type) {
	case dedegoTraceMutexElement:
		elem.exec = true
		elem.suc = suc
		elem.timerPost = GetDedegoCounter()
		currentGoRoutine().Trace[index] = elem
	default:
		panic("DedegoPostTry called on non mutex")
	}
}

// ============================= WaitGroup ===========================

type dedegoTraceWaitGroupElement struct {
	id        uint64    // id of the waitgroup
	exec      bool      // set true if the operation was successfully finished
	op        operation // operation
	delta     int       // delta of the waitgroup
	val       int32     // value of the waitgroup after the operation
	file      string    // file where the operation was called
	line      int       // line where the operation was called
	timerPre  uint64    // global timer
	timerPost uint64    // global timer
}

func (elem dedegoTraceWaitGroupElement) isDedegoTraceElement() {}

/*
 * Get the file where the element was called
 * Return:
 * 	file where the element was called
 */
func (elem dedegoTraceWaitGroupElement) getFile() string {
	return elem.file
}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "W,'tpre','tpost','id','op','exec','delta','val','file':'line'"
 *    'tpre' (number): global before the operation
 *    'tpost' (number): global after the operation
 *    'id' (number): id of the mutex
 *	  'op' (A/W): A if it is an add or Done, W if it is a wait
 *	  'exec' (e/o): e if the operation was successfully finished, o otherwise
 *	  'delta' (number): delta of the waitgroup, positive for add, negative for done, 0 for wait
 *	  'val' (number): value of the waitgroup after the operation
 *    'file' (string): file where the operation was called
 *    'line' (number): line where the operation was called
 */
func (elem dedegoTraceWaitGroupElement) toString() string {
	res := "W,"
	res += uint64ToString(elem.timerPre) + "," + uint64ToString(elem.timerPost) + ","
	res += uint64ToString(elem.id) + ","
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
		val: val, file: file, line: line, timerPre: timer, timerPost: timer}
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
		line: line, timerPre: timer}
	return insertIntoTrace(elem)
}

// ============================= Channel =============================

type dedegoTraceChannelElement struct {
	id         uint64    // id of the channel
	exec       bool      // set true if the operation was successfully finished
	op         operation // operation
	qSize      uint32    // size of the channel, 0 for unbuffered
	qCountPre  uint32    // number of elements in the queue before the operation
	qCountPost uint32    // number of elements in the queue after the operation
	opId       uint64    // id of the operation
	file       string    // file where the operation was called
	line       int       // line where the operation was called
	timerPre   uint64    // global timer before the operation
	timerPost  uint64    // global timer after the operation
}

func (elem dedegoTraceChannelElement) isDedegoTraceElement() {}

/*
 * Get the file where the element was called
 * Return:
 * 	file where the element was called
 */
func (elem dedegoTraceChannelElement) getFile() string {
	return elem.file
}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "C,'tpre','tpost','id','op','exec','pId','file':'line'"
 *    'tpre' (number): global timer before the operation
 *    'tpost' (number): global timer after the operation
 *    'id' (number): id of the channel
 *	  'op' (S/R/C): S if it is a send, R if it is a receive, C if it is a close
 *	  'exec' (e/o): e if the operation was successfully finished, o otherwise
 *	  'pId' (number): id of the channel with witch the communication took place
 *    'file' (string): file where the operation was called
 *    'line' (number): line where the operation was called
 */
func (elem dedegoTraceChannelElement) toString() string {
	res := "C,"
	res += uint64ToString(elem.timerPre) + "," + uint64ToString(elem.timerPost) + ","
	res += uint64ToString(elem.id) + ","

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
	res += "," + uint32ToString(elem.qSize)
	res += "," + uint32ToString(elem.qCountPre)
	res += "," + uint32ToString(elem.qCountPost)
	res += "," + elem.file + ":" + intToString(elem.line)
	return res
}

/*
 * Add a channel send to the trace. If the channel send was created by an atomic
 * operation, add this to the trace as well
 * Args:
 * 	id: id of the channel
 * 	opId: id of the operation
 * 	qSize: size of the channel, 0 for unbuffered
 * 	qCount: number of elements in the queue before the operation
 * Return:
 * 	index of the operation in the trace
 */
var dedegoCounterAtomic uint64

func DedegoChanSendPre(id uint64, opId uint64, qSize uint, qCount uint) int {
	_, file, line, _ := Caller(3)
	if isSuffix(file, "dedego_atomic.go") {
		dedegoCounterAtomic += 1
		DedegoAtomic(dedegoCounterAtomic)
	}
	timer := GetDedegoCounter()
	elem := dedegoTraceChannelElement{id: id, op: opChanSend, opId: opId,
		file: file, line: line, timerPre: timer, qSize: uint32(qSize),
		qCountPre: uint32(qCount)}
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
 * 	qCount: number of elements in the channel before the receive operation
 * Return:
 * 	index of the operation in the trace
 */
func DedegoChanRecvPre(id uint64, opId uint64, qSize uint, qCount uint) int {
	_, file, line, _ := Caller(3)
	timer := GetDedegoCounter()
	elem := dedegoTraceChannelElement{id: id, op: opChanRecv, opId: opId,
		file: file, line: line, timerPre: timer, qSize: uint32(qSize),
		qCountPre: uint32(qCount)}
	return insertIntoTrace(elem)
}

/*
 * Add a channel close to the trace
 * Args:
 * 	id: id of the channel
 * Return:
 * 	index of the operation in the trace
 */
func DedegoChanClose(id uint64, qSize uint, qCount uint) int {
	_, file, line, _ := Caller(2)
	timer := GetDedegoCounter()
	elem := dedegoTraceChannelElement{id: id, op: opChanClose, file: file,
		exec: true,
		line: line, timerPre: timer, timerPost: timer, qSize: uint32(qSize),
		qCountPre: uint32(qCount), qCountPost: uint32(qCount)}
	return insertIntoTrace(elem)
}

/*
 * Set the size of the channel after the operation
 * Args:
 * 	index: index of the operation in the trace
 * 	qSize: size of the channel
 */
func DedegoChanPostQCount(index int, qCount uint) {
	if index == -1 {
		return
	}
	elem := currentGoRoutine().Trace[index].(dedegoTraceChannelElement)
	elem.qCountPost = uint32(qCount)
	currentGoRoutine().Trace[index] = elem
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
	elem := currentGoRoutine().Trace[index].(dedegoTraceChannelElement)
	elem.exec = true
	elem.timerPost = GetDedegoCounter()
	currentGoRoutine().Trace[index] = elem
}

// ============================= Select ==============================

type dedegoTraceSelectElement struct {
	id         uint64   // id of the select
	cases      []string // cases of the select
	send       []bool   // true if the case is a send, false if it is a receive
	nsend      int      // number of send cases
	chosen     int      // index of the chosen case
	chosenChan *hchan   // channel chosen
	exec       bool     // set true if the operation was successfully finished
	opId       uint64   // id of the operation
	defa       bool     // set true if a default case exists
	file       string   // file where the operation was called
	line       int      // line where the operation was called
	timerPre   uint64   // global timer before the operation
	timerPost  uint64   // global timer after the operation
}

func (elem dedegoTraceSelectElement) isDedegoTraceElement() {}

/*
 * Get the file where the element was called
 * Return:
 * 	file where the element was called
 */
func (elem dedegoTraceSelectElement) getFile() string {
	return elem.file
}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "S,'tpre','tpost','id','cases','exec','chosen','opId','file':'line'"
 *    'tpre' (number): global timer before the operation
 *    'tpost' (number): global timer after the operation
 *    'id' (number): id of the mutex
 *	  'cases' (string): cases of the select, id and r/s, separated by '.', d for default
 *	  'exec' (e/o): e if the operation was successfully finished, o otherwise
 *    'chosen' (number): index of the chosen case in cases (0 indexed, -1 for default)
 *	  'opId' (number): id of the operation on the channel
 *    'file' (string): file where the operation was called
 *    'line' (number): line where the operation was called
 */
func (elem dedegoTraceSelectElement) toString() string {
	res := "S," + uint64ToString(elem.id) + ","
	res += uint64ToString(elem.timerPre) + "," + uint64ToString(elem.timerPost) + ","

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

	if elem.exec {
		res += ",e"
	} else {
		res += ",o"
	}

	res += "," + intToString(elem.chosen) + "," + uint64ToString(elem.opId)
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
func DedegoSelectPre(cases *[]scase, nsends int, block bool,
	lockOrder []uint16) int {
	if cases == nil {
		return -1
	}
	id := GetDedegoObjectId()
	casesStr := make([]string, len(*cases))
	for i, ca := range *cases {
		if ca.c == nil {
			casesStr[i] = "-"
		} else {
			casesStr[i] = uint64ToString(ca.c.id)
		}
	}

	send := make([]bool, len(lockOrder))
	for i, lo := range lockOrder {
		send[i] = (lo < uint16(nsends))
	}

	_, file, line, _ := Caller(2)
	timer := GetDedegoCounter()
	elem := dedegoTraceSelectElement{id: id, cases: casesStr, nsend: nsends,
		send: send,
		defa: !block, file: file, line: line, timerPre: timer}
	return insertIntoTrace(elem)
}

/*
* Add a new select element to the trace if the select has exactly one non-default case
 */
func DedegoSelectPreOneNonDef(c *hchan, send bool) int {
	id := GetDedegoObjectId()

	casesStr := []string{uint64ToString(c.id)}

	nSend := 0
	pSend := []bool{send}
	if send {
		nSend = 1
	}

	_, file, line, _ := Caller(2)
	timer := GetDedegoCounter()
	elem := dedegoTraceSelectElement{id: id, cases: casesStr, nsend: nSend,
		send: pSend,
		defa: true, file: file, line: line, timerPre: timer}
	return insertIntoTrace(elem)
}

/*
 * Add the chosen case to the select
 * Args:
 * 	index: index of the operation in the trace
 * 	chosen: index of the chosen case
 * 	chosenChan: chosen channel
 */
func DedegoSelectPost1(index int, chosen int, chosenChan *hchan) {
	// internal elements are not in the trace
	if index == -1 {
		return
	}

	elem := currentGoRoutine().Trace[index].(dedegoTraceSelectElement)

	elem.exec = true
	elem.chosen = chosen
	elem.chosenChan = chosenChan
	elem.timerPost = GetDedegoCounter()

	lock(currentGoRoutine().lock)
	defer unlock(currentGoRoutine().lock)

	currentGoRoutine().Trace[index] = elem
}

/*
 * Add the lock order to the select
 * Args:
 * 	index: index of the operation in the trace
 * 	lockOrder: lock order of the select
 */
func DedegoSelectPost2(index int, lockOrder []uint16) {
	// internal elements are not in the trace
	if index == -1 {
		return
	}

	elem := currentGoRoutine().Trace[index].(dedegoTraceSelectElement)
	send := make([]bool, len(lockOrder))
	for i, lo := range lockOrder {
		send[i] = (lo < uint16(elem.nsend))
	}
	elem.send = send

	if elem.chosenChan != nil {
		if send[elem.chosen] {
			elem.chosenChan.numberSend++
			elem.opId = elem.chosenChan.numberSend
		} else {
			elem.chosenChan.numberRecv++
			elem.opId = elem.chosenChan.numberRecv
		}
	}

	lock(currentGoRoutine().lock)
	defer unlock(currentGoRoutine().lock)

	currentGoRoutine().Trace[index] = elem
}

/*
 * Add the selected case for a select with one non-default and one default case
 * Args:
 * 	index: index of the operation in the trace
 * 	res: 0 for the non-default case, -1 for the default case
 */
func DedegoSelectPostOneNonDef(index int, res bool) {
	if index == -1 {
		return
	}

	elem := currentGoRoutine().Trace[index].(dedegoTraceSelectElement)

	elem.exec = true
	elem.chosen = -1
	if res {
		elem.chosen = 0
	}
	elem.timerPost = GetDedegoCounter()

	lock(currentGoRoutine().lock)
	defer unlock(currentGoRoutine().lock)

	currentGoRoutine().Trace[index] = elem
}

// ============================= Atomic ================================
type dedegoTraceAtomicElement struct {
	timer     uint64 // global timer
	index     uint64 // index of the atomic event in dedegoAtomicMap
	operation int    // type of operation
}

func (elem dedegoTraceAtomicElement) isDedegoTraceElement() {}

/*
 * Get the file where the element was called
 * Return:
 * 	file where the element was called
 */
func (elem dedegoTraceAtomicElement) getFile() string {
	return ""
}

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

// DEDEGO-FILE-END
