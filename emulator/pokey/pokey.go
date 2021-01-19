package pokey

import "github.com/laullon/b2t80s/emulator"

func NewPokey() emulator.PortManager {
	return &pokey{}
}

type pokey struct {
}

func (p *pokey) ReadPort(port uint16) (byte, bool) { return 0xff, false }
func (p *pokey) WritePort(port uint16, data byte)  {}
