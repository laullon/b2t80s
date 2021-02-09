package m6502

import (
	"fmt"

	"github.com/laullon/b2t80s/emulator"
)

type Bus interface {
	Write(addr uint16, data uint8)
	Read(addr uint16) uint8
	RegisterPort(mask emulator.PortMask, manager emulator.PortManager)
}

type bus struct {
	ports map[emulator.PortMask]emulator.PortManager
}

func NewBus() Bus {
	bus := &bus{
		ports: make(map[emulator.PortMask]emulator.PortManager),
	}
	return bus
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
		panic(fmt.Sprintf("[writePort]-(no PM)-> port:0x%04X data:%v\n", addr, data))
	}
}

func (bus *bus) Read(addr uint16) uint8 {
	skip := false
	data := uint8(0)
	// fmt.Printf(fmt.Sprintf("[readPort]-> port:0x%04X pc:0x%04X \n", port, cpu.regs.PC))
	for portMask, portManager := range bus.ports {
		if (addr & portMask.Mask) == portMask.Value {
			// fmt.Printf("[readPort] port:0x%04X (0x%04X)(0x%04X) \n", addr, addr&portMask.Mask, portMask.Value)
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
	if _, found := bus.ports[mask]; found {
		delete(bus.ports, mask)
	}
	bus.ports[mask] = manager
}

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
