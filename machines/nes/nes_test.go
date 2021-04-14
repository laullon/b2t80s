package nes

import (
	"fmt"
	"image/png"
	"os"
	"strings"
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

func TestCPU(t *testing.T) {
	*emulator.CartFile = string("tests/nestest.nes")
	nes := NewNES().(*nes)

	nes.apu.onKeyEvent(&glfw.KeyEvent{Name: glfw.Key2})
	nes.Clock().RunFor(4)

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
	*emulator.CartFile = string("tests/cpu_interrupts.nes")
	nes := NewNES().(*nes)

	nes.Clock().RunFor(4)

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

func TestCPU2(t *testing.T) {
	*emulator.CartFile = string("tests/nestest.nes")
	nes := NewNES().(*nes)

	f, err := os.Create("tests/nestest.out")
	if err != nil {
		assert.FailNowf(t, "Error on CPU 2 test", "%v", err)
	}

	tracer := &tracer{nes, f, 0, ""}
	nes.cpu.SetTracer(tracer)
	nes.cpuBus.Write(0xfffc, 0x00)
	nes.cpuBus.Write(0xfffd, 0xc0)

	nes.Clock().AddTicker(0, tracer)

	assert.Panics(t, func() {
		nes.Clock().RunFor(1)
	})

	// assert.FailNow(t, "Error on CPU 2 test")
}

type tracer struct {
	nes            *nes
	f              *os.File
	ticks          int
	regsStatusNext string
}

func (t *tracer) AppendLastOP(op string) {
	if len(t.regsStatusNext) > 0 {
		t.f.WriteString(strings.ToUpper(op))
		t.f.WriteString("                                                 "[len(op):])
		t.f.WriteString(t.regsStatusNext)
		t.f.WriteString("\n")
	}

	t.regsStatusNext = fmt.Sprintf("A:%02X X:%02X Y:%02X SP:%02X CYC:%d",
		t.nes.cpu.Registers().A, t.nes.cpu.Registers().X, t.nes.cpu.Registers().Y,
		t.nes.cpu.Registers().SP, t.ticks,
	)
}

func (t *tracer) SetNextOP(string)                                            {}
func (log *tracer) SetDiss(pc uint16, getMemory func(pc, leng uint16) []byte) {}

func (t *tracer) Tick() {
	t.ticks++
}
