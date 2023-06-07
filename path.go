package xo

import (
	"os"
	"path/filepath"
	"runtime"
)

func RelativePathOf(fp string) string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	callerDir := filepath.Dir(filepath.FromSlash(file))

	return filepath.FromSlash(filepath.Join(callerDir, fp))
}

func RelativePathBasedOnPwdOf(fp string) string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	return filepath.FromSlash(filepath.Join(dir, fp))
}
