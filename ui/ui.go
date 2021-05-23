package ui

import (
	"image/color"

	"github.com/laullon/b2t80s/gui"
)

type RegText struct {
	label gui.Label
	value gui.Label
}

func NewRegText(label string) *RegText {
	rt := &RegText{
		label: gui.NewLabel(label),
		value: gui.NewLabel(""),
	}
	return rt
}

func (rt *RegText) Update(text string) {
	if rt.value.GetText() != text {
		rt.value.SetText(text)
		rt.value.SetForeground(color.RGBA{0, 0, 0xff, 0xff})
	} else {
		rt.value.SetForeground(color.RGBA{0, 0, 0, 0xff})
	}
}
