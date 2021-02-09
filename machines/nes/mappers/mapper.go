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
	if file.header.fourPages {
		ram := &m6502.BasicRam{Data: make([]byte, 0x1000), Mask: 0x0fff}
		ppuBus.RegisterPort(emulator.PortMask{Mask: 0b1111_000000000000, Value: 0b0010_000000000000}, ram)
		ppuBus.RegisterPort(emulator.PortMask{Mask: 0b1111_000000000000, Value: 0b0011_000000000000}, ram)
	} else if file.header.vMirror {
		ram0 := &m6502.BasicRam{Data: make([]byte, 0x400), Mask: 0x03ff}
		ppuBus.RegisterPort(emulator.PortMask{Mask: 0b1111_1100_00000000, Value: 0b0010_000000000000}, ram0)
		ppuBus.RegisterPort(emulator.PortMask{Mask: 0b1111_1100_00000000, Value: 0b0010_100000000000}, ram0)

		ram1 := &m6502.BasicRam{Data: make([]byte, 0x400), Mask: 0x03ff}
		ppuBus.RegisterPort(emulator.PortMask{Mask: 0b1111_1100_00000000, Value: 0b0010_010000000000}, ram1)
		ppuBus.RegisterPort(emulator.PortMask{Mask: 0b1111_1100_00000000, Value: 0b0010_110000000000}, ram1)
	} else {
		ram0 := &m6502.BasicRam{Data: make([]byte, 0x400), Mask: 0x03ff}
		ppuBus.RegisterPort(emulator.PortMask{Mask: 0b1111_1100_00000000, Value: 0b0010_000000000000}, ram0)
		ppuBus.RegisterPort(emulator.PortMask{Mask: 0b1111_1100_00000000, Value: 0b0010_010000000000}, ram0)

		ram1 := &m6502.BasicRam{Data: make([]byte, 0x400), Mask: 0x03ff}
		ppuBus.RegisterPort(emulator.PortMask{Mask: 0b1111_1100_00000000, Value: 0b0010_100000000000}, ram1)
		ppuBus.RegisterPort(emulator.PortMask{Mask: 0b1111_1100_00000000, Value: 0b0010_110000000000}, ram1)
	}
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
