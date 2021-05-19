package gui

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var windows = make(map[uint32]*window)

func PoolEvents() {
	for running := true; running; {
		for _, win := range windows {
			win.sdlWin.GLMakeCurrent(win.context)
			gl.ClearColor(0, 0, 0, 1)
			gl.Clear(gl.COLOR_BUFFER_BIT)
			win.ui.Render()
			gl.Flush()
			win.sdlWin.GLSwap()
		}
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			// fmt.Printf("%T\n", e)
			switch event := e.(type) {
			case *sdl.WindowEvent:
				if event.Event == sdl.WINDOWEVENT_CLOSE {
					running = false
					sdl.Quit()
				} else if event.Event == sdl.WINDOWEVENT_RESIZED {
					win := windows[event.WindowID]
					win.ui.Resize(Rect{0, 0, event.Data1, event.Data2})
				}

			case *sdl.QuitEvent:
				running = false
				sdl.Quit()

			case *sdl.KeyboardEvent:
				if event.Repeat == 0 {
					win := windows[event.WindowID]
					win.onKey(event.Keysym.Scancode)
				}

			case *sdl.MouseMotionEvent:
				win := windows[event.WindowID]
				_, h := win.sdlWin.GetSize()
				p := Point{event.X, h - event.Y}
				for _, obj := range win.mouseListeners {
					obj.OnMouseOver(obj.Rect().In(p))
				}

			case *sdl.MouseButtonEvent:
				win := windows[event.WindowID]
				_, h := win.sdlWin.GetSize()
				p := Point{event.X, h - event.Y}
				for _, obj := range win.mouseListeners {
					if obj.Rect().In(p) {
						obj.OnMouseClick(event.State == 0)
					}
				}
			}
		}
	}
}
