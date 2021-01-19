package atetris

import (
	"testing"

	"github.com/laullon/b2t80s/cpu/m6502"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tetris := NewATetris()

	if testing.Short() {
		println("skipping logs in short mode.")
	} else {
		tetris.(*atetris).cpu.SetDebuger(m6502.NewDebugger(tetris.(*atetris).cpu, nil, tetris.(*atetris).clock))
	}

	tetris.Clock().RunFor(20)
	assert.FailNow(t, "xxx")
}
