package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
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
	exists, err := mainMethodExists(*fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if !exists {
		fmt.Fprintln(os.Stderr, "Main Method not found in file")
		os.Exit(1)
	}

	addOverhead(*fileName)
}
func mainMethodExists(fileName string) (bool, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return false, err
	}
	defer file.Close()
	regexStr := "func main\\(\\) {"
	regex, err := regexp.Compile(regexStr)
	if err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if regex.MatchString(line) {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
func addOverhead(fileName string) {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		if strings.Contains(line, "import \"") {
			lines = append(lines, "import \"advocate\"")
		} else if strings.Contains(line, "import (") {
			lines = append(lines, "\t\"advocate\"")
		}

		if strings.Contains(line, "func main() {") {
			lines = append(lines, `	// ======= Preamble Start =======
	advocate.InitTracing(0)
	defer advocate.Finish()
	// ======= Preamble End =======`)
		}
	}

	file.Truncate(0)
	file.Seek(0, 0)
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()
}
