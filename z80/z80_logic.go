package z80

import (
	"fmt"
)

func (cpu *z80) writePort(port uint16, data byte) {
	// log.Println(fmt.Sprintf("[writePort]-> port:0x%04X data:%v pc:0x%04X", port, data, cpu.pc))
	ok := false
	for portMask, portManager := range cpu.ports {
		// log.Printf("[writePort] (0x%04X) port:0x%04X (0x%04X)(0x%04X) data:%v", cpu.pc, port, port&portMask.Mask, portMask.Value, data)
		if port&portMask.Mask == portMask.Value {
			// println(reflect.TypeOf(portManager).String())
			portManager.WritePort(port, data)
			ok = true
		}
	}
	if !ok {
		// log.Println(fmt.Sprintf("[writePort]-(no PM)-> port:0x%04X data:%v pc:0x%04X", port, data, cpu.pc))
		// panic("--")
	}
}

func (cpu *z80) readPort(port uint16) byte {
	// log.Println(fmt.Sprintf("[readPort]-> port:0x%04X pc:0x%04X", port, cpu.pc))
	for portMask, portManager := range cpu.ports {
		if port&portMask.Mask == portMask.Value {
			// log.Printf("[readPort] (0x%04X) port:0x%04X (0x%04X)(0x%04X)", cpu.pc, port, port&portMask.Mask, portMask.Value)
			// println(reflect.TypeOf(portManager).Elem().Name())
			data, skip := portManager.ReadPort(port)
			if !skip {
				cpu.f.S = data&0x0080 != 0
				cpu.f.Z = data == 0
				cpu.f.H = false
				cpu.f.P = parityTable[data]
				cpu.f.N = false
				return data
			}
		}
	}
	// panic(fmt.Sprintf("[readPort]-(no PM)-> port:0x%04X pc:0x%04X", port, cpu.pc))
	fmt.Printf("[readPort]-(no PM)-> port:0x%04X pc:0x%04X \n", port, cpu.pc)
	return 0xff
}

func (cpu *z80) getIXn(n byte) uint16 {
	i := int16(int8(n))
	ix := getRR(cpu.ixh, cpu.ixl)
	ix = uint16(int16(ix) + i)
	return ix
}

func (cpu *z80) getIYn(n byte) uint16 {
	i := int16(int8(n))
	iy := getRR(cpu.iyh, cpu.iyl)
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
	cpu.f.S = v&0x0080 != 0
	cpu.f.Z = v == 0
	cpu.f.H = true
	cpu.f.P = parityTable[v]
	cpu.f.N = false
}

func (cpu *z80) adc(s byte) {
	res := int16(cpu.a) + int16(s)
	if cpu.f.C {
		res++
	}
	lookup := ((cpu.a & 0x88) >> 3) | ((s & 0x88) >> 2) | ((byte(res) & 0x88) >> 1)
	cpu.a = byte(res)
	cpu.f.S = cpu.a&0x80 != 0
	cpu.f.Z = cpu.a == 0
	cpu.f.H = halfcarryAddTable[lookup&0x07]
	cpu.f.P = overflowAddTable[lookup>>4]
	cpu.f.N = false
	cpu.f.C = (res & 0x100) == 0x100
}

func (cpu *z80) adcHL(ss uint16) {
	hl := getRR(cpu.h, cpu.l)
	res := int32(hl) + int32(ss)
	if cpu.f.C {
		res++
	}
	lookup := byte(((hl & 0x8800) >> 11) | ((ss & 0x8800) >> 10) | ((uint16(res) & 0x8800) >> 9))
	hl = uint16(res)
	cpu.h, cpu.l = setRR(hl)
	cpu.f.S = cpu.h&0x80 != 0
	cpu.f.Z = hl == 0
	cpu.f.H = halfcarryAddTable[lookup&0x07]
	cpu.f.P = overflowAddTable[lookup>>4]
	cpu.f.N = false
	cpu.f.C = (res & 0x10000) != 0
}

func (cpu *z80) cp(r byte) byte {
	a := int16(cpu.a)
	result := a - int16(r)
	lookup := ((cpu.a & 0x88) >> 3) | (((r) & 0x88) >> 2) | ((byte(result) & 0x88) >> 1)

	cpu.f.S = result&0x80 != 0
	cpu.f.Z = result == 0
	cpu.f.H = halfcarrySubTable[lookup&0x07]
	cpu.f.P = overflowSubTable[lookup>>4]
	cpu.f.N = true
	cpu.f.C = ((result) & 0x100) == 0x100
	return byte(result)
}

func (cpu *z80) cpd() byte {
	bc := getRR(cpu.b, cpu.c)
	hl := getRR(cpu.h, cpu.l)

	val := cpu.memory.GetByte(hl)
	result := cpu.a - val
	lookup := (cpu.a&0x08)>>3 | (val&0x08)>>2 | (result&0x08)>>1

	bc--
	hl--

	cpu.b, cpu.c = setRR(bc)
	cpu.h, cpu.l = setRR(hl)

	cpu.f.S = result&0x80 != 0
	cpu.f.Z = result == 0
	cpu.f.H = halfcarrySubTable[lookup]
	cpu.f.P = bc != 0
	cpu.f.N = true

	return result
}

func (cpu *z80) cpi() byte {
	bc := getRR(cpu.b, cpu.c)
	hl := getRR(cpu.h, cpu.l)

	val := cpu.memory.GetByte(hl)
	result := cpu.a - val
	lookup := (cpu.a&0x08)>>3 | (val&0x08)>>2 | (result&0x08)>>1

	bc--
	hl++

	cpu.b, cpu.c = setRR(bc)
	cpu.h, cpu.l = setRR(hl)

	cpu.f.S = result&0x80 != 0
	cpu.f.Z = result == 0
	cpu.f.H = halfcarrySubTable[lookup]
	cpu.f.P = bc != 0
	cpu.f.N = true

	return result
}

func (cpu *z80) incR(r *byte) {
	*r++
	cpu.f.S = *r&0x80 != 0
	cpu.f.Z = *r == 0
	cpu.f.H = *r&0x0f == 0
	cpu.f.P = *r == 0x80
	cpu.f.N = false
	// panic(fmt.Sprintf("%08b", *r&0x0f))
}

func (cpu *z80) decR(r *byte) {
	cpu.f.H = *r&0x0f == 0
	*r--
	cpu.f.S = *r&0x80 != 0
	cpu.f.Z = *r == 0
	cpu.f.P = *r == 0x7f
	cpu.f.N = true
}

// ------

func (cpu *z80) rlc(r *byte) {
	*r = (*r << 1) | (*r >> 7)
	cpu.f.C = *r&1 != 0
	cpu.f.S = *r&0x80 != 0
	cpu.f.Z = *r == 0
	cpu.f.P = parityTable[*r]
	cpu.f.H = false
	cpu.f.N = false
}

func (cpu *z80) rrc(r *byte) {
	cpu.f.C = *r&1 != 0
	*r = (*r << 7) | (*r >> 1)
	cpu.f.S = *r&0x80 != 0
	cpu.f.Z = *r == 0
	cpu.f.P = parityTable[*r]
	cpu.f.H = false
	cpu.f.N = false
}

func (cpu *z80) sll(r *byte) {
	cpu.f.C = *r&0x80 != 0
	*r = byte((*r << 1) | 0x01)
	cpu.f.S = *r&0x80 != 0
	cpu.f.Z = *r == 0
	cpu.f.P = parityTable[*r]
	cpu.f.H = false
	cpu.f.N = false
}

func (cpu *z80) srl(r *byte) {
	cpu.f.C = *r&1 != 0
	*r = byte(*r >> 1)
	cpu.f.S = *r&0x80 != 0
	cpu.f.Z = *r == 0
	cpu.f.P = parityTable[*r]
	cpu.f.N = false
	cpu.f.H = false
}

func (cpu *z80) sla(r *byte) {
	cpu.f.C = *r&0x80 != 0
	*r = (*r << 1)
	cpu.f.S = *r&0x80 != 0
	cpu.f.Z = *r == 0
	cpu.f.P = parityTable[*r]
	cpu.f.H = false
	cpu.f.N = false
}

func (cpu *z80) sra(r *byte) {
	cpu.f.C = *r&0x1 != 0
	b7 := *r & 0b10000000
	*r = (*r >> 1) | b7
	cpu.f.S = *r&0x0080 != 0
	cpu.f.Z = *r == 0
	cpu.f.P = parityTable[*r]
	cpu.f.H = false
	cpu.f.N = false
}

func (cpu *z80) rr(r *byte) {
	c := cpu.f.C
	cpu.f.C = *r&0x1 != 0
	*r = (*r >> 1)
	if c {
		*r |= 0b10000000
	}
	cpu.f.S = *r&0x80 != 0
	cpu.f.Z = *r == 0
	cpu.f.P = parityTable[*r]
	cpu.f.H = false
	cpu.f.N = false
}

func (cpu *z80) rl(r *byte) {
	c := cpu.f.C
	cpu.f.C = *r&0x80 != 0
	*r = (*r << 1)
	if c {
		*r |= 0x1
	}
	cpu.f.S = *r&0x80 != 0
	cpu.f.Z = *r == 0
	cpu.f.P = parityTable[*r]
	cpu.f.H = false
	cpu.f.N = false
}

// ------

func (cpu *z80) addA(r byte) {
	a := int16(cpu.a)
	result := a + int16(r)
	lookup := ((cpu.a & 0x88) >> 3) | (((r) & 0x88) >> 2) | ((byte(result) & 0x88) >> 1)
	cpu.a = uint8(result & 0x00ff)

	cpu.f.S = cpu.a&0x80 != 0
	cpu.f.Z = cpu.a == 0
	cpu.f.H = halfcarryAddTable[lookup&0x07]
	cpu.f.P = overflowAddTable[lookup>>4]
	cpu.f.N = false
	cpu.f.C = ((result) & 0x100) != 0
}

func (cpu *z80) subA(r byte) {
	a := int16(cpu.a)
	result := a - int16(r)
	lookup := ((cpu.a & 0x88) >> 3) | (((r) & 0x88) >> 2) | ((byte(result) & 0x88) >> 1)
	cpu.a = uint8(result & 0x00ff)

	cpu.f.S = cpu.a&0x80 != 0
	cpu.f.Z = cpu.a == 0
	cpu.f.H = halfcarrySubTable[lookup&0x07]
	cpu.f.P = overflowSubTable[lookup>>4]
	cpu.f.N = true
	cpu.f.C = ((result) & 0x100) == 0x100
}

func (cpu *z80) jr(j byte) {
	jump := int8(j) + 2
	cpu.pc = uint16(int16(cpu.pc) + int16(jump))
}

func (cpu *z80) xor(s uint8) {
	cpu.a = cpu.a ^ s
	cpu.f.S = int8(cpu.a) < 0
	cpu.f.Z = cpu.a == 0
	cpu.f.H = false
	cpu.f.P = parityTable[cpu.a]
	cpu.f.N = false
	cpu.f.C = false
}

func (cpu *z80) and(s uint8) {
	cpu.a = cpu.a & s
	cpu.f.S = int8(cpu.a) < 0
	cpu.f.Z = cpu.a == 0
	cpu.f.H = true
	cpu.f.P = parityTable[cpu.a]
	cpu.f.N = false
	cpu.f.C = false
}

func (cpu *z80) or(s uint8) {
	// TODO: review p/v flag
	cpu.a = cpu.a | s
	cpu.f.S = int8(cpu.a) < 0
	cpu.f.Z = cpu.a == 0
	cpu.f.H = false
	cpu.f.P = parityTable[cpu.a]
	cpu.f.N = false
	cpu.f.C = false
}

func (cpu *z80) sbc(s byte) {
	res := uint16(cpu.a) - uint16(s)
	if cpu.f.C {
		res--
	}
	lookup := ((cpu.a & 0x88) >> 3) | ((s & 0x88) >> 2) | byte(res&0x88>>1)
	cpu.a = byte(res)
	cpu.f.S = cpu.a&0x0080 != 0
	cpu.f.Z = cpu.a == 0
	cpu.f.H = halfcarrySubTable[lookup&0x07]
	cpu.f.P = overflowSubTable[lookup>>4]
	cpu.f.N = true
	cpu.f.C = (res & 0x100) == 0x100
}

func (cpu *z80) sbcHL(ss uint16) {
	hl := getRR(cpu.h, cpu.l)
	res := uint32(hl) - uint32(ss)
	if cpu.f.C {
		res--
	}
	cpu.h, cpu.l = setRR(uint16(res))

	lookup := byte(((hl & 0x8800) >> 11) | ((ss & 0x8800) >> 10) | ((uint16(res) & 0x8800) >> 9))
	cpu.f.N = true
	cpu.f.S = cpu.h&0x80 != 0 // negative
	cpu.f.Z = res == 0
	cpu.f.C = (res & 0x10000) != 0
	cpu.f.P = overflowSubTable[lookup>>4]
	cpu.f.H = halfcarrySubTable[lookup&0x07]
}

func (cpu *z80) ldd() {
	bc := getRR(cpu.b, cpu.c)
	de := getRR(cpu.d, cpu.e)
	hl := getRR(cpu.h, cpu.l)

	v := cpu.memory.GetByte(hl)
	cpu.memory.PutByte(de, v)

	bc--
	de--
	hl--

	cpu.b, cpu.c = setRR(bc)
	cpu.d, cpu.e = setRR(de)
	cpu.h, cpu.l = setRR(hl)

	cpu.f.P = bc != 0
	cpu.f.H = false
	cpu.f.N = false
}

func (cpu *z80) ldi() {
	bc := getRR(cpu.b, cpu.c)
	de := getRR(cpu.d, cpu.e)
	hl := getRR(cpu.h, cpu.l)

	v := cpu.memory.GetByte(hl)
	cpu.memory.PutByte(de, v)

	bc--
	de++
	hl++

	cpu.b, cpu.c = setRR(bc)
	cpu.d, cpu.e = setRR(de)
	cpu.h, cpu.l = setRR(hl)

	cpu.f.P = bc != 0
	cpu.f.H = false
	cpu.f.N = false
}

//TODO Join these 3
func (cpu *z80) addHL(ss uint16) {
	hl := getRR(cpu.h, cpu.l)
	var result = uint32(hl) + uint32(ss)
	var lookup = byte(((hl & 0x0800) >> 11) | ((ss & 0x0800) >> 10) | ((uint16(result) & 0x0800) >> 9))
	cpu.h, cpu.l = setRR(uint16(result))

	cpu.f.N = false
	cpu.f.H = halfcarryAddTable[lookup]
	cpu.f.C = (result & 0x10000) != 0
}

func (cpu *z80) addIX(ss uint16) {
	ix := getRR(cpu.ixh, cpu.ixl)
	var result = uint32(ix) + uint32(ss)
	var lookup = byte(((ix & 0x0800) >> 11) | ((ss & 0x0800) >> 10) | ((uint16(result) & 0x0800) >> 9))
	cpu.ixh, cpu.ixl = setRR(uint16(result))

	cpu.f.N = false
	cpu.f.H = halfcarryAddTable[lookup]
	cpu.f.C = (result & 0x10000) != 0
}

func (cpu *z80) addIY(ss uint16) {
	iy := getRR(cpu.iyh, cpu.iyl)
	var result = uint32(iy) + uint32(ss)
	var lookup = byte(((iy & 0x0800) >> 11) | ((ss & 0x0800) >> 10) | ((uint16(result) & 0x0800) >> 9))
	cpu.iyh, cpu.iyl = setRR(uint16(result))

	cpu.f.N = false
	cpu.f.H = halfcarryAddTable[lookup]
	cpu.f.C = (result & 0x10000) != 0
}
