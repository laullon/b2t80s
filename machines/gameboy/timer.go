package gameboy

import "github.com/laullon/b2t80s/cpu"

type timer struct {
	bus            cpu.Bus
	div            uint16
	tima, tma, tac byte
}

var divMods = []uint16{1024, 16, 64, 256}

func newTimer(bus cpu.Bus) *timer {
	return &timer{bus: bus}
}

func (t *timer) Tick() {
	t.div++

	if t.tac&0b100 != 0 {
		if t.div%divMods[t.tac&3] == 0 {
			t.tima++
			if t.tima == 0 {
				t.tima = t.tma
				t.bus.Write(0xff0f, 0b100)
			}
		}
		// println("t.div:", t.div, "(", t.div%divMods[t.tac&3], ")", "t.tima:", t.tima)
	}
}

func (t *timer) WritePort(addr uint16, data byte) {
	switch addr {
	case 0xff04:
		t.div = uint16(data) << 8
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
