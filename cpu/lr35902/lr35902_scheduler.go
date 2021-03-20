package lr35902

import (
	"github.com/pkg/errors"
)

type lr35902op interface {
	tick(cpu *lr35902)
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

type lr35902f func(*lr35902)

// -------------------------------------------------------------

type fetch struct {
	basicOp
	table []*opCode
}

func (ops *fetch) tick(cpu *lr35902) {
	ops.t++
	// println("> [fetch]", ops.t, "pc:", fmt.Sprintf("0x%04X", cpu.regs.PC))
	switch ops.t {
	case 1:
		cpu.regs.M1 = true
		cpu.bus.SetAddr(cpu.regs.PC)
		cpu.regs.PC++
		cpu.regs.R = cpu.regs.R&0x80 | ((cpu.regs.R + 1) & 0x7f)
	case 3:
		cpu.regs.M1 = false
		cpu.bus.Read()
		d := cpu.bus.GetData()
		cpu.bus.Release()
		cpu.fetched.prefix = cpu.fetched.prefix << 8
		cpu.fetched.prefix |= uint16(cpu.fetched.opCode)
		cpu.fetched.opCode = d
	case 4:
		cpu.fetched.op = ops.table[cpu.fetched.opCode]
		if cpu.fetched.op == nil {
			panic(errors.Errorf("opCode '%X - %X' not found", cpu.fetched.prefix, cpu.fetched.opCode))
		}
		for _, op := range cpu.fetched.op.ops {
			op.reset()
		}
		// fmt.Printf("opCode '%X - %X'\n", cpu.fetched.prefix, cpu.fetched.opCode)
		cpu.scheduler.append(cpu.fetched.op.ops...)
		if cpu.fetched.op.onFetch != nil {
			cpu.fetched.op.onFetch(cpu)
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

type mrNpc struct {
	basicOp
	f lr35902f
}

func (ops *mrNpc) tick(cpu *lr35902) {
	ops.t++
	// println("> [mrNpc]", ops.t, "pc:", fmt.Sprintf("0x%04X", cpu.regs.PC))
	switch ops.t {
	case 1:
		cpu.bus.SetAddr(cpu.regs.PC)
		cpu.regs.PC++
	case 3:
		cpu.bus.Read()
		d := cpu.bus.GetData()
		cpu.bus.Release()
		cpu.fetched.n = d
		if ops.f != nil {
			ops.f(cpu)
		}
		ops.done = true
	}
}

// -------------------------------------------------------------

type mrNNpc struct {
	basicOp
	f lr35902f
}

func (ops *mrNNpc) tick(cpu *lr35902) {
	ops.t++
	// println("> [mrNNpc]", ops.t, "pc:", fmt.Sprintf("0x%04X", cpu.regs.PC))
	switch ops.t {
	case 1:
		cpu.bus.SetAddr(cpu.regs.PC)
		cpu.regs.PC++
	case 3:
		cpu.bus.Read()
		d := cpu.bus.GetData()
		cpu.bus.Release()
		cpu.fetched.n = d
	case 4:
		cpu.bus.SetAddr(cpu.regs.PC)
		cpu.regs.PC++
	case 6:
		cpu.bus.Read()
		d := cpu.bus.GetData()
		cpu.bus.Release()
		cpu.fetched.n2 = d
		cpu.fetched.nn = uint16(cpu.fetched.n) | (uint16(cpu.fetched.n2) << 8)
		if ops.f != nil {
			ops.f(cpu)
		}
		ops.done = true
	}
}

//
// -------------------------------------------------------------

type lr35902MRf func(cpu *lr35902, data byte)
type mr struct {
	basicOp
	f    lr35902MRf
	from uint16
}

func (ops *mr) tick(cpu *lr35902) {
	ops.t++
	switch ops.t {
	case 1:
		cpu.bus.SetAddr(ops.from)
	case 3:
		cpu.bus.Read()
		d := cpu.bus.GetData()
		cpu.bus.Release()
		if ops.f != nil {
			ops.f(cpu, d)
		}
		ops.done = true
	}
}

var mrsPool = newObjectPool(func() interface{} { return &mr{} })

func newMR(from uint16, f lr35902MRf) *mr {
	mr := mrsPool.next().(*mr)
	mr.reset()
	mr.from = from
	mr.f = f
	return mr
}

//
// -------------------------------------------------------------

type mw struct {
	basicOp
	addr uint16
	data uint8
	f    func(cpu *lr35902)
}

func (ops *mw) tick(cpu *lr35902) {
	ops.t++
	switch ops.t {
	case 1:
		cpu.bus.SetAddr(ops.addr)
	case 2:
		cpu.bus.SetData(ops.data)
	case 3:
		cpu.bus.Write()
		cpu.bus.Release()
		if ops.f != nil {
			ops.f(cpu)
		}
		ops.done = true
	}
}

var mwsPool = newObjectPool(func() interface{} { return &mw{} })

func newMW(addr uint16, data uint8, f func(*lr35902)) *mw {
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
	f    func(*lr35902)
}

func (ops *out) tick(cpu *lr35902) {
	ops.t++
	switch ops.t {
	case 1:
		cpu.bus.SetAddr(ops.addr)
		cpu.bus.SetData(ops.data)
		cpu.bus.Write()
		cpu.bus.Release()
	case 3:
		if ops.f != nil {
			ops.f(cpu)
		}
		ops.done = true
	}
}

// -------------------------------------------------------------

type lr35902EXECf func(cpu *lr35902)
type exec struct {
	basicOp
	l byte
	f lr35902EXECf
}

func (ops *exec) tick(cpu *lr35902) {
	ops.t++
	if ops.t == ops.l {
		if ops.f != nil {
			ops.f(cpu)
		}
		ops.done = true
	}
}
