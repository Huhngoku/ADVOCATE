package cobufi

import (
	"os"
	"runtime"
)

/*
 * Write the trace of the program to a file.
 * The trace is written in the file named file_name.
 * The trace is written in the format of CoBufi.
 */
func CreateTrace(file_name string) {
	runtime.DisableTrace()

	os.Remove(file_name)
	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	numRout := runtime.GetNumberOfRoutines()
	for i := 1; i <= numRout; i++ {
		cobufiChan := make(chan string)
		go func() {
			runtime.TraceToStringByIdChannel(i, cobufiChan)
			close(cobufiChan)
		}()
		for trace := range cobufiChan {
			if _, err := file.WriteString(trace); err != nil {
				panic(err)
			}
		}
		if _, err := file.WriteString("\n"); err != nil {
			panic(err)
		}
	}
}
