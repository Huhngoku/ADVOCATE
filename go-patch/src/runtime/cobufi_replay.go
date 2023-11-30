package runtime

type ReplayOperation int

const (
	CobufiNone ReplayOperation = iota
	CobufiReplaySpawn

	CobufiReplayChannelSend
	CobufiReplayChannelRecv
	CobufiReplayChannelClose

	CobufiReplayMutexLock
	CobufiReplayMutexUnlock
	CobufiReplayMutexTryLock
	CobufiReplayRWMutexLock
	CobufiReplayRWMutexUnlock
	CobufiReplayRWMutexTryLock
	CobufiReplayRWMutexRLock
	CobufiReplayRWMutexRUnlock
	CobufiReplayRWMutexTryRLock

	CobufiReplayOnce

	CobufiReplayWaitgroupAddDone
	CobufiReplayWaitgroupWait

	CobufiReplaySelect
	CobufiReplaySelectCase
	CobufiReplaySelectDefault
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

type CobufiReplayTrace []ReplayElement

var replayEnabled bool = false
var replayLock mutex
var replayIndex int = 0

var replayData CobufiReplayTrace = make(CobufiReplayTrace, 0)

func (t CobufiReplayTrace) Print() {
	for _, e := range t {
		println(e.Op, e.Time, e.File, e.Line, e.Blocked, e.Suc)
	}
}

func EnableReplay(trace CobufiReplayTrace) {
	replayData = trace
	if len(replayData) == 0 {
		println("Trace is empty. Disable replay")
	}
	replayEnabled = true
}

func DisableReplay() {
	replayEnabled = false
}

func WaitForReplayFinish() {
	for {
		lock(&replayLock)
		if replayIndex >= len(replayData) {
			unlock(&replayLock)
			break
		}
		unlock(&replayLock)
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

	if IsIgnoredFile(file) {
		return false, ReplayElement{}
	}

	println("WaitForReplayPath", op, file, line)
	for {
		next := getNextReplayElement()

		if next.Time != 0 { // if next == ReplayElement{}
			if (next.Op != op && !correctSelect(next.Op, op)) ||
				next.File != file || next.Line != line {
				// TODO: sleep here to not waste CPU
				continue
			}
		}

		lock(&replayLock)
		replayIndex++
		unlock(&replayLock)
		println("Replay: ", next.Time, op, file, line)
		return true, next
	}
}

func correctSelect(next ReplayOperation, op ReplayOperation) bool {
	if op != CobufiReplaySelect {
		return false
	}

	if next != CobufiReplaySelectCase && next != CobufiReplaySelectDefault {
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
	}
	return replayData[replayIndex]
}
