package cpu

import "strings"

type RegPair struct {
	H, L *byte
}

func (reg *RegPair) Get() uint16 {
	return uint16(*reg.H)<<8 | uint16(*reg.L)
}

func (reg *RegPair) Set(hl uint16) {
	*reg.H = byte(hl >> 8)
	*reg.L = byte(hl & 0x00ff)
}

type Log interface {
	AddEntry(entry string)
	Print() string
}

type logTail struct {
	idx     uint8
	entries []string
}

func NewLogTail() Log {
	return &logTail{
		entries: make([]string, 0x100),
	}
}

func (log *logTail) AddEntry(entry string) {
	log.entries[log.idx] = entry
	log.idx++
}

func (log *logTail) Print() string {
	var sb strings.Builder
	sb.WriteString(strings.Join(log.entries[log.idx:], "\n"))
	sb.WriteString(strings.Join(log.entries[:log.idx], "\n"))
	return sb.String()
}
