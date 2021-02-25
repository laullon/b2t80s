package atetris

import (
	"testing"

	canvas "fyne.io/fyne/v2/canvas"
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	emulator.Debug = new(bool)

	tetris := NewATetris().(*atetris)
	tetris.sos2.monitor = &dummyMonitor{}

	if testing.Short() {
		println("skipping logs in short mode.")
	} else {
		tetris.cpu.SetDebuger(m6502.NewDebugger(tetris.cpu, nil, tetris.clock))
	}

	defer func() {
		if r := recover(); r != nil {
			assert.FailNowf(t, "Panic on '%s'", tetris.cpu.CurrentOP())
		}
	}()

	tetris.Clock().RunFor(1)
	assert.FailNow(t, "xxx")
}

type dummyMonitor struct{}

func (m *dummyMonitor) Canvas() *canvas.Image { return nil }
func (m *dummyMonitor) FrameDone()            {}
func (m *dummyMonitor) FPS() float64          { return 0 }
