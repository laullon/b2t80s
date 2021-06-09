package z80

import (
	"github.com/laullon/b2t80s/cpu"
)

type Bus interface {
	SetAddr(uint16)
	GetAddr() uint16

	SetData(byte)
	GetData() byte

	Release()

	ReadMemory()
	WriteMemory()

	ReadPort()
	WritePort()
}

func NewBus(mem cpu.Bus, ports cpu.Bus) Bus {
	return &z80bus{
		mem:   mem,
		ports: ports,
	}
}

type z80bus struct {
	addr uint16
	data uint8

	mem   cpu.Bus
	ports cpu.Bus
}

func (bus *z80bus) SetAddr(addr uint16) { bus.addr = addr }
func (bus *z80bus) GetAddr() uint16     { return bus.addr }

func (bus *z80bus) SetData(data byte) { bus.data = data }
func (bus *z80bus) GetData() byte     { return bus.data }

func (bus *z80bus) Release() { bus.addr = 0xffff; bus.data = 0xff }

func (bus *z80bus) ReadMemory()  { bus.data = bus.mem.Read(bus.addr) }
func (bus *z80bus) WriteMemory() { bus.mem.Write(bus.addr, bus.data) }

func (bus *z80bus) ReadPort()  { bus.data = bus.ports.Read(bus.addr) }
func (bus *z80bus) WritePort() { bus.ports.Write(bus.addr, bus.data) }
