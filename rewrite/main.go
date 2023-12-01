package main

import (
	"flag"
	"rewrite/io"
	"rewrite/rewriter"
)

func main() {
	tracePath := flag.String("t", "", "Path to the trace file")
	resultPath := flag.String("e", "", "Path to the file containing the analysis results")
	resultIndex := flag.Int("i", 1, "Index of the result to create a trace for, 1 based")
	outputPath := flag.String("o", "new_trace.log", "Path to the file to write the trace to")
	flag.Parse()

	validInput := true

	if *tracePath == "" {
		println("No trace path specified. Set with -t <path>")
		validInput = false
	}

	if *resultPath == "" {
		println("No analysis result path specified. Set with -e <path>")
		validInput = false
	}

	if *resultIndex < 1 {
		println("Invalid result index. Set with -i <index>")
		validInput = false
	}

	if !validInput {
		println("Invalid input. Exiting...")
		return
	}

	actual, bug := io.ReadAnalysisResults(*resultPath, *resultIndex)
	if actual {
		// copy the file of the tracePath to the outputPath
		io.CopyFile(*tracePath, *outputPath)
		println("Trace created")
		return
	}

	numberRoutines := io.ReadTrace(*tracePath)

	rewriter.CreateNewTrace(bug)

	io.WriteTrace(*outputPath, numberRoutines)

}
