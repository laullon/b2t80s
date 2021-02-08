package mappers

import (
	"fmt"

	"github.com/laullon/b2t80s/cpu/m6502"
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
