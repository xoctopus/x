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
	// Matcher is an alias for internal.Matcher.
	Matcher[Actual any] = internal.Matcher[Actual]

	// NormalizedExpectedMatcher is an alias for internal.NormalizedExpectedMatcher.
	NormalizedExpectedMatcher = internal.NormalizedExpectedMatcher

	// MatchFunc defines a function that evaluates an actual value and returns a boolean indicating success.
	MatchFunc[Actual any] func(Actual) bool

	// ComparableMatchFunc defines a function that compares an actual value against an expected value.
	ComparableMatchFunc[Actual any, Expect any] func(Actual, Expect) bool
)

// NewMatcher creates a new Matcher with the given name and matching function.
func NewMatcher[Actual any](name string, matcher MatchFunc[Actual]) Matcher[Actual] {
	return internal.NewMatcher(name, matcher)
}

// NewComparedMatcher creates a matcher constructor that binds an expected value to a ComparableMatchFunc.
func NewComparedMatcher[Actual any, Expect any](name string, matcher ComparableMatchFunc[Actual, Expect]) internal.MatcherNewer[Actual, Expect] {
	return func(expect Expect) internal.Matcher[Actual] {
		return internal.NewComparedMatcher(name, matcher)(expect)
	}
}

// Not negates the result of the provided matcher.
func Not[Actual any](matcher Matcher[Actual]) Matcher[Actual] {
	return internal.Not(matcher)
}

// BeNil creates a matcher that asserts the actual value is nil.
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

// NotBeNil creates a matcher that asserts the actual value is not nil.
func NotBeNil[Actual any]() Matcher[Actual] {
	return internal.Not(BeNil[Actual]())
}

// BeTrue creates a matcher that asserts the actual boolean value is true.
func BeTrue() Matcher[bool] {
	return Be(true)
}

// BeFalse creates a matcher that asserts the actual boolean value is false.
func BeFalse() Matcher[bool] {
	return Be(false)
}

// IsZero creates a matcher that asserts the actual value is the zero value for its type.
func IsZero[Actual any]() Matcher[Actual] {
	return NewMatcher[Actual](
		"IsZero",
		func(a Actual) bool { return reflectx.IsZero(a) },
	)
}

// IsNotZero creates a matcher that asserts the actual value is not the zero value for its type.
func IsNotZero[Actual any]() Matcher[Actual] {
	return Not(IsZero[Actual]())
}

// Be creates a matcher that asserts the actual value is strictly identical to the expected value.
func Be[T any](expect T) Matcher[T] {
	return NewComparedMatcher[T, T]("Be", func(actual, expect T) bool {
		return any(actual) == any(expect)
	})(expect)
}

// NotBe creates a matcher that asserts the actual value is not strictly identical to the expected value.
func NotBe[T any](expect T) Matcher[T] {
	return Not(Be[T](expect))
}

// Equal creates a matcher that asserts the actual value is deeply equal to the expected value.
func Equal[T any](expect T) Matcher[T] {
	return NewComparedMatcher[T, T]("Equal", func(actual, expect T) bool {
		return reflect.DeepEqual(actual, expect)
	})(expect)
}

// NotEqual creates a matcher that asserts the actual value is not deeply equal to the expected value.
func NotEqual[T any](expect T) Matcher[T] {
	return Not(Equal[T](expect))
}

// BeGt creates a matcher that asserts the actual value is greater than the expected value.
func BeGt[T cmp.Ordered](expect T) Matcher[T] {
	return NewComparedMatcher[T, T]("BeGt", func(actual, expect T) bool {
		return actual > expect
	})(expect)
}

// BeGte creates a matcher that asserts the actual value is greater than or equal to the expected value.
func BeGte[T cmp.Ordered](expect T) Matcher[T] {
	return NewComparedMatcher[T, T]("BeGte", func(actual, expect T) bool {
		return actual >= expect
	})(expect)
}

// BeLt creates a matcher that asserts the actual value is less than the expected value.
func BeLt[T cmp.Ordered](expect T) Matcher[T] {
	return NewComparedMatcher[T, T]("BeLt", func(actual, expect T) bool {
		return actual < expect
	})(expect)
}

// BeLte creates a matcher that asserts the actual value is less than or equal to the expected value.
func BeLte[T cmp.Ordered](expect T) Matcher[T] {
	return NewComparedMatcher[T, T]("BeLte", func(actual, expect T) bool {
		return actual <= expect
	})(expect)
}

// HaveCap creates a matcher that asserts the actual collection has the expected capacity.
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

// HaveLen creates a matcher that asserts the actual collection has the expected length.
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

// HaveKey creates a matcher that asserts the actual map contains the expected key.
func HaveKey[K comparable, V any, M ~map[K]V](key K) Matcher[M] {
	return NewComparedMatcher("HaveKey", func(m M, k K) bool {
		_, ok := m[k]
		return ok
	})(key)
}

// HavePrefix creates a matcher that asserts the actual string starts with the expected prefix.
func HavePrefix(prefix string) Matcher[string] {
	return NewComparedMatcher[string, string]("HavePrefix", strings.HasPrefix)(prefix)
}

// HaveSuffix creates a matcher that asserts the actual string ends with the expected suffix.
func HaveSuffix(suffix string) Matcher[string] {
	return NewComparedMatcher[string, string]("HaveSuffix", strings.HasSuffix)(suffix)
}

// ContainsSubString creates a matcher that asserts the actual string contains the expected substring.
func ContainsSubString(sub string) Matcher[string] {
	return NewComparedMatcher[string, string]("ContainsSubString", strings.Contains)(sub)
}

// MatchRegexp creates a matcher that asserts the actual string matches the expected regular expression pattern.
func MatchRegexp(pattern string) Matcher[string] {
	return NewComparedMatcher[string, string]("MatchRegexp", func(actual, pattern string) bool {
		matched, err := regexp.MatchString(pattern, actual)
		return err == nil && matched
	})(pattern)
}

// Contains creates a matcher that asserts the actual slice contains the expected element.
func Contains[E comparable, S ~[]E](v E) Matcher[S] {
	return NewComparedMatcher("Contains", slices.Contains[S, E])(v)
}

// EquivalentSlice creates a matcher that asserts the actual slice has the same elements as the expected slice, regardless of order.
func EquivalentSlice[E comparable, S ~[]E](expect S) Matcher[S] {
	return NewComparedMatcher("EqualElements", slicex.Equivalent[E, S])(expect)
}

// ConsistOfSlice creates a matcher that asserts the actual slice has exactly the same elements in the same order as the expected slice.
func ConsistOfSlice[E comparable, S ~[]E](expect S) Matcher[S] {
	return NewComparedMatcher("ConsistOfSlice", slicex.Equivalent[E, S])(expect)
}

// BeAssignableTo creates a matcher that asserts the actual value is assignable to the type T.
func BeAssignableTo[T any]() Matcher[any] {
	typ := reflect.TypeFor[T]()
	return NewMatcher[any](
		fmt.Sprintf("BeAssignableTo[%s]", typ),
		func(actual any) bool {
			return actual != nil && reflect.TypeOf(actual).AssignableTo(typ)
		},
	)
}

// BeConvertibleTo creates a matcher that asserts the actual value is convertible to the type T.
func BeConvertibleTo[T any]() Matcher[any] {
	typ := reflect.TypeFor[T]()
	return NewMatcher[any](
		fmt.Sprintf("BeConvertibleTo[%s]", typ),
		func(actual any) bool {
			return actual != nil && reflect.TypeOf(actual).ConvertibleTo(typ)
		},
	)
}

// IsType creates a matcher that asserts the actual value is exactly of the type T.
func IsType[T any]() Matcher[any] {
	typ := reflect.TypeFor[T]()
	return NewMatcher[any](
		fmt.Sprintf("IsType[%s]", typ),
		func(actual any) bool {
			return actual != nil && reflect.TypeOf(actual) == typ
		},
	)
}

// IsError creates a matcher that asserts the actual error is the expected error (using errors.Is).
func IsError(expect error) Matcher[error] {
	return NewComparedMatcher[error, error]("IsError", errors.Is)(expect)
}

// AsError creates a matcher that asserts the actual error can be assigned to the expected error target (using errors.As).
func AsError[T any](expect *T) Matcher[error] {
	return NewComparedMatcher[error, any]("AsError", errors.As)(expect)
}

// AsErrorType creates a matcher that asserts the actual error can be assigned to a specific error type T.
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

// IsCodeError creates a matcher that asserts the actual error implements codex.Error and matches the expected code.
func IsCodeError[Code codex.Code](expect Code) Matcher[error] {
	return NewComparedMatcher[error, Code](
		fmt.Sprintf("IsCodeError[%s]", reflect.TypeFor[Code]()),
		codex.IsCode,
	)(expect)
}

// ErrorEqual creates a matcher that asserts the actual error's message exactly matches the expected string.
func ErrorEqual(expect string) Matcher[error] {
	return NewComparedMatcher(
		"ErrorEqual",
		func(actual error, expect string) bool {
			return actual != nil && actual.Error() == expect
		},
	)(expect)
}

// ErrorContains creates a matcher that asserts the actual error's message contains the expected substring.
func ErrorContains(sub string) Matcher[error] {
	return NewComparedMatcher(
		"ErrorContains",
		func(actual error, sub string) bool {
			return actual != nil && strings.Contains(actual.Error(), sub)
		},
	)(sub)
}

// Succeed creates a matcher that asserts the actual error is nil.
func Succeed() Matcher[error] {
	return NewMatcher("Succeed", func(e error) bool { return e == nil })
}

// Failed creates a matcher that asserts the actual error is not nil.
func Failed() Matcher[error] {
	return Not(Succeed())
}
