package gameboy

import (
	"testing"

	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/emulator"
	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {

	clock := emulator.NewCLock(4_194_304, 64)

	timer0 := newTimer(&timerBus{})
	timer0.tac = 0b100

	timer1 := newTimer(&timerBus{})
	timer1.tac = 0b101

	timer2 := newTimer(&timerBus{})
	timer2.tac = 0b110

	timer3 := newTimer(&timerBus{})
	timer3.tac = 0b111

	cpu := &timerTicker{}
	clock.AddTicker(0, cpu)
	clock.AddTicker(0, timer0)
	clock.AddTicker(0, timer1)
	clock.AddTicker(0, timer2)
	clock.AddTicker(0, timer3)
	clock.RunFor(1)

	assert.Equal(t, 0, cpu.count)
	assert.Equal(t, 16, timer0.bus.(*timerBus).count)
	assert.Equal(t, 1024, timer1.bus.(*timerBus).count)
	assert.Equal(t, 256, timer2.bus.(*timerBus).count)
	assert.Equal(t, 64, timer3.bus.(*timerBus).count)
}

type timerBus struct {
	count int
}

func (bus *timerBus) Write(addr uint16, data uint8)                                        { bus.count++ }
func (bus *timerBus) Read(addr uint16) uint8                                               { panic(-1) }
func (bus *timerBus) RegisterPort(name string, mask cpu.PortMask, manager cpu.PortManager) { panic(-1) }
func (bus *timerBus) DumpMap() string                                                      { panic(-1) }
func (bus *timerBus) GetDumplables() map[string]cpu.Dumpable                               { panic(-1) }

type timerTicker struct {
	count int
}

func (t *timerTicker) Tick() {
	t.count++
}
