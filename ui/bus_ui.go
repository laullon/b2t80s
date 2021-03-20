package ui

import (
	"encoding/hex"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/laullon/b2t80s/cpu"
)

type busUI struct {
	widget   *fyne.Container
	bus      cpu.Bus
	win      fyne.Window
	text     *widget.Label
	selected string
}

func NewBusUI(name string, bus cpu.Bus) Control {
	ctrl := &busUI{
		bus: bus,
	}

	if len(name) > 0 {
		name += " "
	}

	show := widget.NewButton(name+"mem.", ctrl.doShow)

	ctrl.widget = container.New(layout.NewHBoxLayout(),
		widget.NewToolbarSeparator().ToolbarObject(),
		show,
	)

	return ctrl
}

func (ui *busUI) Widget() fyne.CanvasObject {
	return ui.widget
}

func (ui *busUI) Update() {
}

func (ui *busUI) doShow() {
	if ui.win == nil {
		ui.text = &widget.Label{}
		// dump.Color = color.Black
		// dump.TextSize = fyne.CurrentApp().Settings().Theme().Size("text")
		ui.text.TextStyle = fyne.TextStyle{Monospace: true}

		keys := []string{}
		for k := range ui.bus.GetDumplables() {
			keys = append(keys, k)
		}

		selector := widget.NewSelect(keys, ui.dumplableChanged)
		if len(keys) > 0 {
			selector.SetSelected(keys[0])
			ui.selected = keys[0]
		}

		container := container.New(layout.NewBorderLayout(selector, nil, nil, nil), selector, container.NewVScroll(ui.text))

		ui.win = App.NewWindow("Memory")
		ui.win.SetContent(container)
		ui.win.Show()

		wait := time.Duration(3 * time.Second)
		ticker := time.NewTicker(wait)
		go func() {
			for range ticker.C {
				if len(ui.selected) != 0 {
					ui.text.Text = hex.Dump(ui.bus.GetDumplables()[ui.selected].Memory())
					container.Refresh()
				}
			}
		}()
	}
}
func (ui *busUI) dumplableChanged(name string) {
	ui.selected = name
	ui.text.Text = hex.Dump(ui.bus.GetDumplables()[ui.selected].Memory())
	ui.widget.Refresh()
}
