package m6502

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type operation interface {
	tick(cpu *m6502) (done bool)
	setup(opCode uint8)
	setPC(pc uint16)
	getPC() uint16
	String() string
	Clone() operation
}

type basicop struct {
	ins    string
	opCode uint8

	pc uint16
	t  uint
}

func (op *basicop) setPC(pc uint16) {
	op.pc = pc
}

func (op *basicop) getPC() uint16 {
	return op.pc
}

type reset struct {
	basicop
}

func (op *reset) Clone() operation {
	return &reset{}
}

func (op *reset) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 7:
		cpu.regs.PC = uint16(cpu.bus.Read(0xfffc))
		cpu.regs.PC |= uint16(cpu.bus.Read(0xfffd)) << 8
		cpu.regs.SP = 0xfd
		cpu.preFetch()
		done = true
	}
	op.t++
	return
}

func (op *reset) setup(opCode uint8) {
	op.opCode = opCode
}

func (op *reset) String() string {
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb)
	writeOP(sb, "RESET")
	return sb.String()
}

// -----
type brk struct {
	basicop
	vector   uint16
	irq, imm bool
}

func (op *brk) Clone() operation {
	return &brk{}
}

func (op *brk) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 0:
		if cpu.debugger != nil {
			cpu.debugger.EvalInterrupt()
		}
		if op.imm || op.irq {
			cpu.bus.Read(cpu.regs.PC)
		} else {
			cpu.preFetch()
			cpu.regs.PC++
		}
	case 1:
		cpu.push(uint8(cpu.regs.PC >> 8))
	case 2:
		cpu.push(uint8(cpu.regs.PC))
		if op.imm {
			op.vector = 0xFFFA
		} else {
			op.vector = 0xFFFE
		}
	case 3:
		if !op.imm && !op.irq {
			cpu.regs.PS.B = true
		} else {
			cpu.regs.PS.B = false
		}
		cpu.regs.PS.X = true
		cpu.push(cpu.regs.PS.Get())
		cpu.regs.PS.X = false
	case 4:
		cpu.regs.PC = uint16(cpu.bus.Read(op.vector))
		cpu.regs.PS.I = true
	case 5:
		cpu.regs.PC |= uint16(cpu.bus.Read(op.vector+1)) << 8
		cpu.preFetch()
		done = true
	}
	op.t++
	return
}

func (op *brk) setup(opCode uint8) {
	op.opCode = opCode
}

func (op *brk) String() string {
	mod := ""
	if op.irq {
		mod = "-IRQ-"
	} else if op.imm {
		mod = "-IMM-"
	}
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode)
	writeOP(sb, "BRK ", mod)
	return sb.String()
}

// -----
type indirectX struct {
	basicop
	r, w, rw bool
	f        interface{}
	addr     uint16
	addrZ    uint8
}

func (op *indirectX) Clone() operation {
	return &indirectX{basicop: op.basicop, f: op.f, r: op.r, w: op.w, rw: op.rw}
}

func (op *indirectX) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 1:
		op.addrZ = cpu.bus.Read(cpu.regs.PC)
		cpu.regs.PC++
	case 2:
		op.addrZ += cpu.regs.X
	case 3:
		op.addr = uint16(cpu.bus.Read(uint16(op.addrZ)))
	case 4:
		op.addr |= uint16(cpu.bus.Read(uint16(op.addrZ+1))) << 8
	case 5:
		exec(cpu, op.f, op.addr, op.r, op.w, op.rw)
		done = true
	}
	op.t++
	return
}

func (op *indirectX) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
	op.r, op.w, op.rw = getRWRW(op.f)
}

func (op *indirectX) String() string {
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode, op.addrZ)
	writeOP(sb, op.ins, " ($", toHex8(op.addrZ), ", X)")
	return sb.String()
}

// -----
type indirectY struct {
	basicop
	r, w, rw   bool
	f          interface{}
	addr       uint16
	targetAddr uint16
	addrZ      uint8
}

func (op *indirectY) Clone() operation {
	return &indirectY{basicop: op.basicop, f: op.f, r: op.r, w: op.w, rw: op.rw}
}

func (op *indirectY) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 1:
		op.addrZ = cpu.bus.Read(cpu.regs.PC)
		cpu.regs.PC++
	case 2:
		op.addr = uint16(cpu.bus.Read(uint16(op.addrZ)))
	case 3:
		op.addr |= uint16(cpu.bus.Read(uint16(op.addrZ+1))) << 8
	case 4:
		op.targetAddr = op.addr + uint16(cpu.regs.Y)
		if (op.targetAddr&0xff00 == op.addr&0xff00) && !op.w { // no page change
			exec(cpu, op.f, op.targetAddr, op.r, op.w, op.rw)
			done = true
		}
	case 5:
		exec(cpu, op.f, op.targetAddr, op.r, op.w, op.rw)
		done = true
	}
	op.t++
	return
}

func (op *indirectY) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
	op.r, op.w, op.rw = getRWRW(op.f)
}

func (op *indirectY) String() string {
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode, op.addrZ)
	writeOP(sb, op.ins, " ($", toHex8(op.addrZ), "), Y")
	return sb.String()
}

// -----
type implicit struct {
	basicop
	f func(cpu *m6502)
	a bool
}

func (op *implicit) Clone() operation {
	return &implicit{basicop: op.basicop, f: op.f, a: op.a}
}

func (op *implicit) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 1:
		cpu.preFetch()
		op.f(cpu)
		done = true
	}
	op.t++
	return
}

func (op *implicit) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
}

func (op *implicit) String() string {
	mod := ""
	if op.a {
		mod = " a"
	}
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode)
	writeOP(sb, op.ins, mod)
	return sb.String()
}

// -----

type immediate struct {
	basicop
	f    func(cpu *m6502, data uint8)
	data uint8
}

func (op *immediate) Clone() operation {
	return &immediate{basicop: op.basicop, f: op.f}
}

func (op *immediate) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 1:
		op.data = cpu.bus.Read(cpu.regs.PC)
		cpu.regs.PC++
		cpu.preFetch()
		op.f(cpu, op.data)
		done = true
	}
	op.t++
	return
}

func (op *immediate) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
}

func (op *immediate) String() string {
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode, op.data)
	writeOP(sb, op.ins, " #$", toHex8(op.data))
	return sb.String()
}

// -----

type relative struct {
	basicop
	f      func(cpu *m6502) bool
	off    int8
	target uint16
}

func (op *relative) Clone() operation {
	return &relative{basicop: op.basicop, f: op.f}
}

func (op *relative) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 1:
		op.off = int8(cpu.bus.Read(cpu.regs.PC))
		cpu.regs.PC++
		op.target = cpu.regs.PC + uint16(op.off)
		done = !op.f(cpu)
		if done {
			cpu.preFetch()
		}
	case 2:
		if (cpu.regs.PC & 0xff00) == (op.target & 0xff00) { // no page change
			cpu.regs.PC = op.target
			cpu.preFetch()
			done = true
		}
	case 3:
		cpu.regs.PC = op.target
		cpu.preFetch()
		done = true
	}
	op.t++
	return
}

func (op *relative) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
}

func (op *relative) String() string {
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode, uint8(op.off))
	writeOP(sb, op.ins, " $", toHex16(op.target))
	return sb.String()
}

// -----

type absoluteJMP struct {
	basicop
	readAddr uint16
}

func (op *absoluteJMP) Clone() operation {
	return &absoluteJMP{basicop: op.basicop}
}

func (op *absoluteJMP) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 1:
		op.readAddr = uint16(cpu.bus.Read(cpu.regs.PC))
		cpu.regs.PC++
	case 2:
		op.readAddr |= uint16(cpu.bus.Read(cpu.regs.PC)) << 8
		cpu.regs.PC = op.readAddr
		cpu.preFetch()
		done = true
	}
	op.t++
	return
}

func (op *absoluteJMP) setup(opCode uint8) {
	op.opCode = opCode
}

func (op *absoluteJMP) String() string {
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode, uint8(op.readAddr&0x0ff), uint8(op.readAddr>>8))
	writeOP(sb, "jmp $", toHex16(op.readAddr))
	return sb.String()
}

// -----

type absoluteJSR struct {
	basicop
	readAddr uint16
}

func (op *absoluteJSR) Clone() operation {
	return &absoluteJSR{basicop: op.basicop}
}

func (op *absoluteJSR) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 1:
		op.readAddr = uint16(cpu.bus.Read(cpu.regs.PC))
		cpu.regs.PC++
	case 2:
		cpu.bus.Read(uint16(cpu.regs.SP))
	case 3:
		cpu.push(uint8((cpu.regs.PC) >> 8))
	case 4:
		cpu.push(uint8((cpu.regs.PC)))
	case 5:
		op.readAddr |= uint16(cpu.bus.Read(cpu.regs.PC)) << 8
		cpu.regs.PC = op.readAddr
		cpu.preFetch()
		done = true
	}
	op.t++
	return
}

func (op *absoluteJSR) setup(opCode uint8) {
	op.opCode = opCode
}

func (op *absoluteJSR) String() string {
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode, uint8(op.readAddr&0x0ff), uint8(op.readAddr>>8))
	writeOP(sb, "jsr $", toHex16(op.readAddr))
	return sb.String()
}

// -----

type rts struct {
	basicop
	targetAddr uint16
}

func (op *rts) Clone() operation {
	return &rts{basicop: op.basicop}
}

func (op *rts) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 1:
		cpu.bus.Read(cpu.regs.PC)
	case 2:
		cpu.bus.Read(uint16(cpu.regs.SP))
	case 3:
		op.targetAddr = uint16(cpu.pop())
	case 4:
		op.targetAddr |= uint16(cpu.pop()) << 8
	case 5:
		cpu.regs.PC = op.targetAddr + 1
		cpu.preFetch()
		done = true
	}
	op.t++
	return
}

func (op *rts) setup(opCode uint8) {
	op.opCode = opCode
}

func (op *rts) String() string {
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode)
	writeOP(sb, "rts")
	return sb.String()
}

// -----

type rti struct {
	basicop
}

func (op *rti) Clone() operation {
	return &rti{basicop: op.basicop}
}

func (op *rti) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 0:
		cpu.bus.Read(cpu.regs.PC)
	case 1:
		cpu.bus.Read(uint16(cpu.regs.SP))
	case 2:
		plp(cpu, cpu.pop())
	case 3:
		cpu.regs.PC = uint16(cpu.pop())
	case 4:
		cpu.regs.PC |= uint16(cpu.pop()) << 8
	case 5:
		cpu.onNMI = false
		cpu.preFetch()
		done = true
	}
	op.t++
	return
}

func (op *rti) setup(opCode uint8) {
	op.opCode = opCode
}

func (op *rti) String() string {
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode)
	writeOP(sb, "rti")
	return sb.String()
}

// -----

type push struct {
	basicop
	targetAddr uint16
	f          func(*m6502) byte
}

func (op *push) Clone() operation {
	return &push{basicop: op.basicop, f: op.f}
}

func (op *push) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 1:
		cpu.bus.Read(cpu.regs.PC)
	case 2:
		cpu.push(op.f(cpu))
		cpu.preFetch()
		done = true
	}
	op.t++
	return
}

func (op *push) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
}

func (op *push) String() string {
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode)
	writeOP(sb, op.ins)
	return sb.String()
}

// -----

type pull struct {
	basicop
	f func(*m6502, uint8)
}

func (op *pull) Clone() operation {
	return &pull{basicop: op.basicop, f: op.f}
}

func (op *pull) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 1:
		cpu.bus.Read(cpu.regs.PC)
	case 2:
		cpu.bus.Read(uint16(cpu.regs.SP))
	case 3:
		cpu.preFetch()
		op.f(cpu, cpu.pop())
		done = true
	}
	op.t++
	return
}

func (op *pull) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
}

func (op *pull) String() string {
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode)
	writeOP(sb, op.ins)
	return sb.String()
}

// -----

type indirectJMP struct {
	basicop
	readAddr uint16
}

func (op *indirectJMP) Clone() operation {
	return &indirectJMP{basicop: op.basicop}
}

func (op *indirectJMP) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 1:
		op.readAddr = uint16(cpu.bus.Read(cpu.regs.PC))
		cpu.regs.PC++
	case 2:
		op.readAddr |= uint16(cpu.bus.Read(cpu.regs.PC)) << 8
	case 3:
		cpu.regs.PC = uint16(cpu.bus.Read(op.readAddr))
	case 4:
		if op.readAddr&0x00ff == 0xff {
			cpu.regs.PC |= uint16(cpu.bus.Read(op.readAddr&0xff00)) << 8
		} else {
			cpu.regs.PC |= uint16(cpu.bus.Read(op.readAddr+1)) << 8
		}
		cpu.preFetch()
		done = true
	}
	op.t++
	return
}

func (op *indirectJMP) setup(opCode uint8) {
	op.opCode = opCode
}

func (op *indirectJMP) String() string {
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode, uint8(op.readAddr&0x0ff), uint8(op.readAddr>>8))
	writeOP(sb, "jmp ($", toHex16(op.readAddr), ")")
	return sb.String()
}

// -----

type absolute struct {
	basicop
	x, y       bool
	r, w, rw   bool
	f          interface{}
	readAddr   uint16
	targetAddr uint16
}

func (op *absolute) Clone() operation {
	return &absolute{basicop: op.basicop, f: op.f, x: op.x, y: op.y, r: op.r, w: op.w, rw: op.rw}
}

func (op *absolute) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 1:
		op.readAddr = uint16(cpu.bus.Read(cpu.regs.PC))
		cpu.regs.PC++
	case 2:
		op.readAddr |= uint16(cpu.bus.Read(cpu.regs.PC)) << 8
		op.targetAddr = op.readAddr
		cpu.regs.PC++
	case 3:
		if op.x {
			op.targetAddr = op.readAddr + uint16(cpu.regs.X)
		} else if op.y {
			op.targetAddr = op.readAddr + uint16(cpu.regs.Y)
		}
		if op.w || op.r {
			exec(cpu, op.f, op.targetAddr, op.r, op.w, op.rw)
			done = true
		}
	case 5:
		if op.y || op.x {
			if (op.targetAddr & 0xff00) == (op.readAddr & 0xff00) { // page change ?
				exec(cpu, op.f, op.targetAddr, op.r, op.w, op.rw)
				done = true
			}
		} else {
			exec(cpu, op.f, op.targetAddr, op.r, op.w, op.rw)
			done = true
		}
	case 6:
		exec(cpu, op.f, op.targetAddr, op.r, op.w, op.rw)
		done = true
	}
	op.t++
	return
}

func (op *absolute) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
	op.r, op.w, op.rw = getRWRW(op.f)
}

func (op *absolute) String() string {
	mod := ""
	if op.x {
		mod = ", X"
	} else if op.y {
		mod = ", Y"
	}
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode, uint8(op.readAddr&0x0ff), uint8(op.readAddr>>8))
	writeOP(sb, op.ins, " $", toHex16(op.readAddr), mod)
	return sb.String()
}

// -----

type zeropage struct {
	basicop
	x, y     bool
	r, w, rw bool
	f        interface{}
	addr     uint8
	tagert   uint8
}

func (op *zeropage) Clone() operation {
	return &zeropage{basicop: op.basicop, f: op.f, x: op.x, y: op.y, r: op.r, w: op.w, rw: op.rw}
}

func (op *zeropage) tick(cpu *m6502) (done bool) {
	switch op.t {
	case 1:
		op.addr = cpu.bus.Read(cpu.regs.PC)
		cpu.regs.PC++
	case 2:
		if op.x {
			op.tagert = op.addr + cpu.regs.X
		} else if op.y {
			op.tagert = op.addr + cpu.regs.Y
		} else {
			op.tagert = op.addr
		}
		if op.w {
			f := op.f.(func(*m6502) uint8)
			cpu.bus.Write(uint16(op.tagert), f(cpu))
			cpu.preFetch()
			done = true
		} else if op.r {
			f := op.f.(func(*m6502, uint8))
			f(cpu, cpu.bus.Read(uint16(op.tagert)))
			cpu.preFetch()
			done = true
		}
	case 3:
		if op.x {
			op.tagert = op.addr + cpu.regs.X
		} else if op.y {
			op.tagert = op.addr + cpu.regs.Y
		} else {
			op.tagert = op.addr
		}
	case 4:
		exec(cpu, op.f, uint16(op.tagert), op.r, op.w, op.rw)
		done = true
	}
	op.t++
	return
}

func (op *zeropage) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
	op.r, op.w, op.rw = getRWRW(op.f)
}

func (op zeropage) String() string {
	mod := ""
	if op.x {
		mod = ", X"
	} else if op.y {
		mod = ", Y"
	}
	sb := &strings.Builder{}
	writePC(sb, op.pc)
	writeMemory(sb, op.opCode, op.addr)
	writeOP(sb, op.ins, " $", toHex8(op.addr), mod)
	return sb.String()
}

func getRWRW(f interface{}) (r, w, rw bool) {
	switch f.(type) {
	case func(*m6502, uint8):
		r = true
	case func(*m6502) uint8:
		w = true
	case func(*m6502, uint8) uint8:
		rw = true
	default:
		panic(-1)
	}
	return
}

// -----

type unsupported struct {
	basicop
}

func (op *unsupported) Clone() operation {
	return &unsupported{}
}

func (op *unsupported) tick(cpu *m6502) (done bool) {
	panic(fmt.Sprintf("opCode: 0x%X NOT FOUND !!! pc:0x%04X", op.opCode, op.pc))
}

func (op *unsupported) setup(opCode uint8) {
	op.opCode = opCode
}

func (op *unsupported) String() string {
	op.tick(nil)
	return "---"
}

//-----------

func getFunctionName(f interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	idx := strings.LastIndex(name, ".") + 1
	return strings.ToUpper(name[idx : idx+3])
}

func exec(cpu *m6502, f interface{}, addr uint16, r, w, rw bool) {
	cpu.preFetch()
	if r {
		ff := f.(func(*m6502, uint8))
		// fmt.Printf("read  0x%04X 0x%02X\n", addr, cpu.bus.Read(addr))
		ff(cpu, cpu.bus.Read(addr))
	} else if w {
		ff := f.(func(*m6502) uint8)
		r := ff(cpu)
		// fmt.Printf("write 0x%04X 0x%02X\n", addr, r)
		cpu.bus.Write(addr, r)
	} else if rw {
		ff := f.(func(*m6502, uint8) uint8)
		// fmt.Printf("read 0x%04X 0x%02X\n", addr, cpu.bus.Read(addr))
		r := ff(cpu, cpu.bus.Read(addr))
		// fmt.Printf("write 0x%04X 0x%02X\n", addr, r)
		cpu.bus.Write(addr, r)
	} else {
		panic((-1))
	}
}

//-----------

func writePC(sb *strings.Builder, pc uint16) {
	sb.WriteString(strings.ToUpper(toHex16(pc)))
	sb.WriteString(": ")
}

func writeMemory(sb *strings.Builder, bytes ...uint8) {
	// for _, b := range bytes {
	// 	sb.WriteString(" ")
	// 	sb.WriteString(toHex8(b))
	// }
	// sb.WriteString("             "[:10-(len(bytes)*3)+1])
}

func writeOP(sb *strings.Builder, strs ...string) {
	for _, str := range strs {
		sb.WriteString(strings.ToLower(str))
	}
}

func toHex8(v uint8) string {
	n := "0" + strconv.FormatUint(uint64(v), 16)
	return n[len(n)-2:]
}

func toHex16(v uint16) string {
	n := "000" + strconv.FormatUint(uint64(v), 16)
	return n[len(n)-4:]
}
