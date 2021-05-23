package emulator

import (
	"github.com/laullon/b2t80s/gui"
	"github.com/veandco/go-sdl2/sdl"
)

type Game interface {
	SetStatus(txt string)
	SetOnKey(onKey func(sdl.Scancode))
}

type game struct {
	window gui.Window
	status gui.Label
}

func NewGame(name string, machine Machine) Game {
	game := &game{
		window: gui.NewWindow(name, gui.Size{800, 600}),
	}

	machine.Monitor().SetRedraw(func() {}) // TODO: need it?

	img := gui.NewDisplayViewer(machine.Monitor().Screen())
	game.status = gui.NewLabel("")

	hct := gui.NewVerticalHCT()
	hct.SetCenter(img)
	hct.SetTail(game.status, 30)

	game.window.SetMainUI(hct)

	return game
}

func (game *game) SetStatus(txt string) {
	game.status.SetText(txt)
}

func (game *game) SetOnKey(onKey func(sdl.Scancode)) {
	game.window.SetOnKey(onKey)
}
