package resultx_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/x/resultx"
)

func f() (int, bool, string) { return 100, true, "f" }

func TestAt(t *testing.T) {
	rs := resultx.ResultsOf(f())
	NewWithT(t).Expect(resultx.At[int](rs, 0)).To(Equal(100))
	NewWithT(t).Expect(rs.At(0)).To(Equal(100))
	NewWithT(t).Expect(resultx.At[bool](rs, 1)).To(BeTrue())
	NewWithT(t).Expect(rs.At(1)).To(BeTrue())
	NewWithT(t).Expect(resultx.At[string](rs, 2)).To(Equal("f"))
	NewWithT(t).Expect(rs.At(2)).To(Equal("f"))
}
