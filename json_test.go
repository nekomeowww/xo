package xo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsJSON(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsJSON(`{"foo": "bar"}`))
	assert.True(IsJSON(`["foo", "bar"]`))
	assert.False(IsJSON(`{"foo": "bar"`))
	assert.False(IsJSON(``))
}

func TestIsJSONBytes(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsJSONBytes([]byte(`{"foo": "bar"}`)))
	assert.True(IsJSONBytes([]byte(`["foo", "bar"]`)))
	assert.False(IsJSONBytes([]byte(`{"foo": "bar"`)))
	assert.False(IsJSONBytes([]byte(``)))
}
