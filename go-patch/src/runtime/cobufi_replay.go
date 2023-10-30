package runtime

type ReplayOperation int

const (
	CobufiReplaySpawn ReplayOperation = iota

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
	File    string
	Line    int
	Blocked bool
	Suc     bool
}

type replay []ReplayElement

var replayEnabled bool = false
var replayLock mutex
var replayIndex int = 0

var replayData replay = make(replay, 0)

/*
 * Import the trace.
 * The function creates the replay data structure, that is used to replay the trace.
 * We only store the information that is needed to replay the trace.
 * This includes operations on
 *  - spawn
 * 	- channels
 * 	- mutexes
 * 	- once
 * 	- waitgroups
 * 	- select
 * For now we ignore atomic operations.
 * We only record the relevant information for each operation.
 */
func ImportTrace(trace []string) {
	for _, line := range trace {
		elems := splitString(line, ";")
		for _, elem := range elems {
			var op ReplayOperation
			var file string
			var line int
			var blocked bool = false
			var suc bool = true
			fields := splitString(elem, ",")
			switch fields[0] {
			case "C":
				switch fields[4] {
				case "S":
					op = CobufiReplayChannelSend
				case "R":
					op = CobufiReplayChannelRecv
				case "C":
					op = CobufiReplayChannelClose
				default:
					panic("Unknown channel operation")
				}
				if fields[2] == "0" {
					blocked = true
				}
				pos := splitString(fields[8], ":")
				file = pos[0]
				line = stringToInt(pos[1])
			case "M":
				switch fields[5] {
				case "L":
					op = CobufiReplayMutexLock
				case "U":
					op = CobufiReplayMutexUnlock
				case "T":
					op = CobufiReplayMutexTryLock
				case "R":
					op = CobufiReplayRWMutexRLock
				case "N":
					op = CobufiReplayRWMutexRUnlock
				case "Y":
					op = CobufiReplayRWMutexTryRLock
				default:
					panic("Unknown mutex operation")
				}
				if fields[2] == "0" {
					blocked = true
				}
				if fields[6] == "f" {
					suc = false
				}
				pos := splitString(fields[7], ":")
				file = pos[0]
				line = stringToInt(pos[1])
			case "O":
				if fields[2] == "0" {
					blocked = true
				}
				if fields[6] == "f" {
					suc = false
				}
				pos := splitString(fields[5], ":")
				file = pos[0]
				line = stringToInt(pos[1])
			case "W":
				switch fields[4] {
				case "W":
					op = CobufiReplayWaitgroupWait
				case "A":
					op = CobufiReplayWaitgroupAddDone
				default:
					panic("Unknown waitgroup operation")
				}
				if fields[2] == "0" {
					blocked = true
				}
				pos := splitString(fields[7], ":")
				file = pos[0]
				line = stringToInt(pos[1])
			case "S": // TODO: get correct select case
				// cases := splitString(fields[4], "~")
				// if cases[len(cases)-1] == "D" {
				// 	op = selectDef
				// } else {
				// 	op = selectCase
				// }
				op = CobufiReplaySelect
				if fields[2] == "0" {
					blocked = true
				}
				pos := splitString(fields[5], ":")
				file = pos[0]
				line = stringToInt(pos[1])
			}
			replayData = append(replayData, ReplayElement{
				Op: op, File: file, Line: line, Blocked: blocked, Suc: suc})
		}
	}
	replayEnabled = true
}

/*
 * Wait until the correct operation is about to be executed.
 * Arguments:
 * 	op: the operation type that is about to be executed
 * 	skip: number of stack frames to skip
 * Return:
 * 	bool: true if trace replay is enabled, false otherwise
 * 	replayElement, the corresponding replayElement
 */
func WaitForReplay(op ReplayOperation, skip int) (bool, ReplayElement) {
	if !replayEnabled {
		return false, ReplayElement{}
	}
	_, file, line, _ := Caller(skip)
	println(file, line)
	next := getNextReplayElement()
	for {
		if next.Op != op || next.File != file || next.Line != line {
			// TODO: very stupid sleep, find a better solution,
			// TODO: problem is that both the sleep and syscall packages cannot be used (cyclic import)
			for i := 0; i < 100000; i++ {
				_ = i
			}
			continue
		}
		break
	}
	lock(&replayLock)
	defer unlock(&replayLock)
	replayIndex++
	return true, next
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
