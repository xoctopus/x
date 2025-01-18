package testdata

import (
	"cmp"
	"fmt"
)

type (
	Func              func(x, y string, z ...int) (Boolean, error)
	Curry             func(x String, y ...fmt.Stringer) func() string
	Max[T comparable] func(...T) T
)

type Functions struct {
	Func  Func
	Curry Curry
	Uname func() func() func() string
	Max   Max[int]
}

var (
	CompareInt         = cmp.Compare[int]
	CompareNamedString = cmp.Compare[String]
)

func (v Max[T]) Compute(e ...T) T {
	return v(e...)
}
