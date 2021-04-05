package lr35902

import (
	"fmt"
	"strings"

	"github.com/laullon/b2t80s/cpu"
)

type LR35902Registers struct {
	PC uint16

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

	IME bool
	IE  byte
	IF  byte
}

type fetchedData struct {
	pc     uint16
	prefix byte
	opCode uint8
	n      uint8
	n2     uint8
	nn     uint16
	op     *opCode
}

func (d *fetchedData) getInstruction() string {
	return d.op.Ins
}

func (d *fetchedData) getMemory() string {
	var res strings.Builder
	if d.prefix > 0xff && d.op.Len == 4 {
		res.WriteString(fmt.Sprintf("%02X ", (d.prefix)))
		res.WriteString(fmt.Sprintf("%02X ", (d.prefix & 0xff)))
		res.WriteString(fmt.Sprintf("%02X ", d.n))
		res.WriteString(fmt.Sprintf("%02X", d.opCode))
	} else if d.prefix > 0x00 && d.op.Len == 4 {
		res.WriteString(fmt.Sprintf("%02X ", (d.prefix & 0xff)))
		res.WriteString(fmt.Sprintf("%02X ", d.opCode))
		res.WriteString(fmt.Sprintf("%02X ", d.n))
		res.WriteString(fmt.Sprintf("%02X", d.n2))
	} else if d.prefix > 0x00 && d.op.Len == 3 {
		res.WriteString(fmt.Sprintf("%02X ", (d.prefix & 0xff)))
		res.WriteString(fmt.Sprintf("%02X ", d.opCode))
		res.WriteString(fmt.Sprintf("%02X", d.n))
	} else if d.prefix == 0x00 && d.op.Len == 3 {
		res.WriteString(fmt.Sprintf("%02X ", d.opCode))
		res.WriteString(fmt.Sprintf("%02X ", d.n))
		res.WriteString(fmt.Sprintf("%02X", d.n2))
	} else if d.prefix == 0x00 && d.op.Len == 2 {
		res.WriteString(fmt.Sprintf("%02X ", d.opCode))
		res.WriteString(fmt.Sprintf("%02X", d.n))
	} else {
		res.WriteString(fmt.Sprintf("%02X", d.opCode))
	}
	return res.String()
}

func (d *fetchedData) String() string {
	return fmt.Sprintf("0x%04X: %-11s : %s", d.pc, d.getMemory(), d.op.Ins)
}

type CPUTrap func()

var intVector = []uint16{0x40, 0x48, 0x50, 0x58, 0x60}

type LR35902 interface {
	cpu.CPU
	Registers() *LR35902Registers
	RegisterTrap(pc uint16, trap CPUTrap)
	cpu.PortManager
}

type lr35902 struct {
	bus cpu.Bus

	regs *LR35902Registers

	halt bool
	wait bool

	fetched   *fetchedData
	scheduler *circularBuffer

	log      cpu.CPUTracer
	debugger cpu.DebuggerCallbacks
}

func New(bus cpu.Bus) LR35902 {
	res := &lr35902{
		bus:       bus,
		fetched:   &fetchedData{},
		scheduler: newCircularBuffer(),
		regs: &LR35902Registers{
			PC: 0,
			A:  0,
			S:  0,
			P:  0,
			F: &flags{
				Z: false,
				C: false,
				H: false,
				N: false,
			},
		},
	}

	res.regs.BC = &cpu.RegPair{&res.regs.B, &res.regs.C}
	res.regs.DE = &cpu.RegPair{&res.regs.D, &res.regs.E}
	res.regs.HL = &cpu.RegPair{&res.regs.H, &res.regs.L}
	res.regs.SP = &cpu.RegPair{&res.regs.S, &res.regs.P}

	return res
}

// TODO: remove from CPU interface
func (cpu *lr35902) Interrupt(bool)                       {}
func (cpu *lr35902) NMI(bool)                             {}
func (cpu *lr35902) CurrentOP() string                    { panic(-2) }
func (cpu *lr35902) RegisterTrap(pc uint16, trap CPUTrap) { panic(-2) }

func (cpu *lr35902) Registers() *LR35902Registers {
	return cpu.regs
}

func (cpu *lr35902) Wait(w bool) {
	cpu.wait = w
}

func (cpu *lr35902) Halt() {
	cpu.halt = true
}

func (cpu *lr35902) Reset() {
	cpu.regs.PC = 0
	cpu.scheduler.append(&fetch{table: OPCodes})
}

func (cpu *lr35902) SetTracer(t cpu.CPUTracer)            { cpu.log = t }
func (cpu *lr35902) SetDebugger(db cpu.DebuggerCallbacks) { cpu.debugger = db }

func (cpu *lr35902) Tick() {
	if cpu.wait {
		return
	}

	// if cpu.halt {
	// 	if cpu.regs.IME {
	// 		cpu.halt = false
	// 		cpu.regs.PC++
	// 		cpu.execInterrupt()
	// 	} else {
	// 		return
	// 	}
	// }

	if cpu.scheduler.isEmpty() {
		cpu.prepareNewInstruction()
		if !cpu.execInterrupt() {
			cpu.scheduler.append(&fetch{table: OPCodes})
		} else {
			if cpu.debugger != nil {
				cpu.debugger.EvalInterrupt()
			}
		}
	}

	cpu.scheduler.first().tick(cpu)
	cpu.scheduler.next()
}

func (cpu *lr35902) prepareNewInstruction() {
	if cpu.log != nil {
		cpu.log.SetDiss(cpu.regs.PC, func(pc, l uint16) []byte {
			return cpu.GetBlock(pc, l)
		})
	}

	if cpu.debugger != nil {
		cpu.debugger.Eval(cpu.regs.PC)
	}

	cpu.fetched.n = 0
	cpu.fetched.nn = 0
	cpu.fetched.opCode = 0
	cpu.fetched.prefix = 0
	cpu.fetched.pc = cpu.regs.PC
}

func (cpu *lr35902) execInterrupt() bool {
	if cpu.regs.IME {
		for i := 0; i < 5; i++ {
			bit := byte(1) << i
			if (cpu.regs.IE&bit != 0) && (cpu.regs.IF&bit != 0) {
				cpu.regs.IME = false
				if cpu.debugger != nil {
					cpu.debugger.EvalInterrupt()
				}
				cpu.regs.IF &^= bit

				cpu.scheduler.append(&exec{}, &exec{})
				cpu.pushToStack(cpu.regs.PC, nil)
				cpu.scheduler.append(&exec{f: func(cpu *lr35902) { cpu.regs.PC = intVector[i] }})

				return true
			}
		}
	}
	return false
}

func (cpu *lr35902) WritePort(addr uint16, data byte) {
	switch addr {
	case 0xffff:
		cpu.regs.IE = data
	case 0xff0f:
		cpu.regs.IF = data
	default:
		panic(-1)
	}
}

func (cpu *lr35902) ReadPort(addr uint16) (byte, bool) {
	switch addr {
	}
	panic(-1)

}

func (cpu *lr35902) GetBlock(pc, l uint16) []byte {
	res := make([]byte, l)
	for i := uint16(0); i < l; i++ {
		res[i] = cpu.bus.Read(pc + i)
	}
	return res
}
