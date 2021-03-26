package lr35902

import (
	"fmt"
	"strconv"
	"strings"
)

type opCode struct {
	Ins string
	Len byte
	f   lr35902f
}

var OPCodes = make([]*opCode, 256)
var OPCodesCB = make([]*opCode, 256)

func init() {
	OPCodes[0x00] = &opCode{"NOP", 1, func(cpu *lr35902) {}}
	OPCodes[0x01] = &opCode{"LD BC, nn", 3, ldDDnn}
	OPCodes[0x02] = &opCode{"LD (BC), A", 1, ldBCa}
	OPCodes[0x03] = &opCode{"INC BC", 1, incSS}
	OPCodes[0x04] = &opCode{"INC B", 1, incR}
	OPCodes[0x05] = &opCode{"DEC B", 1, decR}
	OPCodes[0x06] = &opCode{"LD B, n", 2, ldRn}
	OPCodes[0x07] = &opCode{"RLCA", 1, rlca}
	OPCodes[0x08] = &opCode{"LD (nn), SP", 3, ldNNsp}
	OPCodes[0x09] = &opCode{"ADD HL,BC", 1, addHLss}
	OPCodes[0x0A] = &opCode{"LD A,(BC)", 1, ldAbc}
	OPCodes[0x0B] = &opCode{"DEC BC", 1, decSS}
	OPCodes[0x0C] = &opCode{"INC C", 1, incR}
	OPCodes[0x0D] = &opCode{"DEC C", 1, decR}
	OPCodes[0x0E] = &opCode{"LD C, n", 2, ldRn}
	OPCodes[0x0F] = &opCode{"RRCA", 1, rrca}
	OPCodes[0x0F] = &opCode{"STOP", 1, func(cpu *lr35902) { panic(fmt.Sprintf("panic on 0x%04X", cpu.regs.PC)) }}
	OPCodes[0x11] = &opCode{"LD DE, nn", 3, ldDDnn}
	OPCodes[0x12] = &opCode{"LD (DE), A", 1, ldDEa}
	OPCodes[0x13] = &opCode{"INC DE", 1, incSS}
	OPCodes[0x14] = &opCode{"INC D", 1, incR}
	OPCodes[0x15] = &opCode{"DEC D", 1, decR}
	OPCodes[0x16] = &opCode{"LD D, n", 2, ldRn}
	OPCodes[0x17] = &opCode{"RLA", 1, rla}
	OPCodes[0x18] = &opCode{"JR e", 2, jr}
	OPCodes[0x19] = &opCode{"ADD HL,DE", 1, addHLss}
	OPCodes[0x1A] = &opCode{"LD A,(DE)", 1, ldAde}
	OPCodes[0x1B] = &opCode{"DEC DE", 1, decSS}
	OPCodes[0x1C] = &opCode{"INC E", 1, incR}
	OPCodes[0x1D] = &opCode{"DEC E", 1, decR}
	OPCodes[0x1E] = &opCode{"LD E, n", 2, ldRn}
	OPCodes[0x1F] = &opCode{"RRA", 1, rra}
	OPCodes[0x20] = &opCode{"JR NZ, e", 2, jrnz}
	OPCodes[0x21] = &opCode{"LD HL, nn", 3, ldDDnn}
	OPCodes[0x22] = &opCode{"LDI (HL),a", 1, ldiHLa}
	OPCodes[0x23] = &opCode{"INC HL", 1, incSS}
	OPCodes[0x24] = &opCode{"INC H", 1, incR}
	OPCodes[0x25] = &opCode{"DEC H", 1, decR}
	OPCodes[0x26] = &opCode{"LD H, n", 2, ldRn}
	OPCodes[0x27] = &opCode{"DAA", 1, daa}
	OPCodes[0x28] = &opCode{"JR Z, e", 2, jrz}
	OPCodes[0x29] = &opCode{"ADD HL,HL", 1, addHLss}
	OPCodes[0x2A] = &opCode{"LDI A,(HL)", 1, ldiAhl}
	OPCodes[0x2B] = &opCode{"DEC HL", 1, decSS}
	OPCodes[0x2C] = &opCode{"INC L", 1, incR}
	OPCodes[0x2D] = &opCode{"DEC L", 1, decR}
	OPCodes[0x2E] = &opCode{"LD L, n", 2, ldRn}
	OPCodes[0x2F] = &opCode{"CPL", 1, cpl}
	OPCodes[0x30] = &opCode{"JR NC, e", 2, jrnc}
	OPCodes[0x31] = &opCode{"LD SP, nn", 3, ldDDnn}
	OPCodes[0x32] = &opCode{"LDD (HL),a", 1, lddHLa}
	OPCodes[0x33] = &opCode{"INC SP", 1, incSS}
	OPCodes[0x34] = &opCode{"INC (HL)", 1, incHL}
	OPCodes[0x35] = &opCode{"DEC (HL)", 1, decHL}
	OPCodes[0x36] = &opCode{"LD (HL), n", 2, ldToHLn}
	OPCodes[0x37] = &opCode{"SCF", 1, scf}
	OPCodes[0x38] = &opCode{"JR C, e", 2, jrc}
	OPCodes[0x39] = &opCode{"ADD HL,SP", 1, addHLss}
	OPCodes[0x3A] = &opCode{"LDD A,(HL)", 1, lddAhl}
	OPCodes[0x3B] = &opCode{"DEC SP", 1, decSS}
	OPCodes[0x3C] = &opCode{"INC A", 1, incR}
	OPCodes[0x3D] = &opCode{"DEC A", 1, decR}
	OPCodes[0x3E] = &opCode{"LD A, n", 2, ldRn}
	OPCodes[0x3F] = &opCode{"CCF", 1, ccf}
	OPCodes[0x40] = &opCode{"LD B, B", 1, func(cpu *lr35902) {}}
	OPCodes[0x41] = &opCode{"LD B, C", 1, func(cpu *lr35902) { cpu.regs.B = cpu.regs.C }}
	OPCodes[0x42] = &opCode{"LD B, D", 1, func(cpu *lr35902) { cpu.regs.B = cpu.regs.D }}
	OPCodes[0x43] = &opCode{"LD B, E", 1, func(cpu *lr35902) { cpu.regs.B = cpu.regs.E }}
	OPCodes[0x44] = &opCode{"LD B, H", 1, func(cpu *lr35902) { cpu.regs.B = cpu.regs.H }}
	OPCodes[0x45] = &opCode{"LD B, L", 1, func(cpu *lr35902) { cpu.regs.B = cpu.regs.L }}
	OPCodes[0x46] = &opCode{"LD B, (HL)", 1, ldFromHL}
	OPCodes[0x47] = &opCode{"LD B, A", 1, func(cpu *lr35902) { cpu.regs.B = cpu.regs.A }}
	OPCodes[0x48] = &opCode{"LD C, B", 1, func(cpu *lr35902) { cpu.regs.C = cpu.regs.B }}
	OPCodes[0x49] = &opCode{"LD C, C", 1, func(cpu *lr35902) {}}
	OPCodes[0x4A] = &opCode{"LD C, D", 1, func(cpu *lr35902) { cpu.regs.C = cpu.regs.D }}
	OPCodes[0x4B] = &opCode{"LD C, E", 1, func(cpu *lr35902) { cpu.regs.C = cpu.regs.E }}
	OPCodes[0x4C] = &opCode{"LD C, H", 1, func(cpu *lr35902) { cpu.regs.C = cpu.regs.H }}
	OPCodes[0x4D] = &opCode{"LD C, L", 1, func(cpu *lr35902) { cpu.regs.C = cpu.regs.L }}
	OPCodes[0x4E] = &opCode{"LD C, (HL)", 1, ldFromHL}
	OPCodes[0x4F] = &opCode{"LD C, A", 1, func(cpu *lr35902) { cpu.regs.C = cpu.regs.A }}
	OPCodes[0x50] = &opCode{"LD D, B", 1, func(cpu *lr35902) { cpu.regs.D = cpu.regs.B }}
	OPCodes[0x51] = &opCode{"LD D, C", 1, func(cpu *lr35902) { cpu.regs.D = cpu.regs.C }}
	OPCodes[0x52] = &opCode{"LD D, D", 1, func(cpu *lr35902) {}}
	OPCodes[0x53] = &opCode{"LD D, E", 1, func(cpu *lr35902) { cpu.regs.D = cpu.regs.E }}
	OPCodes[0x54] = &opCode{"LD D, H", 1, func(cpu *lr35902) { cpu.regs.D = cpu.regs.H }}
	OPCodes[0x55] = &opCode{"LD D, L", 1, func(cpu *lr35902) { cpu.regs.D = cpu.regs.L }}
	OPCodes[0x56] = &opCode{"LD D, (HL)", 1, ldFromHL}
	OPCodes[0x57] = &opCode{"LD D, A", 1, func(cpu *lr35902) { cpu.regs.D = cpu.regs.A }}
	OPCodes[0x58] = &opCode{"LD E, B", 1, func(cpu *lr35902) { cpu.regs.E = cpu.regs.B }}
	OPCodes[0x59] = &opCode{"LD E, C", 1, func(cpu *lr35902) { cpu.regs.E = cpu.regs.C }}
	OPCodes[0x5A] = &opCode{"LD E, D", 1, func(cpu *lr35902) { cpu.regs.E = cpu.regs.D }}
	OPCodes[0x5B] = &opCode{"LD E, E", 1, func(cpu *lr35902) {}}
	OPCodes[0x5C] = &opCode{"LD E, H", 1, func(cpu *lr35902) { cpu.regs.E = cpu.regs.H }}
	OPCodes[0x5D] = &opCode{"LD E, L", 1, func(cpu *lr35902) { cpu.regs.E = cpu.regs.L }}
	OPCodes[0x5E] = &opCode{"LD E, (HL)", 1, ldFromHL}
	OPCodes[0x5F] = &opCode{"LD E, A", 1, func(cpu *lr35902) { cpu.regs.E = cpu.regs.A }}
	OPCodes[0x60] = &opCode{"LD H, B", 1, func(cpu *lr35902) { cpu.regs.H = cpu.regs.B }}
	OPCodes[0x61] = &opCode{"LD H, C", 1, func(cpu *lr35902) { cpu.regs.H = cpu.regs.C }}
	OPCodes[0x62] = &opCode{"LD H, D", 1, func(cpu *lr35902) { cpu.regs.H = cpu.regs.D }}
	OPCodes[0x63] = &opCode{"LD H, E", 1, func(cpu *lr35902) { cpu.regs.H = cpu.regs.E }}
	OPCodes[0x64] = &opCode{"LD H, H", 1, func(cpu *lr35902) {}}
	OPCodes[0x65] = &opCode{"LD H, L", 1, func(cpu *lr35902) { cpu.regs.H = cpu.regs.L }}
	OPCodes[0x66] = &opCode{"LD H, (HL)", 1, ldFromHL}
	OPCodes[0x67] = &opCode{"LD H, A", 1, func(cpu *lr35902) { cpu.regs.H = cpu.regs.A }}
	OPCodes[0x68] = &opCode{"LD L, B", 1, func(cpu *lr35902) { cpu.regs.L = cpu.regs.B }}
	OPCodes[0x69] = &opCode{"LD L, C", 1, func(cpu *lr35902) { cpu.regs.L = cpu.regs.C }}
	OPCodes[0x6A] = &opCode{"LD L, D", 1, func(cpu *lr35902) { cpu.regs.L = cpu.regs.D }}
	OPCodes[0x6B] = &opCode{"LD L, E", 1, func(cpu *lr35902) { cpu.regs.L = cpu.regs.E }}
	OPCodes[0x6C] = &opCode{"LD L, H", 1, func(cpu *lr35902) { cpu.regs.L = cpu.regs.H }}
	OPCodes[0x6D] = &opCode{"LD L, L", 1, func(cpu *lr35902) {}}
	OPCodes[0x6E] = &opCode{"LD L, (HL)", 1, ldFromHL}
	OPCodes[0x6F] = &opCode{"LD L, A", 1, func(cpu *lr35902) { cpu.regs.L = cpu.regs.A }}
	OPCodes[0x70] = &opCode{"LD (HL), B", 1, ldToHL}
	OPCodes[0x71] = &opCode{"LD (HL), C", 1, ldToHL}
	OPCodes[0x72] = &opCode{"LD (HL), D", 1, ldToHL}
	OPCodes[0x73] = &opCode{"LD (HL), E", 1, ldToHL}
	OPCodes[0x74] = &opCode{"LD (HL), H", 1, ldToHL}
	OPCodes[0x75] = &opCode{"LD (HL), L", 1, ldToHL}
	OPCodes[0x76] = &opCode{"HALT", 1, halt}
	OPCodes[0x77] = &opCode{"LD (HL), A", 1, ldToHL}
	OPCodes[0x78] = &opCode{"LD A, B", 1, func(cpu *lr35902) { cpu.regs.A = cpu.regs.B }}
	OPCodes[0x79] = &opCode{"LD A, C", 1, func(cpu *lr35902) { cpu.regs.A = cpu.regs.C }}
	OPCodes[0x7A] = &opCode{"LD A, D", 1, func(cpu *lr35902) { cpu.regs.A = cpu.regs.D }}
	OPCodes[0x7B] = &opCode{"LD A, E", 1, func(cpu *lr35902) { cpu.regs.A = cpu.regs.E }}
	OPCodes[0x7C] = &opCode{"LD A, H", 1, func(cpu *lr35902) { cpu.regs.A = cpu.regs.H }}
	OPCodes[0x7D] = &opCode{"LD A, L", 1, func(cpu *lr35902) { cpu.regs.A = cpu.regs.L }}
	OPCodes[0x7E] = &opCode{"LD A, (HL)", 1, ldFromHL}
	OPCodes[0x7F] = &opCode{"LD A, A", 1, func(cpu *lr35902) {}}
	OPCodes[0x80] = &opCode{"ADD A, B", 1, addAr}
	OPCodes[0x81] = &opCode{"ADD A, C", 1, addAr}
	OPCodes[0x82] = &opCode{"ADD A, D", 1, addAr}
	OPCodes[0x83] = &opCode{"ADD A, E", 1, addAr}
	OPCodes[0x84] = &opCode{"ADD A, H", 1, addAr}
	OPCodes[0x85] = &opCode{"ADD A, L", 1, addAr}
	OPCodes[0x86] = &opCode{"ADD A, (HL)", 1, addAhl}
	OPCodes[0x87] = &opCode{"ADD A, A", 1, addAr}
	OPCodes[0x88] = &opCode{"ADC A, B", 1, adcAr}
	OPCodes[0x89] = &opCode{"ADC A, C", 1, adcAr}
	OPCodes[0x8A] = &opCode{"ADC A, D", 1, adcAr}
	OPCodes[0x8B] = &opCode{"ADC A, E", 1, adcAr}
	OPCodes[0x8C] = &opCode{"ADC A, H", 1, adcAr}
	OPCodes[0x8D] = &opCode{"ADC A, L", 1, adcAr}
	OPCodes[0x8E] = &opCode{"ADC A, (HL)", 1, adcAhl}
	OPCodes[0x8F] = &opCode{"ADC A, A", 1, adcAr}
	OPCodes[0x90] = &opCode{"SUB A, B", 1, subAr}
	OPCodes[0x91] = &opCode{"SUB A, C", 1, subAr}
	OPCodes[0x92] = &opCode{"SUB A, D", 1, subAr}
	OPCodes[0x93] = &opCode{"SUB A, E", 1, subAr}
	OPCodes[0x94] = &opCode{"SUB A, H", 1, subAr}
	OPCodes[0x95] = &opCode{"SUB A, L", 1, subAr}
	OPCodes[0x96] = &opCode{"SUB A, (HL)", 1, subAhl}
	OPCodes[0x97] = &opCode{"SUB A, A", 1, subAr}
	OPCodes[0x98] = &opCode{"SUC A, B", 1, sbcAr}
	OPCodes[0x99] = &opCode{"SUC A, C", 1, sbcAr}
	OPCodes[0x9A] = &opCode{"SUC A, D", 1, sbcAr}
	OPCodes[0x9B] = &opCode{"SUC A, E", 1, sbcAr}
	OPCodes[0x9C] = &opCode{"SUC A, H", 1, sbcAr}
	OPCodes[0x9D] = &opCode{"SUC A, L", 1, sbcAr}
	OPCodes[0x9E] = &opCode{"SBC A, (HL)", 1, sbcAhl}
	OPCodes[0x9F] = &opCode{"SUC A, A", 1, sbcAr}
	OPCodes[0xA0] = &opCode{"AND B", 1, andAr}
	OPCodes[0xA1] = &opCode{"AND C", 1, andAr}
	OPCodes[0xA2] = &opCode{"AND D", 1, andAr}
	OPCodes[0xA3] = &opCode{"AND E", 1, andAr}
	OPCodes[0xA4] = &opCode{"AND H", 1, andAr}
	OPCodes[0xA5] = &opCode{"AND L", 1, andAr}
	OPCodes[0xA6] = &opCode{"AND (HL)", 1, andAhl}
	OPCodes[0xA7] = &opCode{"AND A", 1, andAr}
	OPCodes[0xA8] = &opCode{"XOR B", 1, xorAr}
	OPCodes[0xA9] = &opCode{"XOR C", 1, xorAr}
	OPCodes[0xAA] = &opCode{"XOR D", 1, xorAr}
	OPCodes[0xAB] = &opCode{"XOR E", 1, xorAr}
	OPCodes[0xAC] = &opCode{"XOR H", 1, xorAr}
	OPCodes[0xAD] = &opCode{"XOR L", 1, xorAr}
	OPCodes[0xAE] = &opCode{"XOR (HL)", 1, xorAhl}
	OPCodes[0xAF] = &opCode{"XOR A", 1, xorAr}
	OPCodes[0xB0] = &opCode{"OR B", 1, orAr}
	OPCodes[0xB1] = &opCode{"OR C", 1, orAr}
	OPCodes[0xB2] = &opCode{"OR D", 1, orAr}
	OPCodes[0xB3] = &opCode{"OR E", 1, orAr}
	OPCodes[0xB4] = &opCode{"OR H", 1, orAr}
	OPCodes[0xB5] = &opCode{"OR L", 1, orAr}
	OPCodes[0xB6] = &opCode{"OR (HL)", 1, orAhl}
	OPCodes[0xB7] = &opCode{"OR A", 1, orAr}
	OPCodes[0xB8] = &opCode{"CP B", 1, cpR}
	OPCodes[0xB9] = &opCode{"CP C", 1, cpR}
	OPCodes[0xBA] = &opCode{"CP D", 1, cpR}
	OPCodes[0xBB] = &opCode{"CP E", 1, cpR}
	OPCodes[0xBC] = &opCode{"CP H", 1, cpR}
	OPCodes[0xBD] = &opCode{"CP L", 1, cpR}
	OPCodes[0xBE] = &opCode{"CP (HL)", 1, cpHl}
	OPCodes[0xBF] = &opCode{"CP A", 1, cpR}
	OPCodes[0xC0] = &opCode{"RET NZ", 1, retCC}
	OPCodes[0xC1] = &opCode{"POP BC", 1, popSS}
	OPCodes[0xC2] = &opCode{"JP NZ, nn", 3, jpCC}
	OPCodes[0xC3] = &opCode{"JP nn", 3, func(cpu *lr35902) { cpu.regs.PC = cpu.fetched.nn }}
	OPCodes[0xC4] = &opCode{"CALL NZ, nn", 3, callCC}
	OPCodes[0xC5] = &opCode{"PUSH BC", 1, pushSS}
	OPCodes[0xC6] = &opCode{"ADD A, n", 2, func(cpu *lr35902) { cpu.addA(cpu.fetched.n) }}
	OPCodes[0xC7] = &opCode{"RST 0x0", 1, rstP}
	OPCodes[0xC8] = &opCode{"RET Z", 1, retCC}
	OPCodes[0xC9] = &opCode{"RET", 1, ret}
	OPCodes[0xCA] = &opCode{"JP Z, nn", 3, jpCC}
	OPCodes[0xCB] = &opCode{"CB", 1, decodeCB}
	OPCodes[0xCC] = &opCode{"CALL Z, nn", 3, callCC}
	OPCodes[0xCD] = &opCode{"CALL nn", 3, call}
	OPCodes[0xCE] = &opCode{"ADC A, n", 2, func(cpu *lr35902) { cpu.adcA(cpu.fetched.n) }}
	OPCodes[0xCF] = &opCode{"RST 0x8", 1, rstP}
	OPCodes[0xD0] = &opCode{"RET NC", 1, retCC}
	OPCodes[0xD1] = &opCode{"POP DE", 1, popSS}
	OPCodes[0xD2] = &opCode{"JP NC, nn", 3, jpCC}
	OPCodes[0xD4] = &opCode{"CALL NC, nn", 3, callCC}
	OPCodes[0xD5] = &opCode{"PUSH DE", 1, pushSS}
	OPCodes[0xD6] = &opCode{"SUB A, n", 2, func(cpu *lr35902) { cpu.subA(cpu.fetched.n) }}
	OPCodes[0xD7] = &opCode{"RST 0x10", 1, rstP}
	OPCodes[0xD8] = &opCode{"RET C", 1, retCC}
	OPCodes[0xD9] = &opCode{"RETI", 1, reti}
	OPCodes[0xDA] = &opCode{"JP C, nn", 3, jpCC}
	OPCodes[0xDC] = &opCode{"CALL C, nn", 3, callCC}
	OPCodes[0xDE] = &opCode{"SBC A, nn", 3, func(cpu *lr35902) { cpu.sbcA(cpu.fetched.n) }}
	OPCodes[0xDF] = &opCode{"RST 0x18", 1, rstP}
	OPCodes[0xE0] = &opCode{"LD (0xff00+n), A", 2, ldhNa}
	OPCodes[0xE1] = &opCode{"POP HL", 1, popSS}
	OPCodes[0xE2] = &opCode{"LD (0xff00+C), A", 1, ldhCa}
	OPCodes[0xE5] = &opCode{"PUSH HL", 1, pushSS}
	OPCodes[0xE6] = &opCode{"AND n", 2, func(cpu *lr35902) { cpu.and(cpu.fetched.n) }}
	OPCodes[0xE7] = &opCode{"RST 0x20", 1, rstP}
	OPCodes[0xE8] = &opCode{"ADD SP,n", 2, addSPn}
	OPCodes[0xE9] = &opCode{"JP HL", 1, func(cpu *lr35902) { cpu.regs.PC = cpu.regs.HL.Get() }}
	OPCodes[0xEA] = &opCode{"LD (nn), A", 3, ldNNa}
	OPCodes[0xEE] = &opCode{"XOR A, n", 2, func(cpu *lr35902) { cpu.xor(cpu.fetched.n) }}
	OPCodes[0xEF] = &opCode{"RST 0x28", 1, rstP}
	OPCodes[0xF0] = &opCode{"LD A, (0xff00+n)", 2, ldhAn}
	OPCodes[0xF1] = &opCode{"POP AF", 1, popSS}
	OPCodes[0xF2] = &opCode{"LD A, (0xff00+C)", 1, ldhAc}
	OPCodes[0xF3] = &opCode{"DI", 1, func(cpu *lr35902) { cpu.regs.IME = false }}
	OPCodes[0xF5] = &opCode{"PUSH AF", 1, pushSS}
	OPCodes[0xF6] = &opCode{"OR n", 2, func(cpu *lr35902) { cpu.or(cpu.fetched.n) }}
	OPCodes[0xF7] = &opCode{"RST 0x30", 1, rstP}
	OPCodes[0xF8] = &opCode{"LD HL,(SP+n)", 2, ldHLspE}
	OPCodes[0xF9] = &opCode{"LD SP, HL", 1, func(cpu *lr35902) { cpu.regs.SP.Set(cpu.regs.HL.Get()) }}
	OPCodes[0xFA] = &opCode{"LD A, (nn)", 3, ldAnn}
	OPCodes[0xFB] = &opCode{"EI", 1, func(cpu *lr35902) { cpu.regs.IME = true }}
	OPCodes[0xFE] = &opCode{"CP A, n", 2, func(cpu *lr35902) { cpu.cp(cpu.fetched.n) }}
	OPCodes[0xFF] = &opCode{"RST 0x38", 1, rstP}

	OPCodesCB[0x00] = &opCode{"RLC B", 1, cbR}
	OPCodesCB[0x01] = &opCode{"RLC C", 1, cbR}
	OPCodesCB[0x02] = &opCode{"RLC D", 1, cbR}
	OPCodesCB[0x03] = &opCode{"RLC E", 1, cbR}
	OPCodesCB[0x04] = &opCode{"RLC H", 1, cbR}
	OPCodesCB[0x05] = &opCode{"RLC L", 1, cbR}
	OPCodesCB[0x06] = &opCode{"RLC (HL)", 1, cbHL}
	OPCodesCB[0x07] = &opCode{"RLC A", 1, cbR}
	OPCodesCB[0x08] = &opCode{"RRC B", 1, cbR}
	OPCodesCB[0x09] = &opCode{"RRC C", 1, cbR}
	OPCodesCB[0x0A] = &opCode{"RRC D", 1, cbR}
	OPCodesCB[0x0B] = &opCode{"RRC E", 1, cbR}
	OPCodesCB[0x0C] = &opCode{"RRC H", 1, cbR}
	OPCodesCB[0x0D] = &opCode{"RRC L", 1, cbR}
	OPCodesCB[0x0E] = &opCode{"RRC (HL)", 1, cbHL}
	OPCodesCB[0x0F] = &opCode{"RRC A", 1, cbR}
	OPCodesCB[0x10] = &opCode{"RL B", 1, cbR}
	OPCodesCB[0x11] = &opCode{"RL C", 1, cbR}
	OPCodesCB[0x12] = &opCode{"RL D", 1, cbR}
	OPCodesCB[0x13] = &opCode{"RL E", 1, cbR}
	OPCodesCB[0x14] = &opCode{"RL H", 1, cbR}
	OPCodesCB[0x15] = &opCode{"RL L", 1, cbR}
	OPCodesCB[0x16] = &opCode{"RL (HL)", 1, cbHL}
	OPCodesCB[0x17] = &opCode{"RL A", 1, cbR}
	OPCodesCB[0x18] = &opCode{"RR B", 1, cbR}
	OPCodesCB[0x19] = &opCode{"RR C", 1, cbR}
	OPCodesCB[0x1A] = &opCode{"RR D", 1, cbR}
	OPCodesCB[0x1B] = &opCode{"RR E", 1, cbR}
	OPCodesCB[0x1C] = &opCode{"RR H", 1, cbR}
	OPCodesCB[0x1D] = &opCode{"RR L", 1, cbR}
	OPCodesCB[0x1E] = &opCode{"RR (HL)", 1, cbHL}
	OPCodesCB[0x1F] = &opCode{"RR A", 1, cbR}
	OPCodesCB[0x20] = &opCode{"SLA B", 1, cbR}
	OPCodesCB[0x21] = &opCode{"SLA C", 1, cbR}
	OPCodesCB[0x22] = &opCode{"SLA D", 1, cbR}
	OPCodesCB[0x23] = &opCode{"SLA E", 1, cbR}
	OPCodesCB[0x24] = &opCode{"SLA H", 1, cbR}
	OPCodesCB[0x25] = &opCode{"SLA L", 1, cbR}
	OPCodesCB[0x26] = &opCode{"SLA (HL)", 1, cbHL}
	OPCodesCB[0x27] = &opCode{"SLA A", 1, cbR}
	OPCodesCB[0x28] = &opCode{"SRA B", 1, cbR}
	OPCodesCB[0x29] = &opCode{"SRA C", 1, cbR}
	OPCodesCB[0x2A] = &opCode{"SRA D", 1, cbR}
	OPCodesCB[0x2B] = &opCode{"SRA E", 1, cbR}
	OPCodesCB[0x2C] = &opCode{"SRA H", 1, cbR}
	OPCodesCB[0x2D] = &opCode{"SRA L", 1, cbR}
	OPCodesCB[0x2E] = &opCode{"SRA (HL)", 1, cbHL}
	OPCodesCB[0x2F] = &opCode{"SRA A", 1, cbR}
	OPCodesCB[0x30] = &opCode{"SWAP B", 1, swap}
	OPCodesCB[0x31] = &opCode{"SWAP C", 1, swap}
	OPCodesCB[0x32] = &opCode{"SWAP D", 1, swap}
	OPCodesCB[0x33] = &opCode{"SWAP E", 1, swap}
	OPCodesCB[0x34] = &opCode{"SWAP H", 1, swap}
	OPCodesCB[0x35] = &opCode{"SWAP L", 1, swap}
	OPCodesCB[0x36] = &opCode{"SWAP (HL)", 1, swapHL}
	OPCodesCB[0x37] = &opCode{"SWAP A", 1, swap}
	OPCodesCB[0x38] = &opCode{"SRL B", 1, cbR}
	OPCodesCB[0x39] = &opCode{"SRL C", 1, cbR}
	OPCodesCB[0x3A] = &opCode{"SRL D", 1, cbR}
	OPCodesCB[0x3B] = &opCode{"SRL E", 1, cbR}
	OPCodesCB[0x3C] = &opCode{"SRL H", 1, cbR}
	OPCodesCB[0x3D] = &opCode{"SRL L", 1, cbR}
	OPCodesCB[0x3E] = &opCode{"SRL (HL)", 1, cbHL}
	OPCodesCB[0x3F] = &opCode{"SRL A", 1, cbR}
	OPCodesCB[0x40] = &opCode{"BIT 0 B", 1, bit}
	OPCodesCB[0x41] = &opCode{"BIT 0 C", 1, bit}
	OPCodesCB[0x42] = &opCode{"BIT 0 D", 1, bit}
	OPCodesCB[0x43] = &opCode{"BIT 0 E", 1, bit}
	OPCodesCB[0x44] = &opCode{"BIT 0 H", 1, bit}
	OPCodesCB[0x45] = &opCode{"BIT 0 L", 1, bit}
	OPCodesCB[0x46] = &opCode{"BIT 0 (HL)", 1, bitHL}
	OPCodesCB[0x47] = &opCode{"BIT 0 A", 1, bit}
	OPCodesCB[0x48] = &opCode{"BIT 1 B", 1, bit}
	OPCodesCB[0x49] = &opCode{"BIT 1 C", 1, bit}
	OPCodesCB[0x4A] = &opCode{"BIT 1 D", 1, bit}
	OPCodesCB[0x4B] = &opCode{"BIT 1 E", 1, bit}
	OPCodesCB[0x4C] = &opCode{"BIT 1 H", 1, bit}
	OPCodesCB[0x4D] = &opCode{"BIT 1 L", 1, bit}
	OPCodesCB[0x4E] = &opCode{"BIT 1 (HL)", 1, bitHL}
	OPCodesCB[0x4F] = &opCode{"BIT 1 A", 1, bit}
	OPCodesCB[0x50] = &opCode{"BIT 2 B", 1, bit}
	OPCodesCB[0x51] = &opCode{"BIT 2 C", 1, bit}
	OPCodesCB[0x52] = &opCode{"BIT 2 D", 1, bit}
	OPCodesCB[0x53] = &opCode{"BIT 2 E", 1, bit}
	OPCodesCB[0x54] = &opCode{"BIT 2 H", 1, bit}
	OPCodesCB[0x55] = &opCode{"BIT 2 L", 1, bit}
	OPCodesCB[0x56] = &opCode{"BIT 2 (HL)", 1, bitHL}
	OPCodesCB[0x57] = &opCode{"BIT 2 A", 1, bit}
	OPCodesCB[0x58] = &opCode{"BIT 3 B", 1, bit}
	OPCodesCB[0x59] = &opCode{"BIT 3 C", 1, bit}
	OPCodesCB[0x5A] = &opCode{"BIT 3 D", 1, bit}
	OPCodesCB[0x5B] = &opCode{"BIT 3 E", 1, bit}
	OPCodesCB[0x5C] = &opCode{"BIT 3 H", 1, bit}
	OPCodesCB[0x5D] = &opCode{"BIT 3 L", 1, bit}
	OPCodesCB[0x5E] = &opCode{"BIT 3 (HL)", 1, bitHL}
	OPCodesCB[0x5F] = &opCode{"BIT 3 A", 1, bit}
	OPCodesCB[0x60] = &opCode{"BIT 4 B", 1, bit}
	OPCodesCB[0x61] = &opCode{"BIT 4 C", 1, bit}
	OPCodesCB[0x62] = &opCode{"BIT 4 D", 1, bit}
	OPCodesCB[0x63] = &opCode{"BIT 4 E", 1, bit}
	OPCodesCB[0x64] = &opCode{"BIT 4 H", 1, bit}
	OPCodesCB[0x65] = &opCode{"BIT 4 L", 1, bit}
	OPCodesCB[0x66] = &opCode{"BIT 4 (HL)", 1, bitHL}
	OPCodesCB[0x67] = &opCode{"BIT 4 A", 1, bit}
	OPCodesCB[0x68] = &opCode{"BIT 5 B", 1, bit}
	OPCodesCB[0x69] = &opCode{"BIT 5 C", 1, bit}
	OPCodesCB[0x6A] = &opCode{"BIT 5 D", 1, bit}
	OPCodesCB[0x6B] = &opCode{"BIT 5 E", 1, bit}
	OPCodesCB[0x6C] = &opCode{"BIT 5 H", 1, bit}
	OPCodesCB[0x6D] = &opCode{"BIT 5 L", 1, bit}
	OPCodesCB[0x6E] = &opCode{"BIT 5 (HL)", 1, bitHL}
	OPCodesCB[0x6F] = &opCode{"BIT 5 A", 1, bit}
	OPCodesCB[0x70] = &opCode{"BIT 6 B", 1, bit}
	OPCodesCB[0x71] = &opCode{"BIT 6 C", 1, bit}
	OPCodesCB[0x72] = &opCode{"BIT 6 D", 1, bit}
	OPCodesCB[0x73] = &opCode{"BIT 6 E", 1, bit}
	OPCodesCB[0x74] = &opCode{"BIT 6 H", 1, bit}
	OPCodesCB[0x75] = &opCode{"BIT 6 L", 1, bit}
	OPCodesCB[0x76] = &opCode{"BIT 6 (HL)", 1, bitHL}
	OPCodesCB[0x77] = &opCode{"BIT 6 A", 1, bit}
	OPCodesCB[0x78] = &opCode{"BIT 7 B", 1, bit}
	OPCodesCB[0x79] = &opCode{"BIT 7 C", 1, bit}
	OPCodesCB[0x7A] = &opCode{"BIT 7 D", 1, bit}
	OPCodesCB[0x7B] = &opCode{"BIT 7 E", 1, bit}
	OPCodesCB[0x7C] = &opCode{"BIT 7 H", 1, bit}
	OPCodesCB[0x7D] = &opCode{"BIT 7 L", 1, bit}
	OPCodesCB[0x7E] = &opCode{"BIT 7 (HL)", 1, bitHL}
	OPCodesCB[0x7F] = &opCode{"BIT 7 A", 1, bit}
	OPCodesCB[0x80] = &opCode{"RES 0 B", 1, res}
	OPCodesCB[0x81] = &opCode{"RES 0 C", 1, res}
	OPCodesCB[0x82] = &opCode{"RES 0 D", 1, res}
	OPCodesCB[0x83] = &opCode{"RES 0 E", 1, res}
	OPCodesCB[0x84] = &opCode{"RES 0 H", 1, res}
	OPCodesCB[0x85] = &opCode{"RES 0 L", 1, res}
	OPCodesCB[0x86] = &opCode{"RES 0 (HL)", 1, resHL}
	OPCodesCB[0x87] = &opCode{"RES 0 A", 1, res}
	OPCodesCB[0x88] = &opCode{"RES 1 B", 1, res}
	OPCodesCB[0x89] = &opCode{"RES 1 C", 1, res}
	OPCodesCB[0x8A] = &opCode{"RES 1 D", 1, res}
	OPCodesCB[0x8B] = &opCode{"RES 1 E", 1, res}
	OPCodesCB[0x8C] = &opCode{"RES 1 H", 1, res}
	OPCodesCB[0x8D] = &opCode{"RES 1 L", 1, res}
	OPCodesCB[0x8E] = &opCode{"RES 1 (HL)", 1, resHL}
	OPCodesCB[0x8F] = &opCode{"RES 1 A", 1, res}
	OPCodesCB[0x90] = &opCode{"RES 2 B", 1, res}
	OPCodesCB[0x91] = &opCode{"RES 2 C", 1, res}
	OPCodesCB[0x92] = &opCode{"RES 2 D", 1, res}
	OPCodesCB[0x93] = &opCode{"RES 2 E", 1, res}
	OPCodesCB[0x94] = &opCode{"RES 2 H", 1, res}
	OPCodesCB[0x95] = &opCode{"RES 2 L", 1, res}
	OPCodesCB[0x96] = &opCode{"RES 2 (HL)", 1, resHL}
	OPCodesCB[0x97] = &opCode{"RES 2 A", 1, res}
	OPCodesCB[0x98] = &opCode{"RES 3 B", 1, res}
	OPCodesCB[0x99] = &opCode{"RES 3 C", 1, res}
	OPCodesCB[0x9A] = &opCode{"RES 3 D", 1, res}
	OPCodesCB[0x9B] = &opCode{"RES 3 E", 1, res}
	OPCodesCB[0x9C] = &opCode{"RES 3 H", 1, res}
	OPCodesCB[0x9D] = &opCode{"RES 3 L", 1, res}
	OPCodesCB[0x9E] = &opCode{"RES 3 (HL)", 1, resHL}
	OPCodesCB[0x9F] = &opCode{"RES 3 A", 1, res}
	OPCodesCB[0xA0] = &opCode{"RES 4 B", 1, res}
	OPCodesCB[0xA1] = &opCode{"RES 4 C", 1, res}
	OPCodesCB[0xA2] = &opCode{"RES 4 D", 1, res}
	OPCodesCB[0xA3] = &opCode{"RES 4 E", 1, res}
	OPCodesCB[0xA4] = &opCode{"RES 4 H", 1, res}
	OPCodesCB[0xA5] = &opCode{"RES 4 L", 1, res}
	OPCodesCB[0xA6] = &opCode{"RES 4 (HL)", 1, resHL}
	OPCodesCB[0xA7] = &opCode{"RES 4 A", 1, res}
	OPCodesCB[0xA8] = &opCode{"RES 5 B", 1, res}
	OPCodesCB[0xA9] = &opCode{"RES 5 C", 1, res}
	OPCodesCB[0xAA] = &opCode{"RES 5 D", 1, res}
	OPCodesCB[0xAB] = &opCode{"RES 5 E", 1, res}
	OPCodesCB[0xAC] = &opCode{"RES 5 H", 1, res}
	OPCodesCB[0xAD] = &opCode{"RES 5 L", 1, res}
	OPCodesCB[0xAE] = &opCode{"RES 5 (HL)", 1, resHL}
	OPCodesCB[0xAF] = &opCode{"RES 5 A", 1, res}
	OPCodesCB[0xB0] = &opCode{"RES 6 B", 1, res}
	OPCodesCB[0xB1] = &opCode{"RES 6 C", 1, res}
	OPCodesCB[0xB2] = &opCode{"RES 6 D", 1, res}
	OPCodesCB[0xB3] = &opCode{"RES 6 E", 1, res}
	OPCodesCB[0xB4] = &opCode{"RES 6 H", 1, res}
	OPCodesCB[0xB5] = &opCode{"RES 6 L", 1, res}
	OPCodesCB[0xB6] = &opCode{"RES 6 (HL)", 1, resHL}
	OPCodesCB[0xB7] = &opCode{"RES 6 A", 1, res}
	OPCodesCB[0xB8] = &opCode{"RES 7 B", 1, res}
	OPCodesCB[0xB9] = &opCode{"RES 7 C", 1, res}
	OPCodesCB[0xBA] = &opCode{"RES 7 D", 1, res}
	OPCodesCB[0xBB] = &opCode{"RES 7 E", 1, res}
	OPCodesCB[0xBC] = &opCode{"RES 7 H", 1, res}
	OPCodesCB[0xBD] = &opCode{"RES 7 L", 1, res}
	OPCodesCB[0xBE] = &opCode{"RES 7 (HL)", 1, resHL}
	OPCodesCB[0xBF] = &opCode{"RES 7 A", 1, res}
	OPCodesCB[0xC0] = &opCode{"SET 0 B", 1, set}
	OPCodesCB[0xC1] = &opCode{"SET 0 C", 1, set}
	OPCodesCB[0xC2] = &opCode{"SET 0 D", 1, set}
	OPCodesCB[0xC3] = &opCode{"SET 0 E", 1, set}
	OPCodesCB[0xC4] = &opCode{"SET 0 H", 1, set}
	OPCodesCB[0xC5] = &opCode{"SET 0 L", 1, set}
	OPCodesCB[0xC6] = &opCode{"SET 0 (HL)", 1, setHL}
	OPCodesCB[0xC7] = &opCode{"SET 0 A", 1, set}
	OPCodesCB[0xC8] = &opCode{"SET 1 B", 1, set}
	OPCodesCB[0xC9] = &opCode{"SET 1 C", 1, set}
	OPCodesCB[0xCA] = &opCode{"SET 1 D", 1, set}
	OPCodesCB[0xCB] = &opCode{"SET 1 E", 1, set}
	OPCodesCB[0xCC] = &opCode{"SET 1 H", 1, set}
	OPCodesCB[0xCD] = &opCode{"SET 1 L", 1, set}
	OPCodesCB[0xCE] = &opCode{"SET 1 (HL)", 1, setHL}
	OPCodesCB[0xCF] = &opCode{"SET 1 A", 1, set}
	OPCodesCB[0xD0] = &opCode{"SET 2 B", 1, set}
	OPCodesCB[0xD1] = &opCode{"SET 2 C", 1, set}
	OPCodesCB[0xD2] = &opCode{"SET 2 D", 1, set}
	OPCodesCB[0xD3] = &opCode{"SET 2 E", 1, set}
	OPCodesCB[0xD4] = &opCode{"SET 2 H", 1, set}
	OPCodesCB[0xD5] = &opCode{"SET 2 L", 1, set}
	OPCodesCB[0xD6] = &opCode{"SET 2 (HL)", 1, setHL}
	OPCodesCB[0xD7] = &opCode{"SET 2 A", 1, set}
	OPCodesCB[0xD8] = &opCode{"SET 3 B", 1, set}
	OPCodesCB[0xD9] = &opCode{"SET 3 C", 1, set}
	OPCodesCB[0xDA] = &opCode{"SET 3 D", 1, set}
	OPCodesCB[0xDB] = &opCode{"SET 3 E", 1, set}
	OPCodesCB[0xDC] = &opCode{"SET 3 H", 1, set}
	OPCodesCB[0xDD] = &opCode{"SET 3 L", 1, set}
	OPCodesCB[0xDE] = &opCode{"SET 3 (HL)", 1, setHL}
	OPCodesCB[0xDF] = &opCode{"SET 3 A", 1, set}
	OPCodesCB[0xE0] = &opCode{"SET 4 B", 1, set}
	OPCodesCB[0xE1] = &opCode{"SET 4 C", 1, set}
	OPCodesCB[0xE2] = &opCode{"SET 4 D", 1, set}
	OPCodesCB[0xE3] = &opCode{"SET 4 E", 1, set}
	OPCodesCB[0xE4] = &opCode{"SET 4 H", 1, set}
	OPCodesCB[0xE5] = &opCode{"SET 4 L", 1, set}
	OPCodesCB[0xE6] = &opCode{"SET 4 (HL)", 1, setHL}
	OPCodesCB[0xE7] = &opCode{"SET 4 A", 1, set}
	OPCodesCB[0xE8] = &opCode{"SET 5 B", 1, set}
	OPCodesCB[0xE9] = &opCode{"SET 5 C", 1, set}
	OPCodesCB[0xEA] = &opCode{"SET 5 D", 1, set}
	OPCodesCB[0xEB] = &opCode{"SET 5 E", 1, set}
	OPCodesCB[0xEC] = &opCode{"SET 5 H", 1, set}
	OPCodesCB[0xED] = &opCode{"SET 5 L", 1, set}
	OPCodesCB[0xEE] = &opCode{"SET 5 (HL)", 1, setHL}
	OPCodesCB[0xEF] = &opCode{"SET 5 A", 1, set}
	OPCodesCB[0xF0] = &opCode{"SET 6 B", 1, set}
	OPCodesCB[0xF1] = &opCode{"SET 6 C", 1, set}
	OPCodesCB[0xF2] = &opCode{"SET 6 D", 1, set}
	OPCodesCB[0xF3] = &opCode{"SET 6 E", 1, set}
	OPCodesCB[0xF4] = &opCode{"SET 6 H", 1, set}
	OPCodesCB[0xF5] = &opCode{"SET 6 L", 1, set}
	OPCodesCB[0xF6] = &opCode{"SET 6 (HL)", 1, setHL}
	OPCodesCB[0xF7] = &opCode{"SET 6 A", 1, set}
	OPCodesCB[0xF8] = &opCode{"SET 7 B", 1, set}
	OPCodesCB[0xF9] = &opCode{"SET 7 C", 1, set}
	OPCodesCB[0xFA] = &opCode{"SET 7 D", 1, set}
	OPCodesCB[0xFB] = &opCode{"SET 7 E", 1, set}
	OPCodesCB[0xFC] = &opCode{"SET 7 H", 1, set}
	OPCodesCB[0xFD] = &opCode{"SET 7 L", 1, set}
	OPCodesCB[0xFE] = &opCode{"SET 7 (HL)", 1, setHL}
	OPCodesCB[0xFF] = &opCode{"SET 7 A", 1, set}
}

func decodeCB(cpu *lr35902) {
	cpu.scheduler.append(newFetch(OPCodesCB))
}

func (op *opCode) Dump(pc uint16, data []byte) string {
	if op == nil {
		return "not found"
	}
	var sb strings.Builder
	sb.WriteString(toHex16(pc))
	sb.WriteString(":")

	if data[0] != 0xcb {
		for i := byte(0); i < op.Len; i++ {
			sb.WriteString(" ")
			sb.WriteString(toHex8_2(data[i]))
		}
		sb.WriteString("             "[op.Len*3+1:])
		sb.WriteString(" : ")

		ins := op.Ins
		switch op.Len {
		case 2:
			ins = strings.Replace(ins, "n", toHex8(data[1]), 1)
			ins = strings.Replace(ins, "e", toHex16(pc+uint16(int8(data[1])+2)), 1)
		case 3:
			ins = strings.Replace(ins, "n", toHex8(data[2]), 1)
			ins = strings.Replace(ins, "n", toHex8_2(data[1]), 1)
		}
		sb.WriteString(ins)
	} else {
		sb.WriteString(" ")
		sb.WriteString(toHex8_2(data[0]))
		sb.WriteString(" ")
		sb.WriteString(toHex8_2(data[1]))
		sb.WriteString("      ")
		sb.WriteString(" : ")
		sb.WriteString(OPCodesCB[data[1]].Ins)
	}

	return sb.String()
}

func toHex8(v uint8) string {
	n := "0" + strconv.FormatUint(uint64(v), 16)
	return "0x" + n[len(n)-2:]
}

func toHex8_2(v uint8) string {
	n := "0" + strconv.FormatUint(uint64(v), 16)
	return n[len(n)-2:]
}

func toHex16(v uint16) string {
	n := "000" + strconv.FormatUint(uint64(v), 16)
	return "0x" + n[len(n)-4:]
}
