package gui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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
Package: dedego-gui
Project: Dynamic Analysis to detect potential deadlocks in
*/

/*
tracerMutex.go
Drop in replacements for (rw)mutex and (Try)(R)lock and (Try)(R-)Unlock
*/

type GuiElements struct {
	PathLab              *widget.Label
	Output               *widget.TextGrid
	OutputScroll         *container.Scroll
	OpenBut              *widget.Button
	StartBut             *widget.Button
	ProgressInst         *widget.ProgressBar
	ProgressBuild        *widget.ProgressBar
	ProgressAna          *widget.ProgressBar
	Progress             *widget.Form
	Settings             *widget.Form
	SettingsMaxRuns      *widget.Entry
	SettingMaxTime       *widget.Entry
	SettingMaxSelectTime *widget.Entry
}

type Status struct {
	Output                  string
	FolderPath              string
	Name                    string
	InstrumentationComplete bool
	SettingsMaxRuns         int
	SettingsMaxFailed       int
	SettingMaxTime          int
	SettingMaxSelectTime    int
}

func (gui *GuiElements) AddToOutput(output string) {
	// add output to the output text grid withm
	gui.Output.SetText(gui.Output.Text() + output + "\n")
	gui.OutputScroll.ScrollToBottom()
}

func (gui *GuiElements) ClearOutput() {
	gui.Output.SetText("")
}
