package lr35902

import cpuUtils "github.com/laullon/b2t80s/cpu"

func decIXYd(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := newMR(addr,
		func(cpu *lr35902, data uint8) {
			r := data
			cpu.regs.F.H = r&0x0f == 0
			r--
			cpu.regs.F.Z = r == 0
			cpu.regs.F.N = true

			mw := newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		},
	)
	cpu.scheduler.append(mr)
}

func ldRn(cpu *lr35902) {
	rIdx := cpu.fetched.opCode >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	*r = cpu.fetched.n
}

func ldRhl(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(), ldR_m1)
	cpu.scheduler.append(mr)
}

func ldRixyD(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := newMR(addr, ldR_m1)
	cpu.scheduler.append(mr)
}

func ldR_m1(cpu *lr35902, data uint8) {
	rIdx := cpu.fetched.opCode >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	*r = data
}

func addAixyD(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := newMR(addr, addAixyD_m1)
	cpu.scheduler.append(mr)
}

func addAixyD_m1(cpu *lr35902, data uint8) { cpu.addA(data) }

func adcAixyD(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := newMR(addr, adcAixyD_m1)
	cpu.scheduler.append(mr)
}

func adcAixyD_m1(cpu *lr35902, data uint8) { cpu.adcA(data) }

func subAixyD(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := newMR(addr, subAixyD_m1)
	cpu.scheduler.append(mr)
}

func subAixyD_m1(cpu *lr35902, data uint8) { cpu.subA(data) }

func sbcAixyD(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := newMR(addr, sbcAixyD_m1)
	cpu.scheduler.append(mr)
}

func sbcAixyD_m1(cpu *lr35902, data uint8) { cpu.sbcA(data) }

func andAixyD(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := newMR(addr, andAixyD_m1)
	cpu.scheduler.append(mr)
}

func andAixyD_m1(cpu *lr35902, data uint8) { cpu.and(data) }

func xorAixyD(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := newMR(addr, xorAixyD_m1)
	cpu.scheduler.append(mr)
}

func xorAixyD_m1(cpu *lr35902, data uint8) { cpu.xor(data) }

func cpAixyD(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := newMR(addr, cpAixyD_m1)
	cpu.scheduler.append(mr)
}

func cpAixyD_m1(cpu *lr35902, data uint8) { cpu.cp(data) }

func orAixyD(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := newMR(addr, orAixyD_m1)
	cpu.scheduler.append(mr)
}

func orAixyD_m1(cpu *lr35902, data uint8) { cpu.or(data) }

func ldHLr(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	mr := newMW(cpu.regs.HL.Get(), *r, nil)
	cpu.scheduler.append(mr)
}

func ldIXYHr(cpu *lr35902) {
	reg := cpu.indexRegs[cpu.indexIdx]
	rIdx := cpu.fetched.opCode & 0b111
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
		r = reg.H
	case 0b101:
		r = reg.L
	case 0b110:
		panic(-1)
	case 0b111:
		r = &cpu.regs.A
	}
	*reg.H = *r
}

func ldIXYLr(cpu *lr35902) {
	reg := cpu.indexRegs[cpu.indexIdx]
	rIdx := cpu.fetched.opCode & 0b111
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
		r = reg.H
	case 0b101:
		r = reg.L
	case 0b110:
		panic(-1)
	case 0b111:
		r = &cpu.regs.A
	}
	*reg.L = *r
}

func ldRr(cpu *lr35902) {
	r1Idx := cpu.fetched.opCode >> 3 & 0b111
	r2Idx := cpu.fetched.opCode & 0b111
	r1 := cpu.getRptr(r1Idx)
	r2 := cpu.getRptr(r2Idx)
	*r1 = *r2
}

func rlca(cpu *lr35902) {
	cpu.regs.A = cpu.regs.A<<1 | cpu.regs.A>>7
	cpu.regs.F.C = cpu.regs.A&0x01 != 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func rla(cpu *lr35902) {
	c := cpu.regs.F.C
	cpu.regs.F.C = cpu.regs.A&0b10000000 != 0
	cpu.regs.A = (cpu.regs.A << 1)
	if c {
		cpu.regs.A |= 1
	}
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

var cbFuncs = []func(cpu *lr35902, r *byte){rlc, rrc, rl, rr, sla, sra, sll, srl}

func cbR(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	fIdx := cpu.fetched.opCode >> 3
	cbFuncs[fIdx](cpu, r)
}

func cbHL(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(),
		func(cpu *lr35902, data uint8) {
			b := data
			fIdx := cpu.fetched.opCode >> 3
			cbFuncs[fIdx](cpu, &b)
			mw := newMW(cpu.regs.HL.Get(), b, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func cbIXYdr(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := newMR(addr,
		func(cpu *lr35902, data uint8) {
			r := data
			fIdx := (cpu.fetched.opCode >> 3) & 0b111
			cbFuncs[fIdx](cpu, &r)

			rIdx := cpu.fetched.opCode & 0b111
			reg := cpu.getRptr(rIdx)
			*reg = r

			mw := newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func cbIXYd(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := newMR(addr,
		func(cpu *lr35902, data uint8) {
			r := data
			fIdx := (cpu.fetched.opCode >> 3) & 0b111
			cbFuncs[fIdx](cpu, &r)
			mw := newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func bit(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	b := (cpu.fetched.opCode >> 3) & 0b111
	cpu.bit(b, *r)
}

func bitIXYd(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	b := (cpu.fetched.opCode >> 3) & 0b111
	mr := newMR(addr,
		func(cpu *lr35902, data uint8) {
			r := data
			cpu.bit(b, r)
		})
	cpu.scheduler.append(mr)
}

func bitHL(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(),
		func(cpu *lr35902, data uint8) {
			v := data
			b := (cpu.fetched.opCode >> 3) & 0b111
			cpu.bit(b, v)
		})
	cpu.scheduler.append(mr)
}

func res(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	b := (cpu.fetched.opCode >> 3) & 0b111
	cpu.res(b, r)
}

func resHL(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(),
		func(cpu *lr35902, data uint8) {
			v := data
			b := (cpu.fetched.opCode >> 3) & 0b111
			cpu.res(b, &v)
			mw := newMW(cpu.regs.HL.Get(), v, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func resIXYdR(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	b := (cpu.fetched.opCode >> 3) & 0b111
	mr := newMR(addr,
		func(cpu *lr35902, data uint8) {
			r := data
			cpu.res(b, &r)

			rIdx := cpu.fetched.opCode & 0b111
			reg := cpu.getRptr(rIdx)
			*reg = r

			mw := newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func resIXYd(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	b := (cpu.fetched.opCode >> 3) & 0b111
	mr := newMR(addr,
		func(cpu *lr35902, data uint8) {
			r := data
			cpu.res(b, &r)
			mw := newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func set(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	b := (cpu.fetched.opCode >> 3) & 0b111
	cpu.set(b, r)
}

func setHL(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(),
		func(cpu *lr35902, data uint8) {
			v := data
			b := (cpu.fetched.opCode >> 3) & 0b111
			cpu.set(b, &v)
			mw := newMW(cpu.regs.HL.Get(), v, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func setIXYdR(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	b := (cpu.fetched.opCode >> 3) & 0b111
	mr := newMR(addr,
		func(cpu *lr35902, data uint8) {
			r := data
			cpu.set(b, &r)

			rIdx := cpu.fetched.opCode & 0b111
			reg := cpu.getRptr(rIdx)
			*reg = r

			mw := newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func setIXYd(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	b := (cpu.fetched.opCode >> 3) & 0b111
	mr := newMR(addr,
		func(cpu *lr35902, data uint8) {
			r := data
			cpu.set(b, &r)
			mw := newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func rrca(cpu *lr35902) {
	cpu.regs.F.C = cpu.regs.A&0x01 != 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
	cpu.regs.A = (cpu.regs.A >> 1) | (cpu.regs.A << 7)
}

func rra(cpu *lr35902) {
	c := cpu.regs.F.C
	cpu.regs.F.C = cpu.regs.A&1 != 0
	cpu.regs.A = (cpu.regs.A >> 1)
	if c {
		cpu.regs.A |= 0b10000000
	}
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func exDEhl(cpu *lr35902) {
	cpu.regs.D, cpu.regs.H = cpu.regs.H, cpu.regs.D
	cpu.regs.E, cpu.regs.L = cpu.regs.L, cpu.regs.E
}

func halt(cpu *lr35902) {
	cpu.halt = true
	cpu.regs.PC--
}

func addHLss(cpu *lr35902) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)

	hl := cpu.regs.HL.Get()
	var result = uint32(hl) + uint32(reg.Get())
	var lookup = byte(((hl & 0x0800) >> 11) | ((reg.Get() & 0x0800) >> 10) | ((uint16(result) & 0x0800) >> 9))
	cpu.regs.HL.Set(uint16(result))

	cpu.regs.F.N = false
	cpu.regs.F.H = halfcarryAddTable[lookup]
	cpu.regs.F.C = (result & 0x10000) != 0
}

func ldAbc(cpu *lr35902) {
	from := cpu.regs.BC.Get()
	mr := newMR(from, ldAbc_m1)
	cpu.scheduler.append(mr)
}

func ldAbc_m1(cpu *lr35902, data uint8) { cpu.regs.A = data }

func ldAde(cpu *lr35902) {
	from := cpu.regs.DE.Get()
	mr := newMR(from, ldAde_m1)
	cpu.scheduler.append(mr)
}

func ldAde_m1(cpu *lr35902, data uint8) { cpu.regs.A = data }

func djnz(cpu *lr35902) {
	cpu.regs.B--
	if cpu.regs.B != 0 {
		cpu.scheduler.append(&exec{l: 5, f: jr})
	}
}

func jrnz(cpu *lr35902) {
	if !cpu.regs.F.Z {
		cpu.scheduler.append(&exec{l: 5, f: jr})
	}
}

func jrnc(cpu *lr35902) {
	if !cpu.regs.F.C {
		cpu.scheduler.append(&exec{l: 5, f: jr})
	}
}

func jrc(cpu *lr35902) {
	if cpu.regs.F.C {
		cpu.scheduler.append(&exec{l: 5, f: jr})
	}
}

func jrz(cpu *lr35902) {
	if cpu.regs.F.Z {
		cpu.scheduler.append(&exec{l: 5, f: jr})
	}
}

func jr(cpu *lr35902) {
	jump := int8(cpu.fetched.n)
	cpu.regs.PC += uint16(jump)
}

func scf(cpu *lr35902) {
	cpu.regs.F.H = false
	cpu.regs.F.N = false
	cpu.regs.F.C = true
}

func ccf(cpu *lr35902) {
	cpu.regs.F.H = cpu.regs.F.C
	cpu.regs.F.N = false
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
		panic(-1)
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

func (cpu *lr35902) getIXYn(n byte) uint16 {
	reg := cpu.indexRegs[cpu.indexIdx]
	i := int16(int8(n))
	ix := reg.Get()
	ix = uint16(int16(ix) + i)
	return ix
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
	res := int16(cpu.regs.A) + int16(s)
	if cpu.regs.F.C {
		res++
	}
	lookup := ((cpu.regs.A & 0x88) >> 3) | ((s & 0x88) >> 2) | ((byte(res) & 0x88) >> 1)
	cpu.regs.A = byte(res)
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = halfcarryAddTable[lookup&0x07]
	cpu.regs.F.N = false
	cpu.regs.F.C = (res & 0x100) == 0x100
}

func (cpu *lr35902) adcHL(ss uint16) {
	hl := cpu.regs.HL.Get()
	res := int32(hl) + int32(ss)
	if cpu.regs.F.C {
		res++
	}
	lookup := byte(((hl & 0x8800) >> 11) | ((ss & 0x8800) >> 10) | ((uint16(res) & 0x8800) >> 9))
	hl = uint16(res)
	cpu.regs.HL.Set(hl)
	cpu.regs.F.Z = hl == 0
	cpu.regs.F.H = halfcarryAddTable[lookup&0x07]
	cpu.regs.F.N = false
	cpu.regs.F.C = (res & 0x10000) != 0
}

func (cpu *lr35902) cp(r byte) {
	a := int16(cpu.regs.A)
	result := a - int16(r)
	lookup := ((cpu.regs.A & 0x88) >> 3) | (((r) & 0x88) >> 2) | ((byte(result) & 0x88) >> 1)

	cpu.regs.F.Z = result == 0
	cpu.regs.F.H = halfcarrySubTable[lookup&0x07]
	cpu.regs.F.N = true
	cpu.regs.F.C = ((result) & 0x100) == 0x100
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
	a := int16(cpu.regs.A)
	result := a + int16(r)
	lookup := ((cpu.regs.A & 0x88) >> 3) | (((r) & 0x88) >> 2) | ((byte(result) & 0x88) >> 1)
	cpu.regs.A = uint8(result & 0x00ff)

	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = halfcarryAddTable[lookup&0x07]
	cpu.regs.F.N = false
	cpu.regs.F.C = ((result) & 0x100) != 0
}

func (cpu *lr35902) subA(r byte) {
	a := int16(cpu.regs.A)
	result := a - int16(r)
	lookup := ((cpu.regs.A & 0x88) >> 3) | (((r) & 0x88) >> 2) | ((byte(result) & 0x88) >> 1)
	cpu.regs.A = uint8(result & 0x00ff)

	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = halfcarrySubTable[lookup&0x07]
	cpu.regs.F.N = true
	cpu.regs.F.C = ((result) & 0x100) == 0x100
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
	res := uint16(cpu.regs.A) - uint16(s)
	if cpu.regs.F.C {
		res--
	}
	lookup := ((cpu.regs.A & 0x88) >> 3) | ((s & 0x88) >> 2) | byte(res&0x88>>1)
	cpu.regs.A = byte(res)
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = halfcarrySubTable[lookup&0x07]
	cpu.regs.F.N = true
	cpu.regs.F.C = (res & 0x100) == 0x100
}

func (cpu *lr35902) sbcHL(ss uint16) {
	hl := cpu.regs.HL.Get()
	res := uint32(hl) - uint32(ss)
	if cpu.regs.F.C {
		res--
	}
	cpu.regs.HL.Set(uint16(res))

	lookup := byte(((hl & 0x8800) >> 11) | ((ss & 0x8800) >> 10) | ((uint16(res) & 0x8800) >> 9))
	cpu.regs.F.N = true
	cpu.regs.F.Z = res == 0
	cpu.regs.F.C = (res & 0x10000) != 0
	cpu.regs.F.H = halfcarrySubTable[lookup&0x07]
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

func ldHLspE(cpu *lr35902) {
	sp := cpu.regs.SP.Get()
	res := sp + uint16(int8(cpu.fetched.n))
	cpu.regs.HL.Set(res)

	cpu.regs.F.C = ((sp & 0xff) + uint16(cpu.fetched.n)&0xf) > 0xff
	cpu.regs.F.H = ((sp & 0x0f) + uint16(cpu.fetched.n)&0x0) > 0x0f
	cpu.regs.F.Z = false
	cpu.regs.F.N = false
}

func ldNNsp(cpu *lr35902) {
	mm := cpu.fetched.nn
	mw1 := newMW(mm, byte(*cpu.regs.SP.L), nil)
	mw2 := newMW(mm+1, byte(*cpu.regs.SP.H), nil)
	cpu.scheduler.append(mw1, mw2)
}

func ldhNa(cpu *lr35902) {
	mm := uint16(cpu.fetched.n) | uint16(0xff00)
	mw1 := newMW(mm, cpu.regs.A, nil)
	cpu.scheduler.append(mw1)
}

func ldhAn(cpu *lr35902) {
	mm := uint16(cpu.fetched.n) | uint16(0xff00)
	mr1 := newMR(mm, ldhAc_m2)
	cpu.scheduler.append(mr1)
}

func ldhAc_m2(cpu *lr35902, data uint8) {
	cpu.regs.A = data
}

func ldhCa(cpu *lr35902) {
	mm := uint16(0xff00) | uint16(cpu.regs.C)
	mw1 := newMW(mm, cpu.regs.A, nil)
	cpu.scheduler.append(mw1)
}

func ldhAc(cpu *lr35902) {
	mm := uint16(0xff00) | uint16(cpu.regs.C)
	mr1 := newMR(mm, ldhAc_m2)
	cpu.scheduler.append(mr1)
}

func ldiAhl(cpu *lr35902) {
	cpu.scheduler.append(newMR(cpu.regs.HL.Get(), ldiAhl_m2))
}

func ldiAhl_m2(cpu *lr35902, data byte) {
	cpu.regs.A = data
	cpu.regs.HL.Set(cpu.regs.HL.Get() + 1)
}
