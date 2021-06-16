package gui

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/go-gl/gl/v3.3-core/gl"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/inconsolata"
)

type LabelAlign int

const (
	Left LabelAlign = iota
	Center
	Right
)

type Label interface {
	GUIObject
	SetText(txt string)
	GetText() string
	SetForeground(c color.RGBA)
}

type label struct {
	rect  Rect
	text  string
	back  *image.Uniform
	fore  color.RGBA
	face  *basicfont.Face
	aling LabelAlign

	needsUpdate bool
	img         *glImage // TODO: replace by our own image

	texture uint32
	frameID uint32
}

func NewLabel(txt string, aling LabelAlign) Label {
	l := &label{
		back:  image.NewUniform(color.RGBA{255, 255, 255, 255}),
		fore:  color.RGBA{0, 0, 0, 255},
		face:  inconsolata.Regular8x16,
		aling: aling,
	}

	l.init()
	l.SetText(txt)
	return l
}

func (*label) GetChildrens() []GUIObject {
	return []GUIObject{}
}

func (l *label) SetForeground(c color.RGBA) {
	l.fore = c
	l.redraw()
}

func (l *label) GetText() string {
	return l.text
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
	drawText(l.text, l.img, l.fore, l.aling)
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
	// println("l:", l.text, "t:", l.texture, "f:", l.frameID)
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
