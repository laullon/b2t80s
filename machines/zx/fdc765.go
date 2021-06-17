package zx

import (
	"fmt"

	"github.com/laullon/b2t80s/emulator/fdc"
)

type fdc765 struct {
	chip fdc.FDC
}

func NewZXFDC765() *fdc765 {
	return &fdc765{
		chip: fdc.New765(),
	}
}

func (fdc *fdc765) ReadPort(port uint16) byte {
	switch port {
	case 0x2ffd:
		return fdc.chip.ReadStatus()
	case 0x3ffd:
		return fdc.chip.ReadData()
	}
	return 0
}

func (fdc *fdc765) WritePort(port uint16, data byte) {
	switch port {
	case 0x1ffd:
		fdc.chip.SetMotor(data&4 != 0)
	case 0x3ffd:
		fdc.chip.WriteData(data)
	default:
		panic(fmt.Sprintf("0x%04X", port))
	}
}
