package iterx_test

import (
	"iter"
	"testing"

	"github.com/xoctopus/x/iterx"
	"github.com/xoctopus/x/testx"
)

func TestFlatten(t *testing.T) {
	input := [][]int{
		{1, 2},
		{3, 4},
		{5},
	}

	seqs := iterx.Map(iterx.Of(input), func(s []int) iter.Seq[int] {
		return iterx.Of(s)
	})

	var out []int
	for v := range iterx.Flatten(seqs) {
		if v == 5 {
			break
		}
		out = append(out, v)
	}
	testx.Expect(t, out, testx.Equal([]int{1, 2, 3, 4}))
}
