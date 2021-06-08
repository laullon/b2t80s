package ui

import (
	"encoding/hex"

	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/gui"
)

type busUI struct {
	selected int
	names    []string
	text     gui.Text
	dumpable map[string]cpu.Dumpable
	tabs     gui.Tabs
}

func NewBusUI(bus cpu.Bus) gui.GUIObject {
	ctrl := &busUI{
		dumpable: bus.GetDumplables(),
		tabs:     gui.NewTabs(),
	}

	ctrl.text = gui.NewScrollText()
	for name := range ctrl.dumpable {
		ctrl.names = append(ctrl.names, name)
		ctrl.tabs.AddTabs(name, ctrl.text)
	}

	ctrl.tabs.SetOnChange(func(i int) {
		ctrl.selected = i
		ctrl.Update()
	})

	return ctrl
}

func (ctl *busUI) GetMouseTargets() []gui.MouseTarget {
	return ctl.tabs.GetMouseTargets()
}

func (ctl *busUI) Render() {
	ctl.tabs.Render()
}

func (ctl *busUI) Resize(r gui.Rect) {
	ctl.tabs.Resize(r)
}

func (ctrl *busUI) Update() {
	dump := hex.Dump(ctrl.dumpable[ctrl.names[ctrl.selected]].Memory())
	ctrl.text.SetText(dump)
}
