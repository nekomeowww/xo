package xo

// ToPtrAny returns a pointer to the given value.
func ToPtrAny(v any) *any {
	return &v
}

// FromPtrAny returns the value from the given pointer.
func FromPtrAny[T any](v *any) T {
	val, _ := (*v).(T)
	return val
}
