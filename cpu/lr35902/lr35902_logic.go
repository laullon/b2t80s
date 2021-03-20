package lr35902

import (
	cpuUtils "github.com/laullon/b2t80s/cpu"
)

func cpi(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(), cpi_m1)
	cpu.scheduler.append(mr)
}

var cpi_result uint8

func cpi_m1(cpu *lr35902, data uint8) {

	val := data
	cpi_result = cpu.regs.A - val
	lookup := (cpu.regs.A&0x08)>>3 | (val&0x08)>>2 | (cpi_result&0x08)>>1
	cpu.regs.F.H = halfcarrySubTable[lookup]

	cpu.scheduler.append(&exec{l: 5, f: cpi_m2})
}

func cpi_m2(cpu *lr35902) {
	bc := cpu.regs.BC.Get()
	bc--
	cpu.regs.BC.Set(bc)
	cpu.regs.HL.Set(cpu.regs.HL.Get() + 1)

	cpu.regs.F.Z = cpi_result == 0
	cpu.regs.F.N = true
	if cpu.fetched.opCode > 0xAF {
		cpu.scheduler.append(&exec{l: 5, f: cpi_m3})
	}
}

func cpi_m3(cpu *lr35902) {
	if (cpu.regs.BC.Get()) != 0 && (cpi_result != 0) {
		cpu.regs.PC = cpu.regs.PC - 2
	}
}

func cpd(cpu *lr35902) {
	hl := cpu.regs.HL.Get()
	mr := newMR(hl, cpd_m1)
	cpu.scheduler.append(mr)
}

func cpd_m1(cpu *lr35902, data uint8) {
	val := data
	result := cpu.regs.A - byte(val)
	lookup := (cpu.regs.A&0x08)>>3 | (val&0x08)>>2 | (result&0x08)>>1

	cpu.regs.F.Z = result == 0
	cpu.regs.F.H = halfcarrySubTable[lookup]

	cpu.scheduler.append(&exec{l: 5, f: cpd_m2})
}

func cpd_m2(cpu *lr35902) {
	bc := cpu.regs.BC.Get()
	hl := cpu.regs.HL.Get()

	bc--
	hl--

	cpu.regs.BC.Set(bc)
	cpu.regs.HL.Set(hl)

	cpu.regs.F.N = true

	if cpu.fetched.opCode > 0xAF {
		cpu.scheduler.append(&exec{l: 5, f: cpd_m3})
	}
}

func cpd_m3(cpu *lr35902) {
	if (cpu.regs.BC.Get() != 0) && (!cpu.regs.F.Z) {
		cpu.regs.PC = cpu.regs.PC - 2
	}
}

func ldi(cpu *lr35902) {
	hl := cpu.regs.HL.Get()
	mr := newMR(hl, ldi_m1)
	cpu.scheduler.append(mr)
}

func ldi_m1(cpu *lr35902, data uint8) {
	v := data
	de := cpu.regs.DE.Get()
	mw := newMW(de, v, ldi_m2)
	cpu.scheduler.append(&exec{l: 2}, mw)
}

func ldi_m2(cpu *lr35902) {
	bc := cpu.regs.BC.Get()
	de := cpu.regs.DE.Get()
	hl := cpu.regs.HL.Get()

	bc--
	de++
	hl++

	cpu.regs.BC.Set(bc)
	cpu.regs.DE.Set(de)
	cpu.regs.HL.Set(hl)

	cpu.regs.F.H = false
	cpu.regs.F.N = false
	if cpu.fetched.opCode > 0xAF {
		cpu.scheduler.append(&exec{l: 5, f: ldi_m3})
	}
}

func ldi_m3(cpu *lr35902) {
	if cpu.regs.BC.Get() != 0 {
		cpu.regs.PC = cpu.regs.PC - 2
	}
}

func ldd(cpu *lr35902) {
	hl := cpu.regs.HL.Get()

	mr := newMR(hl, ldd_m1)
	cpu.scheduler.append(mr)
}

func ldd_m1(cpu *lr35902, data uint8) {
	de := cpu.regs.DE.Get()
	v := data
	mw := newMW(de, v, ldd_m2)
	cpu.scheduler.append(&exec{l: 2}, mw)
}

func ldd_m2(cpu *lr35902) {
	bc := cpu.regs.BC.Get()
	de := cpu.regs.DE.Get()
	hl := cpu.regs.HL.Get()

	bc--
	de--
	hl--

	cpu.regs.BC.Set(bc)
	cpu.regs.DE.Set(de)
	cpu.regs.HL.Set(hl)

	cpu.regs.F.H = false
	cpu.regs.F.N = false
	if cpu.fetched.opCode > 0xAF {
		cpu.scheduler.append(&exec{l: 5, f: ldd_m3})
	}
}

func ldd_m3(cpu *lr35902) {
	if cpu.regs.BC.Get() != 0 {
		cpu.regs.PC = cpu.regs.PC - 2
	}
}

var spv uint16

func exSP(cpu *lr35902) {
	reg := cpu.indexRegs[cpu.indexIdx]
	mr1 := newMR(cpu.regs.SP.Get(), exSP_m1)
	mr2 := newMR(cpu.regs.SP.Get()+1, exSP_m2)
	mw1 := newMW(cpu.regs.SP.Get(), *reg.L, nil)
	mw2 := newMW(cpu.regs.SP.Get()+1, *reg.H, exSP_m3)
	cpu.scheduler.append(mr1, &exec{l: 1}, mr2, mw1, &exec{l: 2}, mw2)
}

func exSP_m1(cpu *lr35902, data uint8) { spv = uint16(data) }
func exSP_m2(cpu *lr35902, data uint8) { spv |= uint16(data) << 8 }
func exSP_m3(cpu *lr35902)             { reg := cpu.indexRegs[cpu.indexIdx]; reg.Set(spv) }

func addIXY(cpu *lr35902) {
	var reg *cpuUtils.RegPair
	regI := cpu.indexRegs[cpu.indexIdx]
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	switch rIdx {
	case 0b00:
		reg = cpu.regs.BC
	case 0b01:
		reg = cpu.regs.DE
	case 0b10:
		reg = regI
	case 0b11:
		reg = cpu.regs.SP
	}

	ix := regI.Get()
	var result = uint32(ix) + uint32(reg.Get())
	var lookup = byte(((ix & 0x0800) >> 11) | ((reg.Get() & 0x0800) >> 10) | ((uint16(result) & 0x0800) >> 9))
	regI.Set(uint16(result))

	cpu.regs.F.N = false
	cpu.regs.F.H = halfcarryAddTable[lookup]
	cpu.regs.F.C = (result & 0x10000) != 0
}

func retCC(cpu *lr35902) {
	ccIdx := cpu.fetched.opCode >> 3 & 0b111
	branch := cpu.checkCondition(ccIdx)
	if branch {
		cpu.popFromStack(func(cpu *lr35902, data uint16) {
			cpu.regs.PC = data
		})
	}
}

func ret(cpu *lr35902) {
	cpu.popFromStack(func(cpu *lr35902, data uint16) {
		cpu.regs.PC = data
	})
}

func rstP(cpu *lr35902) {
	cpu.pushToStack(cpu.regs.PC, rstP_m1)
}

func rstP_m1(cpu *lr35902) {
	newPCs := []uint16{0x00, 0x08, 0x10, 0x18, 0x20, 0x28, 0x30, 0x38}
	pIdx := cpu.fetched.opCode >> 3 & 0b111
	cpu.regs.PC = newPCs[pIdx]
}

func jpCC(cpu *lr35902) {
	ccIdx := cpu.fetched.opCode >> 3 & 0b111
	branch := cpu.checkCondition(ccIdx)
	if branch {
		cpu.regs.PC = cpu.fetched.nn
	}
}

func callCC(cpu *lr35902) {
	ccIdx := cpu.fetched.opCode >> 3 & 0b11
	branch := cpu.checkCondition(ccIdx)
	if branch {
		cpu.scheduler.append(&exec{l: 1, f: callCC_m2})
	}
}

func callCC_m2(cpu *lr35902) { cpu.pushToStack(cpu.regs.PC, call_m1) }

func call(cpu *lr35902) {
	cpu.pushToStack(cpu.regs.PC, call_m1)
}

func call_m1(cpu *lr35902) { cpu.regs.PC = cpu.fetched.nn }

func (cpu *lr35902) checkCondition(ccIdx byte) bool {
	res := false
	switch ccIdx {
	case 0:
		res = !cpu.regs.F.Z
	case 1:
		res = cpu.regs.F.Z
	case 2:
		res = !cpu.regs.F.C
	case 3:
		res = cpu.regs.F.C
	default:
		panic(-1)
	}
	return res
}

var pushF func(cpu *lr35902)

func (cpu *lr35902) pushToStack(data uint16, f func(cpu *lr35902)) {
	pushF = f
	push1 := newMW(cpu.regs.SP.Get()-1, uint8(data>>8), nil)
	push2 := newMW(cpu.regs.SP.Get()-2, uint8(data), push_m1)
	cpu.scheduler.append(push1, push2)
}

func push_m1(cpu *lr35902) {
	cpu.regs.SP.Set(cpu.regs.SP.Get() - 2)
	if pushF != nil {
		pushF(cpu)
	}
}

var popData uint16
var popF func(cpu *lr35902, data uint16)

func (cpu *lr35902) popFromStack(f func(cpu *lr35902, data uint16)) {
	popF = f
	pop1 := newMR(cpu.regs.SP.Get(), pop_m1)
	pop2 := newMR(cpu.regs.SP.Get()+1, pop_m2)
	cpu.scheduler.append(pop1, pop2)
}

func pop_m1(cpu *lr35902, data uint8) { popData = uint16(data) }

func pop_m2(cpu *lr35902, data uint8) {
	popData |= (uint16(data) << 8)
	cpu.regs.SP.Set(cpu.regs.SP.Get() + 2)
	popF(cpu, popData)
}

func popSS(cpu *lr35902) {
	cpu.popFromStack(popSS_m1)
}

func popSS_m1(cpu *lr35902, data uint16) {
	t := cpu.fetched.opCode >> 4 & 0b11
	switch t {
	case 0b00:
		cpu.regs.BC.Set(data)
	case 0b01:
		cpu.regs.DE.Set(data)
	case 0b10:
		cpu.regs.HL.Set(data)
	case 0b11:
		cpu.regs.A = uint8(data >> 8)
		cpu.regs.F.SetByte(uint8(data))
	}
}

func pushSS(cpu *lr35902) {
	t := cpu.fetched.opCode >> 4 & 0b11
	var data uint16
	switch t {
	case 0b00:
		data = cpu.regs.BC.Get()
	case 0b01:
		data = cpu.regs.DE.Get()
	case 0b10:
		data = cpu.regs.HL.Get()
	case 0b11:
		data = uint16(cpu.regs.A) << 8
		data |= uint16(cpu.regs.F.GetByte())
	}
	cpu.pushToStack(data, nil)
}

func ldDDmm(cpu *lr35902) {
	t := cpu.fetched.opCode >> 4 & 0b11
	switch t {
	case 0b00:
		cpu.regs.B = cpu.fetched.n2
		cpu.regs.C = cpu.fetched.n
	case 0b01:
		cpu.regs.D = cpu.fetched.n2
		cpu.regs.E = cpu.fetched.n
	case 0b10:
		cpu.regs.H = cpu.fetched.n2
		cpu.regs.L = cpu.fetched.n
	case 0b11:
		cpu.regs.S = cpu.fetched.n2
		cpu.regs.P = cpu.fetched.n
	}
}

func ldBCa(cpu *lr35902) {
	pos := cpu.regs.BC.Get()
	cpu.scheduler.append(newMW(pos, cpu.regs.A, nil))
}

func ldDEa(cpu *lr35902) {
	pos := cpu.regs.DE.Get()
	cpu.scheduler.append(newMW(pos, cpu.regs.A, nil))
}

func ldNNhl(cpu *lr35902) {
	mm := cpu.fetched.nn
	mw1 := newMW(mm, cpu.regs.L, nil)
	mw2 := newMW(mm+1, cpu.regs.H, nil)
	cpu.scheduler.append(mw1, mw2)
}

func ldNNIXY(cpu *lr35902) {
	reg := cpu.indexRegs[cpu.indexIdx]
	mm := cpu.fetched.nn
	mw1 := newMW(mm, *reg.L, nil)
	mw2 := newMW(mm+1, *reg.H, nil)
	cpu.scheduler.append(mw1, mw2)
}

func ldNNa(cpu *lr35902) {
	mm := cpu.fetched.nn
	mw1 := newMW(mm, cpu.regs.A, nil)
	cpu.scheduler.append(mw1)
}

func ldiHLa(cpu *lr35902) {
	cpu.scheduler.append(newMW(cpu.regs.HL.Get(), cpu.regs.A, ldiHLa_m2))
}

func ldiHLa_m2(cpu *lr35902) {
	cpu.regs.HL.Set(cpu.regs.HL.Get() + 1)
}

func lddAhl(cpu *lr35902) {
	cpu.scheduler.append(newMR(cpu.regs.HL.Get(), lddAhl_m2))
}

func lddAhl_m2(cpu *lr35902, data byte) {
	cpu.regs.A = data
	cpu.regs.HL.Set(cpu.regs.HL.Get() - 1)
}

func lddHLa(cpu *lr35902) {
	cpu.scheduler.append(newMW(cpu.regs.HL.Get(), cpu.regs.A, lddHLa_m2))
}

func lddHLa_m2(cpu *lr35902) {
	cpu.regs.HL.Set(cpu.regs.HL.Get() - 1)
}

var hlv uint8

func rrd(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(), rrd_m1)
	cpu.scheduler.append(mr)
}

func rrd_m1(cpu *lr35902, data uint8) {
	hlv = data
	mw := newMW(cpu.regs.HL.Get(), (cpu.regs.A<<4 | hlv>>4), rrd_m2)
	cpu.scheduler.append(&exec{l: 4}, mw)
}

func rrd_m2(cpu *lr35902) {
	cpu.regs.A = (cpu.regs.A & 0xf0) | (hlv & 0x0f)
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func rld(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(), rld_m1)
	cpu.scheduler.append(mr)
}

func rld_m1(cpu *lr35902, data uint8) {
	hlv = data
	mw := newMW(cpu.regs.HL.Get(), (hlv<<4 | cpu.regs.A&0x0f), rld_m2)
	cpu.scheduler.append(&exec{l: 4}, mw)
}

func rld_m2(cpu *lr35902) {
	cpu.regs.A = (cpu.regs.A & 0xf0) | (hlv >> 4)
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func ldNNdd(cpu *lr35902) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	mm := cpu.fetched.nn
	mw1 := newMW(mm, *reg.L, nil)
	mw2 := newMW(mm+1, *reg.H, nil)
	cpu.scheduler.append(mw1, mw2)
}

func ldDDnn(cpu *lr35902) {
	mm := cpu.fetched.nn
	mr1 := newMR(mm, ldDDnn_m1)
	mr2 := newMR(mm+1, ldDDnn_m2)
	cpu.scheduler.append(mr1, mr2)
}

func ldDDnn_m1(cpu *lr35902, data uint8) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	*reg.L = data
}
func ldDDnn_m2(cpu *lr35902, data uint8) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	*reg.H = data
}

func ldAi(cpu *lr35902) {
	cpu.regs.A = cpu.regs.I
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func ldAr(cpu *lr35902) {
	cpu.regs.A = cpu.regs.R
	cpu.regs.F.Z = cpu.regs.R == 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func ldHLnn(cpu *lr35902) {
	mm := cpu.fetched.nn
	mr1 := newMR(mm, ldHLnn_m1)
	mr2 := newMR(mm+1, ldHLnn_m2)
	cpu.scheduler.append(mr1, mr2)
}

func ldHLnn_m1(cpu *lr35902, data uint8) { cpu.regs.L = data }
func ldHLnn_m2(cpu *lr35902, data uint8) { cpu.regs.H = data }

func ldIXYnn(cpu *lr35902) {
	mm := cpu.fetched.nn
	mr1 := newMR(mm, ldIXYnn_m1)
	mr2 := newMR(mm+1, ldIXYnn_m2)
	cpu.scheduler.append(mr1, mr2)
}

func ldIXYnn_m1(cpu *lr35902, data uint8) {
	reg := cpu.indexRegs[cpu.indexIdx]
	*reg.L = data
}

func ldIXYnn_m2(cpu *lr35902, data uint8) {
	reg := cpu.indexRegs[cpu.indexIdx]
	*reg.H = data
}

func ldAnn(cpu *lr35902) {
	mm := cpu.fetched.nn
	mr1 := newMR(mm, ldAnn_n1)
	cpu.scheduler.append(mr1)
}

func ldAnn_n1(cpu *lr35902, data uint8) { cpu.regs.A = data }

func ldHLn(cpu *lr35902) {
	mw1 := newMW(cpu.regs.HL.Get(), cpu.fetched.n, nil)
	cpu.scheduler.append(mw1)
}

func ldIXYdN(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mw1 := newMW(addr, cpu.fetched.n2, nil)
	cpu.scheduler.append(mw1)
}

func ldIXYdR(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	reg := cpu.getRptr(rIdx)
	addr := cpu.getIXYn(cpu.fetched.n)
	mw1 := newMW(addr, *reg, nil)
	cpu.scheduler.append(mw1)
}

func incSS(cpu *lr35902) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	v := reg.Get()
	v++
	reg.Set(v)
}

func decSS(cpu *lr35902) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	v := reg.Get()
	v--
	reg.Set(v)
}

func incR(cpu *lr35902) {
	rIdx := cpu.fetched.opCode >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	cpu.incR(r)
}

func (cpu *lr35902) incR(r *byte) {
	*r++
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.H = *r&0x0f == 0
	cpu.regs.F.N = false
	// panic(fmt.Sprintf("%08b", *r&0x0f))
}

func incHL(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(),
		func(cpu *lr35902, data uint8) {
			r := data
			r++
			mw := newMW(cpu.regs.HL.Get(), r, nil)
			cpu.regs.F.Z = r == 0
			cpu.regs.F.H = r&0x0f == 0
			cpu.regs.F.N = false

			cpu.scheduler.append(mw)
		},
	)
	cpu.scheduler.append(mr)
}

func incIXYd(cpu *lr35902) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := newMR(addr,
		func(cpu *lr35902, data uint8) {
			r := data
			r++
			mw := newMW(addr, r, nil)
			cpu.regs.F.Z = r == 0
			cpu.regs.F.H = r&0x0f == 0
			cpu.regs.F.N = false

			cpu.scheduler.append(mw)
		},
	)
	cpu.scheduler.append(mr)
}

func decR(cpu *lr35902) {
	rIdx := cpu.fetched.opCode >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	cpu.decR(r)
}

func (cpu *lr35902) decR(r *byte) {
	cpu.regs.F.H = *r&0x0f == 0
	*r--
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.N = true
}

func addAr(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.addA(*r)
}

func adcAr(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.adcA(*r)
}

func subAr(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.subA(*r)
}

func sbcAr(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.sbcA(*r)
}

func andAr(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.and(*r)
}

func orAr(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.or(*r)
}

func xorAr(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.xor(*r)
}

func cpR(cpu *lr35902) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.cp(*r)
}

func addAhl(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(), addAhl_m1)
	cpu.scheduler.append(mr)
}

func addAhl_m1(cpu *lr35902, data uint8) { cpu.addA(data) }

func subAhl(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(), subAhl_m1)
	cpu.scheduler.append(mr)
}

func subAhl_m1(cpu *lr35902, data uint8) { cpu.subA(data) }

func sbcAhl(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(), sbcAhl_m1)
	cpu.scheduler.append(mr)
}

func sbcAhl_m1(cpu *lr35902, data uint8) { cpu.sbcA(data) }

func adcAhl(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(), adcAhl_m1)
	cpu.scheduler.append(mr)
}

func adcAhl_m1(cpu *lr35902, data uint8) { cpu.adcA(data) }

func andAhl(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(), andAhl_m1)
	cpu.scheduler.append(mr)
}

func andAhl_m1(cpu *lr35902, data uint8) { cpu.and(data) }

func orAhl(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(), orAhl_m1)
	cpu.scheduler.append(mr)
}

func orAhl_m1(cpu *lr35902, data uint8) { cpu.or(data) }

func xorAhl(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(), xorAhl_m1)
	cpu.scheduler.append(mr)
}

func xorAhl_m1(cpu *lr35902, data uint8) { cpu.xor(data) }

func cpHl(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(), cpHl_m1)
	cpu.scheduler.append(mr)
}

func cpHl_m1(cpu *lr35902, data uint8) { cpu.cp(data) }

func decHL(cpu *lr35902) {
	mr := newMR(cpu.regs.HL.Get(),
		func(cpu *lr35902, data uint8) {
			r := data
			cpu.regs.F.H = r&0x0f == 0
			r--
			cpu.regs.F.Z = r == 0
			cpu.regs.F.N = true

			mw := newMW(cpu.regs.HL.Get(), r, nil)
			cpu.scheduler.append(mw)
		},
	)
	cpu.scheduler.append(mr)
}
