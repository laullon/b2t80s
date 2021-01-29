package mappers

import (
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
)

type nrom struct {
	file *nesFile
	ram  *ram
	mask uint16
}

func newNROM(file *nesFile) Mapper {
	m := &nrom{file: file}
	if file.header.prgSize == 1 {
		m.mask = 0x3fff
	} else {
		m.mask = 0x7fff
	}
	m.ram = &ram{mem: make([]byte, 0x2000), mask: 0x1fff}
	return m
}

func (m *nrom) Insert(bus m6502.Bus) {
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b01100000_00000000}, m.ram)
	bus.RegisterPort(emulator.PortMask{Mask: 0b10000000_00000000, Value: 0b10000000_00000000}, m)
}

func (m *nrom) ReadPort(addr uint16) (byte, bool) {
	if addr == 0xfffc {
		return 0x00, false
	}
	return m.file.prg[addr&m.mask], false
}
func (m *nrom) WritePort(addr uint16, data byte) { panic(-1) }
