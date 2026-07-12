package iterx_test

import (
	"slices"
	"testing"

	"github.com/xoctopus/x/iterx"
	"github.com/xoctopus/x/testx"
)

func TestConcat(t *testing.T) {
	var out []int
	for v := range iterx.Concat(
		iterx.Of([]int{1, 2}),
		iterx.Of([]int{3, 4}),
		iterx.Of([]int{5}),
	) {
		if v == 5 {
			break
		}
		out = append(out, v)
	}
	testx.Expect(t, out, testx.Equal([]int{1, 2, 3, 4}))

	out = out[0:0]
	for v := range iterx.Concat[int]() {
		out = append(out, v)
	}
	testx.Expect(t, len(out), testx.Equal(0))

	out = out[0:0]
	for v := range iterx.Concat(iterx.Of([]int{1, 2})) {
		out = append(out, v)
	}
	testx.Expect(t, out, testx.Equal([]int{1, 2}))
}

func TestMerge(t *testing.T) {
	seq1 := iterx.Of([]int{1, 2, 3})
	seq2 := iterx.Of([]int{4, 5, 6})

	var out []int
	for v := range iterx.Merge(seq1, seq2) {
		out = append(out, v)
	}
	slices.Sort(out) // Merge order is not guaranteed
	testx.Expect(t, out, testx.Equal([]int{1, 2, 3, 4, 5, 6}))

	out = out[0:0]
	for v := range iterx.Merge(seq1, seq2) {
		out = append(out, v)
		if len(out) == 2 {
			break
		}
	}
	testx.Expect(t, len(out), testx.Equal(2))

	out = out[0:0]
	for v := range iterx.Merge[int]() {
		out = append(out, v)
	}
	testx.Expect(t, len(out), testx.Equal(0))

	out = out[0:0]
	for v := range iterx.Merge(seq1) {
		out = append(out, v)
	}
	testx.Expect(t, out, testx.Equal([]int{1, 2, 3}))
}
