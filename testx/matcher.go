package testx

import (
	"reflect"
	"slices"
	"strings"

	"github.com/pkg/errors"

	"github.com/xoctopus/x/reflectx"
	"github.com/xoctopus/x/testx/internal"
)

type (
	Matcher[A any]            = internal.Matcher[A]
	NormalizedExpectedMatcher = internal.NormalizedExpectedMatcher
)

func NewMatcher[A any](name string, match func(A) bool) Matcher[A] {
	return internal.NewMatcher(name, match)
}

func NewComparedMatcher[A any, E any](name string, match func(A, E) bool) internal.MatcherNewer[A, E] {
	return func(expected E) internal.Matcher[A] {
		return internal.NewComparedMatcher(name, match)(expected)
	}
}

func Not[A any](matcher Matcher[A]) Matcher[A] {
	return internal.Not(matcher)
}

func BeNil[A any]() Matcher[A] {
	return NewMatcher[A]("BeNil", func(a A) bool {
		return any(a) == nil
	})
}

func NotBeNil[A any]() Matcher[A] {
	return Not(BeNil[A]())
}

func BeTrue() Matcher[bool] {
	return NewMatcher[bool]("BeTrue", func(a bool) bool { return a })
}

func BeFalse() Matcher[bool] {
	return NewMatcher[bool]("BeFalse", func(a bool) bool { return !a })
}

func BeEmpty[A any]() Matcher[A] {
	return NewMatcher[A]("BeEmpty", func(a A) bool {
		return reflectx.IsZero(a)
	})
}

func Be[T any](expected T) Matcher[T] {
	return NewComparedMatcher[T, T]("Be", func(a T, e T) bool {
		return any(a) == any(e)
	})(expected)
}

func NotBe[T any](expected T) Matcher[T] {
	return Not(Be[T](expected))
}

func Equal[T any](expected T) Matcher[T] {
	return NewComparedMatcher[T, T]("Equal", func(a T, e T) bool {
		return reflect.DeepEqual(a, e)
	})(expected)
}

func NotEqual[T any](expected T) Matcher[T] {
	return Not(Equal[T](expected))
}

func HaveCap[T any](cap int) Matcher[T] {
	return NewComparedMatcher[T, int]("HaveCap", func(a T, cap int) bool {
		return reflect.ValueOf(a).Cap() == cap
	})(cap)
}

func HaveLen[T any](len int) Matcher[T] {
	return NewComparedMatcher[T, int]("HaveLen", func(a T, len int) bool {
		return reflect.ValueOf(a).Len() == len
	})(len)
}

func HavePrefix(prefix string) Matcher[string] {
	return NewComparedMatcher[string, string](
		"HavePrefix",
		strings.HasPrefix,
	)(prefix)
}

func HaveSuffix(suffix string) Matcher[string] {
	return NewComparedMatcher[string, string](
		"HaveSuffix",
		strings.HasSuffix,
	)(suffix)
}

func ContainsSubString(sub string) Matcher[string] {
	return NewComparedMatcher[string, string](
		"ContainsSubString",
		func(s, sub string) bool {
			return strings.Contains(s, sub)
		},
	)(sub)
}

func Contains[S []E, E comparable](v E) Matcher[S] {
	return NewComparedMatcher(
		"ContainsStringItem",
		slices.Contains[S, E],
	)(v)
}

func BeAssignableTo[E any]() Matcher[any] {
	return NewComparedMatcher[any, E](
		"BeAssignableTo",
		func(a any, e E) bool {
			if a == nil {
				return false
			}
			return reflect.TypeOf(a).AssignableTo(reflect.TypeFor[E]())
		},
	)(*new(E))
}

func BeConvertibleTo[E any]() Matcher[any] {
	return NewComparedMatcher[any, E](
		"BeConvertibleTo",
		func(a any, e E) bool {
			if a == nil {
				return false
			}
			return reflect.TypeOf(a).ConvertibleTo(reflect.TypeFor[E]())
		},
	)(*new(E))
}

func BeType[E any]() Matcher[any] {
	return NewComparedMatcher[any, E](
		"BeAssignableTo",
		func(a any, e E) bool {
			if a == nil || any(e) == nil {
				return false
			}
			return reflect.TypeOf(a) == reflect.TypeFor[E]()
		},
	)(*new(E))
}

func IsError(expect error) Matcher[error] {
	return NewComparedMatcher[error, error](
		"IsError",
		errors.Is,
	)(expect)
}

func ErrorEqual(expect string) Matcher[error] {
	return NewComparedMatcher(
		"ErrorEqual",
		func(actual error, expect string) bool {
			if actual == nil {
				return false
			}
			return actual.Error() == expect
		},
	)(expect)
}

func ErrorContains(sub string) Matcher[error] {
	return NewComparedMatcher(
		"ErrorContains",
		func(actual error, sub string) bool {
			if actual == nil {
				return false
			}
			return strings.Contains(actual.Error(), sub)
		},
	)(sub)
}

func Succeed() Matcher[error] {
	return NewMatcher("Succeed", func(e error) bool { return e == nil })
}

func Failed() Matcher[error] {
	return Not(Succeed())
}
