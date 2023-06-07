package xo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInTestEnvironment(t *testing.T) {
	assert := assert.New(t)

	originalOSArgs := os.Args
	defer func() {
		os.Args = originalOSArgs
	}()

	os.Args = []string{"-test."}

	assert.True(IsInTestEnvironment())

	os.Args = []string{"-v"}

	assert.False(IsInTestEnvironment())
}
