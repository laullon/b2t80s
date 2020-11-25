package z80

import (
	"github.com/laullon/b2t80s/emulator"
)

var overflowAddTable = []bool{false, false, false, true, true, false, false, false}
var overflowSubTable = []bool{false, true, false, false, false, false, true, false}
var halfcarryAddTable = []bool{false, true, true, true, false, false, false, true}
var halfcarrySubTable = []bool{false, false, true, false, true, false, true, true}

var parityTable = make([]bool, 0x100)

type Z80Registers struct {
	PC uint16
	M1 bool

	A byte
	F *flags

	B  byte
	C  byte
	BC *RegPair

	D  byte
	E  byte
	DE *RegPair

	H  byte
	L  byte
	HL *RegPair

	S  byte
	P  byte
	SP *RegPair

	I  byte
	R  byte
	R7 byte

	IFF1 bool
	IFF2 bool

	IXH byte
	IXL byte
	IX  *RegPair

	IYH byte
	IYL byte
	IY  *RegPair

	Aalt byte
	Falt *flags
	Balt byte
	Calt byte
	Dalt byte
	Ealt byte
	Halt byte
	Lalt byte

	InterruptsMode byte
}

type z80 struct {
	memory   emulator.Memory
	debugger emulator.Debugger

	Bus emulator.Bus

	regs *Z80Registers

	halt, haltDone bool

	doInterrupt bool

	fetched   []uint8
	scheduler []z80op

	traps map[uint16]emulator.CPUTrap
	ports map[emulator.PortMask]emulator.PortManager
}

func init() {
	var i int16
	var j, k byte
	var p byte

	for i = 0; i < 0x100; i++ {
		j = byte(i)
		p = 0
		for k = 0; k < 8; k++ {
			p ^= j & 1
			j >>= 1
		}
		if p != 0 {
			parityTable[i] = false
		} else {
			parityTable[i] = true
		}
	}

}

func NewZ80(bus emulator.Bus) emulator.CPU {
	LoadOPCodess()
	cpu := &z80{
		Bus:   bus,
		traps: make(map[uint16]emulator.CPUTrap),
		ports: make(map[emulator.PortMask]emulator.PortManager),
		regs: &Z80Registers{
			PC: 0,
			M1: false,
			A:  0xff,
			S:  0xFF,
			P:  0xFF,
			F: &flags{
				Z: true,
				C: true,
				S: true,
				H: true,
				P: true,
				N: true,
			},
			Falt: &flags{
				Z: true,
				C: true,
				S: true,
				H: true,
				P: true,
				N: true,
			},
			R: 0x01,
		},
	}

	cpu.regs.BC = &RegPair{&cpu.regs.B, &cpu.regs.C}
	cpu.regs.DE = &RegPair{&cpu.regs.D, &cpu.regs.E}
	cpu.regs.HL = &RegPair{&cpu.regs.H, &cpu.regs.L}
	cpu.regs.SP = &RegPair{&cpu.regs.S, &cpu.regs.P}
	cpu.regs.IX = &RegPair{&cpu.regs.IXH, &cpu.regs.IXL}
	cpu.regs.IY = &RegPair{&cpu.regs.IYH, &cpu.regs.IYL}

	cpu.scheduler = append(cpu.scheduler, &fetch{})
	return cpu
}

func (cpu *z80) SetDebuger(debugger emulator.Debugger) {
	cpu.debugger = debugger
}

func (cpu *z80) RegisterPort(mask emulator.PortMask, manager emulator.PortManager) {
	cpu.ports[mask] = manager
}

func (cpu *z80) RegisterTrap(pc uint16, trap emulator.CPUTrap) {
	cpu.traps[pc] = trap
}

func (cpu *z80) Registers() interface{} {
	return cpu.regs
}

func (cpu *z80) Interrupt(i bool) {
	cpu.doInterrupt = i
}

func (cpu *z80) Halt() {
	cpu.halt = true
}

func (cpu *z80) execInterrupt() uint {
	cpu.doInterrupt = false
	if cpu.halt {
		cpu.haltDone = true
		cpu.halt = false
	}
	var ts uint
	if cpu.regs.IFF1 {
		cpu.regs.IFF1 = false
		cpu.regs.IFF2 = false

		// cpu.regs.SP.Push(cpu.regs.PC)

		switch cpu.regs.InterruptsMode {
		case 0, 1:
			ts = 13
			cpu.regs.PC = 0x38
		case 2:
			ts = 19
			pos := uint16(cpu.regs.I)<<8 + 0xff
			cpu.regs.PC = getWord(cpu.memory, pos)
		}
	}
	return ts
}

func (cpu *z80) Tick() {
	if cpu.scheduler[0].isDone() {
		cpu.scheduler = cpu.scheduler[1:]
	}
	cpu.scheduler[0].tick(cpu)
}

// func (cpu *z80) Tick() {
// 	var err error
// 	if cpu.regs.M1 {

// 		// TODO: review
// 		if trap, ok := cpu.traps[cpu.regs.PC]; ok {
// 			res := trap()
// 			switch res {
// 			case emulator.CONTINUE:
// 			case emulator.STOP:
// 				return
// 			default:
// 				cpu.regs.PC = uint16(res)
// 				return
// 			}
// 		}

// 		cpu.regs.M1 = false

// 		if cpu.doInterrupt {
// 			cpu.pendingTicks = cpu.execInterrupt()
// 			cpu.instruction = emulator.Instruction{Instruction: 0xffffff}
// 		}

// 		if cpu.halt { //TODO review
// 			cpu.pendingTicks = 4
// 			cpu.instruction = emulator.Instruction{Instruction: 0xffffff}
// 		}

// 		cpu.instruction, err = GetOpCode(cpu.memory.GetBlock(cpu.regs.PC, 4))
// 		cpu.pendingTicks = cpu.instruction.Tstates
// 		if err != nil {
// 			panic(err)
// 		}

// 		if cpu.debugger != nil {
// 			cpu.debugger.AddLastInstruction(cpu.instruction)
// 		}

// 	}

// 	cpu.pendingTicks--

// 	if cpu.pendingTicks == 0 {
// 		cpu.regs.M1 = true
// 		needPcUpdate := cpu.runSwitch(cpu.instruction)
// 		if needPcUpdate {
// 			cpu.regs.PC += uint16(cpu.instruction.Length)
// 		}
// 	}
// }

// func (cpu *z80) Step() {
// 	if trap, ok := cpu.traps[cpu.regs.PC]; ok {
// 		res := trap()
// 		switch res {
// 		case emulator.CONTINUE:
// 		case emulator.STOP:
// 			return
// 		default:
// 			cpu.regs.PC = uint16(res)
// 			return
// 		}
// 	}

// 	if cpu.doInterrupt {
// 		ts := cpu.execInterrupt()
// 		cpu.clock.AddTStates(ts)
// 	}

// 	if cpu.debugger != nil {
// 		cpu.debugger.AddLastInstruction(ins)
// 	}

// 	var ts uint
// 	if !cpu.halt {
// 		needPcUpdate := cpu.runSwitch(ins)
// 		if needPcUpdate {
// 			cpu.regs.PC += uint16(ins.Length)
// 		}
// 		ts = ins.Tstates
// 	} else {
// 		ts = 4 // halt
// 	}

// 	cpu.clock.AddTStates(ts)
// 	return
// }
