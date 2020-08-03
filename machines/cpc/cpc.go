package cpc

import (
	"fmt"
	"time"

	"fyne.io/fyne"
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/ay8912"
	"github.com/laullon/b2t80s/emulator/files"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
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
	cassette cassette.Cassette
	sound    emulator.SoundSystem
	clock    emulator.Clock

	debugger emulator.Debugger
}

func NewCPC(cpc464 bool) machines.Machine {
	cassette.SpeedAdj = float64(40) / float64(35)

	romFile := "data/roms/cpc6128.rom"
	if cpc464 {
		romFile = "data/roms/cpc464.rom"
	}

	mem := NewCPCMemory()
	rom := data.MustAsset(romFile)
	mem.LoadRom(-1, rom[:0x3fff])
	mem.LoadRom(0, rom[0x4000:])

	if !cpc464 {
		dos := data.MustAsset("data/roms/amsdos.rom")
		mem.LoadRom(7, dos)
	}

	cas := cassette.New()
	if len(*machines.TapFile) > 0 {
		cas.LoadTapFile(*machines.TapFile)
	}

	cpu := z80.NewZ80(mem)

	ay8912 := ay8912.New()
	sound := emulator.NewSoundSystem(CLOCK_CPC / 80)
	sound.AddSource(ay8912)

	crtc := newCRTC(cpu)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x4000, Value: 0x0000}, crtc)

	ppi := newPPI(crtc, cas, ay8912)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x0800, Value: 0x0000}, ppi)
	sound.AddSource(ppi)

	ga := newGateArray(mem, crtc)
	cpu.RegisterPort(emulator.PortMask{Mask: 0xc000, Value: 0x4000}, ga)

	// cpu.RegisterPort(emulator.PortMask{Mask: 0xDF00, Value: 0xDF00}, mem)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x2000, Value: 0x0000}, mem)

	fdc := NewCPCFDC765()
	cpu.RegisterPort(emulator.PortMask{Mask: 0x0400, Value: 0x0000}, fdc)
	if len(*machines.DskAFile) > 0 {
		disc := files.LoadDsk(*machines.DskAFile)
		fdc.chip.SetDiscA(disc)
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
		cassette: cas,
		sound:    sound,
		clock:    emulator.NewCLock(4000000),
		debugger: z80.NewDebugger(cpu, mem),
	}

	cpu.SetClock(cpc.clock)
	mem.SetClock(cpc.clock)

	cpc.clock.AddTicker(4, crtc)
	cpc.clock.AddTicker(4, ga)
	cpc.clock.AddTicker(4, ay8912)
	cpc.clock.AddTicker(80, sound)
	if *machines.LoadSlow {
		cpc.clock.AddTicker(0, cas)
	}

	if !*machines.LoadSlow {
		if cpc464 {
			cpu.RegisterTrap(0x2836, cpc.loadTapeBlockCPC464)
		} else {
			cpu.RegisterTrap(0x29A6, cpc.loadTapeBlockCPC6128)
		}
	}

	return cpc
}

func (m *cpc) loadTapeBlockCPC464() uint16 {
	return m.LoadTapeBlockCPC(0x2872)
}

func (m *cpc) loadTapeBlockCPC6128() uint16 {
	return m.LoadTapeBlockCPC(0x29E2)
}

func (m *cpc) LoadTapeBlockCPC(exit uint16) uint16 {
	data := m.cassette.NextDataBlock()
	if data == nil {
		return emulator.CONTINUE
	}

	regs := m.cpu.Registers().(*z80.Z80Registers)
	requestedLength := regs.DE.Get()
	startAddress := regs.HL.Get()
	t := regs.A
	// fmt.Printf("Loading block to 0x%04x (bl:0x%04x, l:0x%04x, bt:0x%02X, t:0x%02X)\n", startAddress, len(data), requestedLength, data[0], t)
	if t == data[0] {
		for i := uint16(0); i < requestedLength; i++ {
			m.mem.PutByte(startAddress+i, data[i+1])
		}
		regs.F.SetByte(0x45)
		// println("Done")
		// println(hex.Dump(data[:requestedLength]))
		// println(hex.Dump(m.mem.GetBlock(startAddress, requestedLength)))
		// } else {
		// 	println("BAD Block")
	}
	return exit
}

func (m *cpc) Debugger() emulator.Debugger {
	return m.debugger
}

func (m *cpc) OnKeyEvent(event *fyne.KeyEvent) {
	m.ppi.OnKeyEvent(event)
}

func (m *cpc) Monitor() emulator.Monitor {
	return m.ga.monitor
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
			m.Debugger().SetStatus(fmt.Sprintf("frame rate:%6.2f time:%6.2fms (%v) (ear:%v)", frames/runTime.Seconds(), float64(frameTime.Microseconds())/1000, wait, m.cassette.Ear()))
		}
	}()
}

type dummyPortsManager struct{}

func (*dummyPortsManager) ReadPort(port uint16) (byte, bool) { return 0, false }
func (*dummyPortsManager) WritePort(port uint16, data byte)  {}
