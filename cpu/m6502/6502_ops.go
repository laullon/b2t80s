package m6502

var ops []operation

func init() {
	ops = make([]operation, 0x100)

	ops[0x9a] = &implicit{f: tsx}
	ops[0xa2] = &immediate{f: ldx}
	ops[0xa9] = &immediate{f: lda}
	ops[0xd8] = &implicit{f: cld}
}

type operation interface {
	done() bool
	tick(cpu *m6502)
}

type bacicop struct {
	d bool
	t uint
}

func (op *bacicop) done() bool {
	return op.d
}

type reset struct {
	bacicop
}

func (op *reset) tick(cpu *m6502) {
	switch op.t {
	case 0:
		cpu.regs.SP = 0
	case 3, 4, 5:
		cpu.regs.SP--
	case 6:
		cpu.regs.PC = uint16(cpu.mem[0xff00+uint16(cpu.regs.SP)])<<8 | (cpu.regs.PC & 0x00ff)
		cpu.regs.SP--
	case 7:
		cpu.regs.PC = (cpu.regs.PC & 0xff00) | uint16(cpu.mem[0xff00+uint16(cpu.regs.SP)])
		op.d = true
	}
	op.t++
}

// -----

type implicit struct {
	bacicop
	f func(cpu *m6502)
}

func (op *implicit) tick(cpu *m6502) {
	op.f(cpu)
	op.d = true
}

// -----

type immediate struct {
	bacicop
	f func(cpu *m6502, data uint8)
}

func (op *immediate) tick(cpu *m6502) {
	op.f(cpu, cpu.mem[int(cpu.regs.PC)])
	cpu.regs.PC++
	op.d = true
}
