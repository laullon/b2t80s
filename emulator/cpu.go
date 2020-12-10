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
	AddInstruction(uint16, []byte)
	NextInstruction([]byte)

	Tick()

	Stop()
	Continue()
	Step()
	StopNextFrame()
	SetDump(bool)

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
