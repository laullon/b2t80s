package gameboy

import (
	"fmt"
	"image"
	"image/color"
	"math/bits"

	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/ui"
)

type ppu struct {
	debugger cpu.DebuggerCallbacks

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

	display *ui.Display
	monitor emulator.Monitor

	bus  cpu.Bus
	vRAM []byte
	oam  []byte

	bgBuffer       chan byte
	bgMapAddr      uint16
	bgMapC         byte
	bgNextTileAddr uint16

	spriteBuffer    []uint16
	spriteBufferIdx uint16
	spriteCount     uint16

	mode2Ticks int
	mode3Ticks int
}

func newPPU(bus cpu.Bus) *ppu {
	display := ui.NewDisplay(image.Rect(0, 0, 160, 144))
	ppu := &ppu{
		gbp:  []byte{0, 1, 2, 3},
		obp0: []byte{0, 1, 2, 3},
		obp1: []byte{0, 1, 2, 3},
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

	bus.RegisterPort("vram", cpu.PortMask{0b1110_0000_0000_0000, 0b1000_0000_0000_0000}, cpu.NewRAM(ppu.vRAM, 0x1fff))
	bus.RegisterPort("oam", cpu.PortMask{0b1111_1111_0000_0000, 0b1111_1110_0000_0000}, cpu.NewRAM(ppu.oam, 0x00ff))

	return ppu
}

func (ppu *ppu) Tick() {
	if ppu.control&0x80 == 0 {
		return
	}

	ppu.lx++
	if ppu.lx == 456 {
		ppu.drawLine()
		ppu.drawSprites()
		ppu.drawWin()

		ppu.scy, ppu.scx = ppu.scyNew, ppu.scxNew
		ppu.lx = 0
		ppu.mode2Ticks = 0
		ppu.mode3Ticks = 0
		ppu.ly++
		if ppu.debugger != nil {
			if ppu.debugger.EvalLine() {
				for x := 0; x < 160; x++ {
					ppu.display.SetRGBA(x, ppu.ly, color.RGBA{0xff, 0, 0, 0xff})
				}
				ppu.monitor.FrameDone()
			}
		}
		if ppu.ly == 153 {
			ppu.ly = 0
			ppu.monitor.FrameDone()
			if ppu.debugger != nil {
				ppu.debugger.EvalFrame()
			}
		}

		if ppu.ly == ppu.lyc {
			ppu.status |= 0b00000100
			if ppu.status&0b0100_0000 != 0 {
				ppu.bus.Write(0xff0f, 0b00010)
			}
		} else {
			ppu.status &= 0b11111011
		}
	}

	mode := ppu.status & 3
	if ppu.ly > 143 {
		if mode != 1 {
			ppu.bus.Write(0xff0f, 1)
		}
		mode = 1
	} else if ppu.lx < 80 {
		mode = 2
		ppu.mode2Tick()
	} else if ppu.lx < 80+172 { // TODO: review sprite count
		mode = 3
		ppu.mode3Tick()
	} else if ppu.lx < 80+168+208 { // TODO: review sprite count
		mode = 0
	}
	ppu.status = (ppu.status & 0xfc) | mode

	if ppu.doDMA {
		ppu.dmaTick()
	}
}

func (ppu *ppu) mode2Tick() {
	if ppu.mode2Ticks%2 == 1 {
		if ppu.spriteCount < 10 {
			y := int(ppu.oam[ppu.spriteBufferIdx*4])
			x := int(ppu.oam[ppu.spriteBufferIdx*4+1])
			if (x != 0) && (ppu.ly+16 >= y) && (ppu.ly+16 < y+8) { // TODO: 16 sprites
				ppu.spriteBuffer[ppu.spriteCount] = ppu.spriteBufferIdx * 4
				ppu.spriteCount++
			}
		}
		ppu.spriteBufferIdx++
	} else if ppu.mode2Ticks == 0 {
		ppu.spriteBuffer = make([]uint16, 10)
		ppu.spriteBufferIdx = 0
		ppu.spriteCount = 0
	}
	ppu.mode2Ticks++
}

func (ppu *ppu) mode3Tick() {
	switch ppu.mode3Ticks % 8 {
	case 0:
		if ppu.mode3Ticks == 0 {
			r := (uint16(ppu.ly+ppu.scy) >> 3) & 31
			ppu.bgMapC = byte(ppu.scx >> 3)
			ppu.bgMapAddr = 0x1800 + r*32
			if ppu.control&0b0000_1000 != 0 {
				ppu.bgMapAddr += 0x0400
			}
		}
	case 1:
		l := uint16(ppu.ly+ppu.scy) & 0x07
		tileIdx := ppu.vRAM[ppu.bgMapAddr+uint16(ppu.bgMapC&31)]
		area := ppu.control & 0b0001_0000 >> 4
		if area == 1 {
			ppu.bgNextTileAddr = uint16(tileIdx)*16 + l*2
		} else {
			block := tileIdx & 0x80 >> 7
			idx := tileIdx & 0x7f
			if block == 0 {
				ppu.bgNextTileAddr = 0x1000 + uint16(idx)*16 + l*2
			} else {
				ppu.bgNextTileAddr = 0x0800 + uint16(idx)*16 + l*2
			}
		}
		ppu.bgMapC++
	case 3:
		ppu.bgBuffer <- ppu.vRAM[ppu.bgNextTileAddr]
	case 5:
		ppu.bgBuffer <- ppu.vRAM[ppu.bgNextTileAddr+1]
	}
	ppu.mode3Ticks++
}

func (ppu *ppu) drawLine() {
	if len(ppu.bgBuffer) == 0 {
		return
	}

	if ppu.control&0b0000_0001 != 0 {
		scx := ppu.scx & 7
		for c := uint16(0); c < 21; c++ {
			b1 := <-ppu.bgBuffer
			b2 := <-ppu.bgBuffer
			for x_off := 0; x_off < 8; x_off++ {
				color := (b1 & 1) | ((b2 & 1) << 1)
				color = ppu.gbp[color]
				ppu.display.SetRGBA(int(c*8)+(7-x_off)-scx, ppu.ly, ppu.palette[color])
				b1 >>= 1
				b2 >>= 1
			}
		}
	}

	for len(ppu.bgBuffer) != 0 {
		<-ppu.bgBuffer
	}
}

func (ppu *ppu) drawWin() {
	if ppu.control&0b0010_0000 == 0 {
		return
	}
	wy := ppu.ly - ppu.wy
	if wy >= 0 {
		r := uint16(wy >> 3)
		l := uint16(wy) & 7
		for x := 0; x < 160; x++ {
			wy := x - ppu.wx
			if wy >= 0 && wy&7 == 0 {
				c := uint16(wy >> 3)
				mapAddr := 0x1800 + c + r*32
				if ppu.control&0b0100_0000 != 0 {
					mapAddr += 0x400
				}
				tileIdx := uint16(ppu.vRAM[mapAddr])
				tileAddr := uint16(tileIdx)*16 + l*2
				b1 := ppu.vRAM[tileAddr]
				b2 := ppu.vRAM[tileAddr+1]
				for x_off := 0; x_off < 8; x_off++ {
					color := (b1 & 1) | ((b2 & 1) << 1)
					color = ppu.gbp[color]
					ppu.display.SetRGBA((x-7)+(7-x_off), ppu.ly, ppu.palette[color])
					b1 >>= 1
					b2 >>= 1
				}
			}
		}
	}
}

func (ppu *ppu) drawSprites() {
	for i := 0; i < int(ppu.spriteCount); i++ {
		sprite := ppu.spriteBuffer[i]
		x := int(ppu.oam[sprite+1]) - 8
		y := uint16(ppu.ly - (int(ppu.oam[sprite]) - 16))
		f := ppu.oam[sprite+3]
		if f&0b0100_0000 != 0 {
			y = 7 - y
		}

		tileIdx := ppu.oam[sprite+2]
		tileAddr := uint16(tileIdx) * 16
		b1 := ppu.vRAM[tileAddr+y*2]
		b2 := ppu.vRAM[tileAddr+y*2+1]
		if f&0b0010_0000 != 0 {
			b1 = bits.Reverse8(b1)
			b2 = bits.Reverse8(b2)
		}
		for x_off := 0; x_off < 8; x_off++ {
			color := (b1 & 1) | ((b2 & 1) << 1)
			if color != 0 {
				if f&0b0001_0000 == 0 {
					color = ppu.obp0[color]
				} else {
					color = ppu.obp1[color]
				}
				ppu.display.SetRGBA(x+(7-x_off), ppu.ly, ppu.palette[color])
			}
			b1 >>= 1
			b2 >>= 1
		}
	}
}

func (ppu *ppu) dmaTick() {
	// fmt.Printf("ppu.dmaTarget = 0x%04X\n", ppu.dmaTarget)
	if ppu.dmaT == 3 {
		ppu.dmaT = 0
		ppu.bus.Write(ppu.dmaTarget, ppu.bus.Read(ppu.dma))
		ppu.dmaTarget++
		ppu.dma++
		ppu.doDMA = ppu.dmaTarget != 0xfea0
	} else {
		ppu.dmaT++
	}
}

func (ppu *ppu) ReadPort(addr uint16) (byte, bool) {
	switch addr {
	case 0xff40:
		return ppu.control, false

	case 0xff41:
		return ppu.status, false

	case 0xff42:
		return byte(ppu.scy), false

	case 0xff43:
		return byte(ppu.scx), false

	case 0xff44:
		return byte(ppu.ly), false

	case 0xff45:
		return byte(ppu.lyc), false

	case 0xff46:
		return byte(ppu.dma >> 8), false

	case 0xff47:
		res := ppu.gbp[0] << 0
		res |= ppu.gbp[1] << 2
		res |= ppu.gbp[2] << 4
		res |= ppu.gbp[3] << 6
		return res, false

	case 0xff48:
		res := ppu.obp0[0] << 0
		res |= ppu.obp0[1] << 2
		res |= ppu.obp0[2] << 4
		res |= ppu.obp0[3] << 6
		return res, false

	case 0xff49:
		res := ppu.obp1[0] << 0
		res |= ppu.obp1[1] << 2
		res |= ppu.obp1[2] << 4
		res |= ppu.obp1[3] << 6
		return res, false

	case 0xff4A:
		return byte(ppu.wx), false

	case 0xff4B:
		return byte(ppu.wy), false

	case 0xff4c, 0xff4d, 0xff4e, 0xff4f:
		return 0xff, false

	default:
		panic(fmt.Sprintf("[ppu] read invalid addr:0x%04x", addr))
	}
}

func (ppu *ppu) WritePort(addr uint16, data byte) {
	switch addr {
	case 0xff40:
		ppu.control = data

	case 0xff41:
		ppu.status = data

	case 0xff42:
		ppu.scyNew = int(data)

	case 0xff43:
		ppu.scxNew = int(data)

	case 0xff44:

	case 0xff45:
		ppu.lyc = int(data)

	case 0xff46:
		ppu.dma = uint16(data) << 8
		ppu.dmaTarget = 0xfe00
		ppu.doDMA = true

	case 0xff47:
		ppu.gbp[0] = (data & 0b00000011) >> 0
		ppu.gbp[1] = (data & 0b00001100) >> 2
		ppu.gbp[2] = (data & 0b00110000) >> 4
		ppu.gbp[3] = (data & 0b11000000) >> 6

	case 0xff48:
		ppu.obp0[0] = (data & 0b00000011) >> 0
		ppu.obp0[1] = (data & 0b00001100) >> 2
		ppu.obp0[2] = (data & 0b00110000) >> 4
		ppu.obp0[3] = (data & 0b11000000) >> 6

	case 0xff49:
		ppu.obp1[0] = (data & 0b00000011) >> 0
		ppu.obp1[1] = (data & 0b00001100) >> 2
		ppu.obp1[2] = (data & 0b00110000) >> 4
		ppu.obp1[3] = (data & 0b11000000) >> 6

	case 0xff4A:
		ppu.wy = int(data)

	case 0xff4B:
		ppu.wx = int(data)

	case 0xff4c, 0xff4d, 0xff4e, 0xff4f:

	default:
		panic(fmt.Sprintf("[ppu] write invalid addr:0x%04x", addr))
	}
}
