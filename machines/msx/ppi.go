package msx

import (
	"fmt"
	"sync"

	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/emulator/storage/cassette"
)

type ppi struct {
	mem      *memory
	cassette cassette.Cassette

	click    bool
	soundOut []*emulator.SoundData
	mux      sync.Mutex

	keyboardRows []byte

	c byte
}

func newPPI(mem *memory, cassette cassette.Cassette) *ppi {
	ppi := &ppi{
		mem:          mem,
		cassette:     cassette,
		keyboardRows: make([]byte, 0x10),
	}

	for idx := 0; idx < len(ppi.keyboardRows); idx++ {
		ppi.keyboardRows[idx] = 0xff
	}

	return ppi
}

func (ppi *ppi) ReadPort(port uint16) (byte, bool) {
	switch port & 0xff {
	case 0xa8:
		res := byte(0)
		for slot := 0; slot < 4; slot++ {
			res |= ppi.mem.cfg[slot] << (slot * 2)
		}
		// fmt.Printf("[ppi.ReadPort] mem.cfg: %v (0b%08b)\n", ppi.mem.cfg, res)
		return res, false

	case 0xa9:
		res := ppi.keyboardRows[ppi.c&0x0f]
		return res, false

	case 0xaa:
		return ppi.c, false
	}
	panic(fmt.Sprintf("[ReadPort] Unsopported port: 0x%02X", port))
}

func (ppi *ppi) WritePort(port uint16, data byte) {
	switch port & 0xff {
	case 0xa8:
		for slot := 0; slot < 4; slot++ {
			ppi.mem.cfg[slot] = (data >> (slot * 2)) & 3
		}
		// fmt.Printf("mem.cfg: %v (0b%08b)\n", ppi.mem.cfg, data)

	case 0xa9:
		panic(fmt.Sprintf("unsopported port: 0x%02X", port))

	case 0xaa:
		ppi.c = data

	case 0xab:
		if (data & 0x80) == 0 {
			bit := (data >> 1) & 7
			if (data & 1) == 1 {
				ppi.c |= 1 << bit
			} else {
				ppi.c &= ^(1 << bit)
			}
		} else {
			print(fmt.Sprintf("invalid c (%d) value\n", data))
		}
	}
}

func (ppi *ppi) SoundTick() {
	ppi.mux.Lock()
	defer ppi.mux.Unlock()
	v := 0.0
	if ppi.click {
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
