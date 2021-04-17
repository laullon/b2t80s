package emulator

import (
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window interface {
	Run()
	SetOnKey(func(glfw.Key))
}

type window struct {
	mainWin *glfw.Window
	img     *Display

	texture uint32
	fobID   uint32

	onKey func(glfw.Key)
}

func NewWindow(name string, img *Display) Window {
	var err error
	window := &window{
		img: img,
	}

	if err = glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window.mainWin, err = glfw.CreateWindow(800, 600, name, nil, nil)
	if err != nil {
		panic(err)
	}

	window.mainWin.SetSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		window.draw()
	})

	window.mainWin.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action != 2 {
			window.onKey(key)
			println("key:", key, "action:", action)
		}
	})

	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	println("OpenGL version", version)

	window.iniTexture()

	return window
}

func (win *window) iniTexture() {
	win.mainWin.MakeContextCurrent()
	gl.GenTextures(1, &win.texture)
	gl.BindTexture(gl.TEXTURE_2D, win.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB,
		int32(win.img.Rect.Size().X), int32(win.img.Rect.Size().Y),
		0, gl.RGB, gl.UNSIGNED_BYTE,
		nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.GenFramebuffers(1, &win.fobID)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, win.fobID)
	gl.FramebufferTexture2D(gl.READ_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, win.texture, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
}

func (win *window) Run() {
	t := time.Now()

	for !win.mainWin.ShouldClose() {
		// win.draw()
		time.Sleep(time.Second/time.Duration(60) - time.Since(t))
		t = time.Now()
	}
}

func (win *window) draw() {
	win.mainWin.MakeContextCurrent()

	println("- draw -", win.mainWin.ShouldClose())

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	w, h := win.mainWin.GetSize()
	sW, sH := win.mainWin.GetContentScale()
	w *= int(sW)
	h *= int(sH)

	gl.BindTexture(gl.TEXTURE_2D, win.texture)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0,
		0, 0,
		int32(win.img.Rect.Size().X), int32(win.img.Rect.Size().Y),
		gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(win.img.Pix))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	ratioOrg := float64(win.img.Rect.Size().X) / float64(win.img.Rect.Size().Y)
	ratioDst := float64(w) / float64(h)

	var newW, newH int32
	var offX, offY int32
	if ratioDst > ratioOrg {
		// (wi * hs/hi, hs)
		newW = int32(float64(win.img.Rect.Size().X) * float64(h) / float64(win.img.Rect.Size().Y))
		newH = int32(h)
		offX = (int32(w) - newW) / 2
	} else {
		// hi * ws/wi
		newW = int32(w)
		newH = int32(float64(win.img.Rect.Size().Y) * float64(w) / float64(win.img.Rect.Size().X))
		offY = (int32(h) - newH) / 2
	}

	// println(ratioOrg, " - ", ratioDst, " - ", ratioOrg < ratioDst, "  ->  ", w, "x", h)

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, win.fobID)
	gl.EnableVertexAttribArray(0)
	gl.BlitFramebuffer(
		0, 0, int32(win.img.Rect.Size().X), int32(win.img.Rect.Size().Y),
		offX, offY, int32(newW)+offX, int32(newH)+offY,
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)

	glfw.PollEvents()
	win.mainWin.SwapBuffers()
}

func (win *window) SetOnKey(onKey func(glfw.Key)) {
	win.onKey = onKey
}
