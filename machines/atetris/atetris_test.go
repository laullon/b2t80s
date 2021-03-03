package atetris

import (
	"image/png"
	"os"
	"testing"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/ui"
	"github.com/laullon/b2t80s/utils"
	"github.com/stretchr/testify/assert"
)

func init() {
	emulator.Debug = new(bool)
	ui.App = app.NewWithID("io.fyne.test")
}

func TestTestMode(t *testing.T) {

	tetris := NewATetris().(*atetris)
	tetris.pokey1.P7 = true

	tetris.Clock().RunFor(5)

	result, _, err := utils.ImgCompare("tests/testMode_ok.png", tetris.sos2.display)
	assert.NoError(t, err, "Error on CPU/PPU test")
	if result != 0 {
		f, err := os.Create("tests/testMode_err.png")
		if err != nil {
			panic(err)
		}
		png.Encode(f, tetris.sos2.display)
		assert.FailNow(t, "Error on test mode init")
	}
}

type dummyMonitor struct{}

func (m *dummyMonitor) Canvas() *canvas.Image { return nil }
func (m *dummyMonitor) FrameDone()            {}
func (m *dummyMonitor) FPS() float64          { return 0 }
