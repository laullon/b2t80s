package nes

import (
	"fyne.io/fyne/v2"
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines/nes/mappers"
	"github.com/laullon/b2t80s/ui"
)

var ntscClock = uint(1_789_773)
var palClock = uint(1_662_607)
var CartFile *string

type nes struct {
	cpu    m6502.M6502
	ppu    *ppu
	apu    *apu
	cpuBus m6502.Bus
	ppuBus m6502.Bus

	clock    emulator.Clock
	debugger emulator.Debugger
}

func NewNES() emulator.Machine {
	m := &nes{}

	cartridge, ntsc := mappers.CreateMapper(*CartFile)

	m.cpuBus = m6502.NewBus()
	if *emulator.Debug {
		m.cpuBus = m6502.NewWatchableBus(m.cpuBus)
	}

	m6805 := m6502.MewM6502(m.cpuBus)

	apu := newAPU(m6805)

	m.ppuBus = m6502.NewBus()
	ppu := newPPU(m.ppuBus, m6805)

	// DMA
	apu.cpuBus = m.cpuBus
	apu.ppu = ppu

	if ntsc {
		clock := emulator.NewCLock(ntscClock, 60)
		clock.AddTicker(0, m6805)
		clock.AddTicker(2, apu)
		clock.AddTicker(0, ppu)
		m.clock = clock
		ppu.pixelsPerTicks = 3
		ppu.scanLineW = 341
		ppu.scanLineH = 261
	} else {
		clock := emulator.NewCLock(palClock, 50)
		clock.AddTicker(0, m6805)
		clock.AddTicker(2, apu)
		clock.AddTicker(5, ppu)
		m.clock = clock
		ppu.pixelsPerTicks = 16
		ppu.scanLineW = 341
		ppu.scanLineH = 312
	}

	// RAM
	m.cpuBus.RegisterPort("ram", cpu.PortMask{Mask: 0b1110_0000_0000_0000, Value: 0b0000_0000_0000_0000}, &m6502.BasicRam{Data: make([]byte, 0x800), Mask: 0x7ff})

	// PPU
	m.cpuBus.RegisterPort("ppu", cpu.PortMask{Mask: 0b1110_0000_0000_0000, Value: 0b0010_0000_0000_0000}, ppu)

	// APU
	m.cpuBus.RegisterPort("apu", cpu.PortMask{Mask: 0b1111_1111_1110_0000, Value: 0b0100_0000_0000_0000}, apu)

	cartridge.ConnectToCPU(m.cpuBus)
	cartridge.ConnectToPPU(m.ppuBus)

	m.cpu = m6805
	m.ppu = ppu
	m.apu = apu

	// print("cpu bus:\n", m.cpuBus.DumpMap(), "\n")
	// print("ppu bus:\n", m.ppuBus.DumpMap(), "\n")

	return m
}

func (t *nes) UIControls() []ui.Control {
	return []ui.Control{newPalleteControl(t.ppu), ui.NewM6502BusUI(t.cpuBus)}
}

func (t *nes) Debugger() emulator.Debugger          { return t.debugger }
func (t *nes) Monitor() emulator.Monitor            { return t.ppu.monitor }
func (t *nes) Clock() emulator.Clock                { return t.clock }
func (t *nes) GetVolumeControl() func(float64)      { return func(f float64) {} }
func (t *nes) OnKeyEvent(key *fyne.KeyEvent)        { t.apu.onKeyEvent(key) }
func (t *nes) CPUControl() ui.Control               { return ui.NewM6502UI(t.cpu) }
func (t *nes) SetDebugger(db cpu.DebuggerCallbacks) { t.cpu.SetDebugger(db) }

// ----------------------------
