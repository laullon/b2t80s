package emulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestSound(t *testing.T) {
// 	c := &clock{
// 		tStatesPerFrame: 100,
// 	}

// 	ss := &dummySoundSorce{}
// 	s := &dummySound{ss: ss, t: t}
// 	c.AddTicker(0, ss)
// 	c.AddTicker(32, s)

// 	for {
// 		c.tick()
// 		if c.frameDone() {
// 			break
// 		}
// 	}
// }

type dummySound struct {
	ss *dummySoundSorce
	t  *testing.T
}

func (s *dummySound) Tick() {
	assert.Equal(s.t, 32, s.ss.ticks)
	s.ss.SoundTick()
}

type dummySoundSorce struct {
	ticks int
}

func (ss *dummySoundSorce) Tick() {
	ss.ticks++
}

func (ss *dummySoundSorce) SoundTick() {
	ss.ticks = 0
}
