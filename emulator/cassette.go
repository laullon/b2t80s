package emulator

import (
	"fmt"
	"os"
)

type Cassette interface {
	Ticker
	LoadTapFile(rom string)
	NextBlock() *Block
	Ready() bool

	Play() uint16
	Ear() bool
	Motor(bool)
	IsMotorON() bool
}

type pulse struct {
	length uint
	level  bool
}

type tapCassette struct {
	blocks      []*Block
	actualBlock int

	earPulse         pulse
	earPulseDuration uint
	ear              bool
	earChannel       chan pulse

	motor bool
}

type Block struct {
	flag byte
	data []byte

	pilot, pilotLen, sync1, sync2, zero, one uint
	pause                                    uint
	lastBiteLen                              int8
}

func NewTapCassette() Cassette {
	c := &tapCassette{
		earChannel: make(chan pulse, 0xffff*8),
	}

	// ticker := time.NewTicker(time.Duration(5 * time.Second))
	// go func() {
	// 	for range ticker.C {
	// 		if c.motor {
	// 			println("len(c.earChannel)", len(c.earChannel), "c.earPulseDuration", c.earPulseDuration, "-", c.earPulse.length)
	// 		}
	// 	}
	// }()

	return c
}

func (c *tapCassette) LoadTapFile(rom string) {
	fi, err := os.Stat(rom)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(rom)
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
		// vMajor := file[8]
		// vMinor := file[9]
		// fmt.Printf("Loading %s V%d.%d\n", header, vMajor, vMinor)
		file = file[10:]
		for len(file) > 0 {
			file = c.readTzxBlock(file)
		}
	} else {
		for len(file) > 0 {
			block, l := readBlock(file)
			file = file[l:]
			c.blocks = append(c.blocks, block)
		}
	}
}

func (c *tapCassette) Ready() bool {
	return len(c.blocks) != 0
}

func (c *tapCassette) Motor(on bool) {
	if c.motor != on {
		// println("Motor =>", on)
	}
	c.motor = on
}

func (c *tapCassette) IsMotorON() bool {
	return c.motor
}

func (c *tapCassette) NextBlock() *Block {
	defer func() {
		c.actualBlock++
		if c.actualBlock == len(c.blocks) {
			c.actualBlock = 0
		}
	}()
	return c.blocks[c.actualBlock]
}

func (c *tapCassette) String() string {
	str := fmt.Sprintf("blocks:")
	for i, b := range c.blocks {
		str = fmt.Sprintf("%s\n%d. %s", str, i, b)
	}
	return str
}

func (c *tapCassette) readTzxBlock(file []byte) []byte {
	id := file[0]
	file = file[1:]
	// fmt.Printf("id: 0x%02X\n", id)

	switch id {
	case 0x10:
		len := uint32(file[0x02]) | uint32(file[0x03])<<8
		flag := file[0x04]
		pilotLen := uint(8063)
		if flag > 128 {
			pilotLen = 3223
		}
		block := &Block{
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
		c.blocks = append(c.blocks, block)
		return file[len+0x04:]

	case 0x11:
		len := uint32(file[0x0f]) | uint32(file[0x10])<<8 | uint32(file[0x11])<<16
		block := &Block{
			pilot:       uint(file[0x00]) | uint(file[0x01])<<8,
			sync1:       uint(file[0x02]) | uint(file[0x03])<<8,
			sync2:       uint(file[0x04]) | uint(file[0x05])<<8,
			zero:        uint(file[0x06]) | uint(file[0x07])<<8,
			one:         uint(file[0x08]) | uint(file[0x09])<<8,
			pilotLen:    uint(file[0x0A]) | uint(file[0x0B])<<8,
			pause:       uint(file[0x0D]) | uint(file[0x0E])<<8,
			flag:        file[0x12],
			data:        file[0x12 : len+0x12],
			lastBiteLen: int8(file[0x0C]),
		}
		c.blocks = append(c.blocks, block)
		return file[len+0x12:]

	case 0x12: // Pure Tone
		return file[4:]

	case 0x20: // Pause
		block := &Block{
			pause: (uint(file[0x00]) | uint(file[0x01])<<8),
		}
		c.blocks = append(c.blocks, block)
		return file[2:]

	case 0x21: // Group start
		len := file[0]
		// txt := string(file[1 : 1+len])
		// println("Group:", txt)
		return file[1+len:]

	case 0x24: // Loop start
		// loop := (uint(file[0x00]) | uint(file[0x01])<<8)
		// println("Loop:", loop)
		return file[2:]

	case 0x32: // Archive info
		len := uint32(file[0x00]) | uint32(file[0x01])<<8
		return file[2+len:]

	case 0x30: // Text description
		len := file[0]
		// txt := string(file[1 : 1+len])
		// println(txt)
		return file[1+len:]

	default:
		panic(fmt.Sprintf("id: 0x%02X", id))
	}
}

func readBlock(file []byte) (*Block, uint16) {
	length := uint16(file[0]) | (uint16(file[1]) << 8)
	flag := file[0x02]
	pilotLen := uint(8063)
	if flag > 128 {
		pilotLen = 3223
	}
	block := &Block{
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

func (b *Block) Name() string {
	if len(b.data) == 0 {
		return fmt.Sprintf("Pause (%d)", b.pause)
	} else if b.flag == 0 {
		return string(b.data[1:11])
	} else if b.flag == 0x2c {
		return fmt.Sprintf("%s (%d/%v/%v)", string(b.data[1:17]), b.data[23], b.data[24] != 0, b.data[18] != 0)
	}
	return "Data"
}

func (b *Block) Type() byte {
	return b.flag
}

func (b *Block) GetData() []byte {
	return b.data
}

func (b *Block) String() string {
	return fmt.Sprintf("(0x%02X) fileName: '%s' - datalen:%d - zero:%d one:%d", b.Type(), b.Name(), len(b.data), b.zero, b.one)
}

func (c *tapCassette) Ear() bool {
	return c.ear
}

func (c *tapCassette) Tick() {
	if c.motor {
		c.earPulseDuration++
		if c.earPulseDuration >= c.earPulse.length {
			c.earPulseDuration = 0
			if len(c.earChannel) > 0 {
				c.earPulse = <-c.earChannel
				c.ear = c.earPulse.level
			}
		}
	}
}

func (c *tapCassette) Play() uint16 {
	if len(c.blocks) == 0 {
		return 0
	}

	if len(c.blocks[0].data) > 0 {
		block := &Block{
			pause: 5000,
		}
		c.blocks = append([]*Block{block}, c.blocks...)
	}
	go func() {
		// println("play start")
		for _, block := range c.blocks {
			// println("loading", block.Name())
			c.playBlock(block)
		}
		// println("play done")
	}()
	return CONTINUE
}

func (c *tapCassette) playBlock(block *Block) {
	if len(block.GetData()) > 0 {
		for i := uint(0); i < block.pilotLen/2; i++ {
			c.earChannel <- pulse{block.pilot, true}
			c.earChannel <- pulse{block.pilot, false}
		}

		c.earChannel <- pulse{block.sync1, true}
		c.earChannel <- pulse{block.sync2, false}

		for idx, b := range block.GetData() {
			var bits int8 = 8 - 8
			if idx == len(block.GetData())-1 { // last byte?
				bits = 8 - block.lastBiteLen
			}
			for i := int8(7); i >= bits; i-- {
				if (b & (1 << i)) != 0 {
					c.earChannel <- pulse{block.one, true}
					c.earChannel <- pulse{block.one, false}
				} else {
					c.earChannel <- pulse{block.zero, true}
					c.earChannel <- pulse{block.zero, false}
				}
			}
		}
		c.earChannel <- pulse{200, true}
	}
	c.earChannel <- pulse{block.pause * 4000, false}
}
