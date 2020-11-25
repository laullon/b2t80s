package cpc

import (
	"fmt"
)

type bank []byte

type memory struct {
	lowerRom       bank
	lowerRomEnable bool

	upperRoms      []bank
	upperRomEnable bool
	upperRomIdx    byte

	banks       []bank
	activeBanks []byte

	safe bool
}

func NewCPCMemory() *memory {
	res := &memory{
		safe:           true,
		lowerRom:       make(bank, 0x4000),
		lowerRomEnable: true,
		upperRoms:      make([]bank, 256),
		upperRomEnable: false,
		banks:          make([]bank, 8),
		activeBanks:    []byte{0, 1, 2, 3},
	}

	for idx := 0; idx < 8; idx++ {
		res.banks[idx] = make(bank, 0x4000)
	}

	return res
}

func (mem *memory) Paging(config byte) {
	cfg := config & 0x7
	switch cfg {
	case 0:
		mem.activeBanks = []byte{0, 1, 2, 3}
	case 1:
		mem.activeBanks = []byte{0, 1, 2, 7}
	case 2:
		mem.activeBanks = []byte{4, 5, 6, 7}
	case 3:
		mem.activeBanks = []byte{0, 3, 2, 7}
	case 4:
		mem.activeBanks = []byte{0, 4, 2, 3}
	case 5:
		mem.activeBanks = []byte{0, 5, 2, 3}
	case 6:
		mem.activeBanks = []byte{0, 6, 2, 3}
	case 7:
		mem.activeBanks = []byte{0, 7, 2, 3}
	}
	// fmt.Printf("[mem] mem.activeBanks: %v\n", mem.activeBanks)
}

func (mem *memory) decodeAddress(addr uint16) (page byte, bank byte, pos uint16) {
	page = byte(addr >> 14)
	bank = mem.activeBanks[page]
	pos = addr & 0x3fff
	return
}

func (mem *memory) GetBlock(start, length uint16) []byte {
	res := make([]byte, length)
	for i := uint16(0); i < length; i++ {
		res[i] = mem.GetByte(start + i)
	}
	return res
}

func (mem *memory) GetByte(addr uint16) byte {
	page, bank, pos := mem.decodeAddress(addr)

	if page == 0 && mem.lowerRomEnable {
		return mem.lowerRom[pos]
	}
	if page == 3 && mem.upperRomEnable {
		return mem.upperRoms[mem.upperRomIdx][pos]
	}

	return mem.banks[bank][pos]
}

func (mem *memory) getScreenByte(addr uint16) byte {
	_, bank, pos := mem.decodeAddress(addr)
	return mem.banks[bank][pos]
}

func (mem *memory) PutByte(addr uint16, b byte) {
	// fmt.Printf("-> addr:0x%08x b:0x%02x \n", addr, b)
	_, bank, pos := mem.decodeAddress(addr)
	mem.banks[bank][pos] = b
}

func (mem *memory) GetWord(addr uint16) uint16 {
	res := uint16(mem.GetByte(addr+1)) << 8
	res |= uint16(mem.GetByte(addr))
	return res
}

func (mem *memory) PutWord(addr, w uint16) {
	mem.PutByte(addr, uint8(w&0x00ff))
	mem.PutByte(addr+1, uint8(w>>8))
}

func (mem *memory) LoadRom(idx int, rom []byte) {
	if idx == -1 {
		copy(mem.lowerRom, rom)
	} else {
		mem.upperRoms[idx] = make(bank, 0x4000)
		copy(mem.upperRoms[idx], rom)
	}
}

func (mem *memory) ReadPort(port uint16) (byte, bool) { return 0, true }
func (mem *memory) WritePort(port uint16, data byte) {
	switch {
	case port&0x2000 == 0:
		if len(mem.upperRoms[data]) > 0 {
			mem.upperRomIdx = data
		} else {
			mem.upperRomIdx = 0
		}
		// println("upperRomIdx", mem.upperRomIdx, data, mem.upperRomEnable)
	default:
		panic(fmt.Sprintf("[memory] bad port 0x%04X", port))
	}
}
