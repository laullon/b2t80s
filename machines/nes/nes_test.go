package nes

import (
	"testing"

	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines/nes/mappers"
	"github.com/stretchr/testify/assert"
)

func TestInterrupts(t *testing.T) {
	mmc1 := mappers.CreateMapper("tests/cpu_interrupts.nes")

	clock := emulator.NewCLock(palClock, 50)
	bus := m6502.NewBus()
	cpu := m6502.MewM6502(bus)
	apu := newAPU(cpu, palClock)

	if testing.Short() {
		println("skipping logs in short mode.")
	} else {
		debugger := m6502.NewDebugger(cpu, nil, clock)
		debugger.SetDump(true)
		cpu.SetDebuger(debugger)
	}

	testValidator := &fakeRam{
		data: make([]byte, 0x2000),
		mask: 0x1fff,
		t:    t,
	}

	clock.AddTicker(0, cpu)
	clock.AddTicker(2, apu)
	clock.AddTicker(0, testValidator)

	// RAM
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b00000000_00000000}, &ram{data: make([]byte, 0x800), mask: 0x7ff})

	// Fake PPU
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b00100000_00000000}, &ram{data: make([]byte, 0x08), mask: 0x07})

	// APU
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b01000000_00000000}, apu)

	mmc1.ConnectToCPU(bus)

	// hijack mapper ram
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b01100000_00000000}, testValidator)

	clock.RunFor(10)
}

type fakeRam struct {
	data []byte
	mask uint16
	t    *testing.T
}

func (ram *fakeRam) ReadPort(addr uint16) (byte, bool) { return ram.data[addr&ram.mask], false }
func (ram *fakeRam) WritePort(addr uint16, data byte) {
	ram.data[addr&ram.mask] = data
}

func (ram *fakeRam) Tick() {
	if ram.data[3] != 0 {
		if ram.data[0] != 0x80 {
			println(string(ram.data[4:]))
			assert.FailNowf(ram.t, "test done", "code:0x%02X", ram.data[0])
		}
	}
}
