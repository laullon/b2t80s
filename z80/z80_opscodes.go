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

	{"CB", 0xFF, 0xCB, []z80op{}, decodeCB},
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

func decodeCB(cpu *z80, mem []uint8) {
	cpu.scheduler = append([]z80op{&fetch{table: lookupCB}}, cpu.scheduler...)
}

func (o *opCode) String() string {
	return fmt.Sprintf("%s", o.name)
}

var lookup = make([]*opCode, 256)
var lookupCB = make([]*opCode, 256)

func init() {
	for i := 0; i < 256; i++ {
		code := uint8(i)
		for _, op := range z80OpsCodeTable {
			if (code & op.mask) == op.code {
				lookup[code] = op
			}
		}
	}
	println("------")
	println("lookup")
	println("------")
	for code, op := range lookup {
		fmt.Printf("0x%02X - %v\n", code, op)
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
	println("---------")
	println("lookup CB")
	println("---------")
	for code, op := range lookupCB {
		fmt.Printf("0x%02X - %v\n", code, op)
	}
}
