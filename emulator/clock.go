package emulator

import (
	"fmt"
	"time"
)

type Clock interface {
	// AddTStates increment the tStates counter and return true if the frame is not done
	AddTicker(mod uint, t Ticker)
	Run()
	Stats() string
	Pause()
	Resume()
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
	tStatesPerFrame uint
	tStates         uint
	tickers         []*ticker
	lastFrameTime   float64
	pasued          bool
}

func NewCLock(hz int) Clock {
	clock := &clock{
		tStatesPerFrame: uint(hz) / 50,
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

func (c *clock) frameDone() bool {
	if c.tStates >= c.tStatesPerFrame {
		c.tStates -= c.tStatesPerFrame
		return true
	}
	return false
}

func (c *clock) AddTicker(mod uint, t Ticker) {
	if t == nil {
		panic("NIL Ticker")
	}
	c.tickers = append(c.tickers, &ticker{mod: mod, ticker: t})
}

func (c *clock) Stats() string {
	return fmt.Sprintf("%5.2fms", c.lastFrameTime)
}

func (c *clock) Run() {
	wait := time.Duration(20 * time.Millisecond)
	ticker := time.NewTicker(wait)
	go func() {
		for range ticker.C {
			start := time.Now()
			for !c.frameDone() && !c.pasued {
				c.tick()
			}
			c.lastFrameTime = float64(time.Now().Sub(start).Microseconds()) / 1000.0
		}
	}()
}
