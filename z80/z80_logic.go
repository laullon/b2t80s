package z80

import (
	"fmt"
)

func retCC(cpu *z80, mem []uint8) {
	ccIdx := mem[0] >> 3 & 0b111
	branch := cpu.checkCondition(ccIdx)
	if branch {
		cpu.popFromStack(func(cpu *z80, data uint16) {
			cpu.regs.PC = data
		})
	}
}

func ret(cpu *z80, mem []uint8) {
	cpu.popFromStack(func(cpu *z80, data uint16) {
		cpu.regs.PC = data
	})
}

func rstP(cpu *z80, mem []uint8) {
	newPCs := []uint16{0x00, 0x08, 0x10, 0x18, 0x20, 0x28, 0x30, 0x38}
	pIdx := mem[0] >> 3 & 0b111
	println("pc", cpu.regs.PC)
	cpu.pushToStack(cpu.regs.PC, func(cpu *z80) { cpu.regs.PC = newPCs[pIdx]; println("pc", cpu.regs.PC) })
}

func jpCC(cpu *z80, mem []uint8) {
	ccIdx := mem[0] >> 3 & 0b111
	branch := cpu.checkCondition(ccIdx)
	if branch {
		cpu.regs.PC = toWord(mem[1], mem[2])
	}
}

func callCC(cpu *z80, mem []uint8) {
	ccIdx := mem[0] >> 3 & 0b111
	branch := cpu.checkCondition(ccIdx)
	if branch {
		cpu.pushToStack(cpu.regs.PC, func(cpu *z80) {
			cpu.regs.PC = toWord(mem[1], mem[2])
		})
	}
}

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
	cpu.regs.SP.Set(cpu.regs.SP.Get() - 2)

	push1 := &mw{addr: cpu.regs.SP.Get(), data: uint8(data)}
	push2 := &mw{addr: cpu.regs.SP.Get() + 1, data: uint8(data >> 8), f: func(z *z80) {
		if f != nil {
			f(cpu)
		}
	}}
	cpu.scheduler = append(cpu.scheduler, push1, push2)
}

func (cpu *z80) popFromStack(f func(cpu *z80, data uint16)) {
	var data uint16
	pop1 := &mr{from: cpu.regs.SP.Get(), f: func(z *z80, d []uint8) { data = uint16(d[0]) }}
	pop2 := &mr{from: cpu.regs.SP.Get() + 1, f: func(z *z80, d []uint8) {
		data |= (uint16(d[0]) << 8)
		cpu.regs.SP.Set(cpu.regs.SP.Get() + 2)
		f(cpu, data)
	}}
	cpu.scheduler = append(cpu.scheduler, pop1, pop2)
}

func popSS(cpu *z80, mem []uint8) {
	t := mem[0] >> 4 & 0b11
	cpu.popFromStack(func(cpu *z80, data uint16) {
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
	})
}

func pushSS(cpu *z80, mem []uint8) {
	t := mem[0] >> 4 & 0b11
	var data uint16
	switch t {
	case 0b00:
		data = cpu.regs.BC.Get()
	case 0b01:
		data = cpu.regs.DE.Get()
	case 0b10:
		data = cpu.regs.HL.Get()
	case 0b11:
		data = uint16(cpu.regs.A) >> 8
		data |= uint16(cpu.regs.F.GetByte())
	}
	cpu.pushToStack(data, nil)
}

func ldDDmm(cpu *z80, mem []uint8) {
	t := mem[0] >> 4 & 0b11
	switch t {
	case 0b00:
		cpu.regs.B = mem[2]
		cpu.regs.C = mem[1]
	case 0b01:
		cpu.regs.D = mem[2]
		cpu.regs.E = mem[1]
	case 0b10:
		cpu.regs.H = mem[2]
		cpu.regs.L = mem[1]
	case 0b11:
		cpu.regs.S = mem[2]
		cpu.regs.P = mem[1]
	}
}

func ldBCa(cpu *z80, mem []uint8) {
	pos := cpu.regs.BC.Get()
	cpu.scheduler = append(cpu.scheduler, &mw{addr: pos, data: cpu.regs.A})
}

func ldDEa(cpu *z80, mem []uint8) {
	pos := cpu.regs.DE.Get()
	cpu.scheduler = append(cpu.scheduler, &mw{addr: pos, data: cpu.regs.A})
}

func ldNNhl(cpu *z80, mem []uint8) {
	mm := toWord(mem[1], mem[2])
	mw1 := &mw{addr: mm, data: cpu.regs.L}
	mw2 := &mw{addr: mm + 1, data: cpu.regs.H}
	cpu.scheduler = append(cpu.scheduler, mw1, mw2)
}

func ldNNa(cpu *z80, mem []uint8) {
	mm := toWord(mem[1], mem[2])
	mw1 := &mw{addr: mm, data: cpu.regs.A}
	cpu.scheduler = append(cpu.scheduler, mw1)
}

func ldHLnn(cpu *z80, mem []uint8) {
	mm := toWord(mem[1], mem[2])
	mr1 := &mr{from: mm, f: func(z *z80, d []uint8) { cpu.regs.L = d[0] }}
	mr2 := &mr{from: mm + 1, f: func(z *z80, d []uint8) { cpu.regs.H = d[0] }}
	cpu.scheduler = append(cpu.scheduler, mr1, mr2)
}

func ldAnn(cpu *z80, mem []uint8) {
	mm := toWord(mem[1], mem[2])
	mr1 := &mr{from: mm, f: func(z *z80, d []uint8) { cpu.regs.A = d[0] }}
	cpu.scheduler = append(cpu.scheduler, mr1)
}

func ldHLn(cpu *z80, mem []uint8) {
	mw1 := &mw{addr: cpu.regs.HL.Get(), data: mem[1]}
	cpu.scheduler = append(cpu.scheduler, mw1)
}

func incSS(cpu *z80, mem []uint8) {
	rIdx := mem[0] >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	v := reg.Get()
	v++
	reg.Set(v)
}

func decSS(cpu *z80, mem []uint8) {
	rIdx := mem[0] >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)
	v := reg.Get()
	v--
	reg.Set(v)
}

func incR(cpu *z80, mem []uint8) {
	rIdx := mem[0] >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	*r++
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.H = *r&0x0f == 0
	cpu.regs.F.P = *r == 0x80
	cpu.regs.F.N = false
	// panic(fmt.Sprintf("%08b", *r&0x0f))
}

func incHL(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) {
			r := d[0]
			r++
			mw := &mw{addr: cpu.regs.HL.Get(), data: r}
			cpu.regs.F.S = r&0x80 != 0
			cpu.regs.F.Z = r == 0
			cpu.regs.F.H = r&0x0f == 0
			cpu.regs.F.P = r == 0x80
			cpu.regs.F.N = false

			cpu.scheduler = append(cpu.scheduler, mw)
		},
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func decR(cpu *z80, mem []uint8) {
	rIdx := mem[0] >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	cpu.regs.F.H = *r&0x0f == 0
	*r--
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = *r == 0x7f
	cpu.regs.F.N = true
}

func addAr(cpu *z80, mem []uint8) {
	rIdx := mem[0] & 0b111
	r := cpu.getRptr(rIdx)
	cpu.addA(*r)
}

func adcAr(cpu *z80, mem []uint8) {
	rIdx := mem[0] & 0b111
	r := cpu.getRptr(rIdx)
	cpu.adcA(*r)
}

func subAr(cpu *z80, mem []uint8) {
	rIdx := mem[0] & 0b111
	r := cpu.getRptr(rIdx)
	cpu.subA(*r)
}

func sbcAr(cpu *z80, mem []uint8) {
	rIdx := mem[0] & 0b111
	r := cpu.getRptr(rIdx)
	cpu.sbcA(*r)
}

func andAr(cpu *z80, mem []uint8) {
	rIdx := mem[0] & 0b111
	r := cpu.getRptr(rIdx)
	cpu.and(*r)
}

func orAr(cpu *z80, mem []uint8) {
	rIdx := mem[0] & 0b111
	r := cpu.getRptr(rIdx)
	cpu.or(*r)
}

func xorAr(cpu *z80, mem []uint8) {
	rIdx := mem[0] & 0b111
	r := cpu.getRptr(rIdx)
	cpu.xor(*r)
}

func cpR(cpu *z80, mem []uint8) {
	rIdx := mem[0] & 0b111
	r := cpu.getRptr(rIdx)
	cpu.cp(*r)
}

func addAhl(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) { cpu.addA(d[0]) },
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func subAhl(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) { cpu.subA(d[0]) },
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func sbcAhl(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) { cpu.sbcA(d[0]) },
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func adcAhl(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) { cpu.adcA(d[0]) },
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func andAhl(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) { cpu.and(d[0]) },
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func orAhl(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) { cpu.or(d[0]) },
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func xorAhl(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) { cpu.xor(d[0]) },
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func cpHl(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) { cpu.cp(d[0]) },
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func decHL(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) {
			r := d[0]
			cpu.regs.F.H = r&0x0f == 0
			r--
			cpu.regs.F.S = r&0x80 != 0
			cpu.regs.F.Z = r == 0
			cpu.regs.F.P = r == 0x7f
			cpu.regs.F.N = true

			mw := &mw{addr: cpu.regs.HL.Get(), data: r}
			cpu.scheduler = append(cpu.scheduler, mw)
		},
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func ldRn(cpu *z80, mem []uint8) {
	rIdx := mem[0] >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	*r = mem[1]
}

func ldRhl(cpu *z80, mem []uint8) {
	rIdx := mem[0] >> 3 & 0b111
	r := cpu.getRptr(rIdx)
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) {
			*r = d[0]
		},
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func ldHLr(cpu *z80, mem []uint8) {
	rIdx := mem[0] & 0b111
	r := cpu.getRptr(rIdx)
	mr := &mw{addr: cpu.regs.HL.Get(), data: *r}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func ldRr(cpu *z80, mem []uint8) {
	r1Idx := mem[0] >> 3 & 0b111
	r2Idx := mem[0] & 0b111
	r1 := cpu.getRptr(r1Idx)
	r2 := cpu.getRptr(r2Idx)
	*r1 = *r2
}

func rlca(cpu *z80, mem []uint8) {
	cpu.regs.A = cpu.regs.A<<1 | cpu.regs.A>>7
	cpu.regs.F.C = cpu.regs.A&0x01 != 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func rla(cpu *z80, mem []uint8) {
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

func cbR(cpu *z80, mem []uint8) {
	rIdx := mem[0] & 0b111
	r := cpu.getRptr(rIdx)
	fIdx := mem[0] >> 3
	cbFuncs[fIdx](cpu, r)
}

func cbHL(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) {
			b := d[0]
			fIdx := mem[0] >> 3
			cbFuncs[fIdx](cpu, &b)
			mw := &mw{addr: cpu.regs.HL.Get(), data: b}
			cpu.scheduler = append(cpu.scheduler, mw)
		},
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func bit(cpu *z80, mem []uint8) {
	rIdx := mem[0] & 0b111
	r := cpu.getRptr(rIdx)
	b := (mem[0] >> 3) & 0b111
	cpu.bit(b, *r)
}

func bitHL(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) {
			v := d[0]
			b := (mem[0] >> 3) & 0b111
			cpu.bit(b, v)
		},
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func res(cpu *z80, mem []uint8) {
	rIdx := mem[0] & 0b111
	r := cpu.getRptr(rIdx)
	b := (mem[0] >> 3) & 0b111
	cpu.res(b, r)
}

func resHL(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) {
			v := d[0]
			b := (mem[0] >> 3) & 0b111
			cpu.res(b, &v)
			mw := &mw{addr: cpu.regs.HL.Get(), data: v}
			cpu.scheduler = append(cpu.scheduler, mw)
		},
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func set(cpu *z80, mem []uint8) {
	rIdx := mem[0] & 0b111
	r := cpu.getRptr(rIdx)
	b := (mem[0] >> 3) & 0b111
	cpu.set(b, r)
}

func setHL(cpu *z80, mem []uint8) {
	mr := &mr{from: cpu.regs.HL.Get(),
		f: func(z *z80, d []uint8) {
			v := d[0]
			b := (mem[0] >> 3) & 0b111
			cpu.set(b, &v)
			mw := &mw{addr: cpu.regs.HL.Get(), data: v}
			cpu.scheduler = append(cpu.scheduler, mw)
		},
	}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func rrca(cpu *z80, mem []uint8) {
	cpu.regs.F.C = cpu.regs.A&0x01 != 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
	cpu.regs.A = (cpu.regs.A >> 1) | (cpu.regs.A << 7)
}

func rra(cpu *z80, mem []uint8) {
	c := cpu.regs.F.C
	cpu.regs.F.C = cpu.regs.A&1 != 0
	cpu.regs.A = (cpu.regs.A >> 1)
	if c {
		cpu.regs.A |= 0b10000000
	}
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func exafaf(cpu *z80, mem []uint8) {
	cpu.regs.A, cpu.regs.Aalt = cpu.regs.Aalt, cpu.regs.A
	cpu.regs.F, cpu.regs.Falt = cpu.regs.Falt, cpu.regs.F
}

func halt(cpu *z80, mem []uint8) {
	if cpu.haltDone {
		cpu.haltDone = false
	} else {
		cpu.halt = true
	}

}

func addHLss(cpu *z80, mem []uint8) {
	rIdx := mem[0] >> 4 & 0b11
	reg := cpu.getRRptr(rIdx)

	hl := cpu.regs.HL.Get()
	var result = uint32(hl) + uint32(reg.Get())
	var lookup = byte(((hl & 0x0800) >> 11) | ((reg.Get() & 0x0800) >> 10) | ((uint16(result) & 0x0800) >> 9))
	cpu.regs.HL.Set(uint16(result))

	cpu.regs.F.N = false
	cpu.regs.F.H = halfcarryAddTable[lookup]
	cpu.regs.F.C = (result & 0x10000) != 0
}

func ldAbc(cpu *z80, mem []uint8) {
	from := cpu.regs.BC.Get()
	mr := &mr{from: from, f: func(z *z80, data []uint8) { cpu.regs.A = data[0] }}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func ldAde(cpu *z80, mem []uint8) {
	from := cpu.regs.DE.Get()
	mr := &mr{from: from, f: func(z *z80, data []uint8) { cpu.regs.A = data[0] }}
	cpu.scheduler = append(cpu.scheduler, mr)
}

func djnz(cpu *z80, mem []uint8) {
	cpu.regs.B--
	if cpu.regs.B != 0 {
		cpu.scheduler = append(cpu.scheduler, &exec{l: 5, f: jr})
	}
}

func jrnz(cpu *z80, mem []uint8) {
	if !cpu.regs.F.Z {
		cpu.scheduler = append(cpu.scheduler, &exec{l: 5, f: jr})
	}
}

func jrnc(cpu *z80, mem []uint8) {
	if !cpu.regs.F.C {
		cpu.scheduler = append(cpu.scheduler, &exec{l: 5, f: jr})
	}
}

func jrc(cpu *z80, mem []uint8) {
	if cpu.regs.F.C {
		cpu.scheduler = append(cpu.scheduler, &exec{l: 5, f: jr})
	}
}

func jrz(cpu *z80, mem []uint8) {
	if cpu.regs.F.Z {
		cpu.scheduler = append(cpu.scheduler, &exec{l: 5, f: jr})
	}
}

func jr(cpu *z80, mem []uint8) {
	jump := int8(mem[1])
	cpu.regs.PC += uint16(jump)
}

func scf(cpu *z80, mem []uint8) {
	cpu.regs.F.H = false
	cpu.regs.F.N = false
	cpu.regs.F.C = true
}

func ccf(cpu *z80, mem []uint8) {
	cpu.regs.F.H = cpu.regs.F.C
	cpu.regs.F.N = false
	cpu.regs.F.C = !cpu.regs.F.C
}

func daa(cpu *z80, mem []uint8) {
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

func cpl(cpu *z80, mem []uint8) {
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

func (cpu *z80) getRRptr(rIdx byte) *RegPair {
	var reg *RegPair
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

// -------
// TODO review bellow

func (cpu *z80) writePort(port uint16, data byte) {
	// fmt.Printf("[writePort]-> port:0x%04X data:%v pc:0x%04X \n", port, data, cpu.regs.PC)
	ok := false
	for portMask, portManager := range cpu.ports {
		// fmt.Printf("[writePort] (0x%04X) port:0x%04X (0x%04X)(0x%04X) data:%v\n", cpu.regs.PC, port, port&portMask.Mask, portMask.Value, data)
		if port&portMask.Mask == portMask.Value {
			// println(reflect.TypeOf(portManager).String())
			portManager.WritePort(port, data)
			ok = true
		}
	}
	if !ok {
		fmt.Printf("[writePort]-(no PM)-> port:0x%04X data:%v pc:0x%04X\n", port, data, cpu.regs.PC)
		// panic("--")
	}
}

func (cpu *z80) readPort(port uint16) byte {
	// fmt.Printf(fmt.Sprintf("[readPort]-> port:0x%04X pc:0x%04X \n", port, cpu.regs.PC))
	for portMask, portManager := range cpu.ports {
		if port&portMask.Mask == portMask.Value {
			// fmt.Printf("[readPort] (0x%04X) port:0x%04X (0x%04X)(0x%04X) \n", cpu.regs.PC, port, port&portMask.Mask, portMask.Value)
			// println(reflect.TypeOf(portManager).Elem().Name())
			data, skip := portManager.ReadPort(port)
			if !skip {
				cpu.regs.F.S = data&0x0080 != 0
				cpu.regs.F.Z = data == 0
				cpu.regs.F.H = false
				cpu.regs.F.P = parityTable[data]
				cpu.regs.F.N = false
				return data
			}
		}
	}
	// panic(fmt.Sprintf("[readPort]-(no PM)-> port:0x%04X pc:0x%04X", port, cpu.regs.PC))
	// fmt.Printf("[readPort]-(no PM)-> port:0x%04X pc:0x%04X \n", port, cpu.regs.PC)
	return 0xff
}

func (cpu *z80) getIXn(n byte) uint16 {
	i := int16(int8(n))
	ix := cpu.regs.IX.Get()
	ix = uint16(int16(ix) + i)
	return ix
}

func (cpu *z80) getIYn(n byte) uint16 {
	i := int16(int8(n))
	iy := cpu.regs.IY.Get()
	iy = uint16(int16(iy) + i)
	return iy
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

func (cpu *z80) cpd() byte {
	bc := cpu.regs.BC.Get()
	hl := cpu.regs.HL.Get()

	val := cpu.memory.GetByte(hl)
	result := cpu.regs.A - val
	lookup := (cpu.regs.A&0x08)>>3 | (val&0x08)>>2 | (result&0x08)>>1

	bc--
	hl--

	cpu.regs.BC.Set(bc)
	cpu.regs.HL.Set(hl)

	cpu.regs.F.S = result&0x80 != 0
	cpu.regs.F.Z = result == 0
	cpu.regs.F.H = halfcarrySubTable[lookup]
	cpu.regs.F.P = bc != 0
	cpu.regs.F.N = true

	return result
}

func (cpu *z80) cpi() byte {
	bc := cpu.regs.BC.Get()
	hl := cpu.regs.HL.Get()

	val := cpu.memory.GetByte(hl)
	result := cpu.regs.A - val
	lookup := (cpu.regs.A&0x08)>>3 | (val&0x08)>>2 | (result&0x08)>>1

	bc--
	hl++

	cpu.regs.BC.Set(bc)
	cpu.regs.HL.Set(hl)

	cpu.regs.F.S = result&0x80 != 0
	cpu.regs.F.Z = result == 0
	cpu.regs.F.H = halfcarrySubTable[lookup]
	cpu.regs.F.P = bc != 0
	cpu.regs.F.N = true

	return result
}

// ------

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

func (cpu *z80) jr(j byte) {
	jump := int8(j)
	cpu.regs.PC += uint16(jump)
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

func (cpu *z80) ldd() {
	bc := cpu.regs.BC.Get()
	de := cpu.regs.DE.Get()
	hl := cpu.regs.HL.Get()

	v := cpu.memory.GetByte(hl)
	cpu.memory.PutByte(de, v)

	bc--
	de--
	hl--

	cpu.regs.BC.Set(bc)
	cpu.regs.DE.Set(de)
	cpu.regs.HL.Set(hl)

	cpu.regs.F.P = bc != 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func (cpu *z80) ldi() {
	bc := cpu.regs.BC.Get()
	de := cpu.regs.DE.Get()
	hl := cpu.regs.HL.Get()

	v := cpu.memory.GetByte(hl)
	cpu.memory.PutByte(de, v)

	bc--
	de++
	hl++

	cpu.regs.BC.Set(bc)
	cpu.regs.DE.Set(de)
	cpu.regs.HL.Set(hl)

	cpu.regs.F.P = bc != 0
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func (cpu *z80) addIX(ss uint16) {
	ix := cpu.regs.IX.Get()
	var result = uint32(ix) + uint32(ss)
	var lookup = byte(((ix & 0x0800) >> 11) | ((ss & 0x0800) >> 10) | ((uint16(result) & 0x0800) >> 9))
	cpu.regs.IX.Set(uint16(result))

	cpu.regs.F.N = false
	cpu.regs.F.H = halfcarryAddTable[lookup]
	cpu.regs.F.C = (result & 0x10000) != 0
}

func (cpu *z80) addIY(ss uint16) {
	iy := cpu.regs.IY.Get()
	var result = uint32(iy) + uint32(ss)
	var lookup = byte(((iy & 0x0800) >> 11) | ((ss & 0x0800) >> 10) | ((uint16(result) & 0x0800) >> 9))
	cpu.regs.IY.Set(uint16(result))

	cpu.regs.F.N = false
	cpu.regs.F.H = halfcarryAddTable[lookup]
	cpu.regs.F.C = (result & 0x10000) != 0
}
