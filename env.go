//go:build !release

package xo

import (
	"os"
	"strings"
)

// IsInTestEnvironment determines whether the current environment is a test environment.
func IsInTestEnvironment() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}

	return false
}
