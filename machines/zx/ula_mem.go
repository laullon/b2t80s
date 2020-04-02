package zx

func (ula *ula) GetBlock(start, length uint16) []byte {
	return ula.memory.GetBlock(start, length)
}

func (ula *ula) GetByte(addr uint16) byte {
	for ula.io {
		ula.memory.clock.AddTStates(1)
	}
	return ula.memory.GetByte(addr)
}

func (ula *ula) PutByte(addr uint16, b byte) {
	for ula.io {
		ula.memory.clock.AddTStates(1)
	}
	ula.memory.PutByte(addr, b)
}
