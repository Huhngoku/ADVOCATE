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
 *   index (int): The index of the result to create a trace for (0 based)
 * Returns:
 *   bool: true, if the bug was not a possible, but an actually occuring bug
 *   Bug: The bug that was selected
 *   error: An error if the bug could not be processed
 */
func ReadAnalysisResults(filePath string, index int) (bool, bugs.Bug, error) {
	println("Read analysis results from " + filePath + " for index " + strconv.Itoa(index) + "...")

	mb := 1048576 // 1 MB
	maxTokenSize := 1

	bugStr := ""

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
			bugStr = scanner.Text()
			if index == i {
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

	println("Analysis results read")

	actual, bug, err := bugs.ProcessBug(bugStr)
	if err != nil {
		println("Error processing bug")
		println(err.Error())
		return false, bug, err
	}

	bug.Println()

	if actual {
		println("The bug is an actual bug.")
		println("No rewrite needed.")
		return true, bug, nil
	}

	return false, bug, nil

}
