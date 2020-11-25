package emulator

import (
	"time"
)

type Clock interface {
	// AddTStates increment the tStates counter and return true if the frame is not done
	AddTicker(mod uint, t Ticker)
	Run()
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

func (c *clock) Run() {
	wait := time.Duration(20 * time.Millisecond)
	ticker := time.NewTicker(wait)
	go func() {
		for range ticker.C {
			for !c.frameDone() {
				c.tick()
			}
		}
	}()
}
