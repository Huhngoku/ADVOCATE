package logging

import (
	"os"
	"time"
)

var levelDebug int = 0
var start_time = time.Now()

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var blue = "\033[34m"

type level int

const (
	SILENT level = 0
	RESULT level = 1
	ERROR  level = 2
	INFO   level = 3
	DEBUG  level = 4
)

var output_file *os.File
var output = false
var found_bug = false

/*
* Print a log message if the log level is sufficiant
* Args:
*   message: message to print
*   level: level of the message
 */
func Log(message string, level level) {
	// print result to file
	if output && level == RESULT {
		if level == RESULT {
			if output {
				_, err2 := output_file.WriteString(message + "\n")

				if err2 != nil {
					Log(err2.Error(), ERROR)
				}
			}
		}
	}

	// print message to terminal
	if int(level) <= levelDebug {
		if level == RESULT {
			found_bug = true
			println(message)
		} else if level == ERROR {
			println(red + message + reset)
		} else if level == INFO {
			println(green + message + reset)
		} else if level == DEBUG {
			println(blue + message + reset)
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
 */
func InitLogging(level int, out string) {
	if level < 0 {
		level = 0
	}
	levelDebug = level

	if out != "" {
		output_file, _ = os.Create(out)
		output = true
	}

}

func PrintNotFound() {
	if !found_bug {
		Log("No problems found.", RESULT)
	}
}

/*
* Get the current runtime
 */
func GetRuntime() string {
	return time.Since(start_time).String()
}
