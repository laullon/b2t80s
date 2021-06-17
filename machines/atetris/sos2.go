package atetris

import (
	"fmt"
	"image/color"

	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/gui"
)

var irqPerScanline = map[uint16]bool{
	16:  false,
	48:  true,
	80:  false,
	112: true,
	144: false,
	176: true,
	208: false,
	240: true,
}

type sos2 struct {
	cpu     m6502.M6502
	v, h    uint16
	vram    []byte
	color   *colorRam
	rom     []byte
	display *gui.Display
	monitor emulator.Monitor

	hBlank *bool
}

func newSOS2() *sos2 {
	fmt.Printf("%v\n", irqPerScanline)
	return &sos2{
		vram: make([]byte, 0x1000),
		color: &colorRam{
			colors: make([]color.RGBA, 0x0100),
			mem:    make([]byte, 0x0100),
		},
		rom:     loadRom("136066-1101.35a"),
		display: gui.NewDisplay(336, 240),
	}
}

func (d *sos2) Tick() {
	col := d.h >> 3
	row := d.v >> 3
	if row < 30 && col < 42 {
		y := d.v & 7

		charAddr := uint16(row) << 6
		charAddr |= uint16(col)

		charData1 := d.vram[charAddr*2]
		charData2 := d.vram[charAddr*2+1]

		palette := (charData2 & 0xf0)

		char := uint16(charData1) | (uint16(charData2&0x07) << 8)
		char <<= 5
		char |= uint16(y) << 2

		for x := uint16(0); x < 8; x += 2 {
			pixels := d.rom[char]

			cIdx0 := palette | ((pixels >> 4) & 0xf)
			cIdx1 := palette | (pixels & 0xf)

			d.display.Set(d.h+x, d.v, d.color.colors[cIdx0])
			d.display.Set(d.h+x+1, d.v, d.color.colors[cIdx1])
			char++
		}
	}

	d.h += 8
	if d.h == 456 {
		d.h = 0
		d.v++
		if irq, ok := irqPerScanline[d.v]; ok {
			d.cpu.Interrupt(irq)
		}
		if d.v == 262 {
			d.v = 0
			d.monitor.FrameDone()
		}
	}
	*d.hBlank = (d.v > 240)
}

type colorRam struct {
	colors []color.RGBA
	mem    []byte
}

func (ram *colorRam) ReadPort(addr uint16) byte { return ram.mem[addr&0xff] }

func (ram *colorRam) WritePort(addr uint16, data byte) {
	r := data & 0b11100000 >> 5
	g := data & 0b00011100 >> 2
	b := data & 0b00000011
	ram.colors[addr&0xff] = color.RGBA{r * 32, g * 32, b * 64, 0xff}
	ram.mem[addr&0xff] = data
}
