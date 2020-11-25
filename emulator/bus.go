package emulator

type Bus interface {
	SetAddr(uint16)
	GetAddr() uint16

	SetData(byte)
	GetData() byte

	ReadMemory()
	WriteMemory()

	ReadPort()
	WritePort()
}
