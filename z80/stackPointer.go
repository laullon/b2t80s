package z80

import (
	"fmt"

	"github.com/laullon/b2t80s/emulator"
)

type stackPointer struct {
	mem     emulator.Memory
	pointer uint16
}

func NewStackPointer(mem emulator.Memory) emulator.StackPointer {
	return &stackPointer{
		mem:     mem,
		pointer: 0xffff,
	}
}

func (sp *stackPointer) Set(newSP uint16) {
	sp.pointer = newSP
}

func (sp *stackPointer) Get() uint16 {
	return sp.pointer
}

func (sp *stackPointer) Push(w uint16) {
	sp.pointer -= 2
	sp.mem.PutWord(sp.pointer, w)
}

func (sp *stackPointer) Pop() uint16 {
	v := sp.mem.GetWord(sp.pointer)
	sp.pointer += 2
	return v
}

func (sp *stackPointer) Dump(n int) {
	addr := sp.pointer
	for i := 0; i < n; i++ {
		fmt.Printf("0x%04X: 0x%04X\n", addr, sp.mem.GetWord(addr))
		addr += 2
	}
}
