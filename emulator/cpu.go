package emulator

type CPUTrap func()
type CPU interface {
	Interrupt(bool)
	Halt()
	Tick()

	Registers() interface{}

	SetDebuger(debugger Debugger)
	RegisterTrap(pc uint16, trap CPUTrap)
}

type Debugger interface {
	AddInstruction([]byte)

	Tick()
	Stop()
	Continue()
	Step()
	StopNextFrame()

	GetStatus() string
}

type PortMask struct {
	Mask  uint16
	Value uint16
}

type PortManager interface {
	ReadPort(port uint16) (byte, bool)
	WritePort(port uint16, data byte)
}
