package emulator

import (
	"image"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

type Monitor interface {
	Canvas() *canvas.Image
	FrameDone()
	FPS() float64
}

type monitor struct {
	vram    *image.RGBA
	screen  *image.RGBA
	display *canvas.Image
	start   time.Time
	frames  float64
}

func NewMonitor(img *image.RGBA) Monitor {
	monitor := &monitor{
		vram:   img,
		screen: image.NewRGBA(img.Bounds()),
		start:  time.Now(),
	}

	monitor.display = canvas.NewImageFromImage(monitor.screen)
	monitor.display.FillMode = canvas.ImageFillOriginal
	monitor.display.ScaleMode = canvas.ImageScalePixels
	monitor.display.SetMinSize(fyne.NewSize(352*2, 296*2))

	return monitor
}

func (monitor *monitor) Canvas() *canvas.Image {
	return monitor.display
}

func (monitor *monitor) FrameDone() {
	monitor.frames++
	copy(monitor.screen.Pix, monitor.vram.Pix)
	go func() {
		monitor.display.Refresh()
	}()
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
