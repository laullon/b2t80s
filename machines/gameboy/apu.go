package gameboy

type apu struct {
	data []byte
}

func newAPU() *apu {
	return &apu{make([]byte, 0x40)}
}

func (apu *apu) ReadPort(addr uint16) (byte, bool) {
	return apu.data[addr&0xff], false
}

func (apu *apu) WritePort(addr uint16, data byte) {
	apu.data[addr&0xff] = data
}
