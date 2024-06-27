package complete

import "fmt"

/*
 * Check if all program elements are in trace
 * Args:
 * 	resultFolderPath: path to the folder containing the trace files
 * 	progPath: path to the program file
 * Returns:
 * 	error: error if any
 */
func Check(resultFolderPath string, progPath string) error {

	println(resultFolderPath, progPath)
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

	// for elem := range progElems {
	// 	fmt.Println(elem)
	// }
	// print("\n\n\n\n")
	for elem, lines := range traceElems {
		fmt.Println(elem, lines)
	}
	print("\n\n\n\n")

	println("Program elements: ", len(progElems))
	println("Trace elements: ", len(traceElems))
	res := areAllProgElemInTrace(progElems, traceElems)

	for file, lines := range res {
		if len(lines) != 0 {
			fmt.Printf("File %s: lines %v not in trace\n", file, lines)
		}
	}

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

func contains(arr []int, elem int) bool {
	for _, e := range arr {
		if e == elem {
			return true
		}
	}

	return false
}
