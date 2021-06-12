package emulator

import (
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/gui"
	"github.com/veandco/go-sdl2/sdl"
)

var LoadSlow *bool

var Debug *bool
var Breaks *string
var WatchPoints *string

var DskAFile *string
var TapFile *string
var RomFile *string
var Test *bool
var CartFile *string

type Machine interface {
	OnKey(key sdl.Scancode)
	Monitor() Monitor
	Clock() Clock
	UIControls() []gui.GUIObject
	Control() map[string]gui.GUIObject
	GetVolumeControl() func(float64)
	SetDebugger(cpu.DebuggerCallbacks)
	Reset()
}
