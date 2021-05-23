package cpc

import (
	"github.com/laullon/b2t80s/gui"
)

func (cpc *cpc) UIControls() []gui.GUIObject {
	var res []gui.GUIObject
	// if cpc.cassette != nil {
	// 	res = append(res, ui.NewCasseteControl(cpc.cassette, !*emulator.LoadSlow))
	// }
	// res = append(res, ui.NewVolumenControl(cpc.sound.SetVolume))
	return res
}
