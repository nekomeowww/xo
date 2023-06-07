package xo

import (
	"math"
	"math/big"
	"strconv"

	"github.com/shopspring/decimal"
)

// DigitOf returns the n-th digit of val.
//
// Underlying calculation is based on (val % 10^n) .
func DigitOf[I uint | uint16 | uint32 | uint64 | int | int16 | int32 | int64](val I, n int) I {
	if n < 0 {
		return 0
	}

	r := val % I(math.Pow10(n))

	return r / I(math.Pow10(n-1))
}

// ParseUint64FromDecimalString converts a string with a decimal point to a uint64.
func ParseUint64FromDecimalString(decimalStr string, percision int) (uint64, error) {
	priceDecimal, err := decimal.NewFromString(decimalStr)
	if err != nil {
		return 0, err
	}

	priceDecimal = priceDecimal.Mul(decimal.NewFromFloat(math.Pow10(percision)))

	priceDecimalBigInt := priceDecimal.BigInt()
	if priceDecimalBigInt.Cmp(big.NewInt(0)) <= 0 {
		return 0, nil
	}

	return priceDecimalBigInt.Uint64(), nil
}

// ConvertUint64ToDecimalString formats a uint64 to a string with a decimal point.
func ConvertUint64ToDecimalString(amount uint64, prec int) string {
	if amount == 0 {
		strconv.FormatFloat(0, 'f', prec, 64)
	}

	float64Number, _ := decimal.
		NewFromInt(int64(amount)).
		Div(decimal.NewFromInt(100)).
		Float64()

	return strconv.FormatFloat(float64Number, 'f', prec, 64)
}
