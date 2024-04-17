package advocate

import (
	"bufio"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// ============== Recording =================

var backgroundMemoryTestRunning = true
var traceFileCounter = 0

/*
 * Write the trace of the program to a file.
 * The trace is written in the file named file_name.
 * The trace is written in the format of advocate.
 */
func Finish() {
	runtime.DisableTrace()
	backgroundMemoryTestRunning = false

	writeToTraceFiles()
	// deleteEmptyFiles()
}

/*
 * WaitForReplayFinish waits for the replay to finish.
 */
func WaitForReplayFinish() {
	runtime.WaitForReplayFinish()
}

func writeTraceIfFull() {
	var m runtime.MemStats
	var stat syscall.Sysinfo_t
	err := syscall.Sysinfo(&stat)
	if err != nil {
		panic(err)
	}
	totalRAM := stat.Totalram
	for {
		time.Sleep(5 * time.Second)
		// get the amount of free space on ram
		runtime.ReadMemStats(&m)
		heap := m.Alloc
		stack := m.StackSys
		freeRam := stat.Totalram - heap - stack
		println(toGB(heap), toGB(stack), toGB(freeRam), toGB(stat.Totalram/5))
		if freeRam < totalRAM/5 {
			println("Memory full")
			cleanTrace()
			time.Sleep(15 * time.Second)
		}

		// end if background memory test is not running anymore
		if !backgroundMemoryTestRunning {
			break
		}
	}
}

func toGB(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024 / 1024
}

// BUG: crashes bug
func cleanTrace() {
	println("Cleaning trace")
	// // stop new element from been added to the trace
	// runtime.BlockTrace()
	// writeToTraceFiles()
	// runtime.DeleteTrace()
	// runtime.UnblockTrace()
	// runtime.GC()
}

/*
 * Write the trace to a set of files. The traces are written into a folder
 * with name trace. For each routine, a file is created. The file is named
 * trace_routineId.log. The trace of the routine is written into the file.
 * Args:
 * 	- remove: If true, and a file with the same name already exists, the file is removed, before the trace is written.
 *      If false, the trace is appended to the file.
 */
func writeToTraceFiles() {
	numRout := runtime.GetNumberOfRoutines()
	for i := 1; i <= numRout; i++ {
		// write the trace to the file
		writeToTraceFile(i) // TODO: make this parallel
	}
}

/*
 * Write the trace of a routine to a file.
 * The trace is written in the file named trace_routineId.log.
 * The trace is written in the format of advocate.
 * Args:
 * 	- routine: The id of the routine
 */
func writeToTraceFile(routine int) {
	// create the file if it does not exist and open it
	fileName := "trace/trace_" + strconv.Itoa(routine) + ".log"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// get the runtime to send the trace
	advocateChan := make(chan string)
	go func() { // TODO: do not record this routines or
		runtime.TraceToStringByIDChannel(routine, advocateChan)
		close(advocateChan)
	}()

	// receive the trace and write it to the file
	for trace := range advocateChan {
		if _, err := file.WriteString(trace); err != nil {
			panic(err)
		}
	}
}

/*
 * Delete empty files in the trace folder.
 * The function deletes all files in the trace folder that are empty.
 */
func deleteEmptyFiles() {
	files, err := os.ReadDir("trace")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		stat, err := os.Stat("trace/" + file.Name())
		if err != nil {
			continue
		}
		if stat.Size() == 0 {
			err := os.Remove("trace/" + file.Name())
			if err != nil {
				panic(err)
			}
		}
	}

}

/*
 * InitTracing initializes the tracing.
 * The function creates the trace folder and starts the background memory test.
 * Args:
 * 	- size: The size of the channel used for recording atomic events.
 */
func InitTracing(size int) {
	// remove the trace folder if it exists
	err := os.RemoveAll("trace")
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	}

	// create the trace folder
	err = os.Mkdir("trace", 0755)
	if err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}

	// if the program is terminated by the user, the defer in the header
	// is not executed. Therefore capture the signal and write the trace.
	interuptSignal := make(chan os.Signal, 1)
	signal.Notify(interuptSignal, os.Interrupt)
	go func() {
		<-interuptSignal
		Finish()
		os.Exit(0)
	}()

	go writeTraceIfFull()
	runtime.InitAdvocate(size)
}

// ============== Reading =================

var timeout = false
var tracePath = "rewritten_trace"

/*
 * Read the trace from the trace folder.
 * The function reads all files in the trace folder and adds the trace to the runtime.
 * The trace is added to the runtime by calling the AddReplayTrace function.
 */
func EnableReplay() {
	if _, err := os.Stat(tracePath); os.IsNotExist(err) {
		tracePath = "trace"
	}

	// if trace folder does not exist, panic
	if _, err := os.Stat(tracePath); os.IsNotExist(err) {
		panic("Trace folder does not exist.")
	}

	println("Reading trace from " + tracePath)

	// traverse all files in the trace folder
	files, err := os.ReadDir(tracePath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		// if the file is a directory, ignore it
		if file.IsDir() {
			continue
		}

		// if the file is a log file, read the trace
		if strings.HasSuffix(file.Name(), ".log") {
			routineID, trace := readTraceFile(tracePath + "/" + file.Name())
			runtime.AddReplayTrace(uint64(routineID), trace)
		}
	}

	runtime.EnableReplay(timeout)
}

func EnableReplayWithTimeout() {
	timeout = true
	EnableReplay()
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
 * Args:
 * 	- fileName: The name of the file that contains the trace.
 * Returns:
 * 	The routine id
 * 	The trace for this routine
 */
func readTraceFile(fileName string) (int, runtime.AdvocateReplayTrace) {
	mb := 1048576
	maxTokenSize := 1

	// get the routine id from the file name
	routineID, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(fileName, tracePath+"/trace_"), ".log"))
	if err != nil {
		panic(err)
	}

	replayData := make(runtime.AdvocateReplayTrace, 0)
	chanWithoutPartner := make(map[string]int)
	var routine = 0

	for {
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 0, maxTokenSize*mb), maxTokenSize*mb)

		for scanner.Scan() {
			routine += 1
			l := scanner.Text()
			if l == "" {
				continue
			}
			elems := strings.Split(l, ";")
			for _, elem := range elems {
				if elem == "" {
					continue
				}
				var time int
				var op runtime.Operation
				var file string
				var line int
				var pFile string
				var pLine int
				var blocked = false
				var suc = true
				var selIndex int
				fields := strings.Split(elem, ",")
				time, _ = strconv.Atoi(fields[1])
				switch fields[0] {
				case "X": // disable replay
					if fields[2] == "s" {
						op = runtime.OperationReplayStart
					} else {
						op = runtime.OperationReplayEnd
					}
				case "G":
					op = runtime.OperationSpawn
					// time, _ = strconv.Atoi(fields[1])
					pos := strings.Split(fields[3], ":")
					file = pos[0]
					line, _ = strconv.Atoi(pos[1])
				case "C":
					switch fields[4] {
					case "S":
						op = runtime.OperationChannelSend
					case "R":
						op = runtime.OperationChannelRecv
					case "C":
						op = runtime.OperationChannelClose
					default:
						panic("Unknown channel operation " + fields[4] + " in line " + elem + " in file " + fileName + ".")
					}
					// time, _ = strconv.Atoi(fields[2])
					if time == 0 {
						blocked = true
					}
					pos := strings.Split(fields[8], ":")
					file = pos[0]
					line, _ = strconv.Atoi(pos[1])
					if op == runtime.OperationChannelSend || op == runtime.OperationChannelRecv {
						index := findReplayPartner(fields[3], fields[6], len(replayData), chanWithoutPartner)
						if index != -1 {
							pFile = replayData[index].File
							pLine = replayData[index].Line
							replayData[index].PFile = file
							replayData[index].PLine = line
						}
					}
				case "M":
					rw := false
					if fields[4] == "R" {
						rw = true
					}
					// time, _ = strconv.Atoi(fields[2])
					if fields[6] == "f" {
						suc = false
					}
					pos := strings.Split(fields[7], ":")
					file = pos[0]
					line, _ = strconv.Atoi(pos[1])
					switch fields[5] {
					case "L":
						if rw {
							op = runtime.OperationRWMutexLock
						} else {
							op = runtime.OperationMutexLock
							time = swapTimerRwMutex("L", time, file, line, &replayData)
						}
					case "U":
						if rw {
							op = runtime.OperationRWMutexUnlock
						} else {
							op = runtime.OperationMutexUnlock
							time = swapTimerRwMutex("U", time, file, line, &replayData)
						}
					case "T":
						if rw {
							op = runtime.OperationRWMutexTryLock
						} else {
							op = runtime.OperationMutexTryLock
							time = swapTimerRwMutex("T", time, file, line, &replayData)
						}
					case "R":
						op = runtime.OperationRWMutexRLock
					case "N":
						op = runtime.OperationRWMutexRUnlock
					case "Y":
						op = runtime.OperationRWMutexTryRLock
					default:
						panic("Unknown mutex operation")
					}
					if fields[2] == "0" {
						blocked = true
					}

				case "O":
					op = runtime.OperationOnce
					// time, _ = strconv.Atoi(fields[1]) // read tpre to prevent false order
					if time == 0 {
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
						op = runtime.OperationWaitgroupWait
					case "A":
						op = runtime.OperationWaitgroupAddDone
					default:
						panic("Unknown waitgroup operation")
					}
					// time, _ = strconv.Atoi(fields[2])
					if time == 0 {
						blocked = true
					}
					pos := strings.Split(fields[7], ":")
					file = pos[0]
					line, _ = strconv.Atoi(pos[1])
				case "S":
					cases := strings.Split(fields[4], "~")
					if cases[len(cases)-1] == "D" {
						op = runtime.OperationSelectDefault
					} else {
						op = runtime.OperationSelectCase
					}
					// time, _ = strconv.Atoi(fields[2])
					if time == 0 {
						blocked = true
					}
					selIndex, _ = strconv.Atoi(fields[5])
					pos := strings.Split(fields[6], ":")
					file = pos[0]
					line, _ = strconv.Atoi(pos[1])
				case "N":
					switch fields[4] {
					case "W":
						op = runtime.OperationCondWait
					case "S":
						op = runtime.OperationCondSignal
					case "B":
						op = runtime.OperationCondBroadcast
					default:
						panic("Unknown cond operation")
					}
					pos := strings.Split(fields[5], ":")
					file = pos[0]
					line, _ = strconv.Atoi(pos[1])
					if fields[2] == "0" {
						blocked = true
					}

				default:
					panic("Unknown operation " + fields[0] + " in line " + elem + " in file " + fileName + ".")
				}
				if op != runtime.OperationNone && !runtime.AdvocateIgnore(op, file, line) {
					replayData = append(replayData, runtime.ReplayElement{
						Op: op, Routine: routine, Time: time, File: file, Line: line,
						Blocked: blocked, Suc: suc, PFile: pFile, PLine: pLine,
						SelIndex: selIndex})
				}
			}
		}

		if err := scanner.Err(); err != nil {
			if err == bufio.ErrTooLong {
				maxTokenSize *= 2 // max buffer was to short, restart
				println("Increase max token size to " + strconv.Itoa(maxTokenSize) + "MB")
				replayData = make(runtime.AdvocateReplayTrace, 0)
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

	return routineID, replayData
}

// TODO: swap timer for rwmutix.Trylock
// TODO: LOCAL
func swapTimerRwMutex(op string, time int, file string, line int, replayData *runtime.AdvocateReplayTrace) int {
	if op == "L" {
		if !strings.HasSuffix(file, "sync/rwmutex.go") || line != 266 {
			return time
		}

		for i := len(*replayData) - 1; i >= 0; i-- {
			timeNew := (*replayData)[i].Time
			(*replayData)[i].Time = time
			return timeNew
		}
	} else if op == "U" {
		if !strings.HasSuffix(file, "sync/rwmutex.go") {
			return time
		}

		if line == 390 {
			for i := len(*replayData) - 1; i >= 0; i-- {

				timeNew := (*replayData)[i].Time
				(*replayData)[i].Time = time
				return timeNew
			}
		}
	}

	return time
}

/*
 * Find the partner of a channel operation.
 * The partner is the operation that is executed on the other side of the channel.
 * The partner is identified by the channel id and the operation id.
 * The index is the index of the operation in the replay data structure.
 * The function returns the index of the partner operation.
 * If the partner operation is not found, the function returns -1.
 */
func findReplayPartner(cID string, oID string, index int, chanWithoutPartner map[string]int) int {
	opString := cID + ":" + oID
	if index, ok := chanWithoutPartner[opString]; ok {
		delete(chanWithoutPartner, opString)
		return index
	}

	chanWithoutPartner[opString] = index
	return -1

}

/*
 * Sort the replay data structure by time.
 * The function returns the sorted replay data structure.
 */
func sortReplayDataByTime(replayData runtime.AdvocateReplayTrace) runtime.AdvocateReplayTrace {
	sort.Slice(replayData, func(i, j int) bool {
		return replayData[i].Time < replayData[j].Time
	})
	return replayData
}
