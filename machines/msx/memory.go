package msx

import "github.com/laullon/b2t80s/cpu/z80"

type memory struct {
	slots []z80.Memory
	cfg   []byte
}

func NewMemory(rom rom) *memory {
	mem := &memory{cfg: []byte{0, 0, 0, 0}}

	mem.slots = make([]z80.Memory, 4)

	mem.slots[0] = rom
	mem.slots[1] = make(emptySlot, 0)
	mem.slots[2] = make(emptySlot, 0)
	mem.slots[3] = make(ram, 0x010000)

	return mem
}

func (mem *memory) setCartridge1(cart z80.Memory) {
	mem.slots[1] = cart
}

func (mem *memory) GetBlock(start, length uint16) []byte {
	res := make([]byte, length)
	for i := uint16(0); i < length; i++ {
		res[i] = mem.GetByte(start + i)
	}
	return res
}

func (mem *memory) GetByte(addr uint16) byte {
	slot := mem.getSlot(addr)
	return mem.slots[slot].GetByte(addr)
}

func (mem *memory) PutByte(addr uint16, b byte) {
	slot := mem.getSlot(addr)
	mem.slots[slot].PutByte(addr, b)
}

func (mem *memory) getSlot(addr uint16) byte {
	page := int(addr >> 14)
	return mem.cfg[page]
}

func (mem *memory) ReadPort(port uint16) (byte, bool) { return 0, true }
func (mem *memory) WritePort(port uint16, data byte)  {}

// -------
// - ROM -
// -------

type rom []byte

func (rom rom) GetByte(addr uint16) byte {
	if addr < uint16(len(rom)) {
		return rom[addr]
	}
	return 0xff
}

func (rom rom) PutByte(addr uint16, b byte) {
}

// TODO: remove
func (rom rom) ReadPort(port uint16) (byte, bool)    { return 0, true }
func (rom rom) WritePort(port uint16, data byte)     {}
func (rom rom) GetBlock(start, length uint16) []byte { panic("not supported") }

// -------
// - RAM -
// -------

type ram []byte

func (ram ram) GetByte(addr uint16) byte {
	return ram[addr]
}

func (ram ram) PutByte(addr uint16, b byte) {
	ram[addr] = b
}

// TODO: remove
func (ram ram) ReadPort(port uint16) (byte, bool)    { return 0, true }
func (ram ram) WritePort(port uint16, data byte)     {}
func (ram ram) GetBlock(start, length uint16) []byte { panic("not supported") }

// -------
// - Empty Slot -
// -------

type emptySlot []byte

func (es emptySlot) GetByte(addr uint16) byte    { return 0xff }
func (es emptySlot) PutByte(addr uint16, b byte) {}

// TODO: remove
func (es emptySlot) ReadPort(port uint16) (byte, bool)    { return 0, true }
func (es emptySlot) WritePort(port uint16, data byte)     {}
func (es emptySlot) GetBlock(start, length uint16) []byte { panic("not supported") }
