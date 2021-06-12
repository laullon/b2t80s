package a1942

import (
	"image/color"
	"math/bits"

	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/gui"
)

// screen_device &set_raw(u32 pixclock, u16 htotal, u16 hbend, u16 hbstart, u16 vtotal, u16 vbend, u16 vbstart)
// m_screen->set_raw(MASTER_CLOCK/2, 384, 128, 0, 262, 22, 246);   // hsync is 50..77, vsync is 257..259

type video struct {
	m *a1942

	spriteram cpu.RAM
	fgvram    cpu.RAM
	bgvram    cpu.RAM

	display *gui.Display
	x, y    uint

	charsRom []byte
	tilesRom [][]byte

	palette     []color.RGBA
	charPalette []byte
	bgPalette   []byte

	scroll uint16
}

func newVideo(m *a1942) *video {
	v := &video{
		display:   gui.NewDisplay(gui.Size{W: 256 * 2, H: 256}),
		m:         m,
		spriteram: cpu.NewRAM(make([]byte, 0x0800), 0x07ff),
		charsRom:  loadRom("sr-02.f2"),
		fgvram:    cpu.NewRAM(make([]byte, 0x0800), 0x07ff),
		bgvram:    cpu.NewRAM(make([]byte, 0x0800), 0x07ff),
	}

	v.display.ViewPortRect = gui.Rect{X: 0, Y: 16, W: 256, H: 224}
	v.display.ViewSize = gui.Size{W: 256, H: 224}

	v.tilesRom = make([][]byte, 3)
	v.tilesRom[0] = append(v.tilesRom[0], loadRom("sr-08.a1")...)
	v.tilesRom[0] = append(v.tilesRom[0], loadRom("sr-09.a2")...)
	v.tilesRom[1] = append(v.tilesRom[1], loadRom("sr-10.a3")...)
	v.tilesRom[1] = append(v.tilesRom[1], loadRom("sr-11.a4")...)
	v.tilesRom[2] = append(v.tilesRom[2], loadRom("sr-12.a5")...)
	v.tilesRom[2] = append(v.tilesRom[2], loadRom("sr-13.a6")...)

	red := loadRom("sb-5.e8")
	green := loadRom("sb-6.e9")
	blue := loadRom("sb-7.e10")
	v.palette = make([]color.RGBA, 256)

	var bit0, bit1, bit2, bit3 byte

	for i := 0; i < 256; i++ {
		// red component
		bit0 = red[i] >> 0 & 1
		bit1 = red[i] >> 1 & 1
		bit2 = red[i] >> 2 & 1
		bit3 = red[i] >> 3 & 1
		r := 0x0e*bit0 + 0x1f*bit1 + 0x43*bit2 + 0x8f*bit3

		// green component
		bit0 = green[i] >> 0 & 1
		bit1 = green[i] >> 1 & 1
		bit2 = green[i] >> 2 & 1
		bit3 = green[i] >> 3 & 1
		g := 0x0e*bit0 + 0x1f*bit1 + 0x43*bit2 + 0x8f*bit3

		// blue component
		bit0 = blue[i] >> 0 & 1
		bit1 = blue[i] >> 1 & 1
		bit2 = blue[i] >> 2 & 1
		bit3 = blue[i] >> 3 & 1
		b := 0x0e*bit0 + 0x1f*bit1 + 0x43*bit2 + 0x8f*bit3

		v.palette[i] = color.RGBA{r, g, b, 0xff}
	}

	v.charPalette = loadRom("sb-0.f1")
	v.bgPalette = loadRom("sb-4.d6")
	return v
}

func (v *video) Tick() {
	v.x++
	if v.x == 384 {
		v.x = 0
		v.y++
		if v.y == 262 {
			v.y = 0
			v.reDraw()
			v.display.Swap()
		}
		switch v.y {
		case 44:
			// v.m.audioCpu.Interrupt(true)
		case 109:
			v.m.mainCpu.Interrupt(true, 0xcf) /* RST 08h */
			// v.m.audioCpu.Interrupt(true)
		case 175:
			// v.m.audioCpu.Interrupt(true)
		case 240:
			v.m.mainCpu.Interrupt(true, 0xd7) /* RST 10h - vblank */
			// v.m.audioCpu.Interrupt(true)
		}
	}
}

func (v *video) reDraw() {
	for col := 0; col < 16; col++ {
		for row := 0; row < 32; row++ {
			realRow := (row + int(v.scroll/16)) % 0x1f
			tileAddr := uint16(col + realRow*32)
			tileIdx, _ := v.bgvram.ReadPort(tileAddr)
			colorInfo, _ := v.bgvram.ReadPort(tileAddr + 0x10)
			tile := uint16(tileIdx) | (uint16(colorInfo&0x80) << 1)
			palette := ((colorInfo & 0x1f) + 0x20) << 3
			v.drawTile(v.display, row, int(v.scroll%16), col, tile, palette, colorInfo&0x20 != 0, colorInfo&0x40 != 0)
		}
	}

	for row := 0; row < 32; row++ {
		for col := 0; col < 32; col++ {
			tileAddr := uint16(col + row*32)
			tileIdx, _ := v.fgvram.ReadPort(tileAddr)
			colorInfo, _ := v.fgvram.ReadPort(tileAddr + 0x0400)
			tile := uint16(tileIdx) | (uint16(colorInfo&0x80) << 1)
			palette := (colorInfo & 0x3f) << 2
			v.drawChar(v.display, int(col), int(row), int(tile), palette)
		}
	}
}

func (v *video) drawTile(display *gui.Display, col, scroll, row int, tile uint16, palette byte, fx, fy bool) {
	var data1, data2, data3 byte
	for y := uint16(0); y < 16; y++ {
		for i := uint16(0); i < 2; i++ {
			if fx {
				data1 = v.tilesRom[0][tile*32+y+(1-i)*16]
				data2 = v.tilesRom[1][tile*32+y+(1-i)*16]
				data3 = v.tilesRom[2][tile*32+y+(1-i)*16]
				data1 = bits.Reverse8(data1)
				data2 = bits.Reverse8(data2)
				data3 = bits.Reverse8(data3)
			} else {
				data1 = v.tilesRom[0][tile*32+y+i*16]
				data2 = v.tilesRom[1][tile*32+y+i*16]
				data3 = v.tilesRom[2][tile*32+y+i*16]
			}
			for x := 0; x < 8; x++ {
				color := data1 & 0b00000001 << 2
				color |= data2 & 0b00000001 << 1
				color |= data3 & 0b00000001 << 0
				if fy {
					display.SetRGBA((7-x)+int(i*8)+col*16-scroll, 15-int(y)+row*16, v.palette[v.bgPalette[color|palette]])
				} else {
					display.SetRGBA((7-x)+int(i*8)+col*16-scroll, int(y)+row*16, v.palette[v.bgPalette[color|palette]])
				}
				data1 >>= 1
				data2 >>= 1
				data3 >>= 1
			}
		}
	}
}

func (v *video) drawChar(display *gui.Display, col, row, tile int, palette byte) {
	for y := 0; y < 8; y++ {
		for i := 0; i < 2; i++ {
			data := v.charsRom[tile*16+y*2+i]
			for x := 0; x < 4; x++ {
				color := data & 0b00000001 << 1
				color |= data & 0b00010000 >> 4
				if color != 0 {
					display.SetRGBA(((3 - x) + (i * 4) + col*8), row*8+y, v.palette[0x80|v.charPalette[color|palette]])
				}
				data >>= 1
			}
		}
	}
}

func (v *video) ReadPort(port uint16) (byte, bool) { return 0xff, false }
func (v *video) WritePort(port uint16, data byte) {
	switch port {
	case 0xc802:
		v.scroll = v.scroll&0xff00 | uint16(data)
	case 0xc803:
		v.scroll = v.scroll&0xffff | uint16(data)<<8
	}
	// TODO: c802-c803 background scroll
	// TODO: c805      background palette bank selector

}