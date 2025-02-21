package testx_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/xoctopus/x/testx"
)

func TestRecover(t *testing.T) {
	t.Run("CatchError", func(t *testing.T) {
		defer func() {
			testx.AssertRecoverEqual(t, recover(), "any")
		}()
		func() {
			panic(errors.New("any"))
		}()
	})
	t.Run("Catch", func(t *testing.T) {
		defer func() {
			testx.AssertRecoverContains(t, recover(), "1")
		}()
		func() {
			panic(1)
		}()
	})
	t.Run("CatchNothing", func(t *testing.T) {
		defer func() {
			NewWithT(t).Expect(testx.Recover(recover())).To(Equal(""))
		}()
	})
}
