package zx

import (
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/z80"
)

type zx48k struct {
	*zx
}

func NewZX48K(cassette cassette.Cassette) machines.Machine {
	mem := NewMemory(ZX48K)
	rom48, err := data.Asset("data/roms/48.rom")
	if err != nil {
		panic("data/48.rom not found")
	}
	mem.LoadRom(0, rom48)

	clock := emulator.NewCLock(clock48k)
	ula := NewULA(mem, cassette, clock, false)
	cpu := z80.NewZ80(ula, cassette)
	ula.cpu = cpu

	sound := emulator.NewSoundSystem(clock48k / 80)
	sound.AddSource(ula)

	cpu.RegisterPort(emulator.PortMask{Mask: 0x00FF, Value: 0x00FE}, ula)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x00FF, Value: 0x00FF}, ula)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x00e0, Value: 0x0000}, &emulator.Kempston{})

	cpu.SetClock(clock)

	clock.AddTicker(0, ula)
	clock.AddTicker(0, cassette)
	clock.AddTicker(80, sound)

	zx := NewZX(cpu, ula, mem, cassette, sound, nil)

	if !*machines.LoadSlow {
		cpu.RegisterTrap(0x056b, zx.loadDataBlock)
		cpu.RegisterTrap(0x12A9, ula.LoadCommand)
	}

	return &zx48k{
		zx: zx,
	}
}
