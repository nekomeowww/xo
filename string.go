package xo

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

// IsStringPrintable determines whether a string is printable.
func IsStringPrintable(str string) bool {
	for _, v := range str {
		if v == '\n' || v == '\r' || v == '\t' {
			continue
		}
		if !unicode.IsGraphic(v) {
			return false
		}
	}

	return true
}

// IsASCIIPrintable determines whether a string is printable ASCII.
func IsASCIIPrintable(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return false
		}
	}

	return true
}

// IsValidUUID determines whether a string is a valid UUID.
func IsValidUUID(uuidStr string) bool {
	if _, err := uuid.Parse(uuidStr); err != nil {
		return false
	}

	return true
}

// ContainsCJKChar determines whether a string contains CJK characters.
func ContainsCJKChar(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
		if unicode.Is(unicode.Hangul, r) {
			return true
		}
		if unicode.Is(unicode.Hiragana, r) {
			return true
		}
		if unicode.Is(unicode.Katakana, r) {
			return true
		}

		/*
			U+3001  、
			U+3002  。
			U+3003  〃
			U+3008  〈
			U+3009  〉
			U+300A  《
			U+300B  》
			U+300C  「
			U+300D  」
			U+300E  『
			U+300F  』
			U+3010  【
			U+3011  】
			U+3014  〔
			U+3015  〕
			U+3016  〖
			U+3017  〗
			U+3018  〘
			U+3019  〙
			U+301A  〚
			U+301B  〛
			U+301C  〜
			U+301D  〝
			U+301E  〞
			U+301F  〟
			U+3030  〰
			U+303D  〽
		*/
		if r >= 0x3001 && r <= 0x303D {
			return true
		}
	}

	return false
}

func Stringify(v any) string {
	if v == nil {
		return ""
	}
	if lo.IsNil(v) {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case fmt.Stringer:
		if lo.IsNil(val) {
			return ""
		} else if val == nil {
			return ""
		} else {
			return val.String()
		}
	case int:
		return strconv.FormatInt(int64(val), 10)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case complex64:
		return strconv.FormatComplex(complex128(val), 'f', -1, 64)
	case complex128:
		return strconv.FormatComplex(val, 'f', -1, 128)
	case bool:
		return strconv.FormatBool(val)
	case []byte:
		return string(val)
	case []rune:
		return string(val)
	case strings.Builder:
		return val.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

var (
	errFailedToConvertStringToType = func(t any, err error) error { return fmt.Errorf("failed to convert string to type %T: %w", t, err) }
)

func FromString[T any](str string) (T, error) {
	var empty T
	if str == "" {
		switch any(empty).(type) {
		case []byte:
			val, _ := any(make([]byte, 0)).(T)
			return val, nil
		case []rune:
			val, _ := any(make([]rune, 0)).(T)
			return val, nil
		case *strings.Builder:
			val, _ := any(&strings.Builder{}).(T)
			return val, nil
		}

		return empty, nil
	}

	switch any(empty).(type) {
	case string:
		val, _ := any(str).(T)
		return val, nil
	case int:
		val, err := strconv.ParseInt(str, 10, 0)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(int(val)).(T)

		return typeVal, nil
	case int8:
		val, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(int8(val)).(T)

		return typeVal, nil
	case int16:
		val, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(int16(val)).(T)

		return typeVal, nil
	case int32:
		val, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(int32(val)).(T)

		return typeVal, nil
	case int64:
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(val).(T)

		return typeVal, nil
	case uint:
		val, err := strconv.ParseUint(str, 10, 0)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(uint(val)).(T)

		return typeVal, nil
	case uint8:
		val, err := strconv.ParseUint(str, 10, 8)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(uint8(val)).(T)

		return typeVal, nil
	case uint16:
		val, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(uint16(val)).(T)

		return typeVal, nil
	case uint32:
		val, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(uint32(val)).(T)

		return typeVal, nil
	case uint64:
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(val).(T)

		return typeVal, nil
	case float32:
		val, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(float32(val)).(T)

		return typeVal, nil
	case float64:
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(val).(T)

		return typeVal, nil
	case complex64:
		val, err := strconv.ParseComplex(str, 64)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(complex64(val)).(T)

		return typeVal, nil
	case complex128:
		val, err := strconv.ParseComplex(str, 128)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(val).(T)

		return typeVal, nil
	case bool:
		val, err := strconv.ParseBool(str)
		if err != nil {
			return empty, errFailedToConvertStringToType(empty, err)
		}

		typeVal, _ := any(val).(T)

		return typeVal, nil
	case []byte:
		val, _ := any([]byte(str)).(T)
		return val, nil
	case []rune:
		val, _ := any([]rune(str)).(T)
		return val, nil
	case *strings.Builder:
		var sb strings.Builder

		sb.WriteString(str)
		val, _ := any(&sb).(T)

		return val, nil
	default:
		return empty, fmt.Errorf("unsupported type %T", empty)
	}
}

func FromStringOrEmpty[T any](str string) T {
	var empty T

	val, err := FromString[T](str)
	if err != nil {
		return empty
	}

	return val
}

func Substring(str string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = len(str)
	}
	if start > len(str) {
		start = len(str)
	}
	if end > len(str) {
		end = len(str)
	}
	if end < start {
		start, end = end, start
	}

	return str[start:end]
}
