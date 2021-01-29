package mappers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadNesFile(t *testing.T) {
	f := loadFile("../tests/cpu_interrupts.nes")
	assert.Equal(t, byte(1), f.mapper())
}
