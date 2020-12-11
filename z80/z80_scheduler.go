package z80

import "github.com/pkg/errors"

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

type z80f func(*z80)

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
		cpu.bus.SetAddr(cpu.regs.PC)
		cpu.regs.PC++
	case 3:
		cpu.regs.M1 = false
		cpu.bus.ReadMemory()
		d := cpu.bus.GetData()
		cpu.opBytes = append(cpu.opBytes, d)
		cpu.fetched.prefix = cpu.fetched.prefix << 8
		cpu.fetched.prefix |= uint16(cpu.fetched.opCode)
		cpu.fetched.opCode = d
	case 4:
		op := ops.table[cpu.fetched.opCode]
		if op == nil {
			panic(errors.Errorf("opCode '0x%X' not found", cpu.opBytes))
		}
		for _, op := range op.ops {
			op.reset()
		}
		cpu.scheduler.append(op.ops...)
		if op.onFetch != nil {
			op.onFetch(cpu)
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
	f z80f
}

func (ops *mrNpc) tick(cpu *z80) {
	ops.t++
	// println("> [mrNpc]", ops.t, "pc:", fmt.Sprintf("0x%04X", cpu.regs.PC))
	switch ops.t {
	case 1:
		cpu.bus.SetAddr(cpu.regs.PC)
		cpu.regs.PC++
	case 3:
		cpu.bus.ReadMemory()
		d := cpu.bus.GetData()
		cpu.opBytes = append(cpu.opBytes, d)
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
	f z80f
}

func (ops *mrNNpc) tick(cpu *z80) {
	ops.t++
	// println("> [mrNNpc]", ops.t, "pc:", fmt.Sprintf("0x%04X", cpu.regs.PC))
	switch ops.t {
	case 1:
		cpu.bus.SetAddr(cpu.regs.PC)
		cpu.regs.PC++
	case 3:
		cpu.bus.ReadMemory()
		d := cpu.bus.GetData()
		cpu.opBytes = append(cpu.opBytes, d)
		cpu.fetched.n = d
	case 4:
		cpu.bus.SetAddr(cpu.regs.PC)
		cpu.regs.PC++
	case 6:
		cpu.bus.ReadMemory()
		d := cpu.bus.GetData()
		cpu.opBytes = append(cpu.opBytes, d)
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

type z80MRf func(cpu *z80, data byte)
type mr struct {
	basicOp
	f    z80MRf
	from uint16
}

func (ops *mr) tick(cpu *z80) {
	ops.t++
	switch ops.t {
	case 1:
		cpu.bus.SetAddr(ops.from)
	case 3:
		cpu.bus.ReadMemory()
		d := cpu.bus.GetData()
		if ops.f != nil {
			ops.f(cpu, d)
		}
		ops.done = true
	}
}

var mrsPool = newObjectPool(func() interface{} { return &mr{} })

func newMR(from uint16, f z80MRf) *mr {
	mr := mrsPool.next().(*mr)
	mr.reset()
	mr.from = from
	mr.f = f
	return mr
}

//
// -------------------------------------------------------------

type z80INf func(cpu *z80, data byte)

type in struct {
	basicOp
	f    z80INf
	from uint16
}

func (ops *in) tick(cpu *z80) {
	ops.t++
	switch ops.t {
	case 1:
		cpu.bus.SetAddr(ops.from)
	case 4:
		cpu.bus.ReadPort()
		data := cpu.bus.GetData()
		cpu.regs.F.S = data&0x0080 != 0
		cpu.regs.F.Z = data == 0
		cpu.regs.F.H = false
		cpu.regs.F.P = parityTable[data]
		cpu.regs.F.N = false
		if ops.f != nil {
			ops.f(cpu, data)
		}
		ops.done = true
	}
}

// -------------------------------------------------------------

type mw struct {
	basicOp
	addr uint16
	data uint8
	f    func(cpu *z80)
}

func (ops *mw) tick(cpu *z80) {
	ops.t++
	switch ops.t {
	case 1:
		cpu.bus.SetAddr(ops.addr)
	case 2:
		cpu.bus.SetData(ops.data)
	case 3:
		cpu.bus.WriteMemory()
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
		cpu.bus.SetAddr(ops.addr)
		cpu.bus.SetData(ops.data)
		cpu.bus.WritePort()
	case 3:
		if ops.f != nil {
			ops.f(cpu)
		}
		ops.done = true
	}
}

// -------------------------------------------------------------

type z80EXECf func(cpu *z80)
type exec struct {
	basicOp
	l byte
	f z80EXECf
}

func (ops *exec) tick(cpu *z80) {
	ops.t++
	if ops.t == ops.l {
		ops.f(cpu)
		ops.done = true
	}
}
