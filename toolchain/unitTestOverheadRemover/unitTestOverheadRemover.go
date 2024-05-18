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
	testName := flag.String("t", "", "name of the test")
	flag.Parse()
	if *fileName == "" {
		fmt.Println("Please provide a file name")
		fmt.Println("Usage: go run unitTestOverheadRemover -f <file> -t <test name>")
		return
	}
	if *testName == "" {
		fmt.Println("Please provide a test name")
		fmt.Println("Usage: go run unitTestOverheadRemover -f <file> -t <test name>")
		return
	}
	if _, err := os.Stat(*fileName); os.IsNotExist(err) {
		fmt.Printf("File %s does not exist\n", *fileName)
		return
	}
	testExists, err := testExists(*testName, *fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if !testExists {
		fmt.Println("Test Method not found in file")
		return
	}

	removeOverhead(*fileName, *testName)
}

func testExists(testName string, fileName string) (bool, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return false, err
	}
	defer file.Close()
	regexStr := "func " + testName + "\\(*t \\*testing.T*\\) {"
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

func removeOverhead(fileName string, testName string) {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	inPreamble := false
	inImports := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "// ======= Preamble Start =======") {
			inPreamble = true
			continue
		}
		if strings.Contains(line, "// ======= Preamble End =======") {
			inPreamble = false
			continue
		}
		if inPreamble {
			continue
		}
		if strings.Contains(line, "import (") {
			inImports = true
		}
		if inImports && strings.Contains(line, "\"advocate\"") {
			continue
		}
		if strings.Contains(line, ")") {
			inImports = false
		}
		lines = append(lines, line)
	}

	file.Truncate(0)
	file.Seek(0, 0)
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()
}
