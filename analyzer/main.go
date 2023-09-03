package main

import (
	"flag"
	"time"

	"analyzer/debug"
	"analyzer/trace"
)

func main() {
	start_time := time.Now()

	file_path := flag.String("f", "./dedego.log", "Path to the log file")
	level := flag.Int("d", 1, "Debug Level, 0 = no output, 1 = errors, 2 = info, 3 = debug")

	flag.Parse()
	debug.SetDebugLevel(*level)
	trace.CreateTraceFromFile(*file_path)

	runtime := time.Since(start_time)
	debug.Log("Finished analyzis. Total runtime: "+runtime.String(), 2)
}
