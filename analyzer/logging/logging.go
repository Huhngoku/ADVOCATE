package logging

import (
	"os"
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

var output_file *os.File
var output = false
var found_bug = false

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
	if output {
		output_file.WriteString(message + "\n")
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
 */
func InitLogging(level int, out string, noResult bool, noWarn bool) {
	if level < 0 {
		level = 0
	}
	levelDebug = level

	if out != "" {
		output_file, _ = os.Create(out)
		output = true
	}

	if noResult {
		levelResult = int(NONE)
	} else if noWarn {
		levelResult = int(CRITICAL)
	} else {
		levelResult = int(WARNING)
	}

}

func PrintNotFound() {
	if !found_bug && levelDebug != int(SILENT) && levelResult != int(NONE) {
		print("No bug found\n")
	}
}

/*
* Get the current runtime
 */
func GetRuntime() string {
	return time.Since(start_time).String()
}
