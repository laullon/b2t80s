package emulator

import (
	"image"
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne"
)

var palette = []color.RGBA{
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x20, 0x30, 0xc0, 0xff},
	color.RGBA{0xc0, 0x40, 0x10, 0xff},
	color.RGBA{0xc0, 0x40, 0xc0, 0xff},
	color.RGBA{0x40, 0xb0, 0x10, 0xff},
	color.RGBA{0x50, 0xc0, 0xb0, 0xff},
	color.RGBA{0xe0, 0xc0, 0x10, 0xff},
	color.RGBA{0xc0, 0xc0, 0xc0, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x30, 0x40, 0xff, 0xff},
	color.RGBA{0xff, 0x40, 0x30, 0xff},
	color.RGBA{0xff, 0x70, 0xf0, 0xff},
	color.RGBA{0x50, 0xe0, 0x10, 0xff},
	color.RGBA{0x50, 0xe0, 0xff, 0xff},
	color.RGBA{0xff, 0xe8, 0x50, 0xff},
	color.RGBA{0xff, 0xff, 0xff, 0xff},
}

type ULA interface {
	PortManager
	Ticker
	SoundSource

	OnKeyEvent(event *fyne.KeyEvent)

	LoadCommand() uint16

	FrameDone()
	Display() image.Image
}

type ula struct {
	tStates uint

	memory Memory

	keyboardRow  []byte
	borderColour color.RGBA

	frame   byte
	display *image.RGBA

	scanline        uint
	scanlinesBorder []color.RGBA
	scanlinesData   [][]byte
	scanlinesAttr   [][]byte

	cassette       Cassette
	ear, earActive bool
	buzzer         bool
	out            []*SoundData
	mux            sync.Mutex

	// tStatesPerSample uint
}

func NewULA(mem Memory, cassette Cassette) ULA {
	ula := &ula{
		memory:          mem,
		keyboardRow:     make([]byte, 8),
		borderColour:    palette[0],
		scanlinesBorder: make([]color.RGBA, 296),
		scanlinesData:   make([][]byte, 192),
		scanlinesAttr:   make([][]byte, 192),
		display:         image.NewRGBA(image.Rect(0, 0, 352, 296)),
		// player:           SoundContext.NewPlayer(),
		cassette: cassette,
		// tStatesPerSample: uint(math.Round(float64(3500000) / float64(freq))),
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
		ula.scanlinesData[y] = make([]byte, 32)
		ula.scanlinesAttr[y] = make([]byte, 32)
	}

	return ula
}

func (ula *ula) Tick() {
	ula.tStates++

	// EAR
	if ula.cassette != nil {
		ula.ear = ula.cassette.Ear()
	}

	// SCREEN
	scanline := (ula.tStates / 224)
	if ula.scanline == scanline {
		return
	}

	if ula.scanline < 296 {
		ula.scanlinesBorder[ula.scanline] = ula.borderColour
	}

	if ula.scanline > 47 && ula.scanline < 240 {
		y := uint16(ula.scanline - 48)
		addr := uint16(0)
		addr |= ((y & 0b00000111) | 0b01000000) << 8
		addr |= ((y >> 3) & 0b00011000) << 8
		addr |= ((y << 2) & 0b11100000)
		ula.scanlinesData[y] = ula.memory.GetBlock(addr, 32)

		attrAddr := uint16(((y >> 3) * 32) + 0x5800)
		ula.scanlinesAttr[y] = ula.memory.GetBlock(attrAddr, 32)
	}

	ula.scanline = scanline
	return
}

var s = 0

func (ula *ula) FrameDone() {
	ula.tStates = 0

	ula.frame = (ula.frame + 1) & 0x1f
	for y := 0; y < 296; y++ {
		for x := 0; x < 352; x++ {
			ula.display.Set(x, y, ula.getPixle(x, y))
		}
	}
}

func (ula *ula) SoundTick() {
	ula.mux.Lock()
	defer ula.mux.Unlock()
	v := 0.0
	if ula.buzzer || ula.ear {
		v = 1
	}
	ula.out = append(ula.out, &SoundData{L: v, R: v})
}

func (ula *ula) GetBuffer(max int) (res []*SoundData, l int) {
	ula.mux.Lock()
	defer ula.mux.Unlock()

	if len(ula.out) > max {
		res = ula.out[:max]
		ula.out = ula.out[max:]
		l = max
	} else {
		res = ula.out
		ula.out = nil
		l = len(res)
	}
	return
}

func (ula *ula) ReadPort(port uint16) (byte, bool) {
	data := byte(0b00011111)
	if port&0xfe == 0xfe {
		readRow := port >> 8
		for row := 0; row < 8; row++ {
			if (readRow & (1 << row)) == 0 {
				data &= ula.keyboardRow[row]
				// log.Printf("[read] data:0b%08b row:0b%08b (%d)", data, ula.keyboardRow[row], row)
			}
		}
		if ula.earActive && ula.ear {
			data |= 0b11100000
		} else {
			data |= 0b10100000
		}
		// } else {
		// 	log.Printf("[read] port:0x%02x data:0b%08b", port, data)
	}
	// log.Printf("[read] port:0x%02x data:0b%08b", port, data)
	return data, false
}

func (ula *ula) WritePort(port uint16, data byte) {
	if port&0xff == 0xfe {
		ula.borderColour = palette[data&0x07]
		ula.buzzer = ((data & 16) >> 4) != 0
		ula.earActive = (data & 24) != 0
		// println("ula.earActive:", ula.earActive, "ula.buzzer:", ula.buzzer)
	} else {
		// log.Printf("[write] port:0x%02x data:0b%08b", port, data)
	}
	// log.Printf("[write] port:0x%02x data:0b%08b", port, data)
	// ula.keyboardRow[port] = data
}

func (ula *ula) OnKeyEvent(key *fyne.KeyEvent) {
	// fmt.Println("key:", key.Name)
	switch key.Name {

	case fyne.Key1:
		ula.keyboardRow[3] ^= 0b00000001
	case fyne.Key2:
		ula.keyboardRow[3] ^= 0b00000010
	case fyne.Key3:
		ula.keyboardRow[3] ^= 0b00000100
	case fyne.Key4:
		ula.keyboardRow[3] ^= 0b00001000
	case fyne.Key5:
		ula.keyboardRow[3] ^= 0b00010000

	case fyne.Key0:
		ula.keyboardRow[4] ^= 0b00000001
	case fyne.Key9:
		ula.keyboardRow[4] ^= 0b00000010
	case fyne.Key8:
		ula.keyboardRow[4] ^= 0b00000100
	case fyne.Key7:
		ula.keyboardRow[4] ^= 0b00001000
	case fyne.Key6:
		ula.keyboardRow[4] ^= 0b00010000

	case fyne.KeyQ:
		ula.keyboardRow[2] ^= 0b00000001
	case fyne.KeyW:
		ula.keyboardRow[2] ^= 0b00000010
	case fyne.KeyE:
		ula.keyboardRow[2] ^= 0b00000100
	case fyne.KeyR:
		ula.keyboardRow[2] ^= 0b00001000
	case fyne.KeyT:
		ula.keyboardRow[2] ^= 0b00010000

	case fyne.KeyP:
		ula.keyboardRow[5] ^= 0b00000001
	case fyne.KeyO:
		ula.keyboardRow[5] ^= 0b00000010
	case fyne.KeyI:
		ula.keyboardRow[5] ^= 0b00000100
	case fyne.KeyU:
		ula.keyboardRow[5] ^= 0b00001000
	case fyne.KeyY:
		ula.keyboardRow[5] ^= 0b00010000

	case fyne.KeyA:
		ula.keyboardRow[1] ^= 0b00000001
	case fyne.KeyS:
		ula.keyboardRow[1] ^= 0b00000010
	case fyne.KeyD:
		ula.keyboardRow[1] ^= 0b00000100
	case fyne.KeyF:
		ula.keyboardRow[1] ^= 0b00001000
	case fyne.KeyG:
		ula.keyboardRow[1] ^= 0b00010000

	case fyne.KeyReturn:
		ula.keyboardRow[6] ^= 0b00000001
	case fyne.KeyL:
		ula.keyboardRow[6] ^= 0b00000010
	case fyne.KeyK:
		ula.keyboardRow[6] ^= 0b00000100
	case fyne.KeyJ:
		ula.keyboardRow[6] ^= 0b00001000
	case fyne.KeyH:
		ula.keyboardRow[6] ^= 0b00010000

	case "LeftShift":
		ula.keyboardRow[0] ^= 0b00000001
	case fyne.KeyZ:
		ula.keyboardRow[0] ^= 0b00000010
	case fyne.KeyX:
		ula.keyboardRow[0] ^= 0b00000100
	case fyne.KeyC:
		ula.keyboardRow[0] ^= 0b00001000
	case fyne.KeyV:
		ula.keyboardRow[0] ^= 0b00010000

	case "Space":
		ula.keyboardRow[7] ^= 0b00000001
	case "RightSuper":
		ula.keyboardRow[7] ^= 0b00000010
	case fyne.KeyM:
		ula.keyboardRow[7] ^= 0b00000100
	case fyne.KeyN:
		ula.keyboardRow[7] ^= 0b00001000
	case fyne.KeyB:
		ula.keyboardRow[7] ^= 0b00010000

	case "BackSpace":
		ula.keyboardRow[0] ^= 0b00000001
		ula.keyboardRow[4] ^= 0b00000001

	case fyne.KeyUp:
		ula.keyboardRow[0] ^= 0b00000001
		ula.keyboardRow[4] ^= 0b00001000

	case fyne.KeyDown:
		ula.keyboardRow[0] ^= 0b00000001
		ula.keyboardRow[4] ^= 0b00010000
	}
}

var onlyOnce sync.Once

func (ula *ula) LoadCommand() uint16 {
	go onlyOnce.Do(func() {
		time.Sleep(time.Second)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyJ})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyJ})
		ula.OnKeyEvent(&fyne.KeyEvent{Name: "RightSuper"})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyP})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyP})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyP})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyP})
		ula.OnKeyEvent(&fyne.KeyEvent{Name: "RightSuper"})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyReturn})
		time.Sleep(150 * time.Millisecond)
		ula.OnKeyEvent(&fyne.KeyEvent{Name: fyne.KeyReturn})
	})
	return 0
}

func (ula *ula) Display() image.Image {
	return ula.display
}

func (ula *ula) getPixle(rx, ry int) color.Color {
	border := false
	if ry < 48 || ry > 47+192 {
		border = true
	} else if rx < 48 || rx > 47+256 {
		border = true
	}

	if border {
		return ula.scanlinesBorder[ry]
	}

	rx -= 48
	ry -= 48

	x := rx >> 3
	b := rx & 0x07

	attr := ula.scanlinesAttr[ry][x]

	flash := (attr & 0x80) == 0x80
	brg := (attr & 0x40) >> 6
	paper := palette[((attr&0x38)>>3)+(brg*8)]
	ink := palette[(attr&0x07)+(brg*8)]

	data := ula.scanlinesData[ry][x]
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
