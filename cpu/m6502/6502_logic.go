package m6502

var ops []operation

func init() {
	ops = make([]operation, 0x100)

	ops[0x00] = &implicit{F: brk}
	ops[0x08] = &implicit{F: php}
	ops[0x09] = &immediate{F: ora}
	ops[0x10] = &relative{F: bpl}
	ops[0x20] = &absolute{F: jsr}
	ops[0x28] = &implicit{F: plp}
	ops[0x18] = &implicit{F: clc}
	ops[0x30] = &relative{F: bmi}
	ops[0x38] = &implicit{F: sec}
	ops[0x40] = &implicit{F: rti}
	ops[0x46] = &zeropage{F: lsrM}
	ops[0x48] = &implicit{F: pha}
	ops[0x49] = &immediate{F: eor}
	ops[0x4c] = &absolute{F: jmp}
	ops[0x50] = &relative{F: bvc}
	ops[0x58] = &implicit{F: cli}
	ops[0x60] = &implicit{F: rts}
	ops[0x68] = &implicit{F: pla}
	ops[0x69] = &immediate{F: adc}
	ops[0x6c] = &indirect{F: jmp}
	ops[0x70] = &relative{F: bvs}
	ops[0x78] = &implicit{F: sei}
	ops[0x81] = &indirectXY{F: staM, x: true}
	ops[0x84] = &zeropage{F: styM}
	ops[0x85] = &zeropage{F: staM}
	ops[0x86] = &zeropage{F: stxM}
	ops[0x88] = &implicit{F: dey}
	ops[0x8a] = &implicit{F: txa}
	ops[0x8c] = &absolute{F: styM}
	ops[0x8d] = &absolute{F: staM}
	ops[0x8e] = &absolute{F: stxM}
	ops[0x90] = &relative{F: bcc}
	ops[0x91] = &indirectXY{F: staM, y: true}
	ops[0x94] = &zeropage{F: styM, x: true}
	ops[0x95] = &zeropage{F: staM, x: true}
	ops[0x96] = &zeropage{F: stxM, y: true}
	ops[0x98] = &implicit{F: tya}
	ops[0x99] = &absolute{F: staM, y: true}
	ops[0x9a] = &implicit{F: txs}
	ops[0x9d] = &absolute{F: staM, x: true}
	ops[0xa0] = &immediate{F: ldy}
	ops[0xa1] = &indirectXY{F: ldaM, x: true}
	ops[0xa2] = &immediate{F: ldx}
	ops[0xa4] = &zeropage{F: ldyM}
	ops[0xa5] = &zeropage{F: ldaM}
	ops[0xa6] = &zeropage{F: ldxM}
	ops[0xa8] = &implicit{F: tay}
	ops[0xa9] = &immediate{F: lda}
	ops[0xaa] = &implicit{F: tax}
	ops[0xac] = &absolute{F: ldyM}
	ops[0xad] = &absolute{F: ldaM}
	ops[0xae] = &absolute{F: ldxM}
	ops[0xb1] = &indirectXY{F: ldaM, y: true}
	ops[0xb4] = &zeropage{F: ldyM, x: true}
	ops[0xb5] = &zeropage{F: ldaM, x: true}
	ops[0xb6] = &zeropage{F: ldxM, y: true}
	ops[0xba] = &implicit{F: tsx}
	ops[0xb0] = &relative{F: bcs}
	ops[0xb8] = &implicit{F: clv}
	ops[0xb9] = &absolute{F: ldaM, y: true}
	ops[0xbc] = &absolute{F: ldyM, x: true}
	ops[0xbe] = &absolute{F: ldxM, y: true}
	ops[0xbd] = &absolute{F: ldaM, x: true}
	ops[0xc0] = &immediate{F: cpy}
	ops[0xc4] = &zeropage{F: cpyM}
	ops[0xc5] = &zeropage{F: cmpM}
	ops[0xc8] = &implicit{F: iny}
	ops[0xc9] = &immediate{F: cmp}
	ops[0xca] = &implicit{F: dex}
	ops[0xcc] = &absolute{F: cpyM}
	ops[0xcd] = &absolute{F: cmpM}
	ops[0xd0] = &relative{F: bne}
	ops[0xd1] = &indirectXY{F: cmpM, y: true}
	ops[0xd5] = &zeropage{F: cmpM, x: true}
	ops[0xd8] = &implicit{F: cld}
	ops[0xd9] = &absolute{F: cmpM, y: true}
	ops[0xdd] = &absolute{F: cmpM, x: true}
	ops[0xe0] = &immediate{F: cpx}
	ops[0xe4] = &zeropage{F: cpxM}
	ops[0xe8] = &implicit{F: inx}
	ops[0xea] = &implicit{F: nop}
	ops[0xec] = &absolute{F: cpxM}
	ops[0xf0] = &relative{F: beq}
	ops[0xf8] = &implicit{F: sed}

	for opCode, op := range ops {
		if op != nil {
			op.setup(uint8(opCode), getFunctionName(op))
		}
	}
}

func rti(cpu *m6502) {
	cpu.regs.PS.set(cpu.pop())
	addr := uint16(cpu.pop())
	addr |= uint16(cpu.pop()) << 8
	cpu.regs.PC = addr
}

func brk(cpu *m6502) {
	cpu.push(uint8((cpu.regs.PC + 1) >> 8))
	cpu.push(uint8((cpu.regs.PC + 1)))
	cpu.regs.PS.B = true
	cpu.regs.PS.X = true
	cpu.push(cpu.regs.PS.get())
	addr := uint16(cpu.mem[0xfffe])
	addr |= uint16(cpu.mem[0xffff]) << 8
	cpu.regs.PC = addr
	cpu.regs.PS.I = true
}

// TODO: review extra cycles
func bne(cpu *m6502, data int8) {
	if !cpu.regs.PS.Z {
		jmpr(cpu, data)
	}
}

// TODO: review extra cycles
func bcc(cpu *m6502, data int8) {
	if !cpu.regs.PS.C {
		jmpr(cpu, data)
	}
}

// TODO: review extra cycles
func beq(cpu *m6502, data int8) {
	if cpu.regs.PS.Z {
		jmpr(cpu, data)
	}
}

// TODO: review extra cycles
func bpl(cpu *m6502, data int8) {
	if !cpu.regs.PS.N {
		jmpr(cpu, data)
	}
}

// TODO: review extra cycles
func bmi(cpu *m6502, data int8) {
	if cpu.regs.PS.N {
		jmpr(cpu, data)
	}
}

func bvc(cpu *m6502, data int8) {
	if !cpu.regs.PS.V {
		jmpr(cpu, data)
	}
}

func bvs(cpu *m6502, data int8) {
	if cpu.regs.PS.V {
		jmpr(cpu, data)
	}
}

// TODO: review extra cycles
func bcs(cpu *m6502, data int8) {
	if cpu.regs.PS.C {
		jmpr(cpu, data)
	}
}

func jmpr(cpu *m6502, data int8) {
	pc := cpu.regs.PC + uint16(data)
	jmp(cpu, pc)
}

func jmp(cpu *m6502, addr uint16) { cpu.regs.PC = addr }

func jsr(cpu *m6502, addr uint16) {
	cpu.push(uint8((cpu.regs.PC - 1) >> 8))
	cpu.push(uint8((cpu.regs.PC - 1)))
	cpu.regs.PC = addr
}

func rts(cpu *m6502) {
	addr := uint16(cpu.pop())
	addr |= uint16(cpu.pop()) << 8
	cpu.regs.PC = addr + 1
}

func staM(cpu *m6502, addr uint16) { cpu.mem[addr] = cpu.regs.A }

func stxM(cpu *m6502, addr uint16) { cpu.mem[addr] = cpu.regs.X }

func styM(cpu *m6502, addr uint16) { cpu.mem[addr] = cpu.regs.Y }

func ldaM(cpu *m6502, addr uint16) {
	cpu.regs.A = cpu.mem[addr]
	ldzn(cpu, cpu.regs.A)
}

func ldxM(cpu *m6502, addr uint16) {
	cpu.regs.X = cpu.mem[addr]
	ldzn(cpu, cpu.regs.X)
}

func ldyM(cpu *m6502, addr uint16) {
	cpu.regs.Y = cpu.mem[addr]
	ldzn(cpu, cpu.regs.Y)
}

func nop(cpu *m6502) {}

func cld(cpu *m6502) { cpu.regs.PS.D = false }

func clc(cpu *m6502) { cpu.regs.PS.C = false }

func cli(cpu *m6502) { cpu.regs.PS.I = false }

func clv(cpu *m6502) { cpu.regs.PS.V = false }

func sec(cpu *m6502) { cpu.regs.PS.C = true }

func sei(cpu *m6502) { cpu.regs.PS.I = true }

func sed(cpu *m6502) { cpu.regs.PS.D = true }

func tsx(cpu *m6502) {
	cpu.regs.X = cpu.regs.SP
	ldzn(cpu, cpu.regs.X)
}

func txs(cpu *m6502) {
	cpu.regs.SP = cpu.regs.X
}

func tax(cpu *m6502) {
	cpu.regs.X = cpu.regs.A
	ldzn(cpu, cpu.regs.X)
}

func txa(cpu *m6502) {
	cpu.regs.A = cpu.regs.X
	ldzn(cpu, cpu.regs.A)
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

func iny(cpu *m6502) {
	cpu.regs.Y++
	ldzn(cpu, cpu.regs.Y)
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
	cpu.regs.PS.Z = data == 0
	cpu.regs.PS.N = data&0x80 != 0
}

func cmp(cpu *m6502, data uint8) {
	r := cpu.regs.A - data
	ldzn(cpu, r)
	cpu.regs.PS.C = (cpu.regs.A >= data)
}

func cmpM(cpu *m6502, addr uint16) {
	cmp(cpu, cpu.mem[addr])
}

func cpy(cpu *m6502, data uint8) {
	r := cpu.regs.Y - data
	ldzn(cpu, r)
	cpu.regs.PS.C = (cpu.regs.Y >= data)
}

func cpx(cpu *m6502, data uint8) {
	r := cpu.regs.X - data
	ldzn(cpu, r)
	cpu.regs.PS.C = (cpu.regs.X >= data)
}

func cpxM(cpu *m6502, addr uint16) {
	cpx(cpu, cpu.mem[addr])
}

func cpyM(cpu *m6502, addr uint16) {
	cpy(cpu, cpu.mem[addr])
}

func eor(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A ^ data
	ldzn(cpu, cpu.regs.A)
}

func ora(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A | data
	ldzn(cpu, cpu.regs.A)
}

func adc(cpu *m6502, data uint8) {
	r := uint16(cpu.regs.A) + uint16(data)
	if cpu.regs.PS.C {
		r++
	}
	cpu.regs.A = uint8(r)
	ldzn(cpu, cpu.regs.A)
	cpu.regs.PS.C = r > 0xff
}

func pha(cpu *m6502) {
	cpu.push(cpu.regs.A)
}

func php(cpu *m6502) {
	cpu.regs.PS.B = true
	cpu.regs.PS.X = true
	cpu.push(cpu.regs.PS.get())
}

func pla(cpu *m6502) {
	cpu.regs.A = cpu.pop()
	ldzn(cpu, cpu.regs.A)
}

func plp(cpu *m6502) {
	cpu.regs.PS.set(cpu.pop())
}

func lsrM(cpu *m6502, addr uint16) {
	d := cpu.mem[int(addr)]
	cpu.regs.PS.C = d&1 == 1
	d >>= 1
	ldzn(cpu, d)
	cpu.mem[int(addr)] = d
}
