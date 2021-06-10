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

func NewDebugWindow(name string, machine Machine, debuger *debugger) *debug {
	debug := &debug{
		window: gui.NewWindow(name, gui.Size{800, 600}),
	}

	bt1 := gui.NewButton("Stop")
	bt1.SetAction(func() { debuger.Stop() })

	bt2 := gui.NewButton("Stop Interrup")
	bt2.SetAction(func() { debuger.StopNextInterrupt() })

	bt3 := gui.NewButton("Continue")
	bt3.SetAction(func() { debuger.Continue() })

	bt4 := gui.NewButton("Step")
	bt4.SetAction(func() { debuger.Step() })

	bt5 := gui.NewButton("Step Line")
	bt5.SetAction(func() { debuger.StepLine() })

	bt6 := gui.NewButton("Step Frame")
	bt6.SetAction(func() { debuger.StepFrame() })

	grid := gui.NewHGrid(3, 50, 4)
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
	debug.window.AddMouseListeners(hct.GetMouseTargets()...)
	debug.window.SetOnKey(func(s sdl.Scancode) {})
	debug.window.MoveTo(gui.Point{100, 100})
	return debug
}
