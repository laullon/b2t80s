package emulator

import (
	"strings"
	"sync"

	"github.com/laullon/b2t80s/cpu"
)

type Debugger interface {
	cpu.DebuggerCallbacks
	Stop()
	StopNextFrame()

	Step()
	StepNextFrame()

	Continue()
}

type debugger struct {
	clock Clock

	doStop          bool
	doStopInterrupt bool
	breaks          []uint16
}

func NewDebugger(clock Clock, breaks []uint16) Debugger {

	debug := &debugger{
		clock:  clock,
		breaks: breaks,
		doStop: true,
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
		debug.clock.Pause()
	}
}

func (debug *debugger) EvalInterrupt() {
	if debug.doStopInterrupt {
		debug.doStopInterrupt = false
		debug.clock.Pause()
	}
}

func (debug *debugger) Stop() {
	debug.doStop = true
}

func (debug *debugger) StopNextFrame() {
	debug.doStopInterrupt = true
}

func (debug *debugger) StepNextFrame() {
}

func (debug *debugger) Step() {
	debug.doStop = true
	debug.clock.Resume()
}

func (debug *debugger) Continue() {
	debug.clock.Resume()
}

type Log interface {
	AddEntry(entry string)
	Print() string
}

type logTail struct {
	idx     uint8
	entries []string
	mask    uint8
	mu      sync.Mutex
}

func NewShorLogTail() Log {
	return &logTail{
		entries: make([]string, 0x8),
		mask:    0x07,
	}
}

func NewLogTail() Log {
	return &logTail{
		entries: make([]string, 0x10),
		mask:    0x0f,
	}
}

func NewLongLogTail() Log {
	return &logTail{
		entries: make([]string, 0x100),
		mask:    0xff,
	}
}

func (log *logTail) AddEntry(entry string) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.entries[log.idx] = entry
	log.idx = (log.idx + 1) & log.mask
}

func (log *logTail) Print() string {
	log.mu.Lock()
	defer log.mu.Unlock()

	var sb strings.Builder
	sb.WriteString(strings.Join(log.entries[log.idx:], "\n"))
	sb.WriteString("\n")
	sb.WriteString(strings.Join(log.entries[:log.idx], "\n"))
	return strings.Trim(sb.String(), "\n")
}
