package atetris

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (t *atetris) OnKey(key sdl.Scancode) {
	switch key {

	case sdl.SCANCODE_1:
		t.pokey1.P0 = !t.pokey1.P0
	case sdl.SCANCODE_2:
		t.pokey1.P1 = !t.pokey1.P1

	case sdl.SCANCODE_SPACE:
		t.pokey2.P0 = !t.pokey2.P0
	case sdl.SCANCODE_DOWN:
		t.pokey2.P1 = !t.pokey2.P1
	case sdl.SCANCODE_RIGHT:
		t.pokey2.P2 = !t.pokey2.P2
	case sdl.SCANCODE_LEFT:
		t.pokey2.P3 = !t.pokey2.P3

	}
}
