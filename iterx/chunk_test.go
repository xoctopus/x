package iterx_test

import (
	"testing"
	"time"

	"github.com/xoctopus/x/iterx"
	"github.com/xoctopus/x/testx"
)

func TestChunkN(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7}

	seq := iterx.ChunkN(iterx.Of(input), 3)
	var (
		out [][]int
		idx int
	)
	for chunk := range seq {
		if idx == 2 {
			break
		}
		out = append(out, chunk)
		idx++
	}
	testx.Expect(t, out, testx.Equal([][]int{
		{1, 2, 3},
		{4, 5, 6},
	}))

	var out2 [][]int
	for chunk := range seq {
		out2 = append(out2, chunk)
		break
	}
	testx.Expect(t, out2, testx.Equal([][]int{{1, 2, 3}}))

	var (
		exact = []int{1, 2, 3, 4, 5, 6}
		out3  [][]int
	)
	seq = iterx.ChunkN(iterx.Of(exact), 3)
	for chunk := range seq {
		out3 = append(out3, chunk)
	}
	testx.Expect(t, out3, testx.Equal([][]int{
		{1, 2, 3},
		{4, 5, 6},
	}))

	testx.ExpectPanic[string](t,
		func() {
			iterx.ChunkN(iterx.Of(input), 0)
		},
		testx.Equal("cannot be less than 1"),
	)
}

func TestChunkDuration(t *testing.T) {
	// A mock generator that yields fast, then sleeps, then yields fast
	gen := func(yield func(int) bool) {
		yield(1)
		time.Sleep(30 * time.Millisecond) // trigger chunk
		yield(2)
		time.Sleep(30 * time.Millisecond) // trigger chunk
		yield(3)
		yield(3)
		time.Sleep(30 * time.Millisecond) // trigger chunk
		yield(4)
	}

	// Chunk after 15ms of inactivity
	seq := iterx.ChunkDuration(gen, 15*time.Millisecond)

	var (
		out [][]int
		idx = 0
	)
	for chunk := range seq {
		if idx >= 3 {
			break
		}
		out = append(out, chunk)
		idx++
	}

	testx.Expect(t, len(out), testx.Equal(3))
	testx.Expect(t, out[0], testx.Equal([]int{1}))
	testx.Expect(t, out[1], testx.Equal([]int{2}))
	testx.Expect(t, out[2], testx.Equal([]int{3, 3}))

	// Test early return
	var outEarly [][]int
	for chunk := range iterx.ChunkDuration(gen, 15*time.Millisecond) {
		outEarly = append(outEarly, chunk)
		break
	}
	testx.Expect(t, len(outEarly), testx.Equal(1))
}

func TestDebounce(t *testing.T) {
	// A mock generator that yields fast, then sleeps, then yields fast
	gen := func(yield func(int) bool) {
		// Burst 1
		yield(0)
		yield(0)
		yield(1)
		time.Sleep(15 * time.Millisecond) // Let debounce trigger
		// Burst 2. scene: keyboard multi hitting
		yield(0)
		yield(0)
		yield(2)
	}

	// Debounce window of 15ms
	seq := iterx.Debounce(gen, 15*time.Millisecond)

	var out []int
	for v := range seq {
		out = append(out, v)
	}

	// Should get the last item of burst 1 (3), and last item of burst 2 (6)
	testx.Expect(t, out, testx.Equal([]int{1, 2}))

	// Test early return
	var outEarly []int
	for v := range iterx.Debounce(gen, 15*time.Millisecond) {
		outEarly = append(outEarly, v)
		break
	}
	testx.Expect(t, outEarly, testx.Equal([]int{1}))
}
