//go:generate $HOME/go/bin/go-bindata -pkg data -o data/data.go data/...
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines/atetris"
	"github.com/laullon/b2t80s/machines/cpc"
	"github.com/laullon/b2t80s/machines/msx"
	"github.com/laullon/b2t80s/machines/nes"
	"github.com/laullon/b2t80s/machines/zx"
	"github.com/laullon/b2t80s/ui"

	_ "net/http/pprof"
)

func main() {
	nes.CartFile = flag.String("cart", "", "NESncart file to load")
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
		default:
			panic(fmt.Errorf("mode '%s' not valid", *mode))
		}
	}

	emulator.App = app.New()
	emulator.App.Settings().SetTheme(theme.LightTheme())

	w := emulator.App.NewWindow(name + " - b2t80s Emulator")

	debugger := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	status := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	controls := container.New(layout.NewHBoxLayout())
	for _, control := range machine.UIControls() {
		controls.Add(control.Widget())
	}

	statusBar := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, status, controls),
		status,
		controls,
	)

	display := machine.Monitor().Canvas()

	var cpuCtl ui.Control

	if *emulator.Debug {
		var breaks []uint16
		if len(*emulator.Breaks) > 0 {
			bps := strings.Split(*emulator.Breaks, ",")
			for _, bp := range bps {
				n, err := strconv.ParseUint(bp, 0, 16)
				if err != nil {
					panic(err)
				}
				breaks = append(breaks, uint16(n))
			}
		}

		db := emulator.NewDebugger(machine.Clock(), breaks)
		machine.SetDebugger(db)
		cpuCtl = machine.CPUControl()

		debugger := container.New(layout.NewVBoxLayout(),
			widget.NewLabel("Debugger"),
			fyne.NewContainerWithLayout(
				layout.NewGridLayoutWithColumns(3),
				widget.NewButton("Stop", func() {
					db.Stop()
				}),
				widget.NewButton("Continue", func() {
					db.Continue()
				}),
				widget.NewButton("Step", func() {
					db.Step()
				}),
				widget.NewButton("Stop Next Frame", func() {
					db.StopNextFrame()
				}),
				widget.NewButton("Dump 5 Frames", func() {
				}),
				widget.NewCheck("Dump", func(on bool) {
					panic(-1)
				}),
			),
			cpuCtl.Widget(),
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
			if *emulator.Debug {
				debugger.SetText(emulator.MachineStatus.Status())
				cpuCtl.Update()
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
