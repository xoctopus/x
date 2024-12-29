package testdata

import (
	"fmt"
)

type (
	Func  func(x, y string, z ...int) (Boolean, error)
	Curry func(x String, y ...fmt.Stringer) func() string
)

func Sort[T comparable](v ...T) []T { return v }

var (
	SortIntegers = Sort[int]
)

type Functions struct {
	Func  Func
	Curry Curry
	Uname func() func() func() string
}
