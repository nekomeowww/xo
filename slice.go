package xo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

// ToMap converts a slice to a map with key from key getter func and pairs with value.
func ToMap[T any, K comparable](t []T, keyGetter func(T) K) map[K]T {
	grouped := lo.GroupBy(t, keyGetter)

	return lo.MapValues(grouped, func(values []T, key K) T {
		return values[0]
	})
}

// Join returns a string contains items joined with sep by using fmt.Sprintf.
func Join[T any](from []T, sep string) string {
	strs := make([]string, len(from))

	for i, v := range from {
		switch val := any(v).(type) {
		case string:
			strs[i] = val
		case fmt.Stringer:
			strs[i] = val.String()
		case int:
			strs[i] = strconv.FormatInt(int64(val), 10)
		case int8:
			strs[i] = strconv.FormatInt(int64(val), 10)
		case int16:
			strs[i] = strconv.FormatInt(int64(val), 10)
		case int32:
			strs[i] = strconv.FormatInt(int64(val), 10)
		case int64:
			strs[i] = strconv.FormatInt(val, 10)
		case uint:
			strs[i] = strconv.FormatUint(uint64(val), 10)
		case uint8:
			strs[i] = strconv.FormatUint(uint64(val), 10)
		case uint16:
			strs[i] = strconv.FormatUint(uint64(val), 10)
		case uint32:
			strs[i] = strconv.FormatUint(uint64(val), 10)
		case uint64:
			strs[i] = strconv.FormatUint(val, 10)
		case float32:
			strs[i] = strconv.FormatFloat(float64(val), 'f', -1, 32)
		case float64:
			strs[i] = strconv.FormatFloat(val, 'f', -1, 64)
		case bool:
			strs[i] = strconv.FormatBool(val)
		default:
			strs[i] = fmt.Sprintf("%v", val)
		}
	}

	return strings.Join(strs, sep)
}

// JoinWithConverter returns a string contains converted items joined with sep.
func JoinWithConverter[T any](from []T, sep string, convertFunc func(item T) string) string {
	strs := make([]string, len(from))
	for i, v := range from {
		strs[i] = convertFunc(v)
	}

	return strings.Join(strs, sep)
}

// SliceSlices returns a new slice contains slices with maximum length each.
func SliceSlices[T any](from []T, each int) [][]T {
	result := make([][]T, 0, len(from)/each+1)

	for n := 0; ; n += each {
		if n+each >= len(from) {
			result = append(result, from[n:])
			break
		}
		result = append(result, from[n:n+each])
	}

	return result
}

// Intersection returns a new slice contains items that are in both a and b.
func Intersection[T comparable](a, b []T) []T {
	pendingChecks := make(map[T]int)

	for _, v := range a {
		pendingChecks[v] = 1
	}

	for _, v := range b {
		pendingChecks[v] |= 2
	}

	intersectionResult := make([]T, 0, len(pendingChecks))

	for k, v := range pendingChecks {
		if v == 3 {
			intersectionResult = append(intersectionResult, k)
		}
	}

	return intersectionResult
}

// FindDuplicates returns a new slice contains items that are duplicated in a.
func FindDuplicates[T comparable](a []T) []T {
	pendingChecks := make(map[T]int)

	for _, v := range a {
		pendingChecks[v]++
	}

	duplicates := make([]T, 0, len(pendingChecks))

	for k, v := range pendingChecks {
		if v > 1 {
			duplicates = append(duplicates, k)
		}
	}

	return duplicates
}
