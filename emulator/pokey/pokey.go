package pokey

// TODO: remove?
func NewPokey() *Pokey {
	return &Pokey{}
}

type Pokey struct {
	P0 bool
	P1 bool
	P2 bool
	P3 bool
	P4 bool
	P5 bool
	P6 bool
	P7 bool
}

func (p *Pokey) ReadPort(port uint16) (byte, bool) {
	res := byte(0xff)
	if port&0x0f == 0x08 {
		if !p.P0 {
			res ^= 0b00000001
		}
		if !p.P1 {
			res ^= 0b00000010
		}
		if !p.P2 {
			res ^= 0b00000100
		}
		if !p.P3 {
			res ^= 0b00001000
		}
		if !p.P4 {
			res ^= 0b00010000
		}
		if !p.P5 {
			res ^= 0b00100000
		}
		if !p.P6 {
			res ^= 0b01000000
		}
		if !p.P7 {
			res ^= 0b10000000
		}
	}

	// fmt.Printf("[readPort]-> port:0x%04X res:0x%02X \n", port, res)

	return res, false
}

func (p *Pokey) WritePort(port uint16, data byte) {}
