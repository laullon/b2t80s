package zx

import (
	"fmt"
	"image"
	"time"

	"fyne.io/fyne"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
	"github.com/laullon/b2t80s/z80"
)

const (
	clock48k  = 3500000
	clock128k = 3546900
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
	cassette cassette.Cassette
	sound    emulator.SoundSystem
	debugger emulator.Debugger

	onEndFrame func()
}

func NewZX(cpu emulator.CPU, ula *ula, mem emulator.Memory, cassette cassette.Cassette, sound emulator.SoundSystem, onEndFrame func()) *zx {
	zx := &zx{
		ula:        ula,
		cpu:        cpu,
		mem:        mem,
		cassette:   cassette,
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

func (zx *zx) loadDataBlock() uint16 {
	data := zx.cassette.NextDataBlock()
	if data == nil {
		return emulator.CONTINUE
	}

	regs := zx.cpu.Registers().(*z80.Z80Registers)
	requestedLength := regs.DE.Get()
	startAddress := regs.IX.Get()
	// fmt.Printf("Loading block '%s' to 0x%04x (bl:0x%04x, l:0x%04x, bt:%d, a:%d)\n", block.Name(), startAddress, len(block.GetData()), requestedLength, block.Type(), regs._A)

	if regs.Aalt == data[0] {
		if regs.Falt.C {
			checksum := data[0]
			for i := uint16(0); i < requestedLength; i++ {
				loadedByte := data[i+1]
				zx.mem.PutByte(startAddress+i, loadedByte)
				checksum ^= loadedByte
			}
			regs.F.C = checksum == data[requestedLength+1]
		} else {
			regs.F.C = true
		}
		// log.Print("done")
	} else {
		regs.F.C = false
		// log.Print("BAD Block")
	}
	return 0x05e2
}
