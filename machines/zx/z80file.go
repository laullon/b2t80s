package zx

import (
	"log"
	"os"

	"github.com/laullon/b2t80s/cpu/z80"
	"github.com/laullon/b2t80s/machines"
)

func LoadZ80File(fileName string) machines.Machine {
	fi, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	file := make([]byte, fi.Size()+1)
	l, err := f.Read(file)
	if err != nil {
		panic(err)
	}
	file = file[0:l]

	version := 1
	h := file[34]
	model := 48
	if getUint16(file[6], file[7]) == 0 {
		if h == 2 {
			log.Panic("'SamRam' not supported")
		}
		if file[30] == 23 {
			version = 2
			if h == 3 || h == 4 {
				model = 128
			}
		} else {
			version = 3
			if h == 4 || h == 5 || h == 6 {
				model = 128
			}
		}
	}

	log.Printf("Loading z80 file '%s' v:%d h:%dk (%d)", fileName, version, model, h)

	var machine *zx
	switch model {
	case 48:
		machine = NewZX48K().(*zx)
	case 128:
		machine = NewZX128K().(*zx)
	}

	regs := machine.cpu.Registers().(*z80.Z80Registers)
	// TODO: byte 12
	regs.A = file[0]
	regs.F.SetByte(file[1])
	regs.B = file[3]
	regs.C = file[2]
	regs.D = file[14]
	regs.E = file[13]
	regs.H = file[5]
	regs.L = file[4]
	regs.IXH = file[26]
	regs.IXL = file[25]
	regs.IYH = file[23]
	regs.IYL = file[24]
	regs.Aalt = file[21]
	regs.Falt.SetByte(file[22])
	regs.Balt = file[16]
	regs.Calt = file[15]
	regs.Dalt = file[18]
	regs.Ealt = file[17]
	regs.Halt = file[20]
	regs.Lalt = file[19]

	regs.I = file[10]
	regs.R = file[11]
	regs.IFF1 = file[27] != 0
	regs.InterruptsMode = file[29] & 3
	regs.SP.Set(getUint16(file[9], file[8]))

	if version == 1 {
		regs.PC = getUint16(file[7], file[6])
		data := file[30:]
		mem := make(bank, 0xc000)
		copyMemoryBlock(data, uint16(len(data)), &mem)
		copy(machine.mem.(*memory).banks[0], mem[0x0000:])
		copy(machine.mem.(*memory).banks[1], mem[0x4000:])
		copy(machine.mem.(*memory).banks[2], mem[0x8000:])
	} else {
		regs.PC = getUint16(file[33], file[32])
		block := file[30+file[30]+2:]
		for len(block) > 0 {
			len := getUint16(block[1], block[0])
			if len == 0xffff {
				len = 0x4000
			}
			page := block[2]
			data := block[3 : 3+len]

			var bank *bank
			if model == 48 {
				switch page {
				case 4:
					bank = &machine.mem.(*memory).banks[1]
				case 5:
					bank = &machine.mem.(*memory).banks[2]
				case 8:
					bank = &machine.mem.(*memory).banks[0]
				default:
					log.Panicf("-- page '%d' not supported --", page)
				}
			} else if model == 128 {
				bank = &machine.mem.(*memory).banks[page-3]
			}

			copyMemoryBlock(data, len, bank)
			block = block[3+len:]
		}
		if model == 128 {
			// ay := machine.(*zx128k).ay8912
			// for r, b := range file[39:45] {
			// 	ay.WriteRegister(byte(r), b)
			// }
			// ay.SetReg(file[38])
		}
	}

	return machine
}

func copyMemoryBlock(memOrg []byte, len uint16, bank *bank) {
	posScr := uint16(0)
	posDst := uint16(0)
	// log.Printf("copying %d bytes to page 0x%04X\n", len, posDst)
	for posScr < len {
		if memOrg[posScr] == 0xED && memOrg[posScr+1] == 0xED && posScr+3 < len {
			b := memOrg[posScr+3]
			c := uint16(memOrg[posScr+2])
			for i := uint16(0); i < c; i++ {
				(*bank)[posDst+i] = b
			}
			posDst += c
			posScr += 4
		} else {
			(*bank)[posDst] = memOrg[posScr]
			posDst++
			posScr++
		}
	}
}

func getUint16(h, l byte) uint16 {
	return (uint16(h) << 8) | uint16(l)
}
