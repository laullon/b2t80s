package z80

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/laullon/b2t80s/cpu"
	"github.com/stretchr/testify/assert"
)

func TestRegPair(t *testing.T) {
	cpu := NewZ80(nil)
	cpu.Registers().B = 0x0A
	cpu.Registers().C = 0x0B
	assert.Equal(t, uint16(0x0A0B), cpu.Registers().BC.Get())
}

type cpuTest struct {
	name      string
	registers string
	otherRegs auxRegs
	tStates   uint
	memory    []*memoryState
}

type cpuTestResult struct {
	name      string
	registers string
	otherRegs auxRegs
	memory    []*memoryState
	endPC     uint16
}

type auxRegs struct {
	I    byte
	R    byte
	IFF1 bool
	IFF2 bool
	IM   byte
	HALT bool
	TS   int
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
	readTests(t)
	readTestsResults(t)

	logger := &logger{}
	log.SetOutput(logger)

	var idx int
	var test *cpuTest
	bus := &dummyBus{mem: make([]byte, 0xffff)}
	for idx, test = range tests {

		cpu := NewZ80(bus)
		// cpu.SetDebuger(&dumpDebbuger{cpu: cpu.(*z80), log: false})

		result, ok := results[test.name]
		if !ok {
			assert.FailNowf(t, "error", "result for test '%s' not found", test.name)
		}

		logger.Clear()
		// TODO make this test work
		if strings.HasPrefix(test.name, "dd00") ||
			strings.HasPrefix(test.name, "ddfd00") { //???
			continue
		}

		setRegistersStr(cpu.Registers(), test.registers, test.otherRegs)

		for _, mem := range test.memory {
			for i, b := range mem.bytes {
				bus.mem[mem.start+uint16(i)] = b
			}
		}

		log.Printf("\n")
		log.Printf("ready to test '%v' (%v/%v)", test.name, idx, len(tests))
		log.Printf("%s", hex.Dump(bus.mem[0:16]))
		log.Printf("regs: %s", test.registers)
		log.Printf("start test '%v' (endPC:%v)", test.name, result.endPC)
		// fmt.Printf("start test '%v' (endPC:%v)\n", test.name, result.endPC)

		for i := 0; i < result.otherRegs.TS; i++ {
			cpu.Tick()
		}
		regs := cpu.Registers()

		ko := false
		ko = ko || !assert.Equal(t, result.endPC, regs.PC, "test '%s' cpu.PC fail", test.name)
		ko = ko || !assert.Equal(t, result.otherRegs.I, regs.I, "test '%s' cpu.I fail", test.name)
		ko = ko || !assert.Equal(t, result.otherRegs.R, regs.R, "test '%s' cpu.R fail", test.name)
		ko = ko || !assert.Equal(t, result.otherRegs.IFF1, regs.IFF1, "test '%s' cpu.IFF1 fail", test.name)
		ko = ko || !assert.Equal(t, result.otherRegs.IFF2, regs.IFF2, "test '%s' cpu.IFF2 fail", test.name)
		ko = ko || !assert.Equal(t, result.otherRegs.IM, regs.InterruptsMode, "test '%s' cpu.IM fail", test.name)
		if ko {
			logger.Dump()
			return
		}

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

func setRegistersStr(cpu *Z80Registers, line string, otherReg auxRegs) {
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

	cpu.I = otherReg.I
	cpu.R = otherReg.R
	cpu.IFF1 = otherReg.IFF1
	cpu.IFF2 = otherReg.IFF2
	cpu.InterruptsMode = otherReg.IM
}

func parseMemory(str string) (pos uint16, b []byte) {
	return
}

func readTests(t *testing.T) {
	file, err := os.Open("tests/tests.in")
	if err != nil {
		assert.FailNowf(t, "error!!", "%v", err)
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
			regs := strings.Split(regexp.MustCompile(" +").ReplaceAllString(str, " "), " ")
			test.otherRegs.I = parseHexUInt8(regs[0])
			test.otherRegs.R = parseHexUInt8(regs[1])
			test.otherRegs.IFF1 = parseHexUInt8(regs[2]) == 1
			test.otherRegs.IFF2 = parseHexUInt8(regs[3]) == 1
			test.otherRegs.IM = parseHexUInt8(regs[4])
			test.otherRegs.HALT = parseHexUInt8(regs[5]) == 1
			test.otherRegs.TS, _ = strconv.Atoi(regs[6])

		default:
			test.memory = append(test.memory, parseMemoryState(str))
		}
		line++
	}
}

func readTestsResults(t *testing.T) {
	file, err := os.Open("tests/tests.out")
	if err != nil {
		assert.FailNowf(t, "error!!", "%v", err)
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

		case 2:
			regs := strings.Split(regexp.MustCompile(" +").ReplaceAllString(str, " "), " ")
			result.otherRegs.I = byte(parseHexUInt8(regs[0]))
			result.otherRegs.R = byte(parseHexUInt8(regs[1]))
			result.otherRegs.IFF1 = parseHexUInt8(regs[2]) == 1
			result.otherRegs.IFF2 = parseHexUInt8(regs[3]) == 1
			result.otherRegs.IM = byte(parseHexUInt8(regs[4]))
			result.otherRegs.HALT = parseHexUInt8(regs[5]) == 1
			result.otherRegs.TS, _ = strconv.Atoi(regs[6])

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

	f, err := os.Open("tests/zexdocsmall.cim")
	if err != nil {
		assert.FailNowf(t, "error!!", "%v", err)
	}
	zexdoc, err := ioutil.ReadAll(f)
	if err != nil {
		assert.FailNowf(t, "error!!", "%v", err)
	}

	mem := make([]byte, 0x0100)
	mem = append(mem, zexdoc...)
	mem = append(mem, make([]byte, 0x10000-len(mem))...)

	cpu := NewZ80(&dummyBus{mem: mem})
	// cpu.SetDebuger(&dumpDebbuger{cpu: cpu.(*z80)})
	cpu.Registers().PC = uint16(0x100)
	cpu.RegisterTrap(0x5, func() {
		printChar(cpu.Registers(), mem)
	})

	for {
		cpu.Tick()
		if cpu.Registers().PC == 0 {
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
func printChar(regs *Z80Registers, memory []byte) {
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

	newPC := uint16(memory[regs.SP.Get()])
	newPC |= uint16(memory[regs.SP.Get()+1]) << 8
	regs.SP.Set(regs.SP.Get() + 2)
	regs.PC = newPC
}

// ***
// ***

type basicMemory struct {
	memory []byte
}

func (mem *basicMemory) GetByte(pos uint16) byte {
	return mem.memory[pos]
}
func (mem *basicMemory) PutByte(pos uint16, b byte) {
	mem.memory[pos] = b
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
func (bus *dummyBus) Release()     {}

func (bus *dummyBus) ReadPort()                                               { bus.data = uint8(bus.addr >> 8) }
func (bus *dummyBus) WritePort()                                              {}
func (bus *dummyBus) RegisterPort(mask cpu.PortMask, manager cpu.PortManager) {}

func (bus *dummyBus) GetBlock(addr uint16, l uint16) []byte { return nil }
