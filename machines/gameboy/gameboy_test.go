package gameboy

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2/app"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/ui"
	"github.com/stretchr/testify/assert"
)

func init() {
	emulator.Debug = new(bool)
	emulator.CartFile = new(string)
	ui.App = app.NewWithID("io.fyne.test")
}

func TestInstrs(t *testing.T) {
	*emulator.CartFile = string("test/cpu_instrs.gb")

	serial := make(chan byte, 1000)
	gb := New(serial).(*gb)
	// gb.cpu.SetTracer(&tracer{})

	gb.cpu.Registers().PC = 0x0100

	var result strings.Builder
	go func() {
		for i := range gb.serial {
			result.WriteByte(i)
		}
	}()

	if !assert.NotPanics(t, func() { gb.clock.RunFor(1000) }) {
		println("result:", result.String())
	}
}

type tracer struct{}

func (t *tracer) AppendLastOP(op string) { println(op) }
func (t *tracer) SetNextOP(string)       {}
func (log *tracer) SetDiss(string)       {}
