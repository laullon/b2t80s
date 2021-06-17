package gameboy

import (
	"sync"

	"github.com/laullon/b2t80s/emulator"
)

var waveDuties = []float32{0.125, 0.25, 0.50, 0.75}

type channel interface {
	setRegister(r int, data byte)
	getRegister(r int) byte

	tick()
	lengthTick()
	sweepTick()
	tickEnvelope()

	setOn(bool)
	isOn() bool
}

type channelRuntime struct {
	out          byte
	pulseLengths []uint
	ticksleft    uint
	buff         []*emulator.SoundData
	mux          sync.Mutex
}

type envelope struct {
	initialVolume byte
	volume        byte
	inc           bool
	period        byte
	timer         byte
}
type basicChannel struct {
	enable bool
	dac    bool

	soundLength       uint16
	soundLengthEnable bool

	sampleRate uint

	rt channelRuntime
}

type tomeChannel struct {
	basicChannel

	envelope envelope

	waveDuty byte

	frequency11b uint16
	frequency    uint16
	trigger      bool

	sweepOFF     bool
	sweepEnable  bool
	sweepPeriod  byte
	sweepPeriodr byte
	sweepDec     bool
	sweepShift   byte
}

type waveChannel struct {
	basicChannel

	volume byte

	frequency11b uint16
	frequency    float32
	trigger      bool

	ram []byte
}

type noiseChannel struct {
	basicChannel

	envelope envelope

	trigger bool

	initialVolume  byte
	envelopeInc    bool
	envelopePeriod byte

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
	channel2.sweepOFF = true
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
	for _, ch := range apu.channels {
		ch.tick()
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

func (apu *apu) ReadPort(addr uint16) byte {
	var res byte
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
	// fmt.Printf("[apu] read  0x%04X 0x%02X (%d,%d)%04b\n", addr, res, channel, register, apu.getStatus()&0b1111)
	return res
}

func (apu *apu) WritePort(addr uint16, data byte) {
	channel, register := decode(addr)
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
	// fmt.Printf("[apu] write 0x%04X 0x%02X (%d,%d)%04b\n", addr, data, channel, register, apu.getStatus()&0b1111)
}

func (ch *tomeChannel) update() {
	freq := 131072 / float32(2048-ch.frequency)
	l := float32(ch.sampleRate) / freq
	ch.rt.pulseLengths[0] = uint(l * waveDuties[ch.waveDuty])
	ch.rt.pulseLengths[1] = uint(l) - ch.rt.pulseLengths[0]
}

func (ch *tomeChannel) _tick() {
	ch.rt.mux.Lock()
	defer ch.rt.mux.Unlock()

	if !ch.enable {
		ch.rt.buff = append(ch.rt.buff, &emulator.SoundData{L: 0, R: 0})
		return
	}

	ch.update()

	if ch.rt.ticksleft == 0 {
		ch.rt.out = 1 - ch.rt.out
		ch.rt.ticksleft = ch.rt.pulseLengths[int(ch.rt.out)]
	}
	ch.rt.ticksleft--

	ch.rt.buff = append(ch.rt.buff, &emulator.SoundData{L: float64(ch.rt.out) * float64(ch.envelope.volume), R: float64(ch.rt.out) * float64(ch.envelope.volume)})
}

func (ch *tomeChannel) calculateSweep() uint16 {
	r := ch.frequency >> ch.sweepShift
	if ch.sweepDec {
		r = ch.frequency - r
	} else {
		r = ch.frequency + r
	}

	if r > 2047 {
		ch.enable = false
	}

	return r
}

// ****************************
// ****************************
// ****************************

func (ch *noiseChannel) setRegister(r int, data byte) {
	switch r {
	case 1:
		ch.soundLength = 64 - uint16(data)&63

	case 2:
		ch.initialVolume = data >> 4
		ch.envelopeInc = data&0b1000 != 0
		ch.envelopePeriod = data & 7
		if ch.initialVolume == 0 {
			ch.enable = false
		}

	case 3:
		ch.polyClock = data >> 4
		ch.polyStep7bits = data&0b1000 != 0
		ch.polyRatio = data & 7

	case 4:
		ch.trigger = data&0x80 != 0
		ch.soundLengthEnable = data&0x40 != 0
		if ch.trigger && ch.soundLength == 0 {
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
		res |= ch.envelopePeriod

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

// ****************************
// ****************************
// ****************************

func (ch *basicChannel) setOn(on bool) { ch.enable = on }
func (ch *basicChannel) isOn() bool    { return ch.enable }

func (ch *basicChannel) GetBuffer(max int) (res []*emulator.SoundData, l int) {
	ch.rt.mux.Lock()
	defer ch.rt.mux.Unlock()

	// println("len(ch.rt.buff)", len(ch.rt.buff), "max", max)
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

func (ch *basicChannel) tick() {
	ch.rt.mux.Lock()
	defer ch.rt.mux.Unlock()
	ch.rt.buff = append(ch.rt.buff, &emulator.SoundData{L: 0, R: 0})
}
