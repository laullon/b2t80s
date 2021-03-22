package lr35902

import "fmt"

type opCode struct {
	name       string
	mask, code byte
	len        byte
	f          lr35902f
}

var lr35902OpsCodeTable = []*opCode{
	{"ADD HL,ss", 0b11001111, 0b00001001, 1, addHLss},
	{"INC ss", 0b11001111, 0b00000011, 1, incSS},
	{"DEC ss", 0b11001111, 0b00001011, 1, decSS},
	{"POP ss", 0b11001111, 0b11000001, 1, popSS},
	{"PUSH ss", 0b11001111, 0b11000101, 1, pushSS},
	{"LD r, r'", 0b11000000, 0b01000000, 1, ldRr},
	{"LD r, (HL)", 0b11000111, 0b01000110, 1, ldRhl},
	{"LD (HL), r", 0b11111000, 0b01110000, 1, ldHLr},
	{"INC r", 0b11000111, 0b0000100, 1, incR},
	{"DEC r", 0b11000111, 0b0000101, 1, decR},
	{"ADD A, r", 0b11111000, 0b10000000, 1, addAr},
	{"ADC A, r", 0b11111000, 0b10001000, 1, adcAr},
	{"SUB A, r", 0b11111000, 0b10010000, 1, subAr},
	{"SUC A, r", 0b11111000, 0b10011000, 1, sbcAr},
	{"AND r", 0b11111000, 0b10100000, 1, andAr},
	{"OR r", 0b11111000, 0b10110000, 1, orAr},
	{"XOR r", 0b11111000, 0b10101000, 1, xorAr},
	{"CP r", 0b11111000, 0b10111000, 1, cpR},
	{"RET cc", 0b11100111, 0b11000000, 1, retCC},
	{"LD (ff00+C), A", 0xFF, 0xe2, 1, ldhCa},
	{"LDI (HL),a", 0xFF, 0x22, 1, ldiHLa},
	{"LDI A,(HL)", 0xFF, 0x2a, 1, ldiAhl},
	{"LDD A,(HL)", 0xFF, 0x32, 1, lddAhl},
	{"LDD (HL),a", 0xFF, 0x3a, 1, lddHLa},
	{"RST p", 0b11000111, 0b11000111, 1, rstP},
	{"NOP", 0xFF, 0x00, 1, func(cpu *lr35902) {}},
	{"DAA", 0xFF, 0x27, 1, daa},
	{"CPL", 0xFF, 0x2f, 1, cpl},
	{"SCF", 0xFF, 0x37, 1, scf},
	{"CCF", 0xFF, 0x3F, 1, ccf},
	{"HALT", 0xFF, 0x76, 1, halt},
	{"RET", 0xFF, 0xC9, 1, ret},
	{"RETI", 0xFF, 0xD9, 1, reti},
	{"INC (HL)", 0xFF, 0x34, 1, incHL},
	{"DEC (HL)", 0xFF, 0x35, 1, decHL},
	{"ADD A, (HL)", 0xFF, 0x86, 1, addAhl},
	{"ADC A, (HL)", 0xFF, 0x8e, 1, adcAhl},
	{"SUB A, (HL)", 0xFF, 0x96, 1, subAhl},
	{"SBC A, (HL)", 0xFF, 0x9e, 1, sbcAhl},
	{"AND (HL)", 0xFF, 0xA6, 1, andAhl},
	{"OR (HL)", 0xFF, 0xB6, 1, orAhl},
	{"XOR (HL)", 0xFF, 0xAE, 1, xorAhl},
	{"CP (HL)", 0xFF, 0xBE, 1, cpHl},
	{"LD A,(BC)", 0xFF, 0x0A, 1, ldAbc},
	{"LD A,(DE)", 0xFF, 0x1A, 1, ldAde},
	{"LD (BC), A", 0xFF, 0x02, 1, ldBCa},
	{"LD (DE), A", 0xFF, 0x12, 1, ldDEa},
	{"RLCA", 0xFF, 0x07, 1, rlca},
	{"RLA", 0xFF, 0x17, 1, rla},
	{"RRCA", 0xFF, 0x0F, 1, rrca},
	{"RRA", 0xFF, 0x1F, 1, rra},
	{"JP HL", 0xFF, 0xE9, 1, func(cpu *lr35902) { cpu.regs.PC = cpu.regs.HL.Get() }},
	{"DI", 0xFF, 0xF3, 1, func(cpu *lr35902) { cpu.regs.IME = false }},
	{"EI", 0xFF, 0xFb, 1, func(cpu *lr35902) { cpu.regs.IME = true }},
	{"LD SP, HL", 0xFF, 0xF9, 1, func(cpu *lr35902) { cpu.regs.SP.Set(cpu.regs.HL.Get()) }},
	{"CB", 0xFF, 0xCB, 1, decodeCB},

	{"LD r, n", 0b11000111, 0b00000110, 2, ldRn},
	{"LD (ff00+n), A", 0xFF, 0xe0, 2, ldhNa},
	{"LD A, (ff00+n)", 0xFF, 0xf0, 2, ldhAn},
	{"LD HL,(SP+e)", 0xFF, 0xF8, 2, ldHLspE},
	{"ADD A, n", 0xFF, 0xc6, 2, func(cpu *lr35902) { cpu.addA(cpu.fetched.n) }},
	{"ADC A, n", 0xFF, 0xCE, 2, func(cpu *lr35902) { cpu.adcA(cpu.fetched.n) }},
	{"SBC A, (HL)", 0xFF, 0xDE, 2, func(cpu *lr35902) { cpu.sbcA(cpu.fetched.n) }},
	{"SUB n", 0xFF, 0xD6, 2, func(cpu *lr35902) { cpu.subA(cpu.fetched.n) }},
	{"AND n", 0xFF, 0xE6, 2, func(cpu *lr35902) { cpu.and(cpu.fetched.n) }},
	{"OR n", 0xFF, 0xF6, 2, func(cpu *lr35902) { cpu.or(cpu.fetched.n) }},
	{"LD (HL), n", 0xFF, 0x36, 2, ldHLn},
	{"JR e", 0xFF, 0x18, 2, jr},
	{"JRNZ e", 0xFF, 0x20, 2, jrnz},
	{"JR Z, e", 0xFF, 0x28, 2, jrz},
	{"JRNC e", 0xFF, 0x30, 2, jrnc},
	{"JRC e", 0xFF, 0x38, 2, jrc},
	{"XOR *", 0xFF, 0xEE, 2, func(cpu *lr35902) { cpu.xor(cpu.fetched.n) }},
	{"CP n", 0xFF, 0xFe, 2, func(cpu *lr35902) { cpu.cp(cpu.fetched.n) }},

	{"LD dd, mm", 0b11001111, 0b00000001, 3, ldDDmm},
	{"JP cc, nn", 0b11100111, 0b11000010, 3, jpCC},
	{"CALL cc, nn", 0b11100111, 0b11000100, 3, callCC},
	{"LD (nn), A", 0xFF, 0xea, 3, ldNNa},
	{"LD A, (nn)", 0xFF, 0xfa, 3, ldAnn},
	{"LD (nn), SP", 0xFF, 0x08, 3, ldNNsp},
	{"CALL nn", 0xFF, 0xCD, 3, call},
	{"JP nn", 0xFF, 0xC3, 3, func(cpu *lr35902) { cpu.regs.PC = cpu.fetched.nn }},
}

var lr35902OpsCodeTableCB = []*opCode{
	{"RLC r", 0b11111000, 0b00000000, 1, cbR},
	{"RLC (HL)", 0xFF, 0x06, 1, cbHL},
	{"RRC r", 0b11111000, 0b00001000, 1, cbR},
	{"RRC (HL)", 0xFF, 0x0e, 1, cbHL},
	{"RL r", 0b11111000, 0b00010000, 1, cbR},
	{"RL (HL)", 0xFF, 0x16, 1, cbHL},
	{"RR r", 0b11111000, 0b00011000, 1, cbR},
	{"RR (HL)", 0xFF, 0x1e, 1, cbHL},
	{"SLA r", 0b11111000, 0b00100000, 1, cbR},
	{"SLA (HL)", 0xFF, 0x26, 1, cbHL},
	{"SRA r", 0b11111000, 0b00101000, 1, cbR},
	{"SRA (HL)", 0xFF, 0x2e, 1, cbHL},
	{"SRL r", 0b11111000, 0b00111000, 1, cbR},
	{"SRL (HL)", 0xFF, 0x3e, 1, cbHL},
	{"BIT b, r", 0b11000000, 0b01000000, 1, bit},
	{"BIT b, (HL)", 0b11000111, 0b01000110, 1, bitHL},
	{"RES b, r", 0b11000000, 0b10000000, 1, res},
	{"RES b, (HL)", 0b11000111, 0b10000110, 1, resHL},
	{"SET b, r", 0b11000000, 0b11000000, 1, set},
	{"SET b, (HL)", 0b11000111, 0b11000110, 1, setHL},
	{"SWAP r", 0b1111_1000, 0b0011_0000, 1, swap},
}

func decodeCB(cpu *lr35902) {
	cpu.scheduler.append(newFetch(lookupCB))
}

func (o *opCode) String() string {
	if o == nil {
		return "<nil>"
	}
	return o.name
}

var lookup = make([]*opCode, 256)
var lookupCB = make([]*opCode, 256)

func init() {
	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range lr35902OpsCodeTable {
			if (code & op.mask) == op.code {
				lookup[code] = op
			}
		}
	}

	// -----

	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range lr35902OpsCodeTableCB {
			op.len++
			if (code & op.mask) == op.code {
				lookupCB[code] = op
			}
		}
	}

	// -----

	println("---------")
	println("                         CB")
	for code := 0; code < 256; code++ {
		fmt.Printf("0x%02X - %-18v%-18v \n", code, lookup[code], lookupCB[code])
	}
	println("---------")
}
