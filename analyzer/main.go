package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
	bugIndex := flag.Int("i", -1, "Index of the result to use for the reordered trace file. Only needed if -n is set. 1 based")
	ignoreCriticalSection := flag.Bool("c", false, "Ignore happens before relations of critical sections (default false)")
	noRewrite := flag.Bool("x", false, "Do not ask to create a reordered trace file after the analysis (default false)")
	noWarning := flag.Bool("w", false, "Do not print warnings (default false)")
	noPrint := flag.Bool("p", false, "Do not print the results to the terminal (default false). Automatically set -x to true")
	flag.Parse()

	if *pathTrace == "" {
		fmt.Println("Please provide a path to the trace file. Set with -t [file]")
		return
	}

	if *noPrint {
		*noRewrite = true
	}

	folder := filepath.Dir(*pathTrace) + string(os.PathSeparator)

	outMachine := folder + "/results_machine.log"
	outReadable := folder + "/results_readable.log"
	newTrace := folder + "rewritten_trace.log"

	// rewrite the trace file based on given analysis results. No analysis is run
	if *rewrite {
		if *bugIndex == -1 {
			fmt.Println("Please provide the index of the result to use for the reordered trace file. Set with -i [file]")
			return
		}
		numberOfRoutines := reader.CreateTraceFromFile(*pathTrace)
		rewriteTrace(*pathTrace, newTrace, *bugIndex, numberOfRoutines)
		return
	}

	// run the analysis and, if requested, create a reordered trace file
	// based on the analysis results

	logging.InitLogging(*level, outReadable, outMachine)
	numberOfRoutines := reader.CreateTraceFromFile(*pathTrace)
	trace.SetNumberOfRoutines(numberOfRoutines)
	trace.RunAnalysis(*fifo, *ignoreCriticalSection)

	numberOfResults := logging.PrintSummary(*noWarning, *noPrint)

	if numberOfResults != 0 && !*noRewrite {
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
