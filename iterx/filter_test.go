package iterx_test

import (
	"strconv"
	"testing"

	"github.com/xoctopus/x/iterx"
	"github.com/xoctopus/x/testx"
)

func odd(v int) bool {
	return v%2 != 0
}

func even(v int) bool {
	return !odd(v)
}

func TestFilter(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}

	seq := iterx.Filter(iterx.Of(input), even)

	var out []int
	for v := range seq {
		if v > 5 {
			break
		}
		out = append(out, v)
	}

	testx.Expect(t, out, testx.Equal([]int{2, 4}))
}

func TestMap(t *testing.T) {
	input := []int{1, 2, 3}

	seq := iterx.Map(iterx.Of(input), func(v int) string {
		return strconv.Itoa(v)
	})

	var out []string
	for v := range seq {
		if v == "3" {
			break
		}
		out = append(out, v)
	}
	testx.Expect(t, out, testx.Equal([]string{"1", "2"}))
}

func TestMapFilter(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}

	// keep even numbers, multiply by 10 and convert to string
	seq := iterx.MapFilter(iterx.Of(input), func(v int) (string, bool) {
		if v%2 != 0 {
			return "", false
		}
		return strconv.Itoa(v), true
	})

	var out []string
	for v := range seq {
		if v == "6" {
			break
		}
		out = append(out, v)
	}
	testx.Expect(t, out, testx.Equal([]string{"2", "4"}))
}
