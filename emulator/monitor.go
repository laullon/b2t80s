package emulator

import (
	"image"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

type Monitor interface {
	Canvas() *canvas.Image
	FrameDone()
}

type monitor struct {
	vram    *image.RGBA
	screen  *image.RGBA
	display *canvas.Image
}

func NewMonitor(img *image.RGBA) Monitor {
	monitor := &monitor{
		vram:   img,
		screen: image.NewRGBA(img.Bounds()),
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
	copy(monitor.screen.Pix, monitor.vram.Pix)
	go func() {
		monitor.display.Refresh()
	}()
}
