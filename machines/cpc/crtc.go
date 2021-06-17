package cpc

import (
	"fmt"

	"github.com/laullon/b2t80s/cpu/z80"
)

var (
	regMasks = []byte{0xff, 0xff, 0xff, 0xff, 0x7f, 0x1f, 0x7f, 0x7f, 0xff, 0x1f, 0x7f, 0x1f, 0x3f, 0xff, 0x3f, 0xff}
)

type crtcStatus struct {
	vSync  bool
	hSync  bool
	disPen bool
	ma     uint32
	ra     int
}

type crtcCounters struct {
	hcc int // horizontal char counter
	vlc int // vertical line counter
	vcc int // vertical char counter // ra
}

type crtc struct {
	cpu z80.Z80

	status   *crtcStatus
	counters *crtcCounters

	regs      []byte
	selectReg byte

	addr uint32

	vSyncOn, vSyncOff int
	hSyncOn, hSyncOff int
}

func newCRTC(cpu z80.Z80) *crtc {
	crtc := &crtc{
		regs:     []byte{63, 40, 46, 0x8E, 38, 0, 25, 30, 0, 7, 0, 0, 0x20, 0, 0, 0},
		cpu:      cpu,
		status:   &crtcStatus{},
		counters: &crtcCounters{},
	}
	crtc.recalcule()
	return crtc
}

func (crtc *crtc) ReadPort(port uint16) byte { return 0 }

func (crtc *crtc) WritePort(port uint16, data byte) {
	f := port >> 8 & 3
	switch f {
	case 0:
		crtc.selectReg = data & 0x0f

	case 1:
		crtc.regs[crtc.selectReg] = data & regMasks[crtc.selectReg]
		// fmt.Printf("[crtc] reg: %2d = %d\n", crtc.selectReg, data)

		crtc.recalcule()
	default:
		panic(fmt.Sprintf("[crtc] bad port 0x%04X", port))
	}
}

func (crtc *crtc) recalcule() {
	crtc.addr = ((uint32(crtc.regs[12]) << 8) | uint32(crtc.regs[13]))
	// println("[crtc] addr:", crtc.addr)

	crtc.vSyncOn = int(crtc.regs[7])
	vSyncSize := int((crtc.regs[3] >> 4) & 0x0f / 8)
	crtc.vSyncOff = crtc.vSyncOn + vSyncSize
	// println("[crtc] vSyncOn:", crtc.vSyncOn, "vSyncOff:", crtc.vSyncOff)

	crtc.hSyncOn = int(crtc.regs[2])
	hSyncSize := int(crtc.regs[3] & 0x0f)
	crtc.hSyncOff = crtc.hSyncOn + hSyncSize
	// println("[crtc] hSyncOn:", crtc.hSyncOn, "hSyncOff:", crtc.hSyncOff)
}

func (crtc *crtc) Tick() {
	R0 := int(crtc.regs[0])
	R1 := int(crtc.regs[1])
	R4 := int(crtc.regs[4])
	R6 := int(crtc.regs[6])
	R9 := int(crtc.regs[9])

	if crtc.counters.hcc < R0 {
		crtc.counters.hcc++
	} else {
		crtc.counters.hcc = 0
		if crtc.counters.vlc < R9 {
			crtc.counters.vlc++
		} else {
			crtc.counters.vlc = 0
			if crtc.counters.vcc < R4 {
				crtc.counters.vcc++
			} else {
				crtc.counters.vcc = 0
			}
		}
	}

	crtc.status.vSync = (crtc.counters.vcc >= crtc.vSyncOn) && (crtc.counters.vcc <= crtc.vSyncOff)
	crtc.status.hSync = (crtc.counters.hcc >= crtc.hSyncOn) && (crtc.counters.hcc <= crtc.hSyncOff)
	crtc.status.disPen = (crtc.counters.hcc < R1) && (crtc.counters.vcc < R6)

	crtc.status.ma = crtc.addr + uint32(crtc.counters.vcc)*uint32(R1) + uint32(crtc.counters.hcc)
	crtc.status.ra = crtc.counters.vlc

	// if crtc.counters.h == 0 {
	// 	fmt.Printf("=> %+v => %+v => 0x%04X\n", crtc.counters, crtc.status, crtc.status.ma)
	// }
}

func (st *crtcStatus) getAddress() uint16 {
	return uint16(((st.ma & 0x3FF) << 1) | ((uint32(st.ra) & 7) << 11) | ((st.ma & 0x3000) << 2))
}
