package z80

import (
	"github.com/laullon/b2t80s/emulator"
)

func getWord(mem emulator.Memory, addr uint16) uint16 {
	res := uint16(mem.GetByte(addr))
	res |= uint16(mem.GetByte(addr+1)) << 8
	return res
}

func putWord(mem emulator.Memory, addr, w uint16) {
	mem.PutByte(addr, uint8(w&0x00ff))
	mem.PutByte(addr+1, uint8(w>>8))
}

type regPair struct {
	h, l *byte
}

func (reg *regPair) get() uint16 {
	return uint16(*reg.h)<<8 | uint16(*reg.l)
}

func (reg *regPair) set(hl uint16) {
	*reg.h = byte(hl >> 8)
	*reg.l = byte(hl & 0x00ff)
}
