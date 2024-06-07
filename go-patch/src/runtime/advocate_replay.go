package runtime

const (
	ExitCodeDefault        = 0
	ExitCodeStuckFinish    = 10
	ExitCodeStuckWaitElem  = 11
	ExitCodeStuckNoElem    = 12
	ExitCodeElemEmptyTrace = 13
	ExitCodeLeakUnbuf      = 20
	ExitCodeLeakBuf        = 21
	ExitCodeLeakMutex      = 22
	ExitCodeLeakCond       = 23
	ExitCodeLeakWG         = 24
	ExitCodeSendClose      = 30
	ExitCodeRecvClose      = 31
	ExitCodeNegativeWG     = 32
	ExitCodeCyclic         = 41
)

var ExitCodeNames = map[int]string{
	0:  "The replay terminated without finding a Replay element",
	10: "Replay Stuck: Long wait time for finishing replay",
	11: "Replay Stuck: Long wait time for running element",
	12: "Replay Stuck: No traced operation has been executed for approx. 20s",
	13: "The program tried to execute an operation, although all elements in the trace have already been executed.",
	20: "Leak: Leaking unbuffered channel or select was unstuck",
	21: "Leak: Leaking buffered channel was unstuck",
	22: "Leak: Leaking Mutex was unstuck",
	23: "Leak: Leaking Cond was unstuck",
	24: "Leak: Leaking WaitGroup was unstuck",
	30: "Send on close",
	31: "Receive on close",
	32: "Negative WaitGroup counter",
}

/*
 * String representation of the replay operation.
 * Return:
 * 	string: string representation of the replay operation
 */
func (ro Operation) ToString() string {
	switch ro {
	case OperationNone:
		return "OperationNone"
	case OperationSpawn:
		return "OperationSpawn"
	case OperationSpawned:
		return "OperationSpawned"
	case OperationChannelSend:
		return "OperationChannelSend"
	case OperationChannelRecv:
		return "OperationChannelRecv"
	case OperationChannelClose:
		return "OperationChannelClose"
	case OperationMutexLock:
		return "OperationMutexLock"
	case OperationMutexUnlock:
		return "OperationMutexUnlock"
	case OperationMutexTryLock:
		return "OperationMutexTryLock"
	case OperationRWMutexLock:
		return "OperationRWMutexLock"
	case OperationRWMutexUnlock:
		return "OperationRWMutexUnlock"
	case OperationRWMutexTryLock:
		return "OperationRWMutexTryLock"
	case OperationRWMutexRLock:
		return "OperationRWMutexRLock"
	case OperationRWMutexRUnlock:
		return "OperationRWMutexRUnlock"
	case OperationRWMutexTryRLock:
		return "OperationRWMutexTryRLock"
	case OperationOnce:
		return "OperationOnce"
	case OperationWaitgroupAddDone:
		return "OperationWaitgroupAddDone"
	case OperationWaitgroupWait:
		return "OperationWaitgroupWait"
	case OperationSelect:
		return "OperationSelect"
	case OperationSelectCase:
		return "OperationSelectCase"
	case OperationSelectDefault:
		return "OperationSelectDefault"
	case OperationCondSignal:
		return "OperationCondSignal"
	case OperationCondBroadcast:
		return "OperationCondBroadcast"
	case OperationCondWait:
		return "OperationCondWait"
	case OperationReplayEnd:
		return "OperationReplayEnd"
	default:
		return "Unknown"
	}
}

/*
 * The replay data structure.
 * The replay data structure is used to store the trace of the program.
 * op: identifier of the operation
 * time: int (tpre) of the operation
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
	Op       Operation
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
type AdvocateReplayTraces map[uint64]AdvocateReplayTrace // routine -> trace

var replayEnabled bool // replay is on
var replayLock mutex
var replayDone int
var replayDoneLock mutex

// read trace
var replayData = make(AdvocateReplayTraces, 0)
var numberElementsInTrace int
var traceElementPositions = make(map[string][]int) // file -> []line

// timeout
var timeoutLock mutex
var timeoutCounterGlobal = 0
var timeoutMessageCycle = 1000 // approx. 20s
var timeOutCancel = false

// exit code
var replayExitCode bool
var expectedExitCode int

/*
 * Add a replay trace to the replay data.
 * Arguments:
 * 	routine: the routine id
 * 	trace: the replay trace
 */
func AddReplayTrace(routine uint64, trace AdvocateReplayTrace) {
	if _, ok := replayData[routine]; ok {
		panic("Routine already exists")
	}
	replayData[routine] = trace

	numberElementsInTrace += len(trace)

	for _, e := range trace {
		if _, ok := traceElementPositions[e.File]; !ok {
			traceElementPositions[e.File] = make([]int, 0)
		}
		if !containsInt(traceElementPositions[e.File], e.Line) {
			traceElementPositions[e.File] = append(traceElementPositions[e.File], e.Line)
		}
	}
}

/*
 * Print the replay data.
 */
func (t AdvocateReplayTraces) Print() {
	for id, trace := range t {
		println("\nRoutine: ", id)
		trace.Print()
	}
}

/*
 * Print the replay trace for one routine.
 */
func (t AdvocateReplayTrace) Print() {
	for _, e := range t {
		println(e.Op.ToString(), e.Time, e.File, e.Line, e.Blocked, e.Suc)
	}
}

/*
 * Enable the replay.
 * Arguments:
 * 	timeout: true if the replay should be canceled after a timeout, false otherwise
 */
func EnableReplay(timeout bool) {
	timeOutCancel = timeout

	// run a background routine to check for timeout if no operation is executed
	go checkForTimeoutNoOperation()

	replayEnabled = true
	println("Replay enabled")
}

/*
 * Disable the replay. This is called when a stop character in the trace is
 * encountered.
 */
func DisableReplay() {
	lock(&replayLock)
	defer unlock(&replayLock)

	replayEnabled = false
	println("Replay disabled")
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
		if replayDone >= numberElementsInTrace {
			unlock(&replayDoneLock)
			break
		}
		unlock(&replayDoneLock)

		if !replayEnabled {
			break
		}

		// check for timeout
		if timeoutCounter%timeoutMessageCycle == 0 {
			ExitReplayWithCode(ExitCodeStuckFinish)

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
			warningMessage += "If the problem persist, this message will be repeated.\n\n"
			println(warningMessage)
			if timeOutCancel {
				panic("ReplayError: Replay stuck")
			}
		}

		slowExecution()
	}
}

func IsReplayEnabled() bool {
	return replayEnabled
}

/*
 * Wait until the correct operation is about to be executed.
 * Arguments:
 * 	op: the operation type that is about to be executed
 * 	skip: number of stack frames to skip
 * Return:
 * 	bool: true if trace replay is enabled, false otherwise
 * 	bool: true if the wait was released regularly, false if it was released
 * 		because replay is disables, the operation is ignored or the operation is invalid
 * 	chan ReplayElement: channel to receive the next replay element
 */
func WaitForReplay(op Operation, skip int) (bool, bool, ReplayElement) {
	if !replayEnabled {
		return false, false, ReplayElement{}
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

// var lastNextTime int = 0

/*
 * Wait until the correct operation is about to be executed.
 * Arguments:
 * 		op: the operation type that is about to be executed
 * 		file: file in which the operation is executed
 * 		line: line number of the operation
 * Return:
 * 	bool: true if trace replay is enabled, false otherwise
 * 	bool: true if the wait was released regularly, false if it was released
 * 		because replay is disables, the operation is ignored or the operation is invalid
 * 	chan ReplayElement: channel to receive the next replay element
 */
func WaitForReplayPath(op Operation, file string, line int) (bool, bool, ReplayElement) {
	if !replayEnabled {
		return false, false, ReplayElement{}
	}

	if AdvocateIgnoreReplay(op, file, line) {
		return true, false, ReplayElement{}
	}

	// println("Wait: ", op.ToString(), file, line)
	timeoutCounter := 0
	for {
		if !replayEnabled { // check again if disabled by command
			return false, false, ReplayElement{}
		}

		nextRoutine, next := getNextReplayElement()

		if AdvocateIgnoreReplay(next.Op, next.File, next.Line) {
			// println("Igno: ", next.Op.ToString(), next.File, next.Line)
			foundReplayElement(nextRoutine)
			continue
		}

		// disable the replay, if the next operation is the disable replay operation
		if next.Op == OperationReplayEnd {
			ExitReplayWithCode(next.Line)

			println("Stop Character Found. Disable Replay.")
			DisableReplay()
			foundReplayElement(nextRoutine)
			return false, false, ReplayElement{}
		}

		timeoutCounter++
		// all elements in the trace have been executed
		if nextRoutine == -1 {
			println("The program tried to execute an operation, although all elements in the trace have already been executed.\nDisable Replay")
			DisableReplay()
			ExitReplayWithCode(ExitCodeElemEmptyTrace)
			return false, false, ReplayElement{}
		}

		// if lastNextTime != next.Time {
		// 	println("Next: ", next.Time, next.Op.ToString(), next.File, next.Line)
		// 	lastNextTime = next.Time
		// }

		if next.Time != 0 && replayEnabled {
			if (next.Op != op && !correctSelect(next.Op, op)) ||
				next.File != file || next.Line != line {

				checkForTimeout(timeoutCounter, file, line)
				slowExecution()
				continue
			}
		}

		// println("Run : ", next.Time, next.Op.ToString(), next.File, next.Line)
		foundReplayElement(nextRoutine)

		lock(&timeoutLock)
		timeoutCounterGlobal = 0 // reset the global timeout counter
		unlock(&timeoutLock)

		lock(&replayDoneLock)
		replayDone++
		unlock(&replayDoneLock)

		return true, true, next
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
 * Return:
 * 	bool: false
 */
func checkForTimeout(timeoutCounter int, file string, line int) {
	if !replayEnabled {
		return
	}

	messageCauses := "Possible causes are:\n"
	messageCauses += "    - The program was altered between recording and replay\n"
	messageCauses += "    - The program execution path is not deterministic, e.g. its execution path is determined by a random number\n"
	messageCauses += "    - The program execution path depends on the order of not tracked operations\n"
	messageCauses += "    - The program execution depends on outside input, that was not exactly reproduced\n"

	if timeoutCounter == 500 { // ca. 10s
		// res := isPositionInTrace(file, line)
		// if !res {
		// 	errorMessage := "ReplayError: Program tried to execute an operation that is not in the trace:\n"
		// 	errorMessage += "    File: " + file + "\n"
		// 	errorMessage += "    Line: " + intToString(line) + "\n"
		// 	errorMessage += "This means, that the program replay was not successful.\n"
		// 	errorMessage += messageCauses
		// 	errorMessage += "If you suspect, that one of these causes is the reason for the error, you can try to change the program to avoid the problem.\n"
		// 	errorMessage += "If this is not possible, you can try to rerun the replay, hoping the error does not occur again.\n"
		// 	errorMessage += "If this is not possible or does not work, the program replay is currently not possible.\n\n"

		// 	panic(errorMessage)
		// }
		// } else
		if timeoutCounter%timeoutMessageCycle == 0 { // approx. every 20s
			warningMessage := "\nReplayWarning: Long wait time\n"
			warningMessage += "The following operation is taking a long time to execute:\n"
			warningMessage += "    File: " + file + "\n"
			warningMessage += "    Line: " + intToString(line) + "\n"
			warningMessage += "This can be caused by a stuck replay.\n"
			warningMessage += messageCauses
			warningMessage += "If you believe, the program is still running, you can continue to wait.\n"
			warningMessage += "If you believe, the program is stuck, you can cancel the program.\n"
			warningMessage += "If you suspect, that one of these causes is the reason for the long wait time, you can try to change the program to avoid the problem.\n"
			warningMessage += "If the problem persist, this message will be repeated.\n\n"

			println(warningMessage)

			ExitReplayWithCode(ExitCodeStuckWaitElem)

			if timeOutCancel {
				panic("ReplayError: Replay stuck")
			}
		}
	}
}

func checkForTimeoutNoOperation() {
	if !replayEnabled {
		return
	}

	waitTime := 800 // approx. 20s
	warningMessage := "No traced operation has been executed for a long time.\n"
	warningMessage += "This can be caused by a stuck replay.\n"
	warningMessage += "Possible causes are:\n"
	warningMessage += "    - The program was altered between recording and replay\n"
	warningMessage += "    - The program execution path is not deterministic, e.g. its execution path is determined by a random number\n"
	warningMessage += "    - The program execution path depends on the order of not tracked operations\n"
	warningMessage += "    - The program execution depends on outside input, that was not exactly reproduced\n"
	warningMessage += "If you believe, the program is still running, you can continue to wait.\n"
	warningMessage += "If you believe, the program is stuck, you can cancel the program.\n"
	warningMessage += "If you suspect, that one of these causes is the reason for the long wait time, you can try to change the program to avoid the problem.\n"
	warningMessage += "If the problem persist, this message will be repeated.\n\n"

	for {
		lock(&timeoutLock)
		timeoutCounterGlobal++
		timeoutCounter := timeoutCounterGlobal
		unlock(&timeoutLock)

		if !replayEnabled {
			break
		}

		if timeoutCounter%waitTime == 0 {
			message := "\nReplayWarning: Long wait time\n"
			message += warningMessage

			println(message)
			ExitReplayWithCode(ExitCodeStuckNoElem)
			if timeOutCancel {
				panic("ReplayError: Replay stuck")
			}
		}
		slowExecution()
	}
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

func correctSelect(next Operation, op Operation) bool {
	if op != OperationSelect {
		return false
	}

	if next != OperationSelectCase && next != OperationSelectDefault {
		return false
	}

	return true
}

func BlockForever() {
	gopark(nil, nil, waitReasonZero, traceBlockForever, 1)
}

/*
 * Get the next replay element.
 * Return:
 * 	uint64: the routine of the next replay element or -1 if the trace is empty
 * 	ReplayElement: the next replay element
 */
func getNextReplayElement() (int, ReplayElement) {
	lock(&replayLock)
	defer unlock(&replayLock)

	routine := -1
	// set mintTime to max int
	var minTime int = -1

	for id, trace := range replayData {
		if len(trace) == 0 {
			continue
		}
		elem := trace[0]
		if minTime == -1 || elem.Time < minTime {
			minTime = elem.Time
			routine = int(id)
		}
	}

	if routine == -1 {
		return -1, ReplayElement{}
	}

	return routine, replayData[uint64(routine)][0]
}

/*
 * Check if the next element in the trace is a replay end element with the given code.
 * Args:
 * 	code: the code of the replay end element
 * 	runExit: true if the program should exit with the given code, false otherwise
 *  overwrite: if true, also exit if the next element is not a replay end element but the code is the expected exit code
 * Return:
 * 	bool: true if the next element is a replay end element with the given code or id overwrite is set and the code is the expected code, false otherwise
 */
func IsNextElementReplayEnd(code int, runExit bool, overwrite bool) bool {
	_, next := getNextReplayElement()

	if overwrite && code == expectedExitCode {
		ExitReplayWithCode(code)
		return true
	}

	if next.Op != OperationReplayEnd || next.Line != code {
		return false
	}

	if runExit {
		ExitReplayWithCode(code)
	}

	return true
}

func foundReplayElement(routine int) {
	lock(&replayLock)
	defer unlock(&replayLock)

	// remove the first element from the trace for the routine
	replayData[uint64(routine)] = replayData[uint64(routine)][1:]
}

func SetExitCode(code bool) {
	replayExitCode = code
}

func SetExpectedExitCode(code int) {
	expectedExitCode = code
}

func ExitReplayWithCode(code int) {
	if replayExitCode {
		println("Exit Replay with code ", code, ExitCodeNames[code])
		exit(int32(code))
	}
}
