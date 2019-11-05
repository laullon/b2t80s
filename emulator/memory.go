package emulator

type Memory interface {
	PortManager

	LoadRom(idx int, rom []byte)

	GetBlock(start, length uint16) []byte

	GetByte(pos uint16) byte
	PutByte(pos uint16, b byte)

	GetWord(pos uint16) uint16
	PutWord(pos, w uint16)

	DisableSafeMode()

	SetClock(c Clock)
}
