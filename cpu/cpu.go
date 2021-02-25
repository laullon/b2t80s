package cpu

type CPUTrap func()

type CPU interface {
	Interrupt(bool)
	NMI(bool)
	Halt()
	Wait(bool)
	Reset()
	Tick()

	// RegisterTrap(pc uint16, trap CPUTrap)

	CurrentOP() string

	SetTracer(CPUTracer)
	SetDebugger(DebuggerCallbacks)
}

type DebuggerCallbacks interface {
	Eval()
	EvalInterrupt()
}

type PortMask struct {
	Mask  uint16
	Value uint16
}

type PortManager interface {
	ReadPort(port uint16) (byte, bool)
	WritePort(port uint16, data byte)
}

type CPUTracer interface {
	AppendLastOP(string)
}
