package fdc

import (
	"github.com/laullon/b2t80s/emulator/files"
)

type FDC interface {
	ReadData() byte
	ReadStatus() byte
	WriteData(val byte)
	SetMotor(bool)

	SetDiscA(files.DSK)
	// SetDiscB(files.DSK)
}
