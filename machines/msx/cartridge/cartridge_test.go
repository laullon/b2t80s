package cartridge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKonami(t *testing.T) {
	rom := make([]byte, 16*0x2000)
	for i := 0; i < 16; i++ {
		rom[i*0x2000] = byte(i)
	}
	cart := NewKonami(rom)

	assert.Equal(t, byte(0xff), cart.Read(0x0000))
	assert.Equal(t, byte(0xff), cart.Read(0x2000))

	assert.Equal(t, byte(0x00), cart.Read(0x4000))
	assert.Equal(t, byte(0x01), cart.Read(0x6000))
	assert.Equal(t, byte(0x02), cart.Read(0x8000))
	assert.Equal(t, byte(0x03), cart.Read(0xa000))

	for b := 0; b < 4; b++ {
		for i := 0; i < 16; i++ {
			cart.Write(0x4542, byte(i+(16*b)))
			assert.Equal(t, byte(i), cart.Read(0x4000))
			cart.Write(0x6c63, byte(i+(16*b)))
			assert.Equal(t, byte(i), cart.Read(0x6000))
			cart.Write(0x8c63, byte(i+(16*b)))
			assert.Equal(t, byte(i), cart.Read(0x8000))
			cart.Write(0xac63, byte(i+(16*b)))
			assert.Equal(t, byte(i), cart.Read(0xa000))
		}

	}

	assert.Equal(t, byte(0xff), cart.Read(0xc000))
	assert.Equal(t, byte(0xff), cart.Read(0xe000))
	assert.Equal(t, byte(0xff), cart.Read(0xffff))
}
