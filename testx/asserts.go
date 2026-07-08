package testx

import (
	"reflect"
	"testing"

	"github.com/xoctopus/x/testx/internal"
)

func Expect[A any](t testing.TB, actual A, matchers ...Matcher[A]) {
	t.Helper()
	for i := range matchers {
		internal.Expect(t, actual, matchers[i])
	}
}

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
