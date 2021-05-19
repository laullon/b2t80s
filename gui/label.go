package gui

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/go-gl/gl/all-core/gl"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

type Label interface {
	GUIObject
	SetText(txt string)
}

type label struct {
	rect Rect
	text string
	back *image.Uniform
	fore *image.Uniform
	face *basicfont.Face

	needsUpdate bool
	img         *glImage // TODO: replace by our own image

	texture uint32
	frameID uint32
}

func NewLabel(txt string, rect Rect) Label {
	l := &label{
		rect: rect,
		back: image.NewUniform(color.RGBA{255, 255, 255, 255}),
		fore: image.NewUniform(color.RGBA{0, 0, 0, 255}),
		face: inconsolata.Regular8x16,
	}

	l.init()
	l.SetText(txt)
	return l
}

func (l *label) SetText(txt string) {
	l.text = txt
	l.redraw()
}

func (l *label) Resize(r Rect) {
	l.rect = r
	l.init()
	l.redraw()
}

func (l *label) redraw() {

	draw.Draw(l.img, l.img.Bounds(), l.back, image.Point{}, draw.Src)

	r, _ := font.BoundString(l.face, l.text)
	y := l.rect.H/2 - int32(r.Min.Y.Ceil())/2
	x := l.rect.W/2 - int32(r.Max.X.Ceil())/2
	p := fixed.P(int(x), int(y))
	// println("int32(r.Min.Y.Ceil())", int32(r.Min.Y.Ceil()))
	// fmt.Printf("r:%v\n", r)
	// fmt.Printf("r:%v\n", l.rect)
	// println(l.rect.H, "+", int32(r.Min.Y.Ceil()))
	// fmt.Printf("p:%v\n", p)

	d := &font.Drawer{
		Dst:  l.img,
		Src:  l.fore,
		Face: l.face,
		Dot:  p,
	}

	d.DrawString(l.text)
	l.needsUpdate = true
}

func (l *label) init() {
	l.img = newImage(Size{l.rect.W, l.rect.H})

	if l.texture != 0 {
		gl.DeleteTextures(1, &l.texture)
	}
	// Texture
	gl.GenTextures(1, &l.texture)
	gl.BindTexture(gl.TEXTURE_2D, l.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB,
		l.rect.W, l.rect.H,
		0, gl.RGB, gl.UNSIGNED_BYTE,
		nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// BUFFER
	if l.frameID != 0 {
		gl.DeleteFramebuffers(1, &l.frameID)
	}
	gl.GenFramebuffers(1, &l.frameID)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, l.frameID)
	gl.FramebufferTexture2D(gl.READ_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, l.texture, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
}

func (l *label) Render() {
	// UPDATE TEXTURE
	if l.needsUpdate {
		l.needsUpdate = false
		gl.BindTexture(gl.TEXTURE_2D, l.texture)
		gl.TexSubImage2D(gl.TEXTURE_2D, 0,
			0, 0, l.rect.W, l.rect.H,
			gl.RGBA, gl.UNSIGNED_BYTE,
			gl.Ptr(l.img.Pix))
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	}

	// RENDER
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, l.frameID)
	gl.EnableVertexAttribArray(0)
	gl.BlitFramebuffer(
		0, 0, l.rect.W, l.rect.H,
		l.rect.X, l.rect.Y, l.rect.W+l.rect.X, l.rect.H+l.rect.Y,
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)
}
