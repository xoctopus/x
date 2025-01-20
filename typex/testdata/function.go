package testdata

import (
	"cmp"
	"fmt"
)

type (
	Func                           func(x, y string, z ...int) (Boolean, error)
	Curry                          func(x String, y ...fmt.Stringer) func() string
	Max[T comparable]              func(...T) T
	Contains[S ~[]E, E comparable] func(s S, v E) bool
	ContainsFunc[S ~[]E, E any]    func(s S, f func(E) bool) bool
)

type Functions struct {
	Func            Func
	Curry           Curry
	Uname           func() func() func() string
	Max             Max[int]
	ContainsInt     Contains[[]int, int]
	ContainsIntFunc ContainsFunc[[]int, int]
}

var (
	CompareInt         = cmp.Compare[int]
	CompareNamedString = cmp.Compare[String]
)

func (v Max[T]) Compute(e ...T) T {
	return v(e...)
}
