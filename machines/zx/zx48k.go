package zx

import (
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
)

func NewZX48K() emulator.Machine {
	mem := NewMemory(ZX48K)
	mem.LoadRom(0, data.MustAsset("data/roms/48.rom"))
	// mem.LoadRom(0, data.MustAsset("data/roms/zx_testrom.bin"))

	zx := NewZX(mem, false, true, false)
	// if !*emulator.LoadSlow {
	// 	zx.cpu.RegisterTrap(0x056b, zx.loadDataBlock)
	// 	zx.cpu.RegisterTrap(0x12A9, zx.ula.LoadCommand)
	// }

	return zx
}
