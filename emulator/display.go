package emulator

import (
	"image"
	"image/color"
)

type Display struct {
	image.RGBA
}

func NewDisplay(r image.Rectangle) *Display {
	res := &Display{}
	res.Pix = make([]uint8, uint64(4)*uint64(r.Max.X)*uint64(r.Max.Y)*2)
	res.Stride = 4 * r.Dx()
	res.Rect = r
	return res
}

func (p *Display) SetRGBA(x, y int, c color.RGBA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, p.Rect.Max.Y-y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c.R
	s[1] = c.G
	s[2] = c.B
	s[3] = c.A
}
