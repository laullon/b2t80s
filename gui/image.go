package gui

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Image interface {
	GUIObject
}

type imgViewer struct {
	display *Display

	rect Rect

	texture uint32
	frameID uint32
}

func NewDisplayViewer(display *Display) Image {
	i := &imgViewer{
		display: display,
	}
	i.init()
	return i
}

func (i *imgViewer) Resize(r Rect) {
	i.rect = r
}

func (_ *imgViewer) GetMouseTargets() []MouseTarget {
	return []MouseTarget{}
}

func (i *imgViewer) init() {
	if i.texture != 0 {
		gl.DeleteTextures(1, &i.texture)
	}
	// Texture
	gl.GenTextures(1, &i.texture)
	gl.BindTexture(gl.TEXTURE_2D, i.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB,
		i.display.size.W, i.display.size.H,
		0, gl.RGB, gl.UNSIGNED_BYTE,
		nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// BUFFER
	if i.frameID != 0 {
		gl.DeleteFramebuffers(1, &i.frameID)
	}
	gl.GenFramebuffers(1, &i.frameID)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, i.frameID)
	gl.FramebufferTexture2D(gl.READ_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, i.texture, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
}

func (i *imgViewer) Render() {
	gl.BindTexture(gl.TEXTURE_2D, i.texture)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0,
		0, 0, i.display.size.W, i.display.size.H,
		gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(i.display.Pix()))

	w, h := i.rect.W, i.rect.H

	ratioOrg := float64(i.display.ViewSize.W) / float64(i.display.ViewSize.H)
	ratioDst := float64(w) / float64(h)

	var newW, newH int32
	var offX, offY int32
	if ratioDst > ratioOrg {
		// (wi * hs/hi, hs)
		newW = int32(float64(i.display.ViewSize.W) * float64(h) / float64(i.display.ViewSize.H))
		newH = int32(h)
		offX = (int32(w) - newW) / 2
	} else {
		// hi * ws/wi
		newW = int32(w)
		newH = int32(float64(i.display.ViewSize.H) * float64(w) / float64(i.display.ViewSize.W))
		offY = (int32(h) - newH) / 2
	}

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, i.frameID)
	gl.BlitFramebuffer(
		// SRC
		i.display.ViewPortRect.X, i.display.ViewPortRect.Y,
		i.display.ViewPortRect.X+i.display.ViewPortRect.W, i.display.ViewPortRect.Y+i.display.ViewPortRect.H,
		// DST
		i.rect.X+offX, i.rect.Y+offY,
		i.rect.X+int32(newW)+offX, i.rect.Y+int32(newH)+offY,
		// MODE
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)

	// i.ui.Render()

	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	gl.Flush()
}
