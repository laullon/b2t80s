package mappers

import (
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
)

type nrom struct {
	file     *nesFile
	ram      *ram
	mask     uint16
	pattern0 *rom
	pattern1 *rom
}

func newNROM(file *nesFile) Mapper {
	m := &nrom{file: file}
	if file.header.prgSize == 1 {
		m.mask = 0x3fff
	} else {
		m.mask = 0x7fff
	}
	m.ram = &ram{mem: make([]byte, 0x2000), mask: 0x1fff}
	if m.file.header.chrSize == 1 {
		m.pattern0 = &rom{mem: m.file.chr[:0x1000], mask: 0x0fff}
		m.pattern1 = &rom{mem: m.file.chr[0x1000:], mask: 0x0fff}
	} else if m.file.header.chrSize == 0 {
		m.pattern0 = &rom{mem: make([]byte, 0x1000), mask: 0x0fff}
		m.pattern1 = &rom{mem: make([]byte, 0x1000), mask: 0x0fff}
	} else {
		panic(-1)
	}
	return m
}

func (m *nrom) ConnectToPPU(bus m6502.Bus) {
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111_000000000000, Value: 0b0000_000000000000}, m.pattern0)
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111_000000000000, Value: 0b0001_000000000000}, m.pattern1)
	setPPUMemory(m.file, bus)
}

func (m *nrom) ConnectToCPU(bus m6502.Bus) {
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b01100000_00000000}, m.ram)
	bus.RegisterPort(emulator.PortMask{Mask: 0b10000000_00000000, Value: 0b10000000_00000000}, m)
}

func (m *nrom) ReadPort(addr uint16) (byte, bool) {
	return m.file.prg[addr&m.mask], false
}
func (m *nrom) WritePort(addr uint16, data byte) { panic(-1) }
