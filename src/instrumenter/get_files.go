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
get_files.go
Get all files from the input path
*/

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ErikKassubek/deadlockDetectorGo/src/gui"
)

// get all files in in and write them into file_names
// create folder structure in out
/*
Function to get all files in the in folder and add there names to file_names.
The function also copies the folder structure into the output folder.
@param path string: path to the input folder
@param status *gui.Status: Status object
@return []string: list of file names
@return error: Error or nil
*/
func getAllFiles(status *gui.Status) ([]string, error) {
	// remove old output folder
	err := os.RemoveAll(status.Output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to remove old output folder %s.\n", status.Output)
		return make([]string, 0), err
	}

	var file_names []string = make([]string, 0)

	// get all file names
	err = filepath.Walk(status.FolderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not walk through file path %s.", path)
				return err
			}
			// only save go, mod and sum files
			if len(path) >= 4 && (path[len(path)-3:] == ".go" ||
				path[len(path)-4:] == ".mod" || path[len(path)-4:] == ".sum") {
				file_names = append(file_names, path)
			}
			return nil
		})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to walk through file path.\n")
		return make([]string, 0), err
	}

	// get folder structure in path and copy it to out
	folders := make([]string, 0)
	in_split := strings.Split(status.FolderPath, string(os.PathSeparator))
	folders = append(folders, status.Output+in_split[len(in_split)-2])
	err = filepath.WalkDir(status.FolderPath,
		func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not walk through dir path %s.\n", path)
				return err
			}

			if info.IsDir() && path[:1] != "." {
				folders = append(folders, status.Output+get_relative_path(path, status))
			}
			return nil
		})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to walk through dir path.\n")
		return make([]string, 0), err
	}

	for _, folder := range folders {
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create folder %s.\n", folder)
			return make([]string, 0), err
		}
	}

	return file_names, nil
}
