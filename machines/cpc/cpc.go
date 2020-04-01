package cpc

import (
	"fmt"
	"image"
	"time"

	"fyne.io/fyne"
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/ay8912"
	"github.com/laullon/b2t80s/emulator/fdc"
	"github.com/laullon/b2t80s/emulator/files"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/z80"
)

const (
	CLOCK_CPC = 4000000
)

type CPC interface {
	Debugger() emulator.Debugger
	OnKeyEvent(event *fyne.KeyEvent)
	LoadZ80File(fileName string)
}

type cpc struct {
	cpu      emulator.CPU
	mem      emulator.Memory
	ga       *gatearray
	ppi      *ppi
	cassette emulator.Cassette
	sound    emulator.SoundSystem
	clock    emulator.Clock

	debugger emulator.Debugger
}

func NewCPC(cpc464 bool, cassette emulator.Cassette) machines.Machine {
	romFile := "data/roms/cpc6128.rom"
	if cpc464 {
		romFile = "data/roms/cpc464.rom"
	}

	rom, err := data.Asset(romFile)
	if err != nil {
		panic(romFile)
	}

	mem := NewCPCMemory()
	mem.LoadRom(-1, rom[:0x3fff])
	mem.LoadRom(0, rom[0x4000:])

	if !cpc464 {
		dosFile := "data/roms/amsdos.rom"
		dos, err := data.Asset(dosFile)
		if err != nil {
			panic(dosFile)
		}
		mem.LoadRom(7, dos)
	}

	cpu := z80.NewZ80(mem, cassette)

	ay8912 := ay8912.New()
	sound := emulator.NewSoundSystem(CLOCK_CPC / 80)
	sound.AddSource(ay8912)

	crtc := newCRTC(cpu)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x4000, Value: 0x0000}, crtc)

	ppi := newPPI(crtc, cassette, ay8912)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x0800, Value: 0x0000}, ppi)

	ga := newGateArray(mem.(*memory), crtc)
	cpu.RegisterPort(emulator.PortMask{Mask: 0xc000, Value: 0x4000}, ga)

	// cpu.RegisterPort(emulator.PortMask{Mask: 0xDF00, Value: 0xDF00}, mem)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x2000, Value: 0x0000}, mem)

	fdc := fdc.New765()
	cpu.RegisterPort(emulator.PortMask{Mask: 0x0400, Value: 0x0000}, fdc)
	if len(*machines.DskAFile) > 0 {
		disc := files.LoadDsk(*machines.DskAFile)
		fdc.SetDiscA(disc)
		// fmt.Printf("%v\n", disc)
	}
	// if len(*dskBFile) > 0 {
	// 	disc := files.LoadDsk(*dskBFile)
	// 	fdc.SetDiscB(disc)
	// }

	cpc := &cpc{
		cpu:      cpu,
		mem:      mem,
		ga:       ga,
		ppi:      ppi,
		cassette: cassette,
		sound:    sound,
		clock:    emulator.NewCLock(4000000),
		debugger: z80.NewDebugger(cpu, mem),
	}

	cpu.SetClock(cpc.clock)
	mem.SetClock(cpc.clock)

	cpc.clock.AddTicker(0, crtc)
	cpc.clock.AddTicker(0, cassette)
	cpc.clock.AddTicker(0, ga)
	cpc.clock.AddTicker(2, ay8912)
	cpc.clock.AddTicker(80, sound)

	// if *machines.LoadSlow {
	// 	// cpu.RegisterTrap(0x2bbb, cassette.Play)
	// 	go func() {
	// 		cassette.Play()
	// 	}()
	// } else {
	// }

	if *machines.LoadSlow {
		// if cpc464 {
		// 	cpu.RegisterTrap(0x2836, cassette.Play)
		// } else {
		// 	cpu.RegisterTrap(0x29A6, cassette.Play)
		// }
	} else {
		if cassette.Ready() {
			if cpc464 {
				cpu.RegisterTrap(0x2836, cpc.loadTapeBlockCPC464)
			} else {
				cpu.RegisterTrap(0x29A6, cpc.loadTapeBlockCPC6128)
			}
		}
	}

	cpc.cassette.Play()
	return cpc
}

func (m *cpc) loadTapeBlockCPC464() uint16 {
	return m.cpu.LoadTapeBlockCPC(0x2872)
}

func (m *cpc) loadTapeBlockCPC6128() uint16 {
	return m.cpu.LoadTapeBlockCPC(0x29E2)
}

func (m *cpc) Debugger() emulator.Debugger {
	return m.debugger
}

func (m *cpc) OnKeyEvent(event *fyne.KeyEvent) {
	m.ppi.OnKeyEvent(event)
}

func (m *cpc) Display() image.Image {
	return m.ga.displayScaled
}

func (m *cpc) GetVolumeControl() func(float64) {
	return m.sound.SetVolume
}

func (m *cpc) Run() {
	wait := time.Duration(20 * time.Millisecond)
	runStart := time.Now()
	frames := float64(0)

	ticker := time.NewTicker(wait)
	go func() {
		for range ticker.C {
			frameStart := time.Now()

			m.Debugger().NextFrame()

			err := m.cpu.RunFrame()
			if err != nil {
				panic(err)
			}
			m.ga.FrameEnded()

			frames++
			frameTime := time.Now().Sub(frameStart)
			runTime := time.Now().Sub(runStart)
			m.Debugger().SetStatus(fmt.Sprintf("frame rate:%6.2f time:%6.2fms (%v) (ear:%v) (crtc.sl:%d)", frames/runTime.Seconds(), float64(frameTime.Microseconds())/1000, wait, m.cassette.Ear(), m.ga.crtc.cycles))
		}
	}()
}

type dummyPortsManager struct{}

func (*dummyPortsManager) ReadPort(port uint16) (byte, bool) { return 0, false }
func (*dummyPortsManager) WritePort(port uint16, data byte)  {}
