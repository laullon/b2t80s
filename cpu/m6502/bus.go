package m6502

import (
	"fmt"
	"strings"

	"github.com/laullon/b2t80s/emulator"
)

type Bus interface {
	Write(addr uint16, data uint8)
	Read(addr uint16) uint8
	RegisterPort(name string, mask emulator.PortMask, manager emulator.PortManager)
	DumpMap() string
}

type busEntry struct {
	name    string
	mask    emulator.PortMask
	manager emulator.PortManager
}

type bus struct {
	ports map[string]*busEntry
}

func NewBus() Bus {
	bus := &bus{
		ports: make(map[string]*busEntry),
	}
	return bus
}

func (bus *bus) Write(addr uint16, data uint8) {
	// fmt.Printf("[writePort]-> port:0x%04X data:%v  \n", addr, data)
	ok := false
	for _, entry := range bus.ports {
		if (addr & entry.mask.Mask) == entry.mask.Value {
			// println(reflect.TypeOf(portManager).String())
			entry.manager.WritePort(addr, data)
			ok = true
		}
	}
	if !ok {
		panic(fmt.Sprintf("[writePort]-(no PM)-> port:0x%04X data:%v\n", addr, data))
	}
}

func (bus *bus) Read(addr uint16) uint8 {
	skip := false
	data := uint8(0)
	// fmt.Printf(fmt.Sprintf("[readPort]-> port:0x%04X pc:0x%04X \n", port, cpu.regs.PC))
	for _, entry := range bus.ports {
		if (addr & entry.mask.Mask) == entry.mask.Value {
			// fmt.Printf("[readPort] port:0x%04X (0x%04X)(0x%04X) \n", addr, addr&portMask.Mask, portMask.Value)
			// println(reflect.TypeOf(portManager).Elem().Name())
			data, skip = entry.manager.ReadPort(addr)
			if !skip {
				return data
			}
		}
	}
	panic(fmt.Sprintf("[readPort]-(no PM)-> port:0x%04X", addr))
}

func (bus *bus) RegisterPort(name string, mask emulator.PortMask, manager emulator.PortManager) {
	if _, found := bus.ports[name]; found {
		panic(fmt.Sprintf("port '%s' already registered", name))
	}
	bus.ports[name] = &busEntry{name, mask, manager}
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

////------------

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
