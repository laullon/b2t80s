package emulator

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Kempston struct {
}

func (k *Kempston) ReadPort(port uint16) (byte, bool) {
	axes := glfw.GetJoystickAxes(glfw.Joystick1)
	buttons := glfw.GetJoystickButtons(glfw.Joystick1)
	res := byte(0)

	if len(buttons) > 2 && len(axes) > 2 {

		res |= buttons[1] << 4
		if axes[0] == -1 {
			res |= 0b10
		} else if axes[0] == 1 {
			res |= 0b1
		}

		if axes[1] == -1 {
			res |= 0b1000
		} else if axes[1] == 1 {
			res |= 0b100
		}
	}

	// log.Printf("-> %08b %v %v", res, axes, buttons)
	return res, false
}

func (k *Kempston) WritePort(port uint16, data byte) {
}
