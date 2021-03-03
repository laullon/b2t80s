package ui

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/laullon/b2t80s/cpu/m6502"
)

type m6502UI struct {
	regs   *m6502.Registers
	widget *fyne.Container

	a, x, y *regText
	sp      *regText
	pc      *regText
	ps      *regText

	logTxt *widget.Label
	log    []string
	nextOP string
}

func NewM6502UI(cpu m6502.M6502) Control {
	ui := &m6502UI{regs: cpu.Registers()}
	cpu.SetTracer(ui)

	ui.logTxt = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	ui.a = NewRegText("ui.a")
	ui.x = NewRegText("ui.x")
	ui.y = NewRegText("ui.y")
	ui.sp = NewRegText("ui.sp")
	ui.pc = NewRegText("ui.pc")
	ui.ps = NewRegText("ui.ps")

	a := NewRegText("A:")
	x := NewRegText("X:")
	y := NewRegText("Y:")
	sp := NewRegText("SP:")
	pc := NewRegText("PC:")
	ps := NewRegText("PS:")

	c1 := container.New(layout.NewFormLayout(),
		a.txt, ui.a.txt,
		x.txt, ui.x.txt,
		y.txt, ui.y.txt,
	)

	c2 := container.New(layout.NewFormLayout(),
		sp.txt, ui.sp.txt,
		pc.txt, ui.pc.txt,
		ps.txt, ui.ps.txt,
	)

	regs := container.New(layout.NewGridLayoutWithColumns(2), c1, c2)

	ui.widget = container.New(layout.NewVBoxLayout(), regs, ui.logTxt)

	return ui
}

func (ui *m6502UI) Widget() fyne.CanvasObject {
	return ui.widget
}

func (ui *m6502UI) Update() {
	ui.a.update(toHex8(ui.regs.A))
	ui.x.update(toHex8(ui.regs.X))
	ui.y.update(toHex8(ui.regs.Y))
	ui.sp.update(toHex8(ui.regs.SP))
	ui.pc.update(toHex16(ui.regs.PC))
	ui.ps.update(ui.regs.PS.String())
	ui.logTxt.Text = strings.Join(append(ui.log, "\n", ui.nextOP), "\n")
	ui.widget.Refresh()
}

func (ui *m6502UI) AppendLastOP(op string) {
	println(op)
	log := append(ui.log, op)
	if len(log) > 10 {
		ui.log = log[1:]
	} else {
		ui.log = log
	}
}

func (ui *m6502UI) SetNextOP(op string) {
	ui.nextOP = op
}

func toHex8(v uint8) string {
	n := "0" + strconv.FormatUint(uint64(v), 16)
	return "0x" + n[len(n)-2:]
}
func toHex16(v uint16) string {
	n := "000" + strconv.FormatUint(uint64(v), 16)
	return "0x" + n[len(n)-4:]
}
