package ui

type VolumenControl struct {
}

func NewVolumenControl(setVolume func(float64)) *VolumenControl {
	vc := &VolumenControl{}

	// vc.sel = widget.NewButtonWithIcon("", fyne.NewStaticResource("pp", data.MustAsset("data/icons/volume-control-full.png")), vc.do)
	// vc.slider = widget.NewSlider(0, 1)
	// vc.slider.Step = 0.05
	// vc.slider.OnChanged = setVolume

	// vc.ui = container.New(layout.NewHBoxLayout(),
	// 	widget.NewToolbarSeparator().ToolbarObject(),
	// 	vc.sel,
	// 	vc.slider,
	// )

	return vc
}

func (ui *VolumenControl) GetRegisters() string { return "" }
func (ui *VolumenControl) GetOutput() string    { return "" }

// func (vc *VolumenControl) Widget() fyne.CanvasObject {
// 	return vc.ui
// }

func (vc *VolumenControl) do() {
}

func (vc *VolumenControl) Update() {
}
