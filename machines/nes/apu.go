package nes

import (
	"github.com/laullon/b2t80s/emulator"
)

type apu struct {
	cpu emulator.CPU

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
	if apu.frameCount == 0 {
		println("---")
		apu.cpu.Interrupt(true)
	}

	apu.frameCount++
	if apu.frameCount == apu.frameLength {
		apu.frameCount = 0
	}
}

func (apu *apu) ReadPort(addr uint16) (byte, bool) { return apu.data[addr&apu.mask], false }
func (apu *apu) WritePort(addr uint16, data byte)  { apu.data[addr&apu.mask] = data }
