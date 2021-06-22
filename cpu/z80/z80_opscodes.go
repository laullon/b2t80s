package z80

import "fmt"

type opCode struct {
	name       string
	mask, code byte
	len        byte
	ops        []z80op
	onFetch    z80f
}

func (cpu *z80) initOpsCodes() {
	z80OpsCodeTable := []*opCode{
		{"LD dd, mm", 0b11001111, 0b00000001, 3, []z80op{&mrNNpc{f: ldDDmm}}, nil},
		{"ADD HL,ss", 0b11001111, 0b00001001, 1, []z80op{&exec{l: 7, f: addHLss}}, nil},
		{"INC ss", 0b11001111, 0b00000011, 1, []z80op{&exec{l: 2, f: incSS}}, nil},
		{"DEC ss", 0b11001111, 0b00001011, 1, []z80op{&exec{l: 2, f: decSS}}, nil},
		{"POP ss", 0b11001111, 0b11000001, 1, []z80op{}, popSS},
		{"PUSH ss", 0b11001111, 0b11000101, 1, []z80op{&exec{l: 1, f: pushSS}}, nil},

		{"LD r, n", 0b11000111, 0b00000110, 2, []z80op{&mrNpc{f: ldRn}}, nil},
		{"LD r, r'", 0b11000000, 0b01000000, 1, []z80op{}, ldRr},
		{"LD r, (HL)", 0b11000111, 0b01000110, 1, []z80op{}, ldRhl},
		{"LD (HL), r", 0b11111000, 0b01110000, 1, []z80op{}, ldHLr},
		{"INC r", 0b11000111, 0b0000100, 1, []z80op{}, incR},
		{"DEC r", 0b11000111, 0b0000101, 1, []z80op{}, decR},
		{"ADD A, r", 0b11111000, 0b10000000, 1, []z80op{}, addAr},
		{"ADC A, r", 0b11111000, 0b10001000, 1, []z80op{}, adcAr},
		{"SUB A, r", 0b11111000, 0b10010000, 1, []z80op{}, subAr},
		{"SUC A, r", 0b11111000, 0b10011000, 1, []z80op{}, sbcAr},
		{"AND r", 0b11111000, 0b10100000, 1, []z80op{}, andAr},
		{"OR r", 0b11111000, 0b10110000, 1, []z80op{}, orAr},
		{"XOR r", 0b11111000, 0b10101000, 1, []z80op{}, xorAr},
		{"CP r", 0b11111000, 0b10111000, 1, []z80op{}, cpR},

		{"RET cc", 0b11000111, 0b11000000, 1, []z80op{&exec{l: 1, f: retCC}}, nil},
		{"JP cc, nn", 0b11000111, 0b11000010, 3, []z80op{&mrNNpc{f: jpCC}}, nil},
		{"CALL cc, nn", 0b11000111, 0b11000100, 3, []z80op{&mrNNpc{f: callCC}}, nil},
		{"RST p", 0b11000111, 0b11000111, 1, []z80op{&exec{l: 1, f: rstP}}, nil},
		{"CALL nn", 0xFF, 0xCD, 3, []z80op{&mrNNpc{}, &exec{l: 1, f: call}}, nil},

		// {"", 0xFF, 0x1,,[]z80op{},nil},
		{"NOP", 0xFF, 0x00, 1, []z80op{}, nil},
		{"DAA", 0xFF, 0x27, 1, []z80op{}, daa},
		{"CPL", 0xFF, 0x2f, 1, []z80op{}, cpl},
		{"SCF", 0xFF, 0x37, 1, []z80op{}, scf},
		{"CCF", 0xFF, 0x3F, 1, []z80op{}, ccf},
		{"HALT", 0xFF, 0x76, 1, []z80op{}, halt},
		{"RET", 0xFF, 0xC9, 1, []z80op{}, ret},

		{"INC (HL)", 0xFF, 0x34, 1, []z80op{&exec{l: 1, f: incHL}}, nil},
		{"DEC (HL)", 0xFF, 0x35, 1, []z80op{&exec{l: 1, f: decHL}}, nil},
		{"ADD A, (HL)", 0xFF, 0x86, 1, []z80op{}, addAhl},
		{"ADC A, (HL)", 0xFF, 0x8e, 1, []z80op{}, adcAhl},
		{"SUB A, (HL)", 0xFF, 0x96, 1, []z80op{}, subAhl},
		{"SBC A, (HL)", 0xFF, 0x9e, 1, []z80op{}, sbcAhl},
		{"AND (HL)", 0xFF, 0xA6, 1, []z80op{}, andAhl},
		{"OR (HL)", 0xFF, 0xB6, 1, []z80op{}, orAhl},
		{"XOR (HL)", 0xFF, 0xAE, 1, []z80op{}, xorAhl},
		{"CP (HL)", 0xFF, 0xBE, 1, []z80op{}, cpHl},
		{"ADD A, n", 0xFF, 0xc6, 2, []z80op{&mrNpc{f: func(cpu *z80) { cpu.addA(cpu.fetched.n) }}}, nil},
		{"ADC A, (HL)", 0xFF, 0xCE, 2, []z80op{&mrNpc{f: func(cpu *z80) { cpu.adcA(cpu.fetched.n) }}}, nil},
		{"SBC A, (HL)", 0xFF, 0xDE, 2, []z80op{&mrNpc{f: func(cpu *z80) { cpu.sbcA(cpu.fetched.n) }}}, nil},
		{"SUB n", 0xFF, 0xD6, 2, []z80op{&mrNpc{f: func(cpu *z80) { cpu.subA(cpu.fetched.n) }}}, nil},
		{"AND n", 0xFF, 0xE6, 2, []z80op{&mrNpc{f: func(cpu *z80) { cpu.and(cpu.fetched.n) }}}, nil},
		{"OR n", 0xFF, 0xF6, 2, []z80op{&mrNpc{f: func(cpu *z80) { cpu.or(cpu.fetched.n) }}}, nil},

		{"LD A,(BC)", 0xFF, 0x0A, 1, []z80op{}, ldAbc},
		{"LD A,(DE)", 0xFF, 0x1A, 1, []z80op{}, ldAde},
		{"LD (BC), A", 0xFF, 0x02, 1, []z80op{}, ldBCa},
		{"LD (BC), A", 0xFF, 0x12, 1, []z80op{}, ldDEa},
		{"LD (nn), HL", 0xFF, 0x22, 3, []z80op{&mrNNpc{f: ldNNhl}}, nil},
		{"LD (nn), A", 0xFF, 0x32, 3, []z80op{&mrNNpc{f: ldNNa}}, nil},
		{"LD HL, (nn)", 0xFF, 0x2a, 3, []z80op{&mrNNpc{f: ldHLnn}}, nil},
		{"LD (HL), n", 0xFF, 0x36, 2, []z80op{&mrNpc{f: ldHLn}}, nil},
		{"LD A, (nn)", 0xFF, 0x3a, 3, []z80op{&mrNNpc{f: ldAnn}}, nil},

		{"EX AF, AF'", 0xFF, 0x08, 1, []z80op{}, exafaf},
		{"EXX'", 0xFF, 0xD9, 1, []z80op{}, exx},

		{"DJNZ e", 0xFF, 0x10, 2, []z80op{&mrNpc{}, &exec{l: 1, f: djnz}}, nil},
		{"JR e", 0xFF, 0x18, 2, []z80op{&mrNpc{}, &exec{l: 5, f: jr}}, nil},
		{"JRNZ e", 0xFF, 0x20, 2, []z80op{&mrNpc{f: jrnz}}, nil},
		{"JRZ e", 0xFF, 0x28, 2, []z80op{&mrNpc{f: jrz}}, nil},
		{"JRNC e", 0xFF, 0x30, 2, []z80op{&mrNpc{f: jrnc}}, nil},
		{"JRC e", 0xFF, 0x38, 2, []z80op{&mrNpc{f: jrc}}, nil},

		{"JP nn", 0xFF, 0xC3, 3, []z80op{&mrNNpc{f: func(cpu *z80) { cpu.regs.PC = cpu.fetched.nn }}}, nil},

		{"RLCA", 0xFF, 0x07, 1, []z80op{}, rlca},
		{"RLA", 0xFF, 0x17, 1, []z80op{}, rla},
		{"RRCA", 0xFF, 0x0F, 1, []z80op{}, rrca},
		{"RRA", 0xFF, 0x1F, 1, []z80op{}, rra},

		{"OUT (n), A", 0xFF, 0xD3, 2, []z80op{&mrNpc{}, &exec{l: 1, f: outNa}}, nil},
		{"IN A, (n)", 0xFF, 0xDB, 2, []z80op{&mrNpc{f: inAn}}, nil},

		{"EX (SP), IX", 0xFF, 0xE3, 1, []z80op{}, exSP},
		{"JP HL", 0xFF, 0xE9, 1, []z80op{}, func(cpu *z80) { cpu.regs.PC = cpu.regs.HL.Get() }},
		{"EX DE, HL", 0xFF, 0xEB, 1, []z80op{}, exDEhl},

		{"XOR *", 0xFF, 0xEE, 2, []z80op{&mrNpc{f: func(cpu *z80) { cpu.xor(cpu.fetched.n) }}}, nil},
		{"DI", 0xFF, 0xF3, 1, []z80op{}, func(cpu *z80) { cpu.regs.IFF1 = false; cpu.regs.IFF2 = false }},
		{"EI", 0xFF, 0xFb, 1, []z80op{}, func(cpu *z80) { cpu.regs.IFF1 = true; cpu.regs.IFF2 = true }},
		{"LD SP, HL", 0xFF, 0xF9, 1, []z80op{&exec{l: 2, f: func(cpu *z80) { cpu.regs.SP.Set(cpu.regs.HL.Get()) }}}, nil},
		{"CP *", 0xFF, 0xFe, 2, []z80op{&mrNpc{f: func(cpu *z80) { cpu.cp(cpu.fetched.n) }}}, nil},

		{"CB", 0xFF, 0xCB, 1, []z80op{}, decodeCB},
		{"DD", 0xFF, 0xDD, 1, []z80op{}, decodeDD},
		{"ED", 0xFF, 0xED, 1, []z80op{}, decodeED},
		{"ED", 0xFF, 0xFD, 1, []z80op{}, decodeFD},
	}

	z80OpsCodeTableCB := []*opCode{
		{"RLC r", 0b11111000, 0b00000000, 1, []z80op{}, cbR},
		{"RLC (HL)", 0xFF, 0x06, 1, []z80op{&exec{l: 1, f: cbHL}}, nil},
		{"RRC r", 0b11111000, 0b00001000, 1, []z80op{}, cbR},
		{"RRC (HL)", 0xFF, 0x0e, 1, []z80op{&exec{l: 1, f: cbHL}}, nil},

		{"RL r", 0b11111000, 0b00010000, 1, []z80op{}, cbR},
		{"RL (HL)", 0xFF, 0x16, 1, []z80op{&exec{l: 1, f: cbHL}}, nil},
		{"RR r", 0b11111000, 0b00011000, 1, []z80op{}, cbR},
		{"RR (HL)", 0xFF, 0x1e, 1, []z80op{&exec{l: 1, f: cbHL}}, nil},

		{"SLA r", 0b11111000, 0b00100000, 1, []z80op{}, cbR},
		{"SLA (HL)", 0xFF, 0x26, 1, []z80op{&exec{l: 1, f: cbHL}}, nil},
		{"SRA r", 0b11111000, 0b00101000, 1, []z80op{}, cbR},
		{"SRA (HL)", 0xFF, 0x2e, 1, []z80op{&exec{l: 1, f: cbHL}}, nil},

		{"SLL r", 0b11111000, 0b00110000, 1, []z80op{}, cbR},
		{"SLL (HL)", 0xFF, 0x36, 1, []z80op{&exec{l: 1, f: cbHL}}, nil},
		{"SRL r", 0b11111000, 0b00111000, 1, []z80op{}, cbR},
		{"SRL (HL)", 0xFF, 0x3e, 1, []z80op{&exec{l: 1, f: cbHL}}, nil},

		{"BIT b, r", 0b11000000, 0b01000000, 1, []z80op{}, bit},
		{"BIT b, (HL)", 0b11000111, 0b01000110, 1, []z80op{&exec{l: 1, f: bitHL}}, nil},

		{"RES b, r", 0b11000000, 0b10000000, 1, []z80op{}, res},
		{"RES b, (HL)", 0b11000111, 0b10000110, 1, []z80op{&exec{l: 1, f: resHL}}, nil},

		{"SET b, r", 0b11000000, 0b11000000, 1, []z80op{}, set},
		{"SET b, (HL)", 0b11000111, 0b11000110, 1, []z80op{&exec{l: 1, f: setHL}}, nil},
	}

	z80OpsCodeTableDD := []*opCode{
		{"LD r, r'", 0b11000000, 0b01000000, 1, []z80op{}, ldRr},
		{"ADD IX, rr", 0b11001111, 0b00001001, 1, []z80op{&exec{l: 7, f: addIXY}}, nil},
		{"LD IX, nn", 0xFF, 0x21, 3, []z80op{&mrNNpc{f: func(cpu *z80) { cpu.regs.IXH = cpu.fetched.n2; cpu.regs.IXL = cpu.fetched.n }}}, nil},
		{"LD (nn), IX", 0xFF, 0x22, 3, []z80op{&mrNNpc{f: ldNNIXY}}, nil},
		{"INC IX", 0xFF, 0x23, 1, []z80op{&exec{l: 2, f: func(cpu *z80) { cpu.regs.IX.Set(cpu.regs.IX.Get() + 1) }}}, nil},
		{"INC IXH", 0xFF, 0x24, 1, []z80op{}, func(cpu *z80) { cpu.incR(&cpu.regs.IXH) }},
		{"DEC IXH", 0xFF, 0x25, 1, []z80op{}, func(cpu *z80) { cpu.decR(&cpu.regs.IXH) }},
		{"LD IXH, n", 0xFF, 0x26, 2, []z80op{&mrNpc{f: func(cpu *z80) { cpu.regs.IXH = cpu.fetched.n }}}, nil},
		{"LD IX, nn", 0xFF, 0x2A, 3, []z80op{&mrNNpc{f: ldIXYnn}}, nil},
		{"DEC IX", 0xFF, 0x2B, 1, []z80op{&exec{l: 2, f: func(cpu *z80) { cpu.regs.IX.Set(cpu.regs.IX.Get() - 1) }}}, nil},
		{"INC IXL", 0xFF, 0x2C, 1, []z80op{}, func(cpu *z80) { cpu.incR(&cpu.regs.IXL) }},
		{"DEC IXL", 0xFF, 0x2D, 1, []z80op{}, func(cpu *z80) { cpu.decR(&cpu.regs.IXL) }},
		{"LD IXL, n", 0xFF, 0x2E, 2, []z80op{&mrNpc{f: func(cpu *z80) { cpu.regs.IXL = cpu.fetched.n }}}, nil},
		{"INC (IX+d)", 0xFF, 0x34, 2, []z80op{&mrNpc{}, &exec{l: 6, f: incIXYd}}, nil},
		{"DEC (IX+d)", 0xFF, 0x35, 2, []z80op{&mrNpc{}, &exec{l: 6, f: decIXYd}}, nil},
		{"LD (IX+d), n", 0xFF, 0x36, 3, []z80op{&mrNNpc{}, &exec{l: 2, f: ldIXYdN}}, nil},

		{"LD B, IXH", 0xFF, 0x44, 1, []z80op{}, func(cpu *z80) { cpu.regs.B = cpu.regs.IXH }},
		{"LD B, IXL", 0xFF, 0x45, 1, []z80op{}, func(cpu *z80) { cpu.regs.B = cpu.regs.IXL }},
		{"LD C, IXH", 0xFF, 0x4C, 1, []z80op{}, func(cpu *z80) { cpu.regs.C = cpu.regs.IXH }},
		{"LD C, IXL", 0xFF, 0x4D, 1, []z80op{}, func(cpu *z80) { cpu.regs.C = cpu.regs.IXL }},
		{"LD D, IXH", 0xFF, 0x54, 1, []z80op{}, func(cpu *z80) { cpu.regs.D = cpu.regs.IXH }},
		{"LD D, IXL", 0xFF, 0x55, 1, []z80op{}, func(cpu *z80) { cpu.regs.D = cpu.regs.IXL }},
		{"LD E, IXH", 0xFF, 0x5C, 1, []z80op{}, func(cpu *z80) { cpu.regs.E = cpu.regs.IXH }},
		{"LD E, IXL", 0xFF, 0x5D, 1, []z80op{}, func(cpu *z80) { cpu.regs.E = cpu.regs.IXL }},
		{"LD A, IXH", 0xFF, 0x7C, 1, []z80op{}, func(cpu *z80) { cpu.regs.A = cpu.regs.IXH }},
		{"LD A, IXL", 0xFF, 0x7D, 1, []z80op{}, func(cpu *z80) { cpu.regs.A = cpu.regs.IXL }},

		{"LD IXH, r", 0b11111000, 0b01100000, 1, []z80op{}, ldIXYHr},
		{"LD IXH, r", 0b11111000, 0b01101000, 1, []z80op{}, ldIXYLr},
		{"LD r, (IX+d)", 0b11000111, 0b01000110, 2, []z80op{&mrNpc{}, &exec{l: 5, f: ldRixyD}}, nil},
		{"LD (IX+d), r", 0b11111000, 0b01110000, 2, []z80op{&mrNpc{}, &exec{l: 5, f: ldIXYdR}}, nil},

		{"ADD A, IXH", 0xFF, 0x84, 1, []z80op{}, func(cpu *z80) { cpu.addA(cpu.regs.IXH) }},
		{"ADD A, IXL", 0xFF, 0x85, 1, []z80op{}, func(cpu *z80) { cpu.addA(cpu.regs.IXL) }},
		{"ADC A, IXH", 0xFF, 0x8C, 1, []z80op{}, func(cpu *z80) { cpu.adcA(cpu.regs.IXH) }},
		{"ADC A, IXL", 0xFF, 0x8D, 1, []z80op{}, func(cpu *z80) { cpu.adcA(cpu.regs.IXL) }},
		{"SUB A, IXH", 0xFF, 0x94, 1, []z80op{}, func(cpu *z80) { cpu.subA(cpu.regs.IXH) }},
		{"SUB A, IXL", 0xFF, 0x95, 1, []z80op{}, func(cpu *z80) { cpu.subA(cpu.regs.IXL) }},
		{"SBC A, IXH", 0xFF, 0x9C, 1, []z80op{}, func(cpu *z80) { cpu.sbcA(cpu.regs.IXH) }},
		{"SBC A, IXL", 0xFF, 0x9D, 1, []z80op{}, func(cpu *z80) { cpu.sbcA(cpu.regs.IXL) }},
		{"AND A, IXH", 0xFF, 0xA4, 1, []z80op{}, func(cpu *z80) { cpu.and(cpu.regs.IXH) }},
		{"AND A, IXL", 0xFF, 0xA5, 1, []z80op{}, func(cpu *z80) { cpu.and(cpu.regs.IXL) }},
		{"XOR A, IXH", 0xFF, 0xAC, 1, []z80op{}, func(cpu *z80) { cpu.xor(cpu.regs.IXH) }},
		{"XOR A, IXL", 0xFF, 0xAD, 1, []z80op{}, func(cpu *z80) { cpu.xor(cpu.regs.IXL) }},
		{"OR A, IXH", 0xFF, 0xB4, 1, []z80op{}, func(cpu *z80) { cpu.or(cpu.regs.IXH) }},
		{"OR A, IXL", 0xFF, 0xB5, 1, []z80op{}, func(cpu *z80) { cpu.or(cpu.regs.IXL) }},
		{"CP A, IXH", 0xFF, 0xBC, 1, []z80op{}, func(cpu *z80) { cpu.cp(cpu.regs.IXH) }},
		{"CP A, IXL", 0xFF, 0xBD, 1, []z80op{}, func(cpu *z80) { cpu.cp(cpu.regs.IXL) }},

		{"ADD A, (IX+d)", 0xFF, 0x86, 2, []z80op{&mrNpc{}, &exec{l: 5, f: addAixyD}}, nil},
		{"ADC A, (IX+d)", 0xFF, 0x8E, 2, []z80op{&mrNpc{}, &exec{l: 5, f: adcAixyD}}, nil},
		{"SUB A, (IX+d)", 0xFF, 0x96, 2, []z80op{&mrNpc{}, &exec{l: 5, f: subAixyD}}, nil},
		{"SBC A, (IX+d)", 0xFF, 0x9E, 2, []z80op{&mrNpc{}, &exec{l: 5, f: sbcAixyD}}, nil},
		{"AND A, (IX+d)", 0xFF, 0xA6, 2, []z80op{&mrNpc{}, &exec{l: 5, f: andAixyD}}, nil},
		{"XOR A, (IX+d)", 0xFF, 0xAE, 2, []z80op{&mrNpc{}, &exec{l: 5, f: xorAixyD}}, nil},
		{"OR A, (IX+d)", 0xFF, 0xB6, 2, []z80op{&mrNpc{}, &exec{l: 5, f: orAixyD}}, nil},
		{"CP A, (IX+d)", 0xFF, 0xBE, 2, []z80op{&mrNpc{}, &exec{l: 5, f: cpAixyD}}, nil},

		{"CB", 0xFF, 0xCB, 2, []z80op{&mrNpc{f: decodeDDCB}}, nil},

		{"POP IX", 0xFF, 0xE1, 1, []z80op{}, func(cpu *z80) { cpu.popFromStack(func(cpu *z80, data uint16) { cpu.regs.IX.Set(data) }) }},
		{"EX (SP), IX", 0xFF, 0xE3, 1, []z80op{}, exSP},
		{"PUSH IX", 0xFF, 0xE5, 1, []z80op{&exec{l: 1, f: func(cpu *z80) { cpu.pushToStack(cpu.regs.IX.Get(), nil) }}}, nil},
		{"JP IX", 0xFF, 0xE9, 1, []z80op{}, func(cpu *z80) { cpu.regs.PC = cpu.regs.IX.Get() }},
		{"LD SP, IX", 0xFF, 0xF9, 1, []z80op{&exec{l: 2, f: func(cpu *z80) { cpu.regs.SP.Set(cpu.regs.IX.Get()) }}}, nil},
	}

	z80OpsCodeTableFD := []*opCode{
		{"LD r, r'", 0b11000000, 0b01000000, 1, []z80op{}, ldRr},
		{"ADD IY, rr", 0b11001111, 0b00001001, 1, []z80op{&exec{l: 7, f: addIY}}, nil},
		{"LD IY, nn", 0xFF, 0x21, 3, []z80op{&mrNNpc{f: func(cpu *z80) { cpu.regs.IYH = cpu.fetched.n2; cpu.regs.IYL = cpu.fetched.n }}}, nil},
		{"LD (nn), IY", 0xFF, 0x22, 3, []z80op{&mrNNpc{f: ldNNIXY}}, nil},
		{"INC IY", 0xFF, 0x23, 1, []z80op{&exec{l: 2, f: func(cpu *z80) { cpu.regs.IY.Set(cpu.regs.IY.Get() + 1) }}}, nil},
		{"INC IYH", 0xFF, 0x24, 1, []z80op{}, func(cpu *z80) { cpu.incR(&cpu.regs.IYH) }},
		{"DEC IYH", 0xFF, 0x25, 1, []z80op{}, func(cpu *z80) { cpu.decR(&cpu.regs.IYH) }},
		{"LD IYH, n", 0xFF, 0x26, 2, []z80op{&mrNpc{f: func(cpu *z80) { cpu.regs.IYH = cpu.fetched.n }}}, nil},
		{"LD IY, nn", 0xFF, 0x2A, 3, []z80op{&mrNNpc{f: ldIXYnn}}, nil},
		{"DEC IY", 0xFF, 0x2B, 1, []z80op{&exec{l: 2, f: func(cpu *z80) { cpu.regs.IY.Set(cpu.regs.IY.Get() - 1) }}}, nil},
		{"INC IYL", 0xFF, 0x2C, 1, []z80op{}, func(cpu *z80) { cpu.incR(&cpu.regs.IYL) }},
		{"DEC IYL", 0xFF, 0x2D, 1, []z80op{}, func(cpu *z80) { cpu.decR(&cpu.regs.IYL) }},
		{"LD IYL, n", 0xFF, 0x2E, 2, []z80op{&mrNpc{f: func(cpu *z80) { cpu.regs.IYL = cpu.fetched.n }}}, nil},
		{"INC (IY+d)", 0xFF, 0x34, 2, []z80op{&mrNpc{}, &exec{l: 6, f: incIXYd}}, nil},
		{"DEC (IY+d)", 0xFF, 0x35, 2, []z80op{&mrNpc{}, &exec{l: 6, f: decIXYd}}, nil},
		{"LD (IY+d), n", 0xFF, 0x36, 3, []z80op{&mrNNpc{}, &exec{l: 2, f: ldIXYdN}}, nil},

		{"LD B, IYH", 0xFF, 0x44, 1, []z80op{}, func(cpu *z80) { cpu.regs.B = cpu.regs.IYH }},
		{"LD B, IYL", 0xFF, 0x45, 1, []z80op{}, func(cpu *z80) { cpu.regs.B = cpu.regs.IYL }},
		{"LD C, IYH", 0xFF, 0x4C, 1, []z80op{}, func(cpu *z80) { cpu.regs.C = cpu.regs.IYH }},
		{"LD C, IYL", 0xFF, 0x4D, 1, []z80op{}, func(cpu *z80) { cpu.regs.C = cpu.regs.IYL }},
		{"LD D, IYH", 0xFF, 0x54, 1, []z80op{}, func(cpu *z80) { cpu.regs.D = cpu.regs.IYH }},
		{"LD D, IYL", 0xFF, 0x55, 1, []z80op{}, func(cpu *z80) { cpu.regs.D = cpu.regs.IYL }},
		{"LD E, IYH", 0xFF, 0x5C, 1, []z80op{}, func(cpu *z80) { cpu.regs.E = cpu.regs.IYH }},
		{"LD E, IYL", 0xFF, 0x5D, 1, []z80op{}, func(cpu *z80) { cpu.regs.E = cpu.regs.IYL }},
		{"LD A, IYH", 0xFF, 0x7C, 1, []z80op{}, func(cpu *z80) { cpu.regs.A = cpu.regs.IYH }},
		{"LD A, IYL", 0xFF, 0x7D, 1, []z80op{}, func(cpu *z80) { cpu.regs.A = cpu.regs.IYL }},

		{"LD IYH, r", 0b11111000, 0b01100000, 1, []z80op{}, ldIXYHr},
		{"LD IYL, r", 0b11111000, 0b01101000, 1, []z80op{}, ldIXYLr},
		{"LD r, (IY+d)", 0b11000111, 0b01000110, 2, []z80op{&mrNpc{}, &exec{l: 5, f: ldRixyD}}, nil},
		{"LD (IY+d), r", 0b11111000, 0b01110000, 2, []z80op{&mrNpc{}, &exec{l: 5, f: ldIXYdR}}, nil},

		{"ADD A, IYH", 0xFF, 0x84, 1, []z80op{}, func(cpu *z80) { cpu.addA(cpu.regs.IYH) }},
		{"ADD A, IYL", 0xFF, 0x85, 1, []z80op{}, func(cpu *z80) { cpu.addA(cpu.regs.IYL) }},
		{"ADC A, IYH", 0xFF, 0x8C, 1, []z80op{}, func(cpu *z80) { cpu.adcA(cpu.regs.IYH) }},
		{"ADC A, IYL", 0xFF, 0x8D, 1, []z80op{}, func(cpu *z80) { cpu.adcA(cpu.regs.IYL) }},
		{"SUB A, IYH", 0xFF, 0x94, 1, []z80op{}, func(cpu *z80) { cpu.subA(cpu.regs.IYH) }},
		{"SUB A, IYL", 0xFF, 0x95, 1, []z80op{}, func(cpu *z80) { cpu.subA(cpu.regs.IYL) }},
		{"SBC A, IYH", 0xFF, 0x9C, 1, []z80op{}, func(cpu *z80) { cpu.sbcA(cpu.regs.IYH) }},
		{"SBC A, IYL", 0xFF, 0x9D, 1, []z80op{}, func(cpu *z80) { cpu.sbcA(cpu.regs.IYL) }},
		{"AND A, IYH", 0xFF, 0xA4, 1, []z80op{}, func(cpu *z80) { cpu.and(cpu.regs.IYH) }},
		{"AND A, IYL", 0xFF, 0xA5, 1, []z80op{}, func(cpu *z80) { cpu.and(cpu.regs.IYL) }},
		{"XOR A, IYH", 0xFF, 0xAC, 1, []z80op{}, func(cpu *z80) { cpu.xor(cpu.regs.IYH) }},
		{"XOR A, IYL", 0xFF, 0xAD, 1, []z80op{}, func(cpu *z80) { cpu.xor(cpu.regs.IYL) }},
		{"OR A, IYH", 0xFF, 0xB4, 1, []z80op{}, func(cpu *z80) { cpu.or(cpu.regs.IYH) }},
		{"OR A, IYL", 0xFF, 0xB5, 1, []z80op{}, func(cpu *z80) { cpu.or(cpu.regs.IYL) }},
		{"CP A, IYH", 0xFF, 0xBC, 1, []z80op{}, func(cpu *z80) { cpu.cp(cpu.regs.IYH) }},
		{"CP A, IYL", 0xFF, 0xBD, 1, []z80op{}, func(cpu *z80) { cpu.cp(cpu.regs.IYL) }},

		{"ADD A, (IY+d)", 0xFF, 0x86, 2, []z80op{&mrNpc{}, &exec{l: 5, f: addAixyD}}, nil},
		{"ADC A, (IY+d)", 0xFF, 0x8E, 2, []z80op{&mrNpc{}, &exec{l: 5, f: adcAixyD}}, nil},
		{"SUB A, (IY+d)", 0xFF, 0x96, 2, []z80op{&mrNpc{}, &exec{l: 5, f: subAixyD}}, nil},
		{"SBC A, (IY+d)", 0xFF, 0x9E, 2, []z80op{&mrNpc{}, &exec{l: 5, f: sbcAixyD}}, nil},
		{"AND A, (IY+d)", 0xFF, 0xA6, 2, []z80op{&mrNpc{}, &exec{l: 5, f: andAixyD}}, nil},
		{"XOR A, (IY+d)", 0xFF, 0xAE, 2, []z80op{&mrNpc{}, &exec{l: 5, f: xorAixyD}}, nil},
		{"OR A, (IY+d)", 0xFF, 0xB6, 2, []z80op{&mrNpc{}, &exec{l: 5, f: orAixyD}}, nil},
		{"CP A, (IY+d)", 0xFF, 0xBE, 2, []z80op{&mrNpc{}, &exec{l: 5, f: cpAixyD}}, nil},

		{"CB", 0xFF, 0xCB, 2, []z80op{&mrNpc{f: decodeFDCB}}, nil},

		{"POP IY", 0xFF, 0xE1, 1, []z80op{}, func(cpu *z80) { cpu.popFromStack(func(cpu *z80, data uint16) { cpu.regs.IY.Set(data) }) }},
		{"EX (SP), IY", 0xFF, 0xE3, 1, []z80op{}, exSP},
		{"PUSH IY", 0xFF, 0xE5, 1, []z80op{&exec{l: 1, f: func(cpu *z80) { cpu.pushToStack(cpu.regs.IY.Get(), nil) }}}, nil},
		{"JP IY", 0xFF, 0xE9, 1, []z80op{}, func(cpu *z80) { cpu.regs.PC = cpu.regs.IY.Get() }},
		{"LD SP, IY", 0xFF, 0xF9, 1, []z80op{&exec{l: 2, f: func(cpu *z80) { cpu.regs.SP.Set(cpu.regs.IY.Get()) }}}, nil},
	}

	z80OpsCodeTableDDCB := []*opCode{
		{"RLC (IX+d), r", 0b11111000, 0b00000000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"RLC (IX+d)", 0xFF, 0x06, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},
		{"RRC (IX+d), r", 0b11111000, 0b00001000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"RRC (IX+d)", 0xFF, 0x0e, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},

		{"RL (IX+d), r", 0b11111000, 0b00010000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"RL (IX+d)", 0xFF, 0x16, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},
		{"RR (IX+d), r", 0b11111000, 0b00011000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"RR (IX+d)", 0xFF, 0x1e, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},

		{"SLA (IX+d), r", 0b11111000, 0b00100000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"SLA (IX+d)", 0xFF, 0x26, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},
		{"SRA (IX+d), r", 0b11111000, 0b00101000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"SRA (IX+d)", 0xFF, 0x2e, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},

		{"SLL (IX+d), r", 0b11111000, 0b00110000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"SLL (IX+d)", 0xFF, 0x36, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},
		{"SRL (IX+d), r", 0b11111000, 0b00111000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"SRL (IX+d)", 0xFF, 0x3e, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},

		{"BIT b, (IX+d), r", 0b11000000, 0b01000000, 1, []z80op{&exec{l: 2, f: bitIXYd}}, nil},
		{"BIT b, (IX+d)", 0b11000111, 0b01000110, 1, []z80op{&exec{l: 2, f: bitIXYd}}, nil},

		{"RES b, (IX+d), r", 0b11000000, 0b10000000, 1, []z80op{&exec{l: 2, f: resIXYdR}}, nil},
		{"RES b, (IX+d)", 0b11000111, 0b10000110, 1, []z80op{&exec{l: 2, f: resIXYd}}, nil},

		{"SET b, (IX+d), r", 0b11000000, 0b11000000, 1, []z80op{&exec{l: 2, f: setIXYdR}}, nil},
		{"SET b, (IX+d)", 0b11000111, 0b11000110, 1, []z80op{&exec{l: 2, f: setIXYd}}, nil},
	}

	z80OpsCodeTableFDCB := []*opCode{
		{"RLC (IY+d), r", 0b11111000, 0b00000000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"RLC (IY+d)", 0xFF, 0x06, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},
		{"RRC (IY+d), r", 0b11111000, 0b00001000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"RRC (IY+d)", 0xFF, 0x0e, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},

		{"RL (IY+d), r", 0b11111000, 0b00010000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"RL (IY+d)", 0xFF, 0x16, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},
		{"RR (IY+d), r", 0b11111000, 0b00011000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"RR (IY+d)", 0xFF, 0x1e, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},

		{"SLA (IY+d), r", 0b11111000, 0b00100000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"SLA (IY+d)", 0xFF, 0x26, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},
		{"SRA (IY+d), r", 0b11111000, 0b00101000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"SRA (IY+d)", 0xFF, 0x2e, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},

		{"SLL (IY+d), r", 0b11111000, 0b00110000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"SLL (IY+d)", 0xFF, 0x36, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},
		{"SRL (IY+d), r", 0b11111000, 0b00111000, 1, []z80op{&exec{l: 2, f: cbIXYdr}}, nil},
		{"SRL (IY+d)", 0xFF, 0x3e, 1, []z80op{&exec{l: 2, f: cbIXYd}}, nil},

		{"BIT b, (IY+d), r", 0b11000000, 0b01000000, 1, []z80op{&exec{l: 2, f: bitIXYd}}, nil},
		{"BIT b, (IY+d)", 0b11000111, 0b01000110, 1, []z80op{&exec{l: 2, f: bitIXYd}}, nil},

		{"RES b, (IY+d), r", 0b11000000, 0b10000000, 1, []z80op{&exec{l: 2, f: resIXYdR}}, nil},
		{"RES b, (IY+d)", 0b11000111, 0b10000110, 1, []z80op{&exec{l: 2, f: resIXYd}}, nil},

		{"SET b, (IY+d), r", 0b11000000, 0b11000000, 1, []z80op{&exec{l: 2, f: setIXYdR}}, nil},
		{"SET b, (IY+d)", 0b11000111, 0b11000110, 1, []z80op{&exec{l: 2, f: setIXYd}}, nil},
	}

	z80OpsCodeTableED := []*opCode{
		{"IN r, (c)", 0b11000111, 0b01000000, 1, []z80op{}, inRc},
		{"IN (c)", 0xFF, 0x70, 1, []z80op{}, inC},
		{"OUT (c), r", 0b11000111, 0b01000001, 1, []z80op{&exec{l: 1, f: outCr}}, nil},
		{"OUT (c), 0", 0xFF, 0x71, 1, []z80op{&exec{l: 1, f: outC0}}, nil},

		{"SBC HL, BC", 0xFF, 0x42, 1, []z80op{&exec{l: 7, f: func(cpu *z80) { cpu.sbcHL(cpu.regs.BC.Get()) }}}, nil},
		{"SBC HL, DE", 0xFF, 0x52, 1, []z80op{&exec{l: 7, f: func(cpu *z80) { cpu.sbcHL(cpu.regs.DE.Get()) }}}, nil},
		{"SBC HL, HL", 0xFF, 0x62, 1, []z80op{&exec{l: 7, f: func(cpu *z80) { cpu.sbcHL(cpu.regs.HL.Get()) }}}, nil},
		{"SBC HL, SP", 0xFF, 0x72, 1, []z80op{&exec{l: 7, f: func(cpu *z80) { cpu.sbcHL(cpu.regs.SP.Get()) }}}, nil},

		{"LD (nn), dd", 0b11001111, 0b01000011, 3, []z80op{&mrNNpc{f: ldNNdd}}, nil},
		{"NEG", 0b11000111, 0b01000100, 1, []z80op{}, func(cpu *z80) { n := cpu.regs.A; cpu.regs.A = 0; cpu.subA(n) }},
		{"RETN", 0b11000111, 0b01000101, 1, []z80op{}, func(cpu *z80) { cpu.regs.IFF1 = cpu.regs.IFF2; ret(cpu) }},

		{"IM 0", 0xFF, 0x46, 1, []z80op{}, func(cpu *z80) { cpu.regs.InterruptsMode = 0 }},
		{"IM 0", 0xFF, 0x66, 1, []z80op{}, func(cpu *z80) { cpu.regs.InterruptsMode = 0 }},
		{"IM 1", 0xFF, 0x56, 1, []z80op{}, func(cpu *z80) { cpu.regs.InterruptsMode = 1 }},
		{"IM 2", 0xFF, 0xE5, 1, []z80op{}, func(cpu *z80) { cpu.regs.InterruptsMode = 2 }},
		{"IM 0/1", 0xFF, 0x4E, 1, []z80op{}, func(cpu *z80) { cpu.regs.InterruptsMode = 0 }},
		{"IM 2", 0xFF, 0x5E, 1, []z80op{}, func(cpu *z80) { cpu.regs.InterruptsMode = 2 }},
		{"IM 0/1", 0xFF, 0x6E, 1, []z80op{}, func(cpu *z80) { cpu.regs.InterruptsMode = 0 }},
		{"IM 1", 0xFF, 0x76, 1, []z80op{}, func(cpu *z80) { cpu.regs.InterruptsMode = 1 }},
		{"IM 2", 0xFF, 0x7E, 1, []z80op{}, func(cpu *z80) { cpu.regs.InterruptsMode = 2 }},

		{"LD I, A", 0xFF, 0x47, 1, []z80op{&exec{l: 1, f: func(cpu *z80) { cpu.regs.I = cpu.regs.A }}}, nil},
		{"LD R, A", 0xFF, 0x4F, 1, []z80op{&exec{l: 1, f: func(cpu *z80) { cpu.regs.R = cpu.regs.A }}}, nil},

		{"LD A, I", 0xFF, 0x57, 1, []z80op{&exec{l: 1, f: ldAi}}, nil},
		{"LD A, R", 0xFF, 0x5F, 1, []z80op{&exec{l: 1, f: ldAr}}, nil},

		{"ADC HL, BC", 0xFF, 0x4a, 1, []z80op{&exec{l: 7, f: func(cpu *z80) { cpu.adcHL(cpu.regs.BC.Get()) }}}, nil},
		{"ADC HL, DE", 0xFF, 0x5a, 1, []z80op{&exec{l: 7, f: func(cpu *z80) { cpu.adcHL(cpu.regs.DE.Get()) }}}, nil},
		{"ADC HL, HL", 0xFF, 0x6a, 1, []z80op{&exec{l: 7, f: func(cpu *z80) { cpu.adcHL(cpu.regs.HL.Get()) }}}, nil},
		{"ADC HL, SP", 0xFF, 0x7a, 1, []z80op{&exec{l: 7, f: func(cpu *z80) { cpu.adcHL(cpu.regs.SP.Get()) }}}, nil},

		{"LD (nn), dd", 0b11001111, 0b01001011, 3, []z80op{&mrNNpc{f: ldDDnn}}, nil},

		{"RDD", 0xFF, 0x67, 1, []z80op{}, rrd},
		{"RDD", 0xFF, 0x6f, 1, []z80op{}, rld},

		{"LDI", 0xFF, 0xA0, 1, []z80op{}, ldi},
		{"CPI", 0xFF, 0xA1, 1, []z80op{}, cpi},
		{"INI", 0xFF, 0xA2, 1, []z80op{}, ini},
		{"OUTI", 0xFF, 0xA3, 1, []z80op{}, outi},

		{"LDD", 0xFF, 0xA8, 1, []z80op{}, ldd},
		{"CPD", 0xFF, 0xA9, 1, []z80op{}, cpd},
		{"IND", 0xFF, 0xAA, 1, []z80op{}, ind},
		{"OUTD", 0xFF, 0xAB, 1, []z80op{}, outd},

		{"LDIR", 0xFF, 0xB0, 1, []z80op{}, ldi},
		{"CPIR", 0xFF, 0xB1, 1, []z80op{}, cpi},
		{"INIR", 0xFF, 0xB2, 1, []z80op{}, ini},
		{"OTIR", 0xFF, 0xB3, 1, []z80op{}, outi},

		{"LDDR", 0xFF, 0xB8, 1, []z80op{}, ldd},
		{"CPDR", 0xFF, 0xB9, 1, []z80op{}, cpd},
		{"INDR", 0xFF, 0xBA, 1, []z80op{}, ind},
		{"OTDR", 0xFF, 0xBB, 1, []z80op{}, outd},
	}

	cpu.lookup = make([]*opCode, 256)
	cpu.lookupCB = make([]*opCode, 256)
	cpu.lookupDD = make([]*opCode, 256)
	cpu.lookupED = make([]*opCode, 256)
	cpu.lookupFD = make([]*opCode, 256)
	cpu.lookupDDCB = make([]*opCode, 256)
	cpu.lookupFDCB = make([]*opCode, 256)

	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTable {
			if (code & op.mask) == op.code {
				cpu.lookup[code] = op
			}
		}
	}

	// -----

	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTableCB {
			op.len++
			if (code & op.mask) == op.code {
				cpu.lookupCB[code] = op
			}
		}
	}

	// -----
	for _, op := range z80OpsCodeTableDD {
		op.len++
	}
	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTableDD {
			if (code & op.mask) == op.code {
				cpu.lookupDD[code] = op
			}
		}
	}

	// -----
	for _, op := range z80OpsCodeTableED {
		op.len++
	}
	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTableED {
			if (code & op.mask) == op.code {
				cpu.lookupED[code] = op
			}
		}
	}

	// -----
	for _, op := range z80OpsCodeTableFD {
		op.len++
	}
	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTableFD {
			if (code & op.mask) == op.code {
				cpu.lookupFD[code] = op
			}
		}
	}

	// -----
	for _, op := range z80OpsCodeTableDDCB {
		op.len += 2
	}
	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTableDDCB {
			if (code & op.mask) == op.code {
				cpu.lookupDDCB[code] = op
			}
		}
	}

	// -----
	for _, op := range z80OpsCodeTableFDCB {
		op.len += 2
	}
	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTableFDCB {
			if (code & op.mask) == op.code {
				cpu.lookupFDCB[code] = op
			}
		}
	}

	// -----

	// println("---------")
	// println("                         CB                DD                DDCB              ED                FD                FDCB")
	// for code := 0; code < 256; code++ {
	// 	fmt.Printf("0x%02X - %-18v%-18v%-18v%-18v%-18v%-18v%-18v\n", code, cpu.lookup[code], cpu.lookupCB[code], cpu.lookupDD[code], cpu.lookupDDCB[code], cpu.lookupED[code], cpu.lookupFD[code], cpu.lookupFDCB[code])
	// }
	// println("---------")}
}

func decodeCB(cpu *z80) {
	cpu.scheduler.append(cpu.newFetch(cpu.lookupCB))
}

func decodeDD(cpu *z80) {
	cpu.indexIdx = 1
	cpu.scheduler.append(cpu.newFetch(cpu.lookupDD))
}

func decodeED(cpu *z80) {
	cpu.scheduler.append(cpu.newFetch(cpu.lookupED))
}

func decodeFD(cpu *z80) {
	cpu.indexIdx = 2
	cpu.scheduler.append(cpu.newFetch(cpu.lookupFD))
}

func decodeDDCB(cpu *z80) {
	cpu.regs.R--
	cpu.scheduler.append(cpu.newFetch(cpu.lookupDDCB))
}

func decodeFDCB(cpu *z80) {
	cpu.regs.R--
	cpu.scheduler.append(cpu.newFetch(cpu.lookupFDCB))
}

func (o *opCode) String() string {
	if o == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s", o.name)
}
