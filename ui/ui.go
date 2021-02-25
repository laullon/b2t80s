package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Control interface {
	Widget() fyne.CanvasObject
	Update()
}

type regText struct {
	txt *canvas.Text
}

func NewRegText(txt string) *regText {
	rt := &regText{txt: &canvas.Text{}}
	rt.txt.Text = txt
	rt.txt.Color = color.Black
	rt.txt.TextSize = fyne.CurrentApp().Settings().Theme().Size("text")
	rt.txt.TextStyle = fyne.TextStyle{Monospace: true}
	return rt
}

func (rt *regText) update(text string) {
	if rt.txt.Text != text {
		rt.txt.Text = text
		rt.txt.Color = color.RGBA{0x00, 0x00, 0xff, 0xff}
	} else {
		rt.txt.Color = color.Black
	}
}
