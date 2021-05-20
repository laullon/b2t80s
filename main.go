//go:generate $HOME/go/bin/go-bindata -pkg data -o data/data.go data/...
//go:generate $HOME/go/bin/go-bindata -fs -prefix debug -pkg debug -o debug/data.go debug/...

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/gui"
	"github.com/laullon/b2t80s/machines/atetris"
	"github.com/laullon/b2t80s/machines/cpc"
	"github.com/laullon/b2t80s/machines/gameboy"
	"github.com/laullon/b2t80s/machines/msx"
	"github.com/laullon/b2t80s/machines/nes"
	"github.com/laullon/b2t80s/machines/zx"
	"github.com/veandco/go-sdl2/sdl"

	_ "net/http/pprof"
)

func main() {

	emulator.CartFile = flag.String("cart", "", "NESncart file to load")
	emulator.TapFile = flag.String("tap", "", "tap file to load")
	emulator.RomFile = flag.String("rom", "", "msx1 rom file to load - format: [mapper::]filename - Mappers:konami")
	z80File := flag.String("z80", "", "z80 file to load")
	mode := flag.String("mode", "48k", "Spectrum model to emulate [48k|128k|plus3|cpc464|cpc6128|msx1]")
	emulator.Debug = flag.Bool("debug", false, "shows debugger")
	// turbo := flag.Bool("turbo", false, "run faster")

	emulator.Breaks = flag.String("bp", "", "Breakpoints [0xXXXX[,0xXXXX,...]]")
	emulator.WatchPoints = flag.String("wp", "", "Memory Watch Points [0xXXXX[,0xXXXX,...]]")

	emulator.LoadSlow = flag.Bool("slow", false, "Real Spectrum loading process")
	emulator.DskAFile = flag.String("dskA", "", "disc file to load on drive A")

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	var machine emulator.Machine
	var name string

	if len(*z80File) > 0 {
		machine = zx.LoadZ80File(*z80File)
		name = "ZX Spectrum"
	} else {
		switch *mode {
		case "48k":
			machine = zx.NewZX48K()
			name = "ZX Spectrum 48k"
		case "128k":
			machine = zx.NewZX128K()
			name = "ZX Spectrum 128k"
		case "plus3":
			machine = zx.NewZXPlus3()
			name = "ZX Spectrum +3"
		case "cpc464", "cpc":
			machine = cpc.NewCPC(true)
			name = "Amstrad CPC 464"
		case "cpc6128":
			machine = cpc.NewCPC(false)
			name = "Amstrad CPC 6128"
		case "msx":
			machine = msx.NewMSX()
			name = "MSX 1"
		case "atetris":
			machine = atetris.NewATetris()
			name = "Tetris"
		case "nes":
			machine = nes.NewNES()
			name = "Nes"
		case "gb":
			machine = gameboy.New()
			name = "GameBoy"
		default:
			panic(fmt.Errorf("mode '%s' not valid", *mode))
		}
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	if err := gl.Init(); err != nil {
		panic(err)
	}

	game := emulator.NewGame(name, machine)
	game.SetOnKey(machine.OnKey)

	log.Printf("opengl version %s", gl.GoStr(gl.GetString(gl.VERSION)))

	if *emulator.Debug {
		emulator.NewDebugWindow(name, machine)
	}

	wait := time.Duration(time.Second)
	ticker := time.NewTicker(wait)
	go func() {
		for range ticker.C {
			str := fmt.Sprintf("time: %s - FPS: %03.2f", machine.Clock().Stats(), machine.Monitor().FPS())
			println(str)
			game.SetStatus(str)
		}
	}()

	go func() {
		machine.Clock().Run()
	}()

	gui.PoolEvents()

	// ui.App = app.New()

	// w := ui.App.NewWindow(name + " - b2t80s Emulator")

	// status := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	// statusControl := container.New(layout.NewHBoxLayout())
	// for _, control := range machine.UIControls() {
	// 	statusControl.Add(control.Widget())
	// }

	// display := machine.Monitor().Canvas()

	// var debugTabs *container.AppTabs
	// var controls map[string]ui.Control

	// if *emulator.Debug {
	// 	var breaks []uint16
	// 	if len(*emulator.Breaks) > 0 {
	// 		bps := strings.Split(*emulator.Breaks, ",")
	// 		for _, bp := range bps {
	// 			n, err := strconv.ParseUint(bp, 0, 16)
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 			breaks = append(breaks, uint16(n))
	// 		}
	// 	}

	// 	db := emulator.NewDebugger(machine.Clock(), breaks)
	// 	machine.SetDebugger(db)
	// 	controls = machine.Control()

	// 	debugTabs = container.NewAppTabs()
	// 	for n, ctl := range controls {
	// 		debugTabs.Append(container.NewTabItem(n, ctl.Widget()))
	// 	}
	// 	debugTabs.SelectTabIndex(0)

	// 	debugger := container.New(layout.NewBorderLayout(db.UI(), nil, nil, nil),
	// 		db.UI(),
	// 		debugTabs,
	// 	)

	// 	statusBar := fyne.NewContainerWithLayout(
	// 		layout.NewBorderLayout(nil, nil, status, statusControl),
	// 		status,
	// 		statusControl,
	// 	)

	// 	w.SetContent(
	// 		fyne.NewContainerWithLayout(
	// 			layout.NewBorderLayout(nil, statusBar, nil, debugger),
	// 			display,
	// 			debugger,
	// 			statusBar,
	// 		),
	// 	)
	// 	ui.App.Settings().SetTheme(theme.LightTheme())
	// } else {
	// 	ui.App.Settings().SetTheme(theme.DarkTheme())
	// 	w.SetContent(
	// 		fyne.NewContainerWithLayout(
	// 			layout.NewBorderLayout(nil, nil, nil, nil),
	// 			display,
	// 		),
	// 	)
	// }

	// w.Canvas().(desktop.Canvas).SetOnKeyDown(machine.OnKeyEvent)
	// w.Canvas().(desktop.Canvas).SetOnKeyUp(machine.OnKeyEvent)

	// if *emulator.Debug {
	// wait := time.Duration(20 * time.Millisecond)
	// ticker := time.NewTicker(wait)
	// go func() {
	// 	for range ticker.C {
	// 		controls[debugTabs.CurrentTab().Text].Update()
	// 		status.SetText(fmt.Sprintf("time: %s - FPS: %03.2f", machine.Clock().Stats(), machine.Monitor().FPS()))
	// 	}
	// }()
	// }

	// // w.CenterOnScreen()
	// w.ShowAndRun()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

}
