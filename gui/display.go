package gui

import (
	"image"
	"image/color"
)

type Display struct {
	size         Size  // image size
	Start        Point // image 0x0 position
	ViewPortRect Rect  // rect to display
	ViewSize     Size  // size to display

	back, front *glImage        // TODO: no need it, just pix
	rect        image.Rectangle // just for bounds()

	Trans DisplayTransform
}

type DisplayTransform func(x int, y int) (int, int)

func NewDisplay(s Size) *Display {
	res := &Display{}
	res.back = newImage(s)
	res.front = newImage(s)
	res.ViewPortRect = Rect{0, 0, s.W, s.H}
	res.rect = image.Rect(0, 0, int(s.W), int(s.H))
	res.ViewSize = s
	res.size = s
	res.Start = Point{0, 0}
	res.Trans = func(x, y int) (int, int) { return x, y }
	return res
}

func (dis *Display) ColorModel() color.Model { return color.RGBAModel }
func (dis *Display) Bounds() image.Rectangle { return dis.rect }
func (dis *Display) At(x, y int) color.Color { return dis.front.At(x, y) }
func (dis *Display) Pix() []uint8            { return dis.back.Pix }

func (dis *Display) Swap() {
	dis.front, dis.back = dis.back, dis.front
}

func (dis *Display) Set(x, y int, c color.Color) {
	x, y = dis.addjustXY(x, y)
	dis.front.Set(x, y, c)
}

func (dis *Display) SetRGBA(x, y int, c color.RGBA) {
	x, y = dis.Trans(x, y)
	x, y = dis.addjustXY(x, y)
	dis.front.Set(x, y, c)
}

func (dis *Display) addjustXY(x, y int) (int, int) {
	newX, newY := int32(x), int32(y)
	if dis.Start.X != 0 || dis.Start.Y != 0 {
		newX += dis.Start.X
		if newX > dis.size.W {
			newX -= dis.size.W
		}

		newY += dis.Start.Y
		if newY > dis.size.H {
			newY -= dis.size.H
		}
	}

	return int(newX), int(newY)
}
