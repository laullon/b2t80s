package gui

type Tabs interface {
	GUIObject
	Tabs() []MouseTarget
	AddTabs(name string, panel GUIObject)
}

type tabs struct {
	ui     HCT
	tabs   []MouseTarget
	panels []GUIObject
	bar    Grid
	rect   Rect
}

func NewTabs() Tabs {
	tabs := &tabs{}
	tabs.bar = NewHGrid(6, 50)
	tabs.ui = NewVerticalHCT()
	tabs.ui.SetHead(tabs.bar, 50)
	return tabs
}

func (tabs *tabs) AddTabs(name string, panel GUIObject) {
	tabID := len(tabs.tabs)
	bt := NewButton(name)
	bt.SetAction(func() { tabs.setTab(tabID) })
	tabs.tabs = append(tabs.tabs, bt)
	tabs.bar.Add(bt)
	tabs.panels = append(tabs.panels, panel)
	if tabID == 0 {
		tabs.setTab(tabID)
	}
}

func (tabs *tabs) Resize(r Rect) {
	tabs.ui.Resize(r)
	tabs.rect = r
}

func (tabs *tabs) setTab(tab int) {
	tabs.ui.SetCenter(tabs.panels[tab])
	tabs.ui.Resize(tabs.rect)
}

func (tabs *tabs) Render()             { tabs.ui.Render() }
func (tabs *tabs) Tabs() []MouseTarget { return tabs.tabs }
