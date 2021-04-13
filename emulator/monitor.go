package emulator

import (
	"time"

	"fyne.io/fyne/v2/canvas"
	"github.com/laullon/b2t80s/ui"
)

type Monitor interface {
	Canvas() *canvas.Image
	Screen() *ui.Display
	FrameDone()
	FPS() float64
}

type monitor struct {
	vram    *ui.Display
	screen  *ui.Display
	display *canvas.Image
	start   time.Time
	frames  float64
}

func NewMonitor(img *ui.Display) Monitor {
	monitor := &monitor{
		vram:   img,
		screen: ui.NewDisplay(img.Bounds()),
		start:  time.Now(),
	}

	// monitor.display = canvas.NewImageFromImage(monitor.screen)
	// monitor.display.FillMode = canvas.ImageFillOriginal
	// monitor.display.ScaleMode = canvas.ImageScalePixels
	// monitor.display.SetMinSize(fyne.NewSize(352*2, 296*2))

	return monitor
}

func (monitor *monitor) Screen() *ui.Display {
	return monitor.screen
}

func (monitor *monitor) Canvas() *canvas.Image {
	return monitor.display
}

func (monitor *monitor) FrameDone() {
	monitor.frames++
	copy(monitor.screen.Pix, monitor.vram.Pix)
	// go func() {
	// monitor.display.Refresh()
	// }()
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
