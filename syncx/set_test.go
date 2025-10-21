package syncx_test

import (
	"testing"

	"github.com/xoctopus/x/syncx"
	. "github.com/xoctopus/x/testx"
)

func TestNewSet(t *testing.T) {
	s := syncx.NewSet[int]()

	Expect(t, s.Exists(1), BeFalse())

	s.Store(1)
	Expect(t, s.Exists(1), BeTrue())
	Expect(t, s.Keys(), Equal([]int{1}))
	Expect(t, s.Len(), Equal(1))
	s.Store(2)
	Expect(t, s.Exists(1), BeTrue())
	Expect(t, s.Len(), Equal(2))

	s.Delete(1)
	Expect(t, s.Exists(1), BeFalse())
	Expect(t, s.Len(), Equal(1))
	Expect(t, s.Keys(), Equal([]int{2}))

	s.Clear()
	Expect(t, s.Exists(1), BeFalse())
	Expect(t, s.Exists(2), BeFalse())
	Expect(t, s.Len(), Equal(0))
	Expect(t, s.Keys(), Equal([]int{}))
}
