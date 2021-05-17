package zx

import (
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func (ula *ula) OnKey(key sdl.Scancode) {
	switch key {

	case sdl.SCANCODE_1:
		ula.keyboardRow[3] ^= 0b00000001
	case sdl.SCANCODE_2:
		ula.keyboardRow[3] ^= 0b00000010
	case sdl.SCANCODE_3:
		ula.keyboardRow[3] ^= 0b00000100
	case sdl.SCANCODE_4:
		ula.keyboardRow[3] ^= 0b00001000
	case sdl.SCANCODE_5:
		ula.keyboardRow[3] ^= 0b00010000

	case sdl.SCANCODE_0:
		ula.keyboardRow[4] ^= 0b00000001
	case sdl.SCANCODE_9:
		ula.keyboardRow[4] ^= 0b00000010
	case sdl.SCANCODE_8:
		ula.keyboardRow[4] ^= 0b00000100
	case sdl.SCANCODE_7:
		ula.keyboardRow[4] ^= 0b00001000
	case sdl.SCANCODE_6:
		ula.keyboardRow[4] ^= 0b00010000

	case sdl.SCANCODE_Q:
		ula.keyboardRow[2] ^= 0b00000001
	case sdl.SCANCODE_W:
		ula.keyboardRow[2] ^= 0b00000010
	case sdl.SCANCODE_E:
		ula.keyboardRow[2] ^= 0b00000100
	case sdl.SCANCODE_R:
		ula.keyboardRow[2] ^= 0b00001000
	case sdl.SCANCODE_T:
		ula.keyboardRow[2] ^= 0b00010000

	case sdl.SCANCODE_P:
		ula.keyboardRow[5] ^= 0b00000001
	case sdl.SCANCODE_O:
		ula.keyboardRow[5] ^= 0b00000010
	case sdl.SCANCODE_I:
		ula.keyboardRow[5] ^= 0b00000100
	case sdl.SCANCODE_U:
		ula.keyboardRow[5] ^= 0b00001000
	case sdl.SCANCODE_Y:
		ula.keyboardRow[5] ^= 0b00010000

	case sdl.SCANCODE_A:
		ula.keyboardRow[1] ^= 0b00000001
	case sdl.SCANCODE_S:
		ula.keyboardRow[1] ^= 0b00000010
	case sdl.SCANCODE_D:
		ula.keyboardRow[1] ^= 0b00000100
	case sdl.SCANCODE_F:
		ula.keyboardRow[1] ^= 0b00001000
	case sdl.SCANCODE_G:
		ula.keyboardRow[1] ^= 0b00010000

	case sdl.SCANCODE_RETURN:
		ula.keyboardRow[6] ^= 0b00000001
	case sdl.SCANCODE_L:
		ula.keyboardRow[6] ^= 0b00000010
	case sdl.SCANCODE_K:
		ula.keyboardRow[6] ^= 0b00000100
	case sdl.SCANCODE_J:
		ula.keyboardRow[6] ^= 0b00001000
	case sdl.SCANCODE_H:
		ula.keyboardRow[6] ^= 0b00010000

	case sdl.SCANCODE_LSHIFT, sdl.SCANCODE_RSHIFT:
		ula.keyboardRow[0] ^= 0b00000001
	case sdl.SCANCODE_Z:
		ula.keyboardRow[0] ^= 0b00000010
	case sdl.SCANCODE_X:
		ula.keyboardRow[0] ^= 0b00000100
	case sdl.SCANCODE_C:
		ula.keyboardRow[0] ^= 0b00001000
	case sdl.SCANCODE_V:
		ula.keyboardRow[0] ^= 0b00010000

	case sdl.SCANCODE_SPACE:
		ula.keyboardRow[7] ^= 0b00000001
	case sdl.SCANCODE_LCTRL, sdl.SCANCODE_RCTRL:
		ula.keyboardRow[7] ^= 0b00000010
	case sdl.SCANCODE_M:
		ula.keyboardRow[7] ^= 0b00000100
	case sdl.SCANCODE_N:
		ula.keyboardRow[7] ^= 0b00001000
	case sdl.SCANCODE_B:
		ula.keyboardRow[7] ^= 0b00010000

	case sdl.SCANCODE_BACKSPACE:
		ula.keyboardRow[0] ^= 0b00000001
		ula.keyboardRow[4] ^= 0b00000001

	case sdl.SCANCODE_UP:
		ula.keyboardRow[0] ^= 0b00000001
		ula.keyboardRow[4] ^= 0b00001000

	case sdl.SCANCODE_DOWN:
		ula.keyboardRow[0] ^= 0b00000001
		ula.keyboardRow[4] ^= 0b00010000
	}
}

var onlyOnce sync.Once

func (ula *ula) LoadCommand() {
	go onlyOnce.Do(func() {
		time.Sleep(time.Second)
		ula.OnKey(sdl.SCANCODE_J)
		time.Sleep(150 * time.Millisecond)
		ula.OnKey(sdl.SCANCODE_J)
		ula.OnKey(sdl.SCANCODE_LCTRL)
		time.Sleep(150 * time.Millisecond)
		ula.OnKey(sdl.SCANCODE_P)
		time.Sleep(150 * time.Millisecond)
		ula.OnKey(sdl.SCANCODE_P)
		time.Sleep(150 * time.Millisecond)
		ula.OnKey(sdl.SCANCODE_P)
		time.Sleep(150 * time.Millisecond)
		ula.OnKey(sdl.SCANCODE_P)
		ula.OnKey(sdl.SCANCODE_LCTRL)
		time.Sleep(150 * time.Millisecond)
		ula.OnKey(sdl.SCANCODE_RETURN)
		time.Sleep(150 * time.Millisecond)
		ula.OnKey(sdl.SCANCODE_RETURN)
	})
}

func (ula *ula) LoadCommand128() {
	go onlyOnce.Do(func() {
		time.Sleep(time.Second)
		ula.OnKey(sdl.SCANCODE_RETURN)
		time.Sleep(150 * time.Millisecond)
		ula.OnKey(sdl.SCANCODE_RETURN)
	})
}
