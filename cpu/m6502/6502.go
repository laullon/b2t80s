package m6502

import (
	"fmt"
	"strings"

	cpuUtils "github.com/laullon/b2t80s/cpu"
)

type Registers struct {
	A, X, Y uint8
	SP      uint8
	PC      uint16
	P       Flags
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
	if f.V {
		sb.WriteString("V")
	} else {
		sb.WriteString("-")
	}
	if f.X {
		sb.WriteString("X")
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
	return fmt.Sprintf("A:0x%02X X:0x%02X Y:0x%02X SP:0x%02X PC:0x%04X P:%v", r.A, r.X, r.Y, r.SP, r.PC, r.P)
}

type m6502 struct {
	regs Registers
	AB   *cpuUtils.RegPair

	mem []uint8

	op operation
}

func newM6502(mem []byte) *m6502 {
	return &m6502{
		mem: mem,
		// op:  &reset{},
		AB: &cpuUtils.RegPair{L: new(uint8), H: new(uint8)},
	}
}

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

	if cpu.op.done() {
		fmt.Printf("%-30v%v\n", cpu.op, cpu.regs)
	}
}

func (cpu *m6502) push(data uint8) {
	cpu.mem[0x0100+uint16(cpu.regs.SP)] = data
	cpu.regs.SP--
}

func (cpu *m6502) pop() uint8 {
	cpu.regs.SP++
	return cpu.mem[0x0100+uint16(cpu.regs.SP)]
}
