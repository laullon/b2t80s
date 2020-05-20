package emulator

import (
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

type SoundSystem interface {
	Ticker
	AddSource(source SoundSource)
	SetVolume(float64)
}

type SoundSource interface {
	SoundTick()
	GetChannels() []SoundChannel
}

type SoundChannel interface {
	GetBuffer(int) ([]*SoundData, int)
}

type SoundData struct {
	L, R float64
}

type soundSystem struct {
	sources []SoundSource
	stream  *beep.Mixer
	volume  *effects.Volume
}

type channel struct {
	source SoundChannel
}

func NewSoundSystem(sampleRate int) SoundSystem {
	ss := &soundSystem{}

	sr := beep.SampleRate(sampleRate)
	speaker.Init(sr, sr.N(time.Second/30))
	ss.stream = &beep.Mixer{}

	ss.volume = &effects.Volume{
		Streamer: ss.stream,
		Base:     2,
		Volume:   -7,
	}

	speaker.Play(ss.volume)

	return ss
}

func (ss *soundSystem) AddSource(source SoundSource) {
	channels := source.GetChannels()
	for _, ch := range channels {
		ss.stream.Add(&channel{ch})
	}
	ss.sources = append(ss.sources, source)
}

func (ss *soundSystem) Tick() {
	for _, source := range ss.sources {
		source.SoundTick()
	}
}

func (ss *soundSystem) SetVolume(volume float64) {
	ss.volume.Volume = -3 - ((1 - volume) * 4)
	ss.volume.Silent = volume == 0
	println(volume, ss.volume.Volume)
}

func (ch *channel) Stream(samples [][2]float64) (int, bool) {
	buff, _ := ch.source.GetBuffer(len(samples))
	for i, data := range buff {
		samples[i][0] = data.L
		samples[i][1] = data.R
	}
	return len(samples), true
}

func (ss *channel) Err() error {
	return nil
}
