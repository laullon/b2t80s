package z80

import (
	"fmt"
)

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
		panic("--")
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

func (cpu *z80) adc(s byte) {
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

func (cpu *z80) cp(r byte) byte {
	a := int16(cpu.regs.A)
	result := a - int16(r)
	lookup := ((cpu.regs.A & 0x88) >> 3) | (((r) & 0x88) >> 2) | ((byte(result) & 0x88) >> 1)

	cpu.regs.F.S = result&0x80 != 0
	cpu.regs.F.Z = result == 0
	cpu.regs.F.H = halfcarrySubTable[lookup&0x07]
	cpu.regs.F.P = overflowSubTable[lookup>>4]
	cpu.regs.F.N = true
	cpu.regs.F.C = ((result) & 0x100) == 0x100
	return byte(result)
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

func (cpu *z80) incR(r *byte) {
	*r++
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.H = *r&0x0f == 0
	cpu.regs.F.P = *r == 0x80
	cpu.regs.F.N = false
	// panic(fmt.Sprintf("%08b", *r&0x0f))
}

func (cpu *z80) decR(r *byte) {
	cpu.regs.F.H = *r&0x0f == 0
	*r--
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = *r == 0x7f
	cpu.regs.F.N = true
}

// ------

func (cpu *z80) rlc(r *byte) {
	*r = (*r << 1) | (*r >> 7)
	cpu.regs.F.C = *r&1 != 0
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func (cpu *z80) rrc(r *byte) {
	cpu.regs.F.C = *r&1 != 0
	*r = (*r << 7) | (*r >> 1)
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func (cpu *z80) sll(r *byte) {
	cpu.regs.F.C = *r&0x80 != 0
	*r = byte((*r << 1) | 0x01)
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func (cpu *z80) srl(r *byte) {
	cpu.regs.F.C = *r&1 != 0
	*r = byte(*r >> 1)
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.N = false
	cpu.regs.F.H = false
}

func (cpu *z80) sla(r *byte) {
	cpu.regs.F.C = *r&0x80 != 0
	*r = (*r << 1)
	cpu.regs.F.S = *r&0x80 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func (cpu *z80) sra(r *byte) {
	cpu.regs.F.C = *r&0x1 != 0
	b7 := *r & 0b10000000
	*r = (*r >> 1) | b7
	cpu.regs.F.S = *r&0x0080 != 0
	cpu.regs.F.Z = *r == 0
	cpu.regs.F.P = parityTable[*r]
	cpu.regs.F.H = false
	cpu.regs.F.N = false
}

func (cpu *z80) rr(r *byte) {
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

func (cpu *z80) rl(r *byte) {
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
	jump := int8(j) + 2
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

func (cpu *z80) sbc(s byte) {
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

//TODO Join these 3
func (cpu *z80) addHL(ss uint16) {
	hl := cpu.regs.HL.Get()
	var result = uint32(hl) + uint32(ss)
	var lookup = byte(((hl & 0x0800) >> 11) | ((ss & 0x0800) >> 10) | ((uint16(result) & 0x0800) >> 9))
	cpu.regs.HL.Set(uint16(result))

	cpu.regs.F.N = false
	cpu.regs.F.H = halfcarryAddTable[lookup]
	cpu.regs.F.C = (result & 0x10000) != 0
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
