package machines

import (
	"fyne.io/fyne"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/ui"
)

var LoadSlow *bool
var DskAFile *string
var TapFile *string
var RomFile *string

type Machine interface {
	OnKeyEvent(event *fyne.KeyEvent)
	Run()
	Debugger() emulator.Debugger
	Monitor() emulator.Monitor

	UIControls() []ui.Control

	GetVolumeControl() func(float64)
}
