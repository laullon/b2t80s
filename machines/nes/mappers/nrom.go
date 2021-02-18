package mappers

import (
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
)

type nrom struct {
	file     *nesFile
	ram      *ram
	rom0     *rom
	rom1     *rom
	pattern0 *rom
	pattern1 *rom
}

func newNROM(file *nesFile) Mapper {
	m := &nrom{file: file}

	if file.header.prgSize == 1 {
		m.rom0 = &rom{mem: m.file.prg, mask: 0x3fff}
		m.rom1 = m.rom0
	} else {
		m.rom0 = &rom{mem: m.file.prg[:0x4000], mask: 0x3fff}
		m.rom1 = &rom{mem: m.file.prg[0x4000:], mask: 0x3fff}
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
	bus.RegisterPort("cart.pattern0", emulator.PortMask{Mask: 0b1111_0000_0000_0000, Value: 0b0000_0000_0000_0000}, m.pattern0)
	bus.RegisterPort("cart.pattern1", emulator.PortMask{Mask: 0b1111_0000_0000_0000, Value: 0b0001_0000_0000_0000}, m.pattern1)
	setPPUMemory(m.file, bus)
}

func (m *nrom) ConnectToCPU(bus m6502.Bus) {
	bus.RegisterPort("cart.ram", emulator.PortMask{Mask: 0b1110_0000_0000_0000, Value: 0b0110_0000_0000_0000}, m.ram)
	bus.RegisterPort("cart.rom0", emulator.PortMask{Mask: 0b1100_0000_0000_0000, Value: 0b1000_0000_0000_0000}, m.rom0)
	bus.RegisterPort("cart.rom1", emulator.PortMask{Mask: 0b1100_0000_0000_0000, Value: 0b1100_0000_0000_0000}, m.rom1)
}
