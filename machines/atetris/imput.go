package atetris

import "github.com/go-gl/glfw/v3.3/glfw"

func (t *atetris) OnKey(key glfw.Key) {
	// fmt.Println("key:", key.Name)
	switch key {

	case glfw.Key1:
		t.pokey1.P0 = !t.pokey1.P0
	case glfw.Key2:
		t.pokey1.P1 = !t.pokey1.P1

	case glfw.KeySpace:
		t.pokey2.P0 = !t.pokey2.P0
	case glfw.KeyDown:
		t.pokey2.P1 = !t.pokey2.P1
	case glfw.KeyRight:
		t.pokey2.P2 = !t.pokey2.P2
	case glfw.KeyLeft:
		t.pokey2.P3 = !t.pokey2.P3

	}
}
