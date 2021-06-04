package gameboy

// ****************************
// ****************************
// ****************************

func (ch *tomeChannel) setRegister(r int, data byte) {
	switch r {
	case 0:
		if !ch.sweepOFF {
			ch.sweepPeriod = (data >> 4) & 7
			ch.sweepDec = data&0b1000 != 0
			ch.sweepShift = data & 7
		}

	case 1:
		ch.waveDuty = data >> 6
		ch.soundLength = 64 - uint16(data)&63

	case 2:
		ch.envelope.set(data)
		if ch.envelope.initialVolume != 0 {
			ch.enable = false
		}

	case 3:
		ch.frequency11b = ch.frequency11b&0xff00 | uint16(data)

	case 4:
		ch.frequency11b = ch.frequency11b&0x00ff | uint16(data&0b111)<<8
		ch.trigger = data&0x80 != 0
		ch.soundLengthEnable = data&0x40 != 0

		if ch.trigger {
			ch.enable = true
			ch.frequency = ch.frequency11b

			ch.envelope.reset()

			ch.sweepPeriodr = ch.sweepPeriod
			ch.sweepEnable = ch.sweepShift != 0 || ch.sweepPeriod != 0
			if ch.sweepShift != 0 {
				ch.calculateSweep()
			}

			if ch.soundLength == 0 {
				ch.soundLength = 64
			}

			if !ch.soundLengthEnable && ch.envelope.initialVolume == 0 {
				ch.enable = false
			}
		}
	}
}

func (ch *tomeChannel) getRegister(r int) (res byte) {
	switch r {
	case 0:
		if !ch.sweepOFF {
			res = 0x80
			res |= ch.sweepPeriod & 7 << 4
			res |= ch.sweepShift & 7
			if ch.sweepDec {
				res |= 0b1000
			}
		} else {
			res = 0xff
		}

	case 1:
		res = 0x3f
		res |= ch.waveDuty << 6

	case 2:
		res = ch.envelope.initialVolume << 4
		if ch.envelope.inc {
			res |= 0b1000
		}
		res |= ch.envelope.period

	case 4:
		res = 0xbf
		if ch.soundLengthEnable {
			res |= 0x40
		}

	case 3, 5:
		res = 0xff
	}

	return
}

// ****************************
// ****************************
// ****************************

func (ch *waveChannel) setRegister(r int, data byte) {
	switch r {
	case 0:
		ch.dac = data&0x80 != 0

	case 1:
		ch.soundLength = 256 - uint16(data)

	case 2:
		ch.volume = data >> 5

	case 3:
		ch.frequency11b = ch.frequency11b&0xff00 | uint16(data)

	case 4:
		ch.frequency11b = ch.frequency11b&0x00ff | uint16(data&0b111)<<8
		ch.trigger = data&0x80 != 0
		ch.soundLengthEnable = data&0x40 != 0
		if ch.trigger && ch.soundLength == 0 {
			ch.soundLength = 0x100
		}
	}
}

func (ch *waveChannel) getRegister(r int) (res byte) {
	switch r {
	case 0:
		res = 0x7f
		if ch.dac {
			res |= 0x80
		}

	case 1:
		res = 0xff

	case 2:
		res = 0x9F
		res |= ch.volume << 5

	case 3:
		res = 0xff

	case 4:
		res = 0xbf
		if ch.soundLengthEnable {
			res |= 0x40
		}
	}
	return
}
