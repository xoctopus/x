package slicex_test

import (
	"errors"
	"strings"
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
		key   string
		other int
	}

	ss := []X{
		{"a", 1},
		{"b", 2},
		{"c", 3},
		{"d", 4},
		{"a", 1},
		{"b", 2},
	}

	keys := slicex.UniqueKeys(ss, func(e X) string { return e.key })
	Expect(t, keys, EquivalentSlice([]string{"a", "b", "c", "d"}))

	values := slicex.UniqueValues(ss, func(e X) string { return e.key })
	Expect(t, len(values), Equal(4))

	others := slicex.UniqueKeys(values, func(e X) int { return e.other })
	Expect(t, others, EquivalentSlice([]int{1, 2, 3, 4}))
}

type MyError struct {
	idx string
}

func (e *MyError) Error() string { return e.idx }

func TestM(t *testing.T) {
	errs := []*MyError{{"1"}, {"2"}, {"3"}}

	mapped := slicex.M(errs, func(e *MyError) error { return e })
	Expect(t, len(errs), Equal(len(mapped)))

	joined := errors.Join(mapped...)
	Expect(t, joined.Error(), Equal(strings.Join([]string{"1", "2", "3"}, "\n")))
}
