package instrumenter

/*
Copyright (c) 2023, Erik Kassubek
All rights reserved.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: dedego-instrumenter
Project: Dynamic Analysis to detect potential deadlocks in concurrent Go programs
*/

/*
instrumenter.go
instrument files to work with the
"github.com/ErikKassubek/deadlockDetectorGo/src/dedego" libraries
*/

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ErikKassubek/deadlockDetectorGo/src/gui"
)

const MAX_TOTAL_WAITING_TIME_SEC = "20"
const SELECT_WAITING_TIME string = "2 * time.Second"

/*
Function to perform instrumentation of all list of files
@param file_paths []string: list of file names to instrument
@param gui *gui.GuiElements: gui elements to display output
@param status *gui.Status: status of the program
@return map[int]int: map of select id to size
@return error: error or nil
*/
func InstrumentFiles(elements *gui.GuiElements,
	status *gui.Status) (map[int]int, error) {
	file_paths, err := getAllFiles(status)
	if err != nil {
		elements.AddToOutput("Failed to get all files\n")
		return nil, err
	}

	for i, file := range file_paths {
		elements.Output.SetText(elements.Output.Text() +
			fmt.Sprintf("Instrumenting file %s.\n", file))
		elements.OutputScroll.ScrollToBottom()

		// instrument the file
		err := instrument_file(file, status)
		if err != nil {
			e := "Could not intrument file" + file + ".\n" + err.Error()
			elements.AddToOutput(e)
			return nil, err
		}
		prog := float64(i+1) / float64(len(file_paths))
		elements.ProgressInst.SetValue(prog)
	}

	// create turn select_ops into map
	select_map := make(map[int]int)
	for _, op := range select_ops {
		select_map[op.id] = op.size
	}

	return select_map, nil
}

/*
Function to instrument a given file.
@param file_path string: path to the file
@return error: error or nil
*/
func instrument_file(file_path string, status *gui.Status) error {
	// create output file
	output_file, err := os.Create(status.Output + get_relative_path(file_path, status))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output file %s.\n",
			file_path)
		return err
	}
	defer output_file.Close()

	// copy mod and sum files
	if file_path[len(file_path)-3:] != ".go" {
		content, err := ioutil.ReadFile(file_path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read file %s.\n", file_path)
			return err
		}
		_, err = output_file.Write(content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to output file %s.\n",
				file_path)
			return err
		}
		return nil
	}

	// instrument go files
	err = instrument_go_file(file_path, status)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not instrument %s\n", file_path)
		return err
	}

	return nil
}

/*
Function to instrument a go file.
@param file_path string: path to the file
@return error: error or nil
*/
func instrument_go_file(file_path string, status *gui.Status) error {
	// get the ASP of the file
	astSet := token.NewFileSet()

	f, err := parser.ParseFile(astSet, file_path, nil, 0)
	if err != nil {
		return err
	}

	instrument_chan(f, astSet)
	instrument_mutex(f)

	// print changed ast to output file
	output_path := status.Output + get_relative_path(file_path, status)
	output_file, err := os.OpenFile(output_path, os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open output file %s\n", output_path)
		return err
	}
	defer output_file.Close()

	if err := printer.Fprint(output_file, astSet, f); err != nil {
		return err
	}

	return nil
}

func get_relative_path(file string, status *gui.Status) string {
	totalPathLen := len(strings.Split(status.FolderPath, string(os.PathSeparator)))
	pathSplit := strings.Split(file, string(os.PathSeparator))
	return string(os.PathSeparator) + strings.Join(pathSplit[totalPathLen-1:],
		string(os.PathSeparator))
}
