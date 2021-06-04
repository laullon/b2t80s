package a1942

import (
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/cpu/z80"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/ay8912"
	"github.com/laullon/b2t80s/gui"
	"github.com/veandco/go-sdl2/sdl"
)

type a1942 struct {
	mainBus z80.Bus
	mainCpu z80.Z80

	audioBus z80.Bus
	audioCpu z80.Z80

	romBank cpu.ROM

	clock    emulator.Clock
	debugger emulator.Debugger
	monitor  emulator.Monitor

	ay1 ay8912.AY8912
	ay2 ay8912.AY8912

	display *gui.Display
}

func New1942() emulator.Machine {
	m := &a1942{
		display: gui.NewDisplay(gui.Size{336, 240}),
	}

	m.clock = emulator.NewCLock(12_000_000, 60)

	mainMem := cpu.NewBus("mainMem")
	mainPorts := cpu.NewBus("mainPorts", &unused{})
	m.mainBus = z80.NewBus(mainMem, mainPorts)
	m.mainCpu = z80.NewZ80(m.mainBus)

	audioMem := cpu.NewBus("audioMem", &unused{})
	// audioPorts := cpu.NewBus("audioPorts")
	m.audioBus = z80.NewBus(audioMem, nil)
	m.audioCpu = z80.NewZ80(m.audioBus)

	m.monitor = emulator.NewMonitor(m.display)

	m.ay1 = ay8912.New()
	m.ay2 = ay8912.New()

	m.romBank = cpu.NewROM(loadRom("srb-05.m5"), 0x3fff)

	ayControl := &ayControl{m}
	latch := &latch{m: m}

	// MAIN
	mainMem.RegisterPort("srb-03.m3", cpu.PortMask{Mask: 0b1100_0000_0000_0000, Value: 0x0000}, cpu.NewROM(loadRom("srb-03.m3"), 0x3fff))
	mainMem.RegisterPort("srb-04.m4", cpu.PortMask{Mask: 0b1100_0000_0000_0000, Value: 0x4000}, cpu.NewROM(loadRom("srb-04.m4"), 0x3fff))
	mainMem.RegisterPort("romBank", cpu.PortMask{Mask: 0b1100_0000_0000_0000, Value: 0x8000}, m.romBank)

	mainMem.RegisterPort("RAM", cpu.PortMask{Mask: 0b1111_0000_0000_0000, Value: 0xe000}, cpu.NewRAM(make([]byte, 0x1000), 0x0fff))
	mainMem.RegisterPort("unused", cpu.PortMask{Mask: 0b1111_0000_0000_0000, Value: 0xF000}, &unused{})
	mainMem.RegisterPort("ports", cpu.PortMask{Mask: 0b1111_1111_1111_1100, Value: 0xc000}, m)
	mainMem.RegisterPort("ports", cpu.PortMask{Mask: 0b1111_1111_1111_1111, Value: 0xc004}, m)
	mainMem.RegisterPort("fgvram", cpu.PortMask{Mask: 0b1111_1000_0000_0000, Value: 0xd000}, cpu.NewRAM(make([]byte, 0x0800), 0x07ff))
	mainMem.RegisterPort("bgvram", cpu.PortMask{Mask: 0b1111_1100_0000_0000, Value: 0xd800}, cpu.NewRAM(make([]byte, 0x0800), 0x07ff))

	// AUDIO
	audioMem.RegisterPort("sr-01.c11", cpu.PortMask{Mask: 0b1100_0000_0000_0000, Value: 0x0000}, cpu.NewROM(loadRom("sr-01.c11"), 0x3fff))

	audioMem.RegisterPort("RAM", cpu.PortMask{Mask: 0b1111_1000_0000_0000, Value: 0x4000}, cpu.NewRAM(make([]byte, 0x0800), 0x07ff))
	audioMem.RegisterPort("latch", cpu.PortMask{Mask: 0b1111_1111_1111_1111, Value: 0x6000}, latch)
	audioMem.RegisterPort("AY1", cpu.PortMask{Mask: 0b1111_1111_1111_1110, Value: 0x8000}, ayControl)
	audioMem.RegisterPort("AY2", cpu.PortMask{Mask: 0b1111_1111_1111_1110, Value: 0xc000}, ayControl)

	print("main bus:\n", mainMem.DumpMap(), "\n")
	print("audio bus:\n", audioMem.DumpMap(), "\n")

	m.clock.AddTicker(3, m.mainCpu)  // 4Mhz
	m.clock.AddTicker(4, m.audioCpu) // 3Mhz

	return m
}

func (t *a1942) Reset() {
}

func (t *a1942) Debugger() emulator.Debugger { return t.debugger }

func (t *a1942) Monitor() emulator.Monitor {
	return t.monitor
}

func (t *a1942) Control() map[string]gui.GUIObject {
	return nil //map[string]gui.GUIObject{"CPU": ui.NewM6502UI(t.cpu)}
}

func (t *a1942) SetDebugger(db cpu.DebuggerCallbacks) {
	t.mainCpu.SetDebugger(db)
	t.audioCpu.SetDebugger(db)
}

func (t *a1942) Clock() emulator.Clock           { return t.clock }
func (t *a1942) UIControls() []gui.GUIObject     { return nil } // []gui.GUIObject{ui.NewM6502BusUI("", t.bus)} }
func (t *a1942) GetVolumeControl() func(float64) { return func(f float64) {} }
func (t *a1942) OnKey(key sdl.Scancode)          {}

func (t *a1942) ReadPort(port uint16) (byte, bool) { return 0xff, false }
func (t *a1942) WritePort(port uint16, data byte)  { panic(-1) }

// *******
type unused struct{}

func (_ *unused) ReadPort(port uint16) (byte, bool) { return 0xff, false }
func (_ *unused) WritePort(port uint16, data byte)  {}

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
