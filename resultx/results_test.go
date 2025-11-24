package resultx_test

import (
	"testing"

	"github.com/xoctopus/x/resultx"
	. "github.com/xoctopus/x/testx"
)

func f() (int, bool, string) { return 100, true, "f" }

func TestAt(t *testing.T) {
	rs := resultx.ResultsOf(f())
	Expect(t, resultx.At[int](rs, 0), Equal(100))
	Expect(t, rs.At(0), Equal[any](100))
	Expect(t, resultx.At[bool](rs, 1), BeTrue())
	Expect(t, rs.At(1), Equal[any](true))
	Expect(t, resultx.At[string](rs, 2), Equal("f"))
	Expect(t, rs.At(2), Equal[any]("f"))
}
