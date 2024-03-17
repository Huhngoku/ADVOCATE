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
	pathTrace := flag.String("t", "", "Path to the trace folder to analyze or rewrite")
	level := flag.Int("d", 1, "Debug Level, 0 = silent, 1 = errors, 2 = info, 3 = debug (default 1)")
	fifo := flag.Bool("f", false, "Assume a FIFO ordering for buffered channels (default false)")
	rewrite := flag.Bool("n", false, "Create a reordered trace file from a given analysis "+
		"result without running the analysis. -r and -i are required. If not set, a rewritten trace can be created from the current analysis results")
	bugIndex := flag.Int("i", -1, "Index of the result to use for the reordered trace file. Only needed if -n is set. 1 based")
	ignoreCriticalSection := flag.Bool("c", false, "Ignore happens before relations of critical sections (default false)")
	noRewrite := flag.Bool("x", false, "Do not ask to create a reordered trace file after the analysis (default false)")
	noWarning := flag.Bool("w", false, "Do not print warnings (default false)")
	noPrint := flag.Bool("p", false, "Do not print the results to the terminal (default false). Automatically set -x to true")
	resultFolder := flag.String("r", "", "Path to where the result file should be saved. If not set, it is saved in the trace folder")

	scenarios := flag.String("s", "", "Select which analysis scenario to run, e.g. -s srd for the option s, r and d. Options:\n"+
		"\ts: Send on closed channel\n"+
		"\tr: Receive on closed channel\n"+
		"\tw: Done before add on waitGroup\n"+
		"\tc: Close of closed channel\n"+
		"\tb: Concurrent receive on channel\n"+
		"\tl: Leaking routine\n"+
		"\tu: Select case without partner (not implemented yet)\n"+
		"\tc: Cyclic deadlock\n",
	// "\tm: Mixed deadlock\n"
	)

	flag.Parse()

	if *pathTrace == "" {
		fmt.Println("Please provide a path to the trace file. Set with -t [file]")
		return
	}

	if *noPrint {
		*noRewrite = true
	}

	folder := filepath.Dir(*pathTrace) + string(os.PathSeparator)
	if *resultFolder != "" {
		folder = *resultFolder
		if folder[len(folder)-1] != os.PathSeparator {
			folder += string(os.PathSeparator)
		}
	}

	analysisCases, err := parseAnalysisCases(*scenarios)
	if err != nil {
		panic(err)
	}

	outMachine := folder + "results_machine.log"
	outReadable := folder + "results_readable.log"
	newTrace := folder + "rewritten_trace/"

	// rewrite the trace file based on given analysis results. No analysis is run
	if *rewrite {
		if *bugIndex == -1 {
			fmt.Println("Please provide the index of the result to use for the reordered trace file. Set with -i [file]")
			return
		}
		numberOfRoutines, err := reader.CreateTraceFromFiles(*pathTrace)
		if err != nil {
			panic(err)
		}

		if err := rewriteTrace(*pathTrace, newTrace, *bugIndex, numberOfRoutines); err != nil {
			panic(err)
		}
		return
	}

	// run the analysis and, if requested, create a reordered trace file
	// based on the analysis results

	logging.InitLogging(*level, outReadable, outMachine)
	numberOfRoutines, err := reader.CreateTraceFromFiles(*pathTrace)
	if err != nil {
		panic(err)
	}
	trace.SetNumberOfRoutines(numberOfRoutines)
	trace.RunAnalysis(*fifo, *ignoreCriticalSection, analysisCases)

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
			if err := rewriteTrace(outMachine, newTrace, resultIndex, numberOfRoutines); err != nil {
				panic(err)
			}
		}
	}
}

/*
 * Rewrite the trace file based on given analysis results
 * Args:
 *   outMachine (string): The path to the analysis result file
 *   newTrace (string): The path where the new traces folder will be created
 *   resultIndex (int): The index of the result to use for the reordered trace file
 *   numberOfRoutines (int): The number of routines in the trace
 * Returns:
 *   error: An error if the trace file could not be created
 */
func rewriteTrace(outMachine string, newTrace string, resultIndex int,
	numberOfRoutines int) error {
	actual, bug, err := io.ReadAnalysisResults(outMachine, resultIndex)
	if err != nil {
		return err
	}

	if actual {
		// copy the file of the tracePath to the outputPath
		io.CopyFile(*&outMachine, newTrace)
		println("Trace created")
		return nil
	}

	err = rewriter.RewriteTrace(bug)
	if err != nil {
		return err
	}

	err = io.WriteTrace(newTrace, numberOfRoutines)
	if err != nil {
		return err
	}

	return nil
}

/*
 * Parse the given analysis cases
 * Args:
 *   cases (string): The string of analysis cases to parse
 * Returns:
 *   map[string]bool: A map of the analysis cases and if they are set
 *   error: An error if the cases could not be parsed
 */
func parseAnalysisCases(cases string) (map[string]bool, error) {
	analysisCases := map[string]bool{
		"sendOnClosed":         false,
		"receiveOnClosed":      false,
		"doneBeforeAdd":        false,
		"closeOnClosed":        false,
		"concurrentReceive":    false,
		"leak":                 false,
		"selectWithoutPartner": false,
		"cyclicDeadlock":       false,
		"mixedDeadlock":        false,
	}

	if cases == "" {
		analysisCases["sendOnClosed"] = true
		analysisCases["receiveOnClosed"] = true
		analysisCases["doneBeforeAdd"] = true
		analysisCases["closeOnClosed"] = true
		analysisCases["concurrentReceive"] = true
		analysisCases["leak"] = true
		analysisCases["selectWithoutPartner"] = true
		analysisCases["cyclicDeadlock"] = true
		// analysisCases["mixedDeadlock"] = true
	}

	for _, c := range cases {
		switch c {
		case 's':
			analysisCases["sendOnClosed"] = true
		case 'r':
			analysisCases["receiveOnClosed"] = true
		case 'w':
			analysisCases["doneBeforeAdd"] = true
		case 'c':
			analysisCases["closeOnClosed"] = true
		case 'b':
			analysisCases["concurrentReceive"] = true
		case 'l':
			analysisCases["leak"] = true
		case 'u':
			analysisCases["selectWithoutPartner"] = true
		case 'd':
			analysisCases["cyclicDeadlock"] = true
		// case 'm':
		// analysisCases["mixedDeadlock"] = true
		default:
			return nil, fmt.Errorf("Invalid analysis case: %c", c)
		}
	}
	return analysisCases, nil
}
