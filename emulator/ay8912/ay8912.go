package ay8912

import (
	"math/rand"
	"sync"

	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/emulator"
)

type AY8912 interface {
	cpu.PortManager
	emulator.Ticker
	emulator.SoundSource
	ReadRegister(reg byte) byte
	WriteRegister(reg byte, data byte)
}

type channel struct {
	volume   uint16
	pitch    uint16
	envelope bool
	tone     bool
	noise    bool
	count    uint16
	output   bool
	out      uint16
	buff     []*emulator.SoundData
	mux      sync.Mutex
	ticks    int
}

type envelope struct {
	pitch  uint16
	count  uint16
	pos    byte
	shape  byte
	volume [][]uint16
}

type noise struct {
	pitch  byte
	count  byte
	output bool
}

type ay8912 struct {
	regs        []byte
	selectedReg byte

	channels      []*channel
	soundChannels []emulator.SoundChannel
	envelope      *envelope
	noise         *noise
}

var regMasks = []byte{0xff, 0x0f, 0xff, 0x0f, 0xff, 0x0f, 0x1f, 0xff, 0x1f, 0x1f, 0x1f, 0xff, 0xff, 0x0f, 0xff, 0xff}

func New() AY8912 {
	ay := &ay8912{
		regs:     make([]byte, 16),
		envelope: &envelope{},
		noise:    &noise{},
	}

	ay.regs[14] = 0xff
	ay.regs[15] = 0xff

	for i := 0; i < 3; i++ {
		ch := &channel{}
		ay.channels = append(ay.channels, ch)
		ay.soundChannels = append(ay.soundChannels, ch)
	}

	for env := 0; env < 16; env++ {
		ay.envelope.volume = append(ay.envelope.volume, make([]uint16, 128))
	}

	dir := 0
	vol := 0
	for env := 0; env < 16; env++ {
		hold := false
		if (env & 4) != 0 {
			dir = 1
			vol = -1
		} else {
			dir = -1
			vol = 32
		}
		for pos := 0; pos < 64; pos++ {
			if !hold {
				vol += dir
				if vol < 0 || vol >= 16 {
					if env&8 != 0 {
						if env&2 != 0 {
							dir = -dir
						}
						if dir > 0 {
							vol = 0
						} else {
							vol = 15
						}

						if env&1 != 0 {
							hold = true
							if dir > 0 {
								vol = 15
							} else {
								vol = 0
							}
						}
					} else {
						vol = 0
						hold = true
					}
				}
			}
			ay.envelope.volume[env][pos] = uint16(vol)
		}
	}

	return ay
}

// TODO: remove... create a wrapper for each machine
func (ay *ay8912) ReadPort(port uint16) byte {
	return ay.regs[ay.selectedReg]
}

func (ay *ay8912) WritePort(port uint16, data byte) {
	if port == 0xfffd {
		ay.selectedReg = data
	} else {
		ay.WriteRegister(ay.selectedReg, data)
	}
}

func (ay *ay8912) ReadRegister(reg byte) byte {
	ay.selectedReg = reg
	return ay.regs[ay.selectedReg]
}

func (ay *ay8912) WriteRegister(reg byte, data byte) {
	// fmt.Printf("[ay8912] reg:%d data:%d\n", reg, data)
	ay.selectedReg = reg & 0x0f
	ay.regs[ay.selectedReg] = data & regMasks[ay.selectedReg]
	ay.update(ay.selectedReg)
}

func (ay *ay8912) Tick() {
	for _, channel := range ay.channels {
		channel.count++
		if channel.count >= channel.pitch {
			channel.count = 0
			channel.output = !channel.output
		}
	}

	ay.noise.count++
	if ay.noise.count >= ay.noise.pitch {
		ay.noise.count = 0
		ay.noise.output = rand.Float64() > 0.5
	}

	ay.envelope.count++
	if ay.envelope.count >= ay.envelope.pitch {
		ay.envelope.count = 0
		ay.envelope.pos++
		if ay.envelope.pos > 64 {
			ay.envelope.pos = 32
		}
	}

	for _, channel := range ay.channels {
		if (channel.output && channel.tone) || (ay.noise.output && channel.noise) {
			if channel.envelope {
				channel.out += ay.envelope.volume[ay.envelope.shape][ay.envelope.pos]
			} else {
				channel.out += channel.volume
			}
		}
		channel.ticks++
	}
}

func (ay *ay8912) SoundTick() {
	for _, channel := range ay.channels {
		channel.soundTick()
	}
}

func (ay *ay8912) GetChannels() []emulator.SoundChannel {
	return ay.soundChannels
}

func (ch *channel) soundTick() {
	ch.mux.Lock()
	defer ch.mux.Unlock()
	v := float64(ch.out) / float64(16*ch.ticks)
	ch.buff = append(ch.buff, &emulator.SoundData{L: v, R: v})
	ch.ticks = 0
	ch.out = 0
}

func (ch *channel) GetBuffer(max int) (res []*emulator.SoundData, l int) {
	ch.mux.Lock()
	defer ch.mux.Unlock()

	if len(ch.buff) > max {
		res = ch.buff[:max]
		ch.buff = ch.buff[max:]
		l = max
	} else {
		res = ch.buff
		ch.buff = nil
		l = len(res)
	}
	return
}

func (ay *ay8912) update(reg byte) {
	switch reg {

	case 0, 1, 2, 3, 4, 5:
		ch := reg / 2
		chr := ch * 2
		ay.channels[ch].pitch = uint16(ay.regs[chr]) | (uint16(ay.regs[chr+1]) << 8)
		ay.channels[ch].pitch *= 16

	case 6:
		ay.noise.pitch = ay.regs[6]

	case 7:
		val := ay.regs[7]
		ay.channels[0].tone = val&0b00000001 == 0
		ay.channels[1].tone = val&0b00000010 == 0
		ay.channels[2].tone = val&0b00000100 == 0
		ay.channels[0].noise = val&0b00001000 == 0
		ay.channels[1].noise = val&0b00010000 == 0
		ay.channels[2].noise = val&0b00100000 == 0

	case 8, 9, 10:
		ch := reg - 8
		val := ay.regs[reg]
		ay.channels[ch].envelope = val&0x10 != 0
		if !ay.channels[ch].envelope {
			ay.channels[ch].volume = uint16(val & 0x0f)
		}

	case 11, 12:
		ay.envelope.pitch = uint16(ay.regs[11]) & (uint16(ay.regs[12]) << 8)

	case 13:
		ay.envelope.shape = ay.regs[13]
		ay.envelope.count = 0
		ay.envelope.pos = 0
	}

	// fmt.Printf("[ay8912] regs -> %+v\n", ay.regs)
	// for i, ch := range ay.channels {
	// 	fmt.Printf("%d -> %+v\n", i, ch)
	// }
}
