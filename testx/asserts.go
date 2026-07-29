package testx

import (
	"reflect"
	"testing"

	"github.com/xoctopus/x/testx/internal"
)

// Expect asserts that the actual value satisfies all the provided matchers.
// It fails the test if any matcher returns false.
func Expect[A any](t testing.TB, actual A, matchers ...Matcher[A]) {
	t.Helper()
	for i := range matchers {
		internal.Expect(t, actual, matchers[i])
	}
}

// ExpectPanic asserts that the provided function f panics with a value of type A,
// and optionally asserts that the recovered panic value satisfies the provided matchers.
func ExpectPanic[A any](t testing.TB, f func(), matchers ...Matcher[A]) {
	t.Helper()

	var recovered any

	func() {
		t.Helper()

		defer func() {
			t.Helper()
			recovered = recover()
		}()

		f()
	}()

	if recovered == nil {
		t.Fatalf("expect a panic of `%s`, but f returned normally", reflect.TypeFor[A]())
		return
	}
	if x, ok := recovered.(A); ok {
		if len(matchers) > 0 {
			Expect(t, x, matchers...)
		} else {
			Expect(t, recovered, NotBeNil[any]())
		}
		return
	}
	t.Fatalf("expect a panic of `%s`, but got %s", reflect.TypeFor[A](), reflect.TypeOf(recovered))
}
