package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/laullon/b2t80s/cpu/lr35902"
	"github.com/laullon/b2t80s/gui"
)

type lr35902UI struct {
	regs *lr35902.LR35902Registers

	a, f, b, c, d, e, h, l *RegText
	af, bc, de, hl         *RegText
	sp, pc, flag           *RegText
	ier, ifr, ime          *RegText

	ui  gui.HCT
	out gui.Text

	log       []string
	nextOP    string
	lastPC    uint16
	getMemory func(pc, leng uint16) []byte

	traceFile *os.File
}

func NewLR35902UI(cpu lr35902.LR35902) gui.GUIObject {
	ctl := &lr35902UI{
		regs: cpu.Registers(),
		log:  make([]string, 10),
	}

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
	ctl.sp = NewRegText("SP:")
	ctl.pc = NewRegText("PC:")
	ctl.ier = NewRegText("IE:")
	ctl.ifr = NewRegText("IF:")
	ctl.ime = NewRegText("IME:")
	ctl.flag = NewRegText("FLAG:")
	flag := NewRegText("")
	flag.Update("ZNHC")

	regs := []*RegText{
		ctl.a, ctl.f, ctl.af, ctl.pc, ctl.ier,
		ctl.b, ctl.c, ctl.bc, ctl.sp, ctl.ifr,
		ctl.d, ctl.e, ctl.de, ctl.flag, ctl.ime,
		ctl.h, ctl.l, ctl.hl, flag,
	}

	grid := gui.NewHGrid(10, 20, 0)
	for _, reg := range regs {
		grid.Add(reg.Label, reg.Value)
	}

	ctl.out = gui.NewText("")

	ctl.ui = gui.NewVerticalHCT()
	ctl.ui.SetHead(grid, 80)
	ctl.ui.SetCenter(ctl.out)

	// dump := widget.NewCheck("Dump", func(on bool) {
	// 	ui.doTrace(on)
	// })

	return ctl
}

func (ctl *lr35902UI) GetChildrens() []gui.GUIObject {
	return ctl.ui.GetChildrens()
}

func (ctl *lr35902UI) Render() {
	ctl.ui.Render()
}

func (ctl *lr35902UI) Resize(r gui.Rect) {
	ctl.ui.Resize(r)
}

func (ctl *lr35902UI) Update() {
	af := toHex16(uint16(ctl.regs.A)<<8 | uint16(ctl.regs.F.GetByte()))
	ctl.a.Update(toHex8(ctl.regs.A))
	ctl.f.Update(toHex8(ctl.regs.F.GetByte()))
	ctl.b.Update(toHex8(ctl.regs.B))
	ctl.c.Update(toHex8(ctl.regs.C))
	ctl.d.Update(toHex8(ctl.regs.D))
	ctl.e.Update(toHex8(ctl.regs.E))
	ctl.h.Update(toHex8(ctl.regs.H))
	ctl.l.Update(toHex8(ctl.regs.L))
	ctl.af.Update(af)
	ctl.bc.Update(toHex16(ctl.regs.BC.Get()))
	ctl.de.Update(toHex16(ctl.regs.DE.Get()))
	ctl.hl.Update(toHex16(ctl.regs.HL.Get()))
	ctl.sp.Update(toHex16(ctl.regs.SP.Get()))
	ctl.pc.Update(toHex16(ctl.regs.PC))
	ctl.ifr.Update(fmt.Sprintf("%08b", ctl.regs.IF))
	ctl.ier.Update(fmt.Sprintf("%08b", ctl.regs.IE))
	ctl.ime.Update(fmt.Sprintf("%v", ctl.regs.IME))
	ctl.flag.Update(fmt.Sprintf("%04b", ctl.regs.F.GetByte()>>4))

	ctl.out.SetText(ctl.getOutput())
}

func (ui *lr35902UI) getOutput() string {
	var sb strings.Builder
	sb.WriteString(strings.Join(ui.log, "\n"))

	sb.WriteString("\n\n")
	sb.WriteString(ui.nextOP)
	sb.WriteString("\n\n")

	pc := ui.lastPC
	if ui.getMemory != nil {
		data := ui.getMemory(pc, 40)
		diss := make([]string, 10)
		for i := 0; (len(data) > 4) && (i < 10); i++ {
			op := lr35902.OPCodes[data[0]]
			if op != nil {
				diss[i] = op.Dump(pc, data)
				pc += uint16(op.Len)
				data = data[op.Len:]
			}
		}
		sb.WriteString(strings.Join(diss, "\n"))
	}
	return sb.String()
}

func (ctl *lr35902UI) DoTrace(on bool) { // TODO: implement
}

func (ctl *lr35902UI) AppendLastOP(op string) {
	if ctl.traceFile != nil {
		ctl.traceFile.WriteString(op)
		ctl.traceFile.WriteString("\n")
	}
	// println(op)
	// println()
	nLog := append(ctl.log, op)
	ctl.log = nLog[1:]
}

func (ctl *lr35902UI) SetNextOP(op string) {
	ctl.nextOP = op
}

func (ctl *lr35902UI) SetDiss(pc uint16, getMemory func(pc, leng uint16) []byte) {
	ctl.AppendLastOP(ctl.nextOP)

	data := getMemory(pc, 4)

	op := lr35902.OPCodes[data[0]]
	ctl.nextOP = op.Dump(pc, data)
	pc += uint16(op.Len)
	data = data[op.Len:]

	ctl.lastPC = pc
	ctl.getMemory = getMemory
}

func (ctl *lr35902UI) doTrace(on bool) {
	if on {
		f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		ctl.traceFile = f
	} else {
		ctl.traceFile.Close()
		ctl.traceFile = nil
	}
}
