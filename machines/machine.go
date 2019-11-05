package machines

import (
	"image"

	"fyne.io/fyne"
	"github.com/laullon/b2t80s/emulator"
)

var LoadSlow *bool
var DskAFile *string

type Machine interface {
	OnKeyEvent(event *fyne.KeyEvent)
	Run()
	Debugger() emulator.Debugger
	Display() image.Image
}
