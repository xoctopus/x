package iterx_test

import (
	"fmt"
	"iter"
	"reflect"
	"strconv"

	"github.com/xoctopus/x/iterx"
)

// pushes 1..5 to downstream until all numbers pushed or them don't need more
func Push() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 1; i <= 5; i++ {
			fmt.Printf("push-iter: yield(%d)\n", i)
			if !yield(i) {
				fmt.Println("push-iter: yield returned false, stop pushing")
				return
			}
		}
	}
}

func Example_iter_break() {
	// pulls values from upstream until meet 3
	pull := func() {
		for v := range Push() {
			fmt.Printf("pull-range: got %d\n", v)
			if v == 3 {
				fmt.Println("pull-range: break. no more is needed")
				break
			}
		}
	}
	pull()

	// Output:
	// push-iter: yield(1)
	// pull-range: got 1
	// push-iter: yield(2)
	// pull-range: got 2
	// push-iter: yield(3)
	// pull-range: got 3
	// pull-range: break. no more is needed
	// push-iter: yield returned false, stop pushing
}

func ExampleMap() {
	seq := iterx.Map(iterx.Slice([]int{0, 1, 2}), func(x int) int { return x + 2 })
	for v := range seq {
		fmt.Println(v)
		if v == 3 {
			break
		}
	}
	// Output:
	// 2
	// 3
}

func ExampleMapSlice() {
	seq := iterx.MapSlice([]int{0, 1, 2}, func(x int) int { return x + 1 })
	for v := range seq {
		fmt.Println(v)
		if v == 2 {
			break
		}
	}
	// Output:
	// 1
	// 2
}

func ExampleFilter() {
	seq := iterx.Filter(iterx.Slice([]int{0, 1, 2}), func(x int) bool { return x > 0 })
	for v := range seq {
		fmt.Println(v)
		if v == 1 {
			break
		}
	}
	// Output:
	// 1
}

func ExampleFilterSlice() {
	seq := iterx.FilterSlice([]int{0, 1, 2}, func(x int) bool { return x > 0 })
	for v := range seq {
		fmt.Println(v)
		if v == 1 {
			break
		}
	}
	// Output:
	// 1
}

func ExampleSlice() {
	seq := iterx.Slice([]int{0, 1, 2})
	for v := range seq {
		fmt.Println(v)
		if v == 1 {
			break
		}
	}
	// Output:
	// 0
	// 1
}

func ExampleSliceSeq() {
	seq := iterx.SliceSeq([]int{0, 1, 2})
	for k, v := range seq {
		fmt.Printf("%d:%d\n", k, v)
		if k == 0 {
			break
		}
	}
	// Output:
	// 0:0
}

func ExampleMapSeq() {
	seq := iterx.MapSeq(map[int]string{1: "a"})
	for k, v := range seq {
		fmt.Printf("%d:%s\n", k, v)
		if k == 1 {
			break
		}
	}
	// Output:
	// 1:a
}

func Example() {
	slices := []int{-1, 0, 1, 2, -1}

	seq := iterx.Map(
		// filter negative numbers
		iterx.FilterSlice(slices, func(x int) bool { return x >= 0 }),
		// add 1 and convert int to string
		func(x int) string { return strconv.Itoa(x + 1) },
	)

	for v := range seq {
		fmt.Println(reflect.TypeOf(v), v)
	}

	fmt.Println(iterx.Values(seq))

	// Output:
	// string 1
	// string 2
	// string 3
	// [1 2 3]
}
