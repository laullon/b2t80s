package a2600

import (
	"image"

	"fyne.io/fyne"
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/ui"
	"github.com/laullon/b2t80s/utils"
)

type a2600 struct {
	clock    emulator.Clock
	cpu      emulator.CPU
	debugger emulator.Debugger
	tia      *tia
}

func NewA2600() machines.Machine {
	bus := m6502.NewBus()
	display := image.NewRGBA(image.Rect(0, 0, 160, 192))
	monitor := emulator.NewMonitor(display)

	m := &a2600{
		cpu:   m6502.MewM6502(bus),
		clock: emulator.NewCLock(3_584_160/3, 60),
		tia:   &tia{data: make([]byte, 0x80), monitor: monitor, display: display},
	}

	// TIA
	m.tia.cpu = m.cpu
	bus.RegisterPort(emulator.PortMask{Mask: 0b1_1111_1000_0000, Value: 0b0_0000_0000_0000}, m.tia)

	// RAM
	bus.RegisterPort(emulator.PortMask{Mask: 0b1_1111_1000_0000, Value: 0b0_0000_1000_0000}, &ram{data: make([]byte, 0x80)})

	// SP
	bus.RegisterPort(emulator.PortMask{Mask: 0b1_1111_0000_0000, Value: 0b0_0001_0000_0000}, &ram{data: make([]byte, 0x100)})

	// RIOT
	bus.RegisterPort(emulator.PortMask{Mask: 0b1_1111_0000_0000, Value: 0b0_0010_0000_0000}, &ram{data: make([]byte, 0x100)})

	// ROM
	bin := utils.ReadFile("games/a2600/Space_Invaders.bin")
	// bin := utils.ReadFile("games/a2600/kernel_13.bin")
	// bin := utils.ReadFile("games/a2600/kernel_15.bin")
	bus.RegisterPort(emulator.PortMask{Mask: 0b1000000000000, Value: 0b1000000000000}, &rom{data: bin})

	if *machines.Debug {
		m.debugger = m6502.NewDebugger(m.cpu, nil, m.clock)
		m.cpu.SetDebuger(m.debugger)
	}

	m.clock.AddTicker(0, m.cpu)
	m.clock.AddTicker(0, m.tia)

	return m
}

func (t *a2600) Debugger() emulator.Debugger     { return t.debugger }
func (t *a2600) Monitor() emulator.Monitor       { return t.tia.monitor }
func (t *a2600) Clock() emulator.Clock           { return t.clock }
func (t *a2600) UIControls() []ui.Control        { return nil }
func (t *a2600) GetVolumeControl() func(float64) { return func(f float64) {} }
func (t *a2600) OnKeyEvent(key *fyne.KeyEvent)   {}

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
}

func (ram *ram) ReadPort(addr uint16) (byte, bool) { return ram.data[addr&0x7f], false }
func (ram *ram) WritePort(addr uint16, data byte)  { ram.data[addr&0x7f] = data }
