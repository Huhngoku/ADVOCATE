/*
Package reader provides functions for reading and processing log files.
*/
package reader

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"

	"analyzer/logging"
	"analyzer/trace"
)

func CreateTraceFromFiles(filePath string) (int, error) {
	maxTokenSize := 4
	numberIds := 0

	// traverse all files in the folder
	files, err := os.ReadDir(filePath)
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		routine, err := getRoutineFromFileName(file.Name())
		if err != nil {
			return 0, nil
		}
		numberIds = max(numberIds, routine)

		maxTokenSize, err = CreateTraceFromFile(filePath+"/"+file.Name(), routine, maxTokenSize)
		if err != nil {
			return 0, err
		}

	}

	trace.Sort()

	return numberIds, nil
}

/*
 * Read and build the trace from a file
 * Args:
 *   filePath (string): The path to the log file
 *   routine (int): The routine id
 *   maxTokenSize (int): The max token size
 * Returns:
 *   int: The max token size
 */
func CreateTraceFromFile(filePath string, routine int, maxTokenSize int) (int, error) {
	logging.Debug("Create trace from file "+filePath+"...", logging.INFO)
	mb := 1048576 // 1 MB

	for {
		file, err := os.Open(filePath)
		if err != nil {
			logging.Debug("Error opening file: "+filePath, logging.ERROR)
			return maxTokenSize, err
		}

		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 0, maxTokenSize*mb), maxTokenSize*mb)

		alreadyRead := false
		for scanner.Scan() {
			if alreadyRead {
				return maxTokenSize, errors.New("Log file contains more than one line")
			}

			line := scanner.Text()
			processLine(line, routine)
			alreadyRead = true
		}

		file.Close()

		if err := scanner.Err(); err != nil {
			if err == bufio.ErrTooLong {
				maxTokenSize *= 2 // max buffer was to short, restart
				println("Increase max file size to " + strconv.Itoa(maxTokenSize) + "MB")
			} else {
				return maxTokenSize, err
			}
		} else {
			break
		}
	}

	return maxTokenSize, nil
}

/*
 * Process one line from the log file.
 * Args:
 *   line (string): The line to process
 *   routine (int): The routine id, equal to the line number
 * Returns:
 *   error: An error if the line could not be processed
 */
func processLine(line string, routine int) error {
	logging.Debug("Read routine "+strconv.Itoa(routine), logging.DEBUG)
	elements := strings.Split(line, ";")
	for _, element := range elements {
		err := processElement(element, routine)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
 * Process one element from the log file.
 * Args:
 *   element (string): The element to process
 *   routine (int): The routine id, equal to the line number
 * Returns:
 *   error: An error if the element could not be processed
 */
func processElement(element string, routine int) error {
	if element == "" {
		logging.Debug("Routine "+strconv.Itoa(routine)+" is empty", logging.DEBUG)
		return errors.New("Element is empty")
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
		return errors.New("Unknown element type in: " + element)
	}

	if err != nil {
		return err
	}

	return nil
}

func getRoutineFromFileName(fileName string) (int, error) {
	// the file name is "trace_routineID.log"
	// remove the .log at the end
	fileName1 := strings.TrimSuffix(fileName, ".log")
	if fileName1 == fileName {
		return 0, errors.New("File name does not end with .log")
	}

	fileName2 := strings.TrimPrefix(fileName1, "trace_")
	if fileName2 == fileName1 {
		return 0, errors.New("File name does not start with trace_")
	}

	routine, err := strconv.Atoi(fileName2)
	if err != nil {
		return 0, err
	}

	return routine, nil
}
