package nes

import (
	"image"
	"image/png"
	"os"
	"testing"
)

func TestDumpPages(t *testing.T) {
	// cartridge := mappers.CreateMapper("games/nes/Donkey Kong Classics (U).nes")

	// cpuBus := m6502.NewBus()
	// cpu := m6502.MewM6502(cpuBus)

	// ppuBus := m6502.NewBus()
	// ppu := newPPU(ppuBus, cpu)

	// cartridge.ConnectToPPU(ppuBus)

}
func TestDumpPallete(t *testing.T) {

	img := image.NewRGBA(image.Rect(0, 0, 160, 40))
	for y := 0; y < 4; y++ {
		for x := 0; x < 0x10; x++ {
			c := colors[y<<4|x]
			for dx := 0; dx < 10; dx++ {
				for dy := 0; dy < 10; dy++ {
					img.Set(x*10+dx, y*10+dy, c)
				}
			}
		}
	}

	f, err := os.Create("tests/palette.png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)

}
