package msx

import "fyne.io/fyne/v2"

func (ppi *ppi) OnKeyEvent(key *fyne.KeyEvent) {
	// println("key:", key.Name)

	switch key.Name {

	case fyne.Key0:
		ppi.keyboardRows[0] ^= 0b00000001
	case fyne.Key1:
		ppi.keyboardRows[0] ^= 0b00000010
	case fyne.Key2:
		ppi.keyboardRows[0] ^= 0b00000100
	case fyne.Key3:
		ppi.keyboardRows[0] ^= 0b00001000
	case fyne.Key4:
		ppi.keyboardRows[0] ^= 0b00010000
	case fyne.Key5:
		ppi.keyboardRows[0] ^= 0b00100000
	case fyne.Key6:
		ppi.keyboardRows[0] ^= 0b01000000
	case fyne.Key7:
		ppi.keyboardRows[0] ^= 0b10000000

	case fyne.Key8:
		ppi.keyboardRows[1] ^= 0b00000001
	case fyne.Key9:
		ppi.keyboardRows[1] ^= 0b00000010
	// case fyne.KeyF7:
	// 	ppi.keyboardRows[1] ^= 0b00000100
	// case fyne.KeyF8:
	// 	ppi.keyboardRows[1] ^= 0b00001000
	// case fyne.KeyF5:
	// 	ppi.keyboardRows[1] ^= 0b00010000
	// case fyne.KeyF1:
	// 	ppi.keyboardRows[1] ^= 0b00100000
	// case fyne.KeyF2:
	// 	ppi.keyboardRows[1] ^= 0b01000000
	// case fyne.KeyF10:
	// 	ppi.keyboardRows[1] ^= 0b10000000

	// case fyne.KeyUp:
	// 	ppi.keyboardRows[2] ^= 0b00000001
	// case fyne.KeyLeftBracket:
	// 	ppi.keyboardRows[2] ^= 0b00000010
	// case fyne.KeyReturn:
	// 	ppi.keyboardRows[2] ^= 0b00000100
	// case fyne.KeyRightBracket:
	// 	ppi.keyboardRows[2] ^= 0b00001000
	// case fyne.KeyF4:
	// 	ppi.keyboardRows[2] ^= 0b00010000
	// case "LeftShift", "RightShift":
	// 	ppi.keyboardRows[2] ^= 0b00100000
	case fyne.KeyA:
		ppi.keyboardRows[2] ^= 0b01000000
	case fyne.KeyB:
		ppi.keyboardRows[2] ^= 0b10000000

	case fyne.KeyC:
		ppi.keyboardRows[3] ^= 0b00000001
	case fyne.KeyD:
		ppi.keyboardRows[3] ^= 0b00000010
	case fyne.KeyE:
		ppi.keyboardRows[3] ^= 0b00000100
	case fyne.KeyF:
		ppi.keyboardRows[3] ^= 0b00001000
	case fyne.KeyG:
		ppi.keyboardRows[3] ^= 0b00010000
	case fyne.KeyH:
		ppi.keyboardRows[3] ^= 0b00100000
	case fyne.KeyI:
		ppi.keyboardRows[3] ^= 0b01000000
	case fyne.KeyJ:
		ppi.keyboardRows[3] ^= 0b10000000

	case fyne.KeyK:
		ppi.keyboardRows[4] ^= 0b00000001
	case fyne.KeyL:
		ppi.keyboardRows[4] ^= 0b00000010
	case fyne.KeyM:
		ppi.keyboardRows[4] ^= 0b00000100
	case fyne.KeyN:
		ppi.keyboardRows[4] ^= 0b00001000
	case fyne.KeyO:
		ppi.keyboardRows[4] ^= 0b00010000
	case fyne.KeyP:
		ppi.keyboardRows[4] ^= 0b00100000
	case fyne.KeyQ:
		ppi.keyboardRows[4] ^= 0b01000000
	case fyne.KeyR:
		ppi.keyboardRows[4] ^= 0b10000000

	case fyne.KeyS:
		ppi.keyboardRows[5] ^= 0b00000001
	case fyne.KeyT:
		ppi.keyboardRows[5] ^= 0b00000010
	case fyne.KeyU:
		ppi.keyboardRows[5] ^= 0b00000100
	case fyne.KeyV:
		ppi.keyboardRows[5] ^= 0b00001000
	case fyne.KeyW:
		ppi.keyboardRows[5] ^= 0b00010000
	case fyne.KeyX:
		ppi.keyboardRows[5] ^= 0b00100000
	case fyne.KeyY:
		ppi.keyboardRows[5] ^= 0b01000000
	case fyne.KeyZ:
		ppi.keyboardRows[5] ^= 0b10000000

	case "LeftShift", "RightShift":
		ppi.keyboardRows[6] ^= 0b00000001
	case "LeftControl", "RightControl":
		ppi.keyboardRows[6] ^= 0b00000010
	// case fyne.KeyR:
	// 	ppi.keyboardRows[6] ^= 0b00000100
	// case fyne.KeyT:
	// 	ppi.keyboardRows[6] ^= 0b00001000
	// case fyne.KeyG:
	// 	ppi.keyboardRows[6] ^= 0b00010000
	case fyne.KeyF1:
		ppi.keyboardRows[6] ^= 0b00100000
	case fyne.KeyF2:
		ppi.keyboardRows[6] ^= 0b01000000
	case fyne.KeyF3:
		ppi.keyboardRows[6] ^= 0b10000000

	case fyne.KeyF4:
		ppi.keyboardRows[7] ^= 0b00000001
	case fyne.KeyF5:
		ppi.keyboardRows[7] ^= 0b00000010
	case fyne.KeyEscape:
		ppi.keyboardRows[7] ^= 0b00000100
		// case fyne.KeyF6:
		// ppi.keyboardRows[7] ^= 0b00001000
		// case fyne.Key:
		// 	ppi.keyboardRows[7] ^= 0b00010000
		// case fyne.KeyD:
		// 	ppi.keyboardRows[7] ^= 0b00100000
		// case fyne.KeyC:
		// 	ppi.keyboardRows[7] ^= 0b01000000
	case fyne.KeyReturn:
		ppi.keyboardRows[7] ^= 0b10000000

	case fyne.KeySpace:
		ppi.keyboardRows[8] ^= 0b00000001
	case fyne.KeyHome:
		ppi.keyboardRows[8] ^= 0b00000010
	case fyne.KeyInsert:
		ppi.keyboardRows[8] ^= 0b00000100
	case fyne.KeyDelete:
		ppi.keyboardRows[8] ^= 0b00001000
	case fyne.KeyLeft:
		ppi.keyboardRows[8] ^= 0b00010000
	case fyne.KeyUp:
		ppi.keyboardRows[8] ^= 0b00100000
	case fyne.KeyDown:
		ppi.keyboardRows[8] ^= 0b01000000
	case fyne.KeyRight:
		ppi.keyboardRows[8] ^= 0b10000000

		// case fyne.KeyUp: JOY1
		// 	ppi.keyboardRows[9] ^= 0b00000001
		// case fyne.KeyRight:
		// 	ppi.keyboardRows[9] ^= 0b00000010
		// case fyne.KeyDown:
		// 	ppi.keyboardRows[9] ^= 0b00000100
		// case fyne.KeyF9:
		// 	ppi.keyboardRows[9] ^= 0b00001000
		// case fyne.KeyF6:
		// 	ppi.keyboardRows[9] ^= 0b00010000
		// case fyne.KeyF3:
		// 	ppi.keyboardRows[9] ^= 0b00100000
		// case fyne.KeyEnter:
		// 	ppi.keyboardRows[9] ^= 0b01000000
		// case fyne.KeyDelete, "BackSpace":
		// 	ppi.keyboardRows[9] ^= 0b10000000

		// default:
		// fmt.Println("key:", key.Name)
	}
}
