package ui

import (
	"encoding/hex"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/laullon/b2t80s/cpu/m6502"
)

type m6502BusUI struct {
	widget   *fyne.Container
	bus      m6502.Bus
	win      fyne.Window
	selected string
}

func NewM6502BusUI(bus m6502.Bus) Control {
	ctrl := &m6502BusUI{
		bus: bus,
	}

	show := widget.NewButton("memory", ctrl.doShow)

	ctrl.widget = container.New(layout.NewHBoxLayout(),
		widget.NewToolbarSeparator().ToolbarObject(),
		show,
	)

	return ctrl
}

func (ui *m6502BusUI) Widget() fyne.CanvasObject {
	return ui.widget
}

func (ui *m6502BusUI) Update() {
}

func (ui *m6502BusUI) doShow() {
	if ui.win == nil {
		dump := &widget.Label{}
		// dump.Color = color.Black
		// dump.TextSize = fyne.CurrentApp().Settings().Theme().Size("text")
		dump.TextStyle = fyne.TextStyle{Monospace: true}

		keys := []string{}
		for k := range ui.bus.GetDumplables() {
			keys = append(keys, k)
		}

		selector := widget.NewSelect(keys, ui.dumplableChanged)
		if len(keys) > 0 {
			selector.SetSelected(keys[0])
			ui.selected = keys[0]
		}

		container := container.New(layout.NewBorderLayout(selector, nil, nil, nil), selector, container.NewVScroll(dump))

		ui.win = App.NewWindow("Memory")
		ui.win.SetContent(container)
		ui.win.Show()

		wait := time.Duration(3 * time.Second)
		ticker := time.NewTicker(wait)
		go func() {
			for range ticker.C {
				if len(ui.selected) != 0 {
					dump.Text = hex.Dump(ui.bus.GetDumplables()[ui.selected].Memory())
					container.Refresh()
				}
			}
		}()
	}
}
func (ui *m6502BusUI) dumplableChanged(name string) {
	ui.selected = name
}
