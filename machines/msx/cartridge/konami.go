package cartridge

import (
	"github.com/laullon/b2t80s/emulator"
)

type konami struct {
	rom       []byte
	banks     []uint32
	banksMask byte
}

func NewKonami(rom []byte) emulator.Memory {
	cart := &konami{
		rom: rom,
	}

	for i := 0; i < 4; i++ {
		cart.banks = append(cart.banks, uint32(i)*0x2000)
	}

	// println("len(cart.rom)", len(cart.rom), (len(cart.rom) / 0x2000))
	cart.banksMask = byte((len(cart.rom) / 0x2000) - 1)

	return cart
}

func (cart *konami) GetByte(addr uint16) byte {
	bank, off, ok := decodeAddr(addr)
	if ok {
		return cart.rom[cart.banks[bank]+off]
	}
	return 0xff
}

func (cart *konami) PutByte(addr uint16, data byte) {
	bank, _, ok := decodeAddr(addr)
	// fmt.Printf("[konami] PutByte(0x%04X, %d(%d)(%d)) (bank:%d) (base:0x%08X)\n", addr, data, data&cart.banksMask, cart.banksMask, bank, uint32(data&cart.banksMask)*0x2000)
	if ok {
		cart.banks[bank] = uint32(data&cart.banksMask) * 0x2000
	}
}

func decodeAddr(addr uint16) (bank byte, offset uint32, ok bool) {
	offset = uint32(addr) & 0x1fff
	ok = true
	switch addr & 0xe000 {
	case 0x4000:
		bank = 0
	case 0x6000:
		bank = 1
	case 0x8000:
		bank = 2
	case 0xa000:
		bank = 3
	default:
		ok = false
	}
	return
}

// TODO: remove
func (cart *konami) ReadPort(port uint16) (byte, bool)    { return 0, true }
func (cart *konami) WritePort(port uint16, data byte)     {}
func (cart *konami) GetBlock(start, length uint16) []byte { panic("not supported") }
