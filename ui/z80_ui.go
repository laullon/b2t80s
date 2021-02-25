package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/laullon/b2t80s/cpu/z80"
)

type z80UI struct {
	regs   *z80.Z80Registers
	widget *fyne.Container

	a, x, y *canvas.Text
	sp      *canvas.Text
	pc      *canvas.Text
	ps      *canvas.Text
}

func NewZ80UI(r *z80.Z80Registers) Control {
	ui := &z80UI{regs: r}

	ui.a = canvas.NewText("ui.a", color.Black)
	ui.x = canvas.NewText("ui.x", color.Black)
	ui.y = canvas.NewText("ui.y", color.Black)
	ui.sp = canvas.NewText("ui.sp", color.Black)
	ui.pc = canvas.NewText("ui.pc", color.Black)
	ui.ps = canvas.NewText("ui.ps", color.Black)

	a := canvas.NewText("A:", color.Black)
	x := canvas.NewText("X:", color.Black)
	y := canvas.NewText("Y:", color.Black)
	sp := canvas.NewText("SP:", color.Black)
	pc := canvas.NewText("PC:", color.Black)
	ps := canvas.NewText("PS:", color.Black)

	ui.widget = container.New(layout.NewFormLayout(),
		a, ui.a,
		x, ui.x,
		y, ui.y,
		sp, ui.sp,
		pc, ui.pc,
		ps, ui.ps)

	return ui
}

func (ui *z80UI) Widget() fyne.CanvasObject {
	return ui.widget
}

func (ui *z80UI) Update() {
}
