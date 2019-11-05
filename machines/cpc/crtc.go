package cpc

import (
	"fmt"

	"github.com/laullon/b2t80s/emulator"
)

var (
	regMasks = []byte{0xff, 0xff, 0xff, 0xff, 0x7f, 0x1f, 0x7f, 0x7f, 0xff, 0x1f, 0x7f, 0x1f, 0x3f, 0xff, 0x3f, 0xff}
)

type crtcStatus struct {
	vSync  bool
	hSync  bool
	disPen bool
	ma     uint16
}

type crtcCounters struct {
	h      int
	sl     int
	row    int
	raster int
}

type crtc struct {
	cpu emulator.CPU

	status   *crtcStatus
	counters *crtcCounters

	regs      []byte
	selectReg byte

	cycles uint32
	clock  uint32

	addr uint32

	vSyncOn, vSyncOff int
	hSyncOn, hSyncOff int
}

func newCRTC(cpu emulator.CPU) *crtc {
	crtc := &crtc{
		regs:     []byte{63, 40, 46, 0x8E, 38, 0, 25, 30, 0, 7, 0, 0, 0x20, 0, 0, 0},
		cpu:      cpu,
		status:   &crtcStatus{},
		counters: &crtcCounters{},
	}
	crtc.recalcule()
	return crtc
}

func (crtc *crtc) ReadPort(port uint16) (byte, bool) { return 0, false }

func (crtc *crtc) WritePort(port uint16, data byte) {
	f := port >> 8 & 3
	switch f {
	case 0:
		crtc.selectReg = data & 0x0f

	case 1:
		crtc.regs[crtc.selectReg] = data & regMasks[crtc.selectReg]
		fmt.Printf("[crtc] reg: %2d = %d\n", crtc.selectReg, data)

		crtc.recalcule()
	default:
		panic(fmt.Sprintf("[crtc] bad port 0x%04X", port))
	}
}

func (crtc *crtc) recalcule() {
	crtc.addr = ((uint32(crtc.regs[12]) << 8) | uint32(crtc.regs[13]))
	// println("[crtc] addr:", crtc.addr)

	crtc.vSyncOn = int(crtc.regs[7])
	vSyncSize := int((crtc.regs[3] >> 4) & 0x0f)
	crtc.vSyncOff = crtc.vSyncOn + vSyncSize
	// println("[crtc] vSyncOn:", crtc.vSyncOn, "vSyncOff:", crtc.vSyncOff)

	crtc.hSyncOn = int(crtc.regs[2])
	hSyncSize := int(crtc.regs[3] & 0x0f)
	crtc.hSyncOff = crtc.hSyncOn + hSyncSize
	// println("[crtc] hSyncOn:", crtc.hSyncOn, "hSyncOff:", crtc.hSyncOff)

}

func (crtc *crtc) Tick() {
	clock := crtc.cycles / 4
	crtc.cycles++

	if crtc.clock == clock {
		return
	}
	crtc.clock = clock

	R0 := int(crtc.regs[0])
	R1 := int(crtc.regs[1])
	R4 := int(crtc.regs[4])
	R6 := int(crtc.regs[6])
	R9 := int(crtc.regs[9])

	if crtc.counters.h < R0 {
		crtc.counters.h++
	} else {
		crtc.counters.h = 0
		if crtc.counters.sl < R9 {
			crtc.counters.sl++
		} else {
			crtc.counters.sl = 0
			if crtc.counters.row < R4 {
				crtc.counters.row++
			} else {
				crtc.counters.row = 0
			}
		}
		crtc.counters.raster = int(crtc.counters.row*(R9+1)) + int(crtc.counters.sl)
	}

	if crtc.counters.h == crtc.hSyncOff {
		if crtc.counters.raster%52 == 0 {
			crtc.cpu.Interrupt(true)
		}
	}

	crtc.status.vSync = (crtc.counters.row >= crtc.vSyncOn) && (crtc.counters.row <= crtc.vSyncOff)
	crtc.status.hSync = (crtc.counters.h >= crtc.hSyncOn) && (crtc.counters.h <= crtc.hSyncOff)
	crtc.status.disPen = (crtc.counters.h < R1) && (crtc.counters.row < R6)

	MA := crtc.addr + uint32(crtc.counters.row)*uint32(R1) + uint32(crtc.counters.h)
	crtc.status.ma = uint16(((MA & 0x3FF) << 1) | ((uint32(crtc.counters.raster) & 7) << 11) | ((MA & 0x3000) << 2))

	// if crtc.counters.h == 0 {
	// 	fmt.Printf("=> %+v => %+v => 0x%04X\n", crtc.counters, crtc.status, crtc.status.ma)
	// }
}

func (crtc *crtc) FrameEnded() {
	crtc.cycles = 0
}
