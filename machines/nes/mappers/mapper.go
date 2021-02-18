package mappers

import (
	"fmt"

	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
)

type Mapper interface {
	ConnectToCPU(bus m6502.Bus)
	ConnectToPPU(bus m6502.Bus)
}

func CreateMapper(fileName string) Mapper {
	file := loadFile(fileName)

	switch file.header.mapper {
	case 0:
		return newNROM(file)
	case 1:
		return newMMC1(file)
	case 3:
		return newCNROM(file)

	default:
		panic(fmt.Sprintf("mapper type '%d' not supported", file.header.mapper))
	}
}

func setPPUMemory(file *nesFile, ppuBus m6502.Bus) {
	var nt0 *m6502.BasicRam
	var nt1 *m6502.BasicRam
	var nt2 *m6502.BasicRam
	var nt3 *m6502.BasicRam
	if file.header.fourPages {
		nt0 = &m6502.BasicRam{Data: make([]byte, 0x0400), Mask: 0x03ff}
		nt1 = &m6502.BasicRam{Data: make([]byte, 0x0400), Mask: 0x03ff}
		nt2 = &m6502.BasicRam{Data: make([]byte, 0x0400), Mask: 0x03ff}
		nt3 = &m6502.BasicRam{Data: make([]byte, 0x0400), Mask: 0x03ff}
	} else if file.header.vMirror {
		nt0 = &m6502.BasicRam{Data: make([]byte, 0x0400), Mask: 0x03ff}
		nt1 = &m6502.BasicRam{Data: make([]byte, 0x0400), Mask: 0x03ff}
		nt2 = nt0
		nt3 = nt1
	} else {
		nt0 = &m6502.BasicRam{Data: make([]byte, 0x0400), Mask: 0x03ff}
		nt1 = nt0
		nt2 = &m6502.BasicRam{Data: make([]byte, 0x0400), Mask: 0x03ff}
		nt3 = nt2
	}

	ppuBus.RegisterPort("NameTable_0", emulator.PortMask{Mask: 0b1111_1100_0000_0000, Value: 0b0010_0000_0000_0000}, nt0)
	ppuBus.RegisterPort("NameTable_1", emulator.PortMask{Mask: 0b1111_1100_0000_0000, Value: 0b0010_0100_0000_0000}, nt1)
	ppuBus.RegisterPort("NameTable_2", emulator.PortMask{Mask: 0b1111_1100_0000_0000, Value: 0b0010_1000_0000_0000}, nt2)
	ppuBus.RegisterPort("NameTable_3", emulator.PortMask{Mask: 0b1111_1100_0000_0000, Value: 0b0010_1100_0000_0000}, nt3)

	// ppuBus.RegisterPort("NameTable_0m", emulator.PortMask{Mask: 0b1111_1100_0000_0000, Value: 0b0011_0000_0000_0000}, nt0)
	// ppuBus.RegisterPort("NameTable_1m", emulator.PortMask{Mask: 0b1111_1100_0000_0000, Value: 0b0011_0100_0000_0000}, nt1)
	// ppuBus.RegisterPort("NameTable_2m", emulator.PortMask{Mask: 0b1111_1100_0000_0000, Value: 0b0011_1000_0000_0000}, nt2)
	// ppuBus.RegisterPort("NameTable_3m1", emulator.PortMask{Mask: 0b1111_1111_0000_0000, Value: 0b0011_1100_0000_0000}, nt3)
	// ppuBus.RegisterPort("NameTable_3m2", emulator.PortMask{Mask: 0b1111_1111_0000_0000, Value: 0b0011_1110_0000_0000}, nt3)
	// ppuBus.RegisterPort("NameTable_3m3", emulator.PortMask{Mask: 0b1111_1111_0000_0000, Value: 0b0011_1101_0000_0000}, nt3)

}

// ----------------------------

type rom struct {
	mem   []byte
	mask  uint16
	write func(addr uint16, data byte)
}

func (rom *rom) ReadPort(addr uint16) (byte, bool) { return rom.mem[addr&rom.mask], false }
func (rom *rom) WritePort(addr uint16, data byte)  { rom.mem[addr&rom.mask] = data }

// ----------------------------

type ram struct {
	mem  []byte
	mask uint16
}

func (ram *ram) ReadPort(addr uint16) (byte, bool) { return ram.mem[addr&ram.mask], false }
func (ram *ram) WritePort(addr uint16, data byte)  { ram.mem[addr&ram.mask] = data }
