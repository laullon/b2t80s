package z80

import (
	"fmt"

	"github.com/laullon/b2t80s/emulator"
)

func (cpu *z80) runSwitch(ins emulator.Instruction) bool {
	needPcUpdate := true
	switch ins.Instruction {
	case 0x00: // NOP

	case 0x17: // RLA
		c := cpu.regs.F.C
		cpu.regs.F.C = cpu.regs.A&0b10000000 != 0
		cpu.regs.A = (cpu.regs.A << 1)
		if c {
			cpu.regs.A |= 1
		}
		cpu.regs.F.H = false
		cpu.regs.F.N = false

	case 0x0f: // RRCA
		cpu.regs.F.C = cpu.regs.A&0x01 != 0
		cpu.regs.F.H = false
		cpu.regs.F.N = false
		cpu.regs.A = (cpu.regs.A >> 1) | (cpu.regs.A << 7)

	case 0x1f: // RRA
		c := cpu.regs.F.C
		cpu.regs.F.C = cpu.regs.A&1 != 0
		cpu.regs.A = (cpu.regs.A >> 1)
		if c {
			cpu.regs.A |= 0b10000000
		}
		cpu.regs.F.H = false
		cpu.regs.F.N = false

	case 0xed67: // RRD
		hl := cpu.regs.HL.Get()
		hlv := cpu.memory.GetByte(hl)
		cpu.memory.PutByte(hl, (cpu.regs.A<<4 | hlv>>4))
		cpu.regs.A = (cpu.regs.A & 0xf0) | (hlv & 0x0f)

		cpu.regs.F.S = cpu.regs.A&0x80 != 0
		cpu.regs.F.Z = cpu.regs.A == 0
		cpu.regs.F.P = parityTable[cpu.regs.A]
		cpu.regs.F.H = false
		cpu.regs.F.N = false

	case 0xed6f: // RLD
		hl := cpu.regs.HL.Get()
		hlv := cpu.memory.GetByte(hl)
		cpu.memory.PutByte(hl, (hlv<<4 | cpu.regs.A&0x0f))
		cpu.regs.A = (cpu.regs.A & 0xf0) | (hlv >> 4)

		cpu.regs.F.S = cpu.regs.A&0x80 != 0
		cpu.regs.F.Z = cpu.regs.A == 0
		cpu.regs.F.P = parityTable[cpu.regs.A]
		cpu.regs.F.H = false
		cpu.regs.F.N = false

	case 0x3f: // CCF
		cpu.regs.F.H = cpu.regs.F.C
		cpu.regs.F.N = false
		cpu.regs.F.C = !cpu.regs.F.C

	case 0xcd: // CALL NN
		nn := toWord(ins.Mem[1], ins.Mem[2])
		cpu.regs.PC += uint16(ins.Length)
		// cpu.regs.SP.Push(cpu.regs.PC)
		cpu.regs.PC = nn
		needPcUpdate = false

	case 0xec: // CALL PE,NN
		if cpu.regs.F.P {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.regs.PC += uint16(ins.Length)
			// cpu.regs.SP.Push(cpu.regs.PC)
			cpu.regs.PC = nn
			needPcUpdate = false
		}

	case 0xe4: // CALL PO,NN
		if !cpu.regs.F.P {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.regs.PC += uint16(ins.Length)
			// cpu.regs.SP.Push(cpu.regs.PC)
			cpu.regs.PC = nn
			needPcUpdate = false
		}

	case 0xd4: // CALL NC,NN
		if cpu.regs.F.C == false {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.regs.PC += uint16(ins.Length)
			// cpu.regs.SP.Push(cpu.regs.PC)
			cpu.regs.PC = nn
			needPcUpdate = false
		}

	case 0xdc: // CALL C,NN
		if cpu.regs.F.C {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.regs.PC += uint16(ins.Length)
			// cpu.regs.SP.Push(cpu.regs.PC)
			cpu.regs.PC = nn
			needPcUpdate = false
		}

	case 0xcc: // CALL Z,NN
		if cpu.regs.F.Z == true {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.regs.PC += uint16(ins.Length)
			// cpu.regs.SP.Push(cpu.regs.PC)
			cpu.regs.PC = nn
			needPcUpdate = false
		}

	case 0xc4: // CALL NZ,NN
		if cpu.regs.F.Z == false {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.regs.PC += uint16(ins.Length)
			// cpu.regs.SP.Push(cpu.regs.PC)
			cpu.regs.PC = nn
			needPcUpdate = false
		}

	case 0xf4: // CALL P,NN
		if !cpu.regs.F.S {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.regs.PC += uint16(ins.Length)
			// cpu.regs.SP.Push(cpu.regs.PC)
			cpu.regs.PC = nn
			needPcUpdate = false
		}

	case 0xfc: // CALL M,NN
		if cpu.regs.F.S {
			nn := toWord(ins.Mem[1], ins.Mem[2])
			cpu.regs.PC += uint16(ins.Length)
			// cpu.regs.SP.Push(cpu.regs.PC)
			cpu.regs.PC = nn
			needPcUpdate = false
		}

	case 0xf0: // RET P
		if !cpu.regs.F.S {
			// cpu.regs.PC = cpu.regs.SP.Pop()
			needPcUpdate = false
		}

	case 0xf8: // RET M
		if cpu.regs.F.S {
			// cpu.regs.PC = cpu.regs.SP.Pop()
			needPcUpdate = false
		}

	case 0xc8: // RET Z
		if cpu.regs.F.Z {
			// cpu.regs.PC = cpu.regs.SP.Pop()
			needPcUpdate = false
		}

	case 0xc0: // RET NZ
		if !cpu.regs.F.Z {
			// cpu.regs.PC = cpu.regs.SP.Pop()
			needPcUpdate = false
		}

	case 0xd8: // RET C
		if cpu.regs.F.C {
			// cpu.regs.PC = cpu.regs.SP.Pop()
			needPcUpdate = false
		}

	case 0xe8: // RET PE
		if cpu.regs.F.P {
			// cpu.regs.PC = cpu.regs.SP.Pop()
			needPcUpdate = false
		}

	case 0xe0: // RET PO
		if !cpu.regs.F.P {
			// cpu.regs.PC = cpu.regs.SP.Pop()
			needPcUpdate = false
		}

	case 0xd0: // RET NC
		if !cpu.regs.F.C {
			// cpu.regs.PC = cpu.regs.SP.Pop()
			needPcUpdate = false
		}

	case 0xc9: // RET
		// cpu.regs.PC = cpu.regs.SP.Pop()
		needPcUpdate = false

	case 0xed4d: // RETI
		// cpu.regs.PC = cpu.regs.SP.Pop()
		needPcUpdate = false

	case 0xed45: // RETN
		// TODO IFF1=IFF2 ??
		// cpu.regs.PC = cpu.regs.SP.Pop()
		needPcUpdate = false

	case 0xc7: // RST 0
		// cpu.regs.SP.Push(cpu.regs.PC + ins.Length)
		cpu.regs.PC = 0
		needPcUpdate = false

	case 0xcf: // RST 08H
		// cpu.regs.SP.Push(cpu.regs.PC + ins.Length)
		cpu.regs.PC = 0x08
		needPcUpdate = false

	case 0xd7: // RST 10H
		// cpu.regs.SP.Push(cpu.regs.PC + ins.Length)
		cpu.regs.PC = 0x10
		needPcUpdate = false

	case 0xdf: // RST 18H
		// cpu.regs.SP.Push(cpu.regs.PC + ins.Length)
		cpu.regs.PC = 0x18
		needPcUpdate = false

	case 0xe7: // RST 20H
		// cpu.regs.SP.Push(cpu.regs.PC + ins.Length)
		cpu.regs.PC = 0x20
		needPcUpdate = false

	case 0xef: // RST 28H
		// cpu.regs.SP.Push(cpu.regs.PC + ins.Length)
		cpu.regs.PC = 0x28
		needPcUpdate = false

	case 0xf7: // RST 30H
		// cpu.regs.SP.Push(cpu.regs.PC + ins.Length)
		cpu.regs.PC = 0x30
		needPcUpdate = false

	case 0xff: // RST 38H
		// cpu.regs.SP.Push(cpu.regs.PC + ins.Length)
		cpu.regs.PC = 0x38
		needPcUpdate = false

	// case 0xf5: // PUSH AF
	// 	af := uint16(cpu.regs.A)<<8 | uint16(cpu.regs.F.GetByte())
	// 	// cpu.regs.SP.Push(af)

	// case 0xc5: // PUSH BC
	// 	bc := cpu.regs.BC.Get()
	// 	// cpu.regs.SP.Push(bc)

	// case 0xd5: // PUSH DE
	// 	de := cpu.regs.DE.Get()
	// 	// cpu.regs.SP.Push(de)

	// case 0xe5: // PUSH HL
	// 	hl := cpu.regs.HL.Get()
	// 	// cpu.regs.SP.Push(hl)

	// case 0xdde5: // PUSH IX
	// 	ix := cpu.regs.IX.Get()
	// 	// cpu.regs.SP.Push(ix)

	// case 0xfde5: // PUSH IY
	// 	iy := cpu.regs.IY.Get()
	// 	// cpu.regs.SP.Push(iy)

	// case 0xf1: // POP AF
	// 	// af := cpu.regs.SP.Pop()
	// 	cpu.regs.A = byte(af >> 8)
	// 	cpu.regs.F.SetByte(byte(af & 0xff))

	// case 0xc1: // POP BC
	// 	// bc := cpu.regs.SP.Pop()
	// 	cpu.regs.BC.Set(bc)

	// case 0xd1: // POP DE
	// 	// de := cpu.regs.SP.Pop()
	// 	cpu.regs.DE.Set(de)

	// case 0xe1: // POP HL
	// 	// hl := cpu.regs.SP.Pop()
	// 	cpu.regs.HL.Set(hl)

	// case 0xdde1: // POP IX
	// 	// ix := cpu.regs.SP.Pop()
	// 	cpu.regs.IX.Set(ix)

	// case 0xfde1: // POP IY
	// 	// iy := cpu.regs.SP.Pop()
	// 	cpu.regs.IY.Set(iy)

	case 0xF3: // DI
		cpu.regs.IFF1 = false
		cpu.regs.IFF2 = false

	case 0xfb: // EI
		cpu.regs.IFF1 = true
		cpu.regs.IFF2 = true

	case 0xed46: // IM 0
		cpu.regs.InterruptsMode = 0

	case 0xed56: // IM 1
		cpu.regs.InterruptsMode = 1

	case 0xed5e: // IM 2
		cpu.regs.InterruptsMode = 2

	case 0x27: // DAA
		c := cpu.regs.F.C
		add := byte(0)
		if cpu.regs.F.H || ((cpu.regs.A & 0x0f) > 9) {
			add = 6
		}
		if c || (cpu.regs.A > 0x99) {
			add |= 0x60
		}
		if cpu.regs.A > 0x99 {
			c = true
		}
		if cpu.regs.F.N {
			cpu.subA(add)
		} else {
			cpu.addA(add)
		}
		cpu.regs.F.S = int8(cpu.regs.A) < 0
		cpu.regs.F.Z = cpu.regs.A == 0
		cpu.regs.F.P = parityTable[cpu.regs.A]
		cpu.regs.F.C = c

	case 0xAF: // XOR A
		cpu.xor(cpu.regs.A)

	case 0xa8: // XOR B
		cpu.xor(cpu.regs.B)

	case 0xa9: // XOR C
		cpu.xor(cpu.regs.C)

	case 0xaa: // XOR D
		cpu.xor(cpu.regs.D)

	case 0xab: // XOR E
		cpu.xor(cpu.regs.E)

	case 0xac: // XOR H
		cpu.xor(cpu.regs.H)

	case 0xad: // XOR L
		cpu.xor(cpu.regs.L)

	case 0xae: // XOR (HL)
		hl := cpu.regs.HL.Get()
		cpu.xor(cpu.memory.GetByte(hl))

	case 0xee: // XOR N
		cpu.xor(ins.Mem[1])

	case 0x02: // LD (BC),A
		pos := cpu.regs.BC.Get()
		cpu.memory.PutByte(pos, cpu.regs.A)

	case 0x12: // LD (DE),A
		pos := cpu.regs.DE.Get()
		cpu.memory.PutByte(pos, cpu.regs.A)

	case 0x7f: // LD A,A

	case 0x78: // LD A,B
		cpu.regs.A = cpu.regs.B

	case 0x79: // LD A,C
		cpu.regs.A = cpu.regs.C

	case 0x7a: // LD A,D
		cpu.regs.A = cpu.regs.D

	case 0x7b: // LD A,E
		cpu.regs.A = cpu.regs.E

	case 0x7c: // LD A,H
		cpu.regs.A = cpu.regs.H

	case 0x7d: // LD A,L
		cpu.regs.A = cpu.regs.L

	case 0xed57: // LD A,I
		cpu.regs.A = cpu.regs.I
		cpu.regs.F.S = cpu.regs.A&0x80 != 0
		cpu.regs.F.Z = cpu.regs.A == 0
		cpu.regs.F.H = false
		cpu.regs.F.P = cpu.regs.IFF2
		cpu.regs.F.N = false

	case 0xed5f: // LD A,R TODO: review this and its test
		cpu.regs.A = (cpu.regs.R & 0x7f) | (cpu.regs.R7 & 0x80)
		cpu.regs.F.S = cpu.regs.A&0x80 != 0
		cpu.regs.F.Z = cpu.regs.A == 0
		cpu.regs.F.H = false
		cpu.regs.F.P = cpu.regs.IFF2
		cpu.regs.F.N = false

	case 0x3e: // LD A,n
		cpu.regs.A = ins.Mem[1]

	case 0x47: // LD B,A
		cpu.regs.B = cpu.regs.A

	case 0x40: // LD B,B
		// cpu.regs.B = cpu.regs.B

	case 0x41: // LD B,C
		cpu.regs.B = cpu.regs.C

	case 0x42: // LD B,D
		cpu.regs.B = cpu.regs.D

	case 0x43: // LD B,E
		cpu.regs.B = cpu.regs.E

	case 0x44: // LD B,H
		cpu.regs.B = cpu.regs.H

	case 0x45: // LD B,L
		cpu.regs.B = cpu.regs.L

	case 0x46: // LD B,(HL)
		cpu.regs.B = cpu.memory.GetByte(cpu.regs.HL.Get())

	case 0x4f: // LD C,A
		cpu.regs.C = cpu.regs.A

	case 0x48: // LD C,B
		cpu.regs.C = cpu.regs.B

	case 0x49: // LD C,C

	case 0x4a: // LD C,D
		cpu.regs.C = cpu.regs.D

	case 0x4b: // LD C,E
		cpu.regs.C = cpu.regs.E

	case 0x4c: // LD C,H
		cpu.regs.C = cpu.regs.H

	case 0x4d: // LD C,L
		cpu.regs.C = cpu.regs.L

	case 0x57: // LD D,A
		cpu.regs.D = cpu.regs.A

	case 0x50: // LD D,B
		cpu.regs.D = cpu.regs.B

	case 0x51: // LD D,C
		cpu.regs.D = cpu.regs.C

	case 0x52: // LD D,D

	case 0x53: // LD D,E
		cpu.regs.D = cpu.regs.E

	case 0x55: // LD D,L
		cpu.regs.D = cpu.regs.L

	case 0x54: // LD D,H
		cpu.regs.D = cpu.regs.H

	case 0x58: // LD E,B
		cpu.regs.E = cpu.regs.B

	case 0x59: // LD E,C
		cpu.regs.E = cpu.regs.C

	case 0x5a: // LD E,D
		cpu.regs.E = cpu.regs.D

	case 0x5b: // LD E,E

	case 0x5c: // LD E,H
		cpu.regs.E = cpu.regs.H

	case 0x5d: // LD E,L
		cpu.regs.E = cpu.regs.L

	case 0x1e: // LD E,N
		cpu.regs.E = ins.Mem[1]

	case 0x60: // LD H,B
		cpu.regs.H = cpu.regs.B

	case 0x61: // LD H,C
		cpu.regs.H = cpu.regs.C

	case 0x63: // LD H,E
		cpu.regs.H = cpu.regs.E

	case 0x26: // LD H,N
		cpu.regs.H = ins.Mem[1]

	case 0x68: // LD L,B
		cpu.regs.L = cpu.regs.B

	case 0x64: // LD H,H

	case 0x65: // LD H,L
		cpu.regs.H = cpu.regs.L

	case 0x66: // LD H,(HL)
		cpu.regs.H = cpu.memory.GetByte(cpu.regs.HL.Get())

	case 0x69: // LD L,C
		cpu.regs.L = cpu.regs.C

	case 0x6a: // LD L,D
		cpu.regs.L = cpu.regs.D

	case 0x6c: // LD L,H
		cpu.regs.L = cpu.regs.H

	case 0x6d: // LD L,L

	case 0x2e: // LD L,N
		cpu.regs.L = ins.Mem[1]

	case 0x6e: // LD L,(HL)
		cpu.regs.L = cpu.memory.GetByte(cpu.regs.HL.Get())

	case 0xed4f: // LD R,A
		cpu.regs.R = cpu.regs.A

	case 0x01: // LD BC,nn
		cpu.regs.B = ins.Mem[2]
		cpu.regs.C = ins.Mem[1]

	case 0x0e: // LD C,N
		cpu.regs.C = ins.Mem[1]

	case 0x06: // LD B,N
		cpu.regs.B = ins.Mem[1]

	case 0x16: // LD D,N
		// TODO one case fot LD r,n
		cpu.regs.D = ins.Mem[1]

	case 0x5f: // LD E,A
		cpu.regs.E = cpu.regs.A

	case 0x11: // LD DE,nn
		cpu.regs.D = ins.Mem[2]
		cpu.regs.E = ins.Mem[1]

	case 0xed5b: // LD DE,(NN)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		cpu.regs.D = cpu.memory.GetByte(nn + 1)
		cpu.regs.E = cpu.memory.GetByte(nn)

	case 0x0a: // LD A,(BC)
		bc := cpu.regs.BC.Get()
		cpu.regs.A = cpu.memory.GetByte(bc)

	case 0x1a: // LD A,(DE)
		de := cpu.regs.DE.Get()
		cpu.regs.A = cpu.memory.GetByte(de)

	case 0x7e: // LD A,(HL)
		hl := cpu.regs.HL.Get()
		cpu.regs.A = cpu.memory.GetByte(hl)

	case 0x4e: // LD C,(HL)
		// TODO join all LD r,(HL)
		hl := cpu.regs.HL.Get()
		cpu.regs.C = cpu.memory.GetByte(hl)

	case 0x56: // LD D,(HL)
		hl := cpu.regs.HL.Get()
		cpu.regs.D = cpu.memory.GetByte(hl)

	case 0x5e: // LD E,(HL)
		hl := cpu.regs.HL.Get()
		cpu.regs.E = cpu.memory.GetByte(hl)

	case 0x21: // LD HL,nn
		cpu.regs.H = ins.Mem[2]
		cpu.regs.L = ins.Mem[1]

	case 0x3a: // LD A,(NN)
		nn := toWord(ins.Mem[1], ins.Mem[2])
		cpu.regs.A = cpu.memory.GetByte(nn)

	case 0x2a: // LD HL,(nn)
		nn := toWord(ins.Mem[1], ins.Mem[2])
		cpu.regs.L = cpu.memory.GetByte(nn)
		cpu.regs.H = cpu.memory.GetByte(nn + 1)

	case 0x70: // LD (HL),B
		cpu.memory.PutByte(cpu.regs.HL.Get(), cpu.regs.B)

	case 0x71: // LD (HL),C
		cpu.memory.PutByte(cpu.regs.HL.Get(), cpu.regs.C)

	case 0x75: // LD (HL),L
		cpu.memory.PutByte(cpu.regs.HL.Get(), cpu.regs.L)

	case 0x76: // HALT
		if cpu.haltDone {
			cpu.haltDone = false
		} else {
			cpu.halt = true
			needPcUpdate = false
		}

	case 0xdd77: // LD (IX+N),A
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.regs.A)

	case 0xdd70: // LD (IX+N),B
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.regs.B)

	case 0xdd71: // LD (IX+N),C
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.regs.C)

	case 0xdd72: // LD (IX+N),D
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.regs.D)

	case 0xdd73: // LD (IX+N),E
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.regs.E)

	case 0xdd74: // LD (IX+N),H
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.regs.H)

	case 0xdd75: // LD (IX+N),L
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, cpu.regs.L)

	case 0xfd77: // LD (IY+N),A
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.regs.A)

	case 0xfd70: // LD (IY+N),B
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.regs.B)

	case 0xfd71: // LD (IY+N),C
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.regs.C)

	case 0xfd72: // LD (IY+N),D
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.regs.D)

	case 0xfd73: // LD (IY+N),E
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.regs.E)

	case 0xfd74: // LD (IY+N),H
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.regs.H)

	case 0xfd75: // LD (IY+N),L
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, cpu.regs.L)

	case 0xfd36: // LD (IY+N),N
		iy := cpu.getIYn(ins.Mem[2])
		cpu.memory.PutByte(iy, ins.Mem[3])

	case 0xdd36: // LD (IX+N),N
		ix := cpu.getIXn(ins.Mem[2])
		cpu.memory.PutByte(ix, ins.Mem[3])

	case 0xdd7e: // LD A,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.regs.A = cpu.memory.GetByte(ix)

	case 0xdd46: // LD B,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.regs.B = cpu.memory.GetByte(ix)

	case 0xdd4e: // LD C,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.regs.C = cpu.memory.GetByte(ix)

	case 0xdd56: // LD D,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.regs.D = cpu.memory.GetByte(ix)

	case 0xdd5e: // LD E,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.regs.E = cpu.memory.GetByte(ix)

	case 0xdd66: // LD H,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.regs.H = cpu.memory.GetByte(ix)

	case 0xfd66: // LD H,(IY+N)
		IY := cpu.getIYn(ins.Mem[2])
		cpu.regs.H = cpu.memory.GetByte(IY)

	case 0xdd6e: // LD L,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.regs.L = cpu.memory.GetByte(ix)

	case 0xfd7e: // LD A,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.regs.A = cpu.memory.GetByte(iy)

	case 0xfd4e: // LD C,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.regs.C = cpu.memory.GetByte(iy)

	case 0xfd56: // LD D,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.regs.D = cpu.memory.GetByte(iy)

	case 0xfd6e: // LD L,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.regs.L = cpu.memory.GetByte(iy)

	case 0xfd46: // LD B,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.regs.B = cpu.memory.GetByte(iy)

	case 0xfd5e: // LD E,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.regs.E = cpu.memory.GetByte(iy)

	case 0xca: // JP Z,$NN
		if cpu.regs.F.Z {
			cpu.regs.PC = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xc2: // JP NZ,$NN
		if cpu.regs.F.Z == false {
			cpu.regs.PC = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xf2: // JP P,$NN
		if !cpu.regs.F.S {
			cpu.regs.PC = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xfa: // JP M,$NN
		if cpu.regs.F.S {
			cpu.regs.PC = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xea: // JP PE,$NN
		if cpu.regs.F.P {
			cpu.regs.PC = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xe2: // JP PO,$NN
		if !cpu.regs.F.P == true {
			cpu.regs.PC = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xd2: // JP NC,$NN
		if !cpu.regs.F.C {
			cpu.regs.PC = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xda: // JP C,$NN
		if cpu.regs.F.C {
			cpu.regs.PC = toWord(ins.Mem[1], ins.Mem[2])
			needPcUpdate = false
		}

	case 0xC3: // JP nn
		cpu.regs.PC = toWord(ins.Mem[1], ins.Mem[2])
		needPcUpdate = false

	case 0xdde9: // JP (IX)
		cpu.regs.PC = cpu.regs.IX.Get()
		needPcUpdate = false

	case 0xfde9: // JP (IY)
		cpu.regs.PC = cpu.regs.IY.Get()
		needPcUpdate = false

	case 0xe9: // JP (HL)
		hl := cpu.regs.HL.Get()
		cpu.regs.PC = hl
		needPcUpdate = false

	case 0xed47: // LD I,A
		cpu.regs.I = cpu.regs.A

	case 0x67: // LD H,A
		cpu.regs.H = cpu.regs.A

	case 0x62: // LD H,D
		cpu.regs.H = cpu.regs.D

	case 0x6f: // LD L,A
		cpu.regs.L = cpu.regs.A

	case 0x6b: // LD L,E
		cpu.regs.L = cpu.regs.E

	case 0x36: // LD (HL),n
		hl := cpu.regs.HL.Get()
		v := ins.Mem[1]
		cpu.memory.PutByte(hl, v)

	case 0x77: // LD (HL),A
		hl := cpu.regs.HL.Get()
		cpu.memory.PutByte(hl, cpu.regs.A)

	case 0x72: // LD (HL),D
		hl := cpu.regs.HL.Get()
		cpu.memory.PutByte(hl, cpu.regs.D)

	case 0x73: // LD (HL),E
		hl := cpu.regs.HL.Get()
		cpu.memory.PutByte(hl, cpu.regs.E)

	case 0x74: // LD (HL),H
		hl := cpu.regs.HL.Get()
		cpu.memory.PutByte(hl, cpu.regs.H)

	case 0x32: // LD (nn),A
		nn := toWord(ins.Mem[1], ins.Mem[2])
		cpu.memory.PutByte(nn, cpu.regs.A)

	case 0xed43: // LD (nn),BC
		w := cpu.regs.BC.Get()
		nn := toWord(ins.Mem[2], ins.Mem[3])
		putWord(cpu.memory, nn, w)

	case 0xed53: // LD (nn),DE
		w := cpu.regs.DE.Get()
		nn := toWord(ins.Mem[2], ins.Mem[3])
		putWord(cpu.memory, nn, w)

	case 0x22: // LD (nn),HL
		w := cpu.regs.HL.Get()
		nn := toWord(ins.Mem[1], ins.Mem[2])
		putWord(cpu.memory, nn, w)

	case 0xdd22: // LD (NN),IX
		w := cpu.regs.IX.Get()
		nn := toWord(ins.Mem[2], ins.Mem[3])
		putWord(cpu.memory, nn, w)

	case 0xfd22: // LD (NN),IY
		w := cpu.regs.IY.Get()
		nn := toWord(ins.Mem[2], ins.Mem[3])
		putWord(cpu.memory, nn, w)

	case 0xed73: // LD (NN),SP
		nn := toWord(ins.Mem[2], ins.Mem[3])
		putWord(cpu.memory, nn, cpu.regs.SP.Get())

	case 0xdd21: // LD IX,NN
		cpu.regs.IXH = ins.Mem[3]
		cpu.regs.IXL = ins.Mem[2]

	case 0xdd26: // LD IXH,N
		cpu.regs.IXH = ins.Mem[2]

	case 0xdd2e: // LD IXL,N
		cpu.regs.IXL = ins.Mem[2]

	case 0xfd26: // LD IYH,N
		cpu.regs.IYH = ins.Mem[2]

	case 0xfd2e: // LD IYL,N
		cpu.regs.IYL = ins.Mem[2]

	case 0xed4b: // LD BC,(NN)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		cpu.regs.B = cpu.memory.GetByte(nn + 1)
		cpu.regs.C = cpu.memory.GetByte(nn)

	case 0xdd2a: // LD IX,(NN)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		cpu.regs.IXH = cpu.memory.GetByte(nn + 1)
		cpu.regs.IXL = cpu.memory.GetByte(nn)

	case 0xfd2a: // LD IY,(NN)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		cpu.regs.IYH = cpu.memory.GetByte(nn + 1)
		cpu.regs.IYL = cpu.memory.GetByte(nn)

	case 0xed7b: // LD SP,(NN)
		nn := toWord(ins.Mem[2], ins.Mem[3])
		sp := uint16(cpu.memory.GetByte(nn+1))<<8 | uint16(cpu.memory.GetByte(nn))
		cpu.regs.SP.Set(sp)

	case 0xfd21: // LD IY,nn
		cpu.regs.IYH = ins.Mem[3]
		cpu.regs.IYL = ins.Mem[2]

	case 0xf9: // LD SP,HL
		cpu.regs.SP.Set(cpu.regs.HL.Get())

	case 0xddf9: // LD SP,IX
		cpu.regs.SP.Set(cpu.regs.IX.Get())

	case 0xfdf9: // LD SP,IY
		cpu.regs.SP.Set(cpu.regs.IY.Get())

	case 0x31: // LD SP,NN
		cpu.regs.SP.Set(toWord(ins.Mem[1], ins.Mem[2]))

	case 0xeda0: // LDI
		// cpu.ldi()

	case 0xedb0: // LDIR
		// cpu.ldi()
		// bc := cpu.regs.BC.Get()
		// if bc != 0 {
		// 	needPcUpdate = false
		// }

	case 0xeda8: // LDD
		cpu.ldd()

	case 0xedb8: // LDDR
		cpu.ldd()
		bc := cpu.regs.BC.Get()
		if bc != 0 {
			needPcUpdate = false
		}

	case 0xeda9: // CPD
		cpu.cpd()

	case 0x2f: // CPL
		cpu.regs.A = ^cpu.regs.A
		cpu.regs.F.H = true
		cpu.regs.F.N = true

	case 0xed44: // NEG
		n := cpu.regs.A
		cpu.regs.A = 0
		cpu.subA(n)

	case 0x2b: // DEC HL
		// TODO join all DEC rr
		hl := cpu.regs.HL.Get()
		hl--
		cpu.regs.HL.Set(hl)

	case 0x0b: // DEC BC
		bc := cpu.regs.BC.Get()
		bc--
		cpu.regs.BC.Set(bc)

	case 0x1b: // DEC DE
		de := cpu.regs.DE.Get()
		de--
		cpu.regs.DE.Set(de)

	case 0xdd2b: // DEC IX
		ix := cpu.regs.IX.Get()
		ix--
		cpu.regs.IX.Set(ix)

	case 0xfd2b: // DEC IY
		iy := cpu.regs.IY.Get()
		iy--
		cpu.regs.IY.Set(iy)

	case 0x03: // INC BC
		bc := cpu.regs.BC.Get()
		bc++
		cpu.regs.BC.Set(bc)

	case 0x13: // INC DE
		de := cpu.regs.DE.Get()
		de++
		cpu.regs.DE.Set(de)

	case 0x23: // INC HL
		hl := cpu.regs.HL.Get()
		hl++
		cpu.regs.HL.Set(hl)

	case 0x33: // INC SP
		sp := cpu.regs.SP.Get()
		cpu.regs.SP.Set(sp + 1)

	case 0xdd23: // INC IX
		ix := cpu.regs.IX.Get()
		ix++
		cpu.regs.IX.Set(ix)

	case 0xfd23: // INC IY
		iy := cpu.regs.IY.Get()
		iy++
		cpu.regs.IY.Set(iy)

	case 0xdd34: // INC (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		ixv := cpu.memory.GetByte(ix)
		// cpu.incR(&ixv)
		cpu.memory.PutByte(ix, ixv)

	case 0xfd34: // INC (IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		iyv := cpu.memory.GetByte(iy)
		// cpu.incR(&iyv)
		cpu.memory.PutByte(iy, iyv)

	case 0xdd35: // DEC (IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		ixv := cpu.memory.GetByte(ix)
		// cpu.decR(&ixv)
		cpu.memory.PutByte(ix, ixv)

	case 0xfd35: // DEC (IY+d)
		iy := cpu.getIYn(ins.Mem[2])
		iyv := cpu.memory.GetByte(iy)
		// cpu.decR(&iyv)
		cpu.memory.PutByte(iy, iyv)

	case 0x3b: // DEC SP
		sp := cpu.regs.SP.Get()
		cpu.regs.SP.Set(sp - 1)

	case 0x34: // INC (HL)
		hl := cpu.regs.HL.Get()
		b := cpu.memory.GetByte(hl)
		// cpu.incR(&b)
		cpu.memory.PutByte(hl, b)

	case 0x35: // DEC (HL)
		hl := cpu.regs.HL.Get()
		b := cpu.memory.GetByte(hl)
		// cpu.decR(&b)
		cpu.memory.PutByte(hl, b)

	case 0xbf: // CP A
		cpu.cp(cpu.regs.A)

	case 0xb8: // CP B
		cpu.cp(cpu.regs.B)

	case 0xb9: // CP C
		cpu.cp(cpu.regs.C)

	case 0xba: // CP D
		cpu.cp(cpu.regs.D)

	case 0xbb: // CP E
		cpu.cp(cpu.regs.E)

	case 0xbc: // CP H
		cpu.cp(cpu.regs.H)

	case 0xbd: // CP L
		cpu.cp(cpu.regs.L)

	case 0xfe: // CP N
		n := ins.Mem[1]
		cpu.cp(n)

	case 0xbe: // CP (HL)
		cpu.cp(cpu.memory.GetByte(cpu.regs.HL.Get()))

	case 0xeda1: // CPI
		// cpu.cpi()

	case 0xedb1: // CPIR
		// diff := cpu.cpi()
		// bc := cpu.regs.BC.Get()
		// if (bc != 0) && (diff != 0) {
		// 	needPcUpdate = false
		// }

	case 0xedb9: // CPDR
		diff := cpu.cpd()
		bc := cpu.regs.BC.Get()
		if (bc == 0) || (diff == 0) {
			cpu.regs.F.P = bc != 0
		} else {
			needPcUpdate = false
		}

	case 0x20: // JR NZ,d
		if cpu.regs.F.Z == false {
			cpu.jr(ins.Mem[1])
			needPcUpdate = false
		}

	case 0x28: // JR Z,d
		if cpu.regs.F.Z {
			cpu.jr(ins.Mem[1])
			needPcUpdate = false
		}

	case 0x30: // JR NC,d
		if cpu.regs.F.C != true {
			cpu.jr(ins.Mem[1])
			needPcUpdate = false
		}

	case 0x38: // JR C,$N+2
		if cpu.regs.F.C == true {
			cpu.jr(ins.Mem[1])
			needPcUpdate = false
		}

	case 0x37: // SCF
		cpu.regs.F.H = false
		cpu.regs.F.N = false
		cpu.regs.F.C = true

	case 0x18: // JR $N+2
		cpu.jr(ins.Mem[1])
		needPcUpdate = false

	case 0x10: // DJNZ $+2

	case 0xa7: // AND A
		cpu.and(cpu.regs.A)

	case 0xa0: // AND B
		cpu.and(cpu.regs.B)

	case 0xa1: // AND C
		cpu.and(cpu.regs.C)

	case 0xa2: // AND D
		cpu.and(cpu.regs.D)

	case 0xa3: // AND E
		cpu.and(cpu.regs.E)

	case 0xa4: // AND H
		cpu.and(cpu.regs.H)

	case 0xa5: // AND L
		cpu.and(cpu.regs.L)

	case 0xe6: // AND N
		cpu.and(ins.Mem[1])

	case 0xa6: // AND (HL)
		cpu.and(cpu.memory.GetByte(cpu.regs.HL.Get()))

	case 0xb7: // OR A
		cpu.or(cpu.regs.A)

	case 0xb0: // OR B
		cpu.or(cpu.regs.B)

	case 0xb1: // OR C
		cpu.or(cpu.regs.C)

	case 0xb2: // OR D
		cpu.or(cpu.regs.D)

	case 0xb3: // OR E
		cpu.or(cpu.regs.E)

	case 0xb4: // OR H
		cpu.or(cpu.regs.H)

	case 0xb5: // OR L
		cpu.or(cpu.regs.L)

	case 0xb6: // OR (HL)
		cpu.or(cpu.memory.GetByte(cpu.regs.HL.Get()))

	case 0xf6: // OR N
		cpu.or(ins.Mem[1])

	case 0xde: // SBC A,N
		cpu.sbcA(ins.Mem[1])

	case 0xed52: // SBC HL,DE
		cpu.sbcHL(cpu.regs.DE.Get())

	case 0xed42: // SBC HL,BC
		cpu.sbcHL(cpu.regs.BC.Get())

	case 0xed62: // SBC HL,HL
		cpu.sbcHL(cpu.regs.HL.Get())

	case 0xed72: // SBC HL,SP
		cpu.sbcHL(cpu.regs.SP.Get())

	case 0x87: // ADD A,A
		cpu.addA(cpu.regs.A)

	case 0x80: // ADD A,B
		cpu.addA(cpu.regs.B)

	case 0x81: // ADD A,C
		cpu.addA(cpu.regs.C)

	case 0x82: // ADD A,D
		cpu.addA(cpu.regs.D)

	case 0x83: // ADD A,E
		cpu.addA(cpu.regs.E)

	case 0x84: // ADD A,H
		cpu.addA(cpu.regs.H)

	case 0x85: // ADD A,L
		cpu.addA(cpu.regs.L)

	case 0xc6: // ADD A,N
		cpu.addA(ins.Mem[1])

	case 0x86: // ADD A,(HL)
		cpu.addA(cpu.memory.GetByte(cpu.regs.HL.Get()))

	case 0xfd86: // ADD A,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.addA(cpu.memory.GetByte(iy))

	case 0x8f: // ADC A,A
		cpu.adcA(cpu.regs.A)

	case 0x88: // ADC A,B
		cpu.adcA(cpu.regs.B)

	case 0x89: // ADC A,C
		cpu.adcA(cpu.regs.C)

	case 0x8a: // ADC A,D
		cpu.adcA(cpu.regs.D)

	case 0x8b: // ADC A,E
		cpu.adcA(cpu.regs.E)

	case 0x8c: // ADC A,H
		cpu.adcA(cpu.regs.H)

	case 0x8d: // ADC A,L
		cpu.adcA(cpu.regs.L)

	case 0x8e: // ADC A,(HL)
		cpu.adcA(cpu.memory.GetByte(cpu.regs.HL.Get()))

	case 0xce: // ADC A,N
		cpu.adcA(ins.Mem[1])

	case 0xed4a: // ADC HL,BC
		cpu.adcHL(cpu.regs.BC.Get())

	case 0xed5a: // ADC HL,DE
		cpu.adcHL(cpu.regs.DE.Get())

	case 0xed6a: // ADC HL,HL
		cpu.adcHL(cpu.regs.HL.Get())

	case 0xed7a: // ADC HL,SP
		cpu.adcHL(cpu.regs.SP.Get())

	case 0x97: // SUB A
		cpu.subA(cpu.regs.A)

	case 0x90: // SUB B
		cpu.subA(cpu.regs.B)

	case 0x91: // SUB C
		cpu.subA(cpu.regs.C)

	case 0x92: // SUB D
		cpu.subA(cpu.regs.D)

	case 0x93: // SUB E
		cpu.subA(cpu.regs.E)

	case 0x94: // SUB H
		cpu.subA(cpu.regs.H)

	case 0x95: // SUB L
		cpu.subA(cpu.regs.L)

	case 0x96: // SUB (HL)
		cpu.subA(cpu.memory.GetByte(cpu.regs.HL.Get()))

	case 0xd6: // SUB N
		cpu.subA(ins.Mem[1])

	case 0x9f: // SBC A
		cpu.sbcA(cpu.regs.A)

	case 0x98: // SBC B
		cpu.sbcA(cpu.regs.B)

	case 0x99: // SBC C
		cpu.sbcA(cpu.regs.C)

	case 0x9a: // SBC D
		cpu.sbcA(cpu.regs.D)

	case 0x9b: // SBC E
		cpu.sbcA(cpu.regs.E)

	case 0x9c: // SBC H
		cpu.sbcA(cpu.regs.H)

	case 0x9d: // SBC L
		cpu.sbcA(cpu.regs.L)

	case 0x9e: // SBC (HL)
		cpu.sbcA(cpu.memory.GetByte(cpu.regs.HL.Get()))

	case 0x08: // EX AF,AF'
		ta := cpu.regs.A
		cpu.regs.A = cpu.regs.Aalt
		cpu.regs.Aalt = ta

		tf := cpu.regs.F
		cpu.regs.F = cpu.regs.Falt
		cpu.regs.Falt = tf

	case 0xeb: // EX DE,HL
		td := cpu.regs.D
		cpu.regs.D = cpu.regs.H
		cpu.regs.H = td

		te := cpu.regs.E
		cpu.regs.E = cpu.regs.L
		cpu.regs.L = te

	case 0xdde3: // EX (SP),IX
		spv := getWord(cpu.memory, cpu.regs.SP.Get())
		ix := cpu.regs.IX.Get()
		putWord(cpu.memory, cpu.regs.SP.Get(), ix)
		cpu.regs.IX.Set(spv)

	case 0xfde3: // EX (SP),IY
		spv := getWord(cpu.memory, cpu.regs.SP.Get())
		iy := cpu.regs.IY.Get()
		putWord(cpu.memory, cpu.regs.SP.Get(), iy)
		cpu.regs.IY.Set(spv)

	case 0xe3: // EX (SP),HL
		spv := getWord(cpu.memory, cpu.regs.SP.Get())
		hl := cpu.regs.HL.Get()
		putWord(cpu.memory, cpu.regs.SP.Get(), hl)
		cpu.regs.HL.Set(spv)

	case 0xd9: // EXX
		tb := cpu.regs.B
		cpu.regs.B = cpu.regs.Balt
		cpu.regs.Balt = tb

		tc := cpu.regs.C
		cpu.regs.C = cpu.regs.Calt
		cpu.regs.Calt = tc

		td := cpu.regs.D
		cpu.regs.D = cpu.regs.Dalt
		cpu.regs.Dalt = td

		te := cpu.regs.E
		cpu.regs.E = cpu.regs.Ealt
		cpu.regs.Ealt = te

		th := cpu.regs.H
		cpu.regs.H = cpu.regs.Halt
		cpu.regs.Halt = th

		tl := cpu.regs.L
		cpu.regs.L = cpu.regs.Lalt
		cpu.regs.Lalt = tl

	case 0xeda3: // OUTI
		hl := cpu.regs.HL.Get()
		b := cpu.memory.GetByte(hl)
		cpu.regs.B--
		cpu.writePort(cpu.regs.BC.Get(), b)
		hl++
		cpu.regs.HL.Set(hl)
		cpu.regs.F.Z = cpu.regs.B == 0
		cpu.regs.F.S = cpu.regs.B&0x80 != 0
		cpu.regs.F.N = cpu.regs.B&0x80 == 0
		cpu.regs.F.H = true
		cpu.regs.F.P = parityTable[cpu.regs.B]

	case 0xedab: // OUTD
		hl := cpu.regs.HL.Get()
		b := cpu.memory.GetByte(hl)
		cpu.regs.B--
		cpu.writePort(cpu.regs.BC.Get(), b)
		hl--
		cpu.regs.HL.Set(hl)
		cpu.regs.F.Z = cpu.regs.B == 0
		cpu.regs.F.S = cpu.regs.B&0x80 != 0
		cpu.regs.F.N = cpu.regs.B&0x80 == 0
		cpu.regs.F.H = true
		cpu.regs.F.P = parityTable[cpu.regs.B]

	case 0xd3: // OUT (n),A
		port := toWord(ins.Mem[1], cpu.regs.A)
		cpu.writePort(port, cpu.regs.A)

	case 0xed79: // OUT (C),A
		cpu.writePort(cpu.regs.BC.Get(), cpu.regs.A)

	case 0xed41: // OUT (C),B
		cpu.writePort(cpu.regs.BC.Get(), cpu.regs.B)

	case 0xed49: // OUT (C),C
		cpu.writePort(cpu.regs.BC.Get(), cpu.regs.C)

	case 0xed51: // OUT (C),D
		cpu.writePort(cpu.regs.BC.Get(), cpu.regs.D)

	case 0xed59: // OUT (C),E
		cpu.writePort(cpu.regs.BC.Get(), cpu.regs.E)

	case 0xed61: // OUT (C),H
		cpu.writePort(cpu.regs.BC.Get(), cpu.regs.H)

	case 0xed69: // OUT (C),L
		cpu.writePort(cpu.regs.BC.Get(), cpu.regs.L)

	case 0xed71: // OUT (C),0
		cpu.writePort(cpu.regs.BC.Get(), 0)

	case 0xed78: // IN A,(C)
		cpu.regs.A = cpu.readPort(cpu.regs.BC.Get())

	case 0xed40: // IN B,(C)
		cpu.regs.B = cpu.readPort(cpu.regs.BC.Get())

	case 0xed48: // IN C,(C)
		cpu.regs.C = cpu.readPort(cpu.regs.BC.Get())

	case 0xed50: // IN D,(C)
		cpu.regs.D = cpu.readPort(cpu.regs.BC.Get())

	case 0xed58: // IN E,(C)
		cpu.regs.E = cpu.readPort(cpu.regs.BC.Get())

	case 0xed60: // IN H,(C)
		cpu.regs.H = cpu.readPort(cpu.regs.BC.Get())

	case 0xed68: // IN L,(C)
		cpu.regs.L = cpu.readPort(cpu.regs.BC.Get())

	case 0xed70: // IN (C)
		cpu.readPort(cpu.regs.BC.Get())

	case 0xeda2: // INI
		hl := cpu.regs.HL.Get()
		v := cpu.readPort(cpu.regs.BC.Get())
		cpu.memory.PutByte(hl, v)
		cpu.regs.B--
		hl++
		cpu.regs.F.N = v&0x80 != 0
		cpu.regs.F.Z = v == 0
		// cpu.regs.F.H = v == 0
		// cpu.regs.F.C = false
		// cpu.regs.F.P = false
		cpu.regs.HL.Set(hl)

	case 0xdb: // IN A,(N)
		f := cpu.regs.F.GetByte()
		port := toWord(ins.Mem[1], cpu.regs.A)
		cpu.regs.A = cpu.readPort(port)
		cpu.regs.F.SetByte(f)

	case 0xedb3: // OTIR
		hl := cpu.regs.HL.Get()
		c := cpu.memory.GetByte(hl)
		cpu.readPort(uint16(c) << 8)
		cpu.regs.B--
		hl++
		cpu.regs.HL.Set(hl)
		if cpu.regs.B != 0 {
			cpu.regs.PC -= 2
		} else {
			cpu.regs.F.Z = true
			cpu.regs.F.N = true
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

	case 0xdd86: // ADD A,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.addA(cpu.memory.GetByte(ix))

	case 0xdd8e: // ADC A,(IX+N)
		ix := cpu.getIXn(ins.Mem[2])
		cpu.adcA(cpu.memory.GetByte(ix))

	case 0xfd8e: // ADC A,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.adcA(cpu.memory.GetByte(iy))

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
		cpu.sbcA(cpu.memory.GetByte(ix))

	case 0xfd9e: // SBC A,(IY+N)
		iy := cpu.getIYn(ins.Mem[2])
		cpu.sbcA(cpu.memory.GetByte(iy))

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
		cpu.regs.B = cpu.regs.IXH

	case 0xFD44: //  LD B,IYH
		cpu.regs.B = cpu.regs.IYH

	case 0xDD45: //  LD B,IXL
		cpu.regs.B = cpu.regs.IXL

	case 0xFD45: //  LD B,IYL
		cpu.regs.B = cpu.regs.IYL

	case 0xDD4C: //  LD C,IXH
		cpu.regs.C = cpu.regs.IXH

	case 0xFD4C: //  LD C,IYH
		cpu.regs.C = cpu.regs.IYH

	case 0xDD4D: //  LD C,IXL
		cpu.regs.C = cpu.regs.IXL

	case 0xFD4D: //  LD C,IYL
		cpu.regs.C = cpu.regs.IYL

	case 0xDD54: //  LD D,IXH
		cpu.regs.D = cpu.regs.IXH

	case 0xFD54: //  LD D,IYH
		cpu.regs.D = cpu.regs.IYH

	case 0xDD55: //  LD D,IXL
		cpu.regs.D = cpu.regs.IXL

	case 0xFD55: //  LD D,IYL
		cpu.regs.D = cpu.regs.IYL

	case 0xDD5C: //  LD E,IXH
		cpu.regs.E = cpu.regs.IXH

	case 0xFD5C: //  LD E,IYH
		cpu.regs.E = cpu.regs.IYH

	case 0xDD5D: //  LD E,IXL
		cpu.regs.E = cpu.regs.IXL

	case 0xFD5D: //  LD E,IYL
		cpu.regs.E = cpu.regs.IYL

	case 0xDD7D: //  LD A,IXL
		cpu.regs.A = cpu.regs.IXL

	case 0xDD7C: //  LD A,IXH
		cpu.regs.A = cpu.regs.IXH

	case 0xFD7C: //  LD A,IYH
		cpu.regs.A = cpu.regs.IYH

	case 0xFD7D: //  LD A,IYL
		cpu.regs.A = cpu.regs.IYL

	case 0xDD60: // LD IXH,B
		cpu.regs.IXH = cpu.regs.B

	case 0xFD60: // LD IYH,B
		cpu.regs.IYH = cpu.regs.B

	case 0xDD61: // LD IXH,C
		cpu.regs.IXH = cpu.regs.C

	case 0xFD61: // LD IYH,C
		cpu.regs.IYH = cpu.regs.C

	case 0xDD62: // LD IXH,D
		cpu.regs.IXH = cpu.regs.D

	case 0xFD62: // LD IYH,D
		cpu.regs.IYH = cpu.regs.D

	case 0xDD63: // LD IXH,E
		cpu.regs.IXH = cpu.regs.E

	case 0xFD63: // LD IYH,E
		cpu.regs.IYH = cpu.regs.E

	case 0xDD67: // LD IXH,A
		cpu.regs.IXH = cpu.regs.A

	case 0xFD67: // LD IYH,A
		cpu.regs.IYH = cpu.regs.A

	case 0xDD68: // LD IXL,B
		cpu.regs.IXL = cpu.regs.B

	case 0xFD68: // LD IYL,B
		cpu.regs.IYL = cpu.regs.B

	case 0xDD69: // LD IXL,C
		cpu.regs.IXL = cpu.regs.C

	case 0xFD69: // LD IYL,C
		cpu.regs.IYL = cpu.regs.C

	case 0xDD6A: // LD IXL,D
		cpu.regs.IXL = cpu.regs.D

	case 0xFD6A: // LD IYL,D
		cpu.regs.IYL = cpu.regs.D

	case 0xDD6B: // LD IXL,E
		cpu.regs.IXL = cpu.regs.E

	case 0xFD6B: // LD IYL,E
		cpu.regs.IYL = cpu.regs.E

	case 0xDD6F: // LD IXL,A
		cpu.regs.IXL = cpu.regs.A

	case 0xFD6F: // LD IYL,A
		cpu.regs.IYL = cpu.regs.A

	case 0xDD64: // LD IXH,IXH
		cpu.regs.IXH = cpu.regs.IXH

	case 0xFD64: // LD IYH,IYH
		cpu.regs.IYH = cpu.regs.IYH

	case 0xDD65: // LD IXH,IXL
		cpu.regs.IXH = cpu.regs.IXL

	case 0xFD65: // LD IYH,IYL
		cpu.regs.IYH = cpu.regs.IYL

	case 0xDD6C: // LD IXL,IXH
		cpu.regs.IXL = cpu.regs.IXH

	case 0xFD6C: // LD IYL,IYH
		cpu.regs.IYL = cpu.regs.IYH

	case 0xDD6D: // LD IXL,IXL
		cpu.regs.IXL = cpu.regs.IXL

	case 0xFD6D: // LD IYL,IYL
		cpu.regs.IYL = cpu.regs.IYL

	case 0xDD84: // ADD A,IXH
		cpu.addA(cpu.regs.IXH)

	case 0xFD84: // ADD A,IYH
		cpu.addA(cpu.regs.IYH)

	case 0xDD85: // ADD A,IXL
		cpu.addA(cpu.regs.IXL)

	case 0xFD85: // ADD A,IYL
		cpu.addA(cpu.regs.IYL)

	case 0xDD8C: // ADC A,IXH
		cpu.adcA(cpu.regs.IXH)

	case 0xFD8C: // ADC A,IYH
		cpu.adcA(cpu.regs.IYH)

	case 0xDD8D: // ADC A,IXL
		cpu.adcA(cpu.regs.IXL)

	case 0xFD8D: // ADC A,IYL
		cpu.adcA(cpu.regs.IYL)

	case 0xDD94: // SUB IXH
		cpu.subA(cpu.regs.IXH)

	case 0xFD94: // SUB IYH
		cpu.subA(cpu.regs.IYH)

	case 0xDD95: // SUB IXL
		cpu.subA(cpu.regs.IXL)

	case 0xFD95: // SUB IYL
		cpu.subA(cpu.regs.IYL)

	case 0xDD9C: // SBC A,IXH
		cpu.sbcA(cpu.regs.IXH)

	case 0xFD9C: // SBC A,IYH
		cpu.sbcA(cpu.regs.IYH)

	case 0xDD9D: // SBC A,IXL
		cpu.sbcA(cpu.regs.IXL)

	case 0xFD9D: // SBC A,IYL
		cpu.sbcA(cpu.regs.IYL)

	case 0xDDA4: // AND IXH
		cpu.and(cpu.regs.IXH)

	case 0xFDA4: // AND IYH
		cpu.and(cpu.regs.IYH)

	case 0xDDA5: // AND IXL
		cpu.and(cpu.regs.IXL)

	case 0xFDA5: // AND IYL
		cpu.and(cpu.regs.IYL)

	case 0xDDAC: // XOR IXH
		cpu.xor(cpu.regs.IXH)

	case 0xFDAC: // XOR IYH
		cpu.xor(cpu.regs.IYH)

	case 0xDDAD: // XOR IXL
		cpu.xor(cpu.regs.IXL)

	case 0xFDAD: // XOR IYL
		cpu.xor(cpu.regs.IYL)

	case 0xDDB4: // OR IXH
		cpu.or(cpu.regs.IXH)

	case 0xFDB4: // OR IYH
		cpu.or(cpu.regs.IYH)

	case 0xDDB5: // OR IXL
		cpu.or(cpu.regs.IXL)

	case 0xFDB5: // OR IYL
		cpu.or(cpu.regs.IYL)

	case 0xDDBC: // CP IXH
		cpu.cp(cpu.regs.IXH)

	case 0xFDBC: // CP IYH
		cpu.cp(cpu.regs.IYH)

	case 0xDDBD: // CP IXL
		cpu.cp(cpu.regs.IXL)

	case 0xFDBD: // CP IYL
		cpu.cp(cpu.regs.IYL)

	default:
		panic(fmt.Sprintf("\n----\nopt code '0x%02x: // %s'(%db) not supported\npc: 0x%04x\n----\n", ins.Instruction, ins.Opcode, ins.Length, cpu.regs.PC))
	}
	return needPcUpdate
}

func toWord(a, b byte) uint16 {
	return uint16(a) | uint16(b)<<8
}
