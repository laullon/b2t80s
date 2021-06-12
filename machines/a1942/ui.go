package a1942

import (
	"fmt"
	"image/color"

	"github.com/laullon/b2t80s/gui"
)

type characters struct {
	v   *video
	img *gui.Display
	ui  gui.Image
}

var colors = []color.RGBA{
	{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
	{R: 0xff, G: 0x00, B: 0x00, A: 0xff},
	{R: 0xff, G: 0xff, B: 0x00, A: 0xff},
	{R: 0xff, G: 0x00, B: 0xff, A: 0xff},
	{R: 0x00, G: 0xff, B: 0x00, A: 0xff},
	{R: 0x00, G: 0xff, B: 0xff, A: 0xff},
	{R: 0x00, G: 0x00, B: 0xff, A: 0xff},
	{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
}

func newCharactersUI(v *video) *characters {
	ctrl := &characters{}
	ctrl.v = v
	ctrl.img = gui.NewDisplay(gui.Size{W: 16 * 16, H: 16 * 32})
	ctrl.ui = gui.NewDisplayViewer(ctrl.img)

	for row := 0; row < 32; row++ {
		for col := 0; col < 16; col++ {
			tileIdx := col + row*16
			charAddr := tileIdx * 32
			for y := 0; y < 16; y++ {
				for i := 0; i < 2; i++ {
					data1 := v.tilesRom[0][charAddr+y+i*16]
					data2 := v.tilesRom[1][charAddr+y+i*16]
					data3 := v.tilesRom[2][charAddr+y+i*16]
					for x := 0; x < 8; x++ {
						color := data1 & 0b00000001 << 2
						color |= data2 & 0b00000001 << 1
						color |= data3 & 0b00000001 << 0
						ctrl.img.Set(((7 - x) + (i * 8) + col*16), row*16+y, colors[color])
						data1 >>= 1
						data2 >>= 1
						data3 >>= 1
					}
				}
			}
		}
	}
	ctrl.img.Swap()
	return ctrl
}

func RGN_FRAC(num, den uint32) {
	n := (0x80000000 | (((num) & 0x0f) << 27) | (((den) & 0x0f) << 23))
	fmt.Printf("RGN_FRAC(%d,%d) = %032b\n", num, den, n)
}

func (ctrl *characters) Update()                            {}
func (ctrl *characters) Render()                            { ctrl.ui.Render() }
func (ctrl *characters) Resize(r gui.Rect)                  { ctrl.ui.Resize(r) }
func (ctrl *characters) GetMouseTargets() []gui.MouseTarget { return []gui.MouseTarget{} }
