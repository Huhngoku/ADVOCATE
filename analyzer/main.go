package main

import (
	"analyzer/reader"
	"analyzer/trace"
)

func main() {
	file_path := "./dedego.log"
	reader.CreateTraceFromFile(file_path)

	res := trace.CheckTrace()
	if !res {
		print("Error END")
	}

}
