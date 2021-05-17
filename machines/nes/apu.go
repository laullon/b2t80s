package nes

import (
	"fmt"

	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/veandco/go-sdl2/sdl"
)

type apu struct {
	cpu cpu.CPU

	doIRQ bool
	//
	data []byte
	mask uint16

	frameLength byte
	frameCount  byte
	frameIRQ    bool

	ctrlStrobe bool
	ctrl0idx   byte
	ctrl0      byte
	ctrl1idx   byte
	ctrl1      byte

	dmaSrc uint16
	doDMA  bool
	cpuBus m6502.Bus
	ppu    *ppu

	channels       [5]channel
	channelsEnable [5]bool
}

type channel [4]byte

func newAPU(cpu cpu.CPU) *apu {
	return &apu{
		frameLength: 4,
		cpu:         cpu,
		data:        make([]byte, 0x08),
		mask:        0x07,
	}
}

// mode 0:    mode 1:       function
// ---------  -----------  -----------------------------
//  - - - f    - - - - -    IRQ (if bit 6 is clear)
//  - l - l    l - l - -    Length counter and sweep
//  e e e e    e e e e -    Envelope and linear counter

func (apu *apu) Tick() {
	apu.frameCount = (apu.frameCount + 1) % apu.frameLength
	// println("frameCount:", apu.frameCount, "frameLength:", apu.frameLength, "frameIRQ:", apu.frameIRQ)
	if apu.frameLength == 4 {
		switch apu.frameCount {
		case 3:
			if apu.frameIRQ {
				apu.cpu.Interrupt(true)
			}
		}
	} else {

	}

	if apu.doDMA {
		apu.cpu.Wait(true)
		apu.ppu.oam[apu.dmaSrc&0xff] = apu.cpuBus.Read(apu.dmaSrc)
		apu.dmaSrc++
		if apu.dmaSrc&0xff == 0 {
			apu.doDMA = false
			apu.cpu.Wait(false)
		}
	}
}

func (apu *apu) ReadPort(addr uint16) (res byte, skip bool) {
	if addr < 0x4014 {
		chIdx := addr >> 2 & 0x07
		idx := addr & 0x03
		res = apu.channels[chIdx][idx]
	} else {
		switch addr & 0x1f {
		case 0x15: // TODO: DNT21 & DMC
			if apu.frameIRQ {
				res |= 0x40
				apu.frameIRQ = false
			}

		case 0x16:
			res = (apu.ctrl0 >> apu.ctrl0idx) & 0x1
			if !apu.ctrlStrobe {
				apu.ctrl0idx = (apu.ctrl0idx + 1) & 0x07
			}

		case 0x17:
			res = (apu.ctrl1 >> apu.ctrl1idx) & 0x1
			if !apu.ctrlStrobe {
				apu.ctrl1idx = (apu.ctrl1idx + 1) & 0x07
			}

		default:
			// panic(fmt.Sprintf("[apu] read  0x%04X\n", addr&0x1f))
		}
	}
	return
}

func (apu *apu) WritePort(addr uint16, data byte) {
	if addr < 0x4014 {
		chIdx := addr >> 2 & 0x07
		idx := addr & 0x03
		apu.channels[chIdx][idx] = data
	} else {
		switch addr {
		case 0x4014:
			apu.dmaSrc = uint16(data) << 8
			apu.doDMA = true

		case 0x4015:
			for i := 0; i < 5; i++ {
				apu.channelsEnable[i] = (data & 0x01) == 1
				data >>= 1
			}

		case 0x4016:
			apu.ctrlStrobe = data&1 == 1
			apu.ctrl0idx = 0
			apu.ctrl1idx = 0

		case 0x4017:
			apu.frameLength = 4 + (data>>7)&1
			apu.frameIRQ = (data>>6)&1 == 0

		default:
			panic(fmt.Sprintf("[apu] write 0x%04X 0x%02x\n", addr, data))
		}
	}
}

func (apu *apu) OnKey(key sdl.Scancode) {
	switch key {
	case sdl.SCANCODE_Z: // A
		apu.ctrl0 ^= 0b00000001
	case sdl.SCANCODE_X: // B
		apu.ctrl0 ^= 0b00000010
	case sdl.SCANCODE_1: //select
		apu.ctrl0 ^= 0b00000100
	case sdl.SCANCODE_2: // start
		apu.ctrl0 ^= 0b00001000
	case sdl.SCANCODE_UP:
		apu.ctrl0 ^= 0b00010000
	case sdl.SCANCODE_DOWN:
		apu.ctrl0 ^= 0b00100000
	case sdl.SCANCODE_LEFT:
		apu.ctrl0 ^= 0b01000000
	case sdl.SCANCODE_RIGHT:
		apu.ctrl0 ^= 0b10000000
	}
}
