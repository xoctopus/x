package testx

import (
	"fmt"
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

	defer func() {
		recovered := recover()
		if len(matchers) == 0 {
			Expect(t, recovered, NotBeNil[any]())
		} else {
			if x, ok := recovered.(A); ok {
				Expect(t, x, matchers...)
			} else {
				panic(fmt.Sprintf("caught panic `%T`, is not `%s`", recovered, reflect.TypeFor[A]()))
			}
		}
	}()

	f()
}
