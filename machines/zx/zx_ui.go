package zx

import (
	"github.com/laullon/b2t80s/machines"
	"github.com/laullon/b2t80s/ui"
)

func (zx *zx) UIControls() []ui.Control {
	var res []ui.Control
	if zx.cassette != nil {
		res = append(res, ui.NewCasseteControl(zx.cassette, !*machines.LoadSlow))
	}
	res = append(res, ui.NewVolumenControl(zx.sound))
	return res
}
