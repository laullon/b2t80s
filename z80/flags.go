package z80

//	SZXHXPNC
//	|||||||`- Carry
//	||||||`-- Add/subtract
//	|||||`--- Parity/overflow
//	||||`---- Undocumented f3
//	|||`----- Half carry
//	||`------ Undocumented f5
//	|`------- Zero
//	`-------- Sign

type flags struct {
	C  bool
	N  bool
	P  bool
	F3 bool
	H  bool
	F5 bool
	Z  bool
	S  bool
}

func (f *flags) GetByte() byte {
	res := byte(0)
	if f.C {
		res |= 0b00000001
	}
	if f.N {
		res |= 0b00000010
	}
	if f.P {
		res |= 0b00000100
	}
	if f.F3 {
		res |= 0b00001000
	}
	if f.H {
		res |= 0b00010000
	}
	if f.F5 {
		res |= 0b00100000
	}
	if f.Z {
		res |= 0b01000000
	}
	if f.S {
		res |= 0b10000000
	}
	return res
}

func (f *flags) SetByte(b byte) {
	f.C = b&0b00000001 != 0
	f.N = b&0b00000010 != 0
	f.P = b&0b00000100 != 0
	f.F3 = b&0b00001000 != 0
	f.H = b&0b00010000 != 0
	f.F5 = b&0b00100000 != 0
	f.Z = b&0b01000000 != 0
	f.S = b&0b10000000 != 0
}
