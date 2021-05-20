package emulator

import (
	"time"

	"github.com/laullon/b2t80s/gui"
)

type Monitor interface {
	Screen() *gui.Display
	FrameDone()
	FPS() float64
	SetRedraw(redraw func())
}

type monitor struct {
	display *gui.Display
	start   time.Time
	frames  float64
	redraw  func()
}

func NewMonitor(img *gui.Display) Monitor {
	monitor := &monitor{
		display: img,
		start:   time.Now(),
	}

	return monitor
}

func (monitor *monitor) SetRedraw(redraw func()) {
	monitor.redraw = redraw
}

func (monitor *monitor) Screen() *gui.Display {
	return monitor.display
}

func (monitor *monitor) FrameDone() {
	monitor.frames++
	monitor.display.Swap()
	go func() {
		monitor.redraw()
	}()
}

func (monitor *monitor) FPS() float64 {
	seconds := time.Since(monitor.start).Seconds()
	res := monitor.frames / seconds
	if seconds > 2 {
		monitor.frames = 0
		monitor.start = time.Now()
	}
	return res
}
