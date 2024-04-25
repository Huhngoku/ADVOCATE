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

	OperationReplayStart
	OperationReplayEnd
)

type prePost int // enum for pre/post
const (
	pre prePost = iota
	post
	none
)

// type advocateTraceElement interface {
// 	isAdvocateTraceElement()
// 	toString() string
// 	getOperation() Operation
// 	getFile() string
// 	getLine() int
// }

type advocateAtomicMapElem struct {
	addr      uint64
	operation int
}

var advocateDisabled = true
var advocateAtomicMap = make(map[uint64]advocateAtomicMapElem)
var advocateAtomicMapToID = make(map[uint64]uint64)
var advocateAtomicMapIDCounter uint64 = 1
var advocateAtomicMapLock mutex
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
		res += elem
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
func traceToString(trace *[]string) string {
	res := ""
	for i, elem := range *trace {
		if i != 0 {
			res += ";"
		}
		res += elem
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
func insertIntoTrace(elem string) int {
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

			if elem[0] == 'A' {
				elem = addAtomicInfo(elem)
			}

			res += elem

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
	disableAtomics := false
	if size < 0 {
		size = 0
		disableAtomics = true
	}
	chanSize := (size + 1) * 10000000
	// link runtime with atomic via channel to receive information about
	// atomic events
	c := make(chan at.AtomicElem, chanSize)

	if !disableAtomics {
		at.AdvocateAtomicLink(c)
	}

	go func() {
		for atomic := range c {
			AdvocateAtomicPost(atomic)

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
