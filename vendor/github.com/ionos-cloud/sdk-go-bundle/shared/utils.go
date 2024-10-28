package shared

import (
	"reflect"
	"unicode/utf8"
)

// ToPtr - returns a pointer to the given value.
func ToPtr[T any](v T) *T {
	return &v
}

// ToValue - returns the value of the bool pointer passed in
func ToValue[T any](ptr *T) T {
	return *ptr
}

// ToValueDefault - returns the value of the pointer passed in, or the default type value if the pointer is nil
func ToValueDefault[T any](ptr *T) T {
	var defaultVal T
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

func SliceToValueDefault[T any](ptrSlice *[]T) []T {
	return append([]T{}, *ptrSlice...)
}

type Nullable[T any] struct {
	value *T
	isSet bool
}

func (v Nullable[T]) Get() *T {
	return v.value
}

func (v *Nullable[T]) Set(val *T) {
	v.value = val
	v.isSet = true
}

func (v Nullable[T]) IsSet() bool {
	return v.isSet
}

func (v *Nullable[T]) Unset() {
	v.value = nil
	v.isSet = false
}

func Strlen(s string) int {
	return utf8.RuneCountInString(s)
}

// IsNil checks if an input is nil
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	case reflect.Array:
		return reflect.ValueOf(i).IsZero()
	}
	return false
}
