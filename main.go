//go:generate $HOME/go/bin/go-bindata -pkg data -o data/data.go data/...
package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/machines/atetris"
	"github.com/laullon/b2t80s/machines/cpc"
	"github.com/laullon/b2t80s/machines/msx"
	"github.com/laullon/b2t80s/machines/nes"
	"github.com/laullon/b2t80s/machines/zx"

	_ "net/http/pprof"
)

func init() { runtime.LockOSThread() }

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

	go func() {
		machine.Clock().Run()
	}()

	if err := glfw.Init(); err != nil {
		log.Fatalln("Fallo al inicializar glfw:", err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(
		800,
		600,
		name,
		nil,
		nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL versión", version)
	gl.ClearColor(.5, 1, 0, 0.0)
	gl.Enable(gl.TEXTURE_2D)

	// f, err := os.Open("/Users/glaullon/go/src/github.com/laullon/b2t80s/machines/atetris/tests/testMode_ok.png")
	// if err != nil {
	// 	panic(err)
	// }

	// img1, err := png.Decode(f)
	// if err != nil {
	// 	panic(err)
	// }

	tex, err := newTexture(machine.Monitor().Screen())
	if err != nil {
		panic(err)
	}

	for !window.ShouldClose() {
		// rgba := img1.(*image.RGBA)
		rgba := machine.Monitor().Screen()

		gl.BindTexture(gl.TEXTURE_2D, tex)
		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			int32(rgba.Rect.Size().X),
			int32(rgba.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix))

		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Begin(gl.QUADS)
		gl.TexCoord2f(0, 0)
		gl.Vertex2f(-0.9, 0.9)
		gl.TexCoord2f(1, 0)
		gl.Vertex2f(0.9, 0.9)
		gl.TexCoord2f(1, 1)
		gl.Vertex2f(0.9, -0.9)
		gl.TexCoord2f(0, 1)
		gl.Vertex2f(-0.9, -0.9)
		gl.End()

		window.SwapBuffers()
		glfw.PollEvents()
	}

	for !window.ShouldClose() {
		glfw.PollEvents()

		gl.BindTexture(gl.TEXTURE_2D, tex)

		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Begin(gl.QUADS) // Cada coordenada con su coordenada de textura
		gl.TexCoord2f(0, 0)
		gl.Vertex2f(-0.5, 0.5)
		gl.TexCoord2f(1, 0)
		gl.Vertex2f(0.5, 0.5)
		gl.TexCoord2f(1, 1)
		gl.Vertex2f(0.5, -0.5)
		gl.TexCoord2f(0, 1)
		gl.Vertex2f(-0.5, -0.5)
		gl.End()

		window.SwapBuffers()
	}

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
	// 	wait := time.Duration(20 * time.Millisecond)
	// 	ticker := time.NewTicker(wait)
	// 	go func() {
	// 		for range ticker.C {
	// 			controls[debugTabs.CurrentTab().Text].Update()
	// 			status.SetText(fmt.Sprintf("time: %s - FPS: %03.2f", machine.Clock().Stats(), machine.Monitor().FPS()))
	// 		}
	// 	}()
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

const (
	vertexShaderSource = `
		#version 400
		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 400
		out vec4 frag_colour;
		void main() {
  			frag_colour = vec4(1, 1, 1, 1.0);
		}
	` + "\x00"
)

func initGlfw(name string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 300, name, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func newTexture(rgba *image.RGBA) (uint32, error) {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		nil)

	return texture, nil
}
