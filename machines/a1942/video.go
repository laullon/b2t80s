package a1942

import (
	"github.com/laullon/b2t80s/cpu"
	"github.com/laullon/b2t80s/gui"
)

// screen_device &set_raw(u32 pixclock, u16 htotal, u16 hbend, u16 hbstart, u16 vtotal, u16 vbend, u16 vbstart)
// m_screen->set_raw(MASTER_CLOCK/2, 384, 128, 0, 262, 22, 246);   // hsync is 50..77, vsync is 257..259

type video struct {
	m *a1942

	spriteram cpu.RAM

	display *gui.Display
	x, y    uint
}

func newVideo(m *a1942) *video {
	v := &video{
		display:   gui.NewDisplay(gui.Size{W: 256, H: 224}),
		m:         m,
		spriteram: cpu.NewRAM(make([]byte, 0x0800), 0x07ff),
	}
	return v
}

func (v *video) Tick() {
	v.x++
	if v.x == 384 {
		v.x = 0
		v.y++
		if v.y == 262 {
			v.y = 0
			v.display.Swap()
		}
		switch v.y {
		case 44:
			// v.m.audioCpu.Interrupt(true)
		case 109:
			v.m.mainCpu.Interrupt(true, 0xcf) /* RST 08h */
			// v.m.audioCpu.Interrupt(true)
		case 175:
			// v.m.audioCpu.Interrupt(true)
		case 240:
			v.m.mainCpu.Interrupt(true, 0xd7) /* RST 10h - vblank */
			// v.m.audioCpu.Interrupt(true)
		}
	}
}

func (v *video) ReadPort(port uint16) (byte, bool) { return 0xff, false }
func (v *video) WritePort(port uint16, data byte) {
	// TODO: c802-c803 background scroll
	// TODO: c805      background palette bank selector

}
