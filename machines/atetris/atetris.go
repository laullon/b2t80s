package atetris

import (
	"encoding/hex"
	"os"
	"os/signal"
	"syscall"

	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/pokey"
	"github.com/laullon/b2t80s/ui"
)

type atetris struct {
	clock    emulator.Clock
	bus      m6502.Bus
	cpu      m6502.M6502
	sos2     *sos2
	debugger emulator.Debugger
	monitor  emulator.Monitor
	pokey1   *pokey.Pokey
	pokey2   *pokey.Pokey
}

func NewATetris() emulator.Machine {
	bus := m6502.NewBus()

	m := &atetris{
		cpu:    m6502.MewM6502(bus),
		bus:    bus,
		clock:  emulator.NewCLock(14318181/8, 60),
		pokey1: pokey.NewPokey(),
		pokey2: pokey.NewPokey(),
		sos2:   newSOS2(),
	}

	m.monitor = emulator.NewMonitor(m.sos2.display)
	m.sos2.monitor = m.monitor

	rom := loadRom("136066-1100.45f")
	status := &status{}

	wd := &watchdog{cpu: m.cpu}
	wd.start()

	eeprom := newEERPROM()

	cpuRAM := &ram{mem: make([]byte, 0x1000), mask: 0x0fff, trace: false}

	bus.RegisterPort("ram", cpu.PortMask{Mask: 0b1111000000000000, Value: 0b0000000000000000}, cpuRAM)
	bus.RegisterPort("vRam", cpu.PortMask{Mask: 0b1111_0000_0000_0000, Value: 0b0001_0000_0000_0000}, &ram{mem: m.sos2.vram, mask: 0x0fff})
	bus.RegisterPort("color", cpu.PortMask{Mask: 0b1111_1100_0000_0000, Value: 0b0010_0000_0000_0000}, m.sos2.color)

	bus.RegisterPort("eeprom", cpu.PortMask{Mask: 0b1111_1110_0000_0000, Value: 0b0010_0100_0000_0000}, eeprom)

	bus.RegisterPort("pokey0", cpu.PortMask{Mask: 0b1111_1100_0011_0000, Value: 0b0010_1000_0000_0000}, m.pokey1)
	bus.RegisterPort("pokey1", cpu.PortMask{Mask: 0b1111_1100_0011_0000, Value: 0b0010_1000_0001_0000}, m.pokey2)

	bus.RegisterPort("slapstic", cpu.PortMask{Mask: 0b1100_0000_0000_0000, Value: 0b0100_0000_0000_0000}, newSlapstic(rom))
	bus.RegisterPort("rom", cpu.PortMask{Mask: 0b1000_0000_0000_0000, Value: 0b1000_0000_0000_0000}, &fixedROM{rom: rom})

	bus.RegisterPort("eeprom.status", cpu.PortMask{Mask: 0b1111_1100_0000_0000, Value: 0b0011_0100_0000_0000}, eeprom.status)

	bus.RegisterPort("status", cpu.PortMask{Mask: 0b1111_1100_0000_0000, Value: 0b0011_1100_0000_0000}, status)

	bus.RegisterPort("clearIRQ", cpu.PortMask{Mask: 0b1111_1100_0000_0000, Value: 0b0011_1000_0000_0000}, &clearIRQ{cpu: m.cpu})

	bus.RegisterPort("watchdog", cpu.PortMask{Mask: 0b1111110000000000, Value: 0b0011000000000000}, wd)

	m.pokey1.P7 = false

	m.sos2.hBlank = &m.pokey1.P6

	m.sos2.cpu = m.cpu

	m.clock.AddTicker(0, m.cpu)
	m.clock.AddTicker(2, m.sos2)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGALRM)
	go func() {
		s := <-sigc
		println(hex.Dump(cpuRAM.mem))
		println(s)
	}()

	print("bus:\n", bus.DumpMap(), "\n")
	return m
}

func (t *atetris) Reset() {
}

func (t *atetris) Debugger() emulator.Debugger { return t.debugger }

func (t *atetris) Monitor() emulator.Monitor {
	return t.monitor
}

func (t *atetris) Control() map[string]ui.Control {
	return map[string]ui.Control{"CPU": ui.NewM6502UI(t.cpu)}
}

func (t *atetris) Clock() emulator.Clock                { return t.clock }
func (t *atetris) UIControls() []ui.Control             { return []ui.Control{ui.NewM6502BusUI("", t.bus)} }
func (t *atetris) GetVolumeControl() func(float64)      { return func(f float64) {} }
func (t *atetris) SetDebugger(db cpu.DebuggerCallbacks) { t.cpu.SetDebugger(db) }
