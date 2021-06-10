package a1942

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (m *a1942) OnKey(key sdl.Scancode) {
	switch key {
	case sdl.SCANCODE_P:
		m.p2 ^= 0x01
	case sdl.SCANCODE_I:
		m.p2 ^= 0x02
	case sdl.SCANCODE_0:
		m.p2 ^= 0x04
	case sdl.SCANCODE_O:
		m.p2 ^= 0x08
	case sdl.SCANCODE_Q: // A
		m.p2 ^= 0x10
	case sdl.SCANCODE_W: // B
		m.p2 ^= 0x20

	case sdl.SCANCODE_RIGHT:
		m.p1 ^= 0x01
	case sdl.SCANCODE_LEFT:
		m.p1 ^= 0x02
	case sdl.SCANCODE_DOWN:
		m.p1 ^= 0x04
	case sdl.SCANCODE_UP:
		m.p1 ^= 0x08
	case sdl.SCANCODE_Z: // A
		m.p1 ^= 0x10
	case sdl.SCANCODE_X: // B
		m.p1 ^= 0x20

	case sdl.SCANCODE_3: //IPT_START1
		m.sys ^= 0x01
	case sdl.SCANCODE_4: //IPT_START2
		m.sys ^= 0x02
	case sdl.SCANCODE_ESCAPE: //IPT_SERVICE1
		m.sys ^= 0x10
	case sdl.SCANCODE_1: //IPT_COIN2
		m.sys ^= 0x40
	case sdl.SCANCODE_2: //IPT_COIN1
		m.sys ^= 0x80
	}
}
