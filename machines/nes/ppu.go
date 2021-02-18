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

	writeLacht byte

	addr    uint16
	addrInc uint16

	scrollXv uint8
	scrollYv uint8
	scrollX  uint8
	scrollY  uint8

	redLine bool

	nameTableBase byte
	patternBase   uint16

	palette *m6502.BasicRam

	enableNMI bool

	oam     []byte
	oamAddr byte

	lastWrite byte

	vblank bool

	spriteBase uint16
	sprite0hit bool

	charAddrs [][]uint16
	attrAddrs [][]uint16
	blocks    [][]byte
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

	ppu.charAddrs = make([][]uint16, 64)
	ppu.attrAddrs = make([][]uint16, 64)
	ppu.blocks = make([][]byte, 64)

	for x := 0; x < 64; x++ {
		ppu.charAddrs[x] = make([]uint16, 64)
		ppu.attrAddrs[x] = make([]uint16, 64)
		ppu.blocks[x] = make([]byte, 64)

		for y := 0; y < 64; y++ {
			page := (((y&0x20)>>4 | (x&0x20)>>5) * 4) << 8
			addr := 0x2000 | page

			off := (y&0x1f)<<5 | (x & 0x1f)
			chrAddr := addr | off
			ppu.charAddrs[x][y] = uint16(chrAddr)

			off = ((y >> 2 & 0x07) << 3) | ((x >> 2) & 0x0f)
			attrAddr := addr | 0x03c0 | off
			ppu.attrAddrs[x][y] = uint16(attrAddr)

			off = (y & 2) | (x&2)>>1
			ppu.blocks[x][y] = byte(off)

			// fmt.Printf("0x%04X ", attrAddr)
		}
		// println()
	}

	// panic(-1)
	// palette
	bus.RegisterPort("palette", emulator.PortMask{Mask: 0b1111_1111_0000_0000, Value: 0b0011_1111_0000_0000}, ppu.palette)

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
		if ppu.h == 257 {
			ppu.scrollX = ppu.scrollXv
		}
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
				ppu.scrollY = ppu.scrollYv
				// fmt.Printf("0x%02x 0b%08b %03d\n", ppu.scrollY, ppu.scrollY, ppu.scrollY)
			}

			if ppu.v == 312 {
				ppu.v = 0
				ppu.drawSprites()
				ppu.monitor.FrameDone()
				// println("----")
				// panic(-1)
			}
		}
	}
}

func (ppu *ppu) drawLine() {
	if ppu.mask&0x08 == 0 {
		return
	}
	row := uint8(ppu.v>>3) & 0x1f
	if row >= 31 {
		return
	}
	row = ((ppu.scrollY >> 3) + row)
	if row > 29 {
		row += 2
	}
	row += (ppu.nameTableBase >> 1) * 32

	yOff := uint16(ppu.v) & 0x007

	for x := -8; x < 256+8; x += 8 {
		col := uint8((x>>3)+(int(ppu.scrollX)>>3)) & 0x3f
		col += ((ppu.nameTableBase & 1) * 32)
		col &= 0x3f

		charAddr := ppu.charAddrs[col][row&63]
		char := uint16(ppu.bus.Read(charAddr))

		patternAddr := ppu.patternBase | char<<4 | yOff
		pattern0 := ppu.bus.Read(patternAddr)
		pattern1 := ppu.bus.Read(patternAddr | 0x08)

		attrAddr := ppu.attrAddrs[col][row&63]
		b := ppu.blocks[col][row&63]
		attr := ppu.bus.Read(attrAddr)
		palette := (attr >> (b * 2)) & 0x03

		for i := 0; i < 8; i++ {
			c := uint16(((pattern0 & 0x80) >> 7) | ((pattern1 & 0x80) >> 6))
			colorIdx := uint16(0x3f00) | (uint16(palette) << 2) | c
			pattern0 <<= 1
			pattern1 <<= 1
			imgX := x + i - (int(ppu.scrollX) & 0x07)
			imgY := ppu.v - (int(ppu.scrollY) & 0x07)
			rgb := colors[ppu.bus.Read(colorIdx)&0x3f]
			if ppu.redLine {
				rgb = color.RGBA{0xff, 0, 0, 0xff}
			}
			ppu.display.SetRGBA(imgX, imgY, rgb)
		}
	}
	ppu.redLine = false
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
		ppu.writeLacht = 0

	case 4:
		res = ppu.oam[ppu.oamAddr]

	case 6:
		res = uint8(ppu.addr)

	case 7:
		// fmt.Printf("[ppu] write -> addr:0x%04X data:%v  \n", ppu.addr, data)
		res = ppu.bus.Read(ppu.addr & 0x3FFF)
		ppu.addr += ppu.addrInc

	default:
		// panic(fmt.Sprintf("[ppu] read register %d (0x%04X)\n", addr&0x7, addr))
	}
	return
}

func (ppu *ppu) WritePort(addr uint16, data byte) {
	ppu.lastWrite = data
	switch addr & 0xff {
	case 0:
		ppu.nameTableBase = data & 0x3
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
		if ppu.writeLacht == 0 {
			ppu.scrollXv = data
			ppu.writeLacht = 1
			fmt.Printf("X:%03d 0x%02X %08b  v:%03d \n", data, data, data, ppu.v)
			ppu.redLine = true
		} else {
			ppu.scrollYv = data
			ppu.writeLacht = 0
			// fmt.Printf("V:%03d 0x%02X %08b  v:%03d \n", data, data, data, ppu.v)
		}

	case 6:
		if ppu.writeLacht == 0 {
			ppu.addr = (uint16(data) << 8) | (ppu.addr & 0x00ff)
			ppu.writeLacht = 1
		} else {
			ppu.addr = (ppu.addr & 0xff00) | uint16(data)
			ppu.writeLacht = 0
		}

	case 7:
		ppu.bus.Write(ppu.addr&0x3FFF, data)
		ppu.addr += ppu.addrInc

	default:
		panic(fmt.Sprintf("[ppu] write 0x%04X 0x%02x\n", addr, data))
	}
}
