package mappers

import (
	"fmt"

	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
)

type mmc1 struct {
	file *nesFile
	ram  *ram
	rom  []*rom
	sr   sr

	control byte

	mirror byte
	prgmod byte
	chrmod byte
}

func newMMC1(file *nesFile) Mapper {
	m := &mmc1{file: file}
	m.ram = &ram{mem: make([]byte, 0x2000), mask: 0x1fff}
	m.rom = []*rom{
		{mem: file.prg[:0x4000], mask: 0x3fff, write: m.write},
		{mem: file.prg[0x4000*uint16(file.header.prgSize-1):], mask: 0x3fff, write: m.write},
	}
	return m
}

func (m *mmc1) Insert(bus m6502.Bus) {
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b01100000_00000000}, m.ram)
	bus.RegisterPort(emulator.PortMask{Mask: 0b11000000_00000000, Value: 0b10000000_00000000}, m.rom[0])
	bus.RegisterPort(emulator.PortMask{Mask: 0b11000000_00000000, Value: 0b11000000_00000000}, m.rom[1])

	// CPU $6000-$7FFF: 8 KB PRG RAM bank, (optional)
	// CPU $8000-$BFFF: 16 KB PRG ROM bank, either switchable or fixed to the first bank
	// CPU $C000-$FFFF: 16 KB PRG ROM bank, either fixed to the last bank or switchable

}

func (m *mmc1) write(addr uint16, data byte) {
	fmt.Printf("0x%02X => 0x%02X (%08b) \n", addr, data, data)
	if data&0x80 != 0 {
		m.sr.reset()
		m.writeControl(m.control | 0xc0)
	} else {
		if m.sr.tick(data&0x01 != 0) == 5 {
			// fmt.Printf("sr: %v\n", m.sr)
			v := (m.sr.data & 0b11111000) >> 3
			fmt.Printf("v: %05b\n", v)
			m.sr.reset()

			switch addr & 0xe000 {
			case 0x8000:
				m.writeControl(v)
			case 0xE000:
				m.writePRG(v)
			default:
				panic(-1)
			}
		}
	}
}

func (m *mmc1) writeControl(data byte) {
	m.control = data
	m.mirror = data & 0b00011
	m.prgmod = data & 0b01100 >> 2
	m.chrmod = data & 0b10000 >> 4
}

func (m *mmc1) writePRG(data byte) {
	switch m.prgmod {
	case 0, 1:
		b := uint16(data>>1) * 0x8000
		m.rom[0].mem = m.file.prg[b : b+0x4000]
		m.rom[1].mem = m.file.prg[b+0x4000 : b+0x8000]
	case 2:
		b := uint16(data>>1) * 0x4000
		m.rom[0].mem = m.file.prg[:0x4000]
		m.rom[1].mem = m.file.prg[b : b+0x4000]
	case 3:
		b := uint16(data>>1) * 0x4000
		m.rom[0].mem = m.file.prg[b : b+0x4000]
		m.rom[1].mem = m.file.prg[0x4000*uint16(m.file.header.prgSize-1):]
	}
}

// ----------------------------

type sr struct {
	data  byte
	count byte
}

func (sr *sr) tick(d bool) byte {
	sr.data >>= 1
	sr.count++
	if d {
		sr.data |= 0x80
	}
	return sr.count
}

func (sr *sr) reset() {
	sr.data = 0
	sr.count = 0
}
