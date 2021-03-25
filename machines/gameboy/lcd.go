package gameboy

import (
	"fmt"
	"image"
	"image/color"

	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/emulator"
)

type lcd struct {
	lx, ly, lyc    int
	scx, scy       int
	scxNew, scyNew int

	wy, wx int

	control byte
	status  byte

	dma       uint16
	dmaTarget uint16
	dmaT      uint16
	doDMA     bool

	gbp  []byte
	obp0 []byte
	obp1 []byte

	palette []color.RGBA

	display *image.RGBA
	monitor emulator.Monitor

	bus  cpu.Bus
	vRAM cpu.RAM
	oam  cpu.RAM
}

func newLCD(bus cpu.Bus) *lcd {
	display := image.NewRGBA(image.Rect(0, 0, 160, 144))
	return &lcd{
		gbp:  make([]byte, 4),
		obp0: make([]byte, 4),
		obp1: make([]byte, 4),
		palette: []color.RGBA{
			{0x9b, 0xbc, 0x0f, 0xff},
			{0x8b, 0xac, 0x0f, 0xff},
			{0x30, 0x62, 0x30, 0xff},
			{0x0f, 0x38, 0x0f, 0xff},
		},
		display: display,
		monitor: emulator.NewMonitor(display),
		bus:     bus,
		vRAM:    cpu.NewRAM(make([]byte, 0x2000), 0x1fff),
		oam:     cpu.NewRAM(make([]byte, 0x0100), 0x00ff), //TODO: real size is 0xa0
	}
}

func (lcd *lcd) Tick() {
	if lcd.control&0x80 == 0 {
		return
	}

	lcd.lx++
	if lcd.lx == 456 {
		lcd.drawLine()
		lcd.scy, lcd.scx = lcd.scyNew, lcd.scxNew
		lcd.lx = 0
		lcd.ly++
		if lcd.ly == 153 {
			lcd.ly = 0
			lcd.monitor.FrameDone()
		}

		if lcd.ly == lcd.lyc {
			lcd.status |= 0b00000100
		} else {
			lcd.status &= 0b11111011
		}
	}

	mode := lcd.status & 3
	if lcd.ly > 143 {
		if mode != 1 {
			lcd.bus.Write(0xff0f, 1)
		}
		mode = 1
	} else if lcd.lx < 80 {
		mode = 2
	} else if lcd.lx < 80+168 { // TODO: review sprite count
		mode = 3
	} else if lcd.lx < 80+168+208 { // TODO: review sprite count
		mode = 0
	}
	lcd.status = (lcd.status & 0xfc) | mode

	if lcd.doDMA {
		lcd.dmaTick()
	}
}

func (lcd *lcd) drawLine() {
	r := uint16((lcd.ly + lcd.scy) >> 3)
	l := uint16((lcd.ly + lcd.scy) & 0x07)
	mapBase := uint16(0x1800)
	if lcd.control&0b00001000 != 0 {
		mapBase += 0x0400
	}
	tileBase := uint16(0)
	if lcd.control&0b00010000 == 0 {
		tileBase += 0x0800
	}
	for c := uint16(0); c < 20; c++ {
		mapAddr := mapBase + c + r*32
		tileIdx, _ := lcd.vRAM.ReadPort(mapAddr)
		tileAddr := uint16(tileIdx) * 16
		b1, _ := lcd.vRAM.ReadPort(tileBase + tileAddr + l*2)
		b2, _ := lcd.vRAM.ReadPort(tileBase + tileAddr + l*2 + 1)
		for x_off := 0; x_off < 8; x_off++ {
			color := (b1 & 1) | ((b2 & 1) << 1)
			lcd.display.SetRGBA(int(c*8)+(7-x_off), lcd.ly, lcd.palette[color])
			b1 >>= 1
			b2 >>= 1
		}
	}
}

func (lcd *lcd) dmaTick() {
	// fmt.Printf("lcd.dmaTarget = 0x%04X\n", lcd.dmaTarget)
	if lcd.dmaT == 3 {
		lcd.dmaT = 0
		lcd.bus.Write(lcd.dmaTarget, lcd.bus.Read(lcd.dma))
		lcd.dmaTarget++
		lcd.dma++
		lcd.doDMA = lcd.dmaTarget != 0xfea0
	} else {
		lcd.dmaT++
	}
}

func (lcd *lcd) ReadPort(addr uint16) (byte, bool) {
	switch addr {
	case 0xff40:
		return lcd.control, false

	case 0xff41:
		return lcd.status, false

	case 0xff42:
		return byte(lcd.scy), false

	case 0xff43:
		return byte(lcd.scx), false

	case 0xff44:
		return byte(lcd.ly), false

	case 0xff45:
		return byte(lcd.lyc), false

	case 0xff46:
		return byte(lcd.dma >> 8), false

	case 0xff47:
		res := lcd.gbp[0] << 0
		res |= lcd.gbp[1] << 2
		res |= lcd.gbp[2] << 4
		res |= lcd.gbp[3] << 6
		return res, false

	case 0xff48:
		res := lcd.obp0[0] << 0
		res |= lcd.obp0[1] << 2
		res |= lcd.obp0[2] << 4
		res |= lcd.obp0[3] << 6
		return res, false

	case 0xff49:
		res := lcd.obp1[0] << 0
		res |= lcd.obp1[1] << 2
		res |= lcd.obp1[2] << 4
		res |= lcd.obp1[3] << 6
		return res, false

	case 0xff4A:
		return byte(lcd.wx), false

	case 0xff4B:
		return byte(lcd.wy), false

	case 0xff4c, 0xff4d, 0xff4e, 0xff4f:
		return 0xff, false

	default:
		panic(fmt.Sprintf("[lcd] read invalid addr:0x%04x", addr))
	}
}

func (lcd *lcd) WritePort(addr uint16, data byte) {
	switch addr {
	case 0xff40:
		lcd.control = data

	case 0xff41:
		lcd.status = data

	case 0xff42:
		lcd.scyNew = int(data)

	case 0xff43:
		lcd.scxNew = int(data)

	case 0xff44:

	case 0xff45:
		lcd.lyc = int(data)

	case 0xff46:
		lcd.dma = uint16(data) << 8
		lcd.dmaTarget = 0xfe00
		lcd.doDMA = true

	case 0xff47:
		lcd.gbp[0] = (data & 0b00000011) >> 0
		lcd.gbp[1] = (data & 0b00001100) >> 2
		lcd.gbp[2] = (data & 0b00110000) >> 4
		lcd.gbp[3] = (data & 0b11000000) >> 6

	case 0xff48:
		lcd.obp0[0] = (data & 0b00000011) >> 0
		lcd.obp0[1] = (data & 0b00001100) >> 2
		lcd.obp0[2] = (data & 0b00110000) >> 4
		lcd.obp0[3] = (data & 0b11000000) >> 6

	case 0xff49:
		lcd.obp1[0] = (data & 0b00000011) >> 0
		lcd.obp1[1] = (data & 0b00001100) >> 2
		lcd.obp1[2] = (data & 0b00110000) >> 4
		lcd.obp1[3] = (data & 0b11000000) >> 6

	case 0xff4A:
		lcd.wy = int(data)

	case 0xff4B:
		lcd.wx = int(data)

	case 0xff4c, 0xff4d, 0xff4e, 0xff4f:

	default:
		panic(fmt.Sprintf("[lcd] write invalid addr:0x%04x", addr))
	}
}
