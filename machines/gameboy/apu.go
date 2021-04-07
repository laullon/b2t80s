package gameboy

import (
	"fmt"
	"sync"

	"github.com/laullon/b2t80s/emulator"
)

var waveDuties = []float32{0.125, 0.25, 0.50, 0.75}

type channel interface {
	setRegister(r int, data byte)
	getRegister(r int) byte
	tick()
	lengthTick()
	setOn(bool)
	isOn() bool
}

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

	soundLength       uint16
	soundLengthEnable bool

	sampleRate uint

	rt channelRuntime
}

type tomeChannel struct {
	basicChannel

	waveDuty byte

	frequency11b uint16
	frequency    float32
	restart      bool

	initialVolume  byte
	envelopeInc    bool
	envelopeSweeps byte

	sweepDissabled bool
	sweepTime      byte
	sweepDec       bool
	sweepNum       byte
}

type waveChannel struct {
	basicChannel

	play   bool
	volume byte

	frequency11b uint16
	frequency    float32
	restart      bool

	ram []byte
}

type noiseChannel struct {
	basicChannel

	restart bool

	initialVolume  byte
	envelopeInc    bool
	envelopeSweeps byte

	polyClock     byte
	polyStep7bits bool
	polyRatio     byte
}

type apu struct {
	channels      []channel
	soundChannels []emulator.SoundChannel

	ticks          uint
	sampleRate     uint
	sequencerTicks uint
	sequencerCount uint

	outputTerminal byte
	outputControl  byte
	on             bool
}

func newAPU(sampleRate uint) *apu {
	apu := &apu{}

	apu.sampleRate = sampleRate
	apu.sequencerTicks = sampleRate / 512

	channel1 := &tomeChannel{}
	channel2 := &tomeChannel{}
	channel3 := &waveChannel{}
	channel4 := &noiseChannel{}

	channel1.rt.pulseLengths = make([]uint, 2)
	channel2.rt.pulseLengths = make([]uint, 2)
	channel2.sweepDissabled = true
	channel3.ram = make([]byte, 0x10)

	channel1.sampleRate = sampleRate
	channel2.sampleRate = sampleRate
	channel3.sampleRate = sampleRate
	channel4.sampleRate = sampleRate

	apu.channels = append(apu.channels, channel1, channel2, channel3, channel4)
	apu.soundChannels = append(apu.soundChannels, channel1, channel2, channel3, channel4)
	return apu
}

func (apu *apu) SoundTick() {
}

func (apu *apu) Tick() {
	apu.sequencerCount++
	apu.sequencerCount &= 7
	switch apu.sequencerCount {
	case 0, 4:
		for _, ch := range apu.channels {
			ch.lengthTick()
		}
	case 7:
		// volTick++
	case 2, 6:
		for _, ch := range apu.channels {
			ch.lengthTick()
		}
		// sweepTick++
	}
}

func (apu *apu) GetChannels() []emulator.SoundChannel {
	return apu.soundChannels
}

func decode(addr uint16) (channel, register int) {
	channel = int((addr&0xff - 0x10) / 5)
	register = int(addr&0xff-0x10) % 5
	return
}

func (apu *apu) getStatus() (res byte) {
	for c := 0; c < 4; c++ {
		if apu.channels[c].isOn() {
			res |= 1 << c
		}
	}
	return
}

func (apu *apu) ReadPort(addr uint16) (res byte, skip bool) {
	channel, register := decode(addr)
	if channel < 4 {
		res = apu.channels[channel].getRegister(register)
	} else if addr == 0xff24 {
		res = apu.outputControl
	} else if addr == 0xff25 {
		res = apu.outputTerminal
	} else if addr == 0xff26 {
		res = 0x70 | apu.getStatus()
		if apu.on {
			res |= 0x80
		}

	} else if addr >= 0xFF27 && addr <= 0xFF2F {
		res = 0xff
	} else {
		res = apu.channels[2].(*waveChannel).ram[addr&0x000f]
	}
	fmt.Printf("[apu] read  0x%04X 0x%02X (%d,%d)%04b\n", addr, res, channel, register, apu.getStatus()&0b1111)
	return
}

func (apu *apu) WritePort(addr uint16, data byte) {
	channel, register := decode(addr)
	fmt.Printf("[apu] write 0x%04X 0x%02X (%d,%d)%04b\n", addr, data, channel, register, apu.getStatus()&0b1111)
	if channel < 4 && apu.on {
		apu.channels[channel].setRegister(register, data)
	} else if addr == 0xff24 && apu.on {
		apu.outputControl = data
	} else if addr == 0xff25 && apu.on {
		apu.outputTerminal = data
	} else if addr == 0xff26 {
		apu.on = data&0x80 != 0
		if !apu.on {
			for c := 0; c < 4; c++ {
				for r := 0; r < 5; r++ {
					apu.channels[c].setRegister(r, 0)
				}
			}
			apu.outputControl = 0
			apu.outputTerminal = 0
		}
		for c := 0; c < 4; c++ {
			apu.channels[c].setOn(data&0x80 != 0)
		}
	} else if addr >= 0xFF27 && addr <= 0xFF2F {
	} else {
		apu.channels[2].(*waveChannel).ram[addr&0x000f] = data
	}
}

// 	case 0xff24:
// 		apu.outputControl = data

// 	case 0xff25:
// 		apu.outputTerminal = data

// 	case 0xff26:
// 		apu.channel1.enable = data&0x80 != 0
// 		apu.channel2.enable = data&0x80 != 0
// 		apu.channel3.enable = data&0x80 != 0
// 		apu.channel4.enable = data&0x80 != 0

// func (ch *tomeChannel) update() {
// 	ch.frequency = 131072 / float32(2048-ch.frequency7b)
// 	fmt.Printf("-> f7b:%d f:%f (%f)(%f)\n", ch.frequency7b, ch.frequency, float32(ch.sampleRate)/ch.frequency, 1/ch.frequency)
// 	l := float32(ch.sampleRate) / ch.frequency
// 	ch.rt.pulseLengths[0] = uint(l * waveDuties[ch.waveDuty])
// 	ch.rt.pulseLengths[1] = uint(l) - ch.rt.pulseLengths[0]
// 	fmt.Printf("-> ch.initialVolume: %v ch.envelopeInc: %v ch.envelopeSweeps: %v\n", ch.initialVolume, ch.envelopeInc, ch.envelopeSweeps)

// 	ch.rt.volumen = float64(ch.initialVolume) / 15
// 	ch.rt.volumenDec = !ch.envelopeInc
// 	if ch.rt.volumenDec {
// 		ch.rt.volumenStep = ch.rt.volumen / float64(ch.envelopeSweeps)
// 	}
// 	ch.rt.volumenTicks = uint(ch.envelopeSweeps) * (ch.sampleRate / 64)
// 	fmt.Printf("-> ch.rt.volumen: %v\n", ch.rt.volumen)
// }

// func (ch *tomeChannel) tick() {
// 	ch.rt.mux.Lock()
// 	defer ch.rt.mux.Unlock()

// 	if ch.rt.volumen > 0 {
// 		if ch.rt.ticksleft == 0 {
// 			ch.rt.out = 1 - ch.rt.out
// 			ch.rt.ticksleft = ch.rt.pulseLengths[int(ch.rt.out)]
// 		}
// 		ch.rt.ticksleft--

// 		if ch.rt.volumenTicks > 0 && ch.rt.volumenTicks < 1 {
// 			if ch.rt.volumenDec {
// 				ch.rt.volumen -= ch.rt.volumenStep
// 			} else {
// 				ch.rt.volumen += ch.rt.volumenStep
// 			}
// 			ch.rt.volumenTicks = uint(ch.envelopeSweeps) * (ch.sampleRate / 64)
// 		}
// 		ch.rt.volumenTicks--
// 	}

// 	ch.rt.buff = append(ch.rt.buff, &emulator.SoundData{L: float64(ch.rt.out) * ch.rt.volumen, R: float64(ch.rt.out) * ch.rt.volumen})
// }

// func (ch *tomeChannel) GetBuffer(max int) (res []*emulator.SoundData, l int) {
// 	ch.rt.mux.Lock()
// 	defer ch.rt.mux.Unlock()

// 	if len(ch.rt.buff) > max {
// 		res = ch.rt.buff[:max]
// 		ch.rt.buff = ch.rt.buff[max:]
// 		l = max
// 	} else {
// 		res = ch.rt.buff
// 		ch.rt.buff = nil
// 		l = len(res)
// 	}
// 	return
// }

func (ch *tomeChannel) setRegister(r int, data byte) {
	switch r {
	case 0:
		if !ch.sweepDissabled {
			ch.sweepTime = (data >> 4) & 7
			ch.sweepDec = data&0b1000 != 0
			ch.sweepNum = data & 7
		}

	case 1:
		ch.waveDuty = data >> 6
		ch.soundLength = 64 - uint16(data)&63

	case 2:
		ch.initialVolume = data >> 4
		ch.envelopeInc = data&0b1000 != 0
		ch.envelopeSweeps = data & 7
		if ch.initialVolume == 0 {
			ch.enable = false
		}

	case 3:
		ch.frequency11b = ch.frequency11b&0xff00 | uint16(data)

	case 4:
		ch.frequency11b = ch.frequency11b&0x00ff | uint16(data&0b111)<<8
		ch.restart = data&0x80 != 0
		ch.soundLengthEnable = data&0x40 != 0
		if ch.restart && ch.soundLength == 0 {
			ch.soundLength = 64
		}
		ch.enable = ch.soundLengthEnable
	}
}

func (ch *tomeChannel) getRegister(r int) (res byte) {
	switch r {
	case 0:
		if !ch.sweepDissabled {
			res = 0x80
			res |= ch.sweepTime & 7 << 4
			res |= ch.sweepNum & 7
			if ch.sweepDec {
				res |= 0b1000
			}
		} else {
			res = 0xff
		}

	case 1:
		res = 0x3f
		res |= ch.waveDuty << 6

	case 2:
		res = ch.initialVolume << 4
		if ch.envelopeInc {
			res |= 0b1000
		}
		res |= ch.envelopeSweeps

	case 4:
		res = 0xbf
		if ch.soundLengthEnable {
			res |= 0x40
		}

	case 3, 5:
		res = 0xff
	}

	return
}

func (ch *tomeChannel) tick()                                                {}
func (ch *tomeChannel) GetBuffer(max int) (res []*emulator.SoundData, l int) { return }

func (ch *waveChannel) setRegister(r int, data byte) {
	switch r {
	case 0:
		ch.enable = data&0x80 != 0

	case 1:
		ch.soundLength = 256 - uint16(data)

	case 2:
		ch.volume = data >> 5
		if ch.volume == 0 {
			ch.enable = false
		}

	case 3:
		ch.frequency11b = ch.frequency11b&0xff00 | uint16(data)

	case 4:
		ch.frequency11b = ch.frequency11b&0x00ff | uint16(data&0b111)<<8
		ch.restart = data&0x80 != 0
		ch.soundLengthEnable = data&0x40 != 0
		if ch.restart && ch.soundLength == 0 {
			ch.soundLength = 0x100
		}
		ch.enable = ch.soundLengthEnable
	}
}

func (ch *waveChannel) getRegister(r int) (res byte) {
	switch r {
	case 0:
		res = 0x7f
		if ch.play {
			res |= 0x80
		}

	case 1:
		res = 0xff

	case 2:
		res = 0x9F
		res |= ch.volume << 5

	case 3:
		res = 0xff

	case 4:
		res = 0xbf
		if ch.soundLengthEnable {
			res |= 0x40
		}
	}
	return
}

func (ch *waveChannel) tick()                                                {}
func (ch *waveChannel) GetBuffer(max int) (res []*emulator.SoundData, l int) { return }

func (ch *noiseChannel) setRegister(r int, data byte) {
	switch r {
	case 1:
		ch.soundLength = 64 - uint16(data)&63

	case 2:
		ch.initialVolume = data >> 4
		ch.envelopeInc = data&0b1000 != 0
		ch.envelopeSweeps = data & 7
		if ch.initialVolume == 0 {
			ch.enable = false
		}

	case 3:
		ch.polyClock = data >> 4
		ch.polyStep7bits = data&0b1000 != 0
		ch.polyRatio = data & 7

	case 4:
		ch.restart = data&0x80 != 0
		ch.soundLengthEnable = data&0x40 != 0
		if ch.restart && ch.soundLength == 0 {
			ch.soundLength = 64
		}
		ch.enable = ch.soundLengthEnable
	}
}

func (ch *noiseChannel) getRegister(r int) (res byte) {
	switch r {
	case 0, 1:
		res = 0xff

	case 2:
		res = ch.initialVolume << 4
		if ch.envelopeInc {
			res |= 0b1000
		}
		res |= ch.envelopeSweeps

	case 3:
		res = ch.polyClock << 4
		if ch.polyStep7bits {
			res |= 0b1000
		}
		res |= ch.polyRatio

	case 4:
		res = 0xbf
		if ch.soundLengthEnable {
			res |= 0x40
		}
	}
	return
}

func (ch *noiseChannel) tick()                                                {}
func (ch *noiseChannel) GetBuffer(max int) (res []*emulator.SoundData, l int) { return }

// ****************************
// ****************************
// ****************************

func (ch *basicChannel) setOn(on bool) { ch.enable = on }
func (ch *basicChannel) isOn() bool    { return ch.enable }

func (ch *basicChannel) lengthTick() {
	if ch.soundLengthEnable && ch.soundLength != 0 {
		ch.soundLength--
		if ch.soundLength == 0 {
			ch.enable = false
		}
		println("ch.soundLength:", ch.soundLength, ch.enable)
	}
}
