package z80

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
	cpu    *z80
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
		cpu:    cpu.(*z80),
		memory: mem,
		clock:  clock,
		doStop: false,
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
		println(debug.cpu.regs.dump())
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
	debug.ts++
	if debug.doStop {
		debug.doStop = false
		debug.clock.Pause()
	}

	if debug.cpu.doInterrupt && debug.doStopInterrupt {
		debug.doStopInterrupt = false
		debug.clock.Pause()
	}

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
	var res strings.Builder
	regs := debug.cpu.Registers().(*Z80Registers)
	res.WriteString(fmt.Sprintf("  A:0x%02X    F:0x%02X  AF:0x%04X    SP:0x%04X\n", regs.A, regs.F.GetByte(), uint16(regs.A)<<8|uint16(regs.F.GetByte()), regs.SP.Get()))
	res.WriteString(fmt.Sprintf("  B:0x%02X    C:0x%02X  BC:0x%04X    ---------\n", regs.B, regs.C, uint16(regs.B)<<8|uint16(regs.C)))
	res.WriteString(fmt.Sprintf("  D:0x%02X    E:0x%02X  DE:0x%04X    0x%04X\n", regs.D, regs.E, uint16(regs.D)<<8|uint16(regs.E), getWord(debug.memory, regs.SP.Get()+0)))
	res.WriteString(fmt.Sprintf("  H:0x%02X    L:0x%02X  HL:0x%04X    0x%04X\n", regs.H, regs.L, uint16(regs.H)<<8|uint16(regs.L), getWord(debug.memory, regs.SP.Get()+2)))
	res.WriteString(fmt.Sprintf("IXH:0x%02X  IXL:0x%02X  IX:0x%04X    0x%04X\n", regs.IXH, regs.IXL, uint16(regs.IXH)<<8|uint16(regs.IXL), getWord(debug.memory, regs.SP.Get()+4)))
	res.WriteString(fmt.Sprintf("IYH:0x%02X  IYL:0x%02X  IY:0x%04X    0x%04X\n", regs.IYH, regs.IYL, uint16(regs.IYH)<<8|uint16(regs.IYL), getWord(debug.memory, regs.SP.Get()+6)))
	res.WriteString(fmt.Sprintf("SZ5H3PNC            PC:0x%04X\n%08b", regs.PC, regs.F.GetByte()))
	return res.String()
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

func decode(mem []byte) *opCode {
	var op *opCode
	switch mem[0] {
	case 0xCB:
		op = lookupCB[mem[1]]
	case 0xDD:
		if mem[1] == 0xCB {
			op = lookupDDCB[mem[2]]
		} else {
			op = lookupDD[mem[1]]
		}
	case 0xED:
		op = lookupED[mem[1]]
	case 0xFD:
		if mem[1] == 0xCB {
			op = lookupFDCB[mem[2]]
		} else {
			op = lookupFD[mem[1]]
		}
	default:
		op = lookup[mem[0]]
	}
	return op
}