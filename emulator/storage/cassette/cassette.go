package cassette

import (
	"fmt"
	"path/filepath"
)

type Cassette interface {
	Tick()

	LoadTapFile(rom string)

	Name() string

	Ear() bool
	Motor(bool)
	IsMotorON() bool

	NextDataBlock() []byte
}

type pulse struct {
	length uint
	level  bool
}

type cassette struct {
	tap              *tap
	earPulse         uint
	earPulseDuration uint
	ear              bool
	earChannel       chan *pulse
	motor            bool
	nextBlogIdx      int
	name             string
}

func New() Cassette {
	c := &cassette{
		tap:        &tap{},
		earChannel: make(chan *pulse, 0xffff*8),
	}
	return c
}

func (c *cassette) NextDataBlock() []byte {
	if c.nextBlogIdx == len(c.tap.blocks) {
		return nil
	}

	b, ok := c.tap.blocks[c.nextBlogIdx].(*dataBlock)
	c.nextBlogIdx++
	if ok && len(b.data) > 0 {
		return b.data
	}
	return c.NextDataBlock()
}

func (c *cassette) Name() string {
	return c.name
}

func (c *cassette) LoadTapFile(path string) {
	c.tap = &tap{}
	c.tap.load(path)
	c.name = filepath.Base(path)
}

func (c *cassette) Motor(on bool) {
	c.motor = on
}

func (c *cassette) IsMotorON() bool {
	return c.motor
}

func (c *cassette) String() string {
	str := fmt.Sprintf("blocks:")
	for i, b := range c.tap.blocks {
		str = fmt.Sprintf("%s\n%d. %s", str, i, b)
	}
	return str
}

func (c *cassette) Ear() bool {
	return c.ear
}

var SpeedAdj = float64(1)

func adj(t uint) uint {
	v := float64(t) * SpeedAdj //1 //.1428571429
	return uint(v)
}

func (c *cassette) Tick() {
	if c.motor {
		if len(c.earChannel) == 0 {
			if len(c.tap.blocks) == c.nextBlogIdx {
				c.motor = false
				return
			}
			block := c.tap.blocks[c.nextBlogIdx]
			// fmt.Printf("%d - playing: %v \n", c.nextBlogIdx, block)
			if dataBlock, ok := block.(*dataBlock); ok {
				c.playDataBlock(dataBlock)
			} else if pulseSeqBlock, ok := block.(*pulseSeqBlock); ok {
				c.playPulseSeqBlock(pulseSeqBlock)
			} else {
				panic(block)
			}
			c.nextBlogIdx++
		} else {
			if c.earPulse == 0 {
				c.earPulse = adj((<-c.earChannel).length)
			}
			c.earPulseDuration++
			if c.earPulseDuration == c.earPulse {
				c.earPulseDuration = 0
				c.earPulse = 0
				c.ear = !c.ear //c.earPulse.level
			}
		}
	}
}

func (c *cassette) playPulseSeqBlock(block *pulseSeqBlock) {
	for _, pluse := range block.pulses {
		c.earChannel <- &pulse{pluse, true}
		c.earChannel <- &pulse{pluse, false}
	}
}

func (c *cassette) playDataBlock(block *dataBlock) {
	if block.pilotLen > 0 {
		for i := uint(0); i < block.pilotLen/2; i++ {
			c.earChannel <- &pulse{block.pilot, true}
			c.earChannel <- &pulse{block.pilot, false}
		}
	}

	if block.sync1 > 0 && block.sync2 > 0 {
		c.earChannel <- &pulse{block.sync1, true}
		c.earChannel <- &pulse{block.sync2, false}
	}

	if len(block.data) > 0 {
		for idx, b := range block.data {
			var bits int8 = 0
			if idx == len(block.data)-1 { // last byte?
				bits = 8 - block.lastBiteLen
			}
			for i := int8(7); i >= bits; i-- {
				if (b & (1 << i)) != 0 {
					c.earChannel <- &pulse{block.one, true}
					c.earChannel <- &pulse{block.one, false}
				} else {
					c.earChannel <- &pulse{block.zero, true}
					c.earChannel <- &pulse{block.zero, false}
				}
			}
		}
	}

	if block.pause > 0 {
		c.earChannel <- &pulse{block.pause * 4000, false}
	}
}
