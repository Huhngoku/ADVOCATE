package io

import (
	"analyzer/trace"
	"io"
	"os"
	"strconv"
)

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

/*
 * Write the trace to a file
 * Args:
 *   path (string): The path to the file to write to
 *   numberRoutines (int): The number of routines in the trace
 */
func WriteTrace(path string, numberRoutines int) error {
	// delete folder if exists
	if _, err := os.Stat(path); err == nil {
		println(path + " already exists. Delete folder " + path)
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	// create new folder
	println("Create new folder " + path + "...")
	if err := os.Mkdir(path, 0755); err != nil {
		return err
	}

	// open file
	for i := 1; i <= numberRoutines; i++ {
		fileName := path + "trace_" + strconv.Itoa(i) + ".log"
		println("Create new file " + fileName + "...")
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		// write trace
		println("Write trace to " + path + "...")
		trace := trace.GetTraceFromId(i)
		for index, element := range trace {
			elementString := element.ToString()
			if _, err := file.WriteString(elementString); err != nil {
				panic(err)
			}
			if index < len(trace)-1 {
				if _, err := file.WriteString(";"); err != nil {
					panic(err)
				}
			}
		}
		if _, err := file.WriteString("\n"); err != nil {
			panic(err)
		}
	}
	println("Trace written")
	return nil
}
