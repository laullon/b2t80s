package m6502

var ops []operation

func init() {
	ops = make([]operation, 0x100)

	ops[0x00] = &brk{}
	ops[0x06] = &zeropage{f: aslM}
	ops[0x08] = &implicit{f: php}
	ops[0x0a] = &implicit{f: asl}
	ops[0x0e] = &absolute{f: aslM}
	ops[0x10] = &relative{f: bpl}
	ops[0x16] = &zeropage{f: aslM, x: true}
	ops[0x18] = &implicit{f: clc}
	ops[0x1e] = &absolute{f: aslM, x: true}
	ops[0x20] = &absoluteJSR{}
	ops[0x24] = &zeropage{f: bitM}
	ops[0x26] = &zeropage{f: rolM}
	ops[0x28] = &implicit{f: plp}
	ops[0x2a] = &implicit{f: rol}
	ops[0x2c] = &absolute{f: bitM}
	ops[0x2e] = &absolute{f: rolM}
	ops[0x30] = &relative{f: bmi}
	ops[0x36] = &zeropage{f: rolM, x: true}
	ops[0x38] = &implicit{f: sec}
	ops[0x3e] = &absolute{f: rolM, x: true}
	ops[0x40] = &implicit{f: rti}
	ops[0x46] = &zeropage{f: lsrM}
	ops[0x48] = &implicit{f: pha}
	ops[0x4a] = &implicit{f: lsr}
	ops[0x4c] = &absoluteJMP{}
	ops[0x4e] = &absolute{f: lsrM}
	ops[0x50] = &relative{f: bvc}
	ops[0x56] = &zeropage{f: lsrM, x: true}
	ops[0x58] = &implicit{f: cli}
	ops[0x5e] = &absolute{f: lsrM, x: true}
	ops[0x60] = &implicit{f: rts}
	ops[0x66] = &zeropage{f: rorM}
	ops[0x68] = &implicit{f: pla}
	ops[0x6a] = &implicit{f: ror}
	ops[0x6c] = &indirectJMP{}
	ops[0x6e] = &absolute{f: rorM}
	ops[0x70] = &relative{f: bvs}
	ops[0x76] = &zeropage{f: rorM, x: true}
	ops[0x78] = &implicit{f: sei}
	ops[0x7e] = &absolute{f: rorM, x: true}
	ops[0x81] = &indirectXY{f: staM, x: true}
	ops[0x84] = &zeropage{f: styM}
	ops[0x85] = &zeropage{f: staM}
	ops[0x86] = &zeropage{f: stxM}
	ops[0x88] = &implicit{f: dey}
	ops[0x8a] = &implicit{f: txa}
	ops[0x8c] = &absolute{f: styM}
	ops[0x8d] = &absolute{f: staM}
	ops[0x8e] = &absolute{f: stxM}
	ops[0x90] = &relative{f: bcc}
	ops[0x91] = &indirectXY{f: staM, y: true}
	ops[0x94] = &zeropage{f: styM, x: true}
	ops[0x95] = &zeropage{f: staM, x: true}
	ops[0x96] = &zeropage{f: stxM, y: true}
	ops[0x98] = &implicit{f: tya}
	ops[0x99] = &absolute{f: staM, y: true}
	ops[0x9a] = &implicit{f: txs}
	ops[0x9d] = &absolute{f: staM, x: true}
	ops[0xa0] = &immediate{f: ldy}
	ops[0xa1] = &indirectXY{f: ldaM, x: true}
	ops[0xa2] = &immediate{f: ldx}
	ops[0xa4] = &zeropage{f: ldyM}
	ops[0xa5] = &zeropage{f: ldaM}
	ops[0xa6] = &zeropage{f: ldxM}
	ops[0xa8] = &implicit{f: tay}
	ops[0xa9] = &immediate{f: lda}
	ops[0xaa] = &implicit{f: tax}
	ops[0xac] = &absolute{f: ldyM}
	ops[0xad] = &absolute{f: ldaM}
	ops[0xae] = &absolute{f: ldxM}
	ops[0xb1] = &indirectXY{f: ldaM, y: true}
	ops[0xb4] = &zeropage{f: ldyM, x: true}
	ops[0xb5] = &zeropage{f: ldaM, x: true}
	ops[0xb6] = &zeropage{f: ldxM, y: true}
	ops[0xba] = &implicit{f: tsx}
	ops[0xb0] = &relative{f: bcs}
	ops[0xb8] = &implicit{f: clv}
	ops[0xb9] = &absolute{f: ldaM, y: true}
	ops[0xbc] = &absolute{f: ldyM, x: true}
	ops[0xbe] = &absolute{f: ldxM, y: true}
	ops[0xbd] = &absolute{f: ldaM, x: true}
	ops[0xc0] = &immediate{f: cpy}
	ops[0xc1] = &indirectXY{f: cmpM, x: true}
	ops[0xc4] = &zeropage{f: cpyM}
	ops[0xc5] = &zeropage{f: cmpM}
	ops[0xc6] = &zeropage{f: decM}
	ops[0xc8] = &implicit{f: iny}
	ops[0xc9] = &immediate{f: cmp}
	ops[0xce] = &absolute{f: decM}
	ops[0xca] = &implicit{f: dex}
	ops[0xcc] = &absolute{f: cpyM}
	ops[0xcd] = &absolute{f: cmpM}
	ops[0xd0] = &relative{f: bne}
	ops[0xd1] = &indirectXY{f: cmpM, y: true}
	ops[0xd5] = &zeropage{f: cmpM, x: true}
	ops[0xd6] = &zeropage{f: decM, x: true}
	ops[0xd8] = &implicit{f: cld}
	ops[0xd9] = &absolute{f: cmpM, y: true}
	ops[0xdd] = &absolute{f: cmpM, x: true}
	ops[0xde] = &absolute{f: decM, x: true}
	ops[0xe0] = &immediate{f: cpx}
	ops[0xe4] = &zeropage{f: cpxM}
	ops[0xe6] = &zeropage{f: incM}
	ops[0xe8] = &implicit{f: inx}
	ops[0xea] = &implicit{f: nop}
	ops[0xec] = &absolute{f: cpxM}
	ops[0xee] = &absolute{f: incM}
	ops[0xf0] = &relative{f: beq}
	ops[0xf6] = &zeropage{f: incM, x: true}
	ops[0xf8] = &implicit{f: sed}
	ops[0xfe] = &absolute{f: incM, x: true}

	ops[0x29] = &immediate{f: and}
	ops[0x25] = &zeropage{f: andM}
	ops[0x35] = &zeropage{f: andM, x: true}
	ops[0x2d] = &absolute{f: andM}
	ops[0x3d] = &absolute{f: andM, x: true}
	ops[0x39] = &absolute{f: andM, y: true}
	ops[0x21] = &indirectXY{f: andM, x: true}
	ops[0x31] = &indirectXY{f: andM, y: true}

	ops[0x49] = &immediate{f: eor}
	ops[0x45] = &zeropage{f: eorM}
	ops[0x55] = &zeropage{f: eorM, x: true}
	ops[0x4d] = &absolute{f: eorM}
	ops[0x5d] = &absolute{f: eorM, x: true}
	ops[0x59] = &absolute{f: eorM, y: true}
	ops[0x41] = &indirectXY{f: eorM, x: true}
	ops[0x51] = &indirectXY{f: eorM, y: true}

	ops[0x09] = &immediate{f: ora}
	ops[0x05] = &zeropage{f: oraM}
	ops[0x15] = &zeropage{f: oraM, x: true}
	ops[0x0d] = &absolute{f: oraM}
	ops[0x1d] = &absolute{f: oraM, x: true}
	ops[0x19] = &absolute{f: oraM, y: true}
	ops[0x01] = &indirectXY{f: oraM, x: true}
	ops[0x11] = &indirectXY{f: oraM, y: true}

	ops[0x69] = &immediate{f: adc}
	ops[0x65] = &zeropage{f: adcM}
	ops[0x75] = &zeropage{f: adcM, x: true}
	ops[0x6d] = &absolute{f: adcM}
	ops[0x7d] = &absolute{f: adcM, x: true}
	ops[0x79] = &absolute{f: adcM, y: true}
	ops[0x61] = &indirectXY{f: adcM, x: true}
	ops[0x71] = &indirectXY{f: adcM, y: true}

	ops[0xe9] = &immediate{f: sbc}
	ops[0xe5] = &zeropage{f: sbcM}
	ops[0xf5] = &zeropage{f: sbcM, x: true}
	ops[0xed] = &absolute{f: sbcM}
	ops[0xfd] = &absolute{f: sbcM, x: true}
	ops[0xf9] = &absolute{f: sbcM, y: true}
	ops[0xe1] = &indirectXY{f: sbcM, x: true}
	ops[0xf1] = &indirectXY{f: sbcM, y: true}

	for opCode, op := range ops {
		if op != nil {
			op.setup(uint8(opCode))
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

func staM(cpu *m6502) uint8 { return cpu.regs.A }

func stxM(cpu *m6502) uint8 { return cpu.regs.X }

func styM(cpu *m6502) uint8 { return cpu.regs.Y }

func ldaM(cpu *m6502, data uint8) {
	cpu.regs.A = data
	ldzn(cpu, cpu.regs.A)
}

func ldxM(cpu *m6502, data uint8) {
	cpu.regs.X = data
	ldzn(cpu, cpu.regs.X)
}

func ldyM(cpu *m6502, data uint8) {
	cpu.regs.Y = data
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

func incM(cpu *m6502, data uint8) uint8 {
	data++
	ldzn(cpu, data)
	return data
}

func decM(cpu *m6502, data uint8) uint8 {
	data--
	ldzn(cpu, data)
	return data
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

func cmpM(cpu *m6502, data uint8) {
	cmp(cpu, data)
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

func cpxM(cpu *m6502, data uint8) {
	cpx(cpu, data)
}

func cpyM(cpu *m6502, data uint8) {
	cpy(cpu, data)
}

func eor(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A ^ data
	ldzn(cpu, cpu.regs.A)
}

func eorM(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A ^ data
	ldzn(cpu, cpu.regs.A)
}

func and(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A & data
	ldzn(cpu, cpu.regs.A)
}

func andM(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A & data
	ldzn(cpu, cpu.regs.A)
}

func ora(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A | data
	ldzn(cpu, cpu.regs.A)
}

func oraM(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A | data
	ldzn(cpu, cpu.regs.A)
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

func adcM(cpu *m6502, data uint8) {
	adc(cpu, data)
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

func sbcM(cpu *m6502, data uint8) {
	sbc(cpu, data)
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

func lsrM(cpu *m6502, data uint8) uint8 {
	cpu.regs.PS.C = data&1 == 1
	data >>= 1
	ldzn(cpu, data)
	return data
}

func asl(cpu *m6502) {
	cpu.regs.PS.C = cpu.regs.A&0x80 == 0x80
	cpu.regs.A <<= 1
	ldzn(cpu, cpu.regs.A)
}

func aslM(cpu *m6502, data uint8) uint8 {
	cpu.regs.PS.C = data&0x80 == 0x80
	data <<= 1
	ldzn(cpu, data)
	return data
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

func rolM(cpu *m6502, data uint8) uint8 {
	c := cpu.regs.PS.C
	cpu.regs.PS.C = data&0x80 == 0x80
	data <<= 1
	if c {
		data |= 1
	}
	ldzn(cpu, data)
	return data
}

func rorM(cpu *m6502, data uint8) uint8 {
	c := cpu.regs.PS.C
	cpu.regs.PS.C = data&0x01 == 0x01
	data >>= 1
	if c {
		data |= 0x80
	}
	ldzn(cpu, data)
	return data
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

func bitM(cpu *m6502, data uint8) {
	cpu.regs.PS.Z = (cpu.regs.A & data) == 0
	cpu.regs.PS.V = data&0x40 != 0
	cpu.regs.PS.N = data&0x80 != 0
}
