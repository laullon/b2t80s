package cpu

type RegPair struct {
	H, L *byte
}

func (reg *RegPair) Get() uint16 {
	return uint16(*reg.H)<<8 | uint16(*reg.L)
}

func (reg *RegPair) Set(hl uint16) {
	*reg.H = byte(hl >> 8)
	*reg.L = byte(hl & 0x00ff)
}
