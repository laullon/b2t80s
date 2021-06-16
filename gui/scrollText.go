package gui

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type ScrollText interface {
	Text
	GUIObject
	ScrollTarget
}

type scrollText struct {
	text *text
	bar  *scrollBar
	ui   HCT
	rect Rect
}

func NewScrollText() ScrollText {
	sc := &scrollText{}
	sc.text = NewText("").(*text)

	sc.bar = &scrollBar{text: sc.text}

	sc.ui = NewHorizontalHCT()
	sc.ui.SetCenter(sc.text)
	sc.ui.SetTail(sc.bar, 20)

	return sc
}

func (sc *scrollText) OnScroll(x, y int32) {
	sc.text.startLine += y
	sc.text.redraw()
	sc.bar.redraw()
}

func (sc *scrollText) Resize(r Rect) {
	sc.rect = r
	sc.ui.Resize(r)
}

func (sc *scrollText) GetChildrens() []GUIObject {
	return append([]GUIObject{sc.bar}, sc.ui.GetChildrens()...)
}

func (sc *scrollText) SetText(txt string)          { sc.text.SetText(txt) }
func (sc *scrollText) SetForeground(c color.Color) { sc.text.SetForeground(c) }

func (sc *scrollText) Render() { sc.ui.Render() }

func (sc *scrollText) Rect() Rect        { return sc.rect }
func (sc *scrollText) OnMouseOver(bool)  {}
func (sc *scrollText) OnMouseClick(bool) {}

// *********

type scrollBar struct {
	text        *text
	rect        Rect
	img         *glImage
	needsUpdate bool
	texture     uint32
	frameID     uint32
}

func (sc *scrollBar) GetChildrens() []GUIObject {
	return []GUIObject{}
}

func (b *scrollBar) Resize(r Rect) {
	b.rect = r
	b.init()
	b.redraw()
}

func (b *scrollBar) redraw() {
	w, h := int(b.rect.W), int(b.rect.H)
	lineH := float32(h) / float32(b.text.nLines)
	barY := lineH * float32(b.text.startLine)
	barH := lineH * float32(b.text.nLinesVisible)

	c := &scrollBarRender{size: Size{b.rect.W, b.rect.H}, fill: Rect{0, int32(barY), b.rect.W, int32(barH)}}
	b.img = newImage(Size{b.rect.W, b.rect.H})

	draw.Draw(b.img, image.Rect(0, 0, w, h), c, image.Point{0, 0}, draw.Src)

	b.needsUpdate = true
}

type scrollBarRender struct {
	size Size
	fill Rect
}

func (*scrollBarRender) ColorModel() color.Model {
	return color.AlphaModel
}

func (p *scrollBarRender) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(p.size.W), int(p.size.H))
}

func (p *scrollBarRender) At(x, y int) color.Color {
	if p.fill.In(Point{X: int32(x), Y: int32(y)}) {
		return color.Black
	} else if x%2 == y%2 {
		return color.Black
	}
	return color.White
}

func (b *scrollBar) init() {
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

func (b *scrollBar) Render() {
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
