package gui

import (
	"image"
	"image/color"
	"image/draw"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

type Text interface {
	GUIObject
	SetText(txt string)
	SetForeground(c color.Color)
}

type text struct {
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

func NewText(t string) Text {
	txt := &text{
		back: image.NewUniform(color.RGBA{255, 255, 255, 255}),
		fore: image.NewUniform(color.RGBA{0, 0, 0, 255}),
		face: inconsolata.Regular8x16,
	}

	txt.init()
	txt.SetText(t)
	return txt
}

func (txt *text) SetForeground(c color.Color) {
	txt.fore = image.NewUniform(c)
	txt.redraw()
}

func (txt *text) SetText(text string) {
	txt.text = text
	txt.redraw()
}

func (txt *text) Resize(r Rect) {
	txt.rect = r
	txt.init()
	txt.redraw()
}

func (txt *text) redraw() {
	draw.Draw(txt.img, txt.img.Bounds(), txt.back, image.Point{}, draw.Src)

	p := fixed.P(0, 0)
	for _, line := range strings.Split(txt.text, "\n") {
		p.Y += txt.face.Metrics().Height
		d := &font.Drawer{
			Dst:  txt.img,
			Src:  txt.fore,
			Face: txt.face,
			Dot:  p,
		}

		d.DrawString(line)
	}
	txt.needsUpdate = true
}

func (txt *text) init() {
	txt.img = newImage(Size{txt.rect.W, txt.rect.H})

	if txt.texture != 0 {
		gl.DeleteTextures(1, &txt.texture)
	}
	// Texture
	gl.GenTextures(1, &txt.texture)
	gl.BindTexture(gl.TEXTURE_2D, txt.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB,
		txt.rect.W, txt.rect.H,
		0, gl.RGB, gl.UNSIGNED_BYTE,
		nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// BUFFER
	if txt.frameID != 0 {
		gl.DeleteFramebuffers(1, &txt.frameID)
	}
	gl.GenFramebuffers(1, &txt.frameID)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, txt.frameID)
	gl.FramebufferTexture2D(gl.READ_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, txt.texture, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
}

func (txt *text) Render() {
	// println("l:", txt.text, "t:", txt.texture, "f:", txt.frameID)
	// UPDATE TEXTURE
	if txt.needsUpdate {
		txt.needsUpdate = false
		gl.BindTexture(gl.TEXTURE_2D, txt.texture)
		gl.TexSubImage2D(gl.TEXTURE_2D, 0,
			0, 0, txt.rect.W, txt.rect.H,
			gl.RGBA, gl.UNSIGNED_BYTE,
			gl.Ptr(txt.img.Pix))
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	}

	// RENDER
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, txt.frameID)
	gl.EnableVertexAttribArray(0)
	gl.BlitFramebuffer(
		0, 0, txt.rect.W, txt.rect.H,
		txt.rect.X, txt.rect.Y, txt.rect.W+txt.rect.X, txt.rect.H+txt.rect.Y,
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)
}
