package ui

import (
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/gui"
)

type m6502BusUI struct {
	bus      m6502.Bus
	selected string
}

func NewM6502BusUI(name string, bus m6502.Bus) gui.GUIObject {
	ctrl := &m6502BusUI{
		bus: bus,
	}

	if len(name) > 0 {
		name += " "
	}

	// show := widget.NewButton(name+"mem.", ctrl.doShow)

	// ctrl.widget = container.New(layout.NewHBoxLayout(),
	// 	widget.NewToolbarSeparator().ToolbarObject(),
	// 	show,
	// )

	return ctrl
}

func (ui *m6502BusUI) GetChildrens() []gui.GUIObject {
	return []gui.GUIObject{}
}

func (ui *m6502BusUI) Render() {
}

func (ui *m6502BusUI) Resize(r gui.Rect) {
}

func (ui *m6502BusUI) GetRegisters() string { return "" }
func (ui *m6502BusUI) GetOutput() string    { return "" }

func (ui *m6502BusUI) Update() {
}

func (ui *m6502BusUI) doShow() {
	// if ui.win == nil {
	// 	ui.text = &widget.Label{}
	// 	// dump.Color = color.Black
	// 	// dump.TextSize = fyne.CurrentApp().Settings().Theme().Size("text")
	// 	ui.text.TextStyle = fyne.TextStyle{Monospace: true}

	// 	keys := []string{}
	// 	for k := range ui.bus.GetDumplables() {
	// 		keys = append(keys, k)
	// 	}

	// 	selector := widget.NewSelect(keys, ui.dumplableChanged)
	// 	if len(keys) > 0 {
	// 		selector.SetSelected(keys[0])
	// 		ui.selected = keys[0]
	// 	}

	// container := container.New(layout.NewBorderLayout(selector, nil, nil, nil), selector, container.NewVScroll(ui.text))

	// ui.win = App.NewWindow("Memory")
	// ui.win.SetContent(container)
	// ui.win.Show()

	// wait := time.Duration(3 * time.Second)
	// ticker := time.NewTicker(wait)
	// go func() {
	// 	for range ticker.C {
	// 		if len(ui.selected) != 0 {
	// 			ui.text.Text = hex.Dump(ui.bus.GetDumplables()[ui.selected].Memory())
	// 			container.Refresh()
	// 		}
	// 	}
	// }()
}
func (ui *m6502BusUI) dumplableChanged(name string) {
	ui.selected = name
	// ui.text.Text = hex.Dump(ui.bus.GetDumplables()[ui.selected].Memory())
	// ui.widget.Refresh()
}
