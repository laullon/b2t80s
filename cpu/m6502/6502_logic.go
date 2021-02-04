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
	ops[0x24] = &zeropage{f: bit}
	ops[0x26] = &zeropage{f: rolM}
	ops[0x28] = &implicit{f: plp}
	ops[0x2a] = &implicit{f: rol}
	ops[0x2c] = &absolute{f: bit}
	ops[0x2e] = &absolute{f: rolM}
	ops[0x30] = &relative{f: bmi}
	ops[0x36] = &zeropage{f: rolM, x: true}
	ops[0x38] = &implicit{f: sec}
	ops[0x3e] = &absolute{f: rolM, x: true}
	ops[0x40] = &implicit{f: rti}
	ops[0x48] = &implicit{f: pha}
	ops[0x4c] = &absoluteJMP{}
	ops[0x50] = &relative{f: bvc}
	ops[0x58] = &implicit{f: cli}
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
	ops[0x88] = &implicit{f: dey}
	ops[0x8a] = &implicit{f: txa}
	ops[0x90] = &relative{f: bcc}
	ops[0x98] = &implicit{f: tya}
	ops[0x9a] = &implicit{f: txs}
	ops[0xa8] = &implicit{f: tay}
	ops[0xaa] = &implicit{f: tax}
	ops[0xba] = &implicit{f: tsx}
	ops[0xb0] = &relative{f: bcs}
	ops[0xb8] = &implicit{f: clv}
	ops[0xc6] = &zeropage{f: decM}
	ops[0xce] = &absolute{f: decM}
	ops[0xca] = &implicit{f: dex}
	ops[0xd0] = &relative{f: bne}
	ops[0xd6] = &zeropage{f: decM, x: true}
	ops[0xd8] = &implicit{f: cld}
	ops[0xde] = &absolute{f: decM, x: true}
	ops[0xea] = &implicit{f: nop}
	ops[0xf0] = &relative{f: beq}
	ops[0xf8] = &implicit{f: sed}

	ops[0x4a] = &implicit{f: lsr}
	ops[0x46] = &zeropage{f: lsrM}
	ops[0x56] = &zeropage{f: lsrM, x: true}
	ops[0x4e] = &absolute{f: lsrM}
	ops[0x5e] = &absolute{f: lsrM, x: true}

	ops[0x29] = &immediate{f: and}
	ops[0x25] = &zeropage{f: and}
	ops[0x35] = &zeropage{f: and, x: true}
	ops[0x2d] = &absolute{f: and}
	ops[0x3d] = &absolute{f: and, x: true}
	ops[0x39] = &absolute{f: and, y: true}
	ops[0x21] = &indirectX{f: and}
	ops[0x31] = &indirectY{f: and}

	ops[0x49] = &immediate{f: eor}
	ops[0x45] = &zeropage{f: eor}
	ops[0x55] = &zeropage{f: eor, x: true}
	ops[0x4d] = &absolute{f: eor}
	ops[0x5d] = &absolute{f: eor, x: true}
	ops[0x59] = &absolute{f: eor, y: true}
	ops[0x41] = &indirectX{f: eor}
	ops[0x51] = &indirectY{f: eor}

	ops[0x09] = &immediate{f: ora}
	ops[0x05] = &zeropage{f: ora}
	ops[0x15] = &zeropage{f: ora, x: true}
	ops[0x0d] = &absolute{f: ora}
	ops[0x1d] = &absolute{f: ora, x: true}
	ops[0x19] = &absolute{f: ora, y: true}
	ops[0x01] = &indirectX{f: ora}
	ops[0x11] = &indirectY{f: ora}

	ops[0x69] = &immediate{f: adc}
	ops[0x65] = &zeropage{f: adc}
	ops[0x75] = &zeropage{f: adc, x: true}
	ops[0x6d] = &absolute{f: adc}
	ops[0x7d] = &absolute{f: adc, x: true}
	ops[0x79] = &absolute{f: adc, y: true}
	ops[0x61] = &indirectX{f: adc}
	ops[0x71] = &indirectY{f: adc}

	ops[0xe9] = &immediate{f: sbc}
	ops[0xe5] = &zeropage{f: sbc}
	ops[0xf5] = &zeropage{f: sbc, x: true}
	ops[0xed] = &absolute{f: sbc}
	ops[0xfd] = &absolute{f: sbc, x: true}
	ops[0xf9] = &absolute{f: sbc, y: true}
	ops[0xe1] = &indirectX{f: sbc}
	ops[0xf1] = &indirectY{f: sbc}

	ops[0xc9] = &immediate{f: cmp}
	ops[0xc5] = &zeropage{f: cmp}
	ops[0xd5] = &zeropage{f: cmp, x: true}
	ops[0xcd] = &absolute{f: cmp}
	ops[0xdd] = &absolute{f: cmp, x: true}
	ops[0xd9] = &absolute{f: cmp, y: true}
	ops[0xc1] = &indirectX{f: cmp}
	ops[0xd1] = &indirectY{f: cmp}

	ops[0xa9] = &immediate{f: lda}
	ops[0xa5] = &zeropage{f: lda}
	ops[0xb5] = &zeropage{f: lda, x: true}
	ops[0xad] = &absolute{f: lda}
	ops[0xbd] = &absolute{f: lda, x: true}
	ops[0xb9] = &absolute{f: lda, y: true}
	ops[0xa1] = &indirectX{f: lda}
	ops[0xb1] = &indirectY{f: lda}

	ops[0xa2] = &immediate{f: ldx}
	ops[0xa6] = &zeropage{f: ldx}
	ops[0xb6] = &zeropage{f: ldx, y: true}
	ops[0xae] = &absolute{f: ldx}
	ops[0xbe] = &absolute{f: ldx, y: true}

	ops[0xa0] = &immediate{f: ldy}
	ops[0xa4] = &zeropage{f: ldy}
	ops[0xb4] = &zeropage{f: ldy, x: true}
	ops[0xac] = &absolute{f: ldy}
	ops[0xbc] = &absolute{f: ldy, x: true}

	ops[0x85] = &zeropage{f: sta}
	ops[0x95] = &zeropage{f: sta, x: true}
	ops[0x8d] = &absolute{f: sta}
	ops[0x9d] = &absolute{f: sta, x: true}
	ops[0x99] = &absolute{f: sta, y: true}
	ops[0x81] = &indirectX{f: sta}
	ops[0x91] = &indirectY{f: sta}

	ops[0x86] = &zeropage{f: stx}
	ops[0x96] = &zeropage{f: stx, y: true}
	ops[0x8e] = &absolute{f: stx}

	ops[0x84] = &zeropage{f: sty}
	ops[0x94] = &zeropage{f: sty, x: true}
	ops[0x8c] = &absolute{f: sty}

	ops[0xc0] = &immediate{f: cpy}
	ops[0xc4] = &zeropage{f: cpy}
	ops[0xcc] = &absolute{f: cpy}

	ops[0xe0] = &immediate{f: cpx}
	ops[0xe4] = &zeropage{f: cpx}
	ops[0xec] = &absolute{f: cpx}

	ops[0xc8] = &implicit{f: iny}
	ops[0xe8] = &implicit{f: inx}

	ops[0xe6] = &zeropage{f: incM}
	ops[0xf6] = &zeropage{f: incM, x: true}
	ops[0xee] = &absolute{f: incM}
	ops[0xfe] = &absolute{f: incM, x: true}

	for opCode, op := range ops {
		if op != nil {
			op.setup(uint8(opCode))
		}
	}
}

func rti(cpu *m6502) {
	plp(cpu)
	addr := uint16(cpu.pop())
	addr |= uint16(cpu.pop()) << 8
	cpu.regs.PC = addr
	cpu.preFetch()
}

func bne(cpu *m6502) bool { return !cpu.regs.PS.Z }
func beq(cpu *m6502) bool { return cpu.regs.PS.Z }
func bcc(cpu *m6502) bool { return !cpu.regs.PS.C }
func bcs(cpu *m6502) bool { return cpu.regs.PS.C }
func bpl(cpu *m6502) bool { return !cpu.regs.PS.N }
func bmi(cpu *m6502) bool { return cpu.regs.PS.N }
func bvc(cpu *m6502) bool { return !cpu.regs.PS.V }
func bvs(cpu *m6502) bool { return cpu.regs.PS.V }

func rts(cpu *m6502) {
	addr := uint16(cpu.pop())
	addr |= uint16(cpu.pop()) << 8
	cpu.regs.PC = addr + 1
	cpu.preFetch()
}

func sta(cpu *m6502) uint8 { return cpu.regs.A }
func stx(cpu *m6502) uint8 { return cpu.regs.X }
func sty(cpu *m6502) uint8 { return cpu.regs.Y }

func nop(cpu *m6502) {}

func cld(cpu *m6502) { cpu.regs.PS.D = false }
func sed(cpu *m6502) { cpu.regs.PS.D = true }

func sec(cpu *m6502) { cpu.regs.PS.C = true }
func clc(cpu *m6502) { cpu.regs.PS.C = false }

func sei(cpu *m6502) { cpu.regs.PS.I = true }
func cli(cpu *m6502) { cpu.regs.PS.I = false }

func clv(cpu *m6502) { cpu.regs.PS.V = false }

func tsx(cpu *m6502) {
	cpu.regs.X = cpu.regs.SP
	setZN(cpu, cpu.regs.X)
}

func txs(cpu *m6502) {
	cpu.regs.SP = cpu.regs.X
}

func tax(cpu *m6502) {
	cpu.regs.X = cpu.regs.A
	setZN(cpu, cpu.regs.X)
}

func txa(cpu *m6502) {
	cpu.regs.A = cpu.regs.X
	setZN(cpu, cpu.regs.A)
}

func tay(cpu *m6502) {
	cpu.regs.Y = cpu.regs.A
	setZN(cpu, cpu.regs.Y)
}

func tya(cpu *m6502) {
	cpu.regs.A = cpu.regs.Y
	setZN(cpu, cpu.regs.A)
}

func dex(cpu *m6502) {
	cpu.regs.X--
	setZN(cpu, cpu.regs.X)
}

func inx(cpu *m6502) {
	cpu.regs.X++
	setZN(cpu, cpu.regs.X)
}

func iny(cpu *m6502) {
	cpu.regs.Y++
	setZN(cpu, cpu.regs.Y)
}

func incM(cpu *m6502, data uint8) uint8 {
	data++
	setZN(cpu, data)
	return data
}

func decM(cpu *m6502, data uint8) uint8 {
	data--
	setZN(cpu, data)
	return data
}

func dey(cpu *m6502) {
	cpu.regs.Y--
	setZN(cpu, cpu.regs.Y)
}

func ldy(cpu *m6502, data uint8) {
	cpu.regs.Y = data
	setZN(cpu, data)
}

func ldx(cpu *m6502, data uint8) {
	cpu.regs.X = data
	setZN(cpu, data)
}

func lda(cpu *m6502, data uint8) {
	cpu.regs.A = data
	setZN(cpu, data)
}

func setZN(cpu *m6502, data uint8) {
	cpu.regs.PS.Z = data == 0
	cpu.regs.PS.N = data&0x80 != 0
}

func cmp(cpu *m6502, data uint8) { cp(cpu, cpu.regs.A, data) }
func cpy(cpu *m6502, data uint8) { cp(cpu, cpu.regs.Y, data) }
func cpx(cpu *m6502, data uint8) { cp(cpu, cpu.regs.X, data) }

func cp(cpu *m6502, v, d uint8) {
	r := v - d
	setZN(cpu, r)
	cpu.regs.PS.C = (v >= d)
}

func eor(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A ^ data
	setZN(cpu, cpu.regs.A)
}

func and(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A & data
	setZN(cpu, cpu.regs.A)
}

func ora(cpu *m6502, data uint8) {
	cpu.regs.A = cpu.regs.A | data
	setZN(cpu, cpu.regs.A)
}

func adc(cpu *m6502, data uint8) {
	r := uint16(cpu.regs.A) + uint16(data)
	if cpu.regs.PS.C {
		r++
	}
	cpu.regs.PS.C = r > 0xff
	cpu.regs.PS.V = (((uint16(cpu.regs.A) ^ r) & 0x80) != 0) && (((cpu.regs.A ^ data) & 0x80) == 0)
	cpu.regs.A = uint8(r)
	setZN(cpu, cpu.regs.A)
}

func sbc(cpu *m6502, data uint8) {
	r := uint16(cpu.regs.A) - uint16(data)
	if !cpu.regs.PS.C {
		r--
	}
	cpu.regs.PS.C = int16(r) >= 0
	cpu.regs.PS.V = (((uint16(cpu.regs.A) ^ r) & 0x80) != 0) && (((cpu.regs.A ^ data) & 0x80) != 0)
	cpu.regs.A = uint8(r)
	setZN(cpu, cpu.regs.A)
}

func pha(cpu *m6502) {
	cpu.push(cpu.regs.A)
}

func php(cpu *m6502) {
	cpu.push(cpu.regs.PS.get() | 0b00110000)
}

func pla(cpu *m6502) {
	cpu.regs.A = cpu.pop()
	setZN(cpu, cpu.regs.A)
}

func plp(cpu *m6502) {
	cpu.regs.PS.set(cpu.pop() & 0b11001111)
}

func lsrM(cpu *m6502, data uint8) uint8 {
	cpu.regs.PS.C = data&1 == 1
	data >>= 1
	setZN(cpu, data)
	return data
}

func asl(cpu *m6502) {
	cpu.regs.PS.C = cpu.regs.A&0x80 == 0x80
	cpu.regs.A <<= 1
	setZN(cpu, cpu.regs.A)
}

func aslM(cpu *m6502, data uint8) uint8 {
	cpu.regs.PS.C = data&0x80 == 0x80
	data <<= 1
	setZN(cpu, data)
	return data
}

func rol(cpu *m6502) {
	c := cpu.regs.PS.C
	cpu.regs.PS.C = cpu.regs.A&0x80 == 0x80
	cpu.regs.A <<= 1
	if c {
		cpu.regs.A |= 1
	}
	setZN(cpu, cpu.regs.A)
}

func rolM(cpu *m6502, data uint8) uint8 {
	c := cpu.regs.PS.C
	cpu.regs.PS.C = data&0x80 == 0x80
	data <<= 1
	if c {
		data |= 1
	}
	setZN(cpu, data)
	return data
}

func rorM(cpu *m6502, data uint8) uint8 {
	c := cpu.regs.PS.C
	cpu.regs.PS.C = data&0x01 == 0x01
	data >>= 1
	if c {
		data |= 0x80
	}
	setZN(cpu, data)
	return data
}

func ror(cpu *m6502) {
	c := cpu.regs.PS.C
	cpu.regs.PS.C = cpu.regs.A&0x01 == 0x01
	cpu.regs.A >>= 1
	if c {
		cpu.regs.A |= 0x80
	}
	setZN(cpu, cpu.regs.A)
}

func lsr(cpu *m6502) {
	cpu.regs.PS.C = cpu.regs.A&0x01 == 0x01
	cpu.regs.A >>= 1
	setZN(cpu, cpu.regs.A)
}

func bit(cpu *m6502, data uint8) {
	cpu.regs.PS.Z = (cpu.regs.A & data) == 0
	cpu.regs.PS.V = data&0x40 != 0
	cpu.regs.PS.N = data&0x80 != 0
}
