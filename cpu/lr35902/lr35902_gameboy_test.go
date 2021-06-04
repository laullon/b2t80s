package lr35902_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines/gameboy"
	"github.com/stretchr/testify/assert"
)

func init() {
	emulator.Debug = new(bool)
	emulator.CartFile = new(string)
}

func TestInstrs(t *testing.T) {
	*emulator.CartFile = string("/Users/glaullon/Downloads/gb-test-roms-master/cpu_instrs/individual/11-op a,(hl).gb")
	*gameboy.Bios = string("../../bios/gb_bios.bin")

	serial := make(chan byte, 1000)
	gb := gameboy.New(serial)
	gb.Reset()

	var result strings.Builder
	go func() {
		for i := range serial {
			result.WriteByte(i)
		}
	}()

	assert.NotPanics(t, func() { gb.Clock().RunFor(20) })
	println("result:", result.String())

	re := regexp.MustCompile(`(\d.):(\w.)`)
	results := re.FindAllSubmatch([]byte(result.String()), -1)
	assert.NotEqual(t, 0, len(results), "tests no executed")
	for _, res := range results {
		assert.Equalf(t, "ok", string(res[2]), "error on test %s error %s", res[1], res[2])
	}
}

func _TestInstrsTiming(t *testing.T) {
	*emulator.CartFile = string("/Users/glaullon/Downloads/instr_timing.gb")

	serial := make(chan byte, 1000)
	gb := gameboy.New(serial)

	var result strings.Builder
	go func() {
		for i := range serial {
			result.WriteByte(i)
		}
	}()

	if !assert.NotPanics(t, func() { gb.Clock().RunFor(100) }) {
		println("result:", result.String())
	}

	println("result:", result.String())
}
