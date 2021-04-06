package gameboy

import (
	"fmt"
	"sync"

	"github.com/laullon/b2t80s/emulator"
)

var waveDuties = []float32{0.125, 0.25, 0.50, 0.75}

type channelRuntime struct {
	out          byte
	pulseLengths []uint
	ticksleft    uint
	buff         []*emulator.SoundData
	mux          sync.Mutex

	volumen      float64
	volumenDec   bool
	volumenStep  float64
	volumenTicks uint
}

type basicChannel struct {
	enable bool

	sampleRate uint

	initialVolume  byte
	envelopeInc    bool
	envelopeSweeps byte

	rt channelRuntime
}

type tomeChannel struct {
	basicChannel

	waveDuty    byte
	soundLength byte

	frequency7b uint16
	frequency   float32
	restart     bool
	consecutive bool

	sweepTime byte
	sweepDec  bool
	sweepNum  byte
}

type noiseChannel struct {
	basicChannel

	restart     bool
	consecutive bool
}

type apu struct {
	data []byte

	channel1 *tomeChannel
	channel2 *tomeChannel
	channel3 *basicChannel
	channel4 *noiseChannel

	outputTerminal byte
	outputControl  byte
}

func newAPU(sampleRate uint) *apu {
	apu := &apu{
		data:     make([]byte, 0x40),
		channel1: &tomeChannel{},
		channel2: &tomeChannel{},
		channel3: &basicChannel{},
		channel4: &noiseChannel{},
	}

	apu.channel1.rt.pulseLengths = make([]uint, 2)
	apu.channel2.rt.pulseLengths = make([]uint, 2)

	apu.channel1.sampleRate = sampleRate
	apu.channel2.sampleRate = sampleRate
	apu.channel3.sampleRate = sampleRate
	apu.channel4.sampleRate = sampleRate

	return apu
}

func (apu *apu) SoundTick() {
	apu.channel1.tick()
	apu.channel2.tick()
}

func (apu *apu) GetChannels() []emulator.SoundChannel {
	return []emulator.SoundChannel{
		apu.channel1,
		apu.channel2,
	}
}

func (apu *apu) ReadPort(addr uint16) (byte, bool) {
	fmt.Printf("[apu] read 0x%04X \n", addr)
	switch addr {
	default:
		// os.Exit(-1)
	}
	return apu.data[addr&0xff], false
}

func (apu *apu) WritePort(addr uint16, data byte) {
	apu.data[addr&0xff] = data
	fmt.Printf("[apu] write 0x%04X 0x%02X\n", addr, data)
	switch addr {
	case 0xff10:
		apu.channel1.sweepTime = (data >> 4) & 7
		apu.channel1.sweepDec = data&0b1000 != 0
		apu.channel1.sweepNum = data & 7

	case 0xff11:
		apu.channel1.waveDuty = data >> 6
		apu.channel1.soundLength = data & 63

	case 0xff12:
		apu.channel1.initialVolume = data >> 4
		apu.channel1.envelopeInc = data&0b1000 != 0
		apu.channel1.envelopeSweeps = data & 7

	case 0xff13:
		apu.channel1.frequency7b = apu.channel1.frequency7b&0xff00 | uint16(data)

	case 0xff14:
		apu.channel1.frequency7b = apu.channel1.frequency7b&0x00ff | uint16(data&0b111)<<8
		apu.channel1.restart = data&0x80 != 0
		apu.channel1.consecutive = data&0x40 != 0

	case 0xff17:
		apu.channel2.initialVolume = data >> 4
		apu.channel2.envelopeInc = data&0b1000 != 0
		apu.channel2.envelopeSweeps = data & 7

	case 0xff19:
		apu.channel2.frequency7b = apu.channel2.frequency7b&0x00ff | uint16(data&0b111)<<8
		apu.channel2.restart = data&0x80 != 0
		apu.channel2.consecutive = data&0x40 != 0

	case 0xff21:
		apu.channel4.initialVolume = data >> 4
		apu.channel4.envelopeInc = data&0b1000 != 0
		apu.channel4.envelopeSweeps = data & 7

	case 0xff23:
		apu.channel4.restart = data&0x80 != 0
		apu.channel4.consecutive = data&0x40 != 0

	case 0xff24:
		apu.outputControl = data

	case 0xff25:
		apu.outputTerminal = data

	case 0xff26:
		apu.channel1.enable = data&0x80 != 0
		apu.channel2.enable = data&0x80 != 0
		apu.channel3.enable = data&0x80 != 0
		apu.channel4.enable = data&0x80 != 0

	default:
		// os.Exit(-1)
	}

	if addr < 0xff15 {
		apu.channel1.update()
	} else if addr < 0xff1A {
		apu.channel2.update()
	}
}

func (ch *tomeChannel) update() {
	ch.frequency = 131072 / float32(2048-ch.frequency7b)
	fmt.Printf("-> f7b:%d f:%f (%f)(%f)\n", ch.frequency7b, ch.frequency, float32(ch.sampleRate)/ch.frequency, 1/ch.frequency)
	l := float32(ch.sampleRate) / ch.frequency
	ch.rt.pulseLengths[0] = uint(l * waveDuties[ch.waveDuty])
	ch.rt.pulseLengths[1] = uint(l) - ch.rt.pulseLengths[0]
	fmt.Printf("-> ch.initialVolume: %v ch.envelopeInc: %v ch.envelopeSweeps: %v\n", ch.initialVolume, ch.envelopeInc, ch.envelopeSweeps)

	ch.rt.volumen = float64(ch.initialVolume) / 15
	ch.rt.volumenDec = !ch.envelopeInc
	if ch.rt.volumenDec {
		ch.rt.volumenStep = ch.rt.volumen / float64(ch.envelopeSweeps)
	}
	ch.rt.volumenTicks = uint(ch.envelopeSweeps) * (ch.sampleRate / 64)
	fmt.Printf("-> ch.rt.volumen: %v\n", ch.rt.volumen)
}

func (ch *tomeChannel) tick() {
	ch.rt.mux.Lock()
	defer ch.rt.mux.Unlock()

	if ch.rt.volumen > 0 {
		if ch.rt.ticksleft == 0 {
			ch.rt.out = 1 - ch.rt.out
			ch.rt.ticksleft = ch.rt.pulseLengths[int(ch.rt.out)]
		}
		ch.rt.ticksleft--

		if ch.rt.volumenTicks > 0 && ch.rt.volumenTicks < 1 {
			if ch.rt.volumenDec {
				ch.rt.volumen -= ch.rt.volumenStep
			} else {
				ch.rt.volumen += ch.rt.volumenStep
			}
			ch.rt.volumenTicks = uint(ch.envelopeSweeps) * (ch.sampleRate / 64)
		}
		ch.rt.volumenTicks--
	}

	ch.rt.buff = append(ch.rt.buff, &emulator.SoundData{L: float64(ch.rt.out) * ch.rt.volumen, R: float64(ch.rt.out) * ch.rt.volumen})
}

func (ch *tomeChannel) GetBuffer(max int) (res []*emulator.SoundData, l int) {
	ch.rt.mux.Lock()
	defer ch.rt.mux.Unlock()

	if len(ch.rt.buff) > max {
		res = ch.rt.buff[:max]
		ch.rt.buff = ch.rt.buff[max:]
		l = max
	} else {
		res = ch.rt.buff
		ch.rt.buff = nil
		l = len(res)
	}
	return
}
