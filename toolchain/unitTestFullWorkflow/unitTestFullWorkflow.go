package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	advocateRoot := flag.String("a", "", "Path to the advocate root")
	projectRoot := flag.String("f", "", "Path to the project root folder")
	moduleModeString := flag.String("m", "", "Module mode")
	testName := flag.String("t", "", "Name of the test")
	testFile := flag.String("tf", "", "Path to the test file")
	flag.Parse()
	moduleMode := false
	if *moduleModeString == "true" {
		moduleMode = true
	}
	//check if argument is missing else print usage
	if *advocateRoot == "" || *projectRoot == "" || *testName == "" || *testFile == "" {
		fmt.Println("Please provide all the arguments")
		fmt.Println("Usage: go run unitTestFullWorkflow -a <advocate folder> -f <project root folder> -t <test name> -tf <test file>")
		//print current arguments
		fmt.Println("Advocate Folder:", *advocateRoot)
		fmt.Println("Project Root Folder:", *projectRoot)
		fmt.Println("Test Name:", *testName)
		fmt.Println("Test File:", *testFile)
		//exit program with 1
		return
	}
	pathToAnalyzer := filepath.Join(*advocateRoot, "analyzer", "analyzer")
	pathToPatchedGoRuntime := filepath.Join(*advocateRoot, "go-patch", "bin", "go")
	pathToGoRoot := filepath.Join(*advocateRoot, "go-patch")
	pathToOverheadInserter := filepath.Join(*advocateRoot, "toolchain", "unitTestOverheadInserter", "unitTestOverheadInserter")
	pathToOverheadRemover := filepath.Join(*advocateRoot, "toolchain", "unitTestOverheadRemover", "unitTestOverheadRemover")
	//initial setup
	initialSetup(projectRoot, pathToGoRoot, testFile, testName)
	fmt.Println("Remove Overhead just in case")
	removeOverhead(pathToOverheadRemover, testFile, testName)
	fmt.Println("Add Overhead")
	addOverhead(pathToOverheadInserter, testFile, testName)
	fmt.Println("Run test")
	runTest(moduleMode, pathToPatchedGoRuntime, nil, testFile, testName, pathToOverheadRemover)

}
func initialSetup(projectRoot *string, pathToGoRoot string, testFile *string, testName *string) {
	err := os.Chdir(*projectRoot)
	if err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}
	fmt.Println("In directory:", *projectRoot)
	err = os.Setenv("GOROOT", pathToGoRoot)
	if err != nil {
		fmt.Println("Error setting GOROOT:", err)
		return
	}
	fmt.Println("GOROOT set to:", pathToGoRoot)
	advocateCommandLog, err := os.Create("advocateCommand.log")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	_, err = advocateCommandLog.WriteString(*testFile + "\n")
	_, err = advocateCommandLog.WriteString(*testName + "\n")
}
func addOverhead(pathToOverheadInserter string, testFile *string, testName *string) {
	//Add overhead: Usage $pathToOverheadInserter -f $file -t $testName
	fmt.Println("Add Overhead")
	cmd := exec.Command(pathToOverheadInserter, "-f", *testFile, "-t", *testName)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error adding overhead:", err)
		return
	}
}
func removeOverhead(pathToOverheadRemover string, testFile *string, testName *string) {
	//Remove overhead: Usage $pathToOverheadRemover -f $file -t $testName
	fmt.Println("Remove Overhead")
	cmd := exec.Command(pathToOverheadRemover, "-f", *testFile, "-t", *testName)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error removing overhead:", err)
		return
	}
}
func runTest(moduleMode bool, pathToPatchedGoRuntime string, advocateCommandLog *os.File, testFile *string, testName *string, pathToOverheadRemover string) {
	//check for module mode
	if moduleMode {
		// echo "$pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod ./$package" >>advocateCommand.log
		// $pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod "./$package"
		_, err := advocateCommandLog.WriteString(pathToPatchedGoRuntime + " test -count=1 -run=" + *testName + " -mod=mod ./" + *testFile + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		cmd := exec.Command(pathToPatchedGoRuntime, "test", "-count=1", "-run="+*testName, "-mod=mod", "./"+*testFile)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Remove Overhead")
			cmd = exec.Command(pathToOverheadRemover, "-f", *testFile, "-t", *testName)
			fmt.Println("Error in running therefor overhead removed and full workflow stopped", err)
			return
		}
	} else {
		//do the same without -mod=mod
		_, err := advocateCommandLog.WriteString(pathToPatchedGoRuntime + " test -count=1 -run=" + *testName + " ./" + *testFile + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		cmd := exec.Command(pathToPatchedGoRuntime, "test", "-count=1", "-run="+*testName, "./"+*testFile)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Remove Overhead")
			cmd = exec.Command(pathToOverheadRemover, "-f", *testFile, "-t", *testName)
			fmt.Println("Error in running therefor overhead removed and full workflow stopped", err)
			return
		}
	}

}
