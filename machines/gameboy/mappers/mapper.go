package mappers

import (
	"fmt"

	"github.com/laullon/b2t80s/cpu"
)

type Mapper interface {
	ConnectToCPU(bus cpu.Bus)
}

func CreateMapper(fileName string) Mapper {
	file := loadFile(fileName)
	switch file.header.mapper {
	case 0:
		return &rom{file}

	case 1:
		return newMBC1(file)

	default:
		panic(fmt.Sprintf("mapper type '%d' not supported", file.header.mapper))
	}
}

type rom struct {
	file *gbFile
}

func (rom *rom) ConnectToCPU(bus cpu.Bus) {
	bus.RegisterPort("rom", cpu.PortMask{0b1000_0000_0000_0000, 0b0000_0000_0000_0000}, cpu.NewROM(rom.file.data, 0x7fff))
}
