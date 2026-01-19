package bdd

import (
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

func Be[A any](expect, actual A) Checker {
	return AsChecker(testx.Be[A](expect), actual)
}

func NotBe[A any](expect, actual A) Checker {
	return NegativeChecker[A](Be(expect, actual))
}

func Equal[A any](expect, actual A) Checker {
	return AsChecker(testx.Equal(expect), actual)
}

func NotEqual[A any](expect, actual A) Checker {
	return NegativeChecker[A](Equal(expect, actual))
}

func HaveCap[A any](a A, cap int) Checker {
	return AsChecker(testx.HaveCap[A](cap), a)
}

func HaveLen[A any](a A, len int) Checker {
	return AsChecker(testx.HaveLen[A](len), a)
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

func Contains[E comparable, S ~[]E](s S, v E) Checker {
	return AsChecker(testx.Contains[E, S](v), s)
}

func EquivalentSlice[E comparable, S ~[]E](expect, actual S) Checker {
	return AsChecker(testx.EquivalentSlice[E, S](expect), actual)
}

func BeAssignableTo[E any](actual any) Checker {
	return AsChecker(testx.BeAssignableTo[E](), actual)
}

func BeConvertibleTo[E any](actual any) Checker {
	return AsChecker(testx.BeConvertibleTo[E](), actual)
}

func BeType[E any](actual any) Checker {
	return AsChecker(testx.BeType[E](), actual)
}

func IsError(expect, actual error) Checker {
	return AsChecker(testx.IsError(expect), actual)
}

func IsCodeError[Code codex.Code](expect Code, actual error) Checker {
	return AsChecker(testx.IsCodeError(expect), actual)
}

func ErrorEqual(expect string, actual error) Checker {
	return AsChecker(testx.ErrorEqual(expect), actual)
}

func ErrorContains(sub string, err error) Checker {
	return AsChecker(testx.ErrorContains(sub), err)
}

func Succeed(err error) Checker {
	return AsChecker(testx.Succeed(), err)
}

func Failed(err error) Checker {
	return AsChecker(testx.Failed(), err)
}
