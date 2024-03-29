package lr35902

func retCC(cpu *lr35902) {
	cpu.scheduler.append(&exec{f: func(cpu *lr35902) {
		ccIdx := cpu.fetched.opCode >> 3 & 0b111
		branch := cpu.checkCondition(ccIdx)
		if branch {
			cpu.popFromStack(func(cpu *lr35902, data uint16) {
				cpu.scheduler.append(&exec{f: func(cpu *lr35902) {
					cpu.regs.PC = data
				}})
			})
		}
	}})
}

func ret(cpu *lr35902) {
	cpu.popFromStack(func(cpu *lr35902, data uint16) {
		cpu.scheduler.append(&exec{f: func(cpu *lr35902) {
			cpu.regs.PC = data
		}})
	})
}

func reti(cpu *lr35902) {
	cpu.scheduler.append(&exec{f: func(cpu *lr35902) {
		cpu.popFromStack(func(cpu *lr35902, data uint16) {
			cpu.regs.PC = data
			cpu.regs.IME = true
		})
	}})
}

func rstP(cpu *lr35902) {
	cpu.scheduler.append(&exec{f: func(cpu *lr35902) {
		cpu.pushToStack(cpu.regs.PC, rstP_m1)
	}})
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
		jpNN(cpu)
	}
}

func jpNN(cpu *lr35902) {
	cpu.scheduler.append(&exec{f: func(cpu *lr35902) {
		cpu.regs.PC = cpu.fetched.nn
	}})
}

func call(cpu *lr35902) {
	ccIdx := cpu.fetched.opCode >> 3 & 0b11
	branch := cpu.checkCondition(ccIdx)
	if branch {
		cpu.scheduler.append(&exec{f: call_m2})
	}
}

func call_m2(cpu *lr35902) { cpu.pushToStack(cpu.regs.PC, call_m3) }

func call_m3(cpu *lr35902) { cpu.regs.PC = cpu.fetched.nn }

func (cpu *lr35902) checkCondition(ccIdx byte) bool {
	if cpu.fetched.opCode == 0xCD || cpu.fetched.opCode == 0xC3 {
		return true
	}

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
	}
	return res
}

var pushF func(cpu *lr35902)

func (cpu *lr35902) pushToStack(data uint16, f func(cpu *lr35902)) {
	pushF = f
	push1 := &mw{to: cpu.regs.SP.Get() - 1, d: uint8(data >> 8), f: nil}
	push2 := &mw{to: cpu.regs.SP.Get() - 2, d: uint8(data), f: push_m1}
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
	pop1 := &mr{from: cpu.regs.SP.Get(), f: pop_m1}
	pop2 := &mr{from: cpu.regs.SP.Get() + 1, f: pop_m2}
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
	cpu.scheduler.append(&exec{f: func(cpu *lr35902) { cpu.pushToStack(data, nil) }})
}

func ldDDnn(cpu *lr35902) {
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
	cpu.scheduler.append(&mw{to: pos, d: cpu.regs.A, f: nil})
}

func ldDEa(cpu *lr35902) {
	pos := cpu.regs.DE.Get()
	cpu.scheduler.append(&mw{to: pos, d: cpu.regs.A, f: nil})
}

func ldNNa(cpu *lr35902) {
	mm := cpu.fetched.nn
	mw1 := &mw{to: mm, d: cpu.regs.A, f: nil}
	cpu.scheduler.append(mw1)
}

func ldiHLa(cpu *lr35902) {
	cpu.scheduler.append(&mw{to: cpu.regs.HL.Get(), d: cpu.regs.A, f: ldiHLa_m2})
}

func ldiHLa_m2(cpu *lr35902) {
	cpu.regs.HL.Set(cpu.regs.HL.Get() + 1)
}

func lddAhl(cpu *lr35902) {
	cpu.scheduler.append(&mr{from: cpu.regs.HL.Get(), f: lddAhl_m2})
}

func lddAhl_m2(cpu *lr35902, data byte) {
	cpu.regs.A = data
	cpu.regs.HL.Set(cpu.regs.HL.Get() - 1)
}

func lddHLa(cpu *lr35902) {
	cpu.scheduler.append(&mw{to: cpu.regs.HL.Get(), d: cpu.regs.A, f: lddHLa_m2})
}

func lddHLa_m2(cpu *lr35902) {
	cpu.regs.HL.Set(cpu.regs.HL.Get() - 1)
}

func ldAnn(cpu *lr35902) {
	mm := cpu.fetched.nn
	mr1 := &mr{from: mm, f: ldAnn_n1}
	cpu.scheduler.append(mr1)
}

func ldAnn_n1(cpu *lr35902, data uint8) { cpu.regs.A = data }

func incSS(cpu *lr35902) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	v := reg.Get()
	v++
	cpu.scheduler.append(&exec{f: func(cpu *lr35902) {
		reg.Set(v)
	}})
}

func decSS(cpu *lr35902) {
	rIdx := cpu.fetched.opCode >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	v := reg.Get()
	v--
	*reg.L = byte(v)
	cpu.scheduler.append(&exec{f: func(cpu *lr35902) {
		*reg.H = byte(v >> 8)
	}})
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
	mr := &mr{from: cpu.regs.HL.Get(), f: func(cpu *lr35902, data uint8) {
		r := data
		r++
		mw := &mw{to: cpu.regs.HL.Get(), d: r, f: nil}
		cpu.regs.F.Z = r == 0
		cpu.regs.F.H = r&0x0f == 0
		cpu.regs.F.N = false

		cpu.scheduler.append(mw)
	},
	}
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
	mr := &mr{from: cpu.regs.HL.Get(), f: addAhl_m1}
	cpu.scheduler.append(mr)
}

func addAhl_m1(cpu *lr35902, data uint8) { cpu.addA(data) }

func subAhl(cpu *lr35902) {
	mr := &mr{from: cpu.regs.HL.Get(), f: subAhl_m1}
	cpu.scheduler.append(mr)
}

func subAhl_m1(cpu *lr35902, data uint8) { cpu.subA(data) }

func sbcAhl(cpu *lr35902) {
	mr := &mr{from: cpu.regs.HL.Get(), f: sbcAhl_m1}
	cpu.scheduler.append(mr)
}

func sbcAhl_m1(cpu *lr35902, data uint8) { cpu.sbcA(data) }

func adcAhl(cpu *lr35902) {
	mr := &mr{from: cpu.regs.HL.Get(), f: adcAhl_m1}
	cpu.scheduler.append(mr)
}

func adcAhl_m1(cpu *lr35902, data uint8) { cpu.adcA(data) }

func andAhl(cpu *lr35902) {
	mr := &mr{from: cpu.regs.HL.Get(), f: andAhl_m1}
	cpu.scheduler.append(mr)
}

func andAhl_m1(cpu *lr35902, data uint8) { cpu.and(data) }

func orAhl(cpu *lr35902) {
	mr := &mr{from: cpu.regs.HL.Get(), f: orAhl_m1}
	cpu.scheduler.append(mr)
}

func orAhl_m1(cpu *lr35902, data uint8) { cpu.or(data) }

func xorAhl(cpu *lr35902) {
	mr := &mr{from: cpu.regs.HL.Get(), f: xorAhl_m1}
	cpu.scheduler.append(mr)
}

func xorAhl_m1(cpu *lr35902, data uint8) { cpu.xor(data) }

func cpHl(cpu *lr35902) {
	mr := &mr{from: cpu.regs.HL.Get(), f: cpHl_m1}
	cpu.scheduler.append(mr)
}

func cpHl_m1(cpu *lr35902, data uint8) { cpu.cp(data) }

func decHL(cpu *lr35902) {
	mr := &mr{from: cpu.regs.HL.Get(), f: func(cpu *lr35902, data uint8) {
		r := data
		cpu.regs.F.H = r&0x0f == 0
		r--
		cpu.regs.F.Z = r == 0
		cpu.regs.F.N = true

		mw := &mw{to: cpu.regs.HL.Get(), d: r, f: nil}
		cpu.scheduler.append(mw)
	},
	}
	cpu.scheduler.append(mr)
}
