package cpc

import (
	"sync"

	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/ay8912"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	psgINACTIVE byte = 0
	psgREAD     byte = 1
	psgWRITE    byte = 2
	psgSELECT   byte = 3
)

type ppi struct {
	crtc     *crtc
	cassette cassette.Cassette

	keyboardRows []byte
	keyboardLine byte

	psg        ay8912.AY8912
	psgReg     byte
	psgControl byte
	psgReg15   byte

	a, b, c byte

	aInput, bInput, clInput, chInput bool

	soundOut []*emulator.SoundData
	mux      sync.Mutex
}

func newPPI(crtc *crtc, cassette cassette.Cassette, psg ay8912.AY8912) *ppi {
	ppi := &ppi{
		psg:          psg,
		crtc:         crtc,
		cassette:     cassette,
		keyboardRows: make([]byte, 0x10),
	}

	for idx := 0; idx < len(ppi.keyboardRows); idx++ {
		ppi.keyboardRows[idx] = 0xff
	}

	return ppi
}

func (ppi *ppi) ReadPort(port uint16) byte {
	t := (port>>8)&0xf - 4

	if t == 0x0 {
		return ppi.readA()
	} else if t == 0x1 {
		return ppi.readB()
	} else if t == 0x2 {
		return ppi.readC()
	} else {
		panic(t)
	}

}

func (ppi *ppi) WritePort(port uint16, data byte) {
	t := (port>>8)&0xf - 4

	if t == 0x0 {
		ppi.writeA(data)
	} else if t == 0x1 {
		ppi.writeB(data)
	} else if t == 0x2 {
		ppi.writeC(data)
	} else if t == 0x3 {
		ppi.writeControl(data)
	} else {
		panic(t)
	}
}

func (ppi *ppi) writeControl(value byte) {
	if (value & 0x80) != 0 { // change PPI configuration
		ppi.aInput = value&0b00010000 != 0
		ppi.bInput = value&0b00000010 != 0
		ppi.clInput = value&0b00000001 != 0
		ppi.chInput = value&0b00001000 != 0
		ppi.a = 0 // clear data for all ports
		ppi.b = 0
		ppi.c = 0
	} else { // bit manipulation of port C data
		bit := (value >> 1) & 7 // isolate bit to set
		if (value & 1) == 1 {
			ppi.c |= 1 << bit // set requested bit
		} else {
			ppi.c &= ^(1 << bit) // reset requested bit
		}

		if !ppi.clInput { // output lower half?
			ppi.keyboardLine = ppi.c
		}

		if !ppi.chInput { // output upper half?
			ppi.cassette.Motor(value&0x10 != 0)
			ppi.psgControl = (value & 0xc0) >> 6
			ppi.psgWrite(ppi.a)
		}
	}
}

func (ppi *ppi) writeA(value byte) {
	ppi.a = value
	if !ppi.aInput {
		ppi.psgWrite(value)
	}
}

func (ppi *ppi) writeB(value byte) {
	ppi.b = value
}

func (ppi *ppi) writeC(value byte) {
	ppi.c = value
	if !ppi.clInput { // output lower half?
		ppi.keyboardLine = value
	}
	if !ppi.chInput { // output upper half?
		ppi.cassette.Motor(value&0x10 != 0)
		ppi.psgControl = (value & 0xc0) >> 6
		ppi.psgWrite(ppi.a)
	}
}

func (ppi *ppi) psgWrite(data byte) {
	if ppi.psgControl == psgSELECT {
		ppi.psgReg = data
	} else if ppi.psgControl == psgWRITE {
		if ppi.psgReg < 16 {
			if ppi.psgReg == 14 {
				ppi.keyboardLine = data
			} else {
				ppi.psg.WriteRegister(ppi.psgReg, data)
			}
		} else {
			panic(ppi.psgReg)
		}
	}
}

func (ppi *ppi) readA() (res byte) {
	if ppi.aInput {
		if ppi.psgControl == psgREAD {
			if ppi.psgReg < 16 {
				if ppi.psgReg == 14 {
					res = ppi.keyboardRows[ppi.keyboardLine&0x0f]
				} else {
					res = ppi.psg.ReadRegister(ppi.psgReg)
				}
			} else {
				panic(ppi.psgReg)
			}
		}
	} else {
		res = ppi.a
	}
	return res
}

func (ppi *ppi) readB() byte {
	if ppi.bInput {
		res := byte(0b00011110)
		if ppi.crtc.status.vSync {
			res |= 1
		}
		if ppi.cassette.Ear() {
			res |= 0b10000000
		}
		return res
	}
	return ppi.b
}

func (ppi *ppi) readC() (res byte) {
	res = ppi.c

	if ppi.clInput || ppi.chInput { // either half set to input?
		if ppi.chInput {
			res &= 0x0f         // blank out upper half
			val := ppi.c & 0xc0 // isolate PSG control bits
			if val == 0xc0 {    // PSG specify register?
				val = 0x80 // change to PSG write register
			}
			res |= val | 0x20 // casette write data is always set
			if ppi.cassette.IsMotorON() {
				res |= 0x10 // set the bit if the tape motor is running
			}
		}
		if !ppi.clInput { // lower half set to output?
			res |= 0x0f // invalid - set all bits
		}
	}
	return res
}

func (ppi *ppi) OnKey(key sdl.Scancode) {
	switch key {

	case sdl.SCANCODE_UP:
		ppi.keyboardRows[0] ^= 0b00000001
	case sdl.SCANCODE_RIGHT:
		ppi.keyboardRows[0] ^= 0b00000010
	case sdl.SCANCODE_DOWN:
		ppi.keyboardRows[0] ^= 0b00000100
	case sdl.SCANCODE_F9:
		ppi.keyboardRows[0] ^= 0b00001000
	case sdl.SCANCODE_F6:
		ppi.keyboardRows[0] ^= 0b00010000
	case sdl.SCANCODE_F3:
		ppi.keyboardRows[0] ^= 0b00100000
	case sdl.SCANCODE_RETURN:
		ppi.keyboardRows[0] ^= 0b01000000
	case sdl.SCANCODE_F11:
		ppi.keyboardRows[0] ^= 0b10000000

	case sdl.SCANCODE_LEFT:
		ppi.keyboardRows[1] ^= 0b00000001
	// case sdl.SCANCODE_\U: COPY
	// 	ppi.keyboardRows[1] ^= 0b00000010
	case sdl.SCANCODE_F7:
		ppi.keyboardRows[1] ^= 0b00000100
	case sdl.SCANCODE_F8:
		ppi.keyboardRows[1] ^= 0b00001000
	case sdl.SCANCODE_F5:
		ppi.keyboardRows[1] ^= 0b00010000
	case sdl.SCANCODE_F1:
		ppi.keyboardRows[1] ^= 0b00100000
	case sdl.SCANCODE_F2:
		ppi.keyboardRows[1] ^= 0b01000000
	case sdl.SCANCODE_F10:
		ppi.keyboardRows[1] ^= 0b10000000

	// case sdl.SCANCODE_UP: CLR
	// 	ppi.keyboardRows[2] ^= 0b00000001
	case sdl.SCANCODE_LEFTBRACKET:
		ppi.keyboardRows[2] ^= 0b00000010
	case sdl.SCANCODE_KP_ENTER:
		ppi.keyboardRows[2] ^= 0b00000100
	case sdl.SCANCODE_RIGHTBRACKET:
		ppi.keyboardRows[2] ^= 0b00001000
	case sdl.SCANCODE_F4:
		ppi.keyboardRows[2] ^= 0b00010000
	case sdl.SCANCODE_LSHIFT, sdl.SCANCODE_RSHIFT:
		ppi.keyboardRows[2] ^= 0b00100000
	case sdl.SCANCODE_BACKSLASH:
		ppi.keyboardRows[2] ^= 0b01000000
	case sdl.SCANCODE_LCTRL, sdl.SCANCODE_RCTRL:
		ppi.keyboardRows[2] ^= 0b10000000

	// case sdl.SCANCODE_UP: ^
	// 	ppi.keyboardRows[3] ^= 0b00000001
	case sdl.SCANCODE_MINUS:
		ppi.keyboardRows[3] ^= 0b00000010
	case sdl.SCANCODE_APOSTROPHE:
		ppi.keyboardRows[3] ^= 0b00000100
	case sdl.SCANCODE_P:
		ppi.keyboardRows[3] ^= 0b00001000
	case sdl.SCANCODE_SEMICOLON:
		ppi.keyboardRows[3] ^= 0b00010000
	// case fyne.Ke: :
	// 	ppi.keyboardRows[3] ^= 0b00100000
	case sdl.SCANCODE_SLASH:
		ppi.keyboardRows[3] ^= 0b01000000
	case sdl.SCANCODE_PERIOD:
		ppi.keyboardRows[3] ^= 0b10000000

	case sdl.SCANCODE_0:
		ppi.keyboardRows[4] ^= 0b00000001
	case sdl.SCANCODE_9:
		ppi.keyboardRows[4] ^= 0b00000010
	case sdl.SCANCODE_O:
		ppi.keyboardRows[4] ^= 0b00000100
	case sdl.SCANCODE_I:
		ppi.keyboardRows[4] ^= 0b00001000
	case sdl.SCANCODE_L:
		ppi.keyboardRows[4] ^= 0b00010000
	case sdl.SCANCODE_K:
		ppi.keyboardRows[4] ^= 0b00100000
	case sdl.SCANCODE_M:
		ppi.keyboardRows[4] ^= 0b01000000
	case sdl.SCANCODE_COMMA:
		ppi.keyboardRows[4] ^= 0b10000000

	case sdl.SCANCODE_8:
		ppi.keyboardRows[5] ^= 0b00000001
	case sdl.SCANCODE_7:
		ppi.keyboardRows[5] ^= 0b00000010
	case sdl.SCANCODE_U:
		ppi.keyboardRows[5] ^= 0b00000100
	case sdl.SCANCODE_Y:
		ppi.keyboardRows[5] ^= 0b00001000
	case sdl.SCANCODE_H:
		ppi.keyboardRows[5] ^= 0b00010000
	case sdl.SCANCODE_J:
		ppi.keyboardRows[5] ^= 0b00100000
	case sdl.SCANCODE_N:
		ppi.keyboardRows[5] ^= 0b01000000
	case sdl.SCANCODE_SPACE:
		ppi.keyboardRows[5] ^= 0b10000000

	case sdl.SCANCODE_6:
		ppi.keyboardRows[6] ^= 0b00000001
	case sdl.SCANCODE_5:
		ppi.keyboardRows[6] ^= 0b00000010
	case sdl.SCANCODE_R:
		ppi.keyboardRows[6] ^= 0b00000100
	case sdl.SCANCODE_T:
		ppi.keyboardRows[6] ^= 0b00001000
	case sdl.SCANCODE_G:
		ppi.keyboardRows[6] ^= 0b00010000
	case sdl.SCANCODE_F:
		ppi.keyboardRows[6] ^= 0b00100000
	case sdl.SCANCODE_B:
		ppi.keyboardRows[6] ^= 0b01000000
	case sdl.SCANCODE_V:
		ppi.keyboardRows[6] ^= 0b10000000

	case sdl.SCANCODE_4:
		ppi.keyboardRows[7] ^= 0b00000001
	case sdl.SCANCODE_3:
		ppi.keyboardRows[7] ^= 0b00000010
	case sdl.SCANCODE_E:
		ppi.keyboardRows[7] ^= 0b00000100
	case sdl.SCANCODE_W:
		ppi.keyboardRows[7] ^= 0b00001000
	case sdl.SCANCODE_S:
		ppi.keyboardRows[7] ^= 0b00010000
	case sdl.SCANCODE_D:
		ppi.keyboardRows[7] ^= 0b00100000
	case sdl.SCANCODE_C:
		ppi.keyboardRows[7] ^= 0b01000000
	case sdl.SCANCODE_X:
		ppi.keyboardRows[7] ^= 0b10000000

	case sdl.SCANCODE_1:
		ppi.keyboardRows[8] ^= 0b00000001
	case sdl.SCANCODE_2:
		ppi.keyboardRows[8] ^= 0b00000010
	case sdl.SCANCODE_ESCAPE:
		ppi.keyboardRows[8] ^= 0b00000100
	case sdl.SCANCODE_Q:
		ppi.keyboardRows[8] ^= 0b00001000
	case sdl.SCANCODE_TAB:
		ppi.keyboardRows[8] ^= 0b00010000
	case sdl.SCANCODE_A:
		ppi.keyboardRows[8] ^= 0b00100000
	// case sdl.SCANCODE_\U: CAPSLOCK
	// ppi.keyboardRows[8] ^= 0b01000000
	case sdl.SCANCODE_Z:
		ppi.keyboardRows[8] ^= 0b10000000

	// case sdl.SCANCODE_UP: JOY1
	// 	ppi.keyboardRows[9] ^= 0b00000001
	// case sdl.SCANCODE_RIGHT:
	// 	ppi.keyboardRows[9] ^= 0b00000010
	// case sdl.SCANCODE_DOWN:
	// 	ppi.keyboardRows[9] ^= 0b00000100
	// case sdl.SCANCODE_F9:
	// 	ppi.keyboardRows[9] ^= 0b00001000
	// case sdl.SCANCODE_F6:
	// 	ppi.keyboardRows[9] ^= 0b00010000
	// case sdl.SCANCODE_F3:
	// 	ppi.keyboardRows[9] ^= 0b00100000
	// case sdl.SCANCODE_ENTER:
	// 	ppi.keyboardRows[9] ^= 0b01000000
	case sdl.SCANCODE_DELETE, sdl.SCANCODE_BACKSPACE:
		ppi.keyboardRows[9] ^= 0b10000000
	}
}

func (ppi *ppi) SoundTick() {
	ppi.mux.Lock()
	defer ppi.mux.Unlock()
	v := 0.0
	if ppi.cassette.Ear() {
		v = 1
	}
	ppi.soundOut = append(ppi.soundOut, &emulator.SoundData{L: v, R: v})
}

func (ppi *ppi) GetBuffer(max int) (res []*emulator.SoundData, l int) {
	ppi.mux.Lock()
	defer ppi.mux.Unlock()

	if len(ppi.soundOut) > max {
		res = ppi.soundOut[:max]
		ppi.soundOut = ppi.soundOut[max:]
		l = max
	} else {
		res = ppi.soundOut
		ppi.soundOut = nil
		l = len(res)
	}
	return
}

func (ppi *ppi) GetChannels() []emulator.SoundChannel {
	return []emulator.SoundChannel{ppi}
}
