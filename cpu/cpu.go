package cpu

type CPU interface {
	Interrupt(bool) // TODO:move this to each CPU interface
	NMI(bool)
	Halt()
	Wait(bool)
	Reset()
	Tick()

	CurrentOP() string

	SetTracer(CPUTracer)
	SetDebugger(DebuggerCallbacks)
}

type DebuggerCallbacks interface {
	Eval(pc uint16)
	EvalInterrupt()

	EvalLine()
	EvalFrame()
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
	SetNextOP(string)
	SetDiss(string)
}
