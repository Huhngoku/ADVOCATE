package logging

import (
	"os"
	"strings"
	"time"
)

var levelDebug int = 0
var levelResult int = 0
var start_time = time.Now()

var reset = "\033[0m"
var red = "\033[31m"
var orange = "\033[33m"
var green = "\033[32m"
var blue = "\033[34m"

type debug_level int

const (
	SILENT debug_level = iota
	ERROR
	INFO
	DEBUG
)

type result_level int

const (
	NONE = iota
	CRITICAL
	WARNING
)

var output_file string
var output = false
var found_bug = false
var no_sum = false
var resultsWarning []string
var resultCritical []string

/*
* Print a debug log message if the log level is sufficiant
* Args:
*   message: message to print
*   level: level of the message
 */
func Debug(message string, level debug_level) {
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
func Result(message string, level result_level) {
	found_bug = true
	if level == WARNING {
		resultsWarning = append(resultsWarning, message)
	} else if level == CRITICAL {
		resultCritical = append(resultCritical, message)
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
*   out: path to the output file, no output file if empty
*   noResult: true if no result should be printed
*   noWarn: true if no warnings should be printed
*   NoSum: true if no summary should be printed
 */
func InitLogging(level int, out string, result bool, noSum bool) {
	if level < 0 {
		level = 0
	}
	levelDebug = level

	if out != "" {
		output_file = out
		output = true
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
	output = false
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
	res += "\n=================================================\n"
	res += "Total runtime: " + GetRuntime() + "\n"
	res += "=================================================\n"

	if !no_sum {
		print("\n\n")
		println(res)
	}

	if output {
		res = strings.ReplaceAll(res, red, "")
		res = strings.ReplaceAll(res, orange, "")
		res = strings.ReplaceAll(res, reset, "")

		file, err := os.OpenFile(output_file, os.O_CREATE|os.O_WRONLY, 0644)
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
* Get the current runtime
 */
func GetRuntime() string {
	return time.Since(start_time).String()
}
