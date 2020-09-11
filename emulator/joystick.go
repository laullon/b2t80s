package emulator

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Joystick struct {
	ON                bool
	U, D, R, L, F, F2 bool
}

var j1 = &Joystick{}
var j2 = &Joystick{}

var js = []*Joystick{j1, j2}
var usb = []glfw.Joystick{glfw.Joystick1, glfw.Joystick2}

func ReadJoystick() (*Joystick, *Joystick) {
	for idx, j := range usb {
		axes := j.GetAxes()
		buttons := j.GetButtons()

		if len(buttons) > 0 && len(axes) > 1 {
			js[idx].ON = true
			js[idx].F = buttons[1] != 0
			js[idx].F2 = buttons[2] != 0
			js[idx].L = axes[0] == -1
			js[idx].R = axes[0] == 1
			js[idx].U = axes[1] == -1
			js[idx].D = axes[1] == 1
		} else {
			js[idx].ON = false

		}
		// fmt.Printf("j%d -> %v (%v %v)\n", idx, js[0], axes, buttons)
	}

	return j1, j2
}
