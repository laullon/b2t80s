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
	setup(opCode uint8, ins string)
	setPC(pc uint16)
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

func (op *basicop) setup(opCode uint8, ins string) {
	op.ins = ins
	op.opCode = opCode
}

func (op *basicop) setPC(pc uint16) {
	op.pc = pc
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
		cpu.regs.PC = uint16(cpu.mem[0xff00+uint16(cpu.regs.SP)])<<8 | (cpu.regs.PC & 0x00ff)
		cpu.regs.SP--
	case 7:
		cpu.regs.PC = (cpu.regs.PC & 0xff00) | uint16(cpu.mem[0xff00+uint16(cpu.regs.SP)])
		op.d = true
	}
	op.t++
}

// -----
type indirect struct {
	basicop
	F    func(cpu *m6502, addr uint16)
	addr uint16
}

func (op *indirect) tick(cpu *m6502) {
	switch op.t {
	case 0:
		*cpu.AB.L = cpu.mem[cpu.regs.PC]
		cpu.regs.PC++
	case 1:
		*cpu.AB.H = cpu.mem[cpu.regs.PC]
		cpu.regs.PC++
	case 2:
		op.addr = uint16(cpu.mem[cpu.AB.Get()])
	case 3:
		op.addr |= uint16(cpu.mem[cpu.AB.Get()+1]) << 8
	case 4:
		op.F(cpu, op.addr)
		op.addr = cpu.AB.Get()
		op.d = true
	}
	op.t++
}

func (op *indirect) String() string {
	mod := ""
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X %02X %02X", op.opCode, op.addr&0x0ff, op.addr>>8),
		fmt.Sprintf("%s (0x%04X)%s", op.ins, op.addr, mod),
	)
}

// -----
type indirectXY struct {
	basicop
	x, y  bool
	F     func(cpu *m6502, addr uint16)
	addr  uint16
	addrZ uint8
}

func (op *indirectXY) tick(cpu *m6502) {
	switch op.t {
	case 0:
		op.addrZ = cpu.mem[cpu.regs.PC]
		if op.x {
			op.addrZ += cpu.regs.X
		}
		cpu.regs.PC++
	case 1:
		*cpu.AB.L = cpu.mem[op.addrZ]
	case 2:
		*cpu.AB.H = cpu.mem[op.addrZ+1]
	case 3:
		op.addr = cpu.AB.Get()
		if op.y {
			op.addr += uint16(cpu.regs.Y)
		}
	case 4:
		op.F(cpu, op.addr)
		op.d = true
	}
	op.t++
}

func (op *indirectXY) String() string {
	mod := ""
	if op.x {
		mod = "X"
	} else if op.y {
		mod = "Y"
	}
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X %02X", op.opCode, op.addrZ),
		fmt.Sprintf("%s (0x%02X), %s", op.ins, op.addrZ, mod),
	)
}

// -----
type implicit struct {
	basicop
	F func(cpu *m6502)
}

func (op *implicit) tick(cpu *m6502) {
	op.F(cpu)
	op.d = true
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
	F    func(cpu *m6502, data uint8)
	data uint8
}

func (op *immediate) tick(cpu *m6502) {
	op.data = cpu.mem[int(cpu.regs.PC)]
	op.F(cpu, op.data)
	cpu.regs.PC++
	op.d = true
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
	F    func(cpu *m6502, data int8)
	data int8
}

func (op *relative) tick(cpu *m6502) {
	op.data = int8(cpu.mem[cpu.regs.PC])
	cpu.regs.PC++
	op.F(cpu, op.data)
	op.d = true
}

func (op *relative) String() string {
	return fmt.Sprintf(debugFMT,
		op.pc,
		fmt.Sprintf("%02X %02X", op.opCode, uint8(op.data)),
		fmt.Sprintf("%s %d", op.ins, op.data),
	)
}

// -----

type absolute struct {
	basicop
	x, y bool
	F    func(cpu *m6502, addr uint16)
	addr uint16
}

func (op *absolute) tick(cpu *m6502) {
	switch op.t {
	case 0:
		*cpu.AB.L = cpu.mem[cpu.regs.PC]
		cpu.regs.PC++
	case 1:
		*cpu.AB.H = cpu.mem[cpu.regs.PC]
		cpu.regs.PC++
	case 2:
		addr := cpu.AB.Get()
		op.addr = addr
		if op.x {
			addr += uint16(cpu.regs.X)
		} else if op.y {
			addr += uint16(cpu.regs.Y)
		}
		op.F(cpu, addr)
		op.d = true
	}
	op.t++
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
		fmt.Sprintf("%02X %02X %02X", op.opCode, op.addr&0x0ff, op.addr>>8),
		fmt.Sprintf("%s 0x%04X%s", op.ins, op.addr, mod),
	)
}

// -----

type zeropage struct {
	basicop
	x, y bool
	F    func(cpu *m6502, addr uint16)
	addr uint8
}

func (op *zeropage) tick(cpu *m6502) {
	switch op.t {
	case 0:
		op.addr = cpu.mem[cpu.regs.PC]
		cpu.regs.PC++
	case 1:
		if op.x {
			op.addr += cpu.regs.X
		} else if op.y {
			op.addr += cpu.regs.Y
		}
		op.F(cpu, uint16(op.addr))
		op.d = true
	}
	op.t++
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

func getFunctionName(op interface{}) string {
	r := reflect.ValueOf(op)
	f := reflect.Indirect(r).FieldByName("F")
	if !f.IsValid() {
		return "Error"
	}
	name := runtime.FuncForPC(f.Pointer()).Name()
	idx := strings.LastIndex(name, ".") + 1
	return strings.ToUpper(name[idx : idx+3])
}
