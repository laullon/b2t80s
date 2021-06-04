package nes

import (
	"image"
	"image/draw"
	"strconv"

	"github.com/laullon/b2t80s/ui"
)

type ppuDebugControl struct {
	ppu     *ppu
	display *image.RGBA

	v, h, x, y, p *ui.RegText
}

func newPalleteControl(ppu *ppu) *ppuDebugControl {
	ctrl := &ppuDebugControl{
		ppu:     ppu,
		display: image.NewRGBA(image.Rect(0, 0, (64*8)+4+84, (64*8)+2)),
	}

	// img := canvas.NewImageFromImage(ctrl.display)
	// img.SetMinSize(fyne.NewSize((64*8)+4+84, (64*8)+2))
	// img.ScaleMode = canvas.ImageScalePixels

	ctrl.v = ui.NewRegText("V:")
	ctrl.h = ui.NewRegText("H:")
	ctrl.x = ui.NewRegText("Scroll X:")
	ctrl.y = ui.NewRegText("Scroll Y:")
	ctrl.p = ui.NewRegText("Page:")

	// c1 := container.New(layout.NewFormLayout(),
	// 	ctrl.h.Label, ctrl.h.Value,
	// 	ctrl.v.Label, ctrl.v.Value,
	// )

	// c2 := container.New(layout.NewFormLayout(),
	// 	ctrl.x.Label, ctrl.x.Value,
	// 	ctrl.y.Label, ctrl.y.Value,
	// )

	// c3 := container.New(layout.NewFormLayout(),
	// 	ctrl.p.Label, ctrl.p.Value,
	// )

	// regs := container.New(layout.NewGridLayoutWithColumns(3), c1, c2, c3)

	// ctrl.ui = fyne.NewContainerWithLayout(layout.NewBorderLayout(regs, nil, nil, nil), regs, img)

	return ctrl
}

func (ui *ppuDebugControl) GetRegisters() string { return "" }
func (ui *ppuDebugControl) GetOutput() string    { return "" }

func (ctrl *ppuDebugControl) Update() {
	ctrl.v.Update(strconv.Itoa(ctrl.ppu.v))
	ctrl.h.Update(strconv.Itoa(ctrl.ppu.h))
	ctrl.x.Update(strconv.Itoa(int(ctrl.ppu.scrollX)))
	ctrl.y.Update(strconv.Itoa(int(ctrl.ppu.scrollY)))
	ctrl.p.Update(strconv.Itoa(int(ctrl.ppu.nameTableBase)))

	draw.Draw(ctrl.display, image.Rect((64*8)+4, 0, (64*8)+4+84, 20), &image.Uniform{colors[ctrl.ppu.bus.Read(0x3f00)&0x3f]}, image.ZP, draw.Src)
	for palette := uint16(0); palette < 6; palette++ {
		for color := 0; color < 4; color++ {
			y := int((palette + 1) * 22)
			x := color * 22
			c := uint16(0x3f00) | (palette << 2) | uint16(color)
			draw.Draw(ctrl.display, image.Rect((64*8)+4+x, y, (64*8)+4+x+20, y+20), &image.Uniform{colors[ctrl.ppu.bus.Read(c)&0x3f]}, image.ZP, draw.Src)
		}
	}

	for row := 0; row < 64; row++ {
		for y := 0; y < 8; y++ {
			for col := 0; col < 64; col++ {
				charAddr := ctrl.ppu.charAddrs[col][row]
				char := uint16(ctrl.ppu.bus.Read(charAddr))

				patternAddr := ctrl.ppu.patternBase | char<<4 | uint16(y)
				pattern0 := ctrl.ppu.bus.Read(patternAddr)
				pattern1 := ctrl.ppu.bus.Read(patternAddr | 0x08)

				attrAddr := ctrl.ppu.attrAddrs[col][row]
				b := ctrl.ppu.blocks[col][row]
				attr := ctrl.ppu.bus.Read(attrAddr)
				palette := (attr >> (b * 2)) & 0x03

				for x := 0; x < 8; x++ {
					c := uint16(((pattern0 & 0x80) >> 7) | ((pattern1 & 0x80) >> 6))
					color := uint16(0x3f00)
					if c != 0 {
						color |= uint16(palette)<<2 | c
					}
					pattern0 <<= 1
					pattern1 <<= 1
					imgX := int(col*8) + x
					if col > 31 {
						imgX += 2
					}
					imgY := int(row*8) + y
					if row > 31 {
						imgY += 2
					}
					ctrl.display.SetRGBA(imgX, imgY, colors[ctrl.ppu.bus.Read(color)&0x3f])
				}
			}
		}
	}
}
