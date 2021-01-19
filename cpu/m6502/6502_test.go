package m6502

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	cpuUtils "github.com/laullon/b2t80s/cpu"
	"github.com/stretchr/testify/assert"
)

// func TestReset(t *testing.T) { // TODO: review
// 	cpu := newM6502(mem)

// 	for i := 0; i < 8; i++ {
// 		cpu.Tick()
// 		// fmt.Printf("%v\n", cpu.regs)
// 	}

// 	assert.Equal(t, uint16(0x37a3), cpu.regs.PC, "Bad PC")
// 	assert.Equal(t, uint8(0xfc), cpu.regs.SP, "Bad SP")
// }

func TestFunctionalTests(t *testing.T) {
	f, err := os.Open("functional_test/6502_functional_test.bin")
	if err != nil {
		log.Fatal(err)
	}

	mem, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	cpu := newM6502(mem)
	if testing.Short() {
		println("skipping logs in short mode.")
	} else {
		cpu.log = cpuUtils.NewLogTail()
	}

	cpu.regs.PC = 0x0400
	cpu.op = nil
	for i := 0; ; i++ {
		cpu.Tick()
		if cpu.regs.PC > 0xfff0 {
			if !assert.NotEqual(t, uint16(0xffff), cpu.regs.PC, "ERROR !!!") {
				if cpu.log != nil {
					println(cpu.log.Print())
				}
			}
			return
		}
	}
}

func TestTiming(t *testing.T) {
	f, err := os.Open("timingtest/timingtest-1.bin")
	if err != nil {
		log.Fatal(err)
	}

	mem, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	mem = append(make([]byte, 0x1000), mem...)
	mem = append(mem, make([]byte, 0x1000)...)

	cpu := newM6502(mem)

	cpu.regs.PC = 0x1000
	cpu.op = nil
	ticks := 0
	cpu.log = &logPrinter{ticks: &ticks}
	for i := 0; ; i++ {
		ticks++
		cpu.Tick()
		if cpu.regs.PC == 0x1000 {
			assert.Equal(t, 1141, ticks, "wrong number of ticks: %d", ticks)
			return
		}
	}
}

type logPrinter struct {
	ticks     *int
	prevTicks int
}

func (log *logPrinter) AddEntry(entry string) {
	fmt.Printf("%5d (%d) - %s\n", *log.ticks, *log.ticks-log.prevTicks, entry)
	log.prevTicks = *log.ticks
}
func (log *logPrinter) Print() string { return "" }
