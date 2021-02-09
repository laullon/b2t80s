package nes

import (
	"testing"

	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines/nes/mappers"
	"github.com/stretchr/testify/assert"
)

func TestCPU(t *testing.T) {
	mmc1 := mappers.CreateMapper("tests/nestest.nes")

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

	testValidator := &testValidator{
		data: make([]byte, 0x2000),
		mask: 0x1fff,
		t:    t,
	}
	testValidator.validator = func(t *testing.T, ram []byte) {
		res := (uint16(ram[2]) << 8) | uint16(ram[3])
		if res != 0 {
			assert.FailNowf(t, "error", "error: %0x", res)
		}
	}

	clock.AddTicker(0, cpu)
	clock.AddTicker(2, apu)
	clock.AddTicker(0, testValidator)

	// RAM
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b00000000_00000000}, &m6502.BasicRam{Data: make([]byte, 0x800), Mask: 0x7ff})

	// Fake PPU
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b00100000_00000000}, &m6502.BasicRam{Data: make([]byte, 0x08), Mask: 0x07})

	// APU
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b01000000_00000000}, apu)

	mmc1.ConnectToCPU(bus)

	// hijack mapper ram
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b01100000_00000000}, testValidator)

	clock.RunFor(1)
}

type testValidator struct {
	data      []byte
	mask      uint16
	t         *testing.T
	validator func(*testing.T, []byte)
}

func (tv *testValidator) ReadPort(addr uint16) (byte, bool) { return tv.data[addr&tv.mask], false }
func (tv *testValidator) WritePort(addr uint16, data byte) {
	tv.data[addr&tv.mask] = data
}

func (tv *testValidator) Tick() {
	tv.validator(tv.t, tv.data)
}
