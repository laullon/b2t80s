package msx

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne"
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/ay8912"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/machines/msx/cartridge"
	"github.com/laullon/b2t80s/ui"
	"github.com/laullon/b2t80s/utils"
	"github.com/laullon/b2t80s/z80"
)

const (
	speed = 3579545
)

type msx struct {
	cpu      emulator.CPU
	mem      emulator.Memory
	cassette cassette.Cassette
	sound    emulator.SoundSystem
	debugger emulator.Debugger
	clock    emulator.Clock

	ay8912 ay8912.AY8912
	ayReg  byte
	joy2   bool

	ppi *ppi
	vdp *tms9918
}

func NewMSX() machines.Machine {
	rom := data.MustAsset("data/roms/msx/cbios_main_msx1_eu.rom")
	rom = append(rom, data.MustAsset("data/roms/msx/cbios_logo_msx1.rom")...)
	// rom := data.MustAsset("data/roms/msx/MSX System v1.0 + MSX BASIC (1983)(Microsoft)[MSX.ROM].rom")

	mem := NewMemory(rom)

	if len(*machines.RomFile) > 0 {
		romType := "plain"
		romFile := *machines.RomFile
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

	cpu := z80.NewZ80(mem)
	clock := emulator.NewCLock(speed)

	sound := emulator.NewSoundSystem(speed / 80)

	cas := cassette.New()

	ppi := newPPI(mem, cas)
	sound.AddSource(ppi)

	vdp := newTMS9918(cpu)
	clock.AddTicker(2, vdp)

	ay8912 := ay8912.New()
	sound.AddSource(ay8912)
	clock.AddTicker(2, ay8912)

	cpu.SetClock(clock)
	clock.AddTicker(80, sound)

	msx := &msx{
		cpu:      cpu,
		mem:      mem,
		sound:    sound,
		clock:    clock,
		ppi:      ppi,
		vdp:      vdp,
		ay8912:   ay8912,
		debugger: z80.NewDebugger(cpu, mem),
	}

	cpu.RegisterPort(emulator.PortMask{Mask: 0x0000, Value: 0x0000}, msx)

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

func (msx *msx) Run() {
	wait := time.Duration(20 * time.Millisecond)
	runStart := time.Now()
	frames := float64(0)

	ticker := time.NewTicker(wait)
	go func() {
		for range ticker.C {
			frameStart := time.Now()

			msx.Debugger().NextFrame()
			err := msx.cpu.RunFrame()
			if err != nil {
				panic(err)
			}

			frames++

			frameTime := time.Now().Sub(frameStart)
			runTime := time.Now().Sub(runStart)
			msx.Debugger().SetStatus(fmt.Sprintf("frame rate:%6.2f time:%6.2fms (%v)", frames/runTime.Seconds(), float64(frameTime.Microseconds())/1000, wait))
		}
	}()
}

func (msx *msx) Debugger() emulator.Debugger {
	return msx.debugger
}

func (msx *msx) OnKeyEvent(event *fyne.KeyEvent) {
	msx.ppi.OnKeyEvent(event)
}

func (msx *msx) Monitor() emulator.Monitor {
	return msx.vdp.monitor
}

func (msx *msx) UIControls() []ui.Control {
	var res []ui.Control
	res = append(res, ui.NewVolumenControl(msx.sound))
	res = append(res, newSpriteControl(msx.vdp, msx.debugger))
	return res
}

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
