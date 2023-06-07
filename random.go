package xo

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
)

// RandBytes generates bytes according to the given length, defaults to 32.
func RandBytes(length ...int) ([]byte, error) {
	b := make([]byte, 32)
	if len(length) != 0 {
		b = make([]byte, length[0])
	}

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// RandomBase64Token generates the URL-safe Base64 string based on the given byte length,
// the length is 32 by default, the length is the length of the original byte data, not
// the actual length of the Base64 string, the actual length is about 44 by default in
// the case of 32.
func RandomBase64Token(length ...int) (string, error) {
	b, err := RandBytes(length...)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// RandomHashString generates a random SHA256 string with the maximum length of 64.
func RandomHashString(length ...int) string {
	b, _ := RandBytes(1024)
	if len(length) != 0 {
		sliceLength := length[0]
		if length[0] > 64 {
			sliceLength = 64
		}
		if length[0] <= 0 {
			sliceLength = 64
		}

		return fmt.Sprintf("%x", sha256.Sum256(b))[:sliceLength]
	}

	return fmt.Sprintf("%x", sha256.Sum256(b))
}

// RandomInt64 generates a random integer.
func RandomInt64(max ...int64) int64 {
	innerMax := int64(0)
	if len(max) == 0 || (len(max) > 0 && max[0] <= 0) {
		innerMax = 9999999999
	} else {
		innerMax = max[0]
	}

	nBig, _ := rand.Int(rand.Reader, big.NewInt(innerMax))
	n := nBig.Int64()

	return n
}

// RandomInt64InRange generates a random integer in the range.
func RandomInt64InRange(min, max int64) int64 {
	if min >= max {
		panic("min must be less than max")
	}
	if max <= 0 {
		panic("max must be greater than 0")
	}

	nBig, _ := rand.Int(rand.Reader, big.NewInt(max-min))
	n := nBig.Int64()

	return n + min
}

// RandomInt64String generates a random integer string.
func RandomInt64String(digits int64) string {
	max := big.NewInt(0)
	min := big.NewInt(0)

	max.SetString(strings.Repeat("9", int(digits)), 10)
	min.SetString(fmt.Sprintf("%s%s", "1", strings.Repeat("0", int(digits)-1)), 10)

	nBig, _ := rand.Int(rand.Reader, new(big.Int).Sub(max, min))

	return new(big.Int).Add(nBig, min).String()
}
