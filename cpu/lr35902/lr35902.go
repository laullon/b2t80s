package lr35902

import (
	"fmt"
	"strings"

	"github.com/laullon/b2t80s/cpu"
)

var overflowAddTable = []bool{false, false, false, true, true, false, false, false}
var overflowSubTable = []bool{false, true, false, false, false, false, true, false}
var halfcarryAddTable = []bool{false, true, true, true, false, false, false, true}
var halfcarrySubTable = []bool{false, false, true, false, true, false, true, true}

var parityTable = make([]bool, 0x100)

type LR35902Registers struct {
	PC uint16
	M1 bool

	A byte
	F *flags

	B  byte
	C  byte
	BC *cpu.RegPair

	D  byte
	E  byte
	DE *cpu.RegPair

	H  byte
	L  byte
	HL *cpu.RegPair

	S  byte
	P  byte
	SP *cpu.RegPair

	I  byte
	R  byte
	R7 byte

	IFF1 bool
	IFF2 bool

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

func (d *fetchedData) String() string {
	return fmt.Sprintf("0x%04X: %-11s : %s", d.pc, d.getMemory(), d.op)
}

type CPUTrap func()

type LR35902 interface {
	cpu.CPU
	Registers() *LR35902Registers
	RegisterTrap(pc uint16, trap CPUTrap)
}

type lr35902 struct {
	bus *genericBus

	regs      *LR35902Registers
	indexRegs []*cpu.RegPair
	indexIdx  int

	halt bool
	wait bool

	doInterrupt bool

	fetched   *fetchedData
	scheduler *circularBuffer

	traps map[uint16]CPUTrap

	log      cpu.CPUTracer
	debugger cpu.DebuggerCallbacks
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

func New(bus cpu.Bus) LR35902 {
	res := &lr35902{
		bus:       newBus(bus),
		fetched:   &fetchedData{},
		scheduler: newCircularBuffer(),
		traps:     make(map[uint16]CPUTrap),
		regs: &LR35902Registers{
			PC: 0,
			M1: false,
			A:  0xff,
			S:  0xFF,
			P:  0xFF,
			F: &flags{
				Z: true,
				C: true,
				H: true,
				N: true,
			},
			R: 0x01,
		},
	}

	res.regs.BC = &cpu.RegPair{&res.regs.B, &res.regs.C}
	res.regs.DE = &cpu.RegPair{&res.regs.D, &res.regs.E}
	res.regs.HL = &cpu.RegPair{&res.regs.H, &res.regs.L}
	res.regs.SP = &cpu.RegPair{&res.regs.S, &res.regs.P}

	res.scheduler.append(newFetch(lookup))
	return res
}

func (cpu *lr35902) CurrentOP() string { panic(-2) }

func (cpu *lr35902) RegisterTrap(pc uint16, trap CPUTrap) {
	cpu.traps[pc] = trap
}

func (cpu *lr35902) Registers() *LR35902Registers {
	return cpu.regs
}

func (cpu *lr35902) Interrupt(i bool) {
	cpu.doInterrupt = i
}

func (cpu *lr35902) NMI(i bool) {
}

func (cpu *lr35902) Wait(w bool) {
	cpu.wait = w
}

func (cpu *lr35902) Halt() {
	cpu.halt = true
}

func (cpu *lr35902) Reset() {
	panic(-1)
}

func (cpu *lr35902) SetTracer(t cpu.CPUTracer)            { cpu.log = t }
func (cpu *lr35902) SetDebugger(db cpu.DebuggerCallbacks) { cpu.debugger = db }

func (cpu *lr35902) execInterrupt() {
	cpu.prepareForNewInstruction()
	cpu.doInterrupt = false

	if cpu.regs.IFF1 {
		cpu.regs.IFF1 = false
		cpu.regs.IFF2 = false
		switch cpu.regs.InterruptsMode {
		case 0, 1:
			code := &exec{l: 7, f: func(cpu *lr35902) {
				cpu.pushToStack(cpu.regs.PC, func(cpu *lr35902) {
					cpu.regs.PC = 0x0038
				})
			}}
			cpu.scheduler.append(code)
		case 2:
			code := &exec{l: 7, f: func(cpu *lr35902) {
				cpu.pushToStack(cpu.regs.PC, func(cpu *lr35902) {
					pos := uint16(cpu.regs.I)<<8 + 0xff
					mr1 := newMR(pos, func(cpu *lr35902, data byte) {
						cpu.regs.PC = uint16(data) << 8
					})
					mr2 := newMR(pos, func(cpu *lr35902, data byte) {
						cpu.regs.PC |= uint16(data)
					})
					cpu.scheduler.append(mr1, mr2)
				})
			}}
			cpu.scheduler.append(code)
		}
	}
}

func (cpu *lr35902) prepareForNewInstruction() {
	if cpu.debugger != nil {
		cpu.debugger.Eval(cpu.regs.PC)
	}

	cpu.fetched.n = 0
	cpu.fetched.nn = 0
	cpu.fetched.opCode = 0
	cpu.fetched.prefix = 0
	cpu.fetched.pc = cpu.regs.PC
	cpu.indexIdx = 0
}

func (cpu *lr35902) Tick() {
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
			if cpu.log != nil {
				cpu.log.AppendLastOP(cpu.fetched.String())
			}
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

func (cpu *lr35902) newInstruction() {
	cpu.prepareForNewInstruction()
	cpu.doTraps()
	cpu.scheduler.append(newFetch(lookup))
}

func (cpu *lr35902) doTraps() {
	if trap, ok := cpu.traps[cpu.regs.PC]; ok {
		trap()
	}
}

// func (cpu *lr35902) Tick() {
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

// func (cpu *lr35902) Step() {
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
