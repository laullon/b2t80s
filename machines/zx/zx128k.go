package zx

import (
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
)

func NewZX128K() emulator.Machine {
	mem := NewMemory(ZX128K)
	mem.LoadRom(0, data.MustAsset("data/roms/128-0.rom"))
	mem.LoadRom(1, data.MustAsset("data/roms/128-1.rom"))

	zx := NewZX(mem, true, true, true)
	if !*emulator.LoadSlow {
		zx.cpu.RegisterTrap(0x056b, zx.loadDataBlock)
		zx.cpu.RegisterTrap(0x3683, zx.ula.LoadCommand128)
	}

	zx.ula.tsPerRow = 228
	zx.ula.scanlines = 311
	zx.ula.displayStart = 63

	zx.bus.RegisterPort(cpu.PortMask{Mask: 0x8002, Value: 0x0000}, mem)

	return zx
}
