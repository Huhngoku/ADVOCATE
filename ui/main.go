package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Deadlock Go Detector")

	path := widget.NewLabel("Path:")
	// create folder picker to select the folder
	openBut := widget.NewButton("Open", func() {
		fileDialog := dialog.NewFolderOpen(
			func(r fyne.ListableURI, _ error) {
				path.SetText("Path: " + r.Path())
			}, w)
		fileDialog.Show()
	})

	// create progress bars
	progressInst := widget.NewProgressBar()
	progressDet := widget.NewProgressBar()
	progress := widget.NewForm(
		widget.NewFormItem("Instrumentation", progressInst),
		widget.NewFormItem("Detection", progressDet))
	progress.Hidden = true

	// create a cancel button
	cancelBut := widget.NewButton("Cancel", func() {
		progress.Hidden = true
	})
	cancelBut.Hidden = true

	// create a start button
	startBut := widget.NewButton("Start", func() {
		progress.Hidden = false
		cancelBut.Hidden = false
		// create a file dialog
		go func() {
			for i := 0.0; i <= 1.0; i += 0.01 {
				time.Sleep(time.Millisecond * 250)
				progressInst.SetValue(i)
				progressDet.SetValue(i)
			}
		}()

	})

	// create a vertical layout
	leftGrid := container.NewVBox(path, openBut, startBut, cancelBut, progress)

	// create a text widget with black background
	output := widget.NewLabel("")

	// create a grid layout with 2 columns
	grid := container.New(layout.NewGridLayout(2), leftGrid, output)

	w.SetContent(grid)
	w.ShowAndRun()
}
