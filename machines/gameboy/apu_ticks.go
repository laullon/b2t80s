package gameboy

func (apu *apu) Tick() {
	apu.sequencerCount++
	apu.sequencerCount &= 7
	switch apu.sequencerCount {
	case 0, 4:
		for _, ch := range apu.channels {
			ch.lengthTick()
		}
	case 7:
		for _, ch := range apu.channels {
			ch.tickEnvelope()
		}
	case 2, 6:
		for _, ch := range apu.channels {
			ch.lengthTick()
			ch.sweepTick()
		}
	}
}

func (ch *tomeChannel) sweepTick() {
	if ch.sweepPeriodr > 0 {
		ch.sweepPeriodr--
	}

	if ch.sweepPeriod > 0 {
		if ch.sweepEnable && ch.sweepPeriodr == 0 {
			newFrequency := ch.calculateSweep()
			if newFrequency <= 2047 && ch.sweepShift > 0 {
				ch.frequency = newFrequency
				ch.calculateSweep()
			}
			ch.sweepPeriodr = ch.sweepPeriod
		}
		// println("[sweepTick] ch.frequency:", ch.frequency, "ch.frequency11b:", ch.frequency11b, "ch.sweepPeriodr:", ch.sweepPeriodr, "ch.sweepShift:", ch.sweepShift, ch.enable)
	}
}

func (ch *basicChannel) lengthTick() {
	if ch.soundLengthEnable && ch.soundLength != 0 {
		ch.soundLength--
		if ch.soundLength == 0 {
			ch.enable = false
		}
		// println("[lengthTick] ch.soundLength:", ch.soundLength, ch.enable)
	}
}

func (ch *basicChannel) sweepTick() {}

// ***********************
// ** Envelope ***********
// ***********************

func (envelope *envelope) set(data byte) {
	envelope.initialVolume = data >> 4
	envelope.inc = data&0b1000 != 0
	envelope.period = data & 7
	envelope.reset()
}

func (envelope *envelope) reset() {
	envelope.volume = envelope.initialVolume
	envelope.timer = envelope.period
}

func (ch *tomeChannel) tickEnvelope()  { ch.envelope.tick(ch) }
func (ch *noiseChannel) tickEnvelope() { ch.envelope.tick(ch) }
func (ch *waveChannel) tickEnvelope()  {}

func (envelope *envelope) tick(ch channel) {
	if envelope.period > 0 && envelope.volume != 0 {
		if envelope.timer > 0 {
			envelope.timer -= 1
		}
		if envelope.timer == 0 {
			if envelope.inc {
				if envelope.volume < 0xF {
					envelope.volume += 1
				}
			} else {
				if envelope.volume > 0 {
					envelope.volume -= 1
				}
			}
		}
		println("[envelope] volume:", envelope.volume, ",timer:", envelope.timer)
	}
}
