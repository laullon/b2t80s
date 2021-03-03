package atetris

import (
	"fmt"
	"os"

	"github.com/laullon/b2t80s/cpu"
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

type status struct {
	romPage byte
}

func (s *status) ReadPort(addr uint16) (byte, bool) { panic(-1) }
func (s *status) WritePort(addr uint16, data byte)  { s.romPage = data & 0b00000111 }

// ----------------------------
type ram struct {
	mem   []byte
	mask  uint16
	cpu   cpu.CPU
	trace bool
}

func (ram *ram) ReadPort(addr uint16) (byte, bool) {
	if ram.trace {
		fmt.Printf("-> read  0x%04X = 0x%02x\n", addr&ram.mask, ram.mem[addr&ram.mask])
	}
	return ram.mem[addr&ram.mask], false
}

func (ram *ram) WritePort(addr uint16, data byte) {
	if ram.trace {
		fmt.Printf("-> write 0x%04X = 0x%02x\n", addr&ram.mask, data)
	}
	ram.mem[addr&ram.mask] = data
}

func (ram *ram) Memory() []byte {
	return ram.mem
}

// ----------------------------
type eepromStatus struct {
	lock bool
}

func (s *eepromStatus) ReadPort(addr uint16) (byte, bool) { panic(-1) }
func (s *eepromStatus) WritePort(addr uint16, data byte) {
	s.lock = true
}

// ----------------------------
type eeprom struct {
	mem    []byte
	mask   uint16
	status *eepromStatus
}

func newEERPROM() *eeprom {
	data := make([]byte, 0x0200)
	f, err := os.Open("nvram.tetris")
	defer f.Close()
	if err != nil {
		for i := 0; i < len(data); i++ {
			data[i] = 0xff
		}
	} else {
		_, err := f.Read(data)
		if err != nil {
			panic(err)
		}
	}

	println("data:", data)
	if len(data) != 0x0200 {
	}
	eeprom := &eeprom{
		mem:    data,
		mask:   0x01ff,
		status: &eepromStatus{},
	}
	return eeprom
}

func (eeprom *eeprom) ReadPort(addr uint16) (byte, bool) {
	if eeprom.status.lock {
		return 0, false
	}
	return eeprom.mem[addr&eeprom.mask], false
}

func (eeprom *eeprom) WritePort(addr uint16, data byte) {
	if eeprom.status.lock {
		eeprom.mem[addr&eeprom.mask] = data
		eeprom.status.lock = false
		f, err := os.Create("nvram.tetris")
		if err != nil {
			panic(err)
		}
		_, err = f.Write(eeprom.mem)
		if err != nil {
			panic(err)
		}
	}
}

func (eeprom *eeprom) Memory() []byte {
	return eeprom.mem
}

// ----------------------------
type fixedROM struct {
	rom []byte
}

func (rom *fixedROM) ReadPort(addr uint16) (byte, bool) { return rom.rom[addr], false }
func (rom *fixedROM) WritePort(addr uint16, data byte)  { panic(-1) }

// ----------------------------
type clearIRQ struct {
	cpu cpu.CPU
}

func (s *clearIRQ) ReadPort(addr uint16) (byte, bool) { panic(-1) }
func (s *clearIRQ) WritePort(addr uint16, data byte)  { s.cpu.Interrupt(false) }
