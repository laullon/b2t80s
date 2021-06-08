package gui

type Tabs interface {
	GUIObject
	AddTabs(name string, panel GUIObject)
	SetOnChange(func(int))
}

type tabs struct {
	ui       HCT
	tabs     []Button
	panels   []GUIObject
	bar      Grid
	rect     Rect
	onChange func(int)
}

func NewTabs() Tabs {
	tabs := &tabs{}
	tabs.bar = NewHGrid(6, 50, 0)
	tabs.ui = NewVerticalHCT()
	tabs.ui.SetHead(tabs.bar, 50)
	return tabs
}

func (tabs *tabs) AddTabs(name string, panel GUIObject) {
	tabID := len(tabs.tabs)
	bt := NewTab(name)
	bt.SetAction(func() { tabs.setTab(tabID) })
	tabs.tabs = append(tabs.tabs, bt.(*button))
	tabs.bar.Add(bt)
	tabs.panels = append(tabs.panels, panel)
	if tabID == 0 {
		tabs.setTab(tabID)
	}
}

func (tabs *tabs) GetMouseTargets() []MouseTarget {
	var res []MouseTarget
	for _, obj := range tabs.tabs {
		res = append(res, obj.GetMouseTargets()...)
	}
	for _, obj := range tabs.panels {
		res = append(res, obj.GetMouseTargets()...)
	}
	return res
}

func (tabs *tabs) Resize(r Rect) {
	tabs.ui.Resize(r)
	tabs.rect = r
}

func (tabs *tabs) setTab(tabIdx int) {
	for idx, tab := range tabs.tabs {
		tab.(*button).active = idx == tabIdx
	}
	tabs.ui.SetCenter(tabs.panels[tabIdx])
	tabs.ui.Resize(tabs.rect)
	if tabs.onChange != nil {
		tabs.onChange(tabIdx)
	}
}

func (tabs *tabs) SetOnChange(onChange func(int)) {
	tabs.onChange = onChange
}

func (tabs *tabs) Render() { tabs.ui.Render() }
