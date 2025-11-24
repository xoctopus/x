package resultx_test

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/xoctopus/x/resultx"
	. "github.com/xoctopus/x/testx"
)

func TestResult(t *testing.T) {
	Expect(t, resultx.Succeed[int64](0).Unwrap(), Equal(int64(0)))
	Expect(t, resultx.Succeed[int64](0).UnwrapOr(1), Equal(int64(0)))
	Expect(t, resultx.Err[float32](errors.New("any")).UnwrapOr(0.1), Equal(float32(0.1)))

	Expect(t, resultx.OK[byte]('x').Unwrap(), Equal(byte('x')))
	Expect(t, resultx.OK[byte]('x').UnwrapOr('y'), Equal(byte('x')))
	Expect(t, resultx.NotOK[rune]().UnwrapOr('y'), Equal('y'))

	t.Run("Unwrap", func(t *testing.T) {
		Expect(t, resultx.Unwrap(strconv.Atoi("1")), Equal(1))
		t.Run("Panic", func(t *testing.T) {
			ExpectPanic[error](
				t,
				func() { resultx.Unwrap(strconv.Atoi("x")) },
			)
		})
		t.Run("WrapResult", func(t *testing.T) {
			r := resultx.WrapResult(strconv.Atoi("x"))
			ExpectPanic[error](t, func() { r.Unwrap() })
		})
	})

	t.Run("UnwrapB", func(t *testing.T) {
		Expect(t, resultx.UnwrapB(strings.CutPrefix("good morning", "good ")), Equal("morning"))
		t.Run("Panic", func(t *testing.T) {
			ExpectPanic[error](t, func() {
				resultx.UnwrapB(strings.CutPrefix("good morning", "x"))
			})
		})
		t.Run("WrapResultB", func(t *testing.T) {
			r := resultx.WrapResultB(strings.CutPrefix("good morning", "good "))
			Expect(t, r.Unwrap(), Equal("morning"))
		})
	})
}
