package emulator

import (
	"fmt"

	"github.com/laullon/webview"
)

var tabHTML = `<button id="tab_%s" onclick="tabsSelect(\'%s\');">%s</button>`

type Tabs interface {
	Show()
	Selected() string
}

type tabs struct {
	div      string
	web      webview.WebView
	machine  Machine
	selected string
}

func NewTabs(div string, web webview.WebView, machine Machine) Tabs {
	tabs := &tabs{div: div, web: web, machine: machine}

	web.Bind("tabsSelect", tabs.changeSelected)

	return tabs
}

func (tabs *tabs) Show() {
	i := 0 // TODO: Remove this
	keys := make([]string, len(tabs.machine.Control()))
	for k := range tabs.machine.Control() {
		keys[i] = k
		i++
	}

	for _, name := range keys {
		html := fmt.Sprintf(tabHTML, name, name, name)
		println(html)
		tabs.web.Eval(fmt.Sprintf("document.getElementById('%s').innerHTML += '%s'", tabs.div, html))
	}
	tabs.changeSelected(keys[0])
}

func (tabs *tabs) Selected() string {
	return tabs.selected
}

func (tabs *tabs) changeSelected(tab string) {
	if len(tabs.selected) > 0 {
		tabs.web.Eval(fmt.Sprintf("document.getElementById('tab_%s').classList.remove('active');", tabs.selected))
	}
	tabs.selected = tab
	tabs.web.Eval(fmt.Sprintf("document.getElementById('tab_%s').classList.add('active');", tabs.selected))
}
