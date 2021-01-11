package m6502

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mem []byte

func init() {
	f, err := os.Open("6502_functional_test.bin")
	if err != nil {
		log.Fatal(err)
	}

	mem, err = ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
}

func TestReset(t *testing.T) { // TODO: review
	cpu := newM6502(mem)

	for i := 0; i < 8; i++ {
		cpu.Tick()
		// fmt.Printf("%v\n", cpu.regs)
	}

	assert.Equal(t, uint16(0x37a3), cpu.regs.PC, "Bad PC")
	assert.Equal(t, uint8(0xfc), cpu.regs.SP, "Bad SP")
}

func TestFunctionalTests(t *testing.T) { // TODO: review
	cpu := newM6502(mem)

	cpu.regs.PC = 0x0400
	cpu.op = nil
	for i := 0; ; i++ {
		cpu.Tick()
		fmt.Printf("%5d - %v\n", i, cpu.regs)
	}
}
