package emulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClock(t *testing.T) {
	c := &clock{
		tStatesPerFrame: 100,
	}
	for {
		c.AddTStates(1)
		if c.FrameDone() {
			break
		}
	}
	assert.Equal(t, 0, int(c.tStates))
}

func TestClockMods(t *testing.T) {
	var mod0, mod1, mod2, mod3, mod4, mod5, mod64 int

	c := NewCLock(3546900)

	c.AddTicker(0, &dummyTicker{counter: &mod0})
	c.AddTicker(1, &dummyTicker{counter: &mod1})
	c.AddTicker(2, &dummyTicker{counter: &mod2})
	c.AddTicker(3, &dummyTicker{counter: &mod3})
	c.AddTicker(4, &dummyTicker{counter: &mod4})
	c.AddTicker(5, &dummyTicker{counter: &mod5})
	c.AddTicker(64, &dummyTicker{counter: &mod64})

	for i := 0; i < 100; i++ {
		c.AddTStates(1)
	}

	assert.Equal(t, 100, mod0)
	assert.Equal(t, 100/1, mod1)
	assert.Equal(t, 100/2, mod2)
	assert.Equal(t, 100/3, mod3)
	assert.Equal(t, 100/4, mod4)
	assert.Equal(t, 100/5, mod5)
	assert.Equal(t, 100/64, mod64)
}

type dummyTicker struct {
	counter *int
}

func (t *dummyTicker) Tick() {
	*t.counter++
}
