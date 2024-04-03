// ADVOCATE-FILE-START

package runtime

import (
	at "runtime/internal/atomic"
)

type Operation int // enum for operation

const (
	OperationNone Operation = iota
	OperationSpawn
	OperationSpawned

	OperationChannelSend
	OperationChannelRecv
	OperationChannelClose

	OperationMutexLock
	OperationMutexUnlock
	OperationMutexTryLock
	OperationRWMutexLock
	OperationRWMutexUnlock
	OperationRWMutexTryLock
	OperationRWMutexRLock
	OperationRWMutexRUnlock
	OperationRWMutexTryRLock

	OperationOnce

	OperationWaitgroupAddDone
	OperationWaitgroupWait

	OperationSelect
	OperationSelectCase
	OperationSelectDefault

	OperationCondSignal
	OperationCondBroadcast
	OperationCondWait

	OperationAtomic

	OperationDisableReplay
)

type prePost int // enum for pre/post
const (
	pre prePost = iota
	post
	none
)

type advocateTraceElement interface {
	isAdvocateTraceElement()
	toString() string
	getOperation() Operation
	getFile() string
	getLine() int
}

type advocateAtomicMapElem struct {
	addr      uint64
	operation int
}

var advocateDisabled = true
var advocateAtomicMap = make(map[uint64]advocateAtomicMapElem)
var advocateAtomicMapRoutine = make(map[uint64]uint64)
var advocateAtomicMapToID = make(map[uint64]uint64)
var advocateAtomicMapIDCounter uint64 = 1
var advocateAtomicMapLock mutex
var advocateAtomicMapRoutineLock mutex
var advocateAtomicMapToIDLock mutex

var advocateTraceWritingDisabled = false

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
func traceToString(trace *[]advocateTraceElement) string {
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
func insertIntoTrace(elem advocateTraceElement) int {
	return currentGoRoutine().addToTrace(elem)
}

/*
 * Print the trace of the current routines
 */
func PrintTrace() {
	routineID := GetRoutineID()
	println("Routine", routineID, ":", CurrentTraceToString())
}

/*
 * Return the trace of the routine with id 'id'
 * Args:
 * 	id: id of the routine
 * Return:
 * 	string representation of the trace of the routine
 * 	bool: true if the routine exists, false otherwise
 */
func TraceToStringByID(id uint64) (string, bool) {
	println("TraceToStringById", id)
	lock(&AdvocateRoutinesLock)
	defer unlock(&AdvocateRoutinesLock)
	if routine, ok := AdvocateRoutines[id]; ok {
		return traceToString(&routine.Trace), true
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
func TraceToStringByIDChannel(id int, c chan<- string) {
	lock(&AdvocateRoutinesLock)

	if routine, ok := AdvocateRoutines[uint64(id)]; ok {
		unlock(&AdvocateRoutinesLock)
		res := ""
		for i, elem := range routine.Trace {
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
	} else {
		unlock(&AdvocateRoutinesLock)
	}
}

/*
 * Return the trace of all traces
 * Return:
 * 	string representation of the trace of all routines
 */
func AllTracesToString() string {
	// write warning if projectPath is empty
	res := ""
	lock(&AdvocateRoutinesLock)
	defer unlock(&AdvocateRoutinesLock)

	for i := 1; i <= len(AdvocateRoutines); i++ {
		res += ""
		routine := AdvocateRoutines[uint64(i)]
		if routine == nil {
			panic("Trace is nil")
		}
		res += traceToString(&routine.Trace) + "\n"

	}
	return res
}

/*
* PrintAllTraces prints the trace of all routines
 */
func PrintAllTraces() {
	print(AllTracesToString())
}

/*
 * GetNumberOfRoutines returns the number of routines in the trace
 * Return:
 *	number of routines in the trace
 */
func GetNumberOfRoutines() int {
	lock(&AdvocateRoutinesLock)
	defer unlock(&AdvocateRoutinesLock)
	return len(AdvocateRoutines)
}

/*
 * InitAdvocate enables the collection of the trace
 * Args:
 * 	size: size of the channel used to link the atomic recording to the main
 *    recording.
 */
func InitAdvocate(size int) {
	if size < 0 {
		size = 0
	}
	chanSize := (size + 1) * 10000000
	// link runtime with atomic via channel to receive information about
	// atomic events
	c := make(chan at.AtomicElem, chanSize)
	at.AdvocateAtomicLink(c)

	go func() {
		for atomic := range c {
			lock(&advocateAtomicMapLock)
			advocateAtomicMap[atomic.Index] = advocateAtomicMapElem{
				addr:      atomic.Addr,
				operation: atomic.Operation,
			}
			unlock(&advocateAtomicMapLock)
			// go func() {
			// 	WaitForReplayAtomic(atomic.Operation, atomic.Index)
			// 	atomic.ChanReturn <- true
			// }()
		}
	}()

	advocateDisabled = false
}

/*
 * DisableTrace disables the collection of the trace
 */
func DisableTrace() {
	at.AdvocateAtomicUnlink()
	advocateDisabled = true
}

/*
 * BockTrace blocks the trace collection
 * Resume using UnblockTrace
 */
func BlockTrace() {
	advocateTraceWritingDisabled = true
}

/*
 * UnblockTrace resumes the trace collection
 * Block using BlockTrace
 */
func UnblockTrace() {
	advocateTraceWritingDisabled = false
}

/*
 * DeleteTrace removes all trace elements from the trace
 * Do not remove the routine objects them self
 * Make sure to call BlockTrace(), before calling this function
 */
func DeleteTrace() {
	lock(&AdvocateRoutinesLock)
	defer unlock(&AdvocateRoutinesLock)
	for i := range AdvocateRoutines {
		AdvocateRoutines[i].Trace = AdvocateRoutines[i].Trace[:0]
	}
}

// ============================= Routine ===========================

// type to save in the trace for routines
type advocateSpawnElement struct {
	id    uint64 // id of the routine
	timer uint64 // global timer
	file  string // file where the routine was created
	line  int    // line where the routine was created
}

func (elem advocateSpawnElement) isAdvocateTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "G,'id'"
 *    'id' (number): id of the routine
 */
func (elem advocateSpawnElement) toString() string {
	return "G," + uint64ToString(elem.timer) + "," + uint64ToString(elem.id) + "," + elem.file + ":" + intToString(elem.line)
}

/*
 * Get the operation
 */
func (elem advocateSpawnElement) getOperation() Operation {
	return OperationSpawn
}

/*
 * Get the file
 */
func (elem advocateSpawnElement) getFile() string {
	return elem.file
}

/*
 * Get the line
 */
func (elem advocateSpawnElement) getLine() int {
	return elem.line
}

/*
 * AdvocateSpawnCaller adds a routine spawn to the trace
 * Args:
 * 	callerRoutine: routine that created the new routine
 * 	newID: id of the new routine
 * 	file: file where the routine was created
 * 	line: line where the routine was created
 */
func AdvocateSpawnCaller(callerRoutine *AdvocateRoutine, newID uint64, file string, line int32) {
	timer := GetAdvocateCounter()
	callerRoutine.addToTrace(advocateSpawnElement{id: newID, timer: timer,
		file: file, line: int(line)})
}

// type to save in the trace for routines
type advocateTraceSpawnedElement struct {
	id    uint64 // id of the routine
	timer uint64 // global timer
	file  string // file where the routine was created
	line  int    // line where the routine was created
}

func (elem advocateTraceSpawnedElement) isAdvocateTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "G,'id'"
 *    'id' (number): id of the routine
 */
func (elem advocateTraceSpawnedElement) toString() string {
	return "g," + uint64ToString(elem.timer) + "," + uint64ToString(elem.id) + "," + elem.file + ":" + intToString(elem.line)
}

// ============================= Mutex =============================

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

// ============================= WaitGroup ===========================

type advocateWaitGroupElement struct {
	id    uint64    // id of the waitgroup
	op    Operation // operation
	delta int       // delta of the waitgroup
	val   int32     // value of the waitgroup after the operation
	file  string    // file where the operation was called
	line  int       // line where the operation was called
	tPre  uint64    // global timer
	tPost uint64    // global timer
}

func (elem advocateWaitGroupElement) isAdvocateTraceElement() {}

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
func (elem advocateWaitGroupElement) toString() string {
	res := "W,"
	res += uint64ToString(elem.tPre) + "," + uint64ToString(elem.tPost) + ","
	res += uint64ToString(elem.id) + ","
	switch elem.op {
	case OperationWaitgroupAddDone:
		res += "A,"
	case OperationWaitgroupWait:
		res += "W,"
	}

	res += intToString(elem.delta) + "," + int32ToString(elem.val)
	res += "," + elem.file + ":" + intToString(elem.line)
	return res
}

/*
 * Get the operation
 */
func (elem advocateWaitGroupElement) getOperation() Operation {
	return elem.op
}

/*
 * Get the file
 */
func (elem advocateWaitGroupElement) getFile() string {
	return elem.file
}

/*
 * Get the line
 */
func (elem advocateWaitGroupElement) getLine() int {
	return elem.line
}

/*
 * AdvocateWaitGroupAdd adds a waitgroup add or done to the trace
 * Args:
 * 	id: id of the waitgroup
 *  delta: delta of the waitgroup
 * 	val: value of the waitgroup after the operation
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateWaitGroupAdd(id uint64, delta int, val int32) int {
	var file string
	var line int
	if delta > 0 {
		_, file, line, _ = Caller(2)
	} else {
		_, file, line, _ = Caller(3)
	}
	timer := GetAdvocateCounter()
	elem := advocateWaitGroupElement{id: id, op: OperationWaitgroupAddDone,
		delta: delta, val: val, file: file, line: line, tPre: timer, tPost: timer}
	return insertIntoTrace(elem)

}

/*
 * AdvocateWaitGroupWaitPre adds a waitgroup wait to the trace
 * Args:
 * 	id: id of the waitgroup
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateWaitGroupWaitPre(id uint64) int {
	_, file, line, _ := Caller(2)
	timer := GetAdvocateCounter()
	elem := advocateWaitGroupElement{id: id, op: OperationWaitgroupWait,
		file: file, line: line, tPre: timer}
	return insertIntoTrace(elem)
}

// ============================= Channel =============================

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

// ============================= Select ==============================

type advocateSelectElement struct {
	tPre    uint64                   // global timer before the operation
	tPost   uint64                   // global timer after the operation
	id      uint64                   // id of the select
	cases   []advocateChannelElement // cases of the select
	chosen  int                      // index of the chosen case in cases (0 indexed, -1 for default)
	nsend   int                      // number of send cases
	defa    bool                     // set true if a default case exists
	defaSel bool                     // set true if a default case was chosen
	file    string                   // file where the operation was called
	line    int                      // line where the operation was called
}

func (elem advocateSelectElement) isAdvocateTraceElement() {}

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
func (elem advocateSelectElement) toString() string {
	res := "S,"
	res += uint64ToString(elem.tPre) + "," + uint64ToString(elem.tPost) + ","
	res += uint64ToString(elem.id) + ","

	notNil := 0
	for _, ca := range elem.cases { // cases
		if ca.tPre != 0 { // ignore nil cases
			if notNil != 0 {
				res += "~"
			}
			res += ca.toStringSep(".", false)
			notNil++
		}
	}

	if elem.defa { // default
		if notNil != 0 {
			res += "~"
		}
		if elem.defaSel {
			res += "D"
		} else {
			res += "d"
		}
	}

	res += "," + intToString(elem.chosen) // case index

	res += "," + elem.file + ":" + intToString(elem.line)
	return res
}

/*
 * Get the operation
 */
func (elem advocateSelectElement) getOperation() Operation {
	return OperationSelect
}

/*
 * Get the file
 */
func (elem advocateSelectElement) getFile() string {
	return elem.file
}

/*
 * Get the line
 */
func (elem advocateSelectElement) getLine() int {
	return elem.line
}

/*
 * AdvocateSelectPre adds a select to the trace
 * Args:
 * 	cases: cases of the select
 * 	nsends: number of send cases
 * 	block: true if the select is blocking (has no default), false otherwise
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateSelectPre(cases *[]scase, nsends int, block bool) int {
	timer := GetAdvocateCounter()
	if cases == nil {
		return -1
	}

	id := GetAdvocateObjectID()
	caseElements := make([]advocateChannelElement, len(*cases))
	_, file, line, _ := Caller(2)

	for i, ca := range *cases {
		if ca.c != nil { // ignore nil cases
			caseElements[i] = advocateChannelElement{id: ca.c.id,
				op:    OperationChannelRecv,
				qSize: uint32(ca.c.dataqsiz), tPre: timer}
		}
	}

	elem := advocateSelectElement{id: id, cases: caseElements, nsend: nsends,
		defa: !block, file: file, line: line, tPre: timer}
	return insertIntoTrace(elem)
}

/*
 * AdvocateSelectPost adds a post event for select in case of an non-default case
 * Args:
 * 	index: index of the operation in the trace
 * 	c: channel of the chosen case
 * 	chosenIndex: index of the chosen case in the select
 * 	lockOrder: order of the locks
 * 	rClosed: true if the channel was closed at another routine
 */
func AdvocateSelectPost(index int, c *hchan, chosenIndex int,
	lockOrder []uint16, rClosed bool) {

	if index == -1 || c == nil {
		return
	}

	elem := currentGoRoutine().getElement(index).(advocateSelectElement)
	timer := GetAdvocateCounter()

	elem.chosen = chosenIndex
	elem.tPost = timer

	for i, op := range lockOrder {
		opChan := OperationChannelRecv
		if op < uint16(elem.nsend) {
			opChan = OperationChannelSend
		}
		elem.cases[i].op = opChan
	}

	if chosenIndex == -1 { // default case
		elem.defaSel = true
	} else {
		elem.cases[chosenIndex].tPost = timer
		elem.cases[chosenIndex].closed = rClosed

		// set oId
		if elem.cases[chosenIndex].op == OperationChannelSend {
			c.numberSend++
			elem.cases[chosenIndex].opID = c.numberSend
		} else {
			c.numberRecv++
			elem.cases[chosenIndex].opID = c.numberRecv
		}

	}

	currentGoRoutine().updateElement(index, elem)
}

/*
* AdvocateSelectPreOneNonDef adds a new select element to the trace if the
* select has exactly one non-default case and a default case
* Args:
* 	c: channel of the non-default case
* 	send: true if the non-default case is a send, false otherwise
* Return:
* 	index of the operation in the trace
 */
func AdvocateSelectPreOneNonDef(c *hchan, send bool) int {
	if c == nil {
		return -1
	}

	id := GetAdvocateObjectID()
	timer := GetAdvocateCounter()

	opChan := OperationChannelRecv
	if send {
		opChan = OperationChannelSend
	}

	caseElements := make([]advocateChannelElement, 1)
	caseElements[0] = advocateChannelElement{id: c.id,
		qSize: uint32(c.dataqsiz), tPre: timer, op: opChan}

	nSend := 0
	if send {
		nSend = 1
	}

	_, file, line, _ := Caller(2)

	elem := advocateSelectElement{id: id, cases: caseElements, nsend: nSend,
		defa: true, file: file, line: line, tPre: timer}
	return insertIntoTrace(elem)
}

/*
 * AdvocateSelectPostOneNonDef adds the selected case for a select with one
 * non-default and one default case
 * Args:
 * 	index: index of the operation in the trace
 * 	res: 0 for the non-default case, -1 for the default case
 */
func AdvocateSelectPostOneNonDef(index int, res bool, c *hchan) {
	if index == -1 {
		return
	}

	timer := GetAdvocateCounter()
	elem := currentGoRoutine().getElement(index).(advocateSelectElement)

	if res {
		elem.chosen = 0
		elem.cases[0].tPost = timer
		if elem.cases[0].op == OperationChannelSend {
			c.numberSend++
			elem.cases[0].opID = c.numberSend
		} else {
			c.numberRecv++
			elem.cases[0].opID = c.numberRecv
		}
	} else {
		elem.chosen = -1
		elem.defaSel = true
	}
	elem.tPost = timer

	currentGoRoutine().updateElement(index, elem)
}

// ============================= Atomic ================================
type advocateAtomicElement struct {
	timer     uint64 // global timer
	index     uint64 // index of the atomic event in advocateAtomicMap
	operation int    // type of operation
}

func (elem advocateAtomicElement) isAdvocateTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "A,'addr'"
 *    'addr' (number): address of the atomic variable
 */
// enum for atomic operation, must be the same as in advocate_atomic.go
const (
	LoadOp = iota
	StoreOp
	AddOp
	SwapOp
	CompSwapOp
)

func (elem advocateAtomicElement) toString() string {
	lock(&advocateAtomicMapLock)
	mapElement := advocateAtomicMap[elem.index]
	unlock(&advocateAtomicMapLock)
	lock(&advocateAtomicMapToIDLock)
	if _, ok := advocateAtomicMapToID[mapElement.addr]; !ok {
		advocateAtomicMapToID[mapElement.addr] = advocateAtomicMapIDCounter
		advocateAtomicMapIDCounter++
	}
	id := advocateAtomicMapToID[mapElement.addr]
	unlock(&advocateAtomicMapToIDLock)

	res := "A," + uint64ToString(elem.timer) + "," +
		uint64ToString(id) + ","
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
	return res
}

/*
 * Get the operation
 */
func (elem advocateAtomicElement) getOperation() Operation {
	return OperationAtomic
}

/*
 * Get the file
 */
func (elem advocateAtomicElement) getFile() string {
	return ""
}

/*
 * Get the line
 */
func (elem advocateAtomicElement) getLine() int {
	return 0
}

/*
 * Add an atomic operation to the trace
 * Args:
 * 	index: index of the atomic event in advocateAtomicMap
 */
func AdvocateAtomic(index uint64) {
	timer := GetAdvocateCounter()
	elem := advocateAtomicElement{index: index, timer: timer}
	insertIntoTrace(elem)
}

// ======================= Once ============================
type advocateOnceElement struct {
	tpre  uint64 // global timer at the beginning of the execution
	tpost uint64 // global timer at the end of the execution
	id    uint64 // id of the once
	suc   bool   // true if the do on the once was called for the first time
	file  string // file where the operation was called
	line  int    // line where the operation was called
}

func (elem advocateOnceElement) isAdvocateTraceElement() {}

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
func (elem advocateOnceElement) toString() string {
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

/*
 * Get the operation
 */
func (elem advocateOnceElement) getOperation() Operation {
	return OperationOnce
}

/*
 * Get the file
 */
func (elem advocateOnceElement) getFile() string {
	return elem.file
}

/*
 * Get the line
 */
func (elem advocateOnceElement) getLine() int {
	return elem.line
}

/*
 * AdvocateOncePre adds a once to the trace
 * Args:
 * 	id: id of the once
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateOncePre(id uint64) int {
	_, file, line, _ := Caller(2)
	timer := GetAdvocateCounter()
	elem := advocateOnceElement{id: id, tpre: timer, file: file, line: line}
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
	timer := GetAdvocateCounter()
	elem := currentGoRoutine().getElement(index).(advocateOnceElement)
	elem.tpost = timer
	elem.suc = suc

	currentGoRoutine().updateElement(index, elem)
}

// ======================= Cond ============================
type advocateCondElement struct {
	tpre  uint64 // global timer at the beginning of the execution
	tpost uint64 // global timer at the end of the execution
	id    uint64 // id of the cond
	op    Operation
	file  string // file where the operation was called
	line  int    // line where the operation was called
}

func (elem advocateCondElement) isAdvocateTraceElement() {}

/*
 * Return a string representation of the element
 * Return:
 * 	string representation of the element "C,'tpre','tpost','id','op','file':'line"
 *    'tpre' (number): global timer at the beginning of the execution
 *    'tpost' (number): global timer at the end of the execution
 *    'id' (number): id of the cond
 *    'op' (W/S/B): W if it is a wait, S if it is a signal, B if it is a broadcast
 *    'file' (string): file where the operation was called
 *    'line' (string): line where the operation was called
 */
func (elem advocateCondElement) toString() string {
	res := "N,"
	res += uint64ToString(elem.tpre) + ","
	res += uint64ToString(elem.tpost) + ","
	res += uint64ToString(elem.id) + ","
	switch elem.op {
	case OperationCondWait:
		res += "W"
	case OperationCondSignal:
		res += "S"
	case OperationCondBroadcast:
		res += "B"
	}
	res += "," + elem.file + ":" + intToString(elem.line)
	return res
}

/*
 * Get the operation
 */
func (elem advocateCondElement) getOperation() Operation {
	return elem.op
}

/*
 * Get the file
 */
func (elem advocateCondElement) getFile() string {
	return elem.file
}

/*
 * Get the line
 */
func (elem advocateCondElement) getLine() int {
	return elem.line
}

/*
 * AdvocateCondPre adds a cond wait to the trace
 * Args:
 * 	id: id of the cond
 * 	op: 0 for wait, 1 for signal, 2 for broadcast
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateCondPre(id uint64, op int) int {
	_, file, line, _ := Caller(2)
	timer := GetAdvocateCounter()
	var opC Operation
	switch op {
	case 0:
		opC = OperationCondWait
	case 1:
		opC = OperationCondSignal
	case 2:
		opC = OperationCondWait
	default:
		panic("Unknown cond operation")
	}
	elem := advocateCondElement{id: id, tpre: timer, file: file, line: line, op: opC}
	return insertIntoTrace(elem)
}

/*
 * AdvocateCondPost adds the end counter to an operation of the trace
 * Args:
 * 	index: index of the operation in the trace
 */
func AdvocateCondPost(index int) {
	if index == -1 {
		return
	}
	timer := GetAdvocateCounter()
	elem := currentGoRoutine().getElement(index).(advocateCondElement)
	elem.tpost = timer

	currentGoRoutine().updateElement(index, elem)
}

// ====================== Ignore =========================

/*
 * Some operations, like garbage collection and internal operations, can
 * cause the replay to get stuck or are not needed.
 * For this reason, we ignore them.
 * Arguments:
 * 	operation: operation that is about to be executed
 * 	file: file in which the operation is executed
 * 	line: line number of the operation
 * Return:
 * 	bool: true if the operation should be ignored, false otherwise
 */
// TODO: check if all of them are necessary
func AdvocateIgnore(operation Operation, file string, line int) bool {
	if hasSuffix(file, "advocate/advocate.go") ||
		hasSuffix(file, "advocate/advocate_replay.go") ||
		hasSuffix(file, "advocate/advocate_routine.go") ||
		hasSuffix(file, "advocate/advocate_trace.go") ||
		hasSuffix(file, "advocate/advocate_utile.go") ||
		hasSuffix(file, "advocate/advocate_atomic.go") { // internal
		return true
	}

	if hasSuffix(file, "syscall/env_unix.go") {
		return true
	}

	switch operation {
	case OperationSpawn:
		// garbage collection can cause the replay to get stuck
		if hasSuffix(file, "runtime/mgc.go") && line == 1215 {
			return true
		}
	case OperationMutexLock, OperationMutexUnlock:
		// mutex operations in the once can cause the replay to get stuck,
		// if the once was called by the poll/fd_poll_runtime.go init.
		if hasSuffix(file, "sync/once.go") && (line == 115 || line == 116 ||
			line == 121 || line == 125) {
			return true
		}
		// pools
		if hasSuffix(file, "sync/pool.go") && (line == 217 || line == 218 ||
			line == 224 || line == 234) {
			return true
		}
		// mutex in rwmutex
		// if hasSuffix(file, "sync/rwmutex.go") && (line == 270 || line == 396) {
		// 	return true
		// }
	case OperationOnce:
		// once operations in the poll/fd_poll_runtime.go init can cause the replay to get stuck.
		if hasSuffix(file, "internal/poll/fd_poll_runtime.go") && line == 40 {
			return true
		}
	}
	return false
}

// ADVOCATE-FILE-END
