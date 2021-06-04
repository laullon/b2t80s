package lr35902

import (
	"github.com/pkg/errors"
)

type lr35902op interface {
	tick(cpu *lr35902)
}

type lr35902f func(*lr35902)

// -------------------------------------------------------------

type fetch struct {
	table []*opCode
}

func (fetch *fetch) tick(cpu *lr35902) {
	// fmt.Printf("FT 0x%04X\n", cpu.regs.PC)
	cpu.fetched.opCode = cpu.bus.Read(cpu.regs.PC)
	cpu.regs.PC++
	cpu.fetched.op = fetch.table[cpu.fetched.opCode]
	cpu.decode()
}

func (cpu *lr35902) decode() {
	if cpu.fetched.op == nil {
		panic(errors.Errorf("opCode '%X - %X' not found on 0x%04X", cpu.fetched.prefix, cpu.fetched.opCode, cpu.fetched.pc))
	}

	switch cpu.fetched.op.Len {
	case 1:
		cpu.fetched.op.f(cpu)

	case 2:
		cpu.scheduler.append(&mr{from: cpu.regs.PC, f: func(cpu *lr35902, data byte) {
			cpu.fetched.n = data
			cpu.fetched.op.f(cpu)
		}})
		cpu.regs.PC++

	case 3:
		cpu.scheduler.append(&mr{from: cpu.regs.PC, f: func(cpu *lr35902, data byte) {
			cpu.fetched.n = data
		}})
		cpu.scheduler.append(&mr{from: cpu.regs.PC + 1, f: func(cpu *lr35902, data byte) {
			cpu.fetched.n2 = data
			cpu.fetched.nn = uint16(cpu.fetched.n) | (uint16(cpu.fetched.n2) << 8)
			cpu.fetched.op.f(cpu)
		}})
		cpu.regs.PC += 2
	}
}

// -------------------------------------------------------------

type lr35902MRf func(cpu *lr35902, data byte)
type mr struct {
	f    lr35902MRf
	from uint16
}

func (ops *mr) tick(cpu *lr35902) {
	// fmt.Printf("MR 0x%04X\n", ops.from)
	d := cpu.bus.Read(ops.from)
	if ops.f != nil {
		ops.f(cpu, d)
	}
}

// -------------------------------------------------------------

type mw struct {
	to uint16
	d  uint8
	f  func(cpu *lr35902)
}

func (ops *mw) tick(cpu *lr35902) {
	// fmt.Printf("MW 0x%04X 0x%02X\n", ops.to, ops.d)
	cpu.bus.Write(ops.to, ops.d)
	if ops.f != nil {
		ops.f(cpu)
	}
}

// -------------------------------------------------------------

type lr35902EXECf func(cpu *lr35902)
type exec struct {
	f lr35902EXECf
}

func (ops *exec) tick(cpu *lr35902) {
	// fmt.Printf("EX\n")
	if ops.f != nil {
		ops.f(cpu)
	}
}
