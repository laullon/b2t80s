package emulator

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/laullon/b2t80s/debug"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/webview/webview"
)

type debugWindow struct {
	window *sdl.Window
	web    webview.WebView
	tabs   Tabs
}

func NewDebugWindow(name string, machine Machine) Window {
	win := &debugWindow{}

	window, err := sdl.CreateWindow("Debug", 850, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_RESIZABLE|sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	win.window = window

	wmInfo, err := window.GetWMInfo()
	if err != nil {
		panic(err)
	}

	win.web = webview.NewWindow(true, wmInfo.GetWindowsInfo().Window)

	window.UpdateSurface()

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

func (win *debugWindow) SetOnKey(func(sdl.Scancode)) {
}

func (win *debugWindow) Run() {
	win.web.Run()
	win.web.Destroy()
}
