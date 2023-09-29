package main

import (
	"flag"

	"analyzer/logging"
	"analyzer/reader"
	"analyzer/trace"
)

func main() {
	file_path := flag.String("l", "./trace.log", "Path to the log file")
	level := flag.Int("d", 2, "Debug Level, 0 = silent, 1 = results, 2 = errors, 3 = info, 4 = debug (default 2)")
	buffer_size := flag.Int("b", 25, "Size of the buffer for the scanner in MB (default 25))")
	fifo := flag.Bool("f", false, "Assume a FIFO ordering for buffered channels (default false)")

	flag.Parse()
	logging.SetDebugLevel(*level)
	numberOfROutines := reader.CreateTraceFromFile(*file_path, *buffer_size)
	trace.SetNumberOfRoutines(numberOfROutines)
	trace.CalculateVectorClocks(*fifo)

	logging.Log("Finished analyzis.\nTotal runtime: "+logging.GetRuntime(), logging.INFO)
}
