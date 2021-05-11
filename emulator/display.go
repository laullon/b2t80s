package emulator

import (
	"image"
	"image/color"
)

type Display struct {
	ViewPortRect image.Rectangle
	Size         image.Point
	Image        *image.RGBA
}

func NewDisplay(r image.Rectangle) *Display {
	res := &Display{}
	res.Image = image.NewRGBA(r)
	res.ViewPortRect = r
	res.Size = r.Size()
	return res
}

func (p *Display) ColorModel() color.Model        { return p.Image.ColorModel() }
func (p *Display) Bounds() image.Rectangle        { return p.Image.Bounds() }
func (p *Display) At(x, y int) color.Color        { return p.Image.At(x, y) }
func (p *Display) Set(x, y int, c color.Color)    { p.Image.Set(x, p.Image.Bounds().Dy()-y, c) }
func (p *Display) SetRGBA(x, y int, c color.RGBA) { p.Image.SetRGBA(x, p.Image.Bounds().Dy()-y, c) }
