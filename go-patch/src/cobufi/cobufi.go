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
 * Arguments:
 * 	- file_name: The name of the file that contains the trace.
 */
func ReadTrace(file_name string) runtime.CobufiReplayTrace {
	mb := 1048576
	maxTokenSize := 1

	replayData := make(runtime.CobufiReplayTrace, 0)
	chanWithoutPartner := make(map[string]int)

	for {
		file, err := os.Open(file_name)
		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 0, maxTokenSize*mb), maxTokenSize*mb)

		for scanner.Scan() {
			l := scanner.Text()
			if l == "" {
				continue
			}
			elems := strings.Split(l, ";")
			for _, elem := range elems {
				if elem == "" {
					continue
				}
				var op runtime.ReplayOperation
				var file string
				var line int
				var pFile string
				var pLine int
				var blocked = false
				var suc = true
				var selIndex int
				fields := strings.Split(elem, ",")
				time, _ := strconv.Atoi(fields[1])
				if time == 39 {
					println(elem + "\n\n")
				}
				switch fields[0] {
				case "G":
					op = runtime.CobufiReplaySpawn
					pos := strings.Split(fields[3], ":")
					file = pos[0]
					line, _ = strconv.Atoi(pos[1])
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
					if op == runtime.CobufiReplayChannelSend || op == runtime.CobufiReplayChannelRecv {
						index := findReplayPartner(fields[3], fields[6], len(replayData), chanWithoutPartner)
						if index != -1 {
							pFile = replayData[index].File
							pLine = replayData[index].Line
							replayData[index].PFile = file
							replayData[index].PLine = line
						}
					}
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
					op = runtime.CobufiReplayOnce
					if fields[2] == "0" {
						blocked = true
					}
					if fields[4] == "f" {
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
				case "S": // TODO: (cobufi) get correct select case
					cases := strings.Split(fields[4], "~")
					if cases[len(cases)-1] == "D" {
						op = runtime.CobufiReplaySelectDefault
					} else {
						op = runtime.CobufiReplaySelectCase
					}
					if fields[2] == "0" {
						blocked = true
					}
					selIndex, _ = strconv.Atoi(fields[5])
					pos := strings.Split(fields[6], ":")
					file = pos[0]
					line, _ = strconv.Atoi(pos[1])
				}
				if op != runtime.CobufiNone {
					replayData = append(replayData, runtime.ReplayElement{
						Op: op, Time: time, File: file, Line: line,
						Blocked: blocked, Suc: suc, PFile: pFile, PLine: pLine,
						SelIndex: selIndex})
				}
			}
		}

		if err := scanner.Err(); err != nil {
			if err == bufio.ErrTooLong {
				maxTokenSize *= 2 // max buffer was to short, restart
				println("Increase max token size to " + strconv.Itoa(maxTokenSize) + "MB")
				replayData = make(runtime.CobufiReplayTrace, 0)
				chanWithoutPartner = make(map[string]int)
			} else {
				panic(err)
			}
		} else {
			break // read was successful
		}
	}

	// sort data by tpre
	sortReplayDataByTime(replayData)

	// when the program runs the tracer, the once in
	// src/internal/poll/fd_poll_runtime.go was already called by the trace
	// reader. To prevent a block by the following two mutex operations,
	// we change the trace data.
	fixOnceFdPollRuntime(replayData)

	// remove the first 10 elements from the trace. They are part of the go init
	// and are therefore always called, before the program starts.
	// Because we enable the replay in the program, we must ignore them.
	for elem := range replayData {
		println(replayData[elem].Time, replayData[elem].Op, replayData[elem].File, replayData[elem].Line, replayData[elem].Blocked, replayData[elem].Suc)
	}
	println("\n\n")
	return replayData
}

/*
 * Find the partner of a channel operation.
 * The partner is the operation that is executed on the other side of the channel.
 * The partner is identified by the channel id and the operation id.
 * The index is the index of the operation in the replay data structure.
 * The function returns the index of the partner operation.
 * If the partner operation is not found, the function returns -1.
 */
func findReplayPartner(cId string, oId string, index int, chanWithoutPartner map[string]int) int {
	opString := cId + ":" + oId
	if index, ok := chanWithoutPartner[opString]; ok {
		delete(chanWithoutPartner, opString)
		return index
	} else {
		chanWithoutPartner[opString] = index
		return -1
	}
}

/*
 * Sort the replay data structure by time.
 * The function returns the sorted replay data structure.
 */
func sortReplayDataByTime(replayData runtime.CobufiReplayTrace) runtime.CobufiReplayTrace {
	sort.Slice(replayData, func(i, j int) bool {
		return replayData[i].Time < replayData[j].Time
	})
	return replayData
}

/*
 * For reading a file, a global once in `internal/poll/fd_poll_runtime.go` is
 * used. This once is already called by the trace reader for the replay. Because
 * of the mutexes in the once, this can caused a block in the replay, if
 * the once was, in the recorded run, called by the program. To prevent this,
 * we adapt the trace data by removing the mutex operations in the once from
 * the recorded trace.
 * TODO: (cobufi) does this work when we replay ans simultaneously record?
 */
func fixOnceFdPollRuntime(replayData runtime.CobufiReplayTrace) {
	for i := 0; i < len(replayData); i++ {
		if !(replayData[i].Op == runtime.CobufiReplayOnce &&
			strings.HasSuffix(replayData[i].File, "internal/poll/fd_poll_runtime.go") &&
			replayData[i].Line == 39 &&
			replayData[i].Suc) {
			continue
		}
		replayData[i].Suc = false
		for j := i + 1; j < len(replayData); j++ {
			if !(replayData[j].Op == runtime.CobufiReplayMutexLock &&
				strings.HasSuffix(replayData[j].File, "sync/once.go") &&
				replayData[j].Line == 111) {
				continue
			}
			replayData = append(replayData[:j], replayData[j+1:]...)
			for k := j; k < len(replayData); k++ {
				if !(replayData[j].Op == runtime.CobufiReplayMutexUnlock &&
					strings.HasSuffix(replayData[j].File, "sync/once.go") &&
					(replayData[j].Line == 117 || replayData[j].Line == 121)) {
					continue
				}
				replayData = append(replayData[:k], replayData[k+1:]...)
				return
			}
		}
	}
}
