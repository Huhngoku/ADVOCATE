package debug

var levelDebug int = 0

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var blue = "\033[34m"

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

func SetDebugLevel(level int) {
	if level < 0 {
		level = 0
	}
	levelDebug = level
}
