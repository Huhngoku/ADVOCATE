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
	result := flag.Bool("r", false, "Show the result immediately when found (default false)")
	noSummary := flag.Bool("s", false, "Do not show a summary at the end (default false)")
	buffer_size := flag.Int("b", 25, "Size of the buffer for the scanner in MB (default 25))")
	fifo := flag.Bool("f", false, "Assume a FIFO ordering for buffered channels (default false)")
	out := flag.String("o", "", "Print results to file")

	flag.Parse()
	logging.InitLogging(*level, *out, *result, *noSummary)
	numberOfRoutines := reader.CreateTraceFromFile(*file_path, *buffer_size)
	trace.SetNumberOfRoutines(numberOfRoutines)
	trace.RunAnalysis(*fifo)

	logging.PrintSummary()
}
