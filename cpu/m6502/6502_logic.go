package m6502

var ops []operation

func init() {
	ops = make([]operation, 0x100)

	ops[0x10] = &relative{F: bpl}
	ops[0x20] = &absolute{F: jsr}
	ops[0x18] = &implicit{F: clc}
	ops[0x30] = &relative{F: bmi}
	ops[0x4c] = &absolute{F: jmp}
	ops[0x48] = &implicit{F: pha}
	ops[0x49] = &immediate{F: eor}
	ops[0x60] = &implicit{F: rts}
	ops[0x69] = &immediate{F: adc}
	ops[0x88] = &implicit{F: dey}
	ops[0x8d] = &absolute{F: staM}
	ops[0x90] = &relative{F: bcc}
	ops[0x98] = &implicit{F: tya}
	ops[0x9a] = &implicit{F: txs}
	ops[0xa0] = &immediate{F: ldy}
	ops[0xa2] = &immediate{F: ldx}
	ops[0xa8] = &implicit{F: tay}
	ops[0xa9] = &immediate{F: lda}
	ops[0xaa] = &implicit{F: tax}
	ops[0xad] = &absolute{F: ldaM}
	ops[0xba] = &implicit{F: tsx}
	ops[0xb0] = &relative{F: bcs}
	ops[0xbd] = &absolute{F: ldaM, x: true}
	ops[0xc0] = &immediate{F: cpy}
	ops[0xc9] = &immediate{F: cmp}
	ops[0xca] = &implicit{F: dex}
	ops[0xcd] = &absolute{F: cmpM}
	ops[0xd0] = &relative{F: bne}
	ops[0xd8] = &implicit{F: cld}
	ops[0xe0] = &immediate{F: cpx}
	ops[0xe8] = &implicit{F: inx}
	ops[0xea] = &implicit{F: nop}
	ops[0xf0] = &relative{F: beq}

	for opCode, op := range ops {
		if op != nil {
			op.setup(uint8(opCode), getFunctionName(op))
		}
	}
}

// TODO: review extra cycles
func bne(cpu *m6502, data int8) {
	if !cpu.regs.P.Z {
		jmpr(cpu, data)
	}
}

// TODO: review extra cycles
func bcc(cpu *m6502, data int8) {
	if !cpu.regs.P.C {
		jmpr(cpu, data)
	}
}

// TODO: review extra cycles
func beq(cpu *m6502, data int8) {
	if cpu.regs.P.Z {
		jmpr(cpu, data)
	}
}

// TODO: review extra cycles
func bpl(cpu *m6502, data int8) {
	if !cpu.regs.P.N {
		jmpr(cpu, data)
	}
}

// TODO: review extra cycles
func bmi(cpu *m6502, data int8) {
	if cpu.regs.P.N {
		jmpr(cpu, data)
	}
}

// TODO: review extra cycles
func bcs(cpu *m6502, data int8) {
	if cpu.regs.P.C {
		jmpr(cpu, data)
	}
}

func jmpr(cpu *m6502, data int8) {
	pc := cpu.regs.PC + uint16(data)
	jmp(cpu, pc)
}

func jmp(cpu *m6502, addr uint16) { cpu.regs.PC = addr }

func jsr(cpu *m6502, addr uint16) {
	cpu.push(uint8(cpu.regs.PC >> 8))
	cpu.push(uint8(cpu.regs.PC))
	cpu.regs.PC = addr
}

func rts(cpu *m6502) {
	addr := uint16(cpu.pop())
	addr |= uint16(cpu.pop()) << 8
	cpu.regs.PC = addr
}

func staM(cpu *m6502, addr uint16) { cpu.mem[addr] = cpu.regs.A }

func ldaM(cpu *m6502, addr uint16) {
	cpu.regs.A = cpu.mem[addr]
	ldzn(cpu, cpu.regs.A)
}

func nop(cpu *m6502) {}

func cld(cpu *m6502) { cpu.regs.P.D = false }

func clc(cpu *m6502) { cpu.regs.P.C = false }

func tsx(cpu *m6502) {
	cpu.regs.X = cpu.regs.SP
	ldzn(cpu, cpu.regs.X)
}

func txs(cpu *m6502) {
	cpu.regs.SP = cpu.regs.X
	ldzn(cpu, cpu.regs.SP)
}

func tax(cpu *m6502) {
	cpu.regs.X = cpu.regs.A
	ldzn(cpu, cpu.regs.X)
}

func tay(cpu *m6502) {
	cpu.regs.Y = cpu.regs.A
	ldzn(cpu, cpu.regs.Y)
}

func tya(cpu *m6502) {
	cpu.regs.A = cpu.regs.Y
	ldzn(cpu, cpu.regs.A)
}

func dex(cpu *m6502) {
	cpu.regs.X--
	ldzn(cpu, cpu.regs.X)
}

func inx(cpu *m6502) {
	cpu.regs.X++
	ldzn(cpu, cpu.regs.X)
}

func dey(cpu *m6502) {
	cpu.regs.Y--
	ldzn(cpu, cpu.regs.Y)
}

func ldy(cpu *m6502, data uint8) {
	cpu.regs.Y = data
	ldzn(cpu, data)
}

func ldx(cpu *m6502, data uint8) {
	cpu.regs.X = data
	ldzn(cpu, data)
}

func lda(cpu *m6502, data uint8) {
	cpu.regs.A = data
	ldzn(cpu, data)
}

func ldzn(cpu *m6502, data uint8) {
	cpu.regs.P.Z = data == 0
	cpu.regs.P.N = data&0x80 != 0
}

func cmp(cpu *m6502, data uint8) {
	r := cpu.regs.A - data
	ldzn(cpu, r)
	cpu.regs.P.C = (cpu.regs.A >= data)
}

func cmpM(cpu *m6502, addr uint16) {
	cmp(cpu, cpu.mem[addr])
}

func cpy(cpu *m6502, data uint8) {
	r := cpu.regs.Y - data
	ldzn(cpu, r)
	cpu.regs.P.C = (cpu.regs.Y >= data)
}

func cpx(cpu *m6502, data uint8) {
	r := cpu.regs.X - data
	ldzn(cpu, r)
	cpu.regs.P.C = (cpu.regs.X >= data)
}

func eor(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A ^ data
	ldzn(cpu, cpu.regs.A)
}

func adc(cpu *m6502, data uint8) {
	r := uint16(cpu.regs.A) + uint16(data)
	if cpu.regs.P.C {
		r++
	}
	cpu.regs.A = uint8(r)
	ldzn(cpu, cpu.regs.A)
	cpu.regs.P.C = r > 0xff
}

func pha(cpu *m6502) {
	cpu.push(cpu.regs.A)
}
