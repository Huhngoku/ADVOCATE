package logging

import (
	"os"
	"strings"
)

var levelDebug int = 0
var levelResult int = 0

var reset = "\033[0m"
var red = "\033[31m"
var orange = "\033[33m"
var green = "\033[32m"
var blue = "\033[34m"

type debugLevel int

const (
	SILENT debugLevel = iota
	ERROR
	INFO
	DEBUG
)

type resultLevel int

const (
	NONE = iota
	CRITICAL
	WARNING
)

var outputReadableFile string
var outputReadable = false
var outputMachineFile string
var outputMachine = false
var foundBug = false
var no_sum = false
var resultsWarning []string
var resultCritical []string

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
			println(blue + message + reset)
		} else if level == INFO {
			println(green + message + reset)
		} else {
			println(message)
		}
	}
}

/*
Print a result message
Args:

	message: message to print
	level: level of the message
*/
func Result(message string, level resultLevel) {
	foundBug = true
	if level == WARNING {
		if !contains(resultsWarning, message) {
			resultsWarning = append(resultsWarning, message)
		}
	} else if level == CRITICAL {
		if !contains(resultCritical, message) {
			resultCritical = append(resultCritical, message)
		}
	}
	if int(level) <= levelResult {
		if level == CRITICAL {
			println(red + message + reset)
		} else if level == WARNING {
			println(orange + message + reset)
		} else {
			println(message)
		}
	}
}

/*
* Initialize the debug
* Args:
*   level: level of the debug
*   outReadable: path to the output file, no output file if empty
*   outMachine: path to the output file for the reordered trace, no output file if empty
*   noResult: true if no result should be printed
*   noWarn: true if no warnings should be printed
*   NoSum: true if no summary should be printed
 */
func InitLogging(level int, outReadable string, outMachine string, result bool, noSum bool) {
	if level < 0 {
		level = 0
	}
	levelDebug = level

	if outReadable != "" {
		outputReadableFile = outReadable
		outputReadable = true
	}

	if outMachine != "" {
		outputMachineFile = outMachine
		outputMachine = true
	}

	if result {
		levelResult = int(WARNING)
	} else {
		levelResult = int(NONE)
	}

	no_sum = noSum
}

/*
* Disable the output to a file
 */
func DisableOutput() {
	outputReadable = false
}

func PrintSummary() {
	res := "==================== Summary ====================\n\n"
	found := false
	if len(resultCritical) > 0 {
		found = true
		res += "-------------------- Critical -------------------\n"
		for _, result := range resultCritical {
			res += red + result + reset + "\n"
		}
	}
	if len(resultsWarning) > 0 {
		found = true
		res += "-------------------- Warning --------------------\n"
		for _, result := range resultsWarning {
			res += orange + result + reset + "\n"
		}
	}
	if !found {
		res += red + "No bugs found" + reset + "\n"
	}

	if !no_sum {
		print("\n\n")
		println(res)
	}

	res = strings.ReplaceAll(res, red, "")
	res = strings.ReplaceAll(res, orange, "")
	res = strings.ReplaceAll(res, reset, "")

	if outputReadable {
		file, err := os.OpenFile(outputReadableFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		if _, err := file.WriteString(res); err != nil {
			panic(err)
		}
	}

	// remove a given line from res
	res = strings.ReplaceAll(res, "-------------------- Critical -------------------\n", "")
	res = strings.ReplaceAll(res, "-------------------- Warning --------------------\n", "")
	res = strings.ReplaceAll(res, "==================== Summary ====================\n\n", "")
	res = strings.ReplaceAll(res, "\n\n", "\n")

	if outputMachine {
		file, err := os.OpenFile(outputMachineFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		if _, err := file.WriteString(res); err != nil {
			panic(err)
		}
	}
}

/*
* Check if a slice contains an element
* Args:
*   s: slice to check
*   e: element to check
 */
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
