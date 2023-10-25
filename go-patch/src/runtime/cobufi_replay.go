package runtime

type ReplayOperation int

const (
	channelSend ReplayOperation = iota
	channelRecv
	channelClose

	mutexLock
	mutexUnlock
	mutexTryLock
	mutexRLock
	mutexRUnlock
	mutexRTryLock

	onceSuccess
	onceFailure

	waitgroupAddDone
	waitgroupWait

	selectCase
	selectDef
)

type replayElement struct {
	op   ReplayOperation
	file string
	line int
}

type replay []replayElement

var replayEnabled bool = false
var replayLock mutex
var replayIndex int = 0

var replayData replay = make(replay, 0)

/*
 * Import the trace.
 * The function creates the replay data structure, that is used to replay the trace.
 * We only store the information that is needed to replay the trace.
 * This includes operations on
 * 	- channels
 * 	- mutexes
 * 	- once
 * 	- waitgroups
 * 	- select
 * We ignore all other operations, because they do not influence the execution.
 * For now we only record the type and position (file and line number) of the operation.
 */
func ImportTrace(trace []string) {
	for _, line := range trace {
		elems := splitString(line, ";")
		for _, elem := range elems {
			var op ReplayOperation
			var file string
			var line int
			fields := splitString(elem, ",")
			switch fields[0] {
			case "C":
				switch fields[4] {
				case "S":
					op = channelSend
				case "R":
					op = channelRecv
				case "C":
					op = channelClose
				default:
					panic("Unknown channel operation")
				}
				pos := splitString(fields[8], ":")
				file = pos[0]
				line = stringToInt(pos[1])
			case "M":
				switch fields[5] {
				case "L":
					op = mutexLock
				case "U":
					op = mutexUnlock
				case "T":
					op = mutexTryLock
				case "R":
					op = mutexRLock
				case "N":
					op = mutexRUnlock
				case "Y":
					op = mutexRTryLock
				default:
					panic("Unknown mutex operation")
				}
				pos := splitString(fields[7], ":")
				file = pos[0]
				line = stringToInt(pos[1])
			case "O":
				switch fields[4] {
				case "t":
					op = onceSuccess
				case "f":
					op = onceFailure
				default:
					panic("Unknown once operation")
				}
				pos := splitString(fields[5], ":")
				file = pos[0]
				line = stringToInt(pos[1])
			case "W":
				switch fields[4] {
				case "W":
					op = waitgroupWait
				case "A":
					op = waitgroupAddDone
				default:
					panic("Unknown waitgroup operation")
				}
				pos := splitString(fields[7], ":")
				file = pos[0]
				line = stringToInt(pos[1])
			case "S": // TODO: get correct select case
				cases := splitString(fields[4], "~")
				if cases[len(cases)-1] == "D" {
					op = selectDef
				} else {
					op = selectCase
				}
				pos := splitString(fields[5], ":")
				file = pos[0]
				line = stringToInt(pos[1])
			}
			replayData = append(replayData, replayElement{op, file, line})
		}
	}
	replayEnabled = true
}

func WaitForReplay(op ReplayOperation, skip int) {
	if !replayEnabled {
		return
	}
	_, file, line, _ := Caller(skip)
	println(file, line)
	next := getNextReplayElement()
	for {
		if next.op != op || next.file != file || next.line != line {
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
	return
}

/*
 * Get the next replay element.
 * The function returns the next replay element and increments the index.
 */
func getNextReplayElement() replayElement {
	lock(&replayLock)
	defer unlock(&replayLock)
	return replayData[replayIndex]
}

/*
 * Split a string by the seperator
 */
func splitString(line string, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(line); i++ {
		if line[i] == sep[0] {
			result = append(result, line[start:i])
			start = i + 1
		}
	}
	result = append(result, line[start:])
	return result
}

/*
 * Convert a string to an integer
 * Works only with positive integers
 */
func stringToInt(s string) int {
	var result int
	sign := 1
	for i := 0; i < len(s); i++ {
		if s[i] == '-' && i == 0 {
			sign = -1
		} else if s[i] >= '0' && s[i] <= '9' {
			result = result*10 + int(s[i]-'0')
		} else {
			panic("Invalid input")
		}
	}
	return result * sign
}
