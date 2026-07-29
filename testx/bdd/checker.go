package bdd

import (
	"cmp"
	"testing"

	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/testx"
	"github.com/xoctopus/x/testx/internal"
)

// Checker defines the interface for executing a test assertion.
type Checker interface {
	Check(t TB)
}

// AsChecker creates a new Checker from a given matcher and an actual value.
func AsChecker[T any](matcher internal.Matcher[T], actual T) Checker {
	return &checker[T]{
		Matcher: matcher,
		actual:  actual,
	}
}

// AsNegativeChecker creates a new Checker that negates the result of the given matcher.
func AsNegativeChecker[T any](matcher internal.Matcher[T], actual T) Checker {
	return &checker[T]{
		Matcher: internal.Not(matcher),
		actual:  actual,
	}
}

// NegativeChecker takes an existing Checker and negates its internal matcher.
func NegativeChecker[T any](c Checker) Checker {
	u := c.(*checker[T])
	u.Matcher = internal.Not(u.Matcher)
	return u
}

type checker[T any] struct {
	internal.Matcher[T]

	actual T
}

func (c *checker[T]) Check(t TB) {
	switch x := t.(type) {
	case interface{ Unwrap() *testing.T }:
		tt := x.Unwrap()
		tt.Helper()
		internal.Expect(tt, c.actual, c.Matcher)
	case testing.TB:
		x.Helper()
		internal.Expect(x, c.actual, c.Matcher)
	}
}

// BeNil asserts that the actual value is nil.
func BeNil[A any](a A) Checker {
	return AsChecker(testx.BeNil[A](), a)
}

// NotBeNil asserts that the actual value is not nil.
func NotBeNil[A any](a A) Checker {
	return NegativeChecker[A](BeNil(a))
}

// BeTrue asserts that the actual boolean value is true.
func BeTrue(a bool) Checker {
	return AsChecker(testx.BeTrue(), a)
}

// BeFalse asserts that the actual boolean value is false.
func BeFalse(a bool) Checker {
	return NegativeChecker[bool](BeTrue(a))
}

// IsZero asserts that the actual value is the zero value for its type.
func IsZero[A any](expect A) Checker {
	return AsChecker(testx.IsZero[A](), expect)
}

// IsNotZero asserts that the actual value is not the zero value for its type.
func IsNotZero[A any](expect A) Checker {
	return NegativeChecker[A](IsZero(expect))
}

// Be asserts that the actual value is strictly identical to the expected value.
func Be[A any](actual, expect A) Checker {
	return AsChecker(testx.Be[A](expect), actual)
}

// NotBe asserts that the actual value is not strictly identical to the expected value.
func NotBe[A any](actual, expect A) Checker {
	return NegativeChecker[A](Be(expect, actual))
}

// Equal asserts that the actual value is deeply equal to the expected value.
func Equal[A any](actual, expect A) Checker {
	return AsChecker(testx.Equal(expect), actual)
}

// NotEqual asserts that the actual value is not deeply equal to the expected value.
func NotEqual[A any](actual, expect A) Checker {
	return NegativeChecker[A](Equal(expect, actual))
}

// BeGt asserts that the actual value is greater than the expected value.
func BeGt[T cmp.Ordered](actual, expect T) Checker {
	return AsChecker(testx.BeGt(expect), actual)
}

// BeGte asserts that the actual value is greater than or equal to the expected value.
func BeGte[T cmp.Ordered](actual, expect T) Checker {
	return AsChecker(testx.BeGte(expect), actual)
}

// BeLt asserts that the actual value is less than the expected value.
func BeLt[T cmp.Ordered](actual, expect T) Checker {
	return AsChecker(testx.BeLt(expect), actual)
}

// BeLte asserts that the actual value is less than or equal to the expected value.
func BeLte[T cmp.Ordered](actual, expect T) Checker {
	return AsChecker(testx.BeLte(expect), actual)
}

// HaveCap asserts that the actual collection has the expected capacity.
func HaveCap[A any](a A, cap int) Checker {
	return AsChecker(testx.HaveCap[A](cap), a)
}

// HaveLen asserts that the actual collection has the expected length.
func HaveLen[A any](a A, len int) Checker {
	return AsChecker(testx.HaveLen[A](len), a)
}

// HaveKey asserts that the actual map contains the expected key.
func HaveKey[K comparable, V any, M ~map[K]V](m M, k K) Checker {
	return AsChecker[M](testx.HaveKey[K, V, M](k), m)
}

// HavePrefix asserts that the actual string starts with the expected prefix.
func HavePrefix(s, prefix string) Checker {
	return AsChecker(testx.HavePrefix(prefix), s)
}

// HaveSuffix asserts that the actual string ends with the expected suffix.
func HaveSuffix(s, suffix string) Checker {
	return AsChecker(testx.HaveSuffix(suffix), s)
}

// ContainsSubString asserts that the actual string contains the expected substring.
func ContainsSubString(s, sub string) Checker {
	return AsChecker(testx.ContainsSubString(sub), s)
}

// MatchRegexp asserts that the actual string matches the expected regular expression pattern.
func MatchRegexp(pattern string, actual string) Checker {
	return AsChecker(testx.MatchRegexp(pattern), actual)
}

// Contains asserts that the actual slice contains the expected element.
func Contains[E comparable, S ~[]E](s S, v E) Checker {
	return AsChecker(testx.Contains[E, S](v), s)
}

// EquivalentSlice asserts that the actual slice has the same elements as the expected slice, regardless of order.
func EquivalentSlice[E comparable, S ~[]E](expect, actual S) Checker {
	return AsChecker(testx.EquivalentSlice[E, S](expect), actual)
}

// ConsistOfSlice asserts that the actual slice has exactly the same elements in the same order as the expected slice.
func ConsistOfSlice[E comparable, S ~[]E](expect, actual S) Checker {
	return AsChecker(testx.ConsistOfSlice[E, S](expect), actual)
}

// BeAssignableTo asserts that the actual value is assignable to the type E.
func BeAssignableTo[E any](actual any) Checker {
	return AsChecker(testx.BeAssignableTo[E](), actual)
}

// BeConvertibleTo asserts that the actual value is convertible to the type E.
func BeConvertibleTo[E any](actual any) Checker {
	return AsChecker(testx.BeConvertibleTo[E](), actual)
}

// IsType asserts that the actual value is exactly of the type E.
func IsType[E any](actual any) Checker {
	return AsChecker(testx.IsType[E](), actual)
}

// IsError asserts that the actual error is the expected error (using errors.Is).
func IsError(expect, actual error) Checker {
	return AsChecker(testx.IsError(expect), actual)
}

// AsError asserts that the actual error can be assigned to the expected error target (using errors.As).
func AsError(expect *error, actual error) Checker {
	return AsChecker(testx.AsError(expect), actual)
}

// AsErrorType asserts that the actual error can be assigned to a specific error type T.
func AsErrorType[T error](actual error) Checker {
	return AsChecker(testx.AsErrorType[T](), actual)
}

// IsCodeError asserts that the actual error implements codex.Error and matches the expected code.
func IsCodeError[Code codex.Code](actual error, expect Code) Checker {
	return AsChecker(testx.IsCodeError(expect), actual)
}

// ErrorEqual asserts that the actual error's message exactly matches the expected string.
func ErrorEqual(actual error, expect string) Checker {
	return AsChecker(testx.ErrorEqual(expect), actual)
}

// ErrorContains asserts that the actual error's message contains the expected substring.
func ErrorContains(err error, sub string) Checker {
	return AsChecker(testx.ErrorContains(sub), err)
}

// Succeed asserts that the actual error is nil.
func Succeed(err error) Checker {
	return AsChecker(testx.Succeed(), err)
}

// Failed asserts that the actual error is not nil.
func Failed(err error) Checker {
	return AsChecker(testx.Failed(), err)
}
