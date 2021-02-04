package nes

import (
	"fmt"

	"github.com/laullon/b2t80s/emulator"
)

type apu struct {
	cpu emulator.CPU

	doIRQ bool
	//
	data []byte
	mask uint16

	frameLength uint
	frameCount  uint
}

func newAPU(cpu emulator.CPU, clock uint) *apu {
	return &apu{
		frameLength: clock / 2 / 50,
		cpu:         cpu,
		data:        make([]byte, 0x08),
		mask:        0x07,
	}
}

func (apu *apu) Tick() {
	if apu.doIRQ {
		apu.cpu.Interrupt(true)
	}

	if apu.frameCount < 3 {
		apu.doIRQ = true
	}

	apu.frameCount++
	if apu.frameCount == apu.frameLength {
		apu.frameCount = 0
	}
}

func (apu *apu) ReadPort(addr uint16) (byte, bool) {
	fmt.Printf("[apu] read  0x%04X\n", addr&apu.mask)
	return apu.data[addr&apu.mask], false
}

func (apu *apu) WritePort(addr uint16, data byte) {
	fmt.Printf("[apu] write 0x%04X 0x%02x\n", addr&apu.mask, data)
	apu.data[addr&apu.mask] = data
}
