package cassette

import (
	"fmt"
	"os"
)

type tap struct {
	blocks      []interface{}
	actualBlock int
}

type loopBlock struct {
	id     byte
	count  int
	blocks []interface{}
}

type loopEndBlock struct {
	id byte
}

type pulseSeqBlock struct {
	id     byte
	pulses []uint
}

type dataBlock struct {
	id   byte
	flag byte
	data []byte

	pilot, pilotLen, sync1, sync2, zero, one uint
	pause                                    uint
	lastBiteLen                              int8
}

func (d *dataBlock) getID() byte {
	return d.id
}

func (tap *tap) load(path string) {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	file := make([]byte, fi.Size()+1)
	l, err := f.Read(file)
	if err != nil {
		panic(err)
	}
	file = file[:l]

	header := string(file[:7])
	if header == "ZXTape!" {
		var loop *loopBlock
		var inLoop bool

		file = file[10:]
		for len(file) > 0 {
			block, l := readTzxBlock(file)
			file = file[l+1:]

			if b, ok := block.(*loopBlock); ok {
				inLoop = true
				loop = b
			} else if _, ok = block.(*loopEndBlock); ok {
				for i := 0; i < loop.count; i++ {
					tap.blocks = append(tap.blocks, loop.blocks...)
				}
				inLoop = false
			} else if block != nil {
				if inLoop {
					loop.blocks = append(loop.blocks, block)
				} else {
					tap.blocks = append(tap.blocks, block)
				}
			}
		}
	} else {
		for len(file) > 0 {
			block, l := readDefaultBlock(file)
			file = file[l:]
			tap.blocks = append(tap.blocks, block)
		}
	}
}

func readDefaultBlock(file []byte) (interface{}, uint16) {
	length := uint16(file[0]) | (uint16(file[1]) << 8)
	flag := file[0x02]
	pilotLen := uint(8063)
	if flag > 128 {
		pilotLen = 3223
	}
	block := &dataBlock{
		pilot:       2168,
		sync1:       667,
		sync2:       735,
		zero:        855,
		one:         1710,
		pilotLen:    pilotLen,
		pause:       3000,
		flag:        flag,
		lastBiteLen: 8,
		data:        file[2 : length+2],
	}
	return block, 2 + length
}

func readTzxBlock(file []byte) (interface{}, uint32) {
	id := file[0]
	file = file[1:]

	switch id {
	case 0x10:
		len := uint32(file[0x02]) | uint32(file[0x03])<<8
		flag := file[0x04]
		pilotLen := uint(8063)
		if flag > 128 {
			pilotLen = 3223
		}
		block := &dataBlock{
			id:          id,
			pilot:       2168,
			sync1:       667,
			sync2:       735,
			zero:        855,
			one:         1710,
			pilotLen:    pilotLen,
			pause:       (uint(file[0x00]) | uint(file[0x01])<<8),
			flag:        flag,
			data:        file[0x04 : len+0x04],
			lastBiteLen: 8,
		}
		return block, len + 0x04

	case 0x11:
		len := uint32(file[0x0f]) | uint32(file[0x10])<<8 | uint32(file[0x11])<<16
		block := &dataBlock{
			id:          id,
			pilot:       uint(file[0x00]) | uint(file[0x01])<<8,
			sync1:       uint(file[0x02]) | uint(file[0x03])<<8,
			sync2:       uint(file[0x04]) | uint(file[0x05])<<8,
			zero:        uint(file[0x06]) | uint(file[0x07])<<8,
			one:         uint(file[0x08]) | uint(file[0x09])<<8,
			pilotLen:    uint(file[0x0A]) | uint(file[0x0B])<<8,
			lastBiteLen: int8(file[0x0C]),
			pause:       uint(file[0x0D]) | uint(file[0x0E])<<8,
			flag:        file[0x12],
			data:        file[0x12 : len+0x12],
		}
		return block, len + 0x12

	case 0x12: // Pure Tone
		block := &dataBlock{
			id:       id,
			pilot:    (uint(file[0x00]) | uint(file[0x01])<<8),
			pilotLen: (uint(file[0x02]) | uint(file[0x03])<<8),
		}
		return block, 4

	case 0x13: // Pulse sequence
		len := int(file[0])
		block := &pulseSeqBlock{
			id: id,
		}
		for i := 0; i < len*2; i += 2 {
			block.pulses = append(block.pulses, uint(file[0x01+i])|uint(file[0x02+i])<<8)
		}
		return block, uint32(len*2 + 1)

	case 0x14:
		len := uint32(file[0x07]) | uint32(file[0x08])<<8 | uint32(file[0x09])<<16
		block := &dataBlock{
			id:          id,
			zero:        uint(file[0x00]) | uint(file[0x01])<<8,
			one:         uint(file[0x02]) | uint(file[0x03])<<8,
			pause:       uint(file[0x05]) | uint(file[0x06])<<8,
			data:        file[0x0a : len+0x0a],
			lastBiteLen: int8(file[0x04]),
		}
		return block, len + 0x0a

	case 0x20: // Pause (silence) or 'Stop the Tape' command
		block := &dataBlock{
			id:    id,
			pause: uint(file[0x00]) | uint(file[0x01])<<8,
		}
		return block, 2

	case 0x21: // Group start
		len := uint32(file[0])
		return nil, 1 + len

	case 0x22: // Group end
		return nil, 0

	case 0x24: // Loop start
		block := &loopBlock{
			id:    id,
			count: (int(file[0x00]) | int(file[0x01])<<8),
		}
		// println("Loop:", loop)
		return block, 2

	case 0x25: // TODO: Text description
		return &loopEndBlock{}, 0

	case 0x30: // TODO: Text description
		len := uint32(file[0])
		return nil, 1 + len

	case 0x32: // TODO: Archive info
		len := uint32(file[0x00]) | uint32(file[0x01])<<8
		return nil, 2 + len

	case 0x35: // TODO: Custom info block
		len := uint32(file[0x10]) | uint32(file[0x11])<<8 | uint32(file[0x12])<<16 | uint32(file[0x13])<<24
		return nil, 0x14 + len

	default:
		panic(fmt.Sprintf("id: 0x%02X", id))
	}
}

func (b *pulseSeqBlock) String() string {
	return fmt.Sprintf("(0x%02X) Pulse Seq. Block - pulses:%d", b.id, len(b.pulses))
}

func (b *dataBlock) String() string {
	name := b.name()
	return fmt.Sprintf("(0x%02X)(0x%02X) name: '%-9s' - datalen:%d - pilot:%d(%d) - sync:%d(%d) - zero:%d one:%d - pause:%d", b.id, b.flag, name, len(b.data), b.pilot, b.pilotLen, b.sync1, b.sync2, b.zero, b.one, b.pause)
}

func (b *dataBlock) name() string {
	if (b.flag == 0 || b.flag == 0x2c) && len(b.data) > 0 {
		return string(b.data[1:11])
	}
	return "Data"
}
