package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/laullon/b2t80s/cpu/lr35902"
)

type lr35902UI struct {
	regs *lr35902.LR35902Registers

	a, f, b, c, d, e, h, l *RegText
	af, bc, de, hl         *RegText
	sp, pc, flag           *RegText
	ier, ifr, ime          *RegText
	registers              RegTable

	log       []string
	nextOP    string
	lastPC    uint16
	getMemory func(pc, leng uint16) []byte

	traceFile *os.File
}

func NewLR35902UI(cpu lr35902.LR35902) Control {
	ui := &lr35902UI{
		regs: cpu.Registers(),
		log:  make([]string, 10),
	}

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
	ui.sp = NewRegText("SP:")
	ui.pc = NewRegText("PC:")
	ui.ier = NewRegText("IE:")
	ui.ifr = NewRegText("IF:")
	ui.ime = NewRegText("IME:")
	ui.flag = NewRegText("FLAG:")
	flag := NewRegText("")
	flag.Update("ZNHC")

	ui.registers = [][]*RegText{
		{ui.a, ui.f, ui.af, ui.pc, ui.ier},
		{ui.b, ui.c, ui.bc, ui.sp, ui.ifr},
		{ui.d, ui.e, ui.de, ui.flag, ui.ime},
		{ui.h, ui.l, ui.hl, flag},
	}

	// dump := widget.NewCheck("Dump", func(on bool) {
	// 	ui.doTrace(on)
	// })

	return ui
}

func (ui *lr35902UI) GetRegisters() string {
	println("[lr35902UI] GetRegisters()")
	af := toHex16(uint16(ui.regs.A)<<8 | uint16(ui.regs.F.GetByte()))
	ui.a.Update(toHex8(ui.regs.A))
	ui.f.Update(toHex8(ui.regs.F.GetByte()))
	ui.b.Update(toHex8(ui.regs.B))
	ui.c.Update(toHex8(ui.regs.C))
	ui.d.Update(toHex8(ui.regs.D))
	ui.e.Update(toHex8(ui.regs.E))
	ui.h.Update(toHex8(ui.regs.H))
	ui.l.Update(toHex8(ui.regs.L))
	ui.af.Update(af)
	ui.bc.Update(toHex16(ui.regs.BC.Get()))
	ui.de.Update(toHex16(ui.regs.DE.Get()))
	ui.hl.Update(toHex16(ui.regs.HL.Get()))
	ui.sp.Update(toHex16(ui.regs.SP.Get()))
	ui.pc.Update(toHex16(ui.regs.PC))
	ui.ifr.Update(fmt.Sprintf("%08b", ui.regs.IF))
	ui.ier.Update(fmt.Sprintf("%08b", ui.regs.IE))
	ui.ime.Update(fmt.Sprintf("%v", ui.regs.IME))
	ui.flag.Update(fmt.Sprintf("%04b", ui.regs.F.GetByte()>>4))
	return ui.registers.Dump()
}

func (ui *lr35902UI) GetOutput() string {
	println("[lr35902UI] GetOutput()")
	var sb strings.Builder
	sb.WriteString("<pre>")
	sb.WriteString(strings.Join(ui.log, "\n"))
	sb.WriteString("</pre>")

	sb.WriteString("<pre>")
	sb.WriteString(ui.nextOP)
	sb.WriteString("</pre>")

	pc := ui.lastPC
	if ui.getMemory != nil {
		sb.WriteString("<pre>")
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
		sb.WriteString("</pre>")
	}

	return sb.String()
}

func (ui *lr35902UI) DoTrace(on bool) { // TODO: implement
}

func (ui *lr35902UI) AppendLastOP(op string) {
	if ui.traceFile != nil {
		ui.traceFile.WriteString(op)
		ui.traceFile.WriteString("\n")
	}
	// println(op)
	// println()
	nLog := append(ui.log, op)
	ui.log = nLog[1:]
}

func (ui *lr35902UI) SetNextOP(op string) {
	ui.nextOP = op
}

func (ui *lr35902UI) SetDiss(pc uint16, getMemory func(pc, leng uint16) []byte) {
	ui.AppendLastOP(ui.nextOP)

	data := getMemory(pc, 4)

	op := lr35902.OPCodes[data[0]]
	ui.nextOP = op.Dump(pc, data)
	pc += uint16(op.Len)
	data = data[op.Len:]

	ui.lastPC = pc
	ui.getMemory = getMemory
}

func (ui *lr35902UI) doTrace(on bool) {
	if on {
		f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		ui.traceFile = f
	} else {
		ui.traceFile.Close()
		ui.traceFile = nil
	}
}
