package bdd

import (
	"cmp"
	"testing"

	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/testx"
	"github.com/xoctopus/x/testx/internal"
)

type Checker interface {
	Check(t TB)
}

func AsChecker[T any](matcher internal.Matcher[T], actual T) Checker {
	return &checker[T]{
		Matcher: matcher,
		actual:  actual,
	}
}

func AsNegativeChecker[T any](matcher internal.Matcher[T], actual T) Checker {
	return &checker[T]{
		Matcher: internal.Not(matcher),
		actual:  actual,
	}
}

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

func BeNil[A any](a A) Checker {
	return AsChecker(testx.BeNil[A](), a)
}

func NotBeNil[A any](a A) Checker {
	return NegativeChecker[A](BeNil(a))
}

func BeTrue(a bool) Checker {
	return AsChecker(testx.BeTrue(), a)
}

func BeFalse(a bool) Checker {
	return NegativeChecker[bool](BeTrue(a))
}

func IsZero[A any](expect A) Checker {
	return AsChecker(testx.IsZero[A](), expect)
}

func IsNotZero[A any](expect A) Checker {
	return NegativeChecker[A](IsZero(expect))
}

func Be[A any](actual, expect A) Checker {
	return AsChecker(testx.Be[A](expect), actual)
}

func NotBe[A any](actual, expect A) Checker {
	return NegativeChecker[A](Be(expect, actual))
}

func Equal[A any](actual, expect A) Checker {
	return AsChecker(testx.Equal(expect), actual)
}

func NotEqual[A any](actual, expect A) Checker {
	return NegativeChecker[A](Equal(expect, actual))
}

func BeGt[T cmp.Ordered](actual, expect T) Checker {
	return AsChecker(testx.BeGt(expect), actual)
}

func BeGte[T cmp.Ordered](actual, expect T) Checker {
	return AsChecker(testx.BeGte(expect), actual)
}

func BeLt[T cmp.Ordered](actual, expect T) Checker {
	return AsChecker(testx.BeLt(expect), actual)
}

func BeLte[T cmp.Ordered](actual, expect T) Checker {
	return AsChecker(testx.BeLte(expect), actual)
}

func HaveCap[A any](a A, cap int) Checker {
	return AsChecker(testx.HaveCap[A](cap), a)
}

func HaveLen[A any](a A, len int) Checker {
	return AsChecker(testx.HaveLen[A](len), a)
}

func HaveKey[K comparable, V any, M ~map[K]V](m M, k K) Checker {
	return AsChecker[M](testx.HaveKey[K, V, M](k), m)
}

func HavePrefix(s, prefix string) Checker {
	return AsChecker(testx.HavePrefix(prefix), s)
}

func HaveSuffix(s, suffix string) Checker {
	return AsChecker(testx.HaveSuffix(suffix), s)
}

func ContainsSubString(s, sub string) Checker {
	return AsChecker(testx.ContainsSubString(sub), s)
}

func MatchRegexp(pattern string, actual string) Checker {
	return AsChecker(testx.MatchRegexp(pattern), actual)
}

func Contains[E comparable, S ~[]E](s S, v E) Checker {
	return AsChecker(testx.Contains[E, S](v), s)
}

func EquivalentSlice[E comparable, S ~[]E](expect, actual S) Checker {
	return AsChecker(testx.EquivalentSlice[E, S](expect), actual)
}

func ConsistOfSlice[E comparable, S ~[]E](expect, actual S) Checker {
	return AsChecker(testx.ConsistOfSlice[E, S](expect), actual)
}

func BeAssignableTo[E any](actual any) Checker {
	return AsChecker(testx.BeAssignableTo[E](), actual)
}

func BeConvertibleTo[E any](actual any) Checker {
	return AsChecker(testx.BeConvertibleTo[E](), actual)
}

func IsType[E any](actual any) Checker {
	return AsChecker(testx.IsType[E](), actual)
}

func IsError(expect, actual error) Checker {
	return AsChecker(testx.IsError(expect), actual)
}

func AsError(expect *error, actual error) Checker {
	return AsChecker(testx.AsError(expect), actual)
}

func AsErrorType[T error](actual error) Checker {
	return AsChecker(testx.AsErrorType[T](), actual)
}

func IsCodeError[Code codex.Code](actual error, expect Code) Checker {
	return AsChecker(testx.IsCodeError(expect), actual)
}

func ErrorEqual(actual error, expect string) Checker {
	return AsChecker(testx.ErrorEqual(expect), actual)
}

func ErrorContains(err error, sub string) Checker {
	return AsChecker(testx.ErrorContains(sub), err)
}

func Succeed(err error) Checker {
	return AsChecker(testx.Succeed(), err)
}

func Failed(err error) Checker {
	return AsChecker(testx.Failed(), err)
}
