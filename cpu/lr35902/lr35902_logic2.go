package lr35902

import (
	"fmt"

	cpuUtils "github.com/laullon/b2t80s/cpu"
)

func ldRn(cpu *lr35902) {
	rIdx := cpu.fetched.opCode >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	*r = cpu.fetched.n
}

func ldFromHL(cpu *lr35902) {
	mr := &mr{from: cpu.regs.HL.Get(), f: ldFromHL_m1}
	cpu.scheduler.append(mr)
}

func ldFromHL_m1(cpu *lr35902, data uint8) {
	rIdx := cpu.fetched.opCode >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	*r = data
}

func ldToHL(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	mr := &mw{to: cpu.regs.HL.Get(), d: *r, f: nil}
	cpu.scheduler.append(mr)
}

func ldToHLn(cpu *lr35902) {
	mr := &mw{to: cpu.regs.HL.Get(), d: cpu.fetched.n, f: nil}
	cpu.scheduler.append(mr)
}

var cbFuncs = []func(cpu *lr35902, r *byte){rlc, rrc, rl, rr, sla, sra, sll, srl}

func cbR(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	fIdx := cpu.fetched.opCode >> 3
	cbFuncs[fIdx](cpu, r)
	if cpu.fetched.prefix == 0 && rIdx == 0b111 {
		cpu.regs.F.Z = false
	}
}

func cbHL(cpu *lr35902) {
	mr := &mr{from: cpu.regs.HL.Get(), f: func(cpu *lr35902, data uint8) {
		b := data
		fIdx := cpu.fetched.opCode >> 3
		cbFuncs[fIdx](cpu, &b)
		mw := &mw{to: cpu.regs.HL.Get(), d: b, f: nil}
		cpu.scheduler.append(mw)
	}}
	cpu.scheduler.append(mr)
}

func bit(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	b := (cpu.fetched.opCode >> 3) & 0b111
	cpu.bit(b, *r)
}

func bitHL(cpu *lr35902) {
	mr := &mr{from: cpu.regs.HL.Get(), f: func(cpu *lr35902, data uint8) {
		v := data
		b := (cpu.fetched.opCode >> 3) & 0b111
		cpu.bit(b, v)
	}}
	cpu.scheduler.append(mr)
}

func res(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	b := (cpu.fetched.opCode >> 3) & 0b111
	cpu.res(b, r)
}

func resHL(cpu *lr35902) {
	mr := &mr{from: cpu.regs.HL.Get(), f: func(cpu *lr35902, data uint8) {
		v := data
		b := (cpu.fetched.opCode >> 3) & 0b111
		cpu.res(b, &v)
		mw := &mw{to: cpu.regs.HL.Get(), d: v, f: nil}
		cpu.scheduler.append(mw)
	}}
	cpu.scheduler.append(mr)
}

func set(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	b := (cpu.fetched.opCode >> 3) & 0b111
	cpu.set(b, r)
}

func setHL(cpu *lr35902) {
	mr := &mr{from: cpu.regs.HL.Get(), f: func(cpu *lr35902, data uint8) {
		v := data
		b := (cpu.fetched.opCode >> 3) & 0b111
		cpu.set(b, &v)
		mw := &mw{to: cpu.regs.HL.Get(), d: v, f: nil}
		cpu.scheduler.append(mw)
	}}
	cpu.scheduler.append(mr)
}

func halt(cpu *lr35902) {
	if cpu.halt && cpu.haltDone {
		cpu.halt = false
		cpu.haltDone = false
	} else {
		cpu.halt = true
		cpu.regs.PC--
	}
}

func addHLss(cpu *lr35902) {
	cpu.scheduler.append(&exec{f: func(cpu *lr35902) {
		rIdx := cpu.fetched.opCode >> 4 & 0b11
		reg := cpu.getRRptr(rIdx)

		hl := uint32(cpu.regs.HL.Get())
		rr := uint32(reg.Get())

		var result = hl + rr
		var result2 = hl&0xfff + rr&0xfff

		cpu.regs.F.N = false
		cpu.regs.F.C = result > 0xffff
		cpu.regs.F.H = result2 > 0x0fff

		cpu.regs.HL.Set(uint16(result))
	}})
}

func addSPn(cpu *lr35902) {
	cpu.scheduler.append(&exec{f: func(cpu *lr35902) {
		cpu.regs.F.C = ((cpu.regs.SP.Get() & 0xFF) + uint16(cpu.fetched.n)) > 0xFF
		cpu.regs.F.H = ((cpu.regs.SP.Get() & 0x0F) + (uint16(cpu.fetched.n) & 0x0F)) > 0x0F
		cpu.regs.F.N = false
		cpu.regs.F.Z = false
		cpu.scheduler.append(&exec{f: func(cpu *lr35902) {
			cpu.regs.SP.Set(cpu.regs.SP.Get() + uint16(int8(cpu.fetched.n)))
		}})
	}})
}

func ldAbc(cpu *lr35902) {
	from := cpu.regs.BC.Get()
	mr := &mr{from: from, f: ldAbc_m1}
	cpu.scheduler.append(mr)
}

func ldAbc_m1(cpu *lr35902, data uint8) { cpu.regs.A = data }

func ldAde(cpu *lr35902) {
	from := cpu.regs.DE.Get()
	mr := &mr{from: from, f: ldAde_m1}
	cpu.scheduler.append(mr)
}

func ldAde_m1(cpu *lr35902, data uint8) { cpu.regs.A = data }

func jrnz(cpu *lr35902) {
	if !cpu.regs.F.Z {
		jr(cpu)
	}
}

func jrnc(cpu *lr35902) {
	if !cpu.regs.F.C {
		jr(cpu)
	}
}

func jrc(cpu *lr35902) {
	if cpu.regs.F.C {
		jr(cpu)
	}
}

func jrz(cpu *lr35902) {
	if cpu.regs.F.Z {
		jr(cpu)
	}
}

func jr(cpu *lr35902) {
	cpu.scheduler.append(&exec{f: jr_m2})
}

func jr_m2(cpu *lr35902) {
	jump := int8(cpu.fetched.n)
	cpu.regs.PC += uint16(jump)
}

func scf(cpu *lr35902) {
	cpu.regs.F.H = false
	cpu.regs.F.N = false
	cpu.regs.F.C = true
}

func ccf(cpu *lr35902) {
	cpu.regs.F.N = false
	cpu.regs.F.H = false
	cpu.regs.F.C = !cpu.regs.F.C
}

func daa(cpu *lr35902) {
	tmp := int16(cpu.regs.A)
	if !cpu.regs.F.N {
		if cpu.regs.F.H || ((tmp & 0x0f) > 9) {
			tmp += 6
		}
		if cpu.regs.F.C || (tmp > 0x9f) {
			tmp += 0x60
		}
	} else {
		if cpu.regs.F.H {
			tmp -= 6
			if !cpu.regs.F.C {
				tmp &= 0xFF
			}
		}
		if cpu.regs.F.C {
			tmp -= 0x60
		}
	}

	cpu.regs.A = byte(tmp)

	if (tmp & int16(0x0100)) != 0 {
		cpu.regs.F.C = true
	}
	cpu.regs.F.H = false
	cpu.regs.F.Z = cpu.regs.A == 0
}

func cpl(cpu *lr35902) {
	cpu.regs.A = ^cpu.regs.A
	cpu.regs.F.H = true
	cpu.regs.F.N = true
}

// -------
func (cpu *lr35902) getRptr(rIdx byte) *byte {
	var r *byte
	switch rIdx {
	case 0b000:
		r = &cpu.regs.B
	case 0b001:
		r = &cpu.regs.C
	case 0b010:
		r = &cpu.regs.D
	case 0b011:
		r = &cpu.regs.E
	case 0b100:
		r = &cpu.regs.H
	case 0b101:
		r = &cpu.regs.L
	case 0b110:
		panic(fmt.Sprintf("fail on opCode: 0x%02x%02x", cpu.fetched.prefix, cpu.fetched.opCode))
	case 0b111:
		r = &cpu.regs.A
	}
	return r
}

func (cpu *lr35902) getRRptr(rIdx byte) *cpuUtils.RegPair {
	var reg *cpuUtils.RegPair
	switch rIdx {
	case 0b00:
		reg = cpu.regs.BC
	case 0b01:
		reg = cpu.regs.DE
	case 0b10:
		reg = cpu.regs.HL
	case 0b11:
		reg = cpu.regs.SP
	}
	return reg
}

func (cpu *lr35902) res(b byte, v *byte) {
	b = 1 << b
	*v &= ^b
}

func (cpu *lr35902) set(b byte, v *byte) {
	b = 1 << b
	*v |= b
}

func (cpu *lr35902) bit(b, v byte) {
	b = 1 << b
	v &= b
	cpu.regs.F.Z = v == 0
	cpu.regs.F.H = true
	cpu.regs.F.N = false
}

func (cpu *lr35902) adcA(s byte) {
	a := uint16(cpu.regs.A)
	n := uint16(s)
	res := a + n
	if cpu.regs.F.C {
		res++
	}

	res2 := a&0x0f + n&0x0f
	if cpu.regs.F.C {
		res2++
	}

	cpu.regs.A = byte(res)
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.N = false
	cpu.regs.F.C = res > 0xff
	cpu.regs.F.H = res2 > 0x0f
}

func (cpu *lr35902) cp(n byte) {
	result := uint16(cpu.regs.A) - uint16(n)
	result2 := cpu.regs.A&0x0f - n&0x0f

	cpu.regs.F.Z = result == 0
	cpu.regs.F.H = result2 > 0x0f
	cpu.regs.F.N = true
	cpu.regs.F.C = result > 0xff
}

func rlc(cpu *lr35902, r *byte) {
	*r = (*r << 1) | (*r >> 7)
	cpu.regs.F.C = *r&1 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func rrc(cpu *lr35902, r *byte) {
	cpu.regs.F.C = *r&1 != 0
	*r = (*r << 7) | (*r >> 1)
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func sll(cpu *lr35902, r *byte) {
	cpu.regs.F.C = *r&0x80 != 0
	*r = byte((*r << 1) | 0x01)
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func srl(cpu *lr35902, r *byte) {
	cpu.regs.F.C = *r&1 != 0
	*r = byte(*r >> 1)
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.N = false
	cpu.regs.F.H = false
}

func sla(cpu *lr35902, r *byte) {
	cpu.regs.F.C = *r&0x80 != 0
	*r = (*r << 1)
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func sra(cpu *lr35902, r *byte) {
	cpu.regs.F.C = *r&0x1 != 0
	b7 := *r & 0b10000000
	*r = (*r >> 1) | b7
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func rr(cpu *lr35902, r *byte) {
	c := cpu.regs.F.C
	cpu.regs.F.C = *r&0x1 != 0
	*r = (*r >> 1)
	if c {
		*r |= 0b10000000
	}
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func rl(cpu *lr35902, r *byte) {
	c := cpu.regs.F.C
	cpu.regs.F.C = *r&0x80 != 0
	*r = (*r << 1)
	if c {
		*r |= 0x1
	}
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

// ------

func (cpu *lr35902) addA(r byte) {
	a := uint16(cpu.regs.A)
	n := uint16(r)

	result := a + n
	result2 := a&0x0f + n&0x0f

	cpu.regs.A = byte(result)
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = result2 > 0x0f
	cpu.regs.F.N = false
	cpu.regs.F.C = result > 0xff
}

func (cpu *lr35902) subA(r byte) {
	a := uint16(cpu.regs.A)
	n := uint16(r)

	result := a - n
	result2 := a&0x0f - n&0x0f

	cpu.regs.A = byte(result)
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = result2 > 0x0f
	cpu.regs.F.N = true
	cpu.regs.F.C = result > 0xff
}

func (cpu *lr35902) xor(s uint8) {
	cpu.regs.A = cpu.regs.A ^ s
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
	cpu.regs.F.C = false
}

func (cpu *lr35902) and(s uint8) {
	cpu.regs.A = cpu.regs.A & s
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = true
	cpu.regs.F.N = false
	cpu.regs.F.C = false
}

func (cpu *lr35902) or(s uint8) {
	// TODO: review p/v flag
	cpu.regs.A = cpu.regs.A | s
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
	cpu.regs.F.C = false
}

func (cpu *lr35902) sbcA(s byte) {
	a := uint16(cpu.regs.A)
	n := uint16(s)

	result := a - n
	if cpu.regs.F.C {
		result--
	}
	result2 := a&0x0f - n&0x0f
	if cpu.regs.F.C {
		result2--
	}

	cpu.regs.A = byte(result)
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = result2 > 0x0f
	cpu.regs.F.N = true
	cpu.regs.F.C = result > 0xff
}

func swap(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	*r = ((*r & 0b11110000) >> 4) | ((*r & 0b00001111) << 4)
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.H = false
	cpu.regs.F.C = false
	cpu.regs.F.N = false
}

func swapHL(cpu *lr35902) {
	cpu.scheduler.append(&mr{from: cpu.regs.HL.Get(), f: swapHL_m2})
}

func swapHL_m2(cpu *lr35902, data byte) {
	r := ((data & 0b11110000) >> 4) | ((data & 0b00001111) << 4)
	cpu.regs.F.Z = r == 0
	cpu.regs.F.H = false
	cpu.regs.F.C = false
	cpu.regs.F.N = false
	cpu.scheduler.append(&mw{to: cpu.regs.HL.Get(), d: r, f: nil})
}

func ldHLspE(cpu *lr35902) {
	cpu.scheduler.append(&exec{f: func(cpu *lr35902) {
		sp := cpu.regs.SP.Get()
		res := sp + uint16(int8(cpu.fetched.n))
		cpu.regs.HL.Set(res)

		cpu.regs.F.C = ((sp & 0xff) + uint16(cpu.fetched.n)&0xff) > 0xff
		cpu.regs.F.H = ((sp & 0x0f) + uint16(cpu.fetched.n)&0x0f) > 0x0f
		cpu.regs.F.Z = false
		cpu.regs.F.N = false
	}})
}

func ldNNsp(cpu *lr35902) {
	mm := cpu.fetched.nn
	mw1 := &mw{to: mm, d: byte(*cpu.regs.SP.L), f: nil}
	mw2 := &mw{to: mm + 1, d: byte(*cpu.regs.SP.H), f: nil}
	cpu.scheduler.append(mw1, mw2)
}

func ldhNa(cpu *lr35902) {
	mm := uint16(cpu.fetched.n) | uint16(0xff00)
	mw1 := &mw{to: mm, d: cpu.regs.A, f: nil}
	cpu.scheduler.append(mw1)
}

func ldhAn(cpu *lr35902) {
	mm := uint16(cpu.fetched.n) | uint16(0xff00)
	mr1 := &mr{from: mm, f: ldhAc_m2}
	cpu.scheduler.append(mr1)
}

func ldhAc_m2(cpu *lr35902, data uint8) {
	cpu.regs.A = data
}

func ldhCa(cpu *lr35902) {
	mm := uint16(0xff00) | uint16(cpu.regs.C)
	mw1 := &mw{to: mm, d: cpu.regs.A, f: nil}
	cpu.scheduler.append(mw1)
}

func ldhAc(cpu *lr35902) {
	mm := uint16(cpu.regs.C) | uint16(0xff00)
	mr1 := &mr{from: mm, f: ldhAc_m2}
	cpu.scheduler.append(mr1)
}

func ldiAhl(cpu *lr35902) {
	cpu.scheduler.append(&mr{from: cpu.regs.HL.Get(), f: ldiAhl_m2})
}

func ldiAhl_m2(cpu *lr35902, data byte) {
	cpu.regs.A = data
	cpu.regs.HL.Set(cpu.regs.HL.Get() + 1)
}

func ldSPhl(cpu *lr35902) {
	cpu.scheduler.append(&exec{f: func(cpu *lr35902) {
		cpu.regs.SP.Set(cpu.regs.HL.Get())
	}})
}
