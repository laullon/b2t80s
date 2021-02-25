package zx

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/cpu/z80"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/ay8912"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
	"github.com/laullon/b2t80s/ui"
)

const (
	clock48k  = uint(3500000)
	clock128k = uint(3546900)
)

type ZX interface {
	Debugger() emulator.Debugger
	OnKeyEvent(event *fyne.KeyEvent)
	LoadZ80File(fileName string)
}

type zx struct {
	bus      z80.Bus
	ula      *ula
	cpu      z80.Z80
	mem      z80.Memory
	cassette cassette.Cassette
	sound    emulator.SoundSystem
	debugger emulator.Debugger
	clock    emulator.Clock
	ay8912   ay8912.AY8912
}

func NewZX(mem *memory, plus, cas, ay bool) *zx {
	speed := clock48k
	if plus {
		speed = clock128k
	}

	bus := z80.NewBus(mem)

	z80 := z80.NewZ80(bus)
	clock := emulator.NewCLock(speed, 50)

	ula := NewULA(mem, bus, plus)
	sound := emulator.NewSoundSystem(speed / uint(80))

	ula.cpu = z80
	sound.AddSource(ula)

	clock.AddTicker(0, ula)
	clock.AddTicker(80, sound)

	bus.RegisterPort(cpu.PortMask{Mask: 0x00FF, Value: 0x00FE}, ula)
	bus.RegisterPort(cpu.PortMask{Mask: 0x00FF, Value: 0x00FF}, ula)
	bus.RegisterPort(cpu.PortMask{Mask: 0x00e0, Value: 0x0000}, &kempston{})

	zx := &zx{
		bus:   bus,
		ula:   ula,
		cpu:   z80,
		mem:   mem,
		sound: sound,
		clock: clock,
	}

	if ay {
		zx.ay8912 = ay8912.New()
		sound.AddSource(zx.ay8912)
		bus.RegisterPort(cpu.PortMask{Mask: 0xc002, Value: 0xc000}, zx.ay8912)
		bus.RegisterPort(cpu.PortMask{Mask: 0xc002, Value: 0x8000}, zx.ay8912)
		clock.AddTicker(2, zx.ay8912)
	}

	if cas {
		zx.cassette = cassette.New()
		if len(*emulator.TapFile) > 0 {
			zx.cassette.LoadTapFile(*emulator.TapFile)
		}
		clock.AddTicker(0, zx.cassette)
		ula.cassette = zx.cassette
	}

	return zx
}

func (zx *zx) Debugger() emulator.Debugger {
	return zx.debugger
}

func (zx *zx) OnKeyEvent(event *fyne.KeyEvent) {
	zx.ula.OnKeyEvent(event)
}

func (zx *zx) Monitor() emulator.Monitor {
	return zx.ula.monitor
}

func (zx *zx) GetVolumeControl() func(float64) {
	return zx.sound.SetVolume
}

func (zx *zx) Clock() emulator.Clock {
	return zx.clock
}

func (zx *zx) CPUControl() ui.Control               { return ui.NewZ80UI(zx.cpu.Registers()) }
func (zx *zx) SetDebugger(db cpu.DebuggerCallbacks) { zx.cpu.SetDebugger(db) }

func (zx *zx) loadDataBlock() {
	data := zx.cassette.NextDataBlock()
	if data == nil {
		return //emulator.CONTINUE
	}

	regs := zx.cpu.Registers()
	requestedLength := regs.DE.Get()
	startAddress := regs.IX.Get()
	fmt.Printf("Loading block to 0x%04x \n", startAddress)

	zx.cpu.Wait(true)
	go func() {
		if regs.Aalt == data[0] {
			if regs.Falt.C {
				checksum := data[0]
				for i := uint16(0); i < requestedLength; i++ {
					loadedByte := data[i+1]
					zx.mem.PutByte(startAddress+i, loadedByte)
					checksum ^= loadedByte
					if (startAddress == 0x4000) && (i < 0x1b00) {
						time.Sleep(time.Millisecond / 2)
					}
				}
				regs.F.C = checksum == data[requestedLength+1]
			} else {
				regs.F.C = true
			}
			// log.Print("done")
		} else {
			regs.F.C = false
			// log.Print("BAD Block")
		}

		println("-")
		if startAddress == 0x4000 {
			time.Sleep(5 * time.Second)
		}
		regs.PC = 0x05e2
		zx.cpu.Wait(false)
		println("done")
	}()

	return
}

type kempston struct {
}

func (k *kempston) ReadPort(port uint16) (byte, bool) {
	j, _ := emulator.ReadJoystick()
	res := byte(0)
	// 000FUDLR
	if j.ON {
		if j.F {
			res |= 0b00010000
		}
		if j.U {
			res |= 0b00001000
		}
		if j.D {
			res |= 0b00000100
		}
		if j.L {
			res |= 0b00000010
		}
		if j.R {
			res |= 0b00000001
		}
	}

	return res, false
}

func (k *kempston) WritePort(port uint16, data byte) {
}
