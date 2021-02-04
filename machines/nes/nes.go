package nes

import (
	"fyne.io/fyne"
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/machines/nes/mappers"
	"github.com/laullon/b2t80s/ui"
)

var palClock = uint(1_662_607)

type nes struct {
	clock    emulator.Clock
	cpu      emulator.CPU
	debugger emulator.Debugger
}

func NewNES() machines.Machine {
	cartridge := mappers.CreateMapper("tests/cpu_interrupts.nes")

	clock := emulator.NewCLock(palClock, 50)
	bus := m6502.NewBus()
	cpu := m6502.MewM6502(bus)
	apu := newAPU(cpu, palClock)

	debugger := m6502.NewDebugger(cpu, nil, clock)
	// debugger.SetDump(true)
	cpu.SetDebuger(debugger)

	clock.AddTicker(0, cpu)
	clock.AddTicker(2, apu)

	// RAM
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b00000000_00000000}, &ram{data: make([]byte, 0x800), mask: 0x7ff})

	// Fake PPU
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b00100000_00000000}, &ram{data: make([]byte, 0x08), mask: 0x07})

	// APU
	bus.RegisterPort(emulator.PortMask{Mask: 0b11100000_00000000, Value: 0b01000000_00000000}, apu)

	cartridge.Insert(bus)

	m := &nes{
		cpu:      m6502.MewM6502(bus),
		clock:    emulator.NewCLock(3_584_160/3, 60),
		debugger: debugger,
	}

	return m
}

func (t *nes) Debugger() emulator.Debugger     { return t.debugger }
func (t *nes) Monitor() emulator.Monitor       { return nil } //t.tia.monitor }
func (t *nes) Clock() emulator.Clock           { return t.clock }
func (t *nes) UIControls() []ui.Control        { return nil }
func (t *nes) GetVolumeControl() func(float64) { return func(f float64) {} }
func (t *nes) OnKeyEvent(key *fyne.KeyEvent)   {}

// ----------------------------
// type clearIRQ struct {
// 	cpu emulator.CPU
// }

// func (s *clearIRQ) ReadPort(addr uint16) (byte, bool) { panic(-1) }
// func (s *clearIRQ) WritePort(addr uint16, data byte)  { s.cpu.Interrupt(false) }

type rom struct {
	data []byte
}

func (rom *rom) ReadPort(addr uint16) (byte, bool) { return rom.data[addr&0x0fff], false }
func (rom *rom) WritePort(addr uint16, data byte)  { panic(-1) }

type ram struct {
	data []byte
	mask uint16
}

func (ram *ram) ReadPort(addr uint16) (byte, bool) { return ram.data[addr&ram.mask], false }
func (ram *ram) WritePort(addr uint16, data byte)  { ram.data[addr&ram.mask] = data }
