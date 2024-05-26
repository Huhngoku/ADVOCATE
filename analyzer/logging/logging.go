package logging

import (
	"analyzer/utils"
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
			println(Blue + message + Reset)
		} else if level == INFO {
			println(Green + message + Reset)
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
	// check if the message is valis
	resSplit := strings.Split(message, "\n")
	if len(resSplit) < 2 {
		return
	}
	resSplitSplit := strings.Split(resSplit[1], ":")
	if len(resSplitSplit) < 2 || resSplitSplit[1] == "" || resSplitSplit[1] == " " || resSplitSplit[1] == "\n" {
		return
	}

	foundBug = true
	if level == WARNING {
		if !utils.Contains(resultsWarning, message) {
			resultsWarning = append(resultsWarning, message)
		}
	} else if level == CRITICAL {
		if !utils.Contains(resultCritical, message) {
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

	if len(resultCritical) > 0 {
		found = true
		resReadable += "-------------------- Critical -------------------\n\n"

		if !noPrint {
			fmt.Print("-------------------- Critical -------------------\n\n")
		}

		for _, result := range resultCritical {
			resReadable += strconv.Itoa(counter) + " " + result + "\n"
			resMachine += result + "\n"

			if !noPrint {
				fmt.Println(strconv.Itoa(counter) + " " + Red + result + Reset)
			}

			counter++
		}
	}
	if len(resultsWarning) > 0 && !noWarning {
		found = true
		resReadable += "\n-------------------- Warning --------------------\n\n"
		if !noPrint {
			fmt.Print("\n-------------------- Warning --------------------\n\n")
		}

		for _, result := range resultsWarning {
			resReadable += strconv.Itoa(counter) + " " + result + "\n"
			resMachine += result + "\n"

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

	resMachine = strings.ReplaceAll(resMachine, "\n\t\t", ";")
	resMachine = strings.ReplaceAll(resMachine, ": ;", ": ")
	resMachine = strings.ReplaceAll(resMachine, "\n\n", "\n")
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
