package z80

import "fmt"

type opCode struct {
	name       string
	mask, code byte
	ops        []z80op
	onFetch    z80f
}

var z80OpsCodeTable = []*opCode{
	{"LD dd, mm", 0b11001111, 0b00000001, []z80op{&mrPC{}, &mrPC{f: ldDDmm}}, nil},
	{"ADD HL,ss", 0b11001111, 0b00001001, []z80op{}, addHLss},
	{"INC ss", 0b11001111, 0b00000011, []z80op{&exec{l: 2, f: incSS}}, nil},
	{"DEC ss", 0b11001111, 0b00001011, []z80op{&exec{l: 2, f: decSS}}, nil},
	{"POP ss", 0b11001111, 0b11000001, []z80op{}, popSS},
	{"PUSH ss", 0b11001111, 0b11000101, []z80op{}, pushSS},

	{"LD r, n", 0b11000111, 0b00000110, []z80op{&mrPC{f: ldRn}}, nil},
	{"LD r, r", 0b11000000, 0b01000000, []z80op{}, ldRr},
	{"LD r, (HL)", 0b11000111, 0b01000110, []z80op{}, ldRhl},
	{"LD (HL), r", 0b11111000, 0b01110000, []z80op{}, ldHLr},
	{"INC r", 0b11000111, 0b0000100, []z80op{}, incR},
	{"DEC r", 0b11000111, 0b0000101, []z80op{}, decR},
	{"ADD A, r", 0b11111000, 0b10000000, []z80op{}, addAr},
	{"ADC A, r", 0b11111000, 0b10001000, []z80op{}, adcAr},
	{"SUB A, r", 0b11111000, 0b10010000, []z80op{}, subAr},
	{"SUC A, r", 0b11111000, 0b10011000, []z80op{}, sbcAr},
	{"AND r", 0b11111000, 0b10100000, []z80op{}, andAr},
	{"OR r", 0b11111000, 0b10110000, []z80op{}, orAr},
	{"XOR r", 0b11111000, 0b10101000, []z80op{}, xorAr},
	{"CP r", 0b11111000, 0b10111000, []z80op{}, cpR},

	{"RET cc", 0b11000111, 0b11000000, []z80op{&exec{l: 1, f: retCC}}, nil},
	{"JP cc, nn", 0b11000111, 0b11000010, []z80op{&mrPC{}, &mrPC{f: jpCC}}, nil},
	{"CALL cc, nn", 0b11000111, 0b11000100, []z80op{&mrPC{}, &mrPC{f: callCC}}, nil},
	{"RST p", 0b11000111, 0b11000111, []z80op{&exec{l: 1, f: rstP}}, nil},
	{"CALL cc, nn", 0xFF, 0xCD, []z80op{&mrPC{}, &mrPC{f: call}}, nil},

	// {"", 0xFF, 0x,[]z80op{},nil},
	{"NOP", 0xFF, 0x00, []z80op{}, nil},
	{"DAA", 0xFF, 0x27, []z80op{}, daa},
	{"CPL", 0xFF, 0x2f, []z80op{}, cpl},
	{"SCF", 0xFF, 0x37, []z80op{}, scf},
	{"CCF", 0xFF, 0x3F, []z80op{}, ccf},
	{"HALT", 0xFF, 0x76, []z80op{}, halt},
	{"RET", 0xFF, 0xC9, []z80op{&mrPC{}, &mrPC{f: ret}}, nil},

	{"INC (HL)", 0xFF, 0x34, []z80op{&exec{l: 1, f: incHL}}, nil},
	{"DEC (HL)", 0xFF, 0x35, []z80op{&exec{l: 1, f: decHL}}, nil},
	{"ADD A, (HL)", 0xFF, 0x86, []z80op{}, addAhl},
	{"ADC A, (HL)", 0xFF, 0x8e, []z80op{}, adcAhl},
	{"SUB A, (HL)", 0xFF, 0x96, []z80op{}, subAhl},
	{"SBC A, (HL)", 0xFF, 0x9e, []z80op{}, sbcAhl},
	{"AND (HL)", 0xFF, 0xA6, []z80op{}, andAhl},
	{"OR (HL)", 0xFF, 0xB6, []z80op{}, orAhl},
	{"XOR (HL)", 0xFF, 0xAE, []z80op{}, xorAhl},
	{"CP (HL)", 0xFF, 0xBE, []z80op{}, cpHl},
	{"ADD A, n", 0xFF, 0xc6, []z80op{&mrPC{f: func(cpu *z80, data []uint8) { cpu.addA(data[1]) }}}, nil},
	{"ADC A, (HL)", 0xFF, 0xCE, []z80op{&mrPC{f: func(cpu *z80, data []uint8) { cpu.adcA(data[1]) }}}, nil},
	{"SBC A, (HL)", 0xFF, 0xDE, []z80op{&mrPC{f: func(cpu *z80, data []uint8) { cpu.sbcA(data[1]) }}}, nil},
	{"SUB n", 0xFF, 0xD6, []z80op{&mrPC{f: func(cpu *z80, data []uint8) { cpu.subA(data[1]) }}}, nil},
	{"AND n", 0xFF, 0xE6, []z80op{&mrPC{f: func(cpu *z80, data []uint8) { cpu.and(data[1]) }}}, nil},
	{"OR n", 0xFF, 0xF6, []z80op{&mrPC{f: func(cpu *z80, data []uint8) { cpu.or(data[1]) }}}, nil},

	{"LD A,(BC)", 0xFF, 0x0A, []z80op{}, ldAbc},
	{"LD A,(DE)", 0xFF, 0x1A, []z80op{}, ldAde},
	{"LD (BC), A", 0xFF, 0x02, []z80op{}, ldBCa},
	{"LD (BC), A", 0xFF, 0x12, []z80op{}, ldDEa},
	{"LD (nn), HL", 0xFF, 0x22, []z80op{&mrPC{}, &mrPC{f: ldNNhl}}, nil},
	{"LD (nn), A", 0xFF, 0x32, []z80op{&mrPC{}, &mrPC{f: ldNNa}}, nil},
	{"LD HL, (nn)", 0xFF, 0x2a, []z80op{&mrPC{}, &mrPC{f: ldHLnn}}, nil},
	{"LD (HL), n", 0xFF, 0x36, []z80op{&mrPC{f: ldHLn}}, nil},
	{"LD A, (nn)", 0xFF, 0x3a, []z80op{&mrPC{}, &mrPC{f: ldAnn}}, nil},

	{"EX AF, AF'", 0xFF, 0x08, []z80op{}, exafaf},
	{"EXX'", 0xFF, 0xD9, []z80op{}, exx},

	{"DJNZ e", 0xFF, 0x10, []z80op{&mrPC{f: djnz}}, nil},
	{"JR e", 0xFF, 0x18, []z80op{&mrPC{}, &exec{l: 5, f: jr}}, nil},
	{"JRNZ e", 0xFF, 0x20, []z80op{&mrPC{f: jrnz}}, nil},
	{"JRZ e", 0xFF, 0x28, []z80op{&mrPC{f: jrz}}, nil},
	{"JRNC e", 0xFF, 0x30, []z80op{&mrPC{f: jrnc}}, nil},
	{"JRC e", 0xFF, 0x38, []z80op{&mrPC{f: jrc}}, nil},

	{"JP nn", 0xFF, 0xC3, []z80op{&mrPC{}, &mrPC{f: func(cpu *z80, mem []uint8) { cpu.regs.PC = toWord(mem[1], mem[2]) }}}, nil},

	{"RLCA", 0xFF, 0x07, []z80op{}, rlca},
	{"RLA", 0xFF, 0x17, []z80op{}, rla},
	{"RRCA", 0xFF, 0x0F, []z80op{}, rrca},
	{"RRA", 0xFF, 0x1F, []z80op{}, rra},

	{"OUT (n), A", 0xFF, 0xD3, []z80op{&mrPC{f: outNa}}, nil},
	{"IN A, (n)", 0xFF, 0xDB, []z80op{&mrPC{f: inAn}}, nil},

	{"EX (SP), IX", 0xFF, 0xE3, []z80op{}, exSP},
	{"JP HL", 0xFF, 0xE9, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.PC = cpu.regs.HL.Get() }},
	{"EX DE, HL", 0xFF, 0xEB, []z80op{}, exDEhl},

	{"XOR *", 0xFF, 0xEE, []z80op{&mrPC{f: func(cpu *z80, mem []uint8) { cpu.xor(mem[1]) }}}, nil},
	{"DI", 0xFF, 0xF3, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.IFF1 = false; cpu.regs.IFF2 = false }},
	{"EI", 0xFF, 0xFb, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.IFF1 = true; cpu.regs.IFF2 = true }},
	{"LD SP, HL", 0xFF, 0xF9, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.regs.SP.Set(cpu.regs.HL.Get()) }}}, nil},
	{"CP *", 0xFF, 0xFe, []z80op{&mrPC{f: func(cpu *z80, mem []uint8) { cpu.cp(mem[1]) }}}, nil},

	{"CB", 0xFF, 0xCB, []z80op{}, decodeCB},
	{"DD", 0xFF, 0xDD, []z80op{}, decodeDD},
	{"ED", 0xFF, 0xED, []z80op{}, decodeED},
	{"ED", 0xFF, 0xFD, []z80op{}, decodeFD},
}

var z80OpsCodeTableCB = []*opCode{
	{"RLC r", 0b11111000, 0b00000000, []z80op{}, cbR},
	{"RLC (HL)", 0xFF, 0x06, []z80op{}, cbHL},
	{"RRC r", 0b11111000, 0b00001000, []z80op{}, cbR},
	{"RRC (HL)", 0xFF, 0x0e, []z80op{}, cbHL},

	{"RLC r", 0b11111000, 0b00010000, []z80op{}, cbR},
	{"RLC (HL)", 0xFF, 0x16, []z80op{}, cbHL},
	{"RR r", 0b11111000, 0b00011000, []z80op{}, cbR},
	{"RR (HL)", 0xFF, 0x1e, []z80op{}, cbHL},

	{"SLA r", 0b11111000, 0b00100000, []z80op{}, cbR},
	{"SLA (HL)", 0xFF, 0x26, []z80op{}, cbHL},
	{"SRA r", 0b11111000, 0b00101000, []z80op{}, cbR},
	{"SRA (HL)", 0xFF, 0x2e, []z80op{}, cbHL},

	{"SLL r", 0b11111000, 0b00110000, []z80op{}, cbR},
	{"SLL (HL)", 0xFF, 0x36, []z80op{}, cbHL},
	{"SRL r", 0b11111000, 0b00111000, []z80op{}, cbR},
	{"SRL (HL)", 0xFF, 0x3e, []z80op{}, cbHL},

	{"BIT b, r", 0b11000000, 0b01000000, []z80op{}, bit},
	{"BIT b, (HL)", 0b11000111, 0b01000110, []z80op{}, bitHL},

	{"RES b, r", 0b11000000, 0b10000000, []z80op{}, res},
	{"RES b, (HL)", 0b11000111, 0b10000110, []z80op{}, resHL},

	{"SET b, r", 0b11000000, 0b11000000, []z80op{}, set},
	{"SET b, (HL)", 0b11000111, 0b11000110, []z80op{}, setHL},
}

var z80OpsCodeTableDD = []*opCode{
	{"ADD IX, rr", 0b11001111, 0b00001001, []z80op{&exec{l: 7, f: addIXY}}, nil},
	{"LD IX, nn", 0xFF, 0x21, []z80op{&mrPC{}, &mrPC{f: func(cpu *z80, data []uint8) { cpu.regs.IXH = data[2]; cpu.regs.IXL = data[1] }}}, nil},
	{"LD (nn), IX", 0xFF, 0x22, []z80op{&mrPC{}, &mrPC{f: ldNNIXY}}, nil},
	{"INC IX", 0xFF, 0x23, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.regs.IX.Set(cpu.regs.IX.Get() + 1) }}}, nil},
	{"INC IXH", 0xFF, 0x24, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.incR(&cpu.regs.IXH) }}}, nil},
	{"DEC IXH", 0xFF, 0x25, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.decR(&cpu.regs.IXH) }}}, nil},
	{"LD IXH, n", 0xFF, 0x26, []z80op{&mrPC{f: func(cpu *z80, data []uint8) { cpu.regs.IXH = data[1] }}}, nil},
	{"LD IX, nn", 0xFF, 0x2A, []z80op{&mrPC{}, &mrPC{f: ldIXYnn}}, nil},
	{"DEC IX", 0xFF, 0x2B, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.regs.IX.Set(cpu.regs.IX.Get() - 1) }}}, nil},
	{"INC IXL", 0xFF, 0x2C, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.incR(&cpu.regs.IXL) }}}, nil},
	{"DEC IXL", 0xFF, 0x2D, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.decR(&cpu.regs.IXL) }}}, nil},
	{"LD IXL, n", 0xFF, 0x2E, []z80op{&mrPC{f: func(cpu *z80, data []uint8) { cpu.regs.IXL = data[1] }}}, nil},
	{"INC (IX+d)", 0xFF, 0x34, []z80op{&mrPC{}, &exec{l: 7, f: incIXYd}}, nil},
	{"DEC (IX+d)", 0xFF, 0x35, []z80op{&mrPC{}, &exec{l: 7, f: decIXYd}}, nil},
	{"LD (IX+d), n", 0xFF, 0x36, []z80op{&mrPC{}, &mrPC{}, &exec{l: 2, f: ldIXYdN}}, nil},

	{"LD B, IXH", 0xFF, 0x44, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.B = cpu.regs.IXH }},
	{"LD B, IXL", 0xFF, 0x45, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.B = cpu.regs.IXL }},
	{"LD C, IXH", 0xFF, 0x4C, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.C = cpu.regs.IXH }},
	{"LD C, IXL", 0xFF, 0x4D, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.C = cpu.regs.IXL }},
	{"LD D, IXH", 0xFF, 0x54, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.D = cpu.regs.IXH }},
	{"LD D, IXL", 0xFF, 0x55, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.D = cpu.regs.IXL }},
	{"LD E, IXH", 0xFF, 0x5C, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.E = cpu.regs.IXH }},
	{"LD E, IXL", 0xFF, 0x5D, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.E = cpu.regs.IXL }},
	{"LD A, IXH", 0xFF, 0x7C, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.A = cpu.regs.IXH }},
	{"LD A, IXL", 0xFF, 0x7D, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.A = cpu.regs.IXL }},

	{"LD IXH, r", 0b11111000, 0b01100000, []z80op{}, ldIXYHr},
	{"LD IXH, r", 0b11111000, 0b01101000, []z80op{}, ldIXYLr},
	{"LD r, (IX+d)", 0b11000111, 0b01000110, []z80op{&mrPC{}, &exec{l: 5, f: ldRixyD}}, nil},
	{"LD (IX+d), r", 0b11111000, 0b01110000, []z80op{&mrPC{}, &exec{l: 5, f: ldIXYdR}}, nil},

	{"ADD A, IXH", 0xFF, 0x84, []z80op{}, func(cpu *z80, u []uint8) { cpu.addA(cpu.regs.IXH) }},
	{"ADD A, IXL", 0xFF, 0x85, []z80op{}, func(cpu *z80, u []uint8) { cpu.addA(cpu.regs.IXL) }},
	{"ADC A, IXH", 0xFF, 0x8C, []z80op{}, func(cpu *z80, u []uint8) { cpu.adcA(cpu.regs.IXH) }},
	{"ADC A, IXL", 0xFF, 0x8D, []z80op{}, func(cpu *z80, u []uint8) { cpu.adcA(cpu.regs.IXL) }},
	{"SUB A, IXH", 0xFF, 0x94, []z80op{}, func(cpu *z80, u []uint8) { cpu.subA(cpu.regs.IXH) }},
	{"SUB A, IXL", 0xFF, 0x95, []z80op{}, func(cpu *z80, u []uint8) { cpu.subA(cpu.regs.IXL) }},
	{"SBC A, IXH", 0xFF, 0x9C, []z80op{}, func(cpu *z80, u []uint8) { cpu.sbcA(cpu.regs.IXH) }},
	{"SBC A, IXL", 0xFF, 0x9D, []z80op{}, func(cpu *z80, u []uint8) { cpu.sbcA(cpu.regs.IXL) }},
	{"AND A, IXH", 0xFF, 0xA4, []z80op{}, func(cpu *z80, u []uint8) { cpu.and(cpu.regs.IXH) }},
	{"AND A, IXL", 0xFF, 0xA5, []z80op{}, func(cpu *z80, u []uint8) { cpu.and(cpu.regs.IXL) }},
	{"XOR A, IXH", 0xFF, 0xAC, []z80op{}, func(cpu *z80, u []uint8) { cpu.xor(cpu.regs.IXH) }},
	{"XOR A, IXL", 0xFF, 0xAD, []z80op{}, func(cpu *z80, u []uint8) { cpu.xor(cpu.regs.IXL) }},
	{"OR A, IXH", 0xFF, 0xB4, []z80op{}, func(cpu *z80, u []uint8) { cpu.or(cpu.regs.IXH) }},
	{"OR A, IXL", 0xFF, 0xB5, []z80op{}, func(cpu *z80, u []uint8) { cpu.or(cpu.regs.IXL) }},
	{"CP A, IXH", 0xFF, 0xBC, []z80op{}, func(cpu *z80, u []uint8) { cpu.cp(cpu.regs.IXH) }},
	{"CP A, IXL", 0xFF, 0xBD, []z80op{}, func(cpu *z80, u []uint8) { cpu.cp(cpu.regs.IXL) }},

	{"ADD A, (IX+d)", 0xFF, 0x86, []z80op{&mrPC{}, &exec{l: 5, f: addAixyD}}, nil},
	{"ADC A, (IX+d)", 0xFF, 0x8E, []z80op{&mrPC{}, &exec{l: 5, f: adcAixyD}}, nil},
	{"SUB A, (IX+d)", 0xFF, 0x96, []z80op{&mrPC{}, &exec{l: 5, f: subAixyD}}, nil},
	{"SBC A, (IX+d)", 0xFF, 0x9E, []z80op{&mrPC{}, &exec{l: 5, f: sbcAixyD}}, nil},
	{"AND A, (IX+d)", 0xFF, 0xA6, []z80op{&mrPC{}, &exec{l: 5, f: andAixyD}}, nil},
	{"XOR A, (IX+d)", 0xFF, 0xAE, []z80op{&mrPC{}, &exec{l: 5, f: xorAixyD}}, nil},
	{"OR A, (IX+d)", 0xFF, 0xB6, []z80op{&mrPC{}, &exec{l: 5, f: orAixyD}}, nil},
	{"CP A, (IX+d)", 0xFF, 0xBE, []z80op{&mrPC{}, &exec{l: 5, f: cpAixyD}}, nil},

	{"CB", 0xFF, 0xCB, []z80op{&mrPC{f: decodeDDCB}}, nil},

	{"POP IX", 0xFF, 0xE1, []z80op{}, func(cpu *z80, u []uint8) { cpu.popFromStack(func(cpu *z80, data uint16) { cpu.regs.IX.Set(data) }) }},
	{"EX (SP), IX", 0xFF, 0xE3, []z80op{}, exSP},
	{"PUSH IX", 0xFF, 0xE5, []z80op{}, func(cpu *z80, u []uint8) { cpu.pushToStack(cpu.regs.IX.Get(), nil) }},
	{"JP IX", 0xFF, 0xE9, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.PC = cpu.regs.IX.Get() }},
	{"LD SP, IX", 0xFF, 0xF9, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.SP.Set(cpu.regs.IX.Get()) }},
}

var z80OpsCodeTableFD = []*opCode{
	{"ADD IY, rr", 0b11001111, 0b00001001, []z80op{&exec{l: 7, f: addIY}}, nil},
	{"LD IY, nn", 0xFF, 0x21, []z80op{&mrPC{}, &mrPC{f: func(cpu *z80, data []uint8) { cpu.regs.IYH = data[2]; cpu.regs.IYL = data[1] }}}, nil},
	{"LD (nn), IY", 0xFF, 0x22, []z80op{&mrPC{}, &mrPC{f: ldNNIXY}}, nil},
	{"INC IY", 0xFF, 0x23, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.regs.IY.Set(cpu.regs.IY.Get() + 1) }}}, nil},
	{"INC IYH", 0xFF, 0x24, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.incR(&cpu.regs.IYH) }}}, nil},
	{"DEC IYH", 0xFF, 0x25, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.decR(&cpu.regs.IYH) }}}, nil},
	{"LD IYH, n", 0xFF, 0x26, []z80op{&mrPC{f: func(cpu *z80, data []uint8) { cpu.regs.IYH = data[1] }}}, nil},
	{"LD IY, nn", 0xFF, 0x2A, []z80op{&mrPC{}, &mrPC{f: ldIXYnn}}, nil},
	{"DEC IY", 0xFF, 0x2B, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.regs.IY.Set(cpu.regs.IY.Get() - 1) }}}, nil},
	{"INC IYL", 0xFF, 0x2C, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.incR(&cpu.regs.IYL) }}}, nil},
	{"DEC IYL", 0xFF, 0x2D, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.decR(&cpu.regs.IYL) }}}, nil},
	{"LD IYL, n", 0xFF, 0x2E, []z80op{&mrPC{f: func(cpu *z80, data []uint8) { cpu.regs.IYL = data[1] }}}, nil},
	{"INC (IY+d)", 0xFF, 0x34, []z80op{&mrPC{}, &exec{l: 7, f: incIXYd}}, nil},
	{"DEC (IY+d)", 0xFF, 0x35, []z80op{&mrPC{}, &exec{l: 7, f: decIXYd}}, nil},
	{"LD (IY+d), n", 0xFF, 0x36, []z80op{&mrPC{}, &mrPC{}, &exec{l: 2, f: ldIXYdN}}, nil},

	{"LD B, IYH", 0xFF, 0x44, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.B = cpu.regs.IYH }},
	{"LD B, IYL", 0xFF, 0x45, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.B = cpu.regs.IYL }},
	{"LD C, IYH", 0xFF, 0x4C, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.C = cpu.regs.IYH }},
	{"LD C, IYL", 0xFF, 0x4D, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.C = cpu.regs.IYL }},
	{"LD D, IYH", 0xFF, 0x54, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.D = cpu.regs.IYH }},
	{"LD D, IYL", 0xFF, 0x55, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.D = cpu.regs.IYL }},
	{"LD E, IYH", 0xFF, 0x5C, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.E = cpu.regs.IYH }},
	{"LD E, IYL", 0xFF, 0x5D, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.E = cpu.regs.IYL }},
	{"LD A, IYH", 0xFF, 0x7C, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.A = cpu.regs.IYH }},
	{"LD A, IYL", 0xFF, 0x7D, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.A = cpu.regs.IYL }},

	{"LD IYH, r", 0b11111000, 0b01100000, []z80op{}, ldIXYHr},
	{"LD IYL, r", 0b11111000, 0b01101000, []z80op{}, ldIXYLr},
	{"LD r, (IY+d)", 0b11000111, 0b01000110, []z80op{&mrPC{}, &exec{l: 5, f: ldRixyD}}, nil},
	{"LD (IY+d), r", 0b11111000, 0b01110000, []z80op{&mrPC{}, &exec{l: 5, f: ldIXYdR}}, nil},

	{"ADD A, IYH", 0xFF, 0x84, []z80op{}, func(cpu *z80, u []uint8) { cpu.addA(cpu.regs.IYH) }},
	{"ADD A, IYL", 0xFF, 0x85, []z80op{}, func(cpu *z80, u []uint8) { cpu.addA(cpu.regs.IYL) }},
	{"ADC A, IYH", 0xFF, 0x8C, []z80op{}, func(cpu *z80, u []uint8) { cpu.adcA(cpu.regs.IYH) }},
	{"ADC A, IYL", 0xFF, 0x8D, []z80op{}, func(cpu *z80, u []uint8) { cpu.adcA(cpu.regs.IYL) }},
	{"SUB A, IYH", 0xFF, 0x94, []z80op{}, func(cpu *z80, u []uint8) { cpu.subA(cpu.regs.IYH) }},
	{"SUB A, IYL", 0xFF, 0x95, []z80op{}, func(cpu *z80, u []uint8) { cpu.subA(cpu.regs.IYL) }},
	{"SBC A, IYH", 0xFF, 0x9C, []z80op{}, func(cpu *z80, u []uint8) { cpu.sbcA(cpu.regs.IYH) }},
	{"SBC A, IYL", 0xFF, 0x9D, []z80op{}, func(cpu *z80, u []uint8) { cpu.sbcA(cpu.regs.IYL) }},
	{"AND A, IYH", 0xFF, 0xA4, []z80op{}, func(cpu *z80, u []uint8) { cpu.and(cpu.regs.IYH) }},
	{"AND A, IYL", 0xFF, 0xA5, []z80op{}, func(cpu *z80, u []uint8) { cpu.and(cpu.regs.IYL) }},
	{"XOR A, IYH", 0xFF, 0xAC, []z80op{}, func(cpu *z80, u []uint8) { cpu.xor(cpu.regs.IYH) }},
	{"XOR A, IYL", 0xFF, 0xAD, []z80op{}, func(cpu *z80, u []uint8) { cpu.xor(cpu.regs.IYL) }},
	{"OR A, IYH", 0xFF, 0xB4, []z80op{}, func(cpu *z80, u []uint8) { cpu.or(cpu.regs.IYH) }},
	{"OR A, IYL", 0xFF, 0xB5, []z80op{}, func(cpu *z80, u []uint8) { cpu.or(cpu.regs.IYL) }},
	{"CP A, IYH", 0xFF, 0xBC, []z80op{}, func(cpu *z80, u []uint8) { cpu.cp(cpu.regs.IYH) }},
	{"CP A, IYL", 0xFF, 0xBD, []z80op{}, func(cpu *z80, u []uint8) { cpu.cp(cpu.regs.IYL) }},

	{"ADD A, (IY+d)", 0xFF, 0x86, []z80op{&mrPC{}, &exec{l: 5, f: addAixyD}}, nil},
	{"ADC A, (IY+d)", 0xFF, 0x8E, []z80op{&mrPC{}, &exec{l: 5, f: adcAixyD}}, nil},
	{"SUB A, (IY+d)", 0xFF, 0x96, []z80op{&mrPC{}, &exec{l: 5, f: subAixyD}}, nil},
	{"SBC A, (IY+d)", 0xFF, 0x9E, []z80op{&mrPC{}, &exec{l: 5, f: sbcAixyD}}, nil},
	{"AND A, (IY+d)", 0xFF, 0xA6, []z80op{&mrPC{}, &exec{l: 5, f: andAixyD}}, nil},
	{"XOR A, (IY+d)", 0xFF, 0xAE, []z80op{&mrPC{}, &exec{l: 5, f: xorAixyD}}, nil},
	{"OR A, (IY+d)", 0xFF, 0xB6, []z80op{&mrPC{}, &exec{l: 5, f: orAixyD}}, nil},
	{"CP A, (IY+d)", 0xFF, 0xBE, []z80op{&mrPC{}, &exec{l: 5, f: cpAixyD}}, nil},

	{"CB", 0xFF, 0xCB, []z80op{&mrPC{f: decodeFDCB}}, nil},

	{"POP IY", 0xFF, 0xE1, []z80op{}, func(cpu *z80, u []uint8) { cpu.popFromStack(func(cpu *z80, data uint16) { cpu.regs.IY.Set(data) }) }},
	{"EX (SP), IY", 0xFF, 0xE3, []z80op{}, exSP},
	{"PUSH IY", 0xFF, 0xE5, []z80op{}, func(cpu *z80, u []uint8) { cpu.pushToStack(cpu.regs.IY.Get(), nil) }},
	{"JP IY", 0xFF, 0xE9, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.PC = cpu.regs.IY.Get() }},
	{"LD SP, IY", 0xFF, 0xF9, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.SP.Set(cpu.regs.IY.Get()) }},
}

var z80OpsCodeTableDDCB = []*opCode{
	{"RLC (IX+d), r", 0b11111000, 0b00000000, []z80op{}, cbIXYdr},
	{"RLC (IX+d)", 0xFF, 0x06, []z80op{}, cbIXYd},
	{"RRC (IX+d), r", 0b11111000, 0b00001000, []z80op{}, cbIXYdr},
	{"RRC (IX+d)", 0xFF, 0x0e, []z80op{}, cbIXYd},

	{"RLC (IX+d), r", 0b11111000, 0b00010000, []z80op{}, cbIXYdr},
	{"RLC (IX+d)", 0xFF, 0x16, []z80op{}, cbIXYd},
	{"RR (IX+d), r", 0b11111000, 0b00011000, []z80op{}, cbIXYdr},
	{"RR (IX+d)", 0xFF, 0x1e, []z80op{}, cbIXYd},

	{"SLA (IX+d), r", 0b11111000, 0b00100000, []z80op{}, cbIXYdr},
	{"SLA (IX+d)", 0xFF, 0x26, []z80op{}, cbIXYd},
	{"SRA (IX+d), r", 0b11111000, 0b00101000, []z80op{}, cbIXYdr},
	{"SRA (IX+d)", 0xFF, 0x2e, []z80op{}, cbIXYd},

	{"SLL (IX+d), r", 0b11111000, 0b00110000, []z80op{}, cbIXYdr},
	{"SLL (IX+d)", 0xFF, 0x36, []z80op{}, cbIXYd},
	{"SRL (IX+d), r", 0b11111000, 0b00111000, []z80op{}, cbIXYdr},
	{"SRL (IX+d)", 0xFF, 0x3e, []z80op{}, cbIXYd},

	{"BIT b, (IX+d), r", 0b11000000, 0b01000000, []z80op{}, bitIXYd},
	{"BIT b, (IX+d)", 0b11000111, 0b01000110, []z80op{}, bitIXYd},

	{"RES b, (IX+d), r", 0b11000000, 0b10000000, []z80op{}, resIXYdR},
	{"RES b, (IX+d)", 0b11000111, 0b10000110, []z80op{}, resIXYd},

	{"SET b, (IX+d), r", 0b11000000, 0b11000000, []z80op{}, setIXYdR},
	{"SET b, (IX+d)", 0b11000111, 0b11000110, []z80op{}, setIXYd},
}

var z80OpsCodeTableFDCB = []*opCode{
	{"RLC (IY+d), r", 0b11111000, 0b00000000, []z80op{}, cbIXYdr},
	{"RLC (IY+d)", 0xFF, 0x06, []z80op{}, cbIXYd},
	{"RRC (IY+d), r", 0b11111000, 0b00001000, []z80op{}, cbIXYdr},
	{"RRC (IY+d)", 0xFF, 0x0e, []z80op{}, cbIXYd},

	{"RLC (IY+d), r", 0b11111000, 0b00010000, []z80op{}, cbIXYdr},
	{"RLC (IY+d)", 0xFF, 0x16, []z80op{}, cbIXYd},
	{"RR (IY+d), r", 0b11111000, 0b00011000, []z80op{}, cbIXYdr},
	{"RR (IY+d)", 0xFF, 0x1e, []z80op{}, cbIXYd},

	{"SLA (IY+d), r", 0b11111000, 0b00100000, []z80op{}, cbIXYdr},
	{"SLA (IY+d)", 0xFF, 0x26, []z80op{}, cbIXYd},
	{"SRA (IY+d), r", 0b11111000, 0b00101000, []z80op{}, cbIXYdr},
	{"SRA (IY+d)", 0xFF, 0x2e, []z80op{}, cbIXYd},

	{"SLL (IY+d), r", 0b11111000, 0b00110000, []z80op{}, cbIXYdr},
	{"SLL (IY+d)", 0xFF, 0x36, []z80op{}, cbIXYd},
	{"SRL (IY+d), r", 0b11111000, 0b00111000, []z80op{}, cbIXYdr},
	{"SRL (IY+d)", 0xFF, 0x3e, []z80op{}, cbIXYd},

	{"BIT b, (IY+d), r", 0b11000000, 0b01000000, []z80op{}, bitIXYd},
	{"BIT b, (IY+d)", 0b11000111, 0b01000110, []z80op{}, bitIXYd},

	{"RES b, (IY+d), r", 0b11000000, 0b10000000, []z80op{}, resIXYdR},
	{"RES b, (IY+d)", 0b11000111, 0b10000110, []z80op{}, resIXYd},

	{"SET b, (IY+d), r", 0b11000000, 0b11000000, []z80op{}, setIXYdR},
	{"SET b, (IY+d)", 0b11000111, 0b11000110, []z80op{}, setIXYd},
}
var z80OpsCodeTableED = []*opCode{
	{"IN r, (n)", 0b11000111, 0b01000000, []z80op{}, inRc},
	{"IN (c)", 0xFF, 0x70, []z80op{}, inC},
	{"OUT (c), r", 0b11000111, 0b01000001, []z80op{}, outCr},
	{"OUT (c), 0", 0xFF, 0x71, []z80op{}, outC0},

	{"SBC HL, BC", 0xFF, 0x42, []z80op{&exec{l: 4, f: func(cpu *z80, u []uint8) { cpu.sbcHL(cpu.regs.BC.Get()) }}}, nil},
	{"SBC HL, DE", 0xFF, 0x52, []z80op{&exec{l: 4, f: func(cpu *z80, u []uint8) { cpu.sbcHL(cpu.regs.DE.Get()) }}}, nil},
	{"SBC HL, HL", 0xFF, 0x62, []z80op{&exec{l: 4, f: func(cpu *z80, u []uint8) { cpu.sbcHL(cpu.regs.HL.Get()) }}}, nil},
	{"SBC HL, SP", 0xFF, 0x72, []z80op{&exec{l: 4, f: func(cpu *z80, u []uint8) { cpu.sbcHL(cpu.regs.SP.Get()) }}}, nil},

	{"LD (nn), dd", 0b11001111, 0b01000011, []z80op{&mrPC{}, &mrPC{f: ldNNdd}}, nil},
	{"NEG", 0b11000111, 0b01000100, []z80op{}, func(cpu *z80, u []uint8) { n := cpu.regs.A; cpu.regs.A = 0; cpu.subA(n) }},
	{"RETN", 0b11000111, 0b01000101, []z80op{}, func(cpu *z80, u []uint8) { cpu.popFromStack(func(cpu *z80, data uint16) { cpu.regs.PC = data }) }},

	{"IM 0", 0xFF, 0x46, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.InterruptsMode = 0 }},
	{"IM 0", 0xFF, 0x66, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.InterruptsMode = 0 }},
	{"IM 1", 0xFF, 0x56, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.InterruptsMode = 1 }},
	{"IM 2", 0xFF, 0xE5, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.InterruptsMode = 2 }},
	{"IM 0/1", 0xFF, 0x4E, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.InterruptsMode = 1 - cpu.regs.InterruptsMode }},
	{"IM 2", 0xFF, 0x5E, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.InterruptsMode = 2 }},
	{"IM 0/1", 0xFF, 0x6E, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.InterruptsMode = 1 - cpu.regs.InterruptsMode }},
	{"IM 1", 0xFF, 0x76, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.InterruptsMode = 1 }},
	{"IM 2", 0xFF, 0x7E, []z80op{}, func(cpu *z80, u []uint8) { cpu.regs.InterruptsMode = 2 }},

	{"LD I, A", 0xFF, 0x47, []z80op{&exec{l: 4, f: func(cpu *z80, u []uint8) { cpu.regs.I = cpu.regs.A }}}, nil},
	{"LD R, A", 0xFF, 0x4F, []z80op{&exec{l: 4, f: func(cpu *z80, u []uint8) { cpu.regs.R = cpu.regs.A }}}, nil},

	{"LD A, I", 0xFF, 0x57, []z80op{&exec{l: 1, f: ldAi}}, nil},
	{"LD A, R", 0xFF, 0x5F, []z80op{&exec{l: 1, f: ldAr}}, nil},

	{"ADC HL, BC", 0xFF, 0x4a, []z80op{&exec{l: 4, f: func(cpu *z80, u []uint8) { cpu.adcHL(cpu.regs.BC.Get()) }}}, nil},
	{"ADC HL, DE", 0xFF, 0x5a, []z80op{&exec{l: 4, f: func(cpu *z80, u []uint8) { cpu.adcHL(cpu.regs.DE.Get()) }}}, nil},
	{"ADC HL, HL", 0xFF, 0x6a, []z80op{&exec{l: 4, f: func(cpu *z80, u []uint8) { cpu.adcHL(cpu.regs.HL.Get()) }}}, nil},
	{"ADC HL, SP", 0xFF, 0x7a, []z80op{&exec{l: 4, f: func(cpu *z80, u []uint8) { cpu.adcHL(cpu.regs.SP.Get()) }}}, nil},

	{"LD (nn), dd", 0b11001111, 0b01001011, []z80op{&mrPC{}, &mrPC{f: ldDDnn}}, nil},

	{"RDD", 0xFF, 0x67, []z80op{}, rrd},
	{"RDD", 0xFF, 0x6f, []z80op{}, rld},

	{"LDI", 0xFF, 0xA0, []z80op{}, ldi},
	{"CPI", 0xFF, 0xA1, []z80op{}, cpi},
	{"INI", 0xFF, 0xA2, []z80op{}, ini},
	{"OUTI", 0xFF, 0xA3, []z80op{}, outi},

	{"LDD", 0xFF, 0xA8, []z80op{}, ldd},
	{"CPD", 0xFF, 0xA9, []z80op{}, cpd},
	{"IND", 0xFF, 0xAA, []z80op{}, ind},
	{"OUTD", 0xFF, 0xAB, []z80op{}, outd},

	{"LDIR", 0xFF, 0xB0, []z80op{}, ldi},
	{"CPIR", 0xFF, 0xB1, []z80op{}, cpi},
	{"INIR", 0xFF, 0xB2, []z80op{}, ini},
	{"OTIR", 0xFF, 0xB3, []z80op{}, outi},

	{"LDDR", 0xFF, 0xB8, []z80op{}, ldd},
	{"CPDR", 0xFF, 0xB9, []z80op{}, cpd},
	{"INDR", 0xFF, 0xBA, []z80op{}, ind},
	{"OTDR", 0xFF, 0xBB, []z80op{}, outd},
}

func decodeCB(cpu *z80, mem []uint8) {
	cpu.fetched = nil
	cpu.scheduler = append([]z80op{&fetch{table: lookupCB}}, cpu.scheduler...)
}

func decodeDD(cpu *z80, mem []uint8) {
	cpu.fetched = nil
	cpu.indexIdx = 1
	cpu.scheduler = append([]z80op{&fetch{table: lookupDD}}, cpu.scheduler...)
}

func decodeED(cpu *z80, mem []uint8) {
	cpu.fetched = nil
	cpu.scheduler = append([]z80op{&fetch{table: lookupED}}, cpu.scheduler...)
}

func decodeFD(cpu *z80, mem []uint8) {
	cpu.fetched = nil
	cpu.indexIdx = 2
	cpu.scheduler = append([]z80op{&fetch{table: lookupFD}}, cpu.scheduler...)
}

func decodeDDCB(cpu *z80, mem []uint8) {
	cpu.scheduler = append([]z80op{&fetch{table: lookupDDCB}}, cpu.scheduler...)
}

func decodeFDCB(cpu *z80, mem []uint8) {
	cpu.scheduler = append([]z80op{&fetch{table: lookupFDCB}}, cpu.scheduler...)
}

func (o *opCode) String() string {
	if o == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s", o.name)
}

var lookup = make([]*opCode, 256)
var lookupCB = make([]*opCode, 256)
var lookupDD = make([]*opCode, 256)
var lookupED = make([]*opCode, 256)
var lookupFD = make([]*opCode, 256)
var lookupDDCB = make([]*opCode, 256)
var lookupFDCB = make([]*opCode, 256)

func init() {
	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTable {
			if (code & op.mask) == op.code {
				lookup[code] = op
			}
		}
	}

	// -----

	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTableCB {
			if (code & op.mask) == op.code {
				lookupCB[code] = op
			}
		}
	}

	// -----

	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTableDD {
			if (code & op.mask) == op.code {
				lookupDD[code] = op
			}
		}
	}
	// -----
	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTableED {
			if (code & op.mask) == op.code {
				lookupED[code] = op
			}
		}
	}
	// -----
	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTableFD {
			if (code & op.mask) == op.code {
				lookupFD[code] = op
			}
		}
	}
	// -----

	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTableDDCB {
			if (code & op.mask) == op.code {
				lookupDDCB[code] = op
			}
		}
	}

	// -----
	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTableFDCB {
			if (code & op.mask) == op.code {
				lookupFDCB[code] = op
			}
		}
	}

	// -----

	println("---------")
	println("                         CB                DD                DDCB              ED                FD                FDCB")
	for code := 0; code < 256; code++ {
		fmt.Printf("0x%02X - %-18v%-18v%-18v%-18v%-18v%-18v%-18v\n", code, lookup[code], lookupCB[code], lookupDD[code], lookupDDCB[code], lookupED[code], lookupFD[code], lookupFDCB[code])
	}
	println("---------")
}
