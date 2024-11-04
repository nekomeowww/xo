package xo

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeAsRFC1123Name(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Underscore to hyphen", "before_after", "before-after"},
		{"Dots to hyphens", "1.1.1", "1-1-1"},
		{"Trim hyphens", "-middle-", "middle"},
		{"Special chars", "_@invalid.com", "invalid-com"},
		{"Special chars at end", "invalid.com@.", "invalid-com"},
		{"Multiple special chars", "_@invalid.com@.", "invalid-com"},
		{"Only dots", "...", ""},
		{"Non-ASCII chars", "中文中文", ""},
		{"Mixed ASCII and non-ASCII", "中文abcd中文", "abcd"},
		{"Max length 255", strings.Repeat("a", 255), strings.Repeat("a", 255)},
		{"Over max length 256", strings.Repeat("a", 256), strings.Repeat("a", 256)},
		{"Over max length 257", strings.Repeat("a", 257), strings.Repeat("a", 256)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeAsRFC1123Name(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
