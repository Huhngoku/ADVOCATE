package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	fileName := flag.String("f", "", "path to the file")
	flag.Parse()
	if *fileName == "" {
		fmt.Fprintln(os.Stderr, "Please provide a file name")
		fmt.Fprintln(os.Stderr, "Usage: preambleInserter -f <file>")
		os.Exit(1)
	}
	if _, err := os.Stat(*fileName); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File %s does not exist\n", *fileName)
		os.Exit(1)
	}
	file, err := os.Open(*fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	inImportBlock := false
	numberOfLinesToSkip := 0
	for scanner.Scan() {
		line := scanner.Text()
		if numberOfLinesToSkip > 0 {
			numberOfLinesToSkip--
			continue
		} else if strings.Contains(line, "// ======= Preamble Start =======") {
			numberOfLinesToSkip = 3
			continue
		} else if strings.Contains(line, "import (") {
			inImportBlock = true
			lines = append(lines, line)
		} else if inImportBlock && strings.Contains(line, ")") {
			inImportBlock = false
			lines = append(lines, line)
		} else if inImportBlock && strings.Contains(line, "\"advocate\"") {
			continue
		} else if strings.Contains(line, "import \"advocate\"") {
			continue
		} else {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	err = os.WriteFile(*fileName, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		panic(err)
	}
}
