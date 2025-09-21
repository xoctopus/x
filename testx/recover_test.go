package testx_test

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/xoctopus/x/testx"
)

func TestRecover(t *testing.T) {
	t.Run("CatchError", func(t *testing.T) {
		defer func() {
			testx.AssertRecoverEqual(t, recover(), "any")
		}()
		panic(errors.New("any"))
	})
	t.Run("CatchString", func(t *testing.T) {
		defer func() {
			testx.AssertRecoverEqual(t, recover(), "string")
		}()
		panic("string")
	})
	t.Run("Catch", func(t *testing.T) {
		defer func() {
			testx.AssertRecoverContains(t, recover(), "1")
		}()
		panic(1)
	})
}
