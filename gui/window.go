package gui

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Window interface {
	SetMainUI(GUIObject)
	SetOnKey(func(sdl.Scancode))
	MoveTo(Point)
}

type window struct {
	sdlWin  *sdl.Window
	context sdl.GLContext

	ui             GUIObject
	mouseListeners []MouseTarget
	keyListeners   []KeyTarget
	onKey          func(sdl.Scancode)
}

func NewWindow(name string, size Size) Window {
	win := &window{}

	sdlWin, err := sdl.CreateWindow(name, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		size.W, size.H, sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE|sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	win.sdlWin = sdlWin

	context, err := sdlWin.GLCreateContext()
	if err != nil {
		panic(err)
	}
	win.context = context

	id, err := sdlWin.GetID()
	if err != nil {
		panic(err)
	}
	windows[id] = win

	win.sdlWin.Raise()
	return win
}

func (w *window) MoveTo(p Point) {
	w.sdlWin.SetPosition(p.X, p.Y)
}

func (w *window) SetOnKey(onKey func(sdl.Scancode)) {
	w.onKey = onKey
}

func (w *window) SetMainUI(ui GUIObject) {
	w.ui = ui
	wi, he := w.sdlWin.GetSize()
	ui.Resize(Rect{0, 0, wi, he})
	w.addListeners(ui)
}

func (w *window) addListeners(obj GUIObject) {
	fmt.Printf("obj -> %T ", obj)
	switch ctrl := obj.(type) {
	case KeyTarget:
		w.keyListeners = append(w.keyListeners, ctrl)
		print("KeyTarget ")
	}
	switch ctrl := obj.(type) {
	case MouseTarget:
		w.mouseListeners = append(w.mouseListeners, ctrl)
		print("MouseTarget ")
	}
	println()
	for _, chd := range obj.GetChildrens() {
		if chd != nil {
			w.addListeners(chd)
		}
	}
}
