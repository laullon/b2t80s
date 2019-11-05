package emulator

type Clock interface {
	// AddTStates increment the tStates counter and return true if the frame is not done
	AddTStates(uint)
	FrameDone() bool
	ApplyDeplay()
	AddTicker(mod uint, t Ticker)
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
	tStatesDelay    []byte
	tickers         []*ticker
}

func NewCLock(hz uint) Clock {
	clock := &clock{
		tStatesPerFrame: hz / 50,
	}

	clock.tStatesDelay = make([]byte, clock.tStatesPerFrame+100)
	for y := 0; y < 192; y++ {
		for x := 0; x < 128; x += 8 {
			ts := (y+16+48)*224 + x + 16 + 48
			clock.tStatesDelay[ts] = 6
			clock.tStatesDelay[ts+1] = 5
			clock.tStatesDelay[ts+2] = 4
			clock.tStatesDelay[ts+3] = 3
			clock.tStatesDelay[ts+4] = 2
			clock.tStatesDelay[ts+5] = 1
		}
	}

	return clock
}

func (c *clock) AddTStates(ts uint) {
	for i := uint(0); i < ts; i++ {
		for _, t := range c.tickers {
			t.counter++
			if t.counter == t.mod || t.mod < 2 {
				t.counter = 0
				t.ticker.Tick()
			}
		}
	}
	c.tStates += ts
}

func (c *clock) FrameDone() bool {
	if c.tStates >= c.tStatesPerFrame {
		c.tStates -= c.tStatesPerFrame
		return true
	}
	return false
}

func (c *clock) ApplyDeplay() {
	c.AddTStates(uint(c.tStatesDelay[c.tStates]))
}

func (c *clock) AddTicker(mod uint, t Ticker) {
	c.tickers = append(c.tickers, &ticker{mod: mod, ticker: t})
}
