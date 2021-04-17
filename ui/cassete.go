package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
)

type CasseteControl struct {
	ui       *fyne.Container
	play     *widget.Button
	stop     *widget.Button
	sel      *widget.Button
	cassette cassette.Cassette
}

func NewCasseteControl(cassette cassette.Cassette, disable bool) *CasseteControl {
	cas := &CasseteControl{cassette: cassette}

	cas.play = widget.NewButtonWithIcon("", fyne.NewStaticResource("pp", data.MustAsset("data/icons/controls-play.png")), cas.doPlay)
	cas.stop = widget.NewButtonWithIcon("", fyne.NewStaticResource("pp", data.MustAsset("data/icons/controls-stop.png")), cas.doStop)
	cas.sel = widget.NewButtonWithIcon("", fyne.NewStaticResource("pp", data.MustAsset("data/icons/cassette.png")), cas.doName)

	cas.ui = container.New(layout.NewHBoxLayout(),
		widget.NewToolbarSeparator().ToolbarObject(),
		cas.sel,
	)

	if !disable {
		cas.ui.Add(cas.play)
		cas.ui.Add(cas.stop)
	}

	cas.Update()
	return cas
}

func (ui *CasseteControl) HTML() string { return "" }

func (cas *CasseteControl) Widget() fyne.CanvasObject {
	return cas.ui
}

func (cas *CasseteControl) doName() {
	c := fyne.CurrentApp().Driver().CanvasForObject(cas.ui)
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(cas.ui)
	var pop *widget.PopUp
	ok := widget.NewLabel(cas.cassette.Name())
	pop = widget.NewPopUp(ok, c)
	pop.Move(pos)
}

func (cas *CasseteControl) doPlay() {
	if !cas.cassette.IsMotorON() {
		println("play")
		cas.cassette.Motor(true)
	}
	cas.Update()
}

func (cas *CasseteControl) doStop() {
	if cas.cassette.IsMotorON() {
		println("stop")
		cas.cassette.Motor(false)
	}
	cas.Update()
}

func (cas *CasseteControl) Update() {
	if cas.cassette.IsMotorON() {
		cas.play.Disable()
		cas.stop.Enable()
	} else {
		cas.play.Enable()
		cas.stop.Disable()
	}
	// cas.sel.Disable()
}
