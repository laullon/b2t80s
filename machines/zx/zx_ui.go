package zx

import (
	"github.com/laullon/b2t80s/gui"
)

func (zx *zx) UIControls() []gui.GUIObject {
	var res []gui.GUIObject
	// if zx.cassette != nil {
	// 	res = append(res, ui.NewCasseteControl(zx.cassette, !*emulator.LoadSlow))
	// }
	// res = append(res, ui.NewVolumenControl(zx.sound.SetVolume))
	return res
}
