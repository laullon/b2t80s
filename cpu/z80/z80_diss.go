package z80

import (
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

	ops[0x00] = "nop"
	ops[0x01] = "ld   bc,$nn"
	ops[0x02] = "ld   (bc),a"
	ops[0x03] = "inc  bc"
	ops[0x04] = "inc  b"
	ops[0x05] = "dec  b"
	ops[0x06] = "ld   b,$n"
	ops[0x07] = "rlca"
	ops[0x08] = "ex   af,af'"
	ops[0x09] = "add  hl,bc"
	ops[0x0A] = "ld   a,(bc)"
	ops[0x0B] = "dec  bc"
	ops[0x0C] = "inc  c"
	ops[0x0D] = "dec  c"
	ops[0x0E] = "ld   c,$n"
	ops[0x0F] = "rrca"
	ops[0x10] = "djnz $nn"
	ops[0x11] = "ld   de,$nn"
	ops[0x12] = "ld   (de),a"
	ops[0x13] = "inc  de"
	ops[0x14] = "inc  d"
	ops[0x15] = "dec  d"
	ops[0x16] = "ld   d,$n"
	ops[0x17] = "rla"
	ops[0x18] = "jr   $nn"
	ops[0x19] = "add  hl,de"
	ops[0x1A] = "ld   a,(de)"
	ops[0x1B] = "dec  de"
	ops[0x1C] = "inc  e"
	ops[0x1D] = "dec  e"
	ops[0x1E] = "ld   e,$n"
	ops[0x1F] = "rra"
	ops[0x20] = "jr   nz,$nn"
	ops[0x21] = "ld   hl,$nn"
	ops[0x22] = "ld   ($nn),hl"
	ops[0x23] = "inc  hl"
	ops[0x24] = "inc  h"
	ops[0x25] = "dec  h"
	ops[0x26] = "ld   h,$n"
	ops[0x27] = "daa"
	ops[0x28] = "jr   z,$nn"
	ops[0x29] = "add  hl,hl"
	ops[0x2A] = "ld   hl,($nn)"
	ops[0x2B] = "dec  hl"
	ops[0x2C] = "inc  l"
	ops[0x2D] = "dec  l"
	ops[0x2E] = "ld   l,$n"
	ops[0x2F] = "cpl"
	ops[0x30] = "jr   nc,$nn"
	ops[0x31] = "ld   sp,$nn"
	ops[0x32] = "ld   ($nn),a"
	ops[0x33] = "inc  sp"
	ops[0x34] = "inc  (hl)"
	ops[0x35] = "dec  (hl)"
	ops[0x36] = "ld   (hl),$n"
	ops[0x37] = "scf"
	ops[0x38] = "jr   c,$nn"
	ops[0x39] = "add  hl,sp"
	ops[0x3A] = "ld   a,($nn)"
	ops[0x3B] = "dec  sp"
	ops[0x3C] = "inc  a"
	ops[0x3D] = "dec  a"
	ops[0x3E] = "ld   a,$n"
	ops[0x3F] = "ccf"
	ops[0x40] = "ld   b,b"
	ops[0x41] = "ld   b,c"
	ops[0x42] = "ld   b,d"
	ops[0x43] = "ld   b,e"
	ops[0x44] = "ld   b,h"
	ops[0x45] = "ld   b,l"
	ops[0x46] = "ld   b,(hl)"
	ops[0x47] = "ld   b,a"
	ops[0x48] = "ld   c,b"
	ops[0x49] = "ld   c,c"
	ops[0x4A] = "ld   c,d"
	ops[0x4B] = "ld   c,e"
	ops[0x4C] = "ld   c,h"
	ops[0x4D] = "ld   c,l"
	ops[0x4E] = "ld   c,(hl)"
	ops[0x4F] = "ld   c,a"
	ops[0x50] = "ld   d,b"
	ops[0x51] = "ld   d,c"
	ops[0x52] = "ld   d,d"
	ops[0x53] = "ld   d,e"
	ops[0x54] = "ld   d,h"
	ops[0x55] = "ld   d,l"
	ops[0x56] = "ld   d,(hl)"
	ops[0x57] = "ld   d,a"
	ops[0x58] = "ld   e,b"
	ops[0x59] = "ld   e,c"
	ops[0x5A] = "ld   e,d"
	ops[0x5B] = "ld   e,e"
	ops[0x5C] = "ld   e,h"
	ops[0x5D] = "ld   e,l"
	ops[0x5E] = "ld   e,(hl)"
	ops[0x5F] = "ld   e,a"
	ops[0x60] = "ld   h,b"
	ops[0x61] = "ld   h,c"
	ops[0x62] = "ld   h,d"
	ops[0x63] = "ld   h,e"
	ops[0x64] = "ld   h,h"
	ops[0x65] = "ld   h,l"
	ops[0x66] = "ld   h,(hl)"
	ops[0x67] = "ld   h,a"
	ops[0x68] = "ld   l,b"
	ops[0x69] = "ld   l,c"
	ops[0x6A] = "ld   l,d"
	ops[0x6B] = "ld   l,e"
	ops[0x6C] = "ld   l,h"
	ops[0x6D] = "ld   l,l"
	ops[0x6E] = "ld   l,(hl)"
	ops[0x6F] = "ld   l,a"
	ops[0x70] = "ld   (hl),b"
	ops[0x71] = "ld   (hl),c"
	ops[0x72] = "ld   (hl),d"
	ops[0x73] = "ld   (hl),e"
	ops[0x74] = "ld   (hl),h"
	ops[0x75] = "ld   (hl),l"
	ops[0x76] = "halt"
	ops[0x77] = "ld   (hl),a"
	ops[0x78] = "ld   a,b"
	ops[0x79] = "ld   a,c"
	ops[0x7A] = "ld   a,d"
	ops[0x7B] = "ld   a,e"
	ops[0x7C] = "ld   a,h"
	ops[0x7D] = "ld   a,l"
	ops[0x7E] = "ld   a,(hl)"
	ops[0x7F] = "ld   a,a"
	ops[0x80] = "add  a,b"
	ops[0x81] = "add  a,c"
	ops[0x82] = "add  a,d"
	ops[0x83] = "add  a,e"
	ops[0x84] = "add  a,h"
	ops[0x85] = "add  a,l"
	ops[0x86] = "add  a,(hl)"
	ops[0x87] = "add  a,a"
	ops[0x88] = "adc  a,b"
	ops[0x89] = "adc  a,c"
	ops[0x8A] = "adc  a,d"
	ops[0x8B] = "adc  a,e"
	ops[0x8C] = "adc  a,h"
	ops[0x8D] = "adc  a,l"
	ops[0x8E] = "adc  a,(hl)"
	ops[0x8F] = "adc  a,a"
	ops[0x90] = "sub  a,b"
	ops[0x91] = "sub  a,c"
	ops[0x92] = "sub  a,d"
	ops[0x93] = "sub  a,e"
	ops[0x94] = "sub  a,h"
	ops[0x95] = "sub  a,l"
	ops[0x96] = "sub  a,(hl)"
	ops[0x97] = "sub  a,a"
	ops[0x98] = "sbc  a,b"
	ops[0x99] = "sbc  a,c"
	ops[0x9A] = "sbc  a,d"
	ops[0x9B] = "sbc  a,e"
	ops[0x9C] = "sbc  a,h"
	ops[0x9D] = "sbc  a,l"
	ops[0x9E] = "sbc  a,(hl)"
	ops[0x9F] = "sbc  a,a"
	ops[0xA0] = "and  b"
	ops[0xA1] = "and  c"
	ops[0xA2] = "and  d"
	ops[0xA3] = "and  e"
	ops[0xA4] = "and  h"
	ops[0xA5] = "and  l"
	ops[0xA6] = "and  (hl)"
	ops[0xA7] = "and  a"
	ops[0xA8] = "xor  b"
	ops[0xA9] = "xor  c"
	ops[0xAA] = "xor  d"
	ops[0xAB] = "xor  e"
	ops[0xAC] = "xor  h"
	ops[0xAD] = "xor  l"
	ops[0xAE] = "xor  (hl)"
	ops[0xAF] = "xor  a"
	ops[0xB0] = "or   b"
	ops[0xB1] = "or   c"
	ops[0xB2] = "or   d"
	ops[0xB3] = "or   e"
	ops[0xB4] = "or   h"
	ops[0xB5] = "or   l"
	ops[0xB6] = "or   (hl)"
	ops[0xB7] = "or   a"
	ops[0xB8] = "cp   b"
	ops[0xB9] = "cp   c"
	ops[0xBA] = "cp   d"
	ops[0xBB] = "cp   e"
	ops[0xBC] = "cp   h"
	ops[0xBD] = "cp   l"
	ops[0xBE] = "cp   (hl)"
	ops[0xBF] = "cp   a"
	ops[0xC0] = "ret  nz"
	ops[0xC1] = "pop  bc"
	ops[0xC2] = "jp   nz,$nn"
	ops[0xC3] = "jp   $nn"
	ops[0xC4] = "call nz,$nn"
	ops[0xC5] = "push bc"
	ops[0xC6] = "add  a,$n"
	ops[0xC7] = "rst  $00"
	ops[0xC8] = "ret  z"
	ops[0xC9] = "ret"
	ops[0xCA] = "jp   z,$nn"
	ops[0xCC] = "call z,$nn"
	ops[0xCD] = "call $nn"
	ops[0xCE] = "adc  a,$n"
	ops[0xCF] = "rst  $08"
	ops[0xD0] = "ret  nc"
	ops[0xD1] = "pop  de"
	ops[0xD2] = "jp   nc,$nn"
	ops[0xD3] = "out  ($n),a"
	ops[0xD4] = "call nc,$nn"
	ops[0xD5] = "push de"
	ops[0xD6] = "sub  a,$n"
	ops[0xD7] = "rst  $10"
	ops[0xD8] = "ret  c"
	ops[0xD9] = "exx"
	ops[0xDA] = "jp   c,$nn"
	ops[0xDB] = "in   a,($n)"
	ops[0xDC] = "call c,$nn"
	ops[0xDE] = "sbc  a,$n"
	ops[0xDF] = "rst  $18"
	ops[0xE0] = "ret  po"
	ops[0xE1] = "pop  hl"
	ops[0xE2] = "jp   po,$nn"
	ops[0xE3] = "ex   (sp),hl"
	ops[0xE4] = "call po,$nn"
	ops[0xE5] = "push hl"
	ops[0xE6] = "and  $n"
	ops[0xE7] = "rst  $20"
	ops[0xE8] = "ret  pe"
	ops[0xE9] = "jp   (hl)"
	ops[0xEA] = "jp   pe,$nn"
	ops[0xEB] = "ex   de,hl"
	ops[0xEC] = "call pe,$nn"
	ops[0xEE] = "xor  $n"
	ops[0xEF] = "rst  $28"
	ops[0xF0] = "ret  p"
	ops[0xF1] = "pop  af"
	ops[0xF2] = "jp   p,$nn"
	ops[0xF3] = "di"
	ops[0xF4] = "call p,$nn"
	ops[0xF5] = "push af"
	ops[0xF6] = "or   $n"
	ops[0xF7] = "rst  $30"
	ops[0xF8] = "ret  m"
	ops[0xF9] = "ld   sp,hl"
	ops[0xFA] = "jp   m,$nn"
	ops[0xFB] = "ei"
	ops[0xFC] = "call m,$nn"
	ops[0xFE] = "cp   $n"
	ops[0xFF] = "rst  $38"

	opsDD[0x09] = "add  ix,bc"
	opsDD[0x19] = "add  ix,de"
	opsDD[0x21] = "ld   ix,$nn"
	opsDD[0x22] = "ld  ($nn),ix"
	opsDD[0x23] = "inc  ix"
	opsDD[0x24] = "inc  ixh"
	opsDD[0x25] = "dec  ixh"
	opsDD[0x26] = "ld   ixh,$n"
	opsDD[0x29] = "add  ix,ix"
	opsDD[0x2A] = "ld  ix,($nn)"
	opsDD[0x2B] = "dec  ix"
	opsDD[0x2C] = "inc  ixl"
	opsDD[0x2D] = "dec  ixl"
	opsDD[0x2E] = "ld   ixl,$n"
	opsDD[0x34] = "inc  (ix+0)"
	opsDD[0x35] = "dec  (ix+0)"
	opsDD[0x36] = "ld   (ix+0),$n2"
	opsDD[0x39] = "add  ix,sp"
	opsDD[0x44] = "ld   b,ixh"
	opsDD[0x45] = "ld   b,ixl"
	opsDD[0x46] = "ld   b,(ix+0)"
	opsDD[0x4C] = "ld   c,ixh"
	opsDD[0x4D] = "ld   c,ixl"
	opsDD[0x4E] = "ld   c,(ix+0)"
	opsDD[0x54] = "ld   d,ixh"
	opsDD[0x55] = "ld   d,ixl"
	opsDD[0x56] = "ld   d,(ix+0)"
	opsDD[0x5C] = "ld   e,ixh"
	opsDD[0x5D] = "ld   e,ixl"
	opsDD[0x5E] = "ld   e,(ix+0)"
	opsDD[0x60] = "ld   ixh,b"
	opsDD[0x61] = "ld   ixh,c"
	opsDD[0x62] = "ld   ixh,d"
	opsDD[0x63] = "ld   ixh,e"
	opsDD[0x64] = "ld   ixh,ixh"
	opsDD[0x65] = "ld   ixh,ixl"
	opsDD[0x66] = "ld   h,(ix+0)"
	opsDD[0x67] = "ld   ixh,a"
	opsDD[0x68] = "ld   ixl,b"
	opsDD[0x69] = "ld   ixl,c"
	opsDD[0x6A] = "ld   ixl,d"
	opsDD[0x6B] = "ld   ixl,e"
	opsDD[0x6C] = "ld   ixl,ixh"
	opsDD[0x6D] = "ld   ixl,ixl"
	opsDD[0x6E] = "ld   l,(ix+0)"
	opsDD[0x6F] = "ld   ixl,a"
	opsDD[0x70] = "ld   (ix+0),b"
	opsDD[0x71] = "ld   (ix+0),c"
	opsDD[0x72] = "ld   (ix+0),d"
	opsDD[0x73] = "ld   (ix+0),e"
	opsDD[0x74] = "ld   (ix+0),h"
	opsDD[0x75] = "ld   (ix+0),l"
	opsDD[0x77] = "ld   (ix+0),a"
	opsDD[0x7C] = "ld   a,ixh"
	opsDD[0x7D] = "ld   a,ixl"
	opsDD[0x7E] = "ld   a,(ix+0)"
	opsDD[0x84] = "add  a,ixh"
	opsDD[0x85] = "add  a,ixl"
	opsDD[0x86] = "add  a,(ix+0)"
	opsDD[0x8C] = "adc  a,ixh"
	opsDD[0x8D] = "adc  a,ixl"
	opsDD[0x8E] = "adc  a,(ix+0)"
	opsDD[0x94] = "sub  a,ixh"
	opsDD[0x95] = "sub  a,ixl"
	opsDD[0x96] = "sub  a,(ix+0)"
	opsDD[0x9C] = "sbc  a,ixh"
	opsDD[0x9D] = "sbc  a,ixl"
	opsDD[0x9E] = "sbc  a,(ix+0)"
	opsDD[0xA4] = "and  ixh"
	opsDD[0xA5] = "and  ixl"
	opsDD[0xA6] = "and  (ix+0)"
	opsDD[0xAC] = "xor  ixh"
	opsDD[0xAD] = "xor  ixl"
	opsDD[0xAE] = "xor  (ix+0)"
	opsDD[0xB4] = "or   ixh"
	opsDD[0xB5] = "or   ixl"
	opsDD[0xB6] = "or   (ix+0)"
	opsDD[0xBC] = "cp   ixh"
	opsDD[0xBD] = "cp   ixl"
	opsDD[0xBE] = "cp   (ix+0)"
	opsDD[0xE1] = "pop  ix"
	opsDD[0xE3] = "ex   (sp),ix"
	opsDD[0xE5] = "push ix"
	opsDD[0xE9] = "jp   (ix)"

	opsFD[0x09] = "add  iy,bc"
	opsFD[0x19] = "add  iy,de"
	opsFD[0x21] = "ld   iy,$nn"
	opsFD[0x22] = "ld  ($nn),iy"
	opsFD[0x23] = "inc  iy"
	opsFD[0x24] = "inc  iyh"
	opsFD[0x25] = "dec  iyh"
	opsFD[0x26] = "ld   iyh,$n"
	opsFD[0x29] = "add  iy,iy"
	opsFD[0x2A] = "ld  iy,($nn)"
	opsFD[0x2B] = "dec  iy"
	opsFD[0x2C] = "inc  iyl"
	opsFD[0x2D] = "dec  iyl"
	opsFD[0x2E] = "ld   iyl,$n"
	opsFD[0x34] = "inc  (iy+0)"
	opsFD[0x35] = "dec  (iy+0)"
	opsFD[0x36] = "ld   (iy+0),$n2"
	opsFD[0x39] = "add  iy,sp"
	opsFD[0x44] = "ld   b,iyh"
	opsFD[0x45] = "ld   b,iyl"
	opsFD[0x46] = "ld   b,(iy+0)"
	opsFD[0x4C] = "ld   c,iyh"
	opsFD[0x4D] = "ld   c,iyl"
	opsFD[0x4E] = "ld   c,(iy+0)"
	opsFD[0x54] = "ld   d,iyh"
	opsFD[0x55] = "ld   d,iyl"
	opsFD[0x56] = "ld   d,(iy+0)"
	opsFD[0x5C] = "ld   e,iyh"
	opsFD[0x5D] = "ld   e,iyl"
	opsFD[0x5E] = "ld   e,(iy+0)"
	opsFD[0x60] = "ld   iyh,b"
	opsFD[0x61] = "ld   iyh,c"
	opsFD[0x62] = "ld   iyh,d"
	opsFD[0x63] = "ld   iyh,e"
	opsFD[0x64] = "ld   iyh,iyh"
	opsFD[0x65] = "ld   iyh,iyl"
	opsFD[0x66] = "ld   h,(iy+0)"
	opsFD[0x67] = "ld   iyh,a"
	opsFD[0x68] = "ld   iyl,b"
	opsFD[0x69] = "ld   iyl,c"
	opsFD[0x6A] = "ld   iyl,d"
	opsFD[0x6B] = "ld   iyl,e"
	opsFD[0x6C] = "ld   iyl,iyh"
	opsFD[0x6D] = "ld   iyl,iyl"
	opsFD[0x6E] = "ld   l,(iy+0)"
	opsFD[0x6F] = "ld   iyl,a"
	opsFD[0x70] = "ld   (iy+0),b"
	opsFD[0x71] = "ld   (iy+0),c"
	opsFD[0x72] = "ld   (iy+0),d"
	opsFD[0x73] = "ld   (iy+0),e"
	opsFD[0x74] = "ld   (iy+0),h"
	opsFD[0x75] = "ld   (iy+0),l"
	opsFD[0x77] = "ld   (iy+0),a"
	opsFD[0x7C] = "ld   a,iyh"
	opsFD[0x7D] = "ld   a,iyl"
	opsFD[0x7E] = "ld   a,(iy+0)"
	opsFD[0x84] = "add  a,iyh"
	opsFD[0x85] = "add  a,iyl"
	opsFD[0x86] = "add  a,(iy+0)"
	opsFD[0x8C] = "adc  a,iyh"
	opsFD[0x8D] = "adc  a,iyl"
	opsFD[0x8E] = "adc  a,(iy+0)"
	opsFD[0x94] = "sub  a,iyh"
	opsFD[0x95] = "sub  a,iyl"
	opsFD[0x96] = "sub  a,(iy+0)"
	opsFD[0x9C] = "sbc  a,iyh"
	opsFD[0x9D] = "sbc  a,iyl"
	opsFD[0x9E] = "sbc  a,(iy+0)"
	opsFD[0xA4] = "and  iyh"
	opsFD[0xA5] = "and  iyl"
	opsFD[0xA6] = "and  (iy+0)"
	opsFD[0xAC] = "xor  iyh"
	opsFD[0xAD] = "xor  iyl"
	opsFD[0xAE] = "xor  (iy+0)"
	opsFD[0xB4] = "or   iyh"
	opsFD[0xB5] = "or   iyl"
	opsFD[0xB6] = "or   (iy+0)"
	opsFD[0xBC] = "cp   iyh"
	opsFD[0xBD] = "cp   iyl"
	opsFD[0xBE] = "cp   (iy+0)"
	opsFD[0xE1] = "pop  iy"
	opsFD[0xE3] = "ex   (sp),iy"
	opsFD[0xE5] = "push iy"
	opsFD[0xE9] = "jp   (iy)"

	opsCB[0x00] = "rlc  b"
	opsCB[0x01] = "rlc  c"
	opsCB[0x02] = "rlc  d"
	opsCB[0x03] = "rlc  e"
	opsCB[0x04] = "rlc  h"
	opsCB[0x05] = "rlc  l"
	opsCB[0x06] = "rlc  (hl)"
	opsCB[0x07] = "rlc  a"
	opsCB[0x08] = "rrc  b"
	opsCB[0x09] = "rrc  c"
	opsCB[0x0A] = "rrc  d"
	opsCB[0x0B] = "rrc  e"
	opsCB[0x0C] = "rrc  h"
	opsCB[0x0D] = "rrc  l"
	opsCB[0x0E] = "rrc  (hl)"
	opsCB[0x0F] = "rrc  a"
	opsCB[0x10] = "rl   b"
	opsCB[0x11] = "rl   c"
	opsCB[0x12] = "rl   d"
	opsCB[0x13] = "rl   e"
	opsCB[0x14] = "rl   h"
	opsCB[0x15] = "rl   l"
	opsCB[0x16] = "rl   (hl)"
	opsCB[0x17] = "rl   a"
	opsCB[0x18] = "rr   b"
	opsCB[0x19] = "rr   c"
	opsCB[0x1A] = "rr   d"
	opsCB[0x1B] = "rr   e"
	opsCB[0x1C] = "rr   h"
	opsCB[0x1D] = "rr   l"
	opsCB[0x1E] = "rr   (hl)"
	opsCB[0x1F] = "rr   a"
	opsCB[0x20] = "sla  b"
	opsCB[0x21] = "sla  c"
	opsCB[0x22] = "sla  d"
	opsCB[0x23] = "sla  e"
	opsCB[0x24] = "sla  h"
	opsCB[0x25] = "sla  l"
	opsCB[0x26] = "sla  (hl)"
	opsCB[0x27] = "sla  a"
	opsCB[0x28] = "sra  b"
	opsCB[0x29] = "sra  c"
	opsCB[0x2A] = "sra  d"
	opsCB[0x2B] = "sra  e"
	opsCB[0x2C] = "sra  h"
	opsCB[0x2D] = "sra  l"
	opsCB[0x2E] = "sra  (hl)"
	opsCB[0x2F] = "sra  a"
	opsCB[0x30] = "sls  b"
	opsCB[0x31] = "sls  c"
	opsCB[0x32] = "sls  d"
	opsCB[0x33] = "sls  e"
	opsCB[0x34] = "sls  h"
	opsCB[0x35] = "sls  l"
	opsCB[0x36] = "sls  (hl)"
	opsCB[0x37] = "sls  a"
	opsCB[0x38] = "srl  b"
	opsCB[0x39] = "srl  c"
	opsCB[0x3A] = "srl  d"
	opsCB[0x3B] = "srl  e"
	opsCB[0x3C] = "srl  h"
	opsCB[0x3D] = "srl  l"
	opsCB[0x3E] = "srl  (hl)"
	opsCB[0x3F] = "srl  a"
	opsCB[0x40] = "bit  0,b"
	opsCB[0x41] = "bit  0,c"
	opsCB[0x42] = "bit  0,d"
	opsCB[0x43] = "bit  0,e"
	opsCB[0x44] = "bit  0,h"
	opsCB[0x45] = "bit  0,l"
	opsCB[0x46] = "bit  0,(hl)"
	opsCB[0x47] = "bit  0,a"
	opsCB[0x48] = "bit  1,b"
	opsCB[0x49] = "bit  1,c"
	opsCB[0x4A] = "bit  1,d"
	opsCB[0x4B] = "bit  1,e"
	opsCB[0x4C] = "bit  1,h"
	opsCB[0x4D] = "bit  1,l"
	opsCB[0x4E] = "bit  1,(hl)"
	opsCB[0x4F] = "bit  1,a"
	opsCB[0x50] = "bit  2,b"
	opsCB[0x51] = "bit  2,c"
	opsCB[0x52] = "bit  2,d"
	opsCB[0x53] = "bit  2,e"
	opsCB[0x54] = "bit  2,h"
	opsCB[0x55] = "bit  2,l"
	opsCB[0x56] = "bit  2,(hl)"
	opsCB[0x57] = "bit  2,a"
	opsCB[0x58] = "bit  3,b"
	opsCB[0x59] = "bit  3,c"
	opsCB[0x5A] = "bit  3,d"
	opsCB[0x5B] = "bit  3,e"
	opsCB[0x5C] = "bit  3,h"
	opsCB[0x5D] = "bit  3,l"
	opsCB[0x5E] = "bit  3,(hl)"
	opsCB[0x5F] = "bit  3,a"
	opsCB[0x60] = "bit  4,b"
	opsCB[0x61] = "bit  4,c"
	opsCB[0x62] = "bit  4,d"
	opsCB[0x63] = "bit  4,e"
	opsCB[0x64] = "bit  4,h"
	opsCB[0x65] = "bit  4,l"
	opsCB[0x66] = "bit  4,(hl)"
	opsCB[0x67] = "bit  4,a"
	opsCB[0x68] = "bit  5,b"
	opsCB[0x69] = "bit  5,c"
	opsCB[0x6A] = "bit  5,d"
	opsCB[0x6B] = "bit  5,e"
	opsCB[0x6C] = "bit  5,h"
	opsCB[0x6D] = "bit  5,l"
	opsCB[0x6E] = "bit  5,(hl)"
	opsCB[0x6F] = "bit  5,a"
	opsCB[0x70] = "bit  6,b"
	opsCB[0x71] = "bit  6,c"
	opsCB[0x72] = "bit  6,d"
	opsCB[0x73] = "bit  6,e"
	opsCB[0x74] = "bit  6,h"
	opsCB[0x75] = "bit  6,l"
	opsCB[0x76] = "bit  6,(hl)"
	opsCB[0x77] = "bit  6,a"
	opsCB[0x78] = "bit  7,b"
	opsCB[0x79] = "bit  7,c"
	opsCB[0x7A] = "bit  7,d"
	opsCB[0x7B] = "bit  7,e"
	opsCB[0x7C] = "bit  7,h"
	opsCB[0x7D] = "bit  7,l"
	opsCB[0x7E] = "bit  7,(hl)"
	opsCB[0x7F] = "bit  7,a"
	opsCB[0x80] = "res  0,b"
	opsCB[0x81] = "res  0,c"
	opsCB[0x82] = "res  0,d"
	opsCB[0x83] = "res  0,e"
	opsCB[0x84] = "res  0,h"
	opsCB[0x85] = "res  0,l"
	opsCB[0x86] = "res  0,(hl)"
	opsCB[0x87] = "res  0,a"
	opsCB[0x88] = "res  1,b"
	opsCB[0x89] = "res  1,c"
	opsCB[0x8A] = "res  1,d"
	opsCB[0x8B] = "res  1,e"
	opsCB[0x8C] = "res  1,h"
	opsCB[0x8D] = "res  1,l"
	opsCB[0x8E] = "res  1,(hl)"
	opsCB[0x8F] = "res  1,a"
	opsCB[0x90] = "res  2,b"
	opsCB[0x91] = "res  2,c"
	opsCB[0x92] = "res  2,d"
	opsCB[0x93] = "res  2,e"
	opsCB[0x94] = "res  2,h"
	opsCB[0x95] = "res  2,l"
	opsCB[0x96] = "res  2,(hl)"
	opsCB[0x97] = "res  2,a"
	opsCB[0x98] = "res  3,b"
	opsCB[0x99] = "res  3,c"
	opsCB[0x9A] = "res  3,d"
	opsCB[0x9B] = "res  3,e"
	opsCB[0x9C] = "res  3,h"
	opsCB[0x9D] = "res  3,l"
	opsCB[0x9E] = "res  3,(hl)"
	opsCB[0x9F] = "res  3,a"
	opsCB[0xA0] = "res  4,b"
	opsCB[0xA1] = "res  4,c"
	opsCB[0xA2] = "res  4,d"
	opsCB[0xA3] = "res  4,e"
	opsCB[0xA4] = "res  4,h"
	opsCB[0xA5] = "res  4,l"
	opsCB[0xA6] = "res  4,(hl)"
	opsCB[0xA7] = "res  4,a"
	opsCB[0xA8] = "res  5,b"
	opsCB[0xA9] = "res  5,c"
	opsCB[0xAA] = "res  5,d"
	opsCB[0xAB] = "res  5,e"
	opsCB[0xAC] = "res  5,h"
	opsCB[0xAD] = "res  5,l"
	opsCB[0xAE] = "res  5,(hl)"
	opsCB[0xAF] = "res  5,a"
	opsCB[0xB0] = "res  6,b"
	opsCB[0xB1] = "res  6,c"
	opsCB[0xB2] = "res  6,d"
	opsCB[0xB3] = "res  6,e"
	opsCB[0xB4] = "res  6,h"
	opsCB[0xB5] = "res  6,l"
	opsCB[0xB6] = "res  6,(hl)"
	opsCB[0xB7] = "res  6,a"
	opsCB[0xB8] = "res  7,b"
	opsCB[0xB9] = "res  7,c"
	opsCB[0xBA] = "res  7,d"
	opsCB[0xBB] = "res  7,e"
	opsCB[0xBC] = "res  7,h"
	opsCB[0xBD] = "res  7,l"
	opsCB[0xBE] = "res  7,(hl)"
	opsCB[0xBF] = "res  7,a"
	opsCB[0xC0] = "set  0,b"
	opsCB[0xC1] = "set  0,c"
	opsCB[0xC2] = "set  0,d"
	opsCB[0xC3] = "set  0,e"
	opsCB[0xC4] = "set  0,h"
	opsCB[0xC5] = "set  0,l"
	opsCB[0xC6] = "set  0,(hl)"
	opsCB[0xC7] = "set  0,a"
	opsCB[0xC8] = "set  1,b"
	opsCB[0xC9] = "set  1,c"
	opsCB[0xCA] = "set  1,d"
	opsCB[0xCB] = "set  1,e"
	opsCB[0xCC] = "set  1,h"
	opsCB[0xCD] = "set  1,l"
	opsCB[0xCE] = "set  1,(hl)"
	opsCB[0xCF] = "set  1,a"
	opsCB[0xD0] = "set  2,b"
	opsCB[0xD1] = "set  2,c"
	opsCB[0xD2] = "set  2,d"
	opsCB[0xD3] = "set  2,e"
	opsCB[0xD4] = "set  2,h"
	opsCB[0xD5] = "set  2,l"
	opsCB[0xD6] = "set  2,(hl)"
	opsCB[0xD7] = "set  2,a"
	opsCB[0xD8] = "set  3,b"
	opsCB[0xD9] = "set  3,c"
	opsCB[0xDA] = "set  3,d"
	opsCB[0xDB] = "set  3,e"
	opsCB[0xDC] = "set  3,h"
	opsCB[0xDD] = "set  3,l"
	opsCB[0xDE] = "set  3,(hl)"
	opsCB[0xDF] = "set  3,a"
	opsCB[0xE0] = "set  4,b"
	opsCB[0xE1] = "set  4,c"
	opsCB[0xE2] = "set  4,d"
	opsCB[0xE3] = "set  4,e"
	opsCB[0xE4] = "set  4,h"
	opsCB[0xE5] = "set  4,l"
	opsCB[0xE6] = "set  4,(hl)"
	opsCB[0xE7] = "set  4,a"
	opsCB[0xE8] = "set  5,b"
	opsCB[0xE9] = "set  5,c"
	opsCB[0xEA] = "set  5,d"
	opsCB[0xEB] = "set  5,e"
	opsCB[0xEC] = "set  5,h"
	opsCB[0xED] = "set  5,l"
	opsCB[0xEE] = "set  5,(hl)"
	opsCB[0xEF] = "set  5,a"
	opsCB[0xF0] = "set  6,b"
	opsCB[0xF1] = "set  6,c"
	opsCB[0xF2] = "set  6,d"
	opsCB[0xF3] = "set  6,e"
	opsCB[0xF4] = "set  6,h"
	opsCB[0xF5] = "set  6,l"
	opsCB[0xF6] = "set  6,(hl)"
	opsCB[0xF7] = "set  6,a"
	opsCB[0xF8] = "set  7,b"
	opsCB[0xF9] = "set  7,c"
	opsCB[0xFA] = "set  7,d"
	opsCB[0xFB] = "set  7,e"
	opsCB[0xFC] = "set  7,h"
	opsCB[0xFD] = "set  7,l"
	opsCB[0xFE] = "set  7,(hl)"
	opsCB[0xFF] = "set  7,a"

	opsFDCB[0x00] = "rlc (iy+0)->b"
	opsFDCB[0x01] = "rlc (iy+0)->c"
	opsFDCB[0x02] = "rlc (iy+0)->d"
	opsFDCB[0x03] = "rlc (iy+0)->e"
	opsFDCB[0x04] = "rlc (iy+0)->h"
	opsFDCB[0x05] = "rlc (iy+0)->l"
	opsFDCB[0x06] = "rlc  (iy+0)"
	opsFDCB[0x07] = "rlc (iy+0)->a"
	opsFDCB[0x08] = "rrc (iy+0)->b"
	opsFDCB[0x09] = "rrc (iy+0)->c"
	opsFDCB[0x0A] = "rrc (iy+0)->d"
	opsFDCB[0x0B] = "rrc (iy+0)->e"
	opsFDCB[0x0C] = "rrc (iy+0)->h"
	opsFDCB[0x0D] = "rrc (iy+0)->l"
	opsFDCB[0x0E] = "rrc  (iy+0)"
	opsFDCB[0x0F] = "rrc (iy+0)->a"
	opsFDCB[0x10] = "rl  (iy+0)->b"
	opsFDCB[0x11] = "rl  (iy+0)->c"
	opsFDCB[0x12] = "rl  (iy+0)->d"
	opsFDCB[0x13] = "rl  (iy+0)->e"
	opsFDCB[0x14] = "rl  (iy+0)->h"
	opsFDCB[0x15] = "rl  (iy+0)->l"
	opsFDCB[0x16] = "rl   (iy+0)"
	opsFDCB[0x17] = "rl  (iy+0)->a"
	opsFDCB[0x18] = "rr  (iy+0)->b"
	opsFDCB[0x19] = "rr  (iy+0)->c"
	opsFDCB[0x1A] = "rr  (iy+0)->d"
	opsFDCB[0x1B] = "rr  (iy+0)->e"
	opsFDCB[0x1C] = "rr  (iy+0)->h"
	opsFDCB[0x1D] = "rr  (iy+0)->l"
	opsFDCB[0x1E] = "rr   (iy+0)"
	opsFDCB[0x1F] = "rr  (iy+0)->a"
	opsFDCB[0x20] = "sla (iy+0)->b"
	opsFDCB[0x21] = "sla (iy+0)->c"
	opsFDCB[0x22] = "sla (iy+0)->d"
	opsFDCB[0x23] = "sla (iy+0)->e"
	opsFDCB[0x24] = "sla (iy+0)->h"
	opsFDCB[0x25] = "sla (iy+0)->l"
	opsFDCB[0x26] = "sla  (iy+0)"
	opsFDCB[0x27] = "sla (iy+0)->a"
	opsFDCB[0x28] = "sra (iy+0)->b"
	opsFDCB[0x29] = "sra (iy+0)->c"
	opsFDCB[0x2A] = "sra (iy+0)->d"
	opsFDCB[0x2B] = "sra (iy+0)->e"
	opsFDCB[0x2C] = "sra (iy+0)->h"
	opsFDCB[0x2D] = "sra (iy+0)->l"
	opsFDCB[0x2E] = "sra  (iy+0)"
	opsFDCB[0x2F] = "sra (iy+0)->a"
	opsFDCB[0x30] = "sls (iy+0)->b"
	opsFDCB[0x31] = "sls (iy+0)->c"
	opsFDCB[0x32] = "sls (iy+0)->d"
	opsFDCB[0x33] = "sls (iy+0)->e"
	opsFDCB[0x34] = "sls (iy+0)->h"
	opsFDCB[0x35] = "sls (iy+0)->l"
	opsFDCB[0x36] = "sls  (iy+0)"
	opsFDCB[0x37] = "sls (iy+0)->a"
	opsFDCB[0x38] = "srl (iy+0)->b"
	opsFDCB[0x39] = "srl (iy+0)->c"
	opsFDCB[0x3A] = "srl (iy+0)->d"
	opsFDCB[0x3B] = "srl (iy+0)->e"
	opsFDCB[0x3C] = "srl (iy+0)->h"
	opsFDCB[0x3D] = "srl (iy+0)->l"
	opsFDCB[0x3E] = "srl  (iy+0)"
	opsFDCB[0x3F] = "srl (iy+0)->a"
	opsFDCB[0x40] = "bit 0,(iy+0)->b"
	opsFDCB[0x41] = "bit 0,(iy+0)->c"
	opsFDCB[0x42] = "bit 0,(iy+0)->d"
	opsFDCB[0x43] = "bit 0,(iy+0)->e"
	opsFDCB[0x44] = "bit 0,(iy+0)->h"
	opsFDCB[0x45] = "bit 0,(iy+0)->l"
	opsFDCB[0x46] = "bit  0,(iy+0)"
	opsFDCB[0x47] = "bit 0,(iy+0)->a"
	opsFDCB[0x48] = "bit 1,(iy+0)->b"
	opsFDCB[0x49] = "bit 1,(iy+0)->c"
	opsFDCB[0x4A] = "bit 1,(iy+0)->d"
	opsFDCB[0x4B] = "bit 1,(iy+0)->e"
	opsFDCB[0x4C] = "bit 1,(iy+0)->h"
	opsFDCB[0x4D] = "bit 1,(iy+0)->l"
	opsFDCB[0x4E] = "bit  1,(iy+0)"
	opsFDCB[0x4F] = "bit 1,(iy+0)->a"
	opsFDCB[0x50] = "bit 2,(iy+0)->b"
	opsFDCB[0x51] = "bit 2,(iy+0)->c"
	opsFDCB[0x52] = "bit 2,(iy+0)->d"
	opsFDCB[0x53] = "bit 2,(iy+0)->e"
	opsFDCB[0x54] = "bit 2,(iy+0)->h"
	opsFDCB[0x55] = "bit 2,(iy+0)->l"
	opsFDCB[0x56] = "bit  2,(iy+0)"
	opsFDCB[0x57] = "bit 2,(iy+0)->a"
	opsFDCB[0x58] = "bit 3,(iy+0)->b"
	opsFDCB[0x59] = "bit 3,(iy+0)->c"
	opsFDCB[0x5A] = "bit 3,(iy+0)->d"
	opsFDCB[0x5B] = "bit 3,(iy+0)->e"
	opsFDCB[0x5C] = "bit 3,(iy+0)->h"
	opsFDCB[0x5D] = "bit 3,(iy+0)->l"
	opsFDCB[0x5E] = "bit  3,(iy+0)"
	opsFDCB[0x5F] = "bit 3,(iy+0)->a"
	opsFDCB[0x60] = "bit 4,(iy+0)->b"
	opsFDCB[0x61] = "bit 4,(iy+0)->c"
	opsFDCB[0x62] = "bit 4,(iy+0)->d"
	opsFDCB[0x63] = "bit 4,(iy+0)->e"
	opsFDCB[0x64] = "bit 4,(iy+0)->h"
	opsFDCB[0x65] = "bit 4,(iy+0)->l"
	opsFDCB[0x66] = "bit  4,(iy+0)"
	opsFDCB[0x67] = "bit 4,(iy+0)->a"
	opsFDCB[0x68] = "bit 5,(iy+0)->b"
	opsFDCB[0x69] = "bit 5,(iy+0)->c"
	opsFDCB[0x6A] = "bit 5,(iy+0)->d"
	opsFDCB[0x6B] = "bit 5,(iy+0)->e"
	opsFDCB[0x6C] = "bit 5,(iy+0)->h"
	opsFDCB[0x6D] = "bit 5,(iy+0)->l"
	opsFDCB[0x6E] = "bit  5,(iy+0)"
	opsFDCB[0x6F] = "bit 5,(iy+0)->a"
	opsFDCB[0x70] = "bit 6,(iy+0)->b"
	opsFDCB[0x71] = "bit 6,(iy+0)->c"
	opsFDCB[0x72] = "bit 6,(iy+0)->d"
	opsFDCB[0x73] = "bit 6,(iy+0)->e"
	opsFDCB[0x74] = "bit 6,(iy+0)->h"
	opsFDCB[0x75] = "bit 6,(iy+0)->l"
	opsFDCB[0x76] = "bit  6,(iy+0)"
	opsFDCB[0x77] = "bit 6,(iy+0)->a"
	opsFDCB[0x78] = "bit 7,(iy+0)->b"
	opsFDCB[0x79] = "bit 7,(iy+0)->c"
	opsFDCB[0x7A] = "bit 7,(iy+0)->d"
	opsFDCB[0x7B] = "bit 7,(iy+0)->e"
	opsFDCB[0x7C] = "bit 7,(iy+0)->h"
	opsFDCB[0x7D] = "bit 7,(iy+0)->l"
	opsFDCB[0x7E] = "bit  7,(iy+0)"
	opsFDCB[0x7F] = "bit 7,(iy+0)->a"
	opsFDCB[0x80] = "res 0,(iy+0)->b"
	opsFDCB[0x81] = "res 0,(iy+0)->c"
	opsFDCB[0x82] = "res 0,(iy+0)->d"
	opsFDCB[0x83] = "res 0,(iy+0)->e"
	opsFDCB[0x84] = "res 0,(iy+0)->h"
	opsFDCB[0x85] = "res 0,(iy+0)->l"
	opsFDCB[0x86] = "res  0,(iy+0)"
	opsFDCB[0x87] = "res 0,(iy+0)->a"
	opsFDCB[0x88] = "res 1,(iy+0)->b"
	opsFDCB[0x89] = "res 1,(iy+0)->c"
	opsFDCB[0x8A] = "res 1,(iy+0)->d"
	opsFDCB[0x8B] = "res 1,(iy+0)->e"
	opsFDCB[0x8C] = "res 1,(iy+0)->h"
	opsFDCB[0x8D] = "res 1,(iy+0)->l"
	opsFDCB[0x8E] = "res  1,(iy+0)"
	opsFDCB[0x8F] = "res 1,(iy+0)->a"
	opsFDCB[0x90] = "res 2,(iy+0)->b"
	opsFDCB[0x91] = "res 2,(iy+0)->c"
	opsFDCB[0x92] = "res 2,(iy+0)->d"
	opsFDCB[0x93] = "res 2,(iy+0)->e"
	opsFDCB[0x94] = "res 2,(iy+0)->h"
	opsFDCB[0x95] = "res 2,(iy+0)->l"
	opsFDCB[0x96] = "res  2,(iy+0)"
	opsFDCB[0x97] = "res 2,(iy+0)->a"
	opsFDCB[0x98] = "res 3,(iy+0)->b"
	opsFDCB[0x99] = "res 3,(iy+0)->c"
	opsFDCB[0x9A] = "res 3,(iy+0)->d"
	opsFDCB[0x9B] = "res 3,(iy+0)->e"
	opsFDCB[0x9C] = "res 3,(iy+0)->h"
	opsFDCB[0x9D] = "res 3,(iy+0)->l"
	opsFDCB[0x9E] = "res  3,(iy+0)"
	opsFDCB[0x9F] = "res 3,(iy+0)->a"
	opsFDCB[0xA0] = "res 4,(iy+0)->b"
	opsFDCB[0xA1] = "res 4,(iy+0)->c"
	opsFDCB[0xA2] = "res 4,(iy+0)->d"
	opsFDCB[0xA3] = "res 4,(iy+0)->e"
	opsFDCB[0xA4] = "res 4,(iy+0)->h"
	opsFDCB[0xA5] = "res 4,(iy+0)->l"
	opsFDCB[0xA6] = "res  4,(iy+0)"
	opsFDCB[0xA7] = "res 4,(iy+0)->a"
	opsFDCB[0xA8] = "res 5,(iy+0)->b"
	opsFDCB[0xA9] = "res 5,(iy+0)->c"
	opsFDCB[0xAA] = "res 5,(iy+0)->d"
	opsFDCB[0xAB] = "res 5,(iy+0)->e"
	opsFDCB[0xAC] = "res 5,(iy+0)->h"
	opsFDCB[0xAD] = "res 5,(iy+0)->l"
	opsFDCB[0xAE] = "res  5,(iy+0)"
	opsFDCB[0xAF] = "res 5,(iy+0)->a"
	opsFDCB[0xB0] = "res 6,(iy+0)->b"
	opsFDCB[0xB1] = "res 6,(iy+0)->c"
	opsFDCB[0xB2] = "res 6,(iy+0)->d"
	opsFDCB[0xB3] = "res 6,(iy+0)->e"
	opsFDCB[0xB4] = "res 6,(iy+0)->h"
	opsFDCB[0xB5] = "res 6,(iy+0)->l"
	opsFDCB[0xB6] = "res  6,(iy+0)"
	opsFDCB[0xB7] = "res 6,(iy+0)->a"
	opsFDCB[0xB8] = "res 7,(iy+0)->b"
	opsFDCB[0xB9] = "res 7,(iy+0)->c"
	opsFDCB[0xBA] = "res 7,(iy+0)->d"
	opsFDCB[0xBB] = "res 7,(iy+0)->e"
	opsFDCB[0xBC] = "res 7,(iy+0)->h"
	opsFDCB[0xBD] = "res 7,(iy+0)->l"
	opsFDCB[0xBE] = "res  7,(iy+0)"
	opsFDCB[0xBF] = "res 7,(iy+0)->a"
	opsFDCB[0xC0] = "set 0,(iy+0)->b"
	opsFDCB[0xC1] = "set 0,(iy+0)->c"
	opsFDCB[0xC2] = "set 0,(iy+0)->d"
	opsFDCB[0xC3] = "set 0,(iy+0)->e"
	opsFDCB[0xC4] = "set 0,(iy+0)->h"
	opsFDCB[0xC5] = "set 0,(iy+0)->l"
	opsFDCB[0xC6] = "set  0,(iy+0)"
	opsFDCB[0xC7] = "set 0,(iy+0)->a"
	opsFDCB[0xC8] = "set 1,(iy+0)->b"
	opsFDCB[0xC9] = "set 1,(iy+0)->c"
	opsFDCB[0xCA] = "set 1,(iy+0)->d"
	opsFDCB[0xCB] = "set 1,(iy+0)->e"
	opsFDCB[0xCC] = "set 1,(iy+0)->h"
	opsFDCB[0xCD] = "set 1,(iy+0)->l"
	opsFDCB[0xCE] = "set  1,(iy+0)"
	opsFDCB[0xCF] = "set 1,(iy+0)->a"
	opsFDCB[0xD0] = "set 2,(iy+0)->b"
	opsFDCB[0xD1] = "set 2,(iy+0)->c"
	opsFDCB[0xD2] = "set 2,(iy+0)->d"
	opsFDCB[0xD3] = "set 2,(iy+0)->e"
	opsFDCB[0xD4] = "set 2,(iy+0)->h"
	opsFDCB[0xD5] = "set 2,(iy+0)->l"
	opsFDCB[0xD6] = "set  2,(iy+0)"
	opsFDCB[0xD7] = "set 2,(iy+0)->a"
	opsFDCB[0xD8] = "set 3,(iy+0)->b"
	opsFDCB[0xD9] = "set 3,(iy+0)->c"
	opsFDCB[0xDA] = "set 3,(iy+0)->d"
	opsFDCB[0xDB] = "set 3,(iy+0)->e"
	opsFDCB[0xDC] = "set 3,(iy+0)->h"
	opsFDCB[0xDD] = "set 3,(iy+0)->l"
	opsFDCB[0xDE] = "set  3,(iy+0)"
	opsFDCB[0xDF] = "set 3,(iy+0)->a"
	opsFDCB[0xE0] = "set 4,(iy+0)->b"
	opsFDCB[0xE1] = "set 4,(iy+0)->c"
	opsFDCB[0xE2] = "set 4,(iy+0)->d"
	opsFDCB[0xE3] = "set 4,(iy+0)->e"
	opsFDCB[0xE4] = "set 4,(iy+0)->h"
	opsFDCB[0xE5] = "set 4,(iy+0)->l"
	opsFDCB[0xE6] = "set  4,(iy+0)"
	opsFDCB[0xE7] = "set 4,(iy+0)->a"
	opsFDCB[0xE8] = "set 5,(iy+0)->b"
	opsFDCB[0xE9] = "set 5,(iy+0)->c"
	opsFDCB[0xEA] = "set 5,(iy+0)->d"
	opsFDCB[0xEB] = "set 5,(iy+0)->e"
	opsFDCB[0xEC] = "set 5,(iy+0)->h"
	opsFDCB[0xED] = "set 5,(iy+0)->l"
	opsFDCB[0xEE] = "set  5,(iy+0)"
	opsFDCB[0xEF] = "set 5,(iy+0)->a"
	opsFDCB[0xF0] = "set 6,(iy+0)->b"
	opsFDCB[0xF1] = "set 6,(iy+0)->c"
	opsFDCB[0xF2] = "set 6,(iy+0)->d"
	opsFDCB[0xF3] = "set 6,(iy+0)->e"
	opsFDCB[0xF4] = "set 6,(iy+0)->h"
	opsFDCB[0xF5] = "set 6,(iy+0)->l"
	opsFDCB[0xF6] = "set  6,(iy+0)"
	opsFDCB[0xF7] = "set 6,(iy+0)->a"
	opsFDCB[0xF8] = "set 7,(iy+0)->b"
	opsFDCB[0xF9] = "set 7,(iy+0)->c"
	opsFDCB[0xFA] = "set 7,(iy+0)->d"
	opsFDCB[0xFB] = "set 7,(iy+0)->e"
	opsFDCB[0xFC] = "set 7,(iy+0)->h"
	opsFDCB[0xFD] = "set 7,(iy+0)->l"
	opsFDCB[0xFE] = "set  7,(iy+0)"
	opsFDCB[0xFF] = "set 7,(iy+0)->a"

	opsDDCB[0x00] = "rlc (ix+0)->b"
	opsDDCB[0x01] = "rlc (ix+0)->c"
	opsDDCB[0x02] = "rlc (ix+0)->d"
	opsDDCB[0x03] = "rlc (ix+0)->e"
	opsDDCB[0x04] = "rlc (ix+0)->h"
	opsDDCB[0x05] = "rlc (ix+0)->l"
	opsDDCB[0x06] = "rlc  (ix+0)"
	opsDDCB[0x07] = "rlc (ix+0)->a"
	opsDDCB[0x08] = "rrc (ix+0)->b"
	opsDDCB[0x09] = "rrc (ix+0)->c"
	opsDDCB[0x0A] = "rrc (ix+0)->d"
	opsDDCB[0x0B] = "rrc (ix+0)->e"
	opsDDCB[0x0C] = "rrc (ix+0)->h"
	opsDDCB[0x0D] = "rrc (ix+0)->l"
	opsDDCB[0x0E] = "rrc  (ix+0)"
	opsDDCB[0x0F] = "rrc (ix+0)->a"
	opsDDCB[0x10] = "rl  (ix+0)->b"
	opsDDCB[0x11] = "rl  (ix+0)->c"
	opsDDCB[0x12] = "rl  (ix+0)->d"
	opsDDCB[0x13] = "rl  (ix+0)->e"
	opsDDCB[0x14] = "rl  (ix+0)->h"
	opsDDCB[0x15] = "rl  (ix+0)->l"
	opsDDCB[0x16] = "rl   (ix+0)"
	opsDDCB[0x17] = "rl  (ix+0)->a"
	opsDDCB[0x18] = "rr  (ix+0)->b"
	opsDDCB[0x19] = "rr  (ix+0)->c"
	opsDDCB[0x1A] = "rr  (ix+0)->d"
	opsDDCB[0x1B] = "rr  (ix+0)->e"
	opsDDCB[0x1C] = "rr  (ix+0)->h"
	opsDDCB[0x1D] = "rr  (ix+0)->l"
	opsDDCB[0x1E] = "rr   (ix+0)"
	opsDDCB[0x1F] = "rr  (ix+0)->a"
	opsDDCB[0x20] = "sla (ix+0)->b"
	opsDDCB[0x21] = "sla (ix+0)->c"
	opsDDCB[0x22] = "sla (ix+0)->d"
	opsDDCB[0x23] = "sla (ix+0)->e"
	opsDDCB[0x24] = "sla (ix+0)->h"
	opsDDCB[0x25] = "sla (ix+0)->l"
	opsDDCB[0x26] = "sla  (ix+0)"
	opsDDCB[0x27] = "sla (ix+0)->a"
	opsDDCB[0x28] = "sra (ix+0)->b"
	opsDDCB[0x29] = "sra (ix+0)->c"
	opsDDCB[0x2A] = "sra (ix+0)->d"
	opsDDCB[0x2B] = "sra (ix+0)->e"
	opsDDCB[0x2C] = "sra (ix+0)->h"
	opsDDCB[0x2D] = "sra (ix+0)->l"
	opsDDCB[0x2E] = "sra  (ix+0)"
	opsDDCB[0x2F] = "sra (ix+0)->a"
	opsDDCB[0x30] = "sls (ix+0)->b"
	opsDDCB[0x31] = "sls (ix+0)->c"
	opsDDCB[0x32] = "sls (ix+0)->d"
	opsDDCB[0x33] = "sls (ix+0)->e"
	opsDDCB[0x34] = "sls (ix+0)->h"
	opsDDCB[0x35] = "sls (ix+0)->l"
	opsDDCB[0x36] = "sls  (ix+0)"
	opsDDCB[0x37] = "sls (ix+0)->a"
	opsDDCB[0x38] = "srl (ix+0)->b"
	opsDDCB[0x39] = "srl (ix+0)->c"
	opsDDCB[0x3A] = "srl (ix+0)->d"
	opsDDCB[0x3B] = "srl (ix+0)->e"
	opsDDCB[0x3C] = "srl (ix+0)->h"
	opsDDCB[0x3D] = "srl (ix+0)->l"
	opsDDCB[0x3E] = "srl  (ix+0)"
	opsDDCB[0x3F] = "srl (ix+0)->a"
	opsDDCB[0x40] = "bit 0,(ix+0)->b"
	opsDDCB[0x41] = "bit 0,(ix+0)->c"
	opsDDCB[0x42] = "bit 0,(ix+0)->d"
	opsDDCB[0x43] = "bit 0,(ix+0)->e"
	opsDDCB[0x44] = "bit 0,(ix+0)->h"
	opsDDCB[0x45] = "bit 0,(ix+0)->l"
	opsDDCB[0x46] = "bit  0,(ix+0)"
	opsDDCB[0x47] = "bit 0,(ix+0)->a"
	opsDDCB[0x48] = "bit 1,(ix+0)->b"
	opsDDCB[0x49] = "bit 1,(ix+0)->c"
	opsDDCB[0x4A] = "bit 1,(ix+0)->d"
	opsDDCB[0x4B] = "bit 1,(ix+0)->e"
	opsDDCB[0x4C] = "bit 1,(ix+0)->h"
	opsDDCB[0x4D] = "bit 1,(ix+0)->l"
	opsDDCB[0x4E] = "bit  1,(ix+0)"
	opsDDCB[0x4F] = "bit 1,(ix+0)->a"
	opsDDCB[0x50] = "bit 2,(ix+0)->b"
	opsDDCB[0x51] = "bit 2,(ix+0)->c"
	opsDDCB[0x52] = "bit 2,(ix+0)->d"
	opsDDCB[0x53] = "bit 2,(ix+0)->e"
	opsDDCB[0x54] = "bit 2,(ix+0)->h"
	opsDDCB[0x55] = "bit 2,(ix+0)->l"
	opsDDCB[0x56] = "bit  2,(ix+0)"
	opsDDCB[0x57] = "bit 2,(ix+0)->a"
	opsDDCB[0x58] = "bit 3,(ix+0)->b"
	opsDDCB[0x59] = "bit 3,(ix+0)->c"
	opsDDCB[0x5A] = "bit 3,(ix+0)->d"
	opsDDCB[0x5B] = "bit 3,(ix+0)->e"
	opsDDCB[0x5C] = "bit 3,(ix+0)->h"
	opsDDCB[0x5D] = "bit 3,(ix+0)->l"
	opsDDCB[0x5E] = "bit  3,(ix+0)"
	opsDDCB[0x5F] = "bit 3,(ix+0)->a"
	opsDDCB[0x60] = "bit 4,(ix+0)->b"
	opsDDCB[0x61] = "bit 4,(ix+0)->c"
	opsDDCB[0x62] = "bit 4,(ix+0)->d"
	opsDDCB[0x63] = "bit 4,(ix+0)->e"
	opsDDCB[0x64] = "bit 4,(ix+0)->h"
	opsDDCB[0x65] = "bit 4,(ix+0)->l"
	opsDDCB[0x66] = "bit  4,(ix+0)"
	opsDDCB[0x67] = "bit 4,(ix+0)->a"
	opsDDCB[0x68] = "bit 5,(ix+0)->b"
	opsDDCB[0x69] = "bit 5,(ix+0)->c"
	opsDDCB[0x6A] = "bit 5,(ix+0)->d"
	opsDDCB[0x6B] = "bit 5,(ix+0)->e"
	opsDDCB[0x6C] = "bit 5,(ix+0)->h"
	opsDDCB[0x6D] = "bit 5,(ix+0)->l"
	opsDDCB[0x6E] = "bit  5,(ix+0)"
	opsDDCB[0x6F] = "bit 5,(ix+0)->a"
	opsDDCB[0x70] = "bit 6,(ix+0)->b"
	opsDDCB[0x71] = "bit 6,(ix+0)->c"
	opsDDCB[0x72] = "bit 6,(ix+0)->d"
	opsDDCB[0x73] = "bit 6,(ix+0)->e"
	opsDDCB[0x74] = "bit 6,(ix+0)->h"
	opsDDCB[0x75] = "bit 6,(ix+0)->l"
	opsDDCB[0x76] = "bit  6,(ix+0)"
	opsDDCB[0x77] = "bit 6,(ix+0)->a"
	opsDDCB[0x78] = "bit 7,(ix+0)->b"
	opsDDCB[0x79] = "bit 7,(ix+0)->c"
	opsDDCB[0x7A] = "bit 7,(ix+0)->d"
	opsDDCB[0x7B] = "bit 7,(ix+0)->e"
	opsDDCB[0x7C] = "bit 7,(ix+0)->h"
	opsDDCB[0x7D] = "bit 7,(ix+0)->l"
	opsDDCB[0x7E] = "bit  7,(ix+0)"
	opsDDCB[0x7F] = "bit 7,(ix+0)->a"
	opsDDCB[0x80] = "res 0,(ix+0)->b"
	opsDDCB[0x81] = "res 0,(ix+0)->c"
	opsDDCB[0x82] = "res 0,(ix+0)->d"
	opsDDCB[0x83] = "res 0,(ix+0)->e"
	opsDDCB[0x84] = "res 0,(ix+0)->h"
	opsDDCB[0x85] = "res 0,(ix+0)->l"
	opsDDCB[0x86] = "res  0,(ix+0)"
	opsDDCB[0x87] = "res 0,(ix+0)->a"
	opsDDCB[0x88] = "res 1,(ix+0)->b"
	opsDDCB[0x89] = "res 1,(ix+0)->c"
	opsDDCB[0x8A] = "res 1,(ix+0)->d"
	opsDDCB[0x8B] = "res 1,(ix+0)->e"
	opsDDCB[0x8C] = "res 1,(ix+0)->h"
	opsDDCB[0x8D] = "res 1,(ix+0)->l"
	opsDDCB[0x8E] = "res  1,(ix+0)"
	opsDDCB[0x8F] = "res 1,(ix+0)->a"
	opsDDCB[0x90] = "res 2,(ix+0)->b"
	opsDDCB[0x91] = "res 2,(ix+0)->c"
	opsDDCB[0x92] = "res 2,(ix+0)->d"
	opsDDCB[0x93] = "res 2,(ix+0)->e"
	opsDDCB[0x94] = "res 2,(ix+0)->h"
	opsDDCB[0x95] = "res 2,(ix+0)->l"
	opsDDCB[0x96] = "res  2,(ix+0)"
	opsDDCB[0x97] = "res 2,(ix+0)->a"
	opsDDCB[0x98] = "res 3,(ix+0)->b"
	opsDDCB[0x99] = "res 3,(ix+0)->c"
	opsDDCB[0x9A] = "res 3,(ix+0)->d"
	opsDDCB[0x9B] = "res 3,(ix+0)->e"
	opsDDCB[0x9C] = "res 3,(ix+0)->h"
	opsDDCB[0x9D] = "res 3,(ix+0)->l"
	opsDDCB[0x9E] = "res  3,(ix+0)"
	opsDDCB[0x9F] = "res 3,(ix+0)->a"
	opsDDCB[0xA0] = "res 4,(ix+0)->b"
	opsDDCB[0xA1] = "res 4,(ix+0)->c"
	opsDDCB[0xA2] = "res 4,(ix+0)->d"
	opsDDCB[0xA3] = "res 4,(ix+0)->e"
	opsDDCB[0xA4] = "res 4,(ix+0)->h"
	opsDDCB[0xA5] = "res 4,(ix+0)->l"
	opsDDCB[0xA6] = "res  4,(ix+0)"
	opsDDCB[0xA7] = "res 4,(ix+0)->a"
	opsDDCB[0xA8] = "res 5,(ix+0)->b"
	opsDDCB[0xA9] = "res 5,(ix+0)->c"
	opsDDCB[0xAA] = "res 5,(ix+0)->d"
	opsDDCB[0xAB] = "res 5,(ix+0)->e"
	opsDDCB[0xAC] = "res 5,(ix+0)->h"
	opsDDCB[0xAD] = "res 5,(ix+0)->l"
	opsDDCB[0xAE] = "res  5,(ix+0)"
	opsDDCB[0xAF] = "res 5,(ix+0)->a"
	opsDDCB[0xB0] = "res 6,(ix+0)->b"
	opsDDCB[0xB1] = "res 6,(ix+0)->c"
	opsDDCB[0xB2] = "res 6,(ix+0)->d"
	opsDDCB[0xB3] = "res 6,(ix+0)->e"
	opsDDCB[0xB4] = "res 6,(ix+0)->h"
	opsDDCB[0xB5] = "res 6,(ix+0)->l"
	opsDDCB[0xB6] = "res  6,(ix+0)"
	opsDDCB[0xB7] = "res 6,(ix+0)->a"
	opsDDCB[0xB8] = "res 7,(ix+0)->b"
	opsDDCB[0xB9] = "res 7,(ix+0)->c"
	opsDDCB[0xBA] = "res 7,(ix+0)->d"
	opsDDCB[0xBB] = "res 7,(ix+0)->e"
	opsDDCB[0xBC] = "res 7,(ix+0)->h"
	opsDDCB[0xBD] = "res 7,(ix+0)->l"
	opsDDCB[0xBE] = "res  7,(ix+0)"
	opsDDCB[0xBF] = "res 7,(ix+0)->a"
	opsDDCB[0xC0] = "set 0,(ix+0)->b"
	opsDDCB[0xC1] = "set 0,(ix+0)->c"
	opsDDCB[0xC2] = "set 0,(ix+0)->d"
	opsDDCB[0xC3] = "set 0,(ix+0)->e"
	opsDDCB[0xC4] = "set 0,(ix+0)->h"
	opsDDCB[0xC5] = "set 0,(ix+0)->l"
	opsDDCB[0xC6] = "set  0,(ix+0)"
	opsDDCB[0xC7] = "set 0,(ix+0)->a"
	opsDDCB[0xC8] = "set 1,(ix+0)->b"
	opsDDCB[0xC9] = "set 1,(ix+0)->c"
	opsDDCB[0xCA] = "set 1,(ix+0)->d"
	opsDDCB[0xCB] = "set 1,(ix+0)->e"
	opsDDCB[0xCC] = "set 1,(ix+0)->h"
	opsDDCB[0xCD] = "set 1,(ix+0)->l"
	opsDDCB[0xCE] = "set  1,(ix+0)"
	opsDDCB[0xCF] = "set 1,(ix+0)->a"
	opsDDCB[0xD0] = "set 2,(ix+0)->b"
	opsDDCB[0xD1] = "set 2,(ix+0)->c"
	opsDDCB[0xD2] = "set 2,(ix+0)->d"
	opsDDCB[0xD3] = "set 2,(ix+0)->e"
	opsDDCB[0xD4] = "set 2,(ix+0)->h"
	opsDDCB[0xD5] = "set 2,(ix+0)->l"
	opsDDCB[0xD6] = "set  2,(ix+0)"
	opsDDCB[0xD7] = "set 2,(ix+0)->a"
	opsDDCB[0xD8] = "set 3,(ix+0)->b"
	opsDDCB[0xD9] = "set 3,(ix+0)->c"
	opsDDCB[0xDA] = "set 3,(ix+0)->d"
	opsDDCB[0xDB] = "set 3,(ix+0)->e"
	opsDDCB[0xDC] = "set 3,(ix+0)->h"
	opsDDCB[0xDD] = "set 3,(ix+0)->l"
	opsDDCB[0xDE] = "set  3,(ix+0)"
	opsDDCB[0xDF] = "set 3,(ix+0)->a"
	opsDDCB[0xE0] = "set 4,(ix+0)->b"
	opsDDCB[0xE1] = "set 4,(ix+0)->c"
	opsDDCB[0xE2] = "set 4,(ix+0)->d"
	opsDDCB[0xE3] = "set 4,(ix+0)->e"
	opsDDCB[0xE4] = "set 4,(ix+0)->h"
	opsDDCB[0xE5] = "set 4,(ix+0)->l"
	opsDDCB[0xE6] = "set  4,(ix+0)"
	opsDDCB[0xE7] = "set 4,(ix+0)->a"
	opsDDCB[0xE8] = "set 5,(ix+0)->b"
	opsDDCB[0xE9] = "set 5,(ix+0)->c"
	opsDDCB[0xEA] = "set 5,(ix+0)->d"
	opsDDCB[0xEB] = "set 5,(ix+0)->e"
	opsDDCB[0xEC] = "set 5,(ix+0)->h"
	opsDDCB[0xED] = "set 5,(ix+0)->l"
	opsDDCB[0xEE] = "set  5,(ix+0)"
	opsDDCB[0xEF] = "set 5,(ix+0)->a"
	opsDDCB[0xF0] = "set 6,(ix+0)->b"
	opsDDCB[0xF1] = "set 6,(ix+0)->c"
	opsDDCB[0xF2] = "set 6,(ix+0)->d"
	opsDDCB[0xF3] = "set 6,(ix+0)->e"
	opsDDCB[0xF4] = "set 6,(ix+0)->h"
	opsDDCB[0xF5] = "set 6,(ix+0)->l"
	opsDDCB[0xF6] = "set  6,(ix+0)"
	opsDDCB[0xF7] = "set 6,(ix+0)->a"
	opsDDCB[0xF8] = "set 7,(ix+0)->b"
	opsDDCB[0xF9] = "set 7,(ix+0)->c"
	opsDDCB[0xFA] = "set 7,(ix+0)->d"
	opsDDCB[0xFB] = "set 7,(ix+0)->e"
	opsDDCB[0xFC] = "set 7,(ix+0)->h"
	opsDDCB[0xFD] = "set 7,(ix+0)->l"
	opsDDCB[0xFE] = "set  7,(ix+0)"
	opsDDCB[0xFF] = "set 7,(ix+0)->a"

	opsED[0x00] = "mos_quit"
	opsED[0x01] = "mos_cli"
	opsED[0x02] = "mos_byte"
	opsED[0x03] = "mos_word"
	opsED[0x04] = "mos_wrch"
	opsED[0x05] = "mos_rdch"
	opsED[0x06] = "mos_file"
	opsED[0x07] = "mos_args"
	opsED[0x08] = "mos_bget"
	opsED[0x09] = "mos_bput"
	opsED[0x0A] = "mos_gbpb"
	opsED[0x0B] = "mos_find"
	opsED[0x0C] = "mos_ff0c"
	opsED[0x0D] = "mos_ff0d"
	opsED[0x0E] = "mos_ff0e"
	opsED[0x0F] = "mos_ff0f"
	opsED[0x40] = "in   b,(c)"
	opsED[0x41] = "out  (c),b"
	opsED[0x42] = "sbc  hl,bc"
	opsED[0x43] = "ld   ($nn),bc"
	opsED[0x44] = "neg"
	opsED[0x45] = "retn"
	opsED[0x46] = "im   0"
	opsED[0x47] = "ld   i,a"
	opsED[0x48] = "in   c,(c)"
	opsED[0x49] = "out  (c),c"
	opsED[0x4A] = "adc  hl,bc"
	opsED[0x4B] = "ld   bc,($nn)"
	opsED[0x4C] = "[neg]"
	opsED[0x4D] = "reti"
	opsED[0x4E] = "[im0]"
	opsED[0x4F] = "ld   r,a"
	opsED[0x50] = "in   d,(c)"
	opsED[0x51] = "out  (c),d"
	opsED[0x52] = "sbc  hl,de"
	opsED[0x53] = "ld   ($nn),de"
	opsED[0x54] = "[neg]"
	opsED[0x55] = "[retn]"
	opsED[0x56] = "im   1"
	opsED[0x57] = "ld   a,i"
	opsED[0x58] = "in   e,(c)"
	opsED[0x59] = "out  (c),e"
	opsED[0x5A] = "adc  hl,de"
	opsED[0x5B] = "ld   de,($nn)"
	opsED[0x5C] = "[neg]"
	opsED[0x5D] = "[reti]"
	opsED[0x5E] = "im   2"
	opsED[0x5F] = "ld   a,r"
	opsED[0x60] = "in   h,(c)"
	opsED[0x61] = "out  (c),h"
	opsED[0x62] = "sbc  hl,hl"
	opsED[0x63] = "ld   ($nn),hl"
	opsED[0x64] = "[neg]"
	opsED[0x65] = "[retn]"
	opsED[0x66] = "[im0]"
	opsED[0x67] = "rrd"
	opsED[0x68] = "in   l,(c)"
	opsED[0x69] = "out  (c),l"
	opsED[0x6A] = "adc  hl,hl"
	opsED[0x6B] = "ld   hl,($nn)"
	opsED[0x6C] = "[neg]"
	opsED[0x6D] = "[reti]"
	opsED[0x6E] = "[im0]"
	opsED[0x6F] = "rld"
	opsED[0x70] = "in   f,(c)"
	opsED[0x71] = "out  (c),f"
	opsED[0x72] = "sbc  hl,sp"
	opsED[0x73] = "ld   ($nn),sp"
	opsED[0x74] = "[neg]"
	opsED[0x75] = "[retn]"
	opsED[0x76] = "[im1]"
	opsED[0x77] = "[ld i,i?]"
	opsED[0x78] = "in   a,(c)"
	opsED[0x79] = "out  (c),a"
	opsED[0x7A] = "adc  hl,sp"
	opsED[0x7B] = "ld   sp,($nn)"
	opsED[0x7C] = "[neg]"
	opsED[0x7D] = "[reti]"
	opsED[0x7E] = "[im2]"
	opsED[0x7F] = "[ld r,r?]"
	opsED[0xA0] = "ldi"
	opsED[0xA1] = "cpi"
	opsED[0xA2] = "ini"
	opsED[0xA3] = "oti"
	opsED[0xA8] = "ldd"
	opsED[0xA9] = "cpd"
	opsED[0xAA] = "ind"
	opsED[0xAB] = "otd"
	opsED[0xB0] = "ldir"
	opsED[0xB1] = "cpir"
	opsED[0xB2] = "inir"
	opsED[0xB3] = "otir"
	opsED[0xB8] = "lddr"
	opsED[0xB9] = "cpdr"
	opsED[0xBA] = "indr"
	opsED[0xBB] = "otdr"
	opsED[0xF8] = "[z80]"
	opsED[0xF9] = "[z80]"
	opsED[0xFA] = "[z80]"
	opsED[0xFB] = "ed_load"
	opsED[0xFC] = "[z80]"
	opsED[0xFD] = "[z80]"
	opsED[0xFE] = "[z80]"
	opsED[0xFF] = "ed_dos"
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
		// panic(fmt.Sprintf("DISS --> 0x%X (%s)", fd.prefix, fd.op.name))
	}

	if strings.HasPrefix(op, "jr ") || strings.HasPrefix(op, "djnz ") {
		jump := int8(fd.n)
		pc := fd.pc + uint16(jump)
		op = strings.ReplaceAll(op, "$nn", toHex16(pc+2))
	} else {
		if fd.op.len == 3 {
			op = strings.ReplaceAll(op, "$nn", toHex16(fd.nn))
		} else if fd.op.len == 1 {
			op = strings.ReplaceAll(op, "$n", toHex8(fd.n))
		} else if fd.op.len == 4 {
			op = strings.ReplaceAll(op, "+0", "+"+toHex8(fd.n))
			op = strings.ReplaceAll(op, "$n2", toHex8(fd.n2))
		}
	}
	var sb strings.Builder
	sb.Grow(40)
	sb.WriteString(toHex16(fd.pc))
	sb.WriteString(": ")
	sb.WriteString(op)
	return sb.String()
	// return strings.ToLower(fmt.Sprintf("%04x: %s %s", fd.pc, fd.getMemory(), op))
}

var numbers = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

func toHex16(v uint16) string {
	var sb strings.Builder
	sb.Grow(6)
	sb.WriteString("$")
	sb.WriteString(numbers[v&0xF000>>12])
	sb.WriteString(numbers[v&0x0F00>>8])
	sb.WriteString(numbers[v&0x00F0>>4])
	sb.WriteString(numbers[v&0x000F>>0])
	return sb.String()
}

func toHex8(v uint8) string {
	var sb strings.Builder
	sb.Grow(4)
	sb.WriteString("$")
	sb.WriteString(numbers[v&0x00F0>>4])
	sb.WriteString(numbers[v&0x000F>>0])
	return sb.String()
}
