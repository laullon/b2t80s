package m6502

import (
	"fmt"
	"strings"

	cpuUtils "github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/emulator"
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

func (f *Flags) set(v uint8) {
	f.C = v&0b00000001 != 0
	f.Z = v&0b00000010 != 0
	f.I = v&0b00000100 != 0
	f.D = v&0b00001000 != 0
	f.B = v&0b00010000 != 0
	f.X = v&0b00100000 != 0
	f.V = v&0b01000000 != 0
	f.N = v&0b10000000 != 0
}

func (f *Flags) get() uint8 {
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
	return fmt.Sprintf("A:0x%02X X:0x%02X Y:0x%02X SP:0x%02X PC:0x%04X PS:(0x%02X)%v", r.A, r.X, r.Y, r.SP, r.PC, r.PS.get(), r.PS)
}

type m6502 struct {
	regs Registers

	mem []uint8
	log cpuUtils.Log

	op operation

	debugger emulator.Debugger
}

func MewM6502(mem []byte) emulator.CPU {
	return &m6502{
		mem: mem,
		op:  &reset{},
	}
}

func (cpu *m6502) Interrupt(bool)                                {}
func (cpu *m6502) Halt()                                         {}
func (cpu *m6502) Wait(bool)                                     {}
func (cpu *m6502) Registers() interface{}                        { return cpu.regs }
func (cpu *m6502) SetDebuger(debugger emulator.Debugger)         { cpu.debugger = debugger }
func (cpu *m6502) RegisterTrap(pc uint16, trap emulator.CPUTrap) {}

func (cpu *m6502) Tick() {
	if (cpu.op == nil) || cpu.op.done() {
		opCode := cpu.mem[int(cpu.regs.PC)]
		cpu.op = ops[opCode]
		if cpu.op == nil {
			fmt.Printf("opCode: 0x%X NOT FOUND !!!\n", opCode)
			panic(-1)
		}
		cpu.op.reset()
		cpu.op.setPC(cpu.regs.PC)
		cpu.regs.PC++
	} else {
		cpu.op.tick(cpu)
	}

	if cpu.op.done() && (cpu.log != nil) {
		cpu.log.AddEntry(fmt.Sprintf("%-30v%v", cpu.op, cpu.regs))
	}

	if cpu.op.done() && (cpu.debugger != nil) {
		cpu.debugger.AddInstruction(cpu.op.getPC(), "", fmt.Sprintf("%-30v", cpu.op))
	}
}

func (cpu *m6502) push(data uint8) {
	// fmt.Printf("[PUSH] 0x%04X - 0x%02X \n", 0x0100+uint16(cpu.regs.SP), data)
	cpu.mem[0x0100+uint16(cpu.regs.SP)] = data
	cpu.regs.SP--
}

func (cpu *m6502) pop() uint8 {
	cpu.regs.SP++
	return cpu.mem[0x0100+uint16(cpu.regs.SP)]
}
