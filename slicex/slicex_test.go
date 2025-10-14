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
