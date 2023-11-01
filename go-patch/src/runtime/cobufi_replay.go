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
 */
type ReplayElement struct {
	Op      ReplayOperation
	Time    int
	File    string
	Line    int
	Blocked bool
	Suc     bool
}

type CobufiReplayTrace []ReplayElement

var replayEnabled bool = false
var replayLock mutex
var replayIndex int = 0

var replayData CobufiReplayTrace = make(CobufiReplayTrace, 0)

func EnableReplay(trace CobufiReplayTrace) {
	replayData = trace
	replayEnabled = true
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
func WaitForReplay(op ReplayOperation, skip int) (bool, chan ReplayElement) {
	if !replayEnabled {
		return false, nil
	}

	c := make(chan ReplayElement, 1<<16)

	go func() {
		_, file, line, _ := Caller(skip)
		for {
			next := getNextReplayElement()
			println(next.Op, next.File, next.Line)
			if next.Op != op || next.File != file || next.Line != line {
				// TODO: very stupid sleep, find a better solution,
				// TODO: problem is that both the sleep and syscall packages cannot be used (cyclic import)
				for i := 0; i < 100000; i++ {
					_ = i
				}
				continue
			}
			c <- next
			lock(&replayLock)
			replayIndex++
			unlock(&replayLock)
			return
		}
	}()

	return true, c
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
	return replayData[replayIndex]
}
