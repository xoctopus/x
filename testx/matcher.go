package testx

import (
	"errors"
	"reflect"
	"slices"
	"strings"

	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/reflectx"
	"github.com/xoctopus/x/slicex"
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
	return func(expect E) internal.Matcher[A] {
		return internal.NewComparedMatcher(name, match)(expect)
	}
}

func Not[A any](matcher Matcher[A]) Matcher[A] {
	return internal.Not(matcher)
}

func BeNil[A any]() Matcher[A] {
	return NewMatcher[A]("BeNil", func(a A) bool {
		v := reflect.ValueOf(a)
		if !v.IsValid() {
			return true
		}
		switch v.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface,
			reflect.Map, reflect.Pointer, reflect.Slice:
			return v.IsNil()
		default:
			return false
		}
	})
}

func NotBeNil[A any]() Matcher[A] {
	return internal.Not(BeNil[A]())
}

func BeTrue() Matcher[bool] {
	return NewMatcher[bool]("BeTrue", func(a bool) bool { return a })
}

func BeFalse() Matcher[bool] {
	return NewMatcher[bool]("BeFalse", func(a bool) bool { return !a })
}

func IsZero[A any]() Matcher[A] {
	return NewMatcher[A]("IsZero", func(a A) bool {
		return reflectx.IsZero(a)
	})
}

func IsNotZero[A any]() Matcher[A] {
	return Not(IsZero[A]())
}

func Be[E any](expect E) Matcher[E] {
	return NewComparedMatcher[E, E]("Be", func(a, e E) bool {
		return any(a) == any(e)
	})(expect)
}

func NotBe[E any](expect E) Matcher[E] {
	return Not(Be[E](expect))
}

func Equal[E any](expect E) Matcher[E] {
	return NewComparedMatcher[E, E]("Equal", func(a, e E) bool {
		return reflect.DeepEqual(a, e)
	})(expect)
}

func NotEqual[E any](expect E) Matcher[E] {
	return Not(Equal[E](expect))
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

func Contains[E comparable, S ~[]E](v E) Matcher[S] {
	return NewComparedMatcher(
		"Contains",
		slices.Contains[S, E],
	)(v)
}

func EquivalentSlice[E comparable, S ~[]E](expect S) Matcher[S] {
	return NewComparedMatcher(
		"EqualElements",
		slicex.Equivalent[E, S],
	)(expect)
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

func IsCodeError[Code codex.Code](expect Code) Matcher[error] {
	return NewComparedMatcher[error, Code](
		"IsCodeError",
		codex.IsCode,
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
