package z80

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/laullon/b2t80s/emulator"
	"github.com/stretchr/testify/assert"
)

func TestRegPair(t *testing.T) {
	cpu := NewZ80(nil)
	cpu.Registers().(*Z80Registers).B = 0x0A
	cpu.Registers().(*Z80Registers).C = 0x0B
	assert.Equal(t, uint16(0x0A0B), cpu.Registers().(*Z80Registers).BC.Get())
}

// func TestSP(t *testing.T) {
// 	bus := &dummyBus{mem: make([]byte, 0xffff)}
// 	sp := NewStackPointer(memory)

// 	sp.Push(0xaabb)
// 	sp.Push(2)
// 	sp.Push(3)

// 	assert.Equal(t, uint16(3), sp.Pop())
// 	assert.Equal(t, uint16(2), sp.Pop())
// 	assert.Equal(t, uint16(0xaabb), sp.Pop())
// }

type cpuTest struct {
	name      string
	registers string
	otherReg  []byte
	tStates   uint
	memory    []*memoryState
}

type cpuTestResult struct {
	name      string
	registers string
	otherReg  string
	memory    []*memoryState
	endPC     uint16
}

var tests []*cpuTest
var results = make(map[string]*cpuTestResult)

type logger struct {
	log []string
}

func (l *logger) Write(p []byte) (n int, err error) {
	l.log = append(l.log, strings.Trim(string(p), "\n"))
	return 0, nil
}

func (l *logger) Clear() {
	l.log = nil
}

func (l *logger) Dump() {
	for _, msg := range l.log {
		fmt.Println(msg)
	}
}

func TestOPCodes(t *testing.T) {
	readTests()
	readTestsResults()

	logger := &logger{}
	log.SetOutput(logger)

	var idx int
	var test *cpuTest
	for idx, test = range tests {
		bus := &dummyBus{mem: make([]byte, 0xffff)}

		cpu := NewZ80(bus)

		// debugger := NewDebugger(cpu.(*z80), memory)
		// cpu.SetDebuger(debugger)

		cpu.RegisterPort(emulator.PortMask{Mask: 0, Value: 0}, &dummyPM{})

		result, ok := results[test.name]
		if !ok {
			panic(fmt.Sprintf("result for test '%s' not found", test.name))
		}

		logger.Clear()
		// TODO make this test work
		if strings.HasPrefix(test.name, "eda2") ||
			strings.HasPrefix(test.name, "eda3") ||
			strings.HasPrefix(test.name, "edab") ||
			strings.HasPrefix(test.name, "edb2") ||
			strings.HasPrefix(test.name, "edb3") ||
			strings.HasPrefix(test.name, "edba") ||
			strings.HasPrefix(test.name, "edbb") ||
			strings.HasPrefix(test.name, "edaa") ||
			strings.HasPrefix(test.name, "c7") { // TODO: end PC=0 ??
			continue
		}

		setRegistersStr(cpu.Registers().(*Z80Registers), test.registers, test.otherReg)

		for _, mem := range test.memory {
			for i, b := range mem.bytes {
				bus.mem[mem.start+uint16(i)] = b
			}
		}

		_, err := GetOpCode(bus.mem[0:4])
		if err != nil {
			t.Logf("error on test '%v': %v", test.name, err)
			continue
		}

		log.Printf("\n")
		log.Printf("ready to test '%v' (%v/%v)", test.name, idx, len(tests))
		log.Printf("%s", hex.Dump(bus.mem[0:16]))
		log.Printf("start test '%v' (endPC:%v)", test.name, result.endPC)

		fmt.Printf("start test '%v' (endPC:%v)\n", test.name, result.endPC)
		fmt.Printf("%s", hex.Dump(bus.mem[0:16]))

		cpu.(*z80).halt = false

		for {
			if cpu.(*z80).regs.M1 {
				if cpu.(*z80).regs.PC > result.endPC {
					break
				}
			}
			cpu.Tick()
		}

		cpu.Registers().(*Z80Registers).PC-- //???

		log.Printf("%s", hex.Dump(bus.mem[0:16]))
		log.Printf("done test '%v'", test.name)

		for _, ms := range result.memory {
			err, expt, org := ms.check(bus.mem)
			t := assert.Nil(t, err, "test '%s' memory fail", test.name)
			if !t {
				log.Printf("0x%04X  mem: %s", ms.start, org)
				log.Printf("       expt: %s", expt)
				logger.Dump()
				return
			}
		}

		regs := cpu.Registers().(*Z80Registers)
		registers := fmt.Sprintf(
			"%02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %04x %04x",
			regs.A, regs.F.GetByte()&0b11010111, regs.B, regs.C, regs.D, regs.E, regs.H, regs.L,
			regs.Aalt, regs.Falt.GetByte()&0b11010111, regs.Balt, regs.Calt, regs.Dalt, regs.Ealt, regs.Halt, regs.Lalt,
			regs.IXH, regs.IXL, regs.IYH, regs.IYL,
			regs.SP.Get(), regs.PC,
		)

		t := assert.Equal(t, result.registers, registers, "test '%s' registers fail", test.name)
		if !t {
			logger.Dump()
			return
		}
	}
}

func setRegistersStr(cpu *Z80Registers, line string, otherReg []byte) {
	regs := strings.Split(line, " ")
	cpu.A, _ = setRRstr(regs[0])
	cpu.B, cpu.C = setRRstr(regs[1])
	cpu.D, cpu.E = setRRstr(regs[2])
	cpu.H, cpu.L = setRRstr(regs[3])

	cpu.Aalt, _ = setRRstr(regs[4])
	cpu.Balt, cpu.Calt = setRRstr(regs[5])
	cpu.Dalt, cpu.Ealt = setRRstr(regs[6])
	cpu.Halt, cpu.Lalt = setRRstr(regs[7])

	cpu.IXH, cpu.IXL = setRRstr(regs[8])
	cpu.IYH, cpu.IYL = setRRstr(regs[9])

	s, p := setRRstr(regs[10])
	cpu.SP.Set(uint16(s)<<8 | uint16(p))

	p, c := setRRstr(regs[11])
	cpu.PC = uint16(p)<<8 | uint16(c)

	_, f := setRRstr(regs[0])
	cpu.F.SetByte(f)
	_, _f := setRRstr(regs[4])
	cpu.Falt.SetByte(_f)

	cpu.I = otherReg[0]
	cpu.R = otherReg[1]
	cpu.R7 = otherReg[1]
	cpu.IFF2 = otherReg[3] != 0
}

func parseMemory(str string) (pos uint16, b []byte) {
	return
}

func readTests() {
	file, err := os.Open("tests.in")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line := 0
	test := &cpuTest{}
	for scanner.Scan() {
		str := scanner.Text()
		if str == "-1" {
			line = 0
			tests = append(tests, test)
			test = &cpuTest{}
			continue
		}

		if len(str) == 0 {
			continue
		}

		switch line {
		case 0:
			test.name = str
		case 1:
			test.registers = str
		case 2:
			vals := strings.Split(str, " ")
			ts, err := strconv.ParseInt(vals[len(vals)-1], 10, 16)
			if err != nil {
				panic(err)
			}
			orstr := strings.Join(vals[0:len(vals)-1], "")
			test.otherReg, err = hex.DecodeString(orstr)
			if err != nil {
				panic(err)
			}
			test.tStates = uint(ts)
		default:
			test.memory = append(test.memory, parseMemoryState(str))
		}
		line++
	}
}

func readTestsResults() {
	file, err := os.Open("tests.out")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line := 0
	result := &cpuTestResult{}
	for scanner.Scan() {
		str := scanner.Text()
		if len(str) == 0 {
			line = 0
			results[result.name] = result
			result = &cpuTestResult{}
			continue
		}

		if strings.HasPrefix(str, "  ") {
			continue
		}

		switch line {
		case 0:
			result.name = str
		case 1:
			fstr := str[2:4] // igonoring flags 3&5 on the tests result
			bytes, err := hex.DecodeString(fstr)
			if err != nil {
				panic(fmt.Sprintf("str: '%v'(%v) error: %v", line, str, err))
			}
			f := bytes[0] & 0b11010111

			_fstr := str[22:24] // igonoring flags 3&5 on the tests result
			bytes, err = hex.DecodeString(_fstr)
			if err != nil {
				panic(fmt.Sprintf("str: '%v'(%v) error: %v", line, str, err))
			}
			_f := bytes[0] & 0b11010111
			str = fmt.Sprintf("%s%02x%s%02x%s", str[0:2], f, str[4:22], _f, str[24:])
			// fmt.Println(" ->", str)

			result.registers = str
			regs := strings.Split(str, " ")
			pc := regs[len(regs)-1]
			pcVal, err := strconv.ParseInt(pc, 16, 32)
			if err != nil {
				panic(fmt.Sprintf("str: '%v'(%v) error: %v", line, str, err))
			}
			result.endPC = uint16(pcVal)
		case 2: // TODO use it
		default:
			result.memory = append(result.memory, parseMemoryState(str))
		}
		line++
	}
	results[result.name] = result
}

type memoryState struct {
	start uint16
	bytes []byte
}

func parseMemoryState(line string) *memoryState {
	str := strings.ReplaceAll(line, " ", "")
	str = strings.ReplaceAll(str, "-1", "") // halt
	bytes, err := hex.DecodeString(str)
	if err != nil {
		panic(fmt.Sprintf("str: '%v'(%v) error: %v", line, str, err))
	}

	ms := &memoryState{
		start: uint16(bytes[0])<<8 | uint16(bytes[1]),
	}

	for i := 2; i < len(bytes); i++ {
		ms.bytes = append(ms.bytes, bytes[i])
	}

	return ms
}

func (ms *memoryState) check(mem []byte) (error, string, string) {
	for idx, b := range ms.bytes {
		if b != mem[ms.start+uint16(idx)] {
			return fmt.Errorf("error on byte %d", idx),
				hex.Dump(ms.bytes),
				hex.Dump(mem[ms.start : ms.start+uint16(len(ms.bytes))])
		}
	}
	return nil, "", ""
}

var cpmScreen []byte

func TestZEXDoc(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	f, err := os.Open("../tests/zout/zexdocsmall.cim")
	if err != nil {
		log.Fatal(err)
	}
	zexdoc, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	mem := make([]byte, 0x0100)
	mem = append(mem, zexdoc...)
	mem = append(mem, make([]byte, 0x10000-len(mem))...)

	cpu := NewZ80(&dummyBus{mem: mem})
	cpu.Registers().(*Z80Registers).PC = uint16(0x100)
	cpu.RegisterTrap(0x5, func() uint16 {
		return printChar(cpu.Registers().(*Z80Registers), mem)
	})

	for {
		cpu.Tick()
		if cpu.Registers().(*Z80Registers).PC == 0 {
			assert.NotContains(t, cpmScreen, "ERROR")
			return
		}
	}
}

func setRRstr(hl string) (uint8, uint8) {
	decoded, err := hex.DecodeString(hl)
	if err != nil {
		panic(fmt.Sprintf("string: '%v' error: %v", hl, err))
	}
	return decoded[0], decoded[1]
}

// Emulate CP/M call 5; function is in register C.
// Function 2: print char in register E
// Function 9: print $ terminated string pointer in DE
func printChar(regs *Z80Registers, memory []byte) uint16 {
	switch byte(regs.C) {
	case 2:
		cpmScreen = append(cpmScreen, regs.E)
		fmt.Printf("%c", regs.E)
	case 9:
		de := regs.DE.Get()
		for addr := de; ; addr++ {
			ch := memory[addr]
			if ch == '$' {
				break
			}
			cpmScreen = append(cpmScreen, ch)
			fmt.Printf("%c", ch)
		}
	}
	panic(-1)
	return 1 //regs.SP.Pop()
}

// ***
// ***

type basicMemory struct {
	memory  []byte
	tStates *uint
}

func (mem *basicMemory) LoadRom(idx int, rom []byte)       {}
func (mem *basicMemory) Paging(config byte)                {}
func (mem *basicMemory) ReadPort(port uint16) (byte, bool) { return 0, true }
func (mem *basicMemory) WritePort(port uint16, data byte)  {}

func (mem *basicMemory) GetBlock(start, length uint16) []byte {
	return mem.memory[start : start+length]
}
func (mem *basicMemory) GetByte(pos uint16) byte {
	return mem.memory[pos]
}
func (mem *basicMemory) PutByte(pos uint16, b byte) {
	mem.memory[pos] = b
}
func (mem *basicMemory) GetWord(pos uint16) uint16 {
	return uint16(mem.memory[pos+1])<<8 | uint16(mem.memory[pos])
}
func (mem *basicMemory) PutWord(pos, w uint16) {
	mem.memory[pos+1] = byte(w >> 8)
	mem.memory[pos] = byte(w)
}

// ***
// ***

type dummyBus struct {
	mem  []byte
	addr uint16
	data uint8
}

func (bus *dummyBus) SetAddr(addr uint16) { bus.addr = addr }
func (bus *dummyBus) GetAddr() uint16     { return bus.addr }

func (bus *dummyBus) SetData(data byte) { bus.data = data }
func (bus *dummyBus) GetData() byte     { return bus.data }

func (bus *dummyBus) ReadMemory()  { bus.data = bus.mem[bus.addr] }
func (bus *dummyBus) WriteMemory() { bus.mem[bus.addr] = bus.data }

func (bus *dummyBus) ReadPort()  {}
func (bus *dummyBus) WritePort() {}

// ***
// ***

type dummyPM struct{}

func (dummyPM) ReadPort(port uint16) (byte, bool) { return byte(port >> 8), false }
func (dummyPM) WritePort(port uint16, data byte)  {}

// ***
// ***

type dummyClock struct {
	stopAtTSate uint
	ts          uint
}

func (c *dummyClock) AddTStates(ts uint)                    { c.ts += ts }
func (c *dummyClock) GetTStates() uint                      { return c.ts }
func (c *dummyClock) FrameDone() bool                       { return c.ts >= c.stopAtTSate }
func (c *dummyClock) AddTicker(mod uint, t emulator.Ticker) {}
