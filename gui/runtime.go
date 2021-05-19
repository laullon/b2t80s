package gui

import (
	"fmt"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type window struct {
	window  *sdl.Window
	context sdl.GLContext
	render  func()
}

var windows = make([]*window, 0)

func RegisterWindow(w *sdl.Window, ctx sdl.GLContext, render func()) {
	windows = append(windows, &window{w, ctx, render})
}

func PoolEvents() {
	for running := true; running; {
		for _, win := range windows {
			win.window.GLMakeCurrent(win.context)
			win.render()
			gl.Flush()
			win.window.GLSwap()
		}
		// 	_, h := win.window.GetSize()
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			// 		// fmt.Printf("%T\n", e)
			switch event := e.(type) {
			case *sdl.WindowEvent:
				// fmt.Printf("%v\n", event)
				if event.Event == sdl.WINDOWEVENT_CLOSE {
					running = false
					sdl.Quit()
				} else if event.Event == sdl.WINDOWEVENT_RESIZED {
					// win.ui.Resize(gui.Rect{0, 0, event.Data1, event.Data2})
				}

				// 		case *sdl.QuitEvent:
				// 			running = false
				// 			sdl.Quit()

				// 		case *sdl.KeyboardEvent:
				// 			if event.Repeat == 0 {
				// 				// fmt.Printf("%d\n", event.Keysym.Scancode)
				// 				win.onKey(event.Keysym.Scancode)
				// 			}

				// 		case *sdl.MouseMotionEvent:
				// 			p := gui.Point{event.X, h - event.Y}
				// 			// fmt.Printf("%v\n", p)
				// 			for _, obj := range win.ui_mouse {
				// 				obj.OnMouseOver(obj.Rect().In(p))
				// 			}

			case *sdl.MouseButtonEvent:
				p := Point{event.X, event.Y}
				fmt.Printf("win:%d - %v %v %v %v\n", event.WindowID, p, event.Button, event.Clicks, event.State)
				// 			for _, obj := range win.ui_mouse {
				// 				if obj.Rect().In(p) {
				// 					obj.OnMouseClick(event.State == 0)
				// 				}
				// }
			}
		}
	}
}
