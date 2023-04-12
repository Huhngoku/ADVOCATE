package main

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
