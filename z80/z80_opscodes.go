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

	// {"", 0xFF, 0x,[]z80op{},nil},
	{"NOP", 0xFF, 0x00, []z80op{}, nil},
	{"DAA", 0xFF, 0x27, []z80op{}, daa},
	{"CPL", 0xFF, 0x2f, []z80op{}, cpl},
	{"SCF", 0xFF, 0x37, []z80op{}, scf},
	{"CCF", 0xFF, 0x3F, []z80op{}, ccf},
	{"HALT", 0xFF, 0x76, []z80op{}, halt},
	{"RET", 0xFF, 0xC9, []z80op{&mrPC{}, &mrPC{f: ret}}, nil},
	{"CALL cc, nn", 0xFF, 0xCD, []z80op{&mrPC{}, &mrPC{f: call}}, nil},

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
	{"ADD n", 0xFF, 0xE6, []z80op{&mrPC{f: func(cpu *z80, data []uint8) { cpu.addA(data[1]) }}}, nil},
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

	{"CB", 0xFF, 0xCB, []z80op{}, decodeCB},
	{"DD", 0xFF, 0xDD, []z80op{}, decodeDD},
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
	{"ADD IX, rr", 0b11001111, 0b00001001, []z80op{&exec{l: 7, f: addIX}}, nil},
	{"LD IX, nn", 0xFF, 0x21, []z80op{&mrPC{}, &mrPC{f: func(cpu *z80, data []uint8) { cpu.regs.IXH = data[2]; cpu.regs.IXL = data[1] }}}, nil},
	{"LD (nn), IX", 0xFF, 0x22, []z80op{&mrPC{}, &mrPC{f: ldNNIX}}, nil},
	{"INC IX", 0xFF, 0x23, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.regs.IX.Set(cpu.regs.IX.Get() + 1) }}}, nil},
	{"INC IXH", 0xFF, 0x24, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.incR(&cpu.regs.IXH) }}}, nil},
	{"DEC IXH", 0xFF, 0x25, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.decR(&cpu.regs.IXH) }}}, nil},
	{"LD IXH, n", 0xFF, 0x26, []z80op{&mrPC{f: func(cpu *z80, data []uint8) { cpu.regs.IXH = data[1] }}}, nil},
	{"LD IX, nn", 0xFF, 0x2A, []z80op{&mrPC{}, &mrPC{f: ldIXnn}}, nil},
	{"DEC IX", 0xFF, 0x2B, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.regs.IX.Set(cpu.regs.IX.Get() - 1) }}}, nil},
	{"INC IXL", 0xFF, 0x2C, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.incR(&cpu.regs.IXL) }}}, nil},
	{"DEC IXL", 0xFF, 0x2D, []z80op{&exec{l: 2, f: func(cpu *z80, u []uint8) { cpu.decR(&cpu.regs.IXL) }}}, nil},
	{"LD IXL, n", 0xFF, 0x2E, []z80op{&mrPC{f: func(cpu *z80, data []uint8) { cpu.regs.IXL = data[1] }}}, nil},
	{"INC (IX+d)", 0xFF, 0x34, []z80op{&mrPC{}, &exec{l: 7, f: incIXd}}, nil},
	{"DEC (IX+d)", 0xFF, 0x35, []z80op{&mrPC{}, &exec{l: 7, f: decIXd}}, nil},
	{"LD (IX+d), n", 0xFF, 0x36, []z80op{&mrPC{}, &mrPC{}, &exec{l: 2, f: ldIXdN}}, nil},

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

	{"LD IXH, r", 0b11111000, 0b01100000, []z80op{}, ldIXHr},
	{"LD IXH, r", 0b11111000, 0b01101000, []z80op{}, ldIXLr},
	{"LD r, (IX+d)", 0b11000111, 0b01000110, []z80op{&mrPC{}, &exec{l: 5, f: ldRixD}}, nil},
	{"LD (IX+d), r", 0b11111000, 0b01110000, []z80op{&mrPC{}, &exec{l: 5, f: ldIXdR}}, nil},

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

	{"ADD A, (IX+d)", 0xFF, 0x86, []z80op{&mrPC{}, &exec{l: 5, f: addAixD}}, nil},
	{"ADC A, (IX+d)", 0xFF, 0x8E, []z80op{&mrPC{}, &exec{l: 5, f: adcAixD}}, nil},
	{"SUB A, (IX+d)", 0xFF, 0x96, []z80op{&mrPC{}, &exec{l: 5, f: subAixD}}, nil},
	{"SBC A, (IX+d)", 0xFF, 0x9E, []z80op{&mrPC{}, &exec{l: 5, f: sbcAixD}}, nil},
	{"AND A, (IX+d)", 0xFF, 0xA6, []z80op{&mrPC{}, &exec{l: 5, f: andAixD}}, nil},
	{"XOR A, (IX+d)", 0xFF, 0xAE, []z80op{&mrPC{}, &exec{l: 5, f: xorAixD}}, nil},
	{"OR A, (IX+d)", 0xFF, 0xB6, []z80op{&mrPC{}, &exec{l: 5, f: orAixD}}, nil},
	{"CP A, (IX+d)", 0xFF, 0xBE, []z80op{&mrPC{}, &exec{l: 5, f: cpAixD}}, nil},

	{"CB", 0xFF, 0xCB, []z80op{&mrPC{f: decodeDDCB}}, nil},
}

var z80OpsCodeTableDDCB = []*opCode{}

func decodeCB(cpu *z80, mem []uint8) {
	cpu.scheduler = append([]z80op{&fetch{table: lookupCB}}, cpu.scheduler...)
}

func decodeDD(cpu *z80, mem []uint8) {
	cpu.scheduler = append([]z80op{&fetch{table: lookupDD}}, cpu.scheduler...)
}

func decodeDDCB(cpu *z80, mem []uint8) {
	cpu.scheduler = append([]z80op{&fetch{table: lookupDDCB}}, cpu.scheduler...)
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
var lookupDDCB = make([]*opCode, 256)

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
		for _, op := range z80OpsCodeTableDDCB {
			if (code & op.mask) == op.code {
				lookupDDCB[code] = op
			}
		}
	}

	// -----

	println("---------")
	println("                      CB             DD             DDCB")
	for code := 0; code < 256; code++ {
		fmt.Printf("0x%02X - %-15v%-15v%-15v%-15v\n", code, lookup[code], lookupCB[code], lookupDD[code], lookupDDCB[code])
	}
	println("---------")
}
