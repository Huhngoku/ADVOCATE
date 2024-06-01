package logging

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var levelDebug int = 0

var Reset = "\033[0m"
var Red = "\033[31m"
var Orange = "\033[33m"
var Green = "\033[32m"
var Blue = "\033[34m"

type debugLevel int

const (
	SILENT debugLevel = iota
	ERROR
	INFO
	DEBUG
)

type resultLevel int

const (
	NONE resultLevel = iota
	CRITICAL
	WARNING
)

type ResultType string

const (
	Empty ResultType = ""

	// actual
	ASendOnClosed          ResultType = "A1"
	ARecvOnClosed          ResultType = "A2"
	ACloseOnClosed         ResultType = "A3"
	AConcurrentRecv        ResultType = "A4"
	ASelCaseWithoutPartner ResultType = "A5"

	// possible
	PSendOnClosed ResultType = "P1"
	PRecvOnClosed ResultType = "P2"
	PNegWG        ResultType = "P3"

	// leaks
	LUnbufferedWith    = "L1"
	LUnbufferedWithout = "L2"
	LBufferedWith      = "L3"
	LBufferedWithout   = "L4"
	LNilChan           = "L5"
	LSelectWith        = "L6"
	LSelectWithout     = "L7"
	LMutex             = "L8"
	LWaitGroup         = "L9"
	LCond              = "L0"
)

var resultTypeMap = map[ResultType]string{
	ARecvOnClosed:          "Found receive on closed channel:",
	ASendOnClosed:          "Found send on closed channel:",
	ACloseOnClosed:         "Found close on closed channel:",
	AConcurrentRecv:        "Found concurrent Recv on same channel:",
	ASelCaseWithoutPartner: "Found select case without partner or nil case",

	PSendOnClosed: "Possible send on closed channel:",
	PRecvOnClosed: "Possible receive on closed channel:",
	PNegWG:        "Possible negative waitgroup counter:",

	LUnbufferedWith:    "Leak on unbuffered channel or select with possible partner:",
	LUnbufferedWithout: "Leak on unbuffered channel or select with possible partner:",
	LBufferedWith:      "Leak on buffered channel with possible partner:",
	LBufferedWithout:   "Leak on unbuffered channel with possible partner:",
	LNilChan:           "Leak on nil channel:",
	LSelectWith:        "Leak on select with possible partner:",
	LSelectWithout:     "Leak on select without partner or nil case",
	LMutex:             "Leak on mutex:",
	LWaitGroup:         "Leak on wait group:",
	LCond:              "Leak on conditional variable:",
}

var outputReadableFile string
var outputMachineFile string
var foundBug = false
var resultsWarningReadable []string
var resultsCriticalReadable []string
var resultsWarningMachine []string
var resultCriticalMachine []string

/*
* Print a debug log message if the log level is sufficiant
* Args:
*   message: message to print
*   level: level of the message
 */
func Debug(message string, level debugLevel) {
	// print message to terminal
	if int(level) <= levelDebug {
		if level == ERROR {
			println(Blue + message + Reset)
		} else if level == INFO {
			println(Green + message + Reset)
		} else {
			println(message)
		}
	}
}

type ResultElem interface {
	isInvalid() bool
	stringMachine() string
	stringReadable() string
}

type TraceElementResult struct {
	RoutineID int
	ObjID     int
	TPre      int
	ObjType   string
	File      string
	Line      int
}

func (t TraceElementResult) stringMachine() string {
	return fmt.Sprintf("T:%d:%d:%d:%s:%s:%d", t.RoutineID, t.ObjID, t.TPre, t.ObjType, t.File, t.Line)
}

func (t TraceElementResult) stringReadable() string {
	return fmt.Sprintf("%s:%d", t.File, t.Line)
}

func (t TraceElementResult) isInvalid() bool {
	return t.ObjType == ""
}

type SelectCaseResult struct {
	SelID   int
	ObjID   int
	ObjType string
	Routine int
}

func (s SelectCaseResult) stringMachine() string {
	return fmt.Sprintf("S:%d:%s", s.ObjID, s.ObjType)
}

func (s SelectCaseResult) stringReadable() string {
	return fmt.Sprintf("case: %d:%s", s.ObjID, s.ObjType)
}

func (s SelectCaseResult) isInvalid() bool {
	return s.ObjType == ""
}

/*
 * Print a result message
 * Args:
 * 	level: level of the message
 *	message: message to print
 */
func Result(level resultLevel, resType ResultType, argType1 string, arg1 []ResultElem, argType2 string, arg2 []ResultElem) {
	if arg1[0].isInvalid() {
		return
	}

	// ignore signal_unix.go
	if strings.Contains(arg1[0].stringReadable(), "signal_unix.go") {
		return
	}

	foundBug = true

	resultReadable := resultTypeMap[resType] + "\n\t" + argType1 + ": "
	resultMachine := string(resType) + ","

	for i, arg := range arg1 {
		if i != 0 {
			resultReadable += ";"
			resultMachine += ";"
		}
		resultReadable += arg.stringReadable()
		resultMachine += arg.stringMachine()
	}

	resultReadable += "\n"
	if len(arg2) > 0 {
		resultReadable += "\t" + argType2
		resultMachine += ","
		for i, arg := range arg2 {
			if i != 0 {
				resultReadable += ";"
				resultMachine += ";"
			}
			resultReadable += arg.stringReadable()
			resultMachine += arg.stringMachine()
		}

	}

	resultReadable += "\n"
	resultMachine += "\n"

	if level == WARNING {
		resultsWarningReadable = append(resultsWarningReadable, resultReadable)
		resultsWarningMachine = append(resultsWarningMachine, resultMachine)
	} else if level == CRITICAL {
		resultsCriticalReadable = append(resultsCriticalReadable, resultReadable)
		resultCriticalMachine = append(resultCriticalMachine, resultMachine)
	}
}

/*
* Initialize the debug
* Args:
*   level: level of the debug
*   outReadable: path to the output file, no output file if empty
*   outMachine: path to the output file for the reordered trace, no output file if empty
 */
func InitLogging(level int, outReadable string, outMachine string) {
	if level < 0 {
		level = 0
	}
	levelDebug = level

	outputReadableFile = outReadable
	outputMachineFile = outMachine
}

/*
* Print the summary of the analysis
* Args:
*   noWarning: if true, only critical errors will be shown
*   noPrint: if true, no output will be printed to the terminal
* Returns:
*   int: number of bugs found
 */
func PrintSummary(noWarning bool, noPrint bool) int {
	counter := 1
	resMachine := ""
	resReadable := "```\n==================== Summary ====================\n\n"

	if !noPrint {
		fmt.Print("==================== Summary ====================\n\n")
	}

	found := false

	if len(resultsCriticalReadable) > 0 {
		found = true
		resReadable += "-------------------- Critical -------------------\n\n"

		if !noPrint {
			fmt.Print("-------------------- Critical -------------------\n\n")
		}

		for _, result := range resultsCriticalReadable {
			resReadable += strconv.Itoa(counter) + " " + result + "\n"

			if !noPrint {
				fmt.Println(strconv.Itoa(counter) + " " + result)
			}

			counter++
		}

		for _, result := range resultCriticalMachine {
			resMachine += result
		}
	}
	if len(resultsWarningReadable) > 0 && !noWarning {
		found = true
		resReadable += "\n-------------------- Warning --------------------\n\n"
		if !noPrint {
			fmt.Print("\n-------------------- Warning --------------------\n\n")
		}

		for _, result := range resultsWarningReadable {
			resReadable += strconv.Itoa(counter) + " " + result + "\n"

			if !noPrint {
				fmt.Println(strconv.Itoa(counter) + " " + result)
			}

			counter++
		}

		for _, result := range resultsWarningMachine {
			resMachine += result
		}
	}
	if !found {
		resReadable += "No bugs found" + "\n"

		if !noPrint {
			fmt.Println("No bugs found")
		}
	}

	resReadable += "```"

	// write output readable
	if _, err := os.Stat(outputReadableFile); err == nil {
		if err := os.Remove(outputReadableFile); err != nil {
			panic(err)
		}
	}

	file, err := os.OpenFile(outputReadableFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := file.WriteString(resReadable); err != nil {
		panic(err)
	}

	// write output machine
	if _, err := os.Stat(outputMachineFile); err == nil {
		if err := os.Remove(outputMachineFile); err != nil {
			panic(err)
		}
	}

	file, err = os.OpenFile(outputMachineFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := file.WriteString(resMachine); err != nil {
		panic(err)
	}

	return len(resultCriticalMachine) + len(resultsWarningMachine)
}
