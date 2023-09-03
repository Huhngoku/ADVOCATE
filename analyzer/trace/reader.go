package trace

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"analyzer/debug"
)

/*
 * Read and build the trace from a file
 */
func CreateTraceFromFile(file_path string) {
	debug.Log("Create trace from file "+file_path+"...", 2)
	file, err := os.Open(file_path)
	if err != nil {
		debug.Log("Error opening file: "+file_path, 1)
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	routine := 0
	for scanner.Scan() {
		routine++
		line := scanner.Text()
		processLine(line, routine)
	}

	if err := scanner.Err(); err != nil {
		debug.Log("Error reading file line", 1)
		panic(err)
	}

	Sort()        // sort the trace by tpre
	FindPartner() // set all partner
	debug.Log("Trace created", 2)
}

/*
 * Process one line from the log file.
 * Args:
 *   line (string): The line to process
 *   routine (int): The routine id, equal to the line number
 */
func processLine(line string, routine int) {
	debug.Log("Read routine "+strconv.Itoa(routine), 3)
	elements := strings.Split(line, ";")
	for _, element := range elements {
		processElement(element, routine)
	}
}

/*
 * Process one element from the log file.
 * Args:
 *   element (string): The element to process
 *   routine (int): The routine id, equal to the line number
 */
func processElement(element string, routine int) {
	if element == "" {
		debug.Log("Routine "+strconv.Itoa(routine)+" is empty", 3)
		return
	}
	debug.Log("Read element "+element, 3)
	fields := strings.Split(element, ",")
	var err error = nil
	switch fields[0] {
	case "A":
		err = addTraceElementAtomic(routine, fields[1], fields[2], fields[3])
	case "C":
		err = addTraceElementChannel(routine, fields[1], fields[2],
			fields[3], fields[4], fields[5], fields[6], fields[7], fields[8])
	case "M":
		err = addTraceElementMutex(routine, fields[1], fields[2],
			fields[3], fields[4], fields[5], fields[6], fields[7])
	case "G":
		err = addTraceElementRoutine(routine, fields[1], fields[2])
	case "S":
		err = addTraceElementSelect(routine, fields[1], fields[2], fields[3],
			fields[4], fields[5])
	case "W":
		err = addTraceElementWait(routine, fields[1], fields[2], fields[3],
			fields[4], fields[5], fields[6], fields[7])
	default:
		panic("Unknown element type in: " + element)
	}

	if err != nil {
		panic(err)
	}

}
