package gameboy

import (
	"flag"
	"fmt"
	"os"

	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/cpu/lr35902"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/gui"
	"github.com/laullon/b2t80s/machines/gameboy/mappers"
	"github.com/laullon/b2t80s/ui"
	"github.com/laullon/b2t80s/utils"
	"github.com/veandco/go-sdl2/sdl"
)

type gb struct {
	cpu          lr35902.LR35902
	ppu          *ppu
	apu          *apu
	bus          cpu.Bus
	hram         cpu.RAM
	bios         *bios
	timer        *timer
	clock        emulator.Clock
	serial       chan byte
	serialBuffer []byte

	controls     *byte
	pad, buttons byte
}

var Bios = flag.String("bios", "bios/gb_bios.bin", "GameBoy bios")

const clockHz = 4_194_304
const soundDiv = 128

func New(serial ...chan byte) emulator.Machine {
	m := &gb{
		hram:    cpu.NewRAM(make([]byte, 0x0080), 0x007f),
		pad:     0xff,
		buttons: 0xff,
	}

	if len(serial) > 0 {
		m.serial = serial[0]
	} else {
		m.serial = make(chan byte, 1000)
		go func() {
			for i := range m.serial {
				m.serialBuffer = append(m.serialBuffer, i)
				if len(m.serialBuffer) > 100 {
					m.serialBuffer = m.serialBuffer[len(m.serialBuffer)-100:]
				}
			}
		}()
	}

	cartridge := mappers.CreateMapper(*emulator.CartFile)

	m.bus = cpu.NewBus("main", m)
	// if *emulator.Debug {
	// 	m.cpuBus = m6502.NewWatchableBus(m.cpuBus)
	// }

	m.cpu = lr35902.New(m.bus)
	m.ppu = newPPU(m.bus)
	m.apu = newAPU(clockHz / soundDiv)
	m.timer = newTimer(m.bus)

	// BIOS
	if _, err := os.Stat(*Bios); err == nil {
		m.bios = Newbios(utils.ReadFile(*Bios))
		m.bus.RegisterPort("bios/rom", cpu.PortMask{0b1111_1111_0000_0000, 0b0000_0000_0000_0000}, m.bios)
	} else {
		fmt.Printf("Bios not found\n")
	}

	m.bus.RegisterPort("wram", cpu.PortMask{0b1110_0000_0000_0000, 0b1100_0000_0000_0000}, cpu.NewRAM(make([]byte, 0x2000), 0x1fff))

	m.bus.RegisterPort("APU", cpu.PortMask{0b1111_1111_1111_0000, 0b1111_1111_0001_0000}, m.apu)
	m.bus.RegisterPort("APU", cpu.PortMask{0b1111_1111_1111_0000, 0b1111_1111_0010_0000}, m.apu)
	m.bus.RegisterPort("APU", cpu.PortMask{0b1111_1111_1111_0000, 0b1111_1111_0011_0000}, m.apu)

	m.bus.RegisterPort("PPU", cpu.PortMask{0b1111_1111_1111_0000, 0b1111_1111_0100_0000}, m.ppu)

	m.bus.RegisterPort("TIMER", cpu.PortMask{0b1111_1111_1111_1100, 0b1111_1111_0000_0100}, m.timer)

	// m.bus.RegisterPort("hram", cpu.PortMask{0b1111_1111_1000_0000, 0b1111_1111_1000_0000}, m)

	cartridge.ConnectToCPU(m.bus)

	sound := emulator.NewSoundSystem(clockHz / soundDiv)
	sound.AddSource(m.apu)

	clock := emulator.NewCLock(clockHz, 64)
	m.clock = clock
	clock.AddTicker(0, m.ppu)
	clock.AddTicker(0, m.timer)
	clock.AddTicker(4, m.cpu)
	clock.AddTicker(8192, m.apu) // 512Hz
	clock.AddTicker(soundDiv, sound)

	print("cpu bus:\n", m.bus.DumpMap(), "\n")
	// print("ppu bus:\n", m.ppuBus.DumpMap(), "\n")

	return m
}

func (gb *gb) Reset() {
	gb.cpu.Reset()
	if gb.bios == nil {
		gb.cpu.Registers().PC = 0x100
	}
}

func (gb *gb) UIControls() []gui.GUIObject {
	return nil //} // []gui.GUIObject{
	// 	ui.NewBusUI("memory", gb.bus),
	// }
}

func (gb *gb) Control() map[string]gui.GUIObject {
	return map[string]gui.GUIObject{
		"CPU": ui.NewLR35902UI(gb.cpu),
		"PPU": newPPUControl(gb.ppu),
		// 	"SERIAL": newSerialControl(&m.serialBuffer),
		// 	"Sound":  newSoundCtrl(m.apu),
	}
}

func (gb *gb) Monitor() emulator.Monitor       { return gb.ppu.monitor }
func (gb *gb) Clock() emulator.Clock           { return gb.clock }
func (gb *gb) GetVolumeControl() func(float64) { return func(f float64) {} }

func (gb *gb) SetDebugger(db cpu.DebuggerCallbacks) {
	gb.cpu.SetDebugger(db)
	gb.ppu.debugger = db
}

func (gb *gb) ReadPort(addr uint16) (byte, bool) {
	switch addr {
	case 0xff00:
		return *gb.controls, false

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
	case 0xff00:
		if gui.Joystick1.ON {
			gb.buttons = 0xff
			if gui.Joystick1.F { // A
				gb.buttons ^= 0b00000001
			}
			if gui.Joystick1.F2 { // B
				gb.buttons ^= 0b00000010
			}
			if gui.Joystick1.Select { //select
				gb.buttons ^= 0b00000100
			}
			if gui.Joystick1.Start { // start
				gb.buttons ^= 0b00001000
			}

			gb.pad = 0xff
			if gui.Joystick1.R {
				gb.pad ^= 0b00000001
			}
			if gui.Joystick1.L {
				gb.pad ^= 0b00000010
			}
			if gui.Joystick1.U {
				gb.pad ^= 0b00000100
			}
			if gui.Joystick1.D {
				gb.pad ^= 0b00001000
			}
		}
		if data&0b0001_0000 == 0 {
			gb.controls = &gb.pad
		} else if data&0b0010_0000 == 0 {
			gb.controls = &gb.buttons
		}

	case 0xff01:
		if gb.serial != nil {
			gb.serial <- data
		}

	case 0xff02:

	case 0xff50:
		if gb.bios != nil {
			gb.bios.enable = false
		}
	case 0xffff:
		gb.cpu.Registers().IE = data

	case 0xff0f:
		gb.cpu.Registers().IF |= data

	default:
		if addr > 0xff7f {
			// fmt.Printf("[GB][writePort]-> port:0x%04X data:0x%02X  \n", addr, data)
			gb.hram.WritePort(addr, data)
			// } else {
			// 	panic(fmt.Sprintf("Panic on [GB][writePort]-> port:0x%04X data:0x%02X  \n", addr, data))
		}
	}
}

func (gb *gb) OnKey(key sdl.Scancode) {
	switch key {
	case sdl.SCANCODE_Z: // A
		gb.buttons ^= 0b00000001
	case sdl.SCANCODE_X: // B
		gb.buttons ^= 0b00000010
	case sdl.SCANCODE_1: //select
		gb.buttons ^= 0b00000100
	case sdl.SCANCODE_2: // start
		gb.buttons ^= 0b00001000

	case sdl.SCANCODE_RIGHT:
		gb.pad ^= 0b00000001
	case sdl.SCANCODE_LEFT:
		gb.pad ^= 0b00000010
	case sdl.SCANCODE_UP:
		gb.pad ^= 0b00000100
	case sdl.SCANCODE_DOWN:
		gb.pad ^= 0b00001000
	}
}

// ************************
// ************************
// ************************

type bios struct {
	bank   []byte
	enable bool
}

func Newbios(bank []byte) *bios {
	return &bios{bank: bank, enable: true}
}

func (bios *bios) SetBank(bank []byte) {
	bios.bank = bank
}

func (bios *bios) ReadPort(addr uint16) (byte, bool) {
	return bios.bank[addr], !bios.enable
}

func (bios *bios) WritePort(addr uint16, data byte) {
}
