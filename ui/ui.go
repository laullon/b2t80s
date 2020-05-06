package ui

import "fyne.io/fyne"

type Control interface {
	Widget() fyne.CanvasObject
}
