package testx

import (
	"cmp"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/reflectx"
	"github.com/xoctopus/x/slicex"
	"github.com/xoctopus/x/testx/internal"
)

type (
	Matcher[Actual any]       = internal.Matcher[Actual]
	NormalizedExpectedMatcher = internal.NormalizedExpectedMatcher

	MatchFunc[Actual any] func(Actual) bool

	ComparableMatchFunc[Actual any, Expect any] func(Actual, Expect) bool
)

func NewMatcher[Actual any](name string, matcher MatchFunc[Actual]) Matcher[Actual] {
	return internal.NewMatcher(name, matcher)
}

func NewComparedMatcher[Actual any, Expect any](name string, matcher ComparableMatchFunc[Actual, Expect]) internal.MatcherNewer[Actual, Expect] {
	return func(expect Expect) internal.Matcher[Actual] {
		return internal.NewComparedMatcher(name, matcher)(expect)
	}
}

func Not[Actual any](matcher Matcher[Actual]) Matcher[Actual] {
	return internal.Not(matcher)
}

func BeNil[Actual any]() Matcher[Actual] {
	return NewMatcher[Actual]("BeNil", func(actual Actual) bool {
		v := reflect.ValueOf(actual)
		if !v.IsValid() {
			return true
		}
		switch v.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
			return v.IsNil()
		default:
			return false
		}
	})
}

func NotBeNil[Actual any]() Matcher[Actual] {
	return internal.Not(BeNil[Actual]())
}

func BeTrue() Matcher[bool] {
	return Be(true)
}

func BeFalse() Matcher[bool] {
	return Be(false)
}

func IsZero[Actual any]() Matcher[Actual] {
	return NewMatcher[Actual](
		"IsZero",
		func(a Actual) bool { return reflectx.IsZero(a) },
	)
}

func IsNotZero[Actual any]() Matcher[Actual] {
	return Not(IsZero[Actual]())
}

func Be[T any](expect T) Matcher[T] {
	return NewComparedMatcher[T, T]("Be", func(actual, expect T) bool {
		return any(actual) == any(expect)
	})(expect)
}

func NotBe[T any](expect T) Matcher[T] {
	return Not(Be[T](expect))
}

func Equal[T any](expect T) Matcher[T] {
	return NewComparedMatcher[T, T]("Equal", func(actual, expect T) bool {
		return reflect.DeepEqual(actual, expect)
	})(expect)
}

func NotEqual[T any](expect T) Matcher[T] {
	return Not(Equal[T](expect))
}

func BeGt[T cmp.Ordered](expect T) Matcher[T] {
	return NewComparedMatcher[T, T]("BeGt", func(actual, expect T) bool {
		return actual > expect
	})(expect)
}

func BeGte[T cmp.Ordered](expect T) Matcher[T] {
	return NewComparedMatcher[T, T]("BeGte", func(actual, expect T) bool {
		return actual >= expect
	})(expect)
}

func BeLt[T cmp.Ordered](expect T) Matcher[T] {
	return NewComparedMatcher[T, T]("BeLt", func(actual, expect T) bool {
		return actual < expect
	})(expect)
}

func BeLte[T cmp.Ordered](expect T) Matcher[T] {
	return NewComparedMatcher[T, T]("BeLte", func(actual, expect T) bool {
		return actual <= expect
	})(expect)
}

func HaveCap[T any](cap int) Matcher[T] {
	return NewComparedMatcher[T, int]("HaveCap", func(actual T, cap int) bool {
		v := reflect.ValueOf(actual)
		switch v.Kind() {
		case reflect.Array, reflect.Chan, reflect.Slice:
			return v.Cap() == cap
		default:
			return false
		}
	})(cap)
}

func HaveLen[T any](len int) Matcher[T] {
	return NewComparedMatcher[T, int]("HaveLen", func(actual T, len int) bool {
		v := reflect.ValueOf(actual)
		switch v.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			return v.Len() == len
		default:
			return false
		}
	})(len)
}

func HaveKey[K comparable, V any, M ~map[K]V](key K) Matcher[M] {
	return NewComparedMatcher("HaveKey", func(m M, k K) bool {
		_, ok := m[k]
		return ok
	})(key)
}

func HavePrefix(prefix string) Matcher[string] {
	return NewComparedMatcher[string, string]("HavePrefix", strings.HasPrefix)(prefix)
}

func HaveSuffix(suffix string) Matcher[string] {
	return NewComparedMatcher[string, string]("HaveSuffix", strings.HasSuffix)(suffix)
}

func ContainsSubString(sub string) Matcher[string] {
	return NewComparedMatcher[string, string]("ContainsSubString", strings.Contains)(sub)
}

func MatchRegexp(pattern string) Matcher[string] {
	return NewComparedMatcher[string, string]("MatchRegexp", func(actual, pattern string) bool {
		matched, err := regexp.MatchString(pattern, actual)
		return err == nil && matched
	})(pattern)
}

func Contains[E comparable, S ~[]E](v E) Matcher[S] {
	return NewComparedMatcher("Contains", slices.Contains[S, E])(v)
}

func EquivalentSlice[E comparable, S ~[]E](expect S) Matcher[S] {
	return NewComparedMatcher("EqualElements", slicex.Equivalent[E, S])(expect)
}

func ConsistOfSlice[E comparable, S ~[]E](expect S) Matcher[S] {
	return NewComparedMatcher("ConsistOfSlice", slicex.Equivalent[E, S])(expect)
}

func BeAssignableTo[T any]() Matcher[any] {
	typ := reflect.TypeFor[T]()
	return NewMatcher[any](
		fmt.Sprintf("BeAssignableTo[%s]", typ),
		func(actual any) bool {
			return actual != nil && reflect.TypeOf(actual).AssignableTo(typ)
		},
	)
}

func BeConvertibleTo[T any]() Matcher[any] {
	typ := reflect.TypeFor[T]()
	return NewMatcher[any](
		fmt.Sprintf("BeConvertibleTo[%s]", typ),
		func(actual any) bool {
			return actual != nil && reflect.TypeOf(actual).ConvertibleTo(typ)
		},
	)
}

func IsType[T any]() Matcher[any] {
	typ := reflect.TypeFor[T]()
	return NewMatcher[any](
		fmt.Sprintf("IsType[%s]", typ),
		func(actual any) bool {
			return actual != nil && reflect.TypeOf(actual) == typ
		},
	)
}

func IsError(expect error) Matcher[error] {
	return NewComparedMatcher[error, error]("IsError", errors.Is)(expect)
}

func AsError[T any](expect *T) Matcher[error] {
	return NewComparedMatcher[error, any]("AsError", errors.As)(expect)
}

func AsErrorType[T error]() Matcher[error] {
	typ := reflect.TypeFor[T]()
	return NewMatcher[error](
		fmt.Sprintf("AsErrorType[%s]", typ),
		func(actual error) bool {
			if actual == nil {
				return false
			}
			_, ok := errors.AsType[T](actual)
			return ok
		},
	)
}

func IsCodeError[Code codex.Code](expect Code) Matcher[error] {
	return NewComparedMatcher[error, Code](
		fmt.Sprintf("IsCodeError[%s]", reflect.TypeFor[Code]()),
		codex.IsCode,
	)(expect)
}

func ErrorEqual(expect string) Matcher[error] {
	return NewComparedMatcher(
		"ErrorEqual",
		func(actual error, expect string) bool {
			return actual != nil && actual.Error() == expect
		},
	)(expect)
}

func ErrorContains(sub string) Matcher[error] {
	return NewComparedMatcher(
		"ErrorContains",
		func(actual error, sub string) bool {
			return actual != nil && strings.Contains(actual.Error(), sub)
		},
	)(sub)
}

func Succeed() Matcher[error] {
	return NewMatcher("Succeed", func(e error) bool { return e == nil })
}

func Failed() Matcher[error] {
	return Not(Succeed())
}
