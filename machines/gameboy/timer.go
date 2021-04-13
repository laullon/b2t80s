package gameboy

import "github.com/laullon/b2t80s/cpu"

type timer struct {
	bus            cpu.Bus
	div            uint16
	tima, tma, tac byte
}

var divMods = []uint16{
	0b0000_0010_0000_0000,
	0b0000_0000_0000_1000,
	0b0000_0000_0010_0000,
	0b0000_0000_1000_0000,
}

func newTimer(bus cpu.Bus) *timer {
	return &timer{bus: bus}
}

func (t *timer) Tick() {
	t.setDiv(t.div + 1)
}

func (t *timer) setDiv(newDiv uint16) {
	prev := t.div&divMods[t.tac&3] != 0
	t.div = newDiv
	new := t.div&divMods[t.tac&3] != 0

	if t.tac&0b100 != 0 {
		if prev && !new { // 1 -> 0
			if t.tima == 0xff {
				t.tima = t.tma
				t.bus.Write(0xff0f, 0b100)
			} else {
				t.tima++
			}
		}
		// println("t.div:", t.div, "(", t.div%divMods[t.tac&3], ")", "t.tima:", t.tima)
	}
}

func (t *timer) WritePort(addr uint16, data byte) {
	switch addr {
	case 0xff04:
		t.setDiv(uint16(data) << 8)
	case 0xff05:
		t.tima = data
	case 0xff06:
		t.tma = data
	case 0xff07:
		t.tac = data
	default:
		panic(-1)
	}
}

func (t *timer) ReadPort(addr uint16) (byte, bool) {
	switch addr {
	case 0xff04:
		return uint8(t.div >> 8), false
	case 0xff05:
		return t.tima, false
	case 0xff06:
		return t.tma, false
	case 0xff07:
		return t.tac, false
	}
	panic(-1)

}
