package emulator

type Clock interface {
	// AddTStates increment the tStates counter and return true if the frame is not done
	AddTStates(uint)
	FrameDone() bool
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
	tickers         []*ticker
}

func NewCLock(hz uint) Clock {
	clock := &clock{
		tStatesPerFrame: hz / 50,
	}
	return clock
}

func (c *clock) AddTStates(ts uint) {
	for i := uint(0); i < ts; i++ {
		c.tStates++
		for _, t := range c.tickers {
			t.counter++
			if t.counter == t.mod || t.mod < 2 {
				t.counter = 0
				t.ticker.Tick()
			}
		}
	}
}

func (c *clock) FrameDone() bool {
	if c.tStates >= c.tStatesPerFrame {
		c.tStates -= c.tStatesPerFrame
		return true
	}
	return false
}

func (c *clock) AddTicker(mod uint, t Ticker) {
	c.tickers = append(c.tickers, &ticker{mod: mod, ticker: t})
}
