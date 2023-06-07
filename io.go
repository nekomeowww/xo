package xo

import (
	"bytes"
	"io"
	"os"
)

// ReadFileAsBytesBuffer reads a file and returns a bytes.Buffer.
func ReadFileAsBytesBuffer(path string) (*bytes.Buffer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, file)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
