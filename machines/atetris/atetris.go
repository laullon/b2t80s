package atetris

import (
	"image"

	"fyne.io/fyne"
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/pokey"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/ui"
)

type atetris struct {
	clock    emulator.Clock
	cpu      emulator.CPU
	sos2     *sos2
	debugger emulator.Debugger
	monitor  emulator.Monitor
	pokey1   *pokey.Pokey
	pokey2   *pokey.Pokey
}

func NewATetris() machines.Machine {
	image.NewRGBA(image.Rect(0, 0, 456, 262))

	bus := newBus()
	m := &atetris{
		cpu:    m6502.MewM6502(bus),
		clock:  emulator.NewCLock(14318181, 60),
		pokey1: pokey.NewPokey(),
		pokey2: pokey.NewPokey(),
		sos2: &sos2{
			vram:    make([]byte, 0x1000),
			color:   make([]byte, 0x0100),
			rom:     loadRom("136066-1101.35a"),
			display: image.NewRGBA(image.Rect(0, 0, 456, 262)),
		},
	}

	m.monitor = emulator.NewMonitor(m.sos2.display)
	m.sos2.monitor = m.monitor

	bus.RegisterPort(emulator.PortMask{Mask: 0b1111110000000000, Value: 0b0011100000000000}, &clearIRQ{cpu: m.cpu})

	bus.RegisterPort(emulator.PortMask{Mask: 0b1111000000000000, Value: 0b0001000000000000}, &ram{mem: m.sos2.vram, mask: 0x0fff})
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111110000000000, Value: 0b0010000000000000}, &ram{mem: m.sos2.color, mask: 0x00ff})

	//POKEY
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111110000110000, Value: 0b0010100000000000}, m.pokey1)
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111110000110000, Value: 0b0010100000010000}, m.pokey2)

	m.pokey1.P7 = false
	m.pokey1.P0 = false
	m.pokey1.P1 = false
	m.pokey1.P2 = false
	m.pokey1.P3 = false

	m.pokey2.P4 = false
	m.pokey2.P0 = false
	m.pokey2.P1 = false
	m.pokey2.P2 = false
	m.pokey2.P3 = false

	m.sos2.hBlank = &m.pokey1.P6

	m.sos2.cpu = m.cpu

	m.clock.AddTicker(8, m.cpu)
	m.clock.AddTicker(0, m.sos2)

	m.debugger = m6502.NewDebugger(m.cpu, nil, m.clock)
	m.cpu.SetDebuger(m.debugger)

	return m
}

func (t *atetris) OnKeyEvent(event *fyne.KeyEvent) {}
func (t *atetris) Debugger() emulator.Debugger     { return t.debugger }

func (t *atetris) Monitor() emulator.Monitor {
	return t.monitor
}

func (t *atetris) Clock() emulator.Clock           { return t.clock }
func (t *atetris) UIControls() []ui.Control        { return nil }
func (t *atetris) GetVolumeControl() func(float64) { return func(f float64) {} }

// ----------------------------
type clearIRQ struct {
	cpu emulator.CPU
}

// TODO LDA
func (s *clearIRQ) ReadPort(addr uint16) (byte, bool) { return 0x00, false }
func (s *clearIRQ) WritePort(addr uint16, data byte)  { s.cpu.Interrupt(false) }
