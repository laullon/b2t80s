package m6502

func cld(cpu *m6502) { cpu.regs.P.D = false }

func tsx(cpu *m6502) { cpu.regs.SP = cpu.regs.X }

func ldx(cpu *m6502, data uint8) {
	cpu.regs.X = data
	ldzn(cpu, data)
}

func lda(cpu *m6502, data uint8) {
	cpu.regs.A = data
	ldzn(cpu, data)
}

func ldzn(cpu *m6502, data uint8) {
	cpu.regs.P.Z = data == 0
	cpu.regs.P.Z = data&0x80 != 0
}
