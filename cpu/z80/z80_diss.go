package z80

import (
	"fmt"
	"strconv"
	"strings"
)

var ops []string
var opsDD []string
var opsFD []string
var opsCB []string
var opsDDCB []string
var opsFDCB []string
var opsED []string

func init() {
	ops = make([]string, 0x100)
	opsDD = make([]string, 0x100)
	opsFD = make([]string, 0x100)
	opsCB = make([]string, 0x100)
	opsDDCB = make([]string, 0x100)
	opsFDCB = make([]string, 0x100)
	opsED = make([]string, 0x100)

	ops[0x00] = "NOP"
	ops[0x01] = "LD   BC,$nn"
	ops[0x02] = "LD   (BC),A"
	ops[0x03] = "INC  BC"
	ops[0x04] = "INC  B"
	ops[0x05] = "DEC  B"
	ops[0x06] = "LD   B,$n"
	ops[0x07] = "RLCA"
	ops[0x08] = "EX   AF,AF'"
	ops[0x09] = "ADD  HL,BC"
	ops[0x0A] = "LD   A,(BC)"
	ops[0x0B] = "DEC  BC"
	ops[0x0C] = "INC  C"
	ops[0x0D] = "DEC  C"
	ops[0x0E] = "LD   C,$n"
	ops[0x0F] = "RRCA"
	ops[0x10] = "DJNZ $nn"
	ops[0x11] = "LD   DE,$nn"
	ops[0x12] = "LD   (DE),A"
	ops[0x13] = "INC  DE"
	ops[0x14] = "INC  D"
	ops[0x15] = "DEC  D"
	ops[0x16] = "LD   D,$n"
	ops[0x17] = "RLA"
	ops[0x18] = "JR   $nn"
	ops[0x19] = "ADD  HL,DE"
	ops[0x1A] = "LD   A,(DE)"
	ops[0x1B] = "DEC  DE"
	ops[0x1C] = "INC  E"
	ops[0x1D] = "DEC  E"
	ops[0x1E] = "LD   E,$n"
	ops[0x1F] = "RRA"
	ops[0x20] = "JR   NZ,$nn"
	ops[0x21] = "LD   HL,$nn"
	ops[0x22] = "LD   ($nn),HL"
	ops[0x23] = "INC  HL"
	ops[0x24] = "INC  H"
	ops[0x25] = "DEC  H"
	ops[0x26] = "LD   H,$n"
	ops[0x27] = "DAA"
	ops[0x28] = "JR   Z,$nn"
	ops[0x29] = "ADD  HL,HL"
	ops[0x2A] = "LD   HL,($nn)"
	ops[0x2B] = "DEC  HL"
	ops[0x2C] = "INC  L"
	ops[0x2D] = "DEC  L"
	ops[0x2E] = "LD   L,$n"
	ops[0x2F] = "CPL"
	ops[0x30] = "JR   NC,$nn"
	ops[0x31] = "LD   SP,$nn"
	ops[0x32] = "LD   ($nn),A"
	ops[0x33] = "INC  SP"
	ops[0x34] = "INC  (HL)"
	ops[0x35] = "DEC  (HL)"
	ops[0x36] = "LD   (HL),$n"
	ops[0x37] = "SCF"
	ops[0x38] = "JR   C,$nn"
	ops[0x39] = "ADD  HL,SP"
	ops[0x3A] = "LD   A,($nn)"
	ops[0x3B] = "DEC  SP"
	ops[0x3C] = "INC  A"
	ops[0x3D] = "DEC  A"
	ops[0x3E] = "LD   A,$n"
	ops[0x3F] = "CCF"
	ops[0x40] = "LD   B,B"
	ops[0x41] = "LD   B,C"
	ops[0x42] = "LD   B,D"
	ops[0x43] = "LD   B,E"
	ops[0x44] = "LD   B,H"
	ops[0x45] = "LD   B,L"
	ops[0x46] = "LD   B,(HL)"
	ops[0x47] = "LD   B,A"
	ops[0x48] = "LD   C,B"
	ops[0x49] = "LD   C,C"
	ops[0x4A] = "LD   C,D"
	ops[0x4B] = "LD   C,E"
	ops[0x4C] = "LD   C,H"
	ops[0x4D] = "LD   C,L"
	ops[0x4E] = "LD   C,(HL)"
	ops[0x4F] = "LD   C,A"
	ops[0x50] = "LD   D,B"
	ops[0x51] = "LD   D,C"
	ops[0x52] = "LD   D,D"
	ops[0x53] = "LD   D,E"
	ops[0x54] = "LD   D,H"
	ops[0x55] = "LD   D,L"
	ops[0x56] = "LD   D,(HL)"
	ops[0x57] = "LD   D,A"
	ops[0x58] = "LD   E,B"
	ops[0x59] = "LD   E,C"
	ops[0x5A] = "LD   E,D"
	ops[0x5B] = "LD   E,E"
	ops[0x5C] = "LD   E,H"
	ops[0x5D] = "LD   E,L"
	ops[0x5E] = "LD   E,(HL)"
	ops[0x5F] = "LD   E,A"
	ops[0x60] = "LD   H,B"
	ops[0x61] = "LD   H,C"
	ops[0x62] = "LD   H,D"
	ops[0x63] = "LD   H,E"
	ops[0x64] = "LD   H,H"
	ops[0x65] = "LD   H,L"
	ops[0x66] = "LD   H,(HL)"
	ops[0x67] = "LD   H,A"
	ops[0x68] = "LD   L,B"
	ops[0x69] = "LD   L,C"
	ops[0x6A] = "LD   L,D"
	ops[0x6B] = "LD   L,E"
	ops[0x6C] = "LD   L,H"
	ops[0x6D] = "LD   L,L"
	ops[0x6E] = "LD   L,(HL)"
	ops[0x6F] = "LD   L,A"
	ops[0x70] = "LD   (HL),B"
	ops[0x71] = "LD   (HL),C"
	ops[0x72] = "LD   (HL),D"
	ops[0x73] = "LD   (HL),E"
	ops[0x74] = "LD   (HL),H"
	ops[0x75] = "LD   (HL),L"
	ops[0x76] = "HALT"
	ops[0x77] = "LD   (HL),A"
	ops[0x78] = "LD   A,B"
	ops[0x79] = "LD   A,C"
	ops[0x7A] = "LD   A,D"
	ops[0x7B] = "LD   A,E"
	ops[0x7C] = "LD   A,H"
	ops[0x7D] = "LD   A,L"
	ops[0x7E] = "LD   A,(HL)"
	ops[0x7F] = "LD   A,A"
	ops[0x80] = "ADD  A,B"
	ops[0x81] = "ADD  A,C"
	ops[0x82] = "ADD  A,D"
	ops[0x83] = "ADD  A,E"
	ops[0x84] = "ADD  A,H"
	ops[0x85] = "ADD  A,L"
	ops[0x86] = "ADD  A,(HL)"
	ops[0x87] = "ADD  A,A"
	ops[0x88] = "ADC  A,B"
	ops[0x89] = "ADC  A,C"
	ops[0x8A] = "ADC  A,D"
	ops[0x8B] = "ADC  A,E"
	ops[0x8C] = "ADC  A,H"
	ops[0x8D] = "ADC  A,L"
	ops[0x8E] = "ADC  A,(HL)"
	ops[0x8F] = "ADC  A,A"
	ops[0x90] = "SUB  A,B"
	ops[0x91] = "SUB  A,C"
	ops[0x92] = "SUB  A,D"
	ops[0x93] = "SUB  A,E"
	ops[0x94] = "SUB  A,H"
	ops[0x95] = "SUB  A,L"
	ops[0x96] = "SUB  A,(HL)"
	ops[0x97] = "SUB  A,A"
	ops[0x98] = "SBC  A,B"
	ops[0x99] = "SBC  A,C"
	ops[0x9A] = "SBC  A,D"
	ops[0x9B] = "SBC  A,E"
	ops[0x9C] = "SBC  A,H"
	ops[0x9D] = "SBC  A,L"
	ops[0x9E] = "SBC  A,(HL)"
	ops[0x9F] = "SBC  A,A"
	ops[0xA0] = "AND  B"
	ops[0xA1] = "AND  C"
	ops[0xA2] = "AND  D"
	ops[0xA3] = "AND  E"
	ops[0xA4] = "AND  H"
	ops[0xA5] = "AND  L"
	ops[0xA6] = "AND  (HL)"
	ops[0xA7] = "AND  A"
	ops[0xA8] = "XOR  B"
	ops[0xA9] = "XOR  C"
	ops[0xAA] = "XOR  D"
	ops[0xAB] = "XOR  E"
	ops[0xAC] = "XOR  H"
	ops[0xAD] = "XOR  L"
	ops[0xAE] = "XOR  (HL)"
	ops[0xAF] = "XOR  A"
	ops[0xB0] = "OR   B"
	ops[0xB1] = "OR   C"
	ops[0xB2] = "OR   D"
	ops[0xB3] = "OR   E"
	ops[0xB4] = "OR   H"
	ops[0xB5] = "OR   L"
	ops[0xB6] = "OR   (HL)"
	ops[0xB7] = "OR   A"
	ops[0xB8] = "CP   B"
	ops[0xB9] = "CP   C"
	ops[0xBA] = "CP   D"
	ops[0xBB] = "CP   E"
	ops[0xBC] = "CP   H"
	ops[0xBD] = "CP   L"
	ops[0xBE] = "CP   (HL)"
	ops[0xBF] = "CP   A"
	ops[0xC0] = "RET  NZ"
	ops[0xC1] = "POP  BC"
	ops[0xC2] = "JP   NZ,$nn"
	ops[0xC3] = "JP   $nn"
	ops[0xC4] = "CALL NZ,$nn"
	ops[0xC5] = "PUSH BC"
	ops[0xC6] = "ADD  A,$n"
	ops[0xC7] = "RST  $00"
	ops[0xC8] = "RET  Z"
	ops[0xC9] = "RET"
	ops[0xCA] = "JP   Z,$nn"
	ops[0xCC] = "CALL Z,$nn"
	ops[0xCD] = "CALL $nn"
	ops[0xCE] = "ADC  A,$n"
	ops[0xCF] = "RST  $08"
	ops[0xD0] = "RET  NC"
	ops[0xD1] = "POP  DE"
	ops[0xD2] = "JP   NC,$nn"
	ops[0xD3] = "OUT  ($n),A"
	ops[0xD4] = "CALL NC,$nn"
	ops[0xD5] = "PUSH DE"
	ops[0xD6] = "SUB  A,$n"
	ops[0xD7] = "RST  $10"
	ops[0xD8] = "RET  C"
	ops[0xD9] = "EXX"
	ops[0xDA] = "JP   C,$nn"
	ops[0xDB] = "IN   A,($n)"
	ops[0xDC] = "CALL C,$nn"
	ops[0xDE] = "SBC  A,$n"
	ops[0xDF] = "RST  $18"
	ops[0xE0] = "RET  PO"
	ops[0xE1] = "POP  HL"
	ops[0xE2] = "JP   PO,$nn"
	ops[0xE3] = "EX   (SP),HL"
	ops[0xE4] = "CALL PO,$nn"
	ops[0xE5] = "PUSH HL"
	ops[0xE6] = "AND  $n"
	ops[0xE7] = "RST  $20"
	ops[0xE8] = "RET  PE"
	ops[0xE9] = "JP   (HL)"
	ops[0xEA] = "JP   PE,$nn"
	ops[0xEB] = "EX   DE,HL"
	ops[0xEC] = "CALL PE,$nn"
	ops[0xEE] = "XOR  $n"
	ops[0xEF] = "RST  $28"
	ops[0xF0] = "RET  P"
	ops[0xF1] = "POP  AF"
	ops[0xF2] = "JP   P,$nn"
	ops[0xF3] = "DI"
	ops[0xF4] = "CALL P,$nn"
	ops[0xF5] = "PUSH AF"
	ops[0xF6] = "OR   $n"
	ops[0xF7] = "RST  $30"
	ops[0xF8] = "RET  M"
	ops[0xF9] = "LD   SP,HL"
	ops[0xFA] = "JP   M,$nn"
	ops[0xFB] = "EI"
	ops[0xFC] = "CALL M,$nn"
	ops[0xFE] = "CP   $n"
	ops[0xFF] = "RST  $38"

	opsDD[0x09] = "ADD  IX,BC"
	opsDD[0x19] = "ADD  IX,DE"
	opsDD[0x21] = "LD   IX,$nn"
	opsDD[0x22] = "LD  ($nn),IX"
	opsDD[0x23] = "INC  IX"
	opsDD[0x24] = "INC  IXH"
	opsDD[0x25] = "DEC  IXH"
	opsDD[0x26] = "LD   IXH,$n"
	opsDD[0x29] = "ADD  IX,IX"
	opsDD[0x2A] = "LD  IX,($nn)"
	opsDD[0x2B] = "DEC  IX"
	opsDD[0x2C] = "INC  IXL"
	opsDD[0x2D] = "DEC  IXL"
	opsDD[0x2E] = "LD   IXL,$n"
	opsDD[0x34] = "INC  (IX+0)"
	opsDD[0x35] = "DEC  (IX+0)"
	opsDD[0x36] = "LD   (IX+0),$n2"
	opsDD[0x39] = "ADD  IX,SP"
	opsDD[0x44] = "LD   B,IXH"
	opsDD[0x45] = "LD   B,IXL"
	opsDD[0x46] = "LD   B,(IX+0)"
	opsDD[0x4C] = "LD   C,IXH"
	opsDD[0x4D] = "LD   C,IXL"
	opsDD[0x4E] = "LD   C,(IX+0)"
	opsDD[0x54] = "LD   D,IXH"
	opsDD[0x55] = "LD   D,IXL"
	opsDD[0x56] = "LD   D,(IX+0)"
	opsDD[0x5C] = "LD   E,IXH"
	opsDD[0x5D] = "LD   E,IXL"
	opsDD[0x5E] = "LD   E,(IX+0)"
	opsDD[0x60] = "LD   IXH,B"
	opsDD[0x61] = "LD   IXH,C"
	opsDD[0x62] = "LD   IXH,D"
	opsDD[0x63] = "LD   IXH,E"
	opsDD[0x64] = "LD   IXH,IXH"
	opsDD[0x65] = "LD   IXH,IXL"
	opsDD[0x66] = "LD   H,(IX+0)"
	opsDD[0x67] = "LD   IXH,A"
	opsDD[0x68] = "LD   IXL,B"
	opsDD[0x69] = "LD   IXL,C"
	opsDD[0x6A] = "LD   IXL,D"
	opsDD[0x6B] = "LD   IXL,E"
	opsDD[0x6C] = "LD   IXL,IXH"
	opsDD[0x6D] = "LD   IXL,IXL"
	opsDD[0x6E] = "LD   L,(IX+0)"
	opsDD[0x6F] = "LD   IXL,A"
	opsDD[0x70] = "LD   (IX+0),B"
	opsDD[0x71] = "LD   (IX+0),C"
	opsDD[0x72] = "LD   (IX+0),D"
	opsDD[0x73] = "LD   (IX+0),E"
	opsDD[0x74] = "LD   (IX+0),H"
	opsDD[0x75] = "LD   (IX+0),L"
	opsDD[0x77] = "LD   (IX+0),A"
	opsDD[0x7C] = "LD   A,IXH"
	opsDD[0x7D] = "LD   A,IXL"
	opsDD[0x7E] = "LD   A,(IX+0)"
	opsDD[0x84] = "ADD  A,IXH"
	opsDD[0x85] = "ADD  A,IXL"
	opsDD[0x86] = "ADD  A,(IX+0)"
	opsDD[0x8C] = "ADC  A,IXH"
	opsDD[0x8D] = "ADC  A,IXL"
	opsDD[0x8E] = "ADC  A,(IX+0)"
	opsDD[0x94] = "SUB  A,IXH"
	opsDD[0x95] = "SUB  A,IXL"
	opsDD[0x96] = "SUB  A,(IX+0)"
	opsDD[0x9C] = "SBC  A,IXH"
	opsDD[0x9D] = "SBC  A,IXL"
	opsDD[0x9E] = "SBC  A,(IX+0)"
	opsDD[0xA4] = "AND  IXH"
	opsDD[0xA5] = "AND  IXL"
	opsDD[0xA6] = "AND  (IX+0)"
	opsDD[0xAC] = "XOR  IXH"
	opsDD[0xAD] = "XOR  IXL"
	opsDD[0xAE] = "XOR  (IX+0)"
	opsDD[0xB4] = "OR   IXH"
	opsDD[0xB5] = "OR   IXL"
	opsDD[0xB6] = "OR   (IX+0)"
	opsDD[0xBC] = "CP   IXH"
	opsDD[0xBD] = "CP   IXL"
	opsDD[0xBE] = "CP   (IX+0)"
	opsDD[0xE1] = "POP  IX"
	opsDD[0xE3] = "EX   (SP),IX"
	opsDD[0xE5] = "PUSH IX"
	opsDD[0xE9] = "JP   (IX)"

	opsFD[0x09] = "ADD  IY,BC"
	opsFD[0x19] = "ADD  IY,DE"
	opsFD[0x21] = "LD   IY,$nn"
	opsFD[0x22] = "LD  ($nn),IY"
	opsFD[0x23] = "INC  IY"
	opsFD[0x24] = "INC  IYH"
	opsFD[0x25] = "DEC  IYH"
	opsFD[0x26] = "LD   IYH,$n"
	opsFD[0x29] = "ADD  IY,IY"
	opsFD[0x2A] = "LD  IY,($nn)"
	opsFD[0x2B] = "DEC  IY"
	opsFD[0x2C] = "INC  IYL"
	opsFD[0x2D] = "DEC  IYL"
	opsFD[0x2E] = "LD   IYL,$n"
	opsFD[0x34] = "INC  (IY+0)"
	opsFD[0x35] = "DEC  (IY+0)"
	opsFD[0x36] = "LD   (IY+0),$n2"
	opsFD[0x39] = "ADD  IY,SP"
	opsFD[0x44] = "LD   B,IYH"
	opsFD[0x45] = "LD   B,IYL"
	opsFD[0x46] = "LD   B,(IY+0)"
	opsFD[0x4C] = "LD   C,IYH"
	opsFD[0x4D] = "LD   C,IYL"
	opsFD[0x4E] = "LD   C,(IY+0)"
	opsFD[0x54] = "LD   D,IYH"
	opsFD[0x55] = "LD   D,IYL"
	opsFD[0x56] = "LD   D,(IY+0)"
	opsFD[0x5C] = "LD   E,IYH"
	opsFD[0x5D] = "LD   E,IYL"
	opsFD[0x5E] = "LD   E,(IY+0)"
	opsFD[0x60] = "LD   IYH,B"
	opsFD[0x61] = "LD   IYH,C"
	opsFD[0x62] = "LD   IYH,D"
	opsFD[0x63] = "LD   IYH,E"
	opsFD[0x64] = "LD   IYH,IYH"
	opsFD[0x65] = "LD   IYH,IYL"
	opsFD[0x66] = "LD   H,(IY+0)"
	opsFD[0x67] = "LD   IYH,A"
	opsFD[0x68] = "LD   IYL,B"
	opsFD[0x69] = "LD   IYL,C"
	opsFD[0x6A] = "LD   IYL,D"
	opsFD[0x6B] = "LD   IYL,E"
	opsFD[0x6C] = "LD   IYL,IYH"
	opsFD[0x6D] = "LD   IYL,IYL"
	opsFD[0x6E] = "LD   L,(IY+0)"
	opsFD[0x6F] = "LD   IYL,A"
	opsFD[0x70] = "LD   (IY+0),B"
	opsFD[0x71] = "LD   (IY+0),C"
	opsFD[0x72] = "LD   (IY+0),D"
	opsFD[0x73] = "LD   (IY+0),E"
	opsFD[0x74] = "LD   (IY+0),H"
	opsFD[0x75] = "LD   (IY+0),L"
	opsFD[0x77] = "LD   (IY+0),A"
	opsFD[0x7C] = "LD   A,IYH"
	opsFD[0x7D] = "LD   A,IYL"
	opsFD[0x7E] = "LD   A,(IY+0)"
	opsFD[0x84] = "ADD  A,IYH"
	opsFD[0x85] = "ADD  A,IYL"
	opsFD[0x86] = "ADD  A,(IY+0)"
	opsFD[0x8C] = "ADC  A,IYH"
	opsFD[0x8D] = "ADC  A,IYL"
	opsFD[0x8E] = "ADC  A,(IY+0)"
	opsFD[0x94] = "SUB  A,IYH"
	opsFD[0x95] = "SUB  A,IYL"
	opsFD[0x96] = "SUB  A,(IY+0)"
	opsFD[0x9C] = "SBC  A,IYH"
	opsFD[0x9D] = "SBC  A,IYL"
	opsFD[0x9E] = "SBC  A,(IY+0)"
	opsFD[0xA4] = "AND  IYH"
	opsFD[0xA5] = "AND  IYL"
	opsFD[0xA6] = "AND  (IY+0)"
	opsFD[0xAC] = "XOR  IYH"
	opsFD[0xAD] = "XOR  IYL"
	opsFD[0xAE] = "XOR  (IY+0)"
	opsFD[0xB4] = "OR   IYH"
	opsFD[0xB5] = "OR   IYL"
	opsFD[0xB6] = "OR   (IY+0)"
	opsFD[0xBC] = "CP   IYH"
	opsFD[0xBD] = "CP   IYL"
	opsFD[0xBE] = "CP   (IY+0)"
	opsFD[0xE1] = "POP  IY"
	opsFD[0xE3] = "EX   (SP),IY"
	opsFD[0xE5] = "PUSH IY"
	opsFD[0xE9] = "JP   (IY)"

	opsCB[0x00] = "RLC  B"
	opsCB[0x01] = "RLC  C"
	opsCB[0x02] = "RLC  D"
	opsCB[0x03] = "RLC  E"
	opsCB[0x04] = "RLC  H"
	opsCB[0x05] = "RLC  L"
	opsCB[0x06] = "RLC  (HL)"
	opsCB[0x07] = "RLC  A"
	opsCB[0x08] = "RRC  B"
	opsCB[0x09] = "RRC  C"
	opsCB[0x0A] = "RRC  D"
	opsCB[0x0B] = "RRC  E"
	opsCB[0x0C] = "RRC  H"
	opsCB[0x0D] = "RRC  L"
	opsCB[0x0E] = "RRC  (HL)"
	opsCB[0x0F] = "RRC  A"
	opsCB[0x10] = "RL   B"
	opsCB[0x11] = "RL   C"
	opsCB[0x12] = "RL   D"
	opsCB[0x13] = "RL   E"
	opsCB[0x14] = "RL   H"
	opsCB[0x15] = "RL   L"
	opsCB[0x16] = "RL   (HL)"
	opsCB[0x17] = "RL   A"
	opsCB[0x18] = "RR   B"
	opsCB[0x19] = "RR   C"
	opsCB[0x1A] = "RR   D"
	opsCB[0x1B] = "RR   E"
	opsCB[0x1C] = "RR   H"
	opsCB[0x1D] = "RR   L"
	opsCB[0x1E] = "RR   (HL)"
	opsCB[0x1F] = "RR   A"
	opsCB[0x20] = "SLA  B"
	opsCB[0x21] = "SLA  C"
	opsCB[0x22] = "SLA  D"
	opsCB[0x23] = "SLA  E"
	opsCB[0x24] = "SLA  H"
	opsCB[0x25] = "SLA  L"
	opsCB[0x26] = "SLA  (HL)"
	opsCB[0x27] = "SLA  A"
	opsCB[0x28] = "SRA  B"
	opsCB[0x29] = "SRA  C"
	opsCB[0x2A] = "SRA  D"
	opsCB[0x2B] = "SRA  E"
	opsCB[0x2C] = "SRA  H"
	opsCB[0x2D] = "SRA  L"
	opsCB[0x2E] = "SRA  (HL)"
	opsCB[0x2F] = "SRA  A"
	opsCB[0x30] = "SLS  B"
	opsCB[0x31] = "SLS  C"
	opsCB[0x32] = "SLS  D"
	opsCB[0x33] = "SLS  E"
	opsCB[0x34] = "SLS  H"
	opsCB[0x35] = "SLS  L"
	opsCB[0x36] = "SLS  (HL)"
	opsCB[0x37] = "SLS  A"
	opsCB[0x38] = "SRL  B"
	opsCB[0x39] = "SRL  C"
	opsCB[0x3A] = "SRL  D"
	opsCB[0x3B] = "SRL  E"
	opsCB[0x3C] = "SRL  H"
	opsCB[0x3D] = "SRL  L"
	opsCB[0x3E] = "SRL  (HL)"
	opsCB[0x3F] = "SRL  A"
	opsCB[0x40] = "BIT  0,B"
	opsCB[0x41] = "BIT  0,C"
	opsCB[0x42] = "BIT  0,D"
	opsCB[0x43] = "BIT  0,E"
	opsCB[0x44] = "BIT  0,H"
	opsCB[0x45] = "BIT  0,L"
	opsCB[0x46] = "BIT  0,(HL)"
	opsCB[0x47] = "BIT  0,A"
	opsCB[0x48] = "BIT  1,B"
	opsCB[0x49] = "BIT  1,C"
	opsCB[0x4A] = "BIT  1,D"
	opsCB[0x4B] = "BIT  1,E"
	opsCB[0x4C] = "BIT  1,H"
	opsCB[0x4D] = "BIT  1,L"
	opsCB[0x4E] = "BIT  1,(HL)"
	opsCB[0x4F] = "BIT  1,A"
	opsCB[0x50] = "BIT  2,B"
	opsCB[0x51] = "BIT  2,C"
	opsCB[0x52] = "BIT  2,D"
	opsCB[0x53] = "BIT  2,E"
	opsCB[0x54] = "BIT  2,H"
	opsCB[0x55] = "BIT  2,L"
	opsCB[0x56] = "BIT  2,(HL)"
	opsCB[0x57] = "BIT  2,A"
	opsCB[0x58] = "BIT  3,B"
	opsCB[0x59] = "BIT  3,C"
	opsCB[0x5A] = "BIT  3,D"
	opsCB[0x5B] = "BIT  3,E"
	opsCB[0x5C] = "BIT  3,H"
	opsCB[0x5D] = "BIT  3,L"
	opsCB[0x5E] = "BIT  3,(HL)"
	opsCB[0x5F] = "BIT  3,A"
	opsCB[0x60] = "BIT  4,B"
	opsCB[0x61] = "BIT  4,C"
	opsCB[0x62] = "BIT  4,D"
	opsCB[0x63] = "BIT  4,E"
	opsCB[0x64] = "BIT  4,H"
	opsCB[0x65] = "BIT  4,L"
	opsCB[0x66] = "BIT  4,(HL)"
	opsCB[0x67] = "BIT  4,A"
	opsCB[0x68] = "BIT  5,B"
	opsCB[0x69] = "BIT  5,C"
	opsCB[0x6A] = "BIT  5,D"
	opsCB[0x6B] = "BIT  5,E"
	opsCB[0x6C] = "BIT  5,H"
	opsCB[0x6D] = "BIT  5,L"
	opsCB[0x6E] = "BIT  5,(HL)"
	opsCB[0x6F] = "BIT  5,A"
	opsCB[0x70] = "BIT  6,B"
	opsCB[0x71] = "BIT  6,C"
	opsCB[0x72] = "BIT  6,D"
	opsCB[0x73] = "BIT  6,E"
	opsCB[0x74] = "BIT  6,H"
	opsCB[0x75] = "BIT  6,L"
	opsCB[0x76] = "BIT  6,(HL)"
	opsCB[0x77] = "BIT  6,A"
	opsCB[0x78] = "BIT  7,B"
	opsCB[0x79] = "BIT  7,C"
	opsCB[0x7A] = "BIT  7,D"
	opsCB[0x7B] = "BIT  7,E"
	opsCB[0x7C] = "BIT  7,H"
	opsCB[0x7D] = "BIT  7,L"
	opsCB[0x7E] = "BIT  7,(HL)"
	opsCB[0x7F] = "BIT  7,A"
	opsCB[0x80] = "RES  0,B"
	opsCB[0x81] = "RES  0,C"
	opsCB[0x82] = "RES  0,D"
	opsCB[0x83] = "RES  0,E"
	opsCB[0x84] = "RES  0,H"
	opsCB[0x85] = "RES  0,L"
	opsCB[0x86] = "RES  0,(HL)"
	opsCB[0x87] = "RES  0,A"
	opsCB[0x88] = "RES  1,B"
	opsCB[0x89] = "RES  1,C"
	opsCB[0x8A] = "RES  1,D"
	opsCB[0x8B] = "RES  1,E"
	opsCB[0x8C] = "RES  1,H"
	opsCB[0x8D] = "RES  1,L"
	opsCB[0x8E] = "RES  1,(HL)"
	opsCB[0x8F] = "RES  1,A"
	opsCB[0x90] = "RES  2,B"
	opsCB[0x91] = "RES  2,C"
	opsCB[0x92] = "RES  2,D"
	opsCB[0x93] = "RES  2,E"
	opsCB[0x94] = "RES  2,H"
	opsCB[0x95] = "RES  2,L"
	opsCB[0x96] = "RES  2,(HL)"
	opsCB[0x97] = "RES  2,A"
	opsCB[0x98] = "RES  3,B"
	opsCB[0x99] = "RES  3,C"
	opsCB[0x9A] = "RES  3,D"
	opsCB[0x9B] = "RES  3,E"
	opsCB[0x9C] = "RES  3,H"
	opsCB[0x9D] = "RES  3,L"
	opsCB[0x9E] = "RES  3,(HL)"
	opsCB[0x9F] = "RES  3,A"
	opsCB[0xA0] = "RES  4,B"
	opsCB[0xA1] = "RES  4,C"
	opsCB[0xA2] = "RES  4,D"
	opsCB[0xA3] = "RES  4,E"
	opsCB[0xA4] = "RES  4,H"
	opsCB[0xA5] = "RES  4,L"
	opsCB[0xA6] = "RES  4,(HL)"
	opsCB[0xA7] = "RES  4,A"
	opsCB[0xA8] = "RES  5,B"
	opsCB[0xA9] = "RES  5,C"
	opsCB[0xAA] = "RES  5,D"
	opsCB[0xAB] = "RES  5,E"
	opsCB[0xAC] = "RES  5,H"
	opsCB[0xAD] = "RES  5,L"
	opsCB[0xAE] = "RES  5,(HL)"
	opsCB[0xAF] = "RES  5,A"
	opsCB[0xB0] = "RES  6,B"
	opsCB[0xB1] = "RES  6,C"
	opsCB[0xB2] = "RES  6,D"
	opsCB[0xB3] = "RES  6,E"
	opsCB[0xB4] = "RES  6,H"
	opsCB[0xB5] = "RES  6,L"
	opsCB[0xB6] = "RES  6,(HL)"
	opsCB[0xB7] = "RES  6,A"
	opsCB[0xB8] = "RES  7,B"
	opsCB[0xB9] = "RES  7,C"
	opsCB[0xBA] = "RES  7,D"
	opsCB[0xBB] = "RES  7,E"
	opsCB[0xBC] = "RES  7,H"
	opsCB[0xBD] = "RES  7,L"
	opsCB[0xBE] = "RES  7,(HL)"
	opsCB[0xBF] = "RES  7,A"
	opsCB[0xC0] = "SET  0,B"
	opsCB[0xC1] = "SET  0,C"
	opsCB[0xC2] = "SET  0,D"
	opsCB[0xC3] = "SET  0,E"
	opsCB[0xC4] = "SET  0,H"
	opsCB[0xC5] = "SET  0,L"
	opsCB[0xC6] = "SET  0,(HL)"
	opsCB[0xC7] = "SET  0,A"
	opsCB[0xC8] = "SET  1,B"
	opsCB[0xC9] = "SET  1,C"
	opsCB[0xCA] = "SET  1,D"
	opsCB[0xCB] = "SET  1,E"
	opsCB[0xCC] = "SET  1,H"
	opsCB[0xCD] = "SET  1,L"
	opsCB[0xCE] = "SET  1,(HL)"
	opsCB[0xCF] = "SET  1,A"
	opsCB[0xD0] = "SET  2,B"
	opsCB[0xD1] = "SET  2,C"
	opsCB[0xD2] = "SET  2,D"
	opsCB[0xD3] = "SET  2,E"
	opsCB[0xD4] = "SET  2,H"
	opsCB[0xD5] = "SET  2,L"
	opsCB[0xD6] = "SET  2,(HL)"
	opsCB[0xD7] = "SET  2,A"
	opsCB[0xD8] = "SET  3,B"
	opsCB[0xD9] = "SET  3,C"
	opsCB[0xDA] = "SET  3,D"
	opsCB[0xDB] = "SET  3,E"
	opsCB[0xDC] = "SET  3,H"
	opsCB[0xDD] = "SET  3,L"
	opsCB[0xDE] = "SET  3,(HL)"
	opsCB[0xDF] = "SET  3,A"
	opsCB[0xE0] = "SET  4,B"
	opsCB[0xE1] = "SET  4,C"
	opsCB[0xE2] = "SET  4,D"
	opsCB[0xE3] = "SET  4,E"
	opsCB[0xE4] = "SET  4,H"
	opsCB[0xE5] = "SET  4,L"
	opsCB[0xE6] = "SET  4,(HL)"
	opsCB[0xE7] = "SET  4,A"
	opsCB[0xE8] = "SET  5,B"
	opsCB[0xE9] = "SET  5,C"
	opsCB[0xEA] = "SET  5,D"
	opsCB[0xEB] = "SET  5,E"
	opsCB[0xEC] = "SET  5,H"
	opsCB[0xED] = "SET  5,L"
	opsCB[0xEE] = "SET  5,(HL)"
	opsCB[0xEF] = "SET  5,A"
	opsCB[0xF0] = "SET  6,B"
	opsCB[0xF1] = "SET  6,C"
	opsCB[0xF2] = "SET  6,D"
	opsCB[0xF3] = "SET  6,E"
	opsCB[0xF4] = "SET  6,H"
	opsCB[0xF5] = "SET  6,L"
	opsCB[0xF6] = "SET  6,(HL)"
	opsCB[0xF7] = "SET  6,A"
	opsCB[0xF8] = "SET  7,B"
	opsCB[0xF9] = "SET  7,C"
	opsCB[0xFA] = "SET  7,D"
	opsCB[0xFB] = "SET  7,E"
	opsCB[0xFC] = "SET  7,H"
	opsCB[0xFD] = "SET  7,L"
	opsCB[0xFE] = "SET  7,(HL)"
	opsCB[0xFF] = "SET  7,A"

	opsFDCB[0x00] = "rlc (iy+0)->b"
	opsFDCB[0x01] = "rlc (iy+0)->c"
	opsFDCB[0x02] = "rlc (iy+0)->d"
	opsFDCB[0x03] = "rlc (iy+0)->e"
	opsFDCB[0x04] = "rlc (iy+0)->h"
	opsFDCB[0x05] = "rlc (iy+0)->l"
	opsFDCB[0x06] = "RLC  (IY+0)"
	opsFDCB[0x07] = "rlc (iy+0)->a"
	opsFDCB[0x08] = "rrc (iy+0)->b"
	opsFDCB[0x09] = "rrc (iy+0)->c"
	opsFDCB[0x0A] = "rrc (iy+0)->d"
	opsFDCB[0x0B] = "rrc (iy+0)->e"
	opsFDCB[0x0C] = "rrc (iy+0)->h"
	opsFDCB[0x0D] = "rrc (iy+0)->l"
	opsFDCB[0x0E] = "RRC  (IY+0)"
	opsFDCB[0x0F] = "rrc (iy+0)->a"
	opsFDCB[0x10] = "rl  (iy+0)->b"
	opsFDCB[0x11] = "rl  (iy+0)->c"
	opsFDCB[0x12] = "rl  (iy+0)->d"
	opsFDCB[0x13] = "rl  (iy+0)->e"
	opsFDCB[0x14] = "rl  (iy+0)->h"
	opsFDCB[0x15] = "rl  (iy+0)->l"
	opsFDCB[0x16] = "RL   (IY+0)"
	opsFDCB[0x17] = "rl  (iy+0)->a"
	opsFDCB[0x18] = "rr  (iy+0)->b"
	opsFDCB[0x19] = "rr  (iy+0)->c"
	opsFDCB[0x1A] = "rr  (iy+0)->d"
	opsFDCB[0x1B] = "rr  (iy+0)->e"
	opsFDCB[0x1C] = "rr  (iy+0)->h"
	opsFDCB[0x1D] = "rr  (iy+0)->l"
	opsFDCB[0x1E] = "RR   (IY+0)"
	opsFDCB[0x1F] = "rr  (iy+0)->a"
	opsFDCB[0x20] = "sla (iy+0)->b"
	opsFDCB[0x21] = "sla (iy+0)->c"
	opsFDCB[0x22] = "sla (iy+0)->d"
	opsFDCB[0x23] = "sla (iy+0)->e"
	opsFDCB[0x24] = "sla (iy+0)->h"
	opsFDCB[0x25] = "sla (iy+0)->l"
	opsFDCB[0x26] = "SLA  (IY+0)"
	opsFDCB[0x27] = "sla (iy+0)->a"
	opsFDCB[0x28] = "sra (iy+0)->b"
	opsFDCB[0x29] = "sra (iy+0)->c"
	opsFDCB[0x2A] = "sra (iy+0)->d"
	opsFDCB[0x2B] = "sra (iy+0)->e"
	opsFDCB[0x2C] = "sra (iy+0)->h"
	opsFDCB[0x2D] = "sra (iy+0)->l"
	opsFDCB[0x2E] = "SRA  (IY+0)"
	opsFDCB[0x2F] = "sra (iy+0)->a"
	opsFDCB[0x30] = "sls (iy+0)->b"
	opsFDCB[0x31] = "sls (iy+0)->c"
	opsFDCB[0x32] = "sls (iy+0)->d"
	opsFDCB[0x33] = "sls (iy+0)->e"
	opsFDCB[0x34] = "sls (iy+0)->h"
	opsFDCB[0x35] = "sls (iy+0)->l"
	opsFDCB[0x36] = "SLS  (IY+0)"
	opsFDCB[0x37] = "sls (iy+0)->a"
	opsFDCB[0x38] = "srl (iy+0)->b"
	opsFDCB[0x39] = "srl (iy+0)->c"
	opsFDCB[0x3A] = "srl (iy+0)->d"
	opsFDCB[0x3B] = "srl (iy+0)->e"
	opsFDCB[0x3C] = "srl (iy+0)->h"
	opsFDCB[0x3D] = "srl (iy+0)->l"
	opsFDCB[0x3E] = "SRL  (IY+0)"
	opsFDCB[0x3F] = "srl (iy+0)->a"
	opsFDCB[0x40] = "bit 0,(iy+0)->b"
	opsFDCB[0x41] = "bit 0,(iy+0)->c"
	opsFDCB[0x42] = "bit 0,(iy+0)->d"
	opsFDCB[0x43] = "bit 0,(iy+0)->e"
	opsFDCB[0x44] = "bit 0,(iy+0)->h"
	opsFDCB[0x45] = "bit 0,(iy+0)->l"
	opsFDCB[0x46] = "BIT  0,(IY+0)"
	opsFDCB[0x47] = "bit 0,(iy+0)->a"
	opsFDCB[0x48] = "bit 1,(iy+0)->b"
	opsFDCB[0x49] = "bit 1,(iy+0)->c"
	opsFDCB[0x4A] = "bit 1,(iy+0)->d"
	opsFDCB[0x4B] = "bit 1,(iy+0)->e"
	opsFDCB[0x4C] = "bit 1,(iy+0)->h"
	opsFDCB[0x4D] = "bit 1,(iy+0)->l"
	opsFDCB[0x4E] = "BIT  1,(IY+0)"
	opsFDCB[0x4F] = "bit 1,(iy+0)->a"
	opsFDCB[0x50] = "bit 2,(iy+0)->b"
	opsFDCB[0x51] = "bit 2,(iy+0)->c"
	opsFDCB[0x52] = "bit 2,(iy+0)->d"
	opsFDCB[0x53] = "bit 2,(iy+0)->e"
	opsFDCB[0x54] = "bit 2,(iy+0)->h"
	opsFDCB[0x55] = "bit 2,(iy+0)->l"
	opsFDCB[0x56] = "BIT  2,(IY+0)"
	opsFDCB[0x57] = "bit 2,(iy+0)->a"
	opsFDCB[0x58] = "bit 3,(iy+0)->b"
	opsFDCB[0x59] = "bit 3,(iy+0)->c"
	opsFDCB[0x5A] = "bit 3,(iy+0)->d"
	opsFDCB[0x5B] = "bit 3,(iy+0)->e"
	opsFDCB[0x5C] = "bit 3,(iy+0)->h"
	opsFDCB[0x5D] = "bit 3,(iy+0)->l"
	opsFDCB[0x5E] = "BIT  3,(IY+0)"
	opsFDCB[0x5F] = "bit 3,(iy+0)->a"
	opsFDCB[0x60] = "bit 4,(iy+0)->b"
	opsFDCB[0x61] = "bit 4,(iy+0)->c"
	opsFDCB[0x62] = "bit 4,(iy+0)->d"
	opsFDCB[0x63] = "bit 4,(iy+0)->e"
	opsFDCB[0x64] = "bit 4,(iy+0)->h"
	opsFDCB[0x65] = "bit 4,(iy+0)->l"
	opsFDCB[0x66] = "BIT  4,(IY+0)"
	opsFDCB[0x67] = "bit 4,(iy+0)->a"
	opsFDCB[0x68] = "bit 5,(iy+0)->b"
	opsFDCB[0x69] = "bit 5,(iy+0)->c"
	opsFDCB[0x6A] = "bit 5,(iy+0)->d"
	opsFDCB[0x6B] = "bit 5,(iy+0)->e"
	opsFDCB[0x6C] = "bit 5,(iy+0)->h"
	opsFDCB[0x6D] = "bit 5,(iy+0)->l"
	opsFDCB[0x6E] = "BIT  5,(IY+0)"
	opsFDCB[0x6F] = "bit 5,(iy+0)->a"
	opsFDCB[0x70] = "bit 6,(iy+0)->b"
	opsFDCB[0x71] = "bit 6,(iy+0)->c"
	opsFDCB[0x72] = "bit 6,(iy+0)->d"
	opsFDCB[0x73] = "bit 6,(iy+0)->e"
	opsFDCB[0x74] = "bit 6,(iy+0)->h"
	opsFDCB[0x75] = "bit 6,(iy+0)->l"
	opsFDCB[0x76] = "BIT  6,(IY+0)"
	opsFDCB[0x77] = "bit 6,(iy+0)->a"
	opsFDCB[0x78] = "bit 7,(iy+0)->b"
	opsFDCB[0x79] = "bit 7,(iy+0)->c"
	opsFDCB[0x7A] = "bit 7,(iy+0)->d"
	opsFDCB[0x7B] = "bit 7,(iy+0)->e"
	opsFDCB[0x7C] = "bit 7,(iy+0)->h"
	opsFDCB[0x7D] = "bit 7,(iy+0)->l"
	opsFDCB[0x7E] = "BIT  7,(IY+0)"
	opsFDCB[0x7F] = "bit 7,(iy+0)->a"
	opsFDCB[0x80] = "res 0,(iy+0)->b"
	opsFDCB[0x81] = "res 0,(iy+0)->c"
	opsFDCB[0x82] = "res 0,(iy+0)->d"
	opsFDCB[0x83] = "res 0,(iy+0)->e"
	opsFDCB[0x84] = "res 0,(iy+0)->h"
	opsFDCB[0x85] = "res 0,(iy+0)->l"
	opsFDCB[0x86] = "RES  0,(IY+0)"
	opsFDCB[0x87] = "res 0,(iy+0)->a"
	opsFDCB[0x88] = "res 1,(iy+0)->b"
	opsFDCB[0x89] = "res 1,(iy+0)->c"
	opsFDCB[0x8A] = "res 1,(iy+0)->d"
	opsFDCB[0x8B] = "res 1,(iy+0)->e"
	opsFDCB[0x8C] = "res 1,(iy+0)->h"
	opsFDCB[0x8D] = "res 1,(iy+0)->l"
	opsFDCB[0x8E] = "RES  1,(IY+0)"
	opsFDCB[0x8F] = "res 1,(iy+0)->a"
	opsFDCB[0x90] = "res 2,(iy+0)->b"
	opsFDCB[0x91] = "res 2,(iy+0)->c"
	opsFDCB[0x92] = "res 2,(iy+0)->d"
	opsFDCB[0x93] = "res 2,(iy+0)->e"
	opsFDCB[0x94] = "res 2,(iy+0)->h"
	opsFDCB[0x95] = "res 2,(iy+0)->l"
	opsFDCB[0x96] = "RES  2,(IY+0)"
	opsFDCB[0x97] = "res 2,(iy+0)->a"
	opsFDCB[0x98] = "res 3,(iy+0)->b"
	opsFDCB[0x99] = "res 3,(iy+0)->c"
	opsFDCB[0x9A] = "res 3,(iy+0)->d"
	opsFDCB[0x9B] = "res 3,(iy+0)->e"
	opsFDCB[0x9C] = "res 3,(iy+0)->h"
	opsFDCB[0x9D] = "res 3,(iy+0)->l"
	opsFDCB[0x9E] = "RES  3,(IY+0)"
	opsFDCB[0x9F] = "res 3,(iy+0)->a"
	opsFDCB[0xA0] = "res 4,(iy+0)->b"
	opsFDCB[0xA1] = "res 4,(iy+0)->c"
	opsFDCB[0xA2] = "res 4,(iy+0)->d"
	opsFDCB[0xA3] = "res 4,(iy+0)->e"
	opsFDCB[0xA4] = "res 4,(iy+0)->h"
	opsFDCB[0xA5] = "res 4,(iy+0)->l"
	opsFDCB[0xA6] = "RES  4,(IY+0)"
	opsFDCB[0xA7] = "res 4,(iy+0)->a"
	opsFDCB[0xA8] = "res 5,(iy+0)->b"
	opsFDCB[0xA9] = "res 5,(iy+0)->c"
	opsFDCB[0xAA] = "res 5,(iy+0)->d"
	opsFDCB[0xAB] = "res 5,(iy+0)->e"
	opsFDCB[0xAC] = "res 5,(iy+0)->h"
	opsFDCB[0xAD] = "res 5,(iy+0)->l"
	opsFDCB[0xAE] = "RES  5,(IY+0)"
	opsFDCB[0xAF] = "res 5,(iy+0)->a"
	opsFDCB[0xB0] = "res 6,(iy+0)->b"
	opsFDCB[0xB1] = "res 6,(iy+0)->c"
	opsFDCB[0xB2] = "res 6,(iy+0)->d"
	opsFDCB[0xB3] = "res 6,(iy+0)->e"
	opsFDCB[0xB4] = "res 6,(iy+0)->h"
	opsFDCB[0xB5] = "res 6,(iy+0)->l"
	opsFDCB[0xB6] = "RES  6,(IY+0)"
	opsFDCB[0xB7] = "res 6,(iy+0)->a"
	opsFDCB[0xB8] = "res 7,(iy+0)->b"
	opsFDCB[0xB9] = "res 7,(iy+0)->c"
	opsFDCB[0xBA] = "res 7,(iy+0)->d"
	opsFDCB[0xBB] = "res 7,(iy+0)->e"
	opsFDCB[0xBC] = "res 7,(iy+0)->h"
	opsFDCB[0xBD] = "res 7,(iy+0)->l"
	opsFDCB[0xBE] = "RES  7,(IY+0)"
	opsFDCB[0xBF] = "res 7,(iy+0)->a"
	opsFDCB[0xC0] = "set 0,(iy+0)->b"
	opsFDCB[0xC1] = "set 0,(iy+0)->c"
	opsFDCB[0xC2] = "set 0,(iy+0)->d"
	opsFDCB[0xC3] = "set 0,(iy+0)->e"
	opsFDCB[0xC4] = "set 0,(iy+0)->h"
	opsFDCB[0xC5] = "set 0,(iy+0)->l"
	opsFDCB[0xC6] = "SET  0,(IY+0)"
	opsFDCB[0xC7] = "set 0,(iy+0)->a"
	opsFDCB[0xC8] = "set 1,(iy+0)->b"
	opsFDCB[0xC9] = "set 1,(iy+0)->c"
	opsFDCB[0xCA] = "set 1,(iy+0)->d"
	opsFDCB[0xCB] = "set 1,(iy+0)->e"
	opsFDCB[0xCC] = "set 1,(iy+0)->h"
	opsFDCB[0xCD] = "set 1,(iy+0)->l"
	opsFDCB[0xCE] = "SET  1,(IY+0)"
	opsFDCB[0xCF] = "set 1,(iy+0)->a"
	opsFDCB[0xD0] = "set 2,(iy+0)->b"
	opsFDCB[0xD1] = "set 2,(iy+0)->c"
	opsFDCB[0xD2] = "set 2,(iy+0)->d"
	opsFDCB[0xD3] = "set 2,(iy+0)->e"
	opsFDCB[0xD4] = "set 2,(iy+0)->h"
	opsFDCB[0xD5] = "set 2,(iy+0)->l"
	opsFDCB[0xD6] = "SET  2,(IY+0)"
	opsFDCB[0xD7] = "set 2,(iy+0)->a"
	opsFDCB[0xD8] = "set 3,(iy+0)->b"
	opsFDCB[0xD9] = "set 3,(iy+0)->c"
	opsFDCB[0xDA] = "set 3,(iy+0)->d"
	opsFDCB[0xDB] = "set 3,(iy+0)->e"
	opsFDCB[0xDC] = "set 3,(iy+0)->h"
	opsFDCB[0xDD] = "set 3,(iy+0)->l"
	opsFDCB[0xDE] = "SET  3,(IY+0)"
	opsFDCB[0xDF] = "set 3,(iy+0)->a"
	opsFDCB[0xE0] = "set 4,(iy+0)->b"
	opsFDCB[0xE1] = "set 4,(iy+0)->c"
	opsFDCB[0xE2] = "set 4,(iy+0)->d"
	opsFDCB[0xE3] = "set 4,(iy+0)->e"
	opsFDCB[0xE4] = "set 4,(iy+0)->h"
	opsFDCB[0xE5] = "set 4,(iy+0)->l"
	opsFDCB[0xE6] = "SET  4,(IY+0)"
	opsFDCB[0xE7] = "set 4,(iy+0)->a"
	opsFDCB[0xE8] = "set 5,(iy+0)->b"
	opsFDCB[0xE9] = "set 5,(iy+0)->c"
	opsFDCB[0xEA] = "set 5,(iy+0)->d"
	opsFDCB[0xEB] = "set 5,(iy+0)->e"
	opsFDCB[0xEC] = "set 5,(iy+0)->h"
	opsFDCB[0xED] = "set 5,(iy+0)->l"
	opsFDCB[0xEE] = "SET  5,(IY+0)"
	opsFDCB[0xEF] = "set 5,(iy+0)->a"
	opsFDCB[0xF0] = "set 6,(iy+0)->b"
	opsFDCB[0xF1] = "set 6,(iy+0)->c"
	opsFDCB[0xF2] = "set 6,(iy+0)->d"
	opsFDCB[0xF3] = "set 6,(iy+0)->e"
	opsFDCB[0xF4] = "set 6,(iy+0)->h"
	opsFDCB[0xF5] = "set 6,(iy+0)->l"
	opsFDCB[0xF6] = "SET  6,(IY+0)"
	opsFDCB[0xF7] = "set 6,(iy+0)->a"
	opsFDCB[0xF8] = "set 7,(iy+0)->b"
	opsFDCB[0xF9] = "set 7,(iy+0)->c"
	opsFDCB[0xFA] = "set 7,(iy+0)->d"
	opsFDCB[0xFB] = "set 7,(iy+0)->e"
	opsFDCB[0xFC] = "set 7,(iy+0)->h"
	opsFDCB[0xFD] = "set 7,(iy+0)->l"
	opsFDCB[0xFE] = "SET  7,(IY+0)"
	opsFDCB[0xFF] = "set 7,(iy+0)->a"

	opsDDCB[0x00] = "rlc (ix+0)->b"
	opsDDCB[0x01] = "rlc (ix+0)->c"
	opsDDCB[0x02] = "rlc (ix+0)->d"
	opsDDCB[0x03] = "rlc (ix+0)->e"
	opsDDCB[0x04] = "rlc (ix+0)->h"
	opsDDCB[0x05] = "rlc (ix+0)->l"
	opsDDCB[0x06] = "RLC  (ix+0)"
	opsDDCB[0x07] = "rlc (ix+0)->a"
	opsDDCB[0x08] = "rrc (ix+0)->b"
	opsDDCB[0x09] = "rrc (ix+0)->c"
	opsDDCB[0x0A] = "rrc (ix+0)->d"
	opsDDCB[0x0B] = "rrc (ix+0)->e"
	opsDDCB[0x0C] = "rrc (ix+0)->h"
	opsDDCB[0x0D] = "rrc (ix+0)->l"
	opsDDCB[0x0E] = "RRC  (ix+0)"
	opsDDCB[0x0F] = "rrc (ix+0)->a"
	opsDDCB[0x10] = "rl  (ix+0)->b"
	opsDDCB[0x11] = "rl  (ix+0)->c"
	opsDDCB[0x12] = "rl  (ix+0)->d"
	opsDDCB[0x13] = "rl  (ix+0)->e"
	opsDDCB[0x14] = "rl  (ix+0)->h"
	opsDDCB[0x15] = "rl  (ix+0)->l"
	opsDDCB[0x16] = "RL   (ix+0)"
	opsDDCB[0x17] = "rl  (ix+0)->a"
	opsDDCB[0x18] = "rr  (ix+0)->b"
	opsDDCB[0x19] = "rr  (ix+0)->c"
	opsDDCB[0x1A] = "rr  (ix+0)->d"
	opsDDCB[0x1B] = "rr  (ix+0)->e"
	opsDDCB[0x1C] = "rr  (ix+0)->h"
	opsDDCB[0x1D] = "rr  (ix+0)->l"
	opsDDCB[0x1E] = "RR   (ix+0)"
	opsDDCB[0x1F] = "rr  (ix+0)->a"
	opsDDCB[0x20] = "sla (ix+0)->b"
	opsDDCB[0x21] = "sla (ix+0)->c"
	opsDDCB[0x22] = "sla (ix+0)->d"
	opsDDCB[0x23] = "sla (ix+0)->e"
	opsDDCB[0x24] = "sla (ix+0)->h"
	opsDDCB[0x25] = "sla (ix+0)->l"
	opsDDCB[0x26] = "SLA  (ix+0)"
	opsDDCB[0x27] = "sla (ix+0)->a"
	opsDDCB[0x28] = "sra (ix+0)->b"
	opsDDCB[0x29] = "sra (ix+0)->c"
	opsDDCB[0x2A] = "sra (ix+0)->d"
	opsDDCB[0x2B] = "sra (ix+0)->e"
	opsDDCB[0x2C] = "sra (ix+0)->h"
	opsDDCB[0x2D] = "sra (ix+0)->l"
	opsDDCB[0x2E] = "SRA  (ix+0)"
	opsDDCB[0x2F] = "sra (ix+0)->a"
	opsDDCB[0x30] = "sls (ix+0)->b"
	opsDDCB[0x31] = "sls (ix+0)->c"
	opsDDCB[0x32] = "sls (ix+0)->d"
	opsDDCB[0x33] = "sls (ix+0)->e"
	opsDDCB[0x34] = "sls (ix+0)->h"
	opsDDCB[0x35] = "sls (ix+0)->l"
	opsDDCB[0x36] = "SLS  (ix+0)"
	opsDDCB[0x37] = "sls (ix+0)->a"
	opsDDCB[0x38] = "srl (ix+0)->b"
	opsDDCB[0x39] = "srl (ix+0)->c"
	opsDDCB[0x3A] = "srl (ix+0)->d"
	opsDDCB[0x3B] = "srl (ix+0)->e"
	opsDDCB[0x3C] = "srl (ix+0)->h"
	opsDDCB[0x3D] = "srl (ix+0)->l"
	opsDDCB[0x3E] = "SRL  (ix+0)"
	opsDDCB[0x3F] = "srl (ix+0)->a"
	opsDDCB[0x40] = "bit 0,(ix+0)->b"
	opsDDCB[0x41] = "bit 0,(ix+0)->c"
	opsDDCB[0x42] = "bit 0,(ix+0)->d"
	opsDDCB[0x43] = "bit 0,(ix+0)->e"
	opsDDCB[0x44] = "bit 0,(ix+0)->h"
	opsDDCB[0x45] = "bit 0,(ix+0)->l"
	opsDDCB[0x46] = "BIT  0,(ix+0)"
	opsDDCB[0x47] = "bit 0,(ix+0)->a"
	opsDDCB[0x48] = "bit 1,(ix+0)->b"
	opsDDCB[0x49] = "bit 1,(ix+0)->c"
	opsDDCB[0x4A] = "bit 1,(ix+0)->d"
	opsDDCB[0x4B] = "bit 1,(ix+0)->e"
	opsDDCB[0x4C] = "bit 1,(ix+0)->h"
	opsDDCB[0x4D] = "bit 1,(ix+0)->l"
	opsDDCB[0x4E] = "BIT  1,(ix+0)"
	opsDDCB[0x4F] = "bit 1,(ix+0)->a"
	opsDDCB[0x50] = "bit 2,(ix+0)->b"
	opsDDCB[0x51] = "bit 2,(ix+0)->c"
	opsDDCB[0x52] = "bit 2,(ix+0)->d"
	opsDDCB[0x53] = "bit 2,(ix+0)->e"
	opsDDCB[0x54] = "bit 2,(ix+0)->h"
	opsDDCB[0x55] = "bit 2,(ix+0)->l"
	opsDDCB[0x56] = "BIT  2,(ix+0)"
	opsDDCB[0x57] = "bit 2,(ix+0)->a"
	opsDDCB[0x58] = "bit 3,(ix+0)->b"
	opsDDCB[0x59] = "bit 3,(ix+0)->c"
	opsDDCB[0x5A] = "bit 3,(ix+0)->d"
	opsDDCB[0x5B] = "bit 3,(ix+0)->e"
	opsDDCB[0x5C] = "bit 3,(ix+0)->h"
	opsDDCB[0x5D] = "bit 3,(ix+0)->l"
	opsDDCB[0x5E] = "BIT  3,(ix+0)"
	opsDDCB[0x5F] = "bit 3,(ix+0)->a"
	opsDDCB[0x60] = "bit 4,(ix+0)->b"
	opsDDCB[0x61] = "bit 4,(ix+0)->c"
	opsDDCB[0x62] = "bit 4,(ix+0)->d"
	opsDDCB[0x63] = "bit 4,(ix+0)->e"
	opsDDCB[0x64] = "bit 4,(ix+0)->h"
	opsDDCB[0x65] = "bit 4,(ix+0)->l"
	opsDDCB[0x66] = "BIT  4,(ix+0)"
	opsDDCB[0x67] = "bit 4,(ix+0)->a"
	opsDDCB[0x68] = "bit 5,(ix+0)->b"
	opsDDCB[0x69] = "bit 5,(ix+0)->c"
	opsDDCB[0x6A] = "bit 5,(ix+0)->d"
	opsDDCB[0x6B] = "bit 5,(ix+0)->e"
	opsDDCB[0x6C] = "bit 5,(ix+0)->h"
	opsDDCB[0x6D] = "bit 5,(ix+0)->l"
	opsDDCB[0x6E] = "BIT  5,(ix+0)"
	opsDDCB[0x6F] = "bit 5,(ix+0)->a"
	opsDDCB[0x70] = "bit 6,(ix+0)->b"
	opsDDCB[0x71] = "bit 6,(ix+0)->c"
	opsDDCB[0x72] = "bit 6,(ix+0)->d"
	opsDDCB[0x73] = "bit 6,(ix+0)->e"
	opsDDCB[0x74] = "bit 6,(ix+0)->h"
	opsDDCB[0x75] = "bit 6,(ix+0)->l"
	opsDDCB[0x76] = "BIT  6,(ix+0)"
	opsDDCB[0x77] = "bit 6,(ix+0)->a"
	opsDDCB[0x78] = "bit 7,(ix+0)->b"
	opsDDCB[0x79] = "bit 7,(ix+0)->c"
	opsDDCB[0x7A] = "bit 7,(ix+0)->d"
	opsDDCB[0x7B] = "bit 7,(ix+0)->e"
	opsDDCB[0x7C] = "bit 7,(ix+0)->h"
	opsDDCB[0x7D] = "bit 7,(ix+0)->l"
	opsDDCB[0x7E] = "BIT  7,(ix+0)"
	opsDDCB[0x7F] = "bit 7,(ix+0)->a"
	opsDDCB[0x80] = "res 0,(ix+0)->b"
	opsDDCB[0x81] = "res 0,(ix+0)->c"
	opsDDCB[0x82] = "res 0,(ix+0)->d"
	opsDDCB[0x83] = "res 0,(ix+0)->e"
	opsDDCB[0x84] = "res 0,(ix+0)->h"
	opsDDCB[0x85] = "res 0,(ix+0)->l"
	opsDDCB[0x86] = "RES  0,(ix+0)"
	opsDDCB[0x87] = "res 0,(ix+0)->a"
	opsDDCB[0x88] = "res 1,(ix+0)->b"
	opsDDCB[0x89] = "res 1,(ix+0)->c"
	opsDDCB[0x8A] = "res 1,(ix+0)->d"
	opsDDCB[0x8B] = "res 1,(ix+0)->e"
	opsDDCB[0x8C] = "res 1,(ix+0)->h"
	opsDDCB[0x8D] = "res 1,(ix+0)->l"
	opsDDCB[0x8E] = "RES  1,(ix+0)"
	opsDDCB[0x8F] = "res 1,(ix+0)->a"
	opsDDCB[0x90] = "res 2,(ix+0)->b"
	opsDDCB[0x91] = "res 2,(ix+0)->c"
	opsDDCB[0x92] = "res 2,(ix+0)->d"
	opsDDCB[0x93] = "res 2,(ix+0)->e"
	opsDDCB[0x94] = "res 2,(ix+0)->h"
	opsDDCB[0x95] = "res 2,(ix+0)->l"
	opsDDCB[0x96] = "RES  2,(ix+0)"
	opsDDCB[0x97] = "res 2,(ix+0)->a"
	opsDDCB[0x98] = "res 3,(ix+0)->b"
	opsDDCB[0x99] = "res 3,(ix+0)->c"
	opsDDCB[0x9A] = "res 3,(ix+0)->d"
	opsDDCB[0x9B] = "res 3,(ix+0)->e"
	opsDDCB[0x9C] = "res 3,(ix+0)->h"
	opsDDCB[0x9D] = "res 3,(ix+0)->l"
	opsDDCB[0x9E] = "RES  3,(ix+0)"
	opsDDCB[0x9F] = "res 3,(ix+0)->a"
	opsDDCB[0xA0] = "res 4,(ix+0)->b"
	opsDDCB[0xA1] = "res 4,(ix+0)->c"
	opsDDCB[0xA2] = "res 4,(ix+0)->d"
	opsDDCB[0xA3] = "res 4,(ix+0)->e"
	opsDDCB[0xA4] = "res 4,(ix+0)->h"
	opsDDCB[0xA5] = "res 4,(ix+0)->l"
	opsDDCB[0xA6] = "RES  4,(ix+0)"
	opsDDCB[0xA7] = "res 4,(ix+0)->a"
	opsDDCB[0xA8] = "res 5,(ix+0)->b"
	opsDDCB[0xA9] = "res 5,(ix+0)->c"
	opsDDCB[0xAA] = "res 5,(ix+0)->d"
	opsDDCB[0xAB] = "res 5,(ix+0)->e"
	opsDDCB[0xAC] = "res 5,(ix+0)->h"
	opsDDCB[0xAD] = "res 5,(ix+0)->l"
	opsDDCB[0xAE] = "RES  5,(ix+0)"
	opsDDCB[0xAF] = "res 5,(ix+0)->a"
	opsDDCB[0xB0] = "res 6,(ix+0)->b"
	opsDDCB[0xB1] = "res 6,(ix+0)->c"
	opsDDCB[0xB2] = "res 6,(ix+0)->d"
	opsDDCB[0xB3] = "res 6,(ix+0)->e"
	opsDDCB[0xB4] = "res 6,(ix+0)->h"
	opsDDCB[0xB5] = "res 6,(ix+0)->l"
	opsDDCB[0xB6] = "RES  6,(ix+0)"
	opsDDCB[0xB7] = "res 6,(ix+0)->a"
	opsDDCB[0xB8] = "res 7,(ix+0)->b"
	opsDDCB[0xB9] = "res 7,(ix+0)->c"
	opsDDCB[0xBA] = "res 7,(ix+0)->d"
	opsDDCB[0xBB] = "res 7,(ix+0)->e"
	opsDDCB[0xBC] = "res 7,(ix+0)->h"
	opsDDCB[0xBD] = "res 7,(ix+0)->l"
	opsDDCB[0xBE] = "RES  7,(ix+0)"
	opsDDCB[0xBF] = "res 7,(ix+0)->a"
	opsDDCB[0xC0] = "set 0,(ix+0)->b"
	opsDDCB[0xC1] = "set 0,(ix+0)->c"
	opsDDCB[0xC2] = "set 0,(ix+0)->d"
	opsDDCB[0xC3] = "set 0,(ix+0)->e"
	opsDDCB[0xC4] = "set 0,(ix+0)->h"
	opsDDCB[0xC5] = "set 0,(ix+0)->l"
	opsDDCB[0xC6] = "SET  0,(ix+0)"
	opsDDCB[0xC7] = "set 0,(ix+0)->a"
	opsDDCB[0xC8] = "set 1,(ix+0)->b"
	opsDDCB[0xC9] = "set 1,(ix+0)->c"
	opsDDCB[0xCA] = "set 1,(ix+0)->d"
	opsDDCB[0xCB] = "set 1,(ix+0)->e"
	opsDDCB[0xCC] = "set 1,(ix+0)->h"
	opsDDCB[0xCD] = "set 1,(ix+0)->l"
	opsDDCB[0xCE] = "SET  1,(ix+0)"
	opsDDCB[0xCF] = "set 1,(ix+0)->a"
	opsDDCB[0xD0] = "set 2,(ix+0)->b"
	opsDDCB[0xD1] = "set 2,(ix+0)->c"
	opsDDCB[0xD2] = "set 2,(ix+0)->d"
	opsDDCB[0xD3] = "set 2,(ix+0)->e"
	opsDDCB[0xD4] = "set 2,(ix+0)->h"
	opsDDCB[0xD5] = "set 2,(ix+0)->l"
	opsDDCB[0xD6] = "SET  2,(ix+0)"
	opsDDCB[0xD7] = "set 2,(ix+0)->a"
	opsDDCB[0xD8] = "set 3,(ix+0)->b"
	opsDDCB[0xD9] = "set 3,(ix+0)->c"
	opsDDCB[0xDA] = "set 3,(ix+0)->d"
	opsDDCB[0xDB] = "set 3,(ix+0)->e"
	opsDDCB[0xDC] = "set 3,(ix+0)->h"
	opsDDCB[0xDD] = "set 3,(ix+0)->l"
	opsDDCB[0xDE] = "SET  3,(ix+0)"
	opsDDCB[0xDF] = "set 3,(ix+0)->a"
	opsDDCB[0xE0] = "set 4,(ix+0)->b"
	opsDDCB[0xE1] = "set 4,(ix+0)->c"
	opsDDCB[0xE2] = "set 4,(ix+0)->d"
	opsDDCB[0xE3] = "set 4,(ix+0)->e"
	opsDDCB[0xE4] = "set 4,(ix+0)->h"
	opsDDCB[0xE5] = "set 4,(ix+0)->l"
	opsDDCB[0xE6] = "SET  4,(ix+0)"
	opsDDCB[0xE7] = "set 4,(ix+0)->a"
	opsDDCB[0xE8] = "set 5,(ix+0)->b"
	opsDDCB[0xE9] = "set 5,(ix+0)->c"
	opsDDCB[0xEA] = "set 5,(ix+0)->d"
	opsDDCB[0xEB] = "set 5,(ix+0)->e"
	opsDDCB[0xEC] = "set 5,(ix+0)->h"
	opsDDCB[0xED] = "set 5,(ix+0)->l"
	opsDDCB[0xEE] = "SET  5,(ix+0)"
	opsDDCB[0xEF] = "set 5,(ix+0)->a"
	opsDDCB[0xF0] = "set 6,(ix+0)->b"
	opsDDCB[0xF1] = "set 6,(ix+0)->c"
	opsDDCB[0xF2] = "set 6,(ix+0)->d"
	opsDDCB[0xF3] = "set 6,(ix+0)->e"
	opsDDCB[0xF4] = "set 6,(ix+0)->h"
	opsDDCB[0xF5] = "set 6,(ix+0)->l"
	opsDDCB[0xF6] = "SET  6,(ix+0)"
	opsDDCB[0xF7] = "set 6,(ix+0)->a"
	opsDDCB[0xF8] = "set 7,(ix+0)->b"
	opsDDCB[0xF9] = "set 7,(ix+0)->c"
	opsDDCB[0xFA] = "set 7,(ix+0)->d"
	opsDDCB[0xFB] = "set 7,(ix+0)->e"
	opsDDCB[0xFC] = "set 7,(ix+0)->h"
	opsDDCB[0xFD] = "set 7,(ix+0)->l"
	opsDDCB[0xFE] = "SET  7,(ix+0)"
	opsDDCB[0xFF] = "set 7,(ix+0)->a"

	opsED[0x00] = "MOS_QUIT"
	opsED[0x01] = "MOS_CLI"
	opsED[0x02] = "MOS_BYTE"
	opsED[0x03] = "MOS_WORD"
	opsED[0x04] = "MOS_WRCH"
	opsED[0x05] = "MOS_RDCH"
	opsED[0x06] = "MOS_FILE"
	opsED[0x07] = "MOS_ARGS"
	opsED[0x08] = "MOS_BGET"
	opsED[0x09] = "MOS_BPUT"
	opsED[0x0A] = "MOS_GBPB"
	opsED[0x0B] = "MOS_FIND"
	opsED[0x0C] = "MOS_FF0C"
	opsED[0x0D] = "MOS_FF0D"
	opsED[0x0E] = "MOS_FF0E"
	opsED[0x0F] = "MOS_FF0F"
	opsED[0x40] = "IN   B,(C)"
	opsED[0x41] = "OUT  (C),B"
	opsED[0x42] = "SBC  HL,BC"
	opsED[0x43] = "LD   ($nn),BC"
	opsED[0x44] = "NEG"
	opsED[0x45] = "RETN"
	opsED[0x46] = "IM   0"
	opsED[0x47] = "LD   I,A"
	opsED[0x48] = "IN   C,(C)"
	opsED[0x49] = "OUT  (C),C"
	opsED[0x4A] = "ADC  HL,BC"
	opsED[0x4B] = "LD   BC,($nn)"
	opsED[0x4C] = "[neg]"
	opsED[0x4D] = "RETI"
	opsED[0x4E] = "[im0]"
	opsED[0x4F] = "LD   R,A"
	opsED[0x50] = "IN   D,(C)"
	opsED[0x51] = "OUT  (C),D"
	opsED[0x52] = "SBC  HL,DE"
	opsED[0x53] = "LD   ($nn),DE"
	opsED[0x54] = "[neg]"
	opsED[0x55] = "[retn]"
	opsED[0x56] = "IM   1"
	opsED[0x57] = "LD   A,I"
	opsED[0x58] = "IN   E,(C)"
	opsED[0x59] = "OUT  (C),E"
	opsED[0x5A] = "ADC  HL,DE"
	opsED[0x5B] = "LD   DE,($nn)"
	opsED[0x5C] = "[neg]"
	opsED[0x5D] = "[reti]"
	opsED[0x5E] = "IM   2"
	opsED[0x5F] = "LD   A,R"
	opsED[0x60] = "IN   H,(C)"
	opsED[0x61] = "OUT  (C),H"
	opsED[0x62] = "SBC  HL,HL"
	opsED[0x63] = "LD   ($nn),HL"
	opsED[0x64] = "[neg]"
	opsED[0x65] = "[retn]"
	opsED[0x66] = "[im0]"
	opsED[0x67] = "RRD"
	opsED[0x68] = "IN   L,(C)"
	opsED[0x69] = "OUT  (C),L"
	opsED[0x6A] = "ADC  HL,HL"
	opsED[0x6B] = "LD   HL,($nn)"
	opsED[0x6C] = "[neg]"
	opsED[0x6D] = "[reti]"
	opsED[0x6E] = "[im0]"
	opsED[0x6F] = "RLD"
	opsED[0x70] = "IN   F,(C)"
	opsED[0x71] = "OUT  (C),F"
	opsED[0x72] = "SBC  HL,SP"
	opsED[0x73] = "LD   ($nn),SP"
	opsED[0x74] = "[neg]"
	opsED[0x75] = "[retn]"
	opsED[0x76] = "[im1]"
	opsED[0x77] = "[ld i,i?]"
	opsED[0x78] = "IN   A,(C)"
	opsED[0x79] = "OUT  (C),A"
	opsED[0x7A] = "ADC  HL,SP"
	opsED[0x7B] = "LD   SP,($nn)"
	opsED[0x7C] = "[neg]"
	opsED[0x7D] = "[reti]"
	opsED[0x7E] = "[im2]"
	opsED[0x7F] = "[ld r,r?]"
	opsED[0xA0] = "LDI"
	opsED[0xA1] = "CPI"
	opsED[0xA2] = "INI"
	opsED[0xA3] = "OTI"
	opsED[0xA8] = "LDD"
	opsED[0xA9] = "CPD"
	opsED[0xAA] = "IND"
	opsED[0xAB] = "OTD"
	opsED[0xB0] = "LDIR"
	opsED[0xB1] = "CPIR"
	opsED[0xB2] = "INIR"
	opsED[0xB3] = "OTIR"
	opsED[0xB8] = "LDDR"
	opsED[0xB9] = "CPDR"
	opsED[0xBA] = "INDR"
	opsED[0xBB] = "OTDR"
	opsED[0xF8] = "[z80]"
	opsED[0xF9] = "[z80]"
	opsED[0xFA] = "[z80]"
	opsED[0xFB] = "ED_LOAD"
	opsED[0xFC] = "[z80]"
	opsED[0xFD] = "[z80]"
	opsED[0xFE] = "[z80]"
	opsED[0xFF] = "ED_DOS"
}

func (fd *fetchedData) disassemble() string {
	op := ""
	switch fd.prefix {
	case 0:
		op = ops[fd.opCode]

	case 0xed:
		op = opsED[fd.opCode]

	case 0xcb:
		op = opsCB[fd.opCode]

	case 0xdd:
		op = opsDD[fd.opCode]

	case 0xfd:
		op = opsFD[fd.opCode]

	case 0xDDCB:
		op = opsDDCB[fd.opCode]

	case 0xFDCB:
		op = opsFDCB[fd.opCode]

	default:
		panic(fmt.Sprintf("--> 0x%X (%s)", fd.prefix, fd.op.name))
	}

	if strings.HasPrefix(op, "JR ") || strings.HasPrefix(op, "DJNZ ") {
		jump := int8(fd.n)
		pc := fd.pc + uint16(jump)
		op = strings.ReplaceAll(op, "$nn", toHex16(pc+2))
	} else {
		op = strings.ReplaceAll(op, "$nn", toHex16(fd.nn))
		op = strings.ReplaceAll(op, "$n2", toHex8(fd.n2))
		op = strings.ReplaceAll(op, "$n", toHex8(fd.n))
		op = strings.ReplaceAll(op, "+0", "+"+toHex8(fd.n))
	}
	return strings.ToLower(fmt.Sprintf("%04x: %s", fd.pc, op))
	// return strings.ToLower(fmt.Sprintf("%04x: %s %s", fd.pc, fd.getMemory(), op))
}

func toHex16(v uint16) string {
	n := "000" + strconv.FormatUint(uint64(v), 16)
	return "$" + n[len(n)-4:]
}

func toHex8(v uint8) string {
	n := "0" + strconv.FormatUint(uint64(v), 16)
	return "$" + n[len(n)-2:]
}
