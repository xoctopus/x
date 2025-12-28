package slicex_test

import (
	"testing"

	"github.com/xoctopus/x/slicex"
	. "github.com/xoctopus/x/testx"
)

func TestUnique(t *testing.T) {
	ints := []int{1}
	Expect(t, slicex.Unique(ints), Equal([]int{1}))

	ints = []int{1, 1, 4, 2, 3, 2, 3, 4}
	Expect(t, slicex.Unique(ints), Equal([]int{1, 4, 2, 3}))
}

func TestEqualElements(t *testing.T) {
	Expect(t, slicex.Equivalent([]int{1}, []int{1}), BeTrue())
	Expect(t, slicex.Equivalent([]int{}, nil), BeFalse())
	Expect(t, slicex.Equivalent(nil, []int{}), BeFalse())
	Expect(t, slicex.Equivalent([]int{1, 2}, []int{2, 1}), BeTrue())
	Expect(t, slicex.Equivalent([]int{1, 1, 2}, []int{1, 2, 1}), BeTrue())
	Expect(t, slicex.Equivalent([]int{1, 2}, []int{1, 2, 1}), BeFalse())
	Expect(t, slicex.Equivalent([]int{2, 1, 2}, []int{1, 2, 1}), BeFalse())
}

func TestUniqueMapping(t *testing.T) {
	type X struct {
		val string
	}

	ss := []X{
		{"a"}, {"b"}, {"c"}, {"d"}, {"a"}, {"b"},
	}

	vals := slicex.UniqueM(ss, func(e X) string { return e.val })
	Expect(t, vals, EquivalentSlice([]string{"a", "b", "c", "d"}))
}
