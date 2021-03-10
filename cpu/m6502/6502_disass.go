package m6502

import (
	"fmt"
	"strconv"
	"strings"
)

func disassemble(pc uint16, bus *bus) string {
	sb := &strings.Builder{}

	for i := 0; i < 10; i++ {
		opCode := bus.Read(pc)
		nextOp := ops[opCode]
		if nextOp == nil {
			nextOp = &unsupported{}
		}

		op := nextOp.Clone()
		op.setPC(pc)
		pc++

		switch op.(type) {
		case *reset, *implicit, *unsupported, *brk, *push, *pull, *rts, *rti:

		case *immediate, *relative, *zeropage, *indirectY, *indirectX:
			op.setB1(bus.Read(pc))
			pc++

		case *absoluteJMP, *absoluteJSR, *absolute, *indirectJMP:
			op.setB1(bus.Read(pc))
			pc++
			op.setB2(bus.Read(pc))
			pc++

		default:
			panic(fmt.Sprintf("error on type %T", op))
		}

		sb.WriteString(dumpOperation(op))
		sb.WriteString("\n")
	}
	return sb.String()
}

type builder struct {
	buff []rune
	addr int
}

func (b *builder) WriteString(s string) {
	for _, v := range s {
		b.buff[b.addr] = v
		b.addr++
	}
	b.buff[b.addr] = 0
}

func (b *builder) Reset() {
	b.addr = 0
	b.buff[b.addr] = 0
}

func (b *builder) String() string {
	return string(b.buff[:sb.addr])
}

var sb = &builder{buff: make([]rune, 50)}

func dumpOperation(operation operation) string {
	sb.Reset()

	writePC(sb, operation.getPC())

	switch op := operation.(type) {
	case *brk:
		mod := ""
		if op.irq {
			mod = "-IRQ-"
		} else if op.imm {
			mod = "-IMM-"
		}
		writeMemory(sb, op.opCode)
		writeOP(sb, "BRK ", mod)

	case *reset:
		writeMemory(sb)
		writeOP(sb, "RESET")

	case *implicit:
		mod := ""
		if op.a {
			mod = " a"
		}
		writeMemory(sb, op.opCode)
		writeOP(sb, op.ins, mod)

	case *absoluteJMP:
		writeMemory(sb, op.opCode, op.b1, op.b2)
		writeOP(sb, "jmp $", toHex8(op.b2), toHex8(op.b1))

	case *absoluteJSR:
		writeMemory(sb, op.opCode, op.b1, op.b2)
		writeOP(sb, "jsr $", toHex8(op.b2), toHex8(op.b1))

	case *absolute:
		mod := ""
		if op.x {
			mod = ", X"
		} else if op.y {
			mod = ", Y"
		}
		writeMemory(sb, op.opCode, op.b1, op.b2)
		writeOP(sb, op.ins, " $", toHex8(op.b2), toHex8(op.b1), mod)

	case *unsupported:
		writeMemory(sb, op.opCode)
		writeOP(sb, "unsupported")

	case *immediate:
		writeMemory(sb, op.opCode, op.b1)
		writeOP(sb, op.ins, " #$", toHex8(op.b1))

	case *relative:
		writeMemory(sb, op.opCode, op.b1)
		writeOP(sb, op.ins, " $", toHex16(op.pc+uint16(int8(op.b1))))

	case *zeropage:
		mod := ""
		if op.x {
			mod = ", X"
		} else if op.y {
			mod = ", Y"
		}
		writeMemory(sb, op.opCode, op.b1)
		writeOP(sb, op.ins, " $", toHex8(op.b1), mod)

	case *push:
		writeMemory(sb, op.opCode)
		writeOP(sb, op.ins)

	case *pull:
		writeMemory(sb, op.opCode)
		writeOP(sb, op.ins)

	case *rts:
		writeMemory(sb, op.opCode)
		writeOP(sb, "rts")

	case *indirectY:
		writeMemory(sb, op.opCode, op.b1)
		writeOP(sb, op.ins, " ($", toHex8(op.b1), "), Y")

	case *indirectX:
		writeMemory(sb, op.opCode, op.b1)
		writeOP(sb, op.ins, " ($", toHex8(op.b1), ", X)")

	case *rti:
		writeMemory(sb, op.opCode)
		writeOP(sb, "rti")

	case *indirectJMP:
		writeMemory(sb, op.opCode, op.b1, op.b2)
		writeOP(sb, "jmp ($", toHex8(op.b2), toHex8(op.b1), ")")

	default:
		panic(fmt.Sprintf("error on type %T", op))
	}

	return sb.String()
}

func writePC(sb *builder, pc uint16) {
	sb.WriteString(strings.ToUpper(toHex16(pc)))
	sb.WriteString(": ")
}

func writeMemory(sb *builder, bytes ...uint8) {
	for _, b := range bytes {
		sb.WriteString(" ")
		sb.WriteString(toHex8(b))
	}
	sb.WriteString("             "[:10-(len(bytes)*3)+1])
}

func writeOP(sb *builder, strs ...string) {
	for _, str := range strs {
		sb.WriteString(strings.ToLower(str))
	}
}

func toHex8(v uint8) string {
	n := "0" + strconv.FormatUint(uint64(v), 16)
	return n[len(n)-2:]
}

func toHex16(v uint16) string {
	n := "000" + strconv.FormatUint(uint64(v), 16)
	return n[len(n)-4:]
}
