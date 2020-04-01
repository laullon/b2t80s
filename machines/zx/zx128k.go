package zx

import (
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/ay8912"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/z80"
)

type zx128k struct {
	*zx
	ay8912 ay8912.AY8912
}

func NewZX128K(cassette emulator.Cassette) machines.Machine {
	mem := NewMemory(ZX128K)
	mem.LoadRom(0, data.MustAsset("data/roms/128-0.rom"))
	mem.LoadRom(1, data.MustAsset("data/roms/128-1.rom"))

	ay8912 := ay8912.New()

	ula := NewULA(mem, cassette)
	cpu := z80.NewZ80(mem, cassette)
	clock := emulator.NewCLock(CLOCK_128k)

	sound := emulator.NewSoundSystem(CLOCK_128k / 80)
	sound.AddSource(ay8912)
	sound.AddSource(ula)

	cpu.RegisterPort(emulator.PortMask{Mask: 0x00FF, Value: 0x00FE}, ula)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x00e0, Value: 0x0000}, &emulator.Kempston{})
	cpu.RegisterPort(emulator.PortMask{Mask: 0xc002, Value: 0xc000}, ay8912)
	cpu.RegisterPort(emulator.PortMask{Mask: 0xc002, Value: 0x8000}, ay8912)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x8002, Value: 0x0000}, mem)

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
			cpu.RegisterTrap(0x0111, cpu.LoadTapeBlock)
		}
	}

	cpu.SetClock(clock)
	mem.SetClock(clock)

	clock.AddTicker(0, ula)
	clock.AddTicker(2, ay8912)
	clock.AddTicker(80, sound)

	return &zx128k{
		ay8912: ay8912,
		zx:     NewZX(cpu, ula, mem, cassette, sound, nil),
	}
}
