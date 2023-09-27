package reader

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"analyzer/debug"
	"analyzer/trace"
)

/*
 * Read and build the trace from a file
 * Args:
 *   file_path (string): The path to the log file
 * Returns:
 *   int: The number of routines in the trace
 */
func CreateTraceFromFile(file_path string) int {
	debug.Log("Create trace from file "+file_path+"...", 2)
	file, err := os.Open(file_path)
	if err != nil {
		debug.Log("Error opening file: "+file_path, 1)
		panic(err)
	}

	debug.Log("Count number of routines...", 3)
	numberOfRoutines := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		numberOfRoutines++
	}
	file.Close()

	file2, err := os.Open(file_path)
	if err != nil {
		debug.Log("Error opening file: "+file_path, 1)
		panic(err)
	}

	debug.Log("Create trace with "+strconv.Itoa(numberOfRoutines)+" routines...", 3)

	scanner = bufio.NewScanner(file2)
	routine := 0
	for scanner.Scan() {
		routine++
		line := scanner.Text()
		processLine(line, routine, numberOfRoutines)
	}

	if err := scanner.Err(); err != nil {
		debug.Log("Error reading file line", 1)
		panic(err)
	}

	debug.Log("Trace created", 2)
	return numberOfRoutines
}

/*
 * Process one line from the log file.
 * Args:
 *   line (string): The line to process
 *   routine (int): The routine id, equal to the line number
 *   numberOfRoutines (int): The number of routines in the log file
 */
func processLine(line string, routine int, numberOfRoutines int) {
	debug.Log("Read routine "+strconv.Itoa(routine), 3)
	elements := strings.Split(line, ";")
	for _, element := range elements {
		processElement(element, routine, numberOfRoutines)
	}
}

/*
 * Process one element from the log file.
 * Args:
 *   element (string): The element to process
 *   routine (int): The routine id, equal to the line number
 *   numberOfRoutines (int): The number of routines in the log file
 */
func processElement(element string, routine int, numberOfRoutines int) {
	if element == "" {
		debug.Log("Routine "+strconv.Itoa(routine)+" is empty", 3)
		return
	}
	debug.Log("Read element "+element, 3)
	fields := strings.Split(element, ",")
	var err error = nil
	switch fields[0] {
	case "A":
		err = trace.AddTraceElementAtomic(routine, numberOfRoutines, fields[1], fields[2], fields[3])
	case "C":
		err = trace.AddTraceElementChannel(routine, numberOfRoutines, fields[1], fields[2],
			fields[3], fields[4], fields[5], fields[6], fields[7], fields[8])
	case "M":
		err = trace.AddTraceElementMutex(routine, numberOfRoutines, fields[1], fields[2],
			fields[3], fields[4], fields[5], fields[6], fields[7])
	case "G":
		err = trace.AddTraceElementRoutine(routine, numberOfRoutines, fields[1], fields[2])
	case "S":
		err = trace.AddTraceElementSelect(routine, numberOfRoutines, fields[1], fields[2], fields[3],
			fields[4], fields[5])
	case "W":
		err = trace.AddTraceElementWait(routine, numberOfRoutines, fields[1], fields[2], fields[3],
			fields[4], fields[5], fields[6], fields[7])
	default:
		panic("Unknown element type in: " + element)
	}

	if err != nil {
		panic(err)
	}

}
