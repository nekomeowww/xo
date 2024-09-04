package otelzap

import (
	"encoding/base64"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LogSeverityFromZapLevel(level zapcore.Level) log.Severity {
	switch level {
	case zapcore.DebugLevel:
		return log.SeverityDebug
	case zapcore.InfoLevel:
		return log.SeverityInfo
	case zapcore.WarnLevel:
		return log.SeverityWarn
	case zapcore.ErrorLevel:
		return log.SeverityError
	case zapcore.DPanicLevel:
		return log.SeverityFatal1
	case zapcore.PanicLevel:
		return log.SeverityFatal2
	case zapcore.FatalLevel:
		return log.SeverityFatal3
	case zapcore.InvalidLevel:
		return log.SeverityUndefined
	default:
		return log.SeverityUndefined
	}
}

func attributeKey(k string) string {
	return "log.fields." + k
}

func AttributesFromZapField(f zap.Field) []attribute.KeyValue {
	switch f.Type {
	case zapcore.BoolType:
		return []attribute.KeyValue{
			attribute.Bool(attributeKey(f.Key), f.Integer == 1),
		}
	case zapcore.Int8Type, zapcore.Int16Type, zapcore.Int32Type, zapcore.Int64Type,
		zapcore.Uint32Type, zapcore.Uint8Type, zapcore.Uint16Type, zapcore.Uint64Type,
		zapcore.UintptrType:
		return []attribute.KeyValue{
			attribute.Int64(attributeKey(f.Key), f.Integer),
		}
	case zapcore.Float32Type, zapcore.Float64Type:
		return []attribute.KeyValue{
			attribute.Float64(attributeKey(f.Key), math.Float64frombits(uint64(f.Integer))),
		}
	case zapcore.Complex64Type, zapcore.Complex128Type:
		return []attribute.KeyValue{
			attribute.String(attributeKey(f.Key), fmt.Sprintf("%v", f.Interface)),
		}
	case zapcore.StringType:
		return []attribute.KeyValue{
			attribute.String(attributeKey(f.Key), f.String),
		}
	case zapcore.StringerType:
		val, ok := f.Interface.(fmt.Stringer)
		if !ok {
			return []attribute.KeyValue{
				attribute.String(attributeKey(f.Key), fmt.Sprintf("expected fmt.Stringer, got %T, v: %v", f.Interface, f.Interface)),
			}
		}
		if lo.IsNil(val) {
			return []attribute.KeyValue{
				attribute.String(attributeKey(f.Key), "<nil>"),
			}
		}

		return []attribute.KeyValue{
			attribute.String(attributeKey(f.Key), fmt.Sprint(val)),
		}
	case zapcore.BinaryType, zapcore.ByteStringType:
		val, ok := f.Interface.([]byte)
		if !ok {
			return []attribute.KeyValue{
				attribute.String(attributeKey(f.Key), fmt.Sprintf("expected []byte, got %T, v: %v", f.Interface, f.Interface)),
			}
		}
		if lo.IsNil(val) {
			return []attribute.KeyValue{
				attribute.String(attributeKey(f.Key), "<empty>"),
			}
		}

		return []attribute.KeyValue{
			attribute.String(attributeKey(f.Key), base64.StdEncoding.EncodeToString(val)),
		}
	case zapcore.DurationType:
		val, ok := f.Interface.(time.Duration)
		if !ok {
			return []attribute.KeyValue{
				attribute.String(attributeKey(f.Key), fmt.Sprintf("expected time.Duration, got %T, v: %v", f.Interface, f.Interface)),
			}
		}

		return []attribute.KeyValue{
			attribute.String(attributeKey(f.Key), val.String()),
		}
	case zapcore.TimeType, zapcore.TimeFullType:
		val, ok := f.Interface.(time.Time)
		if !ok {
			return []attribute.KeyValue{
				attribute.String(attributeKey(f.Key), fmt.Sprintf("expected time.Time, got %T, v: %v", f.Interface, f.Interface)),
			}
		}

		return []attribute.KeyValue{
			attribute.String(attributeKey(f.Key), val.Format(time.RFC3339Nano)),
		}
	case zapcore.ErrorType:
		err, ok := f.Interface.(error)
		if !ok {
			return []attribute.KeyValue{attribute.String(attributeKey(f.Key), fmt.Sprintf("expected error, got %T", f.Interface))}
		}
		if lo.IsNil(err) {
			return []attribute.KeyValue{
				attribute.String(attributeKey(f.Key), "<nil>"),
			}
		}

		return []attribute.KeyValue{
			semconv.ExceptionTypeKey.String(reflect.TypeOf(err).String()),
			semconv.ExceptionMessageKey.String(err.Error()),
		}
	case zapcore.ObjectMarshalerType, zapcore.InlineMarshalerType:
		if marshaler, ok := f.Interface.(zapcore.ObjectMarshaler); ok {
			if lo.IsNil(marshaler) {
				return []attribute.KeyValue{
					attribute.String(attributeKey(f.Key), "<nil>"),
				}
			}

			encoder := zapcore.NewMapObjectEncoder()
			if err := marshaler.MarshalLogObject(encoder); err == nil {
				return flattenMap(attributeKey(f.Key), encoder.Fields)
			}
		}

		return []attribute.KeyValue{
			attribute.String(attributeKey(f.Key), fmt.Sprintf("%v", f.Interface)),
		}
	case zapcore.ArrayMarshalerType:
		if marshaler, ok := f.Interface.(zapcore.ArrayMarshaler); ok {
			if lo.IsNil(marshaler) {
				return []attribute.KeyValue{
					attribute.String(attributeKey(f.Key), "<nil>"),
				}
			}

			encoder := NewBufferArrayEncoder()
			if err := marshaler.MarshalLogArray(encoder); err == nil {
				return []attribute.KeyValue{
					attribute.StringSlice(attributeKey(f.Key), encoder.Result()),
				}
			}
		}

		return []attribute.KeyValue{
			attribute.String(attributeKey(f.Key), fmt.Sprintf("%v", f.Interface)),
		}
	case zapcore.ReflectType:
		str := fmt.Sprint(f.Interface)

		return []attribute.KeyValue{
			attribute.String(attributeKey(f.Key), str),
		}
	case zapcore.NamespaceType:
		return []attribute.KeyValue{}
	case zapcore.SkipType:
		return []attribute.KeyValue{}
	case zapcore.UnknownType:
		return []attribute.KeyValue{
			attribute.String(attributeKey(f.Key), fmt.Sprintf("%v", f.Interface)),
		}
	default:
		return []attribute.KeyValue{
			attribute.String(attributeKey(f.Key), fmt.Sprintf("%v", f.Interface)),
		}
	}
}

func flattenMap(prefix string, m map[string]any) []attribute.KeyValue {
	if m == nil {
		return make([]attribute.KeyValue, 0)
	}

	attrs := make([]attribute.KeyValue, 0, len(m))

	for k, v := range m {
		key := prefix + "." + k
		switch val := v.(type) {
		case bool:
			attrs = append(attrs, attribute.Bool(key, val))
		case int:
			attrs = append(attrs, attribute.Int(key, val))
		case int64:
			attrs = append(attrs, attribute.Int64(key, val))
		case float64:
			attrs = append(attrs, attribute.Float64(key, val))
		case string:
			attrs = append(attrs, attribute.String(key, val))
		case []interface{}:
			attrs = append(attrs, attribute.StringSlice(key, interfaceSliceToStringSlice(val)))
		case map[string]interface{}:
			attrs = append(attrs, flattenMap(key, val)...)
		default:
			attrs = append(attrs, attribute.String(key, fmt.Sprintf("%v", val)))
		}
	}

	return attrs
}

func interfaceSliceToStringSlice(slice []any) []string {
	if slice == nil {
		return make([]string, 0)
	}

	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = fmt.Sprintf("%v", v)
	}

	return result
}
