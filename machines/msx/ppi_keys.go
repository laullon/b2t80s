package msx

import "github.com/go-gl/glfw/v3.3/glfw"

func (ppi *ppi) OnKey(key glfw.Key) {
	// println("key:", key.Name)

	switch key {

	case glfw.Key0:
		ppi.keyboardRows[0] ^= 0b00000001
	case glfw.Key1:
		ppi.keyboardRows[0] ^= 0b00000010
	case glfw.Key2:
		ppi.keyboardRows[0] ^= 0b00000100
	case glfw.Key3:
		ppi.keyboardRows[0] ^= 0b00001000
	case glfw.Key4:
		ppi.keyboardRows[0] ^= 0b00010000
	case glfw.Key5:
		ppi.keyboardRows[0] ^= 0b00100000
	case glfw.Key6:
		ppi.keyboardRows[0] ^= 0b01000000
	case glfw.Key7:
		ppi.keyboardRows[0] ^= 0b10000000

	case glfw.Key8:
		ppi.keyboardRows[1] ^= 0b00000001
	case glfw.Key9:
		ppi.keyboardRows[1] ^= 0b00000010
	// case glfw.KeyF7:
	// 	ppi.keyboardRows[1] ^= 0b00000100
	// case glfw.KeyF8:
	// 	ppi.keyboardRows[1] ^= 0b00001000
	// case glfw.KeyF5:
	// 	ppi.keyboardRows[1] ^= 0b00010000
	// case glfw.KeyF1:
	// 	ppi.keyboardRows[1] ^= 0b00100000
	// case glfw.KeyF2:
	// 	ppi.keyboardRows[1] ^= 0b01000000
	// case glfw.KeyF10:
	// 	ppi.keyboardRows[1] ^= 0b10000000

	// case glfw.KeyUp:
	// 	ppi.keyboardRows[2] ^= 0b00000001
	// case glfw.KeyLeftBracket:
	// 	ppi.keyboardRows[2] ^= 0b00000010
	// case glfw.KeyReturn:
	// 	ppi.keyboardRows[2] ^= 0b00000100
	// case glfw.KeyRightBracket:
	// 	ppi.keyboardRows[2] ^= 0b00001000
	// case glfw.KeyF4:
	// 	ppi.keyboardRows[2] ^= 0b00010000
	// case "LeftShift", "RightShift":
	// 	ppi.keyboardRows[2] ^= 0b00100000
	case glfw.KeyA:
		ppi.keyboardRows[2] ^= 0b01000000
	case glfw.KeyB:
		ppi.keyboardRows[2] ^= 0b10000000

	case glfw.KeyC:
		ppi.keyboardRows[3] ^= 0b00000001
	case glfw.KeyD:
		ppi.keyboardRows[3] ^= 0b00000010
	case glfw.KeyE:
		ppi.keyboardRows[3] ^= 0b00000100
	case glfw.KeyF:
		ppi.keyboardRows[3] ^= 0b00001000
	case glfw.KeyG:
		ppi.keyboardRows[3] ^= 0b00010000
	case glfw.KeyH:
		ppi.keyboardRows[3] ^= 0b00100000
	case glfw.KeyI:
		ppi.keyboardRows[3] ^= 0b01000000
	case glfw.KeyJ:
		ppi.keyboardRows[3] ^= 0b10000000

	case glfw.KeyK:
		ppi.keyboardRows[4] ^= 0b00000001
	case glfw.KeyL:
		ppi.keyboardRows[4] ^= 0b00000010
	case glfw.KeyM:
		ppi.keyboardRows[4] ^= 0b00000100
	case glfw.KeyN:
		ppi.keyboardRows[4] ^= 0b00001000
	case glfw.KeyO:
		ppi.keyboardRows[4] ^= 0b00010000
	case glfw.KeyP:
		ppi.keyboardRows[4] ^= 0b00100000
	case glfw.KeyQ:
		ppi.keyboardRows[4] ^= 0b01000000
	case glfw.KeyR:
		ppi.keyboardRows[4] ^= 0b10000000

	case glfw.KeyS:
		ppi.keyboardRows[5] ^= 0b00000001
	case glfw.KeyT:
		ppi.keyboardRows[5] ^= 0b00000010
	case glfw.KeyU:
		ppi.keyboardRows[5] ^= 0b00000100
	case glfw.KeyV:
		ppi.keyboardRows[5] ^= 0b00001000
	case glfw.KeyW:
		ppi.keyboardRows[5] ^= 0b00010000
	case glfw.KeyX:
		ppi.keyboardRows[5] ^= 0b00100000
	case glfw.KeyY:
		ppi.keyboardRows[5] ^= 0b01000000
	case glfw.KeyZ:
		ppi.keyboardRows[5] ^= 0b10000000

	case glfw.KeyLeftShift, glfw.KeyRightShift:
		ppi.keyboardRows[6] ^= 0b00000001
	case glfw.KeyLeftControl, glfw.KeyRightControl:
		ppi.keyboardRows[6] ^= 0b00000010
	// case glfw.KeyR:
	// 	ppi.keyboardRows[6] ^= 0b00000100
	// case glfw.KeyT:
	// 	ppi.keyboardRows[6] ^= 0b00001000
	// case glfw.KeyG:
	// 	ppi.keyboardRows[6] ^= 0b00010000
	case glfw.KeyF1:
		ppi.keyboardRows[6] ^= 0b00100000
	case glfw.KeyF2:
		ppi.keyboardRows[6] ^= 0b01000000
	case glfw.KeyF3:
		ppi.keyboardRows[6] ^= 0b10000000

	case glfw.KeyF4:
		ppi.keyboardRows[7] ^= 0b00000001
	case glfw.KeyF5:
		ppi.keyboardRows[7] ^= 0b00000010
	case glfw.KeyEscape:
		ppi.keyboardRows[7] ^= 0b00000100
		// case glfw.KeyF6:
		// ppi.keyboardRows[7] ^= 0b00001000
		// case glfw.Key:
		// 	ppi.keyboardRows[7] ^= 0b00010000
		// case glfw.KeyD:
		// 	ppi.keyboardRows[7] ^= 0b00100000
		// case glfw.KeyC:
		// 	ppi.keyboardRows[7] ^= 0b01000000
	case glfw.KeyEnter:
		ppi.keyboardRows[7] ^= 0b10000000

	case glfw.KeySpace:
		ppi.keyboardRows[8] ^= 0b00000001
	case glfw.KeyHome:
		ppi.keyboardRows[8] ^= 0b00000010
	case glfw.KeyInsert:
		ppi.keyboardRows[8] ^= 0b00000100
	case glfw.KeyDelete:
		ppi.keyboardRows[8] ^= 0b00001000
	case glfw.KeyLeft:
		ppi.keyboardRows[8] ^= 0b00010000
	case glfw.KeyUp:
		ppi.keyboardRows[8] ^= 0b00100000
	case glfw.KeyDown:
		ppi.keyboardRows[8] ^= 0b01000000
	case glfw.KeyRight:
		ppi.keyboardRows[8] ^= 0b10000000

		// case glfw.KeyUp: JOY1
		// 	ppi.keyboardRows[9] ^= 0b00000001
		// case glfw.KeyRight:
		// 	ppi.keyboardRows[9] ^= 0b00000010
		// case glfw.KeyDown:
		// 	ppi.keyboardRows[9] ^= 0b00000100
		// case glfw.KeyF9:
		// 	ppi.keyboardRows[9] ^= 0b00001000
		// case glfw.KeyF6:
		// 	ppi.keyboardRows[9] ^= 0b00010000
		// case glfw.KeyF3:
		// 	ppi.keyboardRows[9] ^= 0b00100000
		// case glfw.KeyEnter:
		// 	ppi.keyboardRows[9] ^= 0b01000000
		// case glfw.KeyDelete, "BackSpace":
		// 	ppi.keyboardRows[9] ^= 0b10000000

		// default:
		// fmt.Println("key:", key.Name)
	}
}
