package cpu

type CPU interface {
	Halt()
	Wait(bool)
	Reset()
	Tick()

	SetTracer(CPUTracer)
	SetDebugger(DebuggerCallbacks)
}

type DebuggerCallbacks interface {
	Eval(pc uint16)
	EvalInterrupt()

	EvalLine() bool
	EvalFrame() bool
}

type PortMask struct {
	Mask  uint16
	Value uint16
}

type PortManager interface {
	ReadPort(port uint16) byte
	WritePort(port uint16, data byte)
}

type CPUTracer interface {
	AppendLastOP(string)
	SetNextOP(string)
	SetDiss(pc uint16, getMemory func(pc, leng uint16) []byte)
}
