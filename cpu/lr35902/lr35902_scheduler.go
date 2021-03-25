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
	mrN   *mrNpc
	mrNN  *mrNNpc
}

func (fetch *fetch) tick(cpu *lr35902) {
	fetch.t++
	// println("> [fetch]", fetch.t, "pc:", fmt.Sprintf("0x%04X", cpu.regs.PC))
	switch fetch.t {
	case 1:
		cpu.bus.SetAddr(cpu.regs.PC)
		cpu.regs.PC++

	case 2:
		cpu.bus.Read()
		d := cpu.bus.GetData()
		cpu.bus.Release()
		cpu.fetched.prefix = cpu.fetched.opCode
		cpu.fetched.opCode = d

	case 3:
		cpu.fetched.op = fetch.table[cpu.fetched.opCode]
		if cpu.fetched.op == nil {
			panic(errors.Errorf("opCode '%X - %X' not found on 0x%04X", cpu.fetched.prefix, cpu.fetched.opCode, cpu.fetched.pc))
		}

		switch cpu.fetched.op.Len {
		case 1:
			cpu.fetched.op.f(cpu)

		case 2:
			fetch.mrN.f = cpu.fetched.op.f
			cpu.scheduler.append(fetch.mrN)

		case 3:
			fetch.mrNN.f = cpu.fetched.op.f
			cpu.scheduler.append(fetch.mrNN)
		}
		fetch.done = true
	}
}

var fetchPool = newObjectPool(func() interface{} {
	return &fetch{
		mrN:  &mrNpc{},
		mrNN: &mrNNpc{},
	}
})

func newFetch(table []*opCode) *fetch {
	fetch := fetchPool.next().(*fetch)
	fetch.reset()
	fetch.mrN.reset()
	fetch.mrNN.reset()
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
		ops.f(cpu)
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
	case 5:
		cpu.bus.SetAddr(cpu.regs.PC)
		cpu.regs.PC++
	case 7:
		cpu.bus.Read()
		d := cpu.bus.GetData()
		cpu.bus.Release()
		cpu.fetched.n2 = d
		cpu.fetched.nn = uint16(cpu.fetched.n) | (uint16(cpu.fetched.n2) << 8)
		ops.f(cpu)
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
