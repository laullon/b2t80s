package emulator

type CPUTrap func()
type CPU interface {
	Interrupt(bool)
	Halt()
	Wait(bool)
	Reset()
	Tick()

	Registers() interface{}

	SetDebuger(debugger Debugger)
	RegisterTrap(pc uint16, trap CPUTrap)

	CurrentOP() string
}

type Debugger interface {
	AddInstruction(pc uint16, mem, instruction string)
	// NextInstruction([]byte)

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
