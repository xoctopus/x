package reflectx

import (
	"fmt"
)

func MustType[T any](v any) T {
	t, ok := v.(T)
	if !ok {
		panic(fmt.Errorf("must type v (%T) is T (%T)", v, *new(T)))
	}
	return t
}

func CanCast[T any](v any) bool {
	_, ok := v.(T)
	return ok
}
