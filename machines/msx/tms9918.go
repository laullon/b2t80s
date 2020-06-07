package msx

// http://www.cs.columbia.edu/~sedwards/papers/TMS9918.pdf

import (
	"fmt"
	"image"
	"image/color"

	"github.com/laullon/b2t80s/emulator"
)

var vdpMasks = []byte{0x03, 0xFB, 0x0F, 0xFF, 0x07, 0x7F, 0x07, 0xFF}

type tms9918 struct {
	cpu emulator.CPU

	status byte

	vramAddr       uint16
	vramByteToRead byte
	vram           []byte

	waitSecondByte bool

	registers []byte

	display *image.RGBA

	m1, m2, m3     bool
	pc, pg, pn     uint16
	sa, sg         uint16
	pcMask, pgMask uint16
	mag, si        bool

	x, y int
}

var palette = []color.RGBA{
	{0, 0, 0, 0xff},
	{0, 0, 0, 0xff},
	{33, 200, 66, 0xff},
	{94, 220, 120, 0xff},
	{84, 85, 237, 0xff},
	{125, 118, 252, 0xff},
	{212, 82, 77, 0xff},
	{66, 235, 245, 0xff},
	{252, 85, 84, 0xff},
	{255, 121, 120, 0xff},
	{212, 193, 84, 0xff},
	{230, 206, 128, 0xff},
	{33, 176, 59, 0xff},
	{201, 91, 186, 0xff},
	{204, 204, 204, 0xff},
	{0xff, 0xff, 0xff, 0xff},
}

func newTMS9918(cpu emulator.CPU) *tms9918 {
	return &tms9918{
		vram:      make([]byte, 0x4000),
		registers: make([]byte, 8),
		display:   image.NewRGBA(image.Rect(-37, -64, 345-37, 313-64)),
		cpu:       cpu,
	}
}

func (vdp *tms9918) ReadPort(port uint16) (res byte, skip bool) {
	skip = false
	switch port & 0xff {
	case 0x98:
		res = vdp.vramByteToRead
		vdp.vramAddr++
		vdp.vramAddr &= 0x3fff
		vdp.vramByteToRead = vdp.vram[vdp.vramAddr]

	case 0x99:
		// println("vdp.status:", vdp.status)
		res = vdp.status
		vdp.status &= 0x3f

	default:
		panic(fmt.Sprintf("[ReadPort] Unsopported port: 0x%02X", port))
	}
	vdp.waitSecondByte = false
	return
}

func (vdp *tms9918) WritePort(port uint16, data byte) {
	switch port & 0xff {
	case 0x98:
		vdp.vram[vdp.vramAddr] = data
		vdp.vramByteToRead = data
		vdp.vramAddr++
		vdp.vramAddr &= 0x3fff
		vdp.waitSecondByte = false

	case 0x99:
		// fmt.Printf("[vdp.writePort]-> port:0x%04X data:0x%04X \n", port, data)
		if vdp.waitSecondByte {
			vdp.vramAddr = ((uint16(data) << 8) | (vdp.vramAddr & 0x00ff)) & 0x3fff
			vdp.waitSecondByte = false
			addrMode := data&0x80 == 0
			if addrMode {
				reading := data&0x40 == 0
				if reading {
					vdp.vramByteToRead = vdp.vram[vdp.vramAddr]
					vdp.vramAddr++
					vdp.vramAddr &= 0x3fff
				}
			} else {
				vdp.registers[data&0x7] = byte(vdp.vramAddr & 0x00ff)
				// println("r:", data&0x7, "=", vdp.registers[data&0x7])
				vdp.update()
			}
		} else {
			vdp.vramAddr = ((vdp.vramAddr & 0xff00) | uint16(data)) & 0x3fff
			vdp.waitSecondByte = true
		}

	default:
		panic(fmt.Sprintf("[WritePort] Unsopported port: 0x%02X", port))
	}
}

func (vdp *tms9918) update() {
	vdp.m1 = vdp.registers[1]&0b00010000 != 0
	vdp.m2 = vdp.registers[1]&0b00001000 != 0
	vdp.m3 = vdp.registers[0]&0b00000010 != 0

	vdp.mag = vdp.registers[1]&0b00000001 != 0
	vdp.si = vdp.registers[1]&0b00000010 != 0

	vdp.sa = uint16(vdp.registers[5]) << 7
	vdp.sg = uint16(vdp.registers[6]) << 11

	vdp.pn = uint16(vdp.registers[2]) << 10

	if vdp.m3 {
		vdp.pc = uint16(vdp.registers[3]&0x80) * 0x40
		vdp.pg = uint16(vdp.registers[4]&0x04) * 0x800
		vdp.pcMask = (uint16(vdp.registers[3]&0x7f) << 3) | 7
		vdp.pgMask = (uint16(vdp.registers[4]&0x03) << 8) | (vdp.pcMask & 0xff)
	} else {
		vdp.pc = uint16(vdp.registers[3]) * 0x40
		vdp.pg = uint16(vdp.registers[4]) * 0x800
	}

	println("[VPD] m1:", vdp.m1, "m2:", vdp.m2, "m3:", vdp.m3)
	fmt.Printf("[VDP] pn:0x%04X(%d) pc:0x%04X(%d) pg:0x%04X(%d)\n", vdp.pn, vdp.registers[2], vdp.pc&0x2000, vdp.registers[3], vdp.pg&0x2000, vdp.registers[4])
}

var c = uint16(0)

func (vdp *tms9918) Tick() {
	for i := 0; i < 3; i++ {
		c := vdp.getRasteColor()
		vdp.display.SetRGBA(vdp.x, vdp.y, palette[c])

		vdp.x++
		if vdp.x == 342-37 {
			vdp.x = -37
			vdp.drawSprites()
			vdp.y++

			if vdp.y == 313-64 {
				vdp.y = -64
				vdp.status = 0x80
				if vdp.registers[1]&0x20 != 0 {
					vdp.cpu.Interrupt(true)
				}
			}
		}
	}
}

func (vdp *tms9918) drawSprites() {
	sprites := make([][]byte, 0)
	height := 8
	if vdp.si {
		height = 16
	}

	for idx := uint16(0); idx < 32; idx++ {
		sprite := vdp.vram[vdp.sa+(idx*4) : vdp.sa+(idx*4)+4]
		if sprite[0] == 209 {
			break
		}

		y := int(sprite[0])
		if (vdp.y >= y) && (vdp.y < y+height) {
			sprites = append([][]byte{sprite}, sprites...)
		}
		// todo: check sprite5
	}

	for _, sprite := range sprites {
		if !vdp.si {
			vdp.drawSprite(sprite)
		} else {
			vdp.drawSpriteSI(sprite)
		}
	}
}

func (vdp *tms9918) drawSprite(sprite []byte) {
	// for y := uint16(0); y < 8; y++ {
	// 	b := vdp.vram[vdp.sg+uint16(sprite[2])<<3+y]
	// 	for x := 0; x < 8; x++ {
	// 		if b&(1<<(7-x)) != 0 {
	// 			vdp.display.SetRGBA(int(sprite[1])+x, int(sprite[0])+int(y), palette[sprite[3]&0x0f])
	// 		}
	// 	}
	// }
}

func (vdp *tms9918) drawSpriteSI(sprite []byte) {
	for i := uint16(0); i < 2; i++ {
		y := uint16(vdp.y - int(sprite[0]))
		b := vdp.vram[vdp.sg+(uint16(sprite[2])&252)<<3+y+(i*16)]
		for x := 0; x < 8; x++ {
			if b&(1<<(7-x)) != 0 {
				vdp.display.SetRGBA(int(sprite[1])+x+int(i*8), int(sprite[0])+int(y), palette[sprite[3]&0x0f])
			}
		}
	}
}

func (vdp *tms9918) getRasteColor() byte {
	col := 0
	row := 0
	bidx := 0
	b := byte(0)
	d := false

	if vdp.x < 0 || vdp.y < 0 {
		return vdp.registers[7] & 0x0f
	}

	switch true {
	case vdp.m1:
		col = vdp.x / 6
		row = vdp.y / 8
		d = col >= 0 && col < 40 && row >= 0 && row < 24
		bidx = 7 - (vdp.x % 6)

	default:
		col = vdp.x / 8
		row = vdp.y / 8
		d = col >= 0 && col < 32 && row >= 0 && row < 24
		bidx = 7 - (vdp.x % 8)
	}

	if !d {
		return vdp.registers[7] & 0x0f
	}

	c := byte(0)
	switch true {
	case vdp.m1: // Text Mode
		pn := uint16(vdp.registers[2])*0x400 + uint16(row*40+col)
		pg := uint16(vdp.registers[4]) * 0x800
		char := uint16(vdp.vram[pn])
		b = vdp.vram[pg+(char*8)+uint16(vdp.y%8)]
		c = vdp.registers[7]

	case vdp.m2:
		panic(1)

	case vdp.m3: // Bitmap mode (Graphics II)
		charPos := uint16(row*32 + col)

		char := uint16(vdp.vram[vdp.pn+charPos]) + uint16(vdp.y>>6)<<8

		pgChar := ((char & vdp.pgMask) << 3) + uint16(vdp.y&0x07)
		pcChar := ((char & vdp.pcMask) << 3) + uint16(vdp.y&0x07)

		b = vdp.vram[vdp.pg|pgChar]
		c = vdp.vram[vdp.pc|pcChar]

	default: // Standard mode (Graphic I)
		charPos := uint16(row*32 + col)
		char := uint16(vdp.vram[vdp.pn+charPos])
		b = vdp.vram[vdp.pg+(char*8)+uint16(vdp.y%8)]
		c = vdp.vram[vdp.pc+(char/8)]
	}

	if b&(1<<bidx) != 0 {
		return c & 0xf0 >> 4
	}
	return c & 0x0f
}
