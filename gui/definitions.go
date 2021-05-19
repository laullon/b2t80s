package gui

import (
	"image"
	"image/color"
)

type GUIObject interface {
	Render()
	Resize(Rect)
}

type MouseTarget interface {
	Rect() Rect
	OnMouseOver(bool)
	OnMouseClick(bool)
}

type Point struct{ X, Y int32 }
type Size struct{ W, H int32 }
type Rect struct{ X, Y, W, H int32 }

func (r Rect) In(p Point) bool {
	x := p.X - r.X
	if x > 0 && x < r.W {
		y := p.Y - r.Y
		if y > 0 && y < r.H {
			return true
		}
	}
	return false
}

// ******************************************************
// ******************************************************
// ******************************************************

type glImage struct {
	Pix  []uint8
	rect image.Rectangle
}

func newImage(size Size) *glImage {
	return &glImage{
		rect: image.Rect(0, 0, int(size.W), int(size.H)),
		Pix:  make([]uint8, uint64(size.H)*uint64(size.W)*4),
	}
}

func (i *glImage) Bounds() image.Rectangle { return i.rect }
func (i *glImage) ColorModel() color.Model { return color.RGBAModel }

func (i *glImage) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(i.rect)) {
		return color.RGBA{}
	}
	idx := uint64(i.rect.Dy()-y-1)*uint64(i.rect.Max.X)*4 + uint64(x)*4
	s := i.Pix[idx : idx+4 : idx+4]
	return color.RGBA{s[0], s[1], s[2], s[3]}
}

func (i *glImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(i.rect)) {
		return
	}
	idx := uint64(i.rect.Dy()-y-1)*uint64(i.rect.Max.X)*4 + uint64(x)*4
	c1 := color.RGBAModel.Convert(c).(color.RGBA)
	s := i.Pix[idx : idx+4 : idx+4] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c1.R
	s[1] = c1.G
	s[2] = c1.B
	s[3] = c1.A
}
