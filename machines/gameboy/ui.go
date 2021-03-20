package gameboy

import (
	"image"
	"strconv"

	"fyne.io/fyne/v2"
	canvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	layout "fyne.io/fyne/v2/layout"
	"github.com/laullon/b2t80s/ui"
)

type lcdDebugControl struct {
	ui      fyne.CanvasObject
	lcd     *lcd
	display *image.RGBA

	x, y, scX, scY *ui.RegText
}

func newLcdControl(lcd *lcd) *lcdDebugControl {
	ctrl := &lcdDebugControl{
		lcd:     lcd,
		display: image.NewRGBA(image.Rect(0, 0, 32*8, 32*8)),
	}

	img := canvas.NewImageFromImage(ctrl.display)
	img.SetMinSize(fyne.NewSize((64*8)+4+84, (64*8)+2))
	img.ScaleMode = canvas.ImageScalePixels

	ctrl.x = ui.NewRegText("X:")
	ctrl.y = ui.NewRegText("Y:")
	ctrl.scX = ui.NewRegText("Scroll X:")
	ctrl.scY = ui.NewRegText("Scroll Y:")

	c1 := container.New(layout.NewFormLayout(),
		ctrl.y.Label, ctrl.y.Value,
		ctrl.x.Label, ctrl.x.Value,
	)

	c2 := container.New(layout.NewFormLayout(),
		ctrl.scX.Label, ctrl.scX.Value,
		ctrl.scY.Label, ctrl.scY.Value,
	)

	regs := container.New(layout.NewGridLayoutWithColumns(3), c1, c2)

	ctrl.ui = fyne.NewContainerWithLayout(layout.NewBorderLayout(regs, nil, nil, nil), regs, img)

	return ctrl
}

func (ctrl *lcdDebugControl) Widget() fyne.CanvasObject {
	return ctrl.ui
}

func (ctrl *lcdDebugControl) Update() {
	ctrl.x.Update(strconv.Itoa(ctrl.lcd.lx))
	ctrl.y.Update(strconv.Itoa(ctrl.lcd.ly))
	ctrl.scX.Update(strconv.Itoa(int(ctrl.lcd.scx)))
	ctrl.scY.Update(strconv.Itoa(int(ctrl.lcd.scy)))

	for r := uint16(0); r < 32; r++ {
		y := int(r * 8)
		for c := uint16(0); c < 32; c++ {
			x := int(c * 8)
			for y_off := uint16(0); y_off < 8; y_off++ {
				tileAddr := c*16 + r*16*32 + y_off*2
				b1, _ := ctrl.lcd.vRAM.ReadPort(tileAddr)
				b2, _ := ctrl.lcd.vRAM.ReadPort(tileAddr + 1)
				for x_off := 0; x_off < 8; x_off++ {
					c := (b1 & 1) | ((b2 & 1) << 1)
					ctrl.display.Set(x+(7-x_off), y+int(y_off), ctrl.lcd.palette[c])
					b1 >>= 1
					b2 >>= 1
				}
			}
		}
	}

	ctrl.ui.Refresh()
}
