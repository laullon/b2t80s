package z80

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/laullon/b2t80s/emulator"
)

type logEntry struct {
	pc  uint16
	ins emulator.Instruction
}

type debugger struct {
	cpu      *z80
	memory   emulator.Memory
	symbols  map[uint16]string
	log      []*logEntry
	stop     bool
	stopNext bool
	dumpNext bool
	dump     bool
	stopAt   uint16
	status   string
}

func NewDebugger(cpu emulator.CPU, mem emulator.Memory) emulator.Debugger {
	debug := &debugger{
		cpu:     cpu.(*z80),
		memory:  mem,
		symbols: make(map[uint16]string),
		stop:    false,
		// dump:    true,
	}
	cpu.SetDebuger(debug)
	return debug
}

func (debug *debugger) SetDump(on bool) {
	debug.dump = on
}

func (debug *debugger) SetBreakPoint(bp uint16) {
	debug.cpu.RegisterTrap(bp, func() uint16 {
		if debug.IsStoped() {
			return emulator.CONTINUE
		}
		debug.Stop()
		return emulator.STOP
	})
}

func (debug *debugger) LoadSymbols(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()
		str = strings.ReplaceAll(str, "\t", " ")
		str = strings.ReplaceAll(str, "  ", " ")
		parts := strings.Split(str, " ")
		addr, err := strconv.ParseUint(strings.Trim(parts[2], "H"), 16, 16)
		if err != nil {
			panic(err)
		}
		debug.symbols[uint16(addr)] = strings.Trim(parts[0], " ")
	}
}

func (debug *debugger) AddLastInstruction(ins emulator.Instruction) {
	debug.log = append(debug.log, &logEntry{ins: ins, pc: debug.cpu.PC()})
	if len(debug.log) > 10 {
		debug.log = debug.log[len(debug.log)-10 : len(debug.log)]
	}
	if debug.dump {
		regs, _, _ := debug.cpu.DumpRegisters()
		fmt.Printf(
			"                                  A:0x%02X F:%08b BC:0x%04X DE:0x%04X HL:0x%04X SP:0x%04X\n",
			uint16(regs[0]), uint16(regs[1]),
			uint16(regs[2])<<8|uint16(regs[3]),
			uint16(regs[4])<<8|uint16(regs[5]),
			uint16(regs[6])<<8|uint16(regs[7]),
			debug.cpu.sp.Get())
		// res.WriteString(fmt.Sprintf("  B:0x%02X    C:0x%02X  BC:0x%04X    ---------\n", regs[2], regs[3], )
		// res.WriteString(fmt.Sprintf("  D:0x%02X    E:0x%02X  DE:0x%04X    0x%04X\n", regs[4], regs[5], , debug.memory.GetWord(sp+0)))
		// res.WriteString(fmt.Sprintf("  H:0x%02X    L:0x%02X  HL:0x%04X    0x%04X\n", regs[6], regs[7], uint16(regs[6])<<8|uint16(regs[7]), debug.memory.GetWord(sp+2)))
		// res.WriteString(fmt.Sprintf("IXH:0x%02X  IXL:0x%02X  IX:0x%04X    0x%04X\n", regs[8], regs[9], uint16(regs[8])<<8|uint16(regs[9]), debug.memory.GetWord(sp+4)))
		// res.WriteString(fmt.Sprintf("IYH:0x%02X  IYL:0x%02X  IY:0x%04X    0x%04X\n", regs[10], regs[10], uint16(regs[10])<<8|uint16(regs[11]), debug.memory.GetWord(sp+6)))
		fmt.Println(ins.Dump(debug.cpu.PC()))
	}
}

func (debug *debugger) DumpNextFrame() {
	debug.dumpNext = true
}

func (debug *debugger) StopNextFrame() {
	debug.stopNext = true
}

var dumpedFrames int

func (debug *debugger) NextFrame() {
	if debug.dumpNext && !debug.dump {
		debug.dump = true
		dumpedFrames = 0
	} else if debug.dumpNext && debug.dump {
		if dumpedFrames == 5 {
			debug.dumpNext = false
			debug.dump = false
		}
		dumpedFrames++
	}

	if debug.stopNext {
		debug.stopNext = false
		debug.stop = true
	}
}

func (debug *debugger) Stop() {
	debug.stop = true
}

func (debug *debugger) Step() {
	if debug.stop {
		debug.cpu.Step()
	}
}

func (debug *debugger) IsStoped() bool {
	return debug.stop
}

func (debug *debugger) Continue() {
	debug.Step()
	debug.stop = false
}

func (debug *debugger) SetStatus(sts string) { debug.status = sts }
func (debug *debugger) GetStatus() string    { return debug.status }

func (debug *debugger) GetNextInstruction() string {
	ins, _ := GetOpCode(debug.memory.GetBlock(debug.cpu.PC(), 4))
	return ins.Dump(debug.cpu.PC())
}

func (debug *debugger) GetFollowingInstruction() string {
	var log []string
	pc := debug.cpu.PC()
	ins, _ := GetOpCode(debug.memory.GetBlock(pc, 4))
	pc += ins.Length
	for len(log) < 10 {
		ins, err := GetOpCode(debug.memory.GetBlock(pc, 4))
		if err == nil {
			log = append(log, ins.Dump(pc))
			pc += ins.Length
		}
	}
	return strings.Join(log, "\n")
}

func (debug *debugger) GetLog() string {
	var log []string
	for _, le := range debug.log {
		if le == nil {
			break
		}
		log = append(log, le.ins.Dump(le.pc))
	}
	return strings.Join(log, "\n")
}

func (debug *debugger) GetRegisters() string {
	var res strings.Builder
	regs, sp, _ := debug.cpu.DumpRegisters()
	res.WriteString(fmt.Sprintf("  A:0x%02X    F:0x%02X  AF:0x%04X    SP:0x%04X\n", regs[0], regs[1], uint16(regs[0])<<8|uint16(regs[1]), sp))
	res.WriteString(fmt.Sprintf("  B:0x%02X    C:0x%02X  BC:0x%04X    ---------\n", regs[2], regs[3], uint16(regs[2])<<8|uint16(regs[3])))
	res.WriteString(fmt.Sprintf("  D:0x%02X    E:0x%02X  DE:0x%04X    0x%04X\n", regs[4], regs[5], uint16(regs[4])<<8|uint16(regs[5]), getWord(debug.memory, sp+0)))
	res.WriteString(fmt.Sprintf("  H:0x%02X    L:0x%02X  HL:0x%04X    0x%04X\n", regs[6], regs[7], uint16(regs[6])<<8|uint16(regs[7]), getWord(debug.memory, sp+2)))
	res.WriteString(fmt.Sprintf("IXH:0x%02X  IXL:0x%02X  IX:0x%04X    0x%04X\n", regs[8], regs[9], uint16(regs[8])<<8|uint16(regs[9]), getWord(debug.memory, sp+4)))
	res.WriteString(fmt.Sprintf("IYH:0x%02X  IYL:0x%02X  IY:0x%04X    0x%04X\n", regs[10], regs[10], uint16(regs[10])<<8|uint16(regs[11]), getWord(debug.memory, sp+6)))
	res.WriteString(fmt.Sprintf("SZ5H3PNC\n%08b", regs[1]))
	return res.String()
}
