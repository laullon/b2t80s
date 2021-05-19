package emulator

import (
	"github.com/laullon/b2t80s/gui"
	"github.com/veandco/go-sdl2/sdl"
)

type debugWindow struct {
	win gui.Window
}

func NewDebugWindow(name string, machine Machine) Window {
	win := &debugWindow{
		win: gui.NewWindow(name, gui.Size{800, 600}),
	}

	bt1 := gui.NewButton("staus 1", gui.Rect{330, 0, 330, 50})
	bt2 := gui.NewButton("staus 2", gui.Rect{330, 0, 330, 50})

	grid := gui.NewHGrid(3, 50)
	grid.Add(bt1, bt2)
	grid.Resize(gui.Rect{0, 0, 800, 600})

	win.win.SetMainUI(grid)
	win.win.AddMouseListeners(bt1, bt2)

	return win
}

func (win *debugWindow) SetStatus(txt string) {
}

func (win *debugWindow) SetOnKey(func(sdl.Scancode)) {
}

func (win *debugWindow) Run() {
	// wait := time.Duration(time.Second)
	// ticker := time.NewTicker(wait)
	// go func() {
	// 	for range ticker.C {
	// 		win.window.GLMakeCurrent(win.context)
	// 		gl.ClearColor(0, rand.Float32(), 0, 1)
	// 		gl.Clear(gl.COLOR_BUFFER_BIT)
	// 		win.ui.Render()
	// 		win.window.GLSwap()
	// 	}
	// }()
}

func (win *debugWindow) Render() {
	// gl.ClearColor(0, rand.Float32(), 0, 1)
	// gl.Clear(gl.COLOR_BUFFER_BIT)
	// win.ui.Render()
}
