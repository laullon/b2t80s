package nes

import (
	"fyne.io/fyne/v2"
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines/nes/mappers"
	"github.com/laullon/b2t80s/ui"
)

var palClock = uint(1_662_607)
var CartFile *string

type nes struct {
	cpu m6502.M6502
	ppu *ppu
	apu *apu

	clock    emulator.Clock
	debugger emulator.Debugger
}

func NewNES() emulator.Machine {
	m := &nes{}

	cartridge := mappers.CreateMapper(*CartFile)

	clock := emulator.NewCLock(palClock, 50)
	cpuBus := m6502.NewBus()
	if *emulator.Debug {
		cpuBus = m6502.NewWatchableBus(cpuBus)
	}

	m6805 := m6502.MewM6502(cpuBus)

	apu := newAPU(m6805, palClock)

	ppuBus := m6502.NewBus()
	ppu := newPPU(ppuBus, m6805)

	// DMA
	apu.cpuBus = cpuBus
	apu.ppu = ppu

	clock.AddTicker(0, m6805)
	clock.AddTicker(2, apu)
	clock.AddTicker(5, ppu)

	// RAM
	cpuBus.RegisterPort("ram", cpu.PortMask{Mask: 0b1110_0000_0000_0000, Value: 0b0000_0000_0000_0000}, &m6502.BasicRam{Data: make([]byte, 0x800), Mask: 0x7ff})

	// PPU
	cpuBus.RegisterPort("ppu", cpu.PortMask{Mask: 0b1110_0000_0000_0000, Value: 0b0010_0000_0000_0000}, ppu)

	// APU
	cpuBus.RegisterPort("apu", cpu.PortMask{Mask: 0b1111_1111_1110_0000, Value: 0b0100_0000_0000_0000}, apu)

	cartridge.ConnectToCPU(cpuBus)
	cartridge.ConnectToPPU(ppuBus)

	m.cpu = m6805
	m.ppu = ppu
	m.apu = apu
	m.clock = clock

	print("cpu bus:\n", cpuBus.DumpMap(), "\n")
	print("ppu bus:\n", ppuBus.DumpMap(), "\n")

	return m
}

func (t *nes) Debugger() emulator.Debugger          { return t.debugger }
func (t *nes) Monitor() emulator.Monitor            { return t.ppu.monitor }
func (t *nes) Clock() emulator.Clock                { return t.clock }
func (t *nes) UIControls() []ui.Control             { return []ui.Control{newPalleteControl(t.ppu)} }
func (t *nes) GetVolumeControl() func(float64)      { return func(f float64) {} }
func (t *nes) OnKeyEvent(key *fyne.KeyEvent)        { t.apu.onKeyEvent(key) }
func (t *nes) CPUControl() ui.Control               { return ui.NewM6502UI(t.cpu) }
func (t *nes) SetDebugger(db cpu.DebuggerCallbacks) { t.cpu.SetDebugger(db) }

// ----------------------------
