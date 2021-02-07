package nes

import (
	"image"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/laullon/b2t80s/cpu/m6502"
)

type spriteControl struct {
	ui     fyne.CanvasObject
	show   *widget.Button
	ppuBus m6502.Bus
}

func newPalleteControl(ppuBus m6502.Bus) *spriteControl {
	ctrl := &spriteControl{ppuBus: ppuBus}

	ctrl.show = widget.NewButton("show palettes", ctrl.doShow)

	ctrl.ui = widget.NewHBox(
		widget.NewToolbarSeparator().ToolbarObject(),
		ctrl.show,
	)

	return ctrl
}

func (ctrl *spriteControl) Widget() fyne.CanvasObject {
	return ctrl.ui
}

func (ctrl *spriteControl) doShow() {
	container := fyne.NewContainerWithLayout(layout.NewGridLayout(8))

	for idx := uint16(0); idx < 32; idx++ {
		display := image.NewRGBA(image.Rect(0, 0, 160, 40))

		for y := 0; y < 2; y++ {
			for x := 0; x < 0x10; x++ {
				c := uint16(0x3f00 | (y << 4) | x)
				for dx := 0; dx < 10; dx++ {
					for dy := 0; dy < 10; dy++ {
						display.Set(x*10+dx, y*10+dy, colors[ctrl.ppuBus.Read(c)])
					}
				}
			}
		}
		img := canvas.NewImageFromImage(display)
		img.SetMinSize(fyne.NewSize(160, 40))
		img.ScaleMode = canvas.ImageScalePixels

		container.AddObject(img)
	}

	c := fyne.CurrentApp().Driver().CanvasForObject(ctrl.ui)
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(ctrl.ui)
	var pop *widget.PopUp
	pop = widget.NewPopUp(container, c)
	pop.Move(pos)

}
