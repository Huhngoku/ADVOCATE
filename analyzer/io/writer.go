package io

import (
	"analyzer/trace"
	"io"
	"os"
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
func WriteTrace(path string, numberRoutines int) {
	// delete file if exists
	if _, err := os.Stat(path); err == nil {
		println("File " + path + " already exists. Deleting...")
		if err := os.Remove(path); err != nil {
			panic(err)
		}
	}

	// open file
	println("Create new file " + path + "...")
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// write trace
	println("Write trace to " + path + "...")
	traces := trace.GetTraces()
	for i := 1; i <= numberRoutines; i++ {
		for index, element := range (*traces)[i] {
			elementString := element.ToString()
			if _, err := file.WriteString(elementString); err != nil {
				panic(err)
			}
			if index < len((*traces)[i])-1 {
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
}
