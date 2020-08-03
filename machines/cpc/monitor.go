package cpc

import (
	"image"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"golang.org/x/image/draw"
)

type monitor struct {
	display       *image.RGBA
	displayScaled *image.RGBA
	canvas        *canvas.Image
}

func NewMonitor() *monitor {
	monitor := &monitor{
		displayScaled: image.NewRGBA(image.Rect(0, 0, 969, 642)),
		display:       image.NewRGBA(image.Rect(0, 0, 960, 312)),
	}

	monitor.canvas = canvas.NewImageFromImage(monitor.displayScaled)
	monitor.canvas.FillMode = canvas.ImageFillOriginal
	monitor.canvas.ScaleMode = canvas.ImageScalePixels
	monitor.canvas.SetMinSize(fyne.NewSize(352*2, 296*2))

	return monitor
}

func (monitor *monitor) Canvas() *canvas.Image {
	return monitor.canvas
}

func (monitor *monitor) FrameDone() {
	// TODO write a custom function to double horizontal lines, no need for this
	draw.NearestNeighbor.Scale(monitor.displayScaled, monitor.displayScaled.Bounds(), monitor.display, monitor.display.Bounds(), draw.Over, nil)
	// copy(monitor.displayScaled.Pix, monitor.display.Pix)

	go func() {
		monitor.canvas.Refresh()
	}()
}
