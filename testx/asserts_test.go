package testx_test

import (
	"fmt"

	. "github.com/xoctopus/x/testx"
	"github.com/xoctopus/x/testx/testutil"
)

type MockCodeErr int

func (MockCodeErr) Message() string { return "mock" }

func crash(i int) {
	switch i {
	case 0:
		return
	case 1:
		panic(fmt.Errorf("any"))
	}
}

func ExampleExpectPanic() {
	t := &testutil.MockTB{}

	t.Reset()
	ExpectPanic[error](t, func() {
		panic("any")
	})
	fmt.Println(t.Output2())

	t.Reset()
	ExpectPanic[string](t, func() {
		crash(1)
	})
	fmt.Println(t.Output2())

	t.Reset()
	ExpectPanic[string](t, func() {})
	fmt.Println(t.Output2())

	// Output:
	// expect a panic of `error`, but got string
	//
	// expect a panic of `string`, but got *errors.errorString
	//
	// expect a panic of `string`, but f returned normally
}
