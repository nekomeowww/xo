package xo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDigitOf(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(0, DigitOf(0, 1))
	assert.Equal(5, DigitOf(5, 1))

	assert.Equal(1, DigitOf(11, 1))
	assert.Equal(5, DigitOf(51, 2))

	assert.Equal(8, DigitOf(128, 1))
	assert.Equal(2, DigitOf(128, 2))
	assert.Equal(1, DigitOf(128, 3))
}

func TestParseDecimalStringIntoUint64(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	uintNumber, err := ParseUint64FromDecimalString("0.01", 2)
	require.NoError(err)
	assert.Equal(uint64(1), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("0.11", 2)
	require.NoError(err)
	assert.Equal(uint64(11), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("0.1111", 2)
	require.NoError(err)
	assert.Equal(uint64(11), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("10.00", 2)
	require.NoError(err)
	assert.Equal(uint64(1000), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("10.01", 2)
	require.NoError(err)
	assert.Equal(uint64(1001), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("10.11", 2)
	require.NoError(err)
	assert.Equal(uint64(1011), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("10.0000", 4)
	require.NoError(err)
	assert.Equal(uint64(100000), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("10.01", 4)
	require.NoError(err)
	assert.Equal(uint64(100100), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("10.1101", 4)
	require.NoError(err)
	assert.Equal(uint64(101101), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("1000000.0000001", 2)
	require.NoError(err)
	assert.Equal(uint64(100000000), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("1000000.01", 2)
	require.NoError(err)
	assert.Equal(uint64(100000001), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("0", 2)
	require.NoError(err)
	assert.Equal(uint64(0), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("0000000", 2)
	require.NoError(err)
	assert.Equal(uint64(0), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("0.00000", 2)
	require.NoError(err)
	assert.Equal(uint64(0), uintNumber)

	uintNumber, err = ParseUint64FromDecimalString("00000.00000", 2)
	require.NoError(err)
	assert.Equal(uint64(0), uintNumber)
}

func TestConvertUint64ToDecimalString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("0.00", ConvertUint64ToDecimalString(0, 2))
	assert.Equal("0", ConvertUint64ToDecimalString(0, -1))

	assert.Equal("19.90", ConvertUint64ToDecimalString(1990, 2))
	assert.Equal("19.9", ConvertUint64ToDecimalString(1990, -1))
	assert.Equal("19.99", ConvertUint64ToDecimalString(1999, 2))

	assert.Equal("1.00", ConvertUint64ToDecimalString(100, 2))
	assert.Equal("1", ConvertUint64ToDecimalString(100, -1))
}
