package gui

import (
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

var palette = []color.RGBA{ // https://www.slideteam.net/blog/9-beautiful-color-palettes-for-designing-powerful-powerpoint-slides
	{0xff, 0xff, 0xff, 0xff},
	{4, 37, 58, 0xff},     //Dark Blue
	{225, 221, 191, 0xff}, //Tan
	{76, 131, 122, 0xff},  //Green
}

type GUIObject interface {
	Render()
	Resize(Rect)
	GetMouseTargets() []MouseTarget
}

type MouseTarget interface {
	Rect() Rect
	OnMouseOver(bool)
	OnMouseClick(bool)
}

type ScrollTarget interface {
	MouseTarget
	OnScroll(x, y int32)
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

func (r Rect) Relative(new Rect) Rect {
	new.X += r.X
	new.Y += r.Y
	return new
}

// ******************************************************
// ******************************************************
// ******************************************************

type circle struct {
	p    image.Point
	r    int
	r2   int
	in   color.RGBA
	out  color.RGBA
	line color.RGBA
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr, rr2 := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r), float64(c.r2)
	if xx*xx+yy*yy < rr2*rr2 {
		return c.in
	} else if xx*xx+yy*yy < rr*rr {
		return c.line
	}

	return c.out
}

// ******************************************************
// ******************************************************
// ******************************************************

type Joystick struct {
	ON                               bool
	U, D, R, L, F, F2, Select, Start bool
}

var Joystick1 = &Joystick{}
var Joystick2 = &Joystick{}
var joysticks = []*Joystick{Joystick1, Joystick2}

// ******************************************************
// ******************************************************
// ******************************************************

func drawText(text string, img *glImage, color color.RGBA, aling LabelAlign) {
	face := inconsolata.Regular8x16
	r, _ := font.BoundString(face, text)

	w, h := r.Max.X.Ceil(), r.Min.Y.Ceil() //face.Height

	y := img.rect.Dy()/2 - h/2
	x := 0
	switch aling {
	case Center:
		x = img.rect.Dx()/2 - w/2
	case Right:
		x = img.rect.Dx() - w
	}
	p := fixed.P(int(x), int(y))

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color),
		Face: face,
		Dot:  p,
	}

	d.DrawString(text)
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
