package lr35902

import "fmt"

type opCode struct {
	name       string
	mask, code byte
	len        byte
	ops        []lr35902op
	onFetch    lr35902f
}

var lr35902OpsCodeTable = []*opCode{
	{"LD dd, mm", 0b11001111, 0b00000001, 3, []lr35902op{&mrNNpc{f: ldDDmm}}, nil},
	{"ADD HL,ss", 0b11001111, 0b00001001, 1, []lr35902op{&exec{l: 7, f: addHLss}}, nil},
	{"INC ss", 0b11001111, 0b00000011, 1, []lr35902op{&exec{l: 2, f: incSS}}, nil},
	{"DEC ss", 0b11001111, 0b00001011, 1, []lr35902op{&exec{l: 2, f: decSS}}, nil},
	{"POP ss", 0b11001111, 0b11000001, 1, []lr35902op{}, popSS},
	{"PUSH ss", 0b11001111, 0b11000101, 1, []lr35902op{&exec{l: 1, f: pushSS}}, nil},

	{"LD r, n", 0b11000111, 0b00000110, 2, []lr35902op{&mrNpc{f: ldRn}}, nil},
	{"LD r, r'", 0b11000000, 0b01000000, 1, []lr35902op{}, ldRr},
	{"LD r, (HL)", 0b11000111, 0b01000110, 1, []lr35902op{}, ldRhl},
	{"LD (HL), r", 0b11111000, 0b01110000, 1, []lr35902op{}, ldHLr},
	{"INC r", 0b11000111, 0b0000100, 1, []lr35902op{}, incR},
	{"DEC r", 0b11000111, 0b0000101, 1, []lr35902op{}, decR},
	{"ADD A, r", 0b11111000, 0b10000000, 1, []lr35902op{}, addAr},
	{"ADC A, r", 0b11111000, 0b10001000, 1, []lr35902op{}, adcAr},
	{"SUB A, r", 0b11111000, 0b10010000, 1, []lr35902op{}, subAr},
	{"SUC A, r", 0b11111000, 0b10011000, 1, []lr35902op{}, sbcAr},
	{"AND r", 0b11111000, 0b10100000, 1, []lr35902op{}, andAr},
	{"OR r", 0b11111000, 0b10110000, 1, []lr35902op{}, orAr},
	{"XOR r", 0b11111000, 0b10101000, 1, []lr35902op{}, xorAr},
	{"CP r", 0b11111000, 0b10111000, 1, []lr35902op{}, cpR},

	// // TODO: review 0xe0 0xe2 0xe4 0xe8 0xea 0xec 0xf0 0xf2 0xf4 0xf8 0xfa 0xfc
	{"RET cc", 0b11100111, 0b11000000, 1, []lr35902op{&exec{l: 1, f: retCC}}, nil},
	{"JP cc, nn", 0b11100111, 0b11000010, 3, []lr35902op{&mrNNpc{f: jpCC}}, nil},
	{"CALL cc, nn", 0b11100111, 0b11000100, 3, []lr35902op{&mrNNpc{f: callCC}}, nil},

	{"LD (ff00+n), A", 0xFF, 0xe0, 2, []lr35902op{&mrNpc{f: ldhNa}}, nil},
	{"LD (ff00+C), A", 0xFF, 0xe2, 1, []lr35902op{}, ldhCa},

	{"LD A, (ff00+n)", 0xFF, 0xf0, 2, []lr35902op{&mrNpc{f: ldhAn}}, nil},

	{"LD (nn), A", 0xFF, 0xea, 3, []lr35902op{&mrNNpc{f: ldNNa}}, nil},
	{"LD A, (nn)", 0xFF, 0xfa, 3, []lr35902op{&mrNNpc{f: ldAnn}}, nil},

	{"LD (nn), SP", 0xFF, 0x08, 3, []lr35902op{&mrNNpc{f: ldNNsp}}, nil},

	{"LDI (HL),a", 0xFF, 0x22, 1, []lr35902op{}, ldiHLa},
	{"LDI A,(HL)", 0xFF, 0x2a, 1, []lr35902op{}, ldiAhl},

	{"LDD A,(HL)", 0xFF, 0x32, 1, []lr35902op{}, lddAhl},
	{"LDD (HL),a", 0xFF, 0x3a, 1, []lr35902op{}, lddHLa},

	{"LD HL,(SP+e)", 0xFF, 0xF8, 2, []lr35902op{&mrNpc{f: ldHLspE}}, nil},

	{"RST p", 0b11000111, 0b11000111, 1, []lr35902op{&exec{l: 1, f: rstP}}, nil},
	{"CALL nn", 0xFF, 0xCD, 3, []lr35902op{&mrNNpc{}, &exec{l: 1, f: call}}, nil},

	// // {"", 0xFF, 0x1,,[]lr35902op{},nil},
	{"NOP", 0xFF, 0x00, 1, []lr35902op{}, nil},
	{"DAA", 0xFF, 0x27, 1, []lr35902op{}, daa},
	{"CPL", 0xFF, 0x2f, 1, []lr35902op{}, cpl},
	{"SCF", 0xFF, 0x37, 1, []lr35902op{}, scf},
	{"CCF", 0xFF, 0x3F, 1, []lr35902op{}, ccf},
	{"HALT", 0xFF, 0x76, 1, []lr35902op{}, halt},
	{"RET", 0xFF, 0xC9, 1, []lr35902op{}, ret},
	{"RETI", 0xFF, 0xD9, 1, []lr35902op{}, reti},

	{"INC (HL)", 0xFF, 0x34, 1, []lr35902op{&exec{l: 1, f: incHL}}, nil},
	{"DEC (HL)", 0xFF, 0x35, 1, []lr35902op{&exec{l: 1, f: decHL}}, nil},
	{"ADD A, (HL)", 0xFF, 0x86, 1, []lr35902op{}, addAhl},
	{"ADC A, (HL)", 0xFF, 0x8e, 1, []lr35902op{}, adcAhl},
	{"SUB A, (HL)", 0xFF, 0x96, 1, []lr35902op{}, subAhl},
	{"SBC A, (HL)", 0xFF, 0x9e, 1, []lr35902op{}, sbcAhl},
	{"AND (HL)", 0xFF, 0xA6, 1, []lr35902op{}, andAhl},
	{"OR (HL)", 0xFF, 0xB6, 1, []lr35902op{}, orAhl},
	{"XOR (HL)", 0xFF, 0xAE, 1, []lr35902op{}, xorAhl},
	{"CP (HL)", 0xFF, 0xBE, 1, []lr35902op{}, cpHl},
	{"ADD A, n", 0xFF, 0xc6, 2, []lr35902op{&mrNpc{f: func(cpu *lr35902) { cpu.addA(cpu.fetched.n) }}}, nil},
	{"ADC A, n", 0xFF, 0xCE, 2, []lr35902op{&mrNpc{f: func(cpu *lr35902) { cpu.adcA(cpu.fetched.n) }}}, nil},
	{"SBC A, (HL)", 0xFF, 0xDE, 2, []lr35902op{&mrNpc{f: func(cpu *lr35902) { cpu.sbcA(cpu.fetched.n) }}}, nil},
	{"SUB n", 0xFF, 0xD6, 2, []lr35902op{&mrNpc{f: func(cpu *lr35902) { cpu.subA(cpu.fetched.n) }}}, nil},
	{"AND n", 0xFF, 0xE6, 2, []lr35902op{&mrNpc{f: func(cpu *lr35902) { cpu.and(cpu.fetched.n) }}}, nil},
	{"OR n", 0xFF, 0xF6, 2, []lr35902op{&mrNpc{f: func(cpu *lr35902) { cpu.or(cpu.fetched.n) }}}, nil},

	{"LD A,(BC)", 0xFF, 0x0A, 1, []lr35902op{}, ldAbc},
	{"LD A,(DE)", 0xFF, 0x1A, 1, []lr35902op{}, ldAde},
	{"LD (BC), A", 0xFF, 0x02, 1, []lr35902op{}, ldBCa},
	{"LD (DE), A", 0xFF, 0x12, 1, []lr35902op{}, ldDEa},
	{"LD (HL), n", 0xFF, 0x36, 2, []lr35902op{&mrNpc{f: ldHLn}}, nil},

	{"JR e", 0xFF, 0x18, 2, []lr35902op{&mrNpc{}, &exec{l: 5, f: jr}}, nil},
	{"JRNZ e", 0xFF, 0x20, 2, []lr35902op{&mrNpc{f: jrnz}}, nil},
	{"JR Z, e", 0xFF, 0x28, 2, []lr35902op{&mrNpc{f: jrz}}, nil},
	{"JRNC e", 0xFF, 0x30, 2, []lr35902op{&mrNpc{f: jrnc}}, nil},
	{"JRC e", 0xFF, 0x38, 2, []lr35902op{&mrNpc{f: jrc}}, nil},

	{"JP nn", 0xFF, 0xC3, 3, []lr35902op{&mrNNpc{f: func(cpu *lr35902) { cpu.regs.PC = cpu.fetched.nn }}}, nil},

	{"RLCA", 0xFF, 0x07, 1, []lr35902op{}, rlca},
	{"RLA", 0xFF, 0x17, 1, []lr35902op{}, rla},
	{"RRCA", 0xFF, 0x0F, 1, []lr35902op{}, rrca},
	{"RRA", 0xFF, 0x1F, 1, []lr35902op{}, rra},

	{"JP HL", 0xFF, 0xE9, 1, []lr35902op{}, func(cpu *lr35902) { cpu.regs.PC = cpu.regs.HL.Get() }},

	{"XOR *", 0xFF, 0xEE, 2, []lr35902op{&mrNpc{f: func(cpu *lr35902) { cpu.xor(cpu.fetched.n) }}}, nil},
	{"DI", 0xFF, 0xF3, 1, []lr35902op{}, func(cpu *lr35902) { cpu.regs.IME = false }},
	{"EI", 0xFF, 0xFb, 1, []lr35902op{}, func(cpu *lr35902) { cpu.regs.IME = true }},
	{"LD SP, HL", 0xFF, 0xF9, 1, []lr35902op{&exec{l: 2, f: func(cpu *lr35902) { cpu.regs.SP.Set(cpu.regs.HL.Get()) }}}, nil},
	{"CP n", 0xFF, 0xFe, 2, []lr35902op{&mrNpc{f: func(cpu *lr35902) { cpu.cp(cpu.fetched.n) }}}, nil},

	{"CB", 0xFF, 0xCB, 1, []lr35902op{}, decodeCB},
}

var lr35902OpsCodeTableCB = []*opCode{
	{"RLC r", 0b11111000, 0b00000000, 1, []lr35902op{}, cbR},
	{"RLC (HL)", 0xFF, 0x06, 1, []lr35902op{&exec{l: 1, f: cbHL}}, nil},
	{"RRC r", 0b11111000, 0b00001000, 1, []lr35902op{}, cbR},
	{"RRC (HL)", 0xFF, 0x0e, 1, []lr35902op{&exec{l: 1, f: cbHL}}, nil},

	{"RL r", 0b11111000, 0b00010000, 1, []lr35902op{}, cbR},
	{"RL (HL)", 0xFF, 0x16, 1, []lr35902op{&exec{l: 1, f: cbHL}}, nil},
	{"RR r", 0b11111000, 0b00011000, 1, []lr35902op{}, cbR},
	{"RR (HL)", 0xFF, 0x1e, 1, []lr35902op{&exec{l: 1, f: cbHL}}, nil},

	{"SLA r", 0b11111000, 0b00100000, 1, []lr35902op{}, cbR},
	{"SLA (HL)", 0xFF, 0x26, 1, []lr35902op{&exec{l: 1, f: cbHL}}, nil},
	{"SRA r", 0b11111000, 0b00101000, 1, []lr35902op{}, cbR},
	{"SRA (HL)", 0xFF, 0x2e, 1, []lr35902op{&exec{l: 1, f: cbHL}}, nil},

	{"SRL r", 0b11111000, 0b00111000, 1, []lr35902op{}, cbR},
	{"SRL (HL)", 0xFF, 0x3e, 1, []lr35902op{&exec{l: 1, f: cbHL}}, nil},

	{"BIT b, r", 0b11000000, 0b01000000, 1, []lr35902op{}, bit},
	{"BIT b, (HL)", 0b11000111, 0b01000110, 1, []lr35902op{&exec{l: 1, f: bitHL}}, nil},

	{"RES b, r", 0b11000000, 0b10000000, 1, []lr35902op{}, res},
	{"RES b, (HL)", 0b11000111, 0b10000110, 1, []lr35902op{&exec{l: 1, f: resHL}}, nil},

	{"SET b, r", 0b11000000, 0b11000000, 1, []lr35902op{}, set},
	{"SET b, (HL)", 0b11000111, 0b11000110, 1, []lr35902op{&exec{l: 1, f: setHL}}, nil},

	{"SWAP r", 0b1111_0000, 0b0011_0000, 1, []lr35902op{}, swap},
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
