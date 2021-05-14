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

func (ppi *ppi) ReadPort(port uint16) (byte, bool) {
	t := (port>>8)&0xf - 4

	if t == 0x0 {
		return ppi.readA(), false
	} else if t == 0x1 {
		return ppi.readB(), false
	} else if t == 0x2 {
		return ppi.readC(), false
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
					if (ppi.keyboardLine & 0x0f) == 9 {
						j, _ := emulator.ReadJoystick()
						res = byte(0xff)
						// 000FRLDU
						if j.ON {
							if j.F {
								res ^= 0b00010000
							}
							if j.R {
								res ^= 0b00001000
							}
							if j.L {
								res ^= 0b00000100
							}
							if j.D {
								res ^= 0b00000010
							}
							if j.U {
								res ^= 0b00000001
							}
						}
					} else {
						res = ppi.keyboardRows[ppi.keyboardLine&0x0f]
					}
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
	// switch key {

	// case glfw.KeyUp:
	// 	ppi.keyboardRows[0] ^= 0b00000001
	// case glfw.KeyRight:
	// 	ppi.keyboardRows[0] ^= 0b00000010
	// case glfw.KeyDown:
	// 	ppi.keyboardRows[0] ^= 0b00000100
	// case glfw.KeyF9:
	// 	ppi.keyboardRows[0] ^= 0b00001000
	// case glfw.KeyF6:
	// 	ppi.keyboardRows[0] ^= 0b00010000
	// case glfw.KeyF3:
	// 	ppi.keyboardRows[0] ^= 0b00100000
	// case glfw.KeyEnter:
	// 	ppi.keyboardRows[0] ^= 0b01000000
	// case glfw.KeyF11:
	// 	ppi.keyboardRows[0] ^= 0b10000000

	// case glfw.KeyLeft:
	// 	ppi.keyboardRows[1] ^= 0b00000001
	// // case glfw.Key: COPY
	// // 	ppi.keyboardRows[1] ^= 0b00000010
	// case glfw.KeyF7:
	// 	ppi.keyboardRows[1] ^= 0b00000100
	// case glfw.KeyF8:
	// 	ppi.keyboardRows[1] ^= 0b00001000
	// case glfw.KeyF5:
	// 	ppi.keyboardRows[1] ^= 0b00010000
	// case glfw.KeyF1:
	// 	ppi.keyboardRows[1] ^= 0b00100000
	// case glfw.KeyF2:
	// 	ppi.keyboardRows[1] ^= 0b01000000
	// case glfw.KeyF10:
	// 	ppi.keyboardRows[1] ^= 0b10000000

	// // case glfw.KeyUp: CLR
	// // 	ppi.keyboardRows[2] ^= 0b00000001
	// case glfw.KeyLeftBracket:
	// 	ppi.keyboardRows[2] ^= 0b00000010
	// case glfw.KeyKPEnter:
	// 	ppi.keyboardRows[2] ^= 0b00000100
	// case glfw.KeyRightBracket:
	// 	ppi.keyboardRows[2] ^= 0b00001000
	// case glfw.KeyF4:
	// 	ppi.keyboardRows[2] ^= 0b00010000
	// case glfw.KeyLeftShift, glfw.KeyRightShift:
	// 	ppi.keyboardRows[2] ^= 0b00100000
	// case glfw.KeyBackslash:
	// 	ppi.keyboardRows[2] ^= 0b01000000
	// case glfw.KeyLeftControl, glfw.KeyRightControl:
	// 	ppi.keyboardRows[2] ^= 0b10000000

	// // case glfw.KeyUp: ^
	// // 	ppi.keyboardRows[3] ^= 0b00000001
	// case glfw.KeyMinus:
	// 	ppi.keyboardRows[3] ^= 0b00000010
	// case glfw.KeyApostrophe:
	// 	ppi.keyboardRows[3] ^= 0b00000100
	// case glfw.KeyP:
	// 	ppi.keyboardRows[3] ^= 0b00001000
	// case glfw.KeySemicolon:
	// 	ppi.keyboardRows[3] ^= 0b00010000
	// // case fyne.Ke: :
	// // 	ppi.keyboardRows[3] ^= 0b00100000
	// case glfw.KeySlash:
	// 	ppi.keyboardRows[3] ^= 0b01000000
	// case glfw.KeyPeriod:
	// 	ppi.keyboardRows[3] ^= 0b10000000

	// case glfw.Key0:
	// 	ppi.keyboardRows[4] ^= 0b00000001
	// case glfw.Key9:
	// 	ppi.keyboardRows[4] ^= 0b00000010
	// case glfw.KeyO:
	// 	ppi.keyboardRows[4] ^= 0b00000100
	// case glfw.KeyI:
	// 	ppi.keyboardRows[4] ^= 0b00001000
	// case glfw.KeyL:
	// 	ppi.keyboardRows[4] ^= 0b00010000
	// case glfw.KeyK:
	// 	ppi.keyboardRows[4] ^= 0b00100000
	// case glfw.KeyM:
	// 	ppi.keyboardRows[4] ^= 0b01000000
	// case glfw.KeyComma:
	// 	ppi.keyboardRows[4] ^= 0b10000000

	// case glfw.Key8:
	// 	ppi.keyboardRows[5] ^= 0b00000001
	// case glfw.Key7:
	// 	ppi.keyboardRows[5] ^= 0b00000010
	// case glfw.KeyU:
	// 	ppi.keyboardRows[5] ^= 0b00000100
	// case glfw.KeyY:
	// 	ppi.keyboardRows[5] ^= 0b00001000
	// case glfw.KeyH:
	// 	ppi.keyboardRows[5] ^= 0b00010000
	// case glfw.KeyJ:
	// 	ppi.keyboardRows[5] ^= 0b00100000
	// case glfw.KeyN:
	// 	ppi.keyboardRows[5] ^= 0b01000000
	// case glfw.KeySpace:
	// 	ppi.keyboardRows[5] ^= 0b10000000

	// case glfw.Key6:
	// 	ppi.keyboardRows[6] ^= 0b00000001
	// case glfw.Key5:
	// 	ppi.keyboardRows[6] ^= 0b00000010
	// case glfw.KeyR:
	// 	ppi.keyboardRows[6] ^= 0b00000100
	// case glfw.KeyT:
	// 	ppi.keyboardRows[6] ^= 0b00001000
	// case glfw.KeyG:
	// 	ppi.keyboardRows[6] ^= 0b00010000
	// case glfw.KeyF:
	// 	ppi.keyboardRows[6] ^= 0b00100000
	// case glfw.KeyB:
	// 	ppi.keyboardRows[6] ^= 0b01000000
	// case glfw.KeyV:
	// 	ppi.keyboardRows[6] ^= 0b10000000

	// case glfw.Key4:
	// 	ppi.keyboardRows[7] ^= 0b00000001
	// case glfw.Key3:
	// 	ppi.keyboardRows[7] ^= 0b00000010
	// case glfw.KeyE:
	// 	ppi.keyboardRows[7] ^= 0b00000100
	// case glfw.KeyW:
	// 	ppi.keyboardRows[7] ^= 0b00001000
	// case glfw.KeyS:
	// 	ppi.keyboardRows[7] ^= 0b00010000
	// case glfw.KeyD:
	// 	ppi.keyboardRows[7] ^= 0b00100000
	// case glfw.KeyC:
	// 	ppi.keyboardRows[7] ^= 0b01000000
	// case glfw.KeyX:
	// 	ppi.keyboardRows[7] ^= 0b10000000

	// case glfw.Key1:
	// 	ppi.keyboardRows[8] ^= 0b00000001
	// case glfw.Key2:
	// 	ppi.keyboardRows[8] ^= 0b00000010
	// case glfw.KeyEscape:
	// 	ppi.keyboardRows[8] ^= 0b00000100
	// case glfw.KeyQ:
	// 	ppi.keyboardRows[8] ^= 0b00001000
	// case glfw.KeyTab:
	// 	ppi.keyboardRows[8] ^= 0b00010000
	// case glfw.KeyA:
	// 	ppi.keyboardRows[8] ^= 0b00100000
	// // case glfw.Key: CAPSLOCK
	// // ppi.keyboardRows[8] ^= 0b01000000
	// case glfw.KeyZ:
	// 	ppi.keyboardRows[8] ^= 0b10000000

	// // case glfw.KeyUp: JOY1
	// // 	ppi.keyboardRows[9] ^= 0b00000001
	// // case glfw.KeyRight:
	// // 	ppi.keyboardRows[9] ^= 0b00000010
	// // case glfw.KeyDown:
	// // 	ppi.keyboardRows[9] ^= 0b00000100
	// // case glfw.KeyF9:
	// // 	ppi.keyboardRows[9] ^= 0b00001000
	// // case glfw.KeyF6:
	// // 	ppi.keyboardRows[9] ^= 0b00010000
	// // case glfw.KeyF3:
	// // 	ppi.keyboardRows[9] ^= 0b00100000
	// // case glfw.KeyEnter:
	// // 	ppi.keyboardRows[9] ^= 0b01000000
	// case glfw.KeyDelete, glfw.KeyBackspace:
	// 	ppi.keyboardRows[9] ^= 0b10000000

	// 	// default:
	// 	// fmt.Println("key:", key.Name)
	// }
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
