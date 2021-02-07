package mappers

import (
	"fmt"

	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
)

type cnrom struct {
	file *nesFile
	mask uint16
	chr  *rom
}

func newCNROM(file *nesFile) Mapper {
	m := &cnrom{file: file}
	if file.header.prgSize == 1 {
		m.mask = 0x3fff
	} else {
		m.mask = 0x7fff
	}
	m.chr = &rom{mem: m.file.chr[:0x2000], mask: 0x1fff}
	return m
}

func (m *cnrom) ConnectToPPU(bus m6502.Bus) {
	bus.RegisterPort(emulator.PortMask{Mask: 0b1110_000000000000, Value: 0b0000_000000000000}, m.chr)
}

func (m *cnrom) ConnectToCPU(bus m6502.Bus) {
	bus.RegisterPort(emulator.PortMask{Mask: 0b10000000_00000000, Value: 0b10000000_00000000}, m)
}

func (m *cnrom) ReadPort(addr uint16) (byte, bool) {
	return m.file.prg[addr&m.mask], false
}

func (m *cnrom) WritePort(addr uint16, data byte) {
	bank := uint16(data&0x03) * 0x2000
	fmt.Printf("[writePort]-> port:0x%04X data:%v bank:0x%04X \n", addr, data, bank)
	m.chr.mem = m.file.chr[bank : bank+0x2000]
}
