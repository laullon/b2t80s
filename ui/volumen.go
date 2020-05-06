package ui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
)

type VolumenControl struct {
	ui     fyne.CanvasObject
	sel    *widget.Button
	slider *widget.Slider
	sound  emulator.SoundSystem
}

func NewVolumenControl(sound emulator.SoundSystem) *VolumenControl {
	vc := &VolumenControl{sound: sound}

	vc.sel = widget.NewButtonWithIcon("", fyne.NewStaticResource("pp", data.MustAsset("data/icons/volume-control-full.png")), vc.do)
	vc.slider = widget.NewSlider(0, 100)
	vc.slider.OnChanged = vc.sound.SetVolume

	vc.ui = widget.NewHBox(
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
