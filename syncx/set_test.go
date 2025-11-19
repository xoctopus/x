package syncx_test

import (
	"testing"

	"github.com/xoctopus/x/syncx"
	. "github.com/xoctopus/x/testx"
)

func TestNewSet(t *testing.T) {
	s := syncx.NewSet[int](0)

	Expect(t, s.Exists(1), BeFalse())

	s.Store(1)
	Expect(t, s.Exists(1), BeTrue())

	Expect(t, s.Keys(), EquivalentSlice([]int{0, 1}))
	s.Store(2)
	Expect(t, s.Exists(1), BeTrue())
	Expect(t, s.Len(), Equal(3))
	Expect(t, s.Keys(), EquivalentSlice([]int{0, 1, 2}))

	s.Delete(1)
	Expect(t, s.Exists(1), BeFalse())
	Expect(t, s.Len(), Equal(2))
	Expect(t, s.Keys(), EquivalentSlice([]int{0, 2}))

	s.Clear()
	Expect(t, s.Exists(0), BeFalse())
	Expect(t, s.Exists(1), BeFalse())
	Expect(t, s.Exists(2), BeFalse())
	Expect(t, s.Len(), Equal(0))
}
