package cpc

import (
	"fmt"

	// "image"
	"image/color"

	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/gui"
)

type gatearray struct {
	mem  *memory
	crtc *crtc

	screenMode byte
	decode     func(byte) []byte
	ppc        int

	border      bool
	borderColor color.RGBA

	pen     byte
	palette []color.RGBA

	monitor emulator.Monitor
	display *gui.Display

	prevHSync, prevVSync           bool
	hSyncCount, hSyncsInVSyncCount byte

	x, y int
}

var colours = []color.RGBA{
	{0x7f, 0x7f, 0x7f, 0xff},
	{0x7f, 0x7f, 0x7f, 0xff},
	{0x00, 0xff, 0x7f, 0xff},
	{0xff, 0xff, 0x7f, 0xff},
	{0x00, 0x00, 0x7f, 0xff},
	{0xff, 0x00, 0x7f, 0xff},
	{0x00, 0x7f, 0x7f, 0xff},
	{0xff, 0x7f, 0x7f, 0xff},
	{0xff, 0x00, 0x7f, 0xff},
	{0xff, 0xff, 0x7f, 0xff},
	{0xff, 0xff, 0x00, 0xff},
	{0xff, 0xff, 0xff, 0xff},
	{0xff, 0x00, 0x00, 0xff},
	{0xff, 0x00, 0xff, 0xff},
	{0xff, 0x7f, 0x00, 0xff},
	{0xff, 0x7f, 0xff, 0xff},
	{0x00, 0x00, 0x7f, 0xff},
	{0x00, 0xff, 0x7f, 0xff},
	{0x00, 0xff, 0x00, 0xff},
	{0x00, 0xff, 0xff, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0xff, 0xff},
	{0x00, 0x7f, 0x00, 0xff},
	{0x00, 0x7f, 0xff, 0xff},
	{0x7f, 0x00, 0x7f, 0xff},
	{0x7f, 0xff, 0x7f, 0xff},
	{0x7f, 0xff, 0x00, 0xff},
	{0x7f, 0xff, 0xff, 0xff},
	{0x7f, 0x00, 0x00, 0xff},
	{0x7f, 0x00, 0xff, 0xff},
	{0x7f, 0x7f, 0x00, 0xff},
	{0x7f, 0x7f, 0xff, 0xff},
}

func newGateArray(mem *memory, crtc *crtc) *gatearray {
	ga := &gatearray{
		mem:     mem,
		crtc:    crtc,
		palette: make([]color.RGBA, 16),
		display: gui.NewDisplay(gui.Size{960, 312}),

		decode: to1bpp,
		ppc:    16,
	}

	ga.display.ViewSize.W = 768
	ga.display.ViewSize.H = 576
	ga.monitor = emulator.NewMonitor(ga.display)

	return ga
}

func (ga *gatearray) Tick() {
	if !ga.prevHSync && ga.crtc.status.hSync {
		ga.x = 0
		ga.y++
	}

	if !ga.prevVSync && ga.crtc.status.vSync {
		ga.y = 0
		ga.display.ViewPortRect.X = (int32(ga.crtc.regs[3]&0x0f) * 8) * 2
		ga.display.ViewPortRect.Y = 34 / 2
		ga.display.ViewPortRect.W = 384 * 2
		ga.display.ViewPortRect.H = 272
		ga.monitor.FrameDone()
	}

	pixles := 16
	x := ga.x * pixles
	if ga.crtc.status.disPen {
		addr := ga.crtc.status.getAddress()
		cs := ga.decode(ga.mem.getScreenByte(addr))
		cs = append(cs, ga.decode(ga.mem.getScreenByte(addr+1))...)
		for off, c := range cs {
			ga.display.SetRGBA(x+off, ga.y, ga.palette[c])
		}
	} else {
		for i := 0; i < pixles; i++ {
			ga.display.SetRGBA(x+i, ga.y, ga.borderColor)
		}

		// if ga.crtc.status.hSync || ga.crtc.status.vSync {
		// 	for i := 0; i < pixles; i += 2 {
		// 		ga.display.SetRGBA(x+i, ga.y, color.RGBA{0x00, 0x00, 0x00, 0xff})
		// 	}
		// }

		// if ga.x == 0 {
		// 	if ga.screenMode == 0 {
		// 		for i := 0; i < pixles; i++ {
		// 			ga.display.SetRGBA(x+i, ga.y, color.RGBA{0x00, 0xff, 0x00, 0xff})
		// 		}
		// 	} else if ga.screenMode == 1 {
		// 		for i := 0; i < pixles; i++ {
		// 			ga.display.SetRGBA(x+i, ga.y, color.RGBA{0x00, 0x00, 0xff, 0xff})
		// 		}
		// 	}
		// }

		// if ga.x == 1 {
		// 	if ga.hSyncCount == 0 {
		// 		for i := 0; i < pixles; i++ {
		// 			ga.display.SetRGBA(x+i, ga.y, color.RGBA{0xff, 0x00, 0x00, 0xff})
		// 		}
		// 	}
		// }
	}

	if !ga.prevVSync && ga.crtc.status.vSync {
		ga.hSyncsInVSyncCount = 0
	}

	if !ga.prevHSync && ga.crtc.status.hSync {
		ga.hSyncCount++
		if ga.hSyncCount == 52 {
			ga.hSyncCount = 0
			ga.crtc.cpu.Interrupt(true)
		}

		ga.hSyncsInVSyncCount++
		if ga.hSyncsInVSyncCount == 2 {
			if ga.hSyncCount >= 32 {
				ga.crtc.cpu.Interrupt(true)
			}
			ga.hSyncCount = 0
		}
	}

	ga.x++

	ga.prevHSync = ga.crtc.status.hSync
	ga.prevVSync = ga.crtc.status.vSync

	// if ga.crtc.status.disPen {
	// 	ga.mem.clock.AddTStates(4)
	// }

}

func bit2value(b, bit, v byte) byte {
	if b&(1<<bit) != 0 {
		return v
	}
	return 0
}

func to4bpp(b byte) []byte {
	pixel1 := bit2value(b, 1, 0b1000) | bit2value(b, 5, 0b0100) | bit2value(b, 3, 0b0010) | bit2value(b, 7, 0b0001)
	pixel2 := bit2value(b, 0, 0b1000) | bit2value(b, 4, 0b0100) | bit2value(b, 2, 0b0010) | bit2value(b, 6, 0b0001)
	return []byte{pixel1, pixel1, pixel1, pixel1, pixel2, pixel2, pixel2, pixel2}
}

func to2bpp(b byte) []byte {
	pixel1 := bit2value(b, 3, 0b10) | bit2value(b, 7, 0b01)
	pixel2 := bit2value(b, 2, 0b10) | bit2value(b, 6, 0b01)
	pixel3 := bit2value(b, 1, 0b10) | bit2value(b, 5, 0b01)
	pixel4 := bit2value(b, 0, 0b10) | bit2value(b, 4, 0b01)
	return []byte{pixel1, pixel1, pixel2, pixel2, pixel3, pixel3, pixel4, pixel4}
}

func to1bpp(b byte) []byte {
	pixel1 := bit2value(b, 7, 1)
	pixel2 := bit2value(b, 6, 1)
	pixel3 := bit2value(b, 5, 1)
	pixel4 := bit2value(b, 4, 1)
	pixel5 := bit2value(b, 3, 1)
	pixel6 := bit2value(b, 2, 1)
	pixel7 := bit2value(b, 1, 1)
	pixel8 := bit2value(b, 0, 1)
	return []byte{pixel1, pixel2, pixel3, pixel4, pixel5, pixel6, pixel7, pixel8}
}

func (ga *gatearray) ReadPort(port uint16) (byte, bool) { return 0, true }

func (ga *gatearray) WritePort(port uint16, data byte) {
	f := data >> 6
	if f == 0 {
		ga.border = data>>4&1 == 1
		if !ga.border {
			ga.pen = data & 0xf
		}
	} else if f == 1 {
		if !ga.border {
			ga.palette[ga.pen] = colours[data&0x1f]
		} else {
			if ga.borderColor != colours[data&0x1f] {
				ga.borderColor = colours[data&0x1f]
				// println("ga.borderColor:", data&0x1f, ga.y)
				// draw.Draw(ga.display, ga.display.Bounds(), &image.Uniform{ga.borderColor}, image.ZP, draw.Src)
			}
		}
	} else if f == 2 {
		ga.mem.lowerRomEnable = data&0b00000100 == 0
		ga.mem.upperRomEnable = data&0b00001000 == 0

		if data&0x10 != 0 {
			ga.crtc.cpu.Interrupt(false)
			ga.hSyncCount = 0
		}

		screenMode := data & 0b00000011
		if ga.screenMode != screenMode {
			ga.screenMode = screenMode
			switch screenMode {
			case 0:
				ga.decode = to4bpp
			case 1:
				ga.decode = to2bpp
			case 2:
				ga.decode = to1bpp
			default:
				// panic(screenMode)
			}
			// println("screenMode", screenMode, ga.y)
		}
		// println("[ga]", "lowerRomEnable:", ga.mem.lowerRomEnable, "upperRomEnable:", ga.mem.upperRomEnable, "screenMode:", ga.screenMode)
	} else if f == 3 {
		ga.mem.Paging(data)
	} else {
		panic(fmt.Sprintf("unsupported f:%d", f))
	}
}
