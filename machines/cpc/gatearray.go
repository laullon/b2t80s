package cpc

import (
	"fmt"
	"image"
	"image/color"

	"golang.org/x/image/draw"
)

type gatearray struct {
	mem  *memory
	crtc *crtc

	cycles uint32
	clock  uint32

	screenMode byte
	decode     func(byte) []byte
	ppc        int

	border      bool
	borderColor color.RGBA

	pen     byte
	palette []color.RGBA

	display       *image.RGBA
	displayFrame  image.Rectangle
	displayScaled *image.RGBA
}

var colours = []color.RGBA{
	color.RGBA{0x7f, 0x7f, 0x7f, 0xff},
	color.RGBA{0x7f, 0x7f, 0x7f, 0xff},
	color.RGBA{0x00, 0xff, 0x7f, 0xff},
	color.RGBA{0xff, 0xff, 0x7f, 0xff},
	color.RGBA{0x00, 0x00, 0x7f, 0xff},
	color.RGBA{0xff, 0x00, 0x7f, 0xff},
	color.RGBA{0x00, 0x7f, 0x7f, 0xff},
	color.RGBA{0xff, 0x7f, 0x7f, 0xff},
	color.RGBA{0xff, 0x00, 0x7f, 0xff},
	color.RGBA{0xff, 0xff, 0x7f, 0xff},
	color.RGBA{0xff, 0xff, 0x00, 0xff},
	color.RGBA{0xff, 0xff, 0xff, 0xff},
	color.RGBA{0xff, 0x00, 0x00, 0xff},
	color.RGBA{0xff, 0x00, 0xff, 0xff},
	color.RGBA{0xff, 0x7f, 0x00, 0xff},
	color.RGBA{0xff, 0x7f, 0xff, 0xff},
	color.RGBA{0x00, 0x00, 0x7f, 0xff},
	color.RGBA{0x00, 0xff, 0x7f, 0xff},
	color.RGBA{0x00, 0xff, 0x00, 0xff},
	color.RGBA{0x00, 0xff, 0xff, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0xff, 0xff},
	color.RGBA{0x00, 0x7f, 0x00, 0xff},
	color.RGBA{0x00, 0x7f, 0xff, 0xff},
	color.RGBA{0x7f, 0x00, 0x7f, 0xff},
	color.RGBA{0x7f, 0xff, 0x7f, 0xff},
	color.RGBA{0x7f, 0xff, 0x00, 0xff},
	color.RGBA{0x7f, 0xff, 0xff, 0xff},
	color.RGBA{0x7f, 0x00, 0x00, 0xff},
	color.RGBA{0x7f, 0x00, 0xff, 0xff},
	color.RGBA{0x7f, 0x7f, 0x00, 0xff},
	color.RGBA{0x7f, 0x7f, 0xff, 0xff},
}

func newGateArray(mem *memory, crtc *crtc) *gatearray {
	return &gatearray{
		mem:           mem,
		crtc:          crtc,
		palette:       make([]color.RGBA, 16),
		displayScaled: image.NewRGBA(image.Rect(0, 0, 768, 576)),
		display:       image.NewRGBA(image.Rect(0, 0, 640, 200)),
		displayFrame:  image.Rect((768-640)/2, (576-400)/2, 640+(768-640)/2, 400+(576-400)/2),

		decode: to2bpp,
		ppc:    8,
	}
}

func (ga *gatearray) Tick() {
	clock := ga.cycles / 4
	ga.cycles++

	if ga.clock == clock {
		return
	}
	ga.clock = clock

	if ga.crtc.status.disPen {
		page := ga.crtc.status.ma >> 14
		pos := ga.crtc.status.ma & 0x3fff
		cs := ga.decode(ga.mem.banks[page][pos])
		cs = append(cs, ga.decode(ga.mem.banks[page][(pos+1)])...)
		for off, c := range cs {
			ga.display.Set((ga.crtc.counters.h*ga.ppc)+off, ga.crtc.counters.raster, ga.palette[c])
		}
	}
}

func (ga *gatearray) FrameEnded() {
	ga.cycles = 0
	draw.NearestNeighbor.Scale(ga.displayScaled, ga.displayFrame, ga.display, ga.display.Bounds(), draw.Over, nil)
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
	return []byte{pixel1, pixel2}
}

func to2bpp(b byte) []byte {
	pixel1 := bit2value(b, 3, 0b10) | bit2value(b, 7, 0b01)
	pixel2 := bit2value(b, 2, 0b10) | bit2value(b, 6, 0b01)
	pixel3 := bit2value(b, 1, 0b10) | bit2value(b, 5, 0b01)
	pixel4 := bit2value(b, 0, 0b10) | bit2value(b, 4, 0b01)
	return []byte{pixel1, pixel2, pixel3, pixel4}
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
				draw.Draw(ga.displayScaled, ga.displayScaled.Bounds(), &image.Uniform{ga.borderColor}, image.ZP, draw.Src)
			}
		}
	} else if f == 2 {
		ga.mem.lowerRomEnable = data&0b00000100 == 0
		ga.mem.upperRomEnable = data&0b00001000 == 0

		if data&0x10 != 0 {
			ga.crtc.cpu.Interrupt(false)
			ga.crtc.counters.sl = 0
		}

		screenMode := data & 0b00000011
		if ga.screenMode != screenMode {
			ga.screenMode = screenMode
			// println("screenMode", screenMode)
			switch screenMode {
			case 0:
				ga.display.Rect = image.Rect(0, 0, 160, 200)
				ga.decode = to4bpp
				ga.ppc = 4
			case 1:
				ga.display.Rect = image.Rect(0, 0, 320, 200)
				ga.decode = to2bpp
				ga.ppc = 8
			case 2:
				ga.display.Rect = image.Rect(0, 0, 640, 200)
				ga.decode = to1bpp
				ga.ppc = 16
			default:
				// panic(screenMode)
			}
		}
		// println("[ga]", "lowerRomEnable:", ga.mem.lowerRomEnable, "upperRomEnable:", ga.mem.upperRomEnable, "screenMode:", ga.screenMode)
	} else if f == 3 {
		ga.mem.Paging(data)
	} else {
		panic(fmt.Sprintf("unsupported f:%d", f))
	}
}
