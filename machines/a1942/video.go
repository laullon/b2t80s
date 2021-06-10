package a1942

import (
	"image/color"

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
	tilesRom []byte

	palette     []color.RGBA
	charPalette []byte
}

func newVideo(m *a1942) *video {
	v := &video{
		display:   gui.NewDisplay(gui.Size{W: 256, H: 256}),
		m:         m,
		spriteram: cpu.NewRAM(make([]byte, 0x0800), 0x07ff),
		charsRom:  loadRom("sr-02.f2"),
		fgvram:    cpu.NewRAM(make([]byte, 0x0800), 0x07ff),
		bgvram:    cpu.NewRAM(make([]byte, 0x0800), 0x07ff),
	}

	v.tilesRom = append(v.tilesRom, loadRom("sr-08.a1")...)
	v.tilesRom = append(v.tilesRom, loadRom("sr-09.a2")...)
	v.tilesRom = append(v.tilesRom, loadRom("sr-10.a3")...)
	v.tilesRom = append(v.tilesRom, loadRom("sr-11.a4")...)
	v.tilesRom = append(v.tilesRom, loadRom("sr-12.a5")...)
	v.tilesRom = append(v.tilesRom, loadRom("sr-13.a6")...)

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
	return v
}

func (v *video) Tick() {
	v.x++
	if v.x == 384 {
		v.x = 0
		v.y++
		if v.y == 262 {
			v.y = 0
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

	if v.x%8 == 0 && v.y%8 == 0 {
		col := uint16(v.x / 8)
		row := uint16(v.y / 8)
		tileAddr := (col + row*32)
		tileIdx, _ := v.fgvram.ReadPort(tileAddr)
		colorInfo, _ := v.fgvram.ReadPort(tileAddr + 0x0400)
		tile := uint16(tileIdx) | (uint16(colorInfo&0x80) << 1)
		palette := (colorInfo & 0x3f) << 2
		v.drawChar(v.display, int(col), int(row), int(tile), palette)
	}
}

func (v *video) drawChar(display *gui.Display, col, row, tile int, palette byte) {
	for y := 0; y < 8; y++ {
		for i := 0; i < 2; i++ {
			data := v.charsRom[tile*16+y*2+i]
			for x := 0; x < 4; x++ {
				color := data & 0b00000001 << 1
				color |= data & 0b00010000 >> 4
				display.SetRGBA(row*8+y, 255-((3-x)+(i*4)+col*8), v.palette[0x80|v.charPalette[color|palette]])
				data >>= 1
			}
		}
	}
}

func (v *video) ReadPort(port uint16) (byte, bool) { return 0xff, false }
func (v *video) WritePort(port uint16, data byte) {
	// TODO: c802-c803 background scroll
	// TODO: c805      background palette bank selector

}
