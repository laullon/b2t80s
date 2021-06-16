package gui

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Button interface {
	GUIObject
	MouseTarget
	SetAction(f func())
}

type button struct {
	rect    Rect
	txt     string
	over    bool
	clicked bool
	active  bool
	action  func()

	needsUpdate bool
	img         *glImage
	back        color.RGBA

	texture uint32
	frameID uint32

	tab bool
}

func NewButton(txt string) Button {
	return newButton(txt, false)
}

func NewTab(txt string) Button {
	return newButton(txt, true)
}

func newButton(txt string, tab bool) Button {
	b := &button{
		txt: txt,
		tab: tab,
	}

	b.redraw()
	b.init()
	return b
}

func (b *button) GetChildrens() []GUIObject {
	return []GUIObject{}
}

func (b *button) OnMouseOver(over bool) {
	if b.over != over {
		b.over = over
		b.redraw()
	}
}

func (b *button) OnMouseClick(up bool) {
	if up {
		if b.clicked {
			if b.action != nil {
				b.action()
			}
		}
		b.clicked = false
	} else {
		b.clicked = true
	}
	b.redraw()
}

func (b *button) redraw() {
	if b.active {
		b.back = palette[3]
	} else if b.over {
		if b.clicked {
			b.back = palette[3]
		} else {
			b.back = palette[2]
		}
	} else {
		b.back = palette[0]
	}

	w, h := int(b.rect.W), int(b.rect.H)
	c := &circle{p: image.Point{6, 6}, r: 6, r2: 5, in: b.back, line: palette[1], out: palette[0]}
	b.img = newImage(Size{b.rect.W, b.rect.H})

	draw.Draw(b.img, image.Rect(0, 0, w, h), image.NewUniform(palette[1]), image.Point{}, draw.Src)
	draw.Draw(b.img, image.Rect(1, 1, w-1, h-1), image.NewUniform(b.back), image.Point{}, draw.Src)

	draw.Draw(b.img, image.Rect(0, 0, 6, 6), c, image.Point{0, 0}, draw.Src)
	draw.Draw(b.img, image.Rect(w-6, 0, w, 6), c, image.Point{6, 0}, draw.Src)
	if !(b.tab) {
		draw.Draw(b.img, image.Rect(0, h-6, 6, h), c, image.Point{0, 6}, draw.Src)
		draw.Draw(b.img, image.Rect(w-6, h-6, w, h), c, image.Point{6, 6}, draw.Src)
	}

	drawText(b.txt, b.img, color.RGBA{0, 0, 0, 255}, Center)

	b.needsUpdate = true

}

func (b *button) Rect() Rect         { return b.rect }
func (b *button) SetAction(f func()) { b.action = f }

func (b *button) Resize(r Rect) {
	b.rect = r
	b.init()
	b.redraw()
}

func (b *button) init() {
	if b.texture != 0 {
		gl.DeleteTextures(1, &b.texture)
	}
	// Texture
	gl.GenTextures(1, &b.texture)
	gl.BindTexture(gl.TEXTURE_2D, b.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB,
		b.rect.W, b.rect.H,
		0, gl.RGB, gl.UNSIGNED_BYTE,
		nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	b.needsUpdate = true

	// BUFFER
	if b.frameID != 0 {
		gl.DeleteFramebuffers(1, &b.frameID)
	}
	gl.GenFramebuffers(1, &b.frameID)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, b.frameID)
	gl.FramebufferTexture2D(gl.READ_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, b.texture, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
}

func (b *button) Render() {
	// println("l:", l.text, "t:", l.texture, "f:", l.frameID)
	// UPDATE TEXTURE
	if b.needsUpdate {
		b.needsUpdate = false
		gl.BindTexture(gl.TEXTURE_2D, b.texture)
		gl.TexSubImage2D(gl.TEXTURE_2D, 0,
			0, 0, b.rect.W, b.rect.H,
			gl.RGBA, gl.UNSIGNED_BYTE,
			gl.Ptr(b.img.Pix))
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	}

	// RENDER
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, b.frameID)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BlitFramebuffer(
		0, 0, b.rect.W, b.rect.H,
		b.rect.X, b.rect.Y, b.rect.W+b.rect.X, b.rect.H+b.rect.Y,
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)

	// b.label.Render()
}
