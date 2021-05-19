package gui

import (
	"image"
	"image/color"
)

type Button interface {
	GUIObject
	MouseTarget
}

type button struct {
	label    *label
	over     bool
	selected bool
}

func NewButton(txt string, rect Rect) Button {
	b := &button{
		label: NewLabel(txt, rect).(*label),
	}

	return b
}

func (b *button) OnMouseOver(over bool) {
	if b.over != over {
		b.over = over
		b.redraw()
	}
}

func (b *button) OnMouseClick(up bool) {
	if up {
		if b.selected {
			println("click")
		}
		b.selected = false
	} else {
		b.selected = true
	}
	b.redraw()
}

func (b *button) redraw() {
	if b.over {
		if b.selected {
			b.label.back = image.NewUniform(color.RGBA{0, 255, 0, 255})
		} else {
			b.label.back = image.NewUniform(color.RGBA{0, 0, 255, 255})
		}
	} else {
		b.label.back = image.NewUniform(color.RGBA{255, 255, 255, 255})
		b.selected = false
	}
	b.label.redraw()
}

func (b *button) Rect() Rect    { return b.label.rect }
func (b *button) Render()       { b.label.Render() }
func (b *button) Resize(r Rect) { b.label.Resize(r) }
