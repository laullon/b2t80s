package emulator

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/ui"
)

var LoadSlow *bool

var Debug *bool
var Breaks *string
var WatchPoints *string

var MachineStatus = &machineStatus{}

var DskAFile *string
var TapFile *string
var RomFile *string
var App fyne.App

type Machine interface {
	OnKeyEvent(event *fyne.KeyEvent)
	Monitor() Monitor
	Clock() Clock
	UIControls() []ui.Control
	CPUControl() ui.Control
	GetVolumeControl() func(float64)
	SetDebugger(cpu.DebuggerCallbacks)
}

type machineStatus struct {
	log       []string
	cpuStatus string
}

func (status *machineStatus) AddInstruction(instruction string) {
	status.log = append(status.log, instruction)
	if len(status.log) == 11 {
		status.log = status.log[1:]
	}
}

func (status *machineStatus) Status() string {
	return fmt.Sprintf("%s\n\n%s", status.cpuStatus, strings.Join(status.log, "\n"))
}
