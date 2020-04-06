package zx

import (
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/z80"
)

type zx48k struct {
	*zx
}

func NewZX48K(cassette emulator.Cassette) machines.Machine {
	mem := NewMemory(ZX48K)
	rom48, err := data.Asset("data/roms/48.rom")
	if err != nil {
		panic("data/48.rom not found")
	}
	mem.LoadRom(0, rom48)

	clock := emulator.NewCLock(CLOCK_48k)
	ula := NewULA(mem, cassette, clock, false)
	cpu := z80.NewZ80(ula, cassette)
	ula.cpu = cpu

	sound := emulator.NewSoundSystem(CLOCK_48k / 80)
	sound.AddSource(ula)

	cpu.RegisterPort(emulator.PortMask{Mask: 0x00FF, Value: 0x00FE}, ula)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x00e0, Value: 0x0000}, &emulator.Kempston{})

	if cassette != nil && cassette.Ready() {
		clock.AddTicker(0, cassette)
		if *machines.LoadSlow {
			cpu.RegisterTrap(0x056b, func() uint16 {
				if !cassette.IsMotorON() {
					cassette.Motor(true)
					cassette.Play()
				}
				return emulator.CONTINUE
			})
		} else {
			cpu.RegisterTrap(0x056b, cpu.LoadTapeBlock)
			cpu.RegisterTrap(0x12A9, ula.LoadCommand)
		}
	}

	cpu.SetClock(clock)

	clock.AddTicker(0, ula)
	clock.AddTicker(80, sound)

	return &zx48k{
		zx: NewZX(cpu, ula, mem, cassette, sound, nil),
	}
}
