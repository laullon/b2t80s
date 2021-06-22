package z80

import (
	cpuUtils "github.com/laullon/b2t80s/cpu"
)

func ini(cpu *z80) { // TODO review tests changes
	in := &in{from: cpu.regs.BC.Get(), f: ini_m1}
	cpu.scheduler.append(in)
}

func ini_m1(cpu *z80, data uint8) {
	mw := cpu.newMW(cpu.regs.HL.Get(), data, ini_m2)
	cpu.scheduler.append(&exec{l: 1}, mw)
	if cpu.fetched.opCode > 0xAF {
		cpu.scheduler.append(&exec{l: 5, f: ini_m3})
	}
}

func ini_m2(cpu *z80) {
	cpu.regs.B--
	cpu.regs.HL.Set(cpu.regs.HL.Get() + 1)
	cpu.regs.F.N = true
	cpu.regs.F.Z = cpu.regs.B == 0
}

func ini_m3(cpu *z80) {
	if cpu.regs.B != 0 {
		cpu.regs.PC = cpu.regs.PC - 2
	}
}

func ind(cpu *z80) { // TODO review tests changes
	in := &in{from: cpu.regs.BC.Get(), f: ind_m1}
	cpu.scheduler.append(&exec{l: 1}, in)
}

func ind_m1(cpu *z80, data uint8) {
	hl := cpu.regs.HL.Get()
	mw := cpu.newMW(hl, data, ind_m2)
	cpu.scheduler.append(mw)
}

func ind_m2(cpu *z80) {
	hl := cpu.regs.HL.Get()
	cpu.regs.B--
	cpu.regs.HL.Set(hl - 1)
	cpu.regs.F.N = true
	cpu.regs.F.Z = cpu.regs.B == 0
	if cpu.fetched.opCode > 0xAF {
		cpu.scheduler.append(&exec{l: 5, f: ind_m3})
	}
}

func ind_m3(cpu *z80) {
	if cpu.regs.B != 0 {
		cpu.regs.PC = cpu.regs.PC - 2
	}
}

func outi(cpu *z80) { // TODO review tests changes
	mr := cpu.newMR(cpu.regs.HL.Get(), outi_m2)
	cpu.scheduler.append(&exec{l: 1}, mr)
}

func outi_m2(cpu *z80, data uint8) {
	cpu.regs.B--
	out := &out{addr: cpu.regs.BC.Get(), data: data, f: outi_m3}
	cpu.scheduler.append(&exec{l: 1}, out)
	if cpu.fetched.opCode > 0xAF {
		cpu.scheduler.append(&exec{l: 5, f: func(cpu *z80) {
			if cpu.regs.B != 0 {
				cpu.regs.PC = cpu.regs.PC - 2
			}
		}})
	}
}

func outi_m3(cpu *z80) {
	cpu.regs.HL.Set(cpu.regs.HL.Get() + 1)
	cpu.regs.F.Z = cpu.regs.B == 0
	cpu.regs.F.S = cpu.regs.B&0x80 != 0
	cpu.regs.F.N = cpu.regs.B&0x80 == 0
	cpu.regs.F.H = true
	cpu.regs.F.P = parityTable[cpu.regs.B]
}

func outd(cpu *z80) { // TODO review tests changes
	mr := cpu.newMR(cpu.regs.HL.Get(), outd_m1)
	cpu.scheduler.append(&exec{l: 1}, mr)
}

func outd_m1(cpu *z80, data uint8) {
	cpu.regs.B--
	out := &out{addr: cpu.regs.BC.Get(), data: data, f: outd_m2}
	cpu.scheduler.append(&exec{l: 1}, out)
	if cpu.fetched.opCode > 0xAF {
		cpu.scheduler.append(&exec{l: 5, f: outd_m3})
	}
}

func outd_m2(cpu *z80) {
	cpu.regs.HL.Set(cpu.regs.HL.Get() - 1)
	cpu.regs.F.Z = cpu.regs.B == 0
	cpu.regs.F.S = cpu.regs.B&0x80 != 0
	cpu.regs.F.N = cpu.regs.B&0x80 == 0
	cpu.regs.F.H = true
	cpu.regs.F.P = parityTable[cpu.regs.B]
}

func outd_m3(cpu *z80) {
	if cpu.regs.B != 0 {
		cpu.regs.PC = cpu.regs.PC - 2
	}
}

func cpi(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(), cpi_m1)
	cpu.scheduler.append(mr)
}

func cpi_m1(cpu *z80, data uint8) {

	val := data
	cpu.cpi_result = cpu.regs.A - val
	lookup := (cpu.regs.A&0x08)>>3 | (val&0x08)>>2 | (cpu.cpi_result&0x08)>>1
	cpu.regs.F.H = halfcarrySubTable[lookup]

	cpu.scheduler.append(&exec{l: 5, f: cpi_m2})
}

func cpi_m2(cpu *z80) {
	bc := cpu.regs.BC.Get()
	bc--
	cpu.regs.BC.Set(bc)
	cpu.regs.HL.Set(cpu.regs.HL.Get() + 1)

	cpu.regs.F.S = cpu.cpi_result&0x80 != 0
	cpu.regs.F.Z = cpu.cpi_result == 0
	cpu.regs.F.P = bc != 0
	cpu.regs.F.N = true
	if cpu.fetched.opCode > 0xAF {
		cpu.scheduler.append(&exec{l: 5, f: cpi_m3})
	}
}

func cpi_m3(cpu *z80) {
	if (cpu.regs.BC.Get()) != 0 && (cpu.cpi_result != 0) {
		cpu.regs.PC = cpu.regs.PC - 2
	}
}

func cpd(cpu *z80) {
	hl := cpu.regs.HL.Get()
	mr := cpu.newMR(hl, cpd_m1)
	cpu.scheduler.append(mr)
}

func cpd_m1(cpu *z80, data uint8) {
	val := data
	result := cpu.regs.A - byte(val)
	lookup := (cpu.regs.A&0x08)>>3 | (val&0x08)>>2 | (result&0x08)>>1

	cpu.regs.F.S = result&0x80 != 0
	cpu.regs.F.Z = result == 0
	cpu.regs.F.H = halfcarrySubTable[lookup]

	cpu.scheduler.append(&exec{l: 5, f: cpd_m2})
}

func cpd_m2(cpu *z80) {
	bc := cpu.regs.BC.Get()
	hl := cpu.regs.HL.Get()

	bc--
	hl--

	cpu.regs.BC.Set(bc)
	cpu.regs.HL.Set(hl)

	cpu.regs.F.P = bc != 0
	cpu.regs.F.N = true

	if cpu.fetched.opCode > 0xAF {
		cpu.scheduler.append(&exec{l: 5, f: cpd_m3})
	}
}

func cpd_m3(cpu *z80) {
	if (cpu.regs.BC.Get() != 0) && (!cpu.regs.F.Z) {
		cpu.regs.PC = cpu.regs.PC - 2
	}
}

func ldi(cpu *z80) {
	hl := cpu.regs.HL.Get()
	mr := cpu.newMR(hl, ldi_m1)
	cpu.scheduler.append(mr)
}

func ldi_m1(cpu *z80, data uint8) {
	v := data
	de := cpu.regs.DE.Get()
	mw := cpu.newMW(de, v, ldi_m2)
	cpu.scheduler.append(&exec{l: 2}, mw)
}

func ldi_m2(cpu *z80) {
	bc := cpu.regs.BC.Get()
	de := cpu.regs.DE.Get()
	hl := cpu.regs.HL.Get()

	bc--
	de++
	hl++

	cpu.regs.BC.Set(bc)
	cpu.regs.DE.Set(de)
	cpu.regs.HL.Set(hl)

	cpu.regs.F.P = bc != 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
	if cpu.fetched.opCode > 0xAF {
		cpu.scheduler.append(&exec{l: 5, f: ldi_m3})
	}
}

func ldi_m3(cpu *z80) {
	if cpu.regs.BC.Get() != 0 {
		cpu.regs.PC = cpu.regs.PC - 2
	}
}

func ldd(cpu *z80) {
	hl := cpu.regs.HL.Get()

	mr := cpu.newMR(hl, ldd_m1)
	cpu.scheduler.append(mr)
}

func ldd_m1(cpu *z80, data uint8) {
	de := cpu.regs.DE.Get()
	v := data
	mw := cpu.newMW(de, v, ldd_m2)
	cpu.scheduler.append(&exec{l: 2}, mw)
}

func ldd_m2(cpu *z80) {
	bc := cpu.regs.BC.Get()
	de := cpu.regs.DE.Get()
	hl := cpu.regs.HL.Get()

	bc--
	de--
	hl--

	cpu.regs.BC.Set(bc)
	cpu.regs.DE.Set(de)
	cpu.regs.HL.Set(hl)

	cpu.regs.F.P = bc != 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
	if cpu.fetched.opCode > 0xAF {
		cpu.scheduler.append(&exec{l: 5, f: ldd_m3})
	}
}

func ldd_m3(cpu *z80) {
	if cpu.regs.BC.Get() != 0 {
		cpu.regs.PC = cpu.regs.PC - 2
	}
}

func exSP(cpu *z80) {
	reg := cpu.indexRegs[cpu.indexIdx]
	mr1 := cpu.newMR(cpu.regs.SP.Get(), exSP_m1)
	mr2 := cpu.newMR(cpu.regs.SP.Get()+1, exSP_m2)
	mw1 := cpu.newMW(cpu.regs.SP.Get(), *reg.L, nil)
	mw2 := cpu.newMW(cpu.regs.SP.Get()+1, *reg.H, exSP_m3)
	cpu.scheduler.append(mr1, &exec{l: 1}, mr2, mw1, &exec{l: 2}, mw2)
}

func exSP_m1(cpu *z80, data uint8) { cpu.spv = uint16(data) }
func exSP_m2(cpu *z80, data uint8) { cpu.spv |= uint16(data) << 8 }
func exSP_m3(cpu *z80)             { reg := cpu.indexRegs[cpu.indexIdx]; reg.Set(cpu.spv) }

func addIXY(cpu *z80) {
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
	result := uint32(ix) + uint32(reg.Get())
	lookup := byte(((ix & 0x0800) >> 11) | ((reg.Get() & 0x0800) >> 10) | ((uint16(result) & 0x0800) >> 9))
	regI.Set(uint16(result))

	cpu.regs.F.N = false
	cpu.regs.F.H = halfcarryAddTable[lookup]
	cpu.regs.F.C = (result & 0x10000) != 0
}

func addIY(cpu *z80) {
	var reg *cpuUtils.RegPair
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	switch rIdx {
	case 0b00:
		reg = cpu.regs.BC
	case 0b01:
		reg = cpu.regs.DE
	case 0b10:
		reg = cpu.regs.IY
	case 0b11:
		reg = cpu.regs.SP
	}

	iy := cpu.regs.IY.Get()
	result := uint32(iy) + uint32(reg.Get())
	lookup := byte(((iy & 0x0800) >> 11) | ((reg.Get() & 0x0800) >> 10) | ((uint16(result) & 0x0800) >> 9))
	cpu.regs.IY.Set(uint16(result))

	cpu.regs.F.N = false
	cpu.regs.F.H = halfcarryAddTable[lookup]
	cpu.regs.F.C = (result & 0x10000) != 0
}

func outNa(cpu *z80) {
	port := toWord(cpu.fetched.n, cpu.regs.A)
	cpu.scheduler.append(&out{addr: port, data: cpu.regs.A})
}

func inAn(cpu *z80) {
	cpu.inAn_f = cpu.regs.F.GetByte()
	port := toWord(cpu.fetched.n, cpu.regs.A)
	cpu.scheduler.append(&in{from: port, f: inAn_m1})
}

func inAn_m1(cpu *z80, data uint8) {
	cpu.regs.A = data
	cpu.regs.F.SetByte(cpu.inAn_f)
}

func inRc(cpu *z80) {
	cpu.scheduler.append(&in{from: cpu.regs.BC.Get(), f: inRc_m1})
}

func inRc_m1(cpu *z80, data uint8) {
	rIdx := cpu.fetched.opCode >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	*r = data
}

func inC(cpu *z80) {
	cpu.scheduler.append(&in{from: cpu.regs.BC.Get(), f: nil})
}

func outCr(cpu *z80) {
	rIdx := cpu.fetched.opCode >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	cpu.scheduler.append(&out{addr: cpu.regs.BC.Get(), data: *r})
}

func outC0(cpu *z80) {
	cpu.scheduler.append(&out{addr: cpu.regs.BC.Get(), data: 0})
}

func retCC(cpu *z80) {
	ccIdx := cpu.fetched.opCode >> 3 & 0b111
	branch := cpu.checkCondition(ccIdx)
	if branch {
		cpu.popFromStack(func(cpu *z80, data uint16) {
			cpu.regs.PC = data
		})
	}
}

func ret(cpu *z80) {
	cpu.popFromStack(func(cpu *z80, data uint16) {
		cpu.regs.PC = data
	})
}

func rstP(cpu *z80) {
	cpu.pushToStack(cpu.regs.PC, rstP_m1)
}

func rstP_m1(cpu *z80) {
	newPCs := []uint16{0x00, 0x08, 0x10, 0x18, 0x20, 0x28, 0x30, 0x38}
	pIdx := cpu.fetched.opCode >> 3 & 0b111
	cpu.regs.PC = newPCs[pIdx]
}

func jpCC(cpu *z80) {
	ccIdx := cpu.fetched.opCode >> 3 & 0b111
	branch := cpu.checkCondition(ccIdx)
	if branch {
		cpu.regs.PC = cpu.fetched.nn
	}
}

func callCC(cpu *z80) {
	ccIdx := cpu.fetched.opCode >> 3 & 0b111
	branch := cpu.checkCondition(ccIdx)
	if branch {
		cpu.scheduler.append(&exec{l: 1, f: callCC_m2})
	}
}

func callCC_m2(cpu *z80) { cpu.pushToStack(cpu.regs.PC, call_m1) }

func call(cpu *z80) {
	cpu.pushToStack(cpu.regs.PC, call_m1)
}

func call_m1(cpu *z80) { cpu.regs.PC = cpu.fetched.nn }

func (cpu *z80) checkCondition(ccIdx byte) bool {
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
	case 4:
		res = !cpu.regs.F.P
	case 5:
		res = cpu.regs.F.P
	case 6:
		res = !cpu.regs.F.S
	case 7:
		res = cpu.regs.F.S
	}
	return res
}

func (cpu *z80) pushToStack(data uint16, f func(cpu *z80)) {
	cpu.pushF = f
	push1 := cpu.newMW(cpu.regs.SP.Get()-1, uint8(data>>8), nil)
	push2 := cpu.newMW(cpu.regs.SP.Get()-2, uint8(data), push_m1)
	cpu.scheduler.append(push1, push2)
}

func push_m1(cpu *z80) {
	cpu.regs.SP.Set(cpu.regs.SP.Get() - 2)
	if cpu.pushF != nil {
		cpu.pushF(cpu)
	}
}

func (cpu *z80) popFromStack(f func(cpu *z80, data uint16)) {
	cpu.popF = f
	pop1 := cpu.newMR(cpu.regs.SP.Get(), pop_m1)
	pop2 := cpu.newMR(cpu.regs.SP.Get()+1, pop_m2)
	cpu.scheduler.append(pop1, pop2)
}

func pop_m1(cpu *z80, data uint8) { cpu.popData = uint16(data) }

func pop_m2(cpu *z80, data uint8) {
	cpu.popData |= (uint16(data) << 8)
	cpu.regs.SP.Set(cpu.regs.SP.Get() + 2)
	cpu.popF(cpu, cpu.popData)
}

func popSS(cpu *z80) {
	cpu.popFromStack(popSS_m1)
}

func popSS_m1(cpu *z80, data uint16) {
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

func pushSS(cpu *z80) {
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

func ldDDmm(cpu *z80) {
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

func ldBCa(cpu *z80) {
	pos := cpu.regs.BC.Get()
	cpu.scheduler.append(cpu.newMW(pos, cpu.regs.A, nil))
}

func ldDEa(cpu *z80) {
	pos := cpu.regs.DE.Get()
	cpu.scheduler.append(cpu.newMW(pos, cpu.regs.A, nil))
}

func ldNNhl(cpu *z80) {
	mm := cpu.fetched.nn
	mw1 := cpu.newMW(mm, cpu.regs.L, nil)
	mw2 := cpu.newMW(mm+1, cpu.regs.H, nil)
	cpu.scheduler.append(mw1, mw2)
}

func ldNNIXY(cpu *z80) {
	reg := cpu.indexRegs[cpu.indexIdx]
	mm := cpu.fetched.nn
	mw1 := cpu.newMW(mm, *reg.L, nil)
	mw2 := cpu.newMW(mm+1, *reg.H, nil)
	cpu.scheduler.append(mw1, mw2)
}

func ldNNa(cpu *z80) {
	mm := cpu.fetched.nn
	mw1 := cpu.newMW(mm, cpu.regs.A, nil)
	cpu.scheduler.append(mw1)
}

func rrd(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(), rrd_m1)
	cpu.scheduler.append(mr)
}

func rrd_m1(cpu *z80, data uint8) {
	cpu.hlv = data
	mw := cpu.newMW(cpu.regs.HL.Get(), (cpu.regs.A<<4 | cpu.hlv>>4), rrd_m2)
	cpu.scheduler.append(&exec{l: 4}, mw)
}

func rrd_m2(cpu *z80) {
	cpu.regs.A = (cpu.regs.A & 0xf0) | (cpu.hlv & 0x0f)
	cpu.regs.F.S = cpu.regs.A&0x80 != 0
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.P = parityTable[cpu.regs.A]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func rld(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(), rld_m1)
	cpu.scheduler.append(mr)
}

func rld_m1(cpu *z80, data uint8) {
	cpu.hlv = data
	mw := cpu.newMW(cpu.regs.HL.Get(), (cpu.hlv<<4 | cpu.regs.A&0x0f), rld_m2)
	cpu.scheduler.append(&exec{l: 4}, mw)
}

func rld_m2(cpu *z80) {
	cpu.regs.A = (cpu.regs.A & 0xf0) | (cpu.hlv >> 4)
	cpu.regs.F.S = cpu.regs.A&0x80 != 0
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.P = parityTable[cpu.regs.A]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func ldNNdd(cpu *z80) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	mm := cpu.fetched.nn
	mw1 := cpu.newMW(mm, *reg.L, nil)
	mw2 := cpu.newMW(mm+1, *reg.H, nil)
	cpu.scheduler.append(mw1, mw2)
}

func ldDDnn(cpu *z80) {
	mm := cpu.fetched.nn
	mr1 := cpu.newMR(mm, ldDDnn_m1)
	mr2 := cpu.newMR(mm+1, ldDDnn_m2)
	cpu.scheduler.append(mr1, mr2)
}

func ldDDnn_m1(cpu *z80, data uint8) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	*reg.L = data
}
func ldDDnn_m2(cpu *z80, data uint8) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	*reg.H = data
}

func ldAi(cpu *z80) {
	cpu.regs.A = cpu.regs.I
	cpu.regs.F.S = cpu.regs.A&0x80 != 0
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = false
	cpu.regs.F.P = cpu.regs.IFF2
	cpu.regs.F.N = false
}

func ldAr(cpu *z80) {
	cpu.regs.A = cpu.regs.R
	cpu.regs.F.S = cpu.regs.R&0x80 != 0
	cpu.regs.F.Z = cpu.regs.R == 0
	cpu.regs.F.H = false
	cpu.regs.F.P = cpu.regs.IFF2
	cpu.regs.F.N = false
}

func ldHLnn(cpu *z80) {
	mm := cpu.fetched.nn
	mr1 := cpu.newMR(mm, ldHLnn_m1)
	mr2 := cpu.newMR(mm+1, ldHLnn_m2)
	cpu.scheduler.append(mr1, mr2)
}

func ldHLnn_m1(cpu *z80, data uint8) { cpu.regs.L = data }
func ldHLnn_m2(cpu *z80, data uint8) { cpu.regs.H = data }

func ldIXYnn(cpu *z80) {
	mm := cpu.fetched.nn
	mr1 := cpu.newMR(mm, ldIXYnn_m1)
	mr2 := cpu.newMR(mm+1, ldIXYnn_m2)
	cpu.scheduler.append(mr1, mr2)
}

func ldIXYnn_m1(cpu *z80, data uint8) {
	reg := cpu.indexRegs[cpu.indexIdx]
	*reg.L = data
}

func ldIXYnn_m2(cpu *z80, data uint8) {
	reg := cpu.indexRegs[cpu.indexIdx]
	*reg.H = data
}

func ldAnn(cpu *z80) {
	mm := cpu.fetched.nn
	mr1 := cpu.newMR(mm, ldAnn_n1)
	cpu.scheduler.append(mr1)
}

func ldAnn_n1(cpu *z80, data uint8) { cpu.regs.A = data }

func ldHLn(cpu *z80) {
	mw1 := cpu.newMW(cpu.regs.HL.Get(), cpu.fetched.n, nil)
	cpu.scheduler.append(mw1)
}

func ldIXYdN(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mw1 := cpu.newMW(addr, cpu.fetched.n2, nil)
	cpu.scheduler.append(mw1)
}

func ldIXYdR(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	reg := cpu.getRptr(rIdx)
	addr := cpu.getIXYn(cpu.fetched.n)
	mw1 := cpu.newMW(addr, *reg, nil)
	cpu.scheduler.append(mw1)
}

func incSS(cpu *z80) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	v := reg.Get()
	v++
	reg.Set(v)
}

func decSS(cpu *z80) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	v := reg.Get()
	v--
	reg.Set(v)
}

func incR(cpu *z80) {
	rIdx := cpu.fetched.opCode >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	cpu.incR(r)
}

func (cpu *z80) incR(r *byte) {
	*r++
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.H = *r&0x0f == 0
	cpu.regs.F.P = *r == 0x80
	cpu.regs.F.N = false
	// panic(fmt.Sprintf("%08b", *r&0x0f))
}

func incHL(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(),
		func(cpu *z80, data uint8) {
			r := data
			r++
			mw := cpu.newMW(cpu.regs.HL.Get(), r, nil)
			cpu.regs.F.S = r&0x80 != 0
			cpu.regs.F.Z = r == 0
			cpu.regs.F.H = r&0x0f == 0
			cpu.regs.F.P = r == 0x80
			cpu.regs.F.N = false

			cpu.scheduler.append(mw)
		},
	)
	cpu.scheduler.append(mr)
}

func incIXYd(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := cpu.newMR(addr,
		func(cpu *z80, data uint8) {
			r := data
			r++
			mw := cpu.newMW(addr, r, nil)
			cpu.regs.F.S = r&0x80 != 0
			cpu.regs.F.Z = r == 0
			cpu.regs.F.H = r&0x0f == 0
			cpu.regs.F.P = r == 0x80
			cpu.regs.F.N = false

			cpu.scheduler.append(mw)
		},
	)
	cpu.scheduler.append(mr)
}

func decR(cpu *z80) {
	rIdx := cpu.fetched.opCode >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	cpu.decR(r)
}

func (cpu *z80) decR(r *byte) {
	cpu.regs.F.H = *r&0x0f == 0
	*r--
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = *r == 0x7f
	cpu.regs.F.N = true
}

func addAr(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.addA(*r)
}

func adcAr(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.adcA(*r)
}

func subAr(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.subA(*r)
}

func sbcAr(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.sbcA(*r)
}

func andAr(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.and(*r)
}

func orAr(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.or(*r)
}

func xorAr(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.xor(*r)
}

func cpR(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	cpu.cp(*r)
}

func addAhl(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(), addAhl_m1)
	cpu.scheduler.append(mr)
}

func addAhl_m1(cpu *z80, data uint8) { cpu.addA(data) }

func subAhl(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(), subAhl_m1)
	cpu.scheduler.append(mr)
}

func subAhl_m1(cpu *z80, data uint8) { cpu.subA(data) }

func sbcAhl(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(), sbcAhl_m1)
	cpu.scheduler.append(mr)
}

func sbcAhl_m1(cpu *z80, data uint8) { cpu.sbcA(data) }

func adcAhl(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(), adcAhl_m1)
	cpu.scheduler.append(mr)
}

func adcAhl_m1(cpu *z80, data uint8) { cpu.adcA(data) }

func andAhl(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(), andAhl_m1)
	cpu.scheduler.append(mr)
}

func andAhl_m1(cpu *z80, data uint8) { cpu.and(data) }

func orAhl(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(), orAhl_m1)
	cpu.scheduler.append(mr)
}

func orAhl_m1(cpu *z80, data uint8) { cpu.or(data) }

func xorAhl(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(), xorAhl_m1)
	cpu.scheduler.append(mr)
}

func xorAhl_m1(cpu *z80, data uint8) { cpu.xor(data) }

func cpHl(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(), cpHl_m1)
	cpu.scheduler.append(mr)
}

func cpHl_m1(cpu *z80, data uint8) { cpu.cp(data) }

func decHL(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(),
		func(cpu *z80, data uint8) {
			r := data
			cpu.regs.F.H = r&0x0f == 0
			r--
			cpu.regs.F.S = r&0x80 != 0
			cpu.regs.F.Z = r == 0
			cpu.regs.F.P = r == 0x7f
			cpu.regs.F.N = true

			mw := cpu.newMW(cpu.regs.HL.Get(), r, nil)
			cpu.scheduler.append(mw)
		},
	)
	cpu.scheduler.append(mr)
}

func decIXYd(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := cpu.newMR(addr,
		func(cpu *z80, data uint8) {
			r := data
			cpu.regs.F.H = r&0x0f == 0
			r--
			cpu.regs.F.S = r&0x80 != 0
			cpu.regs.F.Z = r == 0
			cpu.regs.F.P = r == 0x7f
			cpu.regs.F.N = true

			mw := cpu.newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		},
	)
	cpu.scheduler.append(mr)
}

func ldRn(cpu *z80) {
	rIdx := cpu.fetched.opCode >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	*r = cpu.fetched.n
}

func ldRhl(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(), ldR_m1)
	cpu.scheduler.append(mr)
}

func ldRixyD(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := cpu.newMR(addr, ldR_m1)
	cpu.scheduler.append(mr)
}

func ldR_m1(cpu *z80, data uint8) {
	rIdx := cpu.fetched.opCode >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	*r = data
}

func addAixyD(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := cpu.newMR(addr, addAixyD_m1)
	cpu.scheduler.append(mr)
}

func addAixyD_m1(cpu *z80, data uint8) { cpu.addA(data) }

func adcAixyD(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := cpu.newMR(addr, adcAixyD_m1)
	cpu.scheduler.append(mr)
}

func adcAixyD_m1(cpu *z80, data uint8) { cpu.adcA(data) }

func subAixyD(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := cpu.newMR(addr, subAixyD_m1)
	cpu.scheduler.append(mr)
}

func subAixyD_m1(cpu *z80, data uint8) { cpu.subA(data) }

func sbcAixyD(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := cpu.newMR(addr, sbcAixyD_m1)
	cpu.scheduler.append(mr)
}

func sbcAixyD_m1(cpu *z80, data uint8) { cpu.sbcA(data) }

func andAixyD(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := cpu.newMR(addr, andAixyD_m1)
	cpu.scheduler.append(mr)
}

func andAixyD_m1(cpu *z80, data uint8) { cpu.and(data) }

func xorAixyD(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := cpu.newMR(addr, xorAixyD_m1)
	cpu.scheduler.append(mr)
}

func xorAixyD_m1(cpu *z80, data uint8) { cpu.xor(data) }

func cpAixyD(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := cpu.newMR(addr, cpAixyD_m1)
	cpu.scheduler.append(mr)
}

func cpAixyD_m1(cpu *z80, data uint8) { cpu.cp(data) }

func orAixyD(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := cpu.newMR(addr, orAixyD_m1)
	cpu.scheduler.append(mr)
}

func orAixyD_m1(cpu *z80, data uint8) { cpu.or(data) }

func ldHLr(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	mr := cpu.newMW(cpu.regs.HL.Get(), *r, nil)
	cpu.scheduler.append(mr)
}

func ldIXYHr(cpu *z80) {
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

func ldIXYLr(cpu *z80) {
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

func ldRr(cpu *z80) {
	r1Idx := cpu.fetched.opCode >> 3 & 0b111
	r2Idx := cpu.fetched.opCode & 0b111
	r1 := cpu.getRptr(r1Idx)
	r2 := cpu.getRptr(r2Idx)
	*r1 = *r2
}

func rlca(cpu *z80) {
	cpu.regs.A = cpu.regs.A<<1 | cpu.regs.A>>7
	cpu.regs.F.C = cpu.regs.A&0x01 != 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func rla(cpu *z80) {
	c := cpu.regs.F.C
	cpu.regs.F.C = cpu.regs.A&0b10000000 != 0
	cpu.regs.A = (cpu.regs.A << 1)
	if c {
		cpu.regs.A |= 1
	}
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

var cbFuncs = []func(cpu *z80, r *byte){rlc, rrc, rl, rr, sla, sra, sll, srl}

func cbR(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	fIdx := cpu.fetched.opCode >> 3
	cbFuncs[fIdx](cpu, r)
}

func cbHL(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(),
		func(cpu *z80, data uint8) {
			b := data
			fIdx := cpu.fetched.opCode >> 3
			cbFuncs[fIdx](cpu, &b)
			mw := cpu.newMW(cpu.regs.HL.Get(), b, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func cbIXYdr(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := cpu.newMR(addr,
		func(cpu *z80, data uint8) {
			r := data
			fIdx := (cpu.fetched.opCode >> 3) & 0b111
			cbFuncs[fIdx](cpu, &r)

			rIdx := cpu.fetched.opCode & 0b111
			reg := cpu.getRptr(rIdx)
			*reg = r

			mw := cpu.newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func cbIXYd(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	mr := cpu.newMR(addr,
		func(cpu *z80, data uint8) {
			r := data
			fIdx := (cpu.fetched.opCode >> 3) & 0b111
			cbFuncs[fIdx](cpu, &r)
			mw := cpu.newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func bit(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	b := (cpu.fetched.opCode >> 3) & 0b111
	cpu.bit(b, *r)
}

func bitIXYd(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	b := (cpu.fetched.opCode >> 3) & 0b111
	mr := cpu.newMR(addr,
		func(cpu *z80, data uint8) {
			r := data
			cpu.bit(b, r)
		})
	cpu.scheduler.append(mr)
}

func bitHL(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(),
		func(cpu *z80, data uint8) {
			v := data
			b := (cpu.fetched.opCode >> 3) & 0b111
			cpu.bit(b, v)
		})
	cpu.scheduler.append(mr)
}

func res(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	b := (cpu.fetched.opCode >> 3) & 0b111
	cpu.res(b, r)
}

func resHL(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(),
		func(cpu *z80, data uint8) {
			v := data
			b := (cpu.fetched.opCode >> 3) & 0b111
			cpu.res(b, &v)
			mw := cpu.newMW(cpu.regs.HL.Get(), v, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func resIXYdR(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	b := (cpu.fetched.opCode >> 3) & 0b111
	mr := cpu.newMR(addr,
		func(cpu *z80, data uint8) {
			r := data
			cpu.res(b, &r)

			rIdx := cpu.fetched.opCode & 0b111
			reg := cpu.getRptr(rIdx)
			*reg = r

			mw := cpu.newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func resIXYd(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	b := (cpu.fetched.opCode >> 3) & 0b111
	mr := cpu.newMR(addr,
		func(cpu *z80, data uint8) {
			r := data
			cpu.res(b, &r)
			mw := cpu.newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func set(cpu *z80) {
	rIdx := cpu.fetched.opCode & 0b111
	r := cpu.getRptr(rIdx)
	b := (cpu.fetched.opCode >> 3) & 0b111
	cpu.set(b, r)
}

func setHL(cpu *z80) {
	mr := cpu.newMR(cpu.regs.HL.Get(),
		func(cpu *z80, data uint8) {
			v := data
			b := (cpu.fetched.opCode >> 3) & 0b111
			cpu.set(b, &v)
			mw := cpu.newMW(cpu.regs.HL.Get(), v, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func setIXYdR(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	b := (cpu.fetched.opCode >> 3) & 0b111
	mr := cpu.newMR(addr,
		func(cpu *z80, data uint8) {
			r := data
			cpu.set(b, &r)

			rIdx := cpu.fetched.opCode & 0b111
			reg := cpu.getRptr(rIdx)
			*reg = r

			mw := cpu.newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func setIXYd(cpu *z80) {
	addr := cpu.getIXYn(cpu.fetched.n)
	b := (cpu.fetched.opCode >> 3) & 0b111
	mr := cpu.newMR(addr,
		func(cpu *z80, data uint8) {
			r := data
			cpu.set(b, &r)
			mw := cpu.newMW(addr, r, nil)
			cpu.scheduler.append(mw)
		})
	cpu.scheduler.append(mr)
}

func rrca(cpu *z80) {
	cpu.regs.F.C = cpu.regs.A&0x01 != 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
	cpu.regs.A = (cpu.regs.A >> 1) | (cpu.regs.A << 7)
}

func rra(cpu *z80) {
	c := cpu.regs.F.C
	cpu.regs.F.C = cpu.regs.A&1 != 0
	cpu.regs.A = (cpu.regs.A >> 1)
	if c {
		cpu.regs.A |= 0b10000000
	}
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func exafaf(cpu *z80) {
	cpu.regs.A, cpu.regs.Aalt = cpu.regs.Aalt, cpu.regs.A
	cpu.regs.F, cpu.regs.Falt = cpu.regs.Falt, cpu.regs.F
}

func exDEhl(cpu *z80) {
	cpu.regs.D, cpu.regs.H = cpu.regs.H, cpu.regs.D
	cpu.regs.E, cpu.regs.L = cpu.regs.L, cpu.regs.E
}

func exx(cpu *z80) {
	cpu.regs.B, cpu.regs.Balt = cpu.regs.Balt, cpu.regs.B
	cpu.regs.C, cpu.regs.Calt = cpu.regs.Calt, cpu.regs.C
	cpu.regs.D, cpu.regs.Dalt = cpu.regs.Dalt, cpu.regs.D
	cpu.regs.E, cpu.regs.Ealt = cpu.regs.Ealt, cpu.regs.E
	cpu.regs.H, cpu.regs.Halt = cpu.regs.Halt, cpu.regs.H
	cpu.regs.L, cpu.regs.Lalt = cpu.regs.Lalt, cpu.regs.L
}

func halt(cpu *z80) {
	cpu.halt = true
	cpu.regs.PC--
}

func addHLss(cpu *z80) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)

	hl := cpu.regs.HL.Get()
	result := uint32(hl) + uint32(reg.Get())
	lookup := byte(((hl & 0x0800) >> 11) | ((reg.Get() & 0x0800) >> 10) | ((uint16(result) & 0x0800) >> 9))
	cpu.regs.HL.Set(uint16(result))

	cpu.regs.F.N = false
	cpu.regs.F.H = halfcarryAddTable[lookup]
	cpu.regs.F.C = (result & 0x10000) != 0
}

func ldAbc(cpu *z80) {
	from := cpu.regs.BC.Get()
	mr := cpu.newMR(from, ldAbc_m1)
	cpu.scheduler.append(mr)
}

func ldAbc_m1(cpu *z80, data uint8) { cpu.regs.A = data }

func ldAde(cpu *z80) {
	from := cpu.regs.DE.Get()
	mr := cpu.newMR(from, ldAde_m1)
	cpu.scheduler.append(mr)
}

func ldAde_m1(cpu *z80, data uint8) { cpu.regs.A = data }

func djnz(cpu *z80) {
	cpu.regs.B--
	if cpu.regs.B != 0 {
		cpu.scheduler.append(&exec{l: 5, f: jr})
	}
}

func jrnz(cpu *z80) {
	if !cpu.regs.F.Z {
		cpu.scheduler.append(&exec{l: 5, f: jr})
	}
}

func jrnc(cpu *z80) {
	if !cpu.regs.F.C {
		cpu.scheduler.append(&exec{l: 5, f: jr})
	}
}

func jrc(cpu *z80) {
	if cpu.regs.F.C {
		cpu.scheduler.append(&exec{l: 5, f: jr})
	}
}

func jrz(cpu *z80) {
	if cpu.regs.F.Z {
		cpu.scheduler.append(&exec{l: 5, f: jr})
	}
}

func jr(cpu *z80) {
	jump := int8(cpu.fetched.n)
	cpu.regs.PC += uint16(jump)
}

func scf(cpu *z80) {
	cpu.regs.F.H = false
	cpu.regs.F.N = false
	cpu.regs.F.C = true
}

func ccf(cpu *z80) {
	cpu.regs.F.H = cpu.regs.F.C
	cpu.regs.F.N = false
	cpu.regs.F.C = !cpu.regs.F.C
}

func daa(cpu *z80) {
	c := cpu.regs.F.C
	add := byte(0)
	if cpu.regs.F.H || ((cpu.regs.A & 0x0f) > 9) {
		add = 6
	}
	if c || (cpu.regs.A > 0x99) {
		add |= 0x60
	}
	if cpu.regs.A > 0x99 {
		c = true
	}
	if cpu.regs.F.N {
		cpu.subA(add)
	} else {
		cpu.addA(add)
	}
	cpu.regs.F.S = int8(cpu.regs.A) < 0
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.P = parityTable[cpu.regs.A]
	cpu.regs.F.C = c
}

func cpl(cpu *z80) {
	cpu.regs.A = ^cpu.regs.A
	cpu.regs.F.H = true
	cpu.regs.F.N = true
}

// -------
func (cpu *z80) getRptr(rIdx byte) *byte {
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

func (cpu *z80) getRRptr(rIdx byte) *cpuUtils.RegPair {
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

func (cpu *z80) getIXYn(n byte) uint16 {
	reg := cpu.indexRegs[cpu.indexIdx]
	i := int16(int8(n))
	ix := reg.Get()
	ix = uint16(int16(ix) + i)
	return ix
}

func (cpu *z80) res(b byte, v *byte) {
	b = 1 << b
	*v &= ^b
}

func (cpu *z80) set(b byte, v *byte) {
	b = 1 << b
	*v |= b
}

func (cpu *z80) bit(b, v byte) {
	b = 1 << b
	v &= b
	cpu.regs.F.S = v&0x0080 != 0
	cpu.regs.F.Z = v == 0
	cpu.regs.F.H = true
	cpu.regs.F.P = parityTable[v]
	cpu.regs.F.N = false
}

func (cpu *z80) adcA(s byte) {
	res := int16(cpu.regs.A) + int16(s)
	if cpu.regs.F.C {
		res++
	}
	lookup := ((cpu.regs.A & 0x88) >> 3) | ((s & 0x88) >> 2) | ((byte(res) & 0x88) >> 1)
	cpu.regs.A = byte(res)
	cpu.regs.F.S = cpu.regs.A&0x80 != 0
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = halfcarryAddTable[lookup&0x07]
	cpu.regs.F.P = overflowAddTable[lookup>>4]
	cpu.regs.F.N = false
	cpu.regs.F.C = (res & 0x100) == 0x100
}

func (cpu *z80) adcHL(ss uint16) {
	hl := cpu.regs.HL.Get()
	res := int32(hl) + int32(ss)
	if cpu.regs.F.C {
		res++
	}
	lookup := byte(((hl & 0x8800) >> 11) | ((ss & 0x8800) >> 10) | ((uint16(res) & 0x8800) >> 9))
	hl = uint16(res)
	cpu.regs.HL.Set(hl)
	cpu.regs.F.S = cpu.regs.H&0x80 != 0
	cpu.regs.F.Z = hl == 0
	cpu.regs.F.H = halfcarryAddTable[lookup&0x07]
	cpu.regs.F.P = overflowAddTable[lookup>>4]
	cpu.regs.F.N = false
	cpu.regs.F.C = (res & 0x10000) != 0
}

func (cpu *z80) cp(r byte) {
	a := int16(cpu.regs.A)
	result := a - int16(r)
	lookup := ((cpu.regs.A & 0x88) >> 3) | (((r) & 0x88) >> 2) | ((byte(result) & 0x88) >> 1)

	cpu.regs.F.S = result&0x80 != 0
	cpu.regs.F.Z = result == 0
	cpu.regs.F.H = halfcarrySubTable[lookup&0x07]
	cpu.regs.F.P = overflowSubTable[lookup>>4]
	cpu.regs.F.N = true
	cpu.regs.F.C = ((result) & 0x100) == 0x100
}

func rlc(cpu *z80, r *byte) {
	*r = (*r << 1) | (*r >> 7)
	cpu.regs.F.C = *r&1 != 0
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func rrc(cpu *z80, r *byte) {
	cpu.regs.F.C = *r&1 != 0
	*r = (*r << 7) | (*r >> 1)
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func sll(cpu *z80, r *byte) {
	cpu.regs.F.C = *r&0x80 != 0
	*r = byte((*r << 1) | 0x01)
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func srl(cpu *z80, r *byte) {
	cpu.regs.F.C = *r&1 != 0
	*r = byte(*r >> 1)
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.N = false
	cpu.regs.F.H = false
}

func sla(cpu *z80, r *byte) {
	cpu.regs.F.C = *r&0x80 != 0
	*r = (*r << 1)
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func sra(cpu *z80, r *byte) {
	cpu.regs.F.C = *r&0x1 != 0
	b7 := *r & 0b10000000
	*r = (*r >> 1) | b7
	cpu.regs.F.S = *r&0x0080 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func rr(cpu *z80, r *byte) {
	c := cpu.regs.F.C
	cpu.regs.F.C = *r&0x1 != 0
	*r = (*r >> 1)
	if c {
		*r |= 0b10000000
	}
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func rl(cpu *z80, r *byte) {
	c := cpu.regs.F.C
	cpu.regs.F.C = *r&0x80 != 0
	*r = (*r << 1)
	if c {
		*r |= 0x1
	}
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

// ------

func (cpu *z80) addA(r byte) {
	a := int16(cpu.regs.A)
	result := a + int16(r)
	lookup := ((cpu.regs.A & 0x88) >> 3) | (((r) & 0x88) >> 2) | ((byte(result) & 0x88) >> 1)
	cpu.regs.A = uint8(result & 0x00ff)

	cpu.regs.F.S = cpu.regs.A&0x80 != 0
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = halfcarryAddTable[lookup&0x07]
	cpu.regs.F.P = overflowAddTable[lookup>>4]
	cpu.regs.F.N = false
	cpu.regs.F.C = ((result) & 0x100) != 0
}

func (cpu *z80) subA(r byte) {
	a := int16(cpu.regs.A)
	result := a - int16(r)
	lookup := ((cpu.regs.A & 0x88) >> 3) | (((r) & 0x88) >> 2) | ((byte(result) & 0x88) >> 1)
	cpu.regs.A = uint8(result & 0x00ff)

	cpu.regs.F.S = cpu.regs.A&0x80 != 0
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = halfcarrySubTable[lookup&0x07]
	cpu.regs.F.P = overflowSubTable[lookup>>4]
	cpu.regs.F.N = true
	cpu.regs.F.C = ((result) & 0x100) == 0x100
}

func (cpu *z80) xor(s uint8) {
	cpu.regs.A = cpu.regs.A ^ s
	cpu.regs.F.S = int8(cpu.regs.A) < 0
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = false
	cpu.regs.F.P = parityTable[cpu.regs.A]
	cpu.regs.F.N = false
	cpu.regs.F.C = false
}

func (cpu *z80) and(s uint8) {
	cpu.regs.A = cpu.regs.A & s
	cpu.regs.F.S = int8(cpu.regs.A) < 0
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = true
	cpu.regs.F.P = parityTable[cpu.regs.A]
	cpu.regs.F.N = false
	cpu.regs.F.C = false
}

func (cpu *z80) or(s uint8) {
	// TODO: review p/v flag
	cpu.regs.A = cpu.regs.A | s
	cpu.regs.F.S = int8(cpu.regs.A) < 0
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = false
	cpu.regs.F.P = parityTable[cpu.regs.A]
	cpu.regs.F.N = false
	cpu.regs.F.C = false
}

func (cpu *z80) sbcA(s byte) {
	res := uint16(cpu.regs.A) - uint16(s)
	if cpu.regs.F.C {
		res--
	}
	lookup := ((cpu.regs.A & 0x88) >> 3) | ((s & 0x88) >> 2) | byte(res&0x88>>1)
	cpu.regs.A = byte(res)
	cpu.regs.F.S = cpu.regs.A&0x0080 != 0
	cpu.regs.F.Z = cpu.regs.A == 0
	cpu.regs.F.H = halfcarrySubTable[lookup&0x07]
	cpu.regs.F.P = overflowSubTable[lookup>>4]
	cpu.regs.F.N = true
	cpu.regs.F.C = (res & 0x100) == 0x100
}

func (cpu *z80) sbcHL(ss uint16) {
	hl := cpu.regs.HL.Get()
	res := uint32(hl) - uint32(ss)
	if cpu.regs.F.C {
		res--
	}
	cpu.regs.HL.Set(uint16(res))

	lookup := byte(((hl & 0x8800) >> 11) | ((ss & 0x8800) >> 10) | ((uint16(res) & 0x8800) >> 9))
	cpu.regs.F.N = true
	cpu.regs.F.S = cpu.regs.H&0x80 != 0 // negative
	cpu.regs.F.Z = res == 0
	cpu.regs.F.C = (res & 0x10000) != 0
	cpu.regs.F.P = overflowSubTable[lookup>>4]
	cpu.regs.F.H = halfcarrySubTable[lookup&0x07]
}
