package runtime

type ReplayOperation int

const (
	AdvocateNone ReplayOperation = iota
	AdvocateReplaySpawn
	AdvocateReplaySpawned

	AdvocateReplayChannelSend
	AdvocateReplayChannelRecv
	AdvocateReplayChannelClose

	AdvocateReplayMutexLock
	AdvocateReplayMutexUnlock
	AdvocateReplayMutexTryLock
	AdvocateReplayRWMutexLock
	AdvocateReplayRWMutexUnlock
	AdvocateReplayRWMutexTryLock
	AdvocateReplayRWMutexRLock
	AdvocateReplayRWMutexRUnlock
	AdvocateReplayRWMutexTryRLock

	AdvocateReplayOnce

	AdvocateReplayWaitgroupAddDone
	AdvocateReplayWaitgroupWait

	AdvocateReplaySelect
	AdvocateReplaySelectCase
	AdvocateReplaySelectDefault

	// AdvocateReplayAtomicLoad
	// AdvocateReplayAtomicStore
	// AdvocateReplayAtomicAdd
	// AdvocateReplayAtomicSwap
	// AdvocateReplayAtomicCompareAndSwap
)

/*
 * The replay data structure.
 * The replay data structure is used to store the trace of the program.
 * op: identifier of the operation
 * time: time (tpre) of the operation
 * file: file in which the operation is executed
 * line: line number of the operation
 * blocked: true if the operation is blocked (never finised, tpost=0), false otherwise
 * suc: success of the opeartion
 *     - for mutexes: trylock operations true if the lock was acquired, false otherwise
 * 			for other operations always true
 *     - for once: true if the once was chosen (was the first), false otherwise
 *     - for others: always true
 * PFile: file of the partner (mainly for channel/select)
 * PLine: line of the partner (mainly for channel/select)
 * SelIndex: index of the select case (only for select, otherwise)
 */
type ReplayElement struct {
	Routine  int
	Op       ReplayOperation
	Time     int
	File     string
	Line     int
	Blocked  bool
	Suc      bool
	PFile    string
	PLine    int
	SelIndex int
}

type AdvocateReplayTrace []ReplayElement

var replayEnabled bool = false
var replayLock mutex
var replayIndex int = 0
var replayDone int = 0
var replayDoneLock mutex

var replayData AdvocateReplayTrace = make(AdvocateReplayTrace, 0)

func (t AdvocateReplayTrace) Print() {
	for _, e := range t {
		println(e.Op, e.Time, e.File, e.Line, e.Blocked, e.Suc)
	}
}

func EnableReplay(trace AdvocateReplayTrace) {
	replayData = trace
	replayEnabled = true
	// trace.Print()
}

func WaitForReplayFinish() {
	for {
		lock(&replayDoneLock)
		if replayDone >= len(replayData) {
			unlock(&replayDoneLock)
			break
		}
		unlock(&replayDoneLock)
	}
}

/*
 * Wait until the correct operation is about to be executed.
 * Arguments:
 * 	op: the operation type that is about to be executed
 * 	skip: number of stack frames to skip
 * Return:
 * 	bool: true if trace replay is enabled, false otherwise
 * 	chan ReplayElement: channel to receive the next replay element
 */
func WaitForReplay(op ReplayOperation, skip int) (bool, ReplayElement) {
	if !replayEnabled {
		return false, ReplayElement{}
	}

	_, file, line, _ := Caller(skip)

	return WaitForReplayPath(op, file, line)
}

/*
 * Wait until the correct atomic operation is about to be executed.
 * Args:
 * 	op: the operation type that is about to be executed
 * 	index: index of the atomic operation
 * Return:
 * 	bool: true if trace replay is enabled, false otherwise
 * 	ReplayElement: the next replay element
 */
// func WaitForReplayAtomic(op int, index uint64) (bool, ReplayElement) {
// 	lock(&advocateAtomicMapLock)
// 	routine := advocateAtomicMapRoutine[index]
// 	unlock(&advocateAtomicMapLock)

// 	if !replayEnabled {
// 		return false, ReplayElement{}
// 	}

// 	for {
// 		next := getNextReplayElement()
// 		// print("Replay: ", next.Time, " ", next.Op, " ", op, " ", next.File, " ", file, " ", next.Line, " ", line, "\n")

// 		if next.Time != 0 {
// 			if int(next.Op) != op || uint64(next.Routine) != routine {
// 				continue
// 			}
// 		}

// 		lock(&replayLock)
// 		replayIndex++
// 		unlock(&replayLock)
// 		// println("Replay: ", next.Time, op, file, line)
// 		return true, next
// 	}
// }

/*
 * Wait until the correct operation is about to be executed.
 * Arguments:
 * 		op: the operation type that is about to be executed
 * 		file: file in which the operation is executed
 * 		line: line number of the operation
 * Return:
 * 	bool: true if trace replay is enabled, false otherwise
 * 	chan ReplayElement: channel to receive the next replay element
 */
func WaitForReplayPath(op ReplayOperation, file string, line int) (bool, ReplayElement) {
	if !replayEnabled {
		return false, ReplayElement{}
	}

	if IgnoreInReplay(op, file, line) {
		return true, ReplayElement{}
	}

	// println("WaitForReplayPath", op, file, line)
	for {
		next := getNextReplayElement()
		// print("Replay: ", next.Time, " ", next.Op, " ", op, " ", next.File, " ", file, " ", next.Line, " ", line, "\n")

		if next.Time != 0 {
			if (next.Op != op && !correctSelect(next.Op, op)) ||
				next.File != file || next.Line != line {
				// TODO: sleep here to not waste CPU
				continue
			}
		}

		lock(&replayLock)
		replayIndex++
		unlock(&replayLock)
		// println("Replay: ", next.Time, op, file, line)
		return true, next
	}
}

/*
 * Notify that the operation is done.
 * This function should be called after a waiting operation is done.
 * Used to prevent the program to terminate before the trace is finished, if
 * the main routine would terminate.
 */
func ReplayDone() {
	if !replayEnabled {
		return
	}
	lock(&replayDoneLock)
	defer unlock(&replayDoneLock)
	replayDone++
}

func correctSelect(next ReplayOperation, op ReplayOperation) bool {
	if op != AdvocateReplaySelect {
		return false
	}

	if next != AdvocateReplaySelectCase && next != AdvocateReplaySelectDefault {
		return false
	}

	return true
}

func BlockForever() {
	gopark(nil, nil, waitReasonZero, traceEvNone, 1)
}

/*
 * Get the next replay element.
 * The function returns the next replay element and increments the index.
 */
func getNextReplayElement() ReplayElement {
	lock(&replayLock)
	defer unlock(&replayLock)
	if replayIndex >= len(replayData) {
		return ReplayElement{}
		// panic("Tace to short. The Program was most likely altered between recording and replay.")
	}
	return replayData[replayIndex]
}

/*
 * Some operations, like garbage collection, can cause the replay to get stuck.
 * For this reason, we ignore them.
 * Arguments:
 * 	operation: operation that is about to be executed
 * 	file: file in which the operation is executed
 * 	line: line number of the operation
 * Return:
 * 	bool: true if the operation should be ignored, false otherwise
 */
// TODO: check if all of them are necessary
func IgnoreInReplay(operation ReplayOperation, file string, line int) bool {
	if hasSuffix(file, "syscall/env_unix.go") {
		return true
	}
	switch operation {
	case AdvocateReplaySpawn:
		// garbage collection can cause the replay to get stuck
		if hasSuffix(file, "runtime/mgc.go") && line == 1215 {
			return true
		}
	case AdvocateReplayMutexLock, AdvocateReplayMutexUnlock:
		// mutex operations in the once can cause the replay to get stuck,
		// if the once was called by the poll/fd_poll_runtime.go init.
		if hasSuffix(file, "sync/once.go") && (line == 115 || line == 121 || line == 125) {
			return true
		}
		// pools
		if hasSuffix(file, "sync/pool.go") && (line == 216 || line == 223 || line == 233) {
			return true
		}
		// mutex in rwmutex
		if hasSuffix(file, "sync/rwmutex.go") && (line == 270 || line == 396) {
			return true
		}
	case AdvocateReplayOnce:
		// once operations in the poll/fd_poll_runtime.go init can cause the replay to get stuck.
		if hasSuffix(file, "internal/poll/fd_poll_runtime.go") && line == 39 {
			return true
		}
	}
	return false
}
