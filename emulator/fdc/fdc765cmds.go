package fdc

func fdc_specify(fdc *fdc765) {
	fdc.phase = IDLE_PHASE
}

func fdc_seek(fdc *fdc765) {
	fdc.st0 = fdc.cmd.args[CMD_UNIT] & 0b111
	side := 0 //int((fdc.cmd.args[CMD_UNIT] & 0b100) >> 2)
	if fdc.cmd.args[CMD_UNIT] != 0xff {
		fdc.discActive = fdc.cmd.args[CMD_UNIT] & 0b001
	}
	track := int(fdc.cmd.args[CMD_C])
	if fdc.isDriveReady() {
		fdc.discs[fdc.discActive].SeekTrack(side, track)
		fdc.st0 |= 0x20
	} else {
		fdc.st0 |= 0x48
	}
	fdc.trackActive = fdc.cmd.args[2]
	fdc.seekDone = true
	fdc.phase = IDLE_PHASE
}

func fdc_recalib(fdc *fdc765) {
	fdc.cmd.args = append(fdc.cmd.args, fdc.trackActive)
	fdc_seek(fdc)
}

func fdc_intstat(fdc *fdc765) {
	fdc.cmd.result[RES_ST0] = fdc.discActive
	if fdc.seekDone {
		fdc.cmd.result[RES_ST0] |= 0x20 // seek done
		fdc.seekDone = false
		fdc.cmd.result[1] = fdc.trackActive
	} else {
		fdc.cmd.result[RES_ST0] = 0x80
		fdc.cmd.result = fdc.cmd.result[:1]
	}
	fdc.phase = RESULT_PHASE
}

func fdc_readID(fdc *fdc765) {
	fdc.discActive = fdc.cmd.args[CMD_UNIT] & 0b001
	if fdc.isDriveReady() {
		copy(fdc.cmd.result[RES_C:], fdc.discs[fdc.discActive].ActualSector().CHRN)
	} else {
		fdc.cmd.result[RES_ST0] |= 0x40 // AT
		fdc.cmd.result[RES_ST1] |= 0x01 // Missing AM
	}
	fdc.phase = RESULT_PHASE
}

func fdc_read(fdc *fdc765) {
	// fdc.discActive = fdc.cmd.args[CMD_UNIT] & 0b001
	track := fdc.cmd.args[CMD_C]
	side := fdc.cmd.args[CMD_H]
	firstSID := fdc.cmd.args[CMD_R]
	size := fdc.cmd.args[CMD_N]
	lastSID := fdc.cmd.args[CMD_EOT]

	if fdc.isDriveReady() {
		done := false
		sector := fdc.discs[fdc.discActive].SeekSector([]byte{track, side, firstSID, size})
		if sector != nil {
			fdc.cmd.result[RES_ST1] = sector.ST1
			fdc.cmd.result[RES_ST2] = sector.ST2
			if (fdc.cmd.args[CMD_CODE] & 0xf) == 0x0c { // delete
				fdc.cmd.result[RES_ST2] ^= 0x40
			}
			for !done {
				fdc.cmd.data = append(fdc.cmd.data, sector.Data...)
				done = sector.CHRN[2] == lastSID
				if !done {
					sector = fdc.discs[fdc.discActive].NextSector()
				}
			}
			// fdc.discs[fdc.discActive].NextSector()
			fdc.LOAD_RESULT_WITH_CHRN()
			// fdc.LOAD_RESULT_WITH_STATUS()
			fdc.phase = EXEC_PHASE
		}
		if len(fdc.cmd.data) == 0 {
			fdc.cmd.result[RES_ST0] |= 0x40 // AT
			fdc.cmd.result[RES_ST1] |= 0x01 // Missing AM
			fdc.LOAD_RESULT_WITH_CHRN()
			fdc.phase = RESULT_PHASE
		}
	} else {
		fdc.cmd.result[RES_ST0] = fdc.cmd.args[CMD_UNIT]&0b111 | 0x48
		fdc.LOAD_RESULT_WITH_CHRN()
		fdc.phase = RESULT_PHASE
	}
}

func fdc_drvstat(fdc *fdc765) {
	drive := fdc.cmd.args[CMD_UNIT] & 3
	val := drive
	if drive == 0 { // TODO support 2 drives
		val |= 0x48 // set Write Protect & One Sided
		val |= 0x20 // set Ready
		// val |= 0x10 // set Track 0
	} else {
		val |= 0x80
	}
	fdc.cmd.result[RES_ST0] = val
	fdc.phase = RESULT_PHASE // switch to result phase
}
