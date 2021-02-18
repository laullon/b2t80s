package nes

import (
	"os"

	"fyne.io/fyne"
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/machines/nes/mappers"
	"github.com/laullon/b2t80s/ui"
)

var palClock = uint(1_662_607)

type nes struct {
	cpu emulator.CPU
	ppu *ppu
	apu *apu

	clock    emulator.Clock
	debugger emulator.Debugger
}

func NewNES() machines.Machine {
	// cartridge := mappers.CreateMapper("games/nes/GALAXIAN.NES")
	// cartridge := mappers.CreateMapper("/Users/glaullon/Downloads/palette_pal.nes")
	// cartridge := mappers.CreateMapper("/Users/glaullon/Downloads/allpads.nes")
	// cartridge := mappers.CreateMapper("machines/nes/tests/nestest.nes")
	// cartridge := mappers.CreateMapper("games/nes/Donkey Kong Classics (U).nes")
	cartridge := mappers.CreateMapper(os.Args[len(os.Args)-1])

	clock := emulator.NewCLock(palClock, 50)
	cpuBus := m6502.NewBus()
	cpu := m6502.MewM6502(cpuBus)

	apu := newAPU(cpu, palClock)

	ppuBus := m6502.NewBus()
	ppu := newPPU(ppuBus, cpu)

	// DMA
	apu.cpuBus = cpuBus
	apu.ppu = ppu

	clock.AddTicker(0, cpu)
	clock.AddTicker(2, apu)
	clock.AddTicker(5, ppu)

	// RAM
	cpuBus.RegisterPort("ram", emulator.PortMask{Mask: 0b1110_0000_0000_0000, Value: 0b0000_0000_0000_0000}, &m6502.BasicRam{Data: make([]byte, 0x800), Mask: 0x7ff})

	// PPU
	cpuBus.RegisterPort("ppu", emulator.PortMask{Mask: 0b1110_0000_0000_0000, Value: 0b0010_0000_0000_0000}, ppu)

	// APU
	cpuBus.RegisterPort("apu", emulator.PortMask{Mask: 0b1111_1111_1110_0000, Value: 0b0100_0000_0000_0000}, apu)

	cartridge.ConnectToCPU(cpuBus)
	cartridge.ConnectToPPU(ppuBus)

	m := &nes{
		cpu:   cpu,
		ppu:   ppu,
		apu:   apu,
		clock: clock,
	}

	if *machines.Debug {
		debugger := m6502.NewDebugger(cpu, nil, clock)
		debugger.SetDump(true)
		cpu.SetDebuger(debugger)
		m.debugger = debugger
	}

	print("cpu bus:\n", cpuBus.DumpMap(), "\n")
	print("ppu bus:\n", ppuBus.DumpMap(), "\n")

	return m
}

func (t *nes) Debugger() emulator.Debugger     { return t.debugger }
func (t *nes) Monitor() emulator.Monitor       { return t.ppu.monitor }
func (t *nes) Clock() emulator.Clock           { return t.clock }
func (t *nes) UIControls() []ui.Control        { return []ui.Control{newPalleteControl(t.ppu)} }
func (t *nes) GetVolumeControl() func(float64) { return func(f float64) {} }
func (t *nes) OnKeyEvent(key *fyne.KeyEvent)   { t.apu.onKeyEvent(key) }

// ----------------------------
