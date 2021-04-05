package mappers

import (
	"fmt"

	"github.com/laullon/b2t80s/cpu"
)

type mbc1 struct {
	file *gbFile
	rom1 cpu.ROM
	rom2 cpu.ROM
	ram  cpu.RAM
}

func newMBC1(file *gbFile, ram bool) Mapper {
	mbc1 := &mbc1{
		file: file,
	}

	mbc1.rom1 = cpu.NewROM(file.data[:0x4000], 0x3fff, mbc1.write)
	mbc1.rom2 = cpu.NewROM(file.data[0x4000:0x8000], 0x3fff, mbc1.write)

	if mbc1.file.header.mapper == 2 {
		mbc1.ram = cpu.NewRAM(make([]byte, 0x2000), 0x1fff)
	}

	return mbc1
}

func (mbc1 *mbc1) ConnectToCPU(bus cpu.Bus) {
	bus.RegisterPort("rom_00", cpu.PortMask{0b1100_0000_0000_0000, 0b0000_0000_0000_0000}, mbc1.rom1)
	bus.RegisterPort("rom_01", cpu.PortMask{0b1100_0000_0000_0000, 0b0100_0000_0000_0000}, mbc1.rom2)
	if mbc1.file.header.ramSize != 0 {
		panic(-1)
	} else if mbc1.file.header.mapper == 2 {
		bus.RegisterPort("ram", cpu.PortMask{0b1110_0000_0000_0000, 0b1010_0000_0000_0000}, mbc1.ram)
	}
}

func (mbc1 *mbc1) write(addr uint16, data uint8) {
	if addr < 0x2000 {
		// TODO: RAM Enable

	} else if addr < 0x4000 {
		mask := (mbc1.file.header.romSize * 4) - 1
		bank := uint32(data & mask)
		// fmt.Printf("[mbc1] write bank: %d of %d (%d)\n", bank, len(mbc1.file.data)/0x4000, len(mbc1.file.data))

		start := uint32(0x4000) * bank
		end := start + 0x4000
		mbc1.rom2.SetBank(mbc1.file.data[start:end])
	} else if addr < 0x6000 {
		if mbc1.file.header.ramSize != 0 {
			panic(-1)
		}
	} else {
		panic(fmt.Sprintf("[mbc1] write 0x%04X 0x%02X \n", addr, data))
	}
}
