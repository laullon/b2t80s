package emulator

import (
	"fmt"
)

type Bus interface {
	SetAddr(uint16)
	GetAddr() uint16

	SetData(byte)
	GetData() byte

	ReadMemory()
	WriteMemory()

	RegisterPort(mask PortMask, manager PortManager)
	ReadPort()
	WritePort()

	// debuger
	GetBlock(addr uint16, l uint16) []byte
}

func NewBus(mem Memory) Bus {
	return &genericBus{
		mem:   mem,
		ports: make(map[PortMask]PortManager),
	}
}

type genericBus struct {
	mem   Memory
	addr  uint16
	data  uint8
	ports map[PortMask]PortManager
}

func (bus *genericBus) SetAddr(addr uint16) { bus.addr = addr }
func (bus *genericBus) GetAddr() uint16     { return bus.addr }

func (bus *genericBus) SetData(data byte) { bus.data = data }
func (bus *genericBus) GetData() byte     { return bus.data }

func (bus *genericBus) ReadMemory()  { bus.data = bus.mem.GetByte(bus.addr) }
func (bus *genericBus) WriteMemory() { bus.mem.PutByte(bus.addr, bus.data) }

func (bus *genericBus) RegisterPort(mask PortMask, manager PortManager) {
	bus.ports[mask] = manager
}

func (bus *genericBus) WritePort() {
	// fmt.Printf("[writePort]-> port:0x%04X data:%v  \n", bus.addr, bus.data)
	ok := false
	for portMask, portManager := range bus.ports {
		// fmt.Printf("[writePort] port:0x%04X (0x%04X)(0x%04X) data:%v\n", bus.addr, bus.addr&portMask.Mask, portMask.Value, bus.data)
		if (bus.addr & portMask.Mask) == portMask.Value {
			// println(reflect.TypeOf(portManager).String())
			portManager.WritePort(bus.addr, bus.data)
			ok = true
		}
	}
	if !ok {
		fmt.Printf("[writePort]-(no PM)-> port:0x%04X data:%v\n", bus.addr, bus.data)
		// panic("--")
	}
}

func (bus *genericBus) ReadPort() {
	skip := false
	// fmt.Printf(fmt.Sprintf("[readPort]-> port:0x%04X pc:0x%04X \n", port, cpu.regs.PC))
	for portMask, portManager := range bus.ports {
		if (bus.addr & portMask.Mask) == portMask.Value {
			// fmt.Printf("[readPort] (0x%04X) port:0x%04X (0x%04X)(0x%04X) \n", cpu.regs.PC, port, port&portMask.Mask, portMask.Value)
			// println(reflect.TypeOf(portManager).Elem().Name())
			bus.data, skip = portManager.ReadPort(bus.addr)
			if !skip {
				return
			}
		}
	}
	// panic(fmt.Sprintf("[readPort]-(no PM)-> port:0x%04X pc:0x%04X", port, cpu.regs.PC))
	fmt.Printf("[readPort]-(no PM)-> port:0x%04X  \n", bus.addr)
	bus.data = 0xff
}

func (bus *genericBus) GetBlock(addr uint16, l uint16) []byte {
	var res []byte
	for i := uint16(0); i < l; i++ {
		res = append(res, bus.mem.GetByte(addr+i))
	}
	return res
}
