package emulator

import (
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/laullon/b2t80s/gui"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	runtime.LockOSThread()
}

type Window interface {
	Run()
	SetOnKey(func(sdl.Scancode))
	SetStatus(txt string)
}

type window struct {
	img *Display

	win gui.Window

	status gui.Label

	displayTexture uint32
	displayFrameID uint32
}

func NewWindow(name string, machine Machine) Window {
	win := &window{
		img: machine.Monitor().Screen(),
		win: gui.NewWindow(name, gui.Size{800, 600}),
	}

	win.init()

	machine.Monitor().SetRedraw(func() {}) // TODO: need it?

	win.status = gui.NewLabel("staus", gui.Rect{0, 0, 330, 50})
	bt := gui.NewButton("staus", gui.Rect{330, 0, 330, 50})

	grid := gui.NewHGrid(3, 50)
	grid.Add(bt, win.status)
	grid.Resize(gui.Rect{0, 0, 800, 600})

	win.win.SetMainUI(grid)
	win.win.AddMouseListeners(bt)

	return win
}

func (win *window) SetStatus(txt string) {
	win.status.SetText(txt)
}

func (win *window) Run() {
	// go func() {
	// 	for running := true; running; {
	// 		<-win.redraw
	// 		win.window.GLMakeCurrent(win.context)
	// 		gl.ClearColor(0, 0, 0, 1)
	// 		gl.Clear(gl.COLOR_BUFFER_BIT)
	// 		win.render()
	// 		win.window.GLSwap()
	// 	}
	// }()
}

func (win *window) SetOnKey(onKey func(sdl.Scancode)) {
	win.win.SetOnKey(onKey)
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
}

func (win *window) Render() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.BindTexture(gl.TEXTURE_2D, win.displayTexture)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0,
		0, 0,
		int32(win.img.Image.Rect.Size().X), int32(win.img.Image.Rect.Size().Y),
		gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(win.img.Image.Pix))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	w, h := 10, 10 // win.win.GetSize()
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

	// win.ui.Render()

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	gl.Flush()
}
