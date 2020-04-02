package z80

import (
	"fmt"

	"github.com/laullon/b2t80s/emulator"
)

func (cpu *z80) runSwitch(ins emulator.Instruction) bool {
	needPcUpdate := true
	switch ins.Instruction {
	case 0x00: // NOP

	case 0x07: // RLCA
		cpu.a = cpu.a<<1 | cpu.a>>7
		cpu.f.C = cpu.a&0x01 != 0
		cpu.f.H = false
		cpu.f.N = false

	case 0xcb17: // RL A
		cpu.rl(&cpu.a)

	case 0xcb10: // RL B
		cpu.rl(&cpu.b)

	case 0xcb11: // RL C
		cpu.rl(&cpu.c)

	case 0xcb12: // RL D
		cpu.rl(&cpu.d)

	case 0xcb13: // RL E
		cpu.rl(&cpu.e)

	case 0xcb14: // RL H
		cpu.rl(&cpu.h)

	case 0xcb15: // RL L
		cpu.rl(&cpu.l)

	case 0xcb16: // RL (HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.rl(&b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xfdcb16: // RL (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		b := cpu.memory.GetByte(iy)
		cpu.rl(&b)
		cpu.memory.PutByte(iy, b)

	case 0xcb1f: // RR A
		cpu.rr(&cpu.a)

	case 0xcb18: // RR B
		cpu.rr(&cpu.b)

	case 0xcb19: // RR C
		cpu.rr(&cpu.c)

	case 0xcb1a: // RR D
		cpu.rr(&cpu.d)

	case 0xcb1b: // RR E
		cpu.rr(&cpu.e)

	case 0xcb1c: // RR H
		cpu.rr(&cpu.h)

	case 0xcb1d: // RR L
		cpu.rr(&cpu.l)

	case 0xcb1e: // RR (HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.rr(&b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xfdcb1e: // RR (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		b := cpu.memory.GetByte(iy)
		cpu.rr(&b)
		cpu.memory.PutByte(iy, b)

	case 0xcb2f: // SRA A
		cpu.sra(&cpu.a)

	case 0xcb28: // SRA B
		cpu.sra(&cpu.b)

	case 0xcb29: // SRA C
		cpu.sra(&cpu.c)

	case 0xcb2a: // SRA D
		cpu.sra(&cpu.d)

	case 0xcb2b: // SRA E
		cpu.sra(&cpu.e)

	case 0xcb2c: // SRA H
		cpu.sra(&cpu.h)

	case 0xcb2d: // SRA L
		cpu.sra(&cpu.l)

	case 0xcb2e: // SRA (HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.sra(&b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcb27: // SLA A
		cpu.sla(&cpu.a)

	case 0xcb20: // SLA B
		cpu.sla(&cpu.b)

	case 0xcb21: // SLA C
		cpu.sla(&cpu.c)

	case 0xcb22: // SLA D
		cpu.sla(&cpu.d)

	case 0xcb23: // SLA E
		cpu.sla(&cpu.e)

	case 0xcb24: // SLA H
		cpu.sla(&cpu.h)

	case 0xcb25: // SLA L
		cpu.sla(&cpu.l)

	case 0xcb26: // SLA (HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.sla(&b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xfdcb26: // SLA (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		b := cpu.memory.GetByte(iy)
		cpu.sla(&b)
		cpu.memory.PutByte(iy, b)

	case 0xcb07: // RLC A
		cpu.rlc(&cpu.a)

	case 0xcb00: // RLC B
		cpu.rlc(&cpu.b)

	case 0xcb01: // RLC C
		cpu.rlc(&cpu.c)

	case 0xcb02: // RLC D
		cpu.rlc(&cpu.d)

	case 0xcb03: // RLC E
		cpu.rlc(&cpu.e)

	case 0xcb04: // RLC H
		cpu.rlc(&cpu.h)

	case 0xcb05: // RLC L
		cpu.rlc(&cpu.l)

	case 0xcb06: // RLC (HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.rlc(&b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xfdcb06: // RLC (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		b := cpu.memory.GetByte(iy)
		cpu.rlc(&b)
		cpu.memory.PutByte(iy, b)

	case 0xcb0f: // RRC A
		cpu.rrc(&cpu.a)

	case 0xcb08: // RRC B
		cpu.rrc(&cpu.b)

	case 0xcb09: // RRC C
		cpu.rrc(&cpu.c)

	case 0xcb0a: // RRC D
		cpu.rrc(&cpu.d)

	case 0xcb0b: // RRC E
		cpu.rrc(&cpu.e)

	case 0xcb0c: // RRC H
		cpu.rrc(&cpu.h)

	case 0xcb0d: // RRC L
		cpu.rrc(&cpu.l)

	case 0xcb0e: // RRC (HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.rrc(&b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xfdcb0e: // RRC (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		b := cpu.memory.GetByte(iy)
		cpu.rrc(&b)
		cpu.memory.PutByte(iy, b)

	case 0x17: // RLA
		c := cpu.f.C
		cpu.f.C = cpu.a&0b10000000 != 0
		cpu.a = (cpu.a << 1)
		if c {
			cpu.a |= 1
		}
		cpu.f.H = false
		cpu.f.N = false

	case 0x0f: // RRCA
		cpu.f.C = cpu.a&0x01 != 0
		cpu.f.H = false
		cpu.f.N = false
		cpu.a = (cpu.a >> 1) | (cpu.a << 7)

	case 0x1f: // RRA
		c := cpu.f.C
		cpu.f.C = cpu.a&1 != 0
		cpu.a = (cpu.a >> 1)
		if c {
			cpu.a |= 0b10000000
		}
		cpu.f.H = false
		cpu.f.N = false

	case 0xed67: // RRD
		hl := getRR(cpu.h, cpu.l)
		hlv := cpu.memory.GetByte(hl)
		cpu.memory.PutByte(hl, (cpu.a<<4 | hlv>>4))
		cpu.a = (cpu.a & 0xf0) | (hlv & 0x0f)

		cpu.f.S = cpu.a&0x80 != 0
		cpu.f.Z = cpu.a == 0
		cpu.f.P = parityTable[cpu.a]
		cpu.f.H = false
		cpu.f.N = false

	case 0xed6f: // RLD
		hl := getRR(cpu.h, cpu.l)
		hlv := cpu.memory.GetByte(hl)
		cpu.memory.PutByte(hl, (hlv<<4 | cpu.a&0x0f))
		cpu.a = (cpu.a & 0xf0) | (hlv >> 4)

		cpu.f.S = cpu.a&0x80 != 0
		cpu.f.Z = cpu.a == 0
		cpu.f.P = parityTable[cpu.a]
		cpu.f.H = false
		cpu.f.N = false

	case 0x3f: // CCF
		cpu.f.H = cpu.f.C
		cpu.f.N = false
		cpu.f.C = !cpu.f.C

	case 0xcd: // CALL NN
		nn := toWord(ins.Mem[1], ins.Mem[2])
		cpu.pc += uint16(ins.Length)
		cpu.sp.Push(cpu.pc)
		cpu.pc = nn
		needPcUpdate = false

	case 0xec: // CALL PE,NN
		if cpu.f.P {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.pc += uint16(ins.Length)
			cpu.sp.Push(cpu.pc)
			cpu.pc = nn
			needPcUpdate = false
		}

	case 0xe4: // CALL PO,NN
		if !cpu.f.P {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.pc += uint16(ins.Length)
			cpu.sp.Push(cpu.pc)
			cpu.pc = nn
			needPcUpdate = false
		}

	case 0xd4: // CALL NC,NN
		if cpu.f.C == false {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.pc += uint16(ins.Length)
			cpu.sp.Push(cpu.pc)
			cpu.pc = nn
			needPcUpdate = false
		}

	case 0xdc: // CALL C,NN
		if cpu.f.C {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.pc += uint16(ins.Length)
			cpu.sp.Push(cpu.pc)
			cpu.pc = nn
			needPcUpdate = false
		}

	case 0xcc: // CALL Z,NN
		if cpu.f.Z == true {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.pc += uint16(ins.Length)
			cpu.sp.Push(cpu.pc)
			cpu.pc = nn
			needPcUpdate = false
		}

	case 0xc4: // CALL NZ,NN
		if cpu.f.Z == false {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.pc += uint16(ins.Length)
			cpu.sp.Push(cpu.pc)
			cpu.pc = nn
			needPcUpdate = false
		}

	case 0xf4: // CALL P,NN
		if !cpu.f.S {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.pc += uint16(ins.Length)
			cpu.sp.Push(cpu.pc)
			cpu.pc = nn
			needPcUpdate = false
		}

	case 0xfc: // CALL M,NN
		if cpu.f.S {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.pc += uint16(ins.Length)
			cpu.sp.Push(cpu.pc)
			cpu.pc = nn
			needPcUpdate = false
		}

	case 0xf0: // RET P
		if !cpu.f.S {
			cpu.pc = cpu.sp.Pop()
			needPcUpdate = false
		}

	case 0xf8: // RET M
		if cpu.f.S {
			cpu.pc = cpu.sp.Pop()
			needPcUpdate = false
		}

	case 0xc8: // RET Z
		if cpu.f.Z {
			cpu.pc = cpu.sp.Pop()
			needPcUpdate = false
		}

	case 0xc0: // RET NZ
		if !cpu.f.Z {
			cpu.pc = cpu.sp.Pop()
			needPcUpdate = false
		}

	case 0xd8: // RET C
		if cpu.f.C {
			cpu.pc = cpu.sp.Pop()
			needPcUpdate = false
		}

	case 0xe8: // RET PE
		if cpu.f.P {
			cpu.pc = cpu.sp.Pop()
			needPcUpdate = false
		}

	case 0xe0: // RET PO
		if !cpu.f.P {
			cpu.pc = cpu.sp.Pop()
			needPcUpdate = false
		}

	case 0xd0: // RET NC
		if !cpu.f.C {
			cpu.pc = cpu.sp.Pop()
			needPcUpdate = false
		}

	case 0xc9: // RET
		cpu.pc = cpu.sp.Pop()
		needPcUpdate = false

	case 0xed4d: // RETI
		cpu.pc = cpu.sp.Pop()
		needPcUpdate = false

	case 0xed45: // RETN
		// TODO IFF1=IFF2 ??
		cpu.pc = cpu.sp.Pop()
		needPcUpdate = false

	case 0xc7: // RST 0
		cpu.sp.Push(cpu.pc + ins.Length)
		cpu.pc = 0
		needPcUpdate = false

	case 0xcf: // RST 08H
		cpu.sp.Push(cpu.pc + ins.Length)
		cpu.pc = 0x08
		needPcUpdate = false

	case 0xd7: // RST 10H
		cpu.sp.Push(cpu.pc + ins.Length)
		cpu.pc = 0x10
		needPcUpdate = false

	case 0xdf: // RST 18H
		cpu.sp.Push(cpu.pc + ins.Length)
		cpu.pc = 0x18
		needPcUpdate = false

	case 0xe7: // RST 20H
		cpu.sp.Push(cpu.pc + ins.Length)
		cpu.pc = 0x20
		needPcUpdate = false

	case 0xef: // RST 28H
		cpu.sp.Push(cpu.pc + ins.Length)
		cpu.pc = 0x28
		needPcUpdate = false

	case 0xf7: // RST 30H
		cpu.sp.Push(cpu.pc + ins.Length)
		cpu.pc = 0x30
		needPcUpdate = false

	case 0xff: // RST 38H
		cpu.sp.Push(cpu.pc + ins.Length)
		cpu.pc = 0x38
		needPcUpdate = false

	case 0xf5: // PUSH AF
		af := getRR(cpu.a, cpu.f.getByte())
		cpu.sp.Push(af)

	case 0xc5: // PUSH BC
		bc := getRR(cpu.b, cpu.c)
		cpu.sp.Push(bc)

	case 0xd5: // PUSH DE
		de := getRR(cpu.d, cpu.e)
		cpu.sp.Push(de)

	case 0xe5: // PUSH HL
		hl := getRR(cpu.h, cpu.l)
		cpu.sp.Push(hl)

	case 0xdde5: // PUSH IX
		ix := getRR(cpu.ixh, cpu.ixl)
		cpu.sp.Push(ix)

	case 0xfde5: // PUSH IY
		iy := getRR(cpu.iyh, cpu.iyl)
		cpu.sp.Push(iy)

	case 0xf1: // POP AF
		af := cpu.sp.Pop()
		a, f := setRR(af)
		cpu.a = a
		cpu.f.setByte(f)

	case 0xc1: // POP BC
		bc := cpu.sp.Pop()
		cpu.b, cpu.c = setRR(bc)

	case 0xd1: // POP DE
		de := cpu.sp.Pop()
		cpu.d, cpu.e = setRR(de)

	case 0xe1: // POP HL
		hl := cpu.sp.Pop()
		cpu.h, cpu.l = setRR(hl)

	case 0xdde1: // POP IX
		ix := cpu.sp.Pop()
		cpu.ixh, cpu.ixl = setRR(ix)

	case 0xfde1: // POP IY
		iy := cpu.sp.Pop()
		cpu.iyh, cpu.iyl = setRR(iy)

	case 0xF3: // DI
		cpu.iff1 = false

	case 0xfb: // EI
		cpu.iff1 = true

	case 0xed46: // IM 0
		cpu.interruptsMode = 0

	case 0xed56: // IM 1
		cpu.interruptsMode = 1

	case 0xed5e: // IM 2
		cpu.interruptsMode = 2

	case 0x27: // DAA
		c := cpu.f.C
		add := byte(0)
		if cpu.f.H || ((cpu.a & 0x0f) > 9) {
			add = 6
		}
		if c || (cpu.a > 0x99) {
			add |= 0x60
		}
		if cpu.a > 0x99 {
			c = true
		}
		if cpu.f.N {
			cpu.subA(add)
		} else {
			cpu.addA(add)
		}
		cpu.f.S = int8(cpu.a) < 0
		cpu.f.Z = cpu.a == 0
		cpu.f.P = parityTable[cpu.a]
		cpu.f.C = c

	case 0xAF: // XOR A
		cpu.xor(cpu.a)

	case 0xa8: // XOR B
		cpu.xor(cpu.b)

	case 0xa9: // XOR C
		cpu.xor(cpu.c)

	case 0xaa: // XOR D
		cpu.xor(cpu.d)

	case 0xab: // XOR E
		cpu.xor(cpu.e)

	case 0xac: // XOR H
		cpu.xor(cpu.h)

	case 0xad: // XOR L
		cpu.xor(cpu.l)

	case 0xae: // XOR (HL)
		hl := getRR(cpu.h, cpu.l)
		cpu.xor(cpu.memory.GetByte(hl))

	case 0xee: // XOR N
		cpu.xor(ins.Mem[1])

	case 0x02: // LD (BC),A
		pos := getRR(cpu.b, cpu.c)
		cpu.memory.PutByte(pos, cpu.a)

	case 0x12: // LD (DE),A
		pos := getRR(cpu.d, cpu.e)
		cpu.memory.PutByte(pos, cpu.a)

	case 0x7f: // LD A,A

	case 0x78: // LD A,B
		cpu.a = cpu.b

	case 0x79: // LD A,C
		cpu.a = cpu.c

	case 0x7a: // LD A,D
		cpu.a = cpu.d

	case 0x7b: // LD A,E
		cpu.a = cpu.e

	case 0x7c: // LD A,H
		cpu.a = cpu.h

	case 0x7d: // LD A,L
		cpu.a = cpu.l

	case 0xed57: // LD A,I
		cpu.a = cpu.i
		cpu.f.S = cpu.a&0x80 != 0
		cpu.f.Z = cpu.a == 0
		cpu.f.H = false
		cpu.f.P = cpu.iff2
		cpu.f.N = false

	case 0xed5f: // LD A,R TODO: review this and its test
		cpu.a = (cpu.r & 0x7f) | (cpu.r7 & 0x80)
		cpu.f.S = cpu.a&0x80 != 0
		cpu.f.Z = cpu.a == 0
		cpu.f.H = false
		cpu.f.P = cpu.iff2
		cpu.f.N = false

	case 0x3e: // LD A,n
		cpu.a = ins.Mem[1]

	case 0x47: // LD B,A
		cpu.b = cpu.a

	case 0x40: // LD B,B
		// cpu.b = cpu.b

	case 0x41: // LD B,C
		cpu.b = cpu.c

	case 0x42: // LD B,D
		cpu.b = cpu.d

	case 0x43: // LD B,E
		cpu.b = cpu.e

	case 0x44: // LD B,H
		cpu.b = cpu.h

	case 0x45: // LD B,L
		cpu.b = cpu.l

	case 0x46: // LD B,(HL)
		cpu.b = cpu.memory.GetByte(getRR(cpu.h, cpu.l))

	case 0x4f: // LD C,A
		cpu.c = cpu.a

	case 0x48: // LD C,B
		cpu.c = cpu.b

	case 0x49: // LD C,C

	case 0x4a: // LD C,D
		cpu.c = cpu.d

	case 0x4b: // LD C,E
		cpu.c = cpu.e

	case 0x4c: // LD C,H
		cpu.c = cpu.h

	case 0x4d: // LD C,L
		cpu.c = cpu.l

	case 0x57: // LD D,A
		cpu.d = cpu.a

	case 0x50: // LD D,B
		cpu.d = cpu.b

	case 0x51: // LD D,C
		cpu.d = cpu.c

	case 0x52: // LD D,D

	case 0x53: // LD D,E
		cpu.d = cpu.e

	case 0x55: // LD D,L
		cpu.d = cpu.l

	case 0x54: // LD D,H
		cpu.d = cpu.h

	case 0x58: // LD E,B
		cpu.e = cpu.b

	case 0x59: // LD E,C
		cpu.e = cpu.c

	case 0x5a: // LD E,D
		cpu.e = cpu.d

	case 0x5b: // LD E,E

	case 0x5c: // LD E,H
		cpu.e = cpu.h

	case 0x5d: // LD E,L
		cpu.e = cpu.l

	case 0x1e: // LD E,N
		cpu.e = ins.Mem[1]

	case 0x60: // LD H,B
		cpu.h = cpu.b

	case 0x61: // LD H,C
		cpu.h = cpu.c

	case 0x63: // LD H,E
		cpu.h = cpu.e

	case 0x26: // LD H,N
		cpu.h = ins.Mem[1]

	case 0x68: // LD L,B
		cpu.l = cpu.b

	case 0x64: // LD H,H

	case 0x65: // LD H,L
		cpu.h = cpu.l

	case 0x66: // LD H,(HL)
		cpu.h = cpu.memory.GetByte(getRR(cpu.h, cpu.l))

	case 0x69: // LD L,C
		cpu.l = cpu.c

	case 0x6a: // LD L,D
		cpu.l = cpu.d

	case 0x6c: // LD L,H
		cpu.l = cpu.h

	case 0x6d: // LD L,L

	case 0x2e: // LD L,N
		cpu.l = ins.Mem[1]

	case 0x6e: // LD L,(HL)
		cpu.l = cpu.memory.GetByte(getRR(cpu.h, cpu.l))

	case 0xed4f: // LD R,A
		cpu.r = cpu.a

	case 0x01: // LD BC,nn
		cpu.b = ins.Mem[2]
		cpu.c = ins.Mem[1]

	case 0x0e: // LD C,N
		cpu.c = ins.Mem[1]

	case 0x06: // LD B,N
		cpu.b = ins.Mem[1]

	case 0x16: // LD D,N
		// TODO one case fot LD r,n
		cpu.d = ins.Mem[1]

	case 0x5f: // LD E,A
		cpu.e = cpu.a

	case 0x11: // LD DE,nn
		cpu.d = ins.Mem[2]
		cpu.e = ins.Mem[1]

	case 0xed5b: // LD DE,(NN)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		cpu.d = cpu.memory.GetByte(nn + 1)
		cpu.e = cpu.memory.GetByte(nn)

	case 0x0a: // LD A,(BC)
		bc := getRR(cpu.b, cpu.c)
		cpu.a = cpu.memory.GetByte(bc)

	case 0x1a: // LD A,(DE)
		de := getRR(cpu.d, cpu.e)
		cpu.a = cpu.memory.GetByte(de)

	case 0x7e: // LD A,(HL)
		hl := getRR(cpu.h, cpu.l)
		cpu.a = cpu.memory.GetByte(hl)

	case 0x4e: // LD C,(HL)
		// TODO join all LD r,(HL)
		hl := getRR(cpu.h, cpu.l)
		cpu.c = cpu.memory.GetByte(hl)

	case 0x56: // LD D,(HL)
		hl := getRR(cpu.h, cpu.l)
		cpu.d = cpu.memory.GetByte(hl)

	case 0x5e: // LD E,(HL)
		hl := getRR(cpu.h, cpu.l)
		cpu.e = cpu.memory.GetByte(hl)

	case 0x21: // LD HL,nn
		cpu.h = ins.Mem[2]
		cpu.l = ins.Mem[1]

	case 0x3a: // LD A,(NN)
		nn := toWord(ins.Mem[1], ins.Mem[2])
		cpu.a = cpu.memory.GetByte(nn)

	case 0x2a: // LD HL,(nn)
		nn := toWord(ins.Mem[1], ins.Mem[2])
		cpu.l = cpu.memory.GetByte(nn)
		cpu.h = cpu.memory.GetByte(nn + 1)

	case 0x70: // LD (HL),B
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), cpu.b)

	case 0x71: // LD (HL),C
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), cpu.c)

	case 0x75: // LD (HL),L
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), cpu.l)

	case 0x76: // HALT
		if cpu.haltDone {
			cpu.haltDone = false
		} else {
			cpu.halt = true
			needPcUpdate = false
		}

	case 0xdd77: // LD (IX+N),A
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.a)

	case 0xdd70: // LD (IX+N),B
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.b)

	case 0xdd71: // LD (IX+N),C
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.c)

	case 0xdd72: // LD (IX+N),D
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.d)

	case 0xdd73: // LD (IX+N),E
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.e)

	case 0xdd74: // LD (IX+N),H
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.h)

	case 0xdd75: // LD (IX+N),L
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.l)

	case 0xfd77: // LD (IY+N),A
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.a)

	case 0xfd70: // LD (IY+N),B
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.b)

	case 0xfd71: // LD (IY+N),C
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.c)

	case 0xfd72: // LD (IY+N),D
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.d)

	case 0xfd73: // LD (IY+N),E
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.e)

	case 0xfd74: // LD (IY+N),H
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.h)

	case 0xfd75: // LD (IY+N),L
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.l)

	case 0xfd36: // LD (IY+N),N
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, ins.Mem[3])

	case 0xdd36: // LD (IX+N),N
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, ins.Mem[3])

	case 0xdd7e: // LD A,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.a = cpu.memory.GetByte(ix)

	case 0xdd46: // LD B,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.b = cpu.memory.GetByte(ix)

	case 0xdd4e: // LD C,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.c = cpu.memory.GetByte(ix)

	case 0xdd56: // LD D,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.d = cpu.memory.GetByte(ix)

	case 0xdd5e: // LD E,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.e = cpu.memory.GetByte(ix)

	case 0xdd66: // LD H,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.h = cpu.memory.GetByte(ix)

	case 0xfd66: // LD H,(IY+N)
		IY := cpu.getIYn(ins.Mem[2])
		cpu.h = cpu.memory.GetByte(IY)

	case 0xdd6e: // LD L,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.l = cpu.memory.GetByte(ix)

	case 0xfd7e: // LD A,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.a = cpu.memory.GetByte(iy)

	case 0xfd4e: // LD C,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.c = cpu.memory.GetByte(iy)

	case 0xfd56: // LD D,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.d = cpu.memory.GetByte(iy)

	case 0xfd6e: // LD L,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.l = cpu.memory.GetByte(iy)

	case 0xfd46: // LD B,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.b = cpu.memory.GetByte(iy)

	case 0xfd5e: // LD E,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.e = cpu.memory.GetByte(iy)

	case 0xca: // JP Z,$NN
		if cpu.f.Z {
			cpu.pc = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xc2: // JP NZ,$NN
		if cpu.f.Z == false {
			cpu.pc = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xf2: // JP P,$NN
		if !cpu.f.S {
			cpu.pc = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xfa: // JP M,$NN
		if cpu.f.S {
			cpu.pc = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xea: // JP PE,$NN
		if cpu.f.P {
			cpu.pc = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xe2: // JP PO,$NN
		if !cpu.f.P == true {
			cpu.pc = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xd2: // JP NC,$NN
		if !cpu.f.C {
			cpu.pc = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xda: // JP C,$NN
		if cpu.f.C {
			cpu.pc = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xC3: // JP nn
		cpu.pc = toWord(ins.Mem[1], ins.Mem[2])
		needPcUpdate = false

	case 0xdde9: // JP (IX)
		cpu.pc = getRR(cpu.ixh, cpu.ixl)
		needPcUpdate = false

	case 0xfde9: // JP (IY)
		cpu.pc = getRR(cpu.iyh, cpu.iyl)
		needPcUpdate = false

	case 0xe9: // JP (HL)
		hl := getRR(cpu.h, cpu.l)
		cpu.pc = hl
		needPcUpdate = false

	case 0xed47: // LD I,A
		cpu.i = cpu.a

	case 0x67: // LD H,A
		cpu.h = cpu.a

	case 0x62: // LD H,D
		cpu.h = cpu.d

	case 0x6f: // LD L,A
		cpu.l = cpu.a

	case 0x6b: // LD L,E
		cpu.l = cpu.e

	case 0x36: // LD (HL),n
		hl := getRR(cpu.h, cpu.l)
		v := ins.Mem[1]
		cpu.memory.PutByte(hl, v)

	case 0x77: // LD (HL),A
		hl := getRR(cpu.h, cpu.l)
		cpu.memory.PutByte(hl, cpu.a)

	case 0x72: // LD (HL),D
		hl := getRR(cpu.h, cpu.l)
		cpu.memory.PutByte(hl, cpu.d)

	case 0x73: // LD (HL),E
		hl := getRR(cpu.h, cpu.l)
		cpu.memory.PutByte(hl, cpu.e)

	case 0x74: // LD (HL),H
		hl := getRR(cpu.h, cpu.l)
		cpu.memory.PutByte(hl, cpu.h)

	case 0x32: // LD (nn),A
		nn := toWord(ins.Mem[1], ins.Mem[2])
		cpu.memory.PutByte(nn, cpu.a)

	case 0xed43: // LD (nn),BC
		w := getRR(cpu.b, cpu.c)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		putWord(cpu.memory, nn, w)

	case 0xed53: // LD (nn),DE
		w := getRR(cpu.d, cpu.e)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		putWord(cpu.memory, nn, w)

	case 0x22: // LD (nn),HL
		w := getRR(cpu.h, cpu.l)
		nn := toWord(ins.Mem[1], ins.Mem[2])
		putWord(cpu.memory, nn, w)

	case 0xdd22: // LD (NN),IX
		w := getRR(cpu.ixh, cpu.ixl)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		putWord(cpu.memory, nn, w)

	case 0xfd22: // LD (NN),IY
		w := getRR(cpu.iyh, cpu.iyl)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		putWord(cpu.memory, nn, w)

	case 0xed73: // LD (NN),SP
		nn := toWord(ins.Mem[2], ins.Mem[3])
		putWord(cpu.memory, nn, cpu.sp.Get())

	case 0xdd21: // LD IX,NN
		cpu.ixh = ins.Mem[3]
		cpu.ixl = ins.Mem[2]

	case 0xdd26: // LD IXH,N
		cpu.ixh = ins.Mem[2]

	case 0xdd2e: // LD IXL,N
		cpu.ixl = ins.Mem[2]

	case 0xfd26: // LD IYH,N
		cpu.iyh = ins.Mem[2]

	case 0xfd2e: // LD IYL,N
		cpu.iyl = ins.Mem[2]

	case 0xed4b: // LD BC,(NN)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		cpu.b = cpu.memory.GetByte(nn + 1)
		cpu.c = cpu.memory.GetByte(nn)

	case 0xdd2a: // LD IX,(NN)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		cpu.ixh = cpu.memory.GetByte(nn + 1)
		cpu.ixl = cpu.memory.GetByte(nn)

	case 0xfd2a: // LD IY,(NN)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		cpu.iyh = cpu.memory.GetByte(nn + 1)
		cpu.iyl = cpu.memory.GetByte(nn)

	case 0xed7b: // LD SP,(NN)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		sp := uint16(cpu.memory.GetByte(nn+1))<<8 | uint16(cpu.memory.GetByte(nn))
		cpu.sp.Set(sp)

	case 0xfd21: // LD IY,nn
		cpu.iyh = ins.Mem[3]
		cpu.iyl = ins.Mem[2]

	case 0xf9: // LD SP,HL
		cpu.sp.Set(getRR(cpu.h, cpu.l))

	case 0xddf9: // LD SP,IX
		cpu.sp.Set(getRR(cpu.ixh, cpu.ixl))

	case 0xfdf9: // LD SP,IY
		cpu.sp.Set(getRR(cpu.iyh, cpu.iyl))

	case 0x31: // LD SP,NN
		cpu.sp.Set(toWord(ins.Mem[1], ins.Mem[2]))

	case 0xeda0: // LDI
		cpu.ldi()

	case 0xedb0: // LDIR
		cpu.ldi()
		bc := getRR(cpu.b, cpu.c)
		if bc != 0 {
			needPcUpdate = false
		}

	case 0xeda8: // LDD
		cpu.ldd()

	case 0xedb8: // LDDR
		cpu.ldd()
		bc := getRR(cpu.b, cpu.c)
		if bc != 0 {
			needPcUpdate = false
		}

	case 0xeda9: // CPD
		cpu.cpd()

	case 0x2f: // CPL
		cpu.a = ^cpu.a
		cpu.f.H = true
		cpu.f.N = true

	case 0xed44: // NEG
		n := cpu.a
		cpu.a = 0
		cpu.subA(n)

	case 0x3d: // DEC A
		cpu.decR(&cpu.a)

	case 0x05: // DEC B
		cpu.decR(&cpu.b)

	case 0x0d: // DEC C
		cpu.decR(&cpu.c)

	case 0x15: // DEC D
		cpu.decR(&cpu.d)

	case 0x1d: // DEC E
		cpu.decR(&cpu.e)

	case 0x25: // DEC H
		cpu.decR(&cpu.h)

	case 0x2d: // DEC L
		cpu.decR(&cpu.l)

	case 0x2b: // DEC HL
		// TODO join all DEC rr
		hl := getRR(cpu.h, cpu.l)
		hl--
		cpu.h, cpu.l = setRR(hl)

	case 0x0b: // DEC BC
		bc := getRR(cpu.b, cpu.c)
		bc--
		cpu.b, cpu.c = setRR(bc)

	case 0x1b: // DEC DE
		de := getRR(cpu.d, cpu.e)
		de--
		cpu.d, cpu.e = setRR(de)

	case 0xdd2b: // DEC IX
		ix := getRR(cpu.ixh, cpu.ixl)
		ix--
		cpu.ixh, cpu.ixl = setRR(ix)

	case 0xfd2b: // DEC IY
		iy := getRR(cpu.iyh, cpu.iyl)
		iy--
		cpu.iyh, cpu.iyl = setRR(iy)

	case 0x3c: // INC A
		cpu.incR(&cpu.a)

	case 0x04: // INC B
		cpu.incR(&cpu.b)

	case 0x0c: // INC C
		cpu.incR(&cpu.c)

	case 0x14: // INC D
		cpu.incR(&cpu.d)

	case 0x1c: // INC E
		cpu.incR(&cpu.e)

	case 0x24: // INC H
		cpu.incR(&cpu.h)

	case 0x2c: // INC L
		cpu.incR(&cpu.l)

	case 0x03: // INC BC
		bc := getRR(cpu.b, cpu.c)
		bc++
		cpu.b, cpu.c = setRR(bc)

	case 0x13: // INC DE
		de := getRR(cpu.d, cpu.e)
		de++
		cpu.d, cpu.e = setRR(de)

	case 0x23: // INC HL
		hl := getRR(cpu.h, cpu.l)
		hl++
		cpu.h, cpu.l = setRR(hl)

	case 0x33: // INC SP
		sp := cpu.sp.Get()
		cpu.sp.Set(sp + 1)

	case 0xdd23: // INC IX
		ix := getRR(cpu.ixh, cpu.ixl)
		ix++
		cpu.ixh, cpu.ixl = setRR(ix)

	case 0xfd23: // INC IY
		iy := getRR(cpu.iyh, cpu.iyl)
		iy++
		cpu.iyh, cpu.iyl = setRR(iy)

	case 0xdd34: // INC (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		ixv := cpu.memory.GetByte(ix)
		cpu.incR(&ixv)
		cpu.memory.PutByte(ix, ixv)

	case 0xfd34: // INC (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		iyv := cpu.memory.GetByte(iy)
		cpu.incR(&iyv)
		cpu.memory.PutByte(iy, iyv)

	case 0xdd35: // DEC (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		ixv := cpu.memory.GetByte(ix)
		cpu.decR(&ixv)
		cpu.memory.PutByte(ix, ixv)

	case 0xfd35: // DEC (IY+d)
		iy := cpu.getIYn(ins.Mem[2])
		iyv := cpu.memory.GetByte(iy)
		cpu.decR(&iyv)
		cpu.memory.PutByte(iy, iyv)

	case 0x3b: // DEC SP
		sp := cpu.sp.Get()
		cpu.sp.Set(sp - 1)

	case 0x34: // INC (HL)
		hl := getRR(cpu.h, cpu.l)
		b := cpu.memory.GetByte(hl)
		cpu.incR(&b)
		cpu.memory.PutByte(hl, b)

	case 0x35: // DEC (HL)
		hl := getRR(cpu.h, cpu.l)
		b := cpu.memory.GetByte(hl)
		cpu.decR(&b)
		cpu.memory.PutByte(hl, b)

	case 0xbf: // CP A
		cpu.cp(cpu.a)

	case 0xb8: // CP B
		cpu.cp(cpu.b)

	case 0xb9: // CP C
		cpu.cp(cpu.c)

	case 0xba: // CP D
		cpu.cp(cpu.d)

	case 0xbb: // CP E
		cpu.cp(cpu.e)

	case 0xbc: // CP H
		cpu.cp(cpu.h)

	case 0xbd: // CP L
		cpu.cp(cpu.l)

	case 0xfe: // CP N
		n := ins.Mem[1]
		cpu.cp(n)

	case 0xbe: // CP (HL)
		cpu.cp(cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xeda1: // CPI
		cpu.cpi()

	case 0xedb1: // CPIR
		diff := cpu.cpi()
		bc := getRR(cpu.b, cpu.c)
		if (bc != 0) && (diff != 0) {
			needPcUpdate = false
		}

	case 0xedb9: // CPDR
		diff := cpu.cpd()
		bc := getRR(cpu.b, cpu.c)
		if (bc == 0) || (diff == 0) {
			cpu.f.P = bc != 0
		} else {
			needPcUpdate = false
		}

	case 0x20: // JR NZ,d
		if cpu.f.Z == false {
			cpu.jr(ins.Mem[1])
			needPcUpdate = false
		}

	case 0x28: // JR Z,d
		if cpu.f.Z {
			cpu.jr(ins.Mem[1])
			needPcUpdate = false
		}

	case 0x30: // JR NC,d
		if cpu.f.C != true {
			cpu.jr(ins.Mem[1])
			needPcUpdate = false
		}

	case 0x38: // JR C,$N+2
		if cpu.f.C == true {
			cpu.jr(ins.Mem[1])
			needPcUpdate = false
		}

	case 0x37: // SCF
		cpu.f.H = false
		cpu.f.N = false
		cpu.f.C = true

	case 0x18: // JR $N+2
		cpu.jr(ins.Mem[1])
		needPcUpdate = false

	case 0x10: // DJNZ $+2
		cpu.b--
		if cpu.b != 0 {
			cpu.jr(ins.Mem[1])
			cpu.clock.AddTStates(5)
			needPcUpdate = false
		}

	case 0xa7: // AND A
		cpu.and(cpu.a)

	case 0xa0: // AND B
		cpu.and(cpu.b)

	case 0xa1: // AND C
		cpu.and(cpu.c)

	case 0xa2: // AND D
		cpu.and(cpu.d)

	case 0xa3: // AND E
		cpu.and(cpu.e)

	case 0xa4: // AND H
		cpu.and(cpu.h)

	case 0xa5: // AND L
		cpu.and(cpu.l)

	case 0xe6: // AND N
		cpu.and(ins.Mem[1])

	case 0xa6: // AND (HL)
		cpu.and(cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xb7: // OR A
		cpu.or(cpu.a)

	case 0xb0: // OR B
		cpu.or(cpu.b)

	case 0xb1: // OR C
		cpu.or(cpu.c)

	case 0xb2: // OR D
		cpu.or(cpu.d)

	case 0xb3: // OR E
		cpu.or(cpu.e)

	case 0xb4: // OR H
		cpu.or(cpu.h)

	case 0xb5: // OR L
		cpu.or(cpu.l)

	case 0xb6: // OR (HL)
		cpu.or(cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xf6: // OR N
		cpu.or(ins.Mem[1])

	case 0xde: // SBC A,N
		cpu.sbc(ins.Mem[1])

	case 0xed52: // SBC HL,DE
		cpu.sbcHL(getRR(cpu.d, cpu.e))

	case 0xed42: // SBC HL,BC
		cpu.sbcHL(getRR(cpu.b, cpu.c))

	case 0xed62: // SBC HL,HL
		cpu.sbcHL(getRR(cpu.h, cpu.l))

	case 0xed72: // SBC HL,SP
		cpu.sbcHL(cpu.sp.Get())

	case 0x09: // ADD HL,BC
		cpu.addHL(getRR(cpu.b, cpu.c))

	case 0x19: // ADD HL,DE
		cpu.addHL(getRR(cpu.d, cpu.e))

	case 0x29: // ADD HL,HL
		cpu.addHL(getRR(cpu.h, cpu.l))

	case 0x39: // ADD HL,SP
		cpu.addHL(cpu.sp.Get())

	case 0xdd09: // ADD IX,BC
		cpu.addIX(getRR(cpu.b, cpu.c))

	case 0xdd19: // ADD IX,DE
		cpu.addIX(getRR(cpu.d, cpu.e))

	case 0xdd29: // ADD IX,IX
		cpu.addIX(getRR(cpu.ixh, cpu.ixl))

	case 0xdd39: // ADD IX,SP
		cpu.addIX(cpu.sp.Get())

	case 0xfd09: // ADD IY,BC
		cpu.addIY(getRR(cpu.b, cpu.c))

	case 0xfd19: // ADD IY,DE
		cpu.addIY(getRR(cpu.d, cpu.e))

	case 0xfd29: // ADD IY,IY
		cpu.addIY(getRR(cpu.iyh, cpu.iyl))

	case 0xfd39: // ADD IY,SP
		cpu.addIY(cpu.sp.Get())

	case 0x87: // ADD A,A
		cpu.addA(cpu.a)

	case 0x80: // ADD A,B
		cpu.addA(cpu.b)

	case 0x81: // ADD A,C
		cpu.addA(cpu.c)

	case 0x82: // ADD A,D
		cpu.addA(cpu.d)

	case 0x83: // ADD A,E
		cpu.addA(cpu.e)

	case 0x84: // ADD A,H
		cpu.addA(cpu.h)

	case 0x85: // ADD A,L
		cpu.addA(cpu.l)

	case 0xc6: // ADD A,N
		cpu.addA(ins.Mem[1])

	case 0x86: // ADD A,(HL)
		cpu.addA(cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xfd86: // ADD A,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.addA(cpu.memory.GetByte(iy))

	case 0x8f: // ADC A,A
		cpu.adc(cpu.a)

	case 0x88: // ADC A,B
		cpu.adc(cpu.b)

	case 0x89: // ADC A,C
		cpu.adc(cpu.c)

	case 0x8a: // ADC A,D
		cpu.adc(cpu.d)

	case 0x8b: // ADC A,E
		cpu.adc(cpu.e)

	case 0x8c: // ADC A,H
		cpu.adc(cpu.h)

	case 0x8d: // ADC A,L
		cpu.adc(cpu.l)

	case 0x8e: // ADC A,(HL)
		cpu.adc(cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xce: // ADC A,N
		cpu.adc(ins.Mem[1])

	case 0xed4a: // ADC HL,BC
		cpu.adcHL(getRR(cpu.b, cpu.c))

	case 0xed5a: // ADC HL,DE
		cpu.adcHL(getRR(cpu.d, cpu.e))

	case 0xed6a: // ADC HL,HL
		cpu.adcHL(getRR(cpu.h, cpu.l))

	case 0xed7a: // ADC HL,SP
		cpu.adcHL(cpu.sp.Get())

	case 0x97: // SUB A
		cpu.subA(cpu.a)

	case 0x90: // SUB B
		cpu.subA(cpu.b)

	case 0x91: // SUB C
		cpu.subA(cpu.c)

	case 0x92: // SUB D
		cpu.subA(cpu.d)

	case 0x93: // SUB E
		cpu.subA(cpu.e)

	case 0x94: // SUB H
		cpu.subA(cpu.h)

	case 0x95: // SUB L
		cpu.subA(cpu.l)

	case 0x96: // SUB (HL)
		cpu.subA(cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xd6: // SUB N
		cpu.subA(ins.Mem[1])

	case 0x9f: // SBC A
		cpu.sbc(cpu.a)

	case 0x98: // SBC B
		cpu.sbc(cpu.b)

	case 0x99: // SBC C
		cpu.sbc(cpu.c)

	case 0x9a: // SBC D
		cpu.sbc(cpu.d)

	case 0x9b: // SBC E
		cpu.sbc(cpu.e)

	case 0x9c: // SBC H
		cpu.sbc(cpu.h)

	case 0x9d: // SBC L
		cpu.sbc(cpu.l)

	case 0x9e: // SBC (HL)
		cpu.sbc(cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0x08: // EX AF,AF'
		ta := cpu.a
		cpu.a = cpu._a
		cpu._a = ta

		tf := cpu.f
		cpu.f = cpu._f
		cpu._f = tf

	case 0xeb: // EX DE,HL
		td := cpu.d
		cpu.d = cpu.h
		cpu.h = td

		te := cpu.e
		cpu.e = cpu.l
		cpu.l = te

	case 0xdde3: // EX (SP),IX
		spv := getWord(cpu.memory, cpu.sp.Get())
		ix := getRR(cpu.ixh, cpu.ixl)
		putWord(cpu.memory, cpu.sp.Get(), ix)
		cpu.ixh, cpu.ixl = setRR(spv)

	case 0xfde3: // EX (SP),IY
		spv := getWord(cpu.memory, cpu.sp.Get())
		iy := getRR(cpu.iyh, cpu.iyl)
		putWord(cpu.memory, cpu.sp.Get(), iy)
		cpu.iyh, cpu.iyl = setRR(spv)

	case 0xe3: // EX (SP),HL
		spv := getWord(cpu.memory, cpu.sp.Get())
		hl := getRR(cpu.h, cpu.l)
		putWord(cpu.memory, cpu.sp.Get(), hl)
		cpu.h, cpu.l = setRR(spv)

	case 0xd9: // EXX
		tb := cpu.b
		cpu.b = cpu._b
		cpu._b = tb

		tc := cpu.c
		cpu.c = cpu._c
		cpu._c = tc

		td := cpu.d
		cpu.d = cpu._d
		cpu._d = td

		te := cpu.e
		cpu.e = cpu._e
		cpu._e = te

		th := cpu.h
		cpu.h = cpu._h
		cpu._h = th

		tl := cpu.l
		cpu.l = cpu._l
		cpu._l = tl

	case 0xcb37: // SLL A
		cpu.sll(&cpu.a)

	case 0xcb30: // SLL B
		cpu.sll(&cpu.b)

	case 0xcb31: // SLL C
		cpu.sll(&cpu.c)

	case 0xcb32: // SLL D
		cpu.sll(&cpu.d)

	case 0xcb33: // SLL E
		cpu.sll(&cpu.e)

	case 0xcb34: // SLL H
		cpu.sll(&cpu.h)

	case 0xcb35: // SLL L
		cpu.sll(&cpu.l)

	case 0xcb36: // SLL (HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.sll(&b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcb3f: // SRL A
		cpu.srl(&cpu.a)

	case 0xcb38: // SRL B
		cpu.srl(&cpu.b)

	case 0xcb39: // SRL C
		cpu.srl(&cpu.c)

	case 0xcb3a: // SRL D
		cpu.srl(&cpu.d)

	case 0xcb3b: // SRL E
		cpu.srl(&cpu.e)

	case 0xcb3c: // SRL H
		cpu.srl(&cpu.h)

	case 0xcb3d: // SRL L
		cpu.srl(&cpu.l)

	case 0xcb3e: // SRL (HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.srl(&b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xeda3: // OUTI
		hl := getRR(cpu.h, cpu.l)
		b := cpu.memory.GetByte(hl)
		cpu.b--
		cpu.writePort(getRR(cpu.b, cpu.c), b)
		hl++
		cpu.h, cpu.l = setRR(hl)
		cpu.f.Z = cpu.b == 0
		cpu.f.S = cpu.b&0x80 != 0
		cpu.f.N = cpu.b&0x80 == 0
		cpu.f.H = true
		cpu.f.P = parityTable[cpu.b]

	case 0xedab: // OUTD
		hl := getRR(cpu.h, cpu.l)
		b := cpu.memory.GetByte(hl)
		cpu.b--
		cpu.writePort(getRR(cpu.b, cpu.c), b)
		hl--
		cpu.h, cpu.l = setRR(hl)
		cpu.f.Z = cpu.b == 0
		cpu.f.S = cpu.b&0x80 != 0
		cpu.f.N = cpu.b&0x80 == 0
		cpu.f.H = true
		cpu.f.P = parityTable[cpu.b]

	case 0xd3: // OUT (n),A
		port := toWord(ins.Mem[1], cpu.a)
		cpu.writePort(port, cpu.a)

	case 0xed79: // OUT (C),A
		cpu.writePort(getRR(cpu.b, cpu.c), cpu.a)

	case 0xed41: // OUT (C),B
		cpu.writePort(getRR(cpu.b, cpu.c), cpu.b)

	case 0xed49: // OUT (C),C
		cpu.writePort(getRR(cpu.b, cpu.c), cpu.c)

	case 0xed51: // OUT (C),D
		cpu.writePort(getRR(cpu.b, cpu.c), cpu.d)

	case 0xed59: // OUT (C),E
		cpu.writePort(getRR(cpu.b, cpu.c), cpu.e)

	case 0xed61: // OUT (C),H
		cpu.writePort(getRR(cpu.b, cpu.c), cpu.h)

	case 0xed69: // OUT (C),L
		cpu.writePort(getRR(cpu.b, cpu.c), cpu.l)

	case 0xed71: // OUT (C),0
		cpu.writePort(getRR(cpu.b, cpu.c), 0)

	case 0xed78: // IN A,(C)
		cpu.a = cpu.readPort(getRR(cpu.b, cpu.c))

	case 0xed40: // IN B,(C)
		cpu.b = cpu.readPort(getRR(cpu.b, cpu.c))

	case 0xed48: // IN C,(C)
		cpu.c = cpu.readPort(getRR(cpu.b, cpu.c))

	case 0xed50: // IN D,(C)
		cpu.d = cpu.readPort(getRR(cpu.b, cpu.c))

	case 0xed58: // IN E,(C)
		cpu.e = cpu.readPort(getRR(cpu.b, cpu.c))

	case 0xed60: // IN H,(C)
		cpu.h = cpu.readPort(getRR(cpu.b, cpu.c))

	case 0xed68: // IN L,(C)
		cpu.l = cpu.readPort(getRR(cpu.b, cpu.c))

	case 0xed70: // IN (C)
		cpu.readPort(getRR(cpu.b, cpu.c))

	case 0xeda2: // INI
		hl := getRR(cpu.h, cpu.l)
		v := cpu.readPort(getRR(cpu.b, cpu.c))
		cpu.memory.PutByte(hl, v)
		cpu.b--
		hl++
		cpu.f.N = v&0x80 != 0
		cpu.f.Z = v == 0
		// cpu.f.H = v == 0
		// cpu.f.C = false
		// cpu.f.P = false
		cpu.h, cpu.l = setRR(hl)

	case 0xdb: // IN A,(N)
		f := cpu.f.getByte()
		port := toWord(ins.Mem[1], cpu.a)
		cpu.a = cpu.readPort(port)
		cpu.f.setByte(f)

	case 0xedb3: // OTIR
		hl := getRR(cpu.h, cpu.l)
		c := cpu.memory.GetByte(hl)
		cpu.readPort(uint16(c) << 8)
		cpu.b--
		hl++
		cpu.h, cpu.l = setRR(hl)
		if cpu.b != 0 {
			cpu.pc -= 2
		} else {
			cpu.f.Z = true
			cpu.f.N = true
		}

	case 0xfdcbc6, 0xfdcbce, 0xfdcbee, 0xfdcbde, 0xfdcbd6, 0xfdcbe6, 0xfdcbfe, 0xfdcbf6: // SET b,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		b := (ins.Mem[3] & 0b00111000) >> 3
		b = 1 << b
		data := cpu.memory.GetByte(iy)
		data |= b
		cpu.memory.PutByte(iy, data)

	case 0xfdcb86, 0xfdcb8e, 0xfdcb96, 0xfdcba6, 0xfdcbae, 0xfdcb9e, 0xfdcbbe, 0xfdcbb6: // RES b,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		b := (ins.Mem[3] & 0b00111000) >> 3
		b = ^(1 << b)
		data := cpu.memory.GetByte(iy)
		data &= b
		cpu.memory.PutByte(iy, data)

	case 0xfdcb76, 0xfdcb4e, 0xfdcb46, 0xfdcb6e, 0xfdcb5e, 0xfdcb66, 0xfdcb7e, 0xfdcb56: // BIT b,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		b := (ins.Mem[3] & 0b00111000) >> 3
		cpu.bit(b, cpu.memory.GetByte(iy))

	case 0xcb47: // BIT 0,A
		cpu.bit(0, cpu.a)

	case 0xcb40: // BIT 0,B
		cpu.bit(0, cpu.b)

	case 0xcb41: // BIT 0,C
		cpu.bit(0, cpu.c)

	case 0xcb42: // BIT 0,D
		cpu.bit(0, cpu.d)

	case 0xcb43: // BIT 0,E
		cpu.bit(0, cpu.e)

	case 0xcb44: // BIT 0,H
		cpu.bit(0, cpu.h)

	case 0xcb45: // BIT 0,L
		cpu.bit(0, cpu.l)

	case 0xcb46: // BIT 0,(HL)
		cpu.bit(0, cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xcb4f: // BIT 1,A
		cpu.bit(1, cpu.a)

	case 0xcb48: // BIT 1,B
		cpu.bit(1, cpu.b)

	case 0xcb49: // BIT 1,C
		cpu.bit(1, cpu.c)

	case 0xcb4a: // BIT 1,D
		cpu.bit(1, cpu.d)

	case 0xcb4b: // BIT 1,E
		cpu.bit(1, cpu.e)

	case 0xcb4c: // BIT 1,H
		cpu.bit(1, cpu.h)

	case 0xcb4d: // BIT 1,L
		cpu.bit(1, cpu.l)

	case 0xcb4e: // BIT 1,(HL)
		cpu.bit(1, cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xcb57: // BIT 2,A
		cpu.bit(2, cpu.a)

	case 0xcb50: // BIT 2,B
		cpu.bit(2, cpu.b)

	case 0xcb51: // BIT 2,C
		cpu.bit(2, cpu.c)

	case 0xcb52: // BIT 2,D
		cpu.bit(2, cpu.d)

	case 0xcb53: // BIT 2,E
		cpu.bit(2, cpu.e)

	case 0xcb54: // BIT 2,H
		cpu.bit(2, cpu.h)

	case 0xcb55: // BIT 2,L
		cpu.bit(2, cpu.l)

	case 0xcb56: // BIT 2,(HL)
		cpu.bit(2, cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xcb5f: // BIT 3,A
		cpu.bit(3, cpu.a)

	case 0xcb58: // BIT 3,B
		cpu.bit(3, cpu.b)

	case 0xcb59: // BIT 3,C
		cpu.bit(3, cpu.c)

	case 0xcb5a: // BIT 3,D
		cpu.bit(3, cpu.d)

	case 0xcb5b: // BIT 3,E
		cpu.bit(3, cpu.e)

	case 0xcb5c: // BIT 3,H
		cpu.bit(3, cpu.h)

	case 0xcb5d: // BIT 3,L
		cpu.bit(3, cpu.l)

	case 0xcb5e: // BIT 3,(HL)
		cpu.bit(3, cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xcb67: // BIT 4,A
		cpu.bit(4, cpu.a)

	case 0xcb60: // BIT 4,B
		cpu.bit(4, cpu.b)

	case 0xcb61: // BIT 4,C
		cpu.bit(4, cpu.c)

	case 0xcb62: // BIT 4,D
		cpu.bit(4, cpu.d)

	case 0xcb63: // BIT 4,E
		cpu.bit(4, cpu.e)

	case 0xcb64: // BIT 4,H
		cpu.bit(4, cpu.h)

	case 0xcb65: // BIT 4,L
		cpu.bit(4, cpu.l)

	case 0xcb66: // BIT 4,(HL)
		cpu.bit(4, cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xcb6f: // BIT 5,A
		cpu.bit(5, cpu.a)

	case 0xcb68: // BIT 5,B
		cpu.bit(5, cpu.b)

	case 0xcb69: // BIT 5,C
		cpu.bit(5, cpu.c)

	case 0xcb6a: // BIT 5,D
		cpu.bit(5, cpu.d)

	case 0xcb6b: // BIT 5,E
		cpu.bit(5, cpu.e)

	case 0xcb6c: // BIT 5,H
		cpu.bit(5, cpu.h)

	case 0xcb6d: // BIT 5,L
		cpu.bit(5, cpu.l)

	case 0xcb6e: // BIT 5,(HL)
		cpu.bit(5, cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xcb77: // BIT 6,A
		cpu.bit(6, cpu.a)

	case 0xcb70: // BIT 6,B
		cpu.bit(6, cpu.b)

	case 0xcb71: // BIT 6,C
		cpu.bit(6, cpu.c)

	case 0xcb72: // BIT 6,D
		cpu.bit(6, cpu.d)

	case 0xcb73: // BIT 6,E
		cpu.bit(6, cpu.e)

	case 0xcb74: // BIT 6,H
		cpu.bit(6, cpu.h)

	case 0xcb75: // BIT 6,L
		cpu.bit(6, cpu.l)

	case 0xcb76: // BIT 6,(HL)
		cpu.bit(6, cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xcb7f: // BIT 7,A
		cpu.bit(7, cpu.a)

	case 0xcb78: // BIT 7,B
		cpu.bit(7, cpu.b)

	case 0xcb79: // BIT 7,C
		cpu.bit(7, cpu.c)

	case 0xcb7a: // BIT 7,D
		cpu.bit(7, cpu.d)

	case 0xcb7b: // BIT 7,E
		cpu.bit(7, cpu.e)

	case 0xcb7c: // BIT 7,H
		cpu.bit(7, cpu.h)

	case 0xcb7d: // BIT 7,L
		cpu.bit(7, cpu.l)

	case 0xcb7e: // BIT 7,(HL)
		cpu.bit(7, cpu.memory.GetByte(getRR(cpu.h, cpu.l)))

	case 0xcb80: // RES 0,B
		cpu.res(0, &cpu.b)

	case 0xcb81: // RES 0,C
		cpu.res(0, &cpu.c)

	case 0xcb82: // RES 0,D
		cpu.res(0, &cpu.d)

	case 0xcb83: // RES 0,E
		cpu.res(0, &cpu.e)

	case 0xcb84: // RES 0,H
		cpu.res(0, &cpu.h)

	case 0xcb85: // RES 0,L
		cpu.res(0, &cpu.l)

	case 0xcb86: // RES 0,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.res(0, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcb87: // RES 0,A
		cpu.res(0, &cpu.a)

	case 0xcb88: // RES 1,B
		cpu.res(1, &cpu.b)

	case 0xcb89: // RES 1,C
		cpu.res(1, &cpu.c)

	case 0xcb8a: // RES 1,D
		cpu.res(1, &cpu.d)

	case 0xcb8b: // RES 1,E
		cpu.res(1, &cpu.e)

	case 0xcb8c: // RES 1,H
		cpu.res(1, &cpu.h)

	case 0xcb8d: // RES 1,L
		cpu.res(1, &cpu.l)

	case 0xcb8e: // RES 1,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.res(1, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcb8f: // RES 1,A
		cpu.res(1, &cpu.a)

	case 0xcb90: // RES 2,B
		cpu.res(2, &cpu.b)

	case 0xcb91: // RES 2,C
		cpu.res(2, &cpu.c)

	case 0xcb92: // RES 2,D
		cpu.res(2, &cpu.d)

	case 0xcb93: // RES 2,E
		cpu.res(2, &cpu.e)

	case 0xcb94: // RES 2,H
		cpu.res(2, &cpu.h)

	case 0xcb95: // RES 2,L
		cpu.res(2, &cpu.l)

	case 0xcb96: // RES 2,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.res(2, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcb97: // RES 2,A
		cpu.res(2, &cpu.a)

	case 0xcb98: // RES 3,B
		cpu.res(3, &cpu.b)

	case 0xcb99: // RES 3,C
		cpu.res(3, &cpu.c)

	case 0xcb9a: // RES 3,D
		cpu.res(3, &cpu.d)

	case 0xcb9b: // RES 3,E
		cpu.res(3, &cpu.e)

	case 0xcb9c: // RES 3,H
		cpu.res(3, &cpu.h)

	case 0xcb9d: // RES 3,L
		cpu.res(3, &cpu.l)

	case 0xcb9e: // RES 3,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.res(3, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcb9f: // RES 3,A
		cpu.res(3, &cpu.a)

	case 0xcba0: // RES 4,B
		cpu.res(4, &cpu.b)

	case 0xcba1: // RES 4,C
		cpu.res(4, &cpu.c)

	case 0xcba2: // RES 4,D
		cpu.res(4, &cpu.d)

	case 0xcba3: // RES 4,E
		cpu.res(4, &cpu.e)

	case 0xcba4: // RES 4,H
		cpu.res(4, &cpu.h)

	case 0xcba5: // RES 4,L
		cpu.res(4, &cpu.l)

	case 0xcba6: // RES 4,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.res(4, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcba7: // RES 4,A
		cpu.res(4, &cpu.a)

	case 0xcba8: // RES 5,B
		cpu.res(5, &cpu.b)

	case 0xcba9: // RES 5,C
		cpu.res(5, &cpu.c)

	case 0xcbaa: // RES 5,D
		cpu.res(5, &cpu.d)

	case 0xcbab: // RES 5,E
		cpu.res(5, &cpu.e)

	case 0xcbac: // RES 5,H
		cpu.res(5, &cpu.h)

	case 0xcbad: // RES 5,L
		cpu.res(5, &cpu.l)

	case 0xcbae: // RES 5,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.res(5, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcbaf: // RES 5,A
		cpu.res(5, &cpu.a)

	case 0xcbb0: // RES 6,B
		cpu.res(6, &cpu.b)

	case 0xcbb1: // RES 6,C
		cpu.res(6, &cpu.c)

	case 0xcbb2: // RES 6,D
		cpu.res(6, &cpu.d)

	case 0xcbb3: // RES 6,E
		cpu.res(6, &cpu.e)

	case 0xcbb4: // RES 6,H
		cpu.res(6, &cpu.h)

	case 0xcbb5: // RES 6,L
		cpu.res(6, &cpu.l)

	case 0xcbb6: // RES 6,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.res(6, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcbb7: // RES 6,A
		cpu.res(6, &cpu.a)

	case 0xcbb8: // RES 7,B
		cpu.res(7, &cpu.b)

	case 0xcbb9: // RES 7,C
		cpu.res(7, &cpu.c)

	case 0xcbba: // RES 7,D
		cpu.res(7, &cpu.d)

	case 0xcbbb: // RES 7,E
		cpu.res(7, &cpu.e)

	case 0xcbbc: // RES 7,H
		cpu.res(7, &cpu.h)

	case 0xcbbd: // RES 7,L
		cpu.res(7, &cpu.l)

	case 0xcbbe: // RES 7,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.res(7, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcbbf: // RES 7,A
		cpu.res(7, &cpu.a)

	case 0xcbc0: // SET 0,B
		cpu.set(0, &cpu.b)

	case 0xcbc1: // SET 0,C
		cpu.set(0, &cpu.c)

	case 0xcbc2: // SET 0,D
		cpu.set(0, &cpu.d)

	case 0xcbc3: // SET 0,E
		cpu.set(0, &cpu.e)

	case 0xcbc4: // SET 0,H
		cpu.set(0, &cpu.h)

	case 0xcbc5: // SET 0,L
		cpu.set(0, &cpu.l)

	case 0xcbc6: // SET 0,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.set(0, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcbc7: // SET 0,A
		cpu.set(0, &cpu.a)

	case 0xcbc8: // SET 1,B
		cpu.set(1, &cpu.b)

	case 0xcbc9: // SET 1,C
		cpu.set(1, &cpu.c)

	case 0xcbca: // SET 1,D
		cpu.set(1, &cpu.d)

	case 0xcbcb: // SET 1,E
		cpu.set(1, &cpu.e)

	case 0xcbcc: // SET 1,H
		cpu.set(1, &cpu.h)

	case 0xcbcd: // SET 1,L
		cpu.set(1, &cpu.l)

	case 0xcbce: // SET 1,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.set(1, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcbcf: // SET 1,A
		cpu.set(1, &cpu.a)

	case 0xcbd0: // SET 2,B
		cpu.set(2, &cpu.b)

	case 0xcbd1: // SET 2,C
		cpu.set(2, &cpu.c)

	case 0xcbd2: // SET 2,D
		cpu.set(2, &cpu.d)

	case 0xcbd3: // SET 2,E
		cpu.set(2, &cpu.e)

	case 0xcbd4: // SET 2,H
		cpu.set(2, &cpu.h)

	case 0xcbd5: // SET 2,L
		cpu.set(2, &cpu.l)

	case 0xcbd6: // SET 2,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.set(2, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcbd7: // SET 2,A
		cpu.set(2, &cpu.a)

	case 0xcbd8: // SET 3,B
		cpu.set(3, &cpu.b)

	case 0xcbd9: // SET 3,C
		cpu.set(3, &cpu.c)

	case 0xcbda: // SET 3,D
		cpu.set(3, &cpu.d)

	case 0xcbdb: // SET 3,E
		cpu.set(3, &cpu.e)

	case 0xcbdc: // SET 3,H
		cpu.set(3, &cpu.h)

	case 0xcbdd: // SET 3,L
		cpu.set(3, &cpu.l)

	case 0xcbde: // SET 3,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.set(3, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcbdf: // SET 3,A
		cpu.set(3, &cpu.a)

	case 0xcbe0: // SET 4,B
		cpu.set(4, &cpu.b)

	case 0xcbe1: // SET 4,C
		cpu.set(4, &cpu.c)

	case 0xcbe2: // SET 4,D
		cpu.set(4, &cpu.d)

	case 0xcbe3: // SET 4,E
		cpu.set(4, &cpu.e)

	case 0xcbe4: // SET 4,H
		cpu.set(4, &cpu.h)

	case 0xcbe5: // SET 4,L
		cpu.set(4, &cpu.l)

	case 0xcbe6: // SET 4,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.set(4, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcbe7: // SET 4,A
		cpu.set(4, &cpu.a)

	case 0xcbe8: // SET 5,B
		cpu.set(5, &cpu.b)

	case 0xcbe9: // SET 5,C
		cpu.set(5, &cpu.c)

	case 0xcbea: // SET 5,D
		cpu.set(5, &cpu.d)

	case 0xcbeb: // SET 5,E
		cpu.set(5, &cpu.e)

	case 0xcbec: // SET 5,H
		cpu.set(5, &cpu.h)

	case 0xcbed: // SET 5,L
		cpu.set(5, &cpu.l)

	case 0xcbee: // SET 5,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.set(5, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcbef: // SET 5,A
		cpu.set(5, &cpu.a)

	case 0xcbf0: // SET 6,B
		cpu.set(6, &cpu.b)

	case 0xcbf1: // SET 6,C
		cpu.set(6, &cpu.c)

	case 0xcbf2: // SET 6,D
		cpu.set(6, &cpu.d)

	case 0xcbf3: // SET 6,E
		cpu.set(6, &cpu.e)

	case 0xcbf4: // SET 6,H
		cpu.set(6, &cpu.h)

	case 0xcbf5: // SET 6,L
		cpu.set(6, &cpu.l)

	case 0xcbf6: // SET 6,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.set(6, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcbf7: // SET 6,A
		cpu.set(6, &cpu.a)

	case 0xcbf8: // SET 7,B
		cpu.set(7, &cpu.b)

	case 0xcbf9: // SET 7,C
		cpu.set(7, &cpu.c)

	case 0xcbfa: // SET 7,D
		cpu.set(7, &cpu.d)

	case 0xcbfb: // SET 7,E
		cpu.set(7, &cpu.e)

	case 0xcbfc: // SET 7,H
		cpu.set(7, &cpu.h)

	case 0xcbfd: // SET 7,L
		cpu.set(7, &cpu.l)

	case 0xcbfe: // SET 7,(HL)
		b := cpu.memory.GetByte(getRR(cpu.h, cpu.l))
		cpu.set(7, &b)
		cpu.memory.PutByte(getRR(cpu.h, cpu.l), b)

	case 0xcbff: // SET 7,A
		cpu.set(7, &cpu.a)

	case 0xdd86: // ADD A,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.addA(cpu.memory.GetByte(ix))

	case 0xdd8e: // ADC A,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.adc(cpu.memory.GetByte(ix))

	case 0xfd8e: // ADC A,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.adc(cpu.memory.GetByte(iy))

	case 0xdd96: // SUB (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.subA(b)
		cpu.memory.PutByte(ix, b)

	case 0xfd96: // SUB (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		b := cpu.memory.GetByte(iy)
		cpu.subA(b)
		cpu.memory.PutByte(iy, b)

	case 0xdd9e: // SBC A,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.sbc(cpu.memory.GetByte(ix))

	case 0xfd9e: // SBC A,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.sbc(cpu.memory.GetByte(iy))

	case 0xdda6: // AND (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.and(cpu.memory.GetByte(ix))

	case 0xfda6: // AND (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.and(cpu.memory.GetByte(iy))

	case 0xddae: // XOR (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.xor(cpu.memory.GetByte(ix))

	case 0xfdae: // XOR (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.xor(cpu.memory.GetByte(iy))

	case 0xddb6: // OR (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.or(cpu.memory.GetByte(ix))

	case 0xfdb6: // OR (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.or(cpu.memory.GetByte(iy))

	case 0xddbe: // CP (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.cp(cpu.memory.GetByte(ix))

	case 0xfdbe: // CP (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.cp(cpu.memory.GetByte(iy))

	case 0xddcb06: // RLC (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.rlc(&b)
		cpu.memory.PutByte(ix, b)

	case 0xddcb0e: // RRC (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.rrc(&b)
		cpu.memory.PutByte(ix, b)

	case 0xddcb16: // RL (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.rl(&b)
		cpu.memory.PutByte(ix, b)

	case 0xddcb1e: // RR (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.rr(&b)
		cpu.memory.PutByte(ix, b)

	case 0xddcb26: // SLA (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.sla(&b)
		cpu.memory.PutByte(ix, b)

	case 0xddcb2e: // SRA (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.sra(&b)
		cpu.memory.PutByte(ix, b)

	case 0xfdcb2e: // SRA (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		b := cpu.memory.GetByte(iy)
		cpu.sra(&b)
		cpu.memory.PutByte(iy, b)

	case 0xddcb36: // SLL (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.sll(&b)
		cpu.memory.PutByte(ix, b)

	case 0xfdcb36: // SLL (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		b := cpu.memory.GetByte(iy)
		cpu.sll(&b)
		cpu.memory.PutByte(iy, b)

	case 0xddcb3e: // SRL (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.srl(&b)
		cpu.memory.PutByte(ix, b)

	case 0xfdcb3e: // SRL (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		b := cpu.memory.GetByte(iy)
		cpu.srl(&b)
		cpu.memory.PutByte(iy, b)

	case 0xddcb46: // BIT 0,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.bit(0, b)

	case 0xddcb4e: // BIT 1,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.bit(1, b)

	case 0xddcb56: // BIT 2,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.bit(2, b)

	case 0xddcb5e: // BIT 3,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.bit(3, b)

	case 0xddcb66: // BIT 4,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.bit(4, b)

	case 0xddcb6e: // BIT 5,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.bit(5, b)

	case 0xddcb76: // BIT 6,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.bit(6, b)

	case 0xddcb7e: // BIT 7,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.bit(7, b)

	case 0xddcb86: // RES 0,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.res(0, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcb8e: // RES 1,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.res(1, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcb96: // RES 2,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.res(2, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcb9e: // RES 3,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.res(3, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcba6: // RES 4,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.res(4, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcbae: // RES 5,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.res(5, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcbb6: // RES 6,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.res(6, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcbbe: // RES 7,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.res(7, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcbc6: // SET 0,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.set(0, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcbce: // SET 1,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.set(1, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcbd6: // SET 2,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.set(2, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcbde: // SET 3,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.set(3, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcbe6: // SET 4,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.set(4, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcbee: // SET 5,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.set(5, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcbf6: // SET 6,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.set(6, &b)
		cpu.memory.PutByte(ix, b)

	case 0xddcbfe: // SET 7,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		b := cpu.memory.GetByte(ix)
		cpu.set(7, &b)
		cpu.memory.PutByte(ix, b)

	case 0xDD44: //  LD B,IXH
		cpu.b = cpu.ixh

	case 0xFD44: //  LD B,IYH
		cpu.b = cpu.iyh

	case 0xDD45: //  LD B,IXL
		cpu.b = cpu.ixl

	case 0xFD45: //  LD B,IYL
		cpu.b = cpu.iyl

	case 0xDD4C: //  LD C,IXH
		cpu.c = cpu.ixh

	case 0xFD4C: //  LD C,IYH
		cpu.c = cpu.iyh

	case 0xDD4D: //  LD C,IXL
		cpu.c = cpu.ixl

	case 0xFD4D: //  LD C,IYL
		cpu.c = cpu.iyl

	case 0xDD54: //  LD D,IXH
		cpu.d = cpu.ixh

	case 0xFD54: //  LD D,IYH
		cpu.d = cpu.iyh

	case 0xDD55: //  LD D,IXL
		cpu.d = cpu.ixl

	case 0xFD55: //  LD D,IYL
		cpu.d = cpu.iyl

	case 0xDD5C: //  LD E,IXH
		cpu.e = cpu.ixh

	case 0xFD5C: //  LD E,IYH
		cpu.e = cpu.iyh

	case 0xDD5D: //  LD E,IXL
		cpu.e = cpu.ixl

	case 0xFD5D: //  LD E,IYL
		cpu.e = cpu.iyl

	case 0xDD7D: //  LD A,IXL
		cpu.a = cpu.ixl

	case 0xDD7C: //  LD A,IXH
		cpu.a = cpu.ixh

	case 0xFD7C: //  LD A,IYH
		cpu.a = cpu.iyh

	case 0xFD7D: //  LD A,IYL
		cpu.a = cpu.iyl

	case 0xDD60: // LD IXH,B
		cpu.ixh = cpu.b

	case 0xFD60: // LD IYH,B
		cpu.iyh = cpu.b

	case 0xDD61: // LD IXH,C
		cpu.ixh = cpu.c

	case 0xFD61: // LD IYH,C
		cpu.iyh = cpu.c

	case 0xDD62: // LD IXH,D
		cpu.ixh = cpu.d

	case 0xFD62: // LD IYH,D
		cpu.iyh = cpu.d

	case 0xDD63: // LD IXH,E
		cpu.ixh = cpu.e

	case 0xFD63: // LD IYH,E
		cpu.iyh = cpu.e

	case 0xDD67: // LD IXH,A
		cpu.ixh = cpu.a

	case 0xFD67: // LD IYH,A
		cpu.iyh = cpu.a

	case 0xDD68: // LD IXL,B
		cpu.ixl = cpu.b

	case 0xFD68: // LD IYL,B
		cpu.iyl = cpu.b

	case 0xDD69: // LD IXL,C
		cpu.ixl = cpu.c

	case 0xFD69: // LD IYL,C
		cpu.iyl = cpu.c

	case 0xDD6A: // LD IXL,D
		cpu.ixl = cpu.d

	case 0xFD6A: // LD IYL,D
		cpu.iyl = cpu.d

	case 0xDD6B: // LD IXL,E
		cpu.ixl = cpu.e

	case 0xFD6B: // LD IYL,E
		cpu.iyl = cpu.e

	case 0xDD6F: // LD IXL,A
		cpu.ixl = cpu.a

	case 0xFD6F: // LD IYL,A
		cpu.iyl = cpu.a

	case 0xDD64: // LD IXH,IXH
		cpu.ixh = cpu.ixh

	case 0xFD64: // LD IYH,IYH
		cpu.iyh = cpu.iyh

	case 0xDD65: // LD IXH,IXL
		cpu.ixh = cpu.ixl

	case 0xFD65: // LD IYH,IYL
		cpu.iyh = cpu.iyl

	case 0xDD6C: // LD IXL,IXH
		cpu.ixl = cpu.ixh

	case 0xFD6C: // LD IYL,IYH
		cpu.iyl = cpu.iyh

	case 0xDD6D: // LD IXL,IXL
		cpu.ixl = cpu.ixl

	case 0xFD6D: // LD IYL,IYL
		cpu.iyl = cpu.iyl

	case 0xDD84: // ADD A,IXH
		cpu.addA(cpu.ixh)

	case 0xFD84: // ADD A,IYH
		cpu.addA(cpu.iyh)

	case 0xDD85: // ADD A,IXL
		cpu.addA(cpu.ixl)

	case 0xFD85: // ADD A,IYL
		cpu.addA(cpu.iyl)

	case 0xDD8C: // ADC A,IXH
		cpu.adc(cpu.ixh)

	case 0xFD8C: // ADC A,IYH
		cpu.adc(cpu.iyh)

	case 0xDD8D: // ADC A,IXL
		cpu.adc(cpu.ixl)

	case 0xFD8D: // ADC A,IYL
		cpu.adc(cpu.iyl)

	case 0xDD94: // SUB IXH
		cpu.subA(cpu.ixh)

	case 0xFD94: // SUB IYH
		cpu.subA(cpu.iyh)

	case 0xDD95: // SUB IXL
		cpu.subA(cpu.ixl)

	case 0xFD95: // SUB IYL
		cpu.subA(cpu.iyl)

	case 0xDD9C: // SBC A,IXH
		cpu.sbc(cpu.ixh)

	case 0xFD9C: // SBC A,IYH
		cpu.sbc(cpu.iyh)

	case 0xDD9D: // SBC A,IXL
		cpu.sbc(cpu.ixl)

	case 0xFD9D: // SBC A,IYL
		cpu.sbc(cpu.iyl)

	case 0xDDA4: // AND IXH
		cpu.and(cpu.ixh)

	case 0xFDA4: // AND IYH
		cpu.and(cpu.iyh)

	case 0xDDA5: // AND IXL
		cpu.and(cpu.ixl)

	case 0xFDA5: // AND IYL
		cpu.and(cpu.iyl)

	case 0xDDAC: // XOR IXH
		cpu.xor(cpu.ixh)

	case 0xFDAC: // XOR IYH
		cpu.xor(cpu.iyh)

	case 0xDDAD: // XOR IXL
		cpu.xor(cpu.ixl)

	case 0xFDAD: // XOR IYL
		cpu.xor(cpu.iyl)

	case 0xDDB4: // OR IXH
		cpu.or(cpu.ixh)

	case 0xFDB4: // OR IYH
		cpu.or(cpu.iyh)

	case 0xDDB5: // OR IXL
		cpu.or(cpu.ixl)

	case 0xFDB5: // OR IYL
		cpu.or(cpu.iyl)

	case 0xDDBC: // CP IXH
		cpu.cp(cpu.ixh)

	case 0xFDBC: // CP IYH
		cpu.cp(cpu.iyh)

	case 0xDDBD: // CP IXL
		cpu.cp(cpu.ixl)

	case 0xFDBD: // CP IYL
		cpu.cp(cpu.iyl)

	case 0xDD24: // INC IXH
		cpu.incR(&cpu.ixh)

	case 0xFD24: // INC IYH
		cpu.incR(&cpu.iyh)

	case 0xDD25: // DEC IXH
		cpu.decR(&cpu.ixh)

	case 0xFD25: // DEC IYH
		cpu.decR(&cpu.iyh)

	case 0xDD2C: // INC IXL
		cpu.incR(&cpu.ixl)

	case 0xFD2C: // INC IYL
		cpu.incR(&cpu.iyl)

	case 0xDD2D: // DEC IXL
		cpu.decR(&cpu.ixl)

	case 0xFD2D: // DEC IYL
		cpu.decR(&cpu.iyl)

	default:
		panic(fmt.Sprintf("\n----\nopt code '0x%02x: // %s'(%db) not supported\npc: 0x%04x\n----\n", ins.Instruction, ins.Opcode, ins.Length, cpu.pc))
	}
	return needPcUpdate
}

func toWord(a, b byte) uint16 {
	return uint16(a) | uint16(b)<<8
}
