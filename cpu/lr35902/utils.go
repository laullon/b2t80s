package lr35902

import (
	"fmt"
	"strconv"
)

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
	elemets []lr35902op
	i, e    byte
}

func newCircularBuffer() *circularBuffer {
	return &circularBuffer{
		elemets: make([]lr35902op, 0x10),
	}
}

func (cb *circularBuffer) isEmpty() bool {
	return cb.i == cb.e
}

func (cb *circularBuffer) append(ops ...lr35902op) {
	for _, op := range ops {
		cb.elemets[cb.e] = op
		cb.e++
		cb.e &= 0x0F
	}
}

func (cb *circularBuffer) first() lr35902op {
	return cb.elemets[cb.i]
}

func (cb *circularBuffer) next() {
	cb.i++
	cb.i &= 0x0F
}

// -----------

func (regs *LR35902Registers) dump() string {
	return fmt.Sprintf(
		"A:0x%02X F:%08b BC:0x%04X DE:0x%04X HL:0x%04X SP:0x%04X",
		regs.A, regs.F.GetByte(),
		uint16(regs.B)<<8|uint16(regs.C),
		uint16(regs.D)<<8|uint16(regs.E),
		uint16(regs.H)<<8|uint16(regs.L),
		regs.SP.Get())
}

func parseHexUInt8(num string) uint8 {
	r, err := strconv.ParseInt(num, 16, 0)
	if err != nil {
		panic(err)
	}
	return uint8(r)
}

func parseHexUInt16(num string) uint16 {
	r, err := strconv.ParseInt(num, 16, 0)
	if err != nil {
		panic(err)
	}
	return uint16(r)
}
