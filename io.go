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

type NopIoWriter struct{}
type EmptyIoWriter = NopIoWriter

func NewNopIoWriter() *NopIoWriter {
	return &NopIoWriter{}
}

func NewEmptyIoWriter() *NopIoWriter {
	return &NopIoWriter{}
}

func (w *NopIoWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type NopIoReader struct{}
type EmptyIoReader = NopIoReader

func NewNopIoReader() *NopIoReader {
	return &NopIoReader{}
}

func NewEmptyIoReader() *NopIoReader {
	return &NopIoReader{}
}

func (r *NopIoReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}
