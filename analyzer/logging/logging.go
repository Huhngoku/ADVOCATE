package logging

import (
	"fmt"
	"os"
	"strconv"
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

type resultType string

const (
	Empty resultType = ""

	// actual
	ARecvOnClosed          resultType = "A1"
	ASendOnClosed          resultType = "A2"
	ACloseOnClosed         resultType = "A3"
	AConcurrentRecv        resultType = "A4"
	ASelCaseWithoutPartner resultType = "A5"

	// possible
	PSendOnClosed resultType = "P1"
	PRecvOnClosed resultType = "P2"
	PNegWG        resultType = "P3"

	// leaks
	LUnbufferedWith    = "L1"
	LUnbufferedWithout = "L2"
	LBufferedWith      = "L3"
	LBufferedWithout   = "L4"
	LMutex             = "L5"
	LWaitGroup         = "L6"
	LCond              = "L7"
)

var resultTypeMap = map[resultType]string{
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

type resultElem interface {
	isInvalid() bool
	stringMachine() string
	stringReadable() string
}

type traceElementResult struct {
	routineID int
	objID     int
	tPre      int
	objType   string
	file      string
	line      int
}

func (t traceElementResult) stringMachine() string {
	return fmt.Sprintf("T%d:%d:%d:%s:%s:%d", t.routineID, t.objID, t.tPre, t.objType, t.file, t.line)
}

func (t traceElementResult) stringReadable() string {
	return fmt.Sprintf("%s:%d", t.file, t.line)
}

func (t traceElementResult) isInvalid() bool {
	return t.objType == ""
}

type selectCaseResult struct {
	objID   int
	objType string
}

func (s selectCaseResult) stringMachine() string {
	return fmt.Sprintf("S%d:%s", s.objID, s.objType)
}

func (s selectCaseResult) stringReadable() string {
	return fmt.Sprintf("case: %d:%s", s.objID, s.objType)
}

func (s selectCaseResult) isInvalid() bool {
	return s.objType == ""
}

/*
 * Print a result message
 * Args:
 * 	level: level of the message
 *	message: message to print
 */
func Result(level resultLevel, resType resultType, argType1 string, argType2 string, arg1 []resultElem, arg2 []resultElem) {
	if arg1[0].isInvalid() || arg2[0].isInvalid() {
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

	if len(arg2) > 0 {
		resultReadable += "\n\t" + argType2
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
			resMachine += result + "\n"

			if !noPrint {
				fmt.Println(strconv.Itoa(counter) + " " + Red + result + Reset)
			}

			counter++
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
				fmt.Println(strconv.Itoa(counter) + " " + Orange + result + Reset)
			}

			counter++
		}
	}
	if !found {
		resReadable += "No bugs found" + "\n"

		if !noPrint {
			fmt.Println(Green + "No bugs found" + Reset)
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

	return len(resultsCriticalReadable) + len(resultsWarningReadable)
}
