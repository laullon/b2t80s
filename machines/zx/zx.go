package zx

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"fyne.io/fyne"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/z80"
)

const (
	CLOCK_48k  = 3500000
	CLOCK_128k = 3546900
)

type ZX interface {
	Debugger() emulator.Debugger
	OnKeyEvent(event *fyne.KeyEvent)
	LoadZ80File(fileName string)
}

type zx struct {
	ula      *ula
	cpu      emulator.CPU
	mem      emulator.Memory
	cassete  cassette.Cassette
	sound    emulator.SoundSystem
	debugger emulator.Debugger

	onEndFrame func()
}

func NewZX(cpu emulator.CPU, ula *ula, mem emulator.Memory, cassete cassette.Cassette, sound emulator.SoundSystem, onEndFrame func()) *zx {
	zx := &zx{
		ula:        ula,
		cpu:        cpu,
		mem:        mem,
		cassete:    cassete,
		sound:      sound,
		debugger:   z80.NewDebugger(cpu, mem),
		onEndFrame: onEndFrame,
	}

	return zx
}

func (m *zx) Run() {
	wait := time.Duration(20 * time.Millisecond)
	runStart := time.Now()
	frames := float64(0)

	ticker := time.NewTicker(wait)
	go func() {
		for range ticker.C {
			frameStart := time.Now()

			m.Debugger().NextFrame()
			err := m.cpu.RunFrame()
			if err != nil {
				panic(err)
			}

			frames++

			frameTime := time.Now().Sub(frameStart)
			runTime := time.Now().Sub(runStart)
			m.Debugger().SetStatus(fmt.Sprintf("frame rate:%6.2f time:%6.2fms (%v)", frames/runTime.Seconds(), float64(frameTime.Microseconds())/1000, wait))
			m.ula.FrameDone()
			if m.onEndFrame != nil {
				m.onEndFrame()
			}
		}
	}()
}

func (m *zx) Debugger() emulator.Debugger {
	return m.debugger
}

func (m *zx) OnKeyEvent(event *fyne.KeyEvent) {
	m.ula.OnKeyEvent(event)
}

func (m *zx) Display() image.Image {
	return m.ula.Display()
}

func (m *zx) GetVolumeControl() func(float64) {
	return m.sound.SetVolume
}

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

	// log.Printf("Loading z80 file '%s' v:%d h:%dk (%d)", fileName, version, model, h)

	var machine machines.Machine
	// var cpu emulator.CPU
	var mem emulator.Memory
	switch model {
	case 48:
		machine = NewZX48K(nil)
		// cpu = machine.(*zx48k).cpu
		mem = machine.(*zx48k).mem
	case 128:
		machine = NewZX128K(nil)
		// cpu = machine.(*zx128k).cpu
		mem = machine.(*zx128k).mem
	}

	// TODO: byte 12
	regs := append([]byte{}, file[0], file[1]) // AF
	regs = append(regs, file[3], file[2])      // BC
	regs = append(regs, file[14], file[13])    // DE
	regs = append(regs, file[5], file[4])      // HL
	regs = append(regs, file[26], file[25])    // IX
	regs = append(regs, file[23], file[24])    // IY
	regs = append(regs, file[21], file[22])    // _AF
	regs = append(regs, file[16], file[15])    // _BC
	regs = append(regs, file[18], file[17])    // _DE
	regs = append(regs, file[20], file[19])    // _HL

	// cpu.SetRegisters(regs, file[10], file[11], file[27], file[29]&3)
	// cpu.SP().Set(getUint16(file[9], file[8]))

	if version == 1 {
		// pc := getUint16(file[7], file[6])
		// cpu.SetPC(pc)
		data := file[30:]
		copyMemoryBlock(data, uint16(len(data)), mem, uint16(0x4000))
	} else {
		// pc := getUint16(file[33], file[32])
		// cpu.SetPC(pc)

		block := file[30+file[30]+2:]
		for len(block) > 0 {
			len := getUint16(block[1], block[0])
			if len == 0xffff {
				len = 0x4000
			}
			page := block[2]
			data := block[3 : 3+len]

			posDst := uint16(0)
			if model == 48 {
				switch page {
				case 4:
					posDst = 0x8000
				case 5:
					posDst = 0xc000
				case 8:
					posDst = 0x4000
				default:
					log.Panicf("-- page '%d' not supported --", page)
				}
			} else if model == 128 {
				// mem.Paging(page - 3)
				posDst = 0xc000
			}

			copyMemoryBlock(data, len, mem, posDst)
			block = block[3+len:]
		}
		if model == 128 {
			// mem.Paging(file[35])
			ay := machine.(*zx128k).ay8912
			for r, b := range file[39:45] {
				ay.WriteRegister(byte(r), b)
			}
			// ay.SetReg(file[38])
		}
	}

	return machine
}

func copyMemoryBlock(memOrg []byte, len uint16, memDest emulator.Memory, pos uint16) {
	posScr := uint16(0)
	posDst := pos
	// log.Printf("copying %d bytes to page 0x%04X\n", len, posDst)
	for posScr < len {
		if memOrg[posScr] == 0xED && memOrg[posScr+1] == 0xED && posScr+3 < len {
			b := memOrg[posScr+3]
			c := uint16(memOrg[posScr+2])
			for i := uint16(0); i < c; i++ {
				memDest.PutByte(posDst+i, b)
			}
			posDst += c
			posScr += 4
			// } else if memOrg[posScr] == 0xED {
			// 	memDest.PutByte(posDst, memOrg[posScr+1])
			// 	posDst++
			// 	posScr += 2
		} else {
			memDest.PutByte(posDst, memOrg[posScr])
			posDst++
			posScr++
		}
	}

}

func getUint16(h, l byte) uint16 {
	return (uint16(h) << 8) | uint16(l)
}
