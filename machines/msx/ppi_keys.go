package msx

import "github.com/veandco/go-sdl2/sdl"

func (ppi *ppi) OnKey(key sdl.Scancode) {
	switch key {

	case sdl.SCANCODE_0:
		ppi.keyboardRows[0] ^= 0b00000001
	case sdl.SCANCODE_1:
		ppi.keyboardRows[0] ^= 0b00000010
	case sdl.SCANCODE_2:
		ppi.keyboardRows[0] ^= 0b00000100
	case sdl.SCANCODE_3:
		ppi.keyboardRows[0] ^= 0b00001000
	case sdl.SCANCODE_4:
		ppi.keyboardRows[0] ^= 0b00010000
	case sdl.SCANCODE_5:
		ppi.keyboardRows[0] ^= 0b00100000
	case sdl.SCANCODE_6:
		ppi.keyboardRows[0] ^= 0b01000000
	case sdl.SCANCODE_7:
		ppi.keyboardRows[0] ^= 0b10000000

	case sdl.SCANCODE_8:
		ppi.keyboardRows[1] ^= 0b00000001
	case sdl.SCANCODE_9:
		ppi.keyboardRows[1] ^= 0b00000010
	// case sdl.SCANCODE_F7:
	// 	ppi.keyboardRows[1] ^= 0b00000100
	// case sdl.SCANCODE_F8:
	// 	ppi.keyboardRows[1] ^= 0b00001000
	// case sdl.SCANCODE_F5:
	// 	ppi.keyboardRows[1] ^= 0b00010000
	// case sdl.SCANCODE_F1:
	// 	ppi.keyboardRows[1] ^= 0b00100000
	// case sdl.SCANCODE_F2:
	// 	ppi.keyboardRows[1] ^= 0b01000000
	// case sdl.SCANCODE_F10:
	// 	ppi.keyboardRows[1] ^= 0b10000000

	// case sdl.SCANCODE_UP:
	// 	ppi.keyboardRows[2] ^= 0b00000001
	// case sdl.SCANCODE_LEFTBRACKET:
	// 	ppi.keyboardRows[2] ^= 0b00000010
	// case sdl.SCANCODE_RETURN:
	// 	ppi.keyboardRows[2] ^= 0b00000100
	// case sdl.SCANCODE_RIGHTBRACKET:
	// 	ppi.keyboardRows[2] ^= 0b00001000
	// case sdl.SCANCODE_F4:
	// 	ppi.keyboardRows[2] ^= 0b00010000
	// case "LeftShift", "RightShift":
	// 	ppi.keyboardRows[2] ^= 0b00100000
	case sdl.SCANCODE_A:
		ppi.keyboardRows[2] ^= 0b01000000
	case sdl.SCANCODE_B:
		ppi.keyboardRows[2] ^= 0b10000000

	case sdl.SCANCODE_C:
		ppi.keyboardRows[3] ^= 0b00000001
	case sdl.SCANCODE_D:
		ppi.keyboardRows[3] ^= 0b00000010
	case sdl.SCANCODE_E:
		ppi.keyboardRows[3] ^= 0b00000100
	case sdl.SCANCODE_F:
		ppi.keyboardRows[3] ^= 0b00001000
	case sdl.SCANCODE_G:
		ppi.keyboardRows[3] ^= 0b00010000
	case sdl.SCANCODE_H:
		ppi.keyboardRows[3] ^= 0b00100000
	case sdl.SCANCODE_I:
		ppi.keyboardRows[3] ^= 0b01000000
	case sdl.SCANCODE_J:
		ppi.keyboardRows[3] ^= 0b10000000

	case sdl.SCANCODE_K:
		ppi.keyboardRows[4] ^= 0b00000001
	case sdl.SCANCODE_L:
		ppi.keyboardRows[4] ^= 0b00000010
	case sdl.SCANCODE_M:
		ppi.keyboardRows[4] ^= 0b00000100
	case sdl.SCANCODE_N:
		ppi.keyboardRows[4] ^= 0b00001000
	case sdl.SCANCODE_O:
		ppi.keyboardRows[4] ^= 0b00010000
	case sdl.SCANCODE_P:
		ppi.keyboardRows[4] ^= 0b00100000
	case sdl.SCANCODE_Q:
		ppi.keyboardRows[4] ^= 0b01000000
	case sdl.SCANCODE_R:
		ppi.keyboardRows[4] ^= 0b10000000

	case sdl.SCANCODE_S:
		ppi.keyboardRows[5] ^= 0b00000001
	case sdl.SCANCODE_T:
		ppi.keyboardRows[5] ^= 0b00000010
	case sdl.SCANCODE_U:
		ppi.keyboardRows[5] ^= 0b00000100
	case sdl.SCANCODE_V:
		ppi.keyboardRows[5] ^= 0b00001000
	case sdl.SCANCODE_W:
		ppi.keyboardRows[5] ^= 0b00010000
	case sdl.SCANCODE_X:
		ppi.keyboardRows[5] ^= 0b00100000
	case sdl.SCANCODE_Y:
		ppi.keyboardRows[5] ^= 0b01000000
	case sdl.SCANCODE_Z:
		ppi.keyboardRows[5] ^= 0b10000000

	case sdl.SCANCODE_LSHIFT, sdl.SCANCODE_RSHIFT:
		ppi.keyboardRows[6] ^= 0b00000001
	case sdl.SCANCODE_LCTRL, sdl.SCANCODE_RCTRL:
		ppi.keyboardRows[6] ^= 0b00000010
	// case sdl.SCANCODE_R:
	// 	ppi.keyboardRows[6] ^= 0b00000100
	// case sdl.SCANCODE_T:
	// 	ppi.keyboardRows[6] ^= 0b00001000
	// case sdl.SCANCODE_G:
	// 	ppi.keyboardRows[6] ^= 0b00010000
	case sdl.SCANCODE_F1:
		ppi.keyboardRows[6] ^= 0b00100000
	case sdl.SCANCODE_F2:
		ppi.keyboardRows[6] ^= 0b01000000
	case sdl.SCANCODE_F3:
		ppi.keyboardRows[6] ^= 0b10000000

	case sdl.SCANCODE_F4:
		ppi.keyboardRows[7] ^= 0b00000001
	case sdl.SCANCODE_F5:
		ppi.keyboardRows[7] ^= 0b00000010
	case sdl.SCANCODE_ESCAPE:
		ppi.keyboardRows[7] ^= 0b00000100
		// case sdl.SCANCODE_F6:
		// ppi.keyboardRows[7] ^= 0b00001000
		// case sdl.SCANCODE_\U:
		// 	ppi.keyboardRows[7] ^= 0b00010000
		// case sdl.SCANCODE_D:
		// 	ppi.keyboardRows[7] ^= 0b00100000
		// case sdl.SCANCODE_C:
		// 	ppi.keyboardRows[7] ^= 0b01000000
	case sdl.SCANCODE_RETURN:
		ppi.keyboardRows[7] ^= 0b10000000

	case sdl.SCANCODE_SPACE:
		ppi.keyboardRows[8] ^= 0b00000001
	case sdl.SCANCODE_HOME:
		ppi.keyboardRows[8] ^= 0b00000010
	case sdl.SCANCODE_INSERT:
		ppi.keyboardRows[8] ^= 0b00000100
	case sdl.SCANCODE_DELETE:
		ppi.keyboardRows[8] ^= 0b00001000
	case sdl.SCANCODE_LEFT:
		ppi.keyboardRows[8] ^= 0b00010000
	case sdl.SCANCODE_UP:
		ppi.keyboardRows[8] ^= 0b00100000
	case sdl.SCANCODE_DOWN:
		ppi.keyboardRows[8] ^= 0b01000000
	case sdl.SCANCODE_RIGHT:
		ppi.keyboardRows[8] ^= 0b10000000

		// case sdl.SCANCODE_UP: JOY1
		// 	ppi.keyboardRows[9] ^= 0b00000001
		// case sdl.SCANCODE_RIGHT:
		// 	ppi.keyboardRows[9] ^= 0b00000010
		// case sdl.SCANCODE_DOWN:
		// 	ppi.keyboardRows[9] ^= 0b00000100
		// case sdl.SCANCODE_F9:
		// 	ppi.keyboardRows[9] ^= 0b00001000
		// case sdl.SCANCODE_F6:
		// 	ppi.keyboardRows[9] ^= 0b00010000
		// case sdl.SCANCODE_F3:
		// 	ppi.keyboardRows[9] ^= 0b00100000
		// case sdl.SCANCODE_ENTER:
		// 	ppi.keyboardRows[9] ^= 0b01000000
		// case sdl.SCANCODE_DELETE, "BACKSPACE":
		// 	ppi.keyboardRows[9] ^= 0b10000000

	}
}
