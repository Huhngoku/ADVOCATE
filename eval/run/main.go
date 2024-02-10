package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

const (
	red = "\033[31m"
	grn = "\033[32m"
	end = "\033[0m"
)

var (
	home, _        = os.UserHomeDir()
	pathToAdvocate = home + "/Uni/HiWi/ADVOCATE"
	analyzerPath   = pathToAdvocate + "/analyzer"
	runEval        = pathToAdvocate + "/eval"
	overviewPath   = pathToAdvocate + "/eval/createOverview"

	// name
	// pathToExec
	// pathToProg
	// name of executable
	// params (if any)
	programs = [][]string{}

	resPath = runEval

	resSucc      = []string{}
	resFailure   = []string{}
	numberErrors = 0
)

func addConstructed() {
	for i := 1; i <= 45; i++ {
		programs = append(programs, []string{"constructed " + strconv.Itoa(i), pathToAdvocate + "/examples/constructed/", pathToAdvocate + "/examples/constructed/", "main", "-c", strconv.Itoa(i), "-t", "5"})
	}
}

// TODO: measure time
// TODO: run without advocate
// TODO: run withput trace write
func runExec(pathToExec string, execName string, execArgs []string) (error, string) {
	commandStr := "./" + execName
	cmd := exec.Command(commandStr, execArgs...)
	cmd.Dir = pathToExec

	log(EXEC, cmd.String()+" in "+pathToExec)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return err, stderr.String()
	}

	return nil, ""
}

// TODO: measure time
func runAnalyzer(progName string, pathToTrace string) (error, string) {
	cmdStr := "./analyzer"
	resOut := resPath + "/" + progName + "/"
	cmdArgs := []string{"-t", pathToTrace, "-x", "-p", "-r", resOut}
	cmd := exec.Command(cmdStr, cmdArgs...)
	cmd.Dir = analyzerPath

	log(ANALYZE, cmd.String()+" in "+analyzerPath)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return err, stderr.String()
	}

	return nil, ""
}

func createOverview(progName string, progPath string, execPath string) error {
	cmdStr := "./overview"
	resOut := resPath + "/" + progName
	cmdArgs := []string{"-n", "overview", "-p", progPath, "-s", resOut, "-t",
		execPath + "trace/", "-r", resOut + "/results_readable.log"} // TODO: add time file
	cmd := exec.Command(cmdStr, cmdArgs...)
	cmd.Dir = overviewPath

	log(OVERVIEW, cmd.String()+" in "+overviewPath)
	err := cmd.Run()
	return err
}

func createResFolder(progName string) error {
	progFolder := resPath + "/" + progName

	log(CREATE, progName+" at "+progFolder)

	err := os.MkdirAll(progFolder, 0755)
	if err != nil {
		return err
	}

	return nil
}

func finish() error {
	log(CLEANUP, "Create overview at "+resPath+"/overview.log")
	file, err := os.OpenFile(resPath+"/overview.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("Success\n\n")
	for _, suc := range resSucc {
		file.WriteString(suc + "\n")
	}

	file.WriteString("Failed\n\n")
	for _, fail := range resFailure {
		_, err = file.WriteString(fail + "\n")
	}
	return err
}

func setup() error {
	// create the res folder if needed
	name := resPath + "/results"
	counter := 1

	for {
		nameStr := name + "_" + strconv.Itoa(counter)
		if _, err := os.Stat(nameStr); os.IsNotExist(err) {
			errDir := os.MkdirAll(nameStr, 0755)
			if errDir != nil {
				return errDir
			}
			resPath = nameStr
			break
		}
		counter++
	}

	log(SETUP, "Create res folder "+resPath)

	addConstructed()
	return nil
}

func main() {
	setup()

	for _, program := range programs {
		log(START, program[0])

		name := program[0]
		pathToTrace := program[1] + "trace/"
		execPath := program[1]
		progPath := program[2]
		execName := program[3]
		execArgs := program[4:]

		err := createResFolder(name)
		if err != nil {
			log(ERROR, err.Error())
			log(FAILED, name)
			continue
		}

		err, out := runExec(execPath, execName, execArgs)
		if err != nil {
			log(ERROR, err.Error()+":\n"+out)
			// log(FAILED, name)
			// continue
		}

		err, out = runAnalyzer(name, pathToTrace)
		if err != nil {
			log(ERROR, err.Error()+":\n"+out)
			log(FAILED, name)
			continue
		}

		err = createOverview(name, progPath, execPath)
		if err != nil {
			log(ERROR, err.Error()+":\n"+out)
			log(FAILED, name)
			continue
		}

		log(DONE, name)
	}

	err := finish()
	if err != nil {
		log(ERROR, err.Error())
		log(FAILED, "Failed to create overfiew")
	} else {
		log(DONE, "Created Overview")
	}

	if numberErrors == 0 {
		log(FINISHS, "Finished with 0 errors")
	} else {
		log(FINISHF, "Finished with "+strconv.Itoa(numberErrors)+" errors")
	}
}

type logType int

const (
	START logType = iota
	CREATE
	EXEC
	ANALYZE
	FAILED
	DONE
	ERROR
	SETUP
	FINISHS
	FINISHF
	CLEANUP
	OVERVIEW
)

func log(lt logType, message string) {
	res := ""
	switch lt {
	case START:
		res = "\n[START   ] " + end + message
	case CREATE:
		res = "[CREATE  ] " + end + message
	case EXEC:
		res = "[EXEC    ] " + end + message
	case ANALYZE:
		res = "[ANALYZE ] " + end + message
	case DONE:
		res = grn + "[DONE    ] " + end + message
		resSucc = append(resSucc, message)
	case FAILED:
		res = red + "[FAILED  ] " + end + message
		resFailure = append(resFailure, message)
	case ERROR:
		res = red + "[ERROR   ] " + end + message
		numberErrors++
	case SETUP:
		res = "[SETUP   ] " + message
	case CLEANUP:
		res = "[CLEANUP ] " + message
	case FINISHS:
		res = "\n" + grn + "[FINISH  ] " + end + message
	case FINISHF:
		res = "\n" + red + "[FINISH  ] " + end + message
	case OVERVIEW:
		res = "[OVERVIEW] " + message
	default:
		res = "[UNKOWN  ] "
	}
	fmt.Println(res)
}
