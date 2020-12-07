package zx

import (
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/files"
	"github.com/laullon/b2t80s/machines"
)

func NewZXPlus3() machines.Machine {
	mem := NewMemory(ZXPLUS3)
	mem.LoadRom(0, data.MustAsset("data/roms/plus3/Spectrum+3_Spanish_ROM0.bin"))
	mem.LoadRom(1, data.MustAsset("data/roms/plus3/Spectrum+3_Spanish_ROM1.bin"))
	mem.LoadRom(2, data.MustAsset("data/roms/plus3/Spectrum+3_Spanish_ROM2.bin"))
	mem.LoadRom(3, data.MustAsset("data/roms/plus3/Spectrum+3_Spanish_ROM3.bin"))

	fdc := NewZXFDC765()
	if len(*machines.DskAFile) > 0 {
		disc := files.LoadDsk(*machines.DskAFile)
		fdc.chip.SetDiscA(disc)
	}

	zx := NewZX(mem, true, false, true)

	zx.ula.tsPerRow = 228
	zx.ula.scanlines = 311
	zx.ula.displayStart = 63

	zx.bus.RegisterPort(emulator.PortMask{Mask: 0x8002, Value: 0x0000}, mem)
	zx.bus.RegisterPort(emulator.PortMask{Mask: 0x8002, Value: 0x4000}, mem)
	zx.bus.RegisterPort(emulator.PortMask{Mask: 0xC002, Value: 0x0000}, fdc)

	return zx
}
