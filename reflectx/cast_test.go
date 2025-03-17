package reflectx_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/reflectx"
	"github.com/xoctopus/x/testx"
)

func TestAssertType(t *testing.T) {
	NewWithT(t).Expect(reflectx.MustAssertType[int](100)).To(Equal(100))
	t.Run("Panic", func(t *testing.T) {
		defer func() {
			testx.AssertRecoverContains(t, recover(), "must assert type")
		}()
		reflectx.MustAssertType[int](100.0)
	})
	NewWithT(t).Expect(reflectx.CanCast[int](1)).To(BeTrue())
}
