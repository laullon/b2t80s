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

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/machines/cpc"
	"github.com/laullon/b2t80s/machines/zx"

	_ "net/http/pprof"
)

func init() {
	runtime.GOMAXPROCS(4)
	runtime.LockOSThread()
}

func main() {
	tapFile := flag.String("tap", "", "tap file to load")
	z80File := flag.String("z80", "", "z80 file to load")
	mode := flag.String("mode", "48k", "Spectrum model to emulate [48k|128k|plus3|cpc464|cpc6128]")
	debug := flag.Bool("debug", false, "shows debugger")
	// turbo := flag.Bool("turbo", false, "run faster")

	breaks := flag.String("bp", "", "Breakpoints [0xXXXX[,0xXXXX,...]]")
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

	cassette := emulator.NewTapCassette()
	if len(*tapFile) > 0 {
		cassette.LoadTapFile(*tapFile)
	}
	// log.Print(cassette)

	var machine machines.Machine

	if len(*z80File) > 0 {
		machine = zx.LoadZ80File(*z80File)
	} else {
		switch *mode {
		case "48k":
			machine = zx.NewZX48K(cassette)
		case "128k":
			machine = zx.NewZX128K(cassette)
		case "plus3":
			machine = zx.NewZXPlus3(cassette)
		case "cpc464", "cpc":
			machine = cpc.NewCPC(true, cassette)
		case "cpc6128":
			machine = cpc.NewCPC(false, cassette)
		default:
			panic(fmt.Errorf("mode '%s' not valid", *mode))
		}
	}

	if len(*breaks) > 0 {
		bps := strings.Split(*breaks, ",")
		for _, bp := range bps {
			n, err := strconv.ParseUint(bp, 0, 16)
			if err != nil {
				panic(err)
			}
			machine.Debugger().SetBreakPoint(uint16(n))
		}
	}

	app := app.New()
	display := canvas.NewImageFromImage(machine.Display())
	display.FillMode = canvas.ImageFillOriginal
	display.ScalingFilter = canvas.NearestFilter
	display.SetMinSize(fyne.NewSize(352*2, 296*2))

	w := app.NewWindow("ZX Spectrum")

	reg := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	pas := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	ins := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Bold: true})
	dis := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	status := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	volumen := widget.NewSlider(0, 90)
	volumen.OnChanged = machine.GetVolumeControl()
	volumen.MinSize()

	statusBar := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, status, volumen),
		status,
		volumen,
	)

	if *debug {
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
					machine.Debugger().DumpNextFrame()
				}),
				widget.NewCheck("Dump", func(on bool) {
					machine.Debugger().SetDump(on)
				}),
			),
			reg,
			pas,
			ins,
			dis,
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
			display.Refresh()
			pas.SetText(machine.Debugger().GetLog())
			reg.SetText(machine.Debugger().GetRegisters())
			ins.SetText(machine.Debugger().GetNextInstruction())
			dis.SetText(machine.Debugger().GetFollowingInstruction())
			status.SetText(machine.Debugger().GetStatus())
		}
	}()

	go func() {
		runtime.LockOSThread()
		machine.Run()
	}()

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
