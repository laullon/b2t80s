package z80

import (
	"fmt"
	"strings"

	"github.com/laullon/b2t80s/cpu"
	cpuUtils "github.com/laullon/b2t80s/cpu"
	"github.com/pkg/errors"
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

type CPUTrap func()

type Z80 interface {
	cpu.CPU
	Registers() *Z80Registers
	RegisterTrap(pc uint16, trap CPUTrap)
	Interrupt(bool, ...byte)
}

type z80 struct {
	bus Bus

	regs      *Z80Registers
	indexRegs []*cpuUtils.RegPair
	indexIdx  int

	halt bool
	wait bool

	doInterrupt   bool
	interruptData byte

	fetched   *fetchedData
	scheduler *circularBuffer

	traps map[uint16]CPUTrap

	log      cpu.CPUTracer
	debugger cpu.DebuggerCallbacks

	lookup     []*opCode
	lookupCB   []*opCode
	lookupDD   []*opCode
	lookupED   []*opCode
	lookupFD   []*opCode
	lookupDDCB []*opCode
	lookupFDCB []*opCode

	cpi_result uint8
	spv        uint16
	inAn_f     byte
	pushF      func(cpu *z80)
	popData    uint16
	popF       func(cpu *z80, data uint16)
	hlv        uint8
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

func NewZ80(bus Bus) Z80 {
	cpu := &z80{
		bus:       bus,
		fetched:   &fetchedData{},
		scheduler: newCircularBuffer(),
		traps:     make(map[uint16]CPUTrap),
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

	cpu.initOpsCodes()

	cpu.regs.BC = &cpuUtils.RegPair{&cpu.regs.B, &cpu.regs.C}
	cpu.regs.DE = &cpuUtils.RegPair{&cpu.regs.D, &cpu.regs.E}
	cpu.regs.HL = &cpuUtils.RegPair{&cpu.regs.H, &cpu.regs.L}
	cpu.regs.SP = &cpuUtils.RegPair{&cpu.regs.S, &cpu.regs.P}
	cpu.regs.IX = &cpuUtils.RegPair{&cpu.regs.IXH, &cpu.regs.IXL}
	cpu.regs.IY = &cpuUtils.RegPair{&cpu.regs.IYH, &cpu.regs.IYL}
	cpu.indexRegs = []*cpuUtils.RegPair{cpu.regs.HL, cpu.regs.IX, cpu.regs.IY}

	cpu.scheduler.append(newFetch(cpu.lookup))
	return cpu
}

func (cpu *z80) Status() string { return cpu.regs.dump() }

func (cpu *z80) FullStatus() string {
	var res strings.Builder
	res.WriteString(fmt.Sprintf("  A:0x%02X    F:0x%02X  AF:0x%04X    SP:0x%04X\n", cpu.regs.A, cpu.regs.F.GetByte(), uint16(cpu.regs.A)<<8|uint16(cpu.regs.F.GetByte()), cpu.regs.SP.Get()))
	res.WriteString(fmt.Sprintf("  B:0x%02X    C:0x%02X  BC:0x%04X    ---------\n", cpu.regs.B, cpu.regs.C, uint16(cpu.regs.B)<<8|uint16(cpu.regs.C)))
	res.WriteString(fmt.Sprintf("  D:0x%02X    E:0x%02X  DE:0x%04X    0x%04X\n", cpu.regs.D, cpu.regs.E, uint16(cpu.regs.D)<<8|uint16(cpu.regs.E), "-"))         // getWord(debug.memory, cpu.regs.SP.Get()+0)))
	res.WriteString(fmt.Sprintf("  H:0x%02X    L:0x%02X  HL:0x%04X    0x%04X\n", cpu.regs.H, cpu.regs.L, uint16(cpu.regs.H)<<8|uint16(cpu.regs.L), "-"))         // getWord(debug.memory, cpu.regs.SP.Get()+2)))
	res.WriteString(fmt.Sprintf("IXH:0x%02X  IXL:0x%02X  IX:0x%04X    0x%04X\n", cpu.regs.IXH, cpu.regs.IXL, uint16(cpu.regs.IXH)<<8|uint16(cpu.regs.IXL), "-")) // getWord(debug.memory, cpu.regs.SP.Get()+4)))
	res.WriteString(fmt.Sprintf("IYH:0x%02X  IYL:0x%02X  IY:0x%04X    0x%04X\n", cpu.regs.IYH, cpu.regs.IYL, uint16(cpu.regs.IYH)<<8|uint16(cpu.regs.IYL), "-")) // getWord(debug.memory, cpu.regs.SP.Get()+6)))
	res.WriteString(fmt.Sprintf("SZ5H3PNC            PC:0x%04X\n%08b", cpu.regs.PC, cpu.regs.F.GetByte()))
	return res.String()
}

func (cpu *z80) CurrentOP() string { panic(-2) }

func (cpu *z80) RegisterTrap(pc uint16, trap CPUTrap) {
	cpu.traps[pc] = trap
}

func (cpu *z80) Registers() *Z80Registers {
	return cpu.regs
}

func (cpu *z80) Interrupt(i bool, data ...byte) {
	cpu.doInterrupt = i
	if len(data) > 0 {
		cpu.interruptData = data[0]
	}
}

func (cpu *z80) NMI(i bool) {
}

func (cpu *z80) Wait(w bool) {
	cpu.wait = w
}

func (cpu *z80) Halt() {
	cpu.halt = true
}

func (cpu *z80) Reset() {
	panic(-1)
}

func (cpu *z80) SetTracer(t cpu.CPUTracer)            { cpu.log = t }
func (cpu *z80) SetDebugger(db cpu.DebuggerCallbacks) { cpu.debugger = db }

func (cpu *z80) execInterrupt() {
	if cpu.debugger != nil {
		cpu.debugger.EvalInterrupt()
	}

	cpu.prepareForNewInstruction()
	cpu.doInterrupt = false

	if cpu.regs.IFF1 {
		cpu.regs.IFF1 = false
		cpu.regs.IFF2 = false
		switch cpu.regs.InterruptsMode {
		case 0:
			opCode := cpu.interruptData
			cpu.fetched.opCode = opCode
			op := cpu.lookup[opCode]
			if op == nil {
				panic(errors.Errorf("opCode '%X - %X' not found", 0, opCode))
			}
			for _, op := range op.ops {
				op.reset()
			}
			cpu.scheduler.clear()
			cpu.scheduler.append(op.ops...)

		case 1:
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
			if cpu.log != nil {
				cpu.log.AppendLastOP(cpu.fetched.disassemble()) //fmt.Sprintf("%04x: %s", cpu.fetched.pc, cpu.fetched.getInstruction()))
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

func (cpu *z80) newInstruction() {
	cpu.prepareForNewInstruction()
	cpu.doTraps()
	cpu.scheduler.append(newFetch(cpu.lookup))
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
