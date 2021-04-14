package cpc

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/cpu/z80"
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/ay8912"
	"github.com/laullon/b2t80s/emulator/files"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
	"github.com/laullon/b2t80s/ui"
)

const (
	CLOCK_CPC = 4000000
)

type CPC interface {
	Debugger() emulator.Debugger
	OnKey(key glfw.Key)
	LoadZ80File(fileName string)
}

type cpc struct {
	bus      z80.Bus
	cpu      z80.Z80
	mem      z80.Memory
	ga       *gatearray
	ppi      *ppi
	cassette cassette.Cassette
	sound    emulator.SoundSystem
	clock    emulator.Clock

	debugger emulator.Debugger
}

func NewCPC(cpc464 bool) emulator.Machine {
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
	if len(*emulator.TapFile) > 0 {
		cas.LoadTapFile(*emulator.TapFile)
	}

	bus := z80.NewBus(mem)

	z80 := z80.NewZ80(bus)

	ay8912 := ay8912.New()
	sound := emulator.NewSoundSystem(CLOCK_CPC / 80)
	sound.AddSource(ay8912)

	crtc := newCRTC(z80)
	bus.RegisterPort(cpu.PortMask{Mask: 0x4000, Value: 0x0000}, crtc)

	ppi := newPPI(crtc, cas, ay8912)
	bus.RegisterPort(cpu.PortMask{Mask: 0x0800, Value: 0x0000}, ppi)
	sound.AddSource(ppi)

	ga := newGateArray(mem, crtc)
	bus.RegisterPort(cpu.PortMask{Mask: 0xc000, Value: 0x4000}, ga)

	// bus.RegisterPort(cpu.PortMask{Mask: 0xDF00, Value: 0xDF00}, mem)
	bus.RegisterPort(cpu.PortMask{Mask: 0x2000, Value: 0x0000}, mem)

	fdc := NewCPCFDC765()
	bus.RegisterPort(cpu.PortMask{Mask: 0x0400, Value: 0x0000}, fdc)
	if len(*emulator.DskAFile) > 0 {
		disc := files.LoadDsk(*emulator.DskAFile)
		fdc.chip.SetDiscA(disc)
		// fmt.Printf("%v\n", disc)
	}
	// if len(*dskBFile) > 0 {
	// 	disc := files.LoadDsk(*dskBFile)
	// 	fdc.SetDiscB(disc)
	// }

	clock := emulator.NewCLock(4000000, 50)
	cpc := &cpc{
		bus:      bus,
		cpu:      z80,
		mem:      mem,
		ga:       ga,
		ppi:      ppi,
		cassette: cas,
		sound:    sound,
		clock:    clock,
	}

	cpc.clock.AddTicker(0, z80)
	cpc.clock.AddTicker(4, crtc)
	cpc.clock.AddTicker(4, ga)
	cpc.clock.AddTicker(4, ay8912)
	cpc.clock.AddTicker(80, sound)
	if *emulator.LoadSlow {
		cpc.clock.AddTicker(0, cas)
	}

	if !*emulator.LoadSlow {
		if cpc464 {
			z80.RegisterTrap(0x2836, cpc.loadTapeBlockCPC464)
		} else {
			z80.RegisterTrap(0x29A6, cpc.loadTapeBlockCPC6128)
		}
	}

	return cpc
}

func (m *cpc) Reset() {
}

func (m *cpc) loadTapeBlockCPC464() {
	m.LoadTapeBlockCPC(0x2872)
}

func (m *cpc) loadTapeBlockCPC6128() {
	m.LoadTapeBlockCPC(0x29E2)
}

func (m *cpc) LoadTapeBlockCPC(exit uint16) {
	data := m.cassette.NextDataBlock()
	if data == nil {
		return
	}

	regs := m.cpu.Registers()
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
	return
}

func (m *cpc) Debugger() emulator.Debugger {
	return m.debugger
}

func (m *cpc) OnKey(key glfw.Key) {
	m.ppi.OnKey(key)
}

func (monitor *monitor) Screen() *ui.Display {
	return nil
}

func (m *cpc) Monitor() emulator.Monitor {
	return m.ga.monitor
}

func (m *cpc) Clock() emulator.Clock {
	return m.clock
}

func (m *cpc) GetVolumeControl() func(float64) {
	return m.sound.SetVolume
}

func (m *cpc) Control() map[string]ui.Control {
	return map[string]ui.Control{"CPU": ui.NewZ80UI(m.cpu)}
}

func (m *cpc) SetDebugger(db cpu.DebuggerCallbacks) { m.cpu.SetDebugger(db) }

type dummyPortsManager struct{}

func (*dummyPortsManager) ReadPort(port uint16) (byte, bool) { return 0, false }
func (*dummyPortsManager) WritePort(port uint16, data byte)  {}
