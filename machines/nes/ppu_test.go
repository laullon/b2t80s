package nes

import (
	"image/png"
	"os"
	"testing"

	"fyne.io/fyne/v2/app"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/ui"
	"github.com/laullon/b2t80s/utils"
	"github.com/stretchr/testify/assert"
)

func init() {
	emulator.Debug = new(bool)
	emulator.CartFile = new(string)
	ui.App = app.NewWithID("io.fyne.test")
}

func Test_test_cpu_exec_space_ppuio(t *testing.T) {
	*emulator.CartFile = string("tests/test_cpu_exec_space_ppuio.nes")
	nes := NewNES().(*nes)

	nes.Clock().RunFor(1)

	result, _, _ := utils.ImgCompare("tests/test_cpu_exec_space_ppuio_ok.png", nes.ppu.display)
	if result != 0 {
		f, err := os.Create("tests/test_cpu_exec_space_ppuio_err.png")
		if err != nil {
			panic(err)
		}
		png.Encode(f, nes.ppu.display)
		assert.FailNow(t, "Error on CPU test")
	}
}

// }
