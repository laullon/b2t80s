package ui

import (
	"image/color"

	"github.com/laullon/b2t80s/gui"
)

type RegText struct {
	Label gui.Label
	Value gui.Label
}

func NewRegText(label string) *RegText {
	rt := &RegText{
		Label: gui.NewLabel(label, gui.Right),
		Value: gui.NewLabel("", gui.Left),
	}
	return rt
}

func (rt *RegText) Update(text string) {
	if rt.Value.GetText() != text {
		rt.Value.SetText(text)
		rt.Value.SetForeground(color.RGBA{0, 0, 0xff, 0xff})
	} else {
		rt.Value.SetForeground(color.RGBA{0, 0, 0, 0xff})
	}
}
