package main

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
main.go
main file to run the program
*/

import (
	"github.com/ErikKassubek/deadlockDetectorGo/src/gui"
	"github.com/ErikKassubek/deadlockDetectorGo/src/instrumenter"

	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New()
	window := app.NewWindow("Deadlock Go Detector")
	status := gui.Status{}
	elements := gui.GuiElements{}

	// create elements
	elements.PathLab = widget.NewLabel("Path:")
	elements.Output = widget.NewTextGrid()
	elements.OutputScroll = container.NewScroll(elements.Output)

	// create a scroll container for the output
	elements.OpenBut = widget.NewButton("Open", nil)
	elements.StartBut = widget.NewButton("Start", nil)

	// create progress bars
	elements.ProgressInst = widget.NewProgressBar()
	elements.ProgressBuild = widget.NewProgressBar()
	elements.ProgressDet = widget.NewProgressBar()
	elements.Progress = widget.NewForm(
		widget.NewFormItem("Instrumentation", elements.ProgressInst),
		widget.NewFormItem("Build", elements.ProgressBuild),
		widget.NewFormItem("Detection", elements.ProgressDet),
	)

	// set button functions

	elements.OpenBut.OnTapped = func() {
		// BUG: cancel creates segmentation violation
		fileDialog := dialog.NewFolderOpen(
			func(r fyne.ListableURI, _ error) {
				status.FolderPath = r.Path()
				splitPath := strings.Split(status.FolderPath, string(os.PathSeparator))
				status.Name = splitPath[len(splitPath)-1]
				elements.PathLab.SetText("Path: " + status.FolderPath)
			}, window)
		fileDialog.Show()
	}

	elements.StartBut.OnTapped = func() {
		if status.FolderPath == "" {
			elements.Output.SetText("No folder selected!")
			return
		} else {
			elements.Progress.Hidden = false
			elements.ClearOutput()
			go func() {
				err := instrumenter.Run(status.FolderPath, &elements, &status)
				if err != nil {
					elements.AddToOutput("Analysis Failed")
				} else {
					elements.AddToOutput("Analysis Complete")
				}
			}()
		}
	}

	// set initial state
	elements.Progress.Hidden = true

	// create layout
	leftGrid := container.NewVBox(elements.PathLab, elements.OpenBut,
		elements.StartBut, elements.Progress)
	grid := container.New(layout.NewGridLayout(2), leftGrid,
		elements.OutputScroll)

	// show window
	window.SetContent(grid)
	window.ShowAndRun()
}
