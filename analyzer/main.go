package main

import (
	"flag"

	"analyzer/logging"
	"analyzer/reader"
	"analyzer/trace"
)

func main() {
	file_path := flag.String("l", "trace.log", "Path to the log file")
	level := flag.Int("d", 1, "Debug Level, 0 = silent, 1 = errors, 2 = info, 3 = debug (default 1)")
	noResult := flag.Bool("r", false, "Do not show any results (default false)")
	noWarning := flag.Bool("w", false, "Do not show warnings (default false)")
	buffer_size := flag.Int("b", 25, "Size of the buffer for the scanner in MB (default 25))")
	fifo := flag.Bool("f", false, "Assume a FIFO ordering for buffered channels (default false)")
	out := flag.String("o", "", "Print results to file")

	flag.Parse()
	logging.InitLogging(*level, *out, *noResult, *noWarning)
	numberOfROutines := reader.CreateTraceFromFile(*file_path, *buffer_size)
	trace.SetNumberOfRoutines(numberOfROutines)
	trace.RunAnalysis(*fifo)

	logging.PrintNotFound() // print message, if no bug was found
	logging.Debug("Finished analyzis.\nTotal runtime: "+logging.GetRuntime(), logging.INFO)
}
