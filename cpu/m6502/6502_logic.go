package m6502

var ops []operation

func init() {
	ops = make([]operation, 0x100)

	ops[0x00] = &brk{}
	ops[0x06] = &zeropage{F: aslM}
	ops[0x08] = &implicit{F: php}
	ops[0x0a] = &implicit{F: asl}
	ops[0x0e] = &absolute{F: aslM}
	ops[0x10] = &relative{F: bpl}
	ops[0x16] = &zeropage{F: aslM, x: true}
	ops[0x18] = &implicit{F: clc}
	ops[0x1e] = &absolute{F: aslM, x: true}
	ops[0x20] = &absoluteJSR{}
	ops[0x24] = &zeropage{F: bitM}
	ops[0x26] = &zeropage{F: rolM}
	ops[0x28] = &implicit{F: plp}
	ops[0x2a] = &implicit{F: rol}
	ops[0x2c] = &absolute{F: bitM}
	ops[0x2e] = &absolute{F: rolM}
	ops[0x30] = &relative{F: bmi}
	ops[0x36] = &zeropage{F: rolM, x: true}
	ops[0x38] = &implicit{F: sec}
	ops[0x3e] = &absolute{F: rolM, x: true}
	ops[0x40] = &implicit{F: rti}
	ops[0x46] = &zeropage{F: lsrM}
	ops[0x48] = &implicit{F: pha}
	ops[0x4a] = &implicit{F: lsr}
	ops[0x4c] = &absoluteJMP{}
	ops[0x4e] = &absolute{F: lsrM}
	ops[0x50] = &relative{F: bvc}
	ops[0x56] = &zeropage{F: lsrM, x: true}
	ops[0x58] = &implicit{F: cli}
	ops[0x5e] = &absolute{F: lsrM, x: true}
	ops[0x60] = &implicit{F: rts}
	ops[0x66] = &zeropage{F: rorM}
	ops[0x68] = &implicit{F: pla}
	ops[0x6a] = &implicit{F: ror}
	ops[0x6c] = &indirectJMP{}
	ops[0x6e] = &absolute{F: rorM}
	ops[0x70] = &relative{F: bvs}
	ops[0x76] = &zeropage{F: rorM, x: true}
	ops[0x78] = &implicit{F: sei}
	ops[0x7e] = &absolute{F: rorM, x: true}
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
	ops[0xc1] = &indirectXY{F: cmpM, x: true}
	ops[0xc4] = &zeropage{F: cpyM}
	ops[0xc5] = &zeropage{F: cmpM}
	ops[0xc6] = &zeropage{F: decM}
	ops[0xc8] = &implicit{F: iny}
	ops[0xc9] = &immediate{F: cmp}
	ops[0xce] = &absolute{F: decM}
	ops[0xca] = &implicit{F: dex}
	ops[0xcc] = &absolute{F: cpyM}
	ops[0xcd] = &absolute{F: cmpM}
	ops[0xd0] = &relative{F: bne}
	ops[0xd1] = &indirectXY{F: cmpM, y: true}
	ops[0xd5] = &zeropage{F: cmpM, x: true}
	ops[0xd6] = &zeropage{F: decM, x: true}
	ops[0xd8] = &implicit{F: cld}
	ops[0xd9] = &absolute{F: cmpM, y: true}
	ops[0xdd] = &absolute{F: cmpM, x: true}
	ops[0xde] = &absolute{F: decM, x: true}
	ops[0xe0] = &immediate{F: cpx}
	ops[0xe4] = &zeropage{F: cpxM}
	ops[0xe6] = &zeropage{F: incM}
	ops[0xe8] = &implicit{F: inx}
	ops[0xea] = &implicit{F: nop}
	ops[0xec] = &absolute{F: cpxM}
	ops[0xee] = &absolute{F: incM}
	ops[0xf0] = &relative{F: beq}
	ops[0xf6] = &zeropage{F: incM, x: true}
	ops[0xf8] = &implicit{F: sed}
	ops[0xfe] = &absolute{F: incM, x: true}

	ops[0x29] = &immediate{F: and}
	ops[0x25] = &zeropage{F: andM}
	ops[0x35] = &zeropage{F: andM, x: true}
	ops[0x2d] = &absolute{F: andM}
	ops[0x3d] = &absolute{F: andM, x: true}
	ops[0x39] = &absolute{F: andM, y: true}
	ops[0x21] = &indirectXY{F: andM, x: true}
	ops[0x31] = &indirectXY{F: andM, y: true}

	ops[0x49] = &immediate{F: eor}
	ops[0x45] = &zeropage{F: eorM}
	ops[0x55] = &zeropage{F: eorM, x: true}
	ops[0x4d] = &absolute{F: eorM}
	ops[0x5d] = &absolute{F: eorM, x: true}
	ops[0x59] = &absolute{F: eorM, y: true}
	ops[0x41] = &indirectXY{F: eorM, x: true}
	ops[0x51] = &indirectXY{F: eorM, y: true}

	ops[0x09] = &immediate{F: ora}
	ops[0x05] = &zeropage{F: oraM}
	ops[0x15] = &zeropage{F: oraM, x: true}
	ops[0x0d] = &absolute{F: oraM}
	ops[0x1d] = &absolute{F: oraM, x: true}
	ops[0x19] = &absolute{F: oraM, y: true}
	ops[0x01] = &indirectXY{F: oraM, x: true}
	ops[0x11] = &indirectXY{F: oraM, y: true}

	ops[0x69] = &immediate{F: adc}
	ops[0x65] = &zeropage{F: adcM}
	ops[0x75] = &zeropage{F: adcM, x: true}
	ops[0x6d] = &absolute{F: adcM}
	ops[0x7d] = &absolute{F: adcM, x: true}
	ops[0x79] = &absolute{F: adcM, y: true}
	ops[0x61] = &indirectXY{F: adcM, x: true}
	ops[0x71] = &indirectXY{F: adcM, y: true}

	ops[0xe9] = &immediate{F: sbc}
	ops[0xe5] = &zeropage{F: sbcM}
	ops[0xf5] = &zeropage{F: sbcM, x: true}
	ops[0xed] = &absolute{F: sbcM}
	ops[0xfd] = &absolute{F: sbcM, x: true}
	ops[0xf9] = &absolute{F: sbcM, y: true}
	ops[0xe1] = &indirectXY{F: sbcM, x: true}
	ops[0xf1] = &indirectXY{F: sbcM, y: true}

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

func bne(cpu *m6502) bool {
	return !cpu.regs.PS.Z
}

func bcc(cpu *m6502) bool {
	return !cpu.regs.PS.C
}

func beq(cpu *m6502) bool {
	return cpu.regs.PS.Z
}

func bpl(cpu *m6502) bool {
	return !cpu.regs.PS.N
}

func bmi(cpu *m6502) bool {
	return cpu.regs.PS.N
}
func bvc(cpu *m6502) bool {
	return !cpu.regs.PS.V
}
func bvs(cpu *m6502) bool {
	return cpu.regs.PS.V
}

func bcs(cpu *m6502) bool {
	return cpu.regs.PS.C
}

func rts(cpu *m6502) {
	addr := uint16(cpu.pop())
	addr |= uint16(cpu.pop()) << 8
	cpu.regs.PC = addr + 1
}

func staM(cpu *m6502, data uint8) (discard bool, v uint8) { return false, cpu.regs.A }

func stxM(cpu *m6502, data uint8) (discard bool, v uint8) { return false, cpu.regs.X }

func styM(cpu *m6502, data uint8) (discard bool, v uint8) { return false, cpu.regs.Y }

func ldaM(cpu *m6502, data uint8) (discard bool, v uint8) {
	cpu.regs.A = data
	ldzn(cpu, cpu.regs.A)
	return true, 0
}

func ldxM(cpu *m6502, data uint8) (discard bool, v uint8) {
	cpu.regs.X = data
	ldzn(cpu, cpu.regs.X)
	return true, 0
}

func ldyM(cpu *m6502, data uint8) (discard bool, v uint8) {
	cpu.regs.Y = data
	ldzn(cpu, cpu.regs.Y)
	return true, 0
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

func incM(cpu *m6502, data uint8) (discard bool, v uint8) {
	data++
	ldzn(cpu, data)
	return false, data
}

func decM(cpu *m6502, data uint8) (discard bool, v uint8) {
	data--
	ldzn(cpu, data)
	return false, data
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

func cmpM(cpu *m6502, data uint8) (discard bool, v uint8) {
	cmp(cpu, data)
	return true, 0
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

func cpxM(cpu *m6502, data uint8) (discard bool, v uint8) {
	cpx(cpu, data)
	return true, 0
}

func cpyM(cpu *m6502, data uint8) (discard bool, v uint8) {
	cpy(cpu, data)
	return true, 0
}

func eor(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A ^ data
	ldzn(cpu, cpu.regs.A)
}

func eorM(cpu *m6502, data uint8) (discard bool, v uint8) {
	cpu.regs.A = cpu.regs.A ^ data
	ldzn(cpu, cpu.regs.A)
	return true, 0
}

func and(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A & data
	ldzn(cpu, cpu.regs.A)
}

func andM(cpu *m6502, data uint8) (discard bool, v uint8) {
	cpu.regs.A = cpu.regs.A & data
	ldzn(cpu, cpu.regs.A)
	return true, 0
}

func ora(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A | data
	ldzn(cpu, cpu.regs.A)
}

func oraM(cpu *m6502, data uint8) (discard bool, v uint8) {
	cpu.regs.A = cpu.regs.A | data
	ldzn(cpu, cpu.regs.A)
	return true, 0
}

func adc(cpu *m6502, data uint8) {
	r := uint16(cpu.regs.A) + uint16(data)
	if cpu.regs.PS.C {
		r++
	}
	cpu.regs.PS.C = r > 0xff
	cpu.regs.PS.V = (((uint16(cpu.regs.A) ^ r) & 0x80) != 0) && (((cpu.regs.A ^ data) & 0x80) == 0)
	cpu.regs.A = uint8(r)
	ldzn(cpu, cpu.regs.A)
}

func adcM(cpu *m6502, data uint8) (discard bool, v uint8) {
	adc(cpu, data)
	return true, 0
}

func sbc(cpu *m6502, data uint8) {
	r := uint16(cpu.regs.A) - uint16(data)
	if !cpu.regs.PS.C {
		r--
	}
	cpu.regs.PS.C = int16(r) >= 0
	cpu.regs.PS.V = (((uint16(cpu.regs.A) ^ r) & 0x80) != 0) && (((cpu.regs.A ^ data) & 0x80) != 0)
	cpu.regs.A = uint8(r)
	ldzn(cpu, cpu.regs.A)
}

func sbcM(cpu *m6502, data uint8) (discard bool, v uint8) {
	sbc(cpu, data)
	return true, 0
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

func lsrM(cpu *m6502, data uint8) (discard bool, v uint8) {
	cpu.regs.PS.C = data&1 == 1
	data >>= 1
	ldzn(cpu, data)
	return false, data
}

func asl(cpu *m6502) {
	cpu.regs.PS.C = cpu.regs.A&0x80 == 0x80
	cpu.regs.A <<= 1
	ldzn(cpu, cpu.regs.A)
}

func aslM(cpu *m6502, data uint8) (discard bool, v uint8) {
	cpu.regs.PS.C = data&0x80 == 0x80
	data <<= 1
	ldzn(cpu, data)
	return false, data
}

func rol(cpu *m6502) {
	c := cpu.regs.PS.C
	cpu.regs.PS.C = cpu.regs.A&0x80 == 0x80
	cpu.regs.A <<= 1
	if c {
		cpu.regs.A |= 1
	}
	ldzn(cpu, cpu.regs.A)
}

func rolM(cpu *m6502, data uint8) (discard bool, v uint8) {
	c := cpu.regs.PS.C
	cpu.regs.PS.C = data&0x80 == 0x80
	data <<= 1
	if c {
		data |= 1
	}
	ldzn(cpu, data)
	return false, data
}

func rorM(cpu *m6502, data uint8) (discard bool, v uint8) {
	c := cpu.regs.PS.C
	cpu.regs.PS.C = data&0x01 == 0x01
	data >>= 1
	if c {
		data |= 0x80
	}
	ldzn(cpu, data)
	return false, data
}

func ror(cpu *m6502) {
	c := cpu.regs.PS.C
	cpu.regs.PS.C = cpu.regs.A&0x01 == 0x01
	cpu.regs.A >>= 1
	if c {
		cpu.regs.A |= 0x80
	}
	ldzn(cpu, cpu.regs.A)
}

func lsr(cpu *m6502) {
	cpu.regs.PS.C = cpu.regs.A&0x01 == 0x01
	cpu.regs.A >>= 1
	ldzn(cpu, cpu.regs.A)
}

func bitM(cpu *m6502, data uint8) (discard bool, v uint8) {
	cpu.regs.PS.Z = (cpu.regs.A & data) == 0
	cpu.regs.PS.V = data&0x40 != 0
	cpu.regs.PS.N = data&0x80 != 0
	return true, 0
}
