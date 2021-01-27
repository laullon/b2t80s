package m6502

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

var debugFMT = "0x%04X: %-10s %-20s"

type operation interface {
	done() bool
	tick(cpu *m6502)
	reset()
	setup(opCode uint8)
	setPC(pc uint16)
	getPC() uint16
}

type basicop struct {
	ins    string
	opCode uint8

	pc uint16
	d  bool
	t  uint
}

func (op *basicop) reset() {
	op.d = false
	op.t = 0
}

func (op *basicop) done() bool {
	return op.d
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

func (op *reset) tick(cpu *m6502) {
	switch op.t {
	case 0:
		cpu.regs.SP = 0
	case 3, 4, 5:
		cpu.regs.SP--
	case 6:
		cpu.regs.PC = uint16(cpu.bus.Read(0xff00+uint16(cpu.regs.SP)))<<8 | (cpu.regs.PC & 0x00ff)
		cpu.regs.SP--
	case 7:
		cpu.regs.PC = (cpu.regs.PC & 0xff00) | uint16(cpu.bus.Read(0xff00+uint16(cpu.regs.SP)))
		op.d = true
	}
	op.t++
}

func (op *reset) setup(opCode uint8) {
	op.opCode = opCode
}

func (op *reset) String() string {
	return fmt.Sprintf(debugFMT,
		op.pc,
		"",
		"RESET",
	)
}

// -----
type brk struct {
	basicop
	vector   uint16
	irq, imm bool
}

func (op *brk) tick(cpu *m6502) {
	switch op.t {
	case 0:
		if !op.imm && !op.irq {
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
		cpu.push(cpu.regs.PS.get())
		cpu.regs.PS.X = false
	case 4:
		cpu.regs.PC = uint16(cpu.bus.Read(op.vector))
		cpu.regs.PS.I = true
	case 5:
		cpu.regs.PC |= uint16(cpu.bus.Read(op.vector+1)) << 8
		op.d = true
	}
	op.t++
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
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X", op.opCode),
		fmt.Sprintf("BRK %5s", mod),
	)
}

// -----
type indirectX struct {
	basicop
	r, w, rw bool
	f        interface{}
	addr     uint16
	addrZ    uint8
}

func (op *indirectX) tick(cpu *m6502) {
	switch op.t {
	case 0:
		op.addrZ = cpu.bus.Read(cpu.regs.PC)
		op.addrZ += cpu.regs.X
		cpu.regs.PC++
	case 1:
		op.addr = uint16(cpu.bus.Read(uint16(op.addrZ)))
	case 2:
		op.addr |= uint16(cpu.bus.Read(uint16(op.addrZ+1))) << 8
	case 3:
		exec(cpu, op.f, op.addr, op.r, op.w, op.rw)
		op.d = true
	}
	op.t++
}

func (op *indirectX) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
	op.r, op.w, op.rw = getRWRW(op.f)
}

func (op *indirectX) String() string {
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X %02X", op.opCode, op.addrZ),
		fmt.Sprintf("%s (0x%02X, X)", op.ins, op.addrZ),
	)
}

// -----
type indirectY struct {
	basicop
	r, w, rw bool
	f        interface{}
	addr     uint16
	addrZ    uint8
}

func (op *indirectY) tick(cpu *m6502) {
	switch op.t {
	case 0:
		op.addrZ = cpu.bus.Read(cpu.regs.PC)
		cpu.regs.PC++
	case 1:
		op.addr = uint16(cpu.bus.Read(uint16(op.addrZ)))
	case 2:
		op.addr |= uint16(cpu.bus.Read(uint16(op.addrZ+1))) << 8
	case 3:
		op.addr += uint16(cpu.regs.Y)
	case 4:
		exec(cpu, op.f, op.addr, op.r, op.w, op.rw)
		op.d = true
	}
	op.t++
}

func (op *indirectY) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
	op.r, op.w, op.rw = getRWRW(op.f)
}

func (op *indirectY) String() string {
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X %02X", op.opCode, op.addrZ),
		fmt.Sprintf("%s (0x%02X), Y", op.ins, op.addrZ),
	)
}

// -----
type implicit struct {
	basicop
	f func(cpu *m6502)
}

func (op *implicit) tick(cpu *m6502) {
	op.f(cpu)
	op.d = true
}

func (op *implicit) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
}

func (op *implicit) String() string {
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X", op.opCode),
		fmt.Sprintf("%s", op.ins),
	)
}

// -----

type immediate struct {
	basicop
	f    func(cpu *m6502, data uint8)
	data uint8
}

func (op *immediate) tick(cpu *m6502) {
	op.data = cpu.bus.Read(cpu.regs.PC)
	op.f(cpu, op.data)
	cpu.regs.PC++
	op.d = true
}

func (op *immediate) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
}

func (op *immediate) String() string {
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X %02X", op.opCode, op.data),
		fmt.Sprintf("%s 0x%02X", op.ins, op.data),
	)
}

// -----

type relative struct {
	basicop
	f      func(cpu *m6502) bool
	off    int8
	target uint16
}

func (op *relative) tick(cpu *m6502) {
	switch op.t {
	case 0:
		op.d = !op.f(cpu)
		op.off = int8(cpu.bus.Read(cpu.regs.PC))
		cpu.regs.PC++
	case 1:
		op.target = cpu.regs.PC + uint16(op.off)
		if (cpu.regs.PC & 0xff00) == (op.target & 0xff00) { // no page change
			cpu.regs.PC = op.target
			op.d = true
		}
	case 2:
		cpu.regs.PC = op.target
		op.d = true
	}
	op.t++
}

func (op *relative) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
}

func (op *relative) String() string {
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X %02X", op.opCode, uint8(op.off)),
		fmt.Sprintf("%s %d", op.ins, op.off),
	)
}

// -----

type absoluteJMP struct {
	basicop
	readAddr uint16
}

func (op *absoluteJMP) tick(cpu *m6502) {
	switch op.t {
	case 0:
		op.readAddr = uint16(cpu.bus.Read(cpu.regs.PC))
		cpu.regs.PC++
	case 1:
		op.readAddr |= uint16(cpu.bus.Read(cpu.regs.PC)) << 8
		cpu.regs.PC = op.readAddr
		op.d = true
	}
	op.t++
}

func (op *absoluteJMP) setup(opCode uint8) {
	op.opCode = opCode
}

func (op *absoluteJMP) String() string {
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X %02X %02X", op.opCode, op.readAddr&0x0ff, op.readAddr>>8),
		fmt.Sprintf("jmp 0x%04X", op.readAddr),
	)
}

// -----

type absoluteJSR struct {
	basicop
	readAddr uint16
}

func (op *absoluteJSR) tick(cpu *m6502) {
	switch op.t {
	case 0:
		op.readAddr = uint16(cpu.bus.Read(cpu.regs.PC))
		cpu.regs.PC++
	case 1:
		op.readAddr |= uint16(cpu.bus.Read(cpu.regs.PC)) << 8
		cpu.regs.PC++
	case 2:
		cpu.push(uint8((cpu.regs.PC - 1) >> 8))
	case 3:
		cpu.push(uint8((cpu.regs.PC - 1)))
		cpu.regs.PC = op.readAddr
		op.d = true
	}
	op.t++
}

func (op *absoluteJSR) setup(opCode uint8) {
	op.opCode = opCode
}

func (op *absoluteJSR) String() string {
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X %02X %02X", op.opCode, op.readAddr&0x0ff, op.readAddr>>8),
		fmt.Sprintf("jsr 0x%04X", op.readAddr),
	)
}

// -----

type indirectJMP struct {
	basicop
	readAddr uint16
}

func (op *indirectJMP) tick(cpu *m6502) {
	switch op.t {
	case 0:
		op.readAddr = uint16(cpu.bus.Read(cpu.regs.PC))
		cpu.regs.PC++
	case 1:
		op.readAddr |= uint16(cpu.bus.Read(cpu.regs.PC)) << 8
	case 2:
		cpu.regs.PC = uint16(cpu.bus.Read(op.readAddr))
	case 3:
		cpu.regs.PC |= uint16(cpu.bus.Read(op.readAddr+1)) << 8
		op.d = true
	}
	op.t++
}

func (op *indirectJMP) setup(opCode uint8) {
	op.opCode = opCode
}

func (op *indirectJMP) String() string {
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X %02X %02X", op.opCode, op.readAddr&0x0ff, op.readAddr>>8),
		fmt.Sprintf("jmp 0x%04X", op.readAddr),
	)
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

func (op *absolute) tick(cpu *m6502) {
	switch op.t {
	case 0:
		op.readAddr = uint16(cpu.bus.Read(cpu.regs.PC))
		cpu.regs.PC++
	case 1:
		op.readAddr |= uint16(cpu.bus.Read(cpu.regs.PC)) << 8
		cpu.regs.PC++
	case 2:
		if op.x {
			op.targetAddr = op.readAddr + uint16(cpu.regs.X)
		} else if op.y {
			op.targetAddr = op.readAddr + uint16(cpu.regs.Y)
		} else {
			op.targetAddr = op.readAddr
		}
		if (op.targetAddr & 0xff00) == (op.readAddr & 0xff00) { // page change ?
			exec(cpu, op.f, op.targetAddr, op.r, op.w, op.rw)
			op.d = true
		}
	case 3:
		exec(cpu, op.f, op.targetAddr, op.r, op.w, op.rw)
		op.d = true
	}
	op.t++
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
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X %02X %02X", op.opCode, op.readAddr&0x0ff, op.readAddr>>8),
		fmt.Sprintf("%s 0x%04X%s", op.ins, op.readAddr, mod),
	)
}

// -----

type zeropage struct {
	basicop
	x, y     bool
	r, w, rw bool
	f        interface{}
	addr     uint8
}

func (op *zeropage) tick(cpu *m6502) {
	switch op.t {
	case 0:
		op.addr = cpu.bus.Read(cpu.regs.PC)
		cpu.regs.PC++
	case 1:
		if op.x {
			op.addr += cpu.regs.X
		} else if op.y {
			op.addr += cpu.regs.Y
		}
	case 2:
		exec(cpu, op.f, uint16(op.addr), op.r, op.w, op.rw)
		op.d = true
	}
	op.t++
}

func (op *zeropage) setup(opCode uint8) {
	op.opCode = opCode
	op.ins = getFunctionName(op.f)
	op.r, op.w, op.rw = getRWRW(op.f)
}

func (op *zeropage) String() string {
	mod := ""
	if op.x {
		mod = ", X"
	} else if op.y {
		mod = ", Y"
	}
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X %02X", op.opCode, op.addr),
		fmt.Sprintf("%s 0x%02X%s", op.ins, op.addr, mod),
	)
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

func getFunctionName(f interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	idx := strings.LastIndex(name, ".") + 1
	return strings.ToUpper(name[idx : idx+3])
}

func exec(cpu *m6502, f interface{}, addr uint16, r, w, rw bool) {
	if r {
		ff := f.(func(*m6502, uint8))
		fmt.Printf("read  0x%04X 0x%02X\n", addr, cpu.bus.Read(addr))
		ff(cpu, cpu.bus.Read(addr))
	} else if w {
		ff := f.(func(*m6502) uint8)
		r := ff(cpu)
		fmt.Printf("write 0x%04X 0x%02X\n", addr, r)
		cpu.bus.Write(addr, r)
	} else if rw {
		ff := f.(func(*m6502, uint8) uint8)
		fmt.Printf("read 0x%04X 0x%02X\n", addr, cpu.bus.Read(addr))
		r := ff(cpu, cpu.bus.Read(addr))
		fmt.Printf("write 0x%04X 0x%02X\n", addr, r)
		cpu.bus.Write(addr, r)
	} else {
		panic((-1))
	}
}
