package must_test

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/sincospro/x/misc/must"
)

func ExampleNoError() {
	must.NoError(nil)

	defer func() {
		fmt.Println(recover())
	}()
	must.NoError(errors.New("any"))

	// Output:
	// any
}

func ExampleNoErrorV() {
	fmt.Println(must.NoErrorV(100, nil))

	defer func() {
		fmt.Println(recover())
	}()
	must.NoError(errors.New("any"))

	// Output:
	// 100
	// any
}

func ExampleBeTrue() {
	must.BeTrue(true)

	defer func() {
		fmt.Println(recover())
	}()
	must.BeTrue(false)

	// Output:
	// must ok
}

func ExampleBeTrueV() {
	fmt.Println(must.BeTrueV(float32(100.1), true))

	defer func() {
		fmt.Println(recover())
	}()
	_ = must.BeTrueV(new(float64), false)

	// Output:
	// 100.1
	// must ok
}
