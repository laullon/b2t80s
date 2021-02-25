package cartridge

import "github.com/laullon/b2t80s/cpu/z80"

type ascii16 struct {
	rom   []byte
	mask  byte
	offB1 uint32
	offB2 uint32
}

func NewAscii16(rom []byte) z80.Memory {
	cart := &ascii16{
		rom:   rom,
		mask:  byte(len(rom)/0x4000) - 1,
		offB2: uint32(0x4000),
	}
	return cart
}

func (cart *ascii16) GetByte(addr uint16) byte {
	switch addr >> 14 {
	case 1, 3:
		return cart.rom[cart.offB1+uint32(addr&0x3fff)]
	case 0, 2:
		return cart.rom[cart.offB2+uint32(addr&0x3fff)]
	}
	return 0xff

}

func (cart *ascii16) PutByte(addr uint16, data byte) {
	// fmt.Printf("[ascii16] PutByte(0x%04X, %d) (b:%d) \n", addr, data, data&cart.mask)
	off := uint32(data&cart.mask) * 0x4000
	switch true {
	case addr > 0x5fff && addr < 0x6800:
		cart.offB1 = off
	case addr > 0x6fff && addr < 0x7800:
		cart.offB2 = off
	}
	// bank := byte(addr >> 13)
	// if 0x6000 <= addr && addr < 0xC000 {
	// 	cart.setRom(bank, data)
	// }
}

// TODO: remove
func (cart *ascii16) ReadPort(port uint16) (byte, bool)    { return 0, true }
func (cart *ascii16) WritePort(port uint16, data byte)     {}
func (cart *ascii16) GetBlock(start, length uint16) []byte { panic("not supported") }
