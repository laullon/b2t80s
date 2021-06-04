package emulator

// typedef unsigned char Uint8;
// void SineWave(void *userdata, Uint8 *stream, int len);
import "C"

import (
	"reflect"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
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
}

type channel struct {
	source SoundChannel
}

var ss *soundSystem

func init() {
	ss = &soundSystem{}
}

func NewSoundSystem(sampleRate uint) SoundSystem {

	spec := &sdl.AudioSpec{
		Freq:     int32(sampleRate),
		Format:   sdl.AUDIO_U8,
		Channels: 2,
		Samples:  256,
		Callback: sdl.AudioCallback(C.SineWave),
	}
	if err := sdl.OpenAudio(spec, nil); err != nil {
		panic(-1)
	}
	sdl.PauseAudio(false)

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
}

//export SineWave
func SineWave(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	dst := *(*[]uint8)(unsafe.Pointer(&hdr))
	for i := range dst {
		dst[i] = 0
	}

	src := make([]uint8, n)
	for _, source := range ss.sources {
		for _, ch := range source.GetChannels() {
			buffer, _ := ch.GetBuffer(n / 2)
			for i, data := range buffer {
				src[i*2] = uint8(data.L*127 + 128)
				src[i*2+1] = uint8(data.R*127 + 128)
			}
			sdl.MixAudioFormat(&dst[0], &src[0], sdl.AUDIO_U8, uint32(n), 50)
			// fmt.Printf("src -> %v\n", src[:20])
			// fmt.Printf("dst -> %v\n", dst[:20])
		}
	}
	// println("----")
}

func (ss *channel) Err() error {
	return nil
}
