package ui

import (
	"github.com/laullon/b2t80s/cpu/z80"
	"github.com/laullon/b2t80s/gui"
)

type z80UI struct {
	regs *z80.Z80Registers

	a, f, b, c, d, e, h, l *RegText
	af, bc, de, hl         *RegText
	ixh, ixl, iyh, iyl     *RegText
	ix, iy                 *RegText
	sp, pc, flag           *RegText

	log    []string
	nextOP string
}

func NewZ80UI(cpu z80.Z80) gui.GUIObject {
	ui := &z80UI{regs: cpu.Registers()}
	cpu.SetTracer(ui)

	ui.a = NewRegText("A:")
	ui.f = NewRegText("F:")
	ui.b = NewRegText("B:")
	ui.c = NewRegText("C:")
	ui.d = NewRegText("D:")
	ui.e = NewRegText("E:")
	ui.h = NewRegText("H:")
	ui.l = NewRegText("L:")
	ui.af = NewRegText("AF:")
	ui.bc = NewRegText("BC:")
	ui.de = NewRegText("DE:")
	ui.hl = NewRegText("HL:")
	ui.ixh = NewRegText("IXH:")
	ui.ixl = NewRegText("IXL:")
	ui.iyh = NewRegText("IYH:")
	ui.iyl = NewRegText("IYL:")
	ui.ix = NewRegText("IX:")
	ui.iy = NewRegText("IY:")
	ui.sp = NewRegText("SP:")
	ui.pc = NewRegText("PC:")
	ui.flag = NewRegText("FLAG:")

	// c1 := container.New(layout.NewFormLayout(),
	// 	ui.a.Label, ui.a.Value,
	// 	ui.b.Label, ui.b.Value,
	// 	ui.d.Label, ui.d.Value,
	// 	ui.h.Label, ui.h.Value,
	// 	ui.ixh.Label, ui.ixh.Value,
	// 	ui.iyh.Label, ui.iyh.Value,
	// 	ui.pc.Label, ui.pc.Value,
	// )
	// c2 := container.New(layout.NewFormLayout(),
	// 	ui.f.Label, ui.f.Value,
	// 	ui.c.Label, ui.c.Value,
	// 	ui.e.Label, ui.e.Value,
	// 	ui.l.Label, ui.l.Value,
	// 	ui.ixl.Label, ui.ixl.Value,
	// 	ui.iyl.Label, ui.iyl.Value,
	// 	ui.sp.Label, ui.sp.Value,
	// )
	// c3 := container.New(layout.NewFormLayout(),
	// 	ui.af.Label, ui.af.Value,
	// 	ui.bc.Label, ui.bc.Value,
	// 	ui.de.Label, ui.de.Value,
	// 	ui.hl.Label, ui.hl.Value,
	// 	ui.ix.Label, ui.ix.Value,
	// 	ui.iy.Label, ui.iy.Value,
	// 	ui.flag.Label, ui.flag.Value,
	// )

	// regs := container.New(layout.NewGridLayoutWithColumns(3), c1, c2, c3)

	// ui.logTxt = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	// ui.widget = container.New(layout.NewVBoxLayout(), regs, ui.logTxt)

	return ui
}

func (ui *z80UI) Render() {
}

func (ui *z80UI) Resize(r gui.Rect) {
}

func (ui *z80UI) GetRegisters() string { return "" }
func (ui *z80UI) GetOutput() string    { return "" }

func (ui *z80UI) Update() {
	af := toHex16(uint16(ui.regs.A)<<8 | uint16(ui.regs.F.GetByte()))
	ui.a.Update(toHex8(ui.regs.A))
	ui.f.Update(toHex8(ui.regs.F.GetByte()))
	ui.b.Update(toHex8(ui.regs.B))
	ui.c.Update(toHex8(ui.regs.C))
	ui.d.Update(toHex8(ui.regs.D))
	ui.e.Update(toHex8(ui.regs.E))
	ui.h.Update(toHex8(ui.regs.H))
	ui.l.Update(toHex8(ui.regs.L))
	ui.ixh.Update(toHex8(ui.regs.IXH))
	ui.ixl.Update(toHex8(ui.regs.IXL))
	ui.iyh.Update(toHex8(ui.regs.IYH))
	ui.iyl.Update(toHex8(ui.regs.IYL))
	ui.af.Update(af)
	ui.bc.Update(toHex16(ui.regs.BC.Get()))
	ui.de.Update(toHex16(ui.regs.DE.Get()))
	ui.hl.Update(toHex16(ui.regs.HL.Get()))
	ui.ix.Update(toHex16(ui.regs.IX.Get()))
	ui.iy.Update(toHex16(ui.regs.IY.Get()))
	ui.sp.Update(toHex16(ui.regs.SP.Get()))
	ui.pc.Update(toHex16(ui.regs.PC))
	ui.flag.Update(af)

	// ui.logTxt.Text = strings.Join(append(ui.log, "\n", ui.nextOP), "\n")

}

func (ui *z80UI) DoTrace(on bool) { // TODO: implement
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

func (ui *z80UI) SetDiss(pc uint16, getMemory func(pc, leng uint16) []byte) {
	panic(-1)
}
