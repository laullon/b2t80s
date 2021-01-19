package m6502

import (
	"fmt"
	"strings"

	"github.com/laullon/b2t80s/emulator"
)

type logEntry struct {
	pc  uint16
	mem string
	ins string
	ts  int
}

type debugger struct {
	cpu    *m6502
	clock  emulator.Clock
	memory emulator.Memory

	next *logEntry
	log  []*logEntry
	diss []*logEntry

	doStop          bool
	doStopInterrupt bool
	dump            bool

	ts int
}

func NewDebugger(cpu emulator.CPU, mem emulator.Memory, clock emulator.Clock) emulator.Debugger {
	debug := &debugger{
		cpu:    cpu.(*m6502),
		memory: mem,
		clock:  clock,
		doStop: false,
		dump:   true,
	}
	cpu.SetDebuger(debug)
	return debug
}

func (debug *debugger) SetDump(on bool) {
	debug.dump = on
}

func (debug *debugger) SetBreakPoint(bp uint16) {
}

func (debug *debugger) LoadSymbols(fileName string) {
}

func (debug *debugger) AddInstruction(pc uint16, mem, instruction string) {
	le := &logEntry{ins: instruction, mem: mem, pc: pc, ts: debug.ts}
	debug.log = append(debug.log, le)
	debug.ts = 0

	if debug.dump {
		print(le.String())
		fmt.Printf("%v\n", debug.cpu.regs)
	}

	if len(debug.log) > 10 {
		debug.log = debug.log[len(debug.log)-10 : len(debug.log)]
	}
}

// func (debug *debugger) NextInstruction(mem []byte) {
// 	op := decode(mem)
// 	if op == nil {
// 		return
// 	}
// 	opMem := make([]byte, op.len)
// 	copy(opMem, mem)
// 	debug.next = &logEntry{op: op, mem: opMem, pc: debug.cpu.regs.PC}
// }

func (debug *debugger) DumpNextFrame() {
}

func (debug *debugger) StopNextFrame() {
	debug.doStopInterrupt = true
}

func (debug *debugger) NextFrame() {
}

func (debug *debugger) Tick() {
	// debug.ts++
	// if debug.doStop {
	// 	debug.doStop = false
	// 	debug.clock.Pause()
	// }

	// if debug.cpu.doInterrupt && debug.doStopInterrupt {
	// 	debug.doStopInterrupt = false
	// 	debug.clock.Pause()
	// }
}

func (debug *debugger) Stop() {
	debug.doStop = true
}

func (debug *debugger) Step() {
	debug.doStop = true
	debug.clock.Resume()
}

func (debug *debugger) Continue() {
	debug.clock.Resume()
}

func (debug *debugger) GetStatus() string {
	// println(debug.getNext())
	return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", debug.getRegisters(), debug.getLog(), debug.getNext(), debug.getDiss())
}

func (debug *debugger) getLog() string {
	var log []string
	for _, le := range debug.log {
		if le == nil {
			break
		}
		log = append(log, le.String())
	}
	return strings.Join(log, "\n")
}

func (debug *debugger) getDiss() string {
	var diss []string
	for _, le := range debug.diss {
		if le == nil {
			break
		}
		diss = append(diss, le.String())
	}
	return strings.Join(diss, "\n")
}

func (debug *debugger) getNext() string {
	return debug.next.String()
}

func (debug *debugger) getRegisters() string {
	return fmt.Sprintf("%v", debug.cpu.regs)
}

func (le *logEntry) String() string {
	if le == nil {
		return ""
	}
	return fmt.Sprintf("0x%04X %-12s %-20s (%3d) ", le.pc, le.mem, le.ins, le.ts)
}

func dump(buff []byte) string {
	res := ""
	for _, b := range buff {
		res = fmt.Sprintf("%s%02X ", res, b)
	}
	return res
}
