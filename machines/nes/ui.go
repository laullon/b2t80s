package nes

import (
	"image"
	"image/draw"
	"time"

	"fyne.io/fyne/v2"
	canvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	layout "fyne.io/fyne/v2/layout"
	widget "fyne.io/fyne/v2/widget"
	"github.com/laullon/b2t80s/emulator"
)

type ppuDebugControl struct {
	ui        fyne.CanvasObject
	show      *widget.Button
	ppu       *ppu
	display   *image.RGBA
	container *fyne.Container
}

func newPalleteControl(ppu *ppu) *ppuDebugControl {
	ctrl := &ppuDebugControl{
		ppu:     ppu,
		display: image.NewRGBA(image.Rect(0, 0, (64*8)+4+84, (64*8)+2)),
	}

	ctrl.show = widget.NewButton("ppu debug", ctrl.doShow)

	ctrl.ui = container.New(layout.NewHBoxLayout(),
		widget.NewToolbarSeparator().ToolbarObject(),
		ctrl.show,
	)

	img := canvas.NewImageFromImage(ctrl.display)
	img.SetMinSize(fyne.NewSize((64*8)+4+84, (64*8)+2))
	img.ScaleMode = canvas.ImageScalePixels

	ctrl.container = fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), img)

	return ctrl
}

func (ctrl *ppuDebugControl) Widget() fyne.CanvasObject {
	return ctrl.ui
}

func (ctrl *ppuDebugControl) doShow() {
	secondaryWindow := emulator.App.NewWindow("secondary")
	secondaryWindow.SetContent(ctrl.container)
	secondaryWindow.Show()

	wait := time.Duration(500 * time.Millisecond)
	ticker := time.NewTicker(wait)
	go func() {
		for range ticker.C {
			ctrl.Update()
			ctrl.container.Refresh()
		}
	}()

}

func (ctrl *ppuDebugControl) Update() {
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
					color := uint16(0x3f00) | uint16(palette)<<2 | c
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
