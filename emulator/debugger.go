package emulator

import (
	"strconv"
	"strings"

	"github.com/laullon/b2t80s/cpu"
)

type Debugger interface {
	cpu.DebuggerCallbacks

	// UI() *fyne.Container
}

type debugger struct {
	clock Clock

	doStop          bool
	doStopInterrupt bool
	doStopLine      bool
	doStopFrame     bool
	breaks          []uint16

	// ui, stop, step *fyne.Container
}

func NewDebugger(clock Clock) *debugger {

	debug := &debugger{
		clock: clock,
	}

	if len(*Breaks) > 0 {
		bps := strings.Split(*Breaks, ",")
		for _, bp := range bps {
			n, err := strconv.ParseUint(bp, 0, 16)
			if err != nil {
				panic(err)
			}
			debug.breaks = append(debug.breaks, uint16(n))
		}
	}
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

func (debug *debugger) Continue() {
	debug.clock.Resume()
}

func (debug *debugger) pause() {
	debug.clock.Pause()
}
