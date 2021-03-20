package lr35902

// 76543210
// ZNHC0000

type flags struct {
	C bool
	H bool
	N bool
	Z bool
}

func (f *flags) GetByte() byte {
	res := byte(0)
	if f.Z {
		res |= 0x80
	}
	if f.N {
		res |= 0x40
	}
	if f.H {
		res |= 0x20
	}
	if f.C {
		res |= 0x10
	}
	return res
}

func (f *flags) SetByte(b byte) {
	f.Z = b&0x80 != 0
	f.N = b&0x40 != 0
	f.H = b&0x20 != 0
	f.C = b&0x10 != 0
}
