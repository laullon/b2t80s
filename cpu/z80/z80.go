package z80

import (
	"fmt"
	"strings"

	cpuUtils "github.com/laullon/b2t80s/cpu"
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
	BC *cpuUtils.RegPair

	D  byte
	E  byte
	DE *cpuUtils.RegPair

	H  byte
	L  byte
	HL *cpuUtils.RegPair

	S  byte
	P  byte
	SP *cpuUtils.RegPair

	I  byte
	R  byte
	R7 byte

	IFF1 bool
	IFF2 bool

	IXH byte
	IXL byte
	IX  *cpuUtils.RegPair

	IYH byte
	IYL byte
	IY  *cpuUtils.RegPair

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

type fetchedData struct {
	pc     uint16
	prefix uint16
	opCode uint8
	n      uint8
	n2     uint8
	nn     uint16
	op     *opCode
}

func (d *fetchedData) getInstruction() string {
	return d.op.String()
}

func (d *fetchedData) getMemory() string {
	var res strings.Builder
	if d.prefix > 0xff && d.op.len == 4 {
		res.WriteString(fmt.Sprintf("%02X ", (d.prefix >> 8)))
		res.WriteString(fmt.Sprintf("%02X ", (d.prefix & 0xff)))
		res.WriteString(fmt.Sprintf("%02X ", d.n))
		res.WriteString(fmt.Sprintf("%02X", d.opCode))
	} else if d.prefix > 0x00 && d.op.len == 4 {
		res.WriteString(fmt.Sprintf("%02X ", (d.prefix & 0xff)))
		res.WriteString(fmt.Sprintf("%02X ", d.opCode))
		res.WriteString(fmt.Sprintf("%02X ", d.n))
		res.WriteString(fmt.Sprintf("%02X", d.n2))
	} else if d.prefix > 0x00 && d.op.len == 3 {
		res.WriteString(fmt.Sprintf("%02X ", (d.prefix & 0xff)))
		res.WriteString(fmt.Sprintf("%02X ", d.opCode))
		res.WriteString(fmt.Sprintf("%02X", d.n))
	} else if d.prefix == 0x00 && d.op.len == 3 {
		res.WriteString(fmt.Sprintf("%02X ", d.opCode))
		res.WriteString(fmt.Sprintf("%02X ", d.n))
		res.WriteString(fmt.Sprintf("%02X", d.n2))
	} else if d.prefix == 0x00 && d.op.len == 2 {
		res.WriteString(fmt.Sprintf("%02X ", d.opCode))
		res.WriteString(fmt.Sprintf("%02X", d.n))
	} else {
		res.WriteString(fmt.Sprintf("%02X", d.opCode))
	}
	return res.String()
}

type z80 struct {
	debugger emulator.Debugger

	bus emulator.Bus

	regs      *Z80Registers
	indexRegs []*cpuUtils.RegPair
	indexIdx  int

	halt bool
	wait bool

	doInterrupt bool

	fetched   *fetchedData
	scheduler *circularBuffer

	traps map[uint16]emulator.CPUTrap
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
	cpu := &z80{
		bus:       bus,
		fetched:   &fetchedData{},
		scheduler: newCircularBuffer(),
		traps:     make(map[uint16]emulator.CPUTrap),
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

	cpu.regs.BC = &cpuUtils.RegPair{&cpu.regs.B, &cpu.regs.C}
	cpu.regs.DE = &cpuUtils.RegPair{&cpu.regs.D, &cpu.regs.E}
	cpu.regs.HL = &cpuUtils.RegPair{&cpu.regs.H, &cpu.regs.L}
	cpu.regs.SP = &cpuUtils.RegPair{&cpu.regs.S, &cpu.regs.P}
	cpu.regs.IX = &cpuUtils.RegPair{&cpu.regs.IXH, &cpu.regs.IXL}
	cpu.regs.IY = &cpuUtils.RegPair{&cpu.regs.IYH, &cpu.regs.IYL}
	cpu.indexRegs = []*cpuUtils.RegPair{cpu.regs.HL, cpu.regs.IX, cpu.regs.IY}

	cpu.scheduler.append(newFetch(lookup))
	return cpu
}

func (cpu *z80) CurrentOP() string { panic(-2) }

func (cpu *z80) SetDebuger(debugger emulator.Debugger) {
	cpu.debugger = debugger
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

func (cpu *z80) Wait(w bool) {
	cpu.wait = w
}

func (cpu *z80) Halt() {
	cpu.halt = true
}

func (cpu *z80) execInterrupt() {
	cpu.prepareForNewInstruction()
	cpu.doInterrupt = false

	if cpu.regs.IFF1 {
		cpu.regs.IFF1 = false
		cpu.regs.IFF2 = false
		switch cpu.regs.InterruptsMode {
		case 0, 1:
			code := &exec{l: 7, f: func(cpu *z80) {
				cpu.pushToStack(cpu.regs.PC, func(cpu *z80) {
					cpu.regs.PC = 0x0038
				})
			}}
			cpu.scheduler.append(code)
		case 2:
			code := &exec{l: 7, f: func(cpu *z80) {
				cpu.pushToStack(cpu.regs.PC, func(cpu *z80) {
					pos := uint16(cpu.regs.I)<<8 + 0xff
					mr1 := newMR(pos, func(cpu *z80, data byte) {
						cpu.regs.PC = uint16(data) << 8
					})
					mr2 := newMR(pos, func(cpu *z80, data byte) {
						cpu.regs.PC |= uint16(data)
					})
					cpu.scheduler.append(mr1, mr2)
				})
			}}
			cpu.scheduler.append(code)
		}
	}
}

func (cpu *z80) prepareForNewInstruction() {
	if cpu.debugger != nil { // TODO: add dummy debuger on the test to remove this IF
		cpu.debugger.AddInstruction(cpu.fetched.pc, cpu.fetched.getMemory(), cpu.fetched.getInstruction())
	}

	cpu.fetched.n = 0
	cpu.fetched.nn = 0
	cpu.fetched.opCode = 0
	cpu.fetched.prefix = 0
	cpu.fetched.pc = cpu.regs.PC
	cpu.indexIdx = 0
}

func (cpu *z80) Tick() {
	if cpu.wait {
		return
	}

	if cpu.halt {
		if cpu.doInterrupt {
			cpu.halt = false
			cpu.regs.PC++
			cpu.execInterrupt()
		} else {
			return
		}
	}

	if cpu.scheduler.first().isDone() {
		cpu.scheduler.next()
		if cpu.scheduler.isEmpty() {
			if cpu.doInterrupt {
				cpu.execInterrupt()
			} else {
				cpu.newInstruction()
			}
		}
	}
	cpu.scheduler.first().tick(cpu)
	// if cpu.scheduler.first().isDone() {
	// 	println("[done]", reflect.TypeOf(cpu.scheduler.first()).String())
	// }
}

func (cpu *z80) newInstruction() {
	cpu.prepareForNewInstruction()
	cpu.doTraps()
	cpu.scheduler.append(newFetch(lookup))
}

func (cpu *z80) doTraps() {
	if trap, ok := cpu.traps[cpu.regs.PC]; ok {
		trap()
	}
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
