package main

import (
	"flag"

	"analyzer/debug"
	"analyzer/reader"
	"analyzer/trace"
)

func main() {
	file_path := flag.String("l", "./dedego.log", "Path to the log file")
	level := flag.Int("d", 1, "Debug Level, 0 = no output, 1 = errors, 2 = info, 3 = debug")
	buffer_size := flag.Int("b", 25, "Size of the buffer for the scanner in MB (default 25))")
	fifo := flag.Bool("f", false, "Assume a FIFO ordering for buffered channels (default false)")

	flag.Parse()
	debug.SetDebugLevel(*level)
	numberOfROutines := reader.CreateTraceFromFile(*file_path, *buffer_size)
	trace.SetNumberOfRoutines(numberOfROutines)
	trace.CalculateVectorClocks(*fifo)

	debug.Log("Finished analyzis.\nTotal runtime: "+debug.GetRuntime(), debug.INFO)
}
