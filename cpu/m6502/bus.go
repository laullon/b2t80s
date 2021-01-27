package m6502

import "github.com/laullon/b2t80s/emulator"

type Bus interface {
	Write(addr uint16, data uint8)
	Read(addr uint16) uint8
	RegisterPort(mask emulator.PortMask, manager emulator.PortManager)
}
