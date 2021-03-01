package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/laullon/b2t80s/cpu/z80"
)

type z80UI struct {
	regs   *z80.Z80Registers
	widget *fyne.Container

	a, f, b, c, d, e, h, l *regText
	af, bc, de, hl         *regText
	ixh, ixl, iyh, iyl     *regText
	ix, iy                 *regText
	sp, pc, flag           *regText

	logTxt *widget.Label
	log    []string
	nextOP string
}

func NewZ80UI(cpu z80.Z80) Control {
	ui := &z80UI{regs: cpu.Registers()}
	cpu.SetTracer(ui)

	ui.a = NewRegText("ui.a")
	ui.f = NewRegText("ui.f")
	ui.b = NewRegText("ui.b")
	ui.c = NewRegText("ui.c")
	ui.d = NewRegText("ui.d")
	ui.e = NewRegText("ui.e")
	ui.h = NewRegText("ui.h")
	ui.l = NewRegText("ui.l")
	ui.af = NewRegText("ui.af")
	ui.bc = NewRegText("ui.bc")
	ui.de = NewRegText("ui.de")
	ui.hl = NewRegText("ui.hl")
	ui.ixh = NewRegText("ui.ixh")
	ui.ixl = NewRegText("ui.ixl")
	ui.iyh = NewRegText("ui.iyh")
	ui.iyl = NewRegText("ui.iyl")
	ui.ix = NewRegText("ui.ix")
	ui.iy = NewRegText("ui.iy")
	ui.sp = NewRegText("ui.sp")
	ui.pc = NewRegText("ui.pc")
	ui.flag = NewRegText("ui.flag")

	a := NewRegText("  A:")
	f := NewRegText("  F:")
	b := NewRegText("  B:")
	c := NewRegText("  C:")
	d := NewRegText("  D:")
	e := NewRegText("  E:")
	h := NewRegText("  H:")
	l := NewRegText("  L:")
	af := NewRegText("  AF:")
	bc := NewRegText("  BC:")
	de := NewRegText("  DE:")
	hl := NewRegText("  HL:")
	ixh := NewRegText("IXH:")
	ixl := NewRegText("IXL:")
	iyh := NewRegText("IYH:")
	iyl := NewRegText("IYL:")
	ix := NewRegText("  IX:")
	iy := NewRegText("  IY:")
	sp := NewRegText(" SP:")
	pc := NewRegText(" PC:")
	flag := NewRegText("FLAG:")

	c1 := container.New(layout.NewFormLayout(),
		a.txt, ui.a.txt,
		b.txt, ui.b.txt,
		d.txt, ui.d.txt,
		h.txt, ui.h.txt,
		ixh.txt, ui.ixh.txt,
		iyh.txt, ui.iyh.txt,
		pc.txt, ui.pc.txt,
	)
	c2 := container.New(layout.NewFormLayout(),
		f.txt, ui.f.txt,
		c.txt, ui.c.txt,
		e.txt, ui.e.txt,
		l.txt, ui.l.txt,
		ixl.txt, ui.ixl.txt,
		iyl.txt, ui.iyl.txt,
		sp.txt, ui.sp.txt,
	)
	c3 := container.New(layout.NewFormLayout(),
		af.txt, ui.af.txt,
		bc.txt, ui.bc.txt,
		de.txt, ui.de.txt,
		hl.txt, ui.hl.txt,
		ix.txt, ui.ix.txt,
		iy.txt, ui.iy.txt,
		flag.txt, ui.flag.txt,
	)

	regs := container.New(layout.NewGridLayoutWithColumns(3), c1, c2, c3)

	ui.logTxt = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	ui.widget = container.New(layout.NewVBoxLayout(), regs, ui.logTxt)

	return ui
}

func (ui *z80UI) Widget() fyne.CanvasObject {
	return ui.widget
}

func (ui *z80UI) Update() {
	af := toHex16(uint16(ui.regs.A)<<8 | uint16(ui.regs.F.GetByte()))
	ui.a.update(toHex8(ui.regs.A))
	ui.f.update(toHex8(ui.regs.F.GetByte()))
	ui.b.update(toHex8(ui.regs.B))
	ui.c.update(toHex8(ui.regs.C))
	ui.d.update(toHex8(ui.regs.D))
	ui.e.update(toHex8(ui.regs.E))
	ui.h.update(toHex8(ui.regs.H))
	ui.l.update(toHex8(ui.regs.L))
	ui.ixh.update(toHex8(ui.regs.IXH))
	ui.ixl.update(toHex8(ui.regs.IXL))
	ui.iyh.update(toHex8(ui.regs.IYH))
	ui.iyl.update(toHex8(ui.regs.IYL))
	ui.af.update(af)
	ui.bc.update(toHex16(ui.regs.BC.Get()))
	ui.de.update(toHex16(ui.regs.DE.Get()))
	ui.hl.update(toHex16(ui.regs.HL.Get()))
	ui.ix.update(toHex16(ui.regs.IX.Get()))
	ui.iy.update(toHex16(ui.regs.IY.Get()))
	ui.sp.update(toHex16(ui.regs.SP.Get()))
	ui.pc.update(toHex16(ui.regs.PC))
	ui.flag.update(af)

	ui.logTxt.Text = strings.Join(append(ui.log, "\n", ui.nextOP), "\n")

	ui.widget.Refresh()
}

func (ui *z80UI) AppendLastOP(op string) {
	log := append(ui.log, op)
	if len(log) > 10 {
		ui.log = log[1:]
	} else {
		ui.log = log
	}
}

func (ui *z80UI) SetNextOP(op string) {
	ui.nextOP = op
}
