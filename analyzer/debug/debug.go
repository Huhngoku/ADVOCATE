package debug

import "time"

var levelDebug int = 0
var start_time = time.Now()

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var blue = "\033[34m"

/*
* Print a log message if the log level is sufficiant
* Args:
*   message: message to print
*   level: level of the message
 */
func Log(message string, level int) {
	if level <= levelDebug {
		if level == 1 {
			println(red + message + reset)
		} else if level == 2 {
			println(green + message + reset)
		} else if level == 3 {
			println(blue + message + reset)
		} else {
			println(message)
		}
	}
}

/*
* Initialize the debug
 */
func SetDebugLevel(level int) {
	if level < 0 {
		level = 0
	}
	levelDebug = level
}

/*
* Get the current runtime
 */
func GetRuntime() string {
	return time.Since(start_time).String()
}
