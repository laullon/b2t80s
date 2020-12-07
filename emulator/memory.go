package emulator

type Memory interface {
	GetByte(pos uint16) byte
	PutByte(pos uint16, b byte)
}
