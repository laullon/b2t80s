package gameboy

import (
	"fmt"
	"image"
	"image/color"
	"math/bits"

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
	vRAM []byte
	oam  []byte

	bgBuffer       chan byte
	bgMapAddr      uint16
	bgNextTileAddr uint16

	spriteBuffer    []uint16
	spriteBufferIdx uint16
	spriteCount     uint16

	mode2Ticks int
	mode3Ticks int
}

func newLCD(bus cpu.Bus) *lcd {
	display := image.NewRGBA(image.Rect(0, 0, 160, 144))
	lcd := &lcd{
		gbp:  make([]byte, 4),
		obp0: make([]byte, 4),
		obp1: make([]byte, 4),
		palette: []color.RGBA{
			{0x9b, 0xbc, 0x0f, 0xff},
			{0x8b, 0xac, 0x0f, 0xff},
			{0x30, 0x62, 0x30, 0xff},
			{0x0f, 0x38, 0x0f, 0xff},
		},
		display:  display,
		monitor:  emulator.NewMonitor(display),
		bus:      bus,
		vRAM:     make([]byte, 0x2000),
		oam:      make([]byte, 0x0100),
		bgBuffer: make(chan byte, 100),
	}

	bus.RegisterPort("vram", cpu.PortMask{0b1110_0000_0000_0000, 0b1000_0000_0000_0000}, cpu.NewRAM(lcd.vRAM, 0x1fff))
	bus.RegisterPort("oam", cpu.PortMask{0b1111_1111_0000_0000, 0b1111_1110_0000_0000}, cpu.NewRAM(lcd.oam, 0x00ff))

	return lcd
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
		lcd.mode2Ticks = 0
		lcd.mode3Ticks = 0
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
		lcd.mode2Tick()
	} else if lcd.lx < 80+172 { // TODO: review sprite count
		mode = 3
		lcd.mode3Tick()
	} else if lcd.lx < 80+168+208 { // TODO: review sprite count
		mode = 0
	}
	lcd.status = (lcd.status & 0xfc) | mode

	if lcd.doDMA {
		lcd.dmaTick()
	}
}

func (lcd *lcd) mode2Tick() {
	if lcd.mode2Ticks%2 == 1 {
		if lcd.spriteCount < 10 {
			y := int(lcd.oam[lcd.spriteBufferIdx*4])
			x := int(lcd.oam[lcd.spriteBufferIdx*4+1])
			if (x != 0) && (lcd.ly+16 >= y) && (lcd.ly+16 < y+8) { // TODO: 16 sprites
				lcd.spriteBuffer[lcd.spriteCount] = lcd.spriteBufferIdx * 4
				lcd.spriteCount++
			}
		}
		lcd.spriteBufferIdx++
	} else if lcd.mode2Ticks == 0 {
		lcd.spriteBuffer = make([]uint16, 10)
		lcd.spriteBufferIdx = 0
		lcd.spriteCount = 0
	}
	lcd.mode2Ticks++
}

func (lcd *lcd) mode3Tick() {
	switch lcd.mode3Ticks % 8 {
	case 0:
		if lcd.mode3Ticks == 0 {
			r := uint16(lcd.ly+lcd.scy) >> 3
			lcd.bgMapAddr = 0x1800 + r*32
		}
	case 1:
		l := uint16(lcd.ly+lcd.scy) & 0x07
		tileIdx := lcd.vRAM[lcd.bgMapAddr]
		area := lcd.control & 0b0001_0000
		if area == 1 {
			lcd.bgNextTileAddr = uint16(tileIdx)*16 + l*2
		} else {
			block := tileIdx & 0x80 >> 7
			idx := tileIdx & 0x7f
			if block == 0 {
				lcd.bgNextTileAddr = 0x1000 + uint16(idx)*16 + l*2
			} else {
				lcd.bgNextTileAddr = 0x0800 + uint16(idx)*16 + l*2
			}
		}
		lcd.bgMapAddr++
	case 3:
		lcd.bgBuffer <- lcd.vRAM[lcd.bgNextTileAddr]
	case 5:
		lcd.bgBuffer <- lcd.vRAM[lcd.bgNextTileAddr+1]
	}
	lcd.mode3Ticks++
}

func (lcd *lcd) drawLine() {
	if len(lcd.bgBuffer) == 0 {
		return
	}

	for c := uint16(0); c < 20; c++ {
		b1 := <-lcd.bgBuffer
		b2 := <-lcd.bgBuffer
		for x_off := 0; x_off < 8; x_off++ {
			color := (b1 & 1) | ((b2 & 1) << 1)
			color = lcd.gbp[color]
			lcd.display.SetRGBA(int(c*8)+(7-x_off), lcd.ly, lcd.palette[color])
			b1 >>= 1
			b2 >>= 1
		}
	}

	for len(lcd.bgBuffer) != 0 {
		<-lcd.bgBuffer
	}

	for i := 0; i < int(lcd.spriteCount); i++ {
		sprite := lcd.spriteBuffer[i]
		x := int(lcd.oam[sprite+1]) - 8
		y := uint16(lcd.ly - (int(lcd.oam[sprite]) - 16))
		f := lcd.oam[sprite+3]

		tileIdx := lcd.oam[sprite+2]
		tileAddr := uint16(tileIdx) * 16
		b1 := lcd.vRAM[tileAddr+y*2]
		b2 := lcd.vRAM[tileAddr+y*2+1]
		if f&0b0010_0000 == 0 {
			bits.Reverse8(b1)
			bits.Reverse8(b2)
		}
		for x_off := 0; x_off < 8; x_off++ {
			color := (b1 & 1) | ((b2 & 1) << 1)
			if color != 0 {
				if f&0b0001_0000 == 0 {
					color = lcd.obp0[color]
				} else {
					color = lcd.obp1[color]
				}
				lcd.display.SetRGBA(x+(7-x_off), lcd.ly, lcd.palette[color])
			}
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
