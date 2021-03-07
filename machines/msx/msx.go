package msx

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/cpu/z80"
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/ay8912"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
	"github.com/laullon/b2t80s/machines/msx/cartridge"
	"github.com/laullon/b2t80s/ui"
	"github.com/laullon/b2t80s/utils"
)

const (
	speed = 3579545
)

type msx struct {
	bus      z80.Bus
	cpu      z80.Z80
	mem      z80.Memory
	cassette cassette.Cassette
	sound    emulator.SoundSystem
	clock    emulator.Clock

	ay8912 ay8912.AY8912
	ayReg  byte
	joy2   bool

	ppi *ppi
	vdp *tms9918
}

func NewMSX() emulator.Machine {
	rom := data.MustAsset("data/roms/msx/cbios_main_msx1_eu.rom")
	rom = append(rom, data.MustAsset("data/roms/msx/cbios_logo_msx1.rom")...)
	// rom := data.MustAsset("data/roms/msx/MSX System v1.0 + MSX BASIC (1983)(Microsoft)[MSX.ROM].rom")

	mem := NewMemory(rom)

	if len(*emulator.RomFile) > 0 {
		romType := "plain"
		romFile := *emulator.RomFile
		if strings.Contains(romFile, "::") {
			romInfo := strings.Split(romFile, "::")
			romType = romInfo[0]
			romFile = romInfo[1]
		}

		switch strings.ToLower(romType) {
		case "konami":
			mem.setCartridge1(cartridge.NewKonami(utils.ReadFile(romFile)))
		case "plain":
			mem.setCartridge1(cartridge.NewPlain(utils.ReadFile(romFile)))
		case "ascii16":
			mem.setCartridge1(cartridge.NewAscii16(utils.ReadFile(romFile)))
		default:
			panic(fmt.Sprintf("ron type '%s' not supported", romType))
		}
	}

	clock := emulator.NewCLock(speed, 50)

	bus := z80.NewBus(mem)

	z80 := z80.NewZ80(bus)
	clock.AddTicker(0, z80)

	sound := emulator.NewSoundSystem(speed / 80)

	cas := cassette.New()

	ppi := newPPI(mem, cas)
	sound.AddSource(ppi)

	vdp := newTMS9918(z80)
	clock.AddTicker(2, vdp)

	ay8912 := ay8912.New()
	sound.AddSource(ay8912)
	clock.AddTicker(2, ay8912)

	clock.AddTicker(80, sound)

	msx := &msx{
		bus:    bus,
		cpu:    z80,
		mem:    mem,
		sound:  sound,
		clock:  clock,
		ppi:    ppi,
		vdp:    vdp,
		ay8912: ay8912,
	}

	bus.RegisterPort(cpu.PortMask{Mask: 0x0000, Value: 0x0000}, msx)

	return msx
}

func (msx *msx) ReadPort(port uint16) (byte, bool) {
	if port&0xff < 40 {
		return 0, false
	}

	switch port & 0xff {
	// case 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f:
	// 	return 0, false

	case 0xa8, 0xa9, 0xaa, 0xab:
		return msx.ppi.ReadPort(port)

	// case 0xac, 0xad, 0xae, 0xaf:
	// 	return 0, false

	case 0x98, 0x99:
		return msx.vdp.ReadPort(port)

	case 0xa2:
		if msx.ayReg == 14 {
			return readJoystick(msx.joy2), false
		}
		return msx.ay8912.ReadRegister(msx.ayReg), false

	case 0xc0, 0xc1, 0xc2, 0xc3:
		return 0, false

	default:
		panic(fmt.Sprintf("[ReadPort] Unsopported port: 0x%02X", port))
	}
}

func (msx *msx) WritePort(port uint16, data byte) {
	switch port & 0xff {
	case 0xa8, 0xa9, 0xaa, 0xab:
		msx.ppi.WritePort(port, data)

	case 0x98, 0x99, 0x9A, 0x9B:
		msx.vdp.WritePort(port, data)

	case 0xa0: // TODO: move to a wrapper
		msx.ayReg = data

	case 0xa1:
		if msx.ayReg == 15 {
			msx.joy2 = data&0b01000000 != 0
		} else {
			msx.ay8912.WriteRegister(msx.ayReg, data)
		}

	case 0x2e, 0x2f:
	case 0xc0, 0xc1, 0xc2, 0xc3:
	case 0xfc, 0xfd, 0xfe, 0xff:
	case 0x90, 0x91:

	default:
		panic(fmt.Sprintf("[WritePort] Unsopported port: 0x%02X", port))
	}
}

func (msx *msx) OnKeyEvent(event *fyne.KeyEvent) {
	msx.ppi.OnKeyEvent(event)
}

func (msx *msx) Monitor() emulator.Monitor {
	return msx.vdp.monitor
}

func (msx *msx) Clock() emulator.Clock {
	return msx.clock
}

func (msx *msx) UIControls() []ui.Control {
	var res []ui.Control
	res = append(res, ui.NewVolumenControl(msx.sound.SetVolume))
	res = append(res, newSpriteControl(msx.vdp))
	return res
}

func (msx *msx) Control() map[string]ui.Control {
	return map[string]ui.Control{"CPU": ui.NewZ80UI(msx.cpu)}
}

func (msx *msx) SetDebugger(db cpu.DebuggerCallbacks) { msx.cpu.SetDebugger(db) }

func (msx *msx) GetVolumeControl() func(float64) {
	return msx.sound.SetVolume
}

func readJoystick(joy2 bool) byte {
	j, j2 := emulator.ReadJoystick()
	if joy2 {
		j = j2
	}
	res := byte(0xff)
	if j.ON {
		if j.F2 {
			res ^= 0b000100000
		}
		if j.F {
			res ^= 0b000010000
		}
		if j.R {
			res ^= 0b000001000
		}
		if j.L {
			res ^= 0b000000100
		}
		if j.D {
			res ^= 0b000000010
		}
		if j.U {
			res ^= 0b000000001
		}
	}
	return res
}
