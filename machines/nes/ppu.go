package nes

import (
	"image"
	"image/color"

	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
)

type ppu struct {
	cpu     emulator.CPU
	bus     m6502.Bus
	display *image.RGBA
	monitor emulator.Monitor
	h, v    int

	addr    uint16
	addrInc uint16

	ram           *ram
	nameTableBase uint16
	patternBase   uint16

	palette *ram

	enableNMI bool
}

// 1662607*3,2 / 50 / 341 = 312,043542522
// 32 x 30 = 256 x 240

func newPPU(bus m6502.Bus, cpu emulator.CPU) *ppu {
	display := image.NewRGBA(image.Rect(0, 0, 256, 240))
	ppu := &ppu{
		cpu:     cpu,
		bus:     bus,
		display: display,
		monitor: emulator.NewMonitor(display),
		ram:     &ram{data: make([]byte, 0x1000), mask: 0x0fff},
		palette: &ram{data: make([]byte, 0x20), mask: 0x1f},
	}

	// ram
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111_000000000000, Value: 0b0010_000000000000}, ppu.ram)

	// palette
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111_1111_0000_0000, Value: 0b0011_1111_0000_0000}, ppu.palette)

	return ppu
}

var colors = []color.RGBA{
	{84, 84, 84, 255}, {0, 30, 116, 255}, {8, 16, 144, 255}, {48, 0, 136, 255}, {68, 0, 100, 255}, {92, 0, 48, 255}, {84, 4, 0, 255}, {60, 24, 0, 255}, {32, 42, 0, 255}, {8, 58, 0, 255}, {0, 64, 0, 255}, {0, 60, 0, 255}, {0, 50, 60, 255}, {0, 0, 0, 255}, {0, 0, 0, 255}, {0, 0, 0, 255},
	{152, 150, 152, 255}, {8, 76, 196, 255}, {48, 50, 236, 255}, {92, 30, 228, 255}, {136, 20, 176, 255}, {160, 20, 100, 255}, {152, 34, 32, 255}, {120, 60, 0, 255}, {84, 90, 0, 255}, {40, 114, 0, 255}, {8, 124, 0, 255}, {0, 118, 40, 255}, {0, 102, 120, 255}, {0, 0, 0, 255}, {0, 0, 0, 255}, {0, 0, 0, 255},
	{236, 238, 236, 255}, {76, 154, 236, 255}, {120, 124, 236, 255}, {176, 98, 236, 255}, {228, 84, 236, 255}, {236, 88, 180, 255}, {236, 106, 100, 255}, {212, 136, 32, 255}, {160, 170, 0, 255}, {116, 196, 0, 255}, {76, 208, 32, 255}, {56, 204, 108, 255}, {56, 180, 204, 255}, {60, 60, 60, 255}, {0, 0, 0, 255}, {0, 0, 0, 255},
	{236, 238, 236, 255}, {168, 204, 236, 255}, {188, 188, 236, 255}, {212, 178, 236, 255}, {236, 174, 236, 255}, {236, 174, 212, 255}, {236, 180, 176, 255}, {228, 196, 144, 255}, {204, 210, 120, 255}, {180, 222, 120, 255}, {168, 226, 144, 255}, {152, 226, 180, 255}, {160, 214, 228, 255}, {160, 162, 160, 255}, {0, 0, 0, 255}, {0, 0, 0, 255},
}

func (ppu *ppu) Tick() {
	for i := 0; i < 16; i++ {
		if (ppu.h & 0x07) == 0 {
			cCol := (uint16(ppu.h) >> 3)
			cRow := (uint16(ppu.v) >> 3)

			if (cCol < 32) && (cRow < 30) {
				cCol &= 0x1f
				cRow &= 0x1f

				bCol := cCol >> 2 & 0x07
				bRow := cRow >> 2 & 0x07

				rCol := cCol >> 1 & 0x01
				rRow := cRow >> 1 & 0x01
				region := (rRow << 1) | rCol

				y := uint16(ppu.v) & 0x007

				charAddr := ppu.nameTableBase | (cRow << 5) | cCol
				char := uint16(ppu.bus.Read(charAddr))
				// if cRow == 0 && y == 0 {
				// 	fmt.Printf("0x%04X - %d \n", charAddr, char)
				// }

				patternAddr := ppu.patternBase | char<<4 | y
				pattern0 := ppu.bus.Read(patternAddr)
				pattern1 := ppu.bus.Read(patternAddr | 0x08)

				attrAddr := ppu.nameTableBase | 0x03c0 | (bRow << 3) | bCol
				attr := ppu.bus.Read(attrAddr)
				palette := (attr >> (region * 2)) & 0x03

				for i := 0; i < 8; i++ {
					color := uint16(0x3f00)
					color |= uint16(palette) << 2
					color |= uint16(((pattern0 & 0x80) >> 7) | ((pattern1 & 0x80) >> 6))
					pattern0 <<= 1
					pattern1 <<= 1
					ppu.display.Set(ppu.h+i, ppu.v, colors[ppu.bus.Read(color)])
				}
			}
		}

		ppu.h++
		if ppu.h == 341 {
			ppu.h = 0
			ppu.v++
			if ppu.v == 241 && ppu.enableNMI {
				ppu.cpu.NMI(true)
			}
			if ppu.v == 312 {
				ppu.v = 0
				ppu.monitor.FrameDone()
				// panic(-1)
			}
		}
	}
}

func (ppu *ppu) ReadPort(addr uint16) (byte, bool) {
	// fmt.Printf("[ppu] read 0x%04X \n", addr)
	return 0xff, false
}

func (ppu *ppu) WritePort(addr uint16, data byte) {
	// fmt.Printf("[ppu] write 0x%04X 0x%02x\n", addr, data)
	switch addr & 0x07 {
	case 0:
		ppu.nameTableBase = 0x2000 | (uint16(data&0x3) << 10)
		ppu.patternBase = 0x1000 * (uint16(data&0x10) >> 4)
		// fmt.Printf("[ppu] write -> nameTableBase:0x%04X data:%08b  \n", ppu.nameTableBase, data)
		ppu.enableNMI = data&0x80 == 0x80
		if data&0x04 == 0 {
			ppu.addrInc = 1
		} else {
			ppu.addrInc = 32
		}

	case 6:
		ppu.addr <<= 8
		ppu.addr |= uint16(data)
	case 7:
		// fmt.Printf("[ppu] write -> addr:0x%04X data:%v  \n", ppu.addr, data)
		ppu.bus.Write(ppu.addr, data)
		ppu.addr += ppu.addrInc
	}
}
