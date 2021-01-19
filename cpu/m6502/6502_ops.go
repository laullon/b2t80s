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

func (op *basicop) setup(opCode uint8, ins string) {
	op.ins = ins
	op.opCode = opCode
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
		cpu.regs.PC = uint16(cpu.mem[0xff00+uint16(cpu.regs.SP)])<<8 | (cpu.regs.PC & 0x00ff)
		cpu.regs.SP--
	case 7:
		cpu.regs.PC = (cpu.regs.PC & 0xff00) | uint16(cpu.mem[0xff00+uint16(cpu.regs.SP)])
		op.d = true
	}
	op.t++
}

func (op *reset) String() string {
	return fmt.Sprintf(debugFMT,
		op.pc,
		"",
		"RESET",
	)
}

// -----
type indirect struct {
	basicop
	F      func(cpu *m6502, addr uint16)
	addr   uint16
	target uint16
}

func (op *indirect) tick(cpu *m6502) {
	switch op.t {
	case 0:
		op.addr = uint16(cpu.mem[cpu.regs.PC])
		cpu.regs.PC++
	case 1:
		op.addr |= uint16(cpu.mem[cpu.regs.PC]) << 8
		cpu.regs.PC++
	case 2:
		op.target = uint16(cpu.mem[op.addr])
	case 3:
		op.target |= uint16(cpu.mem[op.addr+1]) << 8
	case 4:
		op.F(cpu, op.target)
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
	F     func(cpu *m6502, data uint8) (discard bool, v uint8)
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
		op.addr = uint16(cpu.mem[op.addrZ])
	case 2:
		op.addr |= uint16(cpu.mem[op.addrZ+1]) << 8
	case 3:
		if op.y {
			op.addr += uint16(cpu.regs.Y)
		}
	case 4:
		discard, v := op.F(cpu, cpu.mem[op.addr])
		if !discard {
			cpu.mem[op.addr] = v
		}
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
	F      func(cpu *m6502) bool
	off    int8
	target uint16
}

func (op *relative) tick(cpu *m6502) {
	switch op.t {
	case 0:
		op.d = !op.F(cpu)
		op.off = int8(cpu.mem[cpu.regs.PC])
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
		op.readAddr = uint16(cpu.mem[cpu.regs.PC])
		cpu.regs.PC++
	case 1:
		op.readAddr |= uint16(cpu.mem[cpu.regs.PC]) << 8
		cpu.regs.PC = op.readAddr
		op.d = true
	}
	op.t++
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
		op.readAddr = uint16(cpu.mem[cpu.regs.PC])
		cpu.regs.PC++
	case 1:
		op.readAddr |= uint16(cpu.mem[cpu.regs.PC]) << 8
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
		op.readAddr = uint16(cpu.mem[cpu.regs.PC])
		cpu.regs.PC++
	case 1:
		op.readAddr |= uint16(cpu.mem[cpu.regs.PC]) << 8
	case 2:
		cpu.regs.PC = uint16(cpu.mem[op.readAddr])
	case 3:
		cpu.regs.PC |= uint16(cpu.mem[op.readAddr+1]) << 8
		op.d = true
	}
	op.t++
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
	F          func(cpu *m6502, data uint8) (discard bool, v uint8)
	readAddr   uint16
	targetAddr uint16
}

func (op *absolute) tick(cpu *m6502) {
	switch op.t {
	case 0:
		op.readAddr = uint16(cpu.mem[cpu.regs.PC])
		cpu.regs.PC++
	case 1:
		op.readAddr |= uint16(cpu.mem[cpu.regs.PC]) << 8
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
			op.exec(cpu)
		}
	case 3:
		op.exec(cpu)
	}
	op.t++
}

func (op *absolute) exec(cpu *m6502) {
	discard, v := op.F(cpu, cpu.mem[op.targetAddr])
	if !discard {
		cpu.mem[op.targetAddr] = v
	}
	op.d = true
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
	x, y bool
	F    func(cpu *m6502, data uint8) (discard bool, v uint8)
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
	case 2:
		discard, v := op.F(cpu, cpu.mem[op.addr])
		if !discard {
			cpu.mem[op.addr] = v
		}
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