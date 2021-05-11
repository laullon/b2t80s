package emulator

import (
	"C"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/laullon/b2t80s/debug"
	"github.com/laullon/webview"
)

type debugWindow struct {
	web  webview.WebView
	tabs Tabs
}

func NewDebugWindow(name string, machine Machine) Window {
	win := &debugWindow{
		web: initDebugWindow(name, machine),
	}

	win.web.Bind("getStatus", func() string {
		return fmt.Sprintf("time: %s - FPS: %03.2f\n", machine.Clock().Stats(), machine.Monitor().FPS())
	})

	win.web.Bind("getRegisters", func() string {
		ui := machine.Control()[win.tabs.Selected()]
		return ui.GetRegisters()
	})

	win.web.Bind("getOutput", func() string {
		ui := machine.Control()[win.tabs.Selected()]
		return ui.GetOutput()
	})

	win.web.Bind("initUI", func() {
		win.tabs.Show()
	})

	http.Handle("/cmd/", win.web)
	http.Handle("/", http.FileServer(debug.AssetFile()))
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	go func() {
		panic(http.Serve(listener, nil))
	}()

	fmt.Println("Using port:", listener.Addr())
	url := fmt.Sprintf("http://0:%d/app/?%d", listener.Addr().(*net.TCPAddr).Port, time.Now().Local().Second())
	println("url:", url)
	win.web.Navigate(url)

	win.tabs = NewTabs("tabs", win.web, machine)

	return win
}

func (win *debugWindow) SetOnKey(onKey func(glfw.Key)) {
	// win.main.onKey = onKey
}

func (win *debugWindow) Run() {
	win.web.Run()
	win.web.Destroy()
}

func initDebugWindow(title string, machine Machine) webview.WebView {
	w := webview.New(title, 1200, 600)
	return w
}
