package atetris

import (
	"fmt"

	"fyne.io/fyne"
)

func (t *atetris) OnKeyEvent(key *fyne.KeyEvent) {
	fmt.Println("key:", key.Name)
	switch key.Name {

	case fyne.Key1:
		t.pokey1.P0 = !t.pokey1.P0
	case fyne.Key2:
		t.pokey1.P1 = !t.pokey1.P1

	case fyne.KeySpace:
		t.pokey2.P0 = !t.pokey2.P0
	case fyne.KeyDown:
		t.pokey2.P1 = !t.pokey2.P1
	case fyne.KeyRight:
		t.pokey2.P2 = !t.pokey2.P2
	case fyne.KeyLeft:
		t.pokey2.P3 = !t.pokey2.P3

	}
}
