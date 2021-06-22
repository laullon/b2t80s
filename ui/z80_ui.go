package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	im, iff                *RegText

	log    []string
	nextOP string

	ui  gui.HCT
	out gui.Text

	traceFile *os.File
}

func NewZ80UI(cpu z80.Z80, trace bool) gui.GUIObject {
	ctl := &z80UI{regs: cpu.Registers()}
	cpu.SetTracer(ctl)

	ctl.a = NewRegText("A:")
	ctl.f = NewRegText("F:")
	ctl.b = NewRegText("B:")
	ctl.c = NewRegText("C:")
	ctl.d = NewRegText("D:")
	ctl.e = NewRegText("E:")
	ctl.h = NewRegText("H:")
	ctl.l = NewRegText("L:")
	ctl.af = NewRegText("AF:")
	ctl.bc = NewRegText("BC:")
	ctl.de = NewRegText("DE:")
	ctl.hl = NewRegText("HL:")
	ctl.ixh = NewRegText("IXH:")
	ctl.ixl = NewRegText("IXL:")
	ctl.iyh = NewRegText("IYH:")
	ctl.iyl = NewRegText("IYL:")
	ctl.ix = NewRegText("IX:")
	ctl.iy = NewRegText("IY:")
	ctl.sp = NewRegText("SP:")
	ctl.pc = NewRegText("PC:")
	ctl.flag = NewRegText("FLAG:")
	ctl.im = NewRegText("IM:")
	ctl.iff = NewRegText("IFF:")

	flag := NewRegText("")
	flag.Update("SZXHXPNC")
	flag.Update("SZXHXPNC")

	regs := []*RegText{
		ctl.a, ctl.f, ctl.af, ctl.pc,
		ctl.b, ctl.c, ctl.bc, ctl.sp,
		ctl.d, ctl.e, ctl.de, ctl.flag,
		ctl.h, ctl.l, ctl.hl, flag,
		ctl.ixh, ctl.ixl, ctl.ix, ctl.im,
		ctl.iyh, ctl.iyl, ctl.iy, ctl.iff,
	}

	grid := gui.NewHGrid(8, 20, 0)
	for _, reg := range regs {
		grid.Add(reg.Label, reg.Value)
	}

	ctl.out = gui.NewText("")

	ctl.ui = gui.NewVerticalHCT()
	ctl.ui.SetHead(grid, 120)
	ctl.ui.SetCenter(ctl.out)

	if trace {
		f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		ctl.traceFile = f
	}

	return ctl
}

func (ctl *z80UI) GetChildrens() []gui.GUIObject {
	return ctl.ui.GetChildrens()
}

func (ctl *z80UI) Render() {
	ctl.ui.Render()
}

func (ctl *z80UI) Resize(r gui.Rect) {
	ctl.ui.Resize(r)
}

func (ui *z80UI) GetRegisters() string { return "" }
func (ui *z80UI) GetOutput() string    { return "" }

func (ctl *z80UI) Update() {
	af := toHex16(uint16(ctl.regs.A)<<8 | uint16(ctl.regs.F.GetByte()))
	ctl.a.Update(toHex8(ctl.regs.A))
	ctl.f.Update(toHex8(ctl.regs.F.GetByte()))
	ctl.b.Update(toHex8(ctl.regs.B))
	ctl.c.Update(toHex8(ctl.regs.C))
	ctl.d.Update(toHex8(ctl.regs.D))
	ctl.e.Update(toHex8(ctl.regs.E))
	ctl.h.Update(toHex8(ctl.regs.H))
	ctl.l.Update(toHex8(ctl.regs.L))
	ctl.ixh.Update(toHex8(ctl.regs.IXH))
	ctl.ixl.Update(toHex8(ctl.regs.IXL))
	ctl.iyh.Update(toHex8(ctl.regs.IYH))
	ctl.iyl.Update(toHex8(ctl.regs.IYL))
	ctl.af.Update(af)
	ctl.bc.Update(toHex16(ctl.regs.BC.Get()))
	ctl.de.Update(toHex16(ctl.regs.DE.Get()))
	ctl.hl.Update(toHex16(ctl.regs.HL.Get()))
	ctl.ix.Update(toHex16(ctl.regs.IX.Get()))
	ctl.iy.Update(toHex16(ctl.regs.IY.Get()))
	ctl.sp.Update(toHex16(ctl.regs.SP.Get()))
	ctl.pc.Update(toHex16(ctl.regs.PC))
	ctl.im.Update(strconv.Itoa(int(ctl.regs.InterruptsMode)))

	if ctl.regs.IFF1 {
		ctl.iff.Update("true")
	} else {
		ctl.iff.Update("false")
	}

	ctl.flag.Update(fmt.Sprintf("%08b", ctl.regs.F.GetByte()))

	ctl.out.SetText(strings.Join(append(ctl.log, "\n", ctl.nextOP), "\n"))

	// ui.logTxt.Text = strings.Join(append(ui.log, "\n", ui.nextOP), "\n")

}

func (ui *z80UI) DoTrace(on bool) { // TODO: implement
}

func (ctl *z80UI) AppendLastOP(op string) {
	if ctl.traceFile != nil {
		ctl.traceFile.WriteString(op)
		ctl.traceFile.WriteString("\n")
	}
	log := append(ctl.log, op)
	if len(log) > 10 {
		ctl.log = log[1:]
	} else {
		ctl.log = log
	}
}

func (ui *z80UI) SetNextOP(op string) {
	ui.nextOP = op
}

func (ui *z80UI) SetDiss(pc uint16, getMemory func(pc, leng uint16) []byte) {
	panic(-1)
}
