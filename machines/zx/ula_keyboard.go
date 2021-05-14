package zx

import "github.com/veandco/go-sdl2/sdl"

func (ula *ula) OnKey(key sdl.Scancode) {
	//fmt.Println("key:", key.Name)
	// 	switch key {

	// 	case glfw.Key1:
	// 		ula.keyboardRow[3] ^= 0b00000001
	// 	case glfw.Key2:
	// 		ula.keyboardRow[3] ^= 0b00000010
	// 	case glfw.Key3:
	// 		ula.keyboardRow[3] ^= 0b00000100
	// 	case glfw.Key4:
	// 		ula.keyboardRow[3] ^= 0b00001000
	// 	case glfw.Key5:
	// 		ula.keyboardRow[3] ^= 0b00010000

	// 	case glfw.Key0:
	// 		ula.keyboardRow[4] ^= 0b00000001
	// 	case glfw.Key9:
	// 		ula.keyboardRow[4] ^= 0b00000010
	// 	case glfw.Key8:
	// 		ula.keyboardRow[4] ^= 0b00000100
	// 	case glfw.Key7:
	// 		ula.keyboardRow[4] ^= 0b00001000
	// 	case glfw.Key6:
	// 		ula.keyboardRow[4] ^= 0b00010000

	// 	case glfw.KeyQ:
	// 		ula.keyboardRow[2] ^= 0b00000001
	// 	case glfw.KeyW:
	// 		ula.keyboardRow[2] ^= 0b00000010
	// 	case glfw.KeyE:
	// 		ula.keyboardRow[2] ^= 0b00000100
	// 	case glfw.KeyR:
	// 		ula.keyboardRow[2] ^= 0b00001000
	// 	case glfw.KeyT:
	// 		ula.keyboardRow[2] ^= 0b00010000

	// 	case glfw.KeyP:
	// 		ula.keyboardRow[5] ^= 0b00000001
	// 	case glfw.KeyO:
	// 		ula.keyboardRow[5] ^= 0b00000010
	// 	case glfw.KeyI:
	// 		ula.keyboardRow[5] ^= 0b00000100
	// 	case glfw.KeyU:
	// 		ula.keyboardRow[5] ^= 0b00001000
	// 	case glfw.KeyY:
	// 		ula.keyboardRow[5] ^= 0b00010000

	// 	case glfw.KeyA:
	// 		ula.keyboardRow[1] ^= 0b00000001
	// 	case glfw.KeyS:
	// 		ula.keyboardRow[1] ^= 0b00000010
	// 	case glfw.KeyD:
	// 		ula.keyboardRow[1] ^= 0b00000100
	// 	case glfw.KeyF:
	// 		ula.keyboardRow[1] ^= 0b00001000
	// 	case glfw.KeyG:
	// 		ula.keyboardRow[1] ^= 0b00010000

	// 	case glfw.KeyEnter:
	// 		ula.keyboardRow[6] ^= 0b00000001
	// 	case glfw.KeyL:
	// 		ula.keyboardRow[6] ^= 0b00000010
	// 	case glfw.KeyK:
	// 		ula.keyboardRow[6] ^= 0b00000100
	// 	case glfw.KeyJ:
	// 		ula.keyboardRow[6] ^= 0b00001000
	// 	case glfw.KeyH:
	// 		ula.keyboardRow[6] ^= 0b00010000

	// 	case glfw.KeyLeftShift, glfw.KeyRightShift:
	// 		ula.keyboardRow[0] ^= 0b00000001
	// 	case glfw.KeyZ:
	// 		ula.keyboardRow[0] ^= 0b00000010
	// 	case glfw.KeyX:
	// 		ula.keyboardRow[0] ^= 0b00000100
	// 	case glfw.KeyC:
	// 		ula.keyboardRow[0] ^= 0b00001000
	// 	case glfw.KeyV:
	// 		ula.keyboardRow[0] ^= 0b00010000

	// 	case glfw.KeySpace:
	// 		ula.keyboardRow[7] ^= 0b00000001
	// 	case glfw.KeyLeftControl, glfw.KeyRightControl:
	// 		ula.keyboardRow[7] ^= 0b00000010
	// 	case glfw.KeyM:
	// 		ula.keyboardRow[7] ^= 0b00000100
	// 	case glfw.KeyN:
	// 		ula.keyboardRow[7] ^= 0b00001000
	// 	case glfw.KeyB:
	// 		ula.keyboardRow[7] ^= 0b00010000

	// 	case glfw.KeyBackspace:
	// 		ula.keyboardRow[0] ^= 0b00000001
	// 		ula.keyboardRow[4] ^= 0b00000001

	// 	case glfw.KeyUp:
	// 		ula.keyboardRow[0] ^= 0b00000001
	// 		ula.keyboardRow[4] ^= 0b00001000

	// 	case glfw.KeyDown:
	// 		ula.keyboardRow[0] ^= 0b00000001
	// 		ula.keyboardRow[4] ^= 0b00010000
	// 	}
}

// var onlyOnce sync.Once

func (ula *ula) LoadCommand() {
	// 	go onlyOnce.Do(func() {
	// 		time.Sleep(time.Second)
	// 		ula.OnKey(glfw.KeyJ)
	// 		time.Sleep(150 * time.Millisecond)
	// 		ula.OnKey(glfw.KeyJ)
	// 		ula.OnKey(glfw.KeyLeftControl)
	// 		time.Sleep(150 * time.Millisecond)
	// 		ula.OnKey(glfw.KeyP)
	// 		time.Sleep(150 * time.Millisecond)
	// 		ula.OnKey(glfw.KeyP)
	// 		time.Sleep(150 * time.Millisecond)
	// 		ula.OnKey(glfw.KeyP)
	// 		time.Sleep(150 * time.Millisecond)
	// 		ula.OnKey(glfw.KeyP)
	// 		ula.OnKey(glfw.KeyLeftControl)
	// 		time.Sleep(150 * time.Millisecond)
	// 		ula.OnKey(glfw.KeyEnter)
	// 		time.Sleep(150 * time.Millisecond)
	// 		ula.OnKey(glfw.KeyEnter)
	// 	})
	// 	return
}

func (ula *ula) LoadCommand128() {
	// go onlyOnce.Do(func() {
	// 	time.Sleep(time.Second)
	// 	ula.OnKey(glfw.KeyEnter)
	// 	time.Sleep(150 * time.Millisecond)
	// 	ula.OnKey(glfw.KeyEnter)
	// })
	// return
}
