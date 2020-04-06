package zx

type bank []byte
type page *bank

type MemoryMode int

const (
	ZX48K   MemoryMode = 0
	ZX128K  MemoryMode = 1
	ZXPLUS3 MemoryMode = 2
)

type memory struct {
	roms []bank
	rom  byte

	screens       []page
	banks         []bank
	pages         []page
	mode          MemoryMode
	pagingDisable bool
}

func NewMemory(mode MemoryMode) *memory {
	res := &memory{
		mode:          mode,
		pagingDisable: false,
	}

	switch mode {
	case ZX48K:
		res.pages = make([]page, 4)
		res.roms = append(res.roms, make(bank, 0x4000))
		res.pages[0] = &res.roms[0]
		for p := 0; p < 3; p++ {
			res.banks = append(res.banks, make(bank, 0x4000))
			res.pages[p+1] = &res.banks[p]
		}

	case ZX128K, ZXPLUS3:
		res.pages = make([]page, 4)
		res.screens = make([]page, 2)

		for p := 0; p < 8; p++ {
			res.banks = append(res.banks, make(bank, 0x4000))
		}
		for p := 0; p < 4; p++ {
			res.roms = append(res.roms, make(bank, 0x4000))
		}

		res.screens[0] = &res.banks[5]
		res.screens[1] = &res.banks[7]

		res.pages[0] = &res.roms[0]
		res.pages[1] = res.screens[0]
		res.pages[2] = &res.banks[2]
		res.pages[3] = &res.banks[0]

	default:
		panic("wrong Memory Mode")
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
	page, pos := mem.decodeAddress(addr)
	return (*mem.pages[page])[pos]
}

func (mem *memory) PutByte(addr uint16, b byte) {
	if addr > 0x3fff { // TODO: review for plus
		page, pos := mem.decodeAddress(addr)
		(*mem.pages[page])[pos] = b
	}
}

func (mem *memory) LoadRom(idx int, rom []byte) {
	copy(mem.roms[idx], rom)
	// log.Printf("loaded %v bytes into ROM memory\n", n)
}

func (mem *memory) decodeAddress(addr uint16) (page, pos uint16) {
	page = addr >> 14
	pos = addr & 0x3fff
	return
}

func (mem *memory) ReadPort(port uint16) (byte, bool) { return 0, true }
func (mem *memory) WritePort(port uint16, data byte) {
	switch port {
	case 0x1ffd:
		mem.secPaging(data)
	case 0x7ffd:
		mem.paging(data)
	}
}

func (mem *memory) secPaging(config byte) {
	if config&1 == 0 { // not Special paging
		if (config & 0b00000100) != 0 {
			mem.rom = (mem.rom & 0b01) | 0b10
		} else {
			mem.rom = (mem.rom & 0b01)
		}
		mem.pages[0] = &mem.roms[mem.rom]
	} else { // Special paging
		cfg := config & 0x6 >> 1
		switch cfg {
		case 0:
			mem.pages[0] = &mem.banks[0]
			mem.pages[1] = &mem.banks[1]
			mem.pages[2] = &mem.banks[2]
			mem.pages[3] = &mem.banks[4]
		case 1:
			mem.pages[0] = &mem.banks[4]
			mem.pages[1] = &mem.banks[5]
			mem.pages[2] = &mem.banks[6]
			mem.pages[3] = &mem.banks[7]
		case 2:
			mem.pages[0] = &mem.banks[4]
			mem.pages[1] = &mem.banks[5]
			mem.pages[2] = &mem.banks[6]
			mem.pages[3] = &mem.banks[3]
		case 3:
			mem.pages[0] = &mem.banks[4]
			mem.pages[1] = &mem.banks[7]
			mem.pages[2] = &mem.banks[6]
			mem.pages[3] = &mem.banks[3]
		default:
			panic("--")
		}
	}
}

func (mem *memory) paging(config byte) {
	if mem.pagingDisable {
		return
	}

	mem.pagingDisable = ((config >> 5) & 0x1) == 1

	if (config & 0b00010000) != 0 {
		mem.rom = (mem.rom & 0b10) | 0b01
	} else {
		mem.rom = (mem.rom & 0b10)
	}
	mem.pages[0] = &mem.roms[mem.rom]

	// screem := (config >> 3) & 0x1
	// mem.pages[1] = mem.screens[screem]

	bank := config & 0b00000111
	mem.pages[3] = &mem.banks[bank]

	// log.Printf("rom:%v screem:%v bank:%v\n", rom, screem, bank)
}
