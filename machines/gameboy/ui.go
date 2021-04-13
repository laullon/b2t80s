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

type ppuDebugControl struct {
	ui      fyne.CanvasObject
	ppu     *ppu
	display *image.RGBA

	x, y, scX, scY, wx, wy *ui.RegText
	status, control        *ui.RegText

	sprites *widget.Label
}

func newPPUControl(ppu *ppu) *ppuDebugControl {
	ctrl := &ppuDebugControl{
		ppu:     ppu,
		display: image.NewRGBA(image.Rect(0, 0, 32*8, 12*8+2)),
	}

	img := canvas.NewImageFromImage(ctrl.display)
	img.ScaleMode = canvas.ImageScalePixels

	ctrl.x = ui.NewRegText("lX:")
	ctrl.y = ui.NewRegText("lY:")
	ctrl.scX = ui.NewRegText("scX:")
	ctrl.scY = ui.NewRegText("scY:")
	ctrl.wx = ui.NewRegText("wX:")
	ctrl.wy = ui.NewRegText("wY:")
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
		ctrl.wx.Label, ctrl.wx.Value,
		ctrl.wy.Label, ctrl.wy.Value,
	)

	c4 := container.New(layout.NewFormLayout(),
		ctrl.control.Label, ctrl.control.Value,
		ctrl.status.Label, ctrl.status.Value,
	)

	regs := container.New(layout.NewGridLayoutWithColumns(4), c1, c2, c3, c4)

	ctrl.ui = container.New(layout.NewBorderLayout(regs, ctrl.sprites, nil, nil), regs, img, ctrl.sprites)

	return ctrl
}

func (ctrl *ppuDebugControl) Widget() fyne.CanvasObject {
	return ctrl.ui
}

func (ctrl *ppuDebugControl) Update() {
	ctrl.x.Update(strconv.Itoa(ctrl.ppu.lx))
	ctrl.y.Update(strconv.Itoa(ctrl.ppu.ly))
	ctrl.wx.Update(strconv.Itoa(ctrl.ppu.wx))
	ctrl.wy.Update(strconv.Itoa(ctrl.ppu.wy))
	ctrl.scX.Update(strconv.Itoa(int(ctrl.ppu.scxNew)))
	ctrl.scY.Update(strconv.Itoa(int(ctrl.ppu.scy)))
	ctrl.control.Update(fmt.Sprintf("%08b", ctrl.ppu.control))
	ctrl.status.Update(fmt.Sprintf("%08b", ctrl.ppu.status))

	var sb strings.Builder
	sb.WriteString("X   Y   Tile Flag        X   Y   Tile Flag\n")
	for i := uint16(0); i < 40; i++ {
		sb.WriteString(fmt.Sprintf("%03d %03d 0x%02X 0b%08b  ",
			ctrl.ppu.oam[i*4+1],
			ctrl.ppu.oam[i*4+0],
			ctrl.ppu.oam[i*4+2],
			ctrl.ppu.oam[i*4+3],
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
				b1 := ctrl.ppu.vRAM[tileAddr]
				b2 := ctrl.ppu.vRAM[tileAddr+1]
				for x_off := 0; x_off < 8; x_off++ {
					c := (b1 & 1) | ((b2 & 1) << 1)
					ctrl.display.Set(x+(7-x_off), y+int(y_off)+block, ctrl.ppu.palette[c])
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

/// *********************************
/// *********************************
/// *********************************

type soundCtrl struct {
	ui  fyne.CanvasObject
	apu *apu

	ch1On, ch2On, ch3On, ch4On *ui.RegText

	regs [][]*ui.RegText
}

func newSoundCtrl(apu *apu) *soundCtrl {
	ctrl := &soundCtrl{
		apu: apu,
	}

	ctrl.ch1On = ui.NewRegText("ch1 On:")
	ctrl.ch2On = ui.NewRegText("ch2 On:")
	ctrl.ch3On = ui.NewRegText("ch3 On:")
	ctrl.ch4On = ui.NewRegText("ch4 On:")

	for c := 0; c < 4; c++ {
		ctrl.regs = append(ctrl.regs, []*ui.RegText{})
		for r := 0; r < 5; r++ {
			ctrl.regs[c] = append(ctrl.regs[c], ui.NewRegText(fmt.Sprintf("MR%d%d:", c, r)))
		}
	}

	cols := []fyne.CanvasObject{
		container.New(layout.NewFormLayout(), ctrl.ch1On.Label, ctrl.ch1On.Value),
		container.New(layout.NewFormLayout(), ctrl.ch2On.Label, ctrl.ch2On.Value),
		container.New(layout.NewFormLayout(), ctrl.ch3On.Label, ctrl.ch3On.Value),
		container.New(layout.NewFormLayout(), ctrl.ch4On.Label, ctrl.ch4On.Value),
	}

	for c := 0; c < 4; c++ {
		for r := 0; r < 5; r++ {
			cols[c].(*fyne.Container).Add(ctrl.regs[c][r].Label)
			cols[c].(*fyne.Container).Add(ctrl.regs[c][r].Value)
		}
	}

	regs := container.New(layout.NewGridLayoutWithColumns(4), cols...)
	ctrl.ui = container.New(layout.NewBorderLayout(regs, nil, nil, nil), regs)

	return ctrl
}

func (ctrl *soundCtrl) Widget() fyne.CanvasObject {
	return ctrl.ui
}

func (ctrl *soundCtrl) Update() {
	if ctrl.apu.channels[0].isOn() {
		ctrl.ch1On.Update("ON")
	} else {
		ctrl.ch1On.Update("Off")
	}
	if ctrl.apu.channels[1].isOn() {
		ctrl.ch2On.Update("ON")
	} else {
		ctrl.ch2On.Update("Off")
	}
	if ctrl.apu.channels[2].isOn() {
		ctrl.ch3On.Update("ON")
	} else {
		ctrl.ch3On.Update("Off")
	}
	if ctrl.apu.channels[3].isOn() {
		ctrl.ch4On.Update("ON")
	} else {
		ctrl.ch4On.Update("Off")
	}

	for c := 0; c < 4; c++ {
		for r := 0; r < 5; r++ {
			ctrl.regs[c][r].Update(fmt.Sprintf("%08b", ctrl.apu.channels[c].getRegister(r)))
		}
	}

	ctrl.ui.Refresh()
}
