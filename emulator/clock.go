package emulator

import (
	"fmt"
	rtdebug "runtime/debug"
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
	tStatesPerBlock uint
	blocks          uint
	tStates         uint
	tickers         []*ticker
	lastFrameTime   float64
	pasued          bool
}

func NewCLock(hz uint, blocks uint) Clock {
	clock := &clock{
		tStatesPerBlock: hz / blocks,
		wait:            time.Duration(time.Second / time.Duration(blocks)),
		blocks:          blocks,
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
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("PANIC!!! -> %v\n%T\n", r, r)
				fmt.Println("stacktrace from panic: \n" + string(rtdebug.Stack()))
			}
		}()

		for range ticker.C {
			if !c.pasued {
				start := time.Now()
				for (c.tStates < c.tStatesPerBlock) && !c.pasued {
					c.tick()
				}
				c.lastFrameTime = float64(time.Since(start).Microseconds()) / 1000.0
				c.tStates = 0

				if !c.pasued {
					hf = 1 - hf
				}
			}
		}
	}()
}

func (c *clock) RunFor(seconds uint) {
	for c.tStates < c.tStatesPerBlock*c.blocks*seconds {
		c.tick()
	}
}
