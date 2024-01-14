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

var replayEnabled bool
var replayLock mutex
var replayIndex int
var replayDone int
var replayDoneLock mutex

var replayData = make(AdvocateReplayTrace, 0)
var traceElementPositions = make(map[string][]int) // file -> []line

var timeoutMessageCycle = 500 // approx. 10s

func (t AdvocateReplayTrace) Print() {
	for _, e := range t {
		println(e.Op, e.Time, e.File, e.Line, e.Blocked, e.Suc)
	}
}

func EnableReplay(trace AdvocateReplayTrace) {
	replayData = trace

	// extract and save positions
	for _, e := range replayData {
		if _, ok := traceElementPositions[e.File]; !ok {
			traceElementPositions[e.File] = make([]int, 1)
		}

		if !containsInt(traceElementPositions[e.File], e.Line) {
			traceElementPositions[e.File] = append(traceElementPositions[e.File], e.Line)
		}
	}

	replayEnabled = true
	// trace.Print()
}

/*
 * Wait until all operations in the trace are executed.
 * This function should be called after the main routine is finished, to prevent
 * the program to terminate before the trace is finished.
 */
func WaitForReplayFinish() {
	timeoutCounter := 0
	for {
		timeoutCounter++
		lock(&replayDoneLock)
		if replayDone >= len(replayData) {
			unlock(&replayDoneLock)
			break
		}
		unlock(&replayDoneLock)

		// check for timeout
		if timeoutCounter%timeoutMessageCycle == 0 {
			waitTime := intToString(int(10 * timeoutCounter / timeoutMessageCycle))
			warningMessage := "\nReplayWarning: Long wait time for finishing replay."
			warningMessage += "The main routine has already finished approx. "
			warningMessage += waitTime
			warningMessage += "s ago, but the trace still contains not executed operations.\n"
			warningMessage += "This can be caused by a stuck replay.\n"
			warningMessage += "Possible causes are:\n"
			warningMessage += "    - The program was altered between recording and replay\n"
			warningMessage += "    - The program execution path is not deterministic, e.g. its execution path is determined by a random number\n"
			warningMessage += "    - The program execution path depends on the order of not tracked operations\n"
			warningMessage += "    - The program execution depends on outside input, that was not exactly reproduced\n"
			warningMessage += "If you believe, the program is still running, you can continue to wait.\n"
			warningMessage += "If you believe, the program is stuck, you can cancel the program.\n"
			warningMessage += "If you suspect, that one of these causes is the reason for the long wait time, you can try to change the program to avoid the problem.\n"
			warningMessage += "If the problem persist, this message will be repeated every approx. 10s.\n\n"
		}

		slowExecution()
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
	timeoutCounter := 0
	for {
		next := getNextReplayElement()
		// print("Replay: ", next.Time, " ", next.Op, " ", op, " ", next.File, " ", file, " ", next.Line, " ", line, "\n")

		if next.Time != 0 {
			if (next.Op != op && !correctSelect(next.Op, op)) ||
				next.File != file || next.Line != line {
				timeoutCounter++

				checkForTimeout(timeoutCounter, file, line)
				slowExecution()
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
 * At specified timeoutCounter values, do some checks or messages.
 * At the first timeout, check if the position is in the trace.
 * At following timeouts, print a warning message.
 * At the last timeout, panic.
 * Args:
 * 	timeoutCounter: the current timeout counter
 * 	file: file in which the operation is executed
 * 	line: line number of the operation
 */
func checkForTimeout(timeoutCounter int, file string, line int) bool {
	messageCauses := "Possible causes are:\n"
	messageCauses += "    - The program was altered between recording and replay\n"
	messageCauses += "    - The program execution path is not deterministic, e.g. its execution path is determined by a random number\n"
	messageCauses += "    - The program execution path depends on the order of not tracked operations\n"
	messageCauses += "    - The program execution depends on outside input, that was not exactly reproduced\n"

	if timeoutCounter == 250 { // ca. 5s
		res := isPositionInTrace(file, line)
		if !res {
			errorMessage := "ReplayError: Program tried to execute an operation that is not in the trace:\n"
			errorMessage += "    File: " + file + "\n"
			errorMessage += "    Line: " + intToString(line) + "\n"
			errorMessage += "This means, that the program replay was not successful.\n"
			errorMessage += messageCauses
			errorMessage += "If you suspect, that one of these causes is the reason for the error, you can try to change the program to avoid the problem.\n"
			errorMessage += "If this is not possible, you can try to rerun the replay, hoping the error does not occur again.\n"
			errorMessage += "If this is not possible or does not work, the program replay is currently not possible.\n\n"

			panic(errorMessage)
		}
	} else if timeoutCounter%timeoutMessageCycle == 0 { // approx. every 10s
		waitTime := intToString(int(10 * timeoutCounter / timeoutMessageCycle))
		warningMessage := "\nReplayWarning: Long wait time of approx. "
		warningMessage += waitTime + "s.\n"
		warningMessage += "The following operation is taking a long time to execute:\n"
		warningMessage += "    File: " + file + "\n"
		warningMessage += "    Line: " + intToString(line) + "\n"
		warningMessage += "This can be caused by a stuck replay.\n"
		warningMessage += messageCauses
		warningMessage += "If you believe, the program is still running, you can continue to wait.\n"
		warningMessage += "If you believe, the program is stuck, you can cancel the program.\n"
		warningMessage += "If you suspect, that one of these causes is the reason for the long wait time, you can try to change the program to avoid the problem.\n"
		warningMessage += "If the problem persist, this message will be repeated every approx. 10s.\n\n"

		println(warningMessage)
	}

	return false
}

/*
 * Check if the position is in the trace.
 * Args:
 * 	file: file in which the operation is executed
 * 	line: line number of the operation
 * Return:
 * 	bool: true if the position is in the trace, false otherwise
 */
func isPositionInTrace(file string, line int) bool {
	if _, ok := traceElementPositions[file]; !ok {
		return false
	}

	if !containsInt(traceElementPositions[file], line) {
		return false
	}

	return true
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
