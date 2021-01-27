package pokey

//Write
const (
	audf1W  = 0x00
	audc1W  = 0x01
	audf2W  = 0x02
	audc2W  = 0x03
	audf3W  = 0x04
	audc3W  = 0x05
	audf4W  = 0x06
	audc4W  = 0x07
	audctlW = 0x08
	stimerW = 0x09
	skrestW = 0x0A
	potgoW  = 0x0B
	seroutW = 0x0D
	irqenW  = 0x0E
	skctlW  = 0x0F
)

// Read
const (
	pot0R   = 0x00
	pot1R   = 0x01
	pot2R   = 0x02
	pot3R   = 0x03
	pot4R   = 0x04
	pot5R   = 0x05
	pot6R   = 0x06
	pot7R   = 0x07
	allpotR = 0x08
	kbcodeR = 0x09
	randomR = 0x0A
	serinR  = 0x0D
	irqstR  = 0x0E
	skstatR = 0x0F
)

type channel struct {
	noise       byte
	forceVolume bool
	volume      byte
	freq        byte
}

func (c *channel) control(control byte) {
	c.noise = control >> 5
	c.forceVolume = control&0x10 != 0
	c.volume = control & 0x0f
}

// TODO: remove?
func NewPokey() *Pokey {
	return &Pokey{}
}

type Pokey struct {
	audc1 channel
	audc2 channel
	audc3 channel
	audc4 channel

	P0 bool
	P1 bool
	P2 bool
	P3 bool
	P4 bool
	P5 bool
	P6 bool
	P7 bool
}

func (p *Pokey) allpotR() byte {
	res := byte(0xff)
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
	return res
}

func (p *Pokey) ReadPort(port uint16) (byte, bool) {
	// fmt.Printf("[pokey.readPort]-> port:0x%04X \n", port)
	var res byte
	switch port & 0xf {
	case allpotR:
		res = p.allpotR()
	case pot0R, pot1R, pot2R, pot3R, pot4R, pot5R, pot6R, pot7R, skstatR, 0xb:
		res = 0
	default:
		panic(port & 0x0f)
	}
	// fmt.Printf("[pokey.readPort]-> port:0x%04X res:0x%02X \n", port, res)
	return res, false
}

func (p *Pokey) WritePort(port uint16, data byte) {
	// fmt.Printf("[pokey.WritePort]-> port:0x%04X data:0x%02X \n", port, data)
	switch port & 0xf {
	case audc1W:
		p.audc1.control(data)
	case audf1W:
		p.audc1.freq = data
	case audc2W:
		p.audc2.control(data)
	case audf2W:
		p.audc2.freq = data
	case audc3W:
		p.audc3.control(data)
	case audf3W:
		p.audc3.freq = data
	case audc4W:
		p.audc4.control(data)
	case audf4W:
		p.audc4.freq = data

	case skctlW, potgoW, 0x08:
	default:
		panic(port & 0x0f)
	}
}
