package internal

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Expect[A any](t testing.TB, actual A, m Matcher[A]) {
	t.Helper()

	matched := m.Match(actual)
	if m.Negative() {
		matched = !matched
	}
	if !matched {
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

func diff(actual any, m any) (x any) {
	if normalizer, ok := m.(NormalizedExpectedMatcher); ok {
		expect := normalizer.NormalizeExpect()
		defer func() {
			if r := recover(); r != nil {
				// some case may cause cmp.Diff panic. catch to compare dumped
				x = cmp.Diff(
					fmt.Sprintf("%T:%+v", expect, expect),
					fmt.Sprintf("%T:%+v", expect, actual),
				)
			}
		}()
		return cmp.Diff(expect, actual)
	}
	return actual
}
