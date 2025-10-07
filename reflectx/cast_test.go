package reflectx_test

import (
	"testing"

	"github.com/xoctopus/x/reflectx"
	. "github.com/xoctopus/x/testx"
)

func TestAssertType(t *testing.T) {
	Expect(t, reflectx.MustType[int](100), Equal(100))
	t.Run("Panic", func(t *testing.T) {
		ExpectPanic(
			t,
			func() { reflectx.MustType[int](100.0) },
			ErrorContains("must type"),
		)
	})
	Expect(t, reflectx.CanCast[int](1), BeTrue())
}
