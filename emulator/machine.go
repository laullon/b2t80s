package emulator

import (
	"fyne.io/fyne/v2"
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/ui"
)

var LoadSlow *bool

var Debug *bool
var Breaks *string
var WatchPoints *string

var DskAFile *string
var TapFile *string
var RomFile *string
var CartFile *string

type Machine interface {
	OnKeyEvent(event *fyne.KeyEvent)
	Monitor() Monitor
	Clock() Clock
	UIControls() []ui.Control
	Control() map[string]ui.Control
	GetVolumeControl() func(float64)
	SetDebugger(cpu.DebuggerCallbacks)
}
