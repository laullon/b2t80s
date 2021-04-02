package emulator

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/laullon/b2t80s/cpu"
)

type Debugger interface {
	cpu.DebuggerCallbacks

	UI() *fyne.Container
}

type debugger struct {
	clock Clock

	doStop          bool
	doStopInterrupt bool
	doStopLine      bool
	doStopFrame     bool
	breaks          []uint16

	ui, stop, step *fyne.Container
}

func NewDebugger(clock Clock, breaks []uint16) Debugger {

	debug := &debugger{
		clock:  clock,
		breaks: breaks,
	}

	debug.stop = fyne.NewContainerWithLayout(
		layout.NewGridLayoutWithColumns(2),
		widget.NewButton("Stop", func() {
			debug.Stop()
		}),
		widget.NewButton("Stop on Interrup", func() {
			debug.StopNextInterrupt()
		}),
	)

	debug.step = fyne.NewContainerWithLayout(
		layout.NewGridLayoutWithColumns(4),
		widget.NewButton("Continue", func() {
			debug.Continue()
		}),
		widget.NewButton("Step", func() {
			debug.Step()
		}),
		widget.NewButton("Step Line", func() {
			debug.StepLine()
		}),
		widget.NewButton("Step Frame", func() {
			debug.StepFrame()
		}),
	)

	debug.ui = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		debug.stop,
		debug.step,
	)

	debug.pause()
	return debug
}

func (debug *debugger) Eval(pc uint16) {
	for _, brk := range debug.breaks {
		if brk == pc {
			debug.Stop()
		}
	}

	if debug.doStop {
		debug.doStop = false
		debug.pause()
	}
}

func (debug *debugger) EvalInterrupt() {
	if debug.doStopInterrupt {
		debug.doStopInterrupt = false
		debug.pause()
	}
}

func (debug *debugger) EvalLine() bool {
	if debug.doStopLine {
		debug.doStopLine = false
		debug.pause()
		return true
	}
	return false
}

func (debug *debugger) EvalFrame() bool {
	if debug.doStopFrame {
		debug.doStopFrame = false
		debug.pause()
		return true
	}
	return false
}

func (debug *debugger) Stop() {
	debug.doStop = true
}

func (debug *debugger) StopNextInterrupt() {
	debug.doStopInterrupt = true
	debug.Continue()
}

func (debug *debugger) StopNextLine() {
	debug.doStopLine = true
	debug.Continue()
}

func (debug *debugger) StopNextFrame() {
	debug.doStopFrame = true
	debug.Continue()
}

func (debug *debugger) Step() {
	debug.doStop = true
	debug.Continue()
}

func (debug *debugger) StepLine() {
	debug.doStopLine = true
	debug.Continue()
}

func (debug *debugger) StepFrame() {
	debug.doStopFrame = true
	debug.Continue()
}

func (debug *debugger) UI() *fyne.Container {
	return debug.ui
}

func (debug *debugger) Continue() {
	for _, b := range debug.step.Objects {
		b.(*widget.Button).Disable()
	}
	for _, b := range debug.stop.Objects {
		b.(*widget.Button).Enable()
	}
	debug.ui.Refresh()
	debug.clock.Resume()
}

func (debug *debugger) pause() {
	for _, b := range debug.step.Objects {
		b.(*widget.Button).Enable()
	}
	for _, b := range debug.stop.Objects {
		b.(*widget.Button).Disable()
	}
	debug.ui.Refresh()
	debug.clock.Pause()
}
