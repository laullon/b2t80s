package atetris

import (
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/pokey"
	"github.com/laullon/b2t80s/ui"
)

type atetris struct {
	clock    emulator.Clock
	cpu      m6502.M6502
	sos2     *sos2
	debugger emulator.Debugger
	monitor  emulator.Monitor
	pokey1   *pokey.Pokey
	pokey2   *pokey.Pokey
}

func NewATetris() emulator.Machine {
	bus := newBus()

	m := &atetris{
		cpu:    m6502.MewM6502(bus),
		clock:  emulator.NewCLock(14318181/8, 60),
		pokey1: pokey.NewPokey(),
		pokey2: pokey.NewPokey(),
		sos2:   newSOS2(),
	}

	m.monitor = emulator.NewMonitor(m.sos2.display)
	m.sos2.monitor = m.monitor

	bus.RegisterPort("clearIRQ", cpu.PortMask{Mask: 0b1111110000000000, Value: 0b0011100000000000}, &clearIRQ{cpu: m.cpu})

	bus.RegisterPort("vRam", cpu.PortMask{Mask: 0b1111000000000000, Value: 0b0001000000000000}, &ram{mem: m.sos2.vram, mask: 0x0fff})
	bus.RegisterPort("color", cpu.PortMask{Mask: 0b1111110000000000, Value: 0b0010000000000000}, m.sos2.color)

	//POKEY
	bus.RegisterPort("pokey1", cpu.PortMask{Mask: 0b1111110000110000, Value: 0b0010100000000000}, m.pokey1)
	bus.RegisterPort("pokey2", cpu.PortMask{Mask: 0b1111110000110000, Value: 0b0010100000010000}, m.pokey2)

	// Watchdog
	wd := &watchdog{cpu: m.cpu}
	wd.start()
	bus.RegisterPort("watchdog", cpu.PortMask{Mask: 0b1111110000000000, Value: 0b0011000000000000}, wd)

	m.pokey1.P7 = false

	m.sos2.hBlank = &m.pokey1.P6

	m.sos2.cpu = m.cpu

	m.clock.AddTicker(0, m.cpu)
	m.clock.AddTicker(2, m.sos2)

	return m
}

func (t *atetris) Debugger() emulator.Debugger { return t.debugger }

func (t *atetris) Monitor() emulator.Monitor {
	return t.monitor
}

func (t *atetris) Clock() emulator.Clock                { return t.clock }
func (t *atetris) UIControls() []ui.Control             { return nil }
func (t *atetris) GetVolumeControl() func(float64)      { return func(f float64) {} }
func (t *atetris) CPUControl() ui.Control               { return ui.NewM6502UI(t.cpu) }
func (t *atetris) SetDebugger(db cpu.DebuggerCallbacks) { t.cpu.SetDebugger(db) }

// ----------------------------
type clearIRQ struct {
	cpu cpu.CPU
}

func (s *clearIRQ) ReadPort(addr uint16) (byte, bool) { panic(-1) }
func (s *clearIRQ) WritePort(addr uint16, data byte)  { s.cpu.Interrupt(false) }
