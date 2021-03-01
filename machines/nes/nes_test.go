package nes

import (
	"image/png"
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/utils"
	"github.com/stretchr/testify/assert"
)

func init() {
	emulator.Debug = new(bool)
	CartFile = new(string)
	emulator.App = app.NewWithID("io.fyne.test")
}

func TestCPU(t *testing.T) {
	*CartFile = string("tests/nestest.nes")
	nes := NewNES().(*nes)

	nes.apu.onKeyEvent(&fyne.KeyEvent{Name: fyne.Key2})
	nes.Clock().RunFor(1)

	result, _, err := utils.ImgCompare("tests/nestest_ok.png", nes.ppu.display)
	assert.NoError(t, err, "Error on CPU/PPU test")
	if result != 0 {
		f, err := os.Create("tests/nestest_err.png")
		if err != nil {
			panic(err)
		}
		png.Encode(f, nes.ppu.display)
		assert.FailNow(t, "Error on CPU/PPU test")
	}
}

func TestInterrupts(t *testing.T) {
	*CartFile = string("tests/cpu_interrupts.nes")
	nes := NewNES().(*nes)

	nes.Clock().RunFor(2)

	result, _, err := utils.ImgCompare("tests/cpu_interrupts_ok.png", nes.ppu.display)
	assert.NoError(t, err, "Error on CPU/PPU test")
	if result != 0 {
		f, err := os.Create("tests/cpu_interrupts_err.png")
		if err != nil {
			panic(err)
		}
		png.Encode(f, nes.ppu.display)
		assert.FailNow(t, "Error on CPU/PPU test")
	}
}
