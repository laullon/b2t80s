package cpc

import (
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/ui"
)

func (cpc *cpc) UIControls() []ui.Control {
	var res []ui.Control
	if cpc.cassette != nil {
		res = append(res, ui.NewCasseteControl(cpc.cassette, !*emulator.LoadSlow))
	}
	res = append(res, ui.NewVolumenControl(cpc.sound.SetVolume))
	return res
}
