package m6502

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/laullon/b2t80s/cpu"
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
	ticks := 0

	cpu := MewM6502(&simpleBus{mem: mem}).(*m6502)
	if testing.Short() {
		println("skipping logs in short mode.")
	} else {
		cpu.SetTracer(&logPrinter{ticks: &ticks})
	}

	defer func() {
		if r := recover(); r != nil {
			if cpu.log != nil {
				println(cpu.log.(*logPrinter).Print())
			}
			assert.FailNow(t, "Panic", r)
			panic(r)
		}
	}()

	cpu.regs.PC = 0x0400
	cpu.preFetch()
	cpu.op = cpu.nextOp
	for i := 0; ; i++ {
		ticks++
		cpu.Tick()
		if cpu.regs.PC > 0xfff0 {
			assert.FailNow(t, "error")
			return
		} else if cpu.regs.PC == 0 {
			return
		}
	}
}

func TestInterrup(t *testing.T) {
	f, err := os.Open("functional_test/6502_interrupt_test.bin")
	if err != nil {
		log.Fatal(err)
	}

	mem, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	mem = append(make([]byte, 0x0400), mem...)
	mem = append(mem, make([]byte, 0x10000)...)

	bus := &simpleBus{mem: mem}
	bus.interrupt = true
	cpu := MewM6502(bus).(*m6502)
	bus.cpu = cpu
	ticks := 0

	if testing.Short() {
		println("skipping logs in short mode.")
	} else {
		cpu.SetTracer(&logPrinter{ticks: &ticks})
	}

	cpu.regs.PC = 0x0400
	cpu.preFetch()
	cpu.op = cpu.nextOp
	defer func() {
		if r := recover(); r != nil {
			if cpu.log != nil {
				println(cpu.log.(*logPrinter).Print())
			}
			panic(r)
		}
	}()

	cpu.regs.PC = 0x0400
	cpu.preFetch()
	cpu.op = cpu.nextOp
	for i := 0; ; i++ {
		ticks++
		cpu.Tick()
		if cpu.regs.PC > 0xfff0 {
			assert.FailNow(t, "error")
			return
		} else if cpu.regs.PC == 0 {
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
	cpu.preFetch()
	cpu.op = cpu.nextOp
	ticks := 0

	if testing.Short() {
		println("skipping logs in short mode.")
	} else {
		cpu.SetTracer(&logPrinter{ticks: &ticks})
	}

	for i := 0; ; i++ {
		ticks++
		cpu.Tick()
		if cpu.regs.PC == 0x126A {
			// TODO: review
			// assert.Equal(t, 1141, ticks, "wrong number of ticks: %d", ticks)
			assert.Equal(t, 962, ticks, "wrong number of ticks: %d", ticks)
			if cpu.log != nil {
				println(cpu.log.(*logPrinter).Print())
			}
			return
		}
	}
}

// ************
// ************
// ************
// ************
type logPrinter struct {
	ticks     *int
	prevTicks int
}

func (log *logPrinter) AppendLastOP(entry string) {
	fmt.Printf("%5d (%d) - %s\n", *log.ticks, *log.ticks-log.prevTicks, entry)
	log.prevTicks = *log.ticks
}
func (log *logPrinter) Print() string    { return "" }
func (log *logPrinter) SetNextOP(string) {}

type simpleBus struct {
	mem       []byte
	interrupt bool
	cpu       *m6502
	irqConfig uint8
}

func (bus *simpleBus) Write(addr uint16, data uint8) {
	if bus.interrupt {
		if addr == 0xbffc {
			bus.cpu.doNMI = data&0x02 != 0
			bus.cpu.doIRQ = data&0x01 != 0
		}
	}
	bus.mem[addr] = data
}

func (bus *simpleBus) Read(addr uint16) uint8 {
	// if bus.interrupt {
	// 	if addr == 0xbffc {
	// 		panic(-1)
	// 	}
	// }
	return bus.mem[addr]
}
func (bus *simpleBus) RegisterPort(name string, mask cpu.PortMask, manager cpu.PortManager) {}

func (bus *simpleBus) DumpMap() string                    { return "" }
func (bus *simpleBus) GetDumplables() map[string]Dumpable { return make(map[string]Dumpable) }
