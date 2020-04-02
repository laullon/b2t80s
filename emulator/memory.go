package emulator

type Memory interface {
	PortManager // TODO: remove

	GetBlock(start, length uint16) []byte // TODO: just for debug ?

	GetByte(pos uint16) byte
	PutByte(pos uint16, b byte)
}
