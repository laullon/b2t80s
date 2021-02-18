package atetris

import (
	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
)

// 0000-0FFF   R/W   xxxxxxxx    Program RAM
// 1000-1FFF   R/W   xxxxxxxx    Playfield RAM
//                   xxxxxxxx       (byte 0: LSB of character code)
//                   -----xxx       (byte 1: MSB of character code)
//                   xxxx----       (byte 1: palette index)
// 2000-20FF   R/W   xxxxxxxx    Palette RAM
//                   xxx----        (red component)
//                   ---xxx--       (green component)
//                   ------xx       (blue component)
// 2400-25FF   R/W   xxxxxxxx    EEPROM
// 2800-280F   R/W   xxxxxxxx    POKEY #1
// 2810-281F   R/W   xxxxxxxx    POKEY #2
// 3000          W   --------    Watchdog
// 3400          W   --------    EEPROM write enable
// 3800          W   --------    IRQ acknowledge
// 3C00          W   --xx-xxx    System status
//               W   --x-----       (right coin counter)
//               W   ---x----       (left coin counter)
//               W   -----xxx       (page rom bank)
// 4000-7FFF   R     xxxxxxxx    Banked program ROM
// 8000-FFFF   R     xxxxxxxx    Program ROM

func newBus() m6502.Bus {
	bus := m6502.NewBus()

	rom := loadRom("136066-1100.45f")
	status := &status{}

	eeprom := &eeprom{
		mem:    make([]byte, 0x0200),
		mask:   0x01ff,
		status: &eepromStatus{},
	}

	// RAM
	bus.RegisterPort("ram", emulator.PortMask{Mask: 0b1111000000000000, Value: 0b0000000000000000}, &ram{mem: make([]byte, 0x1000), mask: 0x0fff})

	// ROM
	bus.RegisterPort("slapstic", emulator.PortMask{Mask: 0b1100000000000000, Value: 0b0100000000000000}, newSlapstic(rom))
	bus.RegisterPort("rom", emulator.PortMask{Mask: 0b1000000000000000, Value: 0b1000000000000000}, &fixedROM{rom: rom})

	// EEPROM
	bus.RegisterPort("eeprom.status", emulator.PortMask{Mask: 0b1111110000000000, Value: 0b0011010000000000}, eeprom.status)
	bus.RegisterPort("eeprom", emulator.PortMask{Mask: 0b1111110000000000, Value: 0b0010010000000000}, eeprom)

	// STATUS
	bus.RegisterPort("status", emulator.PortMask{Mask: 0b1111110000000000, Value: 0b0011110000000000}, status)

	return bus
}

type status struct {
	romPage byte
}

func (s *status) ReadPort(addr uint16) (byte, bool) { panic(-1) }
func (s *status) WritePort(addr uint16, data byte)  { s.romPage = data & 0b00000111 }

// ----------------------------
type ram struct {
	mem  []byte
	mask uint16
	cpu  emulator.CPU
}

func (ram *ram) ReadPort(addr uint16) (byte, bool) { return ram.mem[addr&ram.mask], false }
func (ram *ram) WritePort(addr uint16, data byte)  { ram.mem[addr&ram.mask] = data }

// ----------------------------
type eepromStatus struct {
	lock bool
}

func (s *eepromStatus) ReadPort(addr uint16) (byte, bool) { panic(-1) }
func (s *eepromStatus) WritePort(addr uint16, data byte)  { s.lock = false }

// ----------------------------
type eeprom struct {
	mem    []byte
	mask   uint16
	status *eepromStatus
}

func (eeprom *eeprom) ReadPort(addr uint16) (byte, bool) { return eeprom.mem[addr&eeprom.mask], false }
func (eeprom *eeprom) WritePort(addr uint16, data byte) {
	if !eeprom.status.lock {
		println("----")
		eeprom.mem[addr&eeprom.mask] = data
		eeprom.status.lock = true
	}
}

// ----------------------------
type fixedROM struct {
	rom []byte
}

func (rom *fixedROM) ReadPort(addr uint16) (byte, bool) { return rom.rom[addr], false }
func (rom *fixedROM) WritePort(addr uint16, data byte)  { panic(-1) }
