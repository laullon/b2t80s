package nes

import (
	"fmt"

	"fyne.io/fyne"
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

	readControlers bool
	ctrl0idx       byte
	ctrl0          byte
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
	// if apu.doIRQ {
	// 	apu.cpu.Interrupt(true)
	// }

	// if apu.frameCount < 3 {
	// 	apu.doIRQ = true
	// }

	// apu.frameCount++
	// if apu.frameCount == apu.frameLength {
	// 	apu.frameCount = 0
	// }
}

func (apu *apu) ReadPort(addr uint16) (res byte, skip bool) {
	fmt.Printf("[apu] read  0x%04X\n", addr)
	switch addr {
	case 0x4016:
		res = (apu.ctrl0 >> apu.ctrl0idx) & 0x1
		apu.ctrl0idx = (apu.ctrl0idx + 1) & 0x07
	}

	return
}

func (apu *apu) WritePort(addr uint16, data byte) {
	fmt.Printf("[apu] write 0x%04X 0x%02x\n", addr, data)
	switch addr {
	case 0x4016:
		apu.readControlers = data&1 == 1
	}
}

func (apu *apu) onKeyEvent(key *fyne.KeyEvent) {
	fmt.Println("key:", key.Name)
	switch key.Name {

	case fyne.KeyZ: // A
		apu.ctrl0 ^= 0b00000001
	case fyne.KeyX: // B
		apu.ctrl0 ^= 0b00000010
	case fyne.Key1: //select
		apu.ctrl0 ^= 0b00000100
	case fyne.Key2: // start
		apu.ctrl0 ^= 0b00001000
	case fyne.KeyUp:
		apu.ctrl0 ^= 0b00010000
	case fyne.KeyDown:
		apu.ctrl0 ^= 0b00100000
	case fyne.KeyLeft:
		apu.ctrl0 ^= 0b01000000
	case fyne.KeyRight:
		apu.ctrl0 ^= 0b10000000
	}
}
