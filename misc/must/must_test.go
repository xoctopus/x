package must_test

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/xoctopus/x/misc/must"
)

func ReturnError() error {
	return errors.New("some error")
}

func ReturnNoError() error {
	return nil
}

func ReturnIntError() (int, error) {
	return 100, errors.New("some error")
}

func ReturnIntNoError() (int, error) {
	return 100, nil
}

func ReturnTrue() bool {
	return true
}

func ReturnFalse() bool {
	return false
}

func ReturnIntTrue() (int, bool) {
	return 100, true
}

func ReturnIntFalse() (int, bool) {
	return 100, false
}

func ExampleNoError() {
	must.NoError(ReturnNoError())

	defer func() {
		fmt.Println(recover())
	}()
	must.NoError(ReturnError())

	// Output:
	// some error
}

func ExampleNoErrorV() {
	fmt.Println(must.NoErrorV(ReturnIntNoError()))

	defer func() {
		fmt.Println(recover())
	}()
	must.NoErrorV(ReturnIntError())

	// Output:
	// 100
	// some error
}

func ExampleNoErrorWrap() {
	must.NoErrorWrap(ReturnNoError(), "any")

	defer func() {
		fmt.Println(recover())
	}()
	must.NoErrorWrap(ReturnError(), "some message: %d", 10)

	// Output:
	// some message: 10: some error
}

func ExampleBeTrue() {
	must.BeTrue(ReturnTrue())

	defer func() {
		fmt.Println(recover())
	}()
	must.BeTrue(ReturnFalse())

	// Output:
	// must be true
}

func ExampleBeTrueV() {
	fmt.Println(must.BeTrueV(ReturnIntTrue()))

	defer func() {
		fmt.Println(recover())
	}()
	_ = must.BeTrueV(ReturnIntFalse())

	// Output:
	// 100
	// must be true
}

func ExampleBeTrueWrap() {
	must.BeTrueWrap(ReturnTrue(), "any")

	defer func() {
		fmt.Println(recover())
	}()
	must.BeTrueWrap(ReturnFalse(), "required exists")

	// Output:
	// must be true: required exists
}
