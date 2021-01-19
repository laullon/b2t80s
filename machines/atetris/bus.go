package atetris

import (
	"archive/zip"
	"fmt"

	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/utils"
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
// 3C00          W   --xx----    Coin counters
//               W   --x-----       (right coin counter)
//               W   ---x----       (left coin counter)
// 4000-7FFF   R     xxxxxxxx    Banked program ROM
// 8000-FFFF   R     xxxxxxxx    Program ROM

func newBus() m6502.Bus {
	bus := &bus{
		ports: make(map[emulator.PortMask]emulator.PortManager),
	}

	rom := loadRom()

	bus.RegisterPort(emulator.PortMask{Mask: 0x8000, Value: 0x8000}, &fixedROM{rom: rom})
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
type fixedROM struct {
	rom []byte
}

func (rom *fixedROM) ReadPort(addr uint16) (byte, bool) { return rom.rom[addr], false }
func (rom *fixedROM) WritePort(port uint16, data byte)  { panic(-1) }

func loadRom() []byte {
	zipFile := "../../games/atetris.zip"
	var mem []byte
	zf, err := zip.OpenReader(zipFile)
	if err != nil {
		panic(err)
	}

	for _, file := range zf.File {
		println(file.Name)
		if file.Name == "136066-1100.45f" {
			mem = utils.ReadZipFile(file)
		}
	}

	err = zf.Close()
	if err != nil {
		panic(err)
	}
	return mem
}
