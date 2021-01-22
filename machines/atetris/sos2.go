package atetris

import (
	"image"
	"image/color"

	"github.com/laullon/b2t80s/emulator"
)

type sos2 struct {
	cpu     emulator.CPU
	v, h    int
	vram    []byte
	color   []byte
	rom     []byte
	display *image.RGBA
	monitor emulator.Monitor

	hBlank *bool
}

var ic byte

func (d *sos2) Tick() {
	col := d.h >> 3
	row := d.v >> 3
	x := d.h & 7
	y := d.v & 7
	b := d.h & 1

	charAddr := uint16(row) << 6
	charAddr |= uint16(col)

	// charAddr := row*32 + col

	// if d.h < 16 && d.v%8 == 0 {
	// 	fmt.Printf("h:%03d(%02d) v:%03d(%02d) charAddr:0x%04X \n", d.h, col, d.v, row, charAddr)
	// }

	charData1 := d.vram[charAddr*2]
	charData2 := d.vram[charAddr*2+1]

	char := uint16(charData1) | (uint16(charData2&0x07) << 8)
	char <<= 5
	char |= uint16(y) << 2
	char |= uint16(x) >> 1

	// if d.h < 8 && d.v < 8 {
	// 	fmt.Printf("h:%03d(%02d) v:%03d(%02d) char:0x%04X \n", d.h, col, d.v, row, char)
	// }

	palette := (charData2 & 0xf0)

	pixels := d.rom[char]

	cIdx := palette | ((pixels >> ((1 - b) * 4)) & 0xf)

	// cIdx = byte(col) & 0x0f
	// cIdx = (byte(row) & 0x0f) << 4
	rgb := rgb8b(d.color[cIdx])

	d.display.Set(d.h, d.v, rgb.color())

	if d.v&32 == 32 {
		d.cpu.Interrupt(true)
	}

	d.h++
	if d.h > 456 {
		d.h = 0
		d.v++
		if d.v == 256 { // 262
			d.v = 0
			d.monitor.FrameDone()

			// f, err := os.Create(fmt.Sprintf("tetris/img_%X.png", ic))
			// if err != nil {
			// 	panic(err)
			// }
			// png.Encode(f, d.display)

			ic++
		}
	}
	*d.hBlank = !(d.v > 240)
}

type rgb8b byte

func (c rgb8b) color() color.Color {
	r := uint8(c) & 0b11100000 >> 5
	g := uint8(c) & 0b00011100 >> 2
	b := uint8(c) & 0b00000011
	return color.RGBA{r * 32, g * 32, b * 64, 0xff}
}
