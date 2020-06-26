package msx

import (
	"image"
	"image/color"
	"image/draw"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/laullon/b2t80s/emulator"
)

type sprite struct {
	x, y    int
	h, w    int
	pattern byte
	colour  byte
	ec      bool
}

func newSprite(data []byte) (*sprite, bool) {
	sprt := &sprite{}

	sprt.y = int(data[0])
	sprt.x = int(data[1])
	sprt.pattern = data[2]
	sprt.colour = data[3] & 0x0f
	sprt.ec = data[3]&0x80 != 0

	if sprt.y > 240 {
		sprt.y -= 256
	}

	sprt.y++

	return sprt, data[0] == 208
}

func (vdp *tms9918) drawSprites() {
	sprites := make([]*sprite, 0)
	height := 8
	if vdp.si {
		height = 16
	}

	vdp.s5 = false
	vdp.fs = 0
	for idx := uint16(0); idx < 32; idx++ {
		sprt, done := newSprite(vdp.vram[vdp.sa+(idx*4) : vdp.sa+(idx*4)+4])
		if done {
			break
		}

		if (vdp.y >= sprt.y) && (vdp.y < sprt.y+height) {
			if len(sprites) == 4 {
				vdp.s5 = true
				vdp.fs = byte(idx)
				break
			} else {
				sprites = append([]*sprite{sprt}, sprites...)
			}
		}
	}

	for _, sprt := range sprites {
		sprt.drawSprite(vdp.y, vdp.si, vdp.sg, vdp.vram, vdp.display)
	}
}

func (sprt *sprite) drawSprite(yPos int, si bool, sg uint16, vram []byte, display *image.RGBA) {
	if sprt.colour == 0 {
		return
	}
	cols := uint16(1)
	if si {
		cols = 2
	}
	for i := uint16(0); i < cols; i++ {
		y := uint16(yPos - int(sprt.y))
		b := vram[sg+(uint16(sprt.pattern)&252)<<3+y+(i*16)]
		for x := 0; x < 8; x++ {
			sx := sprt.x + x + int(i*8)
			if sx >= 0 && sx < 256 {
				if b&(1<<(7-x)) != 0 {
					display.SetRGBA(sx, sprt.y+int(y), palette[sprt.colour])
				}
			}
		}
	}
}

type spriteControl struct {
	ui      fyne.CanvasObject
	show    *widget.Button
	vdp     *tms9918
	debuger emulator.Debugger
}

func newSpriteControl(vdp *tms9918, debuger emulator.Debugger) *spriteControl {
	ctrl := &spriteControl{vdp: vdp, debuger: debuger}

	ctrl.show = widget.NewButton("show Sprites", ctrl.doShow)

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
	ctrl.debuger.Stop()

	container := fyne.NewContainerWithLayout(layout.NewGridLayout(8))

	for idx := uint16(0); idx < 32; idx++ {
		sprite, _ := newSprite(ctrl.vdp.vram[ctrl.vdp.sa+(idx*4) : ctrl.vdp.sa+(idx*4)+4])
		sprite.x = 0
		sprite.y = 0

		size := 8
		if ctrl.vdp.si {
			size = 16
		}

		display := image.NewRGBA(image.Rect(0, 0, size, size))
		draw.Draw(display, display.Bounds(), &image.Uniform{color.RGBA{125, 125, 125, 255}}, image.ZP, draw.Src)
		for y := 0; y < size; y++ {
			sprite.drawSprite(y, ctrl.vdp.si, ctrl.vdp.sg, ctrl.vdp.vram, display)
		}

		img := canvas.NewImageFromImage(display)
		img.SetMinSize(fyne.NewSize(size*4, size*4))
		img.ScaleMode = canvas.ImageScalePixels

		container.AddObject(img)
	}

	c := fyne.CurrentApp().Driver().CanvasForObject(ctrl.ui)
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(ctrl.ui)
	var pop *widget.PopUp
	pop = widget.NewPopUp(container, c)
	pop.Move(pos)

}
