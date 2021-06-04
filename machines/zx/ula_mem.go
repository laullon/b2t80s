package zx

func (ula *ula) GetBlock(start, length uint16) []byte {
	return ula.memory.GetBlock(start, length)
}

func (ula *ula) GetByte(addr uint16) byte {
	return ula.memory.Read(addr)
}

func (ula *ula) PutByte(addr uint16, b byte) {
	ula.memory.Write(addr, b)
}
