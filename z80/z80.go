package z80

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/laullon/b2t80s/emulator"
)

var overflowAddTable = []bool{false, false, false, true, true, false, false, false}
var overflowSubTable = []bool{false, true, false, false, false, false, true, false}
var halfcarryAddTable = []bool{false, true, true, true, false, false, false, true}
var halfcarrySubTable = []bool{false, false, true, false, true, false, true, true}

var parityTable = make([]bool, 0x100)

type z80 struct {
	memory   emulator.Memory
	cassette emulator.Cassette
	debugger emulator.Debugger

	halt, haltDone bool
	pc             uint16
	sp             emulator.StackPointer

	interruptsMode byte
	doInterrupt    bool

	a byte
	f *flags

	b byte
	c byte

	d byte
	e byte

	h byte
	l byte

	i  byte
	r  byte
	r7 byte

	iff1 bool
	iff2 bool

	ixh byte
	ixl byte

	iyh byte
	iyl byte

	_a byte
	_f *flags
	_b byte
	_c byte
	_d byte
	_e byte
	_h byte
	_l byte

	actualOPCode int32

	debug   bool
	useMaps bool

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

func NewZ80(mem emulator.Memory, cassette emulator.Cassette) emulator.CPU {
	LoadOPCodess()
	cpu := &z80{
		debug: false,

		pc:       0,
		memory:   mem,
		cassette: cassette,
		traps:    make(map[uint16]emulator.CPUTrap),
		ports:    make(map[emulator.PortMask]emulator.PortManager),

		a: 0xff,
		f: &flags{
			Z: true,
			C: true,
			S: true,
			H: true,
			P: true,
			N: true,
		},
		_f: &flags{
			Z: true,
			C: true,
			S: true,
			H: true,
			P: true,
			N: true,
		},
		r: 0x01,
	}

	cpu.sp = NewStackPointer(cpu.memory)

	return cpu
}

func (cpu *z80) SetClock(clock emulator.Clock) {
	cpu.clock = clock
}

func (cpu *z80) SetDebuger(debugger emulator.Debugger) {
	cpu.debugger = debugger
}

func (cpu *z80) PC() uint16 {
	return cpu.pc
}

func (cpu *z80) SP() emulator.StackPointer {
	return cpu.sp
}

func (cpu *z80) SetPC(pc uint16) {
	cpu.pc = pc
}

func (cpu *z80) RegisterPort(mask emulator.PortMask, manager emulator.PortManager) {
	cpu.ports[mask] = manager
}

func (cpu *z80) RegisterTrap(pc uint16, trap emulator.CPUTrap) {
	cpu.traps[pc] = trap
}

func (cpu *z80) SetRegisters(regs []byte, i, r, iff1, mode byte) {
	cpu.a = regs[0]
	cpu.f.setByte(regs[1])
	cpu.b = regs[2]
	cpu.c = regs[3]
	cpu.d = regs[4]
	cpu.e = regs[5]
	cpu.h = regs[6]
	cpu.l = regs[7]
	cpu.ixh = regs[8]
	cpu.ixl = regs[9]
	cpu.iyh = regs[10]
	cpu.iyl = regs[11]
	cpu._a = regs[12]
	cpu._f.setByte(regs[13])
	cpu._b = regs[14]
	cpu._c = regs[15]
	cpu._d = regs[16]
	cpu._e = regs[17]
	cpu._h = regs[18]
	cpu._l = regs[19]

	cpu.i = i
	cpu.r = r
	cpu.iff1 = iff1 != 0
	cpu.iff2 = !cpu.iff1
	cpu.interruptsMode = mode
}

func (cpu *z80) DumpRegisters() ([]byte, uint16, uint16) {
	return []byte{
		cpu.a, cpu.f.getByte(),
		cpu.b, cpu.c,
		cpu.d, cpu.e,
		cpu.h, cpu.l,
		cpu.ixh, cpu.ixl,
		cpu.iyh, cpu.iyl,
		cpu._a, cpu._f.getByte(),
		cpu._b, cpu._c,
		cpu._d, cpu._e,
		cpu._h, cpu._l,
	}, cpu.sp.Get(), cpu.pc
}

func (cpu *z80) SetRegistersStr(line string, otherReg []byte) {
	regs := strings.Split(line, " ")
	cpu.a, _ = setRRstr(regs[0])
	cpu.b, cpu.c = setRRstr(regs[1])
	cpu.d, cpu.e = setRRstr(regs[2])
	cpu.h, cpu.l = setRRstr(regs[3])

	cpu._a, _ = setRRstr(regs[4])
	cpu._b, cpu._c = setRRstr(regs[5])
	cpu._d, cpu._e = setRRstr(regs[6])
	cpu._h, cpu._l = setRRstr(regs[7])

	cpu.ixh, cpu.ixl = setRRstr(regs[8])
	cpu.iyh, cpu.iyl = setRRstr(regs[9])

	cpu.sp.Set(getRR(setRRstr(regs[10])))
	cpu.pc = getRR(setRRstr(regs[11]))

	_, f := setRRstr(regs[0])
	cpu.f.setByte(f)
	_, _f := setRRstr(regs[4])
	cpu._f.setByte(_f)

	cpu.i = otherReg[0]
	cpu.r = otherReg[1]
	cpu.r7 = otherReg[1]
	cpu.iff2 = otherReg[3] != 0
}

func (cpu *z80) DumpMemory(start, length uint16) {
	// fmt.Printf("%s", hex.Dump(cpu.memory.GetBlock(start, length)))
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
	if cpu.iff1 {
		cpu.iff1 = false
		cpu.iff2 = false

		cpu.sp.Push(cpu.pc)

		switch cpu.interruptsMode {
		case 0, 1:
			ts = 13
			cpu.pc = 0x38
		case 2:
			ts = 19
			pos := uint16(cpu.i)<<8 + 0xff
			cpu.pc = cpu.memory.GetWord(pos)
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
	if trap, ok := cpu.traps[cpu.pc]; ok {
		res := trap()
		switch res {
		case emulator.CONTINUE:
		case emulator.STOP:
			return
		default:
			cpu.pc = uint16(res)
			return
		}
	}

	if cpu.doInterrupt {
		cpu.execInterrupt()
	}

	ins, err := GetOpCode(cpu.memory.GetBlock(cpu.pc, 4))
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
			cpu.pc += uint16(ins.Length)
		}
		ts = ins.Tstates
	} else {
		ts = 4 // halt
	}

	cpu.clock.AddTStates(ts)
	return
}

func (cpu *z80) LoadTapeBlock() uint16 {
	block := cpu.cassette.NextBlock()
	for len(block.GetData()) == 0 {
		block = cpu.cassette.NextBlock()
	}

	requestedLength := getRR(cpu.d, cpu.e)
	startAddress := getRR(cpu.ixh, cpu.ixl)
	fmt.Printf("Loading block '%s' to 0x%04x (bl:0x%04x, l:0x%04x, bt:%d, a:%d)\n", block.Name(), startAddress, len(block.GetData()), requestedLength, block.Type(), cpu._a)
	if cpu._a == block.Type() {
		if cpu._f.C {
			checksum := block.Type()
			data := block.GetData()
			for i := uint16(0); i < requestedLength; i++ {
				loadedByte := data[i+1]
				cpu.memory.PutByte(startAddress+i, loadedByte)
				checksum ^= loadedByte
			}
			cpu.f.C = checksum == data[requestedLength+1]
		} else {
			cpu.f.C = true
		}
		// log.Print("done")
	} else {
		cpu.f.C = false
		// log.Print("BAD Block")
	}
	return 0x05e2
}

// TODO: move out
func (cpu *z80) LoadTapeBlockCPC(exit uint16) uint16 {
	block := cpu.cassette.NextBlock()
	if len(block.GetData()) == 0 {
		block = cpu.cassette.NextBlock()
	}

	requestedLength := getRR(cpu.d, cpu.e)
	startAddress := getRR(cpu.h, cpu.l)
	t := cpu.a
	// fmt.Printf("Loading block '%s' to 0x%04x (bl:0x%04x, l:0x%04x, bt:0x%02X, t:0x%02X)\n", block.Name(), startAddress, len(block.GetData()), requestedLength, block.Type(), t)
	if t == block.Type() {
		data := block.GetData()
		for i := uint16(0); i < requestedLength; i++ {
			cpu.memory.PutByte(startAddress+i, data[i+1])
		}
		cpu.f.setByte(0x45)
		// log.Print("Done")
		// println(hex.Dump(cpu.memory.GetBlock(startAddress, requestedLength)))
	} else {
		// log.Print("BAD Block")
	}
	return exit
}

func (cpu *z80) pause() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (cpu *z80) dumpRegState() {
	// fmt.Printf("\t\t\t\t\t\t\t\taf:0x%02X%02X bc:0x%02X%02X de:0x%02X%02X hl:0x%02X%02X ix:0x%02X%02X iy:0x%02X%02X sp:0x%04X flags: [z:%v c:%v]\n", cpu.a, cpu.f.getByte(), cpu.b, cpu.c, cpu.d, cpu.e, cpu.h, cpu.l, cpu.ixh, cpu.ixl, cpu.iyh, cpu.iyl, cpu.sp.Get(), cpu.f.Z, cpu.f.C)
}

func getRR(h, l byte) uint16 {
	return (uint16(h) << 8) | uint16(l)
}

func setRR(hl uint16) (uint8, uint8) {
	return uint8(hl >> 8), uint8(hl & 0x00ff)
}

func setRRstr(hl string) (uint8, uint8) {
	decoded, err := hex.DecodeString(hl)
	if err != nil {
		panic(fmt.Sprintf("string: '%v' error: %v", hl, err))
	}
	return decoded[0], decoded[1]
}
