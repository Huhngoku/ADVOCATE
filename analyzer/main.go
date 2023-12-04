package main

import (
	"flag"

	"analyzer/logging"
	"analyzer/reader"
	"analyzer/trace"
)

func main() {
	filePath := flag.String("l", "trace.log", "Path to the log file")
	level := flag.Int("d", 1, "Debug Level, 0 = silent, 1 = errors, 2 = info, 3 = debug (default 1)")
	result := flag.Bool("r", false, "Show the result immediately when found (default false)")
	noSummary := flag.Bool("s", false, "Do not show a summary at the end (default false)")
	fifo := flag.Bool("f", false, "Assume a FIFO ordering for buffered channels (default false)")
	outReadable := flag.String("o", "", "Print results to file in a structured form")
	outMachine := flag.String("m", "", "Print results to file, for reordering by machine")

	flag.Parse()
	logging.InitLogging(*level, *outReadable, *outMachine, *result, *noSummary)
	numberOfRoutines := reader.CreateTraceFromFile(*filePath)
	trace.SetNumberOfRoutines(numberOfRoutines)
	trace.RunAnalysis(*fifo)

	logging.PrintSummary()
}
