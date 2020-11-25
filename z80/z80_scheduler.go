package z80

import (
	"github.com/pkg/errors"
)

type z80op interface {
	tick(cpu *z80)
	isDone() bool
	reset() // TODO i hate this
}

type basicOp struct {
	t    uint8
	done bool
}

func (op *basicOp) reset() {
	op.t = 0
	op.done = false
}

func (op *basicOp) isDone() bool {
	return op.done
}

type z80f func(*z80, []uint8)

// -------------------------------------------------------------

type fetch struct {
	basicOp
}

func (ops *fetch) tick(cpu *z80) {
	ops.t++
	switch ops.t {
	case 1:
		cpu.fetched = nil
		cpu.regs.M1 = true
		cpu.Bus.SetAddr(cpu.regs.PC)
		cpu.regs.PC++
	case 3:
		cpu.regs.M1 = false
		cpu.Bus.ReadMemory()
		d := cpu.Bus.GetData()
		cpu.fetched = append(cpu.fetched, d)
	case 4:
		op := lookup[cpu.fetched[0]]
		if op == nil {
			panic(errors.Errorf("opCode '0x%02X' not found", cpu.fetched[0]))
		}
		println("op", op.String(), cpu.fetched[0], cpu.regs.PC)
		for _, op := range op.ops {
			op.reset()
		}
		cpu.scheduler = append(cpu.scheduler, op.ops...)
		cpu.scheduler = append(cpu.scheduler, &fetch{})
		if op.onFetch != nil {
			op.onFetch(cpu, cpu.fetched)
		}
		ops.done = true
	}
}

// -------------------------------------------------------------

type mrPC struct {
	basicOp
	f z80f
}

func (ops *mrPC) tick(cpu *z80) {
	ops.t++
	switch ops.t {
	case 1:
		cpu.Bus.SetAddr(cpu.regs.PC)
		cpu.regs.PC++
	case 3:
		cpu.Bus.ReadMemory()
		d := cpu.Bus.GetData()
		cpu.fetched = append(cpu.fetched, d)
		if ops.f != nil {
			ops.f(cpu, cpu.fetched)
		}
		ops.done = true
	}
}

//
// -------------------------------------------------------------

type mr struct {
	basicOp
	f    z80f
	from uint16
}

func (ops *mr) tick(cpu *z80) {
	ops.t++
	switch ops.t {
	case 1:
		cpu.Bus.SetAddr(ops.from)
	case 3:
		cpu.Bus.ReadMemory()
		d := cpu.Bus.GetData()
		if ops.f != nil {
			ops.f(cpu, []byte{d})
		}
		ops.done = true
	}
}

// -------------------------------------------------------------

type mw struct {
	basicOp
	addr uint16
	data uint8
	f    func(*z80)
}

func (ops *mw) tick(cpu *z80) {
	ops.t++
	switch ops.t {
	case 1:
		cpu.Bus.SetAddr(ops.addr)
	case 2:
		cpu.Bus.SetData(ops.data)
	case 3:
		cpu.Bus.WriteMemory()
		if ops.f != nil {
			ops.f(cpu)
		}
		ops.done = true
	}
}

// -------------------------------------------------------------

type exec struct {
	basicOp
	l byte
	f z80f
}

func (ops *exec) tick(cpu *z80) {
	ops.t++
	if ops.t == ops.l {
		ops.f(cpu, cpu.fetched)
		ops.done = true
	}
}