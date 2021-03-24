package gameboy

import "github.com/laullon/b2t80s/cpu"

type timer struct {
	bus            cpu.Bus
	div            uint16
	tima, tma, tac byte
	overflow       bool
}

var divMods = []uint16{0x03ff, 0x000f, 0x003f, 0x00ff}

func newTimer(bus cpu.Bus) *timer {
	return &timer{bus: bus}
}

func (t *timer) Tick() {
	t.div++

	if t.div&divMods[t.tac&3] == 0 {
		if t.tac&0b100 != 0 {
			t.tima++
		}
	}

	if t.tima == 0 {
		t.tima = t.tma
		if !t.overflow {
			t.bus.Write(0xff0f, 0b100)
			t.overflow = true
		}
	}

}

func (t *timer) timaTick() {
	if t.tac&0b100 != 0 {
		t.tima++
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
