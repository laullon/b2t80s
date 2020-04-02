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

func TestST(t *testing.T) {
	memory := &dummyMemory{mem: make([]byte, 0xffff)}
	sp := NewStackPointer(memory)

	sp.Push(0xaabb)
	sp.Push(2)
	sp.Push(3)

	assert.Equal(t, uint16(3), sp.Pop())
	assert.Equal(t, uint16(2), sp.Pop())
	assert.Equal(t, uint16(0xaabb), sp.Pop())
}

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

	memory := &dummyMemory{mem: make([]byte, 0xffff)}
	memory.SetClock(&dummyClock{})

	cpu := NewZ80(memory, nil)

	debugger := NewDebugger(cpu, memory)
	cpu.SetDebuger(debugger)

	cpu.RegisterPort(emulator.PortMask{Mask: 0, Value: 0}, &dummyPM{})

	var idx int
	var test *cpuTest

	for idx, test = range tests {
		logger.Clear()
		// TODO make this test work
		if strings.HasPrefix(test.name, "eda2") ||
			strings.HasPrefix(test.name, "eda3") ||
			strings.HasPrefix(test.name, "edab") ||
			strings.HasPrefix(test.name, "edb2") ||
			strings.HasPrefix(test.name, "edb3") ||
			strings.HasPrefix(test.name, "edba") ||
			strings.HasPrefix(test.name, "edbb") ||
			strings.HasPrefix(test.name, "edaa") {
			continue
		}

		cpu.SetRegistersStr(test.registers, test.otherReg)
		for _, mem := range test.memory {
			for i, b := range mem.bytes {
				memory.PutByte(mem.start+uint16(i), b)
			}
		}

		_, err := GetOpCode(memory.GetBlock(0, 4))
		if err != nil {
			t.Logf("error on test '%v': %v", test.name, err)
			continue
		}

		log.Printf("\n")
		log.Printf("ready to test '%v' (%v/%v)", test.name, idx, len(tests))
		log.Printf("%s", hex.Dump(memory.GetBlock(0, 16)))
		log.Printf("start test '%v'", test.name)

		cpu.SetClock(&dummyClock{stopAtTSate: test.tStates})
		cpu.(*z80).halt = false
		err = cpu.RunFrame()
		if err != nil {
			t.Logf("error on test '%v' (%v/%v): %v", test.name, idx, len(tests), err)
			continue
		}

		log.Printf("%s", hex.Dump(memory.GetBlock(0, 16)))
		log.Printf("done test '%v'", test.name)
		result, ok := results[test.name]
		if ok {
			for _, ms := range result.memory {
				err, expt, org := ms.check(memory)
				t := assert.Nil(t, err, "test '%s' memory fail", test.name)
				if !t {
					log.Printf("0x%04X  mem: %s", ms.start, org)
					log.Printf("       expt: %s", expt)
					logger.Dump()
					return
				}
			}

			regis, sp, pc := cpu.DumpRegisters()
			registers := fmt.Sprintf(
				"%02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %04x %04x",
				regis[0], regis[1]&0b11010111, regis[2], regis[3], regis[4], regis[5], regis[6], regis[7],
				regis[12], regis[13]&0b11010111, regis[14], regis[15], regis[16], regis[17], regis[18], regis[19],
				regis[8], regis[9], regis[10], regis[11],
				sp, pc,
			)

			t := assert.Equal(t, result.registers, registers, "test '%s' registers fail", test.name)
			if !t {
				logger.Dump()
				return
			}
		} else {
			panic(fmt.Sprintf("result for test '%s' not found", test.name))
		}
	}
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
	str = strings.ReplaceAll(str, "-1", "")
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

func (ms *memoryState) check(mem emulator.Memory) (error, string, string) {
	for idx, b := range ms.bytes {
		if b != mem.GetByte(ms.start+uint16(idx)) {
			return fmt.Errorf("error on byte %d", idx),
				hex.Dump(ms.bytes),
				hex.Dump(mem.GetBlock(ms.start, uint16(len(ms.bytes))))
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

	mem := &basicMemory{}
	mem.memory = make([]byte, 0x0100)
	mem.memory = append(mem.memory, zexdoc...)
	mem.memory = append(mem.memory, make([]byte, 0x10000-len(mem.memory))...)

	cpu := NewZ80(mem, nil)
	cpu.SetClock(&dummyClock{})
	cpu.SetPC(0x100)
	cpu.RegisterTrap(0x5, func() uint16 {
		regs, _, _ := cpu.DumpRegisters()
		return printChar(regs[3], regs[4], regs[5], cpu.SP(), mem)
	})

	for {
		cpu.Step()
		if cpu.PC() == 0 {
			assert.NotContains(t, cpmScreen, "ERROR")
			return
		}
	}
}

// Emulate CP/M call 5; function is in register C.
// Function 2: print char in register E
// Function 9: print $ terminated string pointer in DE
func printChar(cpu_c, cpu_d, cpu_e byte, cpu_sp emulator.StackPointer, memory emulator.Memory) uint16 {
	switch byte(cpu_c) {
	case 2:
		cpmScreen = append(cpmScreen, cpu_e)
		fmt.Printf("%c", cpu_e)
	case 9:
		de := getRR(cpu_d, cpu_e)
		for addr := de; ; addr++ {
			ch := memory.GetByte(addr)
			if ch == '$' {
				break
			}
			cpmScreen = append(cpmScreen, ch)
			fmt.Printf("%c", ch)
		}
	}
	return cpu_sp.Pop()
}

// ***
// ***

type basicMemory struct {
	memory  []byte
	tStates *uint
	clock   emulator.Clock
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

func (mem *basicMemory) SetClock(clock emulator.Clock) {
	mem.clock = clock
}

// ***
// ***

type dummyMemory struct {
	mem []byte
}

func (m *dummyMemory) GetBlock(start, length uint16) []byte { return m.mem[start : start+length] }
func (m *dummyMemory) GetByte(pos uint16) byte              { return m.mem[pos] }
func (m *dummyMemory) PutByte(pos uint16, b byte)           { m.mem[pos] = b }
func (m *dummyMemory) GetWord(pos uint16) uint16 {
	return uint16(m.GetByte(pos)) | (uint16(m.GetByte(pos+1)))<<8
}

func (m *dummyMemory) PutWord(addr, w uint16) {
	m.PutByte(addr, uint8(w&0x00ff))
	m.PutByte(addr+1, uint8(w>>8))
}

func (m *dummyMemory) LoadRom(idx int, rom []byte)       {}
func (m *dummyMemory) SetClock(c emulator.Clock)         {}
func (m *dummyMemory) ReadPort(port uint16) (byte, bool) { return 0, false }
func (m *dummyMemory) WritePort(port uint16, data byte)  {}

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
func (c *dummyClock) ApplyDeplay()                          {}
func (c *dummyClock) FrameDone() bool                       { return c.ts >= c.stopAtTSate }
func (c *dummyClock) AddTicker(mod uint, t emulator.Ticker) {}
