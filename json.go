package xo

import "encoding/json"

// IsJSON determines whether the string is JSON.
func IsJSON(str string) bool {
	if str == "" {
		return false
	}

	var js json.RawMessage

	return json.Unmarshal([]byte(str), &js) == nil
}

// IsJSONBytes determines whether the bytes is JSON.
func IsJSONBytes(bytes []byte) bool {
	if len(bytes) == 0 {
		return false
	}

	var js json.RawMessage

	return json.Unmarshal(bytes, &js) == nil
}
