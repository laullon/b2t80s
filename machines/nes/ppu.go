package nes

import (
	"fmt"
	"image"
	"image/color"
	"math/bits"

	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
)

type ppu struct {
	cpu     emulator.CPU
	bus     m6502.Bus
	display *image.RGBA
	monitor emulator.Monitor
	h, v    int

	mask byte

	addr    uint16
	addrInc uint16

	scroll  uint16
	scrollX byte
	scrollY byte

	nameTableBase uint16
	patternBase   uint16

	palette *m6502.BasicRam

	enableNMI bool

	oam     []byte
	oamAddr byte

	lastWrite byte

	vblank bool

	spriteBase uint16
	sprite0hit bool
}

// 1662607*3,2 / 50 / 341 = 312,043542522
// 32 x 30 = 256 x 240

func newPPU(bus m6502.Bus, cpu emulator.CPU) *ppu {
	display := image.NewRGBA(image.Rect(0, 0, 256, 240))
	ppu := &ppu{
		cpu:     cpu,
		bus:     bus,
		display: display,
		monitor: emulator.NewMonitor(display),
		palette: &m6502.BasicRam{Data: make([]byte, 0x20), Mask: 0x1f},
		oam:     make([]byte, 0x100),
	}

	// palette
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111_1111_0000_0000, Value: 0b0011_1111_0000_0000}, ppu.palette)

	return ppu
}

var colors = []color.RGBA{
	{84, 84, 84, 255}, {0, 30, 116, 255}, {8, 16, 144, 255}, {48, 0, 136, 255}, {68, 0, 100, 255}, {92, 0, 48, 255}, {84, 4, 0, 255}, {60, 24, 0, 255}, {32, 42, 0, 255}, {8, 58, 0, 255}, {0, 64, 0, 255}, {0, 60, 0, 255}, {0, 50, 60, 255}, {0, 0, 0, 255}, {0, 0, 0, 255}, {0, 0, 0, 255},
	{152, 150, 152, 255}, {8, 76, 196, 255}, {48, 50, 236, 255}, {92, 30, 228, 255}, {136, 20, 176, 255}, {160, 20, 100, 255}, {152, 34, 32, 255}, {120, 60, 0, 255}, {84, 90, 0, 255}, {40, 114, 0, 255}, {8, 124, 0, 255}, {0, 118, 40, 255}, {0, 102, 120, 255}, {0, 0, 0, 255}, {0, 0, 0, 255}, {0, 0, 0, 255},
	{236, 238, 236, 255}, {76, 154, 236, 255}, {120, 124, 236, 255}, {176, 98, 236, 255}, {228, 84, 236, 255}, {236, 88, 180, 255}, {236, 106, 100, 255}, {212, 136, 32, 255}, {160, 170, 0, 255}, {116, 196, 0, 255}, {76, 208, 32, 255}, {56, 204, 108, 255}, {56, 180, 204, 255}, {60, 60, 60, 255}, {0, 0, 0, 255}, {0, 0, 0, 255},
	{236, 238, 236, 255}, {168, 204, 236, 255}, {188, 188, 236, 255}, {212, 178, 236, 255}, {236, 174, 236, 255}, {236, 174, 212, 255}, {236, 180, 176, 255}, {228, 196, 144, 255}, {204, 210, 120, 255}, {180, 222, 120, 255}, {168, 226, 144, 255}, {152, 226, 180, 255}, {160, 214, 228, 255}, {160, 162, 160, 255}, {0, 0, 0, 255}, {0, 0, 0, 255},
}

func (ppu *ppu) Tick() {
	if int(ppu.oam[sY]) == ppu.v {
		ppu.sprite0hit = true
	}
	for i := 0; i < 16; i++ {
		ppu.h++
		if ppu.h == 342 {
			ppu.drawLine()
			ppu.h = 0
			ppu.v++
			ppu.vblank = ppu.v > 241
			if ppu.v == 243 && ppu.enableNMI {
				ppu.cpu.NMI(true)
			}

			if ppu.v == 261 {
				ppu.sprite0hit = false
			}

			if ppu.v == 312 {
				ppu.v = 0
				ppu.drawSprites()
				ppu.monitor.FrameDone()
				// panic(-1)

				ppu.scrollX = byte(ppu.scroll >> 8)
				ppu.scrollY = byte(ppu.scroll)

			}
		}
	}
}

func (ppu *ppu) drawLine() {
	if ppu.mask&0x08 == 0 {
		return
	}
	cRow := (uint16(ppu.v) >> 3) & 0x1f
	y := uint16(ppu.v) & 0x007
	for col := uint16(0); col < 32; col++ {
		base := ppu.nameTableBase
		if col < 32 {
			cCol := col + uint16(ppu.scrollX)>>3
			if cCol >= 32 {
				cCol -= 32
				base ^= 0x0400
			}
			if ppu.v == 0 {
				println("col:", col, "+", uint16(ppu.scrollX)>>3, "=", cCol)
			}
			bCol := cCol >> 2 & 0x07
			bRow := cRow >> 2 & 0x07

			rCol := cCol >> 1 & 0x01
			rRow := cRow >> 1 & 0x01
			region := (rRow << 1) | rCol

			charAddr := base | (cRow << 5) | cCol
			char := uint16(ppu.bus.Read(charAddr))

			patternAddr := ppu.patternBase | char<<4 | y
			pattern0 := ppu.bus.Read(patternAddr)
			pattern1 := ppu.bus.Read(patternAddr | 0x08)

			attrAddr := base | 0x03c0 | (bRow << 3) | bCol
			attr := ppu.bus.Read(attrAddr)
			palette := (attr >> (region * 2)) & 0x03

			for i := 0; i < 8; i++ {
				c := uint16(((pattern0 & 0x80) >> 7) | ((pattern1 & 0x80) >> 6))
				color := uint16(0x3f00)
				if c != 0 {
					color |= uint16(palette) << 2
					color |= c
				}
				pattern0 <<= 1
				pattern1 <<= 1
				x := int(col*8) + i - (int(ppu.scrollX) & 0x07)
				ppu.display.SetRGBA(x, ppu.v, colors[ppu.bus.Read(color)&0x3f])
			}
		}
	}
}

const (
	sY = iota
	sID
	sAttr
	sX
)

// TODO:
//		secondary OAM
//		priority
//		V mirror
//		8x16
//		Flip sprite vertically
//		overlap
func (ppu *ppu) drawSprites() {
	if ppu.mask&0x10 == 0 {
		return
	}
	for sIdx := 0; sIdx < 64; sIdx++ {
		sprite := ppu.oam[sIdx*4 : (sIdx*4)+4]
		if sprite[sY] != 0xff {
			for y := 0; y < 8; y++ {
				patternAddr := ppu.spriteBase | uint16(sprite[sID])<<4 | uint16(y)
				pattern0 := ppu.bus.Read(patternAddr)
				pattern1 := ppu.bus.Read(patternAddr | 0x08)
				if sprite[sAttr]&0x40 == 0x40 {
					pattern0 = bits.Reverse8(pattern0)
					pattern1 = bits.Reverse8(pattern1)
				}

				for i := 0; i < 8; i++ {
					c := ((pattern0 & 0x80) >> 7) | ((pattern1 & 0x80) >> 6)
					if c != 0 {
						color := uint16(0x3f10)
						color |= uint16(sprite[sAttr]&0x3) << 2
						color |= uint16(c)
						ppu.display.SetRGBA(int(sprite[sX])+i, int(sprite[sY])+y+1, colors[ppu.bus.Read(color)&0x3f])
					}
					pattern0 <<= 1
					pattern1 <<= 1
				}
			}
		}
	}
}

func (ppu *ppu) ReadPort(addr uint16) (res byte, skip bool) {
	switch addr & 0x07 {
	case 2: // TODO: sprite bits
		res = ppu.lastWrite & 0x0f
		if ppu.vblank {
			res |= 0x80
		}
		if ppu.sprite0hit {
			res |= 0x40
		}

	case 4:
		res = ppu.oam[ppu.oamAddr]

	case 6:
		res = uint8(ppu.addr)

	case 7:
		// fmt.Printf("[ppu] write -> addr:0x%04X data:%v  \n", ppu.addr, data)
		res = ppu.bus.Read(ppu.addr & 0x3FFF)
		ppu.addr += ppu.addrInc

	default:
		panic(fmt.Sprintf("[ppu] read register %d (0x%04X)\n", addr&0x7, addr))
	}
	return
}

func (ppu *ppu) WritePort(addr uint16, data byte) {
	ppu.lastWrite = data
	switch addr & 0xff {
	case 0:
		ppu.nameTableBase = 0x2000 | (uint16(data&0x3) << 10)
		ppu.patternBase = 0x1000 * (uint16(data&0x10) >> 4)
		ppu.spriteBase = 0x1000 * (uint16(data&0x08) >> 3)
		// fmt.Printf("[ppu] write -> nameTableBase:0x%04X data:%08b  \n", ppu.nameTableBase, data)
		ppu.enableNMI = data&0x80 == 0x80
		if data&0x04 == 0 {
			ppu.addrInc = 1
		} else {
			ppu.addrInc = 32
		}

	case 1:
		ppu.mask = data

	case 2:

	case 3:
		ppu.oamAddr = data
	case 4:
		ppu.oam[ppu.oamAddr] = data

	case 5:
		ppu.scroll <<= 8
		ppu.scroll |= uint16(data)

	case 6:
		ppu.addr <<= 8
		ppu.addr |= uint16(data)
	case 7:
		ppu.bus.Write(ppu.addr&0x3FFF, data)
		ppu.addr += ppu.addrInc
	default:
		panic(fmt.Sprintf("[ppu] write 0x%04X 0x%02x\n", addr, data))

	}
}
