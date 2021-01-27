package atetris

import (
	"time"

	"github.com/laullon/b2t80s/emulator"
)

type watchdog struct {
	count int
	cpu   emulator.CPU
}

func (wd *watchdog) ReadPort(addr uint16) (byte, bool) { panic(-1) }
func (wd *watchdog) WritePort(addr uint16, data byte)  { wd.count++ }
func (wd *watchdog) start() {
	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for range ticker.C {
			if wd.count == 0 {
				wd.cpu.Reset()
			} else {
				wd.count = 0
			}
		}
	}()
}
