package gui

import (
	"image/color"
)

type Display struct {
	width, height  uint16 // image size
	startX, startY uint16 // image 0x0 position
	ViewPortRect   Rect   // rect to display
	ViewSize       Size   // size to display

	back, front []uint8
	// rect        image.Rectangle // just for bounds()

	Trans DisplayTransform
}

type DisplayTransform func(x uint16, y uint16) (uint16, uint16)

func NewDisplay(width, height uint16) *Display {
	dis := &Display{}

	dis.width, dis.height = width, height
	dis.back = make([]uint8, int32(width)*int32(height)*4)
	dis.front = make([]uint8, int32(width)*int32(height)*4)

	dis.ViewPortRect = Rect{0, 0, int32(width), int32(height)}
	dis.ViewSize = Size{int32(width), int32(height)}

	dis.Trans = func(x, y uint16) (uint16, uint16) { return x, y }
	return dis
}

// func (dis *Display) ColorModel() color.Model { return color.RGBAModel }
// func (dis *Display) Bounds() image.Rectangle { return dis.rect }
// func (dis *Display) At(x, y int) color.Color { return dis.front.At(x, y) }
func (dis *Display) Pix() []uint8         { return dis.back }
func (dis *Display) SetStart(x, y uint16) { dis.startX, dis.startY = x, y }

func (dis *Display) Swap() {
	dis.front, dis.back = dis.back, dis.front
}

func (dis *Display) Set(x, y uint16, c color.RGBA) {
	x, y = dis.Trans(x, y)
	x, y = dis.addjustXY(x, y)

	if x < dis.width && y < dis.height {
		y = dis.height - y - 1
		idx := uint32(y)*uint32(dis.width)*4 + uint32(x)*4
		s := dis.front[idx : idx+4 : idx+4]

		s[0] = c.R
		s[1] = c.G
		s[2] = c.B
		s[3] = c.A
	}
}

func (dis *Display) addjustXY(x, y uint16) (uint16, uint16) {
	newX, newY := x, y
	if dis.startX != 0 || dis.startY != 0 {
		newX += dis.startX
		if newX > dis.width {
			newX -= dis.width
		}

		newY += dis.startY
		if newY > dis.height {
			newY -= dis.height
		}
	}

	return newX, newY
}
