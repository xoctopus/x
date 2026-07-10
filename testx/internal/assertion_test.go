package internal_test

import (
	"fmt"

	"github.com/xoctopus/x/testx"
	"github.com/xoctopus/x/testx/internal"
	"github.com/xoctopus/x/testx/testutil"
)

type uncomparable struct {
	x []int
}

func ExampleExpect() {
	t := &testutil.MockTB{}

	t.Reset()
	internal.Expect(t, uncomparable{x: []int{1}}, testx.Equal(uncomparable{x: []int{2}}))
	fmt.Println(t.Output2())

	t.Reset()
	internal.Expect(t, uncomparable{x: []int{1}}, testx.NotEqual(uncomparable{x: []int{1}}))
	fmt.Println(t.Output2())

	t.Reset()
	internal.Expect(t, 1.2, testx.BeAssignableTo[int]())
	fmt.Println(t.Output2())

	// Output:
	// should Equal, but got
	//   strings.Join({
	//   	"internal_test.uncomparable:{x:[",
	// - 	"2",
	// + 	"1",
	//   	"]}",
	//   }, "")
	//
	// should not Equal, but got
	// {[1]}
	//
	// should BeAssignableTo[int], but got
	// 1.2
}
