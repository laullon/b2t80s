package cartridge

import "github.com/laullon/b2t80s/cpu/z80"

type plain struct {
	rom []byte
	off uint16
}

func NewPlain(rom []byte) z80.Memory {
	cart := &plain{
		rom: rom,
	}
	if string(rom[0:2]) == "AB" {
		cart.off = 0x4000
	}
	return cart
}

func (cart *plain) GetByte(addr uint16) byte {
	if int(addr-cart.off) < len(cart.rom) {
		return cart.rom[addr-cart.off]
	}
	return 0xff

}

func (cart *plain) PutByte(addr uint16, data byte) {
	// bank := byte(addr >> 13)
	// if 0x6000 <= addr && addr < 0xC000 {
	// 	cart.setRom(bank, data)
	// }
}

// TODO: remove
func (cart *plain) ReadPort(port uint16) (byte, bool)    { return 0, true }
func (cart *plain) WritePort(port uint16, data byte)     {}
func (cart *plain) GetBlock(start, length uint16) []byte { panic("not supported") }
