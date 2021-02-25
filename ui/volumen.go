package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	widget "fyne.io/fyne/v2/widget"
	"github.com/laullon/b2t80s/data"
)

type VolumenControl struct {
	ui     fyne.CanvasObject
	sel    *widget.Button
	slider *widget.Slider
}

func NewVolumenControl(setVolume func(float64)) *VolumenControl {
	vc := &VolumenControl{}

	vc.sel = widget.NewButtonWithIcon("", fyne.NewStaticResource("pp", data.MustAsset("data/icons/volume-control-full.png")), vc.do)
	vc.slider = widget.NewSlider(0, 1)
	vc.slider.Step = 0.05
	vc.slider.OnChanged = setVolume

	vc.ui = container.New(layout.NewHBoxLayout(),
		widget.NewToolbarSeparator().ToolbarObject(),
		vc.sel,
		vc.slider,
	)

	return vc
}

func (vc *VolumenControl) Widget() fyne.CanvasObject {
	return vc.ui
}

func (vc *VolumenControl) do() {
}

func (vc *VolumenControl) Update() {
}
