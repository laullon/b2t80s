package emulator

import (
	"fmt"
	"time"
)

type Clock interface {
	// AddTStates increment the tStates counter and return true if the frame is not done
	AddTicker(mod uint, t Ticker)
	Run()
	RunFor(seconds uint)
	Stats() string
	Pause()
	Resume()
	SetOnFrameCallback(func())
}

type Ticker interface {
	Tick()
}

type ticker struct {
	counter uint
	mod     uint
	ticker  Ticker
}

type clock struct {
	wait            time.Duration
	tStatesPerFrame uint
	tStates         uint
	tickers         []*ticker
	lastFrameTime   float64
	pasued          bool
	callback        func()
}

func NewCLock(hz uint, fps uint) Clock {
	clock := &clock{
		tStatesPerFrame: hz / (fps * 2),
		wait:            time.Duration(time.Second / time.Duration(fps*2)),
	}
	return clock
}

func (c *clock) tick() {
	c.tStates++
	for _, t := range c.tickers {
		t.counter++
		if t.counter == t.mod || t.mod < 2 {
			t.counter = 0
			t.ticker.Tick()
		}
	}
}

func (c *clock) Pause() {
	c.pasued = true
}

func (c *clock) Resume() {
	c.pasued = false
}

func (c *clock) AddTicker(mod uint, t Ticker) {
	if t == nil {
		panic("NIL Ticker")
	}
	c.tickers = append(c.tickers, &ticker{mod: mod, ticker: t})
}

func (c *clock) Stats() string {
	return fmt.Sprintf("%5.2fms (%d)(%.2f)", c.lastFrameTime, c.wait.Milliseconds(), c.lastFrameTime/float64(c.wait.Milliseconds()))
}

func (c *clock) Run() {
	ticker := time.NewTicker(c.wait)
	hf := 0
	go func() {
		for range ticker.C {
			if !c.pasued {
				start := time.Now()
				for (c.tStates < c.tStatesPerFrame) && !c.pasued {
					c.tick()
				}
				c.lastFrameTime = float64(time.Since(start).Microseconds()) / 1000.0
				c.tStates = 0

				if !c.pasued {
					hf = 1 - hf
					if hf == 0 && c.callback != nil {
						c.callback()
					}
				} else {
					c.callback()
				}
			} else {
				c.callback()
			}
		}
	}()
}

func (c *clock) RunFor(seconds uint) {
	for c.tStates < (c.tStatesPerFrame * 50 * seconds) {
		c.tick()
	}
}

func (c *clock) SetOnFrameCallback(f func()) {
	c.callback = f
}
