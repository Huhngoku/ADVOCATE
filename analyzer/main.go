package main

import (
	"analyzer/debug"
	"analyzer/reader"
	"analyzer/trace"
	"flag"
)

func main() {
	file_path := flag.String("f", "./dedego.log", "Path to the log file")
	level := flag.Int("d", 1, "Debug Level, 0 = no output, 1 = errors, 2 = info, 3 = debug")
	test := flag.Bool("t", false, "Test mode, only for testing")

	flag.Parse()
	debug.SetDebugLevel(*level)

	reader.CreateTraceFromFile(*file_path)

	if *test {
		res := trace.CheckTrace()
		if !res {
			println("Error END")
		}
	}

}
