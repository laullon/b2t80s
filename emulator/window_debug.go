package emulator

import (
	"github.com/laullon/b2t80s/gui"
	"github.com/veandco/go-sdl2/sdl"
)

type debugWindow struct {
	window  *sdl.Window
	context sdl.GLContext
	ui      gui.GUIObject
}

func NewDebugWindow(name string, machine Machine) Window {
	win := &debugWindow{}

	window, err := sdl.CreateWindow("Debug", 850, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_RESIZABLE|sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	win.window = window

	context, err := window.GLCreateContext()
	if err != nil {
		panic(err)
	}
	win.context = context

	gui.RegisterWindow(window, context, win.Render)

	bt := gui.NewButton("staus", gui.Rect{330, 0, 330, 50})
	bt2 := gui.NewButton("staus", gui.Rect{330, 0, 330, 50})

	grid := gui.NewHGrid(3, 50)
	grid.Add(bt, bt)
	grid.Add(bt, bt2)
	grid.Resize(gui.Rect{0, 0, 800, 600})

	win.ui = grid

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
	win.ui.Render()
}
