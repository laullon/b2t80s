package z80

type Memory interface {
	Read(pos uint16) byte
	Write(pos uint16, b byte)
}
