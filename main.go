//go:generate $HOME/go/bin/go-bindata -pkg data -o data/data.go data/...
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/machines/atari/a2600"
	"github.com/laullon/b2t80s/machines/atetris"
	"github.com/laullon/b2t80s/machines/cpc"
	"github.com/laullon/b2t80s/machines/msx"
	"github.com/laullon/b2t80s/machines/zx"

	_ "net/http/pprof"
)

func main() {
	machines.TapFile = flag.String("tap", "", "tap file to load")
	machines.RomFile = flag.String("rom", "", "msx1 rom file to load - format: [mapper::]filename - Mappers:konami")
	z80File := flag.String("z80", "", "z80 file to load")
	mode := flag.String("mode", "48k", "Spectrum model to emulate [48k|128k|plus3|cpc464|cpc6128|msx1]")
	machines.Debug = flag.Bool("debug", false, "shows debugger")
	// turbo := flag.Bool("turbo", false, "run faster")

	// breaks := flag.String("bp", "", "Breakpoints [0xXXXX[,0xXXXX,...]]")
	machines.LoadSlow = flag.Bool("slow", false, "Real Spectrum loading process")
	machines.DskAFile = flag.String("dskA", "", "disc file to load on drive A")

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

	var machine machines.Machine
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
		case "a2600":
			machine = a2600.NewA2600()
			name = "Atari 2600"
		default:
			panic(fmt.Errorf("mode '%s' not valid", *mode))
		}
	}

	// if len(*breaks) > 0 {
	// 	bps := strings.Split(*breaks, ",")
	// 	for _, bp := range bps {
	// 		n, err := strconv.ParseUint(bp, 0, 16)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		machine.Debugger().SetBreakPoint(uint16(n))
	// 	}
	// }

	app := app.New()
	app.Settings().SetTheme(theme.LightTheme())

	w := app.NewWindow(name + " - b2t80s Emulator")

	debugger := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	status := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	controls := widget.NewHBox()
	for _, control := range machine.UIControls() {
		controls.Append(control.Widget())
	}

	statusBar := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, status, controls),
		status,
		controls,
	)

	display := machine.Monitor().Canvas()

	if *machines.Debug {
		debugger := widget.NewVBox(
			widget.NewLabel("Debugger"),
			fyne.NewContainerWithLayout(
				layout.NewGridLayoutWithColumns(3),
				widget.NewButton("Stop", func() {
					machine.Debugger().Stop()
				}),
				widget.NewButton("Continue", func() {
					machine.Debugger().Continue()
				}),
				widget.NewButton("Step", func() {
					machine.Debugger().Step()
				}),
				widget.NewButton("Stop Next Frame", func() {
					machine.Debugger().StopNextFrame()
				}),
				widget.NewButton("Dump 5 Frames", func() {
				}),
				widget.NewCheck("Dump", func(on bool) {
					machine.Debugger().SetDump(on)
				}),
			),
			debugger,
		)

		w.SetContent(
			fyne.NewContainerWithLayout(
				layout.NewBorderLayout(nil, statusBar, nil, debugger),
				display,
				debugger,
				statusBar,
			),
		)
	} else {
		w.SetContent(
			fyne.NewContainerWithLayout(
				layout.NewBorderLayout(nil, statusBar, nil, nil),
				display,
				statusBar,
			),
		)
	}

	w.Canvas().(desktop.Canvas).SetOnKeyDown(machine.OnKeyEvent)
	w.Canvas().(desktop.Canvas).SetOnKeyUp(machine.OnKeyEvent)

	wait := time.Duration(20 * time.Millisecond)
	ticker := time.NewTicker(wait)
	go func() {
		for range ticker.C {
			if *machines.Debug {
				debugger.SetText(machine.Debugger().GetStatus())
			}
			status.SetText(fmt.Sprintf("time: %s - FPS: %03.2f", machine.Clock().Stats(), machine.Monitor().FPS()))
		}
	}()

	go func() {
		machine.Clock().Run()
	}()

	// w.CenterOnScreen()
	w.ShowAndRun()

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
