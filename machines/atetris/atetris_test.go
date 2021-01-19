package atetris

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tetris := NewATetris()
	for i := 0; i < 20; i++ {
		tetris.(*atetris).cpu.Tick()
	}
	assert.FailNow(t, "xxx")
}
