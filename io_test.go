package xo

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNopIoWriter_Writer(t *testing.T) {
	buffer := new(bytes.Buffer)

	n, err := fmt.Fprint(buffer, "abcd")
	require.NoError(t, err)

	n2, err2 := fmt.Fprint(&NopIoWriter{}, "abcd")
	require.NoError(t, err2)

	assert.Equal(t, n, n2)
}

func TestNopIoReader_Read(t *testing.T) {
	buffer := new(bytes.Buffer)
	buffer.WriteString("abcd")

	content, err := io.ReadAll(buffer)
	require.NoError(t, err)
	require.NotEmpty(t, content)
	assert.Len(t, content, 4)

	content2, err2 := io.ReadAll(&NopIoReader{})
	require.NoError(t, err2)
	require.Empty(t, content2)
}
