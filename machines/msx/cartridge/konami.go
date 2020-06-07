package cartridge

import (
	"github.com/laullon/b2t80s/emulator"
)

type konami struct {
	rom   []byte
	banks [][]byte
}

func NewKonami(rom []byte) emulator.Memory {
	cart := &konami{
		rom:   rom,
		banks: make([][]byte, 7),
	}

	for i := byte(2); i < 6; i++ {
		cart.setRom(i, i-2)
	}

	return cart
}

func (cart *konami) GetByte(addr uint16) byte {
	if 0x4000 <= addr && addr < 0xC000 {
		bank := byte(addr >> 13)
		addr &= 0x1fff
		mem := cart.banks[bank]
		return mem[addr]
	}
	return 0xff

}

func (cart *konami) PutByte(addr uint16, data byte) {
	bank := byte(addr >> 13)
	if 0x6000 <= addr && addr < 0xC000 {
		cart.setRom(bank, data)
	}
}

func (cart *konami) setRom(bank, data byte) {
	data %= byte(len(cart.rom) / 0x2000)
	s := uint(data) * 0x2000
	if len(cart.rom) > int(s+0x2000) {
		if len(cart.banks[bank]) == 0 {
			cart.banks[bank] = make([]byte, 0x2000)
		}
		copy(cart.banks[bank], cart.rom[s:s+0x2000])
	}
}

// TODO: remove
func (cart *konami) ReadPort(port uint16) (byte, bool)    { return 0, true }
func (cart *konami) WritePort(port uint16, data byte)     {}
func (cart *konami) GetBlock(start, length uint16) []byte { panic("not supported") }
