package emulator

import (
	"time"

	"github.com/laullon/b2t80s/gui"
	"github.com/veandco/go-sdl2/sdl"
)

type DebugControl interface {
	Update()
}

type debug struct {
	window gui.Window
}

func NewDebugWindow(name string, machine Machine) *debug {
	debug := &debug{
		window: gui.NewWindow(name, gui.Size{800, 600}),
	}

	bt1 := gui.NewButton("Stop")
	bt2 := gui.NewButton("Stop Interrup")
	bt3 := gui.NewButton("Continue")
	bt4 := gui.NewButton("Step")
	bt5 := gui.NewButton("Step Line")
	bt6 := gui.NewButton("Step Frame")

	grid := gui.NewHGrid(3, 50)
	grid.Add(bt1, bt2, bt3, bt4, bt5, bt6)

	tabs := gui.NewTabs()

	controls := machine.Control()
	for name, ctl := range controls {
		tabs.AddTabs(name, ctl)
	}

	wait := time.Duration(time.Second / 4)
	ticker := time.NewTicker(wait)
	go func() {
		for range ticker.C {
			for _, ctl := range controls {
				ctl.(DebugControl).Update()
			}
		}
	}()

	hct := gui.NewVerticalHCT()
	hct.SetHead(grid, 100)
	hct.SetCenter(tabs)

	debug.window.SetMainUI(hct)
	debug.window.AddMouseListeners(bt1, bt2, bt3, bt4, bt5, bt6)
	debug.window.AddMouseListeners(tabs.Tabs()...)
	debug.window.SetOnKey(func(s sdl.Scancode) {})
	return debug
}
