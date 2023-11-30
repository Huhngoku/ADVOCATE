package io

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"rewrite/trace"
)

/*
 * Read and build the trace from a file
 * Args:
 *   filePath (string): The path to the log file
 *   bufferSize (int): The size of the buffer for the scanner
 */
func ReadTrace(filePath string) {
	println("Read trace from " + filePath + "...")
	mb := 1048576 // 1 MB
	maxTokenSize := 4

	for {
		file, err := os.Open(filePath)
		if err != nil {
			println("Error opening file: " + filePath)
			panic(err)
		}

		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 0, maxTokenSize*mb), maxTokenSize*mb)
		routine := 0
		for scanner.Scan() {
			routine++
			line := scanner.Text()
			processLine(line, routine)
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

	trace.Sort()

	println("Trace created")
}

/*
 * Process one line from the log file.
 * Args:
 *   line (string): The line to process
 *   routine (int): The routine id, equal to the line number
 */
func processLine(line string, routine int) {
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
	default:
		panic("Unknown element type in: " + element)
	}

	if err != nil {
		panic(err)
	}

}

/*
 * Copy a file from source to dest
 * Args:
 *   source (string): The path to the source file
 *   dest (string): The path to the destination file
 */
func CopyFile(source string, dest string) {
	println("Copy file from " + source + " to " + dest + "...")
	sourceFile, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		panic(err)
	}

	err = destFile.Sync()
	if err != nil {
		panic(err)
	}
}
