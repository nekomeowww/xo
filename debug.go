package xo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// Print formats the output of all incoming values in terms of field, value,
// type, and size.
func Print(inputs ...interface{}) {
	fmt.Println(Sprint(inputs))
}

// Sprint formats the output of all the fields, values, types, and sizes of
// the values passed in and returns the string.
//
// NOTICE: newline control character is included.
func Sprint(inputs ...interface{}) string {
	return spew.Sdump(inputs)
}

// PrintJSON formats the output of all incoming values in JSON format.
func PrintJSON(inputs ...interface{}) {
	fmt.Println(SprintJSON(inputs))
}

// SprintJSON formats the output of all incoming values in JSON format and
//
// NOTICE: newline control character is included.
func SprintJSON(inputs ...interface{}) string {
	strSlice := make([]string, 0)

	for _, v := range inputs {
		b, _ := json.MarshalIndent(v, "", "  ")
		strSlice = append(strSlice, string(b))
	}

	return strings.Join(strSlice, "\n")
}
