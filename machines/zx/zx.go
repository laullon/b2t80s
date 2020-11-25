package zx

import (
	"fyne.io/fyne"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/ay8912"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/z80"
)

const (
	clock48k  = 3500000
	clock128k = 3546900
)

type ZX interface {
	Debugger() emulator.Debugger
	OnKeyEvent(event *fyne.KeyEvent)
	LoadZ80File(fileName string)
}

type zx struct {
	ula      *ula
	cpu      emulator.CPU
	mem      emulator.Memory
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

	cpu := z80.NewZ80(mem)
	clock := emulator.NewCLock(speed)

	ula := NewULA(mem, plus)
	sound := emulator.NewSoundSystem(speed / 80)

	ula.cpu = cpu
	sound.AddSource(ula)

	// clock.AddTicker(0, cpu)
	clock.AddTicker(0, ula)
	clock.AddTicker(80, sound)

	cpu.RegisterPort(emulator.PortMask{Mask: 0x00FF, Value: 0x00FE}, ula)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x00FF, Value: 0x00FF}, ula)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x00e0, Value: 0x0000}, &kempston{})

	zx := &zx{
		ula:      ula,
		cpu:      cpu,
		mem:      mem,
		sound:    sound,
		clock:    clock,
		debugger: z80.NewDebugger(cpu, mem),
	}

	if ay {
		zx.ay8912 = ay8912.New()
		sound.AddSource(zx.ay8912)
		cpu.RegisterPort(emulator.PortMask{Mask: 0xc002, Value: 0xc000}, zx.ay8912)
		cpu.RegisterPort(emulator.PortMask{Mask: 0xc002, Value: 0x8000}, zx.ay8912)
		clock.AddTicker(2, zx.ay8912)
	}

	if cas {
		zx.cassette = cassette.New()
		if len(*machines.TapFile) > 0 {
			zx.cassette.LoadTapFile(*machines.TapFile)
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

func (zx *zx) loadDataBlock() uint16 {
	data := zx.cassette.NextDataBlock()
	if data == nil {
		return emulator.CONTINUE
	}

	regs := zx.cpu.Registers().(*z80.Z80Registers)
	requestedLength := regs.DE.Get()
	startAddress := regs.IX.Get()
	// fmt.Printf("Loading block '%s' to 0x%04x (bl:0x%04x, l:0x%04x, bt:%d, a:%d)\n", block.Name(), startAddress, len(block.GetData()), requestedLength, block.Type(), regs._A)

	if regs.Aalt == data[0] {
		if regs.Falt.C {
			checksum := data[0]
			for i := uint16(0); i < requestedLength; i++ {
				loadedByte := data[i+1]
				zx.mem.PutByte(startAddress+i, loadedByte)
				checksum ^= loadedByte
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
	return 0x05e2
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
