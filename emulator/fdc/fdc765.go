package fdc

import (
	"fmt"

	"github.com/laullon/b2t80s/emulator/files"
)

const (
	IDLE_PHASE   = 0
	CMD_PHASE    = 1
	EXEC_PHASE   = 2
	RESULT_PHASE = 3

	FDC_TO_CPU = 0
	CPU_TO_FDC = 1

	SKIP_flag       byte = 1   // skip sectors with DDAM/DAM
	RNDDE_flag      byte = 8   // simulate random DE sectors
	OVERRUN_flag    byte = 16  // data transfer timed out
	SCAN_flag       byte = 32  // one of the three scan commands is active
	SCANFAILED_flag byte = 64  // memory and sector data does not match
	STATUSDRVA_flag byte = 128 // status change of drive A
	// STATUSDRVB_flag byte = 256 // status change of drive B

	CMD_CODE = 0
	CMD_UNIT = 1
	CMD_C    = 2
	CMD_H    = 3
	CMD_R    = 4
	CMD_N    = 5
	CMD_EOT  = 6
	CMD_GPL  = 7
	CMD_DTL  = 8
	CMD_STP  = 8

	RES_ST0 = 0
	RES_ST1 = 1
	RES_ST2 = 2
	RES_C   = 3
	RES_H   = 4
	RES_R   = 5
	RES_N   = 6

	OVERRUN_TIMEOUT = (128 * 4)
	INITIAL_TIMEOUT = (OVERRUN_TIMEOUT * 4)
)

type fdc765 struct {
	phase   byte
	cmd     *fdc765cmd
	timeout int
	fase    byte

	st0 byte

	motor bool
	led   bool

	discActive  byte
	trackActive byte
	discs       []files.DSK

	seekDone bool
}

func New765() FDC {
	return &fdc765{
		discs: make([]files.DSK, 4),
	}
}

func (fdc *fdc765) SetDiscA(dsk files.DSK) { fdc.discs[0] = dsk }
func (fdc *fdc765) SetDiscB(dsk files.DSK) { fdc.discs[1] = dsk }

var fdc_cmd_table = map[byte]*fdc765cmd{
	0x03: newFDC765cmd(3, 0, FDC_TO_CPU, fdc_specify), // specify
	0x04: newFDC765cmd(2, 1, FDC_TO_CPU, fdc_drvstat), // sense device status
	0x06: newFDC765cmd(9, 7, FDC_TO_CPU, fdc_read),    // read data
	0x07: newFDC765cmd(2, 0, FDC_TO_CPU, fdc_recalib), // recalibrate
	0x08: newFDC765cmd(1, 2, FDC_TO_CPU, fdc_intstat), // sense interrupt status
	0x0A: newFDC765cmd(2, 7, FDC_TO_CPU, fdc_readID),  // read id
	0x0f: newFDC765cmd(3, 0, FDC_TO_CPU, fdc_seek),    // seek

	// 0x42: newFDC765cmd(9, 7, FDC_TO_CPU, fdc_readtrk), // read diagnostic
	// 0x45: newFDC765cmd(9, 7, CPU_TO_FDC, fdc_write),   // write data
	// 0x49: newFDC765cmd(9, 7, CPU_TO_FDC, fdc_write),   // write deleted data
	0x0c: newFDC765cmd(9, 7, FDC_TO_CPU, fdc_read), // read deleted data
	// 0x4d: newFDC765cmd(6, 7, CPU_TO_FDC, fdc_writeID), // write id
	// 0x51: newFDC765cmd(9, 7, CPU_TO_FDC, fdc_scan)    // scan equal
	// 0x59: newFDC765cmd(9, 7, CPU_TO_FDC, fdc_scan)    // scan low or equal
	// 0x5d: newFDC765cmd(9, 7, CPU_TO_FDC, fdc_scan)    // scan high or equal
}

func (fdc *fdc765) WriteData(val byte) {
	// fmt.Printf("[fdc][write] phase:%d val:0x%02x(%d)\n", fdc.phase, val, val)
	switch fdc.phase {
	case IDLE_PHASE:
		cmdID := val & 0x1f
		if cmd, ok := fdc_cmd_table[cmdID]; ok {
			cmd.init()
			fdc.cmd = cmd
			fdc.phase = CMD_PHASE
			fdc.cmd.args = append(fdc.cmd.args, val)
			if len(fdc.cmd.args) == fdc.cmd.length {
				fdc.led = true
				// println("[fdc] ==>", fdc.cmd.String())
				fdc.cmd.handler(fdc)
				// println("[fdc] <==", fdc.cmd.String())
			} else {
				fdc.phase = CMD_PHASE
			}
		} else {
			panic(fmt.Sprintf("new command: 0x%02X", cmdID))
		}
	case CMD_PHASE:
		if fdc.cmd == nil {
			panic("---")
		}
		fdc.cmd.args = append(fdc.cmd.args, val)
		if len(fdc.cmd.args) == fdc.cmd.length {
			fdc.led = true
			// println("[fdc] ==>", fdc.cmd.String())
			fdc.cmd.handler(fdc)
			// println("[fdc] <==", fdc.cmd.String())
		}
	default:
		panic(fdc.phase)
	}
}

func (fdc *fdc765) LOAD_RESULT_WITH_STATUS() {
	fdc.cmd.result[RES_ST0] |= 0x40     /* AT */
	fdc.cmd.result[RES_ST1] |= 0x80     /* End of Cylinder */
	if fdc.cmd.args[CMD_CODE] != 0x42 { /* continue only if not a read track command */
		if (fdc.cmd.result[RES_ST1]&0x7f) != 0 || (fdc.cmd.result[RES_ST2]&0x7f) != 0 { /* any 'error bits' set? */
			fdc.cmd.result[RES_ST1] &= 0x7f                                                 /* mask out End of Cylinder */
			if (fdc.cmd.result[RES_ST1]&0x20) != 0 || (fdc.cmd.result[RES_ST2]&0x20) != 0 { /* DE and/or DD? */
				fdc.cmd.result[RES_ST2] &= 0xbf /* mask out Control Mark */
			} else if fdc.cmd.result[RES_ST2]&0x40 != 0 { /* Control Mark? */
				fdc.cmd.result[RES_ST0] &= 0x3f /* mask out AT */
				fdc.cmd.result[RES_ST1] &= 0x7f /* mask out End of Cylinder */
			}
		}
	}
}

func (fdc *fdc765) LOAD_RESULT_WITH_CHRN() {
	fdc.cmd.result[RES_C] = fdc.cmd.args[CMD_C] /* load result with current CHRN values */
	fdc.cmd.result[RES_H] = fdc.cmd.args[CMD_H]
	fdc.cmd.result[RES_R] = fdc.cmd.args[CMD_R]
	fdc.cmd.result[RES_N] = fdc.cmd.args[CMD_N]
}

func (fdc *fdc765) isDriveReady() bool {
	val := fdc.cmd.args[CMD_UNIT] & 0b111
	// println("isDriveReady", fdc.discActive, fdc.motor)
	if fdc.discs[fdc.discActive] == nil || (!fdc.motor) {
		val |= 0x48 // Abnormal Termination + Not Ready
	}
	if len(fdc.cmd.result) > 0 {
		fdc.cmd.result[RES_ST0] = val
	}
	return (val & 8) == 0
}

func (fdc *fdc765) ReadData() byte {
	val := byte(0xff) // default value
	if fdc.cmd != nil {
		switch fdc.phase {
		case EXEC_PHASE: // in execution phase?
			if fdc.cmd.direction == FDC_TO_CPU {
				// fmt.Printf("[fdc] <== readData EXEC '%v' \n", fdc.cmd)
				fdc.timeout = OVERRUN_TIMEOUT
				val = fdc.cmd.data[0]
				fdc.cmd.data = fdc.cmd.data[1:]
				if len(fdc.cmd.data) == 0 {
					fdc.LOAD_RESULT_WITH_STATUS()
					fdc.LOAD_RESULT_WITH_CHRN()
					fdc.phase = RESULT_PHASE // switch to result phase
				}
			}

		case RESULT_PHASE: // in result phase?
			// fmt.Printf("[fdc] <== readData RESULT '%v' \n", fdc.cmd)
			val = fdc.cmd.result[0]
			fdc.cmd.result = fdc.cmd.result[1:]
			if len(fdc.cmd.result) == 0 {
				fdc.phase = IDLE_PHASE // switch to command phase
				fdc.led = false        // turn the drive LED off
				fdc.cmd = nil
			}
		}
	}
	return val
}

func (fdc *fdc765) ReadStatus() byte {
	val := byte(0x80) // data register ready

	if fdc.phase != IDLE_PHASE {
		val |= 0x10 // FDC is busy
	}

	switch fdc.phase {
	case EXEC_PHASE:
		val |= 0x20 // FDC is executing & busy
		if fdc.cmd.direction == FDC_TO_CPU {
			val |= 0x40 // FDC is sending data to the CPU
		}

	case RESULT_PHASE:
		val |= 0x40 // FDC is sending data to the CPU, and is busy
	}

	// fmt.Printf("[fdc] <== readStatus 0b%08b phase:%d\n", val, fdc.phase)
	return val
}

func (fdc *fdc765) SetMotor(on bool) {
	fdc.motor = on
}
