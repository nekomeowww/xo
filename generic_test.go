package xo

import (
	"reflect"
	"testing"
)

func TestToPtrAny(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  any
	}{
		{
			name:  "string value",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "integer value",
			input: 42,
			want:  42,
		},
		{
			name:  "struct value",
			input: struct{ Name string }{"test"},
			want:  struct{ Name string }{"test"},
		},
		{
			name:  "nil value",
			input: nil,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToPtrAny(tt.input)
			if got == nil {
				t.Fatal("expected non-nil pointer")
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("ToPtrAny() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestFromPtrAny(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  any
	}{
		{
			name:  "string conversion",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "integer conversion",
			input: 42,
			want:  42,
		},
		{
			name:  "failed conversion returns zero value",
			input: "not an int",
			want:  0, // Testing conversion to int
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := ToPtrAny(tt.input)

			switch tt.want.(type) {
			case string:
				got := FromPtrAny[string](ptr)
				if got != tt.want {
					t.Errorf("FromPtrAny[string]() = %v, want %v", got, tt.want)
				}
			case int:
				got := FromPtrAny[int](ptr)
				if got != tt.want {
					t.Errorf("FromPtrAny[int]() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	tests := []struct {
		name  string
		value any
	}{
		{
			name:  "string round trip",
			value: "hello world",
		},
		{
			name:  "integer round trip",
			value: 42,
		},
		{
			name:  "struct round trip",
			value: Person{Name: "John", Age: 30},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := ToPtrAny(tt.value)
			switch v := tt.value.(type) {
			case string:
				got := FromPtrAny[string](ptr)
				if got != v {
					t.Errorf("Round trip failed for string: got %v, want %v", got, v)
				}
			case int:
				got := FromPtrAny[int](ptr)
				if got != v {
					t.Errorf("Round trip failed for int: got %v, want %v", got, v)
				}
			case Person:
				got := FromPtrAny[Person](ptr)
				if !reflect.DeepEqual(got, v) {
					t.Errorf("Round trip failed for struct: got %v, want %v", got, v)
				}
			}
		})
	}
}
