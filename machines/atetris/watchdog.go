package atetris

import (
	"github.com/laullon/b2t80s/cpu"
)

type watchdog struct {
	count int
	cpu   cpu.CPU
}

func (wd *watchdog) ReadPort(port uint16) byte        { panic(-1) }
func (wd *watchdog) WritePort(addr uint16, data byte) { wd.count++ }
func (wd *watchdog) start() {
	// ticker := time.NewTicker(time.Second * 2)
	// go func() {
	// 	for range ticker.C {
	// 		if wd.count == 0 {
	// 			wd.cpu.Reset()
	// 		} else {
	// 			wd.count = 0
	// 		}
	// 	}
	// }()
}
