package emulator

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func init() {
	runtime.LockOSThread()
}

type Window interface {
	Run()
	SetOnKey(func(sdl.Scancode))
}

type window struct {
	img   *Display
	onKey func(sdl.Scancode)

	displayTexture uint32
	displayFrameID uint32

	statusTexture uint32
	statusFrameID uint32

	a, b, c, d, e, f uint32 //// do not remove

	window  *sdl.Window
	context sdl.GLContext

	redraw chan struct{}
}

func NewWindow(name string, machine Machine) Window {
	win := &window{
		img:    machine.Monitor().Screen(),
		redraw: make(chan struct{}),
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Game", 50, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE|sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	win.window = window

	context, err := window.GLCreateContext()
	if err != nil {
		panic(err)
	}
	win.context = context

	gl.Init()
	log.Printf("opengl version %s", gl.GoStr(gl.GetString(gl.VERSION)))

	win.init()

	machine.Monitor().SetRedraw(func() {
		win.redraw <- struct{}{}
	})

	return win
}

func (win *window) Run() {
	println("run")

	for running := true; running; {
		// select {
		// case <-win.redraw:
		<-win.redraw
		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		win.render()
		win.window.GLSwap()
		// default:
		// }
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch event := e.(type) {
			case *sdl.WindowEvent:
				if event.Event == sdl.WINDOWEVENT_CLOSE {
					running = false
					sdl.Quit()
				}
			case *sdl.QuitEvent:
				running = false
				sdl.Quit()
			case *sdl.KeyboardEvent:
				if event.Repeat == 0 {
					win.onKey(event.Keysym.Scancode)
				}
			}
		}
	}
}

func (win *window) SetOnKey(onKey func(sdl.Scancode)) {
	win.onKey = onKey
}

func (win *window) init() {
	// DISPLAY
	gl.GenTextures(1, &win.displayTexture)
	gl.BindTexture(gl.TEXTURE_2D, win.displayTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB,
		int32(win.img.Image.Rect.Size().X), int32(win.img.Image.Rect.Size().Y),
		0, gl.RGB, gl.UNSIGNED_BYTE,
		nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.GenFramebuffers(1, &win.displayFrameID)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, win.displayFrameID)
	gl.FramebufferTexture2D(gl.READ_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, win.displayTexture, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)

	// STATUS
	gl.GenTextures(2, &win.statusTexture)
	gl.BindTexture(gl.TEXTURE_2D, win.statusTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB,
		int32(win.img.Image.Rect.Size().X), int32(win.img.Image.Rect.Size().Y),
		0, gl.RGB, gl.UNSIGNED_BYTE,
		nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.GenFramebuffers(2, &win.statusFrameID)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, win.statusFrameID)
	gl.FramebufferTexture2D(gl.READ_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, win.statusTexture, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
}

func (win *window) render() bool {
	// println("render")

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.BindTexture(gl.TEXTURE_2D, win.displayTexture)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0,
		0, 0,
		int32(win.img.Image.Rect.Size().X), int32(win.img.Image.Rect.Size().Y),
		gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(win.img.Image.Pix))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	//*****
	m := NewDisplay(image.Rect(0, 0, 100, 100))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.Point{}, draw.Src)

	d := &font.Drawer{
		Dst:  m,
		Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(10, 10),
	}
	d.DrawString("hola")

	gl.BindTexture(gl.TEXTURE_2D, win.statusTexture)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0,
		0, 0,
		int32(m.Image.Rect.Size().X), int32(m.Image.Rect.Size().Y),
		gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(m.Image.Pix))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	//*****

	w, h := win.window.GetSize()
	// sW, sH := win.mainWin.GetContentScale()
	// w *= int(sW)
	// h *= int(sH)

	ratioOrg := float64(win.img.Size.X) / float64(win.img.Size.Y)
	ratioDst := float64(w) / float64(h)

	var newW, newH int32
	var offX, offY int32
	if ratioDst > ratioOrg {
		// (wi * hs/hi, hs)
		newW = int32(float64(win.img.Size.X) * float64(h) / float64(win.img.Size.Y))
		newH = int32(h)
		offX = (int32(w) - newW) / 2
	} else {
		// hi * ws/wi
		newW = int32(w)
		newH = int32(float64(win.img.Size.Y) * float64(w) / float64(win.img.Size.X))
		offY = (int32(h) - newH) / 2
	}

	// println(ratioOrg, " - ", ratioDst, " - ", ratioOrg < ratioDst, "  ->  ", w, "x", h)

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, win.displayFrameID)
	gl.EnableVertexAttribArray(0)
	gl.BlitFramebuffer(
		int32(win.img.ViewPortRect.Min.X), int32(win.img.ViewPortRect.Min.Y), int32(win.img.ViewPortRect.Max.X), int32(win.img.ViewPortRect.Max.Y),
		offX, offY, int32(newW)+offX, int32(newH)+offY,
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, win.statusFrameID)
	gl.EnableVertexAttribArray(0)
	gl.BlitFramebuffer(
		int32(m.Image.Rect.Min.X), int32(m.Image.Rect.Min.Y), int32(m.Image.Rect.Max.X), int32(m.Image.Rect.Max.Y),
		0, 0, int32(m.Image.Rect.Max.X), int32(m.Image.Rect.Max.Y),
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)

	gl.Flush()
	return true
}
