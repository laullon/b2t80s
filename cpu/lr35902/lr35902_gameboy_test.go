package lr35902_test

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2/app"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines/gameboy"
	"github.com/laullon/b2t80s/ui"
	"github.com/stretchr/testify/assert"
)

func init() {
	emulator.Debug = new(bool)
	emulator.CartFile = new(string)
	ui.App = app.NewWithID("io.fyne.test")
}

func TestInstrs(t *testing.T) {
	*emulator.CartFile = string("../../machines/gameboy/test/cpu_instrs.gb")

	serial := make(chan byte, 1000)
	gb := gameboy.New(serial)

	var result strings.Builder
	go func() {
		for i := range serial {
			result.WriteByte(i)
		}
	}()

	if !assert.NotPanics(t, func() { gb.Clock().RunFor(1000) }) {
		println("result:", result.String())
	}
}
