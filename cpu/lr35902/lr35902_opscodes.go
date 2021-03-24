package lr35902

type opCode struct {
	ins string
	len byte
	f   lr35902f
}

func (cpu *lr35902) initOpCodes() {
	cpu.opCodes = make([]*opCode, 256)
	cpu.opCodesCB = make([]*opCode, 256)

	cpu.opCodes[0x00] = &opCode{"NOP", 1, func(cpu *lr35902) {}}
	cpu.opCodes[0x01] = &opCode{"LD BC, nn", 3, ldDDnn}
	cpu.opCodes[0x02] = &opCode{"LD (BC), A", 1, ldBCa}
	cpu.opCodes[0x03] = &opCode{"INC BC", 1, incSS}
	cpu.opCodes[0x04] = &opCode{"INC B", 1, incR}
	cpu.opCodes[0x05] = &opCode{"DEC B", 1, decR}
	cpu.opCodes[0x06] = &opCode{"LD B, n", 2, ldRn}
	cpu.opCodes[0x07] = &opCode{"RLCA", 1, rlca}
	cpu.opCodes[0x08] = &opCode{"LD (nn), SP", 3, ldNNsp}
	cpu.opCodes[0x09] = &opCode{"ADD HL,BC", 1, addHLss}
	cpu.opCodes[0x0A] = &opCode{"LD A,(BC)", 1, ldAbc}
	cpu.opCodes[0x0B] = &opCode{"DEC BC", 1, decSS}
	cpu.opCodes[0x0C] = &opCode{"INC C", 1, incR}
	cpu.opCodes[0x0D] = &opCode{"DEC C", 1, decR}
	cpu.opCodes[0x0E] = &opCode{"LD C, n", 2, ldRn}
	cpu.opCodes[0x0F] = &opCode{"RRCA", 1, rrca}
	cpu.opCodes[0x0F] = &opCode{"STOP", 1, func(cpu *lr35902) { panic(-1) }}
	cpu.opCodes[0x11] = &opCode{"LD DE, nn", 3, ldDDnn}
	cpu.opCodes[0x12] = &opCode{"LD (DE), A", 1, ldDEa}
	cpu.opCodes[0x13] = &opCode{"INC DE", 1, incSS}
	cpu.opCodes[0x14] = &opCode{"INC D", 1, incR}
	cpu.opCodes[0x15] = &opCode{"DEC D", 1, decR}
	cpu.opCodes[0x16] = &opCode{"LD D, n", 2, ldRn}
	cpu.opCodes[0x17] = &opCode{"RLA", 1, rla}
	cpu.opCodes[0x18] = &opCode{"JR e", 2, jr}
	cpu.opCodes[0x19] = &opCode{"ADD HL,DE", 1, addHLss}
	cpu.opCodes[0x1A] = &opCode{"LD A,(DE)", 1, ldAde}
	cpu.opCodes[0x1B] = &opCode{"DEC DE", 1, decSS}
	cpu.opCodes[0x1C] = &opCode{"INC E", 1, incR}
	cpu.opCodes[0x1D] = &opCode{"DEC E", 1, decR}
	cpu.opCodes[0x1E] = &opCode{"LD E, n", 2, ldRn}
	cpu.opCodes[0x1F] = &opCode{"RRA", 1, rra}
	cpu.opCodes[0x20] = &opCode{"JR NZ, e", 2, jrnz}
	cpu.opCodes[0x21] = &opCode{"LD HL, nn", 3, ldDDnn}
	cpu.opCodes[0x22] = &opCode{"LDI (HL),a", 1, ldiHLa}
	cpu.opCodes[0x23] = &opCode{"INC HL", 1, incSS}
	cpu.opCodes[0x24] = &opCode{"INC H", 1, incR}
	cpu.opCodes[0x25] = &opCode{"DEC H", 1, decR}
	cpu.opCodes[0x26] = &opCode{"LD H, n", 2, ldRn}
	cpu.opCodes[0x27] = &opCode{"DAA", 1, daa}
	cpu.opCodes[0x28] = &opCode{"JR Z, e", 2, jrz}
	cpu.opCodes[0x29] = &opCode{"ADD HL,HL", 1, addHLss}
	cpu.opCodes[0x2A] = &opCode{"LDI A,(HL)", 1, ldiAhl}
	cpu.opCodes[0x2B] = &opCode{"DEC HL", 1, decSS}
	cpu.opCodes[0x2C] = &opCode{"INC L", 1, incR}
	cpu.opCodes[0x2D] = &opCode{"DEC L", 1, decR}
	cpu.opCodes[0x2E] = &opCode{"LD L, n", 2, ldRn}
	cpu.opCodes[0x2F] = &opCode{"CPL", 1, cpl}
	cpu.opCodes[0x30] = &opCode{"JR NC, e", 2, jrnc}
	cpu.opCodes[0x31] = &opCode{"LD SP, nn", 3, ldDDnn}
	cpu.opCodes[0x32] = &opCode{"LDD (HL),a", 1, lddHLa}
	cpu.opCodes[0x33] = &opCode{"INC SP", 1, incSS}
	cpu.opCodes[0x34] = &opCode{"INC (HL)", 1, incHL}
	cpu.opCodes[0x35] = &opCode{"DEC (HL)", 1, decHL}
	cpu.opCodes[0x36] = &opCode{"LD (HL), n", 2, ldToHLn}
	cpu.opCodes[0x37] = &opCode{"SCF", 1, scf}
	cpu.opCodes[0x38] = &opCode{"JR C, e", 2, jrc}
	cpu.opCodes[0x39] = &opCode{"ADD HL,SP", 1, addHLss}
	cpu.opCodes[0x3A] = &opCode{"LDD A,(HL)", 1, lddAhl}
	cpu.opCodes[0x3B] = &opCode{"DEC SP", 1, decSS}
	cpu.opCodes[0x3C] = &opCode{"INC A", 1, incR}
	cpu.opCodes[0x3D] = &opCode{"DEC A", 1, decR}
	cpu.opCodes[0x3E] = &opCode{"LD A, n", 2, ldRn}
	cpu.opCodes[0x3F] = &opCode{"CCF", 1, ccf}
	cpu.opCodes[0x40] = &opCode{"LD B, B", 1, func(cpu *lr35902) {}}
	cpu.opCodes[0x41] = &opCode{"LD B, C", 1, func(cpu *lr35902) { cpu.regs.B = cpu.regs.C }}
	cpu.opCodes[0x42] = &opCode{"LD B, D", 1, func(cpu *lr35902) { cpu.regs.B = cpu.regs.D }}
	cpu.opCodes[0x43] = &opCode{"LD B, E", 1, func(cpu *lr35902) { cpu.regs.B = cpu.regs.E }}
	cpu.opCodes[0x44] = &opCode{"LD B, H", 1, func(cpu *lr35902) { cpu.regs.B = cpu.regs.H }}
	cpu.opCodes[0x45] = &opCode{"LD B, L", 1, func(cpu *lr35902) { cpu.regs.B = cpu.regs.L }}
	cpu.opCodes[0x46] = &opCode{"LD B, (HL)", 1, ldFromHL}
	cpu.opCodes[0x47] = &opCode{"LD B, A", 1, func(cpu *lr35902) { cpu.regs.B = cpu.regs.A }}
	cpu.opCodes[0x48] = &opCode{"LD C, B", 1, func(cpu *lr35902) { cpu.regs.C = cpu.regs.B }}
	cpu.opCodes[0x49] = &opCode{"LD C, C", 1, func(cpu *lr35902) {}}
	cpu.opCodes[0x4A] = &opCode{"LD C, D", 1, func(cpu *lr35902) { cpu.regs.C = cpu.regs.D }}
	cpu.opCodes[0x4B] = &opCode{"LD C, E", 1, func(cpu *lr35902) { cpu.regs.C = cpu.regs.E }}
	cpu.opCodes[0x4C] = &opCode{"LD C, H", 1, func(cpu *lr35902) { cpu.regs.C = cpu.regs.H }}
	cpu.opCodes[0x4D] = &opCode{"LD C, L", 1, func(cpu *lr35902) { cpu.regs.C = cpu.regs.L }}
	cpu.opCodes[0x4E] = &opCode{"LD C, (HL)", 1, ldFromHL}
	cpu.opCodes[0x4F] = &opCode{"LD C, A", 1, func(cpu *lr35902) { cpu.regs.C = cpu.regs.A }}
	cpu.opCodes[0x50] = &opCode{"LD D, B", 1, func(cpu *lr35902) { cpu.regs.D = cpu.regs.B }}
	cpu.opCodes[0x51] = &opCode{"LD D, C", 1, func(cpu *lr35902) { cpu.regs.D = cpu.regs.C }}
	cpu.opCodes[0x52] = &opCode{"LD D, D", 1, func(cpu *lr35902) {}}
	cpu.opCodes[0x53] = &opCode{"LD D, E", 1, func(cpu *lr35902) { cpu.regs.D = cpu.regs.E }}
	cpu.opCodes[0x54] = &opCode{"LD D, H", 1, func(cpu *lr35902) { cpu.regs.D = cpu.regs.H }}
	cpu.opCodes[0x55] = &opCode{"LD D, L", 1, func(cpu *lr35902) { cpu.regs.D = cpu.regs.L }}
	cpu.opCodes[0x56] = &opCode{"LD D, (HL)", 1, ldFromHL}
	cpu.opCodes[0x57] = &opCode{"LD D, A", 1, func(cpu *lr35902) { cpu.regs.D = cpu.regs.A }}
	cpu.opCodes[0x58] = &opCode{"LD E, B", 1, func(cpu *lr35902) { cpu.regs.E = cpu.regs.B }}
	cpu.opCodes[0x59] = &opCode{"LD E, C", 1, func(cpu *lr35902) { cpu.regs.E = cpu.regs.C }}
	cpu.opCodes[0x5A] = &opCode{"LD E, D", 1, func(cpu *lr35902) { cpu.regs.E = cpu.regs.D }}
	cpu.opCodes[0x5B] = &opCode{"LD E, E", 1, func(cpu *lr35902) {}}
	cpu.opCodes[0x5C] = &opCode{"LD E, H", 1, func(cpu *lr35902) { cpu.regs.E = cpu.regs.H }}
	cpu.opCodes[0x5D] = &opCode{"LD E, L", 1, func(cpu *lr35902) { cpu.regs.E = cpu.regs.L }}
	cpu.opCodes[0x5E] = &opCode{"LD E, (HL)", 1, ldFromHL}
	cpu.opCodes[0x5F] = &opCode{"LD E, A", 1, func(cpu *lr35902) { cpu.regs.E = cpu.regs.A }}
	cpu.opCodes[0x60] = &opCode{"LD H, B", 1, func(cpu *lr35902) { cpu.regs.H = cpu.regs.B }}
	cpu.opCodes[0x61] = &opCode{"LD H, C", 1, func(cpu *lr35902) { cpu.regs.H = cpu.regs.C }}
	cpu.opCodes[0x62] = &opCode{"LD H, D", 1, func(cpu *lr35902) { cpu.regs.H = cpu.regs.D }}
	cpu.opCodes[0x63] = &opCode{"LD H, E", 1, func(cpu *lr35902) { cpu.regs.H = cpu.regs.E }}
	cpu.opCodes[0x64] = &opCode{"LD H, H", 1, func(cpu *lr35902) {}}
	cpu.opCodes[0x65] = &opCode{"LD H, L", 1, func(cpu *lr35902) { cpu.regs.H = cpu.regs.L }}
	cpu.opCodes[0x66] = &opCode{"LD H, (HL)", 1, ldFromHL}
	cpu.opCodes[0x67] = &opCode{"LD H, A", 1, func(cpu *lr35902) { cpu.regs.H = cpu.regs.A }}
	cpu.opCodes[0x68] = &opCode{"LD L, B", 1, func(cpu *lr35902) { cpu.regs.L = cpu.regs.B }}
	cpu.opCodes[0x69] = &opCode{"LD L, C", 1, func(cpu *lr35902) { cpu.regs.L = cpu.regs.C }}
	cpu.opCodes[0x6A] = &opCode{"LD L, D", 1, func(cpu *lr35902) { cpu.regs.L = cpu.regs.D }}
	cpu.opCodes[0x6B] = &opCode{"LD L, E", 1, func(cpu *lr35902) { cpu.regs.L = cpu.regs.E }}
	cpu.opCodes[0x6C] = &opCode{"LD L, H", 1, func(cpu *lr35902) { cpu.regs.L = cpu.regs.H }}
	cpu.opCodes[0x6D] = &opCode{"LD L, L", 1, func(cpu *lr35902) {}}
	cpu.opCodes[0x6E] = &opCode{"LD L, (HL)", 1, ldFromHL}
	cpu.opCodes[0x6F] = &opCode{"LD L, A", 1, func(cpu *lr35902) { cpu.regs.L = cpu.regs.A }}
	cpu.opCodes[0x70] = &opCode{"LD (HL), B", 1, ldToHL}
	cpu.opCodes[0x71] = &opCode{"LD (HL), C", 1, ldToHL}
	cpu.opCodes[0x72] = &opCode{"LD (HL), D", 1, ldToHL}
	cpu.opCodes[0x73] = &opCode{"LD (HL), E", 1, ldToHL}
	cpu.opCodes[0x74] = &opCode{"LD (HL), H", 1, ldToHL}
	cpu.opCodes[0x75] = &opCode{"LD (HL), L", 1, ldToHL}
	cpu.opCodes[0x76] = &opCode{"HALT", 1, halt}
	cpu.opCodes[0x77] = &opCode{"LD (HL), A", 1, ldToHL}
	cpu.opCodes[0x78] = &opCode{"LD A, B", 1, func(cpu *lr35902) { cpu.regs.A = cpu.regs.B }}
	cpu.opCodes[0x79] = &opCode{"LD A, C", 1, func(cpu *lr35902) { cpu.regs.A = cpu.regs.C }}
	cpu.opCodes[0x7A] = &opCode{"LD A, D", 1, func(cpu *lr35902) { cpu.regs.A = cpu.regs.D }}
	cpu.opCodes[0x7B] = &opCode{"LD A, E", 1, func(cpu *lr35902) { cpu.regs.A = cpu.regs.E }}
	cpu.opCodes[0x7C] = &opCode{"LD A, H", 1, func(cpu *lr35902) { cpu.regs.A = cpu.regs.H }}
	cpu.opCodes[0x7D] = &opCode{"LD A, L", 1, func(cpu *lr35902) { cpu.regs.A = cpu.regs.L }}
	cpu.opCodes[0x7E] = &opCode{"LD A, (HL)", 1, ldFromHL}
	cpu.opCodes[0x7F] = &opCode{"LD A, A", 1, func(cpu *lr35902) {}}
	cpu.opCodes[0x80] = &opCode{"ADD A, B", 1, addAr}
	cpu.opCodes[0x81] = &opCode{"ADD A, C", 1, addAr}
	cpu.opCodes[0x82] = &opCode{"ADD A, D", 1, addAr}
	cpu.opCodes[0x83] = &opCode{"ADD A, E", 1, addAr}
	cpu.opCodes[0x84] = &opCode{"ADD A, H", 1, addAr}
	cpu.opCodes[0x85] = &opCode{"ADD A, L", 1, addAr}
	cpu.opCodes[0x86] = &opCode{"ADD A, (HL)", 1, addAhl}
	cpu.opCodes[0x87] = &opCode{"ADD A, A", 1, addAr}
	cpu.opCodes[0x88] = &opCode{"ADC A, B", 1, adcAr}
	cpu.opCodes[0x89] = &opCode{"ADC A, C", 1, adcAr}
	cpu.opCodes[0x8A] = &opCode{"ADC A, D", 1, adcAr}
	cpu.opCodes[0x8B] = &opCode{"ADC A, E", 1, adcAr}
	cpu.opCodes[0x8C] = &opCode{"ADC A, H", 1, adcAr}
	cpu.opCodes[0x8D] = &opCode{"ADC A, L", 1, adcAr}
	cpu.opCodes[0x8E] = &opCode{"ADC A, (HL)", 1, adcAhl}
	cpu.opCodes[0x8F] = &opCode{"ADC A, A", 1, adcAr}
	cpu.opCodes[0x90] = &opCode{"SUB A, B", 1, subAr}
	cpu.opCodes[0x91] = &opCode{"SUB A, C", 1, subAr}
	cpu.opCodes[0x92] = &opCode{"SUB A, D", 1, subAr}
	cpu.opCodes[0x93] = &opCode{"SUB A, E", 1, subAr}
	cpu.opCodes[0x94] = &opCode{"SUB A, H", 1, subAr}
	cpu.opCodes[0x95] = &opCode{"SUB A, L", 1, subAr}
	cpu.opCodes[0x96] = &opCode{"SUB A, (HL)", 1, subAhl}
	cpu.opCodes[0x97] = &opCode{"SUB A, A", 1, subAr}
	cpu.opCodes[0x98] = &opCode{"SUC A, B", 1, sbcAr}
	cpu.opCodes[0x99] = &opCode{"SUC A, C", 1, sbcAr}
	cpu.opCodes[0x9A] = &opCode{"SUC A, D", 1, sbcAr}
	cpu.opCodes[0x9B] = &opCode{"SUC A, E", 1, sbcAr}
	cpu.opCodes[0x9C] = &opCode{"SUC A, H", 1, sbcAr}
	cpu.opCodes[0x9D] = &opCode{"SUC A, L", 1, sbcAr}
	cpu.opCodes[0x9E] = &opCode{"SBC A, (HL)", 1, sbcAhl}
	cpu.opCodes[0x9F] = &opCode{"SUC A, A", 1, sbcAr}
	cpu.opCodes[0xA0] = &opCode{"AND B", 1, andAr}
	cpu.opCodes[0xA1] = &opCode{"AND C", 1, andAr}
	cpu.opCodes[0xA2] = &opCode{"AND D", 1, andAr}
	cpu.opCodes[0xA3] = &opCode{"AND E", 1, andAr}
	cpu.opCodes[0xA4] = &opCode{"AND H", 1, andAr}
	cpu.opCodes[0xA5] = &opCode{"AND L", 1, andAr}
	cpu.opCodes[0xA6] = &opCode{"AND (HL)", 1, andAhl}
	cpu.opCodes[0xA7] = &opCode{"AND A", 1, andAr}
	cpu.opCodes[0xA8] = &opCode{"XOR B", 1, xorAr}
	cpu.opCodes[0xA9] = &opCode{"XOR C", 1, xorAr}
	cpu.opCodes[0xAA] = &opCode{"XOR D", 1, xorAr}
	cpu.opCodes[0xAB] = &opCode{"XOR E", 1, xorAr}
	cpu.opCodes[0xAC] = &opCode{"XOR H", 1, xorAr}
	cpu.opCodes[0xAD] = &opCode{"XOR L", 1, xorAr}
	cpu.opCodes[0xAE] = &opCode{"XOR (HL)", 1, xorAhl}
	cpu.opCodes[0xAF] = &opCode{"XOR A", 1, xorAr}
	cpu.opCodes[0xB0] = &opCode{"OR B", 1, orAr}
	cpu.opCodes[0xB1] = &opCode{"OR C", 1, orAr}
	cpu.opCodes[0xB2] = &opCode{"OR D", 1, orAr}
	cpu.opCodes[0xB3] = &opCode{"OR E", 1, orAr}
	cpu.opCodes[0xB4] = &opCode{"OR H", 1, orAr}
	cpu.opCodes[0xB5] = &opCode{"OR L", 1, orAr}
	cpu.opCodes[0xB6] = &opCode{"OR (HL)", 1, orAhl}
	cpu.opCodes[0xB7] = &opCode{"OR A", 1, orAr}
	cpu.opCodes[0xB8] = &opCode{"CP B", 1, cpR}
	cpu.opCodes[0xB9] = &opCode{"CP C", 1, cpR}
	cpu.opCodes[0xBA] = &opCode{"CP D", 1, cpR}
	cpu.opCodes[0xBB] = &opCode{"CP E", 1, cpR}
	cpu.opCodes[0xBC] = &opCode{"CP H", 1, cpR}
	cpu.opCodes[0xBD] = &opCode{"CP L", 1, cpR}
	cpu.opCodes[0xBE] = &opCode{"CP (HL)", 1, cpHl}
	cpu.opCodes[0xBF] = &opCode{"CP A", 1, cpR}
	cpu.opCodes[0xC0] = &opCode{"RET NZ", 1, retCC}
	cpu.opCodes[0xC1] = &opCode{"POP BC", 1, popSS}
	cpu.opCodes[0xC2] = &opCode{"JP NZ, nn", 3, jpCC}
	cpu.opCodes[0xC3] = &opCode{"JP nn", 3, func(cpu *lr35902) { cpu.regs.PC = cpu.fetched.nn }}
	cpu.opCodes[0xC4] = &opCode{"CALL NZ, nn", 3, callCC}
	cpu.opCodes[0xC5] = &opCode{"PUSH BC", 1, pushSS}
	cpu.opCodes[0xC6] = &opCode{"ADD A, n", 2, func(cpu *lr35902) { cpu.addA(cpu.fetched.n) }}
	cpu.opCodes[0xC7] = &opCode{"RST 0x0", 1, rstP}
	cpu.opCodes[0xC8] = &opCode{"RET Z", 1, retCC}
	cpu.opCodes[0xC9] = &opCode{"RET", 1, ret}
	cpu.opCodes[0xCA] = &opCode{"JP Z, nn", 3, jpCC}
	cpu.opCodes[0xCB] = &opCode{"CB", 1, decodeCB}
	cpu.opCodes[0xCC] = &opCode{"CALL Z, nn", 3, callCC}
	cpu.opCodes[0xCD] = &opCode{"CALL nn", 3, call}
	cpu.opCodes[0xCE] = &opCode{"ADC A, n", 2, func(cpu *lr35902) { cpu.adcA(cpu.fetched.n) }}
	cpu.opCodes[0xCF] = &opCode{"RST 0x8", 1, rstP}
	cpu.opCodes[0xD0] = &opCode{"RET NC", 1, retCC}
	cpu.opCodes[0xD1] = &opCode{"POP DE", 1, popSS}
	cpu.opCodes[0xD2] = &opCode{"JP NC, nn", 3, jpCC}
	cpu.opCodes[0xD4] = &opCode{"CALL NC, nn", 3, callCC}
	cpu.opCodes[0xD5] = &opCode{"PUSH DE", 1, pushSS}
	cpu.opCodes[0xD6] = &opCode{"SUB A, n", 2, func(cpu *lr35902) { cpu.subA(cpu.fetched.n) }}
	cpu.opCodes[0xD7] = &opCode{"RST 0x10", 1, rstP}
	cpu.opCodes[0xD8] = &opCode{"RET C", 1, retCC}
	cpu.opCodes[0xD9] = &opCode{"RETI", 1, reti}
	cpu.opCodes[0xDA] = &opCode{"JP C, nn", 3, jpCC}
	cpu.opCodes[0xDC] = &opCode{"CALL C, nn", 3, callCC}
	cpu.opCodes[0xDE] = &opCode{"SBC A, nn", 3, func(cpu *lr35902) { cpu.sbcA(cpu.fetched.n) }}
	cpu.opCodes[0xDF] = &opCode{"RST 0x18", 1, rstP}
	cpu.opCodes[0xE0] = &opCode{"LD (0xff00+n), A", 2, ldhNa}
	cpu.opCodes[0xE1] = &opCode{"POP HL", 1, popSS}
	cpu.opCodes[0xE2] = &opCode{"LD (0xff00+C), A", 1, ldhCa}
	cpu.opCodes[0xE5] = &opCode{"PUSH HL", 1, pushSS}
	cpu.opCodes[0xE6] = &opCode{"AND n", 2, func(cpu *lr35902) { cpu.and(cpu.fetched.n) }}
	cpu.opCodes[0xE7] = &opCode{"RST 0x20", 1, rstP}
	cpu.opCodes[0xE8] = &opCode{"ADD SP,n", 2, addSPn}
	cpu.opCodes[0xE9] = &opCode{"JP HL", 1, func(cpu *lr35902) { cpu.regs.PC = cpu.regs.HL.Get() }}
	cpu.opCodes[0xEA] = &opCode{"LD (nn), A", 3, ldNNa}
	cpu.opCodes[0xEE] = &opCode{"XOR A, n", 2, func(cpu *lr35902) { cpu.xor(cpu.fetched.n) }}
	cpu.opCodes[0xEF] = &opCode{"RST 0x28", 1, rstP}
	cpu.opCodes[0xF0] = &opCode{"LD A, (0xff00+n)", 2, ldhAn}
	cpu.opCodes[0xF1] = &opCode{"POP AF", 1, popSS}
	cpu.opCodes[0xF2] = &opCode{"LD A, (0xff00+C)", 1, ldhAc}
	cpu.opCodes[0xF3] = &opCode{"DI", 1, func(cpu *lr35902) { cpu.regs.IME = false }}
	cpu.opCodes[0xF5] = &opCode{"PUSH AF", 1, pushSS}
	cpu.opCodes[0xF6] = &opCode{"OR n", 2, func(cpu *lr35902) { cpu.or(cpu.fetched.n) }}
	cpu.opCodes[0xF7] = &opCode{"RST 0x30", 1, rstP}
	cpu.opCodes[0xF8] = &opCode{"LD HL,(SP+n)", 2, ldHLspE}
	cpu.opCodes[0xF9] = &opCode{"LD SP, HL", 1, func(cpu *lr35902) { cpu.regs.SP.Set(cpu.regs.HL.Get()) }}
	cpu.opCodes[0xFA] = &opCode{"LD A, (nn)", 3, ldAnn}
	cpu.opCodes[0xFB] = &opCode{"EI", 1, func(cpu *lr35902) { cpu.regs.IME = true }}
	cpu.opCodes[0xFE] = &opCode{"CP A, n", 2, func(cpu *lr35902) { cpu.cp(cpu.fetched.n) }}
	cpu.opCodes[0xFF] = &opCode{"RST 0x38", 1, rstP}

	cpu.opCodesCB[0x00] = &opCode{"RLC B", 1, cbR}
	cpu.opCodesCB[0x01] = &opCode{"RLC C", 1, cbR}
	cpu.opCodesCB[0x02] = &opCode{"RLC D", 1, cbR}
	cpu.opCodesCB[0x03] = &opCode{"RLC E", 1, cbR}
	cpu.opCodesCB[0x04] = &opCode{"RLC H", 1, cbR}
	cpu.opCodesCB[0x05] = &opCode{"RLC L", 1, cbR}
	cpu.opCodesCB[0x06] = &opCode{"RLC (HL)", 1, cbHL}
	cpu.opCodesCB[0x07] = &opCode{"RLC A", 1, cbR}
	cpu.opCodesCB[0x08] = &opCode{"RRC B", 1, cbR}
	cpu.opCodesCB[0x09] = &opCode{"RRC C", 1, cbR}
	cpu.opCodesCB[0x0A] = &opCode{"RRC D", 1, cbR}
	cpu.opCodesCB[0x0B] = &opCode{"RRC E", 1, cbR}
	cpu.opCodesCB[0x0C] = &opCode{"RRC H", 1, cbR}
	cpu.opCodesCB[0x0D] = &opCode{"RRC L", 1, cbR}
	cpu.opCodesCB[0x0E] = &opCode{"RRC (HL)", 1, cbHL}
	cpu.opCodesCB[0x0F] = &opCode{"RRC A", 1, cbR}
	cpu.opCodesCB[0x10] = &opCode{"RL B", 1, cbR}
	cpu.opCodesCB[0x11] = &opCode{"RL C", 1, cbR}
	cpu.opCodesCB[0x12] = &opCode{"RL D", 1, cbR}
	cpu.opCodesCB[0x13] = &opCode{"RL E", 1, cbR}
	cpu.opCodesCB[0x14] = &opCode{"RL H", 1, cbR}
	cpu.opCodesCB[0x15] = &opCode{"RL L", 1, cbR}
	cpu.opCodesCB[0x16] = &opCode{"RL (HL)", 1, cbHL}
	cpu.opCodesCB[0x17] = &opCode{"RL A", 1, cbR}
	cpu.opCodesCB[0x18] = &opCode{"RR B", 1, cbR}
	cpu.opCodesCB[0x19] = &opCode{"RR C", 1, cbR}
	cpu.opCodesCB[0x1A] = &opCode{"RR D", 1, cbR}
	cpu.opCodesCB[0x1B] = &opCode{"RR E", 1, cbR}
	cpu.opCodesCB[0x1C] = &opCode{"RR H", 1, cbR}
	cpu.opCodesCB[0x1D] = &opCode{"RR L", 1, cbR}
	cpu.opCodesCB[0x1E] = &opCode{"RR (HL)", 1, cbHL}
	cpu.opCodesCB[0x1F] = &opCode{"RR A", 1, cbR}
	cpu.opCodesCB[0x20] = &opCode{"SLA B", 1, cbR}
	cpu.opCodesCB[0x21] = &opCode{"SLA C", 1, cbR}
	cpu.opCodesCB[0x22] = &opCode{"SLA D", 1, cbR}
	cpu.opCodesCB[0x23] = &opCode{"SLA E", 1, cbR}
	cpu.opCodesCB[0x24] = &opCode{"SLA H", 1, cbR}
	cpu.opCodesCB[0x25] = &opCode{"SLA L", 1, cbR}
	cpu.opCodesCB[0x26] = &opCode{"SLA (HL)", 1, cbHL}
	cpu.opCodesCB[0x27] = &opCode{"SLA A", 1, cbR}
	cpu.opCodesCB[0x28] = &opCode{"SRA B", 1, cbR}
	cpu.opCodesCB[0x29] = &opCode{"SRA C", 1, cbR}
	cpu.opCodesCB[0x2A] = &opCode{"SRA D", 1, cbR}
	cpu.opCodesCB[0x2B] = &opCode{"SRA E", 1, cbR}
	cpu.opCodesCB[0x2C] = &opCode{"SRA H", 1, cbR}
	cpu.opCodesCB[0x2D] = &opCode{"SRA L", 1, cbR}
	cpu.opCodesCB[0x2E] = &opCode{"SRA (HL)", 1, cbHL}
	cpu.opCodesCB[0x2F] = &opCode{"SRA A", 1, cbR}
	cpu.opCodesCB[0x30] = &opCode{"SWAP B", 1, swap}
	cpu.opCodesCB[0x31] = &opCode{"SWAP C", 1, swap}
	cpu.opCodesCB[0x32] = &opCode{"SWAP D", 1, swap}
	cpu.opCodesCB[0x33] = &opCode{"SWAP E", 1, swap}
	cpu.opCodesCB[0x34] = &opCode{"SWAP H", 1, swap}
	cpu.opCodesCB[0x35] = &opCode{"SWAP L", 1, swap}
	cpu.opCodesCB[0x36] = &opCode{"SWAP (HL)", 1, swapHL}
	cpu.opCodesCB[0x37] = &opCode{"SWAP A", 1, swap}
	cpu.opCodesCB[0x38] = &opCode{"SRL B", 1, cbR}
	cpu.opCodesCB[0x39] = &opCode{"SRL C", 1, cbR}
	cpu.opCodesCB[0x3A] = &opCode{"SRL D", 1, cbR}
	cpu.opCodesCB[0x3B] = &opCode{"SRL E", 1, cbR}
	cpu.opCodesCB[0x3C] = &opCode{"SRL H", 1, cbR}
	cpu.opCodesCB[0x3D] = &opCode{"SRL L", 1, cbR}
	cpu.opCodesCB[0x3E] = &opCode{"SRL (HL)", 1, cbHL}
	cpu.opCodesCB[0x3F] = &opCode{"SRL A", 1, cbR}
	cpu.opCodesCB[0x40] = &opCode{"BIT 0 B", 1, bit}
	cpu.opCodesCB[0x41] = &opCode{"BIT 0 C", 1, bit}
	cpu.opCodesCB[0x42] = &opCode{"BIT 0 D", 1, bit}
	cpu.opCodesCB[0x43] = &opCode{"BIT 0 E", 1, bit}
	cpu.opCodesCB[0x44] = &opCode{"BIT 0 H", 1, bit}
	cpu.opCodesCB[0x45] = &opCode{"BIT 0 L", 1, bit}
	cpu.opCodesCB[0x46] = &opCode{"BIT 0 (HL)", 1, bitHL}
	cpu.opCodesCB[0x47] = &opCode{"BIT 0 A", 1, bit}
	cpu.opCodesCB[0x48] = &opCode{"BIT 1 B", 1, bit}
	cpu.opCodesCB[0x49] = &opCode{"BIT 1 C", 1, bit}
	cpu.opCodesCB[0x4A] = &opCode{"BIT 1 D", 1, bit}
	cpu.opCodesCB[0x4B] = &opCode{"BIT 1 E", 1, bit}
	cpu.opCodesCB[0x4C] = &opCode{"BIT 1 H", 1, bit}
	cpu.opCodesCB[0x4D] = &opCode{"BIT 1 L", 1, bit}
	cpu.opCodesCB[0x4E] = &opCode{"BIT 1 (HL)", 1, bitHL}
	cpu.opCodesCB[0x4F] = &opCode{"BIT 1 A", 1, bit}
	cpu.opCodesCB[0x50] = &opCode{"BIT 2 B", 1, bit}
	cpu.opCodesCB[0x51] = &opCode{"BIT 2 C", 1, bit}
	cpu.opCodesCB[0x52] = &opCode{"BIT 2 D", 1, bit}
	cpu.opCodesCB[0x53] = &opCode{"BIT 2 E", 1, bit}
	cpu.opCodesCB[0x54] = &opCode{"BIT 2 H", 1, bit}
	cpu.opCodesCB[0x55] = &opCode{"BIT 2 L", 1, bit}
	cpu.opCodesCB[0x56] = &opCode{"BIT 2 (HL)", 1, bitHL}
	cpu.opCodesCB[0x57] = &opCode{"BIT 2 A", 1, bit}
	cpu.opCodesCB[0x58] = &opCode{"BIT 3 B", 1, bit}
	cpu.opCodesCB[0x59] = &opCode{"BIT 3 C", 1, bit}
	cpu.opCodesCB[0x5A] = &opCode{"BIT 3 D", 1, bit}
	cpu.opCodesCB[0x5B] = &opCode{"BIT 3 E", 1, bit}
	cpu.opCodesCB[0x5C] = &opCode{"BIT 3 H", 1, bit}
	cpu.opCodesCB[0x5D] = &opCode{"BIT 3 L", 1, bit}
	cpu.opCodesCB[0x5E] = &opCode{"BIT 3 (HL)", 1, bitHL}
	cpu.opCodesCB[0x5F] = &opCode{"BIT 3 A", 1, bit}
	cpu.opCodesCB[0x60] = &opCode{"BIT 4 B", 1, bit}
	cpu.opCodesCB[0x61] = &opCode{"BIT 4 C", 1, bit}
	cpu.opCodesCB[0x62] = &opCode{"BIT 4 D", 1, bit}
	cpu.opCodesCB[0x63] = &opCode{"BIT 4 E", 1, bit}
	cpu.opCodesCB[0x64] = &opCode{"BIT 4 H", 1, bit}
	cpu.opCodesCB[0x65] = &opCode{"BIT 4 L", 1, bit}
	cpu.opCodesCB[0x66] = &opCode{"BIT 4 (HL)", 1, bitHL}
	cpu.opCodesCB[0x67] = &opCode{"BIT 4 A", 1, bit}
	cpu.opCodesCB[0x68] = &opCode{"BIT 5 B", 1, bit}
	cpu.opCodesCB[0x69] = &opCode{"BIT 5 C", 1, bit}
	cpu.opCodesCB[0x6A] = &opCode{"BIT 5 D", 1, bit}
	cpu.opCodesCB[0x6B] = &opCode{"BIT 5 E", 1, bit}
	cpu.opCodesCB[0x6C] = &opCode{"BIT 5 H", 1, bit}
	cpu.opCodesCB[0x6D] = &opCode{"BIT 5 L", 1, bit}
	cpu.opCodesCB[0x6E] = &opCode{"BIT 5 (HL)", 1, bitHL}
	cpu.opCodesCB[0x6F] = &opCode{"BIT 5 A", 1, bit}
	cpu.opCodesCB[0x70] = &opCode{"BIT 6 B", 1, bit}
	cpu.opCodesCB[0x71] = &opCode{"BIT 6 C", 1, bit}
	cpu.opCodesCB[0x72] = &opCode{"BIT 6 D", 1, bit}
	cpu.opCodesCB[0x73] = &opCode{"BIT 6 E", 1, bit}
	cpu.opCodesCB[0x74] = &opCode{"BIT 6 H", 1, bit}
	cpu.opCodesCB[0x75] = &opCode{"BIT 6 L", 1, bit}
	cpu.opCodesCB[0x76] = &opCode{"BIT 6 (HL)", 1, bitHL}
	cpu.opCodesCB[0x77] = &opCode{"BIT 6 A", 1, bit}
	cpu.opCodesCB[0x78] = &opCode{"BIT 7 B", 1, bit}
	cpu.opCodesCB[0x79] = &opCode{"BIT 7 C", 1, bit}
	cpu.opCodesCB[0x7A] = &opCode{"BIT 7 D", 1, bit}
	cpu.opCodesCB[0x7B] = &opCode{"BIT 7 E", 1, bit}
	cpu.opCodesCB[0x7C] = &opCode{"BIT 7 H", 1, bit}
	cpu.opCodesCB[0x7D] = &opCode{"BIT 7 L", 1, bit}
	cpu.opCodesCB[0x7E] = &opCode{"BIT 7 (HL)", 1, bitHL}
	cpu.opCodesCB[0x7F] = &opCode{"BIT 7 A", 1, bit}
	cpu.opCodesCB[0x80] = &opCode{"RES 0 B", 1, res}
	cpu.opCodesCB[0x81] = &opCode{"RES 0 C", 1, res}
	cpu.opCodesCB[0x82] = &opCode{"RES 0 D", 1, res}
	cpu.opCodesCB[0x83] = &opCode{"RES 0 E", 1, res}
	cpu.opCodesCB[0x84] = &opCode{"RES 0 H", 1, res}
	cpu.opCodesCB[0x85] = &opCode{"RES 0 L", 1, res}
	cpu.opCodesCB[0x86] = &opCode{"RES 0 (HL)", 1, resHL}
	cpu.opCodesCB[0x87] = &opCode{"RES 0 A", 1, res}
	cpu.opCodesCB[0x88] = &opCode{"RES 1 B", 1, res}
	cpu.opCodesCB[0x89] = &opCode{"RES 1 C", 1, res}
	cpu.opCodesCB[0x8A] = &opCode{"RES 1 D", 1, res}
	cpu.opCodesCB[0x8B] = &opCode{"RES 1 E", 1, res}
	cpu.opCodesCB[0x8C] = &opCode{"RES 1 H", 1, res}
	cpu.opCodesCB[0x8D] = &opCode{"RES 1 L", 1, res}
	cpu.opCodesCB[0x8E] = &opCode{"RES 1 (HL)", 1, resHL}
	cpu.opCodesCB[0x8F] = &opCode{"RES 1 A", 1, res}
	cpu.opCodesCB[0x90] = &opCode{"RES 2 B", 1, res}
	cpu.opCodesCB[0x91] = &opCode{"RES 2 C", 1, res}
	cpu.opCodesCB[0x92] = &opCode{"RES 2 D", 1, res}
	cpu.opCodesCB[0x93] = &opCode{"RES 2 E", 1, res}
	cpu.opCodesCB[0x94] = &opCode{"RES 2 H", 1, res}
	cpu.opCodesCB[0x95] = &opCode{"RES 2 L", 1, res}
	cpu.opCodesCB[0x96] = &opCode{"RES 2 (HL)", 1, resHL}
	cpu.opCodesCB[0x97] = &opCode{"RES 2 A", 1, res}
	cpu.opCodesCB[0x98] = &opCode{"RES 3 B", 1, res}
	cpu.opCodesCB[0x99] = &opCode{"RES 3 C", 1, res}
	cpu.opCodesCB[0x9A] = &opCode{"RES 3 D", 1, res}
	cpu.opCodesCB[0x9B] = &opCode{"RES 3 E", 1, res}
	cpu.opCodesCB[0x9C] = &opCode{"RES 3 H", 1, res}
	cpu.opCodesCB[0x9D] = &opCode{"RES 3 L", 1, res}
	cpu.opCodesCB[0x9E] = &opCode{"RES 3 (HL)", 1, resHL}
	cpu.opCodesCB[0x9F] = &opCode{"RES 3 A", 1, res}
	cpu.opCodesCB[0xA0] = &opCode{"RES 4 B", 1, res}
	cpu.opCodesCB[0xA1] = &opCode{"RES 4 C", 1, res}
	cpu.opCodesCB[0xA2] = &opCode{"RES 4 D", 1, res}
	cpu.opCodesCB[0xA3] = &opCode{"RES 4 E", 1, res}
	cpu.opCodesCB[0xA4] = &opCode{"RES 4 H", 1, res}
	cpu.opCodesCB[0xA5] = &opCode{"RES 4 L", 1, res}
	cpu.opCodesCB[0xA6] = &opCode{"RES 4 (HL)", 1, resHL}
	cpu.opCodesCB[0xA7] = &opCode{"RES 4 A", 1, res}
	cpu.opCodesCB[0xA8] = &opCode{"RES 5 B", 1, res}
	cpu.opCodesCB[0xA9] = &opCode{"RES 5 C", 1, res}
	cpu.opCodesCB[0xAA] = &opCode{"RES 5 D", 1, res}
	cpu.opCodesCB[0xAB] = &opCode{"RES 5 E", 1, res}
	cpu.opCodesCB[0xAC] = &opCode{"RES 5 H", 1, res}
	cpu.opCodesCB[0xAD] = &opCode{"RES 5 L", 1, res}
	cpu.opCodesCB[0xAE] = &opCode{"RES 5 (HL)", 1, resHL}
	cpu.opCodesCB[0xAF] = &opCode{"RES 5 A", 1, res}
	cpu.opCodesCB[0xB0] = &opCode{"RES 6 B", 1, res}
	cpu.opCodesCB[0xB1] = &opCode{"RES 6 C", 1, res}
	cpu.opCodesCB[0xB2] = &opCode{"RES 6 D", 1, res}
	cpu.opCodesCB[0xB3] = &opCode{"RES 6 E", 1, res}
	cpu.opCodesCB[0xB4] = &opCode{"RES 6 H", 1, res}
	cpu.opCodesCB[0xB5] = &opCode{"RES 6 L", 1, res}
	cpu.opCodesCB[0xB6] = &opCode{"RES 6 (HL)", 1, resHL}
	cpu.opCodesCB[0xB7] = &opCode{"RES 6 A", 1, res}
	cpu.opCodesCB[0xB8] = &opCode{"RES 7 B", 1, res}
	cpu.opCodesCB[0xB9] = &opCode{"RES 7 C", 1, res}
	cpu.opCodesCB[0xBA] = &opCode{"RES 7 D", 1, res}
	cpu.opCodesCB[0xBB] = &opCode{"RES 7 E", 1, res}
	cpu.opCodesCB[0xBC] = &opCode{"RES 7 H", 1, res}
	cpu.opCodesCB[0xBD] = &opCode{"RES 7 L", 1, res}
	cpu.opCodesCB[0xBE] = &opCode{"RES 7 (HL)", 1, resHL}
	cpu.opCodesCB[0xBF] = &opCode{"RES 7 A", 1, res}
	cpu.opCodesCB[0xC0] = &opCode{"SET 0 B", 1, set}
	cpu.opCodesCB[0xC1] = &opCode{"SET 0 C", 1, set}
	cpu.opCodesCB[0xC2] = &opCode{"SET 0 D", 1, set}
	cpu.opCodesCB[0xC3] = &opCode{"SET 0 E", 1, set}
	cpu.opCodesCB[0xC4] = &opCode{"SET 0 H", 1, set}
	cpu.opCodesCB[0xC5] = &opCode{"SET 0 L", 1, set}
	cpu.opCodesCB[0xC6] = &opCode{"SET 0 (HL)", 1, setHL}
	cpu.opCodesCB[0xC7] = &opCode{"SET 0 A", 1, set}
	cpu.opCodesCB[0xC8] = &opCode{"SET 1 B", 1, set}
	cpu.opCodesCB[0xC9] = &opCode{"SET 1 C", 1, set}
	cpu.opCodesCB[0xCA] = &opCode{"SET 1 D", 1, set}
	cpu.opCodesCB[0xCB] = &opCode{"SET 1 E", 1, set}
	cpu.opCodesCB[0xCC] = &opCode{"SET 1 H", 1, set}
	cpu.opCodesCB[0xCD] = &opCode{"SET 1 L", 1, set}
	cpu.opCodesCB[0xCE] = &opCode{"SET 1 (HL)", 1, setHL}
	cpu.opCodesCB[0xCF] = &opCode{"SET 1 A", 1, set}
	cpu.opCodesCB[0xD0] = &opCode{"SET 2 B", 1, set}
	cpu.opCodesCB[0xD1] = &opCode{"SET 2 C", 1, set}
	cpu.opCodesCB[0xD2] = &opCode{"SET 2 D", 1, set}
	cpu.opCodesCB[0xD3] = &opCode{"SET 2 E", 1, set}
	cpu.opCodesCB[0xD4] = &opCode{"SET 2 H", 1, set}
	cpu.opCodesCB[0xD5] = &opCode{"SET 2 L", 1, set}
	cpu.opCodesCB[0xD6] = &opCode{"SET 2 (HL)", 1, setHL}
	cpu.opCodesCB[0xD7] = &opCode{"SET 2 A", 1, set}
	cpu.opCodesCB[0xD8] = &opCode{"SET 3 B", 1, set}
	cpu.opCodesCB[0xD9] = &opCode{"SET 3 C", 1, set}
	cpu.opCodesCB[0xDA] = &opCode{"SET 3 D", 1, set}
	cpu.opCodesCB[0xDB] = &opCode{"SET 3 E", 1, set}
	cpu.opCodesCB[0xDC] = &opCode{"SET 3 H", 1, set}
	cpu.opCodesCB[0xDD] = &opCode{"SET 3 L", 1, set}
	cpu.opCodesCB[0xDE] = &opCode{"SET 3 (HL)", 1, setHL}
	cpu.opCodesCB[0xDF] = &opCode{"SET 3 A", 1, set}
	cpu.opCodesCB[0xE0] = &opCode{"SET 4 B", 1, set}
	cpu.opCodesCB[0xE1] = &opCode{"SET 4 C", 1, set}
	cpu.opCodesCB[0xE2] = &opCode{"SET 4 D", 1, set}
	cpu.opCodesCB[0xE3] = &opCode{"SET 4 E", 1, set}
	cpu.opCodesCB[0xE4] = &opCode{"SET 4 H", 1, set}
	cpu.opCodesCB[0xE5] = &opCode{"SET 4 L", 1, set}
	cpu.opCodesCB[0xE6] = &opCode{"SET 4 (HL)", 1, setHL}
	cpu.opCodesCB[0xE7] = &opCode{"SET 4 A", 1, set}
	cpu.opCodesCB[0xE8] = &opCode{"SET 5 B", 1, set}
	cpu.opCodesCB[0xE9] = &opCode{"SET 5 C", 1, set}
	cpu.opCodesCB[0xEA] = &opCode{"SET 5 D", 1, set}
	cpu.opCodesCB[0xEB] = &opCode{"SET 5 E", 1, set}
	cpu.opCodesCB[0xEC] = &opCode{"SET 5 H", 1, set}
	cpu.opCodesCB[0xED] = &opCode{"SET 5 L", 1, set}
	cpu.opCodesCB[0xEE] = &opCode{"SET 5 (HL)", 1, setHL}
	cpu.opCodesCB[0xEF] = &opCode{"SET 5 A", 1, set}
	cpu.opCodesCB[0xF0] = &opCode{"SET 6 B", 1, set}
	cpu.opCodesCB[0xF1] = &opCode{"SET 6 C", 1, set}
	cpu.opCodesCB[0xF2] = &opCode{"SET 6 D", 1, set}
	cpu.opCodesCB[0xF3] = &opCode{"SET 6 E", 1, set}
	cpu.opCodesCB[0xF4] = &opCode{"SET 6 H", 1, set}
	cpu.opCodesCB[0xF5] = &opCode{"SET 6 L", 1, set}
	cpu.opCodesCB[0xF6] = &opCode{"SET 6 (HL)", 1, setHL}
	cpu.opCodesCB[0xF7] = &opCode{"SET 6 A", 1, set}
	cpu.opCodesCB[0xF8] = &opCode{"SET 7 B", 1, set}
	cpu.opCodesCB[0xF9] = &opCode{"SET 7 C", 1, set}
	cpu.opCodesCB[0xFA] = &opCode{"SET 7 D", 1, set}
	cpu.opCodesCB[0xFB] = &opCode{"SET 7 E", 1, set}
	cpu.opCodesCB[0xFC] = &opCode{"SET 7 H", 1, set}
	cpu.opCodesCB[0xFD] = &opCode{"SET 7 L", 1, set}
	cpu.opCodesCB[0xFE] = &opCode{"SET 7 (HL)", 1, setHL}
	cpu.opCodesCB[0xFF] = &opCode{"SET 7 A", 1, set}
}

func decodeCB(cpu *lr35902) {
	cpu.scheduler.append(newFetch(cpu.opCodesCB))
}
