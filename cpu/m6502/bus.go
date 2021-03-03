package m6502

import (
	"fmt"
	"strings"

	"github.com/laullon/b2t80s/cpu"
)

type Dumpable interface {
	Memory() []byte
}

type Bus interface {
	Write(addr uint16, data uint8)
	Read(addr uint16) uint8
	RegisterPort(name string, mask cpu.PortMask, manager cpu.PortManager)
	DumpMap() string
	GetDumplables() map[string]Dumpable
}

type busEntry struct {
	name    string
	mask    cpu.PortMask
	manager cpu.PortManager
}

//-----------------------------------------------

type bus struct {
	ports []*busEntry
}

func NewBus() Bus {
	bus := &bus{
		// ports: make([]*busEntry,0),
	}
	return bus
}

func (bus *bus) Write(addr uint16, data uint8) {
	// fmt.Printf("[writePort]-> port:0x%04X data:0x%02X  \n", addr, data)
	for _, entry := range bus.ports {
		if (addr & entry.mask.Mask) == entry.mask.Value {
			// println(reflect.TypeOf(portManager).String())
			entry.manager.WritePort(addr, data)
			return
		}
	}
	panic(fmt.Sprintf("[writePort]-(no PM)-> port:0x%04X data:%v\n", addr, data))
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
	panic(fmt.Sprintf("[readPort]-(no PM)-> port:0x%04X", addr))
}

func (bus *bus) RegisterPort(name string, mask cpu.PortMask, manager cpu.PortManager) {
	for _, port := range bus.ports {
		if port.name == name {
			panic(fmt.Sprintf("port '%s' already registered", name))
		}

	}
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

type BasicRam struct {
	Data  []byte
	Mask  uint16
	Trace bool
}

func (ram *BasicRam) ReadPort(addr uint16) (byte, bool) {
	if ram.Trace {
		fmt.Printf("[ram] read 0x%04X \n", addr)
	}
	return ram.Data[addr&ram.Mask], false
}

func (ram *BasicRam) WritePort(addr uint16, data byte) {
	if ram.Trace {
		fmt.Printf("[ram] write 0x%04X 0x%02x\n", addr, data)
	}
	ram.Data[addr&ram.Mask] = data
}

//-----------------------------------------------

type watchableBus struct {
	bus Bus
	// wps []uint16
}

func NewWatchableBus(bus Bus) Bus {
	// wps := make([]uint16, 0)
	// watchPoints := *cpu.WatchPoints
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
	// 		cpu.DebuggerCTL.Stop()
	// 		println("--r--")
	// 	}
	// }
	bus.bus.Write(addr, data)
}

func (bus *watchableBus) Read(addr uint16) uint8 {
	// for _, wp := range bus.wps {
	// 	if wp == addr {
	// 		// cpu.DebuggerCTL.Stop()
	// 		println("--r--")
	// 	}
	// }
	return bus.bus.Read(addr)
}

func (bus *watchableBus) RegisterPort(name string, mask cpu.PortMask, manager cpu.PortManager) {
	bus.bus.RegisterPort(name, mask, manager)
}

func (bus *watchableBus) DumpMap() string {
	return bus.bus.DumpMap()
}

func (bus *watchableBus) GetDumplables() map[string]Dumpable {
	return bus.bus.GetDumplables()

}
