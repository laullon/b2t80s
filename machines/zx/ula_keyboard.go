package zx

import (
	"sync"
	"time"

	"fyne.io/fyne"
)

func (ula *ula) OnKeyEvent(key *fyne.KeyEvent) {
	//fmt.Println("key:", key.Name)
	switch key.Name {

	case fyne.Key1:
		ula.keyboardRow[3] ^= 0b00000001
	case fyne.Key2:
		ula.keyboardRow[3] ^= 0b00000010
	case fyne.Key3:
		ula.keyboardRow[3] ^= 0b00000100
	case fyne.Key4:
		ula.keyboardRow[3] ^= 0b00001000
	case fyne.Key5:
		ula.keyboardRow[3] ^= 0b00010000

	case fyne.Key0:
		ula.keyboardRow[4] ^= 0b00000001
	case fyne.Key9:
		ula.keyboardRow[4] ^= 0b00000010
	case fyne.Key8:
		ula.keyboardRow[4] ^= 0b00000100
	case fyne.Key7:
		ula.keyboardRow[4] ^= 0b00001000
	case fyne.Key6:
		ula.keyboardRow[4] ^= 0b00010000

	case fyne.KeyQ:
		ula.keyboardRow[2] ^= 0b00000001
	case fyne.KeyW:
		ula.keyboardRow[2] ^= 0b00000010
	case fyne.KeyE:
		ula.keyboardRow[2] ^= 0b00000100
	case fyne.KeyR:
		ula.keyboardRow[2] ^= 0b00001000
	case fyne.KeyT:
		ula.keyboardRow[2] ^= 0b00010000

	case fyne.KeyP:
		ula.keyboardRow[5] ^= 0b00000001
	case fyne.KeyO:
		ula.keyboardRow[5] ^= 0b00000010
	case fyne.KeyI:
		ula.keyboardRow[5] ^= 0b00000100
	case fyne.KeyU:
		ula.keyboardRow[5] ^= 0b00001000
	case fyne.KeyY:
		ula.keyboardRow[5] ^= 0b00010000

	case fyne.KeyA:
		ula.keyboardRow[1] ^= 0b00000001
	case fyne.KeyS:
		ula.keyboardRow[1] ^= 0b00000010
	case fyne.KeyD:
		ula.keyboardRow[1] ^= 0b00000100
	case fyne.KeyF:
		ula.keyboardRow[1] ^= 0b00001000
	case fyne.KeyG:
		ula.keyboardRow[1] ^= 0b00010000

	case fyne.KeyReturn:
		ula.keyboardRow[6] ^= 0b00000001
	case fyne.KeyL:
		ula.keyboardRow[6] ^= 0b00000010
	case fyne.KeyK:
		ula.keyboardRow[6] ^= 0b00000100
	case fyne.KeyJ:
		ula.keyboardRow[6] ^= 0b00001000
	case fyne.KeyH:
		ula.keyboardRow[6] ^= 0b00010000

	case "LeftShift", "RightShift":
		ula.keyboardRow[0] ^= 0b00000001
	case fyne.KeyZ:
		ula.keyboardRow[0] ^= 0b00000010
	case fyne.KeyX:
		ula.keyboardRow[0] ^= 0b00000100
	case fyne.KeyC:
		ula.keyboardRow[0] ^= 0b00001000
	case fyne.KeyV:
		ula.keyboardRow[0] ^= 0b00010000

	case "Space":
		ula.keyboardRow[7] ^= 0b00000001
	case "LeftControl", "RightControl":
		ula.keyboardRow[7] ^= 0b00000010
	case fyne.KeyM:
		ula.keyboardRow[7] ^= 0b00000100
	case fyne.KeyN:
		ula.keyboardRow[7] ^= 0b00001000
	case fyne.KeyB:
		ula.keyboardRow[7] ^= 0b00010000

	case "BackSpace":
		ula.keyboardRow[0] ^= 0b00000001
		ula.keyboardRow[4] ^= 0b00000001

	case fyne.KeyUp:
		ula.keyboardRow[0] ^= 0b00000001
		ula.keyboardRow[4] ^= 0b00001000

	case fyne.KeyDown:
		ula.keyboardRow[0] ^= 0b00000001
		ula.keyboardRow[4] ^= 0b00010000
	}
}

var onlyOnce sync.Once

func (ula *ula) LoadCommand() uint16 {
	go onlyOnce.Do(func() {
		time.Sleep(time.Second)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyJ})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyJ})
		ula.OnKeyEvent(&fyne.KeyEvent{Name: "RightSuper"})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyP})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyP})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyP})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyP})
		ula.OnKeyEvent(&fyne.KeyEvent{Name: "RightSuper"})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyReturn})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyReturn})
	})
	return 0
}
