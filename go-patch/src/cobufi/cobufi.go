package cobufi

import (
	"bufio"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

/*
 * Write the trace of the program to a file.
 * The trace is written in the file named file_name.
 * The trace is written in the format of CoBufi.
 */
func CreateTrace(file_name string) {
	runtime.DisableTrace()

	os.Remove(file_name)
	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	numRout := runtime.GetNumberOfRoutines()
	for i := 1; i <= numRout; i++ {
		cobufiChan := make(chan string)
		go func() {
			runtime.TraceToStringByIdChannel(i, cobufiChan)
			close(cobufiChan)
		}()
		for trace := range cobufiChan {
			if _, err := file.WriteString(trace); err != nil {
				panic(err)
			}
		}
		if _, err := file.WriteString("\n"); err != nil {
			panic(err)
		}
	}
}

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
func ReadTrace(file_name string) runtime.CobufiReplayTrace {
	file, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}

	replayData := make(runtime.CobufiReplayTrace, 0)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		l := scanner.Text()
		elems := strings.Split(l, ";")
		for _, elem := range elems {
			var op runtime.ReplayOperation
			var file string
			var line int
			var blocked bool = false
			var suc bool = true
			fields := strings.Split(elem, ",")
			time, _ := strconv.Atoi(fields[1])
			switch fields[0] {
			case "C":
				switch fields[4] {
				case "S":
					op = runtime.CobufiReplayChannelSend
				case "R":
					op = runtime.CobufiReplayChannelRecv
				case "C":
					op = runtime.CobufiReplayChannelClose
				default:
					panic("Unknown channel operation")
				}
				if fields[2] == "0" {
					blocked = true
				}
				pos := strings.Split(fields[8], ":")
				file = pos[0]
				line, _ = strconv.Atoi(pos[1])
			case "M":
				switch fields[5] {
				case "L":
					op = runtime.CobufiReplayMutexLock
				case "U":
					op = runtime.CobufiReplayMutexUnlock
				case "T":
					op = runtime.CobufiReplayMutexTryLock
				case "R":
					op = runtime.CobufiReplayRWMutexRLock
				case "N":
					op = runtime.CobufiReplayRWMutexRUnlock
				case "Y":
					op = runtime.CobufiReplayRWMutexTryRLock
				default:
					panic("Unknown mutex operation")
				}
				if fields[2] == "0" {
					blocked = true
				}
				if fields[6] == "f" {
					suc = false
				}
				pos := strings.Split(fields[7], ":")
				file = pos[0]
				line, _ = strconv.Atoi(pos[1])
			case "O":
				if fields[2] == "0" {
					blocked = true
				}
				if fields[6] == "f" {
					suc = false
				}
				pos := strings.Split(fields[5], ":")
				file = pos[0]
				line, _ = strconv.Atoi(pos[1])
			case "W":
				switch fields[4] {
				case "W":
					op = runtime.CobufiReplayWaitgroupWait
				case "A":
					op = runtime.CobufiReplayWaitgroupAddDone
				default:
					panic("Unknown waitgroup operation")
				}
				if fields[2] == "0" {
					blocked = true
				}
				pos := strings.Split(fields[7], ":")
				file = pos[0]
				line, _ = strconv.Atoi(pos[1])
			case "S": // TODO: get correct select case
				// cases := splitString(fields[4], "~")
				// if cases[len(cases)-1] == "D" {
				// 	op = selectDef
				// } else {
				// 	op = selectCase
				// }
				op = runtime.CobufiReplaySelect
				if fields[2] == "0" {
					blocked = true
				}
				pos := strings.Split(fields[5], ":")
				file = pos[0]
				line, _ = strconv.Atoi(pos[1])
			}
			replayData = append(replayData, runtime.ReplayElement{
				Op: op, Time: time, File: file, Line: line, Blocked: blocked, Suc: suc})
		}
	}

	// sort data by tpre
	sortReplayDataByTime(replayData)

	return replayData
}

func sortReplayDataByTime(replayData runtime.CobufiReplayTrace) runtime.CobufiReplayTrace {
	sort.Slice(replayData, func(i, j int) bool {
		return replayData[i].Time < replayData[j].Time
	})
	return replayData
}
