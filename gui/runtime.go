package gui

import (
	"reflect"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var windows = make(map[uint32]*window)

func PoolEvents(stop chan struct{}) {
	for _, win := range windows {
		win.sdlWin.GLMakeCurrent(win.context)
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		win.ui.Render()
		gl.Flush()
		win.sdlWin.GLSwap()
	}
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		// fmt.Printf("%T\n", e)
		wid := reflect.ValueOf(e).Elem().FieldByName("WindowID")
		var win *window
		if wid.IsValid() {
			win = windows[uint32(wid.Uint())]
			win.sdlWin.GLMakeCurrent(win.context)
		}

		switch event := e.(type) {
		case *sdl.WindowEvent:
			if event.Event == sdl.WINDOWEVENT_CLOSE {
				sdl.Quit()
				stop <- struct{}{}
			} else if event.Event == sdl.WINDOWEVENT_RESIZED {
				win.ui.Resize(Rect{0, 0, event.Data1, event.Data2})
			}

		case *sdl.QuitEvent:
			sdl.Quit()
			stop <- struct{}{}

		case *sdl.KeyboardEvent:
			if event.Repeat == 0 {
				win.onKey(event.Keysym.Scancode)
			}

		case *sdl.MouseMotionEvent:
			_, h := win.sdlWin.GetSize()
			p := Point{event.X, h - event.Y}
			for _, obj := range win.mouseListeners {
				obj.OnMouseOver(obj.Rect().In(p))
			}

		case *sdl.MouseButtonEvent:
			_, h := win.sdlWin.GetSize()
			p := Point{event.X, h - event.Y}
			for _, obj := range win.mouseListeners {
				if obj.Rect().In(p) {
					obj.OnMouseClick(event.State == 0)
				}
			}

		case *sdl.MouseWheelEvent:
			_, h := win.sdlWin.GetSize()
			x, y, _ := sdl.GetMouseState()
			p := Point{x, h - y}
			for _, obj := range win.mouseListeners {
				if obj.Rect().In(p) {
					if scroll, ok := obj.(ScrollTarget); ok {
						scroll.OnScroll(event.X, event.Y)
						// fmt.Printf(">> event: %v \n", event)
					}
				}
			}

		case *sdl.JoyAxisEvent:
			switch event.Axis {
			case 0:
				joysticks[event.Which].R = event.Value == 32767
				joysticks[event.Which].L = event.Value == -32768
			case 1:
				joysticks[event.Which].D = event.Value == 32767
				joysticks[event.Which].U = event.Value == -32768
			}

		case *sdl.JoyButtonEvent:
			switch true {
			case event.Button == 8:
				joysticks[event.Which].Select = event.State == 1
			case event.Button == 9:
				joysticks[event.Which].Start = event.State == 1
			case event.Button%2 == 0:
				joysticks[event.Which].F = event.State == 1
			case event.Button%2 == 1:
				joysticks[event.Which].F2 = event.State == 1
			}

		case *sdl.JoyDeviceAddedEvent:
			sdl.JoystickOpen(int(event.Which))
			joysticks[event.Which].ON = true

			// default:
			// fmt.Printf(">> event: %T \n", e)
		}
	}
}
