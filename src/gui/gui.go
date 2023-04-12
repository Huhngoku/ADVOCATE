package gui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type GuiElements struct {
	PathLab       *widget.Label
	Output        *widget.TextGrid
	OutputScroll  *container.Scroll
	OpenBut       *widget.Button
	StartBut      *widget.Button
	ProgressInst  *widget.ProgressBar
	ProgressBuild *widget.ProgressBar
	ProgressDet   *widget.ProgressBar
	Progress      *widget.Form
}

type Status struct {
	Output                  string
	FolderPath              string
	Name                    string
	InstrumentationComplete bool
}

func (gui *GuiElements) AddToOutput(output string) {
	// add output to the output text grid withm
	gui.Output.SetText(gui.Output.Text() + output + "\n")
	gui.OutputScroll.ScrollToBottom()
}

func (gui *GuiElements) ClearOutput() {
	gui.Output.SetText("")
}
