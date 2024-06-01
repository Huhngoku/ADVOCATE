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
	replayOverheadString := flag.String("r", "false", "replay overhead")
	replayNum := flag.String("n", "1", "replay number")
	flag.Parse()
	replayOverhead := false
	if *replayOverheadString == "true" {
		replayOverhead = true
	}
	if *testName == "" {
		fmt.Println("Please provide a test name")
		fmt.Println("Usage: go run unitTestOverheadInserter -f <file> -t <test name>")
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
		os.Exit(1)
		return
	}

	addOverhead(*fileName, *testName, replayOverhead, *replayNum)
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

func addOverhead(fileName string, testName string, replayOverhead bool, replayNumber string) {
	// print replay num for debugging
	fmt.Println("Replay number: ", replayNumber)
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
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
		//check for test method
		if strings.Contains(line, "func "+testName) {
			if replayOverhead {
				lines = append(lines, fmt.Sprintf(`	// ======= Preamble Start =======
	advocate.EnableReplay(%s, true)
	defer advocate.WaitForReplayFinish()
	// ======= Preamble End =======`, replayNumber))
			} else {
				lines = append(lines, `	// ======= Preamble Start =======
	advocate.InitTracing(0)
	defer advocate.Finish()
	// ======= Preamble End =======`)
			}
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
