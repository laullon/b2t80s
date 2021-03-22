package gameboy

import "github.com/laullon/b2t80s/cpu"

type timer struct {
	bus            cpu.Bus
	div            uint16
	tima, tma, tac byte
}

func newTimer(bus cpu.Bus) *timer {
	return &timer{bus: bus}
}

func (t *timer) Tick() {
	t.div++
	switch t.tac & 3 {
	case 0:
		if t.div&0x03ff == 0 {
			t.timaTick()
		}
	case 1:
		if t.div&0x000f == 0 {
			t.timaTick()
		}
	case 2:
		if t.div&0x003f == 0 {
			t.timaTick()
		}
	case 3:
		if t.div&0x00ff == 0 {
			t.timaTick()
		}
	}
}

func (t *timer) timaTick() {
	if t.tac&0b100 != 0 {
		t.tima++
		if t.tima == 0 {
			t.tima = t.tma
			t.bus.Write(0xff0f, 0b100)
		}
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
