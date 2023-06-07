package xo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDecimalsPlacesValid(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsDecimalsPlacesValid(1, 0))
	assert.True(IsDecimalsPlacesValid(1.0, 0))
	assert.True(IsDecimalsPlacesValid(1.0, 1))
	assert.True(IsDecimalsPlacesValid(1.0, 2))

	assert.True(IsDecimalsPlacesValid(10, 1))
	assert.True(IsDecimalsPlacesValid(10, 2))
	assert.True(IsDecimalsPlacesValid(10, 3))
	assert.True(IsDecimalsPlacesValid(10, 4))
	assert.True(IsDecimalsPlacesValid(10, 100))

	assert.True(IsDecimalsPlacesValid(1.1, 1))
	assert.True(IsDecimalsPlacesValid(1.11, 2))
	assert.True(IsDecimalsPlacesValid(1.111, 3))
	assert.True(IsDecimalsPlacesValid(1.111, 4))
	assert.True(IsDecimalsPlacesValid(1.111, 100))

	assert.False(IsDecimalsPlacesValid(1.1, 0))
	assert.False(IsDecimalsPlacesValid(1.11, 1))
	assert.False(IsDecimalsPlacesValid(1.111, 2))
}
