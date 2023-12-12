package logging

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var levelDebug int = 0

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
var outputMachineFile string
var foundBug = false
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
* Returns:
*   int: number of bugs found
 */
func PrintSummary() int {
	fmt.Println("Print Summary")
	counter := 1
	resMachine := ""
	resReadable := "==================== Summary ====================\n\n"
	fmt.Print("==================== Summary ====================\n\n")
	found := false
	if len(resultCritical) > 0 {
		found = true
		resReadable += "-------------------- Critical -------------------\n"
		fmt.Print("-------------------- Critical -------------------\n")
		for _, result := range resultCritical {
			resReadable += strconv.Itoa(counter) + " " + result + "\n"
			resMachine += result + "\n"
			fmt.Println(strconv.Itoa(counter) + " " + red + result + reset)
			counter++
		}
	}
	if len(resultsWarning) > 0 {
		found = true
		resReadable += "-------------------- Warning --------------------\n"
		fmt.Print("-------------------- Warning --------------------\n")
		for _, result := range resultsWarning {
			resReadable += strconv.Itoa(counter) + " " + result + "\n"
			resMachine += result + "\n"
			fmt.Println(strconv.Itoa(counter) + " " + orange + result + reset)
			counter++
		}
	}
	if !found {
		resReadable += "No bugs found" + "\n"
		fmt.Println(green + "No bugs found" + reset)
	}

	// write output readable
	file, err := os.OpenFile(outputReadableFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := file.WriteString(resReadable); err != nil {
		panic(err)
	}

	// write output machine
	resReadable = strings.ReplaceAll(resReadable, "\n\n", "\n")
	file, err = os.OpenFile(outputMachineFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := file.WriteString(resMachine); err != nil {
		panic(err)
	}

	return len(resultCritical) + len(resultsWarning)
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
