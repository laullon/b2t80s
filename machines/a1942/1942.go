package a1942

import (
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/cpu/z80"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/ay8912"
	"github.com/laullon/b2t80s/gui"
	"github.com/laullon/b2t80s/ui"
)

type a1942 struct {
	mainMem cpu.Bus
	mainBus z80.Bus
	mainCpu z80.Z80

	audioMem cpu.Bus
	audioBus z80.Bus
	audioCpu z80.Z80

	romBank  cpu.ROM
	romBanks [][]byte

	clock    emulator.Clock
	debugger emulator.Debugger
	monitor  emulator.Monitor
	video    *video

	ay1 ay8912.AY8912
	ay2 ay8912.AY8912

	sys, p1, p2 byte
}

func New1942() emulator.Machine {

	m := &a1942{
		p1:  0xff,
		p2:  0xff,
		sys: 0xff,
	}

	m.video = newVideo(m)

	m.clock = emulator.NewCLock(12_000_000, 60)

	m.mainMem = cpu.NewBus("mainMem", &unused{})
	mainPorts := cpu.NewBus("mainPorts", &unused{})
	m.mainBus = z80.NewBus(m.mainMem, mainPorts)
	m.mainCpu = z80.NewZ80(m.mainBus)

	m.audioMem = cpu.NewBus("audioMem", &unused{})
	m.audioBus = z80.NewBus(m.audioMem, nil)
	m.audioCpu = z80.NewZ80(m.audioBus)

	m.monitor = emulator.NewMonitor(m.video.display)

	m.ay1 = ay8912.New()
	m.ay2 = ay8912.New()

	m.romBanks = append(m.romBanks, loadRom("srb-05.m5"))
	m.romBanks = append(m.romBanks, loadRom("srb-06.m6"))
	m.romBanks = append(m.romBanks, loadRom("srb-07.m7"))
	m.romBank = cpu.NewROM(m.romBanks[0], 0x3fff)

	ayControl := &ayControl{m}
	latch := &latch{m: m}

	// MAIN
	m.mainMem.RegisterPort("srb-03.m3", cpu.PortMask{Mask: 0b1100_0000_0000_0000, Value: 0x0000}, cpu.NewROM(loadRom("srb-03.m3"), 0x3fff))
	m.mainMem.RegisterPort("srb-04.m4", cpu.PortMask{Mask: 0b1100_0000_0000_0000, Value: 0x4000}, cpu.NewROM(loadRom("srb-04.m4"), 0x3fff))
	m.mainMem.RegisterPort("romBank", cpu.PortMask{Mask: 0b1100_0000_0000_0000, Value: 0x8000}, m.romBank)

	m.mainMem.RegisterPort("RAM", cpu.PortMask{Mask: 0b1111_0000_0000_0000, Value: 0xe000}, cpu.NewRAM(make([]byte, 0x1000), 0x0fff))
	m.mainMem.RegisterPort("ports", cpu.PortMask{Mask: 0b1111_1111_1111_1100, Value: 0xc000}, m)
	m.mainMem.RegisterPort("ports", cpu.PortMask{Mask: 0b1111_1111_1111_1111, Value: 0xc004}, m)
	m.mainMem.RegisterPort("background scroll", cpu.PortMask{Mask: 0b1111_1111_1111_1111, Value: 0xC802}, m.video)
	m.mainMem.RegisterPort("background scroll", cpu.PortMask{Mask: 0b1111_1111_1111_1111, Value: 0xC803}, m.video)
	m.mainMem.RegisterPort("0xC804", cpu.PortMask{Mask: 0b1111_1111_1111_1111, Value: 0xC804}, m)
	m.mainMem.RegisterPort("background palette", cpu.PortMask{Mask: 0b1111_1111_1111_1111, Value: 0xC805}, m.video)
	m.mainMem.RegisterPort("bankswitch", cpu.PortMask{Mask: 0b1111_1111_1111_1111, Value: 0xC806}, m)
	m.mainMem.RegisterPort("spriteram", cpu.PortMask{Mask: 0b1111_1111_1000_0000, Value: 0xcc00}, m.video.spriteram)
	m.mainMem.RegisterPort("fgvram", cpu.PortMask{Mask: 0b1111_1000_0000_0000, Value: 0xd000}, m.video.fgvram)
	m.mainMem.RegisterPort("bgvram", cpu.PortMask{Mask: 0b1111_1100_0000_0000, Value: 0xd800}, m.video.bgvram)

	// AUDIO
	m.audioMem.RegisterPort("sr-01.c11", cpu.PortMask{Mask: 0b1100_0000_0000_0000, Value: 0x0000}, cpu.NewROM(loadRom("sr-01.c11"), 0x3fff))

	m.audioMem.RegisterPort("RAM", cpu.PortMask{Mask: 0b1111_1000_0000_0000, Value: 0x4000}, cpu.NewRAM(make([]byte, 0x0800), 0x07ff))
	m.audioMem.RegisterPort("latch", cpu.PortMask{Mask: 0b1111_1111_1111_1111, Value: 0x6000}, latch)
	m.audioMem.RegisterPort("AY1", cpu.PortMask{Mask: 0b1111_1111_1111_1110, Value: 0x8000}, ayControl)
	m.audioMem.RegisterPort("AY2", cpu.PortMask{Mask: 0b1111_1111_1111_1110, Value: 0xc000}, ayControl)

	print("main bus:\n", m.mainMem.DumpMap(), "\n")
	print("audio bus:\n", m.audioMem.DumpMap(), "\n")

	m.clock.AddTicker(3, m.mainCpu)  // 4Mhz
	m.clock.AddTicker(4, m.audioCpu) // 3Mhz
	m.clock.AddTicker(2, m.video)    // 6Mhz

	return m
}

func (t *a1942) Reset() {
}

func (t *a1942) Debugger() emulator.Debugger { return t.debugger }

func (t *a1942) Monitor() emulator.Monitor {
	return t.monitor
}

func (t *a1942) Control() map[string]gui.GUIObject {
	return map[string]gui.GUIObject{
		"Main CPU":     ui.NewZ80UI(t.mainCpu, true),
		"Audio CPU":    ui.NewZ80UI(t.audioCpu, false),
		"Main Memory":  ui.NewBusUI(t.mainMem),
		"Audio Memory": ui.NewBusUI(t.audioMem),
		"char":         newCharactersUI(t.video),
	}
}

func (t *a1942) SetDebugger(db cpu.DebuggerCallbacks) {
	t.mainCpu.SetDebugger(db)
	t.audioCpu.SetDebugger(db)
}

func (t *a1942) ReadPort(port uint16) (byte, bool) {
	switch port {
	case 0xc000:
		return t.sys, false
	case 0xc001:
		return t.p1, false
	case 0xc002:
		return t.p2, false
	case 0xc004:
		if *emulator.Test {
			return 0xF7, false // TEST
		}
	}
	return 0xff, false
}

func (m *a1942) WritePort(port uint16, data byte) {
	switch port {
	case 0xC804:
		// TODO bit 7: flip screen bit 4: cpu B reset bit 0: coin counter *

	case 0xc806:
		m.romBank.SetBank(m.romBanks[data&3])
	}
}

func (t *a1942) Clock() emulator.Clock           { return t.clock }
func (t *a1942) UIControls() []gui.GUIObject     { return nil } // []gui.GUIObject{ui.NewM6502BusUI("", t.bus)} }
func (t *a1942) GetVolumeControl() func(float64) { return func(f float64) {} }

// *******
type unused struct{}

func (*unused) ReadPort(port uint16) (byte, bool) { return 0xff, false }
func (*unused) WritePort(port uint16, data byte)  {}

// *******
type ayControl struct {
	m *a1942
}

func (ayc *ayControl) ReadPort(port uint16) (byte, bool) { panic(-1) }
func (ayc *ayControl) WritePort(port uint16, data byte) {
	switch port {
	case 0x8000, 0x8001:
		ayc.m.ay1.WriteRegister(byte(port&1), data)
	case 0xc000, 0xc001:
		ayc.m.ay2.WriteRegister(byte(port&1), data)
	}
}

// *******
type latch struct {
	m *a1942
	v byte
}

func (l *latch) ReadPort(port uint16) (byte, bool) { return l.v, false }
func (l *latch) WritePort(port uint16, data byte) {
	l.v = data
	l.m.audioCpu.NMI(true)
}
