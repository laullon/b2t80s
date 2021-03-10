package m6502

import (
	"fmt"
	"strings"

	"github.com/laullon/b2t80s/cpu"
)

type Registers struct {
	A, X, Y uint8
	SP      uint8
	PC      uint16
	PS      Flags
}

type Flags struct {
	C bool
	Z bool
	I bool
	D bool
	B bool
	X bool
	V bool
	N bool
}

// 7  bit  0
// ---- ----
// NVss DIZC
// |||| ||||
// |||| |||+- Carry
// |||| ||+-- Zero
// |||| |+--- Interrupt Disable
// |||| +---- Decimal
// ||++------ No CPU effect, see: the B flag
// |+-------- Overflow
// +--------- Negative

func (f *Flags) set(v uint8) {
	f.C = v&0b0000_0001 != 0
	f.Z = v&0b0000_0010 != 0
	f.I = v&0b0000_0100 != 0
	f.D = v&0b0000_1000 != 0
	f.B = v&0b0001_0000 != 0
	f.X = v&0b0010_0000 != 0
	f.V = v&0b0100_0000 != 0
	f.N = v&0b1000_0000 != 0
}

func (f *Flags) Get() uint8 {
	var res uint8
	if f.C {
		res |= 0b00000001
	}
	if f.Z {
		res |= 0b00000010
	}
	if f.I {
		res |= 0b00000100
	}
	if f.D {
		res |= 0b00001000
	}
	if f.B {
		res |= 0b00010000
	}
	if f.X {
		res |= 0b00100000
	}
	if f.V {
		res |= 0b01000000
	}
	if f.N {
		res |= 0b10000000
	}
	return res
}

func (f Flags) String() string {
	var sb strings.Builder
	if f.C {
		sb.WriteString("C")
	} else {
		sb.WriteString("-")
	}
	if f.Z {
		sb.WriteString("Z")
	} else {
		sb.WriteString("-")
	}
	if f.I {
		sb.WriteString("I")
	} else {
		sb.WriteString("-")
	}
	if f.D {
		sb.WriteString("D")
	} else {
		sb.WriteString("-")
	}
	if f.B {
		sb.WriteString("B")
	} else {
		sb.WriteString("-")
	}
	if f.X {
		sb.WriteString("X")
	} else {
		sb.WriteString("-")
	}
	if f.V {
		sb.WriteString("V")
	} else {
		sb.WriteString("-")
	}
	if f.N {
		sb.WriteString("N")
	} else {
		sb.WriteString("-")
	}
	return sb.String()
}

func (r Registers) String() string {
	return fmt.Sprintf("A:0x%02X X:0x%02X Y:0x%02X SP:0x%02X PC:0x%04X PS:(0x%02X)%v", r.A, r.X, r.Y, r.SP, r.PC, r.PS.Get(), r.PS)
}

type M6502 interface {
	cpu.CPU
	Registers() *Registers
}

type m6502 struct {
	regs Registers

	doIRQ   bool
	doNMI   bool
	onNMI   bool
	doReset bool
	doWait  bool

	bus Bus

	op     operation
	nextOp operation

	log      cpu.CPUTracer
	debugger cpu.DebuggerCallbacks
}

func MewM6502(bus Bus) M6502 {
	return &m6502{
		bus: bus,
		op:  &reset{},
	}
}

func (cpu *m6502) Interrupt(i bool)                     { cpu.doIRQ = i }
func (cpu *m6502) NMI(i bool)                           { cpu.doNMI = i }
func (cpu *m6502) Halt()                                {}
func (cpu *m6502) Reset()                               { cpu.doReset = true }
func (cpu *m6502) Wait(w bool)                          { cpu.doWait = w }
func (cpu *m6502) Registers() *Registers                { return &cpu.regs }
func (cpu *m6502) SetTracer(t cpu.CPUTracer)            { cpu.log = t }
func (cpu *m6502) SetDebugger(db cpu.DebuggerCallbacks) { cpu.debugger = db }

// func (cpu *m6502) RegisterTrap(pc uint16, trap cpu.CPUTrap) {}
func (cpu *m6502) CurrentOP() string { return fmt.Sprintf("%v", cpu.op) }

func (cpu *m6502) Status() string     { return cpu.regs.String() }
func (cpu *m6502) FullStatus() string { return cpu.regs.String() }

func (cpu *m6502) Tick() {

	if cpu.doWait {
		return
	}

	done := cpu.op.tick(cpu)

	if done {
		if cpu.log != nil {
			cpu.log.AppendLastOP(dumpOperation(cpu.op))
			cpu.log.SetNextOP(dumpOperation(cpu.nextOp))
			// cpu.log.SetDiss(disassemble(cpu.nextOp.getPC(), cpu.bus.(*bus)))
		}
		if cpu.nextOp != nil {
			if cpu.debugger != nil {
				cpu.debugger.Eval(cpu.nextOp.getPC())
			}
			cpu.op = cpu.nextOp
			cpu.nextOp = nil
		} else {
			fmt.Printf("no nextOp after -> %-30v \n", cpu.op)
			panic(-1)
		}
	}
}

func (cpu *m6502) preFetch() {
	var newOp operation

	if cpu.doReset {
		newOp = &reset{}
		newOp.setPC(0)
		cpu.doReset = false
	} else if cpu.doNMI && !cpu.onNMI {
		newOp = &brk{imm: true}
		newOp.setPC(0)
		cpu.onNMI = true
		cpu.doNMI = false
	} else if cpu.doIRQ && !cpu.regs.PS.I {
		newOp = &brk{irq: true}
		newOp.setPC(0)
		cpu.doIRQ = false
	} else {
		opCode := cpu.bus.Read(cpu.regs.PC)
		newOp = ops[opCode]
		if newOp == nil {
			newOp = &unsupported{}
			newOp.setup(opCode)
		}
		newOp = newOp.Clone()
		newOp.setPC(cpu.regs.PC)
		cpu.regs.PC++
	}

	cpu.nextOp = newOp
}

func (cpu *m6502) push(data uint8) {
	// fmt.Printf("[PUSH] 0x%04X - 0x%02X \n", 0x0100+uint16(cpu.regs.SP), data)
	cpu.bus.Write(0x0100+uint16(cpu.regs.SP), data)
	cpu.regs.SP--
}

func (cpu *m6502) pop() uint8 {
	cpu.regs.SP++
	return cpu.bus.Read(0x0100 + uint16(cpu.regs.SP))
}
