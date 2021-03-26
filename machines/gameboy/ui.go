package gameboy

import (
	"encoding/hex"
	"fmt"
	"image"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	canvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	layout "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/laullon/b2t80s/cpu/lr35902"
	"github.com/laullon/b2t80s/ui"
)

type lcdDebugControl struct {
	ui      fyne.CanvasObject
	lcd     *lcd
	display *image.RGBA

	x, y, scX, scY  *ui.RegText
	status, control *ui.RegText

	sprites *widget.Label
}

func newLcdControl(lcd *lcd) *lcdDebugControl {
	ctrl := &lcdDebugControl{
		lcd:     lcd,
		display: image.NewRGBA(image.Rect(0, 0, 32*8, 12*8+2)),
	}

	img := canvas.NewImageFromImage(ctrl.display)
	img.ScaleMode = canvas.ImageScalePixels

	ctrl.x = ui.NewRegText("X:")
	ctrl.y = ui.NewRegText("Y:")
	ctrl.scX = ui.NewRegText("Scroll X:")
	ctrl.scY = ui.NewRegText("Scroll Y:")
	ctrl.status = ui.NewRegText("Status:")
	ctrl.control = ui.NewRegText("Control:")

	ctrl.sprites = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	c1 := container.New(layout.NewFormLayout(),
		ctrl.x.Label, ctrl.x.Value,
		ctrl.y.Label, ctrl.y.Value,
	)

	c2 := container.New(layout.NewFormLayout(),
		ctrl.scX.Label, ctrl.scX.Value,
		ctrl.scY.Label, ctrl.scY.Value,
	)

	c3 := container.New(layout.NewFormLayout(),
		ctrl.control.Label, ctrl.control.Value,
		ctrl.status.Label, ctrl.status.Value,
	)

	regs := container.New(layout.NewGridLayoutWithColumns(3), c1, c2, c3)

	ctrl.ui = fyne.NewContainerWithLayout(layout.NewBorderLayout(regs, ctrl.sprites, nil, nil), regs, img, ctrl.sprites)

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
	ctrl.control.Update(fmt.Sprintf("%08b", ctrl.lcd.control))
	ctrl.status.Update(fmt.Sprintf("%08b", ctrl.lcd.status))

	var sb strings.Builder
	sb.WriteString("X   Y   Tile Flag        X   Y   Tile Flag\n")
	for i := uint16(0); i < 40; i++ {
		sb.WriteString(fmt.Sprintf("%03d %03d 0x%02X 0b%08b  ",
			ctrl.lcd.oam[i*4+1],
			ctrl.lcd.oam[i*4+0],
			ctrl.lcd.oam[i*4+2],
			ctrl.lcd.oam[i*4+3],
		))
		if i%2 == 1 {
			sb.WriteString("\n")
		}
	}
	ctrl.sprites.Text = sb.String()

	for r := uint16(0); r < 12; r++ {
		y := int(r * 8)
		for c := uint16(0); c < 32; c++ {
			x := int(c * 8)
			for y_off := uint16(0); y_off < 8; y_off++ {
				tileAddr := c*16 + r*16*32 + y_off*2
				block := int(tileAddr >> 11)
				b1 := ctrl.lcd.vRAM[tileAddr]
				b2 := ctrl.lcd.vRAM[tileAddr+1]
				for x_off := 0; x_off < 8; x_off++ {
					c := (b1 & 1) | ((b2 & 1) << 1)
					ctrl.display.Set(x+(7-x_off), y+int(y_off)+block, ctrl.lcd.palette[c])
					b1 >>= 1
					b2 >>= 1
				}
			}
		}
	}

	ctrl.ui.Refresh()
}

/// *********************************
/// *********************************
/// *********************************

type timerDebugControl struct {
	ui    fyne.CanvasObject
	timer *timer

	div, tima, tma, tac *ui.RegText

	cpu ui.Control
}

func newTimerControl(cpu lr35902.LR35902, timer *timer) *timerDebugControl {
	ctrl := &timerDebugControl{
		timer: timer,
		cpu:   ui.NewLR35902UI(cpu),
	}

	ctrl.div = ui.NewRegText("div:")
	ctrl.tima = ui.NewRegText("tima:")
	ctrl.tma = ui.NewRegText("tma:")
	ctrl.tac = ui.NewRegText("tac:")

	c1 := container.New(layout.NewFormLayout(),
		ctrl.div.Label, ctrl.div.Value,
		ctrl.tima.Label, ctrl.tima.Value,
	)

	c2 := container.New(layout.NewFormLayout(),
		ctrl.tma.Label, ctrl.tma.Value,
		ctrl.tac.Label, ctrl.tac.Value,
	)

	regs := container.New(layout.NewGridLayoutWithColumns(3), c1, c2)
	panel := container.New(layout.NewVBoxLayout(), regs, ctrl.cpu.Widget())

	ctrl.ui = container.New(layout.NewBorderLayout(panel, nil, nil, nil), panel)

	return ctrl
}

func (ctrl *timerDebugControl) Widget() fyne.CanvasObject {
	return ctrl.ui
}

func (ctrl *timerDebugControl) Update() {
	ctrl.div.Update(strconv.Itoa(int(ctrl.timer.div)))
	ctrl.tima.Update(strconv.Itoa(int(ctrl.timer.tima)))
	ctrl.tma.Update(strconv.Itoa(int(ctrl.timer.tma)))
	ctrl.tac.Update(strconv.Itoa(int(ctrl.timer.tac)))
	ctrl.cpu.Update()
	ctrl.ui.Refresh()
}

/// *********************************
/// *********************************
/// *********************************

type serialDebugControl struct {
	ui     *fyne.Container
	text   *widget.Label
	buffer *[]byte
}

func newSerialControl(buffer *[]byte) *serialDebugControl {
	ctrl := &serialDebugControl{
		buffer: buffer,
		text:   &widget.Label{},
	}

	// ctrl.text.Color = color.Black
	// ctrl.text.TextSize = fyne.CurrentApp().Settings().Theme().Size("text")
	ctrl.text.TextStyle = fyne.TextStyle{Monospace: true}

	ctrl.ui = container.New(layout.NewBorderLayout(nil, nil, nil, nil), ctrl.text)

	return ctrl
}

func (ctrl *serialDebugControl) Widget() fyne.CanvasObject {
	return ctrl.ui
}

func (ctrl *serialDebugControl) Update() {
	ctrl.text.Text = hex.Dump(*ctrl.buffer)
	ctrl.ui.Refresh()
}
