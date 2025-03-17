package reflectx

import (
	"github.com/pkg/errors"
)

func MustAssertType[T any](v any) T {
	t, ok := AssertType[T](v)
	if !ok {
		panic(errors.Errorf("must assert type v (%T) is T (%T)", v, *new(T)))
	}
	return t
}

func AssertType[T any](v any) (T, bool) {
	t, ok := v.(T)
	return t, ok
}

func CanCast[T any](v any) bool {
	_, ok := AssertType[T](v)
	return ok
}
