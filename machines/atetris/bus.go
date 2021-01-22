package atetris

import (
	"fmt"

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
	bus := &bus{
		ports: make(map[emulator.PortMask]emulator.PortManager),
	}

	rom := loadRom("136066-1100.45f")
	status := &status{}

	//RAM
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111000000000000, Value: 0b0000000000000000}, &ram{mem: make([]byte, 0x1000), mask: 0x0fff})

	// ROM
	bus.RegisterPort(emulator.PortMask{Mask: 0b1100000000000000, Value: 0b0100000000000000}, newSlapstic(rom))
	bus.RegisterPort(emulator.PortMask{Mask: 0b1000000000000000, Value: 0b1000000000000000}, &fixedROM{rom: rom})

	// EEPROM
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111110000000000, Value: 0b0011010000000000}, &eepromStatus{})
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111110000000000, Value: 0b0010010000000000}, &eeprom{mem: make([]byte, 0x01ff), mask: 0x01ff})

	// STATUS
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111110000000000, Value: 0b0011110000000000}, status)

	// Unused
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111110000000000, Value: 0b0010110000000000}, &ram{mem: make([]byte, 0x01fff), mask: 0x00ff, debug: true})

	// Watchdog
	bus.RegisterPort(emulator.PortMask{Mask: 0b1111110000000000, Value: 0b0011000000000000}, &ram{mem: make([]byte, 0x0100), mask: 0x00ff})

	return bus
}

type bus struct {
	ports map[emulator.PortMask]emulator.PortManager
}

func (bus *bus) Write(addr uint16, data uint8) {
	// fmt.Printf("[writePort]-> port:0x%04X data:%v  \n", addr, data)
	ok := false
	for portMask, portManager := range bus.ports {
		// fmt.Printf("[writePort] port:0x%04X (0x%04X)(0x%04X) data:%v\n", addr, addr&portMask.Mask, portMask.Value, data)
		if (addr & portMask.Mask) == portMask.Value {
			// println(reflect.TypeOf(portManager).String())
			portManager.WritePort(addr, data)
			ok = true
		}
	}
	if !ok {
		fmt.Printf("[writePort]-(no PM)-> port:0x%04X data:%v\n", addr, data)
		panic("--")
	}
}

func (bus *bus) Read(addr uint16) uint8 {
	skip := false
	data := uint8(0)
	// fmt.Printf(fmt.Sprintf("[readPort]-> port:0x%04X pc:0x%04X \n", port, cpu.regs.PC))
	for portMask, portManager := range bus.ports {
		if (addr & portMask.Mask) == portMask.Value {
			// fmt.Printf("[readPort] (0x%04X) port:0x%04X (0x%04X)(0x%04X) \n", cpu.regs.PC, port, port&portMask.Mask, portMask.Value)
			// println(reflect.TypeOf(portManager).Elem().Name())
			data, skip = portManager.ReadPort(addr)
			if !skip {
				return data
			}
		}
	}
	panic(fmt.Sprintf("[readPort]-(no PM)-> port:0x%04X", addr))
	data = 0xff
	return data
}

func (bus *bus) RegisterPort(mask emulator.PortMask, manager emulator.PortManager) {
	bus.ports[mask] = manager
}

// ----------------------------
type eepromStatus struct {
	lock bool
}

// TODO: this should panic, but STA instruction read the byte but not use it
//       STA should not read teh byte.
func (s *eepromStatus) ReadPort(addr uint16) (byte, bool) {
	return 0x00, false
}

func (s *eepromStatus) WritePort(addr uint16, data byte) {
	s.lock = !s.lock
}

// ----------------------------
type status struct {
	romPage byte
}

// TODO: See STA comment above
func (s *status) ReadPort(addr uint16) (byte, bool) { return 0x00, false }
func (s *status) WritePort(addr uint16, data byte) {
	s.romPage = data & 0b00000111
	println("romPage:", s.romPage)
}

// ----------------------------
type watchdog struct {
}

func (wd *watchdog) ReadPort(addr uint16) (byte, bool) { panic(-1) }
func (wd *watchdog) WritePort(addr uint16, data byte)  {}

// ----------------------------
type ram struct {
	mem   []byte
	mask  uint16
	debug bool
}

func (ram *ram) ReadPort(addr uint16) (byte, bool) { return ram.mem[addr&ram.mask], false }
func (ram *ram) WritePort(addr uint16, data byte) {
	ram.mem[addr&ram.mask] = data
	if ram.debug {
		fmt.Printf("W (0x%04X)0x%04X -> 0x%02x \n", addr, addr&ram.mask, data)
	}
}

// ----------------------------
type eeprom struct {
	mem   []byte
	mask  uint16
	debug bool
}

func (eeprom *eeprom) ReadPort(addr uint16) (byte, bool) { return eeprom.mem[addr&eeprom.mask], false }
func (eeprom *eeprom) WritePort(addr uint16, data byte) {
	eeprom.mem[addr&eeprom.mask] = data
	if eeprom.debug {
		fmt.Printf("W (0x%04X)0x%04X -> 0x%02x \n", addr, addr&eeprom.mask, data)
	}
}

// ----------------------------
type fixedROM struct {
	rom   []byte
	debug bool
}

func (rom *fixedROM) ReadPort(addr uint16) (byte, bool) {
	if rom.debug {
		fmt.Printf("R 0x%04X \n", addr)
	}
	return rom.rom[addr], false
}
func (rom *fixedROM) WritePort(addr uint16, data byte) { fmt.Printf("0x%04x\n", addr); panic(-1) }
