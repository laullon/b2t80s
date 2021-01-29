package mappers

import (
	"testing"

	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
)

func TestMMC1(t *testing.T) {
	mmc1 := CreateMapper("../tests/cpu_interrupts.nes")

	bus := m6502.NewBus()
	cpu := m6502.MewM6502(bus)

	// RAM
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b00000000_00000000}, &ram{mem: make([]byte, 0x800), mask: 0x7ff})

	// Fake PPU
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b00100000_00000000}, &ram{mem: make([]byte, 0x08), mask: 0x07})

	mmc1.Insert(bus)

	for i := 0; i < 5000000; i++ {
		cpu.Tick()
	}
}
