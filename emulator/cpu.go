package emulator

import (
	"fmt"
	"strings"
)

type CPUTrap func()
type CPU interface {
	Interrupt(bool)

	Registers() interface{}

	SetDebuger(debugger Debugger)

	RegisterTrap(pc uint16, trap CPUTrap)

	// LoadTapeBlock() uint16
	// LoadTapeBlockCPC(uint16) uint16

	Halt()

	Tick()
}

type Debugger interface {
	LoadSymbols(file string)
	AddLastInstruction(ins Instruction)

	GetLog() string
	GetNextInstruction() string
	GetFollowingInstruction() string
	GetRegisters() string

	SetStatus(sts string)
	GetStatus() string

	SetDump(bool)

	Stop()
	StopNextFrame()
	Continue()
	Step()
	DumpNextFrame()

	NextFrame()

	IsStoped() bool

	SetBreakPoint(bp uint16)
}

type PortMask struct {
	Mask  uint16
	Value uint16
}

type PortManager interface {
	ReadPort(port uint16) (byte, bool)
	WritePort(port uint16, data byte)
}

type Instruction struct {
	Instruction int32
	Opcode      string
	Tstates     uint
	Length      uint16
	Mem         []byte
	Valid       bool
}

func (ins *Instruction) Dump(pc uint16) string {
	if ins == nil {
		return "<nil>"
	}

	op := ins.Opcode
	if strings.Index(op, "$NN") != -1 {
		nn := fmt.Sprintf("0x%02X%02X", ins.Mem[2], ins.Mem[1])
		op = strings.Replace(op, "$NN", nn, -1)
	}

	if strings.Index(op, "+N)") != -1 {
		n := fmt.Sprintf("+%d)", ins.Mem[2])
		op = strings.Replace(op, "+N)", n, -1)
	}

	if strings.Index(op, "$N+2") != -1 {
		jump := int8(ins.Mem[1]) + 2
		jpc := uint16(int16(pc) + int16(jump))
		n := fmt.Sprintf("0x%04X", jpc)
		// if symb, ok := symbols[pc]; ok {
		// 	n = symb
		// }
		op = strings.Replace(op, "$N+2", n, -1)
	}

	op = strings.ReplaceAll(op, ",", ", ")

	var m []string
	for _, b := range ins.Mem {
		m = append(m, fmt.Sprintf("%02X", b))
	}

	return fmt.Sprintf("0x%04x %-14.14s %-14.14v", pc, op, strings.Join(m, " "))
}
