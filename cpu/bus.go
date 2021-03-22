package cpu

import (
	"fmt"
	"strings"
)

type Dumpable interface {
	Memory() []byte
}

type Bus interface {
	Write(addr uint16, data uint8)
	Read(addr uint16) uint8
	RegisterPort(name string, mask PortMask, manager PortManager)
	DumpMap() string
	GetDumplables() map[string]Dumpable
}

type busEntry struct {
	name    string
	mask    PortMask
	manager PortManager
}

//-----------------------------------------------

type bus struct {
	ports          []*busEntry
	defaultManager PortManager
}

func NewBus(defaultManager ...PortManager) Bus {
	bus := &bus{}
	if len(defaultManager) > 0 {
		bus.defaultManager = defaultManager[0]
	}
	return bus
}

func (bus *bus) Write(addr uint16, data uint8) {
	// fmt.Printf("[writePort]-> port:0x%04X data:0x%02X  \n", addr, data)
	for _, entry := range bus.ports {
		if (addr & entry.mask.Mask) == entry.mask.Value {
			// println(entry.name, " - ", reflect.TypeOf(entry.manager).String())
			entry.manager.WritePort(addr, data)
			return
		}
	}
	if bus.defaultManager != nil {
		bus.defaultManager.WritePort(addr, data)
	} else {
		panic(fmt.Sprintf("[writePort]-(no PM)-> port:0x%04X data:%v\n", addr, data))
	}
}

func (bus *bus) Read(addr uint16) uint8 {
	for _, entry := range bus.ports {
		if (addr & entry.mask.Mask) == entry.mask.Value {
			// fmt.Printf("[readPort] port:0x%04X (0x%04X)(0x%04X) \n", addr, addr&portMask.Mask, portMask.Value)
			// println(reflect.TypeOf(portManager).Elem().Name())
			data, _ := entry.manager.ReadPort(addr)
			// fmt.Printf(fmt.Sprintf("[readPort]-> port:0x%04X data:0x%02X \n", addr, data))
			return data
		}
	}

	if bus.defaultManager != nil {
		data, _ := bus.defaultManager.ReadPort(addr)
		return data
	} else {
		panic(fmt.Sprintf("[readPort]-(no PM)-> port:0x%04X", addr))
	}
}

func (bus *bus) RegisterPort(name string, mask PortMask, manager PortManager) {
	bus.ports = append(bus.ports, &busEntry{name, mask, manager})
}

func (bus *bus) DumpMap() string {
	addrs := make([]string, 0x10000)
	for addr := uint16(0); ; {
		for _, entry := range bus.ports {
			if (addr & entry.mask.Mask) == entry.mask.Value {
				addrs[addr] = entry.name
			}
		}
		addr++
		if addr == 0 {
			break
		}
	}

	var res strings.Builder
	actualName := addrs[0]
	firstAddr := 0
	lastAddr := 0
	for addr, name := range addrs {
		if actualName != name {
			res.WriteString(fmt.Sprintf("0x%04X - 0x%04X = %s\n", firstAddr, lastAddr, actualName))
			actualName = name
			firstAddr = addr
		}
		lastAddr = addr
	}
	res.WriteString(fmt.Sprintf("0x%04X - 0x%04X = %s\n", firstAddr, lastAddr, actualName))
	return res.String()
}

func (bus *bus) GetDumplables() map[string]Dumpable {
	res := make(map[string]Dumpable)
	for _, entry := range bus.ports {
		if d, ok := entry.manager.(Dumpable); ok {
			res[entry.name] = d
		}
	}
	return res
}

//-----------------------------------------------

type RAM interface {
	PortManager
	SetBank([]byte)
}

type ram struct {
	bank []byte
	mask uint16
}

func NewRAM(bank []byte, mask uint16) RAM {
	ram := &ram{bank: bank, mask: mask}
	return ram
}

func (ram *ram) SetBank(bank []byte) {
	ram.bank = bank
}

func (ram *ram) ReadPort(addr uint16) (byte, bool) {
	return ram.bank[addr&ram.mask], false
}

func (ram *ram) WritePort(addr uint16, data byte) {
	ram.bank[addr&ram.mask] = data
	// fmt.Printf("---> 0x%04X(0x%04X) = 0x%02X (0x%02X) \n", addr, addr&ram.mask, data, ram.bank[addr&ram.mask])
}

func (ram *ram) Memory() []byte {
	return ram.bank
}

//-----------------------------------------------

type RomWrite func(uint16, uint8)

type ROM interface {
	PortManager
	SetBank([]byte)
}

type rom struct {
	bank  []byte
	mask  uint16
	write RomWrite
}

func NewROM(bank []byte, mask uint16, write ...RomWrite) ROM {
	rom := &rom{bank: bank, mask: mask}
	if len(write) > 0 {
		rom.write = write[0]
	}
	return rom
}

func (rom *rom) SetBank(bank []byte) {
	rom.bank = bank
}

func (rom *rom) ReadPort(addr uint16) (byte, bool) {
	return rom.bank[addr&rom.mask], false
}

func (rom *rom) WritePort(addr uint16, data byte) {
	if rom.write != nil {
		rom.write(addr, data)
	} else {
		// panic(fmt.Sprintf("Write no allowed on 0x%04X", addr))
	}
}

func (rom *rom) Memory() []byte {
	return rom.bank
}

//-----------------------------------------------

type watchableBus struct {
	bus Bus
	// wps []uint16
}

func NewWatchableBus(bus Bus) Bus {
	// wps := make([]uint16, 0)
	// watchPoints := *WatchPoints
	// if len(watchPoints) > 0 {
	// 	bps := strings.Split(watchPoints, ",")
	// 	for _, bp := range bps {
	// 		n, err := strconv.ParseUint(bp, 0, 16)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		wps = append(wps, uint16(n))
	// 	}
	// 	fmt.Printf("watchPoints: %v\n", watchPoints)
	// }

	return &watchableBus{
		bus: bus,
		// wps: wps,
	}
}

func (bus *watchableBus) Write(addr uint16, data uint8) {
	// for _, wp := range bus.wps {
	// 	if wp == addr {
	// 		DebuggerCTL.Stop()
	// 		println("--r--")
	// 	}
	// }
	bus.bus.Write(addr, data)
}

func (bus *watchableBus) Read(addr uint16) uint8 {
	// for _, wp := range bus.wps {
	// 	if wp == addr {
	// 		// DebuggerCTL.Stop()
	// 		println("--r--")
	// 	}
	// }
	return bus.bus.Read(addr)
}

func (bus *watchableBus) RegisterPort(name string, mask PortMask, manager PortManager) {
	bus.bus.RegisterPort(name, mask, manager)
}

func (bus *watchableBus) DumpMap() string {
	return bus.bus.DumpMap()
}

func (bus *watchableBus) GetDumplables() map[string]Dumpable {
	return bus.bus.GetDumplables()

}
