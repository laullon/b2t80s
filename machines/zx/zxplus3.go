package zx

import (
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/ay8912"
	"github.com/laullon/b2t80s/emulator/files"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/z80"
)

type zxplus3 struct {
	*zx
	ay8912 ay8912.AY8912
}

func NewZXPlus3(cassette emulator.Cassette) machines.Machine {
	mem := NewMemory(ZXPLUS3)
	mem.LoadRom(0, data.MustAsset("data/roms/plus3/Spectrum+3_Spanish_ROM0.bin"))
	mem.LoadRom(1, data.MustAsset("data/roms/plus3/Spectrum+3_Spanish_ROM1.bin"))
	mem.LoadRom(2, data.MustAsset("data/roms/plus3/Spectrum+3_Spanish_ROM2.bin"))
	mem.LoadRom(3, data.MustAsset("data/roms/plus3/Spectrum+3_Spanish_ROM3.bin"))

	ay8912 := ay8912.New()

	ula := NewULA(mem, cassette)
	cpu := z80.NewZ80(ula, cassette)
	clock := emulator.NewCLock(CLOCK_128k)

	fdc := NewZXFDC765()
	cpu.RegisterPort(emulator.PortMask{Mask: 0xC002, Value: 0x0000}, fdc)
	if len(*machines.DskAFile) > 0 {
		disc := files.LoadDsk(*machines.DskAFile)
		fdc.chip.SetDiscA(disc)
	}

	sound := emulator.NewSoundSystem(CLOCK_128k / 80)
	sound.AddSource(ay8912)
	sound.AddSource(ula)

	cpu.RegisterPort(emulator.PortMask{Mask: 0x00FF, Value: 0x00FE}, ula)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x00e0, Value: 0x0000}, &emulator.Kempston{})
	cpu.RegisterPort(emulator.PortMask{Mask: 0xc002, Value: 0xc000}, ay8912)
	cpu.RegisterPort(emulator.PortMask{Mask: 0xc002, Value: 0x8000}, ay8912)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x8002, Value: 0x0000}, mem)
	cpu.RegisterPort(emulator.PortMask{Mask: 0x8002, Value: 0x4000}, mem)

	// if cassette != nil && cassette.Ready() {
	// 	clock.AddTicker(0, cassette)
	// 	if *machines.LoadSlow {
	// 		cpu.RegisterTrap(0x056b, func() uint16 {
	// 			if !cassette.IsMotorON() {
	// 				cassette.Motor(true)
	// 				cassette.Play()
	// 			}
	// 			return emulator.CONTINUE
	// 		})
	// 	} else {
	// 		cpu.RegisterTrap(0x056b, cpu.LoadTapeBlock)
	// 		cpu.RegisterTrap(0x0111, cpu.LoadTapeBlock)
	// 	}
	// }

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
