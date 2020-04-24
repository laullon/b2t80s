package z80

import (
	"bufio"
	"fmt"
	"os"

	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
)

var overflowAddTable = []bool{false, false, false, true, true, false, false, false}
var overflowSubTable = []bool{false, true, false, false, false, false, true, false}
var halfcarryAddTable = []bool{false, true, true, true, false, false, false, true}
var halfcarrySubTable = []bool{false, false, true, false, true, false, true, true}

var parityTable = make([]bool, 0x100)

type Z80Registers struct {
	PC uint16
	SP emulator.StackPointer

	A byte
	F *flags

	B  byte
	C  byte
	BC *regPair

	D  byte
	E  byte
	DE *regPair

	H  byte
	L  byte
	HL *regPair

	I  byte
	R  byte
	R7 byte

	IFF1 bool
	IFF2 bool

	IXH byte
	IXL byte
	IX  *regPair

	IYH byte
	IYL byte
	IY  *regPair

	_A byte
	_F *flags
	_B byte
	_C byte
	_D byte
	_E byte
	_H byte
	_L byte
}

type z80 struct {
	memory   emulator.Memory
	cassette cassette.Cassette
	debugger emulator.Debugger

	regs *Z80Registers

	halt, haltDone bool

	interruptsMode byte
	doInterrupt    bool

	actualOPCode int32

	traps map[uint16]emulator.CPUTrap
	ports map[emulator.PortMask]emulator.PortManager

	clock emulator.Clock
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

func NewZ80(mem emulator.Memory, cassette cassette.Cassette) emulator.CPU {
	LoadOPCodess()
	cpu := &z80{
		memory:   mem,
		cassette: cassette,
		traps:    make(map[uint16]emulator.CPUTrap),
		ports:    make(map[emulator.PortMask]emulator.PortManager),
		regs: &Z80Registers{
			PC: 0,
			A:  0xff,
			F: &flags{
				Z: true,
				C: true,
				S: true,
				H: true,
				P: true,
				N: true,
			},
			_F: &flags{
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

	cpu.regs.BC = &regPair{&cpu.regs.B, &cpu.regs.C}
	cpu.regs.DE = &regPair{&cpu.regs.D, &cpu.regs.E}
	cpu.regs.HL = &regPair{&cpu.regs.H, &cpu.regs.L}
	cpu.regs.IX = &regPair{&cpu.regs.IXH, &cpu.regs.IXL}
	cpu.regs.IY = &regPair{&cpu.regs.IYH, &cpu.regs.IYL}
	cpu.regs.SP = NewStackPointer(cpu.memory)

	return cpu
}

func (cpu *z80) SetClock(clock emulator.Clock) {
	cpu.clock = clock
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

func (cpu *z80) SetRegisters(regs []byte, i, r, iff1, mode byte) {
	cpu.regs.A = regs[0]
	cpu.regs.F.setByte(regs[1])
	cpu.regs.B = regs[2]
	cpu.regs.C = regs[3]
	cpu.regs.D = regs[4]
	cpu.regs.E = regs[5]
	cpu.regs.H = regs[6]
	cpu.regs.L = regs[7]
	cpu.regs.IXH = regs[8]
	cpu.regs.IXL = regs[9]
	cpu.regs.IYH = regs[10]
	cpu.regs.IYL = regs[11]
	cpu.regs._A = regs[12]
	cpu.regs._F.setByte(regs[13])
	cpu.regs._B = regs[14]
	cpu.regs._C = regs[15]
	cpu.regs._D = regs[16]
	cpu.regs._E = regs[17]
	cpu.regs._H = regs[18]
	cpu.regs._L = regs[19]

	cpu.regs.I = i
	cpu.regs.R = r
	cpu.regs.IFF1 = iff1 != 0
	cpu.regs.IFF2 = !cpu.regs.IFF1
	cpu.interruptsMode = mode
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

		cpu.regs.SP.Push(cpu.regs.PC)

		switch cpu.interruptsMode {
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

func (cpu *z80) RunFrame() error {
	var done bool
	for !done {
		if cpu.debugger.IsStoped() {
			return nil
		}
		cpu.Step()
		done = cpu.clock.FrameDone()
	}
	return nil
}

func (cpu *z80) Step() {
	if trap, ok := cpu.traps[cpu.regs.PC]; ok {
		res := trap()
		switch res {
		case emulator.CONTINUE:
		case emulator.STOP:
			return
		default:
			cpu.regs.PC = uint16(res)
			return
		}
	}

	if cpu.doInterrupt {
		ts := cpu.execInterrupt()
		cpu.clock.AddTStates(ts)
	}

	ins, err := GetOpCode(cpu.memory.GetBlock(cpu.regs.PC, 4))
	if err != nil {
		panic(err)
	}

	if cpu.debugger != nil {
		cpu.debugger.AddLastInstruction(ins)
	}

	var ts uint
	if !cpu.halt {
		needPcUpdate := cpu.runSwitch(ins)
		if needPcUpdate {
			cpu.regs.PC += uint16(ins.Length)
		}
		ts = ins.Tstates
	} else {
		ts = 4 // halt
	}

	cpu.clock.AddTStates(ts)
	return
}

// func (cpu *z80) LoadTapeBlock() uint16 {
// block := cpu.regs.Cassette.NextBlock()
// for len(block.GetData()) == 0 {
// 	block = cpu.regs.Cassette.NextBlock()
// }

// requestedLength := getRR(cpu.regs.D, cpu.regs.E)
// startAddress := getRR(cpu.regs.IXH, cpu.regs.IXL)
// fmt.Printf("Loading block '%s' to 0x%04x (bl:0x%04x, l:0x%04x, bt:%d, a:%d)\n", block.Name(), startAddress, len(block.GetData()), requestedLength, block.Type(), cpu.regs._A)
// if cpu.regs._A == block.Type() {
// 	if cpu.regs._F.C {
// 		checksum := block.Type()
// 		data := block.GetData()
// 		for i := uint16(0); i < requestedLength; i++ {
// 			loadedByte := data[i+1]
// 			cpu.memory.PutByte(startAddress+i, loadedByte)
// 			checksum ^= loadedByte
// 		}
// 		cpu.regs.F.C = checksum == data[requestedLength+1]
// 	} else {
// 		cpu.regs.F.C = true
// 	}
// 	// log.Print("done")
// } else {
// 	cpu.regs.F.C = false
// 	// log.Print("BAD Block")
// }
// return 0x05e2
// }

// TODO: move out
// func (cpu *z80) LoadTapeBlockCPC(exit uint16) uint16 {
// 	block := cpu.regs.Cassette.NextBlock()
// 	if len(block.GetData()) == 0 {
// 		block = cpu.regs.Cassette.NextBlock()
// 	}

// 	requestedLength := getRR(cpu.regs.D, cpu.regs.E)
// 	startAddress := getRR(cpu.regs.H, cpu.regs.L)
// 	t := cpu.regs.A
// 	// fmt.Printf("Loading block '%s' to 0x%04x (bl:0x%04x, l:0x%04x, bt:0x%02X, t:0x%02X)\n", block.Name(), startAddress, len(block.GetData()), requestedLength, block.Type(), t)
// 	if t == block.Type() {
// 		data := block.GetData()
// 		for i := uint16(0); i < requestedLength; i++ {
// 			cpu.memory.PutByte(startAddress+i, data[i+1])
// 		}
// 		cpu.regs.F.setByte(0x45)
// 		// log.Print("Done")
// 		// println(hex.Dump(cpu.memory.GetBlock(startAddress, requestedLength)))
// 	} else {
// 		// log.Print("BAD Block")
// 	}
// 	return exit
// }

func (cpu *z80) pause() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (cpu *z80) dumpRegState() {
	// fmt.Printf("\t\t\t\t\t\t\t\taf:0x%02X%02X bc:0x%02X%02X de:0x%02X%02X hl:0x%02X%02X ix:0x%02X%02X iy:0x%02X%02X sp:0x%04X flags: [z:%v c:%v]\n", cpu.a, cpu.f.getByte(), cpu.b, cpu.c, cpu.d, cpu.e, cpu.h, cpu.l, cpu.ixh, cpu.ixl, cpu.iyh, cpu.iyl, cpu.sp.Get(), cpu.f.Z, cpu.f.C)
}
