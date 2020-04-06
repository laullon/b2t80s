package zx

import "github.com/laullon/b2t80s/emulator"

func (ula *ula) SoundTick() {
	ula.mux.Lock()
	defer ula.mux.Unlock()
	v := 0.0
	if ula.buzzer || ula.ear {
		v = 1
	}
	ula.out = append(ula.out, &emulator.SoundData{L: v, R: v})
}

func (ula *ula) GetBuffer(max int) (res []*emulator.SoundData, l int) {
	ula.mux.Lock()
	defer ula.mux.Unlock()

	if len(ula.out) > max {
		res = ula.out[:max]
		ula.out = ula.out[max:]
		l = max
	} else {
		res = ula.out
		ula.out = nil
		l = len(res)
	}
	return
}
