package zx

import (
	"image/color"
	"sync"

	"github.com/laullon/b2t80s/cpu/z80"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
	"github.com/laullon/b2t80s/gui"
)

var palette = []color.RGBA{
	{0x00, 0x00, 0x00, 0xff},
	{0x20, 0x30, 0xc0, 0xff},
	{0xc0, 0x40, 0x10, 0xff},
	{0xc0, 0x40, 0xc0, 0xff},
	{0x40, 0xb0, 0x10, 0xff},
	{0x50, 0xc0, 0xb0, 0xff},
	{0xe0, 0xc0, 0x10, 0xff},
	{0xc0, 0xc0, 0xc0, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x30, 0x40, 0xff, 0xff},
	{0xff, 0x40, 0x30, 0xff},
	{0xff, 0x70, 0xf0, 0xff},
	{0x50, 0xe0, 0x10, 0xff},
	{0x50, 0xe0, 0xff, 0xff},
	{0xff, 0xe8, 0x50, 0xff},
	{0xff, 0xff, 0xff, 0xff},
}

type ula struct {
	memory *memory
	bus    z80.Bus
	cpu    z80.Z80

	keyboardRow  []byte
	borderColour color.RGBA

	frame   byte
	display *gui.Display
	monitor emulator.Monitor

	col          int
	tsPerRow     int
	row          uint16
	scanlines    uint16
	displayStart uint16

	scanlinesBorder [][]color.RGBA
	pixlesData      [][]byte
	pixlesAttr      [][]byte
	floatingBus     byte

	cassette       cassette.Cassette
	ear, earActive bool
	buzzer         bool
	out            []*emulator.SoundData
	mux            sync.Mutex
}

func NewULA(mem *memory, bus z80.Bus, plus bool) *ula {
	ula := &ula{
		memory:          mem,
		bus:             bus,
		keyboardRow:     make([]byte, 8),
		borderColour:    palette[0],
		scanlinesBorder: make([][]color.RGBA, 313),
		pixlesData:      make([][]byte, 192),
		pixlesAttr:      make([][]byte, 192),
		display:         gui.NewDisplay(352, 296),
	}

	ula.monitor = emulator.NewMonitor(ula.display)

	if !plus {
		// 48k
		ula.tsPerRow = 224
		ula.scanlines = 312
		ula.displayStart = 64
	} else {
		// 128k
		ula.tsPerRow = 228
		ula.scanlines = 311
		ula.displayStart = 63
	}

	ula.keyboardRow[0] = 0x1f
	ula.keyboardRow[1] = 0x1f
	ula.keyboardRow[2] = 0x1f
	ula.keyboardRow[3] = 0x1f
	ula.keyboardRow[4] = 0x1f
	ula.keyboardRow[5] = 0x1f
	ula.keyboardRow[6] = 0x1f
	ula.keyboardRow[7] = 0x1f

	for y := 0; y < 192; y++ {
		ula.pixlesData[y] = make([]byte, 32)
		ula.pixlesAttr[y] = make([]byte, 32)
	}

	for y := uint16(0); y < ula.scanlines; y++ {
		ula.scanlinesBorder[y] = make([]color.RGBA, ula.tsPerRow)
	}

	return ula
}

func (ula *ula) Tick() {
	// EAR
	if ula.cassette != nil {
		ula.ear = ula.cassette.Ear()
	}

	// SCREEN
	draw := false
	io := false
	if ula.col < 128 && ula.row >= ula.displayStart && ula.row < ula.displayStart+192 {
		io = (ula.col % 8) < 6
		draw = ula.col%4 == 0
	} else {
		ula.floatingBus = 0xff
	}

	ula.scanlinesBorder[ula.row][ula.col] = ula.borderColour

	// CPU CLOCK
	if io {
		if ula.bus.GetAddr()>>14 != 1 {
			ula.cpu.Tick()
		}
	} else {
		ula.cpu.Tick()
	}

	if draw {
		y := uint16(ula.row - ula.displayStart)
		x := uint16(ula.col) / 4
		addr := uint16(0)
		addr |= ((y & 0b00000111) | 0b01000000) << 8
		addr |= ((y >> 3) & 0b00011000) << 8
		addr |= ((y << 2) & 0b11100000)
		ula.pixlesData[y][x] = ula.memory.Read(addr + x)

		attrAddr := uint16(((y >> 3) * 32) + 0x5800)
		ula.pixlesAttr[y][x] = ula.memory.Read(attrAddr + x)
		ula.floatingBus = ula.pixlesAttr[y][x]
	}

	ula.col++
	if ula.col == ula.tsPerRow {
		ula.row++
		if ula.row == ula.scanlines {
			ula.row = 0
			ula.cpu.Interrupt(true)
			ula.FrameDone()
		}
		ula.col = 0
	}
}

func (ula *ula) FrameDone() {
	ula.frame = (ula.frame + 1) & 0x1f
	for y := uint16(0); y < 296; y++ {
		for x := uint16(0); x < 352; x++ {
			ula.display.Set(x, y, ula.getPixel(x, y))
		}
	}
	ula.monitor.FrameDone()
}

func (ula *ula) ReadPort(port uint16) byte {
	if port&0xff == 0xfe {
		data := byte(0b00011111)
		readRow := port >> 8
		for row := 0; row < 8; row++ {
			if (readRow & (1 << row)) == 0 {
				data &= ula.keyboardRow[row]
			}
		}
		if ula.earActive && ula.ear {
			data |= 0b11100000
		} else {
			data |= 0b10100000
		}
		return data
	}
	return ula.floatingBus
}

func (ula *ula) WritePort(port uint16, data byte) {
	if port&0xff == 0xfe {
		if ula.borderColour != palette[data&0x07] {
			ula.borderColour = palette[data&0x07]
			// println("------", ula.col, ula.row)
		}
		ula.buzzer = ((data & 16) >> 4) != 0
		ula.earActive = (data & 24) != 0
		// println("ula.earActive:", ula.earActive, "ula.buzzer:", ula.buzzer)
	} else {
		// log.Printf("[write] port:0x%02x data:0b%08b", port, data)
	}
	// log.Printf("[write] port:0x%02x data:0b%08b", port, data)
	// ula.keyboardRow[port] = data
}

func (ula *ula) getPixel(rx, ry uint16) color.RGBA {
	border := false
	if ry < ula.displayStart || ry >= ula.displayStart+192 {
		border = true
	} else if rx < 48 || rx > 47+256 {
		border = true
	}

	if border {
		// if ry == ula.displayStart || ry == 80 {
		// 	return palette[0]
		// }
		return ula.scanlinesBorder[ry][rx/8]
	}

	ry -= ula.displayStart
	rx -= 48

	x := rx >> 3
	b := rx & 0x07

	attr := ula.pixlesAttr[ry][x]

	flash := (attr & 0x80) == 0x80
	brg := (attr & 0x40) >> 6
	paper := palette[((attr&0x38)>>3)+(brg*8)]
	ink := palette[(attr&0x07)+(brg*8)]

	data := ula.pixlesData[ry][x]
	data = data << b
	data &= 0b10000000
	if flash && (ula.frame&0x10 != 0) {
		if data != 0 {
			return paper
		}
		return ink
	}

	if data != 0 {
		return ink
	}
	return paper
}
