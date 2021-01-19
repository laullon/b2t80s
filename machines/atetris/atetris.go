package atetris

import (
	"image"

	"fyne.io/fyne"
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/ui"
)

type atetris struct {
	clock emulator.Clock
	cpu   emulator.CPU
}

func NewATetris() machines.Machine {

	m := &atetris{
		cpu:   m6502.MewM6502(newBus()),
		clock: emulator.NewCLock(14318181 / 8),
	}

	m.clock.AddTicker(0, m.cpu)

	m.cpu.SetDebuger(m6502.NewDebugger(m.cpu, nil, m.clock))

	return m
}

func (t *atetris) OnKeyEvent(event *fyne.KeyEvent) {}
func (t *atetris) Debugger() emulator.Debugger     { return nil }

func (t *atetris) Monitor() emulator.Monitor {
	return emulator.NewMonitor(image.NewRGBA(image.Rect(0, 0, 200, 200)))
}

func (t *atetris) Clock() emulator.Clock           { return t.clock }
func (t *atetris) UIControls() []ui.Control        { return nil }
func (t *atetris) GetVolumeControl() func(float64) { return func(f float64) {} }
