package gui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type GuiElements struct {
	PathLab      *widget.Label
	Output       *widget.TextGrid
	OutputScroll *container.Scroll
	OpenBut      *widget.Button
	StartBut     *widget.Button
	CancelBut    *widget.Button
	ProgressInst *widget.ProgressBar
	ProgressDet  *widget.ProgressBar
	Progress     *widget.Form
}

type Status struct {
	Output       string
	FolderPath   string
	WasCancelled bool
}
