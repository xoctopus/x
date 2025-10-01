package reflectx_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/reflectx"
	"github.com/xoctopus/x/testx"
)

func TestAssertType(t *testing.T) {
	NewWithT(t).Expect(reflectx.MustType[int](100)).To(Equal(100))
	t.Run("Panic", func(t *testing.T) {
		testx.ExpectPanic(
			t,
			func() { reflectx.MustType[int](100.0) },
			testx.ErrorContains("must type"),
		)
	})
	NewWithT(t).Expect(reflectx.CanCast[int](1)).To(BeTrue())
}
