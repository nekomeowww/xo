package otelzap

import (
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"
)

// BufferArrayEncoder implements zapcore.BufferArrayEncoder.
// It represents all added objects to their string values and
// adds them to the stringsSlice buffer.
type BufferArrayEncoder struct {
	stringsSlice []string
}

var _ zapcore.ArrayEncoder = (*BufferArrayEncoder)(nil)

func NewBufferArrayEncoder() *BufferArrayEncoder {
	return &BufferArrayEncoder{
		stringsSlice: make([]string, 0),
	}
}

func (t *BufferArrayEncoder) Result() []string {
	return t.stringsSlice
}

func (t *BufferArrayEncoder) AppendComplex128(v complex128) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendComplex64(v complex64) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendArray(v zapcore.ArrayMarshaler) error {
	enc := &BufferArrayEncoder{}
	err := v.MarshalLogArray(enc)
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", enc.stringsSlice))

	return err
}

func (t *BufferArrayEncoder) AppendObject(v zapcore.ObjectMarshaler) error {
	m := zapcore.NewMapObjectEncoder()
	err := v.MarshalLogObject(m)
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", m.Fields))

	return err
}

func (t *BufferArrayEncoder) AppendReflected(v interface{}) error {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
	return nil
}

func (t *BufferArrayEncoder) AppendBool(v bool) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendByteString(v []byte) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendDuration(v time.Duration) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendFloat64(v float64) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendFloat32(v float32) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendInt(v int) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendInt64(v int64) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendInt32(v int32) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendInt16(v int16) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendInt8(v int8) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendString(v string) {
	t.stringsSlice = append(t.stringsSlice, v)
}

func (t *BufferArrayEncoder) AppendTime(v time.Time) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendUint(v uint) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendUint64(v uint64) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendUint32(v uint32) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendUint16(v uint16) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendUint8(v uint8) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}

func (t *BufferArrayEncoder) AppendUintptr(v uintptr) {
	t.stringsSlice = append(t.stringsSlice, fmt.Sprintf("%v", v))
}
