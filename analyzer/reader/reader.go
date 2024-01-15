/*
Package reader provides functions for reading and processing log files.
*/
package reader

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"analyzer/logging"
	"analyzer/trace"
)

/*
 * Read and build the trace from a file
 * Args:
 *   filePath (string): The path to the log file
 *   bufferSize (int): The size of the buffer for the scanner
 * Returns:
 *   int: The number of routines in the trace
 */
func CreateTraceFromFile(filePath string) int {
	logging.Debug("Create trace from file "+filePath+"...", logging.INFO)
	mb := 1048576 // 1 MB
	maxTokenSize := 4
	numberOfRoutines := 0

	for {
		file, err := os.Open(filePath)
		if err != nil {
			logging.Debug("Error opening file: "+filePath, logging.ERROR)
			panic(err)
		}

		logging.Debug("Count number of routines...", logging.DEBUG)
		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 0, maxTokenSize*mb), maxTokenSize*mb)
		for scanner.Scan() {
			numberOfRoutines++
		}
		file.Close()
		if err := scanner.Err(); err != nil {
			if err == bufio.ErrTooLong {
				maxTokenSize *= 2 // max buffer was to short, restart
				println("Increase max file size to " + strconv.Itoa(maxTokenSize) + "MB")
				numberOfRoutines = 0
			} else {
				panic(err)
			}
		} else {
			break
		}
	}
	logging.Debug("Number of routines: "+strconv.Itoa(numberOfRoutines), logging.INFO)

	file2, err := os.Open(filePath)
	if err != nil {
		logging.Debug("Error opening file: "+filePath, logging.ERROR)
		panic(err)
	}

	logging.Debug("Create trace with "+strconv.Itoa(numberOfRoutines)+" routines...", logging.DEBUG)

	scanner := bufio.NewScanner(file2)
	scanner.Buffer(make([]byte, 0, maxTokenSize*mb), maxTokenSize*mb)
	routine := 0
	for scanner.Scan() {
		routine++
		line := scanner.Text()
		processLine(line, routine)
	}

	if err := scanner.Err(); err != nil {
		logging.Debug("Error reading file line.", logging.ERROR)
		if err.Error() != "token too long" {
			logging.Debug("Reader buffer size to small. Increase with -b.", logging.ERROR)
		}
		panic(err)
	}

	trace.Sort()

	logging.Debug("Trace created", logging.INFO)
	return numberOfRoutines
}

/*
 * Process one line from the log file.
 * Args:
 *   line (string): The line to process
 *   routine (int): The routine id, equal to the line number
 */
func processLine(line string, routine int) {
	logging.Debug("Read routine "+strconv.Itoa(routine), logging.DEBUG)
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
		logging.Debug("Routine "+strconv.Itoa(routine)+" is empty", logging.DEBUG)
		return
	}
	fields := strings.Split(element, ",")
	var err error
	switch fields[0] {
	case "A":
		err = trace.AddTraceElementAtomic(routine, fields[1], fields[2], fields[3])
	case "C":
		err = trace.AddTraceElementChannel(routine, fields[1], fields[2],
			fields[3], fields[4], fields[5], fields[6], fields[7], fields[8])
	case "M":
		err = trace.AddTraceElementMutex(routine, fields[1], fields[2],
			fields[3], fields[4], fields[5], fields[6], fields[7])
	case "G":
		err = trace.AddTraceElementFork(routine, fields[1], fields[2], fields[3])
	case "S":
		err = trace.AddTraceElementSelect(routine, fields[1], fields[2], fields[3],
			fields[4], fields[5], fields[6])
	case "W":
		err = trace.AddTraceElementWait(routine, fields[1], fields[2], fields[3],
			fields[4], fields[5], fields[6], fields[7])
	case "O":
		err = trace.AddTraceElementOnce(routine, fields[1], fields[2], fields[3],
			fields[4], fields[5])
	case "N":
		err = trace.AddTraceElementCond(routine, fields[1], fields[2], fields[3],
			fields[4], fields[5])
	default:
		panic("Unknown element type in: " + element)
	}

	if err != nil {
		panic(err)
	}

}
