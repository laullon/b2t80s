package lr35902

import (
	"github.com/laullon/b2t80s/cpu"
)

func newBus(bus cpu.Bus) *genericBus {
	return &genericBus{bus: bus}
}

type genericBus struct {
	bus  cpu.Bus
	addr uint16
	data uint8
}

func (bus *genericBus) SetAddr(addr uint16) { bus.addr = addr }
func (bus *genericBus) GetAddr() uint16     { return bus.addr }

func (bus *genericBus) SetData(data byte) { bus.data = data }
func (bus *genericBus) GetData() byte     { return bus.data }

func (bus *genericBus) Release() { bus.addr = 0xffff; bus.data = 0xff }

func (bus *genericBus) Write() { bus.bus.Write(bus.addr, bus.data) }
func (bus *genericBus) Read()  { bus.data = bus.bus.Read(bus.addr) }
