package emulator

import (
	"image"
	"image/color"
)

type Display struct {
	ViewPortRect image.Rectangle
	Start        image.Point
	Size         image.Point
	Image        *image.RGBA
}

func NewDisplay(r image.Rectangle) *Display {
	res := &Display{}
	res.Image = image.NewRGBA(r)
	res.ViewPortRect = r
	res.Size = r.Size()
	res.Start = image.Point{0, 0}
	return res
}

func (p *Display) ColorModel() color.Model { return p.Image.ColorModel() }
func (p *Display) Bounds() image.Rectangle { return p.Image.Bounds() }
func (p *Display) At(x, y int) color.Color { return p.Image.At(x, y) }

func (p *Display) Set(x, y int, c color.Color) {
	x, y = p.addjustXY(x, y)
	p.Image.Set(x, y, c)
}

func (p *Display) SetRGBA(x, y int, c color.RGBA) {
	x, y = p.addjustXY(x, y)
	p.Image.SetRGBA(x, y, c)
}

func (p *Display) addjustXY(x, y int) (int, int) {
	if p.Start.X != 0 || p.Start.Y != 0 {
		x += p.Start.X
		if x > p.Size.X {
			x -= p.Size.X
		}

		y += p.Start.Y
		if y > p.Size.Y {
			y -= p.Size.Y
		}
	}

	return x, p.Image.Bounds().Dy() - y
}
