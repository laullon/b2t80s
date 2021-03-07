package ui

import (
	"os"
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

	a, x, y *RegText
	sp      *RegText
	pc      *RegText
	ps      *RegText

	logTxt *widget.Label
	log    []string
	nextOP string

	tracefile *os.File
}

func NewM6502UI(cpu m6502.M6502) Control {
	ui := &m6502UI{regs: cpu.Registers()}
	cpu.SetTracer(ui)

	ui.logTxt = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	ui.a = NewRegText("A:")
	ui.x = NewRegText("X:")
	ui.y = NewRegText("Y:")
	ui.sp = NewRegText("SP:")
	ui.pc = NewRegText("PC:")
	ui.ps = NewRegText("PS:")

	c1 := container.New(layout.NewFormLayout(),
		ui.a.Label, ui.a.Value,
		ui.x.Label, ui.x.Value,
		ui.y.Label, ui.y.Value,
	)

	c2 := container.New(layout.NewFormLayout(),
		ui.sp.Label, ui.sp.Value,
		ui.pc.Label, ui.pc.Value,
		ui.ps.Label, ui.ps.Value,
	)

	regs := container.New(layout.NewGridLayoutWithColumns(2), c1, c2)

	dump := widget.NewCheck("Dump", func(on bool) {
		ui.doTrace(on)
	})

	ui.widget = container.New(layout.NewVBoxLayout(), dump, regs, ui.logTxt)

	return ui
}

func (ui *m6502UI) Widget() fyne.CanvasObject {
	return ui.widget
}

func (ui *m6502UI) doTrace(on bool) {
	if on {
		f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		ui.tracefile = f
	} else {
		ui.tracefile.Close()
		ui.tracefile = nil
	}
}

func (ui *m6502UI) Update() {
	ui.a.Update(toHex8(ui.regs.A))
	ui.x.Update(toHex8(ui.regs.X))
	ui.y.Update(toHex8(ui.regs.Y))
	ui.sp.Update(toHex8(ui.regs.SP))
	ui.pc.Update(toHex16(ui.regs.PC))
	ui.ps.Update(ui.regs.PS.String())
	ui.logTxt.Text = strings.Join(append(ui.log, "\n", ui.nextOP), "\n")
	ui.widget.Refresh()
}

func (ui *m6502UI) AppendLastOP(op string) {
	if ui.tracefile != nil {
		ui.tracefile.WriteString(op)
		ui.tracefile.WriteString("\n")
	}
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
