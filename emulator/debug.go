package emulator

import (
	"github.com/laullon/b2t80s/gui"
)

type debug struct {
	window gui.Window
}

func NewDebugWindow(name string, machine Machine) *debug {
	debug := &debug{
		window: gui.NewWindow(name, gui.Size{800, 600}),
	}

	bt1 := gui.NewButton("staus 1", gui.Rect{330, 0, 330, 50})
	bt2 := gui.NewButton("staus 2", gui.Rect{330, 0, 330, 50})
	bt3 := gui.NewButton("staus 3", gui.Rect{330, 0, 330, 50})

	grid := gui.NewHGrid(3, 50)
	grid.Add(bt1, bt2, bt3)
	grid.Resize(gui.Rect{0, 0, 800, 600})

	debug.window.SetMainUI(grid)
	debug.window.AddMouseListeners(bt1, bt2, bt3)

	return debug
}
