package io

import (
	"analyzer/bugs"
	"bufio"
	"os"
	"strconv"
)

/*
 * Read the fail containing the output of the analysis
 * Extract the needed information to create a trace to replay the selected error
 * Args:
 *   filePath (string): The path to the file containing the analysis results
 *   index (int): The index of the result to create a trace for
 * Returns:
 *   bool: true, if the bug was not a possible, but an actually occuring bug
 *   Bug: The bug that was selected
 */
func ReadAnalysisResults(filePath string, index int) (bool, bugs.Bug) {
	println("Read analysis results from " + filePath + " for index " + strconv.Itoa(index) + "...")

	index = (index - 1) * 3

	mb := 1048576 // 1 MB
	maxTokenSize := 1

	errorType := ""
	argument1 := ""
	argument2 := ""

	for {
		file, err := os.Open(filePath)
		if err != nil {
			println("Error opening file: " + filePath)
			panic(err)
		}

		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 0, maxTokenSize*mb), maxTokenSize*mb)

		i := 0
		for scanner.Scan() {
			if index == i {
				errorType = scanner.Text()
			} else if index+1 == i {
				argument1 = scanner.Text()
			} else if index+2 == i {
				argument2 = scanner.Text()
				break
			}
			i++
		}

		if err := scanner.Err(); err != nil {
			if err == bufio.ErrTooLong {
				maxTokenSize *= 2 // max buffer was to short, restart
				println("Increase max file size to " + strconv.Itoa(maxTokenSize) + "MB")
			} else {
				println("Error reading file line.")
				panic(err)
			}
		} else {
			break
		}
	}

	print("Error type: " + errorType + "\n")
	print("Argument 1: " + argument1 + "\n")
	print("Argument 2: " + argument2 + "\n")

	println("Analysis results read.")

	actual, bug := bugs.ProcessBug(errorType, argument1, argument2)
	if actual {
		println("The bug is an actual bug.")
		println("No rewrite needed.")
		return true, bug
	} else {
		println("The bug is a possible bug.")
		println("Rewrite needed.")
		return false, bug
	}
}
