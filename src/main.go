package main

import (
	"deadlockDetectorGo/instrumenter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type status struct {
	output        string
	folder_path   string
	was_cancelled bool
}

func main() {
	app := app.New()
	window := app.NewWindow("Deadlock Go Detector")
	status := status{}
	
	
	// create elements
	pathLab := widget.NewLabel("Path:")
	output := widget.NewTextGrid()
	outputScroll := container.NewScroll(output)

	// create a scroll container for the output

	
	openBut := widget.NewButton("Open", nil)
	startBut := widget.NewButton("Start", nil)
	cancelBut := widget.NewButton("Cancel", nil)
	
	progressInst := widget.NewProgressBar()
	progressDet := widget.NewProgressBar()
	progress := widget.NewForm(
		widget.NewFormItem("Instrumentation", progressInst),
		widget.NewFormItem("Detection", progressDet),
	)
	
	// set button functions
	
	openBut.OnTapped = func() {
		// BUG: cancel creates segmentation violation
		fileDialog := dialog.NewFolderOpen(
			func(r fyne.ListableURI, _ error) {
				status.folder_path = r.Path()
				pathLab.SetText("Path: " + status.folder_path)
				}, window)
				fileDialog.Show()
			}
			
			startBut.OnTapped = func() {
				if status.folder_path == "" {
					output.SetText("No folder selected!")
					return
				} else {
					progress.Hidden = false
					cancelBut.Hidden = false
					go func() {
						instrumenter.Run(status.folder_path, output, outputScroll, progressInst)
					}()
				}
	}

	cancelBut.OnTapped = func() {
		progress.Hidden = true
	}

	// set initial state
	progress.Hidden = true
	cancelBut.Hidden = true

	// create layout
	leftGrid := container.NewVBox(pathLab, openBut, startBut, cancelBut, progress)
	grid := container.New(layout.NewGridLayout(2), leftGrid, outputScroll)

	// show window
	window.SetContent(grid)
	window.ShowAndRun()
}
