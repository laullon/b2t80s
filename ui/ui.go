package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Control interface {
	Widget() fyne.CanvasObject
	HTML() string
	Update()
}

var App fyne.App

type RegText struct {
	Label *canvas.Text
	Value *canvas.Text
}

func NewRegText(label string) *RegText {
	rt := &RegText{
		Label: &canvas.Text{},
		Value: &canvas.Text{},
	}
	rt.Label.Text = label
	rt.Label.Color = color.Black
	rt.Label.TextSize = fyne.CurrentApp().Settings().Theme().Size("text")
	rt.Label.TextStyle = fyne.TextStyle{Monospace: true}
	rt.Label.Alignment = fyne.TextAlignTrailing

	rt.Value.Color = color.Black
	rt.Value.TextSize = fyne.CurrentApp().Settings().Theme().Size("text")
	rt.Value.TextStyle = fyne.TextStyle{Monospace: true}

	return rt
}

func (rt *RegText) Update(text string) {
	if rt.Value.Text != text {
		rt.Value.Text = text
		rt.Value.Color = color.RGBA{0x00, 0x00, 0xff, 0xff}
	} else {
		rt.Value.Color = color.Black
	}
}
