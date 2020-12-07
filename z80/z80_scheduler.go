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
	table []*opCode
}

func (ops *fetch) tick(cpu *z80) {
	ops.t++
	// println("> [fetch]", ops.t, "pc:", fmt.Sprintf("0x%04X", cpu.regs.PC))
	switch ops.t {
	case 1:
		cpu.regs.M1 = true
		cpu.Bus.SetAddr(cpu.regs.PC)
		cpu.regs.PC++
	case 3:
		cpu.regs.M1 = false
		cpu.Bus.ReadMemory()
		d := cpu.Bus.GetData()
		cpu.fetched = append(cpu.fetched, d)
	case 4:
		op := ops.table[cpu.fetched[len(cpu.fetched)-1]]
		if op == nil {
			panic(errors.Errorf("opCode '0x%02X' not found", cpu.fetched[0]))
		}
		for _, op := range op.ops {
			op.reset()
		}
		cpu.scheduler.append(op.ops...)
		if op.onFetch != nil {
			op.onFetch(cpu, cpu.fetched)
		}
		ops.done = true
	}
}

var fetchPool = newObjectPool(func() interface{} { return &fetch{} })

func newFetch(table []*opCode) *fetch {
	fetch := fetchPool.next().(*fetch)
	fetch.reset()
	fetch.table = table
	return fetch
}

// -------------------------------------------------------------

type mrPC struct {
	basicOp
	f z80f
}

func (ops *mrPC) tick(cpu *z80) {
	ops.t++
	// println("> [mrPC]", ops.t, "pc:", fmt.Sprintf("0x%04X", cpu.regs.PC))
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

var mrsPool = newObjectPool(func() interface{} { return &mr{} })

func newMR(from uint16, f z80f) *mr {
	mr := mrsPool.next().(*mr)
	mr.reset()
	mr.from = from
	mr.f = f
	return mr
}

//
// -------------------------------------------------------------

type in struct {
	basicOp
	f    z80f
	from uint16
}

func (ops *in) tick(cpu *z80) {
	ops.t++
	switch ops.t {
	case 1:
		cpu.Bus.SetAddr(ops.from)
	case 4:
		cpu.Bus.ReadPort()
		data := cpu.Bus.GetData()
		cpu.regs.F.S = data&0x0080 != 0
		cpu.regs.F.Z = data == 0
		cpu.regs.F.H = false
		cpu.regs.F.P = parityTable[data]
		cpu.regs.F.N = false
		if ops.f != nil {
			ops.f(cpu, []byte{data})
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

var mwsPool = newObjectPool(func() interface{} { return &mw{} })

func newMW(addr uint16, data uint8, f func(*z80)) *mw {
	mw := mwsPool.next().(*mw)
	mw.reset()
	mw.addr = addr
	mw.data = data
	mw.f = f
	return mw
}

// -------------------------------------------------------------

type out struct {
	basicOp
	addr uint16
	data uint8
	f    func(*z80)
}

func (ops *out) tick(cpu *z80) {
	ops.t++
	switch ops.t {
	case 1:
		cpu.Bus.SetAddr(ops.addr)
		cpu.Bus.SetData(ops.data)
		cpu.Bus.WriteMemory()
	case 3:
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
