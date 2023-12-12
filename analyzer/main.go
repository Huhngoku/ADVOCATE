package main

import (
	"flag"
	"fmt"

	"analyzer/io"
	"analyzer/logging"
	"analyzer/reader"
	"analyzer/rewriter"
	"analyzer/trace"
)

func main() {
	pathTrace := flag.String("t", "", "Path to the trace file to analyze or rewrite")
	level := flag.Int("d", 1, "Debug Level, 0 = silent, 1 = errors, 2 = info, 3 = debug (default 1)")
	fifo := flag.Bool("f", false, "Assume a FIFO ordering for buffered channels (default false)")
	rewrite := flag.Bool("n", false, "Create a reordered trace file from a given analysis "+
		"result without running the analysis. -r and -i are required. If not set, a rewritten trace can be created from the current analysis results")
	pathResult := flag.String("r", "", "Path to the analysis result file. Only needed if -n is set")
	bugIndex := flag.Int("i", -1, "Index of the result to use for the reordered trace file. Only needed if -n is set. 1 based")
	flag.Parse()

	outMachine := "results_machine.log"
	outReadable := "results_readable.log"
	newTrace := "rewritten_trace.log"

	if *pathTrace == "" {
		fmt.Println("Please provide a path to the trace file. Set with -t [file]")
		return
	}

	// rewrite the trace file based on given analysis results. No analysis is run
	if *rewrite {
		if *pathResult == "" {
			fmt.Println("Please provide a path to analysis result file. Set with -r [file]")
			return
		}
		if *bugIndex == -1 {
			fmt.Println("Please provide the index of the result to use for the reordered trace file. Set with -i [file]")
			return
		}
		numberOfRoutines := reader.CreateTraceFromFile(*pathTrace)
		rewriteTrace(*pathResult, newTrace, *bugIndex, numberOfRoutines)
		return
	}

	// run the analysis and, if requested, create a reordered trace file
	// based on the analysis results
	if !*rewrite {
		logging.InitLogging(*level, outReadable, outMachine)
		numberOfRoutines := reader.CreateTraceFromFile(*pathTrace)
		trace.SetNumberOfRoutines(numberOfRoutines)
		trace.RunAnalysis(*fifo)

		numberOfResults := logging.PrintSummary()

		if numberOfResults != 0 {
			fmt.Println("\n\n\n")
			fmt.Print("Do you want to create a reordered trace file? (y/n): ")

			var answer string
			var createRewrittenFile bool
			resultIndex := -1

			fmt.Scanln(&answer)
			for answer != "y" && answer != "Y" && answer != "n" && answer != "N" {
				fmt.Print("Please enter y or n: ")
				fmt.Scanln(&answer)
			}

			if answer == "y" || answer == "Y" {
				createRewrittenFile = true
				for {
					fmt.Print("Enter the index of the result to use for the reordered trace file: ")
					_, err := fmt.Scanf("%d", &resultIndex)
					if err != nil {
						fmt.Print("Please enter a valid number: ")
						continue
					}
					if resultIndex < 1 || resultIndex > numberOfResults {
						fmt.Print("Please enter a valid number: ")
						continue
					}
					break
				}
			} else {
				createRewrittenFile = false
			}

			if createRewrittenFile {
				rewriteTrace(outMachine, newTrace, resultIndex, numberOfRoutines)
			}
		}
	}
}

/*
 * Rewrite the trace file based on given analysis results
 * Args:
 *   outMachine (string): The path to the analysis result file
 *   newTrace (string): The path where the new traces will be created
 *   resultIndex (int): The index of the result to use for the reordered trace file
 *   numberOfRoutines (int): The number of routines in the trace
 */
func rewriteTrace(outMachine string, newTrace string, resultIndex int,
	numberOfRoutines int) {
	actual, bug := io.ReadAnalysisResults(outMachine, resultIndex)
	if actual {
		// copy the file of the tracePath to the outputPath
		io.CopyFile(*&outMachine, newTrace)
		println("Trace created")
		return
	}

	rewriter.RewriteTrace(bug)

	io.WriteTrace(newTrace, numberOfRoutines)
}
