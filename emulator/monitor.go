package emulator

import (
	"time"
)

type Monitor interface {
	Screen() *Display
	FrameDone()
	FPS() float64
	SetRedraw(redraw func())
}

type monitor struct {
	vram   *Display
	screen *Display
	start  time.Time
	frames float64
	redraw func()
}

func NewMonitor(img *Display) Monitor {
	monitor := &monitor{
		vram:   img,
		screen: NewDisplay(img.Bounds()),
		start:  time.Now(),
	}

	// monitor.display = canvas.NewImageFromImage(monitor.screen)
	// monitor.display.FillMode = canvas.ImageFillOriginal
	// monitor.display.ScaleMode = canvas.ImageScalePixels
	// monitor.display.SetMinSize(fyne.NewSize(352*2, 296*2))

	return monitor
}

func (monitor *monitor) SetRedraw(redraw func()) {
	monitor.redraw = redraw
}

func (monitor *monitor) Screen() *Display {
	return monitor.screen
}

func (monitor *monitor) FrameDone() {
	monitor.frames++
	copy(monitor.screen.Pix, monitor.vram.Pix)
	monitor.redraw()
}

func (monitor *monitor) FPS() float64 {
	seconds := time.Now().Sub(monitor.start).Seconds()
	res := monitor.frames / seconds
	if seconds > 2 {
		monitor.frames = 0
		monitor.start = time.Now()
	}
	return res
}
