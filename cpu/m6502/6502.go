package m6502

import "fmt"

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
	V bool
	N bool
}

type m6502 struct {
	regs Registers
	ba

	mem []byte

	op operation
}

func newM6502(mem []byte) *m6502 {
	return &m6502{
		mem: mem,
		op:  &reset{},
	}
}

func (cpu *m6502) Tick() {
	if (cpu.op == nil) || cpu.op.done() {
		opCode := cpu.mem[int(cpu.regs.PC)]
		cpu.op = ops[opCode]
		fmt.Printf("opCode: 0x%X - op: %v\n", opCode, cpu.op)
		if cpu.op == nil {
			panic(-1)
		}
		cpu.regs.PC++
	} else {
		cpu.op.tick(cpu)
	}
}
