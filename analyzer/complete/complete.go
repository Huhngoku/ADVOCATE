package complete

import (
	"fmt"
	"os"
)

/*
 * Check if all program elements are in trace
 * Args:
 * 	resultFolderPath: path to the folder containing the trace files
 * 	progPath: path to the program file
 * Returns:
 * 	error: error if any
 */
func Check(resultFolderPath string, progPath string) error {
	progElems, err := getProgramElements(progPath)
	if err != nil {
		println("Error in getProgramElements")
		return err
	}

	traceElems, err := getTraceElements(resultFolderPath)
	if err != nil {
		println("Error in getTraceElements")
		return err
	}

	notInTrace := areAllProgElemInTrace(progElems, traceElems)
	notSelectedSelectCase := getNotSelectedSelectCases()

	err = printResultsToFile(notInTrace, notSelectedSelectCase, resultFolderPath)

	return err
}

func areAllProgElemInTrace(progElems map[string][]int, traceElems map[string][]int) map[string][]int {
	res := map[string][]int{}

	for file, lines := range progElems {
		// file not recorded in trace
		if _, ok := traceElems[file]; !ok {
			if _, ok := res[file]; !ok {
				res[file] = make([]int, 0)
			}

			res[file] = append(res[file], -1) // -1 signaling, that no element in file was in trace
			res[file] = append(res[file], lines...)
		}

		for _, line := range lines {
			if !contains(traceElems[file], line) {
				if _, ok := res[file]; !ok {
					res[file] = make([]int, 0)
				}

				res[file] = append(res[file], line)
			}
		}
	}

	return res
}

/*
 * GetNotSelectedSelectCases prints the elements and select cases that were not executed
 * into a file.
 * Args:
 * 	elements: the elements that were not executed
 * 	selects: the select cases that were not selected
 */
func printResultsToFile(elements map[string][]int, selects map[string]map[int][]int,
	path string) error {
	// create file to write results for elements

	path = fmt.Sprintf("%s/AdvocateNotExecuted.txt", path)
	notExecutedFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer notExecutedFile.Close()

	if len(elements) == 0 && len(selects) == 0 {
		notExecutedFile.WriteString("All program elements were executed\n")
		return nil
	}

	// write elements that were not executed
	if len(elements) > 0 {
		notExecutedFile.WriteString("Program elements that were not executed:\n")
		for file, lines := range elements {
			notExecutedFile.WriteString(fmt.Sprintf("%s:[", file))
			for i, line := range lines {
				if line == -1 {
					notExecutedFile.WriteString("No element in file was executed")
					break
				} else {
					notExecutedFile.WriteString(fmt.Sprintf("%d", line))
					if i != len(lines)-1 {
						notExecutedFile.WriteString(",")
					}
				}
			}
			notExecutedFile.WriteString("]\n")
		}
	}

	if len(elements) > 0 && len(selects) > 0 {
		notExecutedFile.WriteString("\n")
	}

	// write select cases that were not selected
	if len(selects) > 0 {

		notExecutedFile.WriteString("Select cases that were not selected:\n")
		for file, lines := range selects {
			for line, cases := range lines {
				notExecutedFile.WriteString(fmt.Sprintf("%s:%d:[", file, line))
				for i, c := range cases {
					if c == -1 {
						notExecutedFile.WriteString("D")
					} else {
						notExecutedFile.WriteString(fmt.Sprintf("%d", c))
					}

					if i != len(cases)-1 {
						notExecutedFile.WriteString(",")
					}
				}
				notExecutedFile.WriteString("]\n")
			}
		}
	}

	return nil
}

func contains(arr []int, elem int) bool {
	for _, e := range arr {
		if e == elem {
			return true
		}
	}

	return false
}
