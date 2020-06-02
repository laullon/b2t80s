package msx

type bank []byte

type memory struct {
	slots     [][]bank
	cfg       []byte
	cartridge bool
}

func NewMemory() *memory {
	res := &memory{cfg: []byte{0, 0, 0, 0}}

	res.slots = make([][]bank, 4)

	// roms
	for p := 0; p < 2; p++ {
		res.slots[0] = append(res.slots[0], make(bank, 0x4000))
	}

	// memory
	for p := 0; p < 4; p++ {
		res.slots[3] = append(res.slots[3], make(bank, 0x4000))
	}

	return res
}

func (mem *memory) GetBlock(start, length uint16) []byte {
	res := make([]byte, length)
	for i := uint16(0); i < length; i++ {
		res[i] = mem.GetByte(start + i)
	}
	return res
}

func (mem *memory) GetByte(addr uint16) byte {
	slot, page, pos := mem.decodeAddress(addr)
	if (slot == 0 && page > 1) || (slot == 1 && !mem.cartridge) || slot == 2 {
		return 0x00
	}
	return mem.slots[slot][page][pos]
}

func (mem *memory) PutByte(addr uint16, b byte) {
	slot, page, pos := mem.decodeAddress(addr)
	if slot != 3 {
		return
	}
	mem.slots[slot][page][pos] = b
}

func (mem *memory) LoadRom(idx int, rom []byte) {
	copy(mem.slots[0][idx], rom)
}

func (mem *memory) LoadCartridge(rom []byte) {
	for p := 0; p < 4; p++ {
		mem.slots[1] = append(mem.slots[1], make(bank, 0x4000))
	}

	for i := 0; i < len(mem.slots[1]); i++ {
		if len(rom) > 0x4000*i {
			copy(mem.slots[1][i+1], rom[0x4000*i:])
		}
	}
	mem.cartridge = true
}

func (mem *memory) decodeAddress(addr uint16) (slot byte, page int, pos uint16) {
	page = int(addr >> 14)
	pos = addr & 0x3fff
	slot = mem.cfg[page]
	// println("mem", slot, page, pos)
	return
}

func (mem *memory) ReadPort(port uint16) (byte, bool) { return 0, true }
func (mem *memory) WritePort(port uint16, data byte)  {}
