package xo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsASCIIPrintable(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		assert := assert.New(t)

		assert.True(IsASCIIPrintable("abcd1234!?@#$%^&*()[]{}<>|\\/\"'`~,."))
		assert.False(IsASCIIPrintable("abc😊"))
		assert.False(IsASCIIPrintable("😊abc"))
		assert.False(IsASCIIPrintable("abc中文"))
		assert.False(IsASCIIPrintable("abc\n"))
	})
	t.Run("Empty still returns True", func(t *testing.T) {
		assert := assert.New(t)

		assert.True(IsASCIIPrintable(""))
		assert.True(IsASCIIPrintable(" "))
		assert.True(IsASCIIPrintable("abc  f   k"))
	})
}

func TestIsValidUUID(t *testing.T) {
	assert := assert.New(t)

	strOk := "93d3ea4c-c66b-47ac-8472-747a24ecc86b"
	strErr := "93d3ea4c-c66b-47ac-8472-747a24ecc86"
	strErr2 := "93d3ea4c-"

	assert.True(IsValidUUID(strOk))
	assert.False(IsValidUUID(strErr))
	assert.False(IsValidUUID(strErr2))
}

func TestSubstring(t *testing.T) {
	abc := Substring("abc", 0, 0)
	assert.Equal(t, "", abc)

	abc = Substring("abc", 0, 1)
	assert.Equal(t, "a", abc)

	abc = Substring("abc", 0, 2)
	assert.Equal(t, "ab", abc)

	abc = Substring("abc", 0, 3)
	assert.Equal(t, "abc", abc)

	abc = Substring("abc", 0, 4)
	assert.Equal(t, "abc", abc)
}
