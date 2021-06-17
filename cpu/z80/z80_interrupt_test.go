package z80

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/laullon/b2t80s/cpu"
	"github.com/stretchr/testify/assert"
)

func TestInterrupt(t *testing.T) {
	f, err := os.Open("tests/Interrupt_test.z80.bin")
	if err != nil {
		log.Fatal(err)
	}

	binFile, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	mem := make([]byte, 0x0200)
	copy(mem, binFile)

	tester := &counterHardware{}

	memBus := cpu.NewBus("mem", &basicMemory{memory: mem})
	portsBus := cpu.NewBus("ports")
	bus := NewBus(memBus, portsBus)
	z80 := NewZ80(bus)
	// z80.SetTracer(&tracer{})
	portsBus.RegisterPort("tester", cpu.PortMask{Mask: 0, Value: 0}, tester)

	count := 0

	wait := time.Duration(20 * time.Millisecond)
	ticker := time.NewTicker(wait)
	go func() {
		for range ticker.C {
			z80.Interrupt(true)
			count++
		}
	}()

	for count <= 5 {
		z80.Tick()
	}

	assert.Equal(t, 5, tester.c)
}

// ------------------

type counterHardware struct {
	c int
}

func (c *counterHardware) ReadPort(port uint16) byte { return 0 }
func (c *counterHardware) WritePort(port uint16, data byte) {
	c.c++
}

//--------------------

type tracer struct{}

func (*tracer) AppendLastOP(op string)                                    { println(op) }
func (*tracer) SetNextOP(string)                                          {}
func (*tracer) SetDiss(pc uint16, getMemory func(pc, leng uint16) []byte) {}

//--------------------

//
// FROM: https://z80project.wordpress.com/2015/04/29/z80-interrupts-and-strings/
//
// .ORG $0000
// START:
//    DI
//    JP MAIN                       ;Jump to the MAIN routine

// .ORG $0100
// MAIN:
//    LD SP,$01ff
//    IM 1                          ;Use interrupt mode 1
//    EI                            ;Enable interrupts

// END_PROGRAM:
//    HALT
//    JP END_PROGRAM

// .ORG $0038
// MODE1_INTERRUPT:
//    DI                            ;Disable interrupts
//    EX AF,AF'                     ;Save register states
//    EXX                           ;Save register states
//    LD BC,0
//    LD A, 1
//    OUT (C), A
//    EXX                           ;Restore register states
//    EX AF,AF'                     ;Restore register states
//    EI                            ;Enable interrupts
//    RET

// .END
