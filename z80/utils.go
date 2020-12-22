package z80

import (
	"fmt"
	"strconv"

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

// TODO cam it be private ?
type RegPair struct {
	h, l *byte
}

func (reg *RegPair) Get() uint16 {
	return uint16(*reg.h)<<8 | uint16(*reg.l)
}

func (reg *RegPair) Set(hl uint16) {
	*reg.h = byte(hl >> 8)
	*reg.l = byte(hl & 0x00ff)
}

func toWord(a, b byte) uint16 {
	return uint16(a) | uint16(b)<<8
}

//*****

type objectPool struct {
	objects []interface{}
	idx     int
}

func (op *objectPool) next() interface{} {
	op.idx++
	op.idx &= 0xf
	return op.objects[op.idx]
}

func newObjectPool(new func() interface{}) *objectPool {
	op := &objectPool{objects: make([]interface{}, 0x10)}
	for i := 0; i < 0x10; i++ {
		op.objects[i] = new()
	}
	return op
}

//******

type circularBuffer struct {
	elemets []z80op
	i, e    byte
}

func newCircularBuffer() *circularBuffer {
	return &circularBuffer{
		elemets: make([]z80op, 0x10),
	}
}

func (cb *circularBuffer) isEmpty() bool {
	return cb.i == cb.e
}

func (cb *circularBuffer) append(ops ...z80op) {
	for _, op := range ops {
		cb.elemets[cb.e] = op
		cb.e++
		cb.e &= 0x0F
	}
}

func (cb *circularBuffer) first() z80op {
	res := cb.elemets[cb.i]
	return res
}

func (cb *circularBuffer) next() {
	cb.i++
	cb.i &= 0x0F
}

// -----------

func (regs *Z80Registers) dump() string {
	return fmt.Sprintf(
		"A:0x%02X F:%08b BC:0x%04X DE:0x%04X HL:0x%04X SP:0x%04X",
		regs.A, regs.F.GetByte(),
		uint16(regs.B)<<8|uint16(regs.C),
		uint16(regs.D)<<8|uint16(regs.E),
		uint16(regs.H)<<8|uint16(regs.L),
		regs.SP.Get())
}

func ParseHexInt8(num string) uint8 {
	r, err := strconv.ParseInt(num, 16, 0)
	if err != nil {
		panic(err)
	}
	return uint8(r)
}
