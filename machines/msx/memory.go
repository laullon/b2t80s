package msx

import (
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/cpu/z80"
)

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
		res[i] = mem.Read(start + i)
	}
	return res
}

func (mem *memory) Read(addr uint16) byte {
	slot := mem.getSlot(addr)
	return mem.slots[slot].Read(addr)
}

func (mem *memory) Write(addr uint16, b byte) {
	slot := mem.getSlot(addr)
	mem.slots[slot].Write(addr, b)
}

func (mem *memory) DumpMap() string {
	panic(-1)
}

func (mem *memory) GetDumplables() map[string]cpu.Dumpable {
	panic(-1)
}

func (mem *memory) RegisterPort(name string, mask cpu.PortMask, manager cpu.PortManager) {
	panic(-1)
}

func (mem *memory) getSlot(addr uint16) byte {
	page := int(addr >> 14)
	return mem.cfg[page]
}

func (mem *memory) ReadPort(port uint16) byte        { return 0 }
func (mem *memory) WritePort(port uint16, data byte) {}

// -------
// - ROM -
// -------

type rom []byte

func (rom rom) Read(addr uint16) byte {
	if addr < uint16(len(rom)) {
		return rom[addr]
	}
	return 0xff
}

func (rom rom) Write(addr uint16, b byte) {
}

// TODO: remove
func (rom rom) ReadPort(port uint16) byte            { return 0 }
func (rom rom) WritePort(port uint16, data byte)     {}
func (rom rom) GetBlock(start, length uint16) []byte { panic("not supported") }

// -------
// - RAM -
// -------

type ram []byte

func (ram ram) Read(addr uint16) byte {
	return ram[addr]
}

func (ram ram) Write(addr uint16, b byte) {
	ram[addr] = b
}

// TODO: remove
func (ram ram) ReadPort(port uint16) byte            { return 0 }
func (ram ram) WritePort(port uint16, data byte)     {}
func (ram ram) GetBlock(start, length uint16) []byte { panic("not supported") }

// -------
// - Empty Slot -
// -------

type emptySlot []byte

func (es emptySlot) Read(addr uint16) byte     { return 0xff }
func (es emptySlot) Write(addr uint16, b byte) {}

// TODO: remove
func (es emptySlot) ReadPort(port uint16) byte            { return 0 }
func (es emptySlot) WritePort(port uint16, data byte)     {}
func (es emptySlot) GetBlock(start, length uint16) []byte { panic("not supported") }
