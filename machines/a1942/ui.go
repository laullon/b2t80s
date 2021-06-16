package a1942

import (
	"image/color"

	"github.com/laullon/b2t80s/gui"
)

type characters struct {
	v   *video
	img *gui.Display
	ui  gui.Image
}

func newCharactersUI(v *video) *characters {
	ctrl := &characters{}
	ctrl.v = v
	ctrl.img = gui.NewDisplay(gui.Size{W: 16 * 1, H: 16 * 1})
	ctrl.ui = gui.NewDisplayViewer(ctrl.img)

	// for row := 0; row < 1; row++ {
	// 	for col := 0; col < 1; col++ {
	col := 0
	row := 0
	tileIdx := 0x3b
	charAddr := tileIdx * 64
	for y := 0; y < 16; y++ {
		for i := 0; i < 4; i++ {
			idx := charAddr + y<<1 + i&0b10<<4 + i&1
			data1 := ctrl.v.spritesRom[0][idx]
			data2 := ctrl.v.spritesRom[1][idx]
			for x := 0; x < 4; x++ {
				c := ((data1 & 0x01) >> 0) << 3
				c |= ((data1 & 0x10) >> 4) << 2
				c |= ((data2 & 0x01) >> 0) << 1
				c |= ((data2 & 0x10) >> 4) << 0
				_x := (3 - x) + (i * 4) + col*16
				_y := row*16 + y
				ctrl.img.Set(_y, 15-_x, color.Gray{Y: c << 4})
				data1 >>= 1
				data2 >>= 1
			}
		}
	}
	// 	}
	// }

	ctrl.img.Swap()
	return ctrl
}

func (ctrl *characters) Update() {

}

func (ctrl *characters) Render()           { ctrl.ui.Render() }
func (ctrl *characters) Resize(r gui.Rect) { ctrl.ui.Resize(r) }

func (ctrl *characters) GetChildrens() []gui.GUIObject { return []gui.GUIObject{} }
