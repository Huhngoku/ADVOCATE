package main

import (
	"flag"

	"analyzer/debug"
	"analyzer/trace"
)

func main() {
	file_path := flag.String("f", "./dedego.log", "Path to the log file")
	level := flag.Int("d", 1, "Debug Level, 0 = no output, 1 = errors, 2 = info, 3 = debug")

	flag.Parse()
	debug.SetDebugLevel(*level)
	trace.CreateTraceFromFile(*file_path)

	debug.Log("Finished analyzis.\nTotal runtime: "+debug.GetRuntime(), 2)
}
