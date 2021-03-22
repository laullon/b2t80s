package gameboy

import (
	"fyne.io/fyne/v2"
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/cpu/lr35902"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines/gameboy/mappers"
	"github.com/laullon/b2t80s/ui"
)

type gb struct {
	cpu    lr35902.LR35902
	lcd    *lcd
	apu    *apu
	bus    cpu.Bus
	hram   cpu.RAM
	timer  *timer
	clock  emulator.Clock
	serial chan byte
}

func New(serial chan byte) emulator.Machine {
	m := &gb{
		serial: serial,
		hram:   cpu.NewRAM(make([]byte, 0x0080), 0x007f),
	}

	cartridge := mappers.CreateMapper(*emulator.CartFile)

	m.bus = cpu.NewBus(m)
	// if *emulator.Debug {
	// 	m.cpuBus = m6502.NewWatchableBus(m.cpuBus)
	// }

	m.cpu = lr35902.New(m.bus)
	m.cpu.Registers().PC = 0x100 // TODO: bios?
	m.lcd = newLCD(m.bus)
	m.apu = newAPU()
	m.timer = newTimer(m.bus)

	m.bus.RegisterPort("vram", cpu.PortMask{0b1110_0000_0000_0000, 0b1000_0000_0000_0000}, m.lcd.vRAM)
	m.bus.RegisterPort("oam", cpu.PortMask{0b1111_1111_0000_0000, 0b1111_1110_0000_0000}, m.lcd.oam)

	m.bus.RegisterPort("wram", cpu.PortMask{0b1110_0000_0000_0000, 0b1100_0000_0000_0000}, cpu.NewRAM(make([]byte, 0x2000), 0x1fff))

	m.bus.RegisterPort("APU", cpu.PortMask{0b1111_1111_1111_0000, 0b1111_1111_0001_0000}, m.apu)
	m.bus.RegisterPort("APU", cpu.PortMask{0b1111_1111_1111_0000, 0b1111_1111_0010_0000}, m.apu)
	m.bus.RegisterPort("APU", cpu.PortMask{0b1111_1111_1111_0000, 0b1111_1111_0011_0000}, m.apu)

	m.bus.RegisterPort("LCD", cpu.PortMask{0b1111_1111_1111_0000, 0b1111_1111_0100_0000}, m.lcd)

	m.bus.RegisterPort("TIMER", cpu.PortMask{0b1111_1111_1111_1100, 0b1111_1111_0000_0100}, m.timer)

	// m.bus.RegisterPort("hram", cpu.PortMask{0b1111_1111_1000_0000, 0b1111_1111_1000_0000}, m)

	cartridge.ConnectToCPU(m.bus)

	clock := emulator.NewCLock(4_190_000, 50)
	m.clock = clock
	clock.AddTicker(0, m.cpu)
	clock.AddTicker(0, m.lcd)
	clock.AddTicker(0, m.timer)

	print("cpu bus:\n", m.bus.DumpMap(), "\n")
	// print("ppu bus:\n", m.ppuBus.DumpMap(), "\n")

	// panic(-1)

	return m
}

func (gb *gb) UIControls() []ui.Control {
	return []ui.Control{
		ui.NewBusUI("memory", gb.bus),
	}
}

func (gb *gb) Control() map[string]ui.Control {
	return map[string]ui.Control{
		"CPU":   ui.NewLR35902UI(gb.cpu),
		"LCD":   newLcdControl(gb.lcd),
		"TIMER": newTimerControl(gb.timer),
	}
}

func (gb *gb) Monitor() emulator.Monitor       { return gb.lcd.monitor }
func (gb *gb) Clock() emulator.Clock           { return gb.clock }
func (gb *gb) GetVolumeControl() func(float64) { return func(f float64) {} }
func (gb *gb) OnKeyEvent(key *fyne.KeyEvent)   {}

func (gb *gb) SetDebugger(db cpu.DebuggerCallbacks) {
	gb.cpu.SetDebugger(db)
}

func (gb *gb) ReadPort(addr uint16) (byte, bool) {
	switch addr {
	case 0xffff:
		return gb.cpu.Registers().IE, false
	case 0xff0f:
		return gb.cpu.Registers().IF, false
	default:
		if addr > 0xff7f {
			return gb.hram.ReadPort(addr)
		}
		// panic(-1)
	}
	return 0xff, false
}

func (gb *gb) WritePort(addr uint16, data byte) {
	switch addr {
	case 0xff01:
		if gb.serial != nil {
			gb.serial <- data
		}
	case 0xff02:

	case 0xffff:
		gb.cpu.Registers().IE = data
	case 0xff0f:
		gb.cpu.Registers().IF |= data
	default:
		if addr > 0xff7f {
			// fmt.Printf("[GB][writePort]-> port:0x%04X data:0x%02X  \n", addr, data)
			gb.hram.WritePort(addr, data)
			// } else {
			// panic(fmt.Sprintf("Panic on [GB][writePort]-> port:0x%04X data:0x%02X  \n", addr, data))
		}
	}
}
