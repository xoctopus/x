package resultx_test

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/resultx"
)

func TestResult(t *testing.T) {
	NewWithT(t).Expect(resultx.Succeed[int64](0).Unwrap()).To(Equal(int64(0)))
	NewWithT(t).Expect(resultx.Succeed[int64](0).UnwrapOr(1)).To(Equal(int64(0)))
	NewWithT(t).Expect(resultx.Err[float32](errors.New("any")).UnwrapOr(0.1)).To(Equal(float32(0.1)))

	NewWithT(t).Expect(resultx.OK[byte]('x').Unwrap()).To(Equal(byte('x')))
	NewWithT(t).Expect(resultx.OK[byte]('x').UnwrapOr('y')).To(Equal(byte('x')))
	NewWithT(t).Expect(resultx.NotOK[rune]().UnwrapOr('y')).To(Equal('y'))

	t.Run("Unwrap", func(t *testing.T) {
		NewWithT(t).Expect(resultx.Unwrap(strconv.Atoi("1"))).To(Equal(1))
		t.Run("Panic", func(t *testing.T) {
			defer func() {
				e := recover().(error)
				NewWithT(t).Expect(e).NotTo(BeNil())
			}()
			resultx.Unwrap(strconv.Atoi("x"))
		})
		t.Run("WrapResult", func(t *testing.T) {
			r := resultx.WrapResult(strconv.Atoi("x"))
			defer func() {
				e := recover().(error)
				NewWithT(t).Expect(e).NotTo(BeNil())
			}()
			r.Unwrap()
		})
	})

	t.Run("UnwrapB", func(t *testing.T) {
		NewWithT(t).Expect(resultx.UnwrapB(strings.CutPrefix("good morning", "good "))).To(Equal("morning"))
		t.Run("Panic", func(t *testing.T) {
			defer func() {
				NewWithT(t).Expect(recover()).To(BeFalse())
			}()
			resultx.UnwrapB(strings.CutPrefix("good morning", "x"))
		})
		t.Run("WrapResultB", func(t *testing.T) {
			r := resultx.WrapResultB(strings.CutPrefix("good morning", "good "))
			NewWithT(t).Expect(r.Unwrap()).To(Equal("morning"))
		})
	})
}
