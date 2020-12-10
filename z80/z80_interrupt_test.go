package z80

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/laullon/b2t80s/emulator"
	"github.com/stretchr/testify/assert"
)

var start time.Time

func TestInterrupt(t *testing.T) {
	f, err := os.Open("Interrupt_test.z80.bin")
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

	bus := emulator.NewBus(&basicMemory{memory: mem})
	cpu := NewZ80(bus)
	cpu.SetDebuger(&dumpDebbuger{cpu: cpu.(*z80)})
	bus.RegisterPort(emulator.PortMask{Mask: 0, Value: 0}, tester)

	count := 0

	wait := time.Duration(20 * time.Millisecond)
	ticker := time.NewTicker(wait)
	go func() {
		for range ticker.C {
			cpu.Interrupt(true)
			count++
		}
	}()

	start = time.Now()
	for count <= 5 {
		cpu.Tick()
	}

	assert.Equal(t, 5, tester.c)
}

// ------------------

type counterHardware struct {
	c int
}

func (c *counterHardware) ReadPort(port uint16) (byte, bool) { return 0, false }
func (c *counterHardware) WritePort(port uint16, data byte)  { c.c++ }

//--------------------

type dumpDebbuger struct {
	cpu *z80
}

func (d *dumpDebbuger) AddInstruction(pc uint16, mem []byte) {
	op := decode(mem)
	if op == nil {
		return
	}
	opMem := make([]byte, op.len)
	copy(opMem, mem)
	le := &logEntry{op: op, mem: opMem, pc: pc}
	print(le.String())
	println(d.cpu.regs.dump())
}

func (d *dumpDebbuger) NextInstruction(mem []byte) {}
func (d *dumpDebbuger) SetDump(bool)               {}
func (d *dumpDebbuger) Tick()                      {}
func (d *dumpDebbuger) Stop()                      {}
func (d *dumpDebbuger) Continue()                  {}
func (d *dumpDebbuger) Step()                      {}
func (d *dumpDebbuger) StopNextFrame()             {}
func (d *dumpDebbuger) GetStatus() string          { return "nil" }

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
