package msx

import (
	"fmt"
	"image"
	"time"

	"fyne.io/fyne"
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/ay8912"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
	"github.com/laullon/b2t80s/machines"
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

	ppi *ppi
	vdp *tms9918
}

func NewMSX() machines.Machine {
	mem := NewMemory()

	// mem.LoadRom(0, data.MustAsset("data/roms/msx/MSX.ROM"))

	mem.LoadRom(0, data.MustAsset("data/roms/msx/cbios_main_msx1_eu.rom"))
	mem.LoadRom(1, data.MustAsset("data/roms/msx/cbios_logo_msx1.rom"))

	// rom := data.MustAsset("data/roms/msx/MSX System v1.0 + MSX BASIC (1983)(Microsoft)[MSX.ROM].rom")
	// mem.LoadRom(0, rom)
	// mem.LoadRom(1, rom[0x4000:])

	if len(*machines.RomFile) > 0 {
		mem.LoadCartridge(utils.ReadFile(*machines.RomFile))
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
	case 0xa8, 0xa9, 0xaa, 0xab:
		return msx.ppi.ReadPort(port)

	case 0x98, 0x99:
		return msx.vdp.ReadPort(port)

	case 0xa2:
		return msx.ay8912.ReadRegister(msx.ayReg), false

	default:
		panic(fmt.Sprintf("[ReadPort] Unsopported port: 0x%02X", port))
	}
}

func (msx *msx) WritePort(port uint16, data byte) {
	switch port & 0xff {
	case 0xa8, 0xa9, 0xaa, 0xab:
		msx.ppi.WritePort(port, data)

	case 0x98, 0x99:
		msx.vdp.WritePort(port, data)

	case 0xa0: // TODO: move to a wrapper
		msx.ayReg = data

	case 0xa1:
		msx.ay8912.WriteRegister(msx.ayReg, data)

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

func (msx *msx) Display() image.Image {
	return msx.vdp.display
}

func (msx *msx) UIControls() []ui.Control {
	var res []ui.Control
	res = append(res, ui.NewVolumenControl(msx.sound))
	return res
}

func (msx *msx) GetVolumeControl() func(float64) {
	return msx.sound.SetVolume
}