package msx

// http://www.cs.columbia.edu/~sedwards/papers/TMS9918.pdf

import (
	"fmt"
	"image"
	"image/color"

	"github.com/laullon/b2t80s/emulator"
)

type tms9918 struct {
	cpu emulator.CPU

	status byte

	vramAddr       uint16
	vramByteToRead byte
	vram           []byte

	waitSecondByte bool

	registers    []byte
	valueToWrite byte

	display *image.RGBA

	m1, m2, m3 bool

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
			vdp.waitSecondByte = false
			addrMode := data&0x80 == 0
			if addrMode {
				vdp.vramAddr = uint16(vdp.valueToWrite) | uint16(data&0x3f)<<8
				reading := data&0x40 == 0
				if reading {
					vdp.vramByteToRead = vdp.vram[vdp.vramAddr]
					vdp.vramAddr++
					vdp.vramAddr &= 0x3fff
				}
			} else {
				vdp.registers[data&0x7] = vdp.valueToWrite
				println("r:", data&0x7, "=", vdp.valueToWrite)
				vdp.update()
			}
		} else {
			vdp.waitSecondByte = true
			vdp.valueToWrite = data
		}

	default:
		panic(fmt.Sprintf("[WritePort] Unsopported port: 0x%02X", port))
	}
}

func (vdp *tms9918) update() {
	vdp.m1 = vdp.registers[1]&0b00010000 != 0
	vdp.m2 = vdp.registers[1]&0b00001000 != 0
	vdp.m3 = vdp.registers[0]&0b00000010 != 0

	println("[VPD] m1:", vdp.m1, "m2:", vdp.m2, "m3:", vdp.m3)
}

var c = uint16(0)

func (vdp *tms9918) Tick() {
	for i := 0; i < 3; i++ {
		c := vdp.getScreenData()
		vdp.display.SetRGBA(vdp.x, vdp.y, palette[c])

		vdp.x++
		if vdp.x == 342-37 {
			vdp.x = -37
			vdp.y++
			if vdp.y == 313-64 {
				vdp.y = -64
				vdp.status = 0x80
				vdp.cpu.Interrupt(true)
			}
		}
	}
}

func (vdp *tms9918) getScreenData() byte {
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
		part := uint16(row / 8)
		pn := (uint16(vdp.registers[2]) * 0x400) + uint16(row*32+col)
		pc := ((uint16(vdp.registers[3]) & 0x80) * 0x40) + (part * 0x800)
		pg := ((uint16(vdp.registers[4]) & 0x40) * 0x800) + (part * 0x800)
		char := uint16(vdp.vram[pn])
		b = vdp.vram[pg+(char*8)+uint16(vdp.y%8)]
		c = vdp.vram[pc+(char*8)]

	default: // Standard mode (Graphic I)
		pn := uint16(vdp.registers[2])*0x400 + uint16(row*32+col)
		pc := uint16(vdp.registers[3]) * 0x40
		pg := uint16(vdp.registers[4]) * 0x800
		char := uint16(vdp.vram[pn])
		b = vdp.vram[pg+(char*8)+uint16(vdp.y%8)]
		c = vdp.vram[pc+(char/8)]
	}

	if b&(1<<bidx) != 0 {
		return c & 0xf0 >> 4
	}
	return c & 0x0f
}
