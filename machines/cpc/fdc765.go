package cpc

import (
	"fmt"

	"github.com/laullon/b2t80s/emulator/fdc"
)

type fdc765 struct {
	chip fdc.FDC
}

func NewCPCFDC765() *fdc765 {
	return &fdc765{
		chip: fdc.New765(),
	}
}

func (fdc *fdc765) ReadPort(port uint16) (byte, bool) {
	if (port & (1 << 10)) == 0 {
		f := ((port & (1 << 8)) >> (8 - 1)) | (port & 0x01)
		switch f {
		case 2:
			return fdc.chip.ReadStatus(), false
		case 3:
			return fdc.chip.ReadData(), false
		default:
			panic(f)
		}
	} else {
		panic(fmt.Sprintf("port: 0x%04X", port))
	}
}

func (fdc *fdc765) WritePort(port uint16, data byte) {
	switch port {
	case 0xFA7E:
		fdc.chip.SetMotor(data == 1)
	case 0xFB7F:
		fdc.chip.WriteData(data)
	default:
		// fmt.Printf("port: 0x%04X data: 0x%02X \n", port, data)
	}
}
