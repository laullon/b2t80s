package emulator

import (
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

type SoundSystem interface {
	Ticker
	AddSource(source SoundSource)
	SetVolume(float64)
}

type SoundSource interface {
	SoundTick()
	GetBuffer(int) ([]*SoundData, int)
}

type SoundData struct {
	L, R float64
}

type soundSystem struct {
	sources []SoundSource
	volume  float64
}

func NewSoundSystem(sampleRate int) SoundSystem {
	ss := &soundSystem{volume: 100}

	sr := beep.SampleRate(sampleRate)
	speaker.Init(sr, sr.N(time.Second/100))
	speaker.Play(ss)

	return ss
}

func (ss *soundSystem) AddSource(source SoundSource) {
	ss.sources = append(ss.sources, source)
}

func (ss *soundSystem) Tick() {
	for _, source := range ss.sources {
		source.SoundTick()
	}
}

func (ss *soundSystem) SetVolume(volume float64) {
	ss.volume = 100 - volume
}

func (ss *soundSystem) Stream(samples [][2]float64) (int, bool) {
	n := 0
	for _, source := range ss.sources {
		buff, len := source.GetBuffer(len(samples))
		n = max(n, len)
		for i, data := range buff {
			samples[i][0] += data.L
			samples[i][1] += data.R
		}
	}

	if n == 0 {
		n = len(samples)
	}

	for i := 0; i < n; i++ {
		samples[i][0] /= float64(len(ss.sources))
		samples[i][1] /= float64(len(ss.sources))

		samples[i][0] /= ss.volume
		samples[i][1] /= ss.volume
	}

	return n, true
}

func (ss *soundSystem) Err() error {
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
