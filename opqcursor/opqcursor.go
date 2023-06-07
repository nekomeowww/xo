package opqcursor

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/samber/lo"
)

// OpaqueCursor is a struct that contains the filters and sorters for a query.
// Generally used in GraphQL based APIs to pass the cursor as a string. But the scenarios
// are not limited to GraphQL only, you could use it in any API that requires a cursor.
type OpaqueCursor[F any, S any] struct {
	Filters F `json:"filters"`
	Sorters S `json:"sorters"`
}

// IsValid checks if the cursor is valid or not. It checks if the sorters are valid or not.
//
// NOTICE: currently the sorters are checked for ASC and DESC only.
func (oc OpaqueCursor[F, S]) IsValid() bool {
	refValue := reflect.ValueOf(oc.Sorters)
	if refValue.Kind() == reflect.Struct {
		for i := 0; i < refValue.NumField(); i++ {
			if lo.Contains([]string{"ASC", "DESC"}, strings.ToUpper(refValue.Field(i).String())) {
				return true
			}
		}
	}

	return false
}

// Unmarshal decodes the cursor string and returns the OpaqueCursor struct.
func Unmarshal[Filters any, Sorters any](cursor string) (*OpaqueCursor[Filters, Sorters], error) {
	var cursorData OpaqueCursor[Filters, Sorters]

	base64Data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return &cursorData, err
	}

	err = json.Unmarshal(base64Data, &cursorData)
	if err != nil {
		return &cursorData, err
	}

	return &cursorData, nil
}

// UnmarshalWithDefaults decodes the cursor string and returns the OpaqueCursor struct with defaults additionally.
func UnmarshalWithDefaults[Filters any, Sorters any](cursor string, defaults OpaqueCursor[Filters, Sorters]) (*OpaqueCursor[Filters, Sorters], error) {
	if cursor == "" {
		return &defaults, nil
	}

	var cursorData OpaqueCursor[Filters, Sorters]

	base64Data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return &cursorData, err
	}

	err = json.Unmarshal(base64Data, &cursorData)
	if err != nil {
		return &defaults, err
	}

	return &cursorData, nil
}

// Marshal encodes the OpaqueCursor struct and returns the cursor string.
func Marshal[Filters any, Sorters any](filters Filters, sorters Sorters) (string, error) {
	cursorData := OpaqueCursor[Filters, Sorters]{
		Filters: filters,
		Sorters: sorters,
	}

	cursor, err := json.Marshal(cursorData)
	if err != nil {
		return "{}", err
	}

	return base64.StdEncoding.EncodeToString(cursor), nil
}
