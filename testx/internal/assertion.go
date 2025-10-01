package internal

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Expect[A any](t testing.TB, actual A, m Matcher[A]) {
	matched := m.Match(actual)
	if m.Negative() {
		matched = !matched
	}
	if !matched {
		t.Helper()
		t.Fatal("\n" + failed(actual, m))
	}
}

func failed[T any](actual T, m Matcher[T]) string {
	var v any = actual
	if normalizer, ok := m.(NormalizedMatcher[T]); ok {
		v = normalizer.NormalizeActual(actual)
	}
	prefix := "should"
	if m.Negative() {
		prefix = "should not"
	}
	return fmt.Sprintf("%s %s, but got\n%v", prefix, m.Action(), diff(v, m))
}

func diff(actual any, m any) any {
	if normalizer, ok := m.(NormalizedExpectedMatcher); ok {
		return cmp.Diff(
			actual,
			normalizer.NormalizeExpect(),
			cmpopts.IgnoreUnexported(),
			cmpopts.EquateErrors(),
		)
	}
	return actual
}
