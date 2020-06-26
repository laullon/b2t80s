package msx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpriteY(t *testing.T) {
	sprt, _ := newSprite([]byte{0, 0, 0, 0})
	assert.Equal(t, 1, sprt.y)

	sprt, _ = newSprite([]byte{255, 0, 0, 0})
	assert.Equal(t, 0, sprt.y)

	sprt, _ = newSprite([]byte{254, 0, 0, 0})
	assert.Equal(t, -1, sprt.y)
}
