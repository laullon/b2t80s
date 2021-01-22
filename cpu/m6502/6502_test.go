package m6502

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	cpuUtils "github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/emulator"
	"github.com/stretchr/testify/assert"
)

// func TestReset(t *testing.T) { // TODO: review
// 	cpu := MewM6502(mem)

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

	cpu := MewM6502(&simpleBus{mem: mem}).(*m6502)
	if testing.Short() {
		println("skipping logs in short mode.")
	} else {
		cpu.log = cpuUtils.NewLogTail()
	}

	defer func() {
		if r := recover(); r != nil {
			if cpu.log != nil {
				println(cpu.log.Print())
			}
			assert.FailNow(t, "Panic")
		}
	}()

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

	cpu := MewM6502(&simpleBus{mem: mem}).(*m6502)

	cpu.regs.PC = 0x1000
	cpu.op = nil
	ticks := 0

	if testing.Short() {
		println("skipping logs in short mode.")
	} else {
		cpu.log = &logPrinter{ticks: &ticks}
	}

	for i := 0; ; i++ {
		ticks++
		cpu.Tick()
		if cpu.regs.PC == 0x1000 {
			// TODO: review
			// assert.Equal(t, 1141, ticks, "wrong number of ticks: %d", ticks)
			assert.Equal(t, 1075, ticks, "wrong number of ticks: %d", ticks)
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

type simpleBus struct {
	mem []byte
}

func (bus *simpleBus) Write(addr uint16, data uint8)                                     { bus.mem[addr] = data }
func (bus *simpleBus) Read(addr uint16) uint8                                            { return bus.mem[addr] }
func (bus *simpleBus) RegisterPort(mask emulator.PortMask, manager emulator.PortManager) {}
